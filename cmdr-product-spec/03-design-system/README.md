---
id: design-system-readme
domain: 03-design-system
status: draft
owner: Design System Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-002
  - REQ-UX-003
  - REQ-UX-004
  - REQ-UX-005
  - REQ-BRAND-003
  - REQ-BRAND-004
  - REQ-BRAND-007
  - REQ-BRAND-008
---

# Design System CMDR

## Mission

Transformer Operational Editorial Modernism et l'architecture d'expérience en fondations, tokens, shells, composants et patterns communs. Le Design System ne possède ni objets métier, ni permissions, ni palettes de marque.

## Statut

Draft. Les métriques, valeurs sémantiques et dimensions de Phase 3 sont des décisions réversibles. `OPEN-001` à `OPEN-004`, `OPEN-010` et `OPEN-016` restent ouvertes.

## Sources normatives

- `../02-brand/` pour marque et palettes ;
- `../04-experience-architecture/` pour structure et navigation ;
- `../00-governance/ownership-register.md` pour propriétaires ;
- `../14-security-permissions-and-trust/` pour autorité ;
- `../15-content-and-language/` pour contenu.

## Ordre de lecture

1. `foundations/tokens.md`
2. `foundations/token-naming.md`
3. `foundations/color.md`
4. `foundations/theme-contract.md`
5. `foundations/typography.md`
6. `foundations/spacing.md`, `grid.md`, `density.md`, `responsive.md`
7. `layouts/global-shell.md` puis les shells d'activité
8. `components/`
9. `patterns/`
10. `examples/`

## Responsabilités

- trois niveaux de tokens ;
- thèmes clair/sombre et sémantique ;
- typographie sans choix final de famille ;
- dimensions, grilles, focus, mouvement et accessibilité ;
- huit shells ;
- composants et patterns partagés ;
- adaptations produit sans Design System concurrent.

## Règle de contribution

Un composant consomme des tokens sémantiques ou composants, jamais un hex de marque. Toute nouvelle variante démontre un besoin fonctionnel, un modèle clavier, les six états applicables, un comportement responsive et un critère Given/When/Then.
