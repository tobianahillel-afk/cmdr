---
id: roadmap-phase-5-studio-and-endpoint
domain: 18-roadmap-and-releases
status: draft
owner: Product Operations Lead
updated: 2026-08-11
source-of-truth: canonical
---
# Phase 5 Studio And Endpoint

## Objectif
Définir Phase 5 Studio and Endpoint pour CMDR tout en conservant Studio et Endpoint comme deux capability domains séparés.

## Périmètre
Document canonique Roadmap. Il ne remplace pas les sources propriétaires et ne transforme aucun execution lot en sous-phase numérotée.

## Propriétaire fonctionnel
Product Operations Lead.

## Objets concernés
- Concepts du document
- Références canoniques liées

## Fonctionnalités historiques prévues
- Skills/Agents/workflows.
- Assurance/control room.
- EDR telemetry/detection/response.
- Fleet administration comme dépendance Settings, sans transfert d'ownership vers Endpoint.

## UX et interactions
- Navigation par liens stables.
- Contenu lisible en thème clair et sombre.
- Aucune duplication des définitions externes.

## Permissions
Les modifications suivent `../14-security-permissions-and-trust/permission-model.md`; aucun execution lot ne finalise implicitement RBAC/ABAC.

## Dépendances
- `../00-governance/source-of-truth-policy.md`
- `dependency-roadmap.md`
- Studio/Settings/Govern/Shared/Endpoint ownership sources.

## STD-1 — verified
- STD-1 — Studio Foundations — Tools, Skills, Library and Ownership Contracts: **PASS AFTER POST-PUBLICATION VERIFICATION — 190/190**.
- `CAP-STD-001..016`: 16 capabilities / 432 sections / 96 mandatory tables.

## STD-2 — current execution
- STD-2 — Workflow Builder & Orchestration: 17 capabilities `CAP-STD-017..033` / 459 sections / 102 mandatory tables.
- scope: Workflow definition/version, Builder, data/graph/composition, deterministic control, errors/retry/partial/compensation, Human Gate boundary, readiness/pre-publish/provenance.
- status: **PASS AFTER POST-PUBLICATION VERIFICATION — 200/200**; fifth functional SHA `655e9ce0ade2d64a7738a6a479572fef9b6f0e2f` verified at 5 ahead / 0 behind from baseline.
- new Screen IDs / detailed rewrites: 0 / 0.
- runtime scheduler / Automation Run lifecycle / publishing-deployment / Endpoint implementation: not included.

## Future lots
- STD-3 — Agents, Human Gates & Runtime Control: NOT STARTED.
- STD-4 — Assurance & Lifecycle: NOT STARTED.
- Endpoint capability specification: NOT STARTED.
- Studio capability specification: PARTIAL.
- Delivery Roadmap Phase 5 overall: PARTIAL.

STD-1/2/3/4 are execution-lot labels only. No Phase 5A/5B/5C/5D exists.

## Implementation boundary
No API, protocol, product code, orchestration language/runtime, final graph/Tool/Tool Call physical schema, final RBAC/ABAC, provider/runtime selection, detailed screen rewrite or Endpoint implementation is introduced by STD-2.

## Questions ouvertes
OPEN-007/013/015 remain open where consumed. OPEN-008 remains Endpoint/platform support. The permission namespace anomaly remains unresolved.

## Next candidate
STD-3 — Agents, Human Gates & Runtime Control. **Do not begin STD-3 implicitly.**

## STD-2 post-publication verification evidence
Baseline `04dcdb43fd7f944a700bf936eebef003546095eb`; fifth functional/build SHA `655e9ce0ade2d64a7738a6a479572fef9b6f0e2f`; five functional commits reachable in order; PR #2 remains Draft/open/unmerged; branch/main README unchanged; CI N/A; `CAP-STD-001..016` intact; Endpoint 0; STD-3/4 not started.

Final STD-2 documentary verdict: **PASS AFTER POST-PUBLICATION VERIFICATION — 200/200**.

