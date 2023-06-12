package fmt

import (
	"fmt"
	"strings"
  "key-value-server/consts"
)

func StringResponse(message string) string {
	return fmt.Sprintf("+%v\r\n%v\r\n", len(message), message)
}

func IntResponse(message string) string {
	return fmt.Sprintf(":%v\r\n", message)
}

func ListResponse(message string) string {
	var result string
	var listLength int

	for _, item := range strings.Split(message, consts.ListDelimiter) {
		result += fmt.Sprintf("$%v\r\n%v\r\n", len(item), item)
		listLength++
	}

	return fmt.Sprintf("*%v\r\n%v", listLength, result)
}

func ErrResponse(message string) string {
	return fmt.Sprintf("-%v\r\n%v\r\n", len(message), message)
}
