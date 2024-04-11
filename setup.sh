#!/bin/bash

# create certificates
echo "Creating certificates..."
source scripts/certs.sh

# create .env file
echo "Creating .env file..."
source scripts/env.sh

# create proto files
echo "Creating proto files..."
source scripts/proto.sh
source ./scripts/proto-buf-mod-update.sh

# lint 
echo "Linting..."
source scripts/lint.sh