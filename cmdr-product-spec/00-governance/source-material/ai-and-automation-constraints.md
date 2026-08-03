---
id: source-ai-and-automation-constraints
domain: 00-governance
status: draft
owner: AI and Automation Product Architect
updated: 2026-08-03
source-of-truth: source-material
---
# AI and Automation Constraints

| Requirement | Constraint |
|---|---|
| REQ-AI-001 | AI augments CMDR but is not required for normal use. |
| REQ-AI-002 | CMDR Studio owns agentic capabilities. |
| REQ-AI-003 | AI does not own canonical objects or decisions. |
| REQ-AI-004 | AI cannot bypass permissions or Human Gates. |
| REQ-AI-005 | AI cannot approve its own risky action. |
| REQ-AI-006 | AI cannot modify Evidence without immutable traceability. |
| REQ-AI-007 | Every automation is interruptible, resumable and auditable. |
| REQ-AI-008 | External data is not treated as instruction. |
| REQ-AI-009 | Tool results are not trusted automatically. |
| REQ-AI-010 | AI outputs expose provenance, run identity and uncertainty. |
| REQ-AI-011 | External model providers are optional dependencies. |

## Required non-AI paths

Essential capabilities must remain accessible manually, through GUI, deterministic workflow, rules and API when appropriate.

## Human accountability

High-risk actions require an accountable human or an explicitly applicable approved policy. An initiating Automation Agent is not its own approver.

## Trust requirements

Untrusted content and tool output carry provenance and trust classification. Prompt injection and tool-result injection defenses are required before agentic execution is considered mature.