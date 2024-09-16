#!/usr/bin/env bash

set -eux

if [ ! -z $KEYGEN_CORE_API_TOKEN ]; then 
  KEYGEN_CORE_API_TOKEN=token
fi

IMAGE_NAME=pastebin/key-gen-service
if [ -z $(docker images -f reference=$IMAGE_NAME --format="{{ .ID }}") ]; then
  docker build -t=$IMAGE_NAME --target=release .
fi

CONTAINER_NAME=key-gen-service
if [ ! -z $(docker ps -a -f name=$CONTAINER_NAME --format="{{ .ID }}") ]; then
  docker container rm -f $CONTAINER_NAME &> /dev/null
fi

docker run --rm -d -it \
  --network="host" \
  --name $CONTAINER_NAME \
  -e KEYGEN_GRPC_ADDRESS=:8080 \
  -e KEYGEN_CORE_API_TOKEN=$KEYGEN_CORE_API_TOKEN \
  -e KEYGEN_CORE_KEY_LENGTH=15 \
  -e KEYGEN_CORE_KEY_BUFFER=20 \
  $IMAGE_NAME

