package types_test

import (
	"testing"
  . "kv-server/types"
	"github.com/stretchr/testify/assert"
)

func TestStringResponse(t *testing.T) {
	var s String
	s = "example"
	got := s.Response()
	want := []byte("+7\r\nexample\r\n")

  assert.Equal(t, want, got)

}

func TestIntResponse(t *testing.T) {
	var s Int
	s = 123
	got := s.Response()
	want := []byte(":123\r\n")

  assert.Equal(t, want, got)
}

func TestListResponse(t *testing.T) {
	var l List
	l = []string{"1", "23", "456"}
	got := l.Response()
	want := []byte("*3\r\n$1\r\n1\r\n$2\r\n23\r\n$3\r\n456\r\n")

  assert.Equal(t, want, got)
}

func TestErrorResponse(t *testing.T) {
	var s Error
	s = NewError("Error")
	got := s.Response()
	want := []byte("-5\r\nError\r\n")

  assert.Equal(t, want, got)
}

func TestSetResponse(t *testing.T) {
	var l Set
	l = NewSet()
  l.Insert("1")
  l.Insert("23")
  l.Insert("456")

	got := l.Response()
	want := []byte("*3\r\n$1\r\n1\r\n$2\r\n23\r\n$3\r\n456\r\n")

  assert.Equal(t, want, got)
}

