package types

import (
	"fmt"
)

type Number float32

const TNumber = "number"

func (i Number) Type() string {
	return TNumber
}

func (i Number) Response() []byte {
	return []byte(fmt.Sprintf(":%v\r\n", fmt.Sprintf("%g", i)))
}
