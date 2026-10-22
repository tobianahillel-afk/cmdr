---
id: platform-settings-capability-map
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-14
source-of-truth: canonical
---
# Platform Settings Capability Map

## Tenant, Environment and Administrative Foundations
| Capability | Canonical responsibility | Primary objects | Primary screen | Dependencies |
|---|---|---|---|---|
| CAP-SET-001 | Tenant lifecycle and isolation boundary | Tenant | SET-TEN-001 | Security, Administrative Audit |
| CAP-SET-002 | Environment lifecycle and Tenant scope | Environment, Tenant ref | SET-TEN-001 | Security, Experience Architecture |
| CAP-SET-003 | administrative change validation and provenance | Tenant, Environment | SET-TEN-001, SET-AUD-001 | Security, Administrative Audit |
| CAP-SET-004 | cross-product Tenant/Environment context semantics | Tenant/Environment refs | existing product screens | Experience Architecture, Design System, Security |

## Identity Administration
| Capability | Canonical responsibility | Primary objects | Primary screen | Dependencies |
|---|---|---|---|---|
| CAP-SET-005 | Principal lifecycle and identity-state administration | Principal, Tenant ref | SET-IAM-001 | Security, Administrative Audit |
| CAP-SET-006 | Role lifecycle/constraints and bounded Principal/Role relation | Role, Principal ref, Tenant ref | SET-IAM-001 | Security, ABAC, SoD |
| CAP-SET-007 | Access Review evidence/provenance and disposition/handoff | Principal/Role refs | SET-IAM-001, SET-AUD-001 | Security, Administrative Audit |

## Secrets & Connections
| Capability | Canonical responsibility | Primary objects | Primary screen | Dependencies |
|---|---|---|---|---|
| CAP-SET-008 | Integration lifecycle, local validation and administrative result projection | Integration, Tenant ref | SET-SEC-001, SET-AUD-001 | Security, Administrative Audit |
| CAP-SET-009 | Secret Reference lifecycle and reference administration | Secret Reference, Tenant ref | SET-SEC-001, SET-AUD-001 | Security, Administrative Audit |

## Models & Providers
| Capability | Canonical responsibility | Primary objects | Primary screen | Dependencies |
|---|---|---|---|---|
| CAP-SET-010 | Model Provider lifecycle/configuration, local validation and sourced availability/health projection | Model Provider, Tenant ref; optional sourced refs | SET-MDL-001, SET-HLT-001, SET-AUD-001 | Security, Administrative Audit, Health, Studio consumers |
| CAP-SET-011 | routing/eligibility/fallback configuration and provenance of sourced effective-selection observations | Model Provider configuration, Tenant ref | SET-MDL-001, SET-AUD-001 | Security, data-policy owners, Studio/runtime consumers, Administrative Audit |

## Sources & Parsers
| Capability | Canonical responsibility | Primary objects | Primary screen | Dependencies |
|---|---|---|---|---|
| CAP-SET-012 | Data Source lifecycle/configuration, scope, freshness, safe disable, local validation and sourced health projection | Data Source, Tenant ref; optional sourced Integration/Secret Reference refs | SET-SRC-001, SET-HLT-001, SET-AUD-001 | Security, Health/Event contracts, Administrative Audit, runtime owner via handoff |
| CAP-SET-013 | Parser lifecycle/versioned transformation contract, fixtures, errors/quality, local validation and source-backed administrative rollback | Parser, Tenant ref | SET-SRC-001, SET-AUD-001 | Security, Parser/Event contracts, Administrative Audit, parser runtime owner via handoff |

No new Screen ID, Permission ID or canonical object is introduced.

Key Sources & Parsers distinctions: Data Source != Parser != Integration; the canonical objects define no direct Data Source→Parser relation; no assignment/compatibility/routing/selection/fallback/precedence is invented. Health/test/error/quality facts derived from runtime are sourced projections. Settings administrative configuration does not become acquisition, collection, ingestion, connector/probe, parser, normalization, stream-processing, schema-registry or storage execution.

Current Settings functional BUILD structure: **13 capabilities / 351 numbered sections / 78 mandatory tables**. Sources & Parsers contributes **8 meaningful GWT**. Requirements remain **122 = 99/20/3/0**. OPEN remains **18**. `CAP-SET-014+` remains unallocated and unreserved. Platform Settings Capability Specification remains **PARTIAL**.
