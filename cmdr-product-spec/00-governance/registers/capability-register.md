---
id: capability-register
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-05
source-of-truth: registry
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-013
  - REQ-PROD-014
  - REQ-PROD-019
---
# Capability Register

The register is split into active shards. IDs are immutable and never recycled.

| Shard | Scope | Count | Defined | Proposed | Delivery mode |
|---|---|---:|---:|---:|---|
| `capability-register-command.md` | Command | 27 | 26 | 1 | 27 planned |
| `capability-register-investigate-foundation.md` | Investigate CAP-INV-001..114 | 22 | 21 | 1 | 22 planned |
| `capability-register-investigate-collection.md` | Investigate CAP-INV-201..215 | 15 | 15 | 0 | 15 planned |
| `capability-register-investigate-analysis-workbench.md` | Investigate CAP-INV-301..313 | 13 | 13 | 0 | 13 planned |
| `capability-register-investigate-dynamic-sandbox.md` | Investigate CAP-INV-314..328 | 15 | 15 | 0 | 15 planned |
| `capability-register-investigate-reverse-debugger.md` | Investigate CAP-INV-329..346 | 18 | 18 | 0 | 18 planned |
| **Total** | **All registered capabilities** | **110** | **108** | **2** | **110 planned** |

Each shard records owner, users, objects, inputs, outputs, actions, no-AI alternatives, criteria, requirements, OPEN and dependencies. CAP-INV-347..399 remain reserved. No CAP-INV-4xx or CAP-INV-5xx is assigned. No capability is marked validated, implemented, native or integrated.
