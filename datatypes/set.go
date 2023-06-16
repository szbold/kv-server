package datatypes

import (
	"fmt"
	"key-value-server/consts"
	"strconv"
	"strings"
)

type void struct{}
type KvSet map[string]void

const TSet string = "set"

func (s KvSet) String() string {
	var values []string

	for val := range s {
		values = append(values, val)
	}

	return strings.Join(values, consts.ListDelimiter)
}

func (s KvSet) Type() string {
	return TSet
}

func (s KvSet) Response() []byte {
	var result []string

	for elem := range s {
		result = append(result, fmt.Sprintf("$%v\r\n%v\r\n", len(elem), elem))
	}

	return []byte("*" + strconv.Itoa(len(s)) + "\r\n" + strings.Join(result, ""))
}

func NewKvSet() KvSet {
  return make(map[string]void)
}

// METHODS FOR SET OPERATIONS
func (s *KvSet) Insert(value string) {
	if _, ok := (*s)[value]; !ok {
		(*s)[value] = void{}
	}
}

func (s KvSet) Has(value string) bool {
	_, ok := s[value]
	return ok
}

func (s *KvSet) Delete(value string) {
	delete(*s, value)
}

func (s *KvSet) Intersection(other KvSet) []string {
	var result []string
	for value := range *s {
		if other.Has(value) {
			result = append(result, value)
		}
	}

	return result
}
