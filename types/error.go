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

func IncorrectCommandError(query string) Error {
	return NewError(consts.IncorrectCommand + " " + query)
}

func MissingKeyError(key string) Error {
	return NewError(fmt.Sprintf("Key '%v' does not exist", key))

}
func IncorrectTypeError(command, datatype string) Error {
	return NewError(fmt.Sprintf("Cannot use %v on %v", command, datatype))
}

func ParseError(field, tWant string) Error {
	return NewError(fmt.Sprintf("%v should be %v found string", field, tWant))
}

func (e Error) Type() string {
	return TError
}

func (e Error) Response() []byte {
	return []byte(fmt.Sprintf("-%v\r\n%v\r\n", strconv.Itoa(len(e.message)), e.message))
}
