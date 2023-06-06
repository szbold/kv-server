package datastore

import (
	"fmt"
	"strconv"
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
func (ds *DataStore) incr(key string) string {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if _, exists := ds.data[key]; exists {
		if ds.data[key].dtype == "int" {
			// error cant occur since data is of int type
			tmp, _ := strconv.Atoi(ds.data[key].value)
			newVal := strconv.Itoa(tmp + 1)
			ds.data[key] = newEntry(newVal)

			return ""
		}
		return fmt.Sprintf("Cannot run incr on value: %v", ds.data[key])
	}

	return fmt.Sprintf("Key \"%v\" does not exist", key)
}
