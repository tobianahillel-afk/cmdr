---
id: investigate-memory-forensics-shared-capabilities
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
# Shared Capabilities consumed

| Shared Capability | Usage Memory Forensics | Données locales | Source canonique |
|---|---|---|---|
| Background Jobs | progression/cancel/partial d’analyses bornées | session, Tool Call relation | `12-shared-capabilities/background-jobs.md` |
| Notifications | résultats, erreurs et restrictions | object references only | notification center |
| Trace / Activity Stream | provenance et actions | semantic events | design/shared trace sources |
| Inspector / Context Bar | sélection, source, policy et return origin | projections | Design System |
| Timeline | Memory Timeline and cross-source correlation | memory events and quality | Timeline Engine |
| Object Linking | Case, Image, observations, Artifacts | typed relations | Object Linking Service |
| Export / Reporting | packages autorisés | selected results, restrictions | Export/Reporting Engine |
| Versioning / Collaboration | annotations, sessions, disputes | local interpretations | Shared services |
| Audit Hooks | reveal/copy/export and mutation events | semantic audit events | Trust / Shared |
| Permission-aware search | retrouver images et sessions autorisées | identifiers and metadata | Global Search |
| Recovery | interruption and context restoration | session snapshot references | Recovery pattern |

Aucune Shared Capability n’est redéfinie.
