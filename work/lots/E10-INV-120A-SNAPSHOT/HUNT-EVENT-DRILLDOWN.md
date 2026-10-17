# E10-INV-120A — explicit selected-Case event drilldown from offline cross-Case hunt

## Analyst journey / functional integration

The existing `-hunt-q` emits permission-filtered cross-Case **aggregate** leads without source metadata. The new, separately selected `-hunt-events-case` mode lets an analyst investigate one matching authorized Case by retrieving the **first 25 authorized matching Timeline Entries** (IDs, observed/inferred classifications, source timestamps, provenance and typed source references) from the same paginated Case/Timeline/source HTTP authorization chain. The receipt tells the analyst the full visible matching count and whether the detailed list was truncated. From an authorized ID, the existing `-inspect-case` + `-inspect-entry` mode can request typed-source inspection with fresh owner permissions.

```sh
# Step 1: find which authorized Case and time windows contain an IOC-like TEXT hint
go run ./cmd/case-timeline-demo -principal=analyst-a -snapshot=/absolute/path/operator.json -hunt-q=ioc.example > aggregate.json

# Step 2: explicitly open one authorized Case's matching events
go run ./cmd/case-timeline-demo -principal=analyst-a -snapshot=/absolute/path/operator.json -hunt-q=ioc.example -hunt-events-case=case-a > selected-sensitive-events.json

# Step 3 (if justified): inspect one of those authorized event IDs through existing owner ACL checks
go run ./cmd/case-timeline-demo -principal=analyst-a -snapshot=/absolute/path/operator.json -inspect-case=case-a -inspect-entry=event-old > selected-sensitive-source.json
```

## Why this approach

Reuse the first-party `collectDemoCaseTriageWithObserverAndQuery`, the existing paginated `/demo-cases` and `/case-timeline?q=` owner handlers, and the existing typed-source inspection CLI. Alternative approaches such as raw JSON snapshot search, a new index, a cross-Case source join or a new permission policy are unnecessary for a bounded local analyst flow and risk bypassing tenant/source isolation. Go standard library only, no new product dependency, database or production identity provider.

## Data, access and risk boundary

- Default cross-Case hunt remains **aggregate-only**. Detailed event metadata requires the analyst to specify **both** `-hunt-q` and `-hunt-events-case` in a separate mode; it never mingles data from two Cases in the detailed receipt.
- The selected Case must be discoverable to the authenticated synthetic principal through the existing owner. Timeline Entry and typed-source authorization is re-evaluated through every paginated handler request. Denied and foreign events never enter the report. Failure on a later page aborts the entire report; no partial JSON is emitted.
- Maximum **25** detailed events in the output, with `total_matches`, `returned`, `truncated` to prevent silently describing a subset as exhaustive. The underlying permission-checked iteration still consumes **all** matching pages (max 5,000 entries/Case) before returning. This is not transactional consistency.
- The new JSON receipt **contains sensitive authorized IDs, provenance, source references and timestamps**. It is local operator-requested output only; protect its file permissions and do not post it to public logs, tickets or GitHub. It is not forensic chain-of-custody evidence.
- `query` is literal printable text, not a malware verdict, rule, IOC enrichment or correlation proof. Synthetic, operator-controlled owner state is **not** production IAM. No network listener or remote data source is introduced in the CLI mode.
- The canonical Shared Timeline Entry ownership, `OPEN-013`, production persistence and product specification are unchanged.

## Functional acceptance

Two-Case disk snapshot + hidden source + foreign tenant; explicit selection of each authorized Case; denied Case fails with no output; a 205-match source paginates fully but exports only 25 authorized entries with `truncated=true`; cancellation and second page revocation fail without output; mutually exclusive modes cannot expose a sensitive report by accident. CI and any independent L2 must be linked to the immutable published SHA.
