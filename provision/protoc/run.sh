#! /bin/bash

DST_DIR=/proto/
DST_MODULE=github.com/dwethmar/lingo/proto

echo "Running protoc"

find /proto -name "*.proto"

protoc --proto_path /  \
    --go_out=$DST_DIR --go_opt=module=$DST_MODULE \
    --go-grpc_out=$DST_DIR --go-grpc_opt=module=$DST_MODULE \
    $(find /proto -name "*.proto")

echo "protoc done"