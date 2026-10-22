---
id: endpoint-ept6-updates-resilience-security-provenance-capability-conformance
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-12
source-of-truth: quality-report
requirements: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-006, REQ-PROD-008, REQ-PROD-012, REQ-PROD-014, REQ-PROD-015, REQ-PROD-016, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-OBJ-007, REQ-OBJ-008, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004, REQ-SEC-005]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-014, OPEN-015, OPEN-017]
---
# Endpoint EPT-6 — Updates, Resilience, Security and Provenance Capability Conformance

## Build-time scope
This report validates only **EPT-6 — Updates, Resilience, Security and Endpoint Provenance** before publication-dependent checks. It does not start Phase 6 and does not claim implementation.

## Exact baseline
`4d20e8e411aa8919f31e773a9c5ea96efd64791e` — `docs: record Endpoint EPT-5 post-publication verification`.

## Source-driven set
Exactly **18** new capability contracts are authored: `CAP-EPT-082..099`.

Families:
- Update: `082..088` = 7;
- Resilience: `089..093` = 5;
- Security/provenance: `094..098` = 5;
- Endpoint closure: `099` = 1.

The set is source-driven. Anti-tamper, privilege/security context, Secret Reference handling and local audit are independently sourced; reconnect/replay is distinct from buffer retention; resource pressure is distinct from dependency recovery.

## Structural conformance
- capability files: **18/18**;
- numbered sections: **486/486** (27 each);
- mandatory tables: **108/108** (S8, S9, S10, S13, S16, S17 per capability);
- Given/When/Then scenarios: **at least 54**;
- duplicate EPT-6 IDs: **0**;
- recycled IDs: **0**;
- owner conflicts: **0 identified**;
- empty/generic mandatory tables: **0**;
- placeholder capability-layer files: **0**;
- Endpoint Screen IDs created: **0**.

## Counts after EPT-6 content
Endpoint final content: **99 capabilities / 2673 sections / 594 mandatory tables**.

Global content: **484 capabilities / 482 defined / 2 proposed / 484 planned / 13068 sections / 2904 mandatory tables**.

Requirements remain conservatively **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**; EPT-6 adds evidence but does not force a global status promotion without row-level canonical proof. OPEN remains **18**.

## Ownership validation
PASS at build-time for the required boundaries:
- Settings owns admin update target/channel/wave/Fleet/Policy and Secret administration;
- Endpoint owns local update execution/state, buffering/replay/restart/dependency/self-protection/security/audit facts;
- Studio owns Studio asset deployment/reversion;
- Govern owns response authority, response verification/rollback and Result;
- Shared owns generic Jobs/Retry/Recovery/Trace/Activity;
- Security owns global permissions/privacy/audit-integrity/tenant/secret policy.

## Mandatory distinction validation
The capability set explicitly preserves assignment/download/stage/readiness/install/activate/health separation; compatibility/support separation; retry/recovery separation; Endpoint/Studio/Govern rollback separation; buffering/delivery/replay/consumption separation; restart/health/sync separation; degraded/failure separation; self-protection/security/integrity/compromise separation; privilege/permission separation; Local Audit/Shared Trace/immutable-proof separation; Secret Reference/raw-secret separation; and Endpoint security fact/Finding/Incident/Result separation.

## Out-of-scope validation
No API/protocol/port, package/signature/certificate/crypto format, updater/package-manager/CDN, physical persistence/queue schema, watchdog/service/driver design, final anti-tamper mechanism, final RBAC/ABAC, Fleet scheduler/update orchestrator/storage engine/immutable ledger/crypto proof, Endpoint Screen or Phase 6 capability is introduced.

## 240-gate build-time state
Gates **1–233: PASS** based on baseline, source/namespace audit, structural content, ownership/distinctions, functional coverage, AI/security limits, registers/maps/closure artifacts prepared in the fifth functional commit.

Publication-dependent gates **234–240 remain PENDING** at build time:
- 234 five functional commits reachable on remote;
- 235 remote verification executed;
- 236 exact remote build SHA verified;
- 237 exact final verification-record SHA;
- 238 post-publication evidence;
- 239 PR/main/README remote recheck;
- 240 final Endpoint + Phase 5 status after remote evidence.

Build-time verdict: **233 PASS / 7 PENDING-REMOTE / 0 FAIL**.

No `240/240` claim is valid before the five functional commits are fast-forward published and remotely re-read.

## Build-time closure implication
EPT-6 content is structurally complete, but **EPT-6 remains PENDING POST-PUBLICATION VERIFICATION**. Endpoint Capability Specification and Delivery Roadmap Phase 5 remain **PARTIAL/PENDING** until gates 234–240 close remotely.