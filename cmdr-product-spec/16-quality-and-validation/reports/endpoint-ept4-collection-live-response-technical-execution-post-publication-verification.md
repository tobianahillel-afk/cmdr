---
id: endpoint-ept4-collection-live-response-technical-execution-post-publication-verification
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-11
source-of-truth: quality-report
requirements: [REQ-PROD-006, REQ-PROD-014, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-PROD-052, REQ-PROD-055, REQ-OBJ-004, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008, OPEN-013, OPEN-014, OPEN-015, OPEN-017]
---
# Endpoint EPT-4 — Post-Publication Verification

## Scope
This companion verifies only **EPT-4 — Collection and Live Response Technical Execution** under Delivery Roadmap Phase 5 — Studio and Endpoint. It modifies no `CAP-EPT-*` contract and does not start EPT-5 or EPT-6.

## Exact baseline and functional chain
Starting baseline: `67ea28d221ed70baae83ff0689048685e1aacf74` — final EPT-3 verification record.

1. `737cdad93bd77556aeb7b65b313ff6db87564fdc` — `docs: establish Endpoint collection and Live Response boundaries`.
2. `bed4900740cf2839d58ff0e84dfb2efd2b76e1fa` — `docs: define Endpoint collection acquisition transfer and completeness semantics`.
3. `0ea89358010abe5619c011240560462ae19306bd` — `docs: specify Endpoint Live Response sessions commands and script execution`.
4. `99eb32aec4338a5d85d68399a032ec98dfa992e2` — `docs: document Endpoint technical execution outputs handoffs and provenance`.
5. `32082487b48779434c9ac730b43efd498009825d` — `docs: update Endpoint collection Live Response traceability and quality gates`.

Baseline → fifth functional/build SHA was verified at **5 ahead / 0 behind**, with the exact EPT-3 baseline as merge base. Publication used a non-forced fast-forward; no rebase, reset, force-push or history rewrite occurred.

## Pre-publication history-preservation audit
Before publication, the first detached version of the fifth commit was rejected from publication because a reconstructed historical Mobile block in the global changelog diverged from the baseline. The branch was still untouched. The final fifth functional commit above was rebuilt after restoring the global `CHANGELOG.md` to the exact baseline blob; EPT-4 keeps its dedicated `CHANGELOG-EPT4.md` and the canonical EPT-4 status/roadmap/quality addenda. No erroneous detached fifth commit became reachable from the PR branch.

## Remote build verification
- remote branch HEAD reached `32082487b48779434c9ac730b43efd498009825d`;
- PR #2 remained **open / Draft / unmerged**, base `main`;
- `main` remained `bc1ec59e5e79ccbaba291e2984b0eb7d5e54129c`;
- branch/main root README remained exactly `# cmdr`, blob `901c74cda52e28b5ff7fc425ddf28ef89f3ad875`;
- build SHA had no commit statuses and no workflow runs, therefore **CI/status = N/A**;
- `CAP-EPT-046` was remotely re-read and preserved;
- `CAP-EPT-047` and `CAP-EPT-064` were remotely re-read from the published build;
- search for `CAP-EPT-065` returned no result.

## Source audit
Before allocation, all **7/7 Endpoint Collection** and **9/9 Endpoint Live Response** historical source documents were read and revalidated. EPT-3 collection-required/pivot/summary boundaries, Investigate collection and acquisition contracts, Govern execution handoff, Studio Tool Call/Automation Run, Settings Secret References/configuration, Shared Jobs/Trace and containment/resilience boundaries were also revalidated.

The sources independently justify all 18 capabilities; no filler was introduced. Process/system-state fresh collection is distinct from EPT-3 cached context, network capture is a bounded acquisition family, session file transfer is autonomous, and interactive terminal/shell execution is distinct from generic command request semantics.

## Exact capability set
`CAP-EPT-047..064` are the complete EPT-4 set: **18 capability files / 486 numbered sections / 108 mandatory tables / at least 54 Given/When/Then scenarios**, all `draft / defined / planned`, with **0 duplicate ID / 0 recycled ID / 0 owner conflict / 0 empty or generic mandatory table**.

Endpoint cumulative: **64 capabilities / 1728 sections / 384 mandatory tables**. Global: **449 capabilities / 447 defined / 2 proposed / 449 planned / 12123 sections / 2694 mandatory tables**.

## Collection Request and authority boundary
Canonical `Collection Request` remains **Investigate-owned**. Endpoint consumes request/version/target/scope/purpose, computes local technical eligibility and constructs a technical intake/plan. `eligible != authorized`, `authorized != started`, and `Collection Request != Collection Operation != Decision != Response Run`.

Endpoint never creates a Govern Decision/Approval/Response Run. An effectful or otherwise governed collection can consume an external authority reference, but technical capability availability never grants authority.

