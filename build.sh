#!/bin/bash

docker build --target prod -t lingo .

# gcloud builds submit --region=europe-west4 --config cloudbuild.yaml