package main

import (
	kvs "kv-server/server"
)

func main() {
  server := kvs.NewKeyValueServer("192.168.1.36:6379")
  server.Run()
}
