---
id: endpoint-ept5-migration-and-source-map
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# Endpoint EPT-5 — Migration and Source Map

Additive migration only; historical files remain preserved.

| Historical source | Historical owner | EPT-5 target | Preserved semantics | Boundary / stale terminology |
|---|---|---|---|---|
| containment/README.md | Endpoint | CAP-EPT-065..080 | Govern authority, scope, rollback, verification | no execution authority |
| containment/account-containment.md | Endpoint | CAP-EPT-081 | local session lock/termination, scope, rollback | directory actions external |
| containment/host-containment.md | Endpoint | CAP-EPT-068/073..080 | composite host impact/dependencies/verification | business impact remains Govern context |
| containment/network-isolation.md | Endpoint | CAP-EPT-068 | isolation, control-path exception concept, verification | no firewall implementation |
| containment/quarantine.md | Endpoint | CAP-EPT-070/071/078 | path/hash, restricted state, restore | secure-store physical design deferred |
| containment/rollback.md | Endpoint | CAP-EPT-077 | reversal support/partial outcome | historical `Result` term maps to technical outcome; canonical Result Govern |
| containment/verification.md | Endpoint | CAP-EPT-074/075 | independent checks/telemetry confirmation | historical `Result object` term not owned by Endpoint |
| live-response/process-actions.md | Endpoint | CAP-EPT-067 | terminate/suspend/resume/tree scope/verification | moved from EPT-4 effect boundary |
| live-response/network-actions.md | Endpoint | CAP-EPT-069 | temporary block/scope/expiry/conflict/verification | distinct from host isolation |
| live-response/file-actions.md | Endpoint | CAP-EPT-070/071 | delete/restore/quarantine integration/hash/path | move/copy general semantics remain EPT-4 file operation boundary unless effect governed |
| live-response/service-actions.md | Endpoint | CAP-EPT-072 | start/stop/restart/disable-like semantics | only when platform capability declares support |

No source is deleted, no owner is silently transferred, and no EPT-6 update/resilience/security implementation is pulled forward.