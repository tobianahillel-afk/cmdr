---
id: studio-std1-source-audit
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-09
source-of-truth: quality-report
---
# Studio STD-1 Source Audit

## Git and roadmap
Verified repository `tobianahillel-afk/cmdr`, canonical branch `docs/cmdr-product-spec-foundation`, PR #2 open/Draft/unmerged on `main`, preflight baseline `e0c23764df80a3d7109c156d1a2ee0962d19cda6`, and canonical roadmap id `roadmap-phase-5-studio-and-endpoint` before capability allocation.

## Governance and registers read
Master product/ownership/boundary sources, source-of-truth/documentation rules, capability/object/dependency/permission/screen registers, unresolved decisions, Requirements evidence, qualitative baseline, STATUS, roadmap namespace convention, Phase 5 preflight and canonical capability template were consulted. The `CAP-STD-*` namespace was confirmed empty of concrete IDs before allocation.

## Studio corpus read
The complete current Studio foundation corpus was inspected across Library, Builder, Skills, Workflows, Automation Agents, Agent Teams, Human Gates, Control Room, Assurance, Evaluations, Simulations, Versions & Deployment, navigation, permissions and product definition. Future STD-2/3/4 sources were used for boundary only.

The 11 active Studio screens were read: `STD-AGT-001`, `STD-ASR-001`, `STD-ATM-001`, `STD-BLD-001`, `STD-CTL-001`, `STD-DEP-001`, `STD-EVL-001`, `STD-LIB-001`, `STD-SIM-001`, `STD-SKL-001`, `STD-WFL-001`. No screen rewrite was performed.

## Adjacent owners read
Platform Settings provider/integration/secret/tenant/environment/health/fleet ownership; Govern execution boundaries; Shared Search/Trace/Jobs/Versioning/Reporting/Export/Notifications/Collaboration ownership; Investigate Tool/Studio boundary; Command Studio dependency projections; Endpoint target boundaries.

## Findings
1. 16-capability STD-1 set is supported without duplicate capability semantics.
2. Tool/Tool Call final object schemas remain future and are not created.
3. `perm.studio.*` and `perm.cmdr-studio.*` coexist; ambiguity is documented, not normalized.
4. OPEN-007/013/015 remain open where consumed; no OPEN created/closed.
5. Endpoint README missing information-architecture remains future Endpoint backlog and is untouched.
6. No STD-2/3/4 or Endpoint capability is introduced.
