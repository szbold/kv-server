package datastore

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (ds *DataStore) set(key, value string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	var t dtype
	_, err := strconv.Atoi(value)

	if err != nil {
		t = t_string
	} else {
		t = t_int
	}

	ds.data[key] = newEntry(value, t)
}

func (ds *DataStore) get(key string) (string, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, exists := ds.data[key]; exists {
		return e.value, nil
	}

	return "", errors.New(fmt.Sprintf("Key \"%v\" does not exist", key))
}

func (ds *DataStore) incr(key string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if _, ok := ds.data[key]; ok {
		if ds.data[key].dataType == t_int {
			// error cant occur since data is of int type
			tmp, _ := strconv.Atoi(ds.data[key].value)
			newVal := strconv.Itoa(tmp + 1)
			ds.data[key] = newEntry(newVal, t_int)

			return nil
		}
		return errors.New(fmt.Sprintf("Cannot run incr on value: %v with type: %v", ds.data[key].value, ds.data[key].dataType))
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
	var strType string
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, ok := ds.data[key]; ok {
		switch e.dataType {
		case t_string:
			strType = "string"
		case t_int:
			strType = "int"
		case t_list:
			strType = "list"
		}
		return strType, nil
	}

	return strType, errors.New(fmt.Sprintf("Key \"%v\" does not exist", key))
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

		go ds.emitTtl(key, ttl)

		return nil
	}

	ds.mu.Unlock()
	return errors.New(fmt.Sprintf("Key \"%v\" does not exist", key))
}

func (ds *DataStore) emitTtl(key string, ttl int) {
	e, _ := ds.data[key]
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

func (ds *DataStore) setexp(key, value, ttlStr string) error {
	ds.set(key, value)
	return ds.expire(key, ttlStr)
}

func (ds *DataStore) ttl(key string) (string, error) {
	if _, ok := ds.data[key]; ok {
		return strconv.Itoa(<-ds.data[key].ttlChan), nil
	}

	return "", errors.New(fmt.Sprintf("Key \"%v\" does not exist", key))
}

func (ds *DataStore) lpush(key string, values []string) {
	strVal := strings.Join(values, ",")

	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, ok := ds.data[key]; ok {
		e.value = strVal + "," + e.value
		ds.data[key] = e
	} else {
		ds.data[key] = newEntry(strVal, t_list)
	}
}

func (ds *DataStore) rpush(key string, values []string) {
	strVal := strings.Join(values, ",")

	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, ok := ds.data[key]; ok {
		e.value = e.value + "," + strVal
		ds.data[key] = e
	} else {
		ds.data[key] = newEntry(strVal, t_list)
	}
}

func (ds *DataStore) llen(key string) (string, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, exists := ds.data[key]; exists {
		if e.dataType == t_list {
			res := strings.Count(e.value, ",") + 1
			return strconv.Itoa(res), nil
		}

		return "", errors.New(fmt.Sprintf("Value of type %v does not have property length", e.dataType.String()))
	}

	return "", errors.New(fmt.Sprintf("Key \"%v\" does not exist", key))
}

func (ds *DataStore) lrange(key, startStr, endStr string) (string, error) {
	start, err := strconv.Atoi(startStr)

	if err != nil {
		return "", errors.New("Start should be a number value")
	}

	end, err := strconv.Atoi(endStr)

	if err != nil {
		return "", errors.New("End should be a number value")
	}

	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, exists := ds.data[key]; exists {
		if e.dataType == t_list {
			list := strings.Split(e.value, ",")
			if start > end {
				start, end = end, start
			}

			if start < 0 {
				start = 0
			}

			if end > len(list) - 1 {
				end = len(list) - 1
			}
			return strings.Join(list[start:end + 1], ","), nil
		}

		return "", errors.New(fmt.Sprintf("Value of type %v does not have property length", e.dataType.String()))
	}

	return "", errors.New(fmt.Sprintf("Key \"%v\" does not exist", key))
}
