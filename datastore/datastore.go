package datastore

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

const DEFAULT_FILE_PATH = "/tmp/kvdata"

type DataStore struct {
	data map[string]string
	mu   sync.Mutex
}

func (ds *DataStore) String() string {
	var result string

	for key, value := range ds.data {
		result += key + "," + value + "\n"
	}

	return result
}

func NewDataStore() DataStore {
	return DataStore{data: make(map[string]string)}
}

func (ds *DataStore) HandleQuery(query string) string {
	q := strings.Split(strings.Trim(query, "\n"), " ")

	if len(q) == 2 && q[0] == "get" {
		return ds.get(q[1])
	} else if len(q) == 3 && q[0] == "set" {
		ds.set(q[1], q[2])
		return "OK"
	}

	return "INCORRECT COMMAND"
}

func (ds *DataStore) set(key, value string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.data[key] = value
}

func (ds *DataStore) get(key string) string {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if value, exists := ds.data[key]; exists {
		return value
	}

	return "nil"
}

func (ds *DataStore) Load() error {
	file, err := os.Open(DEFAULT_FILE_PATH)

	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var line string
	lineIdx := 1
	for scanner.Scan() {
		line = scanner.Text()

		kv := strings.Split(line, ",")

		if len(kv) != 2 {
			return errors.New(fmt.Sprint("Data posssibly corrupted on line ", lineIdx))
		}

		ds.data[kv[0]] = kv[1]
		lineIdx++
	}

	log.Println("Data loaded successfully")

	return nil
}

// TODO change this to maybe incude sequential writes to file
// comparing keys and values and only writing them if they are changed
func (ds *DataStore) Dump() error {
	fmt.Println(ds)
	err := os.WriteFile(DEFAULT_FILE_PATH, []byte(ds.String()), 0644)

	if err != nil {
		return err
	}

	log.Println("Data dumped successfully")

	return nil
}
