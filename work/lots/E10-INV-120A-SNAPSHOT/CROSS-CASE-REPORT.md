# E10-INV-120A-SNAPSHOT — downloadable, authorized cross-Case triage receipt

## Functional user journey

In the loopback local analyst, an operator may select two distinct readable synthetic Cases, press **Comparer les chronologies autorisées**, and then manually save a **local JSON comparison receipt**. The comparison reuses the existing Case Timeline HTTP pagination and typed-source permission chain for **each** Case (5000 entries per Case; no extra endpoint). The newly saved file contains only bounded aggregates: authorized counts, observed/inferred distinction, uncertain timestamps, certain-time bounds/overlap, up to three shared source references and up to three overlapping five-minute windows. No bearer, source preview, individual Timeline Entry row, or entry identifier is written to the receipt.

The feature turns an otherwise on-screen-only cross-Case investigation into a portable, **unsigned local triage note**; the local file is not a cryptographic custody artifact, signed report, production export, live correlation index, proof of causality or a durable system of record.

## Technology comparison

- **Selected:** existing same-origin `/case-timeline` authorized HTTP pages, existing Go Case/Timeline/source ACL owners and the current first-party browser `Blob`/object-URL download pattern. All access decisions stay in the server; a saved receipt is never a credential.
- **Rejected/deferred:** another Case/search API, a new server-side exporter/permission gate, relational/indexing backend, workflow engine, external analytics SDK and new third-party dependencies. A transactional forensic repository and trusted identity provider are separate production milestones.

## Security and scope

- Both Cases and **all** their selected paginated entries must be authorized, complete and consistent before the download button is enabled.
- Case/second Case/bearer/filter changes invalidate the in-memory report. A second-Case denial, pagination inconsistency or invalid record never yields partial saved data.
- JSON is generated from authorized aggregates and saved only after an explicit click; source references are display-capped to 120 code units, not interpreted as HTML. Top sources and event contents are never downloaded.
- All request credentials stay in `Authorization` headers, never query strings or report bytes. The host owns the synthetic principal and original operator-selected read-only JSON snapshot.
- Results remain **best effort** because two sequential page traversals are not transactionally consistent. A shared source string or overlapping time bucket is merely a triage lead. There is no production IAM, live event ingestion, audit log, cryptographic signature or chain-of-custody attestation, and `OPEN-013` remains unresolved.

## Acceptance and evidence

Browser behavioral tests exercise multi-page Case comparison (205 vs 3), exact authorized aggregates, no raw rows/bearer/URLs in the receipt, explicit click-only download, no export on second-Case access denial, invalidation after token/Case changes, and fail-closed incomplete pagination. Existing Go tests cover the disk-to-owner-to-HTTP Case/Timeline/source permission chain; this addition does not modify those security decisions. Check GitHub Actions **on the exact published SHA**, plus independent risk-based L2 for the combined functional delta. Do not interpret historical green tests as approval of this feature.
