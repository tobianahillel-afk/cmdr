---
id: reverse-debugger-workflows
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
---
# Workflows and transitions

| Source | Trigger | Destination | Owner | Context | Error / return |
|---|---|---|---|---|---|
| Case | open Reverse | CAP-INV-329 | Investigate | tenant, Case, Artifact, Hypothesis, objective, return origin | return to Case; no Tool launch |
| Static Analysis | request Reverse | CAP-INV-329 | Investigate | Artifact, Analysis Session, results, Derived Artifacts, annotations, provenance | return Static Analysis |
| Dynamic Analysis | request Reverse | CAP-INV-329 | Investigate | Artifact, Runtime Artifacts, Sandbox Runs, observations, provenance | return Dynamic Sandbox |
| Reverse Intake | explicit create/resume | CAP-INV-330 | Investigate | type, architecture, restrictions, objective, compatible Tools, owner | blocked state if unsupported |
| Reverse Session | select location | CAP-INV-331/332/333 | Investigate | Artifact, location, function candidates, symbols, selection | return same view and focus |
| Disassembly | switch synchronized view | Decompilation | Investigate | location, function, xrefs, annotations, history | partial mapping is visible |
| Function | inspect references/graph | CAP-INV-334/335 | Investigate | callers, callees, data, strings, symbols | unresolved paths preserved |
| Reverse Session | prepare isolated debugger | CAP-INV-339 | Investigate | Artifact/copy, environment, objective, locations, proposed breakpoints, restrictions | real Endpoint target blocked |
| Breakpoint | explicit control | CAP-INV-340/341 | Investigate | session, location, condition, state, initiator, permission | invalid/unresolved remains visible |
| Exception | inspect trace/context | CAP-INV-343 | Investigate | event, thread, location, stack, runtime state, error | gaps preserved |
| Patch Hypothesis | explicit confirmed experiment | CAP-INV-344 | Investigate | isolated copy, change description, risk, authorization, rollback | source remains immutable |
| Reverse/Debug result | prepare Evidence/Finding | CAP-INV-346 | Investigate | sources, sessions, Artifacts, uncertainty, contradictions | qualification remains 107/108/109 |
| Reverse knowledge | prepare future detection | Future 4B.3 | future owner | behavior, functions, strings, structures, conditions, limits, sources | no rule created or deployed |

All transitions preserve tenant, environment, Artifact, selection, permissions, return origin and ownership.