## Collection scope and planning
EPT-4 defines bounded target/subject selection, paths/process/system refs, optional time bounds, conceptual volume/impact estimates, prerequisites, exclusions, planned items and expected outputs. The technical plan cannot silently expand the Investigate request.

## File and filesystem collection
File collection uses exact path/scope where authorized and records metadata, content acquisition state, missing/locked/restricted items, before/after metadata or integrity references when available, changed-during-collection and partial states. `file metadata observation != file acquisition` and `file acquisition != file execution`.

## Process and system-state collection
EPT-4 can acquire a fresh read-only process/system snapshot when EPT-3 context is stale or insufficient. Categories may include process metadata/relations, modules, sessions, services/system facts and other explicitly supported observations. This does not authorize terminate/suspend/resume or service mutation.

## Memory acquisition
Memory acquisition is functionally defined only because it is explicitly sourced. It remains platform/capability/privilege/impact dependent, can be full/process scoped, tracks progress/partial/failure/cancel and returns a neutral technical output reference. No acquisition tool, command or format is selected. **Memory acquired != memory analyzed.**

## Network/packet capture
A bounded network-capture capability is justified by the corpus. EPT-4 records interface/context reference, conceptual duration/volume/filter bounds, start/stop requests and confirmations, loss/partial status and output provenance. No packet-capture command/API/filter language is selected. **Packet capture != network compromise or detection verdict.**

## Collection progress, partial/failure/cancel and resumability
Local collection states preserve preparing/queued/waiting/running/partial/failed/timed-out/cancel-requested/cancelled/completed/transfer-pending/offline/unknown distinctions. Retry/resume are eligibility concepts only. `cancel != rollback`, `timeout != confirmed target-side termination`, `retry != duplicate-free guarantee`, and `resume != new collection automatically`.

## Packaging, completeness and integrity
`Collection Item`, `Collection Package` and `Collected Technical Output` are neutral Endpoint technical concepts. Packages retain expected/actual item references, missing/partial markers, size/time/source/acquisition context and source-attributed integrity metadata. No archive format, hash algorithm, PKI or Evidence package is defined. **Integrity metadata != cryptographic proof automatically.**

## Transfer and Investigate handoff
Acquisition and transfer remain separate states. Transfer can complete/partial/fail independently of acquisition and acknowledgement is separate from content validation. `transfer complete != validated content`. Investigate receives neutral technical references and performs any canonical qualification.

## OPEN-014 — Artifact / Attachment / Evidence
`OPEN-014` remains **OPEN**. Endpoint does not decide that collected bytes are canonical Artifact, Attachment, Evidence or Finding. `Collection Item != Evidence`, `Collection Item != Finding`, `Collection Item != Artifact automatically`, and `Collection Package != Evidence package`. Investigate retains canonical qualification.

## Live Response Session
EPT-4 defines an Endpoint Technical Session with exact target binding, requester/initiator/operator, purpose, tenant/environment, capability/support state, constraints and authority references. It is distinct from the Investigate business Live Session, Studio Automation Run and Govern Response Run.

## Session lifecycle
Requested/validating/connecting/active/idle/waiting/degraded/disconnected/reconnecting/closing/closed/failed/expired/cancelled semantics preserve requested versus confirmed states. `session active != unrestricted authority`, `session close != target state restored` and disconnect/reconnect history remains auditable.

## Command requests and interactive execution
`Technical Command Request` records session/origin/operator/target/purpose/execution-mode/parameter refs, masking, capability and authority. It is not a Tool Call or Govern Decision. An accepted request is not yet started; a started command is not necessarily successful.

Interactive execution is independently justified by the terminal source. It records declared runtime/shell dependency, user/privilege/working context, start/state/output/error/timeout/cancel/transcript provenance. No PowerShell/Bash/Python or other shell/runtime is selected and no command catalogue or offensive command is provided.

## Script execution
Scripts remain references to an authorized source/version with runtime dependency and parameter context. EPT-4 provides no script and selects no universal interpreter. `script uploaded != executed`, `script accepted != safe`, and technical success is not an authorized Govern Result.

## Live Response file operations
Bounded file read/download/transfer is included. A temporary upload is included only because the canonical Investigate source explicitly supports bidirectional session transfer; upload is classified as effectful and requires strict authority. Delete/quarantine/restore and destructive file actions remain EPT-5/Govern boundaries.

## Technical outputs, errors, timeout and cancellation
Technical output preserves source attempt, conceptual stdout/stderr-like references, warnings/errors, partial/truncated/redacted state, completion facts, timeout, cancel/stop requests and target-side termination knowledge. `technical output != Govern Result`, `command output != Evidence automatically`, `stop requested != stopped` and `timeout != confirmed termination`.

