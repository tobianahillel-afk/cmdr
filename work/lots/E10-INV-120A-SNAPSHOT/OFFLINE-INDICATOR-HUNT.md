# E10-INV-120A — offline authorized cross-Case text hunt

## Investigator-visible functional journey

The analyst can use the existing operator-selected, bounded JSON Case owner snapshot to find one indicator-like substring across **all Cases authorized to one synthetic principal**, not merely two manually named Cases:

```sh
go run ./cmd/case-timeline-demo -principal=analyst-a -snapshot=/absolute/path/operator.json -hunt-q=ioc.example > local-hunt.json
```

The command produces one machine-readable local JSON receipt: query, number of authorized Cases actually searched, matching Cases and visible hit counts, and observed/inferred/uncertain classifications for each matching Case. It does **not** export event IDs, source references, evidence content, previews, access tokens, tenant IDs or inaccessible Case identifiers. An empty result means no match **in the selected principal's authorized bounded Cases**, not absence of the indicator in other tenants, sources or times.

## Technology and authorization choice

Reuse Go stdlib `net/url` and `encoding/json`, the existing `/demo-cases` authorization gate and the existing `/case-timeline?q=...` substring-search implementation. Source search happens **after** server-side Case, Timeline Entry and typed-source ACL filtering; the CLI never searches raw imported rows. The existing paginated `collectDemoCaseTriageWithObserver` is reused with an optional `q` query and an unchanged zero-query wrapper, so existing reports still operate identically.

Alternatives considered: adding a second search index, a vendor SIEM client, a new backend endpoint, separate authorization policies or duplicating the browser's cross-Case hunt (draft PR #42). All would add complexity or trust surface without improving this operator-local scripted investigative use case. This is deliberately **different** from PR #42's browser interactive hunt: it produces a bounded machine-readable offline report and can be scripted without a browser.

## Security, resource bounds and limits

- Maximum **8** discovered, server-authorized Cases; maximum **5,000** authorized entries per Case; existing paging is 200 per server request. Above the cap: error, no partial result.
- Query is printable UTF-8 text of 2–96 bytes, with no control characters or leading/trailing whitespace. It is not a regex, detection verdict or entity-resolution engine.
- Every Case is independently rediscovered and its selected pages pass the existing owner permission path. A revoked, invalid or changed page causes complete failure, not a partial output. No listener starts and no new secret is printed.
- Outputs are local, operator-requested and unsigned, and can reveal **authorized** Case names/counts. Protect redirected stdout appropriately.
- The principal and owner snapshot are **synthetic operator-supplied data**: this is not production IAM, remote-source integration, an audit trail, chain-of-custody evidence, transactional multi-Case consistency or a production performance claim. OPEN-013 and all canonical product specifications remain unchanged.

## Acceptance evidence

Integration tests consume an actual on-disk owner snapshot containing two readable same-tenant Cases, a different tenant and a linked denied source; confirm independent matching of distinct source owners, no denied/foreign Case disclosure, complete two-page 205-event search, invalid query/flag rejection, cancellation, second-Case refusal and excessive authorized Case catalog. Exact-SHA CI and risk-based L2 are required before claiming this functional milestone verified.
