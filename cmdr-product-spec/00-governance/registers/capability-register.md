---
id: capability-register
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-13
source-of-truth: registry
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-013
  - REQ-PROD-014
  - REQ-PROD-015
  - REQ-PROD-019
---
# Capability Register

The register is split into active shards. IDs are immutable and never recycled. The global file is the authoritative index and total; each shard carries the row-level evidence needed for owner, role, object, delivery and traceability checks.

## Required row semantics
An active shard entry identifies Capability ID, name, owner product/module, documentary status, `delivery_status`, `delivery_mode`, canonical file, primary roles, primary objects, consumers, Requirement IDs, OPEN decisions, dependencies, supersession and review date. A shard must not use the registry to claim runtime availability unsupported by implementation/release evidence.

Registry validation checks duplicate/recycled IDs, concurrent owners, missing canonical files, capabilities without users/objects/requirements and `planned` capabilities presented as available. Detailed inputs, outputs, actions, states, permissions and acceptance criteria remain in canonical capability files.

| Shard | Scope | Count | Defined | Proposed | Delivery mode |
|---|---|---:|---:|---:|---|
| `capability-register-command.md` | Command | 27 | 26 | 1 | 27 planned |
| `capability-register-investigate-foundation.md` | Investigate CAP-INV-001..114 | 22 | 21 | 1 | 22 planned |
| `capability-register-investigate-collection.md` | CAP-INV-201..215 | 15 | 15 | 0 | 15 planned |
| `capability-register-investigate-analysis-workbench.md` | CAP-INV-301..313 | 13 | 13 | 0 | 13 planned |
| `capability-register-investigate-dynamic-sandbox.md` | CAP-INV-314..328 | 15 | 15 | 0 | 15 planned |
| `capability-register-investigate-reverse-debugger.md` | CAP-INV-329..346 | 18 | 18 | 0 | 18 planned |
| `capability-register-investigate-memory-forensics.md` | CAP-INV-347..362 | 16 | 16 | 0 | 16 planned |
| `capability-register-investigate-disk-filesystem-forensics.md` | CAP-INV-363..379 | 17 | 17 | 0 | 17 planned |
| `capability-register-investigate-network-forensics.md` | CAP-INV-380..397 | 18 | 18 | 0 | 18 planned |
| `capability-register-investigate-detection-authoring.md` | CAP-INV-401..417 | 17 | 17 | 0 | 17 planned |
| `capability-register-investigate-detection-lifecycle.md` | CAP-INV-418..435 | 18 | 18 | 0 | 18 planned |
| `capability-register-investigate-threat-intelligence-foundations.md` | CAP-INV-501..518 | 18 | 18 | 0 | 18 planned |
| `capability-register-investigate-threat-intelligence-analysis-and-products.md` | CAP-INV-519..537 | 19 | 19 | 0 | 19 planned |
| `capability-register-investigate-cloud-analysis.md` | CAP-INV-601..618 | 18 | 18 | 0 | 18 planned |
| `capability-register-investigate-mobile-forensics.md` | CAP-INV-701..719 | 19 | 19 | 0 | 19 planned |
| `capability-register-govern-gov1.md` | Govern GOV-1 CAP-GOV-001..016 | 16 | 16 | 0 | 16 planned |
| `capability-register-govern-gov2.md` | Govern GOV-2 CAP-GOV-017..033 | 17 | 17 | 0 | 17 planned |
| `capability-register-govern-gov3.md` | Govern GOV-3 CAP-GOV-034..047 | 14 | 14 | 0 | 14 planned |
| `capability-register-studio-std1.md` | Studio STD-1 CAP-STD-001..016 | 16 | 16 | 0 | 16 planned |
| `capability-register-studio-std2.md` | Studio STD-2 CAP-STD-017..033 | 17 | 17 | 0 | 17 planned |
| `capability-register-studio-std3.md` | Studio STD-3 CAP-STD-034..051 | 18 | 18 | 0 | 18 planned |
| `capability-register-studio-std4.md` | Studio STD-4 CAP-STD-052..068 | 17 | 17 | 0 | 17 planned |
| `capability-register-endpoint-ept1.md` | Endpoint EPT-1 CAP-EPT-001..014 | 14 | 14 | 0 | 14 planned |
| `capability-register-endpoint-ept2.md` | Endpoint EPT-2 CAP-EPT-015..030 | 16 | 16 | 0 | 16 planned |
| `capability-register-endpoint-ept3.md` | Endpoint EPT-3 CAP-EPT-031..046 | 16 | 16 | 0 | 16 planned |
| `capability-register-endpoint-ept4.md` | Endpoint EPT-4 CAP-EPT-047..064 | 18 | 18 | 0 | 18 planned |
| `capability-register-endpoint-ept5.md` | Endpoint EPT-5 CAP-EPT-065..081 | 17 | 17 | 0 | 17 planned |
| `capability-register-endpoint-ept6.md` | Endpoint EPT-6 CAP-EPT-082..099 | 18 | 18 | 0 | 18 planned |
| `capability-register-settings-tenant-environment-foundations.md` | Settings Tenant/Environment CAP-SET-001..004 | 4 | 4 | 0 | 4 planned |
| `capability-register-settings-identity-administration.md` | Settings Identity CAP-SET-005..007 | 3 | 3 | 0 | 3 planned |
| `capability-register-settings-secrets-and-connections.md` | Settings Secrets & Connections CAP-SET-008..009 | 2 | 2 | 0 | 2 planned |
| `capability-register-settings-models-and-providers.md` | Settings Models & Providers CAP-SET-010..011 | 2 | 2 | 0 | 2 planned |
| **Total** | **All registered capabilities** | **495** | **493** | **2** | **495 planned** |

