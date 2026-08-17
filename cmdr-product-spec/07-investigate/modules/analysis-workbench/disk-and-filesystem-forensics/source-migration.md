---
id: investigate-disk-filesystem-source-migration
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-INV-001
  - REQ-PROD-014
  - REQ-PROD-020
---
# Source migration

| Legacy source | Useful needs migrated | Replacement | Status |
|---|---|---|---|
| `07-investigate/modules/disk-and-artifact-forensics/README.md` | filesystem, persistent system/browser artifacts, carving and super-timeline | CAP-INV-363..379 and module README | deprecated pointer |
| `07-investigate/modules/disk-and-artifact-forensics/read-only-analysis.md` | read-only invariant, filesystem semantics, privacy minimization and provenance | scope/concepts/permissions plus CAP-INV-365..378 | deprecated pointer |
| `07-investigate/modules/disk-and-artifact-forensics/screens/disk-and-artifact-forensics.md` | browse, carve, timeline, Inspector and mandatory states | screen remains active; capability map updated | active screen |

## Acceptance
No competing functional architecture remains; useful needs are preserved; the active screen is not deprecated; replacement, reason, dependents, migration date and acceptance are explicit. Network screens remain unchanged.
