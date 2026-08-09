---
id: studio-automation-and-ai-model
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-10
source-of-truth: canonical
requirements: [REQ-AI-001, REQ-AI-002]
---
# Studio Automation and AI Model

## STD-1 — preserved
AI is optional. Library, Tool/Skill definition, validation, eligibility assessment, Tool Call preparation/lifecycle inspection, Skill discovery/reuse and provenance retain manual or deterministic paths. AI may suggest metadata, mappings/dependencies, explain compatibility/errors, links and provenance. It may not invent permissions/authority, reveal secrets, silently invoke side effects, create Approval/Decision, hide Tool Calls, erase provenance or force providers/runtimes.

## STD-2 — Workflow authoring
AI may propose graph structure, Tool/Skill steps, mappings, conditions, dependency fixes, retry configuration, compensation candidates, Human Gate placement, validation explanations and Workflow summaries.

AI may not publish, deploy or execute a Workflow; create Approval/Decision; invent permission; reveal secret; hide errors; erase provenance; extend tenant/scope silently; convert branch into Govern Decision or Human Gate into Approval.

Every essential STD-2 function has a non-AI path: manual Builder editing plus deterministic graph, mapping, compatibility, retry-safety and readiness checks.
