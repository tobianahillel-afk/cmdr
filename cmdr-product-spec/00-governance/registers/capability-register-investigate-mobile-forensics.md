---
id: capability-register-investigate-mobile-forensics
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-07
source-of-truth: registry
requirements: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-014, REQ-PROD-019, REQ-PROD-020, REQ-INV-001, REQ-INV-006]
open_decisions: [OPEN-005, OPEN-008, OPEN-011, OPEN-013, OPEN-014, OPEN-015]
---
# Capability Register — Investigate Mobile Forensics

Canonical phase: **4B.4B — Mobile Forensics Foundations and Mobile Investigation**.

`defined` / `planned` means documented, not implemented. No mobile platform, acquisition tool/method, engine, API, protocol, proprietary format, unlock/bypass/root/jailbreak method, connector, device action or product code is delivered.

| ID | Title | Primary role | Canonical path | Main local concepts | Action classes | Open decisions | Status | Delivery |
|---|---|---|---|---|---|---|---|---|
| CAP-INV-701 | Mobile Forensics Intake and Preconditions | Mobile Forensics Analyst | `07-investigate/modules/mobile-forensics/capabilities/mobile-forensics-intake-and-preconditions.md` | Intake, scope, readiness, source gaps | 0,1,2 | 005/008/011/013/014/015 | defined | planned |
| CAP-INV-702 | Mobile Investigation Session and Workspace Management | Investigation Lead | `07-investigate/modules/mobile-forensics/capabilities/mobile-investigation-session-and-workspace-management.md` | Mobile Session, workspace/version/activity | 0,1,2 | 011/013/014/015 | defined | planned |
| CAP-INV-703 | Mobile Device, Platform and Scope Context | Mobile Forensics Analyst | `07-investigate/modules/mobile-forensics/capabilities/mobile-device-platform-and-scope-context.md` | device/platform candidates, scope | 0,1,2 | 008/011/013/014 | defined | planned |
| CAP-INV-704 | Mobile Acquisition, Backup and Extraction Context Review | Evidence Reviewer | `07-investigate/modules/mobile-forensics/capabilities/mobile-acquisition-backup-and-extraction-context-review.md` | acquisition context, representation, custody gaps | 0,1,2 | 005/008/011/013/014 | defined | planned |
| CAP-INV-705 | Mobile Evidence Integrity, Completeness and Accessibility Assessment | Evidence Reviewer | `07-investigate/modules/mobile-forensics/capabilities/mobile-evidence-integrity-completeness-and-accessibility-assessment.md` | integrity/completeness/accessibility assessments | 0,1,2 | 011/013/014 | defined | planned |
| CAP-INV-706 | Mobile Filesystem, Partition and Storage Analysis | Mobile Forensics Analyst | `07-investigate/modules/mobile-forensics/capabilities/mobile-filesystem-partition-and-storage-analysis.md` | filesystem/storage observations | 0,1,2 | 005/011/013/014/015 | defined | planned |
| CAP-INV-707 | Mobile Application Inventory and Package Analysis | Mobile Forensics Analyst | `07-investigate/modules/mobile-forensics/capabilities/mobile-application-inventory-and-package-analysis.md` | application/package observations | 0,1,2 | 005/011/013/014/015 | defined | planned |
| CAP-INV-708 | Mobile Application Data, Sandbox and Local Storage Analysis | Mobile Forensics Analyst | `07-investigate/modules/mobile-forensics/capabilities/mobile-application-data-sandbox-and-local-storage-analysis.md` | app-data observations, protected records | 0,1,2 | 005/011/013/014/015 | defined | planned |
| CAP-INV-709 | Mobile Communications, Messaging, Calls and Contacts Analysis | Mobile Forensics Analyst | `07-investigate/modules/mobile-forensics/capabilities/mobile-communications-messaging-calls-and-contacts-analysis.md` | communication/call/contact observations | 0,1,2 | 011/013/014/015 | defined | planned |
| CAP-INV-710 | Mobile Media, Documents and User Content Analysis | Mobile Forensics Analyst | `07-investigate/modules/mobile-forensics/capabilities/mobile-media-documents-and-user-content-analysis.md` | media/document observations | 0,1,2 | 011/013/014/015 | defined | planned |
| CAP-INV-711 | Mobile Location, Movement and Sensor Artifact Analysis | Sensitive Data Reviewer | `07-investigate/modules/mobile-forensics/capabilities/mobile-location-movement-and-sensor-artifact-analysis.md` | location/sensor observations | 0,1,2 | 011/013/014/015 | defined | planned |
| CAP-INV-712 | Mobile Accounts, Tokens, Keys and Sensitive Material Assessment | Sensitive Data Reviewer | `07-investigate/modules/mobile-forensics/capabilities/mobile-accounts-tokens-keys-and-sensitive-material-assessment.md` | account observations, sensitive candidates | 0,1,2 | 011/013/014/015 | defined | planned |
| CAP-INV-713 | Mobile Network, Wireless, SIM and Paired Device Analysis | DFIR Analyst | `07-investigate/modules/mobile-forensics/capabilities/mobile-network-wireless-sim-and-paired-device-analysis.md` | connectivity/SIM/pairing observations | 0,1,2 | 005/008/011/013/014/015 | defined | planned |
| CAP-INV-714 | Mobile Backup, Synchronization and Cross-Device Artifact Analysis | Mobile Forensics Analyst | `07-investigate/modules/mobile-forensics/capabilities/mobile-backup-synchronization-and-cross-device-artifact-analysis.md` | backup/sync/cross-device observations | 0,1,2 | 008/011/013/014/015 | defined | planned |
| CAP-INV-715 | Mobile Deleted, Residual and Recovered Data Analysis | DFIR Analyst | `07-investigate/modules/mobile-forensics/capabilities/mobile-deleted-residual-and-recovered-data-analysis.md` | deleted candidates, recovery results | 0,1,2 | 005/011/013/014/015 | defined | planned |
| CAP-INV-716 | Mobile Timeline and Cross-Source Correlation | DFIR Analyst | `07-investigate/modules/mobile-forensics/capabilities/mobile-timeline-and-cross-source-correlation.md` | Mobile Timeline, correlation candidates | 0,1,2 | 011/013/014/015 | defined | planned |
| CAP-INV-717 | Mobile Anomaly, Persistence and Compromise Hypothesis Management | Investigation Lead | `07-investigate/modules/mobile-forensics/capabilities/mobile-anomaly-persistence-and-compromise-hypothesis-management.md` | anomaly, persistence candidate, Mobile Hypothesis | 0,1,2 | 005/011/013/014/015 | defined | planned |
| CAP-INV-718 | Mobile Artifact Extraction and Handoff to Evidence, Findings and Detection | Investigation Lead | `07-investigate/modules/mobile-forensics/capabilities/mobile-artifact-extraction-and-handoff.md` | Derived Artifact, Evidence Candidate, Finding Draft, handoffs | 0,1,2 | 005/011/013/014/015/017/018 | defined | planned |
| CAP-INV-719 | Mobile Forensics Provenance and Reproducibility | Evidence Reviewer | `07-investigate/modules/mobile-forensics/capabilities/mobile-forensics-provenance-and-reproducibility.md` | provenance package, reproducibility assessment | 0,1,2 | 011/013/014/015 | defined | planned |

## Totals
- capabilities: **19**;
- numbered sections: **513**;
- mandatory tables: **114**;
- delivery: **19 defined / 19 planned**;
- IDs: `CAP-INV-701..719`, continuous and unique;
- implementation claims: **0**.
