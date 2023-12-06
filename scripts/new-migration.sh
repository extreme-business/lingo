#! /bin/bash

# get the name of the migration
name="$1"

# check if the name is empty
if [ -z "$name" ]; then
    echo "Please provide a name for the migration, e.g. ./scripts/new-migration.sh create_users_table"
    exit 1
fi

# Get the directory of the script
script_dir="$(dirname "$BASH_SOURCE")"

# Read the .atlas-version file relative to the script's location
atlasVersion=$(cat "${script_dir}/.atlas-version")

echo "creating migration... with atlas version $atlasVersion"

docker run --rm -v ./migrations:/migrations arigaio/atlas:$atlasVersion migrate new $name