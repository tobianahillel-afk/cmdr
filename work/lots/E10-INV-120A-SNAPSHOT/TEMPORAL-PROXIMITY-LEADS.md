# E10-INV-120A-SNAPSHOT — cross-Case nearest-time investigative leads

## Functional gap and reuse-first choice

The existing operator-supplied Case snapshot can already be read through the first-party Case picker, the owner-checked Timeline, typed source access, a cross-Case aggregate comparison and a full chronological export. The missing analyst operation was to **select actionable candidate event pairs** across two authorized Cases instead of manually scanning up to 10,000 rows or treating coarse five-minute overlap as a precise event correlation.

Chosen: reuse `collectAuthorizedCrossCaseSequence` (both Case discovery, Timeline Entry and source ACL checks, independent paging and deterministic order), Go `time.Parse`, `sort.Search` and `encoding/json`. No new provider, database, broker, IAM, backend index, third-party library, source schema or write path. A general-purpose correlation engine would be both premature and liable to falsely infer canonical cross-Case source identity.

## Operator workflow

`case-timeline-demo -snapshot=/absolute/path/owner.json -principal=analyst-a -compare-case=case-a -with-case=case-c -compare-proximity -proximity-seconds=60`

Read both Cases completely and enforce existing owner permissions before emitting a **single** local JSON receipt. Bound selected time window to 1–300 seconds. For every event in the first Case with a certain RFC3339 timestamp, find its *nearest* certain second-Case event; retain at most the 25 closest qualifying leads, sorted by absolute separation then stable entry identity. This is NOT an exhaustive Cartesian join: one second-Case entry may appear in multiple leads and other candidate pairs may be omitted. Count uncertain-clock rows separately and exclude them from candidate pairing, never promote inference to observation. Original timestamps, event IDs, Case IDs and source-declared observed/inferred kinds are preserved; raw source payload, source references, provenance, tenant identifiers, bearer token and previews are **not included** in the lead receipt. Explicitly mark results synthetic and non-causal.

Algorithm: authorize first, independently sort the known timestamps, then binary search adjacent predecessor/successor timestamps; O(N log N + M log M), at most 5000 entries per Case and 25 exported leads. Avoid Duration signed overflow for extreme chronology. No cross-Case matching by case-owned `source_ref`.

## Tests and evidence boundaries

Go integration tests exercise an *on-disk* snapshot with 207 source-authorized first-Case entries across multiple pages, a separate authorized second Case with one inferred event at 250ms separation, determinism of equal-distance ties, denied/foreign second Case, absent or conflicting CLI options, interrupted operation, unknown timestamps and no partial output. Exact-SHA GitHub Actions results and independent risk-based L2 must be verified separately before claiming acceptance.

Limitations: synthetic loopback memory provider, not production IAM, live telemetry, persistent Cases, an incident correlation verdict, replay, custody-certified report, statistically calibrated match, complete pair enumeration, or a production SLO. The selected principal is an operator-supplied synthetic identity and `OPEN-013` remains unresolved. Output is sensitive because it exposes two authorized Case/Entry IDs and observed/inferred classifications; operator must store and distribute it appropriately. R3/R4 polish must not block next material feature.
