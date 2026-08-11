# CHANGELOG — Endpoint EPT-1

## 2026-08-11 — Enrollment, Inventory, Health and Platform Foundations
- Started from exact baseline `8326a8cf9e9ca3b645395d192c24856058e67034` after Endpoint preflight PASS 100/100.
- Allocated exactly `CAP-EPT-001..014` after revalidating 0 concrete and 0 reserved Endpoint capability IDs.
- Added 14 `draft / defined / planned` provider/platform-neutral capability specifications with **378 sections / 84 mandatory tables** and no Endpoint Screen ID.
- Scope is identity/registration, local enrollment, tenant/environment, platform/OS/architecture, version/build/compatibility, inventory/freshness, health/self-check, heartbeat/connectivity, degraded/offline/stale, technical capability advertisement, Fleet/Policy projection boundaries and provenance/handoff.
- OPEN-008 remains open; Windows/Linux/macOS are referenced candidates only, not delivered/support commitments.
- Settings retains Fleet, enrollment administration, Endpoint Policy, policy assignment, upgrade waves, providers/integrations/credentials/secrets and tenant/environment administration.
- No telemetry/detection/investigation/collection/Live Response/containment/update/resilience/security implementation, API/protocol/PKI/ports/certificates/tokens/code/final RBAC or EPT-2+ capability is introduced.
- Five functional commits are required; the fifth build remains subject to remote post-publication gates before final 190/190 status.