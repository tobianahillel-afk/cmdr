---
id: investigate-automation-and-ai-model
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-07
source-of-truth: canonical
requirements: [REQ-AI-001, REQ-AI-002, REQ-AI-003, REQ-AI-006, REQ-AI-010, REQ-AI-011]
open_decisions: [OPEN-011, OPEN-012, OPEN-013, OPEN-015, OPEN-018, OPEN-019]
---
# Automation and AI model — Investigate through Phase 4B.4B

AI is optional. For Threat Intelligence it may propose analytic questions, competing hypotheses, claim groupings, corroboration/contradiction candidates, assessment drafts, Product Plans, limitations, audience/marking candidates, Watchlist Definition drafts, Operationalization Packages, Change Assessments, collection gaps, corrections and Continuous Improvement Packages.

For Cloud Analysis it may propose a provider-neutral scope, relevant service/resource groups, identity relations, permission paths, event summaries, anomaly candidates, correlations, Hypotheses, Evidence Candidate/Finding Draft content and Detection Engineering packages. Every suggestion preserves source, timestamp, permission, restriction, uncertainty and human disposition.

For Mobile Forensics it may propose a bounded device/source scope, platform candidate, application groupings, candidate relations between records, summaries of authorized communications, duplicate/recovered-record candidates, anomaly candidates, correlations, Mobile Hypotheses, Evidence Candidate/Finding Draft content and Detection/TI packages. It may not silently select a platform, run a Tool, attribute a person, infer authorship or certain human presence, reveal private content, validate/use credentials, unlock/bypass/root/jailbreak a device or mutate a target.

Every essential workflow remains available through forms, hypothesis matrices, evidence-for/against tables, deterministic parsing/validation, viewers, catalogues, search, comparisons, checklists, templates, tables, trees, timelines, graph with tabular alternative, diff, explicit extraction/recovery Tools where authorized, notifications, non-agentic workflows and human review.

Prohibitions: no automatic attribution; no hypothesis or assessment silently confirmed; no conclusion, Release Recommendation or audience silently accepted; no Approval; no marking removed; no external share; no active watchlist/Indicator/rule/block; no Signal deletion; no permission grant; no Cloud target mutation; no mobile device acquisition/mutation; no secret reveal without permission; no secret use; no password/code testing; no unlock/bypass/root/jailbreak; no exploit-path, vulnerability, compromise, persistence, malware, Evidence or Finding confirmation; no history/trace removal; no mandatory chatbot or model dependency.

Mandatory attribution: initiator, producer/agent/model and version, Tools, Tool Calls, Automation Run, sources, parameters, timestamp, status, errors, uncertainty, human owner and accept/modify/reject disposition. Restricted or sensitive content is never sent to a model without explicit permission and policy. AI assessment remains distinct from analyst assessment; anomaly scores are not compromise and confidence is not calibrated probability.
