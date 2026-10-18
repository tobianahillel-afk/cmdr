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


## Reconciliation — consolidated PR #127 retargeted to first-party #120 (2026-10-19)

This additive handoff records GitHub reality without rewriting prior independently
reviewed source SHA or asserting a native approval. The engineering manifest above
was authored before the independently reviewed #125 cohort filter and #127
filtered catalog/lead navigation existed. Do not re-run its original #120/#124
integration plan as though GitHub had no later code.

- GitHub `main` remains `bc1ec59e5e79ccbaba291e2984b0eb7d5e54129c` (product-spec documents; no published production runtime).
- Canonical engineering integration unit: `E10-DFIR-780A-INTEGRATION`, under
  transversal epic `E10-DFIR-780A`, referencing CAP-INV-002 (Event Search),
  CAP-INV-004 (Event Inspection/Pivot), CAP-INV-110 (Case Timeline). The
  normalized Case owner/typed source remains first-party read-only lab code;
  source Event and Timeline Entry ownership remain Shared. Class-2 mutations
  remain blocked by OPEN-013.
- GitHub PR #120, original Wazuh rule+agent pair navigation, **frozen base SHA**
  `db97ebb3299de6f2e9c3c1558553fdf78174bb35`, is the direct base of
  the new *consolidated* GitHub PR #127. Native Git compare from that base
  to former PR #127 HEAD `11607f13c355027cb0de21ce9bbde287860bef44`
  verified **ahead 19, behind 0, 14 changed files**; source code includes
  both #126's Wazuh catalog integration and #125's exact cohort filtering,
  plus #127's source-authorized lead-pair navigation. No source is rewritten.
- The former intermediary PR #126 is still DRAFT at the **frozen original**
  `2859eb9236f50e49588985e81d292e7c28a0b5e4`. Its own original
  `engineering-kernel` run #39175059577 is **FAILURE** due to mandatory
  changed-security coverage **88.87% <90.00%**; never transfer the green
  checks of an unrelated diff or silently waive that failed gate. PR #125,
  #124, #123, #122, #113, #36 and #2 remain independently reviewable and
  unmerged. #120 itself has prior green independent exact-SHA checks.
- PR #127's original smaller-diff HEAD `11607f13c355027cb0de21ce9bbde287860bef44`
  had GitHub Actions `investigation-timeline-runtime` run #39175997382
  **SUCCESS** and `engineering-kernel` run #39175997394 **SUCCESS**.
  **Those pre-retarget runs were assessed against parent PR #126, not
  necessarily against PR #120.** This handoff-only commit deliberately
  triggers fresh `pull_request.synchronize` workflows on the newly
  retargeted PR #127 and its entire actual 19-commit delta. Report their
  exact new HEAD conclusions before any claim of merge-readiness.
- Candidate remains **DRAFT**, unmerged and not production-accepted.
  Real runtime outcomes are operator-local synthetic saved Case investigation
  and authorized Wazuh metadata/cohort navigation, NOT continuous Wazuh
  collection, production IAM, signed custody, product CAP acceptance or
  deployment. Any prior ChatGPT L2 of #126 and #127 reviewed separate,
  smaller diffs; a changed merge base requires proportionate combined-diff
  review but does not invalidate immutable original commit evidence.
- No further branch stacking is justified merely to re-document these same
  features. On successful exact-HEAD mandatory checks, reconcile the READY
  execution cache to an evidence-based status without claiming VERIFIED or
  accepted before reviews/native approvals, then select a canonical
  READY *functional* obligation with its own correct engineering parent.
  On a security check failure, keep DRAFT, identify the exact affected
  security-critical scope and fix genuine behavior/design; do NOT lower the
  mandatory threshold or create near-identical tests solely for coverage.

This reconciliation is additional history, not a retroactive replacement of
the earlier `E10-DFIR-780A` manifests or product requirements.


## Canonical continuation checkpoint — 2026-10-19, after new-base CI

The direct-base integration is now supported by **fresh exact-head GitHub Actions
on the retargeted PR #127**. On the 20-commit, 14-file diff from immutable
PR #120 HEAD `db97ebb3299de6f2e9c3c1558553fdf78174bb35` to
candidate `a72ecbbe77c5972d9b79d995751061470348310b`:

- `investigation-timeline-runtime` run [39177579100](https://github.com/tobianahillel-afk/cmdr/actions/runs/39177579100):
  **COMPLETED SUCCESS**, Linux Go/Node, macOS and Windows.
- `engineering-kernel` run [39177579094](https://github.com/tobianahillel-afk/cmdr/actions/runs/39177579094):
  **COMPLETED SUCCESS**, mandatory `CHK-SECURITY-TESTS` PASS; report
  `coverage_status=measured global_coverage=90.76 global_floor=80
  changed_floor=90 changed_scopes=1 auth_required=8 tenant_required=8`.
  The mandatory floor was NOT changed or waived.
- Historical PR #126's own 88.87% red kernel status REMAINS recorded;
  this is a different, narrower integration-base/diff with its own real
  successful check, not a retroactive green result for #126.
- Code, Wazuh correlation algorithms, typed-source ACLs and L2-frozen
  source commit contents are unchanged. Only PR base and recovery records
  changed. The current integration work unit is therefore **IMPLEMENTED**
  as an auditable DRAFT release candidate, deliberately NOT
  `VERIFIED`, `ACCEPTED`, native-approved, merged or production ready.
- This recovery metadata commit itself is the NEXT head. Its own
  `pull_request.synchronize` mandatory checks must also be checked at
  the exact new SHA; don't transfer the prior run status to it.
- Distinct former ChatGPT L2 work reviewed #126's original base..head,
  #127's earlier 5-file diff and source lead navigation. Results were
  not available in the author conversation; a single independently
  scoped combined-diff review is needed for the new base before
  advancing beyond `IMPLEMENTED`. Do NOT create redundant reviews of
  each immutable parent commit.
- The work graph no longer falsely advertises the completed integration
  as `READY`; the execution cache deliberately has NO preferred READY
  work unit until canonical planning chooses a genuinely admissible
  next behavior and creates its scoped manifest, instead of recycling
  E10-INV-120A-SNAPSHOT or growing this transversal release unit.
- Keep the prior #36 engineering foundation and #2 product-spec PRs
  separate; downstream merge rights and required native approvals
  must be checked independently. No fusion requested or performed.
