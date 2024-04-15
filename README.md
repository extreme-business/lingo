# Lingo
Lingo is a web app. I don't know what it will do, but for now the goals are:
- single binary, many commands with different functionality, like: serving grpc, http and jobs.
- use proto as the main way to design api's.
- use proto to generate https gateways for the grpc servers.
- add all https services to swagger ui so they can be easily tested.

# Setup
The goal with the setup is that minimal tools are needed to run the project.

## prerequisites 
- install [docker](https://docs.docker.com/get-docker/).
- install openssl, version 3 (for generating certificates and keys).
  - install with [brew](https://formulae.brew.sh/formula/openssl@3.0).
- install [atlas](https://atlasgo.io/).

## local environment
- run `setup.sh`. you should be able to run the setup as many times as you want.
  - Specified deps in the `buf.yaml` need to be covered in your `buf.lock` file. If you get an error, run `scripts/proto-buf-mod-update.sh` to generate the `buf.lock` file.
  - resulting generated files are in the `protogen` folder.
  - check changes with the [buf](https://buf.build/) linter: `scripts/proto-lint.sh`.
- run `scripts/certs.sh` (this is not part of the setup, because you may want to choose to add the cert to you computer so you can use swagger-ui without getting self signing ssl errors).
    - trust self signed certificates for [mac](https://tosbourn.com/getting-os-x-to-trust-self-signed-ssl-certificates/).

# Run locally
- setup should have completed.
- run `docker-compose up`.
- To view the [open-api](https://en.wikipedia.org/wiki/Open_API) specs for various services, open [localhost:8090](localhost:8090) in the browser.
- List the services: `docker compose config | yq '.services[]|key + " | " + .image'`

# Develop
The goal is to have a good developer experience. That means that the developer should have to read minimal setup guides en be up and running as fast as possible.

## Debug
- See the [launch.json](.vscode/launch.json) to debug with vscode

## Linting
- to lint, run `./scripts/lint.sh`.

## Testing
Test that require a database can use [testcontainers](https://golang.testcontainers.org/modules/postgres/).
- to run all tests, run `go run test ./...`.

## Database migrations
Migrations are managed by [atlas](https://atlasgo.io).
- to create a new migration, run `./scripts/new-migration.sh <app> <name of migration>`.
- after you have written your migration, run `./scripts/hash-migration.sh <app>`.

## Proto
Proto files are generated to go server and client code with [buf](https://buf.build). 
Buf also generates [Openapiv2](https://swagger.io/specification/v2/) files in yaml and json for every service. 
- build the proto files with `./scripts/proto.sh`. 
- resulting go files are in `proto/gen/go/**/*.go`.
- resulting Openapiv2 files are in `proto/gen/openapiv2/**/*{json,yaml}`
- lint the proto files with `./scipts/proto-lint.sh`.

# Guidelines.

## API's
- follow: https://cloud.google.com/apis/design
