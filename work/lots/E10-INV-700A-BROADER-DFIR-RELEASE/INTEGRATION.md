# E10-INV-700A — Consolidate one mixed-source DFIR deliverable against the earlier green integration base

## Verified provenance and genuine delivery goal

This is a **real GitHub integration of already implemented working product code**, not a new source owner, detection engine, workflow framework, or security policy exception. A native GitHub commit comparison before this handoff confirmed:
- Immutable upstream **PR #81** SHA `c173c34dbe224b6427cfb37ed880da9266504e1c` (security coverage diagnostics; both of its exact-head Actions SUCCESS).
- Immutable full functional **PR #110** SHA `7b4cc01dcb8bbc5cf64e67e543b5ccdf5ea4afe5` (the existing integrated DFIR runtime; both of its exact-head Actions SUCCESS).
- Native GitHub compare `#81..#110`: **111 genuine published commits AHEAD, zero behind, 60 changed files**. A new child integration branch points directly to the same #110 source SHA, with **no rebasing, source rewrites or changes to any earlier L2-reviewed SHA**. This handoff is its ONLY new documentation commit.

### Material integration objective

Previously PR #92 was the direct predecessor of release candidate #110. GitHub Actions for immutable #92 HEAD `fced2c5c89f0458247be9cc0644e516118af2ae6` include a mandatory engineering security failure on its own larger diff, although its runtime is green. Thus the source candidate #110 passing its smaller `#92..#110` changed-line security gate DOES NOT permit merging #92 or claim all upstream dependencies are satisfied.

This candidate deliberately tests the **full `#81..#110` product diff as a SINGLE GitHub PR**, replacing that one blocked dependency hop with a single coherent reviewed candidate **IF AND ONLY IF** its own exact-HEAD CI, risk-based L2 and native GitHub approvals all pass. No source PR #82..#110 is force-merged, closed, rebased or retroactively marked passed. If this larger security scope fails the >=90% gate, it remains DRAFT and cannot replace #92. Other ancestors, particularly PR #80 (its own changed-security coverage 87.90% <90), foundational #36 (763 commits/871 files), and product spec PR #2, remain independently unmerged/unresolved. No claim that the whole stack becomes mergeable.

## Existing integrated first-party executable investigator journey

One operator-local, synthetic, offline Case/Timeline/typed-source graph can combine Wazuh 4.x and generic normalized JSONL events from selected local files. Case-minimized private snapshots are saved explicitly as POSIX 0600 temporary files and atomically published without overwriting existing evidence. Their exact bytes are SHA-256 checked against independently held operator receipts before owner initialization. Original raw input feeds are no longer required to inspect a saved Case.

The same owner-authorized first-party Event Search Jobs search imported generic events. Fully paginated, tenant/environment/principal/Case/Timeline/source-approved Wazuh observations from two **independent** saved owners produce 5-minute cohorts by rule+agent and <=60-second UTC event-linked pairs. The Case-to-Case comparator provides aggregate time-window co-occurrence, NOT malware diagnosis or causal attribution. A selected pair index is ONLY a selector, never an authorization grant: each Wazuh typed source must be independently reopened on its original owner before a private dossier can be written.

The dossier defaults to v2: Wazuh source previews in ONE operator-selected private no-overwrite file, with generic context only Entry/time metadata (up to ten normalized sources per Case within inclusive +/-60s). Explicit `-include-pair-generic-evidence` AND an explicitly selected private file permits v3: generic source previews are reopened through their own first-party owners once, staged only until BOTH Case/source authorizations succeed, then written to that same private file. Rejected ACLs, stale/mismatched SHA, missing owner, cross-tenant/environment mismatch, wrong Case/source or existing destination must never cause partial evidence or raw preview/token on stdout.

No new libraries, paid API, online ingestion, production IAM, durable multi-tenant database, signed chain-of-custody, GPU workload or deployed SIEM. Synthetic operator-supplied input is NOT production forensic truth.

## Required acceptance

1. Run `investigation-timeline-runtime` on THIS precise published HEAD, including Linux Go + Node/browser, macOS and Windows; test real Case import/search, mixed-source correlation, source-specific ACL revocation, typed source-open and v2/v3 dossier publishing using existing tests.
2. Run the **UNCHANGED mandatory** `engineering-kernel` checks with changed security-critical coverage >=90% **against THE NEW PR BASE #81**; never infer this from PR #110's smaller diff or reduce/waive any threshold. Compare precise GitHub Action head/base and logs, not status cache.
3. Independent risk-based read-only L2 only for this NEW materially broader consolidation, referencing the exact published base/head diff. Existing #110 L2 is separately scheduled; do not duplicate it, do not change frozen SHAs. GitHub native approval and upstream checks remain required.
4. Never promote `engineering/state/current.json` cache or `work/graph.json` flags above current git history, CI evidence and canonical `cmdr-product-spec`. All runtime work remains draft until actually merged.
