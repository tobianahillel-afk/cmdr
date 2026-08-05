---
id: investigate-cross-product-links
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-PROD-004
  - REQ-PROD-005
  - REQ-PROD-008
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-015
---
# Cross-product links — Investigate through Phase 4B.2B.3B.1

| Transition | Context | Ownership / return |
|---|---|---|
| Case / Collection result / Artifact → Disk Intake | Case, source/Endpoint, Disk Image, acquisition, custody, objective, Hypothesis, return origin | Investigate; acquisition owner unchanged |
| Disk Intake → Integrity / Session | image, restrictions, support, owner, objective and partiality | Investigate |
| Session → Partition/Volume/Filesystem / Tool | image, candidates, limitations and permissions | Settings owns support/admin; Studio owns Tool/Call/Run |
| Filesystem → Navigation / File / Deleted / Journal | selected candidate, confidence, volume, restrictions and source | Investigate observations |
| File / region → Derived Artifact / Static / Reverse | parent image, source context, extraction, Case and provenance | Investigate; source immutable |
| Journal / artifacts → Disk Timeline | timestamp type, timezone, source, conflicts and gaps | Shared mechanism; Disk/Case timelines remain distinct |
| System artifact → Endpoint/Entity | persistent configuration, age, uncertainty and relation | owner projections; no current-state claim |
| User artifact → Privacy-aware review | profile context, source, timestamps and minimization | Investigate + Security/Privacy |
| Startup artifact → Persistence candidate | file/configuration, evidence for/against and cross-source links | Investigate; no confirmation |
| Disk observation → Memory correlation | source Artifact, temporal quality and relation | Investigate; Disk and Memory remain distinct |
| Disk result → Evidence/Finding | selected observations, contradictions, lineage and uncertainty | CAP-INV-107/108/109 retain qualification |
| Disk knowledge → future Detection Engineering | behavior, files, metadata, persistence conditions and sources | future 4B.3; no rule created |
| Persistent network artifact → future Network Forensics | configuration/address candidate, source and limits; no capture/connection claim | future 4B.2B.3B.2; no full analysis now |
| Any real-target request | Endpoint, requested outcome, risk and Case | block or route to Collection/Live Response and Govern |

Tenant, environment, selection, permissions, masking/minimization and return origin are preserved.
