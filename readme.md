# go-pulse

A lightweight Go service that monitors URLs using a concurrent worker pool — with zero external dependencies.

## Overview

`go-pulse` submits URLs to a pool of worker goroutines that process them concurrently. It also exposes a simple HTTP server with a `/health` endpoint, and shuts down gracefully on `SIGINT` / `SIGTERM`.

## Project Structure

```
go-pulse/
├── cmd/
│   └── server/
│       └── main.go          # Entry point — wires config, pool, and HTTP server
├── internal/
│   ├── config/
│   │   └── config.go        # Env-based configuration
│   └── worker/
│       ├── job.go           # Job type definition
│       └── pool.go          # Worker pool (goroutines + channel)
└── go.mod                   # module go-pulse, Go 1.22
```

## How It Works

1. **Config** is loaded from environment variables (with sensible defaults).
2. A **worker pool** is created with `WorkerCount` goroutines, each reading from a shared `Jobs` channel.
3. **Jobs** (URLs) are submitted to the pool. Each worker logs the URL it picks up.
4. An **HTTP server** starts on the configured port, serving a `/health` endpoint that returns `200 OK`.
5. On receiving `SIGINT` or `SIGTERM`, the server shuts down gracefully (5-second timeout) and the jobs channel is closed.

## Configuration

All options are set via environment variables:

| Variable          | Default | Description                          |
|-------------------|---------|--------------------------------------|
| `PORT`            | `8080`  | HTTP server port                     |
| `WORKER_COUNT`    | `50`    | Number of concurrent worker goroutines |
| `REQUEST_TIMEOUT` | `5`     | Request timeout in seconds (defined in config, available for future use) |

## Requirements

- Go 1.22+
- No external dependencies

## Getting Started

```bash
# Clone
git clone https://github.com/never-code/go-pulse.git
cd go-pulse

# Run with defaults
go run ./cmd/server

# Run with custom config
PORT=9090 WORKER_COUNT=10 go run ./cmd/server

# Build
go build -o go-pulse ./cmd/server
./go-pulse
```

## API

### `GET /health`

Returns a simple liveness check.

```
HTTP/1.1 200 OK

OK
```

## Code Walkthrough

### `cmd/server/main.go`

The entry point. Loads config, starts the pool, submits hardcoded test URLs, starts the HTTP server in a goroutine, then blocks on an OS signal channel for graceful shutdown.

```go
pool := worker.NewPool(cfg.WorkerCount)
pool.Start()

pool.Submit(worker.Job{URL: "https://google.com"})
// ...
```

### `internal/config/config.go`

Reads `PORT`, `WORKER_COUNT`, and `REQUEST_TIMEOUT` from the environment. Falls back to defaults if unset or unparseable.

### `internal/worker/job.go`

Minimal type — a `Job` is just a struct with a `URL` string. Designed to be extended with additional fields (e.g. expected status code, timeout, headers).

### `internal/worker/pool.go`

A `Pool` holds a `WorkerCount` and an unbuffered `Jobs` channel. `Start()` spawns `WorkerCount` goroutines, each ranging over the channel. `Submit()` sends a job onto the channel (blocks until a worker is free).

```go
func (p *Pool) worker(id int) {
    for job := range p.Jobs {
        log.Printf("Worker %d processing URL: %s", id, job.URL)
    }
}
```

> **Note:** Workers currently only log URLs — actual HTTP request logic using `REQUEST_TIMEOUT` is not yet implemented and is the natural next step.

## Potential Next Steps

- Implement actual HTTP health checks inside `worker()` using `http.Client` with the configured timeout
- Return and aggregate results (status code, latency, errors)
- Buffer the `Jobs` channel to avoid blocking on `Submit`
- Add structured logging
- Wire up the `/health` endpoint to reflect pool/check status

## Contributing

1. Fork the repository.
2. Create a feature branch off `develop`.
3. Commit your changes and open a pull request against `develop`.
