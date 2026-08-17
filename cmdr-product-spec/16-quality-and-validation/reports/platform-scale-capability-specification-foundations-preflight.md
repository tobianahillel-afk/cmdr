---
id: platform-scale-capability-specification-foundations-preflight
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-12
source-of-truth: quality-report
---
# Delivery Roadmap Phase 6 — Platform Scale — Capability Specification Foundations Preflight

## 1. Baseline
Exact starting HEAD: `a29069408d76e6366d7913538a49bf956ecf791e` — `docs: synchronize global capability register with Endpoint EPT-6 closure`.
Direct parent: `3f1c775c706396380610ffc1efdaaba888cec9e4`.
This preflight creates no capability, allocates/reserves no capability ID, creates no Screen ID and begins no functional Phase-6 execution.

## 2. Reconciliation validation
`BLOCKER-P6-PREFLIGHT-001 — Global capability register not synchronized with published EPT-6 closure` is **RESOLVED**.
The canonical global register now indexes the EPT-6 shard and current state without deleting the historical EPT-5 snapshots.
Current counts: **484 capabilities / 482 defined / 2 proposed / 484 planned / 13068 sections / 2904 mandatory tables**. Endpoint is **99 / 2673 / 594**.

## 3. PR / Git state
PR #2 is open, Draft, unmerged, base `main`, auto-merge disabled. `main` remains `bc1ec59e5e79ccbaba291e2984b0eb7d5e54129c`. Branch and main root README remain exactly `# cmdr`, blob `901c74cda52e28b5ff7fc425ddf28ef89f3ad875`.

## 4. Phase-5 closure
Command **PASS 27**; Investigate **PASS 243**; Govern **PASS 47**; Studio **PASS 68**; Endpoint **PASS 99**. Delivery Roadmap Phase 5 remains **PASS — capability specification complete**. This is documentary capability completeness, not implementation completeness.

## 5. Current global totals
| Domain | Capabilities |
|---|---:|
| Command | 27 |
| Investigate | 243 |
| Govern | 47 |
| Studio | 68 |
| Endpoint | 99 |
| **Total** | **484** |

Delivery state: **482 defined / 2 proposed / 484 planned**. Structural totals: **13068 sections / 2904 mandatory tables**. Requirements remain **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**. OPEN decisions remain **18**.

## 6. Canonical Phase-6 roadmap identity
Canonical roadmap: `18-roadmap-and-releases/phase-6-platform-scale.md`.
Canonical id: `roadmap-phase-6-platform-scale`.
Canonical title: **Delivery Roadmap Phase 6 — Platform Scale**.
Document status: `draft`. Documentary owner: **Product Operations Lead**.
The roadmap explicitly names MSSP aggregation, SLO/resilience, localization, advanced integrations and compliance. It is intentionally high-level and does not itself allocate capability ownership or IDs.

## 7. Corpus inventory
Audited source families include Platform Settings administration (tenant/environment, users/groups/roles, permissions, secrets/connections, models/providers, Fleet/Policy, health and retention), Shared mechanisms (search, jobs, reporting, export, collaboration, linking, localization, metrics, saved views, correlation/data-quality/entity-resolution/query/timeline and related services), Security permission/trust/privacy/audit/secret constraints, existing owner registers, Screen/permission registers, Requirements and OPEN decisions, and prior closed product capability layers.

## 8. Settings capability-layer audit
Platform Settings has **canonical functional/source documents, administrative object usage, 12 active screens and permission families**, but it does **not** have a registered `CAP-SET-*` (or equivalent) capability-contract layer in the global Capability Register.
Settings screens are not a Settings capability specification. Settings retains administrative ownership of tenant/environment configuration, user/group/role administration, integrations/providers/credentials/secrets, configuration, Endpoint Fleet and Endpoint Policies.

## 9. Shared capability-layer audit
Shared has canonical service/mechanism documents but no registered `CAP-SHR-*` capability-contract layer.

