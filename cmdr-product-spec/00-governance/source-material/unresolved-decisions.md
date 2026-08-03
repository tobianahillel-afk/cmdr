---
id: source-unresolved-decisions
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: source-material
---
# Unresolved Decisions

Generic `À compléter` placeholders are not accepted. Open decisions use this structured register.

| Decision | Requirement | Owner | Target phase | Blocking | Question |
|---|---|---|---|---|---|
| OPEN-001 | REQ-PROD-048 | Brand Design Lead | Phase 2 | false | Select Investigate palette direction. |
| OPEN-002 | REQ-PROD-049 | Brand Design Lead | Phase 2 | false | Select Govern palette direction. |
| OPEN-003 | REQ-PROD-050 | Brand Design Lead | Phase 2 | false | Select Studio palette direction. |
| OPEN-004 | REQ-PROD-051 | Design System Lead | Phase 2 | false | Confirm final typography stack and licensing. |
| OPEN-005 | REQ-PROD-052 | Investigate Product Lead | Phase 4 | true | Classify initial forensic engines. |
| OPEN-006 | REQ-PROD-053 | Command Product Lead | Phase 4 | false | Define customer-delivery deployment applicability. |
| OPEN-007 | REQ-PROD-054 | Security Architecture | Phase 4 | true | Define when a Studio Human Gate requires Govern. |
| OPEN-008 | REQ-PROD-055 | Endpoint Agent Product Lead | Phase 4 | true | Define visible platform support. |
| OPEN-009 | REQ-PROD-056 | Product Architecture | Phase 1 | true | Define capability-status review authority. |
| OPEN-010 | REQ-PROD-057 | UX Architecture | Phase 3 | false | Define role density preferences. |
| OPEN-011 | REQ-PROD-058 | Investigate Product Lead | Phase 4 | false | Define mobile forensic scope. |
| OPEN-012 | REQ-PROD-059 | Investigate Product Lead | Phase 4 | false | Define cloud analysis scope. |
| OPEN-013 | REQ-PROD-060 | Security Architecture | Phase 4 | true | Define default governance for Action Class 2. |
| OPEN-014 | REQ-PROD-061 | Investigate Product Lead | Phase 7 | true | Define Artifact versus Attachment boundary. |
| OPEN-015 | REQ-PROD-062 | CMDR Studio Product Lead | Phase 4 | true | Define Automation Run to Response Run bridge. |

## Required record shape

```yaml
decision_id:
status: open
owner:
target_phase:
blocking: true|false
depends_on:
question:
options:
required_evidence:
affected_files:
```

Each open decision must be expanded with options, required evidence and affected files before its target phase begins.