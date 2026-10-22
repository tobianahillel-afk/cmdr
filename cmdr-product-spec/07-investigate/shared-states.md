---
id: investigate-shared-states
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-005
  - REQ-UX-004
  - REQ-UX-005
---

# États partagés Investigate

Chaque capability configure Loading, Empty, Partial, Error, Offline, Permission denied et Stale sans recopier le contrat Design System.

## Règles métier transversales

- `partial` nomme toujours la source ou relation manquante ;
- `stale` expose source, date et conséquence ;
- `conflicted` préserve les versions et exige une résolution explicite ;
- `proposed` n’est jamais équivalent à `confirmed` ;
- `superseded` conserve le remplaçant et l’historique ;
- `restricted` ne divulgue pas le contenu protégé ;
- aucune absence de donnée n’est convertie en fait négatif certain.

Les états spécifiques restent dans chaque capability et les machines d’état finales restent Phase Objets.
