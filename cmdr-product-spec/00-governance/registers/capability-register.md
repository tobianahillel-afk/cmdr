---
id: capability-register
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-06
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
| `capability-register-investigate-memory-forensics.md` | Investigate CAP-INV-347..362 | 16 | 16 | 0 | 16 planned |
| `capability-register-investigate-disk-filesystem-forensics.md` | Investigate CAP-INV-363..379 | 17 | 17 | 0 | 17 planned |
| `capability-register-investigate-network-forensics.md` | Investigate CAP-INV-380..397 | 18 | 18 | 0 | 18 planned |
| `capability-register-investigate-detection-authoring.md` | Investigate CAP-INV-401..417 | 17 | 17 | 0 | 17 planned |
| **Total** | **All registered capabilities** | **178** | **176** | **2** | **178 planned** |

## Product totals
- Command: **27** capabilities.
- Investigate: **151** capabilities.
- CAP-INV-3xx: **97**.
- CAP-INV-4xx: **17**, CAP-INV-401 through CAP-INV-417.
- No CAP-INV-5xx is assigned.
- No capability is marked validated, implemented, promoted, deployed, active, native or integrated.

Each shard records owner, users, objects, inputs, outputs, actions, no-AI alternatives, criteria, Requirement IDs, OPEN decisions, dependencies, supersession and date.
