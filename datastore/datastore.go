package datastore

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

const _DELIMITER = ";"
const list_delimiter = ","
const (
	missing_command   = "Missing command"
	incorrect_command = "Command incorrect"
)

// think of separating this
func stringResponse(message string) string {
	return fmt.Sprintf("+%v\r\n%v\r\n", len(message), message)
}

func intResponse(message string) string {
	return fmt.Sprintf(":%v\r\n", message)
}

func listResponse(message string) string {
  var result string
  var listLength int

  for _, item := range strings.Split(message, list_delimiter) {
    result += fmt.Sprintf("$%v\r\n%v\r\n", len(item), item)
    listLength++
  }

	return fmt.Sprintf("*%v\r\n%v\r\n", listLength, result)
}

func errResponse(message string) string {
	return fmt.Sprintf("-%v\r\n%v\r\n", len(message), message)
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
  if query == "keys" {
    return listResponse(ds.keys())
  }

	var string_res string
  var int_res string
  var list_res string
	var err error
	q := strings.Split(strings.Trim(query, "\n"), " ")


	if len(q) < 2 {
		return errResponse(incorrect_command + " " + query)
	}

	// might refactor into less returns and just set res and err in every case
	switch q[0] {
	case "get":
		string_res, err = ds.get(q[1])
	case "set":
		if len(q) != 3 {
			return errResponse(incorrect_command + " " + query)
		}
		ds.set(q[1], q[2])
	case "incr":
		err = ds.incr(q[1])
	case "exists":
		exists := ds.exists(q[1])
		string_res = strconv.FormatBool(exists)
	case "del":
		err = ds.del(q[1])
	case "type":
		string_res, err = ds.dtype(q[1])
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
		int_res, err = ds.ttl(q[1])
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
		int_res, err = ds.llen(q[1])
	case "lrange":
		if len(q) != 4 {
			return errResponse(incorrect_command + " " + query)
		}
		list_res, err = ds.lrange(q[1], q[2], q[3])
	case "ltrim":
		if len(q) != 4 {
			return errResponse(incorrect_command + " " + query)
		}
		err = ds.ltrim(q[1], q[2], q[3])
	default:
		return errResponse(incorrect_command + " " + query)
	}

	if err != nil {
		return errResponse(err.Error())
	}
  
  if string_res != "" {
    return stringResponse(string_res)
  }

  if int_res != "" {
    return intResponse(int_res)
  }

  if list_res != "" {
    return listResponse(list_res)
  }

	return stringResponse("OK")
}
