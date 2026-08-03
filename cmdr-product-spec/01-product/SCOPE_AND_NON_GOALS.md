# Scope and Non-goals

## In scope

- Operational incident and alert prioritisation.
- SIEM event search, correlation and timeline reconstruction.
- Case, evidence, entity, hypothesis and finding management.
- Static and dynamic malware analysis, reverse engineering and debugging.
- Memory and disk/artifact forensics.
- Threat-intelligence enrichment.
- Response request, decision, approval, policy gate, playbook, run and rollback management.
- Multi-tenant RBAC/ABAC, immutable audit, integrations, reporting and platform health.
- Cross-console transitions preserving user, tenant, time range, selected entities and investigation context.

## Non-goals

- Replacing every specialist third-party analysis engine; CMDR may orchestrate or embed them.
- Automatically authorising high-impact actions without explicit policy support.
- Treating AI output as evidence or final authority.
- Hiding uncertainty behind a single opaque risk score.
- Mixing customer tenants in search, models, cache, exports or analytics.
- Turning Command Center into a full forensic workbench.
- Turning Investigation Lab into a response approval surface.
- Turning Response & Governance into a second incident dashboard.

## Boundary rule

A capability belongs where its primary user decision occurs. Shared services live under `40-platform/`; shared objects and rules live under `02-domain-model/`.
