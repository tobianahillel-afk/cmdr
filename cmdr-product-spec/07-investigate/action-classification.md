---
id: investigate-action-classification
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-SEC-001
  - REQ-SEC-002
open_decisions:
  - OPEN-013
---
# Action classification — Investigate through Phase 4B.2B.2B

| Class | Meaning | Reverse/Debugger examples |
|---:|---|---|
| 0 | observation | navigate, inspect representations, functions, xrefs, graphs, types, runtime state, threads, stacks, modules, events and comparisons |
| 1 | bounded isolated processing or extraction | explicitly open an authorized isolated Debugger Session, capture snapshot, create/export Derived Artifact, reproduce a session |
| 2 | reversible mutation | session changes, annotations, renames, types, breakpoints, pause/resume/step/stop, Patch Hypothesis and isolated reversible experiment |
| 3/4 | real-target or destructive authority | excluded; block or route to owner capability and Govern |

CAP-INV-329..346 use only classes 0, 1 and 2. OPEN-013 remains open. A reversible isolated experiment is not a production patch; debugger control is not Live Response.
