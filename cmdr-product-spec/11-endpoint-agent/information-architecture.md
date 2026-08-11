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

---

## EPT-3 Detection and Investigation linkage
The historical “Detection/Investigation — future EPT-3” lines above remain the pre-EPT-3 snapshot. EPT-3 now defines `CAP-EPT-031..046` with the conceptual module flow:

**Telemetry → Detection → Investigation → future Collection / Live Response.**

- Telemetry remains EPT-2 source-backed observations and Shared `telemetry-event` projection.
- Detection consumes Investigate-owned Detection Content/version and produces Endpoint-local eligibility/evaluation/match/local-signal-candidate/context/coverage projections. Canonical Command Detection/Signal remain external.
- Investigation provides read-only local process/file/network/user-session/system contexts, contextual timeline/correlation and permission-aware pivots using already available data.
- Detection-to-Investigation expansion stops at `collection-required` whenever new acquisition would be necessary; EPT-4 remains NOT STARTED.
- Endpoint Investigation Summary is a technical handoff only; Investigate retains Case/Evidence/Finding and analyst qualification.
- Shared retains generic Search/Timeline/Linking/Correlation and navigation never grants source permission.

No Endpoint Screen ID, detailed UX, wireframe, button, final filter/column, API/protocol, runtime architecture or platform-support claim is introduced. `OPEN-008` and `OPEN-017` remain open. EPT-4..EPT-6 remain NOT STARTED.

---

## EPT-4 Collection and Live Response linkage
The historical EPT-4 NOT STARTED wording above is preserved as the pre-EPT-4 snapshot. EPT-4 extends the conceptual flow to:

**Telemetry → Detection → Investigation → Collection → Live Response → future Containment.**

- `collection-required` from EPT-3 can become an Investigate-owned Collection Request and then an Endpoint technical intake/operation only after eligibility and authority boundaries.
- Collection provides bounded acquisition, local progress, neutral Collection Items/Package, completeness/integrity metadata and transfer; Investigate performs Artifact/Evidence qualification.
- Live Response provides Endpoint Technical Session, command/script/file execution requests and target-side technical outputs under policy/authority; it is not a Govern Response Run or Studio Automation Run.
- Effectful/destructive process/network/service/file containment actions stop at future EPT-5.
- Shared generic Jobs/Trace/Activity/Export remain mechanisms, not Endpoint object owners.

No Endpoint Screen ID, wireframe, shell/terminal design, button, filter, shortcut, final column, API/protocol/transport, runtime architecture or platform-support claim is introduced. OPEN-008/014/015/017 remain open. EPT-5/EPT-6 remain NOT STARTED.