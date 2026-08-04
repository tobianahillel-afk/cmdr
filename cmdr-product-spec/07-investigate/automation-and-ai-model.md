---
id: investigate-automation-and-ai-model
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-AI-001
  - REQ-AI-002
  - REQ-AI-003
  - REQ-AI-006
  - REQ-AI-010
  - REQ-AI-011
open_decisions:
  - OPEN-007
  - OPEN-013
  - OPEN-015
---
# Automation and AI model — Investigate through Phase 4B.2A

## Authority
| Producteur | Observe | Propose | Mutate class 2 | Execute class 3/4 | Create Decision |
|---|---:|---:|---:|---:|---:|
| Authorized human | yes | yes | by permission/policy | only through Govern contract | no |
| Rule/deterministic engine | yes | yes | only explicit reversible contract | no | no |
| Workflow | yes | yes | deployed and permissioned only | orchestrates after Govern | no |
| Automation Agent | yes | yes | proposal by default | no autonomous sensitive action | no |
| Endpoint Agent | local observe/execute | reports capabilities/results | only authorized operation | executes governed run when contracted | no |
| Govern | consumes context | returns gate/status | owns authority when required | owns Decision/Response Run | yes |

## Phase 4B.2A uses
| Function | Manual | Deterministic | Automatable | AI possible | No-AI alternative |
|---|---:|---:|---:|---:|---|
| Summarize Endpoint state | yes | projection aggregation | yes | attributed summary | raw fields and filters |
| Propose collection scope/profile | yes | policies/profiles/checklists | yes | editable suggestion | forms and profiles |
| Detect overbroad scope | yes | validators/limits | yes | explanation | deterministic warnings |
| Explain error/partial result | yes | error codes/grouping | yes | summary | raw error and inspector |
| Propose next collection | yes | rules/workflow | yes | suggestion | procedures and analyst judgment |
| Summarize transcript/results | yes | aggregation | yes | draft | transcript/search/filter |
| Propose Evidence candidate/Finding draft | yes | source links/checks | yes | proposal only | human qualification/drafting |
| Prepare containment request | yes | completeness/policy | yes | draft | CAP-INV-215/113 forms |
| Open session or execute sensitive operation | explicit human/gate | contract only | workflow after authority | never autonomous | manual authorized action |

## Mandatory provenance
Initiator; producer type/ID/version; Automation Run; Tool Calls; source objects/versions; functional parameters and scope; timestamp/timezone; status/errors; uncertainty; human owner; accept/modify/reject; correlation ID and return origin.

## Prohibitions
No mandatory chatbot; no silent session; no silent sensitive command; no silent scope expansion; no silent Artifact deletion or Evidence qualification; no automatic Finding confirmation; no containment bypass; no self-permission; no hidden Tool Calls or trace removal.

## Govern and Human Gate
OPEN-007 remains open: Human Gate is not automatically a Govern Decision. OPEN-013 remains open for class 2 defaults. OPEN-015 remains open for Automation Run/Response Run. Absence of an AI provider changes suggestions only, never essential operation.