package datastore

import (
	"key-value-server/consts"
	"key-value-server/datatypes"
	"strconv"
	"strings"
)

func (ds *DataStore) HandleQuery(query string) []byte {
	if query == "keys" {
		return ds.keys().Response()
	}

	var res datatypes.Data
	q := strings.Split(strings.Trim(query, "\n"), " ")

	if len(q) < 2 {
		return datatypes.NewKvError(consts.IncorrectCommand + " " + query).Response()
	}

	switch q[0] {
	case "get":
		res = ds.get(q[1])
	case "set":
		if len(q) != 3 {
			res = datatypes.NewIncorrectCommandError(query)
      break
		}

		var value datatypes.Data
		num, err := strconv.Atoi(q[2])

		if err != nil {
			value = datatypes.KvString(q[2])
		} else {
			value = datatypes.KvInt(num)
		}

		res = ds.set(q[1], value)
	case "incr":
		res = ds.incr(q[1])
	case "exists":
		res = ds.exists(q[1])
	case "del":
		res = ds.del(q[1])
	case "type":
		res = ds.dtype(q[1])
	case "expire":
		if len(q) != 3 {
			res = datatypes.NewIncorrectCommandError(query)
      break
		}

		res = ds.expire(q[1], q[2])
	case "setexp":
		if len(q) != 4 {
			res = datatypes.NewIncorrectCommandError(query)
      break
		}

		var value datatypes.Data
		num, err := strconv.Atoi(q[2])

		if err != nil {
			value = datatypes.KvString(q[2])
		} else {
			value = datatypes.KvInt(num)
		}

		res = ds.setexp(q[1], value, q[3])
	case "ttl":
		res = ds.ttl(q[1])
	case "lpush":
		if len(q) < 3 {
			res = datatypes.NewIncorrectCommandError(query)
      break
		}

		res = ds.lpush(q[1], q[2:])
	case "rpush":
		if len(q) < 3 {
			res = datatypes.NewIncorrectCommandError(query)
      break
		}

		res = ds.rpush(q[1], q[2:])
	case "llen":
		res = ds.llen(q[1])
	case "lrange":
		if len(q) != 4 {
			res = datatypes.NewIncorrectCommandError(query)
      break
		}
		res = ds.lrange(q[1], q[2], q[3])
	case "ltrim":
		if len(q) != 4 {
			res = datatypes.NewIncorrectCommandError(query)
      break
		}
		res = ds.ltrim(q[1], q[2], q[3])
	default:
			res = datatypes.NewIncorrectCommandError(query)
      break
	}

	return res.Response()
}
