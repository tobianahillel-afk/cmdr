# Traceability Matrix

## Purpose

Connect product outcomes to canonical documents and implementation evidence without copying requirements.

| Outcome | Canonical specification | Primary owners | Minimum evidence |
|---|---|---|---|
| Shared operational picture | `10-command-center/MISSION_CONTROL.md` | SOC Operations | Reconciled metrics, queue and freshness tests |
| Evidence-led investigation | `20-investigation-lab/CASE_WORKSPACE.md` and `EVIDENCE_BOARD.md` | DFIR | Provenance, integrity and collaboration tests |
| Complete forensic chain | Static, Sandbox, Reverse, Debugger, Memory and Disk pages | DFIR / Malware Research | Isolated execution and reproducibility tests |
| Controlled response | `30-response-governance/ACTION_CENTER.md` | Response Governance | Authority, policy, rationale and target-scope tests |
| Safe execution and rollback | `RUNS_AND_ROLLBACK.md` | Response / Platform | Idempotency, partial failure and rollback tests |
| Tenant isolation | `40-platform/TENANCY_AND_SCOPING.md` | Platform | Cross-tenant negative tests |
| Immutable accountability | `40-platform/AUDIT_AND_IMMUTABILITY.md` | Security Architecture | Completeness and integrity verification |
| Context-preserving hand-offs | `04-cross-console/TRANSITIONS.md` | Product Architecture | End-to-end transition tests |
