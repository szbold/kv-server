package datastore

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

func (ds *DataStore) set(key, value string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.data[key] = newEntry(value, nil)
}

func (ds *DataStore) get(key string) string {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, exists := ds.data[key]; exists {
		return e.value
	}

	return "nil"
}

// a bit hacky because error does not implement stringer for some reaseon ????
func (ds *DataStore) incr(key string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if _, ok := ds.data[key]; ok {
		if ds.data[key].dtype == "int" {
			// error cant occur since data is of int type
			tmp, _ := strconv.Atoi(ds.data[key].value)
			newVal := strconv.Itoa(tmp + 1)
			ds.data[key] = newEntry(newVal, nil)

			return nil
		}
		return errors.New(fmt.Sprintf("Cannot run incr on value: %v with type: %v", ds.data[key].value, ds.data[key].dtype))
	}

	return errors.New(fmt.Sprintf("Key \"%v\" does not exist", key))
}

func (ds *DataStore) exists(key string) bool {
	exists := false

	ds.mu.Lock()
	defer ds.mu.Unlock()

	if _, ok := ds.data[key]; ok {
		exists = true
	}

	return exists
}

func (ds *DataStore) del(key string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if _, ok := ds.data[key]; ok {
		delete(ds.data, key)
	} else {
		return errors.New(fmt.Sprintf("Key \"%v\" does not exist", key))
	}

	return nil
}

func (ds *DataStore) dtype(key string) (string, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if _, ok := ds.data[key]; ok {
		return ds.data[key].dtype, nil
	}

	return "", errors.New(fmt.Sprintf("Key \"%v\" does not exist", key))
}

func (ds *DataStore) expire(key string, ttlStr string) error {
	ttl, err := strconv.Atoi(ttlStr)

	if err != nil {
		return errors.New(fmt.Sprint("TTL should be int found string"))
	}

	if ttl < 1 {
		return errors.New(fmt.Sprint("TTL should be at least 1"))
	}

	ds.mu.Lock()

	if e, ok := ds.data[key]; ok {
    ch := make(chan int, 1)
    e.ttlChan = ch
		ds.data[key] = e
		ds.mu.Unlock()

		go ds.emitTtl(&e, key, ttl)

		return nil
	}

	ds.mu.Unlock()
	return errors.New(fmt.Sprintf("Key \"%v\" does not exist", key))
}

func (ds *DataStore) emitTtl(e *entry, key string, ttl int) {
  defer close(e.ttlChan)
	e.ttlChan <- ttl

	for {
		go func() {
			if len(e.ttlChan) == 1 {
				<-e.ttlChan
			}
		}()

		e.ttlChan <- ttl

		time.Sleep(time.Second)
		ttl--

		if ttl == 0 {
			_ = ds.del(key)
			return
		}
	}
}

func (ds *DataStore) setexp(key, value string, ttlStr string) error {
	ds.set(key, value)
	return ds.expire(key, ttlStr)
}

func (ds *DataStore) ttl(key string) (string, error) {
	if _, ok := ds.data[key]; ok {
		return strconv.Itoa(<-ds.data[key].ttlChan), nil
	}

	return "", errors.New(fmt.Sprintf("Key \"%v\" does not exist", key))
}