## STD-3 current execution addendum
The full STD-1/STD-2 roadmap text above is preserved as the exact pre-STD-3 snapshot; its STD-3 `NOT STARTED` / `Next candidate` text is historical evidence only.

- execution lot: **STD-3 — Agents, Human Gates & Runtime Control**;
- `CAP-STD-034..051`: **18 capabilities / 486 sections / 108 mandatory tables**;
- fifth functional/build SHA: `c658168c9de6bd803941116989bc3aaedf154260`;
- baseline `c472b055ce00fd33efd96ac920b0add5f65ab8f7` → build: **5 ahead / 0 behind**, same merge base;
- scope: Agents/Teams, bounded planning/oversight, Human Gate runtime, Automation Runs/control, queue/scheduling semantics, failures/retries/partial, transient context, Control Room, outcome/handoff and provenance;
- ownership boundaries: Govern retains Approval/Decision/Response Run/Result; Settings retains provider/secret/runtime administration; Shared retains generic Jobs/queue/scheduling/Trace; Endpoint retains technical primitives;
- no new Screen ID/detailed rewrite, agent framework/model/provider/runtime/scheduler/API/protocol/code/final JSON Schema/RBAC;
- post-publication history-preservation correction is being recorded before the final 210-gate verdict.

STD-3 remains **PENDING POST-PUBLICATION VERIFICATION**. STD-4 and Endpoint remain **NOT STARTED**; Delivery Roadmap Phase 5 and global/repository maturity remain PARTIAL.

## STD-3 final verification addendum
The preceding `PENDING` state is preserved as historical evidence. STD-3 final documentary verdict is **PASS AFTER POST-PUBLICATION VERIFICATION — 210/210 gates PASS, 0 PENDING, 0 FAIL**.

- canonical companion: `../16-quality-and-validation/reports/studio-std3-agents-human-gates-runtime-control-post-publication-verification.md`;
- history restoration: `bff197f7cc33296220a211425f62ac6b806a5f7a` changes no capability contract;
- Studio now has **51 capabilities / 1377 sections / 306 mandatory tables** across STD-1/2/3;
- global total: **368 capabilities / 9936 sections / 2208 mandatory tables**;
- STD-4 — Assurance & Lifecycle: **NOT STARTED**;
- Endpoint capability specification: **NOT STARTED**;
- Delivery Roadmap Phase 5 overall remains **PARTIAL**;
- next candidate is STD-4, but **do not begin it implicitly**.

## STD-4 current execution addendum
The complete STD-1/2/3 roadmap text above remains preserved. Earlier `STD-4 NOT STARTED` / `next candidate` text is historical evidence only.

- execution lot: **STD-4 — Assurance & Lifecycle**;
- `CAP-STD-052..068`: **17 capabilities / 459 sections / 102 mandatory tables**;
- content: Evaluation/Suites/Cases, deterministic/assisted evaluation, Simulation, Regression, Reliability/Safety, Boundary Assurance, Evaluation Result, Quality Gates/Readiness, Publishing Candidate/Review, release/promotion/channel, deployment compatibility/lifecycle/health/reversion, deprecation/retirement/migration and lifecycle provenance;
- Studio cumulative content: **68 capabilities / 1836 sections / 408 mandatory tables**;
- global content total: **385 capabilities / 10395 sections / 2310 mandatory tables**;
- content closure audit across Studio is positive;
- build-time gates: **212 PASS / 8 PENDING-REMOTE / 0 FAIL**;
- Studio is **PARTIAL / PENDING POST-PUBLICATION VERIFICATION** until 220/220;
- Endpoint capability specification remains **NOT STARTED / 0**;
- Delivery Roadmap Phase 5 remains **PARTIAL** and must remain PARTIAL even after Studio PASS until Endpoint closes;
- no Phase 5A/B/C/D/E, no implementation, no Endpoint capability.

Next after successful Studio closure: **Endpoint lot 1 — Enrollment, Inventory, Health and Platform Foundations**, but it must not start in STD-4.

---

