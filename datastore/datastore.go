package datastore

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

const _DEFAULT_FILE_PATH = "/tmp/kvdata"
const _DELIMITER = ";"

type entry struct {
	value string
	dtype string
}

type DataStore struct {
	data map[string]entry
	mu   sync.Mutex
}

func (e entry) String() string {
	return e.value
}

func newEntry(val string) entry {
  var dtype string
  _, err := strconv.Atoi(val)

  if err != nil {
    dtype = "string"
  } else {
    dtype = "int"
  }

	return entry{val, dtype}
}

func (ds *DataStore) String() string {
	var result string

	for key, value := range ds.data {
		result += key + _DELIMITER + value.String() + "\n"
	}

	return result
}

func NewDataStore() DataStore {
	return DataStore{data: make(map[string]entry)}
}

func (ds *DataStore) HandleQuery(query string) string {
	q := strings.Split(strings.Trim(query, "\n"), " ")

	switch q[0] {
	case "get":
		return ds.get(q[1])
	case "set":
		ds.set(q[1], q[2])
		return "OK"
	}

	return "INCORRECT COMMAND"
}

func (ds *DataStore) set(key, value string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.data[key] = newEntry(value)
}

func (ds *DataStore) get(key string) string {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if e, exists := ds.data[key]; exists {
		return e.value
	}

	return "nil"
}

func (ds *DataStore) Load() error {
	file, err := os.OpenFile(_DEFAULT_FILE_PATH, os.O_RDONLY|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var line string
	lineIdx := 1
	for scanner.Scan() {
		line = scanner.Text()

		kv := strings.Split(line, _DELIMITER)

		if len(kv) != 2 {
			return errors.New(fmt.Sprintf("Data posssibly corrupted on line %v\n%v", lineIdx, line))
		}

		ds.data[kv[0]] = newEntry(kv[1])
		lineIdx++
	}

	log.Println("Data loaded successfully")

	return nil
}

// TODO change this to maybe incude sequential writes to file
// comparing keys and values and only writing them if they are changed
func (ds *DataStore) Dump() error {
	err := os.WriteFile(_DEFAULT_FILE_PATH, []byte(ds.String()), 0644)

	if err != nil {
		return err
	}

	log.Println("Data dumped successfully")

	return nil
}
