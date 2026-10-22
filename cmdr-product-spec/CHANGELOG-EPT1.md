# CHANGELOG — Endpoint EPT-1

## 2026-08-11 — Enrollment, Inventory, Health and Platform Foundations
- Started from exact baseline `8326a8cf9e9ca3b645395d192c24856058e67034` after Endpoint preflight PASS 100/100.
- Allocated exactly `CAP-EPT-001..014` after revalidating 0 concrete and 0 reserved Endpoint capability IDs.
- Added 14 `draft / defined / planned` provider/platform-neutral capability specifications with **378 sections / 84 mandatory tables** and no Endpoint Screen ID.
- Scope is identity/registration, local enrollment, tenant/environment, platform/OS/architecture, version/build/compatibility, inventory/freshness, health/self-check, heartbeat/connectivity, degraded/offline/stale, technical capability advertisement, Fleet/Policy projection boundaries and provenance/handoff.
- OPEN-008 remains open; Windows/Linux/macOS are referenced candidates only, not delivered/support commitments.
- Settings retains Fleet, enrollment administration, Endpoint Policy, policy assignment, upgrade waves, providers/integrations/credentials/secrets and tenant/environment administration.
- No telemetry/detection/investigation/collection/Live Response/containment/update/resilience/security implementation, API/protocol/PKI/ports/certificates/tokens/code/final RBAC or EPT-2+ capability is introduced.
- Five functional commits were published linearly from the baseline:
  1. `e2805d55f9700296b14e824f28934139099fcc48` — `docs: establish Endpoint identity enrollment and platform boundaries`;
  2. `8e5435499d88836abe4f33b34203a1b93dff704b` — `docs: define Endpoint inventory version and compatibility foundations`;
  3. `5995d3aad377395107abcfae4657f1950d6649c5` — `docs: specify Endpoint health heartbeat and operational state`;
  4. `f5929e17a55cb621b19e51aa328b2f958077579b` — `docs: document Endpoint capability Fleet Policy and provenance boundaries`;
  5. `828b231ec2de4d3b891410a643898577f14cbcc4` — `docs: update Endpoint foundations traceability and quality gates`.
- Baseline → fifth functional/build SHA verified at **5 ahead / 0 behind**, same merge base; build PR/main/README checks passed and CI/status is N/A.
- Post-publication verification companion: `16-quality-and-validation/reports/endpoint-ept1-enrollment-inventory-health-platform-foundations-post-publication-verification.md`.
- A single documentary correction `docs: record Endpoint EPT-1 post-publication verification` records canonical final changelog/status/roadmap/quality evidence and changes no capability contract; exact correction SHA is recorded in PR #2 after remote publication.
- Final verdict becomes **EPT-1 PASS AFTER POST-PUBLICATION VERIFICATION — 190/190 PASS, 0 PENDING, 0 FAIL** after that final remote check. Endpoint Capability Specification remains PARTIAL; EPT-2..EPT-6 remain NOT STARTED.