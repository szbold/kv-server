package main

import (
	kvs "key-value-server/server"
)

func main() {
  // TODO add commads:
  // ttl
  // lists - lpush rpush lrange ltrim
  // ADD SUPPORT FOR INT AND LIST TYPES
  server := kvs.NewKeyValueServer("192.168.1.36:6379")
  server.Run()
}
