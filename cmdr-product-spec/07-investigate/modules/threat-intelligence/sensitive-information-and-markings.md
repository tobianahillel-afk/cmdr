---
id: investigate-threat-intelligence-sensitive-information
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements:
  - REQ-SEC-001
  - REQ-SEC-002
open_decisions:
  - OPEN-014
  - OPEN-018
---
# Sensitive information and markings

## Access levels
1. existence;
2. metadata;
3. masked preview;
4. read;
5. copy;
6. extraction;
7. relation;
8. export;
9. future external sharing.

Each level is independently permissioned. Inaccessible content is not absent; a denied action must not reveal protected data.

## Protected dimensions
Source restrictions, classifications, handling markings, licence/authorized use, personal data, commercial information, secrets, victim/customer/investigation data, cross-tenant data, confidential-source identity, provider identity and raw-content access.

## Invariants
- no default reveal or unrestricted copy;
- no model transfer without explicit permission and policy;
- no export bypassing markings, licence or tenant scope;
- no external sharing in 4B.3B.1;
- no marking loss during extraction, normalization, relation, version or merge review;
- no merge that erases provenance or source-specific restrictions;
- audit actor, purpose, source version, permission, disposition and redaction reason.

Detailed dissemination, release markings and external audience policy remain future 4B.3B.2.
