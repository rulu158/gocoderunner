#!/bin/sh
docker image prune -f
docker rmi --force $(docker images -f dangling=true -q)
