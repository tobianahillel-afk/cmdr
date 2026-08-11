---
id: endpoint-resilience-ept6-capability-links
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-12
source-of-truth: canonical
---
# Resilience — EPT-6 Capability Links

The historical Resilience documents remain source/background material. Normative capability contracts are:

- `CAP-EPT-089` — Offline Buffering, Queueing and Local Retention Semantics;
- `CAP-EPT-090` — Reconnect, Replay, Resumption and Duplicate-Control Boundary;
- `CAP-EPT-091` — Local State Persistence, Restart and Crash Recovery;
- `CAP-EPT-092` — Resource Pressure, Backpressure and Degraded-Operation Semantics;
- `CAP-EPT-093` — Dependency Failure, Capability Reassessment and Recovery.

Cross-links: EPT-2 `CAP-EPT-024/025` for ordering/loss/backpressure, EPT-4 `CAP-EPT-062` for timeout/cancel/lost status, EPT-6 Update `CAP-EPT-087/088`, EPT-6 Security `CAP-EPT-094..098`.

Shared retains generic Jobs/Retry/Recovery/Trace/Activity. No storage engine, physical queue, watchdog, OS service or exactly-once guarantee is introduced.