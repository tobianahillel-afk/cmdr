---
id: studio-std1-screen-capability-map
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-09
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
| STD-CTL-001 | future runtime projection only | STD-3 NOT STARTED |
| STD-AGT-001 | boundary only | STD-3 NOT STARTED |
| STD-ATM-001 | boundary only | STD-3 NOT STARTED |
| STD-ASR-001 | readiness/pre-publish handoff only | STD-4 NOT STARTED |
| STD-SIM-001 | static/no-effect validation != product simulation | STD-4 NOT STARTED |
| STD-EVL-001 | boundary only | STD-4 NOT STARTED |
| STD-DEP-001 | approved-for-publishing-candidate only | STD-4 NOT STARTED |

After STD-2: new Screen IDs = **0**; detailed rewrites/wireframes/final buttons/columns/filters/animations/shortcuts = **0**; Endpoint Screen IDs = **0**.

## STD-3 capability-link addendum

The prior `STD-3 NOT STARTED` cells above are preserved as the exact STD-2 closure snapshot; current STD-3 links follow.

| Screen ID | STD-3 relevance | Disposition |
|---|---|---|
| STD-AGT-001 | primary CAP-STD-034..036/038/039 | existing screen unchanged |
| STD-ATM-001 | primary CAP-STD-037..039 | existing screen unchanged |
| STD-CTL-001 | primary CAP-STD-042..051 plus oversight/gates | existing screen unchanged |
| STD-WFL-001 | Run launch/step/Human Gate references | existing STD-2 screen unchanged |
| STD-BLD-001 | Agent/Team configuration references | existing screen unchanged |
| STD-LIB-001 | Agent/Team/version discovery references | existing screen unchanged |
| STD-SKL-001 | Agent Skill-access references | existing screen unchanged |
| STD-ASR-001 | Agent/runtime evidence consumer boundary | STD-4 NOT STARTED |
| STD-SIM-001 | future simulation of Agent/Run | STD-4 NOT STARTED |
| STD-EVL-001 | future evaluation of Agent/Run | STD-4 NOT STARTED |
| STD-DEP-001 | runtime refs do not imply deployment/promotion | STD-4 NOT STARTED |

After STD-3 content: new Screen IDs = **0**; detailed rewrites/wireframes/final buttons/columns/filters/animations/shortcuts = **0**; Endpoint Screen IDs = **0**.

## STD-4 capability-link addendum

The historical `STD-4 NOT STARTED` cells above remain preserved snapshots. STD-4 adds conceptual links only.

| Screen ID | STD-4 relevance | Disposition |
|---|---|---|
| STD-EVL-001 | primary CAP-STD-052..054/059/060 | existing screen unchanged; “Approve result” does not mean Govern Approval |
| STD-SIM-001 | primary CAP-STD-055..057/059/060 | existing screen unchanged; no sandbox engine selected |
| STD-DEP-001 | primary CAP-STD-061..067 | existing screen unchanged; “Rollback” normalized functionally as Studio Deployment Reversion |
| STD-ASR-001 | primary CAP-STD-052..060/068 | existing screen unchanged; “Approve promotion” does not create Govern Approval |
| STD-LIB-001 | published/deprecated/version discovery | existing screen unchanged |
| STD-SKL-001 | evaluated/released Skill version links | existing screen unchanged |
| STD-WFL-001 | evaluated/released Workflow version links | existing screen unchanged |
| STD-BLD-001 | readiness/publishing handoff refs | existing screen unchanged |
| STD-AGT-001 | evaluated/released Agent refs | existing screen unchanged |
| STD-ATM-001 | team-related assurance refs | existing screen unchanged |
| STD-CTL-001 | deployment/runtime health correlation only | existing screen unchanged; deployment health != Run/business success |

After STD-4 content: new Screen IDs = **0**; detailed rewrites/wireframes/final buttons/columns/filters/animations/shortcuts = **0**; Endpoint Screen IDs = **0**.