## Product totals
- Command: **27** capabilities, 729 sections, 162 mandatory tables.
- Investigate: **243** capabilities, 6561 sections, 1458 mandatory tables.
- Govern GOV-1: **16** capabilities, 432 sections, 96 mandatory tables.
- Govern GOV-2: **17** capabilities, 459 sections, 102 mandatory tables.
- Govern GOV-3: **14** capabilities, 378 sections, 84 mandatory tables.
- Govern cumulative: **47** capabilities, **1269** sections, **282** tables.
- Studio STD-1: **16** capabilities, **432** sections, **96** mandatory tables.
- Studio STD-2: **17** capabilities, **459** sections, **102** mandatory tables.
- Studio STD-3: **18** capabilities, **486** sections, **108** mandatory tables.
- Studio STD-4: **17** capabilities, **459** sections, **102** mandatory tables.
- Studio cumulative: **68** capabilities, **1836** sections, **408** mandatory tables.
- Endpoint EPT-1: **14** capabilities, **378** sections, **84** mandatory tables.
- Endpoint EPT-2: **16** capabilities, **432** sections, **96** mandatory tables.
- Endpoint EPT-3: **16** capabilities, **432** sections, **96** mandatory tables.
- Endpoint EPT-4: **18** capabilities, **486** sections, **108** mandatory tables.
- Endpoint EPT-5: **17** capabilities, **459** sections, **102** mandatory tables.
- Endpoint EPT-6: **18** capabilities, **486** sections, **108** mandatory tables.
- Endpoint cumulative: **99 capabilities / 2673 sections / 594 mandatory tables**; EPT-1/EPT-2/EPT-3/EPT-4/EPT-5/EPT-6 verified PASS after their respective post-publication verification records.
- Settings Tenant/Environment foundations: **4 capabilities / 108 sections / 24 mandatory tables**; **PASS AFTER POST-PUBLICATION VERIFICATION — 160/160**.
- Settings Identity Administration: **3 capabilities / 81 sections / 18 mandatory tables / 12 meaningful GWT**; **PASS AFTER POST-PUBLICATION VERIFICATION — 174/174**.
- Settings Secrets & Connections: **2 capabilities / 54 sections / 12 mandatory tables / 12 meaningful GWT**; **PASS AFTER POST-PUBLICATION VERIFICATION — 172/172**.
- Settings Models & Providers functional build: **2 capabilities / 54 sections / 12 mandatory tables / 13 meaningful GWT**; build-time **198 PASS / 6 PENDING-REMOTE / 0 FAIL**.
- Settings cumulative current build content: **11 capabilities / 297 sections / 66 mandatory tables**.
- CAP-INV-3xx / 4xx / 5xx / 6xx / 7xx: **97 / 35 / 37 / 18 / 19**.
- Detection Engineering: **35 capabilities, 945 sections, 210 tables**.
- Threat Intelligence: **37 capabilities, 999 sections, 222 tables**.
- Cloud Analysis: **18 capabilities, 486 sections, 108 tables**.
- Mobile Forensics: **19 capabilities, 513 sections, 114 mandatory tables**.
- Phase 4B.4 Cloud + Mobile: **37 capabilities, 999 sections and 222 tables**.
- Command + Investigate: **270 capabilities, 7290 sections, 1620 tables**.
- Command + Investigate + Govern: **317 capabilities, 8559 sections, 1902 mandatory tables**.
- Command + Investigate + Govern + Studio: **385 capabilities, 10395 sections, 2310 mandatory tables**.
- Command + Investigate + Govern + Studio + Endpoint EPT-1: **399 capabilities, 10773 sections, 2394 mandatory tables**.
- Command + Investigate + Govern + Studio + Endpoint EPT-1/EPT-2: **415 capabilities, 11205 sections, 2490 mandatory tables**.
- Command + Investigate + Govern + Studio + Endpoint EPT-1/EPT-2/EPT-3: **431 capabilities, 11637 sections, 2586 mandatory tables**.
- Command + Investigate + Govern + Studio + Endpoint EPT-1/EPT-2/EPT-3/EPT-4: **449 capabilities, 12123 sections, 2694 mandatory tables**.
- Command + Investigate + Govern + Studio + Endpoint EPT-1/EPT-2/EPT-3/EPT-4/EPT-5: **466 capabilities, 12582 sections, 2796 mandatory tables**.
- Command + Investigate + Govern + Studio + Endpoint EPT-1/EPT-2/EPT-3/EPT-4/EPT-5/EPT-6: **484 capabilities, 13068 sections, 2904 mandatory tables**.
- Historical current state after first Settings lot: **488 capabilities, 13176 sections, 2928 mandatory tables**.
- Historical Identity Administration build/closure state: **491 capabilities, 13257 sections, 2946 mandatory tables**.
- Historical Secrets & Connections build/closure state: **493 capabilities, 13311 sections, 2958 mandatory tables**.
- Current global functional build content including Settings Models & Providers: **495 capabilities / 493 defined / 2 proposed / 495 planned / 13365 sections / 2970 mandatory tables**.

