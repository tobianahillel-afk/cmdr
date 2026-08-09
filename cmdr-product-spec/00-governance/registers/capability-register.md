---
id: capability-register
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-09
source-of-truth: registry
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-013
  - REQ-PROD-014
  - REQ-PROD-015
  - REQ-PROD-019
---
# Capability Register

The register is split into active shards. IDs are immutable and never recycled. The global file is the authoritative index and total; each shard carries the row-level evidence needed for owner, role, object, delivery and traceability checks.

## Required row semantics
An active shard entry identifies Capability ID, name, owner product/module, documentary status, `delivery_status`, `delivery_mode`, canonical file, primary roles, primary objects, consumers, Requirement IDs, OPEN decisions, dependencies, supersession and review date. A shard must not use the registry to claim runtime availability unsupported by implementation/release evidence.

Registry validation checks duplicate/recycled IDs, concurrent owners, missing canonical files, capabilities without users/objects/requirements and `planned` capabilities presented as available. Detailed inputs, outputs, actions, states, permissions and acceptance criteria remain in canonical capability files.

| Shard | Scope | Count | Defined | Proposed | Delivery mode |
|---|---|---:|---:|---:|---|
| `capability-register-command.md` | Command | 27 | 26 | 1 | 27 planned |
| `capability-register-investigate-foundation.md` | Investigate CAP-INV-001..114 | 22 | 21 | 1 | 22 planned |
| `capability-register-investigate-collection.md` | CAP-INV-201..215 | 15 | 15 | 0 | 15 planned |
| `capability-register-investigate-analysis-workbench.md` | CAP-INV-301..313 | 13 | 13 | 0 | 13 planned |
| `capability-register-investigate-dynamic-sandbox.md` | CAP-INV-314..328 | 15 | 15 | 0 | 15 planned |
| `capability-register-investigate-reverse-debugger.md` | CAP-INV-329..346 | 18 | 18 | 0 | 18 planned |
| `capability-register-investigate-memory-forensics.md` | CAP-INV-347..362 | 16 | 16 | 0 | 16 planned |
| `capability-register-investigate-disk-filesystem-forensics.md` | CAP-INV-363..379 | 17 | 17 | 0 | 17 planned |
| `capability-register-investigate-network-forensics.md` | CAP-INV-380..397 | 18 | 18 | 0 | 18 planned |
| `capability-register-investigate-detection-authoring.md` | CAP-INV-401..417 | 17 | 17 | 0 | 17 planned |
| `capability-register-investigate-detection-lifecycle.md` | CAP-INV-418..435 | 18 | 18 | 0 | 18 planned |
| `capability-register-investigate-threat-intelligence-foundations.md` | CAP-INV-501..518 | 18 | 18 | 0 | 18 planned |
| `capability-register-investigate-threat-intelligence-analysis-and-products.md` | CAP-INV-519..537 | 19 | 19 | 0 | 19 planned |
| `capability-register-investigate-cloud-analysis.md` | CAP-INV-601..618 | 18 | 18 | 0 | 18 planned |
| `capability-register-investigate-mobile-forensics.md` | CAP-INV-701..719 | 19 | 19 | 0 | 19 planned |
| `capability-register-govern-gov1.md` | Govern GOV-1 CAP-GOV-001..016 | 16 | 16 | 0 | 16 planned |
| `capability-register-govern-gov2.md` | Govern GOV-2 CAP-GOV-017..033 | 17 | 17 | 0 | 17 planned |
| `capability-register-govern-gov3.md` | Govern GOV-3 CAP-GOV-034..047 | 14 | 14 | 0 | 14 planned |
| **Total** | **All registered capabilities** | **317** | **315** | **2** | **317 planned** |

## Product totals
- Command: **27** capabilities, 729 sections, 162 mandatory tables.
- Investigate: **243** capabilities, 6561 sections, 1458 mandatory tables.
- Govern GOV-1: **16** capabilities, 432 sections, 96 mandatory tables.
- Govern GOV-2: **17** capabilities, 459 sections, 102 mandatory tables.
- Govern GOV-3: **14** capabilities, 378 sections, 84 mandatory tables.
- Govern cumulative: **47** capabilities, **1269** sections, **282** mandatory tables.
- CAP-INV-3xx / 4xx / 5xx / 6xx / 7xx: **97 / 35 / 37 / 18 / 19**.
- Detection Engineering: **35 capabilities, 945 sections, 210 tables**.
- Threat Intelligence: **37 capabilities, 999 sections, 222 tables**.
- Cloud Analysis: **18 capabilities, 486 sections, 108 tables**.
- Mobile Forensics: **19 capabilities, 513 sections, 114 mandatory tables**.
- Phase 4B.4 Cloud + Mobile: **37 capabilities, 999 sections, 222 tables**.
- Command + Investigate: **270 capabilities, 7290 sections, 1620 tables**.
- Command + Investigate + Govern: **317 capabilities, 8559 sections, 1902 mandatory tables**.

## Command registry state
The Command shard remains 27 unique IDs, 26 defined / 1 proposed / 27 planned, with current native/integrated claims equal to 0. `CAP-CMD-401` remains deployment-dependent under OPEN-006.

## Govern registry state
- GOV-1: 16 unique `CAP-GOV-001..016`, 16 defined/planned, 432 sections, 96 tables; historical 180/180 PASS evidence retained.
- GOV-2: 17 unique `CAP-GOV-017..033`, 17 defined/planned, 459 sections, 102 tables; historical 190/190 PASS evidence retained.
- GOV-3: 14 unique `CAP-GOV-034..047`, 14 defined/planned, 378 sections, 84 tables; **200/200 post-publication PASS evidence recorded**.
- Govern capability specification: **PASS** across 47 capabilities / 1269 sections / 282 tables.

No capability is marked implemented, promoted, deployed, active, native or integrated by this registry; `PASS` is documentary capability-specification conformance only.