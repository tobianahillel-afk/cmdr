---
id: investigate-disk-filesystem-scope
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
---
# Scope

## In scope
Disk intake/session, image integrity and acquisition context, partition/volume/filesystem candidates, read-only navigation, metadata, file identity and relationships, deleted/unallocated data, filesystem journals, persistent system/configuration/user/application/startup artifacts, Disk Timeline, bounded carving/recovery, encrypted/compressed/restricted handling, multi-image comparison, provenance/reproducibility and owner handoffs.

## Out of scope
Acquisition implementation, live filesystem mutation, final image format, engine/plugin/tool choice, low-level structures or offsets, carving algorithm, password attack, encryption bypass, anti-forensics, API, protocol, database, storage design, product code, full Network Forensics, Cloud/Mobile analysis and Detection Engineering deployment.

## Source rule
Disk Forensics consumes Disk Image or equivalent Artifacts produced or imported through owner capabilities. It never acquires or changes the source Endpoint.
