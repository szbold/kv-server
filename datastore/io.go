package datastore

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
  "strings"
)

const _DEFAULT_FILE_PATH = "/tmp/kvdata"

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

