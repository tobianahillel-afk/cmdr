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
  - OPEN-013
  - OPEN-015
---
# Cross-product links — Investigate through Phase 4B.2B.2B

| Transition | Context | Ownership / return |
|---|---|---|
| Case / Static / Dynamic Analysis → Reverse Intake | Case, Artifact, static/dynamic results, Derived/Runtime Artifacts, objective, provenance, return origin | Investigate; return source workspace |
| Reverse Intake → Reverse Session | type, architecture, restrictions, objective, compatible Tools, owner | Investigate |
| Reverse Session → Tool | Artifact, location, scope, parameters and permission | Studio owns Tool/Tool Call/Automation Run |
| Reverse/Debugger → Environment | required capabilities, isolation and restrictions | Platform Settings owns environment administration |
| Reverse Session ↔ code views | location, function, symbols, xrefs, annotations and history | Investigate; selection preserved |
| Reverse Session → Debugger Session | Artifact or isolated copy, environment, objective, locations, proposed breakpoints, restrictions | Investigate; no real Endpoint |
| Debugger events → Trace/Timeline | session, controls, threads, locations, errors and gaps | Shared mechanism; Investigate interpretation |
| Isolated experiment → Derived Artifact | source/copy, change description, before/after, rollback and provenance | Investigate; original immutable |
| Reverse/Debug results → Evidence/Finding | sources, sessions, uncertainty and contradictions | CAP-INV-107/108/109 retain qualification |
| Reverse knowledge → future Detection Engineering | behavior, functions, strings, structures, conditions, limits and sources | future 4B.3 owner; no rule created |
| Any real-target request | target, requested outcome, risk and Case | block or route to Collection/Live Response and Govern |

Tenant, environment, selection, permissions and return origin are preserved.
