# Simple key-value datastore (redis clone)

## Data transfered over tcp with REDIS protocol

## Supported types
* number (stored as float32, but some operations prevent floats from being written)
* string
* list
* set
* sorted set

## Supported commands
### Genral
* get key
* set key val
* exists key
* del key
* dtype key
* expire key ttl
* setexp key val ttl
* ttl key
### Integers
* incr key
* incrby key increment
* decr key
* decrby key decrement
### Lists
* lpush key ...values
* rpush key ...values
* llen key
* lrange key start stop\*
* ltrim key start stop\*
  - start is swapped with end if end is smaller than start
  - if start is negative it starts from the beginning of the list
  - if end is larger than the length of the list then it is sent to the length of the list
### Sets
* sadd key value
* srem key value
* sismember key value
* sinter key other_key
* scard key
### Sorted sets
* zadd key value score
* zrem key value
* zrank key value
* zrange key start end

<!-- saving is temporarily disabled -->
<!-- Server dumps data to /tmp/kvdata -->