## Command registry state
The Command shard remains 27 unique IDs, 26 defined / 1 proposed / 27 planned, with current native/integrated claims equal to 0. `CAP-CMD-401` remains deployment-dependent under OPEN-006.

## Govern registry state
- GOV-1: 16 unique `CAP-GOV-001..016`, 16 defined/planned, 432 sections, 96 tables; historical 180/180 PASS evidence retained.
- GOV-2: 17 unique `CAP-GOV-017..033`, 17 defined/planned, 459 sections, 102 tables; historical 190/190 PASS evidence retained.
- GOV-3: 14 unique `CAP-GOV-034..047`, 14 defined/planned, 378 sections, 84 tables; **200/200 post-publication PASS evidence recorded**.
- Govern capability specification: **PASS** across 47 capabilities / 1269 sections / 282 tables.

## Studio registry state
- STD-1: 16 unique `CAP-STD-001..016`, all `draft / defined / planned`, 432 sections, 96 mandatory tables; **PASS AFTER POST-PUBLICATION VERIFICATION — 190/190**.
- STD-2: 17 unique `CAP-STD-017..033`, all `draft / defined / planned`, 459 sections, 102 mandatory tables; **PASS AFTER POST-PUBLICATION VERIFICATION — 200/200**.
- STD-3: 18 unique `CAP-STD-034..051`, all `draft / defined / planned`, 486 sections, 108 mandatory tables; historical build PENDING retained; current final **PASS AFTER POST-PUBLICATION VERIFICATION — 210/210**.
- STD-4: 17 unique `CAP-STD-052..068`, all `draft / defined / planned`, 459 sections, 102 mandatory tables; build-time PENDING is historical; final companion records **PASS AFTER POST-PUBLICATION VERIFICATION — 220/220**.
- Permission namespace ambiguity `perm.studio.*` vs `perm.cmdr-studio.*` remains documented and unresolved; no atomic namespace is selected.

