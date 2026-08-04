---
id: qualitative-baseline
domain: 00-governance
status: draft
owner: QA and Traceability Lead
updated: 2026-08-04
source-of-truth: source-material
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-013
  - REQ-UX-010
  - REQ-OBJ-001
  - REQ-OBJ-012
---

# Qualitative Baseline

## Baseline repository Phase 0

781 Markdown ; 559 generic skeletons ; 61 Template-level screens ; 57 objects insufficiently formalized ; 461 `À compléter` ; 289 exact repeated placeholders.

## Résultats précédents

- Phase 1 : gouvernance et vision PASS.
- Phase 2 : marque PASS.
- Phase 3 : Experience Architecture et Design System PASS.
- Couverture avant Phase 4A : 99 conform, 20 partial, 3 absent, 0 contradictory.

## Périmètre lu Phase 4A

- 142 fichiers accessibles lus intégralement ;
- 27 fichiers Command actuels lus : 22 actifs et 5 deprecated ;
- 10 écrans Command actifs lus sans réécriture ;
- 5 anciens écrans Work Queue relus et laissés deprecated ;
- 4 chemins d’objets demandés confirmés absents : `service.md`, `exposure.md`, `report.md`, `audit-record.md`.

## Mesures Command avant/après

| Mesure | Avant Phase 4A | Après Phase 4A |
|---|---:|---:|
| Fichiers Command totaux | 27 | 69 |
| Fichiers Command actifs | 22 | 61 |
| Fichiers Command deprecated | 5 | 8 |
| Fichiers Command actifs génériques ou Template-level | 19 | 10 |
| Placeholders `À compléter` actifs dans Command | 19 | 10 |
| Placeholders dans les documents fonctionnels traités hors écrans | 9 | 0 |
| Écrans actifs contenant encore le placeholder reporté | 10 | 10 |
| Capability IDs Command | 0 | 27 |
| Capabilities enregistrées | 0 | 27 |
| Capabilities sans owner | N/A | 0 |
| Capabilities sans utilisateur | N/A | 0 |
| Capabilities sans entrée | N/A | 0 |
| Capabilities sans sortie | N/A | 0 |
| Capabilities sans objet | N/A | 0 |
| Capabilities sans action classée | N/A | 0 |
| Capabilities sans alternative non-IA | N/A | 0 |
| Capabilities sans critères Given/When/Then | N/A | 0 |
| Delivery mode courant `native` | 0 | 0 |
| Delivery mode courant `integrated` | 0 | 0 |
| Delivery mode courant `temporary-integration` | 0 | 0 |
| Delivery mode courant `planned` | 0 | 27 |
| Delivery status `defined` | 0 | 26 |
| Delivery status `proposed` | 0 | 1 |
| Delivery status `out-of-scope` | 0 | 0 |
| Dépendances Phase 4A structurées | 0 | 12 |
| Entrées totales dans le dependency register | 15 | 27 |
| Sources actives concurrentes connues dans le scope | 2 | 0 |
| Écrans détaillés réécrits | 0 | 0 |
| Écrans modifiés | 0 | 0 |
| Objets modifiés | 0 | 1, Task ownership uniquement |
| Sources de permissions atomiques modifiées | 0 | 0 |
| Code produit ajouté | 0 | 0 |
| APIs ou protocoles créés | 0 | 0 |
| Fichiers de police ajoutés | 0 | 0 |
| Nouveaux Requirement IDs | 0 | 0 |
| Questions ouvertes | 15 | 15 |
| Nouveaux OPEN | 0 | 0 |
| Liens locaux introduits cassés | N/A | 0 |
| Fichiers ciblés vides | N/A | 0 |

Les dix fichiers actifs encore génériques sont exclusivement les dix écrans Command existants, volontairement reportés à la phase Écrans. Aucun document de capability, module, registre ou cadre Phase 4A ne conserve un placeholder générique.

## Couverture des exigences

| État | Avant Phase 4A | Après Phase 4A |
|---|---:|---:|
| conform | 99 | 99 |
| partial | 20 | 20 |
| absent | 3 | 3 |
| contradictory | 0 | 0 |
| total | 122 | 122 |

La sous-phase ajoute des preuves Command détaillées sans promouvoir artificiellement les exigences globales qui dépendent encore d’Investigate, Govern, Studio, Settings, Endpoint Agent, Parcours, Objets, Permissions, Technique ou Implémentation.

## Résolutions Phase 4A

- modèle canonique des capabilities à 27 sections ;
- 27 Capability IDs immuables Command ;
- Capability Register exploitable ;
- consolidation `Risk and Coverage` ;
- Customers and Delivery explicitement proposed et deployment-dependent ;
- Task opérationnelle alignée sur Command ; Shared Task Inbox reste consommateur ;
- Work Queue unique et six vues inchangées ;
- actions classées 0–4 ;
- alternatives sans IA pour chaque capability ;
- ownership interproduit explicite ;
- transitions Command vers Investigate, Govern et Studio spécifiées fonctionnellement.

## Limites

Les écrans, parcours, objets complets, permissions atomiques, contrats techniques et logiciel restent incomplets. La Phase 4B n’est pas commencée.
