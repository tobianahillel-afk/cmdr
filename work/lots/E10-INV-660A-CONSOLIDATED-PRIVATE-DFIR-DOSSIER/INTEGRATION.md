# E10-INV-660A — Integrated two-Case private DFIR dossier

## Why this integrated PR exists

GitHub's exact immutable compare PR #103 `5b52e4cdd82281c4ed5537b532b78b5064756df2` → new integration HEAD after implementation confirms the COMPLETE **16 genuine existing/new commits ahead, 0 behind**. The previous PR #105 contained the canonical ±60-second generic contextual lead/source ACL implementation, while #106 added explicit opt-in sensitive preview inclusion and got BOTH exact-HEAD Actions green *against #105*. However, #105 itself was still blocked by mandatory changed security-critical coverage at **84.52% <90%**. A successful child check does NOT prove its parent diff eligible for merge. This integrated PR tests and reviews their combined actual code/diff together **directly against PR #103**, avoiding a false claim of merge-readiness.

## Complete first-party operator workflow

1. Two separately owner-minimized private Case snapshots, including real local Wazuh 4.x and generic JSONL evidence, are verified against independent operator-held exact SHA-256 values.
2. First-party per-Case tenant/environment, principal, complete paginated Timeline, and typed-source authorization runs independently. Wazuh pairs are correlated on the same rule+agent, inclusive 60 seconds, with bounded selected pairs.
3. Operator explicitly requests `-inspect-saved-pair 1 -private-report /absolute/NEW.json`. The existing v2 private report contains independently reopened Wazuh source previews and at most ten source-authorized **metadata-only** near-Wazuh generic links per Case. No new preview disclosure by default.
4. Operator can add **`-include-pair-generic-evidence`**. The v3 report includes up to ten independently reopened, source-authorized generic source previews per Case in the SAME complete private 0600/no-overwrite artifact; no stdout and no listener. Source IDs are selectors, not grants.
5. Any case, digest, source or principal failure aborts without creating/overwriting a partial private dossier. No two sources share permissions or are merged into a cross-tenant owner graph. Source provenance is operator-local/synthetic; no production IAM, live ingest, database, legal custody attestation or signed source origin.

## New material risk-based behavior tested in this consolidation

The new end-to-end test exercises **actual post-save generic source ACL revocation in the FIRST Case**, with the file structurally owner-valid, the updated SHA-256 correctly recomputed, and the Wazuh pair plus an authorized generic source remaining available. The operator executes the full **opted-in v3 CLI**: both authorized Wazuh previews and remaining first/second generic evidence are delivered; the revoked source Entry ID, canonical source ref, principal and confidential payload are absent both from contextual lead metadata and the sensitive artifact. This ensures a valid hash cannot elevate permissions and the cross-Case pairing cannot confer access to a revoked generic source.

A separate integration test revokes the **SECOND selected Wazuh typed source** with a valid updated independent SHA-256 and requests full opt-in v3 dossier: must refuse before any publication and print no source content. Both cases target data-exposure boundaries, not cosmetic test variants.

## Checks and controls

Require exact HEAD GitHub Actions `investigation-timeline-runtime` (Linux Go/Node, macOS, Windows) and `engineering-kernel` security-critical changed coverage **>=90%** on the **combined** PR #103..integration HEAD; NEVER use #106's green child run as evidence for merged chain. Existing review approvals and L2 decisions remain mandatory. This consolidation does NOT modify L2-frozen SHA of PR #103, #104, #105 or #106; it copies their already published commits into one reviewable diff by creating a new PR based directly on #103, with one genuinely added integrated ACL security test. No forced merges and no API paid.

Decision: If the combined mandatory gate fails, keep DRAFT; do not weaken checks or publish another series of microtest PRs. Determine precise risk-based missing branch and solve one root cause or redesign instead. This is one substantive integration milestone eligible for one independent read-only L2 after a stable exact-HEAD run.
