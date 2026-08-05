---
id: investigate-action-classification
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
---
# Action classification — Investigate through Phase 4B.2B.3B.1

| Class | Meaning | Disk and Filesystem Forensics examples |
|---:|---|---|
| 0 | observation | inspect source/integrity, navigate, search, preview without execution, read metadata/journals/artifacts/timelines and compare views |
| 1 | bounded isolated processing or extraction | explicit forensic Tool run, bounded carving/recovery, Derived Artifact extraction, comparison, authorized export and reproduction |
| 2 | reversible analytical mutation | session/filesystem selection, annotation, dispute, link, bookmark, withdrawal from active use, access request and handoff preparation |
| 3/4 | real-target, destructive or authority action | excluded; block or route to Collection/Live Response and Govern |

CAP-INV-363..379 use only classes 0, 1 and 2. OPEN-013 remains open. Disk analysis is not acquisition, live filesystem mutation or an Endpoint action.
