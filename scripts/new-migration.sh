#! /bin/bash

source $(dirname $0)/var_setup.sh

# get the name of the migration
name="$1"

# check if the name is empty
if [ -z "$name" ]; then
    echo "Please provide a name for the migration, e.g. ./scripts/new-migration.sh create_users_table"
    exit 1
fi

# Read the .atlas-version file relative to the script's location
atlasVersion=$(cat "$LINGO_PROJECT_PATH/scripts/.atlas-version")

docker run --rm -v $LINGO_PROJECT_PATH/migrations:/migrations arigaio/atlas:$atlasVersion migrate new $name