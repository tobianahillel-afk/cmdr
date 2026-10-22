---
id: investigate-cloud-analysis-sensitive-information
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements: [REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-012, OPEN-013]
---
# Cloud Analysis sensitive information

## Sensitive families
- principal and identity attributes;
- account, tenant and cross-account trust data;
- network topology and private endpoints;
- storage object identifiers and classifications;
- source payloads and request context;
- keys, tokens, credentials, certificates and secret references;
- customer, regulated or cross-tenant data;
- Tool parameters, query text and exported packages.

## Secret access ladder
1. candidate presence;
2. metadata;
3. masked preview;
4. reveal;
5. copy;
6. export;
7. use — prohibited in Investigate.

Each level has a separate permission, purpose, audit record and masking rule. No level is implied by a previous one except presence required to describe metadata.

## Minimization
Only the smallest authorized projection is displayed or exported. A denied view must not reveal protected existence or value. Sensitive values are never sent to AI without explicit policy and permission.
