#! /bin/bash

source $(dirname $0)/var_setup.sh

# get the name of the app
app="$1"

if [ -z "$app" ]; then
    echo "Please provide the name of the app, e.g. ./scripts/hash-migration.sh <app>"
    exit 1
fi

# Read the .atlas-version file relative to the script's location
atlasVersion=$(cat "$LINGO_PROJECT_PATH/scripts/.atlas-version")

docker run --rm -v $LINGO_PROJECT_PATH/apps/$app/migrations:/migrations arigaio/atlas:$atlasVersion migrate hash