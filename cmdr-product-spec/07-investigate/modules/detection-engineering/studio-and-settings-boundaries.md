---
id: investigate-detection-engineering-studio-settings-boundaries
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
open_decisions:
  - OPEN-008
  - OPEN-015
  - OPEN-017
---
# Studio, Settings and Endpoint boundaries

## CMDR Studio
Studio owns Tool, Tool Call, Skill, Workflow, Automation Agent, Human Gate, Automation Run, generic Dataset/Evaluation and their versions. Investigate owns Detection Engineering context, assessments, proposals and human dispositions. Studio automation cannot approve or execute a governed Detection change silently.

## Platform Settings
Settings owns Data Source, Integration, Parser, schema administration, configured runtimes, deployment targets, environments, tenants, providers, secrets, ingestion/execution health, retention, policies, feature flags and channels. Detection Engineering reads projections, selects proposed targets and prepares requests; it does not administer them.

## Endpoint Agent
Endpoint Agent owns declared telemetry capabilities, version/health projections and authorized local execution according to future runtime architecture. Investigate does not directly update, deploy, suppress, activate, deactivate or roll back Endpoint content.

## Govern
Govern owns Action Request, Decision, Approval, Response Run and Result for applicable class-3/4 changes. Investigate prepares evidence and observes outcomes.

## OPEN-017
Runtime, target language and portability remain undecided. No engine, adapter, vendor syntax or executable format is selected.
