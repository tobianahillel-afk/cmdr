---
id: platform-scale-identity-administration-principals-roles-access-reviews-capability-conformance
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-12
source-of-truth: quality-report
---
# Platform Scale — Identity Administration — Principals, Roles and Access Reviews — Capability Conformance

## Scope

Functional execution lot under **Delivery Roadmap Phase 6 — Platform Scale**. This is not Phase 6B. The lot allocates exactly `CAP-SET-005..007`; `CAP-SET-008+` is neither allocated nor reserved.

## BUILD-TIME EVIDENCE

### Execution baseline

- required and audited remote baseline: `d605265f5b8a2e4350388b4ec9cfe51920a4aa50`;
- baseline title: `docs: record Settings tenant environment foundations post-publication verification`;
- previous Settings lot: **PASS AFTER POST-PUBLICATION VERIFICATION — 160/160**;
- PR #2: open / Draft / unmerged / base `main` / auto-merge disabled at preflight;
- main: `bc1ec59e5e79ccbaba291e2984b0eb7d5e54129c` at preflight;
- branch and main README: exactly `# cmdr`, blob `901c74cda52e28b5ff7fc425ddf28ef89f3ad875` at preflight.

### Source corpus audit

Execution reread the governance source-of-truth/ownership/registers/Requirements/OPEN/product-boundary corpus, Platform Settings product/IA/navigation/permission sources, Users & Roles sources, Administrative Audit sources, Principal/Role/Tenant/Environment objects, Govern object boundaries, mandatory Security sources and established action classification before writes.

Findings:

- canonical Identity module is `users-and-roles`;
- canonical objects are Principal and Role; no canonical Group/Group Membership was found;
- no generic AccessAssignment/RoleAssignment/PermissionAssignment/EffectiveAccess object or engine is sourced;
- Principal states are exactly `pending`, `active`, `suspended`, `revoked`;
- Role states are exactly `draft`, `active`, `deprecated`;
- Role expiry is a condition/constraint, not a state;
- Security owns Permission Model/RBAC/ABAC/authorization/tenant isolation/SoD/step-up;
- Govern owns Approval/Decision/Decision Authority/Response Run/Result;
- `SET-IAM-001` and `SET-AUD-001` are reused; no new screen is needed.

### Namespace audit

The prior Settings shard contains only `CAP-SET-001..004` and explicitly records no reservation of `CAP-SET-005+`. Target canonical paths for `CAP-SET-005`, `006` and `007` were absent at the execution baseline. Build allocation is therefore exactly:

- `CAP-SET-005`;
- `CAP-SET-006`;
- `CAP-SET-007`.

No `CAP-SET-008+` reservation or allocation is made.

### Capability matrix

| Capability | Responsibility | Primary objects | Sections | Mandatory tables | GWT | Delivery |
|---|---|---|---:|---:|---:|---|
| `CAP-SET-005` | Principal administrative lifecycle and identity state | Principal; Tenant ref | 27 | 6 | >=3 | draft / defined / planned |
| `CAP-SET-006` | Role lifecycle/constraints and bounded Principal relation | Role; Principal ref; Tenant ref | 27 | 6 | >=3 | draft / defined / planned |
| `CAP-SET-007` | Access Review evidence and revocation disposition | Principal/Role refs | 27 | 6 | >=3 | draft / defined / planned |
| **Total** | | | **81** | **18** | **>=9** | |

The six mandatory non-empty capability-specific tables are S8 Inputs, S9 Objects Read, S10 Objects Created/Modified, S13 Automation/AI, S16 Outputs and S17 Transitions/Handoffs.

### CAP-SET-007 revocation-mechanics verdict

**Variant (b) is used.** The canonical Principal/Role contracts define typed tenant-scoped relations and auditable mutations, but the execution corpus does not define a concrete generic role-assignment removal mechanic. Therefore a review result `revoke` produces a **revocation disposition / revocation-required handoff**. The lot does not invent `AccessAssignment`, `RoleAssignment`, `PermissionAssignment`, direct Permission revocation or implicit relation deletion.

### Ownership and safety

- Platform Settings owns the three capability semantics and Principal/Role administrative lifecycle.
- Security retains permission semantics and authorization evaluation.
- Govern retains authority-bearing response/approval/decision semantics.
- Principal ≠ Person/Customer/Tenant.
- Human Principal ≠ Service Principal.
- authentication mapping ≠ authorization.
- Role ≠ Permission/Group/job title/Decision Authority.
- Principal/Role relation ≠ permission grant/effective authorization.
- Access Review ≠ Govern Approval/Decision/generic Compliance engine.
- Review evidence ≠ Investigate Evidence automatically.

### Groups / Access Assignment / Effective Access audit

Created or defined by this build:

- Group: **0**;
- Group Membership: **0**;
- AccessAssignment: **0**;
- RoleAssignment: **0**;
- PermissionAssignment: **0**;
- EffectiveAccess object/engine: **0**.

### Screens

