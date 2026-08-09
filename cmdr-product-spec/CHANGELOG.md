# Changelog

## 2026-08-09 — Current Phase 4A Command corrective revalidation
- Revalidated `Capability Specification Phase 4A — Command` from exact remote SHA `edbc67be43f971e4c49e91518e3e21517ef69410` without recreating any CAP-CMD file or modifying historical Phase 4B functional content.
- Re-read all 27 Command capability contracts, Command architecture/module sources, active Command screens and five deprecated Work Queue aliases; retained 729/729 sections and 162/162 mandatory tables.
- Restored five `CAP-CMD-*` evidence ranges in the active Requirements Traceability Matrix while keeping all 122 Requirement IDs and the 99 conform / 20 partial / 3 absent / 0 contradictory distribution unchanged.
- Replaced generic Command capability-register role/dependency summaries with capability-specific roles, objects, consumers and dependency families; retained 26 defined, 1 proposed and 27 planned.
- Added ten Command dependency families to the global Dependency Register while preserving the detailed 161-edge Command functional dependency map.
- Preserved one Work Queue workspace and exactly six system views: All, Incidents, Tasks, Unassigned, SLA Risk and My Work; Team Load remains non-canonical as a system view.
- Preserved CAP-CMD-401 Customers and Delivery as proposed, planned and deployment-dependent; OPEN-006 remains open. OPEN-010 and OPEN-013 remain open.
- Modified no Command capability, screen, domain object, permission source, API, protocol, font or product code; created no new OPEN decision.
- Preserved the current total of 18 open decisions and all later Investigate/Cloud/Mobile historical evidence.
- Added a current 60-control Phase 4A revalidation report; publication-dependent Git/PR/README controls are verified only after canonical squash publication.

## 2026-08-09 — Phase numbering namespace reconciliation
- Recorded two historical phase namespaces: `Capability Specification Phase` for detailed capability-specification execution and `Delivery Roadmap Phase` for the historical product-delivery sequence.
- Preserved all existing phase IDs, filenames, reports, SHAs, Capability IDs and PASS history; no historical phase was renumbered.
- Preserved `phase-4-govern.md`, canonical id `roadmap-phase-4-govern` and title `Phase 4 Govern` as `Delivery Roadmap Phase 4 — Govern`.
- Explicitly prohibited creating or renaming Govern as `Phase 4C Govern` solely to align with Capability Specification Phase 4A/4B.
- Required qualified phase names in future ambiguous contexts and prohibited inferring functional dependencies from numeric proximity across namespaces.
- Confirmed `Capability Specification Phase 4A — Command` PASS and `Capability Specification Phase 4B — Investigate` PASS while `Capability Specification Phase 4` global maturity remains PARTIAL.
- Recorded Govern capability specification as NOT STARTED; created no CAP-GOV ID, capability, object, screen, atomic permission, API, protocol, implementation or new product decision.
- Kept the open-decision count unchanged at 18; no OPEN was created or closed.

## 2026-08-07 — Phase 4B.4B Mobile Forensics closure preparation
- Verified the exact Mobile starting remote head `ed874ea414fc57f24fa61f410f91b7345f4a868a`, PR #2 Draft/open/unmerged, public repository, disabled auto-merge and unchanged root README before modification.
- Added the unique canonical `Phase 4B.4B — Mobile Forensics Foundations and Mobile Investigation` roadmap and provider/platform-neutral module.
- Added exactly 19 canonical capabilities `CAP-INV-701..719`, each `draft` / `defined` / `planned`, with 27 sections and six mandatory tables.
- Added 14 Mobile supporting module documents covering scope, concepts, workflows, states, permissions, privacy, Shared/AI/Collection/Settings/Endpoint/Studio/Govern boundaries, map and source migration.
- Preserved Collection ownership of acquisition requests/jobs/execution/results/custody; Settings ownership of sources/Fleet/credentials/secrets/retention/policies; Endpoint declared-capability ownership; Govern real-device authority; Studio Tools/Runs; Shared generic engines.
- Defined Mobile analysis for device/platform scope, acquisition-context review, integrity/completeness/accessibility, filesystem/storage, applications/data, communications, media, location/sensors, sensitive material, connectivity/SIM/pairing, backups/synchronization, deleted/recovered data, Timeline/correlation, Hypotheses, Derived Artifacts/handoffs and provenance.
- Expanded `OPEN-011 — Mobile Forensics scope and delivery strategy` while keeping it open with no platform/version, acquisition method, Tool, engine, API, protocol, proprietary format or implementation selected.
- Updated dependencies, Object Consumption, action classes, AI boundaries, cross-product links and conceptual screen mapping without creating a canonical Mobile object schema, atomic permission matrix or Screen ID.
- Recalculated **270 global capabilities, 243 Investigate, 268 defined, 2 proposed, 270 planned; 6561 Investigate sections / 1458 tables; 7290 Command + Investigate sections / 1620 tables**.
- Recalculated Phase 4B.4 Cloud + Mobile at **37 capabilities, 999 sections and 222 tables**.
- Preserved 122 Requirement IDs and the **99 conform / 20 partial / 3 absent / 0 contradictory** distribution.
- Prepared a Mobile conformance report with **250 gates**; publication-dependent gates remain pending until the fifth functional commit is published and rechecked remotely.
- Kept Phase 4B.4B, Phase 4B.4 and Phase 4B pending/partial until post-publication verification; Phase 4/global maturity remains PARTIAL regardless.
- Added no mobile integration, connector, API, protocol, platform/tool/acquisition engine, unlock/bypass/root/jailbreak technique, command, script, device mutation, secret use, deployed rule, response or product code.

