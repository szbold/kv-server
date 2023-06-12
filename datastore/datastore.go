package datastore

import (
  "key-value-server/consts"
	"sync"
)

type DataStore struct {
	data map[string]entry
	mu   sync.Mutex
}

func (ds *DataStore) String() string {
	var result string

	for key, value := range ds.data {
		result += key + consts.FileDelimiter + value.String() + "\n"
	}

	return result
}

func NewDataStore() DataStore {
	return DataStore{data: make(map[string]entry)}
}

