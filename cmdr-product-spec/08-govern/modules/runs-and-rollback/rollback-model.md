---
id: govern-runs-and-rollback-rollback-model
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-015, REQ-PROD-016, REQ-PROD-020, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
---
# Rollback Model — GOV-2 Supporting Contract

## Purpose

Support `CAP-GOV-030` and `CAP-GOV-031` with stable functional rollback distinctions. The capability files remain normative for complete inputs, outputs, actions, permissions and acceptance criteria.

## Mandatory separation

- rollback supported ≠ rollback eligible;
- rollback eligible ≠ rollback authorized;
- rollback available ≠ rollback safe;
- Rollback Plan ≠ rollback execution;
- rollback request ≠ rollback started;
- technical rollback completion ≠ verified rollback success;
- rollback success ≠ full recovery guaranteed;
- `rolled-back` ≠ exact original state restored;
- `recovered` ≠ exact restoration;
- compensation ≠ rollback;
- cancellation ≠ rollback;
- rollback failure ≠ automatic retry of the original action.

## Required rollback context

A rollback review preserves by reference:
- original Response Run, Decision, Execution Plan and Playbook/version;
- exact affected effects, targets and scope;
- rollback trigger and authority context;
- current target/readiness state;
- rollback capability owner/version;
- expected restored state and explicitly non-restored elements;
- potential data loss, limitations and residual risk;
- parameter/Secret References only, never raw secrets;
- verification criteria and recovery/manual fallback.

## Ownership

Govern owns rollback eligibility, Rollback Plan, Response Rollback governance, recovery coordination and canonical Result relations. Studio/Endpoint/provider owners retain technical reverse/recovery primitives and raw outputs. Platform Settings retains secrets, connections and runtime configuration. Shared Recovery remains a generic mechanism, not a competing response owner.

## Verification

Rollback is not considered successful solely because the technical executor returned success. GOV-2 must compare expected restored/recovered state with authorized verification observations and preserve partial, failed, inconclusive and residual-impact outcomes.

## Safety

Target drift, expired/missing authority, unknown original effect state, unsupported primitive, data-loss risk, unavailable executor/secret or missing verification can block effectful rollback or require re-decision/manual recovery. No command, bypass, API, protocol or provider/runtime is selected by this document.

## GOV-3 boundary

Rollback provenance and conceptual outcome inputs are emitted for future Audit Trail and Response Metrics. GOV-2 does not create either GOV-3 capability.