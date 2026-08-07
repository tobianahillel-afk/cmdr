---
id: investigate-mobile-forensics-workflows
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-07
source-of-truth: canonical
requirements: [REQ-INV-001, REQ-PROD-014, REQ-PROD-020]
open_decisions: [OPEN-011, OPEN-013, OPEN-014, OPEN-015]
---
# Mobile Forensics Workflows

## Primary workflow
1. Open Mobile Intake from an authorized origin and preserve return origin.
2. Review device, platform, scope, restrictions, acquisition representation and authority.
3. Open or resume a Mobile Forensics Session.
4. Review acquisition context, custody, integrity, completeness and accessibility before interpreting data.
5. Inspect filesystem/storage and application inventory/data according to available representation and permission.
6. Analyze communications, media, location/sensors, account/sensitive material and connectivity only where authorized.
7. Compare backups, synchronized and cross-device observations while preserving local/cloud ambiguity.
8. Review deleted/residual/recovered candidates and recovery quality where available.
9. Build a Mobile Timeline and cross-source correlations with explicit timestamp quality.
10. Create/review Mobile Anomalies and Mobile Hypotheses with supporting, contradicting and missing elements.
11. Select bounded source material for Derived Artifacts and prepare destination-owned handoffs.
12. Review provenance and reproducibility; close, archive or supersede the Session without deleting history.

## Required transitions
| Source | Trigger | Destination | Owner at destination | Context |
|---|---|---|---|---|
| Case / Incident / Finding / Hunt / Signal | mobile analysis need | CAP-INV-701 | source owner retained | objective, device/source refs, period, restrictions, return origin |
| Collection Job / Artifact / Mobile Evidence Package | result available | CAP-INV-701 | Investigate | request/job/package/source, custody, errors, partiality |
| CAP-INV-701 | intake ready | CAP-INV-702/704 | Investigate | scope, permissions, representation, limitations |
| CAP-INV-704 | acquisition context reviewed | CAP-INV-705 | Investigate | method declared, custody, transformations, gaps |
| CAP-INV-702 | analysis starts | CAP-INV-703 | Investigate | Session, device candidates, scope |
| CAP-INV-703 | represented storage available | CAP-INV-706 | Investigate | device/platform/scope + package refs |
| CAP-INV-706 | apps identified | CAP-INV-707/708 | Investigate | filesystem/application area observations |
| CAP-INV-708 | relevant records | CAP-INV-709/710/712/715 | Investigate | records, restrictions, timestamps, app context |
| CAP-INV-703 | location/connectivity context | CAP-INV-711/713 | Investigate | device scope and authorized observations |
| Backup/package | sync relation available | CAP-INV-714 | Investigate | backup/source/version/sync scope |
| Observations | temporal/correlation need | CAP-INV-716 | Investigate | sourced events and timestamp qualities |
| Observations / CAP-INV-716 | interpretation needed | CAP-INV-717 | Investigate | support, contradiction, gaps, confidence |
| Any bounded observation | artifact/handoff need | CAP-INV-718 | Investigate then destination owner | source lineage, restrictions, purpose |
| Missing data | collection gap | CAP-INV-202 | Collection | source, need, authority, bounded scope; no execution by Mobile |
| CAP-INV-718 | technical file produced | Static / Reverse / Dynamic | destination owner | Derived Artifact + lineage, no auto conclusion |
| Mobile result package | Evidence/Finding/Detection/TI need | owner workflow | destination owner | candidate/draft/package only |
| CAP-INV-701..718 | closure/review | CAP-INV-719 | Investigate/Shared provenance | complete lineage, Tools/Runs, errors and human decisions |

## Failure behavior
Missing source, lock/encryption, corruption, unsupported representation, permission denial, stale data, cross-tenant restriction or failed Tool never erases valid context. The workflow becomes partial/blocked/restricted with explicit gaps and a safe return path.
