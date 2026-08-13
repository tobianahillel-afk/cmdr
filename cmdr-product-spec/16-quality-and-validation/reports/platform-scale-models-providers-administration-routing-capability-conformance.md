---
id: platform-scale-models-providers-administration-routing-capability-conformance
domain: 16-quality-and-validation
status: draft
owner: Quality Engineering Lead
updated: 2026-08-13
source-of-truth: validation
---
# Platform Scale — Models & Providers capability conformance

## Scope
Current functional lot: **Models & Providers — Model Provider Administration and Model Routing**.

Capabilities:
- `CAP-SET-010` — Model Provider Administrative Lifecycle, Validation, Model Availability and Health Projection;
- `CAP-SET-011` — Model Routing Configuration, Eligibility, Fallback Constraints and Provider Switch Provenance.

## Functional BUILD structure
- capabilities: **2/2**;
- numbered sections: **54/54**;
- mandatory capability tables: **12/12**;
- meaningful Given/When/Then scenarios: **13**;
- new Screen IDs: **0**;
- new Permission IDs: **0**;
- new canonical objects: **0**;
- new Requirement IDs: **0**;
- new OPEN decisions: **0**;
- `CAP-SET-012+` allocated/reserved: **0 / 0**.

## Canonical invariants
PASS at BUILD:
- Model Provider lifecycle preserved: `configured → validating → active → degraded → disabled`;
- Model Provider and Integration remain distinct;
- no canonical Model or routing-policy object introduced;
- Settings administrative configuration remains distinct from provider/runtime execution;
- availability, health and effective-selection facts are source-attributed projections;
- fallback configuration is not asserted as automatic failover;
- observed provider/model changes are not represented as silent Settings execution;
- `Secret Reference` dependency remains source-dependent and non-mandatory unless canonically sourced;
- Tenant, Security, Studio, Govern, Health and Administrative Audit boundaries preserved;
- existing screens and permissions reused;
- Requirements remain **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**;
- OPEN remains **18**, including OPEN-008/012/013 unresolved;
- previous Settings closures remain intact.

## Gate matrix
Mandatory total: **204 gates**.

| Family | Gates |
|---|---:|
| Git/baseline/PR/main/README/previous closures | 14 |
| source discovery/source-of-truth | 18 |
| namespace | 10 |
| two structural contracts | 12 |
| ownership/objects/terminology | 14 |
| Model Provider semantics | 16 |
| Model/no-Model-object semantics | 8 |
| routing/fallback/provider-switch semantics | 16 |
| provider runtime-executor boundary | 10 |
| Secret Reference dependency | 8 |
| Security/data policy/Tenant/Environment | 16 |
| Studio/consumers/Health/Audit | 12 |
| Screens/Permissions | 12 |
| Action classes/AI/Audit | 10 |
| Requirements/OPEN/migration | 12 |
| Registers/Roadmap/counts/diff/non-regression | 10 |
| remote publication/post-publication | 6 |
| **Total** | **204** |

## BUILD verdict
**198 PASS / 6 PENDING-REMOTE / 0 FAIL**.

The six remote-dependent gates are intentionally not closed before publication. `204/204` and `PASS AFTER POST-PUBLICATION VERIFICATION` are forbidden until the published BUILD, ancestry, remote canonical surfaces, PR/main/README invariants and CI/status/workflow applicability have been independently verified.

Platform Settings Capability Specification: **PARTIAL**. Delivery Roadmap Phase 6 Capability Specification: **PARTIAL**. Global Capability Specification: **PARTIAL**. Repository maturity: **PARTIAL**.