## STD-4 final verification addendum — 2026-08-11
The preceding build-time/PENDING material remains historical evidence. Remote recovery verifies the exact five functional commits from baseline `9babd679f52f3f28458a5f8f4d9c76698ebf875a` to build `216bff304fa389e4814cb610097571a4a83c1c54` at **5 ahead / 0 behind**, then documentary correction `c21ea86cde1bea425d7d9233d9973b867f5ef9e8`.

- canonical companion: `../16-quality-and-validation/reports/studio-std4-assurance-lifecycle-post-publication-verification.md`;
- STD-4 final documentary verdict: **PASS AFTER POST-PUBLICATION VERIFICATION — 220/220 PASS, 0 PENDING, 0 FAIL**;
- Studio Capability Specification: **PASS** across **68 capabilities / 1836 sections / 408 mandatory tables**;
- global capabilities remain **385**, Requirements **122 = 99/20/3/0**, OPEN **18**;
- Endpoint Capability Specification remains **NOT STARTED / 0**;
- Delivery Roadmap Phase 5 remains **PARTIAL**;
- Global Capability Specification / repository maturity remain **PARTIAL**;
- no Phase 5A/B/C/D/E, implementation, Endpoint capability or Endpoint Screen ID is introduced.

The next run after successful final remote publication is **ENDPOINT PREFLIGHT RERUN / CLOSURE**. The previous Endpoint preflight remains historically **BLOCKED 96/100** until that rerun; EPT-1 must not begin before a 100/100 preflight.

---

## Endpoint foundations preflight rerun — closure addendum — 2026-08-11
The previous Endpoint preflight remains preserved historically as **BLOCKED 96/100** at `c21ea86cde1bea425d7d9233d9973b867f5ef9e8`. The current rerun starts from exact Studio-closure HEAD `9030186e7aa12990a3d8fb6f30aa107539e2a117` and is recorded in `../16-quality-and-validation/reports/endpoint-capability-specification-foundations-preflight.md`.

After successful publication and remote verification of the single preflight-closure commit:
- Endpoint preflight: **PASS — 100/100**;
- Endpoint Capability Specification: **NOT STARTED**;
- Endpoint capabilities: **0**;
- concrete/reserved `CAP-EPT-*` IDs: **0 / 0**;
- Endpoint Screen IDs: **0**;
- Studio remains **PASS**;
- Delivery Roadmap Phase 5 remains **PARTIAL**;
- Global Capability Specification / repository maturity remain **PARTIAL**.

Endpoint execution-lot decomposition is retained as EPT-1 through EPT-6. The next functional execution becomes **EPT-1 — Enrollment, Inventory, Health and Platform Foundations** only after this preflight is remotely confirmed 100/100. This addendum does not start EPT-1, allocate `CAP-EPT-001..014`, create a platform-support promise or transfer Fleet/Policy ownership from Settings.

---

## EPT-1 current execution — build-time addendum — 2026-08-11
The preflight section above is preserved as the exact pre-EPT-1 snapshot. EPT-1 is an execution lot only and creates no Phase 5E1/5E2/etc.

- lot: **EPT-1 — Enrollment, Inventory, Health and Platform Foundations**;
- exact baseline: `8326a8cf9e9ca3b645395d192c24856058e67034`;
- `CAP-EPT-001..014`: **14 capabilities / 378 sections / 84 mandatory tables**, all `draft / defined / planned`;
- scope: Agent identity/registration, local enrollment, tenant/environment binding, platform/OS/architecture, version/build/compatibility, inventory/freshness, health/self-check, heartbeat/connectivity, operational states, capability availability, Fleet/Policy boundaries and provenance/handoff;
- OPEN-008 remains open; no Windows/Linux/macOS or other platform is declared officially supported/delivered;
- Platform Settings retains Fleet, enrollment administration, Endpoint Policy/assignment, upgrade waves and administrative configuration;
- Endpoint Screen IDs remain **0**;
- EPT-2 through EPT-6 remain **NOT STARTED**;
- content totals become **399 global capabilities / 397 defined / 2 proposed / 399 planned / 10773 sections / 2394 mandatory tables**;
- Requirements remain **122 = 99/20/3/0**, OPEN remains **18**;
- build-time quality: **184 PASS / 6 PENDING-BUILD-OR-REMOTE / 0 FAIL** before exact fifth-commit and remote checks;
- Endpoint Capability Specification is **PARTIAL** with EPT-1 content; Delivery Roadmap Phase 5/global/repository maturity remain **PARTIAL**.

