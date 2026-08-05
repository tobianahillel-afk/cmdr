---
id: investigate-detection-engineering-shared-capabilities
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
---
# Shared capability consumption

| Shared Capability | Usage Detection Engineering | Données locales | Source canonique |
|---|---|---|---|
| Background Jobs | validation and replay progress | job references and disposition | Shared |
| Notifications | completion, error, permission and review events | product-specific event context | Shared |
| Trace / Activity / Audit Hooks | lineage and action attribution | business relations only | Shared |
| Versioning / Comparison | drafts, metadata, outcomes and results | product versions and diffs | Shared |
| Linking / Search | source and cross-product navigation | linked refs and filters | Shared |
| Export / Reporting | authorized validation and review outputs | package content and restrictions | Shared |
| Collaboration / Comments / Assignments | project contribution and review | product roles and discussion context | Shared |
| Inspector / Context Bar | selected object and scope | authorized projections | Design System / Experience |
| Recovery | drafts, runs and conflict recovery | recovery checkpoints | Shared |

No Shared capability is redefined.
