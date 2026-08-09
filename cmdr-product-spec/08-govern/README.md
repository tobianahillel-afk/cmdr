---
id: 08-govern-readme
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-004, REQ-PROD-008, REQ-PROD-015, REQ-OBJ-005, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-015, OPEN-019]
---
# Govern

## Mission
Govern is CMDR's authority, governed-response and response-review product. GOV-1 governs Action Requests through Decision/Handoff; GOV-2 governs Playbooks, execution planning, Response Runs, verification, rollback/recovery and canonical Result; GOV-3 defines Govern audit interpretation, Govern-specific metrics, continuous-improvement packages and documentary closure without duplicating Shared infrastructure.

## Canonical programme identity
- Parent: **Delivery Roadmap Phase 4 — Govern**.
- Canonical id: `roadmap-phase-4-govern`.
- GOV-1 — Action Requests, Policy, Authorities and Decisions: historical **PASS**, 16 capabilities / 432 sections / 96 tables / 180 gates.
- GOV-2 — Playbooks, Response Runs, Execution, Verification and Rollback: historical **PASS**, 17 / 459 / 102 / 190 gates.
- Current construction: **GOV-3 — Audit Trail, Response Metrics and Govern Closure**.
- GOV-1/GOV-2/GOV-3 are execution lots only. `Phase 4C/4D/4E Govern` do not exist.

## Govern ownership after GOV-3 functional specification
Govern owns the semantics/lifecycle established by GOV-1/GOV-2 plus:
- Govern Audit Event interpretation and lifecycle reconstruction;
- audit completeness/gap/contradiction assessments and Audit Review;
- Audit Evidence Package composition;
- Govern-specific Policy/Exception/Emergency, Approval/Authority/SoD, Decision, Run, verification/rollback/recovery, Result/effectiveness and flow metric semantics;
- Trend and Control Health Assessments;
- Govern Continuous Improvement Package and closure provenance.

## Consumed, not owned
- Shared: Trace, Activity, Search, Metrics Engine, Reporting Engine, Export, Jobs, Notifications, Versioning and generic dashboard primitives.
- Platform Settings: users/roles/tenants/environments, retention/storage, export destinations, providers/integrations/secrets and administrative configuration.
- Security: permission model, tenant isolation, privacy, secure export, legal hold and audit-integrity policy.
- Command: Incident, Work Queue and Command KPIs.
- Investigate: Case, Evidence, Finding and investigation metrics.
- Studio: Workflow, Tool, Tool Call, Human Gate, Automation Run and Studio metrics.
- Endpoint/provider owners: technical execution primitives, local audit records and technical runtime metrics.

## Modules
1. Response Inbox — GOV-1.
2. Action Center — GOV-1.
3. Decision Register — GOV-1.
4. Policy Gates — GOV-1.
5. Approvals & Authorities — GOV-1.
6. Playbooks — GOV-2.
7. Runs & Rollback — GOV-2.
8. **Audit Trail — GOV-3 CAP-GOV-034..038.**
9. **Response Metrics — GOV-3 CAP-GOV-039..047.**

## Audit/metrics invariants
Audit Trail ≠ Trace/Activity. Govern Audit Event ≠ raw log/SIEM event. Reconstruction ≠ execution. Audit completeness ≠ truth completeness. Gap/contradiction ≠ wrongdoing/falsity automatically. Integrity requirement ≠ implemented cryptographic proof. Audit Evidence Package ≠ canonical Evidence.

Metric ≠ objective/Policy/SLO/KPI automatically. Count ≠ quality; throughput ≠ effectiveness; faster Decision ≠ better Decision; Policy block count ≠ prevented incidents; runtime success ≠ verified success; rollback rate ≠ failure rate; Result success ≠ business value. Trend ≠ causal explanation; anomaly ≠ control failure; dashboard ≠ source of truth.

## AI / sensitive data
AI is optional and proposal-only. It may summarize sourced audit chains/metrics or draft hypotheses/packages, but cannot invent events, alter history, declare fraud/violation as fact, mutate Decision/Result, hide contradictions, change thresholds, publish externally or apply an improvement. Raw secret material is excluded. Identity/tenant/target/exception/emergency dimensions remain permission-aware.

## Canonical GOV-3 addenda
- `audit-and-metrics-boundaries.md`
- `capability-map-gov3.md`
- `object-consumption-map-gov3.md`
- `action-classification-gov3.md`
- `automation-and-ai-model-gov3.md`
- `permissions-gov3.md`
- `cross-product-links-gov3.md`
- `screen-capability-map-gov3.md`
- `source-migration-gov3.md`

Historical GOV-1/GOV-2 maps remain valid evidence and are not destructively condensed.

## Current construction counts
- GOV-1: 16 / 432 / 96 — historical PASS.
- GOV-2: 17 / 459 / 102 — historical PASS.
- GOV-3 functional set: **14 / 378 / 84**.
- Govern cumulative: **47 / 1269 / 282**.
- Global: **317 capabilities / 315 defined / 2 proposed / 317 planned / 8559 sections / 1902 tables**.

At this pre-publication construction point, GOV-3 and parent closure remain pending the fifth functional commit plus remote 200-gate verification. Documentary PASS never means implementation complete and does not start Delivery Roadmap Phase 5.