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

type DataStore struct {
	data map[string]entry
	mu   sync.Mutex
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

	// might refactor into less returns and just set res and err in every case
	switch q[0] {
	case "get":
		res, err = ds.get(q[1])
	case "set":
		if len(q) != 3 {
			return errResponse(incorrect_command + " " + query)
		}
		ds.set(q[1], q[2])
	case "incr":
		err = ds.incr(q[1])
	case "exists":
		exists := ds.exists(q[1])
		res = strconv.FormatBool(exists)
	case "del":
		err = ds.del(q[1])
	case "type":
		res, err = ds.dtype(q[1])
	case "expire":
		if len(q) != 3 {
			return errResponse(incorrect_command + " " + query)
		}

		err = ds.expire(q[1], q[2])
	case "setexp":
		if len(q) != 4 {
			return errResponse(incorrect_command + " " + query)
		}

		err = ds.setexp(q[1], q[2], q[3])
	case "ttl":
		res, err = ds.ttl(q[1])
	case "lpush":
		if len(q) < 3 {
			return errResponse(incorrect_command + " " + query)
		}

		ds.lpush(q[1], q[2:])
	case "rpush":
		if len(q) < 3 {
			return errResponse(incorrect_command + " " + query)
		}

		ds.rpush(q[1], q[2:])
	case "llen":
		res, err = ds.llen(q[1])
	case "lrange":
		if len(q) != 4 {
			return errResponse(incorrect_command + " " + query)
		}
		res, err = ds.lrange(q[1], q[2], q[3])
	default:
		return errResponse(incorrect_command + " " + query)
	}

	if err != nil {
		return errResponse(err.Error())
	}
	return okResponse(res)
}
