---
id: investigate-mobile-forensics-scope
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-07
source-of-truth: canonical
requirements: [REQ-INV-001, REQ-PROD-014, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-005, OPEN-008, OPEN-011, OPEN-013, OPEN-014, OPEN-015]
---
# Mobile Forensics Scope

## In scope
- open analysis from Case, Incident, Finding, Hunt, Signal, Collection Job, Artifact or Mobile Evidence Package;
- declare device candidates, platform/version candidates, scope and return origin;
- review declared logical, filesystem, backup, synchronized backup, selected-file, application-export or other evidence representation without prescribing acquisition technique;
- assess integrity, completeness, accessibility, custody and limitations;
- analyze partitions/filesystems/storage, application inventory/data, communications, calls/contacts, media/documents, location/sensors, accounts/sensitive material, connectivity/SIM/eSIM/paired devices, backups/synchronization and deleted/residual/recovered data;
- build Mobile Timeline semantics, correlations, Mobile Anomalies and Mobile Hypotheses;
- select bounded Derived Artifacts and prepare Evidence Candidate, Finding Draft, Detection Engineering, Threat Intelligence and complementary Collection handoffs;
- preserve full source/tool/human provenance and work entirely without AI.

## Explicitly out of scope
No real integration, connector, API, protocol, proprietary format, third-party tool, engine, final acquisition model, unlock, password/code cracking, bypass, jailbreak, rooting, exploit, command, script, offensive extraction procedure, device modification, profile installation, MDM action, wipe, isolation, revocation or product code.

## Platform neutrality
The functional model may record a declared or candidate platform and version but selects no supported platform. Phone/tablet, managed/unmanaged, SIM/eSIM and related wearable records are analysis context only. Wearables may be linked sources; this phase does not create wearable forensics.

## Authority
Access must preserve tenant, environment, purpose, source restriction, classification, collection authority, privacy restrictions and return origin. Cross-device and cross-tenant correlations do not broaden source permission. Missing permission produces a masked/denied state rather than data leakage.

## Closure boundary
OPEN-011 remains open for platform/service/method/tool delivery strategy. Its open status does not invalidate provider-neutral functional coverage after documentary verification, but it blocks implementation decisions that depend on a selected platform or acquisition method.
