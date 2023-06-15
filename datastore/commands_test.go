package datastore

import (
	"github.com/stretchr/testify/assert"
	types "key-value-server/datatypes"
	"testing"
)

var ds DataStore = NewDataStore()

func TestKeys(t *testing.T) {
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

func TestGet(t *testing.T) {
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

func TestSet(t *testing.T) {
	defer delete(ds.data, "key")
	var got types.Data
	var want types.Data
	ds.set("key", types.KvString("value"))
	got = ds.data["key"].value
	want = types.KvString("value")

	assert.Equal(t, want, got)
}

func TestIncr(t *testing.T) {
	defer delete(ds.data, "key")
	var got types.Data
	var want types.Data

	ds.data["key"] = newEntry(types.KvInt(1))
	ds.incr("key")

	got = ds.data["key"].value
	want = types.KvInt(2)

	assert.Equal(t, want, got)
}

func TestExists(t *testing.T) {
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

func TestDelete(t *testing.T) {
	defer delete(ds.data, "key")
	var got bool
	var want bool

	ds.data["key"] = newEntry(types.KvInt(1))
	ds.del("key")

	_, got = ds.data["key"]
	want = false

	assert.Equal(t, want, got)
}

func TestType(t *testing.T) {
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

func TestLpush(t *testing.T) {
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

func TestRpush(t *testing.T) {
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

func TestLlen(t *testing.T) {
	defer delete(ds.data, "key")
	var got types.Data
	var want types.Data

  ds.data["key"] = newEntry(types.KvList([]string{"1", "2"}))

	got = ds.llen("key")
	want = types.KvInt(2)

  assert.Equal(t, want, got)
}

func TestLrange(t *testing.T) {
	defer delete(ds.data, "key")
	var got types.Data
	var want types.Data

  ds.data["key"] = newEntry(types.KvList([]string{"1", "2", "3", "4"}))

	got = ds.lrange("key", "1", "2")
	want = types.KvList{"2", "3"}

  assert.Equal(t, want, got)
}

func TestLtrim(t *testing.T) {
	defer delete(ds.data, "key")
	var got types.Data
	var want types.Data

  ds.data["key"] = newEntry(types.KvList([]string{"1", "2", "3", "4"}))

	ds.ltrim("key", "1", "2")
  got = ds.data["key"].value
	want = types.KvList{"2", "3"}

  assert.Equal(t, want, got)
}
