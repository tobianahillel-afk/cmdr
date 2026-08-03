---
id: cmdr-master-product-brief
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: source-material
---
# CMDR Master Product Brief

## Mission

CMDR is an evidence-first operational cybersecurity platform connecting detection, coordination, investigation, proof, governed decision, execution, result and continuous improvement. It is not reducible to a SIEM, EDR, SOAR, ticketing tool, chatbot, dashboard or third-party wrapper.

## Current program phase

Product → Functional → UX → UI → Brand → User Journeys → Information Architecture → Functional Trust and Governance.

Low-level protocol, infrastructure and production-code choices are later-phase dependencies.

## Decided product principles

| Requirement | Decision |
|---|---|
| REQ-PROD-001 | CMDR is a complete operational platform, not a single-category tool. |
| REQ-PROD-002 | Evidence first. |
| REQ-PROD-003 | Human accountable. |
| REQ-PROD-004 | Actions safely governed. |
| REQ-PROD-005 | Every conclusion traceable. |
| REQ-PROD-006 | No duplicated object. |
| REQ-PROD-007 | Progressive disclosure. |
| REQ-PROD-008 | Context preserved across products. |
| REQ-PROD-009 | One workflow, one owner. |
| REQ-PROD-010 | Essential workflows remain usable without AI. |
| REQ-PROD-011 | Deterministic engines and rules remain first-class. |
| REQ-PROD-012 | Planned capabilities are never described as implemented. |

## Product ownership

Command owns operational coordination and Incident. Investigate owns Case, Hypothesis, Artifact, Evidence and Finding. Govern owns Decision, Approval, Response Run and Result. CMDR Studio owns agentic automation. Platform Settings owns platform administration and Endpoint Agent Fleet. Endpoint Agent is a distinct product component. Shared Capabilities owns Reporting Engine and generic Saved Views.

## Canonical nominal chain

Telemetry Event → Detection → Signal → Alert → Incident → Case → Evidence → Finding → Action Request → Decision → Response Run → Result.

The chain is nominal; not every scenario creates every object.