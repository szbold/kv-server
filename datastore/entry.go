package datastore

import (
	"key-value-server/consts"
	"key-value-server/datatypes"
)

type entry struct {
	value   datatypes.Data
	ttlChan chan datatypes.KvInt
}

func (e entry) String() string {
	return e.value.String() + consts.FileDelimiter + e.value.Type()
}

func newEntry(val datatypes.Data) entry {
	return entry{val, nil}
}
