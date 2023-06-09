package datastore

import (
	"errors"
	"fmt"
)

type dtype int

const (
	t_string dtype = iota
	t_int
	t_list
)

func stringToDtype(input string) (dtype, error) {
  switch input {
  case "string":
    return t_string, nil
  case "int":
    return t_int, nil
  case "list":
    return t_list, nil
  }

  // also should never occur, just a safety measure
  return t_string, errors.New(fmt.Sprintf("%v is not a type", input))
}

func (input dtype) String() string {
  switch input {
  case t_string:
    return "string"
  case t_int:
    return "int"
  case t_list:
    return "list"
  }

  // should never occur
  return ""
}
