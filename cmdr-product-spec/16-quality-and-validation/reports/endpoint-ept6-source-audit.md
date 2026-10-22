---
id: endpoint-ept6-source-audit
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-12
source-of-truth: quality-report
requirements: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-006, REQ-PROD-008, REQ-PROD-012, REQ-PROD-014, REQ-PROD-015, REQ-PROD-016, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-OBJ-007, REQ-OBJ-008, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004, REQ-SEC-005]
open_decisions: [OPEN-004, OPEN-007, OPEN-008, OPEN-013, OPEN-014, OPEN-015, OPEN-017]
---
# Endpoint EPT-6 — Source Audit

## Scope
This audit covers only **EPT-6 — Updates, Resilience, Security and Endpoint Provenance**, the sixth and final Endpoint Capability Specification execution lot under Delivery Roadmap Phase 5 — Studio and Endpoint. It does not start Delivery Roadmap Phase 6 and introduces no product implementation.

## Exact baseline
- remote branch: `docs/cmdr-product-spec-foundation`;
- baseline SHA: `4d20e8e411aa8919f31e773a9c5ea96efd64791e`;
- baseline title: `docs: record Endpoint EPT-5 post-publication verification`;
- PR #2: open / Draft / unmerged, base `main`;
- `main`: `bc1ec59e5e79ccbaba291e2984b0eb7d5e54129c`;
- Endpoint before EPT-6: 81 capabilities / 2187 sections / 486 mandatory tables;
- global before EPT-6: 466 capabilities / 464 defined / 2 proposed / 466 planned / 12582 sections / 2796 mandatory tables;
- Endpoint Screen IDs: 0;
- OPEN decisions: 18.

## Namespace preflight
Repository search for `CAP-EPT-082` returned no result before allocation. `CAP-EPT-001..081` are treated as immutable prior-lot contracts. No concrete or reserved `CAP-EPT-082+` was identified in the active namespace.

## Endpoint Update sources read
- `11-endpoint-agent/resilience/update-and-rollback.md` — signed updates, rollback slot, health gate, recovery mode;
- `11-endpoint-agent/detection/detection-update.md` — signed package, staged rollout, rollback and offline expiry context for detection content;
- `17-implementation-contracts/endpoint-agent-update-contract.md` — signed artifact, compatibility, waves, health gate, rollback;
- EPT-1 `CAP-EPT-005` — current observed version/build and compatibility context while explicitly excluding the update lifecycle;
- Platform Settings Fleet/upgrade-management — administrative waves, compatibility, rollback and offline handling;
- Platform Settings policy assignment — Fleet/environment/staged rollout administration.

### Update ownership result
Platform Settings retains administrative update targeting, channels/waves, desired version, Fleet selection, policy assignment and administrative package/provider configuration. Endpoint owns only local assigned-target projection, package/release reference consumption, download/staging readiness, installation/activation, local progress/failure/defer/retry facts, post-update health/compatibility observation and previous-version recovery facts.

`Settings upgrade wave != Endpoint local update execution`.
`Endpoint Agent Update != Studio Asset Deployment`.
`Endpoint Update Reversion != Studio Deployment Reversion`.
`Endpoint Update Recovery != Govern Response Rollback`.

## Endpoint Resilience sources read
The complete active `11-endpoint-agent/resilience/` family was audited:
- `README.md` — safe operation through connectivity loss, bounded queues, recovery, no duplicate execution;
- `offline-mode.md` — last approved policy, command restriction, local detection and deferred synchronization;
- `queueing-and-retry.md` — persistent queues, idempotency keys, backoff and expiry;
- `crash-recovery.md` — state journals, incomplete-command handling, no replay without idempotency and audit continuity;
- `health-monitoring.md` — self-checks, sensor health, heartbeat and degraded capability report;
- `resource-guardrails.md` — CPU/memory/disk/network budgets, priority classes, bounded stop/kill-switch semantics and throttling telemetry;
- `update-and-rollback.md` — update recovery and health gate.

EPT-2 `CAP-EPT-024/025` were re-read to preserve ordering/dedup/gap/loss and backpressure semantics. EPT-4 `CAP-EPT-062` was re-read so timeout/cancel/lost-status semantics remain distinct from recovery. EPT-5 `CAP-EPT-077` was re-read so local update reversion never becomes Govern Response Rollback.

