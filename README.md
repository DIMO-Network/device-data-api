# device-data-api
Serves the device data in all shapes, historical, rolled up, analytics, mostly from Elastic

## Run locally:
`cp settings.sampe.yaml settings.yaml`
Check settings and make sure make sense for your setup. You may need to be running docker compose elastic etc from devices-api.
`go run ./cmd/device-data-api/`

## Swagger generation

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/device-data-api/main.go --parseDependency --parseInternal --generatedTime true
# optionally add `--parseDepth 2` if have issues
```