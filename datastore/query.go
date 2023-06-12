package datastore

import (
	"key-value-server/consts"
	"key-value-server/fmt"
	"strings"
)

func (ds *DataStore) HandleQuery(query string) string {
  if query == "keys" {
    return fmt.ListResponse(ds.keys())
  }

	var string_res string
  var int_res string
  var list_res string
	var err error
	q := strings.Split(strings.Trim(query, "\n"), " ")


	if len(q) < 2 {
		return fmt.ErrResponse(consts.IncorrectCommand + " " + query)
	}

	switch q[0] {
	case "get":
		string_res, err = ds.get(q[1])
	case "set":
		if len(q) != 3 {
			return fmt.ErrResponse(consts.IncorrectCommand + " " + query)
		}
		ds.set(q[1], q[2])
	case "incr":
		err = ds.incr(q[1])
	case "exists":
		string_res = ds.exists(q[1])
	case "del":
		err = ds.del(q[1])
	case "type":
		string_res, err = ds.dtype(q[1])
	case "expire":
		if len(q) != 3 {
			return fmt.ErrResponse(consts.IncorrectCommand + " " + query)
		}

		err = ds.expire(q[1], q[2])
	case "setexp":
		if len(q) != 4 {
			return fmt.ErrResponse(consts.IncorrectCommand + " " + query)
		}

		err = ds.setexp(q[1], q[2], q[3])
	case "ttl":
		int_res, err = ds.ttl(q[1])
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
		int_res, err = ds.llen(q[1])
	case "lrange":
		if len(q) != 4 {
			return fmt.ErrResponse(consts.IncorrectCommand + " " + query)
		}
		list_res, err = ds.lrange(q[1], q[2], q[3])
	case "ltrim":
		if len(q) != 4 {
			return fmt.ErrResponse(consts.IncorrectCommand + " " + query)
		}
		err = ds.ltrim(q[1], q[2], q[3])
	default:
		return fmt.ErrResponse(consts.IncorrectCommand + " " + query)
	}

	if err != nil {
		return fmt.ErrResponse(err.Error())
	}
  
  if string_res != "" {
    return fmt.StringResponse(string_res)
  }

  if int_res != "" {
    return fmt.IntResponse(int_res)
  }

  if list_res != "" {
    return fmt.ListResponse(list_res)
  }

	return fmt.StringResponse("OK")
}
