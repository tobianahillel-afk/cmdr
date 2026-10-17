# E10-INV-120A — actionable local cross-Case sequence export

## Functional increment
An analyst who has explicitly compared two separately authorized offline Case timelines may inspect either a shared 5-minute UTC interval or a common source-label sequence and **voluntarily export all currently selected, already authorized event rows** as local JSON. The existing viewer's 25-row pagination is only presentation; the downloaded local receipt contains the complete bounded selection. The record contains Case and entry IDs, declared kind, source reference and chronology/uncertainty. It does not contain the bearer, source preview or server-side access capability.

The export is an operator-local, unsigned, sensitive synthetic report. It is **not** production evidence, a stable data snapshot, custody/signature, a live-source query or proof of causality. Source labels from different Cases need not denote the same owner object; rights can change after reads. Freshly opening any source still traverses the existing owner authorization endpoints. A bearer/Case/filter edit invalidates the UI export action, but revocation between an already-authorized read and offline file creation cannot be inferred by the client.

## Reuse-first technical choice
- Reuse existing `analyst.js` independently authorized paginated Case reads, cross-Case chronological sort, embedded HTML, Blob downloads and source-owner open path.
- Avoid a new HTTP endpoint, index/search engine, storage layer, third-party export package or client-side permission engine; these would create new security ownership and unnecessary complexity.
- Max 5,000 visible rows per Case (10,000 combined); explicit click only, no browser persistence or automatic download.
- Behavioral Node DOM tests verify selected 5-minute event export, sensitive reference scoping, 25-row page independence, zero automatic exports, bearer omission and invalidation after token rotation.
- No changes to canonical Product Spec, OPEN-013 or IAM claims.

## Verification
This note describes intended behavior only; the exact SHA CI and L2 findings must be recorded after GitHub confirms results. Do not equate a GitHub Actions green check with production authorization or independent approval.
