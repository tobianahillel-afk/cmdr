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
  - OPEN-007
  - OPEN-013
  - OPEN-015
---
# Cross-product links — Investigate through Phase 4B.2B.1

| Transition | Source | Destination | Context | Ownership/return |
|---|---|---|---|---|
| Case → Static Intake | Artifact open | CAP-INV-301 | Case, Artifact, Hypothesis, objective, return origin | Investigate; return Case |
| Intake → Session → Tool | explicit analyst action | CAP-INV-302/303 → Studio Tool Call | source, restrictions, parameters, permissions | Studio owns Tool Call |
| Result → Derived Artifact | transform/extract | CAP-INV-311 | parent, Tool/version, parameters, provenance | Investigate concept |
| Result → Evidence/Finding | human preparation | CAP-INV-313 → 107/108/109 | sources, contradictions, uncertainty | distinct objects |
| Static → future Dynamic/Reverse | explicit handoff | future phase | Artifact, Derived Artifacts, results, provenance | return Workbench |

Tenant, environment, selection and return origin are preserved.
