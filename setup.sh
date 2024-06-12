#!/bin/bash

cd scripts

# create .env file
echo "Creating .env file..."
source env.sh

# create proto files
echo "Creating proto files..."
source proto.sh

# certs
echo "Creating certs..."
source certs.sh

# lint 
echo "Linting..."
source lint.sh
