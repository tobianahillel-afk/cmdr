# E10-INV-120A — comparative chronology of visible Case source references

## Functional increment

The operator-selected offline Case journey already uses the trusted first-party server path for Case discovery, paginated Timeline data and source-owner access. Its analyst viewer compares exact source references between **two separately reauthorized Cases**. This increment makes the temporal ordering of each common source **explicit and actionable for investigation triage**:

- If the latest *certain visible* timestamp in Case A is strictly before the earliest certain timestamp in Case B, present `case_a_before_case_b` and the minimum observed-envelope gap.
- The inverse produces `case_b_before_case_a`.
- Intersecting/touching certain timestamp envelopes produce `certain_envelopes_overlap`. **This is not evidence of simultaneous events.**
- If either Case has no certain timestamp for that reference, produce `indeterminate`, with no fabricated gap.

The comparison's on-screen labels and its existing explicitly requested local JSON aggregate both carry these values. The numeric gap is in **milliseconds**, including `0` for overlapping envelopes and `null` when indeterminate. All inputs are the already-authorized pages, still capped at 5,000 visible entries per Case. Unknown-time events remain counted but must not establish chronology.

## Design choice / reuse

Reuse JavaScript `Map`, `Date`, the existing `sourcesFor` certain-time envelope, server-authorized `/case-timeline` paging, current comparative report export and its Node VM functional test harness. A new correlation index, temporal database, event source, microservice, identity provider or third-party library would add risk without improving this bounded in-process offline triage.

## Boundaries and verification

Existing Case/Entry/source checks remain server-side. Client comparison never authorizes another Case, navigates to a source, updates data or grants evidence custody. Two-case pages are nontransactional; observations can be missing/filtered/revoked, and sharing a reference proves **neither causality nor the same underlying source object**. The result is intentionally called a *visible certain-time envelope relation*, not chronological proof of a real attack sequence.

Behavioral tests cover A-before-B, B-before-A, intersecting ranges, unknown-only A timestamps, exact bounded export, denied Case and absence of token/raw Event IDs in aggregate. CI and independent review must be verified at the published SHA; this document is not itself evidence that the checks passed.

The local synthetic demo does **not** provide production IAM, durable evidence storage, SIEM ingestion, transactional snapshots, immutable chain of custody or a resolution to OPEN-013. Work unit remains READY pending its completion gates and risk-based independent review.
