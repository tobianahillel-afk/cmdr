# E10-INV-560A — Atomic private inspection of both sources in one correlated Wazuh Case pair

## First-party functional integration

The existing workflow already saves and independently verifies two operator-owned Wazuh 4.x Case archives, runs two-Case owner-authorized full Timeline comparison, and generates stable Wazuh agent+rule event pairs within an inclusive 60-second UTC window (#93/#94). It can also re-open *one selected source at a time* through the first-party typed-source owner into a private report (#95).

**This lot links those capabilities into one analyst action:** select the 1-based pair number (1–10) from the prior bounded cross-Case comparison; produce **ONE** atomic, private, no-overwrite report containing the canonical pair metadata and the two fully re-authorized actual source inspections. No copy/paste of Entry IDs and no partially written first-side receipt if second-side access fails.

```sh
go run ./cmd/case-timeline-demo \
  -principal analyst-a \
  -snapshot /absolute/first-saved-case.json \
  -compare-saved-case CASE_A \
  -expect-snapshot-sha256 <operator-held-hex-A> \
  -compare-saved-with-snapshot /absolute/second-saved-case.json \
  -compare-saved-with-case CASE_B \
  -expect-other-snapshot-sha256 <operator-held-hex-B> \
  -inspect-saved-pair 1 \
  -private-report /absolute/NEW-private-pair-source-report.json
```

Implementation intentionally **reuses** the complete `compareVerifiedSavedCases` path (independent digest checks, separate principal/Case/tenant/environment and full Timeline/typed-source owners, deterministic existing 60s pair algorithm). The ordinal selects an established first-party pair from that verified report. Both snapshots are reloaded independently using the unchanged `loadAuthorizedImportedCase`, the selected imported IDs are confirmed in each *own* Case, the scopes rechecked and **both** typed-source previews separately obtained through `inspectDemoEntryWithData`. Neither Entry ID is an access grant.

The entire operation runs inside existing `writeAuthorizedReceiptPrivate`, which stages bounded JSON and publishes a new private 0600 file via the existing safe atomic no-overwrite path. No evidence preview or bearer appears on stdout; never start an HTTP listener; no new disk store, source parser, auth model, search engine, paid API, dependency or report writer.

The report is sensitive and contains **both** source previews, which a legitimate analyst is already authorized to inspect independently. Treat as private evidence and do not share it publicly. Correlation remains rule+agent/time co-occurrence, **not** identity, causal proof or trusted chain of custody. Digest provides change detection against independently retained bytes, not authenticated provenance or anti-rollback.

## Behavioral security and acceptance

- Use two **actual Wazuh 4.x JSONL** files imported into separate private Case archives and independently fingerprinted; original JSONL files deleted. Generate real cross-Case pair with 20-second offset. Running exact CLI `-inspect-saved-pair 1 -private-report ABSOLUTE_NEW` produces ONE valid private bundle with correct Case/Entry IDs, cohort and both real source previews. No stdout, token, remote listener, partial receipt or overwrite on a repeat.
- Reject wrong first or second digest, modified-but-still-valid second file, foreign principal, missing/foreign Case, absence of selected pair, out-of-range ordinal, missing/relative/private output and mode conflicts. Ensure no private file or stdout produced.
- Reject two individually valid but **different-tenant** owners even when second file has correctly recomputed digest and its own principal grants access.
- First-party owner/access checks remain mandatory; only a Case-authorized source may be included, no cross-Case ACL blending.

## Delivery discipline

Child of consolidated draft parent #96 frozen exact SHA `f16a5056875cb19bb9081712beff567675685089`, base consolidated #92 frozen SHA `fced2c5c89f0458247be9cc0644e516118af2ae6`; all audited ancestor commits remain unchanged. Parent #96 has real runtime CI SUCCESS but engineering-kernel `CHK-SECURITY-TESTS` only **88.22%** against mandatory 90%. **Neither this child nor #96 is mergeable unless the unchanged security gate and required native reviews succeed on the proper integration diff.** This is synthetic offline operator-controlled lab functionality, not production tenant IAM, continuous live ingestion, cryptographic custody, or durable search storage. One independent risk-based read-only L2 on this substantive two-source disclosure boundary after exact-head CI; never duplicate prior L2 per commit.