## Studio STD-3 final verification addendum — preserved historical evidence
The STD-3 build-time `PENDING` state is retained historically. Final documentary status is **PASS AFTER POST-PUBLICATION VERIFICATION — 210/210 gates PASS**. Canonical companion: `16-quality-and-validation/reports/studio-std3-agents-human-gates-runtime-control-post-publication-verification.md`. The exact STD-3 verification-record SHA is `9babd679f52f3f28458a5f8f4d9c76698ebf875a`. At that closure point STD-4 and Endpoint were NOT STARTED.

## Studio STD-4 closure addendum
The content audit finds all required Studio capability families represented across `CAP-STD-001..068`, with **68 / 1836 / 408** cumulative structure, no owner conflict and no capability-layer placeholder. Studio becomes PASS only after STD-4 reaches **220/220** post-publication. Delivery Roadmap Phase 5 remains PARTIAL because Endpoint is not complete.

No capability is marked implemented, promoted, deployed, active, native or integrated by this registry; documentary PASS never proves implementation.

## Studio STD-4 final verification addendum — 2026-08-11
The build-time `PENDING` line and closure condition above remain historical evidence. Canonical post-publication companion: `16-quality-and-validation/reports/studio-std4-assurance-lifecycle-post-publication-verification.md`.

- STD-4 final documentary verdict: **PASS AFTER POST-PUBLICATION VERIFICATION — 220/220 PASS, 0 PENDING, 0 FAIL**;
- Studio Capability Specification: **PASS** across `CAP-STD-001..068`;
- Studio cumulative: **68 capabilities / 1836 sections / 408 mandatory tables**;
- Delivery Roadmap Phase 5: **PARTIAL**;
- no capability is reclassified as implemented/native/integrated by this documentary PASS.

## Endpoint EPT-1 build addendum — 2026-08-11
The historical pre-EPT-1 state had **0 `CAP-EPT-*` capabilities**. After exact baseline/preflight/namespace revalidation, EPT-1 allocates immutable `CAP-EPT-001..014`, all `draft / defined / planned`.

- structure: **14 capabilities / 378 sections / 84 mandatory tables**;
- scope: enrollment/identity/platform/version/inventory/health/heartbeat/state/capability/Fleet-Policy boundaries/provenance foundations only;
- OPEN-008 remains open; no supported platform is declared;
- Endpoint Screen IDs remain **0**;
- EPT-2..EPT-6 remain **NOT STARTED**;
- build-time quality remains pending remote publication gates; final EPT-1 PASS requires 190/190 post-publication verification;
- Endpoint Capability Specification is **PARTIAL** once EPT-1 content exists;
- Delivery Roadmap Phase 5/global/repository maturity remain **PARTIAL**.

## Endpoint EPT-2 build addendum — 2026-08-11
The preceding EPT-2 NOT STARTED statements are historical pre-EPT-2 snapshots. After namespace/source revalidation, EPT-2 allocates immutable `CAP-EPT-015..030`, all `draft / defined / planned`.

- EPT-2: **16 capabilities / 432 sections / 96 mandatory tables / at least 48 GWT**;
- Endpoint cumulative: **30 / 810 / 180**;
- global current content: **415 capabilities / 413 defined / 2 proposed / 415 planned / 11205 sections / 2490 mandatory tables**;
- EPT-1 remains PASS 190/190 and `CAP-EPT-001..014` are intact;
- Shared owns `telemetry-event` and generic normalization; CAP-EPT-011 remains the foundation availability summary;
- `OPEN-008` remains open; no delivered-platform/source claim is made;
- Endpoint Screen IDs remain 0; EPT-3..EPT-6 remain NOT STARTED;
- build-time EPT-2 quality is **194 PASS / 6 PENDING-BUILD-OR-REMOTE / 0 FAIL** and final PASS requires remote post-publication verification;
- no capability is reclassified as implemented/native/integrated.

