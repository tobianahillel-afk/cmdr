# E10-INV-120A-SNAPSHOT — offline authorized source timeline

## Functional delivery

An analyst can now ask the already implemented local operator-snapshot tool for a **single exact source reference** inside an authorized Case, with:

```shell
go run ./cmd/case-timeline-demo -principal=analyst-a -snapshot=/absolute/path/operator.json -source-case=case-a -source-ref=demo-evidence-a
```

The output is one unsigned JSON receipt containing only events from that specific source that passed existing Case membership, Timeline Entry and typed owner-source permissions. It preserves original kind (observed/inferred), source timestamp and uncertainty, provenance, version and accessible supersession. It includes the Case-wide visibility state so a filtered/partial view is not misrepresented as exhaustive. The full authorized Case pagination (up to 5000 rows) must finish before any output. Missing/denied source, unauthorized Case, cancelled request and pagination failures return an error with **zero report output**. No new route, background service, database or persistence layer exists.

## Technology comparison

**Selected:** the existing first-party owner HTTP handler and `collectDemoCaseTriageWithObserver` to reuse one server-side Case/Entry/source ACL and pagination implementation; Go standard-library JSON and a linear exact-source selection over already bounded authorized pages. **Rejected:** a separate source-reference authority, direct JSON scanning without owner authorization, a new search/index database, a second pager, a third-party query engine, or source identifiers in a browser URL. This approach is deliberately a local forensic investigation integration rather than a benchmark or validation-framework extension.

## Trust and limitations

The principal and owner data are **operator supplied and synthetic**; `-principal` is not enterprise IAM. A real production provider, event ingestion, remotely authenticated Case owner, transactional snapshot, forensic signing/chain-of-custody and production SLO remain absent. The report contains authorized entry IDs and source references, and is **sensitive local stdout output**; do not publish it, upload real customer snapshots or treat it as court-ready evidence. No metadata from other or denied sources appears in the report; the explicit source reference itself comes from CLI input. `OPEN-013` remains unresolved. The product spec and Shared-owned object contracts remain unchanged.

## Acceptance

Run the existing source and triage tests plus new end-to-end disk-snapshot, multi-page, denied source, cross-tenant, cancelled-request and writer-failure tests, `go test -race ./...`, `go vet`, and exact-SHA engineering-kernel/Timeline CI. Independent risk-based read-only L2 is required for the material functional capability. Document CI conclusions only after verified.
