package datatypes

import (
	"fmt"
	"strconv"
	"strings"
)

type KvList []string

const TList string = "list"

func (l KvList) Type() string {
	return TList
}

func (l KvList) Response() []byte {
	var result []string

	for _, elem := range l {
		result = append(result, fmt.Sprintf("$%v\r\n%v\r\n", len(elem), elem))
	}

	return []byte("*" + strconv.Itoa(len(l)) + "\r\n" + strings.Join(result, ""))
}
