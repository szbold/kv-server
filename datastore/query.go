package datastore

import (
	. "key-value-server/datatypes"
	"strconv"
	"strings"
)

func (ds *DataStore) HandleQuery(query string) []byte {
	var res Data = NewIncorrectCommandError(query) // default response is error, which can be later changed
	q := strings.Split(strings.Trim(query, "\n"), " ")

	switch len(q) {
	case 0:
	case 1:
		if query == "keys" {
			res = ds.keys()
		}
	case 2:
		switch q[0] {
		case "get":
			res = ds.get(q[1])
		case "incr":
			res = ds.incr(q[1])
		case "decr":
			res = ds.decr(q[1])
		case "exists":
			res = ds.exists(q[1])
		case "del":
			res = ds.del(q[1])
		case "type":
			res = ds.dtype(q[1])
		case "ttl":
			res = ds.ttl(q[1])
		case "llen":
			res = ds.llen(q[1])
		case "scard":
			res = ds.scard(q[1])
		}
	case 3:
		switch q[0] {
		case "set":
			var value Data
			num, err := strconv.Atoi(q[2])

			if err != nil {
				value = KvString(q[2])
			} else {
				value = KvInt(num)
			}

			res = ds.set(q[1], value)
		case "expire":
			res = ds.expire(q[1], q[2])
		case "sadd":
			res = ds.sadd(q[1], q[2])
		case "srem":
			res = ds.srem(q[1], q[2])
		case "sismember":
			res = ds.sismember(q[1], q[2])
		case "sinter":
			res = ds.sinter(q[1], q[2])
    case "incrby":
      res = ds.incrby(q[1], q[2])
    case "decrby":
      res = ds.decrby(q[1], q[2])
		}
	case 4:
		switch q[0] {
		case "setexp":
			var value Data
			num, err := strconv.Atoi(q[2])

			if err != nil {
				value = KvString(q[2])
			} else {
				value = KvInt(num)
			}

			res = ds.setexp(q[1], value, q[3])
		case "lrange":
			res = ds.lrange(q[1], q[2], q[3])
		case "ltrim":
			res = ds.ltrim(q[1], q[2], q[3])
		}
	default:
		switch q[0] {
		case "lpush":
			res = ds.lpush(q[1], q[2:])
		case "rpush":
			res = ds.rpush(q[1], q[2:])
		}
	}

	return res.Response()
}
