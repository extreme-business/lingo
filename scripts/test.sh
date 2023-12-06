#! /bin/bash

source $(dirname $0)/var_setup.sh

# build the docker image if it doesn't exist
if [ -z "$(docker images -q lingo-test 2> /dev/null)" ]; then
    docker build --target test -t lingo-test . 
fi

docker run --rm \
    -v "$LINGO_PROJECT_PATH:/src/lingo" \
   lingo-test -short ./...