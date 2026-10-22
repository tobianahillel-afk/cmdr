---
id: endpoint-ept4-source-audit
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-11
source-of-truth: quality-report
---
# Endpoint EPT-4 — Source Audit

## Scope
Audit read-before-write for **EPT-4 — Collection and Live Response Technical Execution** from exact baseline `67ea28d221ed70baae83ff0689048685e1aacf74`. No EPT-5/EPT-6 capability is included.

## Endpoint Collection — 7/7 read
1. `collection/README.md` — governed acquisition, integrity/custody, chunking/resume, secure transfer.
2. `collection/artifact-collection.md` — logs/registry/browser/execution/triage categories, profiles, minimization, manifest; term `Artifact` treated as historical source language only under OPEN-014.
3. `collection/chain-of-custody.md` — requester, authority, collector, time, method, hashes, transfers, verification, legal hold.
4. `collection/collection-queue.md` — priority, bandwidth, pause/resume, offline queue, expiry.
5. `collection/file-collection.md` — exact path/scope, hashes before/after, locked and sparse/large files.
6. `collection/memory-collection.md` — full/process memory, platform/privilege, encryption, impact preview.
7. `collection/triage-package.md` — versioned profile, collected items, failures, manifest/hashes.

## Endpoint Live Response — 9/9 read
1. `live-response/README.md` — interactive/scripted response, signed commands, target context, per-command audit.
2. `command-execution.md` — allow-list/policy, arguments/environment, timeout, exit/output/result terminology.
3. `file-actions.md` — move/copy/delete/restore/quarantine integration; only bounded read/download/transfer boundary enters EPT-4; destructive/quarantine semantics are deferred.
4. `network-actions.md` — temporary block rules are EPT-5 containment boundary, not EPT-4 authority.
5. `process-actions.md` — terminate/suspend/resume are EPT-5 effect/containment boundary.
6. `script-execution.md` — signed script reference, interpreter/version dependency, inputs, output/cleanup; no script/runtime is selected.
7. `service-actions.md` — start/stop/restart/disable are effectful response boundaries deferred to EPT-5.
8. `session-management.md` — open/close, idle timeout, concurrent sessions, emergency termination.
9. `terminal.md` — PTY/shell context, user/privilege, session timeout, transcript/redaction; justifies a distinct interactive shell capability without choosing shell/protocol.

## EPT-3 boundaries revalidated
`CAP-EPT-043` and `CAP-EPT-044` return `collection-required` when existing local context is insufficient. Pivot != Collection; Context Expansion != acquisition. `CAP-EPT-045` is summary/handoff only and creates no Evidence/Finding/Decision/Result.

## Investigate revalidated
`Collection Request` is canonical **Investigate-owned**. CAP-INV-202 owns business scope/request preparation; CAP-INV-203 keeps Collection Job distinct from Shared Background Job; CAP-INV-204 supports versioned triage profiles; CAP-INV-205 bounded file/directory acquisition; CAP-INV-206 fresh process/system read-only snapshot; CAP-INV-207 memory acquisition request; CAP-INV-208 bounded network capture; CAP-INV-209 business Live Session; CAP-INV-210 interactive operation selection; CAP-INV-211 bidirectional session file transfer including temporary upload; CAP-INV-212 Operation Result != Govern Result; CAP-INV-213 integrity/custody; CAP-INV-214 cross-product provenance.

## OPEN-014
`OPEN-014 — Artifact/Attachment/dataset/material relations/retention` remains open. EPT-4 therefore uses neutral Endpoint concepts `Collection Item`, `Collected Technical Output` and `Collection Package`. No collected output becomes Artifact, Attachment, Evidence or Finding automatically. Investigate performs any canonical qualification.

## OPEN-015
`OPEN-015 — Tool Call/Automation Run and Response Run/cross-product provenance bridge` remains open. Endpoint Technical Execution != Tool Call != Automation Run != Response Run. Cross-product IDs may be correlated without identity/ownership merger.

## Govern / Studio / Settings / Shared
Govern owns Decision, Approval, Response Run, Result, verification/rollback authority and bounded execution handoff. Studio owns Tool Call/Automation Run lifecycle and technical outcomes. Settings owns Fleet/Policy/providers/credentials/Secret References/admin runtime configuration. Shared owns generic Jobs/Trace/Activity/Recovery/Export/Reporting. Endpoint owns target-side technical execution facts only.

## Platform/security/containment/resilience
OPEN-008 remains open: Windows/Linux/macOS are release candidates, not delivered support claims. OPEN-017 remains open: no detection/script universal runtime is selected. Local audit is append-oriented. Containment is explicitly under Govern authority and belongs to EPT-5. Deep resilience/update/anti-tamper/self-protection belong to EPT-6.

## Capability-set decision
The actual corpus independently justifies all 18 candidates `CAP-EPT-047..064`: process/system live snapshot is distinct from EPT-3 cached context; Network Capture Request is autonomous; Session File Transfer is autonomous; terminal/interactive shell is distinct from generic command invocation. No artificial filler is required.

## Explicit exclusions
No isolation, quarantine, process termination as response primitive, service/network mutation, destructive remediation, containment verification, Govern rollback, update/resilience/security implementation, API/protocol/transport, final command catalog, forced PowerShell/Bash/Python runtime, physical output/evidence schema, storage engine, transfer protocol, final RBAC/ABAC or Endpoint Screen ID.