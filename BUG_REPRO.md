# Bug Reproduction

## Bug

The UDP listener shutdown path races with the active read loop and can wait before releasing the blocked read. Empty limiter pruning and denied requests can retain the limiter lock. Reading the rate-limit statistics can also clear the counters.

## Trigger

Run the record's four targeted race-enabled verification commands against the red branch. They exercise UDP shutdown, pruning an empty limiter, a denied limiter request, and two consecutive statistics snapshots.

## Observed errors

```text
WARNING: DATA RACE
close waited before releasing UDP read
empty prune retained lock
denied path retained lock
destructive snapshot: [0 0]
```

Each targeted command exits with status 1 on the red snapshot.
