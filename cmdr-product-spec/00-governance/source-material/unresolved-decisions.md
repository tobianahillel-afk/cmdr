---
id: source-unresolved-decisions
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-04
source-of-truth: source-material
requirements:
  - REQ-PROD-048
  - REQ-PROD-062
  - REQ-BRAND-008
---
# Unresolved Decisions

A functional specification or Draft value is not a decision. Closure requires approved evidence, owner and dependent updates.

## Open decisions — 15
| ID | Decision | Status | Owner | Target phase | Requirement |
|---|---|---|---|---|---|
| OPEN-001 | Investigate palette direction | open | Brand Design Lead | Phase 2 review | REQ-PROD-048 |
| OPEN-002 | Govern palette direction | open | Brand Design Lead | Phase 2 review | REQ-PROD-049 |
| OPEN-003 | CMDR Studio palette direction | open | Brand Design Lead | Phase 2 review | REQ-PROD-050 |
| OPEN-004 | Final typography stack and licensing | open | Design System Lead | Phase 2 review | REQ-PROD-051 |
| OPEN-005 | Initial forensic engines | open | Investigate Product Lead | Phase 4B.2B/8 | REQ-PROD-052 |
| OPEN-006 | Customers and Delivery applicability | open | Command Product Lead | Phase 4A review | REQ-PROD-053 |
| OPEN-007 | Human Gate and Govern relationship | open | Security Architecture | Phase 4C/7 | REQ-PROD-054 |
| OPEN-008 | Endpoint platform support | open | Endpoint Agent Product Lead | Phase 4D/8 | REQ-PROD-055 |
| OPEN-010 | Density by role and activity | open | UX Architecture | Phase 3 review | REQ-PROD-057 |
| OPEN-011 | Mobile forensic scope | open | Investigate Product Lead | Phase 4B.2B | REQ-PROD-058 |
| OPEN-012 | Cloud analysis scope | open | Investigate Product Lead | Phase 4B.2B/4B.3 | REQ-PROD-059 |
| OPEN-013 | Default governance for Action Class 2 | open | Security Architecture | Phase 4C/7 | REQ-PROD-060 |
| OPEN-014 | Artifact versus Attachment | open | Investigate Product Lead | Phase 7 | REQ-PROD-061 |
| OPEN-015 | Automation Run to Response Run bridge | open | CMDR Studio + Govern | Phase 4C/7 | REQ-PROD-062 |
| OPEN-016 | Final wordmark construction and optional symbol | open | Brand Design Lead | Phase 2 review | REQ-BRAND-008 |

## Phase 4B.2A disposition
- **OPEN-007 remains open:** a Studio Human Gate is not automatically a Govern Decision.
- **OPEN-008 remains open:** platform support is not claimed; unsupported states are explicit.
- **OPEN-013 remains open:** no default class-2 step-up or Govern policy is invented.
- **OPEN-015 remains open:** Automation Run, Endpoint Operation and Response Run remain distinct.
- **OPEN-005 remains open:** memory acquisition is specified functionally, but forensic analysis engines belong to Phase 4B.2B.
- **OPEN-011 and OPEN-012 remain open:** mobile and cloud analysis scope are not started.
- OPEN-014 remains open and is only consumed for Attachment/custody relations.
- No new OPEN decision is created and none is closed.

## Resolved history
`OPEN-009` remains the only historically resolved decision: capability delivery-classification authority.

## Maintenance
IDs are never reused. `defined` plus `planned` proves specification only, not implementation or support.
