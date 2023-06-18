package datastore

import (
	"github.com/stretchr/testify/assert"
	. "key-value-server/datatypes"
	"testing"
)

var ds DataStore = NewDataStore()

func TestKeysCommand(t *testing.T) {
	defer delete(ds.data, "key1")
	defer delete(ds.data, "key2")

	var got KvList
	var want KvList

	ds.data["key1"] = newEntry(KvString("value1"))
	ds.data["key2"] = newEntry(KvString("value2"))

	got = ds.keys()
	want = KvList{"key1", "key2"}

	assert.Equal(t, want, got)

}

func TestGetCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got Data
	var want Data

	got = ds.get("key")
	want = NewKvError("Key 'key' does not exist")

	assert.Equal(t, want, got)

	ds.data["key"] = newEntry(KvString("value"))

	got = ds.get("key")
	want = KvString("value")

	assert.Equal(t, want, got)
}

func TestSetCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got Data
	var want Data
	ds.set("key", KvString("value"))
	got = ds.data["key"].value
	want = KvString("value")

	assert.Equal(t, want, got)
}

func TestIncrCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got Data
	var want Data

	ds.data["key"] = newEntry(KvInt(1))
	ds.incr("key")

	got = ds.data["key"].value
	want = KvInt(2)

	assert.Equal(t, want, got)
}

func TestIncrbyCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got Data
	var want Data

	ds.data["key"] = newEntry(KvInt(1))
	ds.incrby("key", "2")

	got = ds.data["key"].value
	want = KvInt(3)

	assert.Equal(t, want, got)
}

func TestDecrCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got Data
	var want Data

	ds.data["key"] = newEntry(KvInt(1))
	ds.decr("key")

	got = ds.data["key"].value
	want = KvInt(0)

	assert.Equal(t, want, got)
}

func TestDecrbyCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got Data
	var want Data

	ds.data["key"] = newEntry(KvInt(1))
	ds.decrby("key", "2")

	got = ds.data["key"].value
	want = KvInt(-1)

	assert.Equal(t, want, got)
}

func TestExistsCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got Data
	var want Data

	got = ds.exists("key")
	want = KvInt(0)

	assert.Equal(t, want, got)

	ds.data["key"] = newEntry(KvInt(1))

	got = ds.exists("key")
	want = KvInt(1)

	assert.Equal(t, want, got)
}

func TestDeleteCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got bool
	var want bool

	ds.data["key"] = newEntry(KvInt(1))
	ds.del("key")

	_, got = ds.data["key"]
	want = false

	assert.Equal(t, want, got)
}

func TestTypeCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got Data
	var want Data

	ds.data["key"] = newEntry(KvInt(1))

	got = ds.dtype("key")
	want = KvString(TInt)

	assert.Equal(t, want, got)

	ds.data["key"] = newEntry(KvString(""))

	got = ds.dtype("key")
	want = KvString(TString)

	assert.Equal(t, want, got)

	ds.data["key"] = newEntry(KvList([]string{}))

	got = ds.dtype("key")
	want = KvString(TList)

	assert.Equal(t, want, got)
}

func TestLpushCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got Data
	var want Data

	ds.lpush("key", []string{"1", "12", "123"})

	got = ds.data["key"].value
	want = KvList{"1", "12", "123"}

	assert.Equal(t, want, got)

	ds.lpush("key", []string{"9", "8", "7"})

	got = ds.data["key"].value
	want = KvList{"9", "8", "7", "1", "12", "123"}

	assert.Equal(t, want, got)
}

func TestRpushCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got Data
	var want Data

	ds.rpush("key", []string{"1", "12", "123"})

	got = ds.data["key"].value
	want = KvList{"1", "12", "123"}

	assert.Equal(t, want, got)

	ds.rpush("key", []string{"9", "8", "7"})

	got = ds.data["key"].value
	want = KvList{"1", "12", "123", "9", "8", "7"}

	assert.Equal(t, want, got)
}

func TestLlenCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got Data
	var want Data

	ds.data["key"] = newEntry(KvList([]string{"1", "2"}))

	got = ds.llen("key")
	want = KvInt(2)

	assert.Equal(t, want, got)
}

func TestLrangeCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got Data
	var want Data

	ds.data["key"] = newEntry(KvList([]string{"1", "2", "3", "4"}))

	got = ds.lrange("key", "1", "2")
	want = KvList{"2", "3"}

	assert.Equal(t, want, got)
}

func TestLtrimCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got Data
	var want Data

	ds.data["key"] = newEntry(KvList([]string{"1", "2", "3", "4"}))

	ds.ltrim("key", "1", "2")
	got = ds.data["key"].value
	want = KvList{"2", "3"}

	assert.Equal(t, want, got)
}

func TestSaddCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got KvSet
	var want KvSet

	ds.sadd("key", "value")
	ds.sadd("key", "value")
	ds.sadd("key", "value2")

	got = ds.data["key"].value.(KvSet)

	want = NewKvSet()

	want.Insert("value")
	want.Insert("value2")

	assert.Equal(t, want, got)
}

func TestSremCommand(t *testing.T) {
	defer delete(ds.data, "key")
	var got KvSet
	var want KvSet

	ds.sadd("key", "value")
	ds.srem("key", "value")
  got = ds.data["key"].value.(KvSet)

	want = NewKvSet()

	assert.Equal(t, want, got)
}

func TestSismember(t *testing.T) {
	defer delete(ds.data, "key")
	var got Data
	var want Data

	ds.sadd("key", "value")
	got = ds.sismember("key", "value")
	want = KvInt(1)

	assert.Equal(t, want, got)
}

func TestSinter(t *testing.T) {
	defer delete(ds.data, "key")
	var got Data
	var want Data

	ds.sadd("key", "2")
	ds.sadd("key", "3")

	ds.sadd("key2", "1")
	ds.sadd("key2", "3")

	got = ds.sinter("key", "key2")
	want = KvList([]string{"3"})

	assert.Equal(t, want, got)
}

func TestScard(t *testing.T) {
	defer delete(ds.data, "key")
	var got Data
	var want Data

	ds.sadd("key", "1")
	ds.sadd("key", "2")

	got = ds.scard("key")
	want = KvInt(2)

	assert.Equal(t, want, got)
}
