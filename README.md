# device-data-api

API & worker that serves the device data in all shapes, historical, rolled up, analytics, mostly from Elastic. 
Also has it's own database that keeps status snapshot of latest data state of devices. 

For an overview of the project, see the [DIMO technical documentation site.](https://docs.dimo.zone/docs/overview/intro)

## Run locally:

`cp settings.sampe.yaml settings.yaml`

### Run dependencies:
`docker compose up -d`
You will also need local psql instance running, eg. via brew services.

If you are working on the data download endpoint, copy and paste the following into the terminal (this creates a bucket on your local s3 instance):

```
aws s3api create-bucket \
    --bucket test-bucket \
    --region us-east-2 \
    --create-bucket-configuration LocationConstraint=us-east-2 \
    --endpoint-url http://localhost:4566
```

Check settings and make sure make sense for your setup. You may need to be running docker compose elastic etc from devices-api.
`go run ./cmd/device-data-api/`

Run migrations:
`go run ./cmd/device-data-api migrate`

### Microservice dependencies

This project depends on other microservices: 
- devices-api 
- device-definitions-api
- users-api
- dex-roles-rights

To run it locally best bet is https://github.com/DIMO-Network/cluster-local

### Authenticating

One of the variables set in `settings.yaml` is `JWT_KEY_SET_URL`. By default this is set to `http://127.0.0.1:5556/dex/keys`. To make use of this, clone the DIMO Dex fork:
```sh
git clone git@github.com:DIMO-Network/dex.git
cd dex
make build examples
./bin/dex serve examples/config-dev.yaml
```
This will start up the Dex identity server on port 5556. Next, start up the example interface by running
```sh
./bin/example-app
```
You can reach this on port 5555. The "Log in with Example" option is probably the easiest. This will give you an ID token you can provide to the [API](#api).


## Swagger generation

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/device-data-api/main.go --parseDependency --parseInternal --generatedTime true
# optionally add `--parseDepth 2` if have issues
```

## NATS

The `/user/device-data/:userDeviceID/export/json/email` endpoint users NATS to create a queue of requests so that the endpoint can prompty return before processing the user request.
Suggested Stream Name and Subject can be found in `settings.sample.yaml`.
If messages are not acknowledged within 5 minutes, they will be resent (this value is also set in settings and can be increased or decreased as needed)

`brew install nats-streaming-server` although this does not seem to work locally with our setup as it always returns No Responders error.

## gRPC library

To regenerate gRPC code, if you make changes to the .proto files:

```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pkg/grpc/*.proto
```

## Linting

`brew install golangci-lint`

`golangci-lint run`

This should use the settings from `.golangci.yml`, which you can override.

## Database ORM

This is using [sqlboiler](https://github.com/volatiletech/sqlboiler). The ORM models are code generated. If the db changes,
you must update the models.

Make sure you have sqlboiler installed:
```bash
go install github.com/volatiletech/sqlboiler/v4@latest
go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest
```

To generate the models:
```bash
sqlboiler psql --no-tests --wipe
```
*Make sure you're running the docker image (ie. docker compose up)*

If you get a command not found error with sqlboiler, make sure your go install is correct.
[Instructions here](https://jimkang.medium.com/install-go-on-mac-with-homebrew-5fa421fc55f5)

## Migrations

To install goose in GO:
```bash
$ go get github.com/pressly/goose/v3/cmd/goose@v3.5.3
export GOOSE_DRIVER=postgres
```

To install goose CLI:
```bash
$ go install github.com/pressly/goose/v3/cmd/goose
export GOOSE_DRIVER=postgres
```

Add a migrations:
`$ goose -dir migrations create <migration_name> sql`

Migrate DB to latest:
`$ go run ./cmd/device-data-api migrate`