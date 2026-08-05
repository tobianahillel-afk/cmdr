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
# Cross-product links — Investigate through Phase 4B.2B.3A

| Transition | Context | Ownership / return |
|---|---|---|
| Case / Collection result / Artifact → Memory Intake | Case, Endpoint, Memory Image, acquisition, custody, objective, Hypothesis, return origin | Investigate; acquisition owner unchanged |
| Memory Intake → Integrity / Session | image, restrictions, support, owner and objective | Investigate |
| Session → Platform/Profile / Tool | image, candidates, limitations, permissions | Settings owns support/admin; Studio owns Tool/Call/Run |
| Process → Regions/Modules/Handles/Network | process candidate, identifiers, timestamps, provenance | Investigate observations |
| Region/module → Anomaly or Extraction | source, relations, restrictions and Tool context | Investigate; source image immutable |
| Derived Artifact → Static/Reverse | parent image, source context, extraction, Case and provenance | Investigate |
| Memory Timeline → Case Timeline | selected events, timestamp quality, uncertainty and return | Shared mechanism; Case Timeline not replaced |
| Sensitive candidate → Security review/export | masked reference, reason, policy, assurance | Security/Shared owner; access audited |
| Memory result → Evidence/Finding | observations, contradictions, lineage and uncertainty | CAP-INV-107/108/109 retain qualification |
| Memory knowledge → future Detection Engineering | behavior, processes, regions, strings, conditions and sources | future 4B.3; no rule created |
| Memory network observations → future Network Forensics | candidates, process links, timestamps and limitations | future 4B.2B.3B; no full analysis now |
| Any real-target request | Endpoint, requested outcome, risk and Case | block or route to Collection/Live Response and Govern |

Tenant, environment, selection, permissions, masked-state and return origin are preserved.
