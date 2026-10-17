# E10-INV-120A-SNAPSHOT — authorized offline entry inspection

**Functional increment.** An operator with a local, explicitly chosen, synthetic read-only owner snapshot can now inspect the owner-provided evidence preview for a *specific authorized Timeline Entry*, directly from the CLI:

```sh
go run ./cmd/case-timeline-demo -principal=analyst-a -snapshot=/absolute/path/owner.json -inspect-case=case-a -inspect-entry=event-old
```

The output is one JSON receipt `cmdr-local-entry-inspection-v1` with Case ID, Timeline Entry ID, canonical source kind, owner object ID and preview content. It is **sensitive output**, emitted only when explicitly requested; operators must protect redirection destinations. An error yields no report, and an ordinary launch still runs the unchanged loopback analyst.

**End-to-end integration.** The command loads the bounded on-disk snapshot, instantiates the existing `NewMemoryCaseTimelineHTTPHandlers`, authenticates with an ephemeral in-process bearer, then requests the existing `/demo-source-open` handler with only Case ID and Entry ID. That handler independently checks Case membership, Timeline Entry visibility, canonical typed-source permissions and the server-resolved source owner identity. The CLI does not parse `source_ref` into a privilege, does not accept permissions from HTTP input and does not access owner preview data directly. It validates response identity and writes after the authorization request completes.

**Reuse-first choice.** Chosen: existing Go standard library `encoding/json`/`net/url`, `queryDemoOwner`, and existing typed-source authorization handler. Rejected: a second source authorization pipeline, vendor API, direct raw JSON preview indexing, ORM, new persistence layer or generic integration framework. This is a single meaningful offline DFIR drilldown in the existing Case -> Timeline -> source flow.

**Tests.** Actual end-to-end in-process HTTP owner dispatch using the built-in fixture and an operator-written disk snapshot whose preview differs from all embedded fixtures; denial for cross-tenant, missing and inaccessible entries, invalid flags, cancellation and writer errors. Run exact-SHA Go race tests, vet and CI before a pass claim.

**Boundary.** The local owner snapshot and selected principal are synthetic, not validated by production IAM. No EDR/SIEM connectors, remote ingestion, durable audit, transactionally consistent Case store, integrity signature or forensic chain-of-custody is implemented. The output exposes authorized source content by explicit choice and must not be shared publicly. `OPEN-013` remains unresolved. This is one increment of READY `E10-INV-120A-SNAPSHOT`, not production completion.
