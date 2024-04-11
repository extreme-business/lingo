#! /bin/bash

source $(dirname $0)/var_setup.sh

# build the docker image if it doesn't exist
if [ -z "$(docker images -q lingo-lint 2> /dev/null)" ]; then
    docker build --target lint -t lingo-lint . 
fi

docker run --rm -v "$LINGO_PROJECT_PATH:/workspace" -w /workspace lingo-lint run /workspace/...