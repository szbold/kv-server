package types

import (
	"fmt"
	"strconv"
)

type Int int

const TInt string = "int"

func (i Int) Type() string {
	return TInt
}

func (i Int) Response() []byte {
	return []byte(fmt.Sprintf(":%v\r\n", strconv.Itoa(int(i))))
}
