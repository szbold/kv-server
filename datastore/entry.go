package datastore

type entry struct {
	value    string
	dataType dtype
	ttlChan  chan int
}

func (e entry) String() string {
	return e.value + _DELIMITER + e.dataType.String()
}

func newEntry(val string, t dtype) entry {
	return entry{val, t, nil}
}
