---
id: requirements-traceability-matrix
domain: 00-governance
status: draft
owner: QA and Traceability Lead
updated: 2026-08-03
source-of-truth: source-material
---
# Requirements Traceability Matrix

## Coverage summary

| State | Count |
|---|---:|
| conform | 0 |
| partial | 68 |
| absent | 47 |
| contradictory | 7 |
| total | 122 |

## Requirement catalog

| Requirement ID | Requirement | Type | Expected canonical area | Current state | Priority |
|---|---|---|---|---|---|
| REQ-PROD-001 | CMDR complete operational platform | decision | `01-product-vision/` | partial | P0 |
| REQ-PROD-002 | Evidence first | constraint | `01-product-vision/` | partial | P0 |
| REQ-PROD-003 | Human accountable | constraint | `01-product-vision/` | partial | P0 |
| REQ-PROD-004 | Actions safely governed | constraint | `01-product-vision/` | partial | P0 |
| REQ-PROD-005 | Every conclusion traceable | constraint | `01-product-vision/` | partial | P0 |
| REQ-PROD-006 | No duplicated object | constraint | `00-governance/` | contradictory | P0 |
| REQ-PROD-007 | Progressive disclosure | preference | `04-experience-architecture/` | partial | P0 |
| REQ-PROD-008 | Context preserved across products | constraint | `04-experience-architecture/` | partial | P0 |
| REQ-PROD-009 | One workflow one owner | constraint | `00-governance/` | partial | P0 |
| REQ-PROD-010 | Essential workflows work without AI | constraint | product modules | absent | P0 |
| REQ-PROD-011 | Deterministic engines remain first-class | constraint | `07-investigate/` | partial | P0 |
| REQ-PROD-012 | Planned capability not implemented | constraint | capability inventory | contradictory | P0 |
| REQ-PROD-013 | Command boundary | decision | `06-command/` | partial | P0 |
| REQ-PROD-014 | Investigate boundary | decision | `07-investigate/` | partial | P0 |
| REQ-PROD-015 | Govern boundary | decision | `08-govern/` | partial | P0 |
| REQ-PROD-016 | Studio boundary | decision | `09-cmdr-studio/` | partial | P0 |
| REQ-PROD-017 | Settings boundary | decision | `10-platform-settings/` | partial | P0 |
| REQ-PROD-018 | Endpoint Agent distinct EDR target | decision | `11-endpoint-agent/` | contradictory | P0 |
| REQ-PROD-019 | Capability delivery classification | constraint | capability inventory | absent | P0 |
| REQ-PROD-020 | Deterministic Detection Engineering | constraint | `07-investigate/` | partial | P0 |
| REQ-PROD-021 | Incident Commander role | decision | personas | absent | P1 |
| REQ-PROD-022 | SOC Analyst L1 role | decision | personas | absent | P1 |
| REQ-PROD-023 | SOC Analyst L2 role | decision | personas | absent | P1 |
| REQ-PROD-024 | Senior Analyst role | decision | personas | absent | P1 |
| REQ-PROD-025 | Threat Hunter role | decision | personas | absent | P1 |
| REQ-PROD-026 | Forensic Analyst role | decision | personas | absent | P1 |
| REQ-PROD-027 | Malware Analyst role | decision | personas | absent | P1 |
| REQ-PROD-028 | Reverse Engineer role | decision | personas | absent | P1 |
| REQ-PROD-029 | Detection Engineer role | decision | personas | absent | P1 |
| REQ-PROD-030 | Response Operator role | decision | personas | absent | P1 |
| REQ-PROD-031 | Technical Approver role | decision | personas | absent | P1 |
| REQ-PROD-032 | Business Owner role | decision | personas | absent | P1 |
| REQ-PROD-033 | Customer Success Manager role | decision | personas | absent | P1 |
| REQ-PROD-034 | Auditor role | decision | personas | absent | P1 |
| REQ-PROD-035 | Platform Administrator role | decision | personas | absent | P1 |
| REQ-PROD-036 | Automation Designer role | decision | personas | absent | P1 |
| REQ-PROD-037 | Classify search | constraint | capability inventory | absent | P1 |
| REQ-PROD-038 | Classify investigation | constraint | capability inventory | absent | P1 |
| REQ-PROD-039 | Classify collection | constraint | capability inventory | absent | P1 |
| REQ-PROD-040 | Classify Live Response | constraint | capability inventory | absent | P1 |
| REQ-PROD-041 | Classify endpoint telemetry | constraint | capability inventory | absent | P1 |
| REQ-PROD-042 | Classify endpoint detection | constraint | capability inventory | absent | P1 |
| REQ-PROD-043 | Classify orchestration | constraint | capability inventory | absent | P1 |
| REQ-PROD-044 | Classify governance | constraint | capability inventory | absent | P1 |
| REQ-PROD-045 | Classify reporting | constraint | capability inventory | absent | P1 |
| REQ-PROD-046 | Classify audit | constraint | capability inventory | absent | P1 |
| REQ-PROD-047 | Classify automation | constraint | capability inventory | absent | P1 |
| REQ-PROD-048 | Investigate palette decision | open-question | unresolved decisions | absent | P1 |
| REQ-PROD-049 | Govern palette decision | open-question | unresolved decisions | absent | P1 |
| REQ-PROD-050 | Studio palette decision | open-question | unresolved decisions | absent | P1 |
| REQ-PROD-051 | Typography licensing | open-question | unresolved decisions | absent | P1 |
| REQ-PROD-052 | Initial forensic engines | open-question | unresolved decisions | absent | P0 |
| REQ-PROD-053 | Customer delivery applicability | open-question | unresolved decisions | absent | P1 |
| REQ-PROD-054 | Human Gate and Govern relation | open-question | unresolved decisions | absent | P0 |
| REQ-PROD-055 | Endpoint platform support | open-question | unresolved decisions | absent | P0 |
| REQ-PROD-056 | Capability review authority | open-question | unresolved decisions | absent | P0 |
| REQ-PROD-057 | Role density preference | open-question | unresolved decisions | absent | P1 |
| REQ-PROD-058 | Mobile forensic scope | open-question | unresolved decisions | absent | P1 |
| REQ-PROD-059 | Cloud analysis scope | open-question | unresolved decisions | absent | P1 |
| REQ-PROD-060 | Action Class 2 defaults | open-question | unresolved decisions | absent | P0 |
| REQ-PROD-061 | Artifact versus Attachment | open-question | unresolved decisions | absent | P0 |
| REQ-PROD-062 | Automation Run versus Response Run | open-question | unresolved decisions | absent | P0 |
| REQ-AI-001 | AI augments but is not required | decision | AI constraints | partial | P0 |
| REQ-AI-002 | Studio owns agentic capabilities | decision | Studio | partial | P0 |
| REQ-AI-003 | AI owns no canonical object or decision | constraint | AI constraints | partial | P0 |
| REQ-AI-004 | AI cannot bypass permission or gate | constraint | AI constraints | partial | P0 |
| REQ-AI-005 | AI cannot self-approve risk | constraint | AI constraints | partial | P0 |
| REQ-AI-006 | Evidence changes remain traceable | constraint | AI constraints | partial | P0 |
| REQ-AI-007 | Automation interruptible and resumable | constraint | Studio | partial | P0 |
| REQ-AI-008 | External data is not instruction | constraint | Studio security | absent | P0 |
| REQ-AI-009 | Tool result not automatically trusted | constraint | Studio security | absent | P0 |
| REQ-AI-010 | AI provenance and uncertainty | constraint | AI UX | partial | P0 |
| REQ-AI-011 | Model provider optional | constraint | Settings and Studio | partial | P0 |
| REQ-OBJ-001 | Incident owner Command | decision | Incident | partial | P0 |
| REQ-OBJ-002 | Case owner Investigate | decision | Case | partial | P0 |
| REQ-OBJ-003 | Evidence owner Investigate | decision | Evidence | partial | P0 |
| REQ-OBJ-004 | Finding owner Investigate | decision | Finding | partial | P0 |
| REQ-OBJ-005 | Decision owner Govern | decision | Decision | partial | P0 |
| REQ-OBJ-006 | Approval owner Govern | decision | Approval | partial | P0 |
| REQ-OBJ-007 | Response Run owner Govern | decision | Response Run | partial | P0 |
| REQ-OBJ-008 | Fleet owner Settings | decision | Endpoint Agent Fleet | partial | P0 |
| REQ-OBJ-009 | Agentic objects owner Studio | decision | Studio objects | contradictory | P0 |
| REQ-OBJ-010 | Reporting Engine owner Shared | decision | Reporting Engine | partial | P0 |
| REQ-OBJ-011 | Generic Saved Views owner Shared | decision | Saved Views | contradictory | P0 |
| REQ-OBJ-012 | Work Queue Saved Views owner Command | decision | Work Queue Saved Views | contradictory | P0 |
| REQ-UX-001 | Page workspace view mode filter distinction | decision | experience architecture | partial | P0 |
| REQ-UX-002 | Single canonical Inspector | decision | design system | partial | P0 |
| REQ-UX-003 | Eight canonical shells | decision | design system layouts | partial | P0 |
| REQ-UX-004 | Workbench panel limits | constraint | workbench shell | absent | P0 |
| REQ-UX-005 | Maximum six visible technical tabs | constraint | workbench shell | absent | P0 |
| REQ-UX-006 | Cross-product context preservation | constraint | navigation | partial | P0 |
| REQ-UX-007 | Back restores real state | constraint | navigation | absent | P0 |
| REQ-UX-008 | Work Queue views not pages | decision | Command | contradictory | P0 |
| REQ-UX-009 | Mission Control views share workspace | decision | Command | partial | P0 |
| REQ-UX-010 | Screens have substantive sections | constraint | screen specs | absent | P0 |
| REQ-BRAND-001 | Operational Editorial Modernism | decision | brand | partial | P1 |
| REQ-BRAND-002 | Forbidden visual directions | constraint | brand | partial | P1 |
| REQ-BRAND-003 | Exact CMDR palette | decision | brand | partial | P1 |
| REQ-BRAND-004 | Exact Command palette | decision | brand | partial | P1 |
| REQ-BRAND-005 | Moss plus Ember signature | decision | brand | partial | P1 |
| REQ-BRAND-006 | Other palettes remain proposals | constraint | brand | partial | P1 |
| REQ-BRAND-007 | Typography direction | preference | brand | partial | P1 |
| REQ-BRAND-008 | Related product identities | constraint | brand | partial | P1 |
| REQ-SEC-001 | Action classes 0–4 | decision | trust model | absent | P0 |
| REQ-SEC-002 | High classes require governance | constraint | trust model | partial | P0 |
| REQ-SEC-003 | UI exposes class impact scope rollback approvals audit | constraint | dangerous-action pattern | absent | P0 |
| REQ-SEC-004 | Strong identity secure communication signed command anti-replay | constraint | Endpoint Agent | partial | P0 |
| REQ-SEC-005 | Offline safety tamper rollback verification least privilege | constraint | Endpoint Agent | partial | P0 |
| REQ-JRN-001 | Endpoint Alert to Result | decision | journeys | partial | P1 |
| REQ-JRN-002 | Malware Analysis to Detection Rule | decision | journeys | absent | P1 |
| REQ-JRN-003 | Incident to Containment | decision | journeys | partial | P1 |
| REQ-JRN-004 | Forensic Acquisition to Finding | decision | journeys | partial | P1 |
| REQ-JRN-005 | Detection Rule Development | decision | journeys | absent | P1 |
| REQ-JRN-006 | Handover | decision | journeys | partial | P1 |
| REQ-JRN-007 | Customer Report | decision | journeys | partial | P1 |
| REQ-JRN-008 | Agentic Investigation | decision | journeys | absent | P1 |
| REQ-INV-001 | Classify forensics | constraint | capability inventory | absent | P1 |
| REQ-INV-002 | Classify static analysis | constraint | capability inventory | absent | P1 |
| REQ-INV-003 | Classify reverse engineering | constraint | capability inventory | absent | P1 |
| REQ-INV-004 | Classify debugger | constraint | capability inventory | absent | P1 |
| REQ-INV-005 | Classify sandbox | constraint | capability inventory | absent | P1 |
| REQ-INV-006 | Classify Detection Engineering | constraint | capability inventory | absent | P1 |

## Contradictions requiring correction

| Requirement | Current conflicting area | Required correction |
|---|---|---|
| REQ-PROD-006 | repeated generic definitions across active files | retain one canonical definition and local usage only |
| REQ-PROD-012 | complete/native wording for unclassified capabilities | label delivery status explicitly |
| REQ-PROD-018 | Endpoint Agent described as complete although contracts are placeholders | describe it as planned complete native target |
| REQ-OBJ-009 | Tool Call and Automation Run absent | add canonical Studio objects in later phases |
| REQ-OBJ-011/012 | Work Queue language copied into unrelated screens | reference Shared Saved Views outside Command |
| REQ-UX-008 | Work Queue variants represented by multiple screens | consolidate into one workspace with saved views |
| source precedence | current policy places ADR before sponsor source material | update governance policy in Phase 1 |

## Default remediation by state

- partial: preserve useful content and complete behavior in its owning phase;
- absent: create canonical coverage in its owning phase;
- contradictory: remove the competing active definition before validation.