package datastore

import (
	. "kv-server/types"
)

type entry struct {
	value   Data
	ttlChan chan Int
}

func newEntry(val Data) entry {
	return entry{val, nil}
}
