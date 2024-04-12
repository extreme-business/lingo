#! /bin/bash

# delete everything in the proto directory
rm -rf protogen/*

source $(dirname $0)/var_setup.sh

docker run --volume "$LINGO_PROJECT_PATH/proto:/workspace" --workdir /workspace bufbuild/buf generate