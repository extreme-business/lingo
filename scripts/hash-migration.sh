#! /bin/bash

# Get the directory of the script
script_dir="$(dirname "$BASH_SOURCE")"

# Read the .atlas-version file relative to the script's location
atlasVersion=$(cat "${script_dir}/.atlas-version")

echo "Hashing migrations... with atlas version $atlasVersion"

docker run --rm -v ./migrations:/migrations arigaio/atlas:$atlasVersion migrate hash