---
id: endpoint-information-architecture
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# Endpoint Agent Information Architecture

This supporting document resolves the pre-existing README navigation gap using only canonical product boundaries and current module structure. It creates no Screen ID, wireframe, final layout, button, filter, shortcut, implementation architecture, protocol or platform-support claim.

## Module hierarchy
1. Foundations — identity/registration, enrollment local state, tenant/environment, platform/version, inventory/freshness, health/connectivity/state, technical capability projection and provenance.
2. Telemetry — future EPT-2.
3. Detection — future EPT-3.
4. Investigation — future EPT-3.
5. Collection — future EPT-4.
6. Live Response — future EPT-4.
7. Containment — future EPT-5.
8. Resilience — future EPT-6.
9. Security — cross-cutting constraints with deep EPT-6 work.

## Source ownership
Endpoint owns individual local technical facts. Platform Settings owns Fleet, enrollment administration, Endpoint Policies/assignment, tenant/environment administration and upgrade waves. Govern owns response authority/run/result/rollback. Studio owns Tool/Automation Agent/Automation Run. Shared owns generic mechanisms.

## Navigation relationships
Product consumers can navigate by stable Agent/Endpoint/Fleet/Policy/tenant/environment references where permission allows. Return origin and source ownership are preserved; navigation never grants source permission.

## EPT-1 linkage
`CAP-EPT-001..014` form the Foundation capability layer. No Endpoint product-specific screen exists at EPT-1; Settings retains `SET-EAF-001` and `SET-EPL-001` for administrative Fleet/Policy surfaces.

## Open support boundary
OPEN-008 remains open. Windows, Linux and macOS remain referenced candidates with actual release support/version NOT DECIDED.

---

## EPT-2 Telemetry linkage
The historical “Telemetry — future EPT-2” line above is the pre-EPT-2 snapshot. EPT-2 now defines the functional Telemetry capability layer `CAP-EPT-015..030` without creating a product screen or technical architecture.

Conceptual navigation/relationships only:
- source/sensor → local observation → Shared `telemetry-event` projection;
- observation families → process/file/network/auth/system/sensor facts;
- normalization → ordering/freshness/gap/loss/rate/sampling/backpressure/privacy states;
- detailed capability declaration/availability → CAP-EPT-011 foundation summary and Settings capability-inventory projection;
- observations → authorized Investigate/Command/Govern consumers without Evidence/Finding/Result promotion;
- capability declarations → Studio consumers without becoming Tools.

No Screen ID, wireframe, final layout, API, protocol, event schema, storage/event-bus/SIEM architecture or platform-support claim is introduced. EPT-3..EPT-6 remain NOT STARTED.
