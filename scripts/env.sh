#!/bin/bash
source $(dirname $0)/var_setup.sh

# Write the keys to a .env file
echo "LINGO_SIGNING_KEY_REGISTRATION=$(openssl rand -hex 32)" > $LINGO_PROJECT_PATH/.env
echo "LINGO_SIGNING_KEY_AUTHENTICATION=$(openssl rand -hex 32)" >> $LINGO_PROJECT_PATH/.env

# system user
echo "LINGO_SYSTEM_USER_PASSWORD=$(openssl rand -hex 32)" >> $LINGO_PROJECT_PATH/.env
