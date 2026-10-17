# E10-INV-120A — one-shot local JSONL to authorized triage

## Functional integration

A local analyst previously needed two commands and a sensitive intermediate full owner JSON snapshot to turn normalized events into a Case report. The new explicit CLI flag `-triage-import` requires `-import-jsonl`, `-into-case` and `-snapshot` and runs the full path in **one process**:

`operator JSONL + owner JSON -> validated in-memory import -> existing authenticated Case discovery and paginated Timeline/source ACL -> aggregate-only Case triage JSON`

Example (synthetic lab inputs only):

```sh
go run ./cmd/case-timeline-demo -principal=analyst-a -snapshot=/absolute/owner.json -import-jsonl=/absolute/events.jsonl -into-case=case-a -triage-import > /private/triage.json
```

This does **not** output the merged owner graph (which may contain other tenant identities and restricted source previews). The original snapshot and JSONL are never mutated. Output has the existing Case triage counts, kind distribution, known time bounds, uncertainty and top-five **authorized** source references; plus the count of new imported entries. No raw messages, source previews, bearer, or unfiltered Timeline Entry IDs. Existing offline import-as-snapshot and ordinary read-only loopback browser routes are unchanged. The summary still contains potentially sensitive authorized Case/source identifiers and must be protected.

## Reuse-first decision

**Chosen**: existing first-party `loadDemoSnapshot`, `readLocalEventFile`, `importLocalEvents`, `newDemoHandlerWithData`, `collectDemoCaseTriage`, Go `encoding/json`, same in-memory ACL chain. No second provider, permission model, correlator, index, ORM, database or paid API. Alternative (persisting an intermediate merged graph, then reloading it) would increase disclosure and I/O risks without a functional benefit.

## Explicit limits / acceptance

The operator JSONL uses strictly normalized local synthetic rows, not untrusted remote SIEM streams. Host CLI arguments are not HTTP parameters; no network listener. Owner-supplied principals/permissions are valid only **inside the local fixture**, not deployed IAM. No durable case DB, external telemetry agent, transactional snapshot, signed chain of custody, production reliability or final UI. No mutation of Shared Timeline Entry or resolution of OPEN-013.

Tests exercise real disk files, actual existing in-process authenticated HTTP Case handler traversal, authorized Case triage, hidden/cross-tenant non-disclosure, malformed JSONL rejection with zero bytes, cancellation and explicit-only CLI modes. Exact-SHA CI and independent risk-based review determine whether this increment is verified; this note is not a CI verdict.
