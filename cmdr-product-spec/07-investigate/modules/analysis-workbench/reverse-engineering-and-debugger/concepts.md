---
id: reverse-debugger-concepts
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
open_decisions:
  - OPEN-005
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Concepts

| Concept | Probable nature | Owner now | Local use | Schema status |
|---|---|---|---|---|
| Reverse Analysis Session | Investigate analytical context | Investigate | objective, Artifacts, views, knowledge, linked debugger sessions | deferred to Objects |
| Code Location | address/offset/location reference | Investigate concept | synchronized navigation and handoff | final representation deferred |
| Function / Symbol / Cross Reference | sourced analytical concepts | Investigate | interpretation, navigation and relationships | fields/cardinality deferred |
| Control Flow Graph / Call Graph | graph projections | Investigate + Shared mechanism | analysis and comparison | graph schema deferred |
| Type / Structure Definition | versioned analytical knowledge | Investigate | local interpretation only | object/state deferred |
| Debugger Session | isolated analytical session | Investigate | controls, snapshots, events and return context | object/state deferred |
| Breakpoint | session-scoped control concept | Investigate | reversible condition/location control | implementation deferred |
| Runtime State Snapshot / Stack Frame | attributed observations | Investigate | inspection and comparison | object/state deferred |
| Memory Region / Module Observation | debugger-scoped observations | Investigate | view and extraction; not Memory Forensics | object/state deferred |
| Exception Event / Debug Trace | attributed runtime events | Investigate + Shared trace mechanism | analysis and handoff | event schema deferred |
| Patch Hypothesis | reversible analytical hypothesis | Investigate | isolated experiment only | object/state deferred |
| Reproducibility Assessment | analytical disposition | Investigate | attempts, gaps and dispute | object/state deferred |
| Tool / Tool Call / Automation Run | automation objects | CMDR Studio | selection, execution and provenance | Studio owner |
| Environment / policies | administered resources | Platform Settings | read/select only | Settings owner |
| Decision / Response Run / Result | governed response objects | Govern | read/link only if real-target work is requested | Govern owner |

No complete schema, JSON Schema, final cardinality, object state machine or internal address format is defined.
