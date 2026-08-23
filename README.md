# chronos-mesh

Chronos Mesh is a pure-Go NTPv4 logical time service with source clustering, holdover, NTS cookie primitives, control API, metrics and graceful shutdown. The UDP data path never writes to a database. PostgreSQL migrations are supplied for control-plane persistence.

## Run

```sh
go run ./cmd/chronos-mesh -config configs/config.yaml
curl http://127.0.0.1:8080/healthz
curl http://127.0.0.1:8080/readyz
```

The default sample sources point at local UDP simulators. When no simulator is available, `/healthz` remains available and `/readyz` accurately reports unsynchronized state. Use `examples/upstream.go` to run deterministic UDP responders.

## Verification

Run `gofmt -w .`, `go vet ./...`, `go test -race ./...`, and `go build ./...`. The smoke script starts the service, checks health and metrics, and terminates it through a trap.

## API

The REST contract is in `api/chronos.yaml`; NTS-KE record details are documented in `api/nts-ke.md`. Endpoints expose logical time, source diagnostics, leap metadata and Prometheus metrics.

## Operations page

The static operations page is in `web/index.html`. Serve it from the same origin as the control API to inspect health, logical time, source samples and leap metadata in a browser.

## Security boundaries

System clock adjustment is intentionally absent. NTS-KE certificate termination is kept behind an injectable adapter; cookie encryption and replay-window primitives are implemented in `internal/nts`. Production deployments must provide TLS certificates, API authentication, upstream allowlists and PostgreSQL credentials through a secret manager.
