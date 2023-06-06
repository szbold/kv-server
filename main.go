package main

import (
	kvs "key-value-server/server"
)

func main() {
  // TODO add commads:
  // incr, exists, del, type, expire, setexp, ttl
  // lists - lpush rpush lrange ltrim
  // ADD SUPPORT FOR INT AND LIST TYPES
  server := kvs.NewKeyValueServer("127.0.0.1:6379")
  server.Run()
}