| Mechanism | Canonical object/source | Capability contract exists | Screen exists | Permission/source evidence | Phase-6 gap |
|---|---|---|---|---|---|
| Search | Global Search / Query Engine | No registered CAP-SHR | No Shared screen | Shared permission/model references | Formal capability specification |
| Linking | Object Linking Service | No | No | Shared/Security model | Formal capability specification |
| Versioning | Saved Views/report/version references | No standalone registered capability | No | source-local + Security | Normalize capability boundary, not object redesign |
| Notifications | Notification Center | No | No | Security permission model | Formal capability specification |
| Jobs | Background Jobs | No | No | Security permission model | Formal capability specification |
| Activity | cross-product activity/audit references | No registered Shared capability | No | cross-product sources | Consolidate source boundary before capability allocation |
| Trace | cross-product Trace references | No registered Shared capability | No | cross-product sources | Consolidate source boundary before capability allocation |
| Reporting | Reporting Engine / Report | No registered CAP-SHR | No | `perm.shared.report.*` | Formal capability specification |
| Export | Export Engine + Reporting | No registered CAP-SHR | No | Shared/Security | Formal capability specification |
| Collaboration | Collaboration Service | No | No | Shared/Security | Formal capability specification |
| Comments | collaboration/task sources | No standalone registered capability | No | Shared/Security | Determine whether separate capability is justified |
| Assignments | Task Inbox/collaboration sources | No standalone registered capability | No | Shared/Security | Determine boundary vs task workflow |
| Inspector | canonical experience/design projection | No Shared capability | No Shared screen | design/experience contracts | Consumer projection; do not duplicate ownership |
| Recovery | generic recovery/retry references + owner-specific recovery | No standalone registered capability | No | Shared/owner-specific sources | Define generic boundary without absorbing owner-specific recovery |

Additional Shared canonical sources relevant to Platform Scale include Business Service Catalog, Correlation Engine, Data Quality Service, Entity Resolution, Localization Engine, Metrics Engine, Mapping Engine, Saved View Engine, Timeline Engine and Telemetry Normalization. Presence of a source does not automatically require one capability per file.

## 10. Security audit
| Security concept | Canonical owner | Phase 6 consumes | Phase 6 owns | Boundary |
|---|---|---|---|---|
| Permission Model | Security Architecture | Yes | No | RBAC + ABAC + Decision Authority + SoD + step-up remain Security-owned |
| Tenant isolation | Security | Yes | No | Settings/Shared capabilities must obey isolation |
| Privacy | Security | Yes | No | data minimization/visibility constraints consumed |
| Audit integrity | Security | Yes | No | generic/platform operations emit evidence; integrity policy stays Security |
| Trust policy | Security | Yes | No | provider/integration trust consumes policy |
| Secret policy | Security | Yes | No | Settings administers references; Security constrains handling |
| Security policy | Security | Yes | No | cross-cutting constraint, not Platform ownership |

## 11. Global capability namespace inventory
Concrete registered capability families remain only:
- `CAP-CMD-*` — 27;
- `CAP-INV-*` — 243;
- `CAP-GOV-*` — 47;
- `CAP-STD-*` — 68;
- `CAP-EPT-*` — 99.
Searches found no concrete `CAP-SET-*`, `CAP-SHR-*`, `CAP-PLT-*`, `CAP-ADM-*`, `CAP-IAM-*`, `CAP-TEN-*`, `CAP-INT-*`, `CAP-CUS-*` or `CAP-DEL-*` family. No Phase-6 namespace is created by this preflight.

## 12. Reserved-ID audit
No Phase-6 capability ID or candidate namespace range is reserved. No `CAP-EPT-100+` is allocated. Future prefixes/ranges must be ratified from ownership and collision audits in the first functional lot, not in this preflight.