After successful remote 190/190 verification, EPT-1 becomes **PASS AFTER POST-PUBLICATION VERIFICATION**. The next candidate is **EPT-2 — Telemetry, Observation and Technical Capability Declaration**, but EPT-2 must not start in this run.

---

## EPT-1 final post-publication verification — 2026-08-11
The preceding build-time state remains historical evidence. Canonical companion: `../16-quality-and-validation/reports/endpoint-ept1-enrollment-inventory-health-platform-foundations-post-publication-verification.md`.

- five functional commits from `8326a8cf9e9ca3b645395d192c24856058e67034` to build `828b231ec2de4d3b891410a643898577f14cbcc4` are verified at **5 ahead / 0 behind**, same merge base;
- build remote PR/main/README checks passed; CI/status = N/A;
- a real post-publication documentary divergence is corrected without changing any capability contract;
- after remote verification of the documentary correction, EPT-1 final verdict is **PASS AFTER POST-PUBLICATION VERIFICATION — 190/190 PASS, 0 PENDING, 0 FAIL**;
- Endpoint Capability Specification remains **PARTIAL**: EPT-1 PASS, EPT-2..EPT-6 NOT STARTED;
- Endpoint remains **14 / 378 / 84** and global remains **399 / 397 defined / 2 proposed / 399 planned / 10773 / 2394**;
- OPEN-008 remains open; no supported-platform delivery claim is created;
- Command/Investigate/Govern/Studio remain PASS;
- Delivery Roadmap Phase 5/global/repository maturity remain PARTIAL;
- exact documentary correction SHA is recorded in PR #2 after publication;
- next candidate is **EPT-2 — Telemetry, Observation and Technical Capability Declaration**, but it remains **NOT STARTED**.

---

## EPT-2 current execution — build-time addendum — 2026-08-11
The preceding EPT-2 NOT STARTED statement remains historical pre-EPT-2 evidence. EPT-2 is an execution lot only; no Phase 5E2 is created.

- exact baseline: `1f8e482f6b7949885bd1bd7ae691215bde187b28`;
- EPT-1 remains **PASS 190/190** with `CAP-EPT-001..014` intact;
- EPT-2 allocates exactly `CAP-EPT-015..030`: **16 capabilities / 432 sections / 96 mandatory tables / at least 48 GWT**, all `draft / defined / planned`;
- Endpoint cumulative: **30 capabilities / 810 sections / 180 mandatory tables**;
- global content: **415 capabilities / 413 defined / 2 proposed / 415 planned / 11205 sections / 2490 tables**;
- Shared retains `telemetry-event` and generic normalization; Settings administration, Investigate Evidence/Finding/Case, Studio Tool/Run and Govern Decision/Response Run/Result remain external owners;
- `OPEN-008` remains open; no platform/source support is declared delivered;
- Endpoint Screen IDs remain 0; EPT-3..EPT-6 remain NOT STARTED;
- build-time quality: **194 PASS / 6 PENDING-BUILD-OR-REMOTE / 0 FAIL**, gates 194–199 pending until fifth-commit/remote verification;
- Delivery Roadmap Phase 5 remains **PARTIAL**; Global Capability Specification and repository maturity remain **PARTIAL**;
- no implementation/API/protocol/port/physical schema/storage/event bus/SIEM/final RBAC or EPT-3+ work is introduced.

EPT-2 may become **PASS AFTER POST-PUBLICATION VERIFICATION — 200/200** only after remote verification. Do not start EPT-3.

---

## EPT-3 current execution — build-time addendum — 2026-08-11
The preceding EPT-3 NOT STARTED statement is historical pre-EPT-3 evidence. EPT-3 is an execution lot only; no Phase 5E3 or new Capability Specification Phase is created.

