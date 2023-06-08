package datastore

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

const _DELIMITER = ";"
const (
	missing_command   = "Missing command"
	incorrect_command = "Command incorrect"
)

// think of separating this
func okResponse(message string) string {
	return fmt.Sprintf("[OK] %v", message)
}

func errResponse(message string) string {
	return fmt.Sprintf("[ERR] %v", message)
}

type entry struct {
	value string
	dtype string
	ttl   int
}

type DataStore struct {
	data map[string]entry
	mu   sync.Mutex
}

func (e entry) String() string {
	return e.value
}

func newEntry(val string) entry {
	var dtype string
	_, err := strconv.Atoi(val)

	if err != nil {
		dtype = "string"
	} else {
		dtype = "int"
	}

	return entry{val, dtype, -1}
}

func newEntryExp(val string, ttl uint) entry {
	var dtype string
	_, err := strconv.Atoi(val)

	if err != nil {
		dtype = "string"
	} else {
		dtype = "int"
	}

	return entry{val, dtype, int(ttl)}
}

func (ds *DataStore) String() string {
	var result string

	for key, value := range ds.data {
		result += key + _DELIMITER + value.String() + "\n"
	}

	return result
}

func NewDataStore() DataStore {
	return DataStore{data: make(map[string]entry)}
}

func (ds *DataStore) HandleQuery(query string) string {
	var res string
	var err error
	q := strings.Split(strings.Trim(query, "\n"), " ")

	if len(q) < 2 {
		return errResponse(incorrect_command + " " + query)
	}

	switch q[0] {
	case "get":
		res = ds.get(q[1])
		return okResponse(res)
	case "set":
		if len(q) != 3 {
			return errResponse(incorrect_command + " " + query)
		}
		ds.set(q[1], q[2])
		return okResponse(res)
	case "incr":
		err = ds.incr(q[1])
		if err != nil {
			return errResponse(err.Error())
		}
		return okResponse(res)
	case "exists":
		exists := ds.exists(q[1])
		return okResponse(strconv.FormatBool(exists))
	case "del":
		err = ds.del(q[1])
		if err != nil {
			return errResponse(err.Error())
		}
		return okResponse(res)
	case "type":
		res, err = ds.dtype(q[1])
		if err != nil {
			return errResponse(err.Error())
		}
		return okResponse(res)
	case "expire":
		if len(q) != 3 {
			return errResponse(incorrect_command + " " + query)
		}

		err = ds.expire(q[1], q[2])

		if err != nil {
			return errResponse(err.Error())
		}
		return okResponse(res)
	}

	return errResponse(incorrect_command + " " + query)
}