### Resilience ownership result
Endpoint owns local buffering, queue state, replay/resumption facts, restart/crash recovery state, degraded operation, resource pressure and local dependency/capability reassessment. Shared retains generic Jobs/Retry/Recovery/Trace/Activity mechanisms and implementation-contract-level retry primitives.

`replay != exactly once`.
`reconnect != recovered`.
`Agent restarted != healthy`.
`local recovery != Govern rollback`.

## Endpoint Security sources read
The complete active `11-endpoint-agent/security/` family was audited:
- `README.md` — anti-tamper, signed commands/updates, mutual authentication, local audit, secure storage;
- `anti-tamper.md` — protected service/config, maintenance mode, tamper detection, alert/local evidence;
- `command-signing.md` — signer identity, payload hash, target/tenant/expiry, replay protection, revocation;
- `least-privilege.md` — privilege separation, brokered privileged actions, no permanent interactive admin, platform sandboxing;
- `local-audit.md` — append-only local audit semantics, actor/signature context, offline queue, integrity verification requirement;
- `mutual-authentication.md` — Agent/server identity, rotation/revocation/pinning requirements;
- `privacy.md` — minimization, sensitive-field policy and export controls;
- `secret-protection.md` — no secret logging, memory/redaction constraints, raw credential material prohibited;
- `secure-storage.md` — protected-at-rest requirement, key-protection requirement, separation of evidence/queue/config, secure-deletion requirement;
- `self-protection-recovery.md` — recovery mode, signed-repair requirement, tamper evidence preservation and no silent reset.

Global Security sources were revalidated: canonical Permission Model, Privacy and Minimization, Audit and Immutability, and Secrets and Key Management. `17-implementation-contracts/local-audit-sync-contract.md` was read for local sequence, integrity, batch, acknowledgement and deduplication boundaries.

### Security ownership result
Endpoint owns only local technical security state, self-protection/tamper observations, runtime privilege context, protected-component state, local sensitive-material handling facts, local audit events and consumer handoff. Security retains global permission/privacy/audit-integrity policy; Settings retains secret/credential administration; Shared Trace remains a distinct generic mechanism.

`privilege != permission`.
`local audit event != Shared Trace`.
`append-only != cryptographic proof`.
`tamper candidate != malicious verdict`.
`Secret Reference != raw secret`.
`Endpoint security fact != Incident/Finding/Result`.

## Studio / Govern / Shared boundaries revalidated
- Studio `versions-and-deployment/deployment.md` and `rollback.md`: Studio owns lifecycle/reversion for Studio assets only.
- Govern `runs-and-rollback/rollback-model.md`: Govern owns rollback eligibility, Rollback Plan, Response Rollback governance, recovery coordination and Result relationships.
- Shared `background-jobs.md` and implementation `offline-and-retry-contract.md`: generic queue/job/retry mechanisms stay outside Endpoint ownership.

## OPEN decisions
All 18 baseline OPEN decisions remain open unless a separate canonical decision proves otherwise. In particular EPT-6 preserves OPEN-008, OPEN-013, OPEN-015 and OPEN-017; OPEN-004 remains unrelated and intact. No OPEN is closed to obtain documentary PASS.

## Actual source-driven capability set
The source audit independently supports the 18 candidate contracts without forcing a quota:
- `CAP-EPT-082..088` — seven update lifecycle/boundary contracts;
- `CAP-EPT-089..093` — five resilience contracts;
- `CAP-EPT-094..098` — five security/provenance contracts;
- `CAP-EPT-099` — one Endpoint capability-closure contract.

Anti-tamper is independently sourced; runtime privilege/security context is independently sourced; local audit is independently sourced; secret handling is independently sourced; reconnect/replay is independently distinguishable from queue retention; resource pressure is independently distinguishable from dependency recovery. No additional mandatory autonomous family was found beyond this set.

## Forbidden implementation claims
EPT-6 introduces no updater code, package manager, CDN, final package/signature/certificate/crypto format, transport protocol, port, API, physical persistence or queue schema, watchdog implementation, OS service/driver design, final anti-tamper mechanism, final RBAC/ABAC, Fleet scheduler, update orchestrator, storage engine, immutable ledger, cryptographic audit proof, Endpoint Screen ID or Phase 6 capability.

## Audit verdict
**SOURCE AUDIT PASS FOR EPT-6 BUILD**. The 18-capability set is source-justified and may be authored as `draft / defined / planned`. This audit is not the EPT-6 240-gate final verdict and does not itself close Endpoint or Phase 5.