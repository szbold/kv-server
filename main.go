package main

import (
	kvs "kv-server/server"
)

func main() {
  server := kvs.NewKeyValueServer("127.0.0.1:6379")
  server.Run()
}
