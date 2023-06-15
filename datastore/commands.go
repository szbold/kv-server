package datastore

import (
	"fmt"
	"key-value-server/consts"
	types "key-value-server/datatypes"
	"strconv"
	"time"
)

func (ds *DataStore) keys() types.KvList {
	var result types.KvList
	ds.mu.Lock()
	defer ds.mu.Unlock()

	for key := range ds.data {
		result = append(result, key)
	}

	return result
}

func (ds *DataStore) set(key string, value types.Data) types.KvString {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.data[key] = newEntry(value)

	return types.KvString(consts.Ok)
}

func (ds *DataStore) get(key string) types.Data {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, exists := ds.data[key]; exists {
		return e.value
	}
	return types.NewKvError(fmt.Sprintf("Key '%v' does not exist", key))
}

func (ds *DataStore) incr(key string) types.Data {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, ok := ds.data[key]; ok {
		if e.value.Type() == types.TInt {
			val := e.value.(types.KvInt)
			ds.data[key] = newEntry(val + types.KvInt(1))
			return types.KvString(consts.Ok)
		}

		return types.NewKvError(fmt.Sprintf("Cannot run incr on value: %v with type: %v", e.value, e.value.Type()))
	}

	return types.NewKvError(fmt.Sprintf("Key '%v' does not exist", key))
}

func (ds *DataStore) exists(key string) types.KvInt {
	var exists types.KvInt = 0

	ds.mu.Lock()
	defer ds.mu.Unlock()

	if _, ok := ds.data[key]; ok {
		exists = 1
	}

	return exists
}

func (ds *DataStore) del(key string) types.Data {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if _, ok := ds.data[key]; ok {
		delete(ds.data, key)
	} else {
		return types.NewKvError(fmt.Sprintf("Key '%v' does not exist", key))
	}

	return types.KvString(consts.Ok)
}

func (ds *DataStore) dtype(key string) types.Data {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, ok := ds.data[key]; ok {
		return types.KvString(e.value.Type())
	}

	return types.NewKvError(fmt.Sprintf("Key '%v' does not exist", key))
}

func (ds *DataStore) expire(key string, ttlStr string) types.Data {
	ttl, err := strconv.Atoi(ttlStr)

	if err != nil {
		return types.NewKvError(fmt.Sprint("TTL should be int found string"))
	}

	if ttl < 1 {
		return types.NewKvError(fmt.Sprint("TTL should be at least 1"))
	}

	ds.mu.Lock()

	if e, ok := ds.data[key]; ok {
		ch := make(chan types.KvInt, 1)
		e.ttlChan = ch
		ds.data[key] = e
		ds.mu.Unlock()

		go ds.emitTtl(key, ttl)

		return types.KvString(consts.Ok)
	}

	ds.mu.Unlock()
	return types.NewKvError(fmt.Sprintf("Key '%v' does not exist", key))
}

func (ds *DataStore) emitTtl(key string, ttl int) {
	e, _ := ds.data[key]
	defer close(e.ttlChan)
	e.ttlChan <- types.KvInt(ttl)

	for {
		go func() {
			if len(e.ttlChan) == 1 {
				<-e.ttlChan
			}
		}()

		e.ttlChan <- types.KvInt(ttl)

		time.Sleep(time.Second)
		ttl--

		if ttl == 0 {
			_ = ds.del(key)
			return
		}
	}
}

func (ds *DataStore) setexp(key string, value types.Data, ttlStr string) types.Data {
	ds.set(key, value)
	return ds.expire(key, ttlStr)
}

func (ds *DataStore) ttl(key string) types.Data {
	if _, ok := ds.data[key]; ok {
		return <-ds.data[key].ttlChan
	}

	return types.NewKvError(fmt.Sprintf("Key '%v' does not exist", key))
}

func (ds *DataStore) lpush(key string, values []string) types.Data {
	var list types.KvList

	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, ok := ds.data[key]; ok {
		if e.value.Type() == types.TList {
			list = e.value.(types.KvList)
			list = append(values, list...)
			e.value = list
			ds.data[key] = e
		} else {
			return types.NewKvError(fmt.Sprintf("Cannot lpush on type %v", e.value.Type()))
		}
	} else {
		ds.data[key] = newEntry(types.KvList(values))
	}

	return types.KvString(consts.Ok)
}

func (ds *DataStore) rpush(key string, values []string) types.Data {
	var list types.KvList

	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, ok := ds.data[key]; ok {
		if e.value.Type() == types.TList {
			list = e.value.(types.KvList)
			list = append(list, values...)
			e.value = list
			ds.data[key] = e
		} else {
			return types.NewKvError(fmt.Sprintf("Cannot rpush on type %v", e.value.Type()))
		}
	} else {
		ds.data[key] = newEntry(list)
	}

	return types.KvString(consts.Ok)
}

func (ds *DataStore) llen(key string) types.Data {
	var list types.KvList
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, exists := ds.data[key]; exists {
		if e.value.Type() == types.TList {
			list = e.value.(types.KvList)
			return types.KvInt(len(list))
		}

		return types.NewKvError(fmt.Sprintf("Value of type %v does not have property length", e.value.Type()))
	}

	return types.NewKvError(fmt.Sprintf("Key '%v' does not exist", key))
}

func (ds *DataStore) lrange(key, startStr, endStr string) types.Data {
	var list types.KvList
	start, err := strconv.Atoi(startStr)

	if err != nil {
		return types.NewKvError("Start should be a number value")
	}

	end, err := strconv.Atoi(endStr)

	if err != nil {
		return types.NewKvError("End should be a number value")
	}

	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, exists := ds.data[key]; exists {
		if e.value.Type() == types.TList {
			list = e.value.(types.KvList)
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

		return types.NewKvError(fmt.Sprintf("Cannot use lrange on %v", e.value.Type()))
	}

	return types.NewKvError(fmt.Sprintf("Key '%v' does not exist", key))
}

func (ds *DataStore) ltrim(key, startStr, endStr string) types.Data {
	var list types.KvList
	start, err := strconv.Atoi(startStr)

	if err != nil {
		return types.NewKvError("Start should be a number value")
	}

	end, err := strconv.Atoi(endStr)

	if err != nil {
		return types.NewKvError("End should be a number value")
	}

	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, exists := ds.data[key]; exists {
		if e.value.Type() == types.TList {
			list = e.value.(types.KvList)
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

			return types.KvString(consts.Ok)
		}

		return types.NewKvError(fmt.Sprintf("Cannot use ltrim on %v", e.value.Type()))
	}

	return types.NewKvError(fmt.Sprintf("Key '%v' does not exist", key))
}
