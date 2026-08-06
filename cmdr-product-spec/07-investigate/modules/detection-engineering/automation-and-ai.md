---
id: investigate-detection-engineering-automation-and-ai
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements:
  - REQ-AI-001
  - REQ-AI-002
  - REQ-AI-003
  - REQ-AI-010
  - REQ-AI-011
open_decisions:
  - OPEN-013
  - OPEN-015
  - OPEN-017
---
# Automation and AI

AI may propose authoring content from CAP-INV-401..417 and lifecycle assistance including Review Package summaries, missing evidence, candidate target compatibility, shadow/canary plans, stop criteria, health explanations, feedback grouping, candidate TP/FP/FN, tuning, bounded suppression/exception drafts, drift candidates, performance summaries, rollback/retirement drafts and Continuous Improvement Packages.

Essential workflows remain available through forms, structured editors, checklists, matrices, deterministic validators, catalogues, controlled fixtures, historical replay, target/version comparisons, health projections, metrics tables, timelines, diff, non-agentic workflows and human review.

Prohibited: silent candidate approval, Approval or Decision creation; silent target selection, promotion, activation, deactivation or rollback; automatic active suppression/exception; silent active-rule modification; Signal/Alert deletion; invented ground truth; self-permission; Govern bypass; trace removal; mandatory chatbot.

Every automated output exposes initiator, agent or engine and version, Automation Run, Tool Calls, sources, parameters, timestamp, status, errors, uncertainty, human owner and accept/modify/reject disposition.