## Operator control and closure
EPT-4 records operator presence/control, takeover/reconnect/close/cancel/stop requests, session timeout, closure context, outstanding operations and sensitive-command visibility. Operator is not automatically a Govern approver. Closure never claims rollback or state restoration.

## OPEN-015 — execution bridge
`OPEN-015` remains **OPEN**. Endpoint Technical Execution is distinct from Studio Tool Call, Studio Automation Run and Govern Response Run. Cross-product IDs may be correlated, but no identity merger or final bridge is selected.

## Govern bridge
A governed effectful path may be documented as `Action Request → Decision → Response Run → Endpoint technical execution → technical output/status → Govern reconciliation → Result`. Endpoint never creates the Decision, never becomes the Response Run and never directly creates the canonical Result.

## Studio / Settings / Shared boundaries
Studio keeps Tool/Tool Call/Workflow/Automation Agent/Automation Run. Settings keeps Fleet, Endpoint Policy, credentials, Secret References, providers and runtime/admin configuration. Shared keeps generic Jobs, Trace, Activity, Recovery, Reporting and Export/transfer mechanisms. Endpoint records only its target-side technical semantics.

## Permissions and security/privacy
EPT-4 documents functional needs for Collection request/intake/execute/cancel/transfer, sensitive outputs, memory/network acquisition, Live Response session controls, non-mutating/effectful command execution, scripts, file transfer/upload, output reads and provenance. No final RBAC/ABAC namespace is selected. Tenant isolation and cross-tenant denial remain mandatory; sensitive collected content, paths, memory, command/script parameters and outputs use classification/masking/least privilege. Raw secrets are never arbitrary fields; Secret References are consumed when applicable.

## AI / no-AI
AI is optional. It may suggest bounded collection scope, summarize outputs, explain failures and suggest diagnostic next steps. It may not authorize collection/session/execution, fabricate bytes/output, create Evidence/Result, bypass Govern, reveal secrets, hide partial states, delete provenance or claim timeout means termination. Every essential path has a manual/deterministic alternative.

## Screens, IA and migration
The Screen Register remains **56 active screens** and Endpoint product Screen IDs remain **0**. Information Architecture is functional only and now expresses `Telemetry → Detection → Investigation → Collection → Live Response → future Containment`; no terminal design, wireframe, button, shortcut, filter or final column is created.

Migration is additive: the 7 Collection and 9 Live Response source files remain preserved; `CAP-EPT-047..064` become the normative capability layer. Historical destructive process/network/service/file actions are retained as source evidence but bounded to future EPT-5 rather than deleted or silently re-owned.

## Requirements, OPEN and non-regression
Requirements remain **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**. OPEN remains **18**. `OPEN-008`, `OPEN-014`, `OPEN-015` and `OPEN-017` remain open; `OPEN-013` remains open where reversible/effectful controls are classified.

EPT-1 remains **PASS 190/190**, EPT-2 **PASS 200/200**, EPT-3 **PASS 210/210** with `CAP-EPT-001..046` intact. Command remains PASS 27; Investigate PASS 243; Govern PASS 47; Studio PASS 68.

## Gates 211–220 final closure
- 211 PASS — conformance report exists.
- 212 PASS — build-time state `214 PASS / 6 PENDING / 0 FAIL` is explicit and preserved historically.
- 213 PASS — remote-dependent gates were explicitly pending before publication.
- 214 PASS — the five functional commits are reachable from the exact baseline at **5 ahead / 0 behind**, same merge base.
- 215 PASS — remote verification was actually executed.
- 216 PASS — exact build SHA `32082487b48779434c9ac730b43efd498009825d` is recorded.
- 217 PASS after this documentation-only verification-record commit is fast-forward published and its exact SHA is recorded in PR #2.
- 218 PASS — this canonical post-publication companion records the evidence.
- 219 PASS after PR/main/README are rechecked on the verification-record HEAD.
- 220 PASS — EPT-5 and EPT-6 remain NOT STARTED.

## Final verdict contract
After this single documentation-only verification-record commit is fast-forward published, its exact SHA and final PR/main/README state are recorded in PR #2 to avoid a self-referential SHA-only commit. With gates 217 and 219 thereby confirmed, **EPT-4 = PASS AFTER POST-PUBLICATION VERIFICATION — 220/220 PASS, 0 PENDING, 0 FAIL**.

Endpoint Capability Specification remains **PARTIAL** because EPT-5 and EPT-6 are NOT STARTED. Delivery Roadmap Phase 5, Global Capability Specification and repository maturity remain PARTIAL. The next candidate is **EPT-5 — Containment, Verification and Governed Response Primitives**, but EPT-5 is not started by this run.