## 13. Ownership matrix
| Concept | Canonical owner | Existing source | Existing registered capability | Phase-6 need | Conflict |
|---|---|---|---|---|---|
| Tenant / Environment | Platform Settings | tenants-and-environments | No Settings CAP | Formalize admin capability | None |
| User / Group / Role | Platform Settings; Security constraints | users-and-roles | No Settings CAP | Formalize admin lifecycle | None |
| Permission | Security | permission-model/register | No Phase-6 CAP | Consume only | None |
| Administrative Authority | Settings + Security constraints | Settings/Security | No | Define admin action boundary | None |
| Endpoint Fleet / Policy | Platform Settings | endpoint-agent-fleet/policy sources | Endpoint has local technical CAPs, not admin ownership | Formalize admin fleet/policy capability | None |
| Integration / Provider | Platform Settings | provider/connection sources | No Settings CAP | Formalize configuration/lifecycle | None |
| Credential / Secret Reference | Settings + Security | secrets/connections + security | Endpoint/Studio consume references | Formalize administration only | None |
| Configuration | Platform Settings | module sources | No | Formalize admin configuration | None |
| Search / Linking / Version | Shared | Shared services | No CAP-SHR | Formalize shared capabilities | None |
| Notification / Job | Shared | notification-center/background-jobs | No CAP-SHR | Formalize shared operations | None |
| Activity / Trace | Shared generic boundary + producing owners | cross-product evidence | No CAP-SHR | Clarify generic contract | No blocking conflict |
| Report / Export | Shared | reporting/export | No CAP-SHR | Formalize shared capability | None |
| Collaboration / Comment / Assignment | Shared | collaboration/task sources | No CAP-SHR | Determine justified capability granularity | None |
| Recovery | Shared generic + owner-specific primitives | Shared/Endpoint/Studio/Govern | owner-specific CAPs already exist | Define only generic shared semantics | None |
| Inspector | Experience/Design System | design/experience sources | No Shared CAP | Consume, do not duplicate | None |
| Customer / Delivery | unresolved / OPEN-006 relevant | roadmap + decision sources | No | Source/owner confirmation required before allocation | Deferred, not preflight blocker |
| Quota / Limit / SLO | Platform-scale coordination with owning services | health/retention/roadmap/Shared services | No Phase-6 CAP | Formalize scale semantics by owner | None yet |
| Provider / Integration Health | Settings + producing provider/integration | Settings health/provider sources | No | Formalize admin health projection | None |

## 14. Screens
Canonical active Screen Register recount = **56**:
- Command 10;
- Investigate 14;
- Govern 9;
- Studio 11;
- Settings 12;
- Endpoint 0;
- Shared 0;
- Other 0.
Five deprecated `CMD-IWQ-*` aliases are excluded. New Screen IDs in this run: **0**.

## 15. Permissions
The global permission model remains Security-owned. Relevant families include canonical/current and historical aliases such as `perm.platform-settings.*` vs `perm.settings.*`, `perm.shared-capabilities.*` vs `perm.shared.*`, and `perm.cmdr-studio.*` vs `perm.studio.*`. These aliases are documented migration/normalization concerns. This preflight performs **0 rename, 0 permission mutation, 0 final RBAC/ABAC design**.

## 16. OPEN audit
All **18 OPEN** remain open; no decision is closed by this preflight.

| OPEN | Topic / Phase-6 relevance | Preflight blocker | Capability-spec blocker | Implementation blocker |
|---|---|---|---|---|
| OPEN-001 | canonical source/architecture decision family | No | Only where directly referenced | Potential |
| OPEN-002 | unresolved product/source decision | No | Conditional | Potential |
| OPEN-003 | unresolved product/source decision | No | Conditional | Potential |
| OPEN-004 | unresolved product/source decision | No | Conditional | Potential |
| OPEN-005 | unresolved product/source decision | No | Conditional | Potential |
| OPEN-006 | Customers and Delivery | No | Yes for Customer/Delivery capability allocation | Yes for that domain |
| OPEN-007 | Human Gate / Govern | No | No for initial Settings lot | Potential cross-product |
| OPEN-008 | Endpoint platform support | No | No; Endpoint closed documentarily | Yes for implementation claims |
| OPEN-010 | density | No | No | UX/implementation |
| OPEN-011 | unresolved cross-product decision | No | Conditional | Potential |
| OPEN-012 | unresolved cross-product decision | No | Conditional | Potential |
| OPEN-013 | default governance Class 2 | No | No for foundations | Yes for governed execution defaults |
| OPEN-014 | collection output/evidence boundary | No | No | Potential |
| OPEN-015 | Automation Run / Response Run bridge | No | No for initial Settings lot | Yes for bridge implementation |
| OPEN-016 | unresolved cross-product decision | No | Conditional | Potential |
| OPEN-017 | detection/runtime boundary | No | No for initial Settings lot | Potential |
| OPEN-018 | TI interoperability/exchange | No | Relevant only if Phase-6 integration scope reaches TI exchange | Yes for interoperability |
| OPEN-019 | dissemination/releasability | No | Relevant to reporting/sharing/customer delivery | Yes for external release |

