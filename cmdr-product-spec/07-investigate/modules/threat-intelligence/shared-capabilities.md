---
id: investigate-threat-intelligence-shared-capabilities
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
---
# Shared capability consumption

| Shared Capability | Usage Threat Intelligence | Données locales | Source canonique |
|---|---|---|---|
| Entity / Entity Resolution | candidate-to-canonical relation proposal | candidate context and uncertainty | Shared |
| Graph / Object Linking | relationship visualization and navigation | typed sourced candidate relation | Shared |
| Timeline | Sightings, versions and lifecycle history | domain events and time context | Shared |
| Search / Global Search | authorized material and knowledge retrieval | filters and project scope | Shared |
| Versioning / Comparison | candidates, assessments, relations and handoffs | product versions and diffs | Shared |
| Jobs / Notifications | extraction, review, lifecycle and handoff progress | business status and recipient context | Shared |
| Trace / Activity / Audit Hooks | full lineage and human dispositions | product semantics only | Shared |
| Export / Reporting | authorized packages and local report inputs | restrictions and redaction context | Shared |
| Collaboration / Comments / Assignments | analyst/reviewer work | roles, comments and disputes | Shared |
| Recovery | draft, run and conflict recovery | checkpoints and return origin | Shared |
| Inspector / Context Bar | selected candidate/source/project | authorized projection | Design System / Experience |
| Threat Intelligence Enrichment | shared enrichment projection | context is not Evidence or truth | Shared |

No Shared capability is redefined. A Threat Entity Candidate does not become a Shared Entity automatically.
