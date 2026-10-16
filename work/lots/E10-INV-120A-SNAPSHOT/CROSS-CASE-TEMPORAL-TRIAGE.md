# E10-INV-120A — meaningful cross-Case temporal and classification triage

## Functional gap and reuse-first decision

The existing authenticated loopback Case comparison displayed only exact `source_ref` intersections. But the real `memoryTimelineProvider` binds every canonical source reference to **one Case**, so cross-Case source intersections from valid operator snapshots generally equal zero. This is correct source isolation, but a poor investigative comparison.

Reuse existing paginated Case/Timeline/typed-source-authorized `/case-timeline` API and the in-browser comparison loop. In the same complete two-Case read, count **observed versus inferred** rows, **uncertain** timestamps and compute the earliest/latest **certain** timestamps of each Case. Report whether these two bounding intervals overlap; "indeterminable" if either Case has no certain timestamps. Also keep exact-ref intersections as a lead, never as an assertion of shared canonical objects or causal link.

Chosen: native JavaScript Map, Date and existing handlers, O(n) additional work, no third-party library, no new service, no cross-tenant index, no new source authority, no persistent writes or new API. A backend correlation/index service is not justified for this 5,000-entry-per-Case, operator-loaded synthetic integration.

## Security/forensic limits

- The *same* existing authenticated, no-store, same-origin requests read both Cases. On either denial, bad timestamp or inconsistent later page, **no partial comparison** is displayed.
- Event categories are counts **only for authorized entries** and never promote inferred records into observed facts.
- Only entries marked `time_uncertain=false` with strictly valid RFC3339 dates contribute to the time bounds. An uncertain entry counts but never stretches a bound. Date.parse normalizes offsets to UTC; JavaScript millisecond resolution is not original nanosecond forensic precision.
- Range overlap means *only* overlap of minimum/maximum visible certain timestamps; it does not mean continuous activity or shared actors. Source references remain Case-owned; there is no generalized cross-Case entity identity.
- No production IAM, durable Case persistence, live SIEM ingestion, custody attestation or Shared Timeline Engine is introduced. `OPEN-013` stays open.

## Required exact-SHA acceptance

Node browser tests exercise 205-entry/3-entry paged comparison, observed/inferred counts, unknown timestamps, no common source IDs while still returning meaningful comparison, and invalid second-Case timestamps failing closed. Existing tests cover denied second Case, token revocation, cancellation; Go tests cover disk JSON to Case, Timeline and source ACL HTTP. Verify GitHub Actions on the exact published SHA before claiming PASS. L2 should review the combined functional integration milestone, not each cosmetic follow-up.
