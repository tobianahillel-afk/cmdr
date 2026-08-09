---
id: studio-std1-action-classification
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-09
source-of-truth: canonical
open_decisions: [OPEN-013]
---
# STD-1 Action Classification

| Class | STD-1 meaning | Representative actions | Boundary |
|---:|---|---|---|
| 0 | observation | Library browse/search/filter, Tool/Skill/version/metadata/provenance inspection | local read under permission |
| 1 | no-effect deterministic assessment | input/dependency/compatibility/eligibility validation, deterministic preview | no production effect |
| 2 | reversible Studio preparation/mutation | Tool/Skill draft metadata, version metadata, deprecation/supersession, Tool Call preparation, explicit low-impact invocation when separately authorized, links/annotations | versioned/reversible; OPEN-013 remains open |
| 3 | authority-bearing production effect | not newly defined by STD-1; any Tool invocation whose underlying effect is class 3 follows exact runtime/Govern authority | external authority/runtime owner |
| 4 | potentially irreversible/destructive effect | not defined by STD-1 | heightened Govern/runtime policy required |

STD-1 never turns Tool visibility, Skill reuse or compatibility into authority. A class-3/4 effect initiated from Studio remains bounded by Tool risk metadata, Security permission, runtime owner and Govern where applicable.
