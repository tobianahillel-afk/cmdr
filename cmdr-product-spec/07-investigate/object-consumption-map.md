---
id: investigate-object-consumption-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-07
source-of-truth: canonical
requirements: [REQ-PROD-006, REQ-PROD-014, REQ-PROD-061, REQ-PROD-062, REQ-INV-001, REQ-INV-006]
open_decisions: [OPEN-008, OPEN-011, OPEN-012, OPEN-013, OPEN-014, OPEN-015, OPEN-017, OPEN-018, OPEN-019]
---
# Object consumption map — Investigate through Phase 4B.4B closure content

| Objet ou concept | Owner actuel | Usage local | Opérations locales | Lacune | Phase propriétaire |
|---|---|---|---|---|---|
| Case / Finding / Evidence / Case Hypothesis / Hunt / Artifact | Investigate | origins, evidence and return context | read/link | final relations | Objects |
| Incident / runtime Detection / Signal / Alert | Command | source/runtime/feedback projections | read/link only | bridge contracts | Command/Objects |
| Detection Content / Detection Hypothesis / Coverage / Gap | Detection Engineering | operationalization and feedback | prepare/read/link in owner module | runtime/language OPEN-017 | Objects/Technique |
| Action Request / Decision / Approval / Response Run / Result | Govern | authority/execution projections | prepare/read/link only | cross-run OPEN-015 | Govern/Objects |
| Data Source / Provider / Connector / Environment / Destination / Policy | Platform Settings | source, access, target and destination projections | read/select/request only | OPEN-008/012/018/019 | Settings/Technique |
| Tool / Tool Call / Workflow / Automation Run / Agent / Human Gate | Studio | attributed automation and monitoring | select/invoke/read/link | OPEN-015 | Studio/Objects |
| Entity / Graph / Timeline / Search / Linking / Versioning / Report / Export / Notification / Trace / Activity / Recovery | Shared | generic mechanisms | consume/emit domain semantics | final contracts | Shared/Technique |
| Intelligence Requirement / Collection Priority / Knowledge Project | Investigate concepts | needs, scope, priority and durable workspace | create/update/pause/close/reopen/supersede | canonical schemas absent | Objects |
| Intelligence Source / Material / Normalized Record | Settings/source owner + Investigate concepts | authorized source/material context | reference/read/assess/version | Artifact/Material relation OPEN-014 | Objects/Settings |
| Observable / Indicator / Threat Entity candidates | Investigate concepts | candidate knowledge | create/update/review/withdraw/supersede | ontology OPEN-018 | Objects |
| Malware / Tool / Capability / Infrastructure knowledge | Investigate concepts | sourced technical knowledge | create/update/version/review | canonical models absent | Objects |
| Activity Cluster / Campaign / Intrusion Set candidates | Investigate concepts | candidate grouping | create/compare/merge proposal/separate | attribution model absent | Objects |
| TTP Mapping / Sighting / Intelligence Relationship | Investigate using Shared | behavior, occurrence and sourced relation | create/review/dispute/supersede | taxonomy/graph absent | Objects/Shared |
| Confidence / Contradiction / Lifecycle / Analysis Handoff Package | Investigate concepts | assessment, lifecycle and handoff | create/review/version/withdraw | no calibrated score/package schema | Objects |
| Intelligence Analysis Session | Investigate concept | bounded analysis context | create/update/pause/close/reopen/archive/supersede | canonical object absent | Objects |
| Intelligence Analysis Question | Investigate concept | analytic question distinct from Requirement | create/update/withdraw/supersede | canonical object absent | Objects |
| Intelligence Hypothesis / Competing Assessment | Investigate concepts | alternative explanations | create/update/review/dispute/withdraw | distinct from Case/Detection hypotheses | Objects |
| Source Fusion Assessment | Investigate concept | dependency, repetition, corroboration and gaps | run/create/review/supersede | algorithm/model not selected | Objects/Technique |
| Threat Actor / Attribution Assessment | Investigate concepts | uncertain actor/attribution judgment | create/update/review/withdraw/supersede | identity/attribution model absent | Objects/Trust |
| Campaign / Activity / Intrusion Set Assessment | Investigate concepts | assessed grouping | create/review/separate/merge proposal | candidate relation future | Objects |
| Malware / Tool / Infrastructure / Capability Assessment | Investigate concepts | technical intelligence assessment | create/update/review/supersede | canonical schemas absent | Objects |
| Intelligence Product Plan / Audience Definition | Investigate concepts | objective, audience, type and constraints | create/update/pause/supersede | audience policy OPEN-019 | Objects/Permissions |
| Intelligence Product Draft / Version | Investigate concepts using Shared Versioning | authored structured content | create/update/compare/restore/archive/supersede | not canonical Report | Objects/Shared |
| Intelligence Review / Release Recommendation | Investigate concepts | quality disposition and recommendation | assign/comment/review/disposition | not Govern Approval | Objects/Trust |
| Releasability Assessment / Dissemination Plan | Investigate concepts | markings, redaction, audience and handling plan | create/update/review/supersede | policy OPEN-019 | Objects/Permissions |
| Internal Publication Record / Consumer Access Record | Investigate concepts using Shared delivery | reversible internal publication/access history | create/publish/suspend/withdraw/supersede by policy | delivery/access contracts future | Shared/Permissions |
| Watchlist Definition | Investigate concept | functional definition and lifecycle | create/update/review/handoff | activation owned by runtime | Objects/Technique |
| Indicator Operationalization Package | Investigate concept | Detection/Command/Settings handoff | prepare/version/withdraw/supersede | not deployed Indicator/rule | Objects/Technique |
| Intelligence Monitoring Plan / Change Assessment | Investigate concepts | authorized source monitoring and significant change | create/update/run-request/review | not target surveillance/Alert | Objects/Studio |
| External Sharing Package | Investigate concept | candidate release package and Action Request | prepare/version/withdraw | no transmission; OPEN-019 | Govern/Technique |
| Consumer Feedback Record / Effectiveness Assessment | Investigate concepts | perceived utility, usage and limitations | create/compare/review/supersede | feedback ≠ ground truth | Objects/Metrics |
| Requirement Satisfaction Assessment / Collection Feedback Package | Investigate concepts | satisfaction and collection gaps | assess/reopen/prepare handoff | no collection execution | Objects/Settings/Studio |
| Correction / Retraction / Product Supersession | Investigate concepts using Shared Versioning | error handling and history | create/review/withdraw/supersede | external recall governed | Objects/Govern |
| Continuous Improvement Package | Investigate concept | complete lifecycle handoff | prepare/version/withdraw | no active mutation | Objects |
| Provenance Record | Shared mechanisms + product semantics | complete source/tool/human lineage | read/emit/link/export by permission | audit contracts future | Shared/Trust |
| Cloud Investigation Session / Cloud Scope Assessment | Investigate concepts | bounded Cloud workspace and selected scope | create/update/close/reopen/version/supersede | no final object/schema | Objects |
| Cloud Organization / Cloud Tenant / Cloud Account / Cloud Subscription / Cloud Project | Platform Settings projections + Investigate scope semantics | distinct provider-neutral hierarchy units | read/select/include/exclude/link | provider mapping OPEN-012 | Settings/Objects |
| Cloud Resource Observation / Configuration Observation | Investigate concepts | observed resource, configuration, state, version and drift | create/compare/dispute/supersede | observed ≠ current | Objects |
| Cloud Identity Observation / Role Observation / Permission Observation / Effective Permission Candidate | Investigate concepts using source projections | identity, assignment, policy and authorization analysis | create/review/dispute/withdraw/supersede | identity/permission final schemas absent | Objects/Permissions |
| Cloud Audit Event / Cloud Activity Observation | source owner + Investigate concept | actor/action/resource/result/time analysis | read/normalize/link/annotate/version | provider event schemas OPEN-012 | Settings/Objects |
| Compute Observation / Workload Observation / Container Observation / Serverless Observation | Investigate concepts | control-plane and available runtime context | create/compare/handoff/supersede | not full Endpoint/container forensics | Objects/Endpoint |
| Cloud Network Observation / Network Exposure Observation | Investigate concepts using Shared Graph | configured and observed connectivity/exposure | create/compare/correlate/dispute | no active scan; not Network Forensics | Objects/Shared |
| Storage Observation / Data Access Observation | Investigate concepts | metadata, access events, retention/versioning and movement candidates | create/compare/dispute/handoff | metadata ≠ content | Objects/Trust |
| Sensitive Material Candidate | Investigate concept; Settings/Security own secret administration | candidate presence, metadata, masked view and audit context | create/review/request access/prepare handoff | reveal/copy/export/use remain separate | Objects/Permissions |
| Cloud Anomaly / Cloud Hypothesis | Investigate concepts | candidate interpretation with support, contradictions and confidence | create/review/dispute/withdraw/supersede | anomaly ≠ Finding/compromise | Objects |
| Cloud Timeline / Cross-source Correlation Candidate | Investigate semantics using Shared Timeline/Linking | Cloud-local chronology and cross-source candidate relations | create/reconstruct/compare/dispute/supersede | not Case Timeline; correlation ≠ causality | Shared/Objects |
| Cloud Reproducibility Assessment / Cloud Analysis Provenance Package | Investigate concepts using Shared Trace/Versioning/Export | source/tool/permission/version lineage and reproducibility conditions | create/compare/export by permission/supersede | no guaranteed replay | Shared/Trust |
| Evidence Candidate Package / Finding Draft / Cloud Detection Gap | Investigate preparation; destination owners qualify/create canonical objects | Evidence, Finding and Detection handoffs | prepare/version/withdraw/link | candidate/draft/package ≠ destination object | Evidence/Detection/Objects |
| Mobile Device / Mobile Device Candidate / Mobile Platform Candidate | source/Settings projections + Investigate concepts | scoped device/platform/version candidates | read/create/review/dispute/supersede | identifier/platform candidate ≠ person/certain support | Objects/Settings |
| Mobile Acquisition Request / Collection Request / Collection Job | Mobile request concept + Collection owner | acquisition need/context/status | prepare/read/link only | acquisition execution remains Collection | Collection/Objects |
| Mobile Evidence Package / Device Backup / Logical Extraction / Filesystem Extraction | source/Collection + Investigate concepts | bounded evidence representations | read/compare/review limitations | package ≠ device; backup ≠ extraction; extraction ≠ original | Objects/Collection |
| Mobile Forensics Session | Investigate concept | durable mobile analytical workspace | create/update/pause/close/reopen/archive/supersede | final schema absent | Objects |
| Integrity Assessment / Completeness Assessment / Accessibility Assessment | Investigate concepts using source/Trust facts | separate trust/coverage/access judgments | create/review/dispute/supersede | integrity ≠ completeness ≠ relevance ≠ usability | Objects/Trust |
| Mobile Filesystem Observation | Investigate concept | represented partitions/filesystems/storage/entries | create/search/compare/dispute/supersede | no final filesystem model | Objects |
| Application Observation / Application Data Observation | Investigate concepts | app/package and local-record analysis | create/compare/review/dispute/supersede | installed ≠ used; data ≠ user intent | Objects |
| Communication Observation / Call Observation / Contact Observation | Investigate concepts | messaging/call/contact analysis | create/review/dispute/supersede | stored record ≠ author/read/conversation/relationship truth | Objects/Trust |
| Media Observation / Document Observation | Investigate concepts | media/document content and metadata | create/compare/review/dispute/supersede | file ≠ user-created; metadata ≠ certain truth | Objects/Trust |
| Location Observation / Sensor Observation | Investigate concepts | device location/activity/sensor context | create/compare/dispute/supersede | device location ≠ user presence; sensor ≠ human action | Objects/Trust |
| Account Observation / Sensitive Material Candidate | Investigate concepts; Settings/Security own secret administration | account context and masked sensitive candidates | create/review/dispute/request access/handoff | account ≠ person; candidate ≠ valid/usable; use prohibited | Objects/Permissions |
| Connectivity Observation / SIM-eSIM Observation / Paired Device Observation | Investigate concepts using source projections | network/wireless/SIM/pairing context | create/compare/correlate/dispute/handoff | record ≠ connection; SIM/paired device ≠ person/owner; no active scan | Objects/Shared |
| Synchronization Observation | Investigate concept | local/synchronized/cloud-backed/cross-device relation | create/compare/dispute/supersede | synchronized/cloud-backed ≠ certain local state/full Cloud analysis | Objects/Cloud |
| Deleted Entry Candidate / Recovery Result | Investigate concepts | deleted/residual/recovered/carved analysis | create/run bounded recovery/review/invalidate/supersede | deleted ≠ intent; recovered ≠ complete original/attribution | Objects/Studio |
| Mobile Anomaly / Mobile Hypothesis | Investigate concepts | candidate anomaly/persistence/compromise interpretation | create/review/dispute/withdraw/supersede | anomaly/persistence/suspicious app ≠ compromise/malware/Finding | Objects |
| Mobile Timeline / Cross-source Correlation Candidate | Investigate semantics using Shared Timeline/Linking | mobile-local chronology and cross-source relations | create/reconstruct/compare/dispute/supersede | Mobile Timeline ≠ Case Timeline; correlation ≠ causality | Shared/Objects |
| Derived Artifact | Investigate canonical Artifact semantics + Mobile derivation | bounded transformed/selected mobile material | create/version/export/withdraw | Derived Artifact ≠ Evidence | Objects/Shared |
| Evidence Candidate Package / Finding Draft / Detection Engineering Package / Intelligence Handoff Package | Investigate preparation; destination owners qualify/create | downstream handoffs | prepare/version/withdraw/link | candidate/draft/package ≠ destination object/action | Evidence/Detection/TI/Objects |
| Reproducibility Assessment / Mobile Analysis Provenance Package | Investigate concepts using Shared Trace/Versioning/Export | source/acquisition/tool/permission/human lineage | create/review/dispute/export/supersede | no guaranteed replay; no trace deletion | Shared/Trust |

No complete schema, JSON Schema, final cardinality/state machine, physical graph, algorithm, provider/platform event model, extraction format, filesystem/application/message model, technical identifier, runtime package or atomic permission is defined.
