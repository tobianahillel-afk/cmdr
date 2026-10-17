# E10-INV-120A — offline cross-Case five-minute event pivot

## Functional capability

After the operator supplies a bounded local owner snapshot, the analyst can
request a **single precise UTC five-minute window** across two independently
authorized Cases, with a single bounded, sensitive local JSON chronology:

```sh
go run ./cmd/case-timeline-demo \
  -principal=analyst-a -snapshot=/absolute/path/owner.json \
  -compare-case=case-a -with-case=case-c \
  -compare-events -pivot-window=2026-10-01T10:00:00Z
```

The code reuses the existing authenticated Case discovery, paginated
Case/Timeline reads, typed-source ACLs, and chronological merge. It visits
*all* authorized pages for each Case and produces **no output** unless each
Case contributes at least one certain, permitted event within the exact
half-open interval [start UTC, start+5 minutes). Unknown-clock events are
not silently promoted to certain chronology. The receipt retains the owning
Case, event ID, classification, original timestamp, source reference and
provenance, plus the already-authorized per-Case aggregate summaries. It is
explicitly **sensitive local stdout**, not an automatic saved file.

## Reuse-first technology decision

Existing first-party `collectAuthorizedCrossCaseSequence` plus the Go
standard library `time.Parse`, `json.Encoder` and immutable in-memory
filter are sufficient for a two-Case bounded 10,000-entry maximum. An
external index, search engine, correlation engine, new route, new IAM layer
or persistence schema would be unjustified. The caller chooses a canonical
RFC3339 UTC five-minute boundary; this is an explicit analyst selection, not
an inferred timestamp window or proof of causality.

The memory owner provider currently uses a **globally unique**
`SourceOwnerIdentity.Reference` within its dataset. A cross-Case join on
an identical source reference would therefore be misleading without a
separately approved owner/identity contract. This increment deliberately
joins by independent UTC time window instead; it does NOT assign common
source identity or source equivalence to two Cases. OPEN-013 is unchanged.

## Evidence and security scope

- Real on-disk two-Case snapshot with 207 permitted events in the first
  Case across two pages; source-owner authorized single event in the second
  Case. Test expects the exact 206+1 window subset, chronological ordering,
  observed/inferred preservation and provenance from the disk-only event.
- Denied/foreign Case or source, zero overlap, unknown second clock,
  malformed/noncanonical window, conflicting CLI modes, canceled work and
  failing output writers must not emit partial receipts.
- Existing `case-timeline-demo` remains loopback-only with a **synthetic**
  principal, host-supplied 8 MiB snapshot and ephemeral bearer. There is
  no production IAM, live data ingestion, durable Case owner, provider audit
  log, transactional snapshot consistency or signed chain of custody.
- This is a local investigator workflow, not full production CAP-INV-110.
  `E10-INV-120A-SNAPSHOT` remains READY pending exact-commit CI and
  a single independent, risk-based L2 for this functional milestone.

## Verification

Check on immutable published SHA: `go test -race ./...`,
`go vet ./...`, strict `gofmt`, workflow `engineering-kernel`
and `investigation-timeline-runtime`. Historic green runs do not prove
the new pivot; exact run IDs and independent findings must be recorded
after they are actually observed.
