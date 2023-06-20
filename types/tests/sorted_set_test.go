package types_test

import (
	. "kv-server/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
  var ss SortedSet

  ss.Insert("value", 1)

  val, err := ss.Get("value")

  assert.Equal(t, err, nil)
  assert.Equal(t, WithScore{Value: "value", Score: 1}, val)
}
