#! /bin/bash

set -e
set -u

rm -rf ./proto/v1/**/*.go

docker build -t lingo-protobuf -f ./provision/protoc/Dockerfile .

docker run --rm -v $(pwd)/proto:/proto lingo-protobuf