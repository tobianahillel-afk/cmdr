---
id: unresolved-decisions
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-06
source-of-truth: canonical
---
# Unresolved decisions

The programme retains **18 open decisions**. OPEN-009 remains the only historically resolved item. Phase 4B.3B.2 closes none.

| Relevant decision | Application | Disposition |
|---|---|---|
| OPEN-005 | forensic engine selection | remains open; not reused for Intelligence |
| OPEN-008 | platform/source availability and support | remains open |
| OPEN-011 | Mobile Forensics scope | remains open and not implemented |
| OPEN-012 | Cloud Analysis scope | remains open and not implemented |
| OPEN-013 | reversible class-2 mutations and authority | remains open |
| OPEN-014 | Artifact/Attachment/dataset/material relations and retention | remains open |
| OPEN-015 | Tool Call/Automation Run and cross-product provenance | remains open |
| OPEN-017 | Detection runtime/language/portability | remains open and Detection-only |
| OPEN-018 | Threat Intelligence ontology/interoperability/exchange | remains open |
| **OPEN-019** | dissemination, releasability, sharing and consumer access | **created open; no option selected** |

## OPEN-018 — Threat intelligence ontology, interoperability and exchange strategy
**Status:** open. No standard, protocol, message format, schema, provider, graph database or external representation is selected. Candidate options remain CMDR canonical ontology with adapters; external-standard alignment; hybrid representations; native source models with projection; progressive combination.

## OPEN-019 — Intelligence dissemination, releasability, sharing and consumer access policy
**Status:** open  
**Owner:** Product Architecture  
**Consumers:** Threat Intelligence, Command, Detection Engineering, Platform Settings, Govern, Shared Reporting/Export/Notifications, Security, Objects, Permissions, Screens, Quality and Technique.

### Decision question
Which policy determines internal audiences, markings, releasability, tenant/environment boundaries, client delivery, external sharing, publication approval, separation of duties, access suspension and recall/revocation while preserving source restrictions and provenance?

### Options retained without selection
1. Internal dissemination strictly role-based.
2. Dissemination based on markings and releasability.
3. Tenant-isolated dissemination with controlled cross-tenant sharing.
4. Client-specific delivery policies.
5. External sharing only through Govern.
6. Hybrid policy by classification, source restriction and audience.

### Constraints
- marking is not permission;
- releasable is not shared;
- internal publication is not external or public publication;
- product access grants no raw-source access;
- cross-tenant/client/external sharing requires explicit authority and destination ownership;
- no actual sharing, publication, protocol or destination is configured by this phase;
- corrections, withdrawals and supersession preserve history and access traces.

## Phase 4B closure consequence
OPEN-011 and OPEN-012 remain open and are not explicitly deferred outside Phase 4B by an approved roadmap decision. Therefore Phase 4B.3B.2, 4B.3B and 4B.3 may close, but Phase 4B remains PARTIAL until Cloud/Mobile scope is formally deferred or completed.
