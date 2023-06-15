package datastore

import (
	"key-value-server/consts"
	"key-value-server/datatypes"
	"key-value-server/fmt"
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
		return fmt.ErrResponse(consts.IncorrectCommand + " " + query)
	}

	switch q[0] {
	case "get":
		res = ds.get(q[1])
	case "set":
		if len(q) != 3 {
			return fmt.ErrResponse(consts.IncorrectCommand + " " + query)
		}

		var value datatypes.Data
		num, err := strconv.Atoi(q[2])

		if err != nil {
			value = datatypes.KvString(q[2])
		} else {
			value = datatypes.KvInt(num)
		}

		ds.set(q[1], value)
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
			return fmt.ErrResponse(consts.IncorrectCommand + " " + query)
		}

		res = ds.expire(q[1], q[2])
	case "setexp":
		if len(q) != 4 {
			return fmt.ErrResponse(consts.IncorrectCommand + " " + query)
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
			return fmt.ErrResponse(consts.IncorrectCommand + " " + query)
		}

		ds.lpush(q[1], q[2:])
	case "rpush":
		if len(q) < 3 {
			return fmt.ErrResponse(consts.IncorrectCommand + " " + query)
		}

		ds.rpush(q[1], q[2:])
	case "llen":
		res = ds.llen(q[1])
	case "lrange":
		if len(q) != 4 {
			return fmt.ErrResponse(consts.IncorrectCommand + " " + query)
		}
		res = ds.lrange(q[1], q[2], q[3])
	case "ltrim":
		if len(q) != 4 {
			return fmt.ErrResponse(consts.IncorrectCommand + " " + query)
		}
		res = ds.ltrim(q[1], q[2], q[3])
	default:
		return fmt.ErrResponse(consts.IncorrectCommand + " " + query)
	}

	return res.Response()
}
