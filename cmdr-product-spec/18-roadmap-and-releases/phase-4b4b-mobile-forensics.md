---
id: phase-4b4b-mobile-forensics
domain: 18-roadmap-and-releases
status: draft
owner: Product Operations Lead
updated: 2026-08-07
source-of-truth: canonical
requirements: [REQ-PROD-012, REQ-PROD-014, REQ-PROD-019, REQ-PROD-020, REQ-INV-001]
open_decisions: [OPEN-005, OPEN-008, OPEN-011, OPEN-013, OPEN-014, OPEN-015]
---
# Phase 4B.4B — Mobile Forensics Foundations and Mobile Investigation

## Canonical identity
`4B.4B` is the unique Mobile Forensics subphase. The roadmap audit found no pre-existing canonical Mobile phase or active competing Mobile module. Phase 4B.4 groups extended-environment analysis: `4B.4A` is Cloud Analysis and `4B.4B` is Mobile Forensics.

## Objective
Define provider-neutral, platform-neutral functional capabilities for authorized analysis of Mobile Evidence Packages, backups, extraction representations and mobile-derived artifacts while preserving acquisition authority, privacy, provenance and owner boundaries.

## Planned scope
- `CAP-INV-701..719`, nineteen documentary capabilities;
- intake, Session, device/platform/scope, acquisition context and integrity review;
- filesystem/storage, applications, communications, media, location/sensors, sensitive material and connectivity;
- backups/synchronization, deleted/recovered data, Mobile Timeline, anomaly/Hypothesis management, Derived Artifacts and handoffs;
- complete provenance and reproducibility assessment.

## Boundaries
No platform is selected. No iOS, Android or other release is imposed. No acquisition tool, commercial product, engine, API, protocol, proprietary format, filesystem model, MDM/EMM integration or technical acquisition method is selected. No unlocking, password/code testing, bypass, rooting, jailbreak, exploit, profile installation, device mutation, isolation, lock, wipe, secret use or product code is created.

Collection owns acquisition requests/jobs/execution/results and collection-time custody. Platform Settings owns configured sources, MDM/EMM providers, connectors, credentials, secrets, fleet, retention and policies. Endpoint Agent may contribute only declared authorized capabilities when present. Govern owns authority and real-device actions. Studio owns Tools/Runs. Shared owns generic mechanisms. Investigate owns Mobile analytical concepts only.

## Closure rule
Phase 4B.4 may become PASS only if both Cloud Analysis and Mobile Forensics pass their documentary gates. Phase 4B may become PASS only after a complete audit of all Phase 4B children confirms no mandatory functional scope or canonical blocker remains. A PASS is documentary functional conformance, never implementation proof.

## Publication gates
The Mobile phase remains pending until its five functional commits are published, canonical IDs and links are rechecked remotely, PR #2 remains Draft/open/unmerged, `main` and the root README remain unchanged, and the final gate catalogue has no failure.
