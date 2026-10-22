---
id: investigate-mobile-forensics-source-migration
domain: 07-investigate
status: draft
owner: Product Architecture
updated: 2026-08-07
source-of-truth: canonical
requirements: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-020]
open_decisions: [OPEN-011]
---
# Mobile Forensics Source Migration

## Audit result
No active canonical `mobile-forensics/` module, Mobile capability shard, Mobile roadmap file, `CAP-INV-7xx` family or Mobile Screen ID existed at the starting SHA `ed874ea414fc57f24fa61f410f91b7345f4a868a`.

Historical/generic Mobile terms were found only as future-scope references or exclusions in existing owner documents. They do not constitute a competing Mobile functional architecture.

## Canonical replacement
`07-investigate/modules/mobile-forensics/` and `CAP-INV-701..719` become the functional Mobile source once published and verified. Existing owner sources remain active:
- Collection and Live Response for acquisition requests/jobs/execution/custody;
- Disk/Filesystem, Network, Memory, Static, Reverse and Dynamic analysis for their specialized domains;
- Cloud Analysis for Cloud scope;
- Platform Settings for sources/Fleet/MDM-like future administration/secrets/retention;
- Endpoint Agent for declared native capabilities when present;
- Govern for authority and real-device actions;
- Studio and Shared for Tools/Runs and generic mechanisms.

## Deprecation disposition
Competing Investigate Mobile functional documents deprecated: **0**. Active screens deprecated: **0**. No prior Mobile source is overwritten or removed.

## Acceptance
A future Mobile-related document that overlaps this module must link to the relevant `CAP-INV-701..719` capability or be explicitly classified as owner-specific, implementation-specific or deprecated. No second active Mobile functional architecture may coexist silently.
