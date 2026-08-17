---
id: 06-command-product-definition
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-003
  - REQ-PROD-008
  - REQ-PROD-010
  - REQ-PROD-013
  - REQ-PROD-021
  - REQ-OBJ-001
  - REQ-OBJ-012
  - REQ-UX-008
  - REQ-UX-009
open_decisions:
  - OPEN-013
---
# Product boundary and functional definition — Command

## Définition

Command est le produit de coordination opérationnelle de CMDR. Il reçoit ou crée le travail Command, maintient owner, priorité, état de travail, impact, SLA, prochaine action et relations, puis transmet le contexte aux produits propriétaires lorsque investigation, autorité ou exécution sont nécessaires.

## Résultats utilisateurs

- savoir ce qui compte maintenant et pourquoi ;
- savoir qui possède chaque prochaine action ;
- travailler dans une Work Queue unique ;
- conserver le contexte entre Command, Investigate, Govern et Studio ;
- distinguer faits, projections, propositions et décisions ;
- préparer une relève ou une escalade sans perte d’information.

## Ownership

| Concept | Owner | Droit Command |
|---|---|---|
| Incident | Command | créer, coordonner, modifier selon permission |
| Task opérationnelle | Command | créer, coordonner, modifier selon permission |
| Case / Evidence / Finding | Investigate | lire, relier, naviguer, demander création/action |
| Action Request / Decision / Response Run / Result | Govern | fournir contexte, lire projection, suivre |
| Workflow / Automation Run | Studio | lancer si autorisé, inspecter, suivre |
| Endpoint Agent Fleet | Platform Settings | lire health/capability projection |
| Shared engines | Shared Capabilities | consommer, jamais dupliquer |

## Frontières

Command ne fournit pas terminal, forensic, debugger, sandbox, Detection Engineering, approval, Response Run execution, agent builder ou fleet administration. Une projection locale ne devient pas une copie propriétaire.

## Actions à effet

- classe 0 : observation, lecture, filtre, navigation ;
- classe 1 : une demande de collecte part vers Investigate ;
- classe 2 : mutations Command réversibles, sous permissions et OPEN-013 ;
- classe 3 : Command prépare/suit la demande, Govern décide et exécute ;
- classe 4 : Command apporte le contexte, Govern gère l’autorité et la réponse.

## Delivery

La cible produit est native pour les capabilities cœur. Le mode courant reste `planned`; un document ne prouve ni implementation, ni endpoint, ni moteur.

## Non-objectifs Phase 4A

- écrire les 27 sections détaillées des écrans ;
- définir toutes les colonnes ;
- formaliser les schémas complets des objets ;
- finaliser permissions, API, protocoles ou stockage ;
- définir les capabilities Investigate, Govern, Studio, Settings ou Endpoint Agent.

## Critère

**Given** une mutation ayant un effet hors Command, **When** l’utilisateur la demande, **Then** Command prépare un package sourcé et transfère au produit propriétaire ; aucune autorité ou exécution n’est simulée localement.
