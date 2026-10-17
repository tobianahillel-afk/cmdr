# E10-INV-120A — functional cross-Case event pivots

## Reuse-first choice

The first-party analyst already implements authorized Case discovery, Timeline text search, per-entry source navigation and owner preview. Draft PR #42 previously listed only total text hits per Case. The new journey exposes up to 3 already-authorized event IDs per Case, allowing the analyst to pivot directly to an exact event search in the ordinary Case Timeline form. That form and the original Go Case/Entry/typed-source ACL chain reauthorize the entry and eventual source. No new search index, provider, database, routing API, persistence, dependency or schema is introduced. Unlike a new global search engine or unverified source-preview cache, this reuses the existing bounded HTTP path and authorization boundary.

## Workflow, bounds and security

Host-provided offline owner snapshot (existing CLI) -> ephemeral loopback bearer -> GET /demo-cases -> all GET /case-timeline?q=<indicator> pages -> bounded Case totals and max 3 event shortcuts per Case -> select Case + event ID in existing reauthorized Timeline search -> user opens typed source via existing reauthorized source workflow. Up to 8 Cases and 5000 matches per Case. A new case or bearer change clears results; a rejected later Case clears ALL partial results. No source reference, preview, bearer or HTML is placed into shortcut text, and no source access is performed in the cross-Case hunt itself. Event IDs are capped at 1024 characters and output uses DOM textContent.

## Functional evidence criteria

Node browser-harness tests exercise two Cases with four hits each, three shortcuts per Case, exact target Case/event, filter clearing, absence of browser-side source fetches or protected reference labels, stale bearer invalidation, later-Case authorization failure and oversized event IDs. Existing Go HTTP integration exercises real server owner ACLs, but must be rerun at the final SHA. Treat CI success as verification of this bounded offline feature, not a production IAM/SIEM ingestion, transactional cross-Case snapshot, IOC verdict, source-proof or audit/chain-of-custody capability. OPEN-013 remains unresolved. Draft PR #42 targets the engineering integration branch; it must not merge to main before base #36 and required reviews.
