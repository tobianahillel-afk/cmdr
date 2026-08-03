---
id: source-user-role-and-journey-inventory
domain: 00-governance
status: draft
owner: Product Management
updated: 2026-08-03
source-of-truth: source-material
---
# User Role and Journey Inventory

## Required roles

| Requirement | Role |
|---|---|
| REQ-PROD-021 | Incident Commander |
| REQ-PROD-022 | SOC Analyst L1 |
| REQ-PROD-023 | SOC Analyst L2 |
| REQ-PROD-024 | Senior Analyst |
| REQ-PROD-025 | Threat Hunter |
| REQ-PROD-026 | Forensic Analyst |
| REQ-PROD-027 | Malware Analyst |
| REQ-PROD-028 | Reverse Engineer |
| REQ-PROD-029 | Detection Engineer |
| REQ-PROD-030 | Response Operator |
| REQ-PROD-031 | Technical Approver |
| REQ-PROD-032 | Business Owner |
| REQ-PROD-033 | Customer Success Manager |
| REQ-PROD-034 | Auditor |
| REQ-PROD-035 | Platform Administrator |
| REQ-PROD-036 | Automation Designer |

Each role specification must define goals, responsibilities, usage frequency, technical level, priority information, information hidden by default, allowed and forbidden actions, principal workspaces, journeys, accessibility needs and density preferences.

## Required end-to-end journeys

| Requirement | Journey |
|---|---|
| REQ-JRN-001 | Endpoint Alert to Result |
| REQ-JRN-002 | Malware Analysis to Detection Rule |
| REQ-JRN-003 | Incident to Containment |
| REQ-JRN-004 | Forensic Acquisition to Finding |
| REQ-JRN-005 | Detection Rule Development |
| REQ-JRN-006 | Handover |
| REQ-JRN-007 | Customer Report |
| REQ-JRN-008 | Agentic Investigation |

Every journey must define trigger, user, preconditions, steps, screens, buttons, objects, permissions, errors, states, transitions, context preservation, result and rollback.