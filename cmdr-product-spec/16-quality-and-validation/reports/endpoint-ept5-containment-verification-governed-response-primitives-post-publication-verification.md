---
id: endpoint-ept5-containment-verification-governed-response-primitives-post-publication-verification
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-11
source-of-truth: quality-report
requirements: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-006, REQ-PROD-008, REQ-PROD-014, REQ-PROD-015, REQ-PROD-016, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-014, OPEN-015, OPEN-017]
---
# Endpoint EPT-5 — Post-Publication Verification

## Scope
This companion verifies only **EPT-5 — Containment, Verification and Governed Response Primitives**. It changes no `CAP-EPT-*` capability contract and does not start EPT-6.

## Exact baseline and functional chain
Starting baseline: `5d576295fa12693ef375a35cfe515d7bdf577f68` — `docs: record Endpoint EPT-4 post-publication verification`.

1. `cdd4517fabc5961f5c32adbf04843cbafae3b972` — `docs: establish Endpoint containment and governed response primitive boundaries`.
2. `fd97302149f398c43a86e070b9234cd4d5eac2ec` — `docs: define Endpoint process network file and system containment primitives`.
3. `6a951cc5db76aed07f04ac46f5ce003748fc9b7f` — `docs: specify Endpoint response verification failure and target-state semantics`.
4. `ceaff139227adb14eb6f65ff821a142dc61a7279` — `docs: document Endpoint reversal reconciliation and response provenance`.
5. `b346491d4f09541b0064db1e1ec4764f113804ce` — `docs: update Endpoint containment response traceability and quality gates`.

Baseline → fifth functional/build SHA was remotely verified at **5 ahead / 0 behind**, with the exact baseline as merge base. Publication used a non-forced fast-forward; no rebase/reset/force-push/history rewrite occurred.

## Remote build verification
- PR #2 remote HEAD reached `b346491d4f09541b0064db1e1ec4764f113804ce`;
- PR #2 remained **open / Draft / unmerged**, base `main`;
- `main` remained `bc1ec59e5e79ccbaba291e2984b0eb7d5e54129c`;
- branch/main root README remained exactly `# cmdr`, blob `901c74cda52e28b5ff7fc425ddf28ef89f3ad875`;
- build SHA had no commit statuses and no workflow runs, therefore **CI/status = N/A**;
- baseline → build changed no `CAP-EPT-001..064` capability file;
- `CAP-EPT-064` was remotely re-read and preserved;
- `CAP-EPT-065` and `CAP-EPT-081` were remotely re-read from the published build;
- repository search for `CAP-EPT-082` returned no result;
- Screen Register remained **56 active screens** with **0 Endpoint Screen IDs**;
- unresolved decisions remained **18 OPEN**, including OPEN-007/008/013/014/015/017.

## Source audit and final capability set
All **7/7 Endpoint Containment** historical documents were read: README, account containment, host containment, network isolation, quarantine, rollback and verification. Effectful historical Live Response process/network/file/service sources were re-read. Govern Action Request/Approval/Decision/Response Run/technical handoff/Verification/Rollback/Result, Studio Human Gate/Tool Call/Automation Run, Settings Policy/Fleet/Secrets/admin, Investigate containment handoff, Shared generic mechanisms, Security and the EPT-6 resilience boundary were revalidated.

The 16 candidate responsibilities are independently supported. `containment/account-containment.md` additionally and independently sources local session termination/lock while declaring directory actions external. Therefore the actual EPT-5 set is **17 capabilities `CAP-EPT-065..081`**, not a forced 16-ID quota.

Structural result: **17 capability files / 459 numbered sections / 102 mandatory tables / at least 51 GWT**, all `draft / defined / planned`, with **0 duplicate / 0 recycled / 0 owner conflict / 0 empty or generic mandatory table**.

Endpoint cumulative: **81 capabilities / 2187 sections / 486 mandatory tables**. Global: **466 capabilities / 464 defined / 2 proposed / 466 planned / 12582 sections / 2796 mandatory tables**. Requirements remain **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**.

## Response Primitive and authority boundary
`Response Primitive` is Endpoint technical vocabulary, not an Action Request. Endpoint computes eligibility/availability/readiness and consumes Decision/Response Run/authority references. `eligible != authorized`, `available != allowed`, authority-reference presence does not create authority, and Endpoint Policy is not authority.

