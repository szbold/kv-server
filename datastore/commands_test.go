package datastore

import (
	"github.com/stretchr/testify/assert"
	types "key-value-server/datatypes"
	"testing"
)

var ds DataStore = NewDataStore()

func TestKeysCommand(t *testing.T) {
	defer delete(ds.data, "key1")
	defer delete(ds.data, "key2")

	var got types.KvList
	var want types.KvList

	ds.data["key1"] = newEntry(types.KvString("value1"))
	ds.data["key2"] = newEntry(types.KvString("value2"))

	got = ds.keys()
	want = types.KvList{"key1", "key2"}

	assert.Equal(t, want, got)

}

func TestGetCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got types.Data
	var want types.Data

	got = ds.get("key")
	want = types.NewKvError("Key 'key' does not exist")

	assert.Equal(t, want, got)

	ds.data["key"] = newEntry(types.KvString("value"))

	got = ds.get("key")
	want = types.KvString("value")

	assert.Equal(t, want, got)
}

func TestSetCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got types.Data
	var want types.Data
	ds.set("key", types.KvString("value"))
	got = ds.data["key"].value
	want = types.KvString("value")

	assert.Equal(t, want, got)
}

func TestIncrCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got types.Data
	var want types.Data

	ds.data["key"] = newEntry(types.KvInt(1))
	ds.incr("key")

	got = ds.data["key"].value
	want = types.KvInt(2)

	assert.Equal(t, want, got)
}

func TestExistsCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got types.Data
	var want types.Data

	got = ds.exists("key")
	want = types.KvInt(0)

	assert.Equal(t, want, got)

	ds.data["key"] = newEntry(types.KvInt(1))

	got = ds.exists("key")
	want = types.KvInt(1)

	assert.Equal(t, want, got)
}

func TestDeleteCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got bool
	var want bool

	ds.data["key"] = newEntry(types.KvInt(1))
	ds.del("key")

	_, got = ds.data["key"]
	want = false

	assert.Equal(t, want, got)
}

func TestTypeCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got types.Data
	var want types.Data

	ds.data["key"] = newEntry(types.KvInt(1))

	got = ds.dtype("key")
	want = types.KvString(types.TInt)

	assert.Equal(t, want, got)

	ds.data["key"] = newEntry(types.KvString(""))

	got = ds.dtype("key")
	want = types.KvString(types.TString)

	assert.Equal(t, want, got)

	ds.data["key"] = newEntry(types.KvList([]string{}))

	got = ds.dtype("key")
	want = types.KvString(types.TList)

	assert.Equal(t, want, got)
}

func TestLpushCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got types.Data
	var want types.Data

	ds.lpush("key", []string{"1", "12", "123"})

	got = ds.data["key"].value
	want = types.KvList{"1", "12", "123"}

	assert.Equal(t, want, got)

	ds.lpush("key", []string{"9", "8", "7"})

	got = ds.data["key"].value
	want = types.KvList{"9", "8", "7", "1", "12", "123"}

	assert.Equal(t, want, got)
}

func TestRpushCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got types.Data
	var want types.Data

	ds.rpush("key", []string{"1", "12", "123"})

	got = ds.data["key"].value
	want = types.KvList{"1", "12", "123"}

	assert.Equal(t, want, got)

	ds.rpush("key", []string{"9", "8", "7"})

	got = ds.data["key"].value
	want = types.KvList{"1", "12", "123", "9", "8", "7"}

	assert.Equal(t, want, got)
}

func TestLlenCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got types.Data
	var want types.Data

  ds.data["key"] = newEntry(types.KvList([]string{"1", "2"}))

	got = ds.llen("key")
	want = types.KvInt(2)

  assert.Equal(t, want, got)
}

func TestLrangeCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got types.Data
	var want types.Data

  ds.data["key"] = newEntry(types.KvList([]string{"1", "2", "3", "4"}))

	got = ds.lrange("key", "1", "2")
	want = types.KvList{"2", "3"}

  assert.Equal(t, want, got)
}

func TestLtrimCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got types.Data
	var want types.Data

  ds.data["key"] = newEntry(types.KvList([]string{"1", "2", "3", "4"}))

	ds.ltrim("key", "1", "2")
  got = ds.data["key"].value
	want = types.KvList{"2", "3"}

  assert.Equal(t, want, got)
}
