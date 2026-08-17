---
id: experience-architecture-readme
domain: 04-experience-architecture
status: draft
owner: UX Architecture Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-001
  - REQ-UX-002
  - REQ-UX-003
  - REQ-UX-004
  - REQ-UX-005
  - REQ-UX-006
  - REQ-UX-007
  - REQ-UX-008
  - REQ-UX-009
  - REQ-UX-010
---

# Architecture d’expérience CMDR

## Mission

Définir comment les personnes se déplacent, conservent leur contexte, organisent un workspace dense et utilisent les mêmes primitives entre Command, Investigate, Govern, CMDR Studio et Platform Settings. Ce domaine possède les règles d'expérience ; `03-design-system/` possède leur traduction visuelle et interactive.

## Statut et limites

Tous les documents restent `draft`. La Phase 3 fixe des décisions UX réversibles et testables, sans réécrire les 61 écrans, sans définir leurs colonnes métier et sans fermer les décisions de marque ou de typographie.

## Sources normatives

1. `../00-governance/source-material/ux-and-navigation-decisions.md`
2. `../01-product-vision/product-principles.md`
3. `../01-product-vision/product-boundaries.md`
4. `../02-brand/operational-editorial-modernism.md`
5. `../00-governance/terminology-rules.md`

## Ordre de lecture

1. `experience-principles.md`
2. `information-architecture.md`
3. `global-navigation.md`
4. `workspace-model.md`
5. `page-view-mode-filter-rules.md`
6. `context-preservation.md`
7. `history-and-back.md`
8. `progressive-disclosure.md`
9. `cognitive-load.md`
10. `cross-product-transitions.md`
11. `role-based-defaults.md`
12. `personalization.md`
13. `accessibility.md`

## Responsabilités

- structure produit, module, page, workspace, vue, mode et filtre ;
- Global Header, navigation primaire et locale ;
- modèle de contexte, deep links et retour ;
- règles de densité informationnelle ;
- transitions interproduits ;
- valeurs par défaut, personnalisation, onboarding et aide ;
- critères d'expérience accessibles et testables.

## Consommateurs

Design System, produits, parcours, écrans, Shared Capabilities, sécurité, qualité et futures implémentations.

## Décisions ouvertes

`OPEN-001`, `OPEN-002`, `OPEN-003`, `OPEN-004`, `OPEN-010` et `OPEN-016` restent ouvertes. Les densités définies ici sont des valeurs d'activité Draft ; elles ne résolvent pas les préférences finales par rôle.

## Contribution

Une règle nouvelle doit citer son Requirement ID, nommer son propriétaire, indiquer ce qui est préservé dans l'historique et fournir au moins un critère Given/When/Then. Une page ne peut être créée pour contourner une vue, un mode, un filtre, un Inspector, un drawer ou une modal.
