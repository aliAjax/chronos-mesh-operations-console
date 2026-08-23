#!/bin/sh
set -eu
go run ./cmd/chronos-mesh -config configs/config.yaml &
pid=$!
trap 'kill "$pid" 2>/dev/null || true' EXIT
sleep 1
curl -fsS http://127.0.0.1:18085/healthz
curl -fsS http://127.0.0.1:18085/metrics
