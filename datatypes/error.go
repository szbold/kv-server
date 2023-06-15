package datatypes

import (
	"fmt"
	"key-value-server/consts"
	"strconv"
)

type KvError struct {
	message string
}

const TError string = "error"

func (e KvError) String() string {
	return e.message
}

func NewKvError(input string) KvError {
  return KvError{message: input}
}

func NewIncorrectCommandError(query string) KvError {
  return NewKvError(consts.IncorrectCommand + " " + query)
}

func (e KvError) Type() string {
	return TError
}

func (e KvError) Response() []byte {
	return []byte(fmt.Sprintf("-%v\r\n%v\r\n", strconv.Itoa(len(e.message)), e.message))
}
