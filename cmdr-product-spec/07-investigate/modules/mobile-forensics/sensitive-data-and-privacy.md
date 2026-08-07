---
id: investigate-mobile-forensics-sensitive-data-and-privacy
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-07
source-of-truth: canonical
requirements: [REQ-SEC-001, REQ-SEC-002, REQ-PROD-020]
open_decisions: [OPEN-011, OPEN-013, OPEN-014]
---
# Mobile Forensics Sensitive Data and Privacy

## Data classes considered
Owner data, third-party data, messages, contacts, photos, videos, location, health/activity, biometric-classified records when present, accounts, tokens, secrets, professional/client/victim data, potentially privileged communications, cross-tenant material, and records classified as belonging to children or vulnerable persons. This document creates no legal determination or jurisdiction-specific rule.

## Access ladder
1. existence;
2. metadata;
3. masked preview;
4. read;
5. reveal;
6. copy;
7. bounded extraction;
8. export;
9. inclusion in an Evidence Candidate Package;
10. future sharing by destination authority.

Each step is a separate functional permission. Earlier access never grants a later step.

## Rules
- collect/analyze the minimum scope necessary for the stated purpose;
- show tenant, device/package scope, period, source and restrictions at access time;
- mask sensitive values by default;
- require explicit access for raw communications, media content, location/health records and sensitive material;
- preserve source classification, markings, custody, retention and handling restrictions in derivatives and handoffs;
- audit reason, actor, permission, source, operation, result and return origin for sensitive access;
- do not send restricted/private content to a model without explicit permission and policy;
- do not reveal protected content in a denial or error message;
- do not use a token, key, certificate or secret discovered during analysis;
- Session permission does not grant permission to the raw package/extraction;
- cross-tenant/cross-device correlation requires permission to each source and does not create a combined access grant;
- Derived Artifacts remain masked/classified according to their source content and purpose.

## Attribution limits
A device identifier, account, contact, phone number, message record, location, SIM/eSIM, paired device or health/activity record is not sufficient by itself to attribute a person, authorship, intent, presence or relationship. Human reviewers must preserve contradictory and missing evidence.

## Retention and custody
Mobile analysis consumes retention/legal-hold/custody projections from their owners. It may record analytical restrictions and handoff needs but does not set global retention, legal authority or sharing policy.
