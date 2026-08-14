---
id: unresolved-decisions
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-14
source-of-truth: canonical
---
# Unresolved decisions

The programme retains **17 open decisions**. `OPEN-006` is resolved on 2026-08-14 by explicit project-owner approval recorded in `ADR-0008`; `OPEN-009` remains the historically resolved item from the earlier programme state.

## Resolved decision record — OPEN-006 — Customers and Delivery deployment model

**Status:** resolved  
**Resolved:** 2026-08-14  
**Owner:** Product Architecture  
**Approval reference:** Hillel Tobiana — explicit project-owner approval in ChatGPT conversation  
**Canonical ADR:** `../adr/ADR-0008-customers-mssp-delivery-deployment-and-cross-tenant-architecture.md`

### Approved authority input — verbatim

**APPROVED_D1_DEPLOYMENT_MODES:**  
CMDR supports Internal, Enterprise multi-tenant, and MSP/MSSP deployment modes.  
MSSP capabilities are deployment-dependent.

**APPROVED_D2_CUSTOMER_SEMANTICS:**  
Customer remains an external/deployment/customer/contract projection.  
Customer is NOT a canonical CMDR object and is NOT an alias of Tenant.

**APPROVED_D3_MSSP_TENANT_MODEL:**  
MSSP operates over independent Tenants.  
No parent/child Tenant hierarchy.  
No ManagedTenant, TenantGroup, Portfolio, CustomerTenant, or equivalent canonical object for the MVP.

**APPROVED_D4_CROSS_TENANT_AUTH_MVP:**  
Security resolves an Authorized Tenant Set as a non-canonical authorization projection.  
The MVP allows read-only aggregation across that explicitly authorized Tenant set and explicit Tenant context switching.  
Cross-tenant mutation, administration, response, delegated administration, and automatic export widening are not allowed.  
Existing object read permissions remain subject to server-side RBAC/ABAC and Tenant-set evaluation.  
No generic cross-tenant.manage permission is created.

**APPROVED_D5_SEARCH_REPORTING_SCOPE:**  
Search is single-selected-Tenant initially.  
Report is single-Tenant initially.  
Export is single-Tenant initially and must never widen visibility.  
Multi-tenant Search, Reporting, and Export are deferred.  
Shared retains Search / Reporting / Export ownership.

**APPROVED_D6_RESPONSE_AUTHORITY:**  
MSSP response requires explicit Tenant selection, Security authorization re-evaluation inside that Tenant, then Govern Decision Authority evaluation inside that Tenant.  
No centralized cross-tenant response authority is introduced.

**APPROVED_D7_CAP_CMD_401_DISPOSITION:**  
CAP-CMD-401 keeps the same ID and owner.  
After the architecture is canonically recorded, it may move from delivery_status: proposed to delivery_status: defined.  
Its scope is narrowed to deployment-aware activation, external Customer/engagement projection, authorized read-only multi-tenant overview, Tenant selector/context switching, portfolio-like UI projection without a Portfolio object, single-Tenant Reporting requests, service-delivery Tasks, contractual SLA projection, and source/freshness/audience visibility.

**Explicitly OUT:**  
canonical Customer lifecycle/admin,  
Tenant hierarchy,  
cross-tenant mutation/admin/response,  
delegated administration,  
initial multi-tenant Search/Report/Export,  
CRM,  
billing,  
customer portal,  
contract mutation.

`OPEN-013` remains OPEN.  
`OPEN-019` remains OPEN.

`OPEN-006`: **APPROVED AS RESOLVED BY D1–D7 ABOVE.**

No new capability ID.  
No CAP-SET-014.  
No CAP-CMD-402.  
No new canonical object.  
No new Permission ID.  
No new Screen ID.

The approved decision set above must be used verbatim as the authority input for the architecture-recording run.

## Complete open-decision audit

