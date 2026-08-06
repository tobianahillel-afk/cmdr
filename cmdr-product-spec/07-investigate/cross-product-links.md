---
id: investigate-cross-product-links
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements:
  - REQ-PROD-004
  - REQ-PROD-005
  - REQ-PROD-008
  - REQ-INV-006
open_decisions:
  - OPEN-008
  - OPEN-013
  - OPEN-015
  - OPEN-017
---
# Cross-product links — Detection Engineering

| Transition | Context | Ownership / return |
|---|---|---|
| Finding / Case Hypothesis / Hunt / Incident / technical analysis → Intake | sources, observations, Evidence, contradictions, limits and return origin | source owner retained; Intake Investigate |
| CAP-INV-401..416 → Review Package | immutable version, validations, tests, replay, match review, coverage, gaps and risks | Investigate |
| Review Package → Release Candidate / Review / Readiness | candidate version, reviewer, targets, dependencies and rollback/observation plan | Investigate; no Approval |
| Readiness / Promotion / Canary / Suppression / Exception / Rollback / Retirement proposal → Govern | class, targets, risks, conditions, stop/rollback and evidence | Govern owns Action Request, Decision, Approval, Response Run and Result |
| Promotion/Deployment coordination → Settings/Endpoint/runtime | approved target/version context and Run/Result projections | Settings/Endpoint/runtime execute/administer; Investigate observes |
| Runtime Detection → Command | version, state, Signals, Alerts, Incidents and dispositions | Command owns runtime and operational objects |
| Command feedback → Runtime Quality / Production Match Review | dispositions, comments, value, noise, outcomes and limitations | Investigate assessment; Command source retained |
| Drift/Performance/Health → Settings or new Draft | changed dependency, metrics, impact and tests | request or new authoring cycle; no silent runtime change |
| Rollback/Retirement Result → Version/Health/Coverage | per-target state, residual risk and gaps | source owner retained |
| Lifecycle Provenance → Continuous Improvement | complete lineage and unresolved items | CAP-INV-401/402/403/406, Settings request or Reporting |
| Future Intelligence | candidate handoff only | 4B.3B not started; no canonical Intelligence object |

Tenant, environment, target, immutable version, authority, permissions, masking, errors, partiality, timestamps and return origin are preserved.
