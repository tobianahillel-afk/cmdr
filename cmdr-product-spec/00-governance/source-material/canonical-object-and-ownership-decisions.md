---
id: source-canonical-object-and-ownership-decisions
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: source-material
---
# Canonical Object and Ownership Decisions

## Nominal chain

Telemetry Event → Detection → Signal → Alert → Incident → Case → Evidence → Finding → Action Request → Decision → Response Run → Result.

The chain is nominal. Not every scenario creates every object.

## Ownership

| Requirement | Concept | Owner |
|---|---|---|
| REQ-OBJ-001 | Incident | Command |
| REQ-OBJ-002 | Task | Command unless explicitly specialized |
| REQ-OBJ-003 | Case | Investigate |
| REQ-OBJ-004 | Hypothesis | Investigate |
| REQ-OBJ-005 | Artifact | Investigate |
| REQ-OBJ-006 | Evidence | Investigate |
| REQ-OBJ-007 | Finding | Investigate |
| REQ-OBJ-008 | Decision | Govern |
| REQ-OBJ-009 | Approval | Govern |
| REQ-OBJ-010 | Response Run | Govern |
| REQ-OBJ-011 | Result | Govern |
| REQ-OBJ-012 | Endpoint Agent Fleet | Platform Settings |
| REQ-OBJ-013 | Endpoint | Shared model, administered by Platform Settings |
| REQ-OBJ-014 | Skill, Tool, Tool Call, Automation Agent, Agent Team, Workflow, Human Gate, Automation Run | CMDR Studio |
| REQ-OBJ-015 | Reporting Engine | Shared Capabilities |
| REQ-OBJ-016 | Generic Saved Views | Shared Capabilities |
| REQ-OBJ-017 | Work Queue Saved Views | Command |
| REQ-OBJ-018 | Inspector | Design System |
| REQ-OBJ-019 | Permission Model | Security, Permissions and Trust |

A product may display a projection of an object it does not own. It must not redefine that object.

## Required distinctions

Telemetry Event ≠ Detection; Detection ≠ Signal; Signal ≠ Alert; Alert ≠ Incident; Incident ≠ Case; Artifact ≠ Evidence; Finding ≠ Result; Action Request ≠ Decision; Decision ≠ Approval; Automation Run ≠ Response Run.