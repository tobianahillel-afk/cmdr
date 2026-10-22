---
id: investigate-disk-filesystem-forensics
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-INV-001
  - REQ-PROD-014
  - REQ-PROD-020
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Disk and Filesystem Forensics

## Mission
Permit an authorized analyst to review a Disk Image in read-only analytical context, understand partitions, volumes, filesystems, entries, persistent artifacts, deleted or unallocated data, temporal evidence and derived outputs without defining acquisition, implementation, engine, format or full Network Forensics.

## Capability range
CAP-INV-363..379 are canonical for Phase 4B.2B.3B.1.

## Ownership
Investigate owns Disk Image, analytical context, observations, annotations, Derived Artifacts and candidates/drafts. Endpoint Agent contributes authorized acquisition results. Platform Settings owns Fleet, policies, storage, retention and health. Studio owns Tool, Tool Call, Workflow and Automation Run. Govern owns authority on real targets. Shared owns generic mechanisms.

## Workbench constraints
One main canvas, one right Inspector, at most two auxiliary panels and six visible technical tabs, optional bottom Console, contextual Artifact Explorer, Automation Tray closed by default, preserved selection/history/return/focus, keyboard resizers and structured alternatives to graphs.

## Invariants
- source image is never modified;
- no content is executed;
- partial, deleted, recovered, inferred and inaccessible states remain explicit;
- no Tool/AI result becomes a conclusion, Evidence, Finding or deployed rule automatically;
- full Network Forensics, Cloud Analysis, Mobile Forensics and Phase 4B.3 remain outside this phase.
