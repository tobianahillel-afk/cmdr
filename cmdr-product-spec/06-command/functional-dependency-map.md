---
id: command-functional-dependency-map
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-013
---
# Functional dependency map — Command

| Capability | Dependency | Type | Failure behavior |
|---|---|---|---|
| CAP-CMD-001 | CAP-CMD-003 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-001 | CAP-CMD-005 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-001 | CAP-CMD-006 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-001 | Shared Global Search | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-001 | Object Linking Service | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-001 | Timeline Engine | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-001 | Business Service Catalog | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-002 | CAP-CMD-104 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-002 | CAP-CMD-105 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-002 | CAP-CMD-204 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-002 | CAP-CMD-205 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-002 | Metrics Engine | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-002 | Audit hooks | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-003 | Timeline Engine | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-003 | Object Linking Service | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-003 | Trace | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-003 | Export Engine | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-003 | CAP-CMD-006 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-004 | CAP-CMD-001 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-004 | CAP-CMD-005 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-004 | CAP-CMD-006 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-004 | Notification Center | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-004 | Collaboration Service | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-004 | Object Linking Service | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-005 | CAP-CMD-110 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-005 | CAP-CMD-107 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-005 | Object Linking Service | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-005 | Notification Center | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-006 | CAP-CMD-001 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-006 | CAP-CMD-003 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-006 | CAP-CMD-106 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-006 | Object Linking Service | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-006 | Timeline Engine | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-006 | Reporting Engine | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-101 | CAP-CMD-102 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-101 | CAP-CMD-104 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-101 | CAP-CMD-105 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-101 | CAP-CMD-108 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-101 | CAP-CMD-109 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-101 | Saved Views | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-101 | Data Table | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-101 | Filters | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-101 | Inspector | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-102 | Identity projection | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-102 | Notification Center | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-102 | Collaboration Service | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-102 | CAP-CMD-103 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-102 | CAP-CMD-110 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-103 | CAP-CMD-102 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-103 | Ownership Register | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-103 | Identity projections | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-103 | Object Linking Service | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-103 | Audit hooks | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-104 | CAP-CMD-002 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-104 | CAP-CMD-105 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-104 | CAP-CMD-204 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-104 | CAP-CMD-205 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-104 | Metrics Engine | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-105 | CAP-CMD-101 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-105 | CAP-CMD-102 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-105 | CAP-CMD-110 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-105 | Notification Center | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-105 | Metrics Engine | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-105 | ADR-0008 / deployment SLA source / CAP-CMD-401 | functional | source-backed contractual SLA context; no silent fallback |
| CAP-CMD-106 | CAP-CMD-102 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-106 | CAP-CMD-104 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-106 | CAP-CMD-107 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-106 | CAP-CMD-110 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-106 | Object Linking Service | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-106 | Timeline Engine | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-106 | Investigate Case lifecycle | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-106 | Govern Action Request lifecycle | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-107 | CAP-CMD-102 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-107 | CAP-CMD-103 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-107 | CAP-CMD-005 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-107 | Object Linking Service | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-107 | Notification Center | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-107 | Task Inbox projection | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-108 | CAP-CMD-101 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-108 | CAP-CMD-102 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-108 | CAP-CMD-104 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-108 | Background Jobs | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-108 | Export Engine | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-108 | Audit hooks | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-109 | Data Quality Service | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-109 | Background Jobs | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-109 | Notification Center | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-109 | CAP-CMD-107 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-109 | CAP-CMD-110 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-110 | CAP-CMD-005 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-110 | CAP-CMD-106 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-110 | Notification Center | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-110 | Collaboration Service | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-110 | Object Linking Service | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-110 | Context preservation | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-201 | Business Service Catalog | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-201 | Object Linking Service | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-201 | CAP-CMD-204 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-201 | CAP-CMD-106 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-201 | CAP-CMD-107 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-202 | Business Service Catalog | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-202 | Object Linking Service | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-202 | Data Quality Service | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-202 | CAP-CMD-201 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-202 | CAP-CMD-203 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-203 | CAP-CMD-201 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-203 | CAP-CMD-301 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-203 | CAP-CMD-303 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-203 | Metrics Engine | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-203 | Business Service Catalog | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-203 | Data Quality Service | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-204 | CAP-CMD-201 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-204 | CAP-CMD-205 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-204 | CAP-CMD-106 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-204 | Business Service Catalog | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-204 | Object Linking Service | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-205 | CAP-CMD-002 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-205 | CAP-CMD-104 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-205 | CAP-CMD-105 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-205 | CAP-CMD-201 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-205 | CAP-CMD-202 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-205 | CAP-CMD-203 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-205 | CAP-CMD-204 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-205 | Metrics Engine | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-301 | CAP-CMD-203 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-301 | CAP-CMD-303 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-301 | CAP-CMD-304 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-301 | CAP-CMD-305 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-301 | Studio Assurance | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-301 | Platform Health | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-301 | Metrics Engine | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-302 | CAP-CMD-301 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-302 | CAP-CMD-303 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-302 | CAP-CMD-304 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-302 | Collaboration Service | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-302 | Notification Center | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-302 | Reporting Engine | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-303 | CAP-CMD-107 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-303 | CAP-CMD-301 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-303 | CAP-CMD-302 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-303 | CAP-CMD-203 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-303 | Object Linking Service | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-304 | CAP-CMD-301 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-304 | CAP-CMD-302 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-304 | CAP-CMD-305 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-304 | Object Linking Service | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-304 | Versioning service | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-304 | Notification Center | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-305 | Capability Register | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-305 | CAP-CMD-301 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-305 | CAP-CMD-303 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-305 | Platform Health | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-305 | Studio Assurance | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-305 | Metrics Engine | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-401 | ADR-0008 / deployment customer-contract source | functional | source-backed deployment context; no silent fallback |
| CAP-CMD-401 | Security Tenant isolation / Authorized Tenant Set | functional | authorized read-only aggregation; explicit single-Tenant context before mutation/response |
| CAP-CMD-401 | Govern Decision Authority | functional | required after explicit single-Tenant selection where response is applicable |
| CAP-CMD-401 | Shared Global Search | functional | initial scope is one selected Tenant; no silent widening |
| CAP-CMD-401 | Reporting Engine | functional | initial report scope is single-Tenant; no silent widening |
| CAP-CMD-401 | Metrics Engine | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-401 | Export Engine | functional | single-Tenant export must not widen visibility |
| CAP-CMD-401 | Business Service Catalog | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-401 | CAP-CMD-105 | functional | visible Partial/blocked; no silent fallback |
| CAP-CMD-401 | CAP-CMD-303 | functional | visible Partial/blocked; no silent fallback |

## Rules

- dependency does not imply ownership;
- every cross-product projection carries source, freshness and permission;
- failure preserves valid Command data and exposes the impact;
- technical endpoint, protocol and SLO remain out of Phase 4A.
