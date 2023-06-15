package datatypes

import (
	"fmt"
	"strconv"
)

type KvError struct {
	Message string
}

const TError string = "error"

func (e KvError) String() string {
	return e.Message
}

func FromString(input string) KvError {
  return KvError{input}
}

func (e KvError) Type() string {
	return TError
}

func (e KvError) Response() []byte {
	return []byte(fmt.Sprintf("-%v\r\n%v\r\n", strconv.Itoa(len(e.Message)), e.Message))
}
