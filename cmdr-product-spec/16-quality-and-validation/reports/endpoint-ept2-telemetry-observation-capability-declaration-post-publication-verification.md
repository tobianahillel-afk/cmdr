---
id: endpoint-ept2-telemetry-observation-capability-declaration-post-publication-verification
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-11
source-of-truth: quality-report
requirements: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-018, REQ-PROD-019, REQ-OBJ-008, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008, OPEN-015]
---
# Endpoint EPT-2 — Post-Publication Verification

## Scope
This companion verifies only **EPT-2 — Telemetry, Observation and Technical Capability Declaration** under Delivery Roadmap Phase 5 — Studio and Endpoint. It modifies no `CAP-EPT-*` contract and does not start EPT-3.

## Exact baseline and functional chain
Starting baseline: `1f8e482f6b7949885bd1bd7ae691215bde187b28` — `docs: record Endpoint EPT-1 post-publication verification`.

1. `c74e9de159671e675450a85e1db2403eb7684dbc` — `docs: establish Endpoint telemetry source and observation boundaries`.
2. `770c71559dc2427bc8bb033e955eba13ce98bce6` — `docs: define Endpoint activity telemetry and normalization semantics`.
3. `5936d2ed123bfc847be3192c7ab51f48f569af1d` — `docs: specify Endpoint telemetry quality privacy and loss handling`.
4. `bc33bbad828ab324cfc02d22b60edf457407f94c` — `docs: document Endpoint technical capability declaration and availability`.
5. `d7698749fb35030295e2c061c2b98596ff87588c` — `docs: update Endpoint telemetry traceability and quality gates`.

Baseline → fifth functional/build SHA is verified at **5 ahead / 0 behind**, same merge base. Publication was a fast-forward with no force, rebase, reset or history rewrite.

## Remote build verification
- remote HEAD reached `d7698749fb35030295e2c061c2b98596ff87588c`;
- PR #2 remained **open / Draft / unmerged**, base `main`;
- `main` remained `bc1ec59e5e79ccbaba291e2984b0eb7d5e54129c`;
- branch/main root README remained exact `# cmdr`, blob `901c74cda52e28b5ff7fc425ddf28ef89f3ad875`;
- build SHA had no commit statuses and no workflow runs, therefore **CI/status = N/A**;
- CAP-EPT-001 remained present/intact and `CAP-EPT-031` remained absent.

## Structural evidence
`CAP-EPT-015..030` are the complete EPT-2 set: **16 capability files / 432 numbered sections / 96 mandatory tables / at least 48 Given/When/Then scenarios**, all `draft / defined / planned`, with **0 duplicate/recycled ID / 0 owner conflict / 0 empty mandatory table**.

Endpoint cumulative: **30 capabilities / 810 sections / 180 mandatory tables**. Global: **415 capabilities / 413 defined / 2 proposed / 415 planned / 11205 sections / 2490 mandatory tables**.

## Source and ownership verification
All 11 Endpoint Telemetry source documents were audited before allocation. Shared remains owner of `telemetry-event` and generic telemetry normalization. Endpoint owns endpoint-local source/observation/quality/declaration/provenance semantics. Platform Settings retains source/Fleet/Policy/provider/secret/tenant administration. Investigate retains Evidence/Finding/Case; Studio retains Tool/Tool Call/Automation Run; Govern retains Decision/Response Run/Result; Shared retains generic Trace/Activity/Jobs/Search/Reporting/Export.

CAP-EPT-011 remains the EPT-1 foundation availability summary. CAP-EPT-027/028 deepen declaration and dynamic availability without duplicating CAP-EPT-011.

## Telemetry semantics verified
EPT-2 covers source/sensor identity and local observation; Shared event projection; source/observation/ingestion time distinctions; process/file/network/auth/system/sensor observations; normalization/classification; ordering/late/duplicate/freshness/gap/loss; volume/sampling/backpressure; privacy/masking; detailed capability declaration/availability/degradation/dependencies; consumer handoff and provenance.

Mandatory distinctions remain explicit: telemetry/observation != Detection/Finding/Evidence/Alert/Result; process/file/network/auth/security-state facts are not malicious/compromise/incident conclusions; emitted != delivered != consumed; ordering != causality; sampling != loss; gap != tampering/sensor failure; backpressure != endpoint failure; declared != available != globally supported; Endpoint capability != Studio Tool; continuous telemetry != forensic Collection; local audit != Shared Trace.

## OPEN, screens and implementation boundary
The programme retains **18 OPEN decisions**. `OPEN-008 — platform/source availability and support` remains **OPEN**. No Windows/Linux/macOS/cloud/container/mobile or universal sensor/source capability is declared delivered.

The Screen Register remains **56 active screens**, with **0 Endpoint product Screen IDs**. No Screen ID, detailed screen design, API, protocol, port, physical event schema, storage engine, event bus, SIEM engine, final query language, final normalization standard, final RBAC/ABAC or product code is introduced.

EPT-3, EPT-4, EPT-5 and EPT-6 remain **NOT STARTED**.

## Non-regression
EPT-1 remains **PASS AFTER POST-PUBLICATION VERIFICATION — 190/190**, with `CAP-EPT-001..014` intact. Command remains PASS 27; Investigate PASS 243; Govern PASS 47; Studio PASS 68. Delivery Roadmap Phase 5, Global Capability Specification and repository maturity remain PARTIAL.

Requirements remain **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**.

## Gates 191–200 final closure
- 191 PASS — conformance report exists.
- 192 PASS — build-time 194/6 state is explicit and preserved historically.
- 193 PASS — remote-dependent gates were explicitly pending before publication.
- 194 PASS — five functional commits are reachable from the exact baseline, 5 ahead / 0 behind.
- 195 PASS — remote verification was actually executed.
- 196 PASS — exact build SHA `d7698749fb35030295e2c061c2b98596ff87588c` is recorded.
- 197 PASS after this documentary verification-record commit is fast-forward published and its exact SHA is recorded in PR #2.
- 198 PASS — this canonical post-publication companion records the remote evidence.
- 199 PASS after PR/main/README are rechecked on the documentary verification-record HEAD.
- 200 PASS — EPT-3..EPT-6 remain NOT STARTED.

## Final verdict contract
After the single documentation-only verification-record commit is fast-forward published, its exact SHA and final PR/main/README state are recorded in PR #2 to avoid a self-referential SHA-only commit. With gates 197 and 199 thereby confirmed, **EPT-2 = PASS AFTER POST-PUBLICATION VERIFICATION — 200/200 PASS, 0 PENDING, 0 FAIL**.

Endpoint Capability Specification remains **PARTIAL** because EPT-3..EPT-6 are NOT STARTED. Next candidate is **EPT-3 — Local Detection and Endpoint Investigation**, but EPT-3 is not started by this run.