| Decision | Current subject | Current relevance | Disposition |
|---|---|---|---|
| OPEN-001 | Investigate product accent/palette | presentation | remains open |
| OPEN-002 | Govern product accent/palette | presentation only; not capability authority | remains open |
| OPEN-003 | Studio product accent/palette | presentation | remains open |
| OPEN-004 | final typography stack | presentation only | remains open |
| OPEN-005 | forensic engine selection | upstream evidence tooling | remains open |
| **OPEN-007** | **Human Gate / Govern Decision-Approval relation** | **Human Gate must remain distinct from Approval and Decision** | **remains open** |
| OPEN-008 | platform/source availability and support | source/target/context availability | remains open |
| OPEN-010 | final density by role/activity | final screen density | remains open |
| OPEN-011 | Mobile Forensics scope/delivery | Mobile platform/tool delivery | remains open |
| OPEN-012 | Cloud Analysis scope/delivery | Cloud provider/service delivery | remains open |
| **OPEN-013** | **default governance/authority for reversible class-2 mutations** | **direct for C2 actions, including delivery Tasks where applicable** | **remains open** |
| OPEN-014 | Artifact/Attachment/dataset/material relations/retention | Evidence/context links | remains open |
| **OPEN-015** | **Tool Call/Automation Run and Response Run/cross-product provenance bridge** | **cross-product provenance** | **remains open** |
| OPEN-016 | final wordmark/symbol construction | brand | remains open |
| OPEN-017 | Detection runtime/language/portability | Detection implementation | remains open |
| OPEN-018 | Threat Intelligence ontology/interoperability/exchange | Intelligence interoperability | remains open |
| OPEN-019 | Intelligence dissemination/releasability/sharing/access | external/client sharing and audience policy | remains open; no sharing policy selected |

## Decision discipline after OPEN-006

- `OPEN-006` is resolved only by the exact D1–D7 approval recorded above and ADR-0008.
- `OPEN-013` remains the unresolved default-governance question for reversible Class-2 mutations; the OPEN-006 resolution does not select its default.
- `OPEN-019` remains relevant to external/client-facing sharing and publication; read-only internal MSSP aggregation does not resolve dissemination policy.
- `OPEN-007` continues to keep Studio Human Gate distinct from Govern Approval and Decision.
- `OPEN-015` continues to govern the Automation Run / Response Run provenance bridge.
- no other OPEN is closed or weakened by the Customers/MSSP/Delivery decision.

## Previously highlighted cross-domain decisions

| Relevant decision | Application | Disposition |
|---|---|---|
| OPEN-005 | forensic engine selection | remains open; not reused to select a Mobile acquisition/analysis engine |
| OPEN-008 | platform/source availability and support | remains open |
| **OPEN-011** | Mobile Forensics scope and delivery strategy | **remains open; provider-neutral functional capabilities verified, platform/tool/acquisition delivery unresolved** |
| **OPEN-012** | Cloud Analysis scope | **remains open; provider-neutral functional capabilities verified, provider and delivery scope unresolved** |
| OPEN-013 | reversible class-2 mutations and authority | remains open |
| OPEN-014 | Artifact/Attachment/dataset/material relations and retention | remains open |
| OPEN-015 | Tool Call/Automation Run and cross-product provenance | remains open |
| OPEN-017 | Detection runtime/language/portability | remains open and Detection-only |
| OPEN-018 | Threat Intelligence ontology/interoperability/exchange | remains open |
| OPEN-019 | dissemination, releasability, sharing and consumer access | remains open; no option selected |

## OPEN-011 — Mobile Forensics scope and delivery strategy
**Status:** open  
**Owner:** Product Architecture  
**Consumers:** Mobile Forensics, Investigate, Collection and Live Response, Platform Settings, Endpoint Agent, Case/Evidence, Static/Dynamic/Reverse, Disk/Network/Memory Forensics, Cloud Analysis, Detection Engineering, Threat Intelligence, Govern, Studio, Shared, Security, Objects, Permissions, Screens, Quality, Technique and Roadmap.