`OPEN-009` is historical resolved evidence and is not counted among the 18 OPEN.

## 17. Requirements audit
Baseline remains **122 total / 99 conform / 20 partial / 3 absent / 0 contradictory**. No Requirement changes status in this preflight.
Platform Scale is expected to provide future evidence particularly for administrative/platform, integration/provider, shared-service, scale/SLO/resilience, localization/compliance and customer/delivery gaps where the canonical Requirement matrix currently remains partial/absent. Exact row promotion requires row-level evidence in the future functional lots; future intent alone is insufficient.

## 18. Source-of-truth matrix
| Future family | Canonical source | Supporting sources | Owner | Existing capability | Gap |
|---|---|---|---|---|---|
| Tenant/environment administration | Platform Settings tenant/environment | Security, object model | Settings | No | Capability contract |
| Identity administration | Platform Settings users/roles | Security permission model | Settings | No | Capability contract |
| Integration/provider/secret administration | Settings provider/secret sources | Security trust/secret policy | Settings | No | Capability contract |
| Shared operations | Shared jobs/notifications/metrics/etc. | Security + producing owners | Shared | No CAP-SHR | Capability contracts and boundaries |
| Shared knowledge/collaboration/localization | Shared search/link/report/export/collab/localization | Experience/Security | Shared | No CAP-SHR | Capability contracts |
| Scale/SLO/limits/resilience | Phase-6 roadmap + Settings health/retention + Shared service sources | Security/owners | owner-aligned, coordinated by Platform Scale | No | Cross-owner capability coverage without ownership collapse |
| Customer/MSSP/Delivery | Phase-6 roadmap + OPEN-006 | reporting/shared/govern/security | Not sufficiently confirmed | No | Genuine source/ownership detail before allocation |

## 19. Migration / legacy
Classification of relevant anomalies:
- corrected global Capability Register stale EPT-5 snapshot: **D — resolved stale reference**;
- permission aliases (`platform-settings/settings`, `shared-capabilities/shared`, `cmdr-studio/studio`): **C/D — historical/normalization debt**, non-blocking for preflight;
- closed product historical NOT STARTED/PENDING snapshots: **C — valid historical evidence**;
- Customer/MSSP/Delivery detail beyond roadmap/OPEN-006: **F — genuine missing source detail**, deferred and blocking only that future capability allocation, not the foundations preflight.
No active competing canonical source blocks the first future lot.

## 20. Canonical Phase-6 scope
Delivery Roadmap Phase 6 — Platform Scale is the capability-specification program for platform administration, shared platform services and scale concerns that remain outside the five already-closed product capability layers. It coordinates owner-aligned work rather than becoming a new monolithic product owner.
Problems covered: tenant/environment administration; identity administration; integrations/providers/secrets; generic shared operations and knowledge/collaboration services; localization; health/SLO/limits/retention/resilience semantics; advanced integration administration; compliance-related constraints/evidence boundaries; and Customer/MSSP/Delivery only after source/owner confirmation.
Dependencies: existing object model, Security permission/trust/privacy/audit/secret policy, Shared services, Settings modules, closed product capability contracts and OPEN decisions.
PASS for Phase 6 will require source-driven capability coverage, owner consistency, traceability, no unsupported implementation claims and a final global completeness audit.

