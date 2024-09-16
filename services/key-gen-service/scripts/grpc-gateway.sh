#!/usr/bin/env bash

set -eux
set -o pipefail

SERVER_ADDR="localhost:$SERVER_PORT"

curl --http2 -iL -w "\n" \
  -H "Accept: application/json" \
  -H "X-Api-Key-Gen-Token: $API_KEY_GEN_TOKEN" \
  -X GET $SERVER_ADDR
