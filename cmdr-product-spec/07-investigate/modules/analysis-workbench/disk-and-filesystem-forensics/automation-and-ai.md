---
id: investigate-disk-filesystem-automation-ai
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-AI-001
  - REQ-AI-002
  - REQ-AI-003
open_decisions:
  - OPEN-015
---
# Automation and AI

## Allowed assistance
Propose filesystem candidates, group files or relations, summarize sessions, explain inconsistencies/errors, propose temporal correlations or persistence candidates, and prepare Hypothesis, Evidence candidate, Finding Draft or future Detection Engineering knowledge.

## Mandatory no-AI path
Deterministic Tools, viewers, trees, tables, timelines, search, filters, parsers, comparators, extractors, explicit rules, checklists, non-agentic workflows and human review cover every essential function.

## Prohibitions
No silent Tool or filesystem selection; no content execution; no encryption bypass or password attack; no secret reuse; no partial recovery as certainty; no user artifact as intent; no automatic persistence/Evidence/Finding/rule/IOC; no self-permission; no trace removal; no mandatory chatbot; no Govern bypass.

## Attribution
Initiator, producer/agent, version, Automation Run, Tool Calls, sources, parameters, timestamp, status, errors, uncertainty, human owner and accept/modify/reject disposition are visible.