Govern retains Action Request, Approval, Decision, Response Run, production-response authority, response-level verification, rollback/recovery governance and canonical Result.

## Precheck and target-state planning
Endpoint Technical Precheck records current target facts, requested technical state, freshness, prerequisites, capability availability, Policy projections, expected technical impact, reversibility information and verification-path references. `precheck PASS != universally safe` and requested state remains distinct from achieved state.

## Process control
CAP-EPT-067 covers source-backed suspend/resume/terminate semantics with exact process identity and optional bounded tree scope. No OS command is selected. Process disappear/reused identity/partial outcomes are explicit. `terminate requested != terminated`, `process absent != threat resolved`, `suspended != terminated`, `resume != rollback`.

## Host isolation and bounded network control
CAP-EPT-068 defines host isolation as a target-side technical state; endpoint unreachable/offline/network failure remain distinct. CAP-EPT-069 separately represents bounded network blocking/unblocking because Live Response sources independently define target/scope/expiry/conflict/partial enforcement. No firewall syntax or implementation detail is introduced.

## File quarantine, delete and restore
CAP-EPT-070 preserves file identity before quarantine and distinguishes requested/applied/missing/changed/partial states. `quarantined != deleted` and `quarantined != malicious verdict`. CAP-EPT-071 separately covers source-backed destructive delete and restore/recovery boundaries with stronger authority, retained provenance and uncertainty; `deleted != guaranteed unrecoverable`, `restored != safe`.

## Service/system component control
CAP-EPT-072 covers only declared/sourced start/stop/restart/disable-like technical semantics where platform capability exists. Unsupported platform is not fabricated as execution failure. `service stopped != remediation complete`. EPT-6 watchdog/persistence/self-protection are not imported.

## Local account-session containment
CAP-EPT-081 is the source-driven seventeenth capability. It covers **local endpoint session lock/termination only**. Directory account disable/delete/reset, cloud-token revocation and provider-side identity mutation remain external. Local session termination does not prove credential revocation or compromise resolution.

## Technical execution outcome
CAP-EPT-073 normalizes target-side states while preserving raw primitive status. `technical success != Response Run success`, `technical failure != response failure automatically`, `Technical Outcome != Result`, `cancel != rollback`, timeout/unknown remain non-binary.

## Endpoint Technical Verification versus Govern Verification
CAP-EPT-074/075 define **Endpoint Technical Verification** only: criteria/source/freshness for technical target-state observations and a source-backed comparison of expected vs observed target state. Govern retains its Verification Plan/Assessment, residual-risk reasoning and response-objective determination.

Endpoint may report process absent, isolation observed, quarantine state, service state, local session state, match/mismatch/partial/unknown/stale/unverifiable. Endpoint never automatically declares incident contained, remediation successful, response complete or business impact resolved.

## Partial, failure, unknown and drift
CAP-EPT-076 preserves partial execution, contradictory/stale/lost/unknown state, verification unavailable/mismatch, target change and post-action drift/regression. Unknown is never coerced to success/failure. Drift does not rewrite historical execution outcome or automatically prove execution failure.

## Technical reversal and Govern rollback
CAP-EPT-077 exposes an inverse technical primitive only where source/capability supports it and after current-state preconditions. It consumes Govern Rollback Plan/authority and reports raw reversal outcome. **Endpoint Technical Reversal != Govern Response Rollback**; local compensation does not complete Govern rollback.

## Containment release / restore to service
CAP-EPT-078 defines technical release of isolation/block/quarantine/control under Govern authority. `release technically complete != restored pre-incident state`, `connectivity restored != secure`, and Endpoint does not decide safe return to production.

## Govern reconciliation and Result
CAP-EPT-079 sends raw technical action/target/requested/observed state, partial/failure/unknown/drift, reversal/release and technical verification refs back to the Govern Response Run. Govern reconciles these facts, performs response-level verification and may create Result. Endpoint never creates Result.

## Provenance and cross-product contracts
CAP-EPT-080 reconstructs `Action Request → Approval where required → Decision → Response Run → primitive → authority/precheck → technical execution → Technical Outcome → target observation/technical verification → optional reversal/release → Govern reconciliation → Result`, retaining the canonical owner at every hop.

