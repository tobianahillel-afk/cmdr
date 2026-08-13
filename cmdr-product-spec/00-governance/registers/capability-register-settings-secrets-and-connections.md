---
id: capability-register-settings-secrets-and-connections
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-13
source-of-truth: registry
requirements:
  - REQ-PROD-006
  - REQ-PROD-009
  - REQ-PROD-010
  - REQ-PROD-011
  - REQ-PROD-012
---
# Capability Register — Settings Secrets & Connections

## Lot identity
Execution lot: **Secrets & Connections — Integration and Secret Reference Administration** under **Delivery Roadmap Phase 6 — Platform Scale**. This is an execution lot, not a new roadmap phase.

## Namespace allocation record
This first functional commit allocates exactly `CAP-SET-008` and `CAP-SET-009` after execution-time repository, register and tree revalidation. Their capability contracts are authored by the following functional commits in the same four-commit build. This allocation record does not claim runtime availability or implementation.

| Capability ID | Intended canonical title | Owner | Canonical object | Contract path | Allocation state |
|---|---|---|---|---|---|
| `CAP-SET-008` | Integration Administrative Lifecycle, Validation and Connection State | Platform Settings Product Lead | Integration | `10-platform-settings/capabilities/cap-set-008-integration-administrative-lifecycle-validation-and-connection-state.md` | allocated; contract follows in this build |
| `CAP-SET-009` | Secret Reference Administrative Lifecycle, Rotation and Revocation | Platform Settings Product Lead | Secret Reference | `10-platform-settings/capabilities/cap-set-009-secret-reference-administrative-lifecycle-rotation-and-revocation.md` | allocated; contract follows in this build |

No `CAP-SET-010+` ID is allocated or reserved by this lot. IDs remain immutable and must never be recycled.

## Canonical ownership boundaries
- Platform Settings owns the administrative lifecycle and Settings-side configuration semantics of canonical `Integration` and `Secret Reference` objects.
- Security retains Permission Model, permission semantics, RBAC/ABAC, tenant isolation, SoD, step-up, authorization enforcement, secret-protection policy and cryptographic/key-management policy.
- Studio retains Tool, Tool Call, Skill, Workflow, Human Gate, Automation Run and agentic/runtime execution.
- Endpoint retains technical endpoint execution; Command retains Command objects; Investigate retains investigation/analysis runtime; Govern retains Action Request, Approval, Decision, Decision Authority, Response Run and Result.
- Consumer projections never transfer canonical ownership.

## Terminology boundaries
- `Integration` is the canonical Platform Settings object. Its `capabilities` metadata is integration metadata and is **not** the CMDR Capability object or a `CAP-*` allocation mechanism.
- `Connection` is module/functional terminology around Integration administration; no Connection object is created.
- `Connector` remains consumer/runtime terminology unless a later canonical source establishes an object; no Connector object is created here.
- `Credential` is not a standalone canonical object in this lot; no Credential object is created.
- `Secret Reference` is the canonical administrative reference object. Raw secret material is protected material and is never exposed as a Settings object.
- `Model Provider`, `Data Source` and `Parser` remain separate canonical objects in their own Settings modules and are outside this lot.

## Hard zero-change gates
This lot targets **0 new Screen IDs, 0 new Permission IDs and 0 new canonical objects**. It reuses `SET-SEC-001`, existing Integration/Secret Reference permissions and existing Security controls. If those existing contracts prove insufficient, execution must block rather than invent an ID or object.

## Environment boundary
Tenant scope is mandatory for both target objects. Environment is consumed only where an exact canonical source establishes the relation; this lot does not create a global Environment requirement or reopen `CAP-SET-001..004`.

## Deferred adjacent lots
Models & Providers and Sources & Parsers remain outside this lot. Customer/MSSP/Delivery and implementation remain outside this run. No `Phase 6B` is created.

## Build closure rule
The final functional commit must replace the transitional allocation wording above with the final active shard state only after both capability contracts exist and structural, ownership, security, traceability and quality gates have been audited.