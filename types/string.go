package types

import (
	"fmt"
	"strconv"
)

type String string

const TString string = "string"

func (s String) Type() string {
	return TString
}

func (s String) Response() []byte {
	return []byte(fmt.Sprintf("+%v\r\n%v\r\n", strconv.Itoa(len(s)), s))
}
