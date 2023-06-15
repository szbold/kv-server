package datatypes

import (
	"fmt"
	"strconv"
)

type KvInt int

const TInt string = "int"

func (i KvInt) String() string {
	return strconv.Itoa(int(i))
}

func (i KvInt) Type() string {
	return TInt
}

func (i KvInt) Response() []byte {
	return []byte(fmt.Sprintf(":%v\r\n", i.String()))
}
