---
id: endpoint-ept2-source-audit
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-11
source-of-truth: quality-report
---
# Endpoint EPT-2 — Source Audit

Baseline: `1f8e482f6b7949885bd1bd7ae691215bde187b28`. EPT-1 is PASS 190/190 and `CAP-EPT-001..014` remain the foundation.

## Telemetry family — 11/11 read
1. `telemetry/README.md` — continuous/contextual telemetry, event time/sequence/source, buffering/privacy/rate controls.
2. `authentication-telemetry.md` — logon/session/user/device facts, no raw secrets, platform differences.
3. `cloud-workload-telemetry.md` — optional workload/container identity/runtime facts; capability conditional.
4. `dns-telemetry.md` — queries/responses, process attribution when available, cache/encrypted-DNS limitations.
5. `file-telemetry.md` — create/modify/delete/rename, path/process, policy-bounded hash/minimization.
6. `network-telemetry.md` — connections/listeners/flows, process/DNS relation and rate controls.
7. `process-telemetry.md` — start/stop, parent-child, executable/hash/signature/user; command metadata restricted.
8. `registry-telemetry.md` — registry activity explicitly platform-dependent/Windows-only as a source declaration, not a global support claim.
9. `script-and-shell-telemetry.md` — script/shell execution facts with policy/redaction and optional source integration.
10. `system-telemetry.md` — boot/services/drivers/tasks/kernel/system facts and capability context.
11. `telemetry-buffering.md` — queue/order/backpressure/retention/encryption source material; no buffer implementation selected.

## Critical boundaries
- Shared owns canonical `telemetry-event` and generic telemetry normalization.
- Endpoint owns endpoint-local source/observation semantics, source health/coverage and detailed capability declarations.
- Detection = future EPT-3 interpretation; Investigation = EPT-3 host inspection; Collection = EPT-4 governed acquisition.
- Settings retains source/Fleet/Policy/provider/secret/tenant administration.
- Investigate retains Evidence/Finding/Case; Studio Tool/Tool Call/Automation Run; Govern Decision/Response Run/Result.
- local audit remains Endpoint-owned and distinct from Shared Trace.
- `OPEN-008` remains open; Windows/Linux/macOS/cloud/container/mobile support is not declared delivered.

## Capability-set decision
`CAP-EPT-015..030` are all justified. Authentication and system/source-specific observations have enough independent source semantics to remain autonomous; DNS, registry, script/shell and cloud-workload semantics are represented inside the corresponding network/system/process/declaration contracts rather than creating empty extra capabilities. `CAP-EPT-027/028` deepen but do not duplicate CAP-EPT-011.
