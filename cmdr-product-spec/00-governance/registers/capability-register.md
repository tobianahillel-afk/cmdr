---
id: capability-register
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-04
source-of-truth: registry
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-013
  - REQ-PROD-014
  - REQ-PROD-019
---
# Capability Register

The register is split into active shards to keep ownership and validation auditable. IDs are immutable and never recycled.

| Shard | Scope | Count | Defined | Proposed | Delivery mode |
|---|---|---:|---:|---:|---|
| `capability-register-command.md` | Command | 27 | 26 | 1 | 27 planned |
| `capability-register-investigate-foundation.md` | Investigate CAP-INV-001..114 | 22 | 21 | 1 | 22 planned |
| `capability-register-investigate-collection.md` | Investigate CAP-INV-201..215 | 15 | 15 | 0 | 15 planned |
| `capability-register-investigate-analysis-workbench.md` | Investigate CAP-INV-301..313 | 13 | 13 | 0 | 13 planned |
| `capability-register-investigate-dynamic-sandbox.md` | Investigate CAP-INV-314..328 | 15 | 15 | 0 | 15 planned |
| **Total** | **All products currently registered** | **92** | **90** | **2** | **92 planned** |

Each shard records ID, name, product, module, status, delivery status/mode, canonical file, roles, objects, consumers, requirements, OPEN, dependencies, supersession and review date. CAP-INV-329..399 remain reserved; no CAP-INV-4xx or CAP-INV-5xx is assigned. No capability is marked validated, implemented, native or integrated.
