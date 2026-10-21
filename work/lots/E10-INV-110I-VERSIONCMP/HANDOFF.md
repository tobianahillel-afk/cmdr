# E10-INV-110I-VERSIONCMP — Bounded canonical Investigation Timeline correction comparison

## Scope and ownership

Canon: CAP-INV-110 §3/11/12/14/19-22/25 read/compare and supersession; the selected Case is Investigate-owned while Timeline Entry and the source revision link are Shared-owned. Engineering parent E10-INV-110 (child E10-INV-110I-VERSIONCMP, dependency E10-INV-110H-UIJOIN IMPLEMENTED). Code is the existing first-party analyst browser `product-runtime/investigation-timeline/cmd/case-timeline-demo/analyst.js`; no new backend or owner, and Git parent is PR #193 SHA `0bea75b8941846d1191f38b08e743ceae53b6310`, NOT the engineering parent.

## Functional progress

The existing Case summary gathers all pages from the existing authenticated GET `/case-timeline` under original Case/Entry/typed-source authorizations. After a successful complete authorized read, source-provided `supersedes` and the permission-derived `superseded` flag are joined in memory by visible ID, with strictly increasing numeric versions and no fetch-by-guessed-ID. Up to 10 actionable comparisons display original vs corrected version, type, source and timestamp, all as textContent only. Opening either selected entry uses the existing source navigation and reauthorizes at the server. Absence of a pair **does not prove there is no correction**: successor may be filtered, inaccessible or excluded by a nontransactional page set. A true superseded flag on an older visible entry without a co-visible successor is not used to infer details. No action creates/supersedes Timeline Entry or exercises OPEN-013 class-2, no export or IAM change.

## Security and CI

The same bearer/Case/filter/revision state gates the button and local comparison. Any invalid indicator shape, nonpositive version, inconsistent visible revision ordering or broken associated successor provenance aborts the ENTIRE aggregate and leaves no report/download/partial action. Tests use real first-party browser controller with synthetic authenticated server response pages; one positive crosses 200-page boundary, one hides successor, and malformed/synthetic owner response must fail closed. No production IAM, provider, real Case store, signed evidence or native GitHub approval is claimed.

## Required reconciliation — not redundant product work

The old graph/cache at PR #193 still erroneously marked `E10-INV-120A-SNAPSHOT READY`, despite the local snapshot runtime/tests already committed on frozen #36. PR #156 `726a8c681e9a3f7876d52cf4ade0010cd72bde10` established `IMPLEMENTED`, **not VERIFIED**, and published a source-backed `RECONCILIATION.json` (blob `eef553bcad685561b3430993ab7aff3ec539e692`). This branch reuses that exact report, updates ONLY that node status and stale execution cache while preserving the separate new work nodes from #187–#193. No previous report or audited SHA is rewritten, no faux acceptance. A fresh agent must pick the next bounded independent canonical obligation, not rebuild the already working local snapshot.

## Next

Confirm PR exact HEAD, CI Timeline Linux/macOS/Windows and engineering-kernel, validate security boundaries and independent L2 if this significant analyst-facing version comparison warrants it. The PR is DRAFT; upstream #193/#191/#190/#36/#2 remain unmerged, which prevents any production or main claim. Do not modify #193 or another SHA under review.
