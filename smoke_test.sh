#!/bin/sh

docker run -d kv-server > id.tmp

if [ ! -z "$(docker ps | grep kv-server)" ]; then
  # container running correctly
  docker stop $(cat id.tmp)
  rm id.tmp
else 
  rm id.tmp
  exit 1
fi
