# go-crypto-product-service

`go-crypto-product-service` is a small backend service in Go for fetching crypto
product data from Coinbase public APIs, applying domain-specific transformations,
and exposing the result through a simple HTTP API.

The project is intentionally scoped as a simplified product cache service. The
goal is to demonstrate clean service structure, clear dependency boundaries,
external API integration, caching considerations, and local operability without
overbuilding a production-grade distributed cache system.

## Current Status

Phase 2 is now focused on the HTTP bootstrap milestone:

- Runnable Go executable in `cmd/server/main.go`
- Standard-library HTTP server with route registration
- Config-based port selection through the `PORT` environment variable
- JSON responses from `GET /` and `GET /health`
- Basic server-side timeouts for safer local operation

Coinbase integration, Redis cache, and Docker Compose setup will be added in
later phases.

## Planned Architecture

- `cmd/server`: application entrypoint and process bootstrap
- `internal/api`: HTTP handlers and route registration
- `internal/service`: domain orchestration and transformation logic
- `internal/client`: Coinbase API client
- `internal/store`: cache access layer, initially Redis-backed
- `internal/model`: request, response, and domain models
- `internal/config`: environment-based configuration loading
- `deployments`: local deployment assets such as Docker Compose
- `scripts`: helper scripts for local development

## Milestones

1. Phase 1: Foundation
   Runnable executable, repo structure, initial documentation.
2. Phase 2: HTTP bootstrap
   Minimal HTTP server with route wiring and config-driven startup.
3. Phase 3: Coinbase integration
   Upstream client with timeout handling and parsed product data.
4. Phase 4: Domain transformation
   Internal models and structured API response shaping.
5. Phase 5: Redis caching
   Simple read-through cache with TTL.
6. Phase 6: Main endpoint
   Endpoint that checks cache, fetches upstream data on miss, transforms it, and
   returns structured JSON.
7. Phase 7: Local operations
   Docker Compose and local run instructions.
8. Phase 8: Quality pass
   Focused tests, improved reliability, documented tradeoffs, and final README.

## Running Locally

Start the service:

```bash
go run ./cmd/server
```

Optionally override the default port if the default port is already in use:

```bash
PORT=9090 go run ./cmd/server
```

Example requests:

```bash
curl http://localhost:8080/
curl http://localhost:8080/health
```

## Scope Notes

This project is meant to represent a clean, well-explained subset of a cache
service rather than a full production simulation. Areas such as advanced retry
policy, cache invalidation strategy, and future scaling behavior may be partly
implemented and partly documented as design tradeoffs.

## AI Usage

AI assistance is being used selectively for scaffolding, editing support, and
iteration on documentation and implementation structure.
