---
id: investigate-mobile-forensics-collection-acquisition-boundaries
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-07
source-of-truth: canonical
requirements: [REQ-INV-001, REQ-PROD-014, REQ-PROD-020]
open_decisions: [OPEN-005, OPEN-008, OPEN-011, OPEN-013, OPEN-014]
---
# Collection and Acquisition Boundaries

## Ownership
Collection and Live Response owns Collection Request, Collection Job business execution context, acquisition execution/result/status/errors, source transfer and collection-time custody. Endpoint Agent or another future authorized collection mechanism may contribute acquisition results only according to declared capability and authority. Mobile Forensics consumes those results.

## Distinctions
- Mobile Acquisition Request ≠ Collection Request ≠ Collection Job.
- Collection Job ≠ Mobile Evidence Package.
- Mobile Evidence Package ≠ Device Backup.
- Device Backup ≠ Filesystem Extraction ≠ Logical Extraction.
- extraction representation ≠ original device and ≠ guaranteed physical image.
- acquisition authorization ≠ analytical permission.
- integrity verification ≠ completeness ≠ accessibility ≠ suitability.

## Mobile behavior
Mobile may review declared acquisition method, initiator, authority, device state, lock/encryption state, period, transformations, transfers, copies, custody, missing elements, restrictions and errors. It may prepare a complementary bounded Collection Request when a gap is identified. It does not execute acquisition or modify the source device.

## Forbidden technical detail
No unlock procedure, bypass, rooting, jailbreak, code/password testing, exploit, command, script, physical extraction procedure, transport, driver, cable protocol, pairing action, profile installation or acquisition engine is defined.

## Partial and failed acquisition
`partial`, `encrypted`, `locked`, `restricted`, `unsupported`, `inaccessible`, `failed` and `disputed` remain explicit. No inference fills absent areas. Failed/partial collection preserves all valid custody and result records.

## Custody
Collection records collection-time custody. Mobile preserves and consumes it, adds analytical lineage and can dispute or annotate gaps without rewriting the source custody record. Evidence qualification remains in the Evidence workflow.
