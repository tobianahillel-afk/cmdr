# E10-INV-120A — offline authenticated Case triage from an operator-selected snapshot

## Functional increment

The existing analyst loopback app can now run in **non-server, one-shot CLI mode**:

```shell
go run ./cmd/case-timeline-demo -principal=analyst-a -snapshot=/absolute/path/operator-owner.json -report-case=forensic-case-on-disk > local-case-triage.json
```

An explicit `-report-case` reads the existing bounded operator JSON snapshot, constructs the existing in-memory owner/authorization providers and executes **the existing authenticated `/demo-cases` and paginated `/case-timeline` HTTP handlers in-process**. There is no listener, UI, new API, authorization-model fork, external service, dependency, index or database. The ordinary demo without `-report-case` is unchanged.

Before output, it verifies the requested Case's server-authorized discovery, pages of at most 200 entries, stable total/offset/status and no duplicated visible IDs. It computes a small triage report over **only already-authorized source entries**: per-kind counts, certain UTC time bounds, count of uncertain timestamps, distinct source references, and top five authorized source-reference frequencies. It never exports raw Timeline Entry IDs, raw previews, denied sources, credentials, tenant/environment or HTTP links. It retains `visibility` (e.g. filtered) to signal partial evidence. Every page is assembled before a single write; cancellation or authorization/pagination failure produces an error without a report.

## Technology decision and scope

**Reuse** Go standard-library JSON, in-process `net/http/httptest` and the existing `NewMemoryCaseTimelineHTTPHandlers` Case/Timeline/typed-source ACL chain. **Do not create** another permission-checker, correlation engine, ORM, index, new source schema or persistence layer. The CLI synthesizes a new ephemeral bearer only for internal dispatch and never prints it. This is an analyst-initiated offline aggregation, **not** an IAM integration, forensic custody artifact or production service. The saved JSON is unsigned, unencrypted, contains authorized source identifiers and must be treated as sensitive. Redirecting stdout to a file is the operator's responsibility; this code does not choose a destination or persist output.

The data is still **operator-supplied synthetic owner data**, not authority in a real customer environment. It does not ingest Wazuh, Sekoia, EVTX or a remote event stream, preserve snapshot-level transactional consistency across changing providers, attest provenance cryptographically, or satisfy OPEN-013. For real SIEM/DFIR integration, first provide a production tenant-authenticated, durable owner and ACL provider and a controlled ingest boundary. The existing canonical Product Spec remains unchanged.

## Tests and acceptance

New behavior tests check default synthetic and actual operator JSON disk datasets, 205-row multi-page Case triage, observed/inferred counts, UTC bounds, denied source and cross-tenant non-disclosure, cancelled command and failed output writer. Reuse the established `go test -race`, `go vet`, `gofmt` and engineering-kernel CI; only an exact-SHA successful run and independent risk-based L2 can qualify the increment. This is a first **CLI investigative workflow**, not the completion of `E10-INV-120A-SNAPSHOT` or CAP-INV-110.
