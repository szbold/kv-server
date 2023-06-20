package types

import (
	"fmt"
	"kv-server/consts"
	"strconv"
)

type Error struct {
	message string
}

const TError string = "error"

func NewError(input string) Error {
  return Error{message: input}
}

func NewIncorrectCommandError(query string) Error {
  return NewError(consts.IncorrectCommand + " " + query)
}

func (e Error) Type() string {
	return TError
}

func (e Error) Response() []byte {
	return []byte(fmt.Sprintf("-%v\r\n%v\r\n", strconv.Itoa(len(e.message)), e.message))
}
