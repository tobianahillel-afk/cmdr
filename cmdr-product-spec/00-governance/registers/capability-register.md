---
id: capability-register
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-07
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
| **Total** | **All registered capabilities** | **270** | **268** | **2** | **270 planned** |

## Product totals
- Command: **27** capabilities, 729 sections, 162 mandatory tables.
- Investigate: **243** capabilities, 6561 sections, 1458 mandatory tables.
- CAP-INV-3xx / 4xx / 5xx / 6xx / 7xx: **97 / 35 / 37 / 18 / 19**.
- Detection Engineering: **35 capabilities, 945 sections, 210 tables**.
- Threat Intelligence: **37 capabilities, 999 sections, 222 tables**.
- Cloud Analysis: **18 capabilities, 486 sections, 108 tables**.
- Mobile Forensics: **19 capabilities, 513 sections, 114 tables**.
- Phase 4B.4 Cloud + Mobile: **37 capabilities, 999 sections, 222 tables**.
- Command + Investigate: **270 capabilities, 7290 sections, 1620 tables**.
- No capability is marked validated, implemented, promoted, deployed, active, native or integrated.