- active Screen Register baseline: **56**;
- Settings screens: **12**;
- reused: `SET-IAM-001`, `SET-AUD-001`;
- new Screen IDs: **0**;
- Screen Register functional change expected: **0**.

### Permissions

Consumed existing families only:

- `perm.settings.identity.read/manage`;
- `perm.platform-settings.principal.read/manage`;
- `perm.platform-settings.role.read/manage`.

The alias/normalization debt is documented but not normalized in bulk. New Permission IDs: **0**.

### AI contract

AI may explain sourced state/constraints, summarize review provenance, explain deterministic denial, highlight a potential source-backed SoD issue for human review and suggest a safe next step. AI cannot autonomously create/suspend/revoke Principal, create/mutate Role or Principal/Role relation, grant/revoke Permission, execute a review disposition, invent Group/effective access, bypass Security, create Decision Authority, resolve OPEN decisions or fabricate provenance. Every path has a deterministic/manual alternative.

### OPEN decisions

All **18** existing OPEN decisions remain open. `OPEN-013` remains directly relevant to default reversible Class-2 governance. No OPEN is created or closed by this lot.

### Requirements

Requirement IDs remain **122**. Build evidence was audited particularly against `REQ-PROD-006`, `REQ-PROD-017`, `REQ-OBJ-001`, `REQ-SEC-001`, `REQ-SEC-002`, `REQ-SEC-006` and relevant privacy/UX constraints. No Requirement state is promoted merely to make this lot pass.

Build distribution remains:

- conform: **99**;
- partial: **20**;
- absent: **3**;
- contradictory: **0**.

### Build totals

| Metric | Baseline | Build | Delta |
|---|---:|---:|---:|
| capabilities | 488 | **491** | +3 |
| defined | 486 | **489** | +3 |
| proposed | 2 | **2** | 0 |
| planned | 488 | **491** | +3 |
| numbered sections | 13,176 | **13,257** | +81 |
| mandatory tables | 2,928 | **2,946** | +18 |

Settings cumulative build state: **7 capabilities / 189 sections / 42 mandatory tables**.

Closed domains remain unchanged: Command 27 PASS; Investigate 243 PASS; Govern 47 PASS; Studio 68 PASS; Endpoint 99 PASS. Delivery Roadmap Phase 5 remains PASS — capability specification complete.

### Build-time gate matrix

| Gates | Count | Family | Build state |
|---|---:|---|---|
| 1–12 | 12 | Git/baseline/PR/main/README/previous closure | PASS |
| 13–30 | 18 | source corpus/source-of-truth | PASS |
| 31–42 | 12 | namespace/free IDs/reservations | PASS |
| 43–60 | 18 | three capability structural contracts | PASS |
| 61–76 | 16 | ownership/object boundaries | PASS |
| 77–92 | 16 | Principal lifecycle/identity semantics | PASS |
| 93–106 | 14 | Role lifecycle/relation/access boundary | PASS |
| 107–118 | 12 | Access Review/evidence/revocation | PASS |
| 119–134 | 16 | Security/permissions/SoD/step-up/AI | PASS |
| 135–144 | 10 | Screens/Settings IA/zero-ID rules | PASS |
| 145–156 | 12 | Requirements/OPEN/migration/history | PASS |
| 157–164 | 8 | registers/traceability/roadmap | PASS |
| 165–168 | 4 | candidate diff/count/non-regression | PASS after pre-publication candidate audit |
| 169–174 | 6 | remote publication/final verification | **PENDING-REMOTE** |
| **Total** | **174** | | **168 PASS / 6 PENDING-REMOTE / 0 FAIL** |

The fifth functional commit SHA is intentionally not predicted inside its own content.

## POST-PUBLICATION EVIDENCE

**PENDING.** This section may be finalized only after the five functional commits are published by non-forced fast-forward and the published remote content, ancestry, PR/main/README and CI/status are independently re-read.

Required final checks:

1. remote BUILD SHA equals the fifth functional commit;
2. baseline → BUILD = 5 ahead / 0 behind with same merge-base;
3. published CAP-SET-005..007 = 3 / 81 / 18 / >=9 GWT;
4. global/Settings counts and zero-ID/object invariants are remotely confirmed;
5. PR #2/main/README state is revalidated;
6. applicable status/workflow state is checked.

Only then may the documentary post-publication record state **PASS AFTER POST-PUBLICATION VERIFICATION — 174/174 PASS, 0 PENDING, 0 FAIL**.

## Build disposition

Identity Administration functional build: **168 PASS / 6 PENDING-REMOTE / 0 FAIL**.

Settings Capability Specification: **PARTIAL**.  
Delivery Roadmap Phase 6 Capability Specification: **PARTIAL**.  
Global Capability Specification: **PARTIAL**.  
Repository maturity: **PARTIAL**.

No product implementation, API, protocol, physical schema, final RBAC/ABAC, authentication protocol, SSO implementation, IdP selection, SCIM implementation or platform-support claim is introduced.
