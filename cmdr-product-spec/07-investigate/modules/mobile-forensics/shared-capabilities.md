---
id: investigate-mobile-forensics-shared-capabilities
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-07
source-of-truth: canonical
requirements: [REQ-PROD-019, REQ-PROD-020, REQ-INV-001]
open_decisions: [OPEN-014, OPEN-015]
---
# Mobile Forensics Shared Capabilities

| Shared capability | Mobile use | Local semantics retained |
|---|---|---|
| Entity / Graph | device/account/contact/paired-device candidate relations | relationship candidate, source, confidence, contradiction; no identity merge |
| Timeline | rendering/order of Mobile Timeline entries | mobile timestamp quality and scope; Mobile Timeline ≠ Case Timeline |
| Search | filesystem/application/record search | source package, scope and permission context |
| Object Linking | Case, Artifact, package, observation and handoff links | typed/source-aware return origin |
| Versioning / Comparison | package, backup, extraction, observation and assessment comparison | partiality, supersession and source version |
| Background Jobs | bounded parsing/extraction/recovery/correlation work | business operation remains Mobile/Collection/Studio owned |
| Notifications | job completion, restriction, review and correction | deep link and tenant/scope preserved |
| Trace / Activity | source/tool/human lineage and access audit | Mobile-specific provenance fields |
| Export | authorized minimized output | masking/classification/source restrictions retained |
| Reporting | render Mobile assessments within owner reports | no new Report object |
| Collaboration | comments/assignments/review | analytical ownership and permissions retained |
| Recovery | UI/workflow recovery of reversible analytical state | no recovery of deleted source or real-device mutation |
| Inspector / Context Bar | selected source/record/scope/restriction/return context | no Mobile-specific competing Inspector |

Shared mechanisms never own Mobile Forensics business decisions. Mobile never recreates a Graph, Timeline engine, Search, Export, Reporting, Trace, Activity, Versioning, Collaboration or Recovery engine.
