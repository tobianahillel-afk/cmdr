---
id: roadmap-phase-4-govern
domain: 18-roadmap-and-releases
status: draft
owner: Product Operations Lead
updated: 2026-08-09
source-of-truth: canonical
---
# Phase 4 Govern

## Canonical phase identity

- namespace: **Delivery Roadmap**;
- canonical id: `roadmap-phase-4-govern`;
- canonical title: `Phase 4 Govern`;
- qualified title: **Delivery Roadmap Phase 4 — Govern**;
- previous delivery phase: `Delivery Roadmap Phase 3 — Investigate`;
- next delivery phase: `Delivery Roadmap Phase 5 — Studio and Endpoint`;
- Govern capability specification status: **PARTIAL**.

`Delivery Roadmap Phase 4 — Govern` and `Capability Specification Phase 4` belong to distinct namespaces. `Phase 4C Govern`, `Phase 4D Govern` and `Capability Specification Phase 4C` **do not exist**. GOV-1/GOV-2/GOV-3 are execution-lot identifiers, not roadmap phases.

## Capability-specification execution lots

| Execution lot | Scope | Status |
|---|---|---|
| **GOV-1** | Action Requests, Policy, Authorities and Decisions | **PASS AFTER POST-PUBLICATION VERIFICATION — 180/180** |
| **GOV-2** | Playbooks, Response Runs, Execution, Verification and Rollback | **PENDING POST-PUBLICATION VERIFICATION** |
| GOV-3 | Audit Trail, Response Metrics and Govern Closure | **NOT STARTED** |

## GOV-1 preserved closure

GOV-1 remains exactly `CAP-GOV-001..016`, 16 defined/planned capabilities, 432 numbered sections and 96 mandatory tables. Its dedicated post-publication record is based on `b8dd93e03443adb9101c7592094a48e358b460e2` (`docs: record Govern GOV-1 post-publication verification`). GOV-2 does not redefine Action Request, Policy/authority/Approval, Decision or Execution Handoff Package semantics.

## GOV-2 scope and chain

GOV-2 starts from a Decision plus Execution Handoff Package and specifies:

`Decision → Execution Handoff Package → Playbook Selection → exact-version Compatibility Review → Execution Plan → Target Resolution/Readiness → Authorization Reconciliation → Response Run → scheduling/control → Response Steps → Studio/Endpoint/provider handoff → Runtime Reconciliation → Error/Retry/Partial Success → Verification Plan → Post-Execution Verification/Residual Risk → Rollback/Recovery when required → canonical Result → cross-product handoff`.

GOV-2 primarily covers modules:
6. Playbooks;
7. Runs & Rollback.

Modules 1–5 remain GOV-1 and are only consumed/re-entered for review when authority changes. Modules 8–9 remain future GOV-3:
8. Audit Trail;
9. Response Metrics.

## GOV-2 capability set

Exactly `CAP-GOV-017..033` are allocated, all `draft` / `defined` / `planned`:
- 17 capabilities;
- 459 numbered sections;
- 102 mandatory S8/S9/S10/S13/S16/S17 tables;
- duplicate/recycled IDs: 0;
- GOV-3 capabilities created: 0.

Govern cumulative after GOV-2 functional specification: **33 capabilities / 891 sections / 198 mandatory tables**.

## Ownership and execution boundary

Govern owns Response Playbook semantics, Execution Plan, Response Run/Step governance, runtime reconciliation, verification, rollback/recovery governance and canonical Result. Studio retains Workflow/Tool/Tool Call/Human Gate/Automation Run. Endpoint/provider owners retain technical primitives/raw outcomes. Settings retains providers/integrations/secrets/credentials/runtime/tenant/environment administration. Command retains Incident/Work Queue; Investigate retains Case/Evidence/Finding/analysis; Shared retains generic Jobs/Trace/Activity/Versioning/Reporting/Recovery mechanisms.

Key distinctions: Playbook != Workflow; Response Run != Automation Run/Job/Tool Call; technical output != canonical Result; retry != reauthorization; cancel/compensation != rollback; runtime success != verification success; Result never rewrites Decision/Evidence/Finding.

## GOV-2 functional commits

Exact baseline:
- `b8dd93e03443adb9101c7592094a48e358b460e2` — `docs: record Govern GOV-1 post-publication verification`.

Published functional commits before the traceability/quality commit:
1. `f2981f5da45c0011390e3b4c9f21b596780758bb` — `docs: establish Govern playbook and execution boundaries`;
2. `9ff1ebb8fcadb5ea8cccef3ed7901c491d1627e7` — `docs: define Govern execution planning readiness and response runs`;
3. `8d109caea41867aaf74794fbcad14896b35987ca` — `docs: specify Govern runtime coordination verification and failure handling`;
4. `c2314c475a75cfc09122917cc72d282f216f2fbd` — `docs: document Govern rollback recovery results and provenance`;
5. `docs: update Govern execution traceability and quality gates` — exact SHA recorded only after publication.

## Requirements / OPEN / implementation

- Requirements remain **122 — 99 conform / 20 partial / 3 absent / 0 contradictory**; GOV-2 adds documentary evidence without global state promotion.
- Open decisions remain **18**; GOV-2 creates/closes 0. OPEN-007/008/013/015/019 remain relevant.
- GOV-2 defines no provider/runtime, command, exploit/bypass, API/protocol, raw secret value, complete object schema, final state machine, final RBAC/ABAC or product code.
- detailed screen rewrites: 0; new Screen IDs: 0.

## Acceptance for GOV-2

GOV-2 can be promoted to `PASS AFTER POST-PUBLICATION VERIFICATION` only after the fifth required functional commit is published and all **190** gates pass, including exact remote SHA/ancestry, PR/readme/main invariants, 17/459/102 structural conformance and GOV-1/Command/Investigate non-regression.

## Stop line

GOV-3 remains NOT STARTED. Govern capability specification and Delivery Roadmap Phase 4 — Govern remain PARTIAL even after successful GOV-2 closure.