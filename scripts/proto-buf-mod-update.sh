#! /bin/bash

source $(dirname $0)/var_setup.sh

docker run --volume "$LINGO_PROJECT_PATH:/workspace" --workdir /workspace bufbuild/buf mod update
