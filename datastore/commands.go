package datastore

import (
	"fmt"
	"kv-server/consts"
	. "kv-server/types"
	"strconv"
	"time"
)

func (ds *DataStore) keys() List {
	var result List
	ds.mu.Lock()
	defer ds.mu.Unlock()

	for key := range ds.data {
		result = append(result, key)
	}

	return result
}

func (ds *DataStore) set(key string, value Data) String {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.data[key] = newEntry(value)

	return String(consts.Ok)
}

func (ds *DataStore) get(key string) Data {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, ok := ds.data[key]; ok {
		return e.value
	}
	return NewError(fmt.Sprintf("Key '%v' does not exist", key))
}

func (ds *DataStore) incr(key string) Data {
	return ds.incrby(key, "1")
}

func (ds *DataStore) incrby(key, incrementStr string) Data {
	increment, err := strconv.Atoi(incrementStr)

	if err != nil {
		return ParseError("Increment", "int")
	}

	ds.mu.Lock()
	defer ds.mu.Unlock()

	e, ok := ds.data[key]

	if !ok {
		return NewError(fmt.Sprintf("Key '%v' does not exist", key))
	}

	if e.value.Type() != TNumber {
		return NewError(fmt.Sprintf("Cannot run incrby on: %v", e.value.Type()))
	}

	val := e.value.(Number)
	ds.data[key] = newEntry(val + Number(increment))
	return String(consts.Ok)
}

func (ds *DataStore) decr(key string) Data {
	return ds.decrby(key, "1")
}

func (ds *DataStore) decrby(key, decrementStr string) Data {
	decrement, err := strconv.Atoi(decrementStr)

	if err != nil {
		return NewError(fmt.Sprint("Decrement should be int found string"))
	}

	ds.mu.Lock()
	defer ds.mu.Unlock()

	e, ok := ds.data[key]

	if !ok {
		return NewError(fmt.Sprintf("Key '%v' does not exist", key))
	}

	if e.value.Type() != TNumber {
		return NewError(fmt.Sprintf("Cannot run decrby on: %v", e.value.Type()))
	}

	val := e.value.(Number)
	ds.data[key] = newEntry(val - Number(decrement))
	return String(consts.Ok)
}

func (ds *DataStore) exists(key string) Number {
	var exists Number = 0

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
		return NewError(fmt.Sprintf("Key '%v' does not exist", key))
	}

	return String(consts.Ok)
}

func (ds *DataStore) dtype(key string) Data {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, ok := ds.data[key]; ok {
		return String(e.value.Type())
	}

	return NewError(fmt.Sprintf("Key '%v' does not exist", key))
}

func (ds *DataStore) expire(key string, ttlStr string) Data {
	ttl, err := strconv.Atoi(ttlStr)

	if err != nil {
		return NewError(fmt.Sprint("TTL should be int found string"))
	}

	if ttl < 1 {
		return NewError(fmt.Sprint("TTL should be at least 1"))
	}

	ds.mu.Lock()

	if e, ok := ds.data[key]; ok {
		ch := make(chan Number, 1)
		e.ttlChan = ch
		ds.data[key] = e
		ds.mu.Unlock()

		go ds.emitTtl(key, ttl)

		return String(consts.Ok)
	}

	ds.mu.Unlock()
	return NewError(fmt.Sprintf("Key '%v' does not exist", key))
}

