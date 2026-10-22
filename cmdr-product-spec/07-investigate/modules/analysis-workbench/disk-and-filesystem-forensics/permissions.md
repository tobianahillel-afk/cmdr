---
id: investigate-disk-filesystem-permissions
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-SEC-001
  - REQ-SEC-002
open_decisions:
  - OPEN-013
  - OPEN-014
---
# Functional permissions

| Need | Capabilities | Risk | Class | Masking / step-up / separation | Owner / future phase |
|---|---|---|---:|---|---|
| Disk Image read / raw / restricted read | 363,365,378 | highly sensitive source | 0 | minimization; step-up for raw/restricted; least privilege | Investigate + Security / Permissions |
| Disk Session create/update/close/reopen | 364 | reversible analytical state | 2 | versioning; owner/reviewer where required | Investigate / Permissions |
| partition/volume/filesystem select | 366 | analytical interpretation | 2 | confidence visible; possible reviewer | Investigate / Permissions |
| filesystem browse / metadata read | 367 | sensitive paths and metadata | 0 | path/content masking | Investigate + Security |
| file preview / sensitive content read | 368,376 | private or secret content | 0 | reveal separately gated; audit | Security / Permissions |
| deleted/unallocated/journal read | 369,370 | recovered private data | 0 | minimization and purpose limitation | Investigate + Security |
| system/user/application/persistence artifacts read | 371..373 | attribution/privacy | 0 | labels and masking; reviewer | Investigate + Privacy |
| timeline read / comparison create | 374,377 | aggregated sensitive data | 0/1 | source-level authorization | Investigate + Shared |
| recovery prepare/run/read/export | 375 | volume, privacy and diffusion | 1/2 | explicit scope, step-up, producer/reviewer | Investigate + Security |
| restricted content request/read | 376 | protected content | 0/2 | no bypass; authorization and audit | Security / Settings |
| Derived Artifact create/read/export | 368,369,375,376 | derivation and export | 1/2 | lineage, redaction, export gate | Investigate + Shared |
| Evidence/Finding/Detection/Network handoff prepare | 379 | premature qualification | 2 | destination review; no auto-confirm | Investigate / owner phases |
| automated disk analysis request | 363..378 | Tool execution | 1 | visible Tool/parameters; human owner | Studio + Investigate |
| cross-tenant analysis | all | tenant isolation | 0/1 | strong step-up and separate approval | Security / Permissions |
| custody/provenance/reproducibility review/export | 365,378 | audit and disclosure | 0/1/2 | immutable trace; export controls | Investigate + Shared |

No atomic matrix, namespace, final RBAC/ABAC, step-up or separation rule is finalized.
