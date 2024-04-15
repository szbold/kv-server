#!/bin/sh

docker run -d $1 > id.tmp

if [ ! -z "$(docker ps | grep $1)" ]; then
  # container running correctly
  docker stop $(cat id.tmp)
  rm id.tmp
else 
  rm id.tmp
  exit 1
fi
