---
id: studio-std1-screen-capability-map
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-10
source-of-truth: canonical
---
# Studio Screen / STD-1 Capability Map

No Studio screen is rewritten by STD-1. Existing 11 Screen IDs stay active/draft. Capability links below are conceptual navigation/consumption only.

| Screen ID | Surface | STD-1 relevance | Disposition |
|---|---|---|---|
| STD-LIB-001 | Library | CAP-STD-001/002/006/014/015/016 | existing screen unchanged |
| STD-SKL-001 | Skill Detail | CAP-STD-010..014/015/016 | existing screen unchanged |
| STD-BLD-001 | Builder | CAP-STD-003/004/010..012 references | existing screen unchanged; no detailed rewrite |
| STD-AGT-001 | Automation Agent Detail | Tool/Skill refs only | future STD-3 semantics not started |
| STD-ATM-001 | Agent Team Detail | Tool/Skill refs only | future STD-3 semantics not started |
| STD-WFL-001 | Workflow Detail | Tool/Skill refs only | future STD-2 semantics not started |
| STD-CTL-001 | Control Room | Tool Call projection ref only | future STD-3 semantics not started |
| STD-ASR-001 | Assurance | version/reference only | future STD-4 semantics not started |
| STD-EVL-001 | Evaluation Detail | Skill/Tool version refs | future STD-4 semantics not started |
| STD-SIM-001 | Simulation Detail | Tool/Skill refs | future STD-4 semantics not started |
| STD-DEP-001 | Deployment | version refs | future STD-4 semantics not started |

New Screen IDs = **0**. Detailed rewrites = **0**. Wireframes/final buttons/columns/filters/animations/shortcuts = **0**. Endpoint Screen IDs created = **0**.

## STD-2 capability-link addendum

The historical STD-1 mapping above is preserved verbatim as evidence. STD-2 updates only conceptual links; it does not rewrite any screen.

| Screen ID | STD-2 relevance | Disposition |
|---|---|---|
| STD-BLD-001 | primary: CAP-STD-017..033 | existing screen unchanged; no detailed rewrite |
| STD-WFL-001 | primary: CAP-STD-017/019..033 | existing screen unchanged; no detailed rewrite |
| STD-LIB-001 | Workflow/version discovery references | existing screen unchanged |
| STD-SKL-001 | CAP-STD-021/023 dependency references | existing screen unchanged |
| STD-CTL-001 | future runtime projection only | STD-3 historical pre-start boundary |
| STD-AGT-001 | boundary only | STD-3 historical pre-start boundary |
| STD-ATM-001 | boundary only | STD-3 historical pre-start boundary |
| STD-ASR-001 | readiness/pre-publish handoff only | STD-4 NOT STARTED |
| STD-SIM-001 | static/no-effect validation != product simulation | STD-4 NOT STARTED |
| STD-EVL-001 | boundary only | STD-4 NOT STARTED |
| STD-DEP-001 | approved-for-publishing-candidate only | STD-4 NOT STARTED |

After STD-2: new Screen IDs = **0**; detailed rewrites/wireframes/final buttons/columns/filters/animations/shortcuts = **0**; Endpoint Screen IDs = **0**.

## STD-3 capability-link addendum

STD-3 updates conceptual capability links only. No screen source, route, buttons, columns, filters, shortcuts, animations or wireframe is rewritten.

| Screen ID | STD-3 relevance | Disposition |
|---|---|---|
| STD-AGT-001 | primary CAP-STD-034..036/038/039 | existing screen unchanged |
| STD-ATM-001 | primary CAP-STD-037..039 | existing screen unchanged |
| STD-CTL-001 | primary CAP-STD-042..051 plus oversight/gates | existing screen unchanged |
| STD-WFL-001 | Run launch/step/Human Gate references | existing STD-2 screen unchanged |
| STD-BLD-001 | Agent/Team config references; no runtime UI rewrite | existing screen unchanged |
| STD-LIB-001 | Agent/Team/version discovery references | existing screen unchanged |
| STD-SKL-001 | Agent Skill-access references | existing screen unchanged |
| STD-ASR-001 | Agent/runtime evidence consumer boundary only | STD-4 NOT STARTED |
| STD-SIM-001 | future simulation of Agent/Run only | STD-4 NOT STARTED |
| STD-EVL-001 | future evaluation of Agent/Run only | STD-4 NOT STARTED |
| STD-DEP-001 | runtime refs do not imply deployment/promotion | STD-4 NOT STARTED |

After STD-3 content: new Screen IDs = **0**; detailed rewrites/wireframes/final buttons/columns/filters/animations/shortcuts = **0**; Endpoint Screen IDs = **0**.