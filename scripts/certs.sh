#!/bin/bash
source $(dirname $0)/var_setup.sh

mkdir -p $LINGO_PROJECT_PATH/certs

openssl req -x509 -newkey rsa:4096 -sha256 -days 3650 -nodes \
    -keyout $LINGO_PROJECT_PATH/certs/grpc-lingo.key -out $LINGO_PROJECT_PATH/certs/grpc-lingo.crt -subj '/CN=lingo' \
    -extensions san \
    -config <(echo '[req]'; echo 'distinguished_name=req'; echo '[san]'; echo 'subjectAltName=DNS:lingo')

openssl req -x509 -newkey rsa:4096 -sha256 -days 3650 -nodes \
    -keyout $LINGO_PROJECT_PATH/certs/http-lingo.key -out $LINGO_PROJECT_PATH/certs/http-lingo.crt -subj '/CN=lingo' \
    -extensions san \
    -config <(echo '[req]'; echo 'distinguished_name=req'; echo '[san]'; echo 'subjectAltName=DNS:lingo')
