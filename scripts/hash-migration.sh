#! /bin/bash

source $(dirname $0)/var_setup.sh

# Read the .atlas-version file relative to the script's location
atlasVersion=$(cat "$LINGO_PROJECT_PATH/scripts/.atlas-version")

docker run --rm -v $LINGO_PROJECT_PATH/migrations:/migrations arigaio/atlas:$atlasVersion migrate hash