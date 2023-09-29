#!/bin/sh
#docker rm $(docker stop $(docker ps -a -q --filter="name=$1" --format="{{.ID}}"))
docker container stop $1
docker container rm $1