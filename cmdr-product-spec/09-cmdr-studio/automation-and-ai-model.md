---
id: studio-std1-automation-and-ai-model
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-AI-001, REQ-AI-002]
---
# STD-1 Automation and AI Model

AI is optional. Library, Tool/Skill definition, validation, eligibility assessment, Tool Call preparation/lifecycle inspection, Skill discovery/reuse and provenance all retain manual or deterministic paths.

## AI may assist
- propose metadata/descriptions/tags;
- suggest input/output mappings or Skill dependencies;
- explain compatibility and Tool errors;
- suggest consumer links;
- summarize provenance while preserving raw references.

## AI may not
- invent permissions or authority;
- reveal or persist raw secrets;
- silently invoke a side-effecting Tool;
- create Govern Approval or Decision;
- turn a Skill into authorization;
- hide Tool Calls or erase provenance;
- force provider/runtime selection;
- become mandatory for any essential STD-1 path.

Deterministic Tool != AI Tool automatically; AI-capable Tool != Automation Agent; non-AI path does not require the same implementation engine.

## STD-2 — Workflow Builder & Orchestration addendum

AI may propose graph structure, Tool/Skill steps, mappings, conditions, dependency fixes, retry configuration, compensation candidates, Human Gate placement, validation explanations and Workflow summaries.

AI may not publish, deploy or execute a Workflow; create Approval/Decision; invent permission; reveal secret; hide errors; erase provenance; extend tenant/scope silently; convert branch into Govern Decision or Human Gate into Approval.

Every essential STD-2 function has a non-AI path: manual Builder editing plus deterministic graph, mapping, compatibility, retry-safety and readiness checks.

STD-2 does not redefine the STD-1 Tool/Tool Call/Skill AI boundaries above.

## STD-3 — Agents, Human Gates & Runtime Control addendum

AI may propose a bounded Agent plan, next-step/Tool/Skill candidate, missing-input or human-review need, stuck-Run candidate, retry/intervention/escalation candidate and source-backed Run summary.

AI may never grant permission, expand scope, authorize a production effect, create Govern Approval/Decision, bypass a Human Gate, reveal a raw secret, restart/retry silently, retry indefinitely, hide errors, alter provenance or create a Govern Result/Evidence/Finding.

Every essential STD-3 management function retains a non-AI path through Workflow-defined behavior, explicit Human Gates, deterministic eligibility/retry/readiness/status checks and human Control Room controls. STD-3 selects no agent framework, model or provider.