## 21. Explicit out-of-scope
No reopening/reownership of Command, Investigate, Govern, Studio or Endpoint contracts. No product implementation, API/protocol, physical schema, final RBAC/ABAC, infrastructure/technology choice, provider/runtime selection, new Screen design, unsupported compliance certification claim, or silent closure of OPEN decisions.

## 22. Future execution-lot decomposition
Source-driven future execution lots (labels only, not roadmap subphases):
1. **Tenant, Environment and Administrative Foundations**.
2. **Identity Administration — Users, Groups, Roles and Access Assignment**.
3. **Integrations, Providers, Credentials and Secrets Administration**.
4. **Shared Platform Operations — Jobs, Notifications, Metrics, Activity, Trace and Recovery Boundaries**.
5. **Shared Knowledge, Collaboration and Localization — Search, Linking, Versioning, Reporting, Export and Collaboration**.
6. **Platform Scale Assurance — Health, SLO, Limits, Retention, Resilience and Compliance Boundaries**.
7. **Customer / MSSP / Delivery** — conditional execution lot only after OPEN-006/source/owner detail is sufficient.
These are execution-lot hypotheses validated for planning; they are not `Phase 6A/6B/...` and allocate no IDs.

## 23. First future functional lot
First source-confirmed lot: **Tenant, Environment and Administrative Foundations**.
Primary owner: Platform Settings Product Lead. Security remains a consumed constraint owner; Shared provides generic mechanisms only where referenced.
Source corpus: Settings tenant/environment/configuration roots, Security permission/isolation constraints, ownership/permission/object registers, existing Settings screens and cross-product tenant context consumers.
Dependencies: canonical tenant/environment objects, permission model, audit/trace mechanisms, Requirements and OPEN inventory.
Out-of-scope: users/groups/roles detailed lifecycle (next lot), providers/integrations/secrets (later lot), Shared service formalization, Customer/MSSP/Delivery, implementation and screens.

## 24. Candidate capability titles — no IDs
Candidate titles for the first future lot only:
- Tenant Administrative Definition and Lifecycle;
- Environment Administrative Definition and Lifecycle;
- Tenant–Environment Administrative Relationship and Scope;
- Administrative Configuration State and Provenance;
- Tenant and Environment Administrative Validation;
- Tenant and Environment Operational Context Projection;
- Administrative Change Audit Handoff;
- Cross-Product Tenant and Environment Context Contract.
These titles are candidates, not allocated capabilities. **No IDs and no reservations exist.**

## 25. Namespace recommendation
Recommend **owner-aligned future namespace families**, with Settings and Shared remaining distinct. A generic `CAP-PLT-*` would incorrectly collapse Platform Settings, Shared and Security ownership and is therefore **not recommended** on current evidence.
Exact prefix spelling/ranges are intentionally not selected or reserved here. Before first allocation, ratify the Settings-owned namespace using naming conventions + collision audit; independently ratify a Shared-owned namespace before Shared lots. Security remains consumed cross-cutting policy rather than being folded into either namespace. Customer/Delivery receives no namespace until owner/source resolution.
Lifecycle should follow immutable IDs, `draft / defined / planned` documentary state and append-only supersession conventions. Collision risk is currently low because candidate families searched are absent; migration impact is limited to permission aliases and future registry additions, not existing capability IDs.

## 26. Blockers
**No blocker prevents closing this foundations preflight.**
Resolved blocker: `BLOCKER-P6-PREFLIGHT-001` global Capability Register stale after EPT-6.
Future scoped blocker (not a preflight blocker): Customer/MSSP/Delivery capability allocation lacks sufficient canonical detail and remains dependent on OPEN-006/source ownership confirmation.

## 27. Non-blockers
The 18 OPEN decisions, permission alias normalization, absence of a pre-existing Phase-6 namespace, absence of Shared/Settings capability IDs, 0 Endpoint screens, high-level roadmap wording, and unresolved implementation/SLO values do not block this documentary preflight. They must remain explicit inputs to future lots and cannot be silently resolved.

