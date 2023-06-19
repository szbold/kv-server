package datastore

import (
	"key-value-server/datatypes"
)

type entry struct {
	value   datatypes.Data
	ttlChan chan datatypes.KvInt
}

func newEntry(val datatypes.Data) entry {
	return entry{val, nil}
}
