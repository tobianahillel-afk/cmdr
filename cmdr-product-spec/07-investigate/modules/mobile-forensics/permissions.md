---
id: investigate-mobile-forensics-permissions
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-07
source-of-truth: canonical
requirements: [REQ-SEC-001, REQ-SEC-002, REQ-INV-001]
open_decisions: [OPEN-011, OPEN-013, OPEN-014]
---
# Mobile Forensics Functional Permissions

This document records functional permission needs, not a final RBAC/ABAC namespace or atomic matrix.

| Need | Capability range | Risk | Class | Masking / step-up | Owner / future phase |
|---|---|---|---:|---|---|
| Mobile Evidence Package read | 701..719 | source exposure | 0 | classification-aware | Investigate / Permissions |
| raw or restricted extraction read | 704..719 | private/raw data | 0 | masking; step-up possible | Security / Permissions |
| Session create/update/close/reopen | 702 | analytical mutation | 2 | normal audit; step-up by policy | Investigate / OPEN-013 |
| platform candidate select/review | 703 | interpretation | 2 | source/confidence visible | Investigate |
| acquisition context read/review | 704 | custody/authority | 0/2 | authority visible | Collection/Investigate |
| integrity assessment read/update | 705 | evidentiary interpretation | 0/2 | separation reviewer possible | Investigate/Trust |
| filesystem/application/data read | 706..708 | private content | 0 | scope + masking | Investigate/Security |
| communications/contacts/calls read | 709 | third-party/private data | 0 | explicit purpose; masking | Security |
| media metadata/content read | 710 | private/sensitive content | 0 | content separate from metadata | Security |
| location/health/sensor data read | 711 | highly sensitive data | 0 | explicit access; step-up possible | Security |
| account metadata read | 712 | identity data | 0 | masking | Security |
| sensitive material presence/metadata/reveal/copy/export | 712 | credential exposure | 0/1/2 | distinct permissions; reveal/export step-up | Security/Settings |
| network/SIM/eSIM/paired-device read | 713 | identity/location linkage | 0 | scope/masking | Security |
| backup/synchronized data read | 714 | cross-device/cloud data | 0 | source permission required | Security/Settings |
| deleted/recovered-data read/recovery | 715 | residual/private data | 0/1 | explicit source scope | Investigate/Security |
| timeline/correlation create | 716 | cross-source inference | 1/2 | cross-tenant gate | Investigate/Shared |
| anomaly/Hypothesis create/review | 717 | analytical conclusion risk | 2 | human disposition | Investigate |
| Derived Artifact create/read/export | 718 | content derivation | 1/2 | export separate; masking retained | Investigate/Shared |
| Evidence Candidate/Finding Draft/Detection/TI handoff prepare | 718 | downstream influence | 2 | destination review required | destination owners |
| automated analysis request | 701..719 | data to Tool/model | 1/2 | source permission + policy | Studio/Security |
| cross-tenant analysis | all | isolation breach | 2 | explicit authority, no implied access | Security/Settings |
| provenance export / reproducibility review | 719 | audit/private metadata | 0/1 | classification retained | Shared/Trust |

## Rules
Read never implies export, execution, reveal, copy, sharing or approval. Session access never grants raw source access. Cross-device/cross-tenant relations never widen permission. A denied action explains the required capability without leaking protected content. All access decisions and sensitive-access attempts are auditable.

Final permission names, ABAC expressions, step-up policy and separation-of-duties matrix remain future Security/Permissions work; OPEN-013 stays open.
