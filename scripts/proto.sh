#! /bin/bash

source $(dirname $0)/var_setup.sh

# delete everything in the proto directory
rm -rf $LINGO_PROJECT_PATH/proto/gen/*

docker run --volume "$LINGO_PROJECT_PATH/proto:/workspace" --workdir /workspace bufbuild/buf generate