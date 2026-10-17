# E10-INV-120A-SNAPSHOT — offline two-Case temporal window drilldown

## Analyst-visible functional increment

A local analyst can now select two separately authorized Cases and an **explicit UTC-aligned five-minute window** and request `-compare-events -compare-window=2026-10-01T10:00:00Z`. The command emits a single, ordered, unsigned local JSON event-level receipt containing only certain events from that window, with owning Case IDs, observed/inferred labels, provenance and source references.

It reuses `collectAuthorizedCrossCaseSequence`, which already completes **both** Case discovery, all pages of Timeline access, and typed-source owner authorization before releasing any rows. There is no new endpoint, IAM role, data reader, persistence mechanism, index, vendor package or generic execution framework. The bounded in-memory filtering is at most 10,000 authorized entries (5,000 per Case) and executes only after both Case traversals validate. A window without a certain event in BOTH Cases produces a generic error with zero report bytes, as do unauthorized/missing Cases, wrong flags, invalid UTC windows, cancellation and malformed snapshots.

### Alternative considered

The existing browser cross-Case comparison offers coarse five-minute co-occurrence and manually opens individual events. The offline CLI previously exported all events or only coarse totals. Rather than create a new temporal-correlation backend or clone existing ACL/paging checks, this increment filters the existing **fully authorized event-level sequence** for the selected window. It does not match source references across Cases; the reference memory provider deliberately enforces Case-scoped source owners.

### Trust and limitations

Operator-supplied *local synthetic* owner snapshot and synthetic analyst only, with a random local token. No provider IAM, customer telemetry, persistent events, transactional two-Case snapshot, source identity resolution across Cases, real incident correlation or causal inference. The event receipt **contains authorized raw entry metadata** and must be treated as sensitive. No output is written on partial owner failure; it is not a forensic chain-of-custody certificate. `OPEN-013` remains open.

Tests: actual two-Case disk fixture with 207 first-Case authorized rows across two pages and one second-Case event, independent ACL checks, correct [start,end) filter, observed/inferred classification, denied source exclusion, UTC boundary, uncertain time exclusion, invalid flags, unauthorized second Case, cancellation, and output writer error. Validate exact SHA through GitHub Actions and one risk-based L2. Do not claim production readiness from local CI.

## Browser integration of the same functional milestone

The loopback analyst UI extends the already authorized two-Case comparison
with an **actionable combined five-minute chronology** for each of the
five leading shared UTC windows. The action operates solely on the rows
returned after BOTH existing `/case-timeline` Case/Timeline/typed-source
owner traversals have completed, preserving the active search, kind and
UTC bounds. It does not fetch a third-party service or add a new HTTP route.

A small first-party renderer is reused for existing cross-Case shared-source
sequences and these new time-window sequences: stable certain-UTC ordering,
Case/entry ID tie-break, 25 entries per page, text-only labels. Selecting a
row switches to the originating Case when needed, then calls the **existing
server-authorized two-stage source navigation and source preview**; client
selection cannot grant access. Changing the bearer, Case, second Case or
filters invalidates stale actions. Known-time windows exclude events
flagged uncertain and do not claim temporal causality. The browser's
coarse per-Case time-window pivot remains available unchanged.

Behavioral Node tests cover complete two-Case owner page retrieval, combined
window list, second-page navigation, owner reauthorization for a peer Case,
absence of `source_ref` or bearer in URLs, and stale action invalidation.
No new SDK or abstraction boundary is introduced: the existing in-memory
DOM and Go ACL boundaries are reused. No production IAM, durable provider,
audit trail, chain-of-custody export or transaction consistency is implied.
This is one material browser integration into the already open temporal
window feature PR; do not schedule a separate L2 for cosmetic follow-ups.

## Functional Case → source pivot within a selected temporal window

The previously introduced two-Case chronological window receipt can now be followed by a **separate owner-reauthorized source-specific receipt**. The operator chooses an already visible `case_id` and its `source_ref`, then runs `-source-case=<case> -source-ref=<source> -source-window=2026-10-01T10:00:00Z` with the same trusted offline snapshot/principal. This produces `cmdr-local-source-timeline-window-v1`, returning only certain entries in `[start,start+5min)` from that exact Case-scoped source. It does **not** accept the old comparison JSON as authority: the existing Case discovery, all paginated Timeline owner checks and typed-source permissions run again.

Reuse choice: the existing `writeDemoSourceTimeline`, `collectDemoCaseTriageWithObserver` and `parseCertainFiveMinuteWindow` are sufficient. No duplicate API, general graph engine, database, new JSON owner or vendor library. Standalone source timeline without a window remains byte-compatible; invalid windows, different tenants, denied references, uncertain or boundary-excluded events, cancel and writer errors yield no JSON receipt. No source preview, bearer or custody assurance is exported. This remains **an explicit sensitive local synthetic output only**, not production IAM, a causal indicator or external SIEM integration.

Behavioral test: disk-backed two-Case HTTP traversal (with 207+1 authorized events) → selection of the actual `event-c` source in the combined time-window receipt → independent authenticated source-focused re-traversal → exact filtered `event-c` receipt. Negative paths exercise a denied principal, foreign Case, hidden reference, malformed/empty windows, uncertainty, exclusive boundary and cancellation. Existing compare and source timeline tests remain unchanged.
