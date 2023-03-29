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
