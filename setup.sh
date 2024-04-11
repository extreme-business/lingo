#!/bin/bash

cd scripts

# create certificates
echo "Creating certificates..."
source certs.sh

# create .env file
echo "Creating .env file..."
source env.sh

# create proto files
echo "Creating proto files..."
source proto.sh
source proto-buf-mod-update.sh

# lint 
echo "Linting..."
source lint.sh