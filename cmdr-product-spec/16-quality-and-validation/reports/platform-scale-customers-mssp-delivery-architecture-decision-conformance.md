---
id: platform-scale-customers-mssp-delivery-architecture-decision-conformance
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-14
source-of-truth: quality-report
---
# Platform Scale — Customers / MSSP / Delivery Architecture Decision — Conformance

## Scope

Architecture-recording/unblocking run under Delivery Roadmap Phase 6 — Platform Scale. This run records explicit human approval for `OPEN-006`, aligns Security and Command, promotes `CAP-CMD-401` from `proposed` to documentary `defined`, and creates no runtime implementation.

Approval reference: **Hillel Tobiana — explicit project-owner approval in ChatGPT conversation**.

Canonical decision: `../../00-governance/adr/ADR-0008-customers-mssp-delivery-deployment-and-cross-tenant-architecture.md`.

The complete approved D1–D7 authority input is preserved **verbatim** in ADR-0008 and in the resolved OPEN-006 record. Equality to that verbatim block is the approval-integrity source; no inferred replacement decision is used.

## Execution baseline

- branch: `docs/cmdr-product-spec-foundation` ;
- baseline: `6206fa322895bface4c11d173afc5c30ceaa472c` — `docs: record Settings Sources and Parsers post-publication verification` ;
- baseline tree: `7cfd4d4e9b283af76cc51d51da0a18d92f9bb57d` ;
- main: `bc1ec59e5e79ccbaba291e2984b0eb7d5e54129c` ;
- PR #2: open / Draft / unmerged / base `main` ;
- baseline README branch/main: exact `# cmdr`, blob `901c74cda52e28b5ff7fc425ddf28ef89f3ad875` ;
- Sources & Parsers predecessor: **316/316 PASS** ;
- `CAP-SET-014+`: 0 allocated / 0 reserved.

## Functional chain before BUILD publication

1. `0cfef4c56825c76865a4aaefa3182c9d382186d5` — `docs: record Customers and Delivery deployment architecture decision` ;
2. `e533346197bf70e90af49c4232cc611fc9a8d83f` — `docs: define MSSP authorized tenant-set and cross-tenant safety boundaries` ;
3. `1eb1bb15de57cc878ec541bb2bb4ef9c724f3fe9` — `docs: align Command Customers and Delivery with approved deployment model` ;
4. this report/traceability commit is the functional/documentary BUILD. Exact BUILD SHA is recorded after commit construction and remote publication verification.

All commits were constructed off-ref. The branch is moved only after the final pre-publication concurrency guard and only by non-forced fast-forward.

## Approved decision conformance

| Decision | Approved architecture | Build evidence | Result |
|---|---|---|---|
| D1 | Internal + Enterprise multi-tenant + MSP/MSSP; MSSP deployment-dependent | ADR-0008 + Customers & Delivery | PASS |
| D2 | Customer external projection; not canonical; not Tenant alias | ADR-0008 + CAP-CMD-401 | PASS |
| D3 | independent Tenants; no hierarchy/ManagedTenant/TenantGroup/Portfolio/CustomerTenant | ADR-0008 + Tenant Isolation | PASS |
| D4 | non-canonical Authorized Tenant Set; read-only aggregation/context switch only | Tenant Isolation + Permission Model | PASS |
| D5 | Search/Report/Export single-selected-Tenant; Shared ownership | ADR-0008 + CAP-CMD-401 + screen | PASS |
| D6 | Tenant selection → Security re-evaluation → Govern Decision Authority | Security docs + CAP-CMD-401 | PASS |
| D7 | CAP-CMD-401 same ID/owner, `defined`, narrowed scope | CAP-CMD-401 + registers | PASS |

`OPEN-013` and `OPEN-019` remain open. Only `OPEN-006` is resolved.

## Object / permission / screen invariants

- new Capability IDs: **0** ;
- `CAP-SET-014`: not allocated / not reserved ;
- `CAP-CMD-402`: not allocated / not reserved ;
- new canonical objects: **0** ;
- new Permission IDs: **0** ;
- new Screen IDs: **0** ;
- active Screens remain **56** ;
- `CMD-CRP-001` is reused and mapped to `customers-and-delivery` ;
- Customer remains external ;
- Authorized Tenant Set remains non-canonical ;
- no Group/Assignment/Effective Access regression ;
- no `cross-tenant.manage`.

## Security invariants

- Tenant-first object resolution ;
- requested Tenant set must be a subset of the Authorized Tenant Set ;
- object permission remains server-side RBAC/ABAC ;
- each aggregate object retains Tenant identity ;
- no fallback Tenant ;
- read does not imply export/admin/execute/approve/response ;
- cross-tenant mutation/admin/response/delegated admin/export widening are forbidden ;
- response is tenant-local after Security and Govern re-evaluation ;
- deny/cross-scope attempts remain auditable.

