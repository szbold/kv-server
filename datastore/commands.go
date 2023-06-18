package datastore

import (
	"fmt"
	"key-value-server/consts"
	. "key-value-server/datatypes"
	"strconv"
	"time"
)

func (ds *DataStore) keys() KvList {
	var result KvList
	ds.mu.Lock()
	defer ds.mu.Unlock()

	for key := range ds.data {
		result = append(result, key)
	}

	return result
}

func (ds *DataStore) set(key string, value Data) KvString {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.data[key] = newEntry(value)

	return KvString(consts.Ok)
}

func (ds *DataStore) get(key string) Data {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, ok := ds.data[key]; ok {
		return e.value
	}
	return NewKvError(fmt.Sprintf("Key '%v' does not exist", key))
}

func (ds *DataStore) incr(key string) Data {
	return ds.incrby(key, "1")
}

func (ds *DataStore) incrby(key, incrementStr string) Data {
	increment, err := strconv.Atoi(incrementStr)

	if err != nil {
		return NewKvError(fmt.Sprint("Increment should be int found string"))
	}

	ds.mu.Lock()
	defer ds.mu.Unlock()

	e, ok := ds.data[key]

	if !ok {
		return NewKvError(fmt.Sprintf("Key '%v' does not exist", key))
	}

	if e.value.Type() != TInt {
		return NewKvError(fmt.Sprintf("Cannot run incrby on value: %v with type: %v", e.value, e.value.Type()))
	}

	val := e.value.(KvInt)
	ds.data[key] = newEntry(val + KvInt(increment))
	return KvString(consts.Ok)
}

func (ds *DataStore) decr(key string) Data {
	return ds.decrby(key, "1")
}

func (ds *DataStore) decrby(key, decrementStr string) Data {
	decrement, err := strconv.Atoi(decrementStr)

	if err != nil {
		return NewKvError(fmt.Sprint("Decrement should be int found string"))
	}

	ds.mu.Lock()
	defer ds.mu.Unlock()

	e, ok := ds.data[key]

	if !ok {
		return NewKvError(fmt.Sprintf("Key '%v' does not exist", key))
	}

	if e.value.Type() != TInt {
		return NewKvError(fmt.Sprintf("Cannot run decrby on value: %v with type: %v", e.value, e.value.Type()))
	}

	val := e.value.(KvInt)
	ds.data[key] = newEntry(val - KvInt(decrement))
	return KvString(consts.Ok)
}

func (ds *DataStore) exists(key string) KvInt {
	var exists KvInt = 0

	ds.mu.Lock()
	defer ds.mu.Unlock()

	if _, ok := ds.data[key]; ok {
		exists = 1
	}

	return exists
}

func (ds *DataStore) del(key string) Data {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if _, ok := ds.data[key]; ok {
		delete(ds.data, key)
	} else {
		return NewKvError(fmt.Sprintf("Key '%v' does not exist", key))
	}

	return KvString(consts.Ok)
}

func (ds *DataStore) dtype(key string) Data {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, ok := ds.data[key]; ok {
		return KvString(e.value.Type())
	}

	return NewKvError(fmt.Sprintf("Key '%v' does not exist", key))
}

func (ds *DataStore) expire(key string, ttlStr string) Data {
	ttl, err := strconv.Atoi(ttlStr)

	if err != nil {
		return NewKvError(fmt.Sprint("TTL should be int found string"))
	}

	if ttl < 1 {
		return NewKvError(fmt.Sprint("TTL should be at least 1"))
	}

	ds.mu.Lock()

	if e, ok := ds.data[key]; ok {
		ch := make(chan KvInt, 1)
		e.ttlChan = ch
		ds.data[key] = e
		ds.mu.Unlock()

		go ds.emitTtl(key, ttl)

		return KvString(consts.Ok)
	}

	ds.mu.Unlock()
	return NewKvError(fmt.Sprintf("Key '%v' does not exist", key))
}

func (ds *DataStore) emitTtl(key string, ttl int) {
	e, _ := ds.data[key]
	defer close(e.ttlChan)
	e.ttlChan <- KvInt(ttl)

	for {
		go func() {
			if len(e.ttlChan) == 1 {
				<-e.ttlChan
			}
		}()

		e.ttlChan <- KvInt(ttl)

		time.Sleep(time.Second)
		ttl--

		if ttl == 0 {
			_ = ds.del(key)
			return
		}
	}
}

func (ds *DataStore) setexp(key string, value Data, ttlStr string) Data {
	ds.set(key, value)
	return ds.expire(key, ttlStr)
}

func (ds *DataStore) ttl(key string) Data {
	if _, ok := ds.data[key]; ok {
		return <-ds.data[key].ttlChan
	}

	return NewKvError(fmt.Sprintf("Key '%v' does not exist", key))
}

func (ds *DataStore) lpush(key string, values []string) Data {
	var list KvList

	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, ok := ds.data[key]; ok {
		if e.value.Type() == TList {
			list = e.value.(KvList)
			list = append(values, list...)
			e.value = list
			ds.data[key] = e
		} else {
			return NewKvError(fmt.Sprintf("Cannot lpush on type %v", e.value.Type()))
		}
	} else {
		ds.data[key] = newEntry(KvList(values))
	}

	return KvString(consts.Ok)
}

