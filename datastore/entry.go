package datastore

import (
	. "kv-server/types"
)

type entry struct {
	value   Data
	ttlChan chan Number
}

func newEntry(val Data) entry {
	return entry{val, nil}
}