### Mobile functional commits
1. `2626cdb3b502110ef7e06a9f66997b81bed5bacd` — `docs: establish Investigate mobile forensics boundaries`.
2. `9f18ea128bb5cd8dbeb00f38a5dce7c05d64c37a` — `docs: define mobile acquisition storage and application analysis`.
3. `b38e40d4385f013e409d2001321e463babe0a2d1` — `docs: specify mobile communications location and sensitive data analysis`.
4. `78d49e3fb8895a4fb9f46bd1f4f7a28fcb4d8a52` — `docs: document mobile correlation evidence handoff and provenance`.
5. `docs: close Mobile Forensics and Investigate capability phases` — exact squash SHA recorded only after publication verification.

## 2026-08-07 — Phase 4B.4A Cloud Analysis verified closure
- Published the fifth functional commit `6baee257a0ce81b4c50b38aba3805a61eb5daa6c` — `docs: update Cloud Analysis traceability and quality gates`.
- Verified the functional Cloud chain at five commits ahead and zero behind `6cc7bf2426c12a6f2f983c62581f9d44e9829990`, with the same merge base.
- Verified PR #2 open, Draft and unmerged; repository public; auto-merge disabled; branch and `main` README both exactly `# cmdr` with the same blob.
- Verified no GitHub Actions workflow run and no commit status on the fifth functional SHA.
- Promoted Phase 4B.4A to `PASS AFTER POST-PUBLICATION VERIFICATION` with 243/243 gates.
- Retained 18 capabilities, 486 sections, 108 tables, 251 global capabilities, 224 Investigate capabilities, 249 defined, 2 proposed and 251 planned.
- Retained OPEN-012 open for provider, service and delivery strategy without selecting a provider, connector, API, protocol, schema, query language, engine or architecture.
- Kept Phase 4B.4, Phase 4B, Phase 4 and global maturity PARTIAL because Mobile Forensics is NOT STARTED.
- Added no capability, code, provider integration, target mutation, secret use, active scan, deployed rule, response or Mobile content.
- Added a fast-forward verification-record correction titled `docs: record Cloud Analysis post-publication verification`; its final SHA is recorded in PR #2 and the final report rather than creating a further SHA-only commit.

## 2026-08-07 — Phase 4B.4A Cloud Analysis closure preparation
- Preserved the four published Cloud functional commits and `CAP-INV-601..618` without changing IDs, owners or functional scope.
- Added the Cloud capability-register shard and recalculated 251 global capabilities, 224 Investigate capabilities, 249 defined, 2 proposed and 251 planned.
- Recalculated Investigate at 6048 sections and 1344 mandatory tables; Command plus Investigate at 6777 sections and 1506 tables.
- Added provider-neutral Cloud dependencies, object-consumption projections, action classes, AI boundaries, cross-product handoffs and conceptual screen mapping.
- Expanded `OPEN-012 — Cloud Analysis` for provider priorities, organization/tenant/account/subscription/project scope, services, multi-cloud, log models, inventory, cloud-native workloads, containers, orchestration, serverless, SaaS, Cloud Evidence, cross-account, cross-tenant, collection, coverage and limitations.
- Kept OPEN-012 open with no provider, engine, connector, protocol, schema, query language or architecture selected.
- Updated the Requirements Matrix without changing the 122 Requirement IDs or the 99 conform / 20 partial / 3 absent / 0 contradictory distribution.
- Preserved the full Threat Intelligence baseline and added the Cloud before/after measures and verified source audit.
- Prepared a 243-gate Cloud conformance report with publication-dependent gates marked pending.
- Recorded the five intended functional commit titles; the fifth commit was this traceability and quality closure lot once squash-published.
- Kept Phase 4B.4A `PENDING POST-PUBLICATION VERIFICATION` until remote checks completed.
- Kept Phase 4B.4, Phase 4B, Phase 4 and global maturity PARTIAL; Mobile Forensics remained NOT STARTED.
- Added no product code, provider integration, connector, command, active scan, secret use, target mutation, deployed rule, response or Mobile content.

### Cloud functional commits
1. `779d39a3198c02c50baebb26648cd2edb3c548d2` — `docs: establish Investigate cloud analysis boundaries`.
2. `7e569536db16f55382d3c393b876423cd00b10c7` — `docs: define cloud scope identity permissions and audit analysis`.
3. `4f5b5aea6447159f0e836e5bd7d5e42e756865f2` — `docs: specify cloud workloads network storage and sensitive data analysis`.
4. `59f69203574a6b1e3dc8c631b3a8a2e4bf5afe42` — `docs: document cloud correlation evidence handoff and provenance`.
5. `6baee257a0ce81b4c50b38aba3805a61eb5daa6c` — `docs: update Cloud Analysis traceability and quality gates`.

