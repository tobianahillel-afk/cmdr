---
id: investigate-phase-4b4-capability-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-07
source-of-truth: canonical
requirements: [REQ-PROD-012, REQ-PROD-014, REQ-PROD-020, REQ-INV-001]
open_decisions: [OPEN-011, OPEN-012]
---
# Phase 4B.4 Capability Map — Cloud and Mobile Analysis

| Subphase | Range | Capabilities | Sections | Mandatory tables | Documentary state before final Mobile verification |
|---|---|---:|---:|---:|---|
| 4B.4A — Cloud Analysis Foundations and Cloud Investigation | CAP-INV-601..618 | 18 | 486 | 108 | PASS AFTER POST-PUBLICATION VERIFICATION |
| 4B.4B — Mobile Forensics Foundations and Mobile Investigation | CAP-INV-701..719 | 19 | 513 | 114 | PENDING POST-PUBLICATION VERIFICATION |
| **Phase 4B.4 total** | **CAP-INV-601..618 + 701..719** | **37** | **999** | **222** | **PARTIAL until Mobile remote verification** |

## Boundaries
Cloud and Mobile remain separate functional modules. Cloud-backed Mobile records do not replace Cloud Analysis; mobile network observations do not replace Network Forensics; acquisition remains Collection-owned; configured providers/sources/Fleet remain Settings-owned; real target/device action remains Govern-owned.

## Closure condition
Phase 4B.4 becomes PASS only when Cloud remains verified PASS and Mobile passes its complete post-publication gate catalogue. OPEN-011 and OPEN-012 may remain open for implementation/delivery strategy if they do not invalidate provider-neutral functional completeness.
