---
id: endpoint-capability-specification-closure
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-12
source-of-truth: quality-report
---
# Endpoint Capability Specification Closure

## Scope
Evaluate capability-layer completeness for the Endpoint Agent across EPT-1 through EPT-6. This document distinguishes content closure from publication verification and from software implementation.

## Lot status entering EPT-6
- EPT-1: PASS AFTER POST-PUBLICATION VERIFICATION — 190/190;
- EPT-2: PASS — 200/200;
- EPT-3: PASS — 210/210;
- EPT-4: PASS — 220/220;
- EPT-5: PASS — 230/230;
- EPT-6: build content complete, post-publication verification pending.

## Final content inventory
- EPT-1 `CAP-EPT-001..014`: 14 / 378 / 84;
- EPT-2 `CAP-EPT-015..030`: 16 / 432 / 96;
- EPT-3 `CAP-EPT-031..046`: 16 / 432 / 96;
- EPT-4 `CAP-EPT-047..064`: 18 / 486 / 108;
- EPT-5 `CAP-EPT-065..081`: 17 / 459 / 102;
- EPT-6 `CAP-EPT-082..099`: 18 / 486 / 108.

Endpoint total: **99 capabilities / 2673 sections / 594 mandatory tables**.

## Mandatory-family completeness
The capability audit finds all required Endpoint families represented:
1. identity/enrollment/platform/version/inventory/health;
2. telemetry/observations/capability declaration;
3. local detection/investigation;
4. collection/live response technical execution;
5. containment/technical verification/governed response primitives;
6. Agent update lifecycle;
7. offline buffering/replay/restart/crash recovery;
8. resource/dependency degradation and recovery;
9. self-protection/anti-tamper and local security context;
10. sensitive-material/Secret Reference handling;
11. Local Audit Event and Endpoint provenance;
12. cross-product security-state handoff and overall capability closure.

No mandatory capability family is missing in the content layer. No blocking owner conflict, recycled ID, placeholder capability or competing canonical capability source was identified in the EPT-6 source/migration audit.

## Ownership closure
Settings retains Fleet/update administration, policy and Secret administration. Studio retains Studio deployment lifecycle. Govern retains response authority/rollback/Result. Shared retains generic mechanisms. Security retains global policy. Endpoint retains individual local technical facts and operations only.

## OPEN decisions
OPEN count remains **18**. OPEN-008/013/015/017 and other historical OPEN decisions remain explicit. OPEN-008 does not block documentary Endpoint closure because EPT contracts remain platform-neutral and make no delivered-platform claim.

## Screens / implementation
Endpoint Screen IDs remain **0**. Capability closure introduces no code/API/protocol/final schema/final RBAC/physical engine or supported-platform promise.

## Build-time verdict
**CONTENT CLOSURE POSITIVE / FINAL VERDICT PENDING REMOTE EPT-6 VERIFICATION.**

Endpoint Capability Specification must remain **PARTIAL/PENDING** until EPT-6 reaches final **240/240** after remote publication. If and only if all EPT-6 publication gates pass with no regression, the final documentary status becomes **PASS**.

`Endpoint Capability Specification PASS != Endpoint implementation complete`.