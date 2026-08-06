---
id: unresolved-decisions
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-06
source-of-truth: canonical
---
# Unresolved decisions

The programme retains **18 open decisions**. OPEN-009 remains the only historically resolved item. Phase 4B.4A closes none.

| Relevant decision | Application | Disposition |
|---|---|---|
| OPEN-005 | forensic engine selection | remains open; not reused for Cloud Analysis |
| OPEN-008 | platform/source availability and support | remains open |
| OPEN-011 | Mobile Forensics scope | remains open and not implemented |
| **OPEN-012** | Cloud Analysis scope | **remains open; functional Cloud capabilities defined, provider and delivery scope unresolved** |
| OPEN-013 | reversible class-2 mutations and authority | remains open |
| OPEN-014 | Artifact/Attachment/dataset/material relations and retention | remains open |
| OPEN-015 | Tool Call/Automation Run and cross-product provenance | remains open |
| OPEN-017 | Detection runtime/language/portability | remains open and Detection-only |
| OPEN-018 | Threat Intelligence ontology/interoperability/exchange | remains open |
| OPEN-019 | dissemination, releasability, sharing and consumer access | remains open; no option selected |

## OPEN-012 — Cloud Analysis scope and delivery strategy
**Status:** open  
**Owner:** Product Architecture  
**Consumers:** Cloud Analysis, Investigate, Platform Settings, Endpoint Agent, Command, Detection Engineering, Threat Intelligence, Govern, Studio, Shared, Security, Objects, Permissions, Screens, Quality, Technique and Roadmap.

### Decision question
Which Cloud providers, account/tenant models, services, audit models, inventory strategies, cloud-native workloads, SaaS boundaries, evidence rules and cross-tenant limits are supported first, and what provider-neutral versus provider-specific contracts are allowed, without collapsing configured, accessible and complete coverage into one state?

### Questions retained
- Which providers are prioritized, and which remain explicitly unsupported or experimental?
- Which organization, tenant, account, subscription, project, folder and region units are supported?
- Which control-plane, data-plane, identity, audit, inventory and configuration sources are required or optional?
- How are multi-cloud and cross-account/cross-tenant relations represented and authorized?
- Which compute, containers, orchestration, Kubernetes-like, serverless and managed execution services are covered?
- Which SaaS analysis belongs to Cloud Analysis versus a future dedicated capability?
- Which Cloud evidence, custody, retention, legal and reproducibility rules apply?
- Which inventory snapshots and activity models are sufficient to claim partial, complete, current or stale coverage?
- Which provider-specific schemas, parsers, query languages, APIs and connectors may be introduced in later implementation phases?

### Options retained without selection
1. Provider-neutral core with adapters introduced only after object/permission contracts.
2. One prioritized provider, then progressive expansion.
3. Multi-cloud control-plane foundation before workload-specific depth.
4. Service-family slices independent of provider priority.
5. Separate Cloud infrastructure and SaaS delivery tracks.
6. Hybrid sequencing by evidence availability, customer need and risk.

### Constraints
- no provider, connector, API, protocol, query language, schema, scanner, command or implementation is selected by Phase 4B.4A;
- configured provider ≠ accessible provider ≠ complete coverage;
- inventory/configuration observation ≠ certain current state;
- Cloud identity/principal ≠ person;
- permission path ≠ exploit path;
- Cloud anomaly ≠ compromise, Finding or vulnerability;
- Cloud Analysis prepares but does not execute collection, rule deployment, permission change, credential revocation or response;
- cross-tenant analysis requires explicit scope and authority;
- source permissions and sensitive-access audit remain mandatory.

### Phase disposition
Phase 4B.4A defines provider-neutral Cloud investigation capabilities and is PASS documentarily. OPEN-012 remains open because provider/service prioritization and implementation/delivery strategy are not selected.

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
OPEN-011 remains open and Mobile Forensics is not started. Cloud Analysis is now documentarily defined and PASS, but Phase 4B.4 and Phase 4B remain PARTIAL until Mobile Forensics is completed or formally deferred by an approved roadmap decision. OPEN-012 remains open for provider/service/delivery strategy and does not invalidate the provider-neutral functional conformance of 4B.4A.
