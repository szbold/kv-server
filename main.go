package main

import (
	kvs "kv-server/server"
)

func main() {
	server := kvs.NewKeyValueServer("0.0.0.0:6379")
	server.Run()
}