- exact baseline: `5d7c037aff6004984416665e7e188a8700e62b2f`;
- EPT-1 remains **PASS 190/190** and EPT-2 remains **PASS 200/200** with `CAP-EPT-001..030` intact;
- EPT-3 allocates exactly `CAP-EPT-031..046`: **16 capabilities / 432 sections / 96 mandatory tables / at least 48 GWT**, all `draft / defined / planned`;
- Endpoint cumulative: **46 capabilities / 1242 sections / 276 mandatory tables**;
- global content: **431 capabilities / 429 defined / 2 proposed / 431 planned / 11637 sections / 2586 tables**;
- scope: local Detection Content consumption/eligibility/evaluation/match/candidate/context/grouping/coverage plus local read-only process/file/network/user-session/system investigation, timeline/correlation, pivots, summary/handoff and provenance;
- Investigate retains Detection Engineering/Case/Evidence/Finding; Command retains canonical Detection/Signal/Alert/Incident; Shared, Settings, Govern and Studio boundaries remain preserved;
- `OPEN-008` and `OPEN-017` remain open; no supported platform/source or final detection runtime/language/model/engine/portability decision;
- Endpoint Screen IDs remain 0; EPT-4..EPT-6 remain **NOT STARTED**;
- build-time quality: **204 PASS / 6 PENDING-BUILD-OR-REMOTE / 0 FAIL**, gates 204–209 pending until fifth-commit/remote verification;
- Delivery Roadmap Phase 5, Global Capability Specification and repository maturity remain **PARTIAL**;
- no Collection, Live Response, containment, response execution, API/protocol, physical schema, final RBAC or implementation is introduced.

EPT-3 may become **PASS AFTER POST-PUBLICATION VERIFICATION — 210/210** only after remote verification. Next candidate after PASS is **EPT-4 — Collection and Live Response Technical Execution**, but **do not begin EPT-4 in this run**.

---

## EPT-4 current execution — build-time addendum — 2026-08-11
The preceding EPT-4 NOT STARTED statements are historical pre-EPT-4 evidence. EPT-4 is an execution lot only and creates no Phase 5E4.

- exact baseline: `67ea28d221ed70baae83ff0689048685e1aacf74` — final EPT-3 verification record;
- EPT-1 **PASS 190/190**, EPT-2 **PASS 200/200**, EPT-3 **PASS 210/210**; `CAP-EPT-001..046` intact;
- EPT-4 allocates exactly `CAP-EPT-047..064`: **18 capabilities / 486 sections / 108 mandatory tables / at least 54 GWT**, all `draft / defined / planned`;
- Endpoint cumulative: **64 capabilities / 1728 sections / 384 mandatory tables**;
- global content: **449 capabilities / 447 defined / 2 proposed / 449 planned / 12123 sections / 2694 tables**;
- Collection Request/Case/Evidence/Finding/Artifact qualification remain Investigate-owned; Endpoint owns technical eligibility/operations/output only;
- Govern retains Decision/Response Run/Result and response authority; Studio retains Tool Call/Automation Run; Settings Fleet/Policy/Secret References; Shared generic Jobs/Trace/Activity/Export;
- OPEN-008/014/015/017 remain open; no platform support, Artifact identity, execution bridge or universal runtime is finalized;
- Endpoint Screen IDs remain 0; EPT-5/EPT-6 remain **NOT STARTED**;
- build-time quality: **214 PASS / 6 PENDING-BUILD-OR-REMOTE / 0 FAIL**, pending gates 214–219 until exact fifth-commit and remote verification;
- effectful containment/remediation actions remain future EPT-5; update/resilience/security depth remains future EPT-6;
- no API/protocol/remote-shell protocol/transport/command catalog/final runtime/physical schema/storage/final RBAC/product implementation.

EPT-4 may become **PASS AFTER POST-PUBLICATION VERIFICATION — 220/220** only after exact five-commit ancestry and remote verification. Endpoint Capability Specification and Delivery Roadmap Phase 5 remain PARTIAL. **Do not begin EPT-5.**