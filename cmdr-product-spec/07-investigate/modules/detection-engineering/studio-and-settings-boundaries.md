---
id: investigate-detection-engineering-studio-settings-boundaries
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
---
# Studio, Settings and Endpoint boundaries

## CMDR Studio
Studio owns Tool, Tool Call, Skill, Workflow, Automation Agent, Human Gate, Automation Run, generic Dataset/Evaluation and their versions. Investigate owns Detection Engineering business context and human dispositions.

## Platform Settings
Settings owns Data Source, Integration, Parser, schema administration, ingestion health, retention, environments, providers and secrets. Detection Engineering reads projections and prepares requests; it does not administer them.

## Endpoint Agent
Endpoint Agent owns declared telemetry capabilities, local observations and future runtime detection under policy. Authoring does not update, deploy, suppress or activate Endpoint rules.

## Govern
Future production change may require Decision or Approval. CAP-INV-417 only prepares a review package.