## 2026-08-06 — Post-publication traceability preservation correction
- Restored the complete pre-4B.3B.2 ADR links, Requirement IDs and detailed Dependency Register rows before retaining the new 4B.3B.2 dependencies.
- Preserved the complete 4B.3B.1 source audit and before/after baseline alongside the 4B.3B.2 measurements.
- Preserved prior requirement evidence ranges while adding CAP-INV-519..537 evidence.
- No capability, status, owner, implementation claim or phase verdict changed.

## 2026-08-06 — Phase 4B.3B.2 Intelligence Analysis, Dissemination and Operationalization
- Added CAP-INV-519 through CAP-INV-537 and the `analysis-and-products/` functional architecture.
- Defined Analysis Sessions, analytic questions, competing hypotheses, multi-source fusion, corroboration and structured analysis.
- Defined Threat Actor/attribution, Campaign/activity and malware/tool/infrastructure/capability assessments without automatic attribution.
- Defined Intelligence Product planning, audience, authoring, versioning, review, quality control and Release Recommendation.
- Defined handling, markings, releasability, Dissemination Plans and reversible internal publication without raw-source permission inheritance.
- Defined functional Watchlist Definitions and Indicator Operationalization Packages without activation, deployment, rule creation, Signal generation or runtime mutation.
- Defined Intelligence monitoring, Sighting/change updates, external-sharing preparation, consumer feedback, effectiveness, Requirement satisfaction, collection feedback, correction, retraction, supersession and lifecycle provenance.
- Created `OPEN-019 — Intelligence dissemination, releasability, sharing and consumer access policy` as open; selected no final policy or destination.
- Preserved Shared Report/Entity/Graph/Reporting, Command runtime objects, Detection Engineering content, Settings administration, Studio Tools/Runs and Govern authority.
- Verified 19/19 capabilities, 513/513 sections, 114/114 mandatory tables and 180/180 gates after publication verification.
- Recalculated Threat Intelligence at 37 capabilities, 999 sections and 222 tables.
- Recalculated Investigate at 206 capabilities, 5562 sections and 1236 tables.
- Recalculated global capabilities at 233; 231 defined, 2 proposed, all planned.
- Marked Phase 4B.3B.2, Phase 4B.3B and Phase 4B.3 PASS.
- Kept Phase 4B PARTIAL because OPEN-011 Mobile Forensics and OPEN-012 Cloud Analysis remained open and not explicitly deferred.
- Added no API, protocol, standard, provider, code, active watchlist, deployed Indicator, rule, block, response, external sharing, Cloud or Mobile content.

## 2026-08-06 — Phase 4B.3B.1 Threat Intelligence Foundations and Knowledge Management
- Added CAP-INV-501 through CAP-INV-518.
- Defined Intelligence Intake, Requirements and priorities, Knowledge Projects, source catalog/access context, source reliability and information credibility.
- Defined Intelligence Material reference/intake/normalization with original, restrictions, ambiguity and Tool/Run provenance.
- Defined Observable/Indicator and Threat Entity candidates, Malware/Tool/Capability and Infrastructure knowledge, Campaign/Activity Cluster/Intrusion Set candidates and TTP mappings.
- Defined Sightings, sourced relationship graph semantics, confidence/contradictions, deduplication/versioning/supersession, expiration/revocation and Analysis Handoff provenance.
- Created `OPEN-018 — Threat intelligence ontology, interoperability and exchange strategy` as open; selected no ontology, standard, protocol, provider or representation.
- Preserved Shared Entity/Graph/Linking, Command Detection/Signal/Alert/Incident, Detection Engineering Content, Settings source administration, Studio Tools/Runs and Govern external-release authority.
- Deprecated no competing Investigate functional source and rewrote no detailed screen.
- Added no API, protocol, STIX/TAXII-like contract, physical schema, provider, scraper, active collection, Indicator deployment, watchlist, rule, block, response, external sharing, 4B.3B.2, Cloud or Mobile content.
- Verified 18/18 capability files, 486/486 sections, 108/108 mandatory tables and 170/170 gates after publication verification.
- Recalculated Investigate at 187 capabilities, 5049 sections and 1122 tables.
- Recalculated Command plus Investigate at 214 capabilities, 5778 sections and 1284 tables.
- Marked Phase 4B.3B.1 PASS after publication verification.

## 2026-08-06 — Phase 4B.3A.2 Detection lifecycle and Phase 4B.3A closure
- Added CAP-INV-418 through CAP-INV-435 and closed Detection Engineering as PASS.
- Created OPEN-017 open without selecting a runtime or language.

## 2026-08-06 — Phase 4B.3A.1 Detection Engineering Foundations, Authoring and Validation
- Added CAP-INV-401 through CAP-INV-417 and completed authoring/validation evidence.

## 2026-08-05 — Phase 4B.2 closure
- Added CAP-INV-380 through CAP-INV-397 and closed Phase 4B.2 as PASS.

Earlier phases remain preserved.