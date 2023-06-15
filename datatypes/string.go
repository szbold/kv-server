package datatypes

import (
	"fmt"
	"strconv"
)

type KvString string

const TString string = "string"

func (s KvString) String() string {
	return string(s)
}

func (s KvString) Type() string {
	return TString
}

func (s KvString) Response() []byte {
	return []byte(fmt.Sprintf("+%v\r\n%v\r\n", strconv.Itoa(len(s)), s.String()))
}
