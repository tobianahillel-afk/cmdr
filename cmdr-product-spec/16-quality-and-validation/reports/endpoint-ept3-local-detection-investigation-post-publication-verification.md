---
id: endpoint-ept3-local-detection-investigation-post-publication-verification
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-11
source-of-truth: quality-report
requirements: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-018, REQ-PROD-019, REQ-OBJ-008, REQ-INV-001, REQ-INV-006, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008, OPEN-015, OPEN-017]
---
# Endpoint EPT-3 — Post-Publication Verification

## Scope
This companion verifies only **EPT-3 — Local Detection and Endpoint Investigation** under Delivery Roadmap Phase 5 — Studio and Endpoint. It modifies no `CAP-EPT-*` contract and does not start EPT-4, EPT-5 or EPT-6.

## Exact baseline and functional chain
Starting baseline: `5d7c037aff6004984416665e7e188a8700e62b2f` — `docs: record Endpoint EPT-2 post-publication verification`.

1. `e735b57e9baed9d1d740e03492ba73514083a08a` — `docs: establish Endpoint local detection and investigation boundaries`.
2. `f1b314b89b666c900c8932dabfe6c0a74bafacbf` — `docs: define Endpoint detection evaluation signals and coverage semantics`.
3. `51c6b5beb8d322e0d8d19c8e9b27ec813dfbba4b` — `docs: specify Endpoint investigation contexts timelines and pivots`.
4. `0938186f0940a689468b2237ae0c48c4a0cebf6b` — `docs: document Endpoint detection investigation handoff and provenance`.
5. `941bfb5da8a3598ca3dd79d135246b0a8865a31a` — `docs: update Endpoint detection investigation traceability and quality gates`.

Baseline → fifth functional/build SHA was verified at **5 ahead / 0 behind**, with the exact baseline as merge base. Publication used a non-forced fast-forward; no rebase, reset or history rewrite occurred.

## Remote build verification
- remote branch HEAD reached `941bfb5da8a3598ca3dd79d135246b0a8865a31a`;
- PR #2 remained **open / Draft / unmerged**, base `main`;
- `main` remained `bc1ec59e5e79ccbaba291e2984b0eb7d5e54129c`;
- branch/main root README remained exactly `# cmdr`, blob `901c74cda52e28b5ff7fc425ddf28ef89f3ad875`;
- build SHA had no commit statuses and no workflow runs, so **CI/status = N/A**;
- `CAP-EPT-001` and `CAP-EPT-030` were remotely re-read and preserved;
- `CAP-EPT-031` and `CAP-EPT-046` were remotely re-read from the build;
- search for `CAP-EPT-047` returned no result.

## Source audit
Before allocation, all **8/8 Endpoint Detection** and **8/8 Endpoint Investigation** source documents were read. They justify the 16-capability set without artificial fillers. Detection sources justify content/rule/model consumption boundaries, local evaluation, behavioral/local correlation, suppression semantics, context/severity/confidence and health/coverage. Investigation sources justify read-only process tree, host timeline, modules/drivers/system, network connections, user sessions and persistence-like context, while host inspection explicitly stops before new acquisition.

`detection-update.md` remains future update/rollout evidence and is not promoted into EPT-3 implementation.

## Exact capability set
`CAP-EPT-031..046` are the complete EPT-3 set: **16 capability files / 432 numbered sections / 96 mandatory tables / at least 48 Given/When/Then scenarios**, all `draft / defined / planned`, with **0 duplicate ID / 0 recycled ID / 0 owner conflict / 0 empty or generic mandatory table**.

Endpoint cumulative: **46 capabilities / 1242 sections / 276 mandatory tables**. Global: **431 capabilities / 429 defined / 2 proposed / 431 planned / 11637 sections / 2586 mandatory tables**.

## Detection Content boundary
Investigate Detection Engineering remains owner of Detection Content authoring, validation and lifecycle. Platform Settings retains administrative runtimes/targets/sources/policies/assignments. Endpoint consumes a content/version reference, calculates local eligibility/applicability and evaluates locally when prerequisites are satisfied. `OPEN-017` remains open, so EPT-3 selects no Sigma/YARA/proprietary final runtime, target language, compilation representation, model runtime or portability model.

Command remains owner of canonical `Detection`, `Signal`, `Alert` and `Incident`. Endpoint Local Detection Evaluation, Detection Match and Local Detection Signal Candidate are local technical concepts only and do not compete with those canonical objects.

## Evaluation, match and signal semantics
EPT-3 defines content/version + eligible inputs → local evaluation → `matched / not-matched / inconclusive / error` conceptual outcomes. Match keeps observation/content/version/time/rationale provenance. Match is not a malicious verdict, Finding, Evidence or Incident; no-match is not benign.

A Local Detection Signal Candidate can be created/repeated/updated/grouped/suppression-projected/expired/invalidated/forwarded. It is not the canonical Command Signal. Suppressed does not mean deleted; grouped does not mean identical; deduplicated does not mean erased.

## Context, severity, confidence and rationale
Detection Context links only available process/file/network/user/system observations. Severity is contextual technical severity, not automatic business impact. Confidence is not certainty or necessarily a probability. Rationale is not proof. Missing/restricted/stale context remains explicit.

## Grouping, suppression, coverage and gaps
Grouping/deduplication retains all source matches and provenance. Suppression is only a projection from an authorized owner/source with scope/reason/expiry; Endpoint invents no suppression policy.

