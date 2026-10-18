# E10-DFIR-780A — Recoverable consolidation of authorized Case investigations

## Correct engineering parent and product ownership

This is **one integration work unit under the transversal engineering epic E10-DFIR-780A**, not a newly invented product capability. Existing canonical behavior is governed by:

- `CAP-INV-110` (Investigate / Cases and Evidence): Case-scope Timeline selection, provenance, source navigation, observed versus inferred, source ACL and chronological context; the underlying Timeline Entry remains Shared-owned.
- `CAP-INV-002` (Investigate / Signals and Hunt): manually initiated authorized telemetry query, run and results; Query / Search Job remain Shared-owned.
- `CAP-INV-004` (Investigate / Signals and Hunt): rendered/raw source inspection and manual source pivot subject to independent permissions; source events remain Shared-owned.

The already implemented Wazuh 60-second rule+agent cohorts are **derived navigation leads**, not a canonical replacement of any owner, telemetry authority, identity provider, forensic custody signature or Evidence qualification. No mutations resolving `OPEN-013` or `OPEN-014` are claimed.

## Git ancestry and why integration is needed

As independently confirmed before branching, `github.com/tobianahillel-afk/cmdr`:

| Source PR | Exact frozen Git HEAD | Role and check state |
|---|---|---|
| #113 | `c11809d41a54851d709ed19f33a10b5771cc2c06` | Parallel broad DFIR candidate; review / CI independence remains, NOT selected as immediate base |
| #120 | `db97ebb3299de6f2e9c3c1558553fdf78174bb35` | Unified Wazuh pair navigations, incorporated through native merge ancestry |
| #122 | `ef0fd7088447f85ed241ef1115645dd3aa8366e0` | Source-authorized cohort catalog implementation from parallel branch, source ported into #124 |
| #123 | `18fbbeb871813e7fb5ac1cc89dec575303f41805` | Parent integration; engineering-kernel FAILURE at run 39172895573, **not mergeable by itself** |
| #124 | `37c9a64714d5a0ae4a08691c69480126500b1bf8` | Consolidates #123 plus original catalog from #122; `investigation-timeline-runtime` run 39173197573 **SUCCESS**, `engineering-kernel` run 39173197541 **SUCCESS** |

GitHub compare `#113 HEAD ... #124 HEAD` returned **ahead 17 / behind 0 / 21 changed files**; however #113's larger candidate has a mandatory security-coverage failure on its OWN diff. The chosen more specific Wazuh navigation parent is **PR #120** at `db97ebb3299de6f2e9c3c1558553fdf78174bb35`: native compare `#120 HEAD ... #124 HEAD` returns 0 behind, and the new integration candidate was confirmed as 12 changed files against #120. This is a coherent cross-lot *integration epic*, not a single implementation task of 21 files; the original independent code and behavioral tests are inherited unchanged, including revocation of canonical typed sources after re-hashing and independent source-owner reauthorization. No force push/rebase, code duplication or rewritten L2 SHA.

This newly created integration branch is a **direct Git copy** of original #124 HEAD. The only new commits are engineering traceability / reconciliation and this handoff. New PR must compare against unchanged #120, not against an arbitrary older fork, so it can consolidate all tested behaviors under one reviewable candidate rather than require either failed broader release PR #113 or #123 to be merged as an intermediate. Prior source PRs stay open and their provenance intact; closing/rebasing/fusing any of them is out of scope.

## Actual delivery evidence and limitations

Supported *synthetic* analyst journey from first-party Go code:

1. Operator opens two separately saved local Case owners with independently held SHA-256 fingerprints and per-Case principal/tenant/environment authorization.
2. Collect Wazuh rule+agent cohorts only after both owners and typed sources are authorized; navigate deterministic 60-second temporal leads.
3. Use catalog `rule_id` and `agent_id` to filter paginated pair leads, select original Case-scoped Entry IDs and reopen owner-authorized source previews or explicitly requested private evidence dossiers.
4. A denied/revoked source disappears from cohort and pair navigation on fresh verification; no ID becomes an access grant.

**Not delivered:** deployment on main, real operator IAM, continuous collectors, production storage, signed custody, mature time-index, general cross-tenant correlation, production security attestation or feature acceptance for all of CAP-INV-002/004/110.

The broad parent #113 and intermediate #123's red CI results are neither waived nor recharacterized as green; the candidate MUST pass its OWN exact SHA workflows after reconciliation commits. The broad #36 foundation remains DRAFT, with native rules/reviews as separate gates. Only after independent L2 and native approval may integration be considered. No merge initiated here.

## Reconciliation of stale execution cache

`engineering/state/current.json` and `work/graph.json` previously kept pointing to `E10-INV-120A-SNAPSHOT READY`, although GitHub history already contains the original bounded snapshot feature plus multiple later independent investigations, imports and correlations. The original manifest's `outputs` and `acceptance_tests` were incrementally extended to many distinct behaviors, exceeding its own **one separately testable behavior** policy. That historic manifest remains intact rather than rewriting its supposed estimates or erasing provenance.

The original node is conservatively changed from READY to **DECOMPOSED**, **not VERIFIED/ACCEPTED**; this stops blind reexecution as one giant unit. A separate transversal epic and one bounded integration task reference actual multi-capability and Git dependencies. No official claim of CAP-INV-110 completion, production or accepted user-facing feature is made.

**Next actions:** Verify GitHub compare and SHA, run all mandatory CI, ask for one context-separated read-only L2 of the significant merge-base/integration delta, preserve native branch protection, then select a fresh bounded functional obligation from canonical product specs. If the original E10-INV-120A unit requires final closeout, decompose its outstanding acceptance obligations into separate correctly owned tasks rather than expanding it further.

## Corrected Git integration choice (same product code, no SHA rewrite)

The original exploratory base #113 produced a 26-file security-critical diff and a concrete kernel gate failure (88.87% below the mandatory 90%). Instead of adding nearly duplicate tests or changing the gate, the actual product-boundary parent is #120, containing the existing authorized Wazuh pair navigation. The explicit Git base is retargeted to `feat/e10-inv-750a-unified-wazuh-pair-navigation` / immutable SHA `db97ebb3299de6f2e9c3c1558553fdf78174bb35`. The candidate still contains source-provenance contributions from #113 and #122 via Git ancestry; their PRs and L2 SHAs are not modified, merged, approved or closed. The exact new-head CI and independent material integration review still must pass. This is a corrected **dependency boundary**, not a security threshold exemption. Future work must separately reconcile upstream foundation #36 and release candidates, and only merge when the real required ancestor gates and native approvals are satisfied.