## Endpoint EPT-3 build addendum — 2026-08-11
The preceding EPT-3 NOT STARTED statement is historical pre-EPT-3 evidence. After source/namespace/ownership revalidation, EPT-3 allocates immutable `CAP-EPT-031..046`, all `draft / defined / planned`.

- EPT-3: **16 capabilities / 432 sections / 96 mandatory tables / at least 48 GWT**;
- Endpoint cumulative: **46 / 1242 / 276**;
- global content: **431 capabilities / 429 defined / 2 proposed / 431 planned / 11637 sections / 2586 mandatory tables**;
- EPT-1 remains PASS 190/190 and EPT-2 remains PASS 200/200; `CAP-EPT-001..030` are intact;
- Investigate retains Detection Engineering/Case/Evidence/Finding; Command retains canonical Detection/Signal/Alert/Incident; Shared/Settings/Govern/Studio boundaries are preserved;
- `OPEN-008` and `OPEN-017` remain open; no delivered-platform/source claim or final detection runtime/language/model is selected;
- Endpoint Screen IDs remain **0**; EPT-4..EPT-6 remain **NOT STARTED**;
- build-time EPT-3 quality is **204 PASS / 6 PENDING-BUILD-OR-REMOTE / 0 FAIL**, pending 204–209 until exact fifth-commit and remote publication verification;
- no acquisition, Live Response, containment, API/protocol, physical schema, final RBAC or implementation is introduced;
- no capability is reclassified as implemented/native/integrated.

## Endpoint EPT-3 final verification addendum — 2026-08-11
The preceding EPT-3 build-time PENDING material remains historical evidence. Canonical companion: `16-quality-and-validation/reports/endpoint-ept3-local-detection-investigation-post-publication-verification.md`.

- final EPT-3 verdict: **PASS AFTER POST-PUBLICATION VERIFICATION — 210/210 PASS, 0 PENDING, 0 FAIL**;
- exact final EPT-3 HEAD before EPT-4: `67ea28d221ed70baae83ff0689048685e1aacf74`;
- `CAP-EPT-001..046` remain immutable `draft / defined / planned` capability contracts;
- EPT-4 was NOT STARTED at that closure point.

## Endpoint EPT-4 build addendum — 2026-08-11
The EPT-4 pre-state is the final EPT-3 HEAD above. After exact namespace/source/ownership audit, EPT-4 allocates `CAP-EPT-047..064`, all `draft / defined / planned`.

- EPT-4 structure: **18 capabilities / 486 sections / 108 mandatory tables / at least 54 GWT**;
- Endpoint cumulative: **64 / 1728 / 384**;
- global content: **449 capabilities / 447 defined / 2 proposed / 449 planned / 12123 sections / 2694 mandatory tables**;
- Collection Request remains Investigate-owned; Endpoint owns local technical eligibility/operation/output only;
- OPEN-014 remains open: Collection Item/Collected Technical Output/Package do not become Artifact/Attachment/Evidence automatically;
- OPEN-015 remains open: Endpoint Technical Execution != Studio Tool Call/Automation Run != Govern Response Run;
- OPEN-008/017 remain open; no supported platform or universal detection/script runtime is selected;
- effectful/destructive containment primitives remain EPT-5 boundary;
- Endpoint Screen IDs remain 0; EPT-5/EPT-6 remain NOT STARTED;
- build-time quality: **214 PASS / 6 PENDING-BUILD-OR-REMOTE / 0 FAIL**, pending 214–219 until exact fifth commit and remote verification;
- no capability is reclassified as implemented/native/integrated and no API/protocol/transport/physical schema/storage/final RBAC/product code is introduced.

## Endpoint EPT-5 build addendum — 2026-08-11
The prior EPT-5 NOT STARTED lines are historical snapshots. Exact baseline: `5d576295fa12693ef375a35cfe515d7bdf577f68`, final EPT-4 verification record.

