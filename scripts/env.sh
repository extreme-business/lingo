#!/bin/bash

# Write the keys to a .env file
echo "LINGO_SIGNING_KEY_REGISTRATION=$(openssl rand -hex 32)" > .env
echo "LINGO_SIGNING_KEY_AUTHENTICATION=$(openssl rand -hex 32)" >> .env

