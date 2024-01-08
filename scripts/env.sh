#!/bin/bash

# Generate a 32-byte (256-bit) random key for AES-256
key=$(openssl rand -hex 32)

# Write the key to a .env file
echo "LINGO_AES_256_KEY=$key" > .env

echo "Generated .env with AES-256 key"