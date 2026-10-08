# ORP Backend — OpenResty Plus Go Control Plane

Project identifier: `orp-backend`.

This directory is an independent Go project and can be built or deployed without the frontend workspace.

## Requirements

- Go 1.26.1
- MySQL 8.4 (schema is initialized/upgraded from `internal/store/migrations/`)
- Redis and Kafka are optional for core API availability; configure them when using dashboard snapshot queueing and Kafka log ingestion.

## Run locally

1. Copy `.env.example` to `.env`.
2. Set a strong `OPENRESTY_ADMIN_PASSWORD` and generate `OPENRESTY_DATA_KEY` with `openssl rand -hex 32`.
3. Configure a reachable MySQL instance, then load `.env` and run `go run ./cmd/control-plane`.

The API listens on `OPENRESTY_HTTP_ADDR` (default `:8081`) and exposes `/healthz`. The standalone Docker image is defined by `Dockerfile`. For the local MySQL/Redis/Kafka and three-node integration environment, run `docker compose up -d mysql redis kafka` from this project directory.

## Checks

- `go test ./...`
- `go vet ./...`
- `go run ./cmd/openapi-check`

The frontend maintains the client OpenAPI document and its own consistency check in the separate frontend project.

## Node operations

The current publish adapter supports the local Compose nodes only. External-node Agent enrollment, state ingestion, artifact transfer, and fixed reload operations are documented in the separate `orp-node-agent` project; check its current README for implementation status.
