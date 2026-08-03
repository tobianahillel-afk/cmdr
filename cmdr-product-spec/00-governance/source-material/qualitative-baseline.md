---
id: qualitative-baseline
domain: 00-governance
status: draft
owner: QA and Traceability Lead
updated: 2026-08-03
source-of-truth: source-material
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-BRAND-001
  - REQ-UX-001
  - REQ-UX-010
  - REQ-OBJ-011
  - REQ-OBJ-012
---

# Qualitative Baseline

## Baseline repository Phase 0

781 Markdown ; 559 generic skeletons ; 61 Template-level screens ; 57 objects insufficiently formalized ; 461 `À compléter` ; 289 exact repeated placeholders.

## Résultats précédents

- Phase 1 : gouvernance et vision PASS.
- Phase 2 : marque PASS ; 88 conform, 25 partial, 6 absent, 3 contradictory.

## Mesures Phase 3

Le scope comprend tous les fichiers actifs initiaux de `03-design-system/` (89) et `04-experience-architecture/` (32), les sources Shared/Command nécessaires à Saved Views, les registres et les cinq fichiers écran legacy migrés.

| Mesure | Avant | Après |
|---|---:|---:|
| Fichiers `03-design-system/` | 89 | 138 |
| Fichiers actifs Draft `03-design-system/` | 89 | 73 |
| Fichiers Deprecated `03-design-system/` | 0 | 65 |
| Generic skeletons dans `03-design-system/` | 80 | 0 actifs |
| Fichiers `04-experience-architecture/` | 32 | 47 |
| Fichiers actifs Draft `04-experience-architecture/` | 32 | 29 |
| Fichiers Deprecated `04-experience-architecture/` | 0 | 18 |
| Generic skeletons dans `04-experience-architecture/` | 30 | 0 actifs |
| Placeholders actifs ciblés | 110 | 0 |
| Fichiers ciblés sans Requirement IDs | 113 | 0 |
| Fichiers ciblés sans statut | 0 | 0 |
| Primitive tokens documentés | 0 | 41 |
| Semantic token roles | 0 | 45 |
| Component token contracts initiaux | 0 | 24 |
| Aliases non résolus | 0 | 5 |
| Layout shells substantiels | 0 / 8 | 8 / 8 |
| Composants canoniques substantiels | 1 / 29 cibles | 29 / 29 |
| Patterns canoniques substantiels | 0 / 14 | 14 / 14 |
| Exemples documentaires | 0 | 5 |
| Sources actives concurrentes dans le scope | nombreuses | 0 connues |
| Palettes ouvertes consommées par token actif | 0 | 0 |
| Fichiers écran métier détaillés réécrits | 0 | 0 |
| Fichiers écran legacy transformés en pointeur | 0 | 5 |
| Code produit ajouté | 0 | 0 |

## Couverture

| État | Après Phase 2 | Après Phase 3 |
|---|---:|---:|
| conform | 88 | 99 |
| partial | 25 | 20 |
| absent | 6 | 3 |
| contradictory | 3 | 0 |
| total | 122 | 122 |

## Résolutions

- Generic Saved Views : Shared Capabilities ;
- Work Queue Saved Views : Command ;
- Saved Views UI : Design System ;
- Work Queue : un workspace et six vues ;
- Inspector : une seule source ;
- contexte : une source UX et un composant visuel ;
- tokens : une architecture à trois niveaux ;
- anciens noms : pointeurs Deprecated avec migration.

## Limites

Les écrans détaillés, fonctions de modules, objets, permissions, parcours, architecture technique et logiciel restent incomplets. Les questions OPEN-001..008 et OPEN-010..016 restent ouvertes, sauf OPEN-009 déjà résolue.
