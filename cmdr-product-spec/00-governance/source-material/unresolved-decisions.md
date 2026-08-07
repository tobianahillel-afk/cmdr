---
id: unresolved-decisions
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-07
source-of-truth: canonical
---
# Unresolved decisions

The programme retains **18 open decisions**. OPEN-009 remains the only historically resolved item. Phase 4B.4B closes none.

| Relevant decision | Application | Disposition |
|---|---|---|
| OPEN-005 | forensic engine selection | remains open; not reused to select a Mobile acquisition/analysis engine |
| OPEN-008 | platform/source availability and support | remains open |
| **OPEN-011** | Mobile Forensics scope and delivery strategy | **remains open; provider-neutral functional capabilities defined, platform/tool/acquisition delivery unresolved** |
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

### Phase disposition before final verification
`CAP-INV-701..719` provide provider-neutral Mobile functional documentation. The phase remains **PENDING POST-PUBLICATION VERIFICATION** until the fifth functional commit and remote checks complete. OPEN-011 remains open because platform/version/tool/acquisition delivery strategy is not selected.

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

## Phase 4B closure consequence before Mobile publication verification
Cloud Analysis remains verified PASS. Mobile functional content is complete but awaits the fifth functional commit and remote publication verification. Phase 4B.4 and Phase 4B therefore remain PARTIAL until those checks pass. OPEN-011 and OPEN-012 remain open for delivery strategy and do not by themselves invalidate provider-neutral functional coverage once their corresponding subphase is verified.
