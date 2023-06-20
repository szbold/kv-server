package types_test

import (
	. "kv-server/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
  ss := NewSortedSet(3, 0.5)

  ss.Insert("value", 1)

  val, err := ss.Get("value")

  assert.Equal(t, err, nil)
  assert.Equal(t, WithScore{Value: "value", Score: 1}, val)
}

func TestDelete(t *testing.T) {
  ss := NewSortedSet(3, 0.5)

  ss.Insert("value", 1)
  ss.Delete("value")

  val, err := ss.Get("value")

  assert.NotEqual(t, err, nil)
  assert.Equal(t, WithScore{}, val)
}
