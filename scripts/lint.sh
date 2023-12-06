#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status.
set -u  # Treat unset variables as an error.

SCRIPT_DIR=$(dirname "$0")
source "$SCRIPT_DIR/var_setup.sh"

# Ensure LINGO_PROJECT_PATH is set
if [ -z "${LINGO_PROJECT_PATH:-}" ]; then
    echo "Error: LINGO_PROJECT_PATH environment variable is not set."
    exit 1
fi

DOCKERFILE="$LINGO_PROJECT_PATH/Dockerfile"

# Check if the Dockerfile exists
if [ ! -f "$DOCKERFILE" ]; then
    echo "Error: Dockerfile not found in $LINGO_PROJECT_PATH"
    exit 1
fi

# Get the last modified time of the Dockerfile
if [[ "$OSTYPE" == "darwin"* ]]; then
    DOCKERFILE_MTIME=$(stat -f "%m" "$DOCKERFILE")
else
    DOCKERFILE_MTIME=$(stat -c %Y "$DOCKERFILE")
fi

# Get the creation time of the Docker image (if it exists)
IMAGE_ID=$(docker images -q lingo-lint 2> /dev/null)

REBUILD_IMAGE=false

if [ -z "$IMAGE_ID" ]; then
    echo "Docker image lingo-lint not found. Building image..."
    REBUILD_IMAGE=true
else
    # Get the image creation time
    IMAGE_CREATED=$(docker inspect -f '{{.Created}}' "$IMAGE_ID")

    if [[ "$OSTYPE" == "darwin"* ]]; then
        IMAGE_CREATED_TIMESTAMP=$(date -jf "%Y-%m-%dT%H:%M:%S" "$(echo $IMAGE_CREATED | cut -d. -f1)" +%s)
    else
        IMAGE_CREATED_TIMESTAMP=$(date -d "$IMAGE_CREATED" +%s)
    fi
    
    # Compare timestamps
    if [ "$DOCKERFILE_MTIME" -gt "$IMAGE_CREATED_TIMESTAMP" ]; then
        echo "Dockerfile is newer than the Docker image. Rebuilding image..."
        REBUILD_IMAGE=true
    else
        echo "Docker image is up to date."
    fi
fi

if [ "$REBUILD_IMAGE" = true ]; then
    docker build --target lint -t lingo-lint .
    echo "Docker image lingo-lint built successfully."
fi

# Run the linter in the docker container
echo "Running lingo-lint in Docker container..."
docker run --rm -v "$LINGO_PROJECT_PATH:/workspace" -w /workspace lingo-lint run --fix /workspace/...
echo "Lingo-lint run completed."