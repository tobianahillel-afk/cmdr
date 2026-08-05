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
---
# Automation and AI

AI may propose hypotheses, sources, fields, logic, sequence, thresholds, test scenarios, Expected Outcomes, explanations, candidate FP/FN, coverage mappings, gaps and review-package drafts.

Essential workflows remain available through forms, structured editors, deterministic validators, field catalogues, controlled fixtures, manual tests, deterministic replay, comparators, explicit rules, checklists, non-agentic workflows and human review.

Prohibited: silent draft/logic/source/mapping selection, automatically accepted oracle or invented ground truth, certain TP/FN qualification, promotion/deployment/activation/deactivation, production exception, Signal/Alert deletion, self-permission, trace removal or mandatory chatbot.

Every automated output exposes initiator, agent or engine and version, Automation Run, Tool Calls, sources, parameters, timestamp, status, errors, uncertainty, human owner and accept/modify/reject disposition.
