#! /bin/bash

DST_DIR=/proto/

echo "Running protoc"

find /proto -name "*.proto"

protoc --proto_path /  \
    --go_out=$DST_DIR --go_opt=module=github.com/dwethmar/lingo/proto \
    --go-grpc_out=$DST_DIR --go-grpc_opt=module=github.com/dwethmar/lingo/proto \
    $(find /proto -name "*.proto")

echo "protoc done"