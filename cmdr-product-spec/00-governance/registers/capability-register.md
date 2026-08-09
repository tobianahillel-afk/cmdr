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

An active shard entry must identify: Capability ID, name, owner product/module, documentary status, `delivery_status`, `delivery_mode`, canonical file, primary roles, primary objects, consumers, Requirement IDs, OPEN decisions, dependencies, supersession and review date. A shard must not use the registry to claim runtime availability that is not supported by implementation/release evidence.

Registry validation checks for duplicate/recycled IDs, concurrent owners, missing canonical files, capabilities without users/objects/requirements, and `planned` capabilities presented as available. Detailed inputs, outputs, actions, states, permissions and acceptance criteria remain in the canonical capability file and are not duplicated into the global index.

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
| **Total** | **All registered capabilities** | **303** | **301** | **2** | **303 planned** |

## Product totals
- Command: **27** capabilities, 729 sections, 162 mandatory tables.
- Investigate: **243** capabilities, 6561 sections, 1458 mandatory tables.
- Govern GOV-1: **16** capabilities, 432 sections, 96 mandatory tables.
- Govern GOV-2: **17** capabilities, 459 sections, 102 mandatory tables.
- Govern cumulative: **33** capabilities, **891** sections, **198** mandatory tables.
- CAP-INV-3xx / 4xx / 5xx / 6xx / 7xx: **97 / 35 / 37 / 18 / 19**.
- Detection Engineering: **35 capabilities, 945 sections, 210 tables**.
- Threat Intelligence: **37 capabilities, 999 sections, 222 tables**.
- Cloud Analysis: **18 capabilities, 486 sections, 108 tables**.
- Mobile Forensics: **19 capabilities, 513 sections, 114 tables**.
- Phase 4B.4 Cloud + Mobile: **37 capabilities, 999 sections, 222 tables**.
- Command + Investigate: **270 capabilities, 7290 sections, 1620 tables**.
- Command + Investigate + Govern: **303 capabilities, 8181 sections, 1818 mandatory tables**.

## Current Command registry revalidation
The Command shard has 27 unique IDs, 27 owners, 27 named user sets, 27 primary-object sets and 27 concrete dependency summaries. It retains **26 defined / 1 proposed / 27 planned**, with current `native` or `integrated` claims equal to **0**. `CAP-CMD-401` remains deployment-dependent under OPEN-006.

## Govern GOV-1 registry state
The GOV-1 shard remains unchanged with 16 unique `CAP-GOV-001..016`, 16 defined / 16 planned, 432 sections and 96 mandatory tables. Its dedicated post-publication record confirms historical **180/180 PASS**. No GOV-1 capability is redefined by GOV-2.

## Govern GOV-2 registry state
The GOV-2 shard has 17 unique `CAP-GOV-017..033`, 17 named user sets, 17 owner/module assignments, 17 primary object/concept sets and concrete Requirement/OPEN/dependency evidence. All are `draft` / `defined` / `planned`; **459 sections / 102 mandatory tables**. No current native, integrated, implemented, deployed or active runtime claim is made. GOV-3 registers **0 capabilities**.

No capability is marked validated, implemented, promoted, deployed, active, native or integrated by the registry.