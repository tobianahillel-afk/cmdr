---
id: endpoint-ept5-source-audit
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-11
source-of-truth: quality-report
requirements: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-006, REQ-PROD-008, REQ-PROD-014, REQ-PROD-015, REQ-PROD-016, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-014, OPEN-015, OPEN-017]
---
# Endpoint EPT-5 — Source Audit

## Scope
Read-before-write audit for **EPT-5 — Containment, Verification and Governed Response Primitives** from exact EPT-4 final baseline `5d576295fa12693ef375a35cfe515d7bdf577f68`. EPT-6 is excluded.

## Endpoint Containment — 7/7 read
1. `containment/README.md` — containment under Govern authority; scope/impact; rollback; verification.
2. `containment/account-containment.md` — local session termination or lock where supported; directory actions explicitly external; scope and rollback.
3. `containment/host-containment.md` — composite host containment, business impact, dependencies and verification.
4. `containment/network-isolation.md` — control-plane allowance concept, policy exceptions, local firewall concept and connectivity verification without implementation syntax.
5. `containment/quarantine.md` — secure-store concept, original path/hash, access restriction and restore.
6. `containment/rollback.md` — eligibility, reverse action and partial rollback; historical `Result` wording is reconciled to Govern-owned canonical Result.
7. `containment/verification.md` — independent checks, telemetry confirmation and residual access; historical `Result object` wording is reconciled to Govern verification/Result ownership.

## Effectful Live Response sources re-read
- `live-response/process-actions.md`: terminate/suspend/resume, tree scope, impact preview and verification.
- `live-response/network-actions.md`: bounded temporary block, scope/expiry/conflict and verification.
- `live-response/file-actions.md`: move/copy/delete/restore, quarantine integration, path/hash and rollback limits.
- `live-response/service-actions.md`: start/stop/restart/disable-like semantics, dependency impact, platform scope and verification.
EPT-4 command/session/output contracts are reused rather than recreated.

## Govern sources revalidated
Action Request, Approval, Decision, Response Run, Response Rollback and Result canonical objects were re-read. GOV-1 request/policy/authority/approval/Decision/handoff contracts and GOV-2 execution handoff, raw technical outcome reconciliation, Verification Plan, post-execution verification, rollback eligibility/execution and Result contracts were revalidated.

The canonical chain remains `Finding/context → Action Request → policy/authority → Decision → Response Run → Endpoint technical primitive → technical status/output → Endpoint technical verification observation → Govern reconciliation/verification → Result`.

## Cross-product boundaries
- Investigate owns Case/Finding/Evidence and Containment Request Preparation; it prepares but does not execute containment.
- Govern owns production-response authority, Approval, Decision, Response Run, response-level verification, rollback governance and canonical Result.
- Endpoint owns technical primitive definitions, local eligibility/availability/precheck, target-side execution facts, target-state observations, technical reversal primitives and local provenance.
- Studio owns Human Gate, Tool Call and Automation Run; `accepted-for-workflow` never becomes Govern authority.
- Settings owns Endpoint Policy, Fleet, tenant/environment administration and Secret References.
- Shared owns generic Jobs, Trace/Activity, Reporting and Recovery mechanisms, not response authority.

## Security and EPT-6 boundary
Security requires least privilege, step-up/authority controls, no hidden delegation and append-oriented local audit/provenance. OPEN-008 keeps supported platforms/releases unresolved. EPT-6 owns updates, deep resilience, anti-tamper, self-protection, watchdog/persistence and update rollback; none is introduced here.

## OPEN decisions
The programme remains at **18 OPEN decisions**. OPEN-007, OPEN-008, OPEN-013, OPEN-014, OPEN-015 and OPEN-017 remain open. EPT-5 closes none.

## Capability-set decision
The corpus independently supports the 16 candidate capabilities `CAP-EPT-065..080`. In addition, `containment/account-containment.md` independently sources **local session termination/lock** while declaring directory actions external. That responsibility is not equivalent to process control, service mutation, host isolation or directory identity administration, so EPT-5 adds `CAP-EPT-081 — Local Account Session Lock and Termination Primitive`.

Final EPT-5 set: **17 capabilities `CAP-EPT-065..081`**. The extra capability is source-driven, not quota-driven.

## Expected structure and recalculated metrics
EPT-5 target: **17 capabilities / 459 numbered sections / 102 mandatory tables / at least 51 GWT**.
Endpoint cumulative after EPT-5: **81 / 2187 / 486**.
Global after EPT-5: **466 capabilities / 464 defined / 2 proposed / 466 planned / 12582 sections / 2796 mandatory tables**.

## Explicit exclusions
No directory account disable/reset, Endpoint update/upgrade, self-update, anti-tamper, self-protection, watchdog, persistence, deep crash/offline resilience, API/protocol, native command syntax, PowerShell/shell, firewall/process-kill syntax, final RBAC/ABAC, final Policy/Approval/Verification/Rollback engine, physical schema, Screen ID, detailed UX or implementation.