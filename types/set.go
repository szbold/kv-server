package types

import (
	"fmt"
	"strconv"
	"strings"
)

type void struct{}
type Set map[string]void

const TSet string = "set"

func (s Set) Type() string {
	return TSet
}

func (s Set) Response() []byte {
	var result []string

	for elem := range s {
		result = append(result, fmt.Sprintf("$%v\r\n%v\r\n", len(elem), elem))
	}

	return []byte("*" + strconv.Itoa(len(s)) + "\r\n" + strings.Join(result, ""))
}

func NewSet() Set {
  return make(map[string]void)
}

// METHODS FOR SET OPERATIONS
func (s *Set) Insert(value string) {
	if _, ok := (*s)[value]; !ok {
		(*s)[value] = void{}
	}
}

func (s Set) Has(value string) bool {
	_, ok := s[value]
	return ok
}

func (s *Set) Delete(value string) {
	delete(*s, value)
}

func (s *Set) Intersection(other Set) []string {
	var result []string
	for value := range *s {
		if other.Has(value) {
			result = append(result, value)
		}
	}

	return result
}
