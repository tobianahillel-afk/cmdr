---
id: source-product-capability-inventory
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: source-material
---
# Product Capability Inventory

## Decided ownership

| Capability | Owner | Current classification |
|---|---|---|
| Operational coordination | Command | native product capability |
| Investigation lifecycle | Investigate | native product capability |
| Governed decision and response | Govern | native product capability |
| Agentic automation design | CMDR Studio | native product capability |
| Platform administration | Platform Settings | native product capability |
| Endpoint Agent EDR | Endpoint Agent | planned complete native target |
| Reporting Engine | Shared Capabilities | native shared capability |
| Generic Saved Views | Shared Capabilities | native shared capability |
| Work Queue Saved Views | Command | native product-specific capability |

## Delivery classification required before implementation claims

| Requirement | Capability | Phase 0 state |
|---|---|---|
| REQ-PROD-045 | search | proposed classification required |
| REQ-PROD-046 | investigation | proposed classification required |
| REQ-PROD-047 | collection | proposed classification required |
| REQ-PROD-048 | Live Response | proposed classification required |
| REQ-PROD-049 | endpoint telemetry | proposed classification required |
| REQ-PROD-050 | endpoint detection | proposed classification required |
| REQ-INV-001 | forensics | open classification |
| REQ-INV-002 | static analysis | open classification |
| REQ-INV-003 | reverse engineering | open classification |
| REQ-INV-004 | debugger | open classification |
| REQ-INV-005 | sandbox | open classification |
| REQ-INV-006 | Detection Engineering | open classification |
| REQ-PROD-051 | orchestration | proposed classification required |
| REQ-PROD-052 | governance | proposed classification required |
| REQ-PROD-053 | reporting | proposed classification required |
| REQ-PROD-054 | audit | proposed classification required |
| REQ-PROD-055 | automation | proposed classification required |

A wrapper is never automatically native. Temporary integrations require an explicit replacement or review plan.