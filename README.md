# Simple key-value datastore (redis clone)

## Data transfered over tcp with REDIS protocol

## Supported types
* int
* string
* list
* set

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

<!-- saving is temporarily disabled -->
<!-- Server dumps data to /tmp/kvdata -->
