#! /bin/bash

# delete everything in the proto directory
rm -rf protogen/*

docker run --volume "$(pwd):/workspace" --workdir /workspace bufbuild/buf generate