- complete source audit: all 7 Containment docs, effectful Live Response process/network/file/service sources, Govern authority/Run/Verification/Rollback/Result, Studio/Settings/Investigate/Shared/Security/EPT-6 boundaries;
- final source-driven set: **`CAP-EPT-065..081` — 17 capabilities / 459 sections / 102 mandatory tables / at least 51 GWT**;
- `CAP-EPT-081` is independently justified local-session lock/termination; directory identity actions remain external;
- Endpoint cumulative: **81 / 2187 / 486**;
- global content: **466 capabilities / 464 defined / 2 proposed / 466 planned / 12582 sections / 2796 mandatory tables**;
- Requirements remain **122 = 99/20/3/0**; OPEN remains **18**;
- Govern owns authority/Approval/Decision/Response Run/response verification/rollback/Result; Endpoint owns technical primitive/readiness/effect/target observations/reversal only;
- OPEN-007/008/013/014/015/017 remain open;
- Endpoint Screen IDs remain 0; EPT-6 remains NOT STARTED;
- build-time quality: **224 PASS / 6 PENDING-BUILD-OR-REMOTE / 0 FAIL**, pending gates 224–229;
- no capability is reclassified as implemented/native/integrated and no API/protocol/native command/shell/physical schema/final RBAC/final response engine/product code is introduced.

## Endpoint EPT-6 reconciliation addendum — 2026-08-12
The preceding EPT-5 totals and EPT-6 NOT STARTED statements remain preserved as historical pre-EPT-6 evidence. The published EPT-6 shard and closure records supersede them for current-state indexing only.

- EPT-6 immutable set: **`CAP-EPT-082..099` — 18 capabilities / 486 sections / 108 mandatory tables / at least 54 GWT**;
- all 18 are `draft / defined / planned`; no capability contract is modified by this reconciliation;
- Endpoint cumulative current state: **99 capabilities / 2673 sections / 594 mandatory tables**;
- global current state at EPT-6 closure: **484 capabilities / 482 defined / 2 proposed / 484 planned / 13068 sections / 2904 mandatory tables**;
- Requirements remain **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**;
- OPEN decisions remain **18**;
- Endpoint Screen IDs remain **0**;
- Command **27**, Investigate **243**, Govern **47**, Studio **68**, Endpoint **99** remain the current closed-domain counts;
- Endpoint Capability Specification remains **PASS** and Delivery Roadmap Phase 5 remains **PASS — capability specification complete**;
- Delivery Roadmap Phase 6 — Platform Scale was **NOT STARTED** at this historical closure point;
- Phase 6 capability created/reserved, new capability namespace, new Screen and implementation were **0** at that historical closure point.

This reconciliation corrected only the stale canonical global register summary after EPT-6 publication. It allocated or reserved no ID and began no Phase 6 functional work.

## Settings Tenant, Environment and Administrative Foundations build addendum — 2026-08-12
The preceding Phase-6 NOT STARTED statements remain historical pre-build evidence. After the 120/120 Platform Scale foundations preflight and execution-time namespace revalidation, this first functional Phase-6 lot allocates exactly `CAP-SET-001..004`.

- all four are `draft / defined / planned` and owned uniquely by Platform Settings Product Lead;
- structure: **4 capabilities / 108 numbered sections / 24 mandatory capability tables / at least 12 GWT**;
- current global build state: **488 capabilities / 486 defined / 2 proposed / 488 planned / 13176 sections / 2928 mandatory tables**;
- Settings capability specification: **PARTIAL**;
- Delivery Roadmap Phase 6 Capability Specification: **PARTIAL / PENDING POST-PUBLICATION VERIFICATION**;
- Requirements remain **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**; no Requirement ID is added or removed;
- OPEN decisions remain **18**; this lot creates/closes 0;
- new Permission IDs: **0**; new Screen IDs: **0**; new canonical objects: **0**;
- `CAP-SET-004` capability owner is Platform Settings Product Lead only; Experience Architecture remains dependency/mechanism owner;
- `CAP-SET-005+` remains unallocated/unreserved;
- build-time quality target is **154 PASS / 6 PENDING-REMOTE / 0 FAIL** and final PASS requires 160/160 after remote verification;
- Command 27, Investigate 243, Govern 47, Studio 68 and Endpoint 99 capability contracts remain unchanged;
- no implementation/API/protocol/physical schema/final RBAC/platform-support claim is introduced.