## 28. Global completeness analysis
| Domain | Current status | Capability count | Remaining capability gap | Phase-6 responsibility |
|---|---|---:|---|---|
| Command | PASS | 27 | None known at capability layer | Preserve |
| Investigate | PASS | 243 | None known at capability layer | Preserve |
| Govern | PASS | 47 | None known at capability layer | Preserve |
| Studio | PASS | 68 | None known at capability layer | Preserve |
| Endpoint | PASS | 99 | None known at capability layer | Preserve |
| Settings | Functional/source layer present | 0 registered owner-specific capability contracts | Formal capability specification | Yes |
| Shared | Canonical mechanisms present | 0 registered CAP-SHR contracts | Formal capability specification/granularity | Yes |
| Security | Cross-cutting canonical policy | Not a separate registered capability layer | Preserve/consume, do not absorb | Cross-cutting |
| Platform Scale | Roadmap draft | 0 new Phase-6 capabilities | Coordinate remaining owner-aligned coverage | Yes |
| Customer/Delivery | Source/ownership incomplete | 0 | Resolve scope/owner before allocation | Conditional Phase-6 concern |

Completing the six source-confirmed lots is **not by itself sufficient to promise Global Capability Specification PASS** until the Phase-6 closure audit resolves whether Customer/MSSP/Delivery (and any independently discovered domain) requires a formal capability lot. Therefore current Global Capability Specification remains PARTIAL.

## 29. Closure criteria and 120-gate matrix
The 120-gate foundations matrix was re-run against the reconciled baseline. Gate groups below enumerate all gates; each individual gate in the inclusive range is PASS.

| Gates | Count | Validation family | Result |
|---|---:|---|---|
| 1–12 | 12 | exact Git baseline, parent, title, PR/base/draft/unmerged/auto-merge, main/README non-regression | PASS |
| 13–24 | 12 | register reconciliation, domain/global counts, defined/proposed/planned, sections/tables, Requirements, OPEN | PASS |
| 25–38 | 14 | Phase-6 canonical identity, roadmap and full owner/source-corpus discovery | PASS |
| 39–50 | 12 | concrete namespace inventory, candidate-family searches, reserved-ID and CAP-EPT-100+ checks | PASS |
| 51–60 | 10 | Settings source/function/screen/object/permission/roadmap capability-layer classification | PASS |
| 61–72 | 12 | Shared mechanisms, objects, source coverage and formal capability gaps | PASS |
| 73–84 | 12 | Security ownership, consumption boundaries and cross-domain ownership conflicts | PASS |
| 85–94 | 10 | Screen Register recount, zero new screens, permission namespaces/aliases and no RBAC mutation | PASS |
| 95–104 | 10 | all OPEN preserved, relevance classification, Requirements baseline/relevance/no promotion | PASS |
| 105–112 | 8 | source-of-truth matrix, legacy classification, previous blocker resolution, canonical scope/out-of-scope | PASS |
| 113–118 | 6 | lot decomposition, first lot, title candidates without IDs, namespace recommendation, global completeness | PASS |
| 119–120 | 2 | report consistency, zero-functional-work stop-line and final verdict | PASS |
| **Total** | **120** | **Mandatory foundations preflight gates** | **120/120 PASS** |

Verdict: **Platform Scale Foundations Preflight — PASS — 120/120**.
This validates readiness to prepare the first functional execution runbook. It does **not** start Delivery Roadmap Phase 6 Capability Specification.

## 30. Exact next functional run
After this preflight is published and remotely revalidated, the next run is to **prepare the exact execution runbook for `Tenant, Environment and Administrative Foundations`**. That future run may allocate capability IDs only after its owner namespace, source audit, reserved-ID audit and source-driven capability set are confirmed. Do not execute that functional lot during this preflight.

### Final stop-line
- New Phase-6 capability: **0**.
- Allocated Phase-6 IDs: **0**.
- Reserved Phase-6 IDs: **0**.
- New Screen IDs: **0**.
- Implementation: **0**.
- Functional Phase-6 lots started: **0**.
- Delivery Roadmap Phase 6 Capability Specification: **NOT STARTED**.
- Global Capability Specification: **PARTIAL**.
- Repository maturity: **PARTIAL**.
