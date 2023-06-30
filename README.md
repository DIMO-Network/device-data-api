# device-data-api

Serves the device data in all shapes, historical, rolled up, analytics, mostly from Elastic

## Run locally:

`cp settings.sampe.yaml settings.yaml`
`docker compose up -d`
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


## gRPC library

To regenerate gRPC code, if you make changes to the .proto files:

```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pkg/grpc/*.proto
```