# Simple key-value datastore (redis clone)

## Data transfered over tcp with REDIS protocol

## Supported types
* int
* string
* list

### Supported commands
* get key
* set key val
* increment key
* exists key
* delete key
* type key
* expire key ttl
* setexp key val ttl
* ttl key
* lpush key ...values
* rpush key ...values
* llen key
* lrange key start stop\*
* ltrim key start stop\*
  - start is swapped with end if end is smaller than start
  - if start is negative it starts from the beginning of the list
  - if end is larger than the length of the list then it is sent to the length of the list

<!-- saving is temporarily disabled -->
<!-- Server dumps data to /tmp/kvdata -->
