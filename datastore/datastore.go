package datastore

import (
	"sync"
)

type DataStore struct {
	data map[string]entry
	mu   sync.Mutex
}

func NewDataStore() DataStore {
	return DataStore{data: make(map[string]entry)}
}
