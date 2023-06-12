package datastore

import "key-value-server/consts"

type entry struct {
	value    string
	dataType dtype
	ttlChan  chan int
}

func (e entry) String() string {
	return e.value + consts.FileDelimiter + e.dataType.String()
}

func newEntry(val string, t dtype) entry {
	return entry{val, t, nil}
}
