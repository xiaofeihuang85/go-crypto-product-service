# go-crypto-product-service

`go-crypto-product-service` is a small Go backend service that fetches product
data from Coinbase public APIs, transforms it into a service-owned response
model, and exposes it through a simple HTTP API.

The project is intentionally scoped as a simplified product cache service. The
goal is to demonstrate clean service structure, clear dependency boundaries,
external API integration, caching behavior, and local operability without
overbuilding a production-grade distributed cache system.

## Current Status

Phase 8 is focused on the final quality pass:

- Runnable Go executable in `cmd/server/main.go`
- Standard-library HTTP server with route registration
- Config-based port selection through the `PORT` environment variable
- JSON responses from `GET /`, `GET /health`, and `GET /products/{product_id}`
- Thin Coinbase client for fetching a single product from the public market
  products endpoint
- Transforms Coinbase data into a cleaner market view with derived fields such as `market_pair`, `is_trading_enabled`, and `source`
- Read-through Redis cache with a short TTL for product lookups
- Cache is best-effort: Redis failures fall back to live Coinbase data instead of
  failing the request
- Polished primary endpoint response with `cache_status` and `retrieved_at`
- Structured API errors with consistent error codes and request paths
- Basic server-side and upstream HTTP timeouts for safer local operation
- Dockerfile and Docker Compose setup for running the service with Redis locally
- Focused tests around service behavior and handler error mapping

The remaining work in this phase is final verification and submission polish.

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

```text
client -> api -> service -> redis cache
                      \
                       -> coinbase public api
```

## Request Flow

1. A client calls `GET /products/{product_id}`.
2. The handler validates the path and delegates to the service layer.
3. The service checks Redis using a read-through cache pattern.
4. On cache miss, the service calls Coinbase's public market products endpoint.
5. The Coinbase response is transformed into a service-owned product response.
6. The transformed response is cached in Redis with a short TTL and returned to the client.

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
curl http://localhost:8080/products/BTC-USD
```

Optional cache-related environment variables:

```bash
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
CACHE_TTL=60s
```

## Run With Docker Compose

Start the app and Redis together:

```bash
docker compose -f deployments/docker-compose.yml up --build
```

Stop the stack:

```bash
docker compose -f deployments/docker-compose.yml down
```

Once the stack is running, the API is available at:

```bash
http://localhost:8080
```

Example product response:

```json
{
  "product_id": "BTC-USD",
  "market_pair": "BTC/USD",
  "product_name": "Bitcoin",
  "base_currency": "BTC",
  "quote_currency": "USD",
  "status": "online",
  "is_trading_enabled": true,
  "price": "70813.85",
  "price_change_24h": "-2.66118725176928",
  "cache_status": "miss",
  "retrieved_at": "2026-04-12T17:22:46Z",
  "source": "coinbase"
}
```

The product response is intentionally service-owned rather than a raw pass-through
of the upstream Coinbase schema. This keeps the API easier to reason about and
helps stabilize the contract around the cache layer.

The cache strategy is intentionally simple: `GET /products/{product_id}` uses a
read-through flow with a short TTL and a cache key shaped like `product:{product_id}`.
More advanced invalidation, refresh behavior, and multi-node cache coordination are
future-scaling concerns that can be described and discussed without fully implementing
them in this version of the project.

Example cache flow:

1. First request for `GET /products/BTC-USD` returns `cache_status: "miss"` and stores the transformed response in Redis.
2. A repeated request for the same product returns `cache_status: "hit"` from Redis while keeping `source: "coinbase"` to show the origin of the data.

You can test that flow locally with:

```bash
curl http://localhost:8080/products/BTC-USD
```

Example error response:

```json
{
  "code": "product_not_found",
  "error": "product not found: XXX-XXX",
  "path": "/products/XXX-XXX"
}
```

## Testing

Run the focused unit tests with:

```bash
go test ./...
```

The current test suite focuses on:

- product service cache hit and miss behavior
- product transformation and normalization
- handler status-code and error-response mapping

## Assumptions And Tradeoffs

- Coinbase's public market products endpoint is used as the upstream data source.
- Redis is treated as a best-effort cache. Cache failures should degrade to live Coinbase fetches rather than fail the request immediately.
- Cache invalidation is intentionally simplified to a short TTL.
- The response contract is service-owned rather than a raw Coinbase pass-through so the cache layer can remain stable even if the upstream schema changes.

## Future Scaling Notes

Areas intentionally left out of this implementation but worth discussing in an interview:

- retry and backoff strategy for upstream failures
- active cache invalidation and background refresh
- request collapsing to reduce duplicate upstream calls during cache misses
- metrics, tracing, and structured logging
- multi-node cache coordination and broader production deployment concerns

## Scope Notes

This project is meant to represent a clean, well-explained subset of a cache
service rather than a full production simulation. Areas such as advanced retry
policy, cache invalidation strategy, and future scaling behavior may be partly
implemented and partly documented as design tradeoffs.

## AI Usage

AI assistance was used selectively for scaffolding, editing support, unit-test
drafting, and iteration on documentation and implementation structure.

It primarily helped speed up repetitive setup work, compare implementation
options, and tighten the codebase during later quality passes. Final architecture,
scope, API design, caching behavior, and testing choices were reviewed and
directed manually.
