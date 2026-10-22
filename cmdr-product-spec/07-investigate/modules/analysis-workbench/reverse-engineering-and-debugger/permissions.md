---
id: reverse-debugger-permissions
domain: 07-investigate
status: draft
owner: Security Architecture
updated: 2026-08-05
source-of-truth: canonical
open_decisions:
  - OPEN-013
---
# Functional permission needs

| Permission need | Capabilities | Risk | Class | Step-up / separation | Owner | Finalization |
|---|---|---|---:|---|---|---|
| Reverse Analysis read / raw or sensitive Artifact read | 329–338 | sensitive binaries and IP | 0 | policy and tenant isolation | Investigate/Security | Permissions phase |
| Reverse Session create/update/close/reopen | 329–330 | shared analytical context | 2 | possible step-up; reviewer separation | Investigate/Security | OPEN-013 |
| disassembly/decompilation/function/xref/control-flow view | 331–335 | sensitive technical content | 0 | read policy | Investigate | Permissions phase |
| annotation, rename, type create/update | 334–337 | reversible interpretation | 2 | conflict and reviewer controls | Investigate/Security | OPEN-013 |
| binary comparison | 338 | multi-Artifact exposure | 0/1/2 | scope confirmation and export controls | Investigate/Security | Permissions phase |
| Debugger Session prepare/open/close | 339 | isolated execution | 2/1/2 | environment policy and explicit confirmation | Investigate/Settings | Permissions phase |
| breakpoint and execution controls | 340 | runtime state change | 2 | step-up may apply | Investigate/Security | OPEN-013 |
| runtime/thread/stack/memory/module read | 341–342 | highly sensitive runtime state | 0 | redaction and restricted environment | Investigate/Security | Permissions phase |
| runtime snapshot / Derived Artifact create/export | 341–342 | capture and diffusion | 1 | export and retention policy | Investigate/Settings | Permissions phase |
| exception/trace read and annotation | 343 | sensitive trace and interpretation | 0/2 | reviewer separation | Investigate | Permissions phase |
| Patch Hypothesis create/experiment/revert | 344 | reversible isolated modification | 2 | explicit confirmation; no production | Investigate/Security | OPEN-013 |
| Evidence/Finding/Detection handoff prepare | 346 | premature qualification/deployment | 2 | receiver review mandatory | Investigate | Permissions phase |
| automated reverse request | 329–346 | tool execution and attribution | 2 | Studio provenance and human owner | Investigate/Studio | OPEN-015 |
| restricted environment / cross-tenant use | 339–344 | isolation breach | 1/2 | deny by default and step-up | Settings/Security | Permissions phase |

No atomic matrix, namespace, RBAC/ABAC decision or final step-up rule is defined.
