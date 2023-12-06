#!/bin/bash

# Get the directory of the script that sourced this one
script_dir="$(dirname "$BASH_SOURCE")"

# Normalize the path to ensure it's absolute
script_dir="$(realpath "$script_dir")"

# Export the variables so they are available in scripts that source this file
export LINGO_PROJECT_PATH="$(realpath "$script_dir/../")"