## Shared boundaries

- Global Search remains Shared-owned and single-selected-Tenant initially ;
- Reporting Engine remains Shared-owned and Report remains single-Tenant initially ;
- Export remains Shared-owned, single-Tenant initially and cannot widen visibility ;
- multi-tenant Search/Report/Export are deferred ;
- external/client sharing remains fenced by `OPEN-019`.

## Command boundaries

`CAP-CMD-401` is `draft / defined / planned`, deployment-dependent, same ID and same Command Product Lead owner.

Included: deployment-aware activation, external Customer/engagement projection, authorized read-only multi-tenant overview, Tenant selector/context switching, portfolio-like View without Portfolio object, single-Tenant Reporting requests, service-delivery Tasks, contractual SLA projection, source/freshness/audience visibility.

Excluded: Customer lifecycle/admin, Tenant hierarchy, cross-tenant mutation/admin/response, delegated administration, initial multi-tenant Search/Report/Export, CRM, billing, customer portal, contract mutation.

`CAP-CMD-105` no longer depends on OPEN-006; contractual SLA remains an external/deployment projection and `OPEN-013` remains applicable to Class-2 policy.

## Counters

- global: **497 capabilities / 496 defined / 1 proposed / 497 planned / 13,419 sections / 2,982 mandatory tables** ;
- Command: **27 / 27 defined / 0 proposed / 27 planned / 729 / 162** ;
- Settings: **13 / 351 / 78** ;
- Requirements: **122 = 99 conform / 20 partial / 3 absent / 0 contradictory** ;
- OPEN: **17** ;
- active Screens: **56**.

No Requirement ID or Requirement state is changed.

## Expected functional surfaces

The functional BUILD is expected to differ from baseline only on these architecture/traceability surfaces:

1. `00-governance/adr/ADR-0008-customers-mssp-delivery-deployment-and-cross-tenant-architecture.md` — new ;
2. `00-governance/decision-log.md` ;
3. `00-governance/adr/ADR-0001-product-separation.md` ;
4. `00-governance/source-material/unresolved-decisions.md` ;
5. `14-security-permissions-and-trust/tenant-isolation.md` ;
6. `14-security-permissions-and-trust/permission-model.md` ;
7. `06-command/modules/customers-and-delivery/README.md` ;
8. `06-command/modules/customers-and-delivery/capabilities/customers-and-delivery-context.md` ;
9. `06-command/modules/customer-and-reports/README.md` ;
10. `06-command/modules/customer-and-reports/screens/customer-overview.md` ;
11. `06-command/modules/incidents-and-work-queue/capabilities/sla-tracking.md` ;
12. `00-governance/registers/capability-register-command.md` ;
13. `00-governance/registers/capability-register.md` ;
14. `00-governance/registers/screen-register.md` ;
15. `00-governance/source-material/requirements-traceability-matrix.md` ;
16. `18-roadmap-and-releases/phase-6-platform-scale.md` ;
17. this conformance report — new ;
18. `16-quality-and-validation/quality-index-platform-scale-customers-mssp-delivery-architecture.md` — new ;
19. `16-quality-and-validation/validation-status-platform-scale-customers-mssp-delivery-architecture.md` — new.

No root README, `main`, object file, permission catalog/register, Settings capability contract or CAP namespace is expected to change.

## 51 mandatory gates

| Block | Count | Build-time result |
|---|---:|---|
| A — baseline / Git | 6 | PASS |
| B — approved D1–D7 integrity | 7 | PASS |
| C — source-of-truth / OPEN / objects | 6 | PASS |
| D — Security / Identity | 8 | PASS |
| E — Command / Shared / UX | 7 | PASS |
| F — registries / Requirements / roadmap | 7 | PASS |
| G1–G3 — BUILD quality / expected files / pre-publication concurrency | 3 | PASS when the final concurrency guard confirms the baseline |
| G4–G10 — publication / remote / CI applicability / closure / final | 7 | PENDING-REMOTE until actually executed |
| **Total** | **51** | **44 PASS / 7 PENDING-REMOTE / 0 FAIL** before publication |

Final `51/51` is forbidden before actual publication, remote reread, CI/status/check/workflow applicability evidence, documentary closure and final ancestry/blob verification.

## Roadmap preservation

Target and build audit: **REMOVED 0 / WEAKENED 0 / UNKNOWN 0**. Historical Phase 6 lots and their exact closure evidence remain present.

## Build-time verdict

**44 PASS / 7 PENDING-REMOTE / 0 FAIL** once the pre-publication concurrency guard passes.

This is an architecture/documentary build. It does not claim MSSP runtime support, authorization-engine implementation, Customer integration, Search/Reporting runtime changes, CRM, billing, portal, tenant hierarchy, delegated administration or production deployment.
