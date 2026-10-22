---
id: endpoint-ept3-source-audit
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-11
source-of-truth: quality-report
---
# Endpoint EPT-3 — Source Audit

## Baseline
Exact starting HEAD: `5d7c037aff6004984416665e7e188a8700e62b2f` — `docs: record Endpoint EPT-2 post-publication verification`. PR #2 was open/Draft/unmerged on `main`; root README branch/main remained exact `# cmdr`; repository public; auto-merge disabled.

## Endpoint Detection corpus — 8/8
All eight current files are active `draft / canonical`, owner Endpoint Agent Product Lead: `README.md`, `behavioral-detection.md`, `detection-suppression.md`, `detection-update.md`, `local-correlation.md`, `local-detection-engine.md`, `model-inference.md`, `rule-model.md`.

The corpus supports local rule evaluation, behavioral correlation, governed model hooks, sequence/windowing, severity/confidence, suppression scope/expiry/owner/audit, cross-sensor correlation, partial telemetry and result provenance. `detection-update.md` is retained as a future update/rollout source and is not implemented by EPT-3.

## Endpoint Investigation corpus — 8/8
All eight current files are active `draft / canonical`, owner Endpoint Agent Product Lead: `README.md`, `host-inspection.md`, `host-timeline.md`, `modules-and-drivers.md`, `network-connections.md`, `persistence.md`, `process-tree.md`, `user-sessions.md`.

The corpus supports read-only host context, process ancestry, network/process/user linkage, user/session privacy, service/module/driver/system context, persistence-like observations, local sequence/clock health/gaps and explicit escalation to future Collection rather than acquisition inside EPT-3.

## EPT-1/EPT-2 non-regression
EPT-1 remains 190/190 PASS with `CAP-EPT-001..014`. EPT-2 remains 200/200 PASS with `CAP-EPT-015..030`. EPT-3 consumes telemetry/source/normalization/quality/capability-declaration contracts and does not redefine them.

## Cross-product ownership
- Investigate owns Case, Evidence, Finding and Detection Engineering authoring/validation/lifecycle.
- Command owns canonical Detection, Signal, Alert and Incident.
- Endpoint owns local technical evaluation, match, local signal-candidate/context and local investigation projections only.
- Settings owns administrative sources, parsers, configured runtimes/targets, Endpoint Policy, assignments, environments and rollout configuration.
- Govern owns Action Request, Decision, Approval, Response Run, Result and production response authority.
- Shared owns generic Search, Timeline, Linking, Correlation, Trace, Activity, Jobs and Reporting mechanisms.
- Studio owns Tool, Tool Call, Skill, Workflow and Automation Run.

## OPEN decisions
The programme retains 18 OPEN decisions. `OPEN-008` remains open for platform/source availability and support. `OPEN-017` remains open for Detection runtime, target language and portability. EPT-3 closes neither and selects no final rule/model/runtime/engine.

## Namespace
Current capability register ends at `CAP-EPT-030`. Searches for `CAP-EPT-031`, `CAP-EPT-046` and reserved CAP-EPT ranges returned no existing allocation/reservation. The candidate range `CAP-EPT-031..046` is therefore available subject to this audit.

## Migration assessment
No Detection or Investigation source is deprecated, duplicated as a canonical capability layer, or safe to delete. The legacy files are concise module sources rather than competing capability contracts. EPT-3 therefore uses an additive migration strategy: preserve all 16 source files and make `CAP-EPT-031..046` the normative capability layer while retaining source links and ownership.

## Set justification
All 16 candidate capabilities are independently justified. In particular, suppression has a dedicated Detection source; user/session and system/module/driver context have dedicated Investigation sources; Timeline/Correlation and Pivot/Handoff address distinct functional contracts. No artificial capability is required to retain the 16-count set.

## Stop boundary
EPT-4 Collection/Live Response, EPT-5 containment/governed response and EPT-6 update/resilience/security remain NOT STARTED. No acquisition, Live Response, containment, response execution, API/protocol, physical schema, final RBAC or product code is introduced by this audit.