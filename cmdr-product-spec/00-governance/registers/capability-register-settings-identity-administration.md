---
id: capability-register-settings-identity-administration
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-12
source-of-truth: registry
---
# Capability Register — Settings Identity Administration

This shard allocates exactly `CAP-SET-005..007` for **Identity Administration — Principals, Roles and Access Reviews** under Delivery Roadmap Phase 6 — Platform Scale. IDs are immutable and never recycled. `CAP-SET-008+` is neither allocated nor reserved by this lot.

| Capability ID | Name | Owner product/module | Status | delivery_status | delivery_mode | Canonical file | Primary roles | Primary objects | Consumers | Requirement IDs | OPEN | Dependencies |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| `CAP-SET-005` | Principal Administrative Lifecycle and Identity State | Platform Settings / users-and-roles | draft | defined | planned | `10-platform-settings/capabilities/cap-set-005-principal-administrative-lifecycle-and-identity-state.md` | Platform Administrator; Security Administrator | Principal; Tenant ref | Settings; Security; Command; Investigate; Govern; Studio; Endpoint; Shared | REQ-PROD-006; REQ-OBJ-001; REQ-SEC-001; REQ-SEC-002 | OPEN-013 | Principal; Tenant; Users & Roles; Permission Model; ABAC; tenant isolation; SoD; step-up; Administrative Audit |
| `CAP-SET-006` | Role Administrative Lifecycle, Constraints and Principal-Relation Boundary | Platform Settings / users-and-roles | draft | defined | planned | `10-platform-settings/capabilities/cap-set-006-role-administrative-lifecycle-constraints-and-principal-relation-boundary.md` | Platform Administrator; Security Administrator | Role; Principal ref; Tenant ref | Settings; Security; Govern; other authorized consumers | REQ-PROD-006; REQ-OBJ-001; REQ-SEC-001; REQ-SEC-002; REQ-SEC-006 | OPEN-013 | Role; Principal; Tenant; Permission Model/catalog; ABAC; SoD; step-up; Decision Authority; Administrative Audit |
| `CAP-SET-007` | Access Review, Evidence and Revocation Disposition | Platform Settings / users-and-roles | draft | defined | planned | `10-platform-settings/capabilities/cap-set-007-access-review-evidence-and-revocation-disposition.md` | Review owner; Platform Administrator; Security Administrator | Principal ref; Role ref | Settings; Security; Govern when separately required; other authorized consumers | REQ-PROD-006; REQ-OBJ-001; REQ-SEC-001; REQ-SEC-002; REQ-SEC-006 | OPEN-007; OPEN-013; OPEN-014; OPEN-015; OPEN-019 | Access Reviews; Principal; Role; Administrative Audit; Permission Model; ABAC; SoD; step-up |

## Structural allocation

- capabilities: **3**;
- defined: **3**;
- proposed: **0**;
- planned: **3**;
- numbered sections: **81**;
- mandatory tables: **18**;
- minimum meaningful Given/When/Then scenarios: **9**;
- new canonical object contracts: **0**;
- new Permission IDs: **0**;
- new Screen IDs: **0**.

## Hard boundaries

- Principal and Role remain Platform Settings-owned canonical objects.
- Security retains Permission Model, RBAC/ABAC, authorization evaluation, tenant/environment isolation, SoD, step-up and permission enforcement.
- Govern retains Action Request, Approval, Decision, Decision Authority lifecycle, Response Run and Result.
- Groups, Group Membership, generic Access Assignment, RoleAssignment, PermissionAssignment and Effective Access are outside this shard.
- Role expiry is a condition/constraint, not a Role state.
- `CAP-SET-007` uses **revocation disposition/handoff** because the current canonical corpus does not define a concrete generic Principal/Role assignment-removal mechanic.

## Reservation state

No `CAP-SET-008+` ID or range is reserved.
