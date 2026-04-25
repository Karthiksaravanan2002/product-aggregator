# Product Aggregator API

A Go backend service that aggregates product data from multiple providers, processes it, and exposes REST APIs.

## Project Structure

```
cmd/aggregator/main.go       → Application entry point
internal/
  config/config.go           → YAML configuration loader
  model/model.go             → Shared data structures
  provider/provider.go       → 3 mock product providers
  service/service.go         → Business logic (search, dedup, sort)
  storage/storage.go         → Thread-safe in-memory store
  handler/handler.go         → HTTP request handlers
  server/server.go           → HTTP server and routing
env/application.yaml         → App configuration
```

## Setup

```bash
go mod tidy
```

## Run

```bash
make run
# or
go run cmd/aggregator/main.go
```

## Test

```bash
make test
# or
go test ./internal/... -v
```

## API

### Search Products
```bash
curl "http://localhost:8080/search?q=mouse"
```
Response:
```json
{
  "query": "mouse",
  "products": [
    {"sku": "SKU-COMMON1", "name": "Gaming Mouse Pad", "price": 15.99, "currency": "USD", "provider": "provider_a"},
    {"sku": "SKU-C1", "name": "Bluetooth Mouse", "price": 19.99, "currency": "USD", "provider": "provider_c"},
    {"sku": "SKU-A1", "name": "Wireless Mouse", "price": 25.99, "currency": "USD", "provider": "provider_a"},
    {"sku": "SKU-B1", "name": "Gaming Mouse", "price": 45.00, "currency": "USD", "provider": "provider_b"}
  ],
  "failed_providers": []
}
```

### Search History
```bash
curl http://localhost:8080/history
curl http://localhost:8080/history/1
```

### Health Check
```bash
curl http://localhost:8080/health
```

## Design Decisions

**Layered Architecture (Handler → Service → Provider/Storage)**
The code is split into layers where each layer does one job. Handlers deal with HTTP, the service holds the business logic, providers fetch data, and storage saves history. No layer reaches into another's responsibilities.

**Provider Interface (Strategy + Registry Pattern)**
All providers implement the same interface, so the service doesn't care which provider it's talking to. New providers can be added without touching existing code. The registry collects all enabled providers and the service fans out to all of them concurrently.

**Standard Library Router (Go 1.22+ ServeMux)**
No third-party router needed. Go 1.22 added support for method-based routing and path parameters to the standard library, which covers everything this project needs.

**In-Memory Storage with sync.RWMutex**
Search history is stored in memory using a read-write mutex. Multiple goroutines can read history at the same time, but writes are serialized to prevent race conditions.

**Buffered Channel for Concurrent Results**
Each provider runs in its own goroutine. Results are sent into a buffered channel so no goroutine blocks waiting for another. Each provider also gets its own timeout so a slow provider doesn't hold up the rest.

## Tradeoffs

| Choice | Benefit | Tradeoff |
|--------|---------|----------|
| In-memory storage | Fast, no setup | Data lost on restart |
| Standard library only (no chi/gin) | No extra HTTP dependencies needed | Routing syntax is less convenient than third-party routers |
| Per-provider context timeout | A slow provider won't delay the others | Each provider creates its own goroutine and context |
| First-seen dedup strategy | Simple and predictable | If the same product appears in multiple providers, we keep the first one found instead of comparing prices |
| YAML config | Human readable, structured and decoupled | Adds one external dependency |
