---
id: command-automation-and-ai-model
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-AI-001
  - REQ-AI-002
  - REQ-AI-003
  - REQ-AI-005
  - REQ-AI-010
  - REQ-PROD-010
  - REQ-PROD-013
---
# Automation and AI model — Command

## Authority ladder

| Producer | Can observe | Can propose | Can mutate Command | Can decide Govern | Evidence |
|---|---:|---:|---:|---:|---|
| Human authorized | yes | yes | according to permission/action class | only with Govern authority | actor + justification |
| Rule | yes | yes | class 2 only if explicitly policy-authorized | no | rule/version/factors |
| Deterministic engine | yes | yes | validation/aggregation; mutation only by approved contract | no | engine/version/input |
| Workflow | yes | yes | explicit deployed steps and permissions | no self-approval | workflow/version/run |
| Automation Agent | yes | yes | proposal by default; no silent mutation | no | agent/version/run/tools |
| Govern | yes | n/a | returns Decision/Run/Result projections | yes within authority | policy/authority/audit |

## Command uses

- summary, priority suggestion, assignment suggestion, missing-data detection and handover draft ;
- deterministic aggregation, deduplication checks, SLA calculation and freshness ;
- workflow launch only for deployed, authorized versions ;
- Automation Tray remains contextual and closed by default.

## Required alternative without AI

Every capability must expose manual actions and deterministic/rule behavior sufficient for its essential outcome. Missing model/provider changes only suggestion availability.

## Prohibited

- AI creates or approves a Decision ;
- AI changes effective priority or owner silently ;
- AI hides source data or uncertainty ;
- AI becomes Work Queue or Incident Commander ;
- agent grants itself permission or bypasses Govern.

## Audit

Every proposal records producer type, ID/version, run, sources, factors, time, disposition and human actor when accepted.