func (ds *DataStore) emitTtl(key string, ttl int) {
	e, _ := ds.data[key]
	defer close(e.ttlChan)
	e.ttlChan <- Number(ttl)

	for {
		go func() {
			if len(e.ttlChan) == 1 {
				<-e.ttlChan
			}
		}()

		e.ttlChan <- Number(ttl)

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

	return NewError(fmt.Sprintf("Key '%v' does not exist", key))
}

func (ds *DataStore) lpush(key string, values []string) Data {
	var list List

	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, ok := ds.data[key]; ok {
		if e.value.Type() == TList {
			list = e.value.(List)
			list = append(values, list...)
			e.value = list
			ds.data[key] = e
		} else {
			return NewError(fmt.Sprintf("Cannot lpush on type %v", e.value.Type()))
		}
	} else {
		ds.data[key] = newEntry(List(values))
	}

	return String(consts.Ok)
}

func (ds *DataStore) rpush(key string, values []string) Data {
	var list List

	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, ok := ds.data[key]; ok {
		if e.value.Type() == TList {
			list = e.value.(List)
			list = append(list, values...)
			e.value = list
			ds.data[key] = e
		} else {
			return NewError(fmt.Sprintf("Cannot rpush on type %v", e.value.Type()))
		}
	} else {
		ds.data[key] = newEntry(List(values))
	}

	return String(consts.Ok)
}

func (ds *DataStore) llen(key string) Data {
	var list List
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, exists := ds.data[key]; exists {
		if e.value.Type() == TList {
			list = e.value.(List)
			return Number(len(list))
		}

		return NewError(fmt.Sprintf("Value of type %v does not have property length", e.value.Type()))
	}

	return NewError(fmt.Sprintf("Key '%v' does not exist", key))
}

func (ds *DataStore) lrange(key, startStr, endStr string) Data {
	var list List
	start, err := strconv.Atoi(startStr)

	if err != nil {
		return NewError("Start should be a number value")
	}

	end, err := strconv.Atoi(endStr)

	if err != nil {
		return NewError("End should be a number value")
	}

	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, exists := ds.data[key]; exists {
		if e.value.Type() == TList {
			list = e.value.(List)
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

		return NewError(fmt.Sprintf("Cannot use lrange on %v", e.value.Type()))
	}

	return NewError(fmt.Sprintf("Key '%v' does not exist", key))
}

func (ds *DataStore) ltrim(key, startStr, endStr string) Data {
	var list List
	start, err := strconv.Atoi(startStr)

	if err != nil {
		return ParseError("Start", "int")
	}

	end, err := strconv.Atoi(endStr)

	if err != nil {
		return ParseError("End", "int")
	}

	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, exists := ds.data[key]; exists {
		if e.value.Type() == TList {
			list = e.value.(List)
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

			return String(consts.Ok)
		}

		return NewError(fmt.Sprintf("Cannot use ltrim on %v", e.value.Type()))
	}

	return NewError(fmt.Sprintf("Key '%v' does not exist", key))
}

func (ds *DataStore) sadd(key, value string) Data {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	e, exists := ds.data[key]

	if !exists {
		e = newEntry(NewSet())
	}

	if e.value.Type() != TSet {
		return NewError(fmt.Sprintf("Cannot use sadd on %v", e.value.Type()))
	}

	set := e.value.(Set)
	set.Insert(value)
	ds.data[key] = e
	return String(consts.Ok)
}

func (ds *DataStore) srem(key, value string) Data {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	e, exists := ds.data[key]

	if !exists {
		return NewError(fmt.Sprintf("Key '%v' does not exist", key))
	}

	if e.value.Type() != TSet {
		return NewError(fmt.Sprintf("Cannot use srem on %v", e.value.Type()))
	}

	set := e.value.(Set)
	set.Delete(value)
	e.value = set
	ds.data[key] = e
	return String(consts.Ok)
}

func (ds *DataStore) sismember(key, value string) Data {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	e, exists := ds.data[key]

	if !exists {
		return NewError(fmt.Sprintf("Key '%v' does not exist", key))
	}

	if e.value.Type() != TSet {
		return NewError(fmt.Sprintf("Cannot use sismember on %v", e.value.Type()))
	}

	set := e.value.(Set)

	if set.Has(value) {
		return Number(1)
	}
	return Number(0)
}

func (ds *DataStore) sinter(key, other string) Data {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	keyEntry, keyExists := ds.data[key]
	otherEntry, otherExists := ds.data[other]

	if !(keyExists && otherExists) {
		return NewError(fmt.Sprintf("One of the keys '%v' or '%v' does not exist", key, other))
	}

	if !(keyEntry.value.Type() == TSet && otherEntry.value.Type() == TSet) {
		return NewError(fmt.Sprintf("Canno perform sinter on '%v' or '%v'", key, other))
	}

	keySet := keyEntry.value.(Set)
	otherSet := otherEntry.value.(Set)

	return List(keySet.Intersection(otherSet))
}

func (ds *DataStore) scard(key string) Data {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	e, exists := ds.data[key]

	if !exists {
		return MissingKeyError(key)
	}

	if e.value.Type() != TSet {
		return IncorrectTypeError("scard", e.value.Type())
	}

	set := e.value.(Set)
	return Number(len(set))
}

func (ds *DataStore) zadd(key, value, scoreStr string) Data {
	var sset SortedSet
	score, err := strconv.ParseFloat(scoreStr, 32)

	if err != nil {
		return ParseError("Score", "float")
	}

	ds.mu.Lock()
	defer ds.mu.Unlock()
	e, exists := ds.data[key]

	if !exists {
		sset = NewSortedSet(consts.SortedSetLevels, consts.SortedSetLevelProbability)
		sset.Insert(value, float32(score))

		ds.data[key] = newEntry(sset)
		return String(consts.Ok)
	}

	if e.value.Type() != TSortedSet {
		return IncorrectTypeError("zadd", e.value.Type())
	}

	sset = e.value.(SortedSet)
	sset.Insert(value, float32(score))
	e.value = sset
	ds.data[key] = e
	return String(consts.Ok)
}

func (ds *DataStore) zrem(key, value string) Data {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	e, exists := ds.data[key]

	if !exists {
		return MissingKeyError(key)
	}

	if e.value.Type() != TSortedSet {
		return IncorrectTypeError("zadd", e.value.Type())
	}

	sset := e.value.(SortedSet)
	sset.Delete(value)
	e.value = sset
	ds.data[key] = e
	return String(consts.Ok)
}

func (ds *DataStore) zrank(key, value string) Data {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	e, exists := ds.data[key]

	if !exists {
		return MissingKeyError(key)
	}

	if e.value.Type() != TSortedSet {
		return IncorrectTypeError("zadd", e.value.Type())
	}

	result, err := e.value.(SortedSet).Get(value)

	if err != nil {
		return NewError(fmt.Sprintf("No member %v in %v", value, key))
	}

	return Number(result.Score)
}

func (ds *DataStore) zrange(key, startStr, endStr string) Data {
	start, err := strconv.Atoi(startStr)

	if err != nil {
		return ParseError("Start", "int")
	}

	end, err := strconv.Atoi(endStr)

	if err != nil {
		return ParseError("End", "int")
	}

	ds.mu.Lock()
	defer ds.mu.Unlock()
	e, exists := ds.data[key]

	if !exists {
		return MissingKeyError(key)
	}

	if e.value.Type() != TSortedSet {
		return IncorrectTypeError("zadd", e.value.Type())
	}

	rangeResults := e.value.(SortedSet).Range(start, end)
	var results []string

  for i := range rangeResults {
    results = append(results, rangeResults[i].Value, fmt.Sprintf("%g", rangeResults[i].Score))
	}

  return List(results)
}