Human Gate remains Studio-owned and not Approval/Decision. Endpoint action remains distinct from Tool Call/Automation Run/Response Run. Shared Trace/Activity/Jobs/Recovery/Reporting remain mechanisms, not response authority or fact owners.

## OPEN decisions
All **18 OPEN** decisions remain open unless historically resolved. EPT-5 closes none. OPEN-007 Human Gate/Govern, OPEN-008 platform/source support, OPEN-013 Class-2 default governance, OPEN-014 Artifact/Attachment, OPEN-015 Automation Run/Response Run bridge and OPEN-017 detection runtime/language remain open.

## Action classes, permissions and AI
Class 0 = inspect; Class 1 = eligibility/precheck/no-effect verification; Class 2 only for genuine no-effect/bounded preparation where existing rules allow; Class 3 default for significant target effects; potentially irreversible file effects may require Class 4 treatment. OPEN-013 prevents arbitrary downgrading.

Functional permissions distinguish primitive read, precheck, process suspend/resume/terminate, host isolation/release, network block/unblock, quarantine/release, file delete/restore, service/system control, local session control, technical verification, reversal and sensitive provenance. Step-up, SoD, masking, tenant isolation, Govern dependency and Settings dependency are explicit. No final RBAC/ABAC is selected.

AI is optional. It can explain readiness, suggest source-backed candidates, summarize technical output/verification mismatch/drift/reversal. It cannot authorize, create Approval/Decision, execute effectful primitives autonomously, bypass Govern, invent state/verification/Result, hide failure, destroy provenance or silently choose rollback/release. Deterministic/manual alternatives exist.

## Screens, IA and migration
Endpoint Screen IDs remain **0**. IA extends conceptually to `Telemetry → Detection → Investigation → Collection → Live Response → Containment → Verification → Govern reconciliation → future EPT-6`. No kill-process UI, isolation final screen, buttons, terminal, wireframe, filter, columns or detailed UX is defined.

Migration is additive. Historical Containment and Live Response sources remain. Stale historical `Result` wording in Endpoint rollback/verification sources is reconciled to **Technical Outcome/Observation**, because canonical Result is Govern-owned.

## Non-regression
- EPT-1: **190/190 PASS**;
- EPT-2: **200/200 PASS**;
- EPT-3: **210/210 PASS**;
- EPT-4: **220/220 PASS**;
- `CAP-EPT-001..064`: intact;
- Command: **PASS 27**;
- Investigate: **PASS 243**;
- Govern: **PASS 47**;
- Studio: **PASS 68**;
- EPT-6: **NOT STARTED**.

## Gates 221–230 final closure contract
- 221 PASS — conformance report exists.
- 222 PASS — historical build-time state `224 PASS / 6 PENDING / 0 FAIL` is explicit.
- 223 PASS — remote-dependent gates were explicitly pending before publication.
- 224 PASS — five functional commits are remotely reachable from the exact baseline at 5 ahead / 0 behind, same merge base.
- 225 PASS — remote verification was actually executed.
- 226 PASS — exact build SHA `b346491d4f09541b0064db1e1ec4764f113804ce` is recorded.
- 227 PASS after this documentation-only verification-record commit is fast-forward published and its exact SHA is recorded in PR #2.
- 228 PASS — this companion is the post-publication evidence.
- 229 PASS after PR/main/README are rechecked on the verification-record HEAD.
- 230 PASS — EPT-6 remains NOT STARTED.

## Final verdict contract
After this single documentation-only verification-record commit is fast-forward published, its exact SHA and final PR/main/README state are recorded in PR #2 to avoid a self-referential SHA-only commit. With gates 227 and 229 confirmed, **EPT-5 = PASS AFTER POST-PUBLICATION VERIFICATION — 230/230 PASS, 0 PENDING, 0 FAIL**.

Endpoint Capability Specification remains **PARTIAL** because EPT-6 is NOT STARTED. Delivery Roadmap Phase 5, Global Capability Specification and repository maturity remain PARTIAL. Next candidate: **EPT-6 — Updates, Resilience, Security and Endpoint Provenance**; EPT-6 is not started by this run.