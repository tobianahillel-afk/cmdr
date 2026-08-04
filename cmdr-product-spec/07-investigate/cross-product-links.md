---
id: investigate-cross-product-links
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
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
# Cross-product links — Investigate through Phase 4B.2B.2A

| Transition | Context | Ownership / return |
|---|---|---|
| Case or Static Analysis → Dynamic Intake | Case, Artifact, Derived Artifacts, objective, provenance, return origin | Investigate; return source workspace |
| Session → Environment | requirements, restrictions and profile | Settings owns environment |
| Session → Sandbox Run → Tools | Artifact, environment/version, profile, limits, permissions | Investigate Run context; Studio Tool Calls |
| Run → timeline/process/file/network/runtime results | Run, sources, timestamps, errors and partials | Investigate interpretation |
| Runtime Artifact → Static Analysis | lineage, restrictions, Case and return origin | Investigate |
| Unsafe environment → Settings | environment/version, symptoms, Runs, stop/reset status | Settings administration |
| Dynamic results → Evidence/Finding | observations, Runtime Artifacts, contradictions and provenance | 107/108/109 retain qualification |
| Future Reverse/Debugger | Artifact, static/dynamic results and provenance only | phase 4B.2B.2B not started |

Tenant, environment, selection, permissions and return origin are preserved.
