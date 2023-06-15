package fmt

import (
	"fmt"
)

func ErrResponse(message string) []byte {
	return []byte(fmt.Sprintf("-%v\r\n%v\r\n", len(message), message))
}
