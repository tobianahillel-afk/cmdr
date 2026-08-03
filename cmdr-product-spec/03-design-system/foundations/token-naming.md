---
id: token-naming
domain: 03-design-system
status: draft
owner: Design System Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-BRAND-003
  - REQ-BRAND-004
  - REQ-UX-003
---

# Nommage des tokens

## Convention unique

`category.element.property.variant.state`

Les segments inutiles sont omis sans changer l'ordre. Noms en minuscules, séparés par des points ; aucun nom de composant dans un token sémantique.

Exemples : `color.text.primary`, `color.surface.canvas`, `space.300`, `radius.control`, `motion.duration.fast`, `component.button.primary.background.default`.

## Niveaux

- primitive : valeur brute ou alias vers marque ;
- semantic : intention fonctionnelle indépendante du composant ;
- component : décision locale d'un composant.

## Métadonnées obligatoires

`id`, `category`, `role`, `value-or-alias`, `source`, `status`, `modes`, `products`, `accessibility`, `allowed`, `forbidden`, `owner`, `updated`.

## Interdictions

Abréviations opaques, hex dans un composant, mélange kebab/dot, `primary-green`, nom de fournisseur, produit dans un statut, et alias circulaire.
