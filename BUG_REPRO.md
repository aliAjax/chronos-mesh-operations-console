# Bug Reproduction

## Bug

Selection and quality helpers reuse caller-owned slice storage. Sorting or filtering one result can reorder an input list, overwrite a retained decision, or change a weight snapshot returned by an earlier call.

## Trigger

From the repository root, run:

```sh
go test ./internal/selection -run '^TestSelectPreservesCallerSamplesS006$' -count=1
go test ./internal/selection -run '^TestWeightsRemainStableAcrossCallsS006$' -count=1
go test ./internal/selection -run '^TestAcceptedDoesNotRewriteDecisionS006$' -count=1
go test ./internal/selection/../source -run '^TestMedianOffsetPreservesSampleOrderQ006$' -count=1
```

## Observed Errors

```text
--- FAIL: TestSelectPreservesCallerSamplesS006 (0.00s)
    record006_test.go:15: caller samples rewritten
--- FAIL: TestWeightsRemainStableAcrossCallsS006 (0.00s)
    record006_test.go:23: weight snapshot reused
--- FAIL: TestAcceptedDoesNotRewriteDecisionS006 (0.00s)
    record006_test.go:32: decision rewritten
--- FAIL: TestMedianOffsetPreservesSampleOrderQ006 (0.00s)
    record006_test.go:16: median rewrote input
```
