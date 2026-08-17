---
id: dependency-register-endpoint-ept6
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-12
source-of-truth: registry
---
# Dependency Register — Endpoint EPT-6

| Dependency | Owner | EPT-6 use | Boundary | Status |
|---|---|---|---|---|
| CAP-EPT-001..081 | Endpoint | identity/version/health/telemetry/execution/response foundations | consumed, not rewritten | preserved |
| Settings Endpoint Agent Fleet | Platform Settings | assignment, desired version, waves/channels | admin ownership remains Settings | required |
| Settings Endpoint Policy | Platform Settings | effective constraints | policy administration remains Settings | required |
| Settings Secret References | Platform Settings / Security | opaque sensitive refs | no raw secret ownership | required where used |
| Studio Deployment/Reversion | CMDR Studio | non-overlap reference | Endpoint update != Studio deployment | boundary |
| Govern Response Rollback/Result | Govern | non-overlap / consumer refs | update recovery != response rollback | boundary |
| Shared Jobs/Retry/Recovery | Shared Capabilities | generic mechanism references | local resilience remains Endpoint fact | boundary |
| Security Permission/Privacy/Audit | Security Architecture | global policy constraints | local security facts do not absorb policy | required |
| OPEN-008 | Product Architecture | platform/support uncertainty | remains open | open |
| OPEN-013 | Product Architecture | Class-2 governance | remains open | open |
| OPEN-015 | Product Architecture | Automation/Response bridge | remains open | open |
| OPEN-017 | Product Architecture | detection runtime/language | remains open | open |

## Dependency invariants
- Settings target/version/wave configuration is input, not Endpoint-owned administration.
- Shared retry/recovery contracts are reusable mechanisms, not Endpoint fact ownership.
- Studio deployment and Govern rollback remain separate lifecycle domains.
- Security global policy remains external to Endpoint local observations.
- EPT-6 closes no OPEN decision and introduces no physical implementation dependency.