# E10-INV-120A — one-command normalized JSONL → authorized two-Case comparison

## Usable operator-local investigation path

Prior capabilities were separate: JSONL could be imported into a sensitive owner snapshot; a separate command then compared two Cases. A one-shot triage import could summarize only one Case. The new option combines those already implemented capabilities **without writing or exporting the merged owner dataset**:

```sh
go run ./cmd/case-timeline-demo \
  -principal=analyst-a \
  -snapshot=/absolute/path/owner.json \
  -import-jsonl=/absolute/path/normalized-events.jsonl \
  -into-case=case-a -triage-import -with-case=case-c \
  > /private/imported-case-comparison.json
```

The output is an aggregate-only unsigned JSON receipt: imported event count, both independently authorized Case triage totals, observed/inferred counts, certain UTC time overlap, and bounded shared five-minute windows. Raw event messages, source previews, owner principals, bearer and hidden source IDs are never exported. **Authorized** source references can appear in top-source aggregates and are still sensitive. Neither input file is modified.

## Reuse-first implementation

Reuse `readLocalEventFile`, `loadDemoSnapshot`, `importLocalEvents`, the already existing `newDemoHandlerWithData` owner ACL chain, `collectDemoCaseTriage` paginated HTTP traversal, and `buildCaseComparison` extracted from the existing comparison function. Go stdlib only. No second correlation engine, index, ORM, database, network API, new permission model, or third-party dependency.

Comparison of alternatives: materialize an intermediate merged snapshot on disk (unnecessary disclosure and I/O); create a bespoke JSONL correlator (duplicates Case/Entry/source authorization, high risk); reuse existing memory-backed comparison (selected for bounded local integration).

## Security and limitations

This path remains a synthetic **local operator-controlled lab**: the operator-supplied owner snapshot defines fixture ACLs, and the CLI principal is a synthetic selection. Neither is independently authenticated production IAM. Each Case is validated independently through the existing HTTP Case/Timeline/source owner checks; a denied second Case or invalid JSONL must return an error with **zero output bytes**. The result is not a replay, live SIEM ingestion, durable evidence store, production forensics custody proof or real-time transactional snapshot. Temporal co-occurrence is not causation. `OPEN-013` and Shared Timeline Entry ownership remain unresolved.

Acceptance: operator disk snapshot with 207 authorized original first-Case rows plus 2 imported JSONL observations, independently authorized second Case with one inferred row, protected third tenant and restricted source. Verify count 209 vs 1, observed 106 vs 0, inferred 103 vs 1; deny foreign Case and malformed JSONL without output; compare across multiple pages; original file unchanged. Independently verify GitHub Actions CI on exact source SHA and risk-based L2 only after the functional increment is stable.
