# E10-INV-110B — owner / permission integration gate (not released)

## Status verified against the code at `af88dcdbff81757314b2a0d0db2dd3b5013d1217`

The implemented module is a bounded read-only Timeline projection and has opt-in dependency-injected interfaces. It is **not** a production authorization service. The preceding security change removed untrusted child-output text from probe errors. Its CI and the one-shot L2 must be evaluated at their exact code SHA; this document is not an L2 verdict.

## Source-backed ownership and permission comparison

The canonical `engineering/implementation/investigation-timeline/contract.json` (A-CONTRACT design snapshot) lists independent obligations:

| Source owner / object | Canonical read obligation | Implementation status |
|---|---|---|
| Investigate Case | `perm.investigate.case.read` | `AccessChecker` dispatches explicit CaseID |
| Shared Timeline Entry | `perm.shared-capabilities.timeline-entry.read` | Per-entry `AccessChecker`, selected by exact ID |
| Investigate Evidence | `perm.investigate.evidence.read` | **No concrete Evidence provider or typed policy resolver** |
| Shared telemetry event | `perm.shared-capabilities.telemetry-event.read` | **No concrete telemetry provider or typed policy resolver** |
| Govern Decision | `perm.govern.decision.read` | **No concrete Govern provider or typed policy resolver** |
| Command Incident | `perm.command.incident.read` | **No concrete Command provider or typed policy resolver** |

Both `ProjectWithAccessChecker` and `ProjectFromOwnerSnapshots` currently send linked-source checks using generic `Resource="source"` and `Permission="source-owner-read"`. This is an internal obligation placeholder, **not** one of the four canonical owner-specific permissions listed above. An actual integration must not interpret the generic placeholder as blanket access. It must resolve each source's *authoritative owner/type and Case relationship* independently and enforce the matching real permission; source type must **not** be trusted merely because a requester put a prefix into `SourceRef`.

## Minimum readiness contract for a future real adapter

1. Authenticate principal and bind tenant, environment and CaseID on every provider call. Reject unknown or mismatched identity rather than deriving permissions from caller-supplied Scope or Entry flags.
2. Resolve canonical Shared Timeline Entry ownership under the same Case and verify `perm.shared-capabilities.timeline-entry.read`; any denied Entry must never drive a source lookup or expose ID/provenance.
3. Resolve the linked source through the owning service's trusted object-reference mechanism (typed object + id + version/tenant as defined by that owner), enforce its own real canonical read permission, check Case relationship, and fail closed for unknown/ambiguous kind.
4. Account for ACL revocation and snapshots becoming stale between checks. The current three-phase sequence is bounded but **not an atomic security transaction**; a provider consistency/version strategy and adversarial race tests are prerequisites for production assurance.
5. Add independent integration tests with actual first-party provider adapters (no mock-only promotion), denied/missing/foreign/misbound cases for all four source types, cancellation, provider failure and rollback/non-disclosure. Record exact provider version, CI and source SHA.
6. Review the owner-specific design and permission routing before replacing the generic placeholder. No new schema, storage backend, vendor index or production SLO is selected by this note.

## Decisions and next authorized actions

**Reuse before inventing:** prefer the existing canonical object-reference and permission-register contracts, Go's current `AccessChecker`/`EntryResolver` boundaries, and the owning service's future real provider. Parsing arbitrary strings or building a new global ACL database would duplicate ownership and risk an authorization bypass. There is no source-backed basis to choose such a vendor/framework at this stage.

**Not verified:** real server ACL provider, principal authentication, canonical source-type-specific permission dispatch, atomic ACL consistency, navigation, UI, Shared Timeline Engine, production indexing/performance. `OPEN-013` remains unresolved. The work graph must retain `E10-INV-110B-RUNTIME = READY`, `E10-INV-110C-NAV = BLOCKED` until genuinely qualifying implementation and independent review evidence exist.

This is a **documentation-only, independently scoped follow-up**. It does not rewrite the earlier immutable code snapshot scheduled for L2, and does not authorize promotion.
