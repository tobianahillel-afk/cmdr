---
id: investigate-threat-intelligence-studio-settings-govern-boundaries
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
---
# Studio, Settings and Govern boundaries

## Platform Settings
Settings owns providers, integrations, configured feeds, connectors, credentials, secrets, tenants, environments, storage, retention, health, classifications, access policies and source administration. Investigate reads projections and prepares access requests; it never configures a provider or secret.

## CMDR Studio
Studio owns Tool, Tool Call, Skill, Workflow, Automation Agent, Human Gate, Automation Run, catalogue and versions. Investigate owns Intelligence purpose, candidates, assessments and human dispositions. Automated extraction or enrichment remains attributed.

## Govern
Govern owns Decision, Approval, future information release, external sharing, Response Run and Result. No content is externally shared in this phase. Future external publication, exchange, active watchlist or production action requires the applicable owner and Govern path.

## Boundary guarantee
A request, proposal, Tool output or handoff is never an access grant, Approval, publication or executed action.
