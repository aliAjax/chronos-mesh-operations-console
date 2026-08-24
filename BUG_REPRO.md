# Bug Reproduction

## Bug

Concurrent NTS operations can accept the same replay nonce more than once, expose mutable key and capability snapshots, and race while replacing and reading a connection handler.

## Trigger

From the repository root, run:

```sh
go test -race ./internal/ratelimit/../nts -run '^TestReplayAcceptPublishesAtomically$' -count=1
go test -race ./internal/ratelimit/../nts -run '^TestRotatedKeyCopyStaysDetached$' -count=1
go test -race ./internal/ratelimit/../nts -run '^TestKEServerHandlerSnapshotConcurrent$' -count=1
go test -race ./internal/ratelimit/../nts -run '^TestNegotiatedCapabilitiesStayDetached$' -count=1
```

## Observed Errors

```text
--- FAIL: TestReplayAcceptPublishesAtomically (0.00s)
    record009_test.go:30: replay accepted 3 concurrent copies
--- FAIL: TestRotatedKeyCopyStaysDetached (0.00s)
    record009_test.go:39: mutating key snapshot changed key ring
WARNING: DATA RACE
--- FAIL: TestKEServerHandlerSnapshotConcurrent (0.00s)
    testing.go:1399: race detected during execution of test
--- FAIL: TestNegotiatedCapabilitiesStayDetached (0.00s)
    record009_test.go:60: capability snapshot aliases internal slices
```
