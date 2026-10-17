# E10-INV-120A — One-step JSONL import to authorized Case proximity

Functional path: operator-local normalized event JSONL plus a read-only owner snapshot → validated in-memory append to one already readable Case → independently authenticated and ACL-filtered Case/Timeline/typed-source traversal of BOTH Cases → bounded, synthetic nearest-known-time cross-Case leads.

Example (synthetic lab only):

```sh
go run ./cmd/case-timeline-demo -principal=analyst-a -snapshot=/absolute/owner.json -import-jsonl=/absolute/events.jsonl -into-case=case-a -triage-import -with-case=case-c -compare-proximity -proximity-seconds=30 > /private/temporal-leads.json
```

The existing aggregate `-triage-import -with-case` report is unchanged. New explicit `-compare-proximity` mode requires `-with-case`, uses the existing 1..300-second window and emits an explicitly sensitive, unsigned `cmdr-local-imported-case-proximity-v1` receipt. Its `proximity` field is the original `localCaseProximity` format, with a maximum of 25 nearest deterministic leads, only from already source-authorized entries. Chronologically unknown entries do not create leads; timestamps alone never prove common origin or causation. The result is written only after both independently authorized Case reads and full calculation succeed.

Technical selection: reuse `readLocalEventFile`, `importLocalEvents`, `newDemoHandlerWithData`, `collectAuthorizedCrossCaseSequence`, `selectCaseProximity` and Go standard library, and a small shared import preparation helper. Alternative of a new database, SIEM ingest agent, correlation engine, permission provider, intermediate owner export or private forked matcher adds complexity/security risk without meeting a new requirement. No new dependency or canonical Product Spec mutation.

Risk-based checks: disk-backed two-Case synthetic scenario with a new imported event that produces an authorized 125ms temporal lead; reauthorized foreign/denied second Case yields zero bytes; malformed JSONL, ambiguous CLI, cancellation and writer errors yield no successful report. The operator source snapshot must remain byte-identical and no raw message, forbidden source or bearer may leak in output.

**Limits**: explicit local synthetic demo, not a deployed DFIR workflow, real live IAM/ACL or multi-tenant persistent SIEM source. The input owner identity and permissions are operator-controlled. The JSON contains authorized IDs and is potentially sensitive; store privately. The comparison has no transactional snapshot isolation, durable chain of custody, replay, stream ingestion, case write, production SLO or evidence of causality. `OPEN-013` stays unresolved.

CI and L2 evidence must reference the exact published SHA; no passing test should be inferred merely from publication.
