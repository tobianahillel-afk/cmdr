---
id: endpoint-ept6-update-capability-links
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-12
source-of-truth: canonical
---
# Endpoint Update — EPT-6 Capability Links

Normative Agent-update capability contracts:
- `CAP-EPT-082` — Settings assignment/channel/wave boundary;
- `CAP-EPT-083` — package/release/compatibility context;
- `CAP-EPT-084` — download/staging/readiness;
- `CAP-EPT-085` — installation/activation/local lifecycle;
- `CAP-EPT-086` — progress/deferral/failure/retry;
- `CAP-EPT-087` — post-update health/compatibility/technical verification;
- `CAP-EPT-088` — update reversion/previous-version recovery.

Sources include `resilience/update-and-rollback.md`, `detection/detection-update.md`, EPT-1 version/health contracts and `17-implementation-contracts/endpoint-agent-update-contract.md`.

Platform Settings retains Fleet targeting, desired version, update channels/waves and policy assignment. Studio retains Studio asset deployment/reversion. Govern retains response rollback. No package format, signature scheme, update transport, updater implementation or support-platform decision is introduced.