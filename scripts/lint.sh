#! /bin/bash

# build the docker image if it doesn't exist
if [ -z "$(docker images -q lingo-lint 2> /dev/null)" ]; then
    docker build --target lint -t lingo-lint . 
fi

docker run --rm -it lingo-lint run /src/lingo/...