Detection Coverage/Health compares eligible content requirements with available telemetry/capabilities and evaluation health. Detection health is distinct from Endpoint health. Healthy runtime/source does not prove complete coverage. Coverage does not mean complete visibility. A missing signal does not prove absence of malicious activity.

## Local Endpoint Investigation
- process context: parent/children/ancestry/executable/user/session/modules/file/network refs using existing observations only;
- file context: paths/metadata/hashes only when already observed, with process/detection chronology;
- network context: connection/listener/direction/process/DNS/time refs, without packet capture;
- user/session/auth context: session/auth/process/privilege facts without credential material or compromise verdict;
- service/module/driver/system context: source-backed component/system facts and persistence-like context without persistence conclusion;
- Endpoint Contextual Timeline: local projection preserving late events, gaps and clock uncertainty; timeline/correlation never proves causality.

No process ancestry, file hash, network relation, authentication anomaly or service/driver observation is automatically malicious evidence.

## Pivots, context expansion and EPT-4 boundary
EPT-3 supports permission-aware pivots among already available process/file/network/user/session/system/signal/source/timeline context. **Pivot ≠ Collection** and **context expansion ≠ forensic acquisition**. If required bytes, memory, packets or other unavailable data would require acquisition, EPT-3 returns a `collection-required` boundary and stops. No Collection request, Live Response or command is executed.

## Summary and consumer handoff
Endpoint Investigation Summary contains local context, related observations/detection refs/timeline/entity refs, uncertainty, missing visibility, limitations and provenance. It is **not** Evidence, Finding, Case, Decision or Result. Investigate retains analyst qualification. Handoff never grants source permissions or creates a Case/Finding/Evidence automatically.

## Provenance and cross-product contracts
EPT-3 preserves the chain EPT-2 Observation → Detection Content/version → eligibility → local evaluation → match/candidate → context → related observations → local timeline/investigation → summary → consumer. Shared keeps generic Search/Timeline/Linking/Correlation/Trace/Activity/Jobs/Reporting; Studio keeps Tool/Tool Call/Skill/Automation Run; Govern keeps Decision/Approval/Response Run/Result and production-response authority.

## Permissions, AI and security
EPT-3 identifies functional read needs for restricted Detection Content/rationale, local signals, sensitive process command metadata, file paths/hashes, network/DNS, user/session/privilege, system/module/driver context, local timeline, summary and provenance. Cross-tenant reads are denied. No final RBAC/ABAC namespace is selected.

AI is optional. It may explain a local signal, summarize context/timeline, suggest permission-safe pivots, explain coverage gaps and propose an explicitly labelled hypothesis. It cannot declare maliciousness with authority, create Evidence/Finding/Case/Incident/Decision/Result, execute response, hide a signal, remove provenance, invent observations/context or claim causality as fact. Every essential path has a deterministic/manual alternative.

Sensitive command metadata, user/session data, file paths and network identifiers remain classification/masking/least-privilege controlled. No raw secret is introduced.

## Screens, IA and migration
The Screen Register remains **56 active screens** and Endpoint product Screen IDs remain **0**. EPT-3 only enriches functional IA to **Telemetry → Detection → Investigation → future Collection/Live Response** and capability/pivot relationships; no detailed UX, wireframe, button, final filter or column is created.

All 8 Detection and 8 Investigation source files remain preserved. Migration is additive: `CAP-EPT-031..046` become the normative capability layer while source documents remain canonical module/reference sources. No mass deletion or competing physical schema is introduced.

## Requirements, OPEN and non-regression
Requirements remain **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**. OPEN remains **18**. `OPEN-008` and `OPEN-017` remain open; `OPEN-015` remains open where cross-product provenance is referenced.

EPT-1 remains **PASS 190/190** and EPT-2 remains **PASS 200/200** with `CAP-EPT-001..030` intact. Command remains PASS 27; Investigate PASS 243; Govern PASS 47; Studio PASS 68.

## Gates 201–210 final closure
- 201 PASS — conformance report exists.
- 202 PASS — build-time state `204 PASS / 6 PENDING / 0 FAIL` is explicit and preserved historically.
- 203 PASS — remote-dependent gates were explicitly pending before publication.
- 204 PASS — the five functional commits are reachable from exact baseline, **5 ahead / 0 behind**, same merge base.
- 205 PASS — remote verification was actually executed.
- 206 PASS — exact build SHA `941bfb5da8a3598ca3dd79d135246b0a8865a31a` is recorded.
- 207 PASS after this documentation-only verification-record commit is fast-forward published and its exact SHA is recorded in PR #2.
- 208 PASS — this canonical post-publication companion records the evidence.
- 209 PASS after PR/main/README are rechecked on the verification-record HEAD.
- 210 PASS — EPT-4, EPT-5 and EPT-6 remain NOT STARTED.

## Final verdict contract
After this single documentation-only verification-record commit is fast-forward published, its exact SHA and final PR/main/README state are recorded in PR #2 to avoid a self-referential SHA-only commit. With gates 207 and 209 thereby confirmed, **EPT-3 = PASS AFTER POST-PUBLICATION VERIFICATION — 210/210 PASS, 0 PENDING, 0 FAIL**.

Endpoint Capability Specification remains **PARTIAL** because EPT-4..EPT-6 are NOT STARTED. Delivery Roadmap Phase 5, Global Capability Specification and repository maturity remain PARTIAL. The next candidate is **EPT-4 — Collection and Live Response Technical Execution**, but EPT-4 is not started by this run.