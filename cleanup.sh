#!/bin/sh

docker container prune -f

docker image prune -f
