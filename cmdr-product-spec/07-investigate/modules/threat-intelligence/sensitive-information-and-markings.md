---
id: investigate-threat-intelligence-sensitive-information-and-markings
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
open_decisions: [OPEN-018, OPEN-019]
---
# Sensitive information, markings and releasability

## Distinct rights
1. existence; 2. metadata; 3. masked preview; 4. read; 5. copy; 6. extraction; 7. inclusion in a product; 8. internal publication; 9. export; 10. external-sharing preparation; 11. future external sharing; 12. future recall/revocation.

## Preserved context
Source restrictions, licence, classification, tenant/environment, customer/victim/personal data, confidential-source identity, provider identity, investigation information, secrets, raw access, derivative handling, markings, releasability, versions, consumers and access traces survive every transformation and handoff.

## Invariants
- product permission ≠ source permission;
- internal publication ≠ export;
- export ≠ external sharing;
- marking ≠ permission;
- releasable ≠ shared;
- internal ≠ cross-tenant/client/external/public;
- masking ≠ deletion;
- anonymization candidate ≠ guaranteed anonymization;
- withdrawn source ≠ automatic product deletion;
- retraction ≠ erasure of historical access.

CAP-INV-528 prepares handling/releasability; CAP-INV-529 publishes internally only under policy; CAP-INV-533 prepares external sharing only. OPEN-019 governs unresolved audience, cross-tenant, client, external and recall policy.