func (ds *DataStore) rpush(key string, values []string) Data {
	var list KvList

	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, ok := ds.data[key]; ok {
		if e.value.Type() == TList {
			list = e.value.(KvList)
			list = append(list, values...)
			e.value = list
			ds.data[key] = e
		} else {
			return NewKvError(fmt.Sprintf("Cannot rpush on type %v", e.value.Type()))
		}
	} else {
		ds.data[key] = newEntry(KvList(values))
	}

	return KvString(consts.Ok)
}

func (ds *DataStore) llen(key string) Data {
	var list KvList
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, exists := ds.data[key]; exists {
		if e.value.Type() == TList {
			list = e.value.(KvList)
			return KvInt(len(list))
		}

		return NewKvError(fmt.Sprintf("Value of type %v does not have property length", e.value.Type()))
	}

	return NewKvError(fmt.Sprintf("Key '%v' does not exist", key))
}

func (ds *DataStore) lrange(key, startStr, endStr string) Data {
	var list KvList
	start, err := strconv.Atoi(startStr)

	if err != nil {
		return NewKvError("Start should be a number value")
	}

	end, err := strconv.Atoi(endStr)

	if err != nil {
		return NewKvError("End should be a number value")
	}

	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, exists := ds.data[key]; exists {
		if e.value.Type() == TList {
			list = e.value.(KvList)
			if start > end {
				start, end = end, start
			}

			if start < 0 {
				start = 0
			}

			if end > len(list)-1 {
				end = len(list) - 1
			}

			return list[start : end+1]
		}

		return NewKvError(fmt.Sprintf("Cannot use lrange on %v", e.value.Type()))
	}

	return NewKvError(fmt.Sprintf("Key '%v' does not exist", key))
}

func (ds *DataStore) ltrim(key, startStr, endStr string) Data {
	var list KvList
	start, err := strconv.Atoi(startStr)

	if err != nil {
		return NewKvError("Start should be a number value")
	}

	end, err := strconv.Atoi(endStr)

	if err != nil {
		return NewKvError("End should be a number value")
	}

	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, exists := ds.data[key]; exists {
		if e.value.Type() == TList {
			list = e.value.(KvList)
			if start > end {
				start, end = end, start
			}

			if start < 0 {
				start = 0
			}

			if end > len(list)-1 {
				end = len(list) - 1
			}

			e.value = list[start : end+1]
			ds.data[key] = e

			return KvString(consts.Ok)
		}

		return NewKvError(fmt.Sprintf("Cannot use ltrim on %v", e.value.Type()))
	}

	return NewKvError(fmt.Sprintf("Key '%v' does not exist", key))
}

func (ds *DataStore) sadd(key, value string) Data {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	e, exists := ds.data[key]

	if !exists {
		e = newEntry(NewKvSet())
	}

	if e.value.Type() != TSet {
		return NewKvError(fmt.Sprintf("Cannot use sadd on %v", e.value.Type()))
	}

	set := e.value.(KvSet)
	set.Insert(value)
	ds.data[key] = e
	return KvString(consts.Ok)
}

func (ds *DataStore) srem(key, value string) Data {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	e, exists := ds.data[key]

	if !exists {
		return NewKvError(fmt.Sprintf("Key '%v' does not exist", key))
	}

	if e.value.Type() != TSet {
		return NewKvError(fmt.Sprintf("Cannot use srem on %v", e.value.Type()))
	}

	set := e.value.(KvSet)
	set.Delete(value)
	e.value = set
	ds.data[key] = e
	return KvString(consts.Ok)
}

func (ds *DataStore) sismember(key, value string) Data {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	e, exists := ds.data[key]

	if !exists {
		return NewKvError(fmt.Sprintf("Key '%v' does not exist", key))
	}

	if e.value.Type() != TSet {
		return NewKvError(fmt.Sprintf("Cannot use sismember on %v", e.value.Type()))
	}

	set := e.value.(KvSet)

	if set.Has(value) {
		return KvInt(1)
	}
	return KvInt(0)
}

func (ds *DataStore) sinter(key, other string) Data {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	keyEntry, keyExists := ds.data[key]
	otherEntry, otherExists := ds.data[other]

	if !(keyExists && otherExists) {
		return NewKvError(fmt.Sprintf("One of the keys '%v' or '%v' does not exist", key, other))
	}

	if !(keyEntry.value.Type() == TSet && otherEntry.value.Type() == TSet) {
		return NewKvError(fmt.Sprintf("Canno perform sinter on '%v' or '%v'", key, other))
	}

	keySet := keyEntry.value.(KvSet)
	otherSet := otherEntry.value.(KvSet)

	return KvList(keySet.Intersection(otherSet))
}

func (ds *DataStore) scard(key string) Data {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	e, exists := ds.data[key]

	if !exists {
		return NewKvError(fmt.Sprintf("Key '%v' does not exist", key))
	}

	if e.value.Type() != TSet {
		return NewKvError(fmt.Sprintf("Cannot use scard on %v", e.value.Type()))
	}

	set := e.value.(KvSet)
	return KvInt(len(set))
}