No capability is reclassified as implemented, native, integrated, active or deployed by this documentary build state.

## Settings Tenant, Environment and Administrative Foundations post-publication addendum — 2026-08-12
The preceding build-time PENDING material remains historical evidence. Canonical companion: `16-quality-and-validation/reports/platform-scale-tenant-environment-administrative-foundations-post-publication-verification.md`.

- final first-lot verdict: **PASS AFTER POST-PUBLICATION VERIFICATION — 160/160 PASS, 0 PENDING, 0 FAIL**;
- exact baseline for Identity Administration is `d605265f5b8a2e4350388b4ec9cfe51920a4aa50`;
- `CAP-SET-001..004` remain immutable `draft / defined / planned`;
- no `CAP-SET-005+` was allocated or reserved at that closure point;
- Settings and Delivery Roadmap Phase 6 Capability Specification remained **PARTIAL**.

## Settings Identity Administration build addendum — 2026-08-12
After exact execution-time preflight and source/namespace/ownership audit, Identity Administration — Principals, Roles and Access Reviews allocates exactly `CAP-SET-005..007`.

- all three are `draft / defined / planned`, owned uniquely by Platform Settings Product Lead;
- structure: **3 capabilities / 81 numbered sections / 18 mandatory tables / at least 9 meaningful GWT**;
- Settings cumulative build state: **7 capabilities / 189 sections / 42 mandatory tables**;
- global build state: **491 capabilities / 489 defined / 2 proposed / 491 planned / 13257 sections / 2946 mandatory tables**;
- Requirements remain **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**; no Requirement ID or state changes are introduced by this documentary evidence;
- OPEN remains **18**; no decision is created or closed; OPEN-013 remains unresolved;
- Principal states remain `pending/active/suspended/revoked`; Role states remain `draft/active/deprecated`; Role expiry is a condition/constraint only;
- Groups, Group Membership, generic Access Assignment and Effective Access remain outside scope;
- CAP-SET-007 uses a **revocation disposition/handoff**, not direct generic assignment removal, because concrete assignment-removal mechanics are not canonically sourced;
- new Permission IDs: **0**; new Screen IDs: **0**; new canonical objects: **0**;
- existing `SET-IAM-001` and `SET-AUD-001` are reused; Security retains authorization semantics and Govern retains authority/Decision/Response objects;
- build-time conformance: **168 PASS / 6 PENDING-REMOTE / 0 FAIL**; final PASS requires remote publication verification;
- Command 27, Investigate 243, Govern 47, Studio 68 and Endpoint 99 closed-domain capability contracts remain unchanged;
- no implementation/API/protocol/physical schema/final RBAC/ABAC/authentication protocol/platform-support claim is introduced;
- `CAP-SET-008+` is neither allocated nor reserved at this historical build snapshot.

No capability is reclassified as implemented, native, integrated, active or deployed by this documentary build state.

## Settings Identity Administration post-publication addendum — 2026-08-13
The preceding Identity build-time PENDING material remains historical evidence. Canonical post-publication verification closed all six remote-dependent gates without changing `CAP-SET-005..007`.

- final Identity verdict: **PASS AFTER POST-PUBLICATION VERIFICATION — 174/174 PASS, 0 PENDING, 0 FAIL**;
- exact Identity final HEAD / Secrets & Connections baseline: `a5c450528ad25f090832adcda4ae37c79bea072b`;
- Identity structure remains **3 capabilities / 81 sections / 18 mandatory tables / 12 meaningful GWT**;
- global state at Identity closure remains **491 / 489 defined / 2 proposed / 491 planned / 13257 sections / 2946 tables**;
- Settings remains **PARTIAL** and Phase 6 remains **PARTIAL**.

## Settings Secrets & Connections functional build addendum — 2026-08-13
After exact execution-time branch/PR/main/README/source/namespace/security/ownership revalidation, **Secrets & Connections — Integration and Secret Reference Administration** allocates exactly `CAP-SET-008..009`.

