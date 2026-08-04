---
id: investigate-action-classification
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-003
  - REQ-PROD-004
  - REQ-SEC-001
  - REQ-SEC-002
open_decisions:
  - OPEN-013
---

# Action classification — Investigate

| Classe | Investigate peut | Exemples 4B.1 | Autorité |
|---:|---|---|---|
| 0 | observer, rechercher, inspecter, filtrer, pivoter, comparer | Event Search, read Evidence, inspect Finding | produit/source selon permission |
| 1 | préparer ou demander une collecte | collection request draft | exécution 4B.2 et Govern selon risque |
| 2 | modifier réversiblement un objet ou lien Investigate | Case, Hypothesis, relation, Artifact, qualification Evidence, Finding draft | `OPEN-013` reste ouverte |
| 3 | préparer une demande de containment | Action Request fondée sur Finding/Evidence | Govern requis |
| 4 | fournir preuve et contexte pour irréversible/destructif | cible, impact, alternatives, rollback | Govern obligatoire |

## Invariants

- la classe décrit l’effet, pas le bouton ;
- une action classe 3 ou 4 n’est jamais exécutée dans Investigate ;
- une automatisation ne s’accorde pas sa propre permission ;
- bulk linking reste classe 2 et permission-aware ;
- la création d’une Evidence est réversible au niveau de l’usage, mais sa trace historique n’est pas supprimée ;
- `OPEN-013` n’est pas fermée.
