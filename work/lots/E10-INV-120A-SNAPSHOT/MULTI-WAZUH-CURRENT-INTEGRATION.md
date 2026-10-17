# E10-INV-120A — Multi-Wazuh authorized local investigation on current foundation

## Functional increment

Rebase-by-port, not force-merge, from historical draft PR #48 onto the current verified engineering foundation. This preserves existing rule-ID, agent, time-window triage and browser cohort capabilities while adding **2–4 distinct Wazuh 4.x alert JSONL feeds** into one existing authorized Case. The browser provides the existing authenticated Case → Timeline → typed-source-preview journey. Separate offline modes expose:

- `-serve-import`: authorized localhost Case viewer for the selected imported feeds.
- `-triage-import`: per-Case aggregate receipt, or, with `-with-case`, independently authorized cross-Case aggregate comparison.
- `-triage-import -with-case=<case> -compare-proximity`: closest certain-time cross-Case leads within an explicit 1..300-second window.
- One `-wazuh-alerts` flag retains existing single-feed triage, comparison, proximity and viewer behavior; repeat the flag 2..4 times to opt into multi-feed handling.

Example:

```sh
cd product-runtime/investigation-timeline
go run ./cmd/case-timeline-demo \
  -principal=analyst-a -snapshot=/absolute/path/owner.json \
  -wazuh-alerts=/absolute/path/sensor-a.json \
  -wazuh-alerts=/absolute/path/sensor-b.json \
  -into-case=case-a -triage-import -with-case=case-c
```

## Reuse and trust choice

Chosen: first-party Wazuh normalizer → bounded local JSONL importer → in-memory Case owner → existing authenticated Case/Timeline/typed-source handlers → existing aggregate comparison and proximity algorithms. Simpler and safer than a fresh Wazuh SDK, real-time queue, database, search index, synthetic correlation engine or generic authorizer. All authorization flows through the *existing synthetic owner*, not alert claims or browser flags. Feed ordinals namespace identical raw alert IDs, but **are not cryptographic source identities** and are tied to operator-selected ordering.

Every file is validated before listener/token issuance or report output; duplicate file identities and invalid data fail closed; no intermediate owner export is written. Independent Case authorization is required on comparison/proximity. No source preview, bearer, local path or cross-tenant owner snapshot appears in JSON reports. Files are read-only.

## Explicit limits and acceptance

**Not production**: synthetic local principal/ACL, no production SSO/RBAC integration, no live Wazuh connector, persistent Case storage, forensic custody chain, signed origin or verified causal correlation. `OPEN-013` remains unresolved and canonical Product Spec is unchanged.

Tests exercise two and four feeds with colliding Wazuh IDs, real loopback typed-source previews and unauthorized reads, independently authorized cross-Case comparison/proximity, triage, bad second file and zero premature output. Reuse existing Go race/vet/gofmt and security checks. The old PR #48 is dirty against its changed base and must not be merged blindly; this new child branch keeps its own exact-SHA CI and risk-based L2. The changes remain a **draft** until CI and independent review are confirmed. No changes to main, and no general production-ready claim.