- `CAP-SET-008` — Integration Administrative Lifecycle, Validation and Connection State;
- `CAP-SET-009` — Secret Reference Administrative Lifecycle, Rotation and Revocation;
- both are `draft / defined / planned`, owned uniquely by Platform Settings Product Lead;
- structure: **2 capabilities / 54 numbered sections / 12 mandatory tables / 12 meaningful GWT**;
- Settings cumulative build state: **9 capabilities / 243 sections / 54 mandatory tables**;
- global build state: **493 capabilities / 491 defined / 2 proposed / 493 planned / 13311 sections / 2958 mandatory tables**;
- `CAP-SET-010+` is neither allocated nor reserved;
- new Permission IDs: **0**; new Screen IDs: **0**; new canonical objects: **0**;
- canonical Integration states remain `draft/validating/active/degraded/disabled/error`;
- canonical Secret Reference states remain `pending/active/rotating/expired/revoked`;
- Integration `capabilities` metadata is explicitly not CMDR Capability/CAP-* allocation;
- no Connection/Connector/Credential/raw Secret object is created;
- current sources do not assign Settings the technical external connection-probe executor; CAP-SET-008 therefore owns administrative test request/preconditions/status/result projection/provenance and handoff only for external testing;
- current sources do not assign Settings underlying-secret generation/write/rotation or external credential revocation; CAP-SET-009 remains reference-only and owns administrative lifecycle/reference mutation/handoff/result projection only;
- Tenant is mandatory and Environment source-dependent only;
- Requirements remain **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**;
- OPEN remains **18**; `OPEN-008`, `OPEN-012`, `OPEN-013` remain open as applicable;
- Command 27, Investigate 243, Govern 47, Studio 68, Endpoint 99 and `CAP-SET-001..007` remain unchanged;
- Models & Providers, Sources & Parsers, Customer/MSSP/Delivery and implementation remain outside this lot;
- build-time quality: **166 PASS / 6 PENDING-REMOTE / 0 FAIL**; final 172/172 is forbidden before actual remote verification;
- no implementation/API/protocol/physical schema/provider support/Vault/KMS/HSM/final RBAC claim is introduced.

No capability is reclassified as implemented, native, integrated, active or deployed by this documentary build state.

## Settings Models & Providers functional build addendum — 2026-08-13
After execution-time revalidation, **Models & Providers — Model Provider Administration and Model Routing** allocates exactly `CAP-SET-010..011`.

- `CAP-SET-010` — Model Provider Administrative Lifecycle, Validation, Model Availability and Health Projection;
- `CAP-SET-011` — Model Routing Configuration, Eligibility, Fallback Constraints and Provider Switch Provenance;
- both are `draft / defined / planned`, owned uniquely by Platform Settings Product Lead;
- structure: **2 capabilities / 54 numbered sections / 12 mandatory tables / 13 meaningful GWT**;
- Settings cumulative build state: **11 capabilities / 297 sections / 66 mandatory tables**;
- global build state: **495 capabilities / 493 defined / 2 proposed / 495 planned / 13365 sections / 2970 mandatory tables**;
- `CAP-SET-012+` remains unallocated and unreserved;
- new Permission IDs: **0**; new Screen IDs: **0**; new canonical objects: **0**;
- Model Provider remains distinct from Integration and no canonical Model or routing-policy object is created;
- provider administration/routing configuration remains distinct from provider execution;
- availability, health and effective selection remain source-attributed projections;
- fallback configuration does not assert automatic failover;
- Requirements remain **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**;
- OPEN remains **18**; `OPEN-008`, `OPEN-012`, `OPEN-013` remain open;
- Command 27, Investigate 243, Govern 47, Studio 68, Endpoint 99 and `CAP-SET-001..009` remain unchanged;
- Sources & Parsers, Customer/MSSP/Delivery and implementation remain outside this lot;
- build-time quality: **198 PASS / 6 PENDING-REMOTE / 0 FAIL**; final 204/204 requires actual remote verification.

No capability is reclassified as implemented, native, integrated, active or deployed by this documentary build state.