package main

import (
	kvs "key-value-server/server"
)

func main() {
  // TODO add commads:
  // lists - lrange ltrim
  // ADD SUPPORT FOR INT AND LIST TYPES
  server := kvs.NewKeyValueServer("192.168.1.18:6379")
  server.Run()
}