### Decision question
Which mobile platforms and versions, phone/tablet classes, managed/unmanaged states, acquisition representations, backup/synchronization sources, device/access states and future tools/methods are supported first, and how are privacy, cross-device, cross-tenant, custody and Evidence limitations enforced without making provider-neutral functional coverage depend on a selected vendor or bypass technique?

### Questions retained
- Which mobile platforms and versions are prioritized, and which remain unsupported or experimental?
- Which phone/tablet/device classes are included and how are managed versus unmanaged devices represented?
- How are SIM/eSIM, local backups, synchronized backups, logical extraction, filesystem extraction and other functional representations supported?
- How are locked/unlocked state, encryption, device state, partial/inaccessible material and platform-version uncertainty represented?
- Which application, application-data, messaging, media, location, sensor, account, token and paired-device sources are supported by future implementation?
- How are wearables represented only as linked sources without creating a complete wearable-forensics product?
- Which MDM/EMM integrations and managed-device relations belong to Platform Settings and Govern rather than Mobile analysis?
- Which cross-device and cross-tenant correlations are allowed, and how is source permission preserved?
- Which future methods and Tools may implement parsing, extraction, recovery and comparison without becoming mandatory architecture?
- Which consent/authorization/source restrictions, privacy classifications, retention, Evidence and custody rules apply?

### Scope retained without selection
- phones and tablets;
- managed and unmanaged devices;
- SIM and eSIM;
- local and synchronized backups;
- logical and filesystem extraction plus other evidence representations;
- locked/unlocked, encrypted, partial, restricted and inaccessible states;
- applications and application data;
- messaging, communications, calls and contacts;
- media, documents and downloads;
- location, movement, sensors and health/activity data where authorized;
- accounts, tokens, keys and sensitive material;
- Wi-Fi, Bluetooth, NFC, paired devices and network metadata;
- wearables only as linked sources;
- MDM/EMM administrative context;
- cross-device and cross-tenant limits;
- deleted/residual/recovered data;
- Evidence/custody/provenance and privacy minimization.

### Constraints
- no mobile platform or version is selected by Phase 4B.4B;
- no acquisition tool, engine, API, protocol, proprietary format or final technical method is selected;
- no unlocking, bypass, rooting, jailbreak, code/password cracking, exploit or offensive extraction is defined;
- Collection owns acquisition requests/jobs/execution/results and collection-time custody;
- Platform Settings owns configured MDM/EMM providers, sources, connectors, credentials, secrets, Fleet, retention and policies;
- Endpoint Agent capability is never assumed for a mobile device;
- Govern owns every real-device action including governed remote collection, profile installation, revocation, isolation, lock and wipe;
- device/account/contact/number/location records do not establish a person or human presence automatically;
- secret presence does not grant reveal, copy, export or use; secret use is prohibited in Investigate;
- cross-device/cross-tenant correlation does not grant access to either source;
- Evidence Candidate and Finding Draft remain distinct from qualified Evidence and confirmed Finding.

### Phase disposition after verification
`CAP-INV-701..719` provide verified provider-neutral Mobile functional documentation. OPEN-011 remains open because platform/version/tool/acquisition delivery strategy is not selected.

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

### Phase disposition after verification
Phase 4B.4A provider-neutral functional documentation is **PASS AFTER POST-PUBLICATION VERIFICATION**. OPEN-012 remains open because provider/service prioritization and implementation/delivery strategy are not selected. The open decision does not invalidate the verified provider-neutral functional scope.

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

## Current capability-programme consequence
Cloud Analysis and Mobile Forensics remain verified PASS. Capability Specification Phase 4B — Investigate remains PASS. `OPEN-006` is resolved by ADR-0008 and is no longer counted among the 17 OPEN decisions. `OPEN-007`, `OPEN-010`, `OPEN-011`, `OPEN-012`, `OPEN-013`, `OPEN-014`, `OPEN-015`, `OPEN-016`, `OPEN-017`, `OPEN-018`, `OPEN-019` and all other current open decisions retain their dispositions. The Customers/MSSP/Delivery architecture creates no new Capability ID, object, Permission ID or Screen ID and does not claim implementation.
