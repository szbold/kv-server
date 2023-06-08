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

	ds.data[key] = newEntry(value)
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
			ds.data[key] = newEntry(newVal)

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

	if entry, ok := ds.data[key]; ok {
		ds.data[key] = newEntryExp(entry.value, uint(ttl))
		ds.mu.Unlock()
		go func() {
			liveTtl := ttl
			for {
				time.Sleep(time.Second)
				liveTtl--

				if liveTtl == 0 {
					_ = ds.del(key)
					return
				}
			}
		}()
		return nil
	}

	ds.mu.Unlock()
	return errors.New(fmt.Sprintf("Key \"%v\" does not exist", key))
}
