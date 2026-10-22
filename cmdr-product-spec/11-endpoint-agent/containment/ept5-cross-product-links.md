---
id: endpoint-ept5-cross-product-links
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# Endpoint EPT-5 — Cross-Product Links

Canonical flow: `Finding/context → Action Request → Policy/Authority → Approval as required → Decision → Response Run/Step → Endpoint Response Primitive → Technical Outcome → Endpoint Technical Verification Observation → Govern Verification/Reconciliation → Result`.

- Investigate owns Case/Finding/Evidence and Containment Request Preparation.
- Govern owns Action Request, Approval, Decision, Response Run, response verification, rollback/recovery governance and Result.
- Endpoint owns primitive eligibility/readiness, target-side effect facts, target-state observation, technical verification observation and technical reversal primitive.
- Studio Human Gate, Tool Call and Automation Run remain distinct references; OPEN-007 and OPEN-015 remain open.
- Settings owns Endpoint Policy, Fleet, tenant/environment, administrative configuration and Secret References.
- Shared owns generic Jobs/Trace/Activity/Reporting/Recovery mechanisms.

Mandatory non-equivalence: Human Gate != Approval/Decision; Response Run != Endpoint execution; Technical Outcome != Result; Endpoint Technical Verification != Govern Verification; Technical Reversal != Govern Response Rollback; Shared Recovery != Govern rollback semantics.