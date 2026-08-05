---
id: investigate-analysis-workbench-dependencies
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-PROD-013
  - REQ-PROD-014
  - REQ-PROD-019
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Analysis Workbench dependencies

| Dependency | Workbench use | Owner retained |
|---|---|---|
| Case, Hypothesis, Artifact, Evidence and Finding | context, source, reasoning and handoff | Investigate owner capabilities |
| Collection Request, Job, integrity, custody and provenance | acquisition source and limitations | Collection / producer projections |
| Dynamic Sandbox | isolated behavior and runtime Artifact correlation | separate Investigate module |
| Static, Reverse, Memory, Disk and Network families | explicit cross-analysis handoffs | each family retains its analytical concepts |
| Tool, Tool Call, Workflow and Automation Run | processing and provenance | CMDR Studio |
| Fleet, sensors, policies, providers, storage, retention, health and time synchronization | administrative projections | Platform Settings |
| Entity, Graph, Timeline, Jobs, Notifications, Trace, Activity, Linking, Search, Export, Reporting, Versioning and Recovery | shared mechanisms | Shared Capabilities |
| Decision, Approval, Response Run and Result | real-target authority and returned outcomes | Govern |
| Detection, Signal, Alert and Incident | operational context only | Command |

No technical API, protocol, engine, command, capture format, packet field model or implementation dependency is selected.
