---
id: investigate-mobile-forensics-automation-and-ai
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-07
source-of-truth: canonical
requirements: [REQ-AI-001, REQ-AI-002, REQ-AI-003, REQ-AI-006, REQ-AI-010, REQ-AI-011]
open_decisions: [OPEN-011, OPEN-013, OPEN-015]
---
# Mobile Forensics Automation and AI

AI is optional. Every essential workflow is fully available through viewers, tables, trees, timelines, graph views with tabular alternatives, filters, search, comparison, deterministic parsers/extractors, explicit rules, checklists, non-agentic workflows and human review.

## Allowed assistance
AI may propose a bounded scope, a platform candidate, application groupings, candidate relations between records, summaries of authorized communications, anomaly candidates, correlations, Hypotheses, duplicate/recovered-record candidates, Evidence Candidate content, Finding Draft content, Detection Engineering packages or Threat Intelligence handoffs.

## Prohibitions
- no silent platform selection or Tool execution;
- no private/sensitive reveal without permission;
- no secret/token/key use, password/code testing or credential validation;
- no device unlock, bypass, rooting, jailbreak, exploit or offensive extraction;
- no automatic attribution of a device/account/contact/number/message to a person;
- no location record presented as certain human presence;
- no anomaly/persistence/suspicious app automatically confirmed as compromise or malware;
- no automatic Evidence qualification or Finding confirmation;
- no rule creation/deployment, response execution, permission grant or trace deletion;
- no mandatory chatbot and no essential model dependency.

## Attribution contract
Every automated output records initiator, producer/agent/model and version, Tool/Tool Calls, Automation Run, sources, parameters, timestamp, status, errors, uncertainty, human owner, accept/modify/reject disposition and trace. Restricted content is not sent to a model without explicit source permission and policy.

## Deterministic priority
Integrity checks, structural parsing, timestamp normalization, duplicate detection, field comparison, bounded extraction/recovery and completeness checks should use deterministic methods where available. AI may explain or propose but does not replace source records or deterministic output.
