---
id: investigate-mobile-forensics-capability-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-07
source-of-truth: canonical
requirements: [REQ-PROD-012, REQ-PROD-014, REQ-INV-001]
open_decisions: [OPEN-005, OPEN-008, OPEN-011, OPEN-013, OPEN-014, OPEN-015]
---
# Mobile Forensics Capability Map

| ID | Capability | Primary role | Local focus |
|---|---|---|---|
| CAP-INV-701 | Mobile Forensics Intake and Preconditions | Mobile Forensics Analyst | origin, authority, representation, scope and readiness |
| CAP-INV-702 | Mobile Investigation Session and Workspace Management | Investigation Lead | durable analytical Session and return context |
| CAP-INV-703 | Mobile Device, Platform and Scope Context | Mobile Forensics Analyst | device/platform candidates, included/excluded scope |
| CAP-INV-704 | Mobile Acquisition, Backup and Extraction Context Review | Evidence Reviewer | declared acquisition representation, custody and limits |
| CAP-INV-705 | Mobile Evidence Integrity, Completeness and Accessibility Assessment | Evidence Reviewer | separate integrity/completeness/accessibility judgments |
| CAP-INV-706 | Mobile Filesystem, Partition and Storage Analysis | Mobile Forensics Analyst | represented partitions/filesystems/storage/metadata |
| CAP-INV-707 | Mobile Application Inventory and Package Analysis | Mobile Forensics Analyst | app/package/version/configuration observations |
| CAP-INV-708 | Mobile Application Data, Sandbox and Local Storage Analysis | Mobile Forensics Analyst | databases/files/caches/preferences/local storage |
| CAP-INV-709 | Mobile Communications, Messaging, Calls and Contacts Analysis | Mobile Forensics Analyst | communication records with attribution limits |
| CAP-INV-710 | Mobile Media, Documents and User Content Analysis | Mobile Forensics Analyst | media/documents/metadata/derivation candidates |
| CAP-INV-711 | Mobile Location, Movement and Sensor Artifact Analysis | Sensitive Data Reviewer | location/activity/sensor records and uncertainty |
| CAP-INV-712 | Mobile Accounts, Tokens, Keys and Sensitive Material Assessment | Sensitive Data Reviewer | masked sensitive candidates; no use |
| CAP-INV-713 | Mobile Network, Wireless, SIM and Paired Device Analysis | DFIR Analyst | connectivity/SIM/eSIM/pairing observations; no scan |
| CAP-INV-714 | Mobile Backup, Synchronization and Cross-Device Artifact Analysis | Mobile Forensics Analyst | versions/sync/cross-device/cloud-backed ambiguity |
| CAP-INV-715 | Mobile Deleted, Residual and Recovered Data Analysis | DFIR Analyst | deleted/residual/recovered/carved candidates |
| CAP-INV-716 | Mobile Timeline and Cross-Source Correlation | DFIR Analyst | time quality and correlations with other domains |
| CAP-INV-717 | Mobile Anomaly, Persistence and Compromise Hypothesis Management | Investigation Lead | candidate interpretation with support/contradiction |
| CAP-INV-718 | Mobile Artifact Extraction and Handoff to Evidence, Findings and Detection | Investigation Lead | bounded Derived Artifacts and destination packages |
| CAP-INV-719 | Mobile Forensics Provenance and Reproducibility | Evidence Reviewer | end-to-end source/tool/human lineage |

## Totals
- capabilities: **19**;
- numbered sections expected: **513**;
- mandatory tables expected: **114**;
- IDs: `CAP-INV-701..719`, continuous and reserved after registry audit;
- delivery target: **19 defined / 19 planned**;
- implementation claims: **0**.
