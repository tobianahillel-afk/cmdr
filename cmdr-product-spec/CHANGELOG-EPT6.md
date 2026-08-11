# CHANGELOG — Endpoint EPT-6

## 2026-08-12 — build-time functional publication set

Execution lot: **EPT-6 — Updates, Resilience, Security and Endpoint Provenance** under Delivery Roadmap Phase 5 — Studio and Endpoint.

Baseline: `4d20e8e411aa8919f31e773a9c5ea96efd64791e` — EPT-5 post-publication verification.

Source-driven capability set: `CAP-EPT-082..099` — **18 capabilities / 486 sections / 108 mandatory tables / at least 54 GWT**.

New coverage:
- administrative-update/Endpoint-local-update boundary;
- package/release/compatibility, download/staging/readiness, install/activation, progress/failure/retry, post-update verification and previous-version recovery;
- offline buffering/queueing, reconnect/replay, persistence/restart/crash recovery, resource/degraded operation and dependency recovery;
- self-protection/anti-tamper, privilege/security context, Secret Reference handling, Local Audit Event/provenance and security-state handoff;
- Endpoint capability-specification closure contract.

Build totals: Endpoint **99 / 2673 / 594**; global **484 capabilities / 482 defined / 2 proposed / 484 planned / 13068 sections / 2904 mandatory tables**. Requirements conservatively remain **122 = 99/20/3/0**; OPEN remains **18**; Endpoint Screen IDs remain **0**.

Build-time quality: **233 PASS / 7 PENDING-REMOTE / 0 FAIL**. EPT-6, Endpoint closure and Phase 5 closure are not final until remote post-publication verification.

No API/protocol/code/package format/crypto scheme/physical schema/final RBAC/Endpoint Screen/Phase 6 capability is introduced.