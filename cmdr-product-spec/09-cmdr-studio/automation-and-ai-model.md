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
