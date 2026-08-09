---
id: studio-action-classification
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-10
source-of-truth: canonical
open_decisions: [OPEN-013]
---
# Studio Action Classification

## STD-1 — preserved
| Class | STD-1 meaning | Representative actions | Boundary |
|---:|---|---|---|
| 0 | observation | Library browse/search/filter, Tool/Skill/version/metadata/provenance inspection | local read under permission |
| 1 | no-effect deterministic assessment | input/dependency/compatibility/eligibility validation, deterministic preview | no production effect |
| 2 | reversible Studio preparation/mutation | Tool/Skill draft metadata, version metadata, deprecation/supersession, Tool Call preparation, explicit low-impact invocation when separately authorized, links/annotations | versioned/reversible; OPEN-013 remains open |
| 3 | authority-bearing production effect | not newly defined by STD-1 | external authority/runtime owner |
| 4 | potentially irreversible/destructive effect | not defined by STD-1 | heightened Govern/runtime policy required |

## STD-2 — Workflow Builder & Orchestration
| Class | STD-2 meaning | Representative actions | Boundary |
|---:|---|---|---|
| 0 | observation | inspect Workflow/graph/version/provenance; compare; view validation | read under permission |
| 1 | no-effect deterministic assessment | graph/compatibility/input/static validation, branch preview, dependency check | never execution or authority |
| 2 | reversible/versioned Studio definition mutation | create/edit Workflow draft, add/remove/reorder steps, bindings/conditions/retry/compensation/Human Gate config, candidate Version, submit/review | OPEN-013 remains open |
| 3 | production effect | not created by STD-2 | future STD-3 + Govern/runtime rules |
| 4 | destructive/irreversible | not created by STD-2 | Govern/runtime authority required |

A Workflow definition or validation result never grants Tool invoke, Govern authority, runtime execution or Endpoint mutation.
