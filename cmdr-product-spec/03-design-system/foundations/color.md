---
id: foundation-color
domain: 03-design-system
status: draft
owner: Design System Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-BRAND-003
  - REQ-BRAND-004
  - REQ-BRAND-006
  - REQ-UX-003
  - REQ-UX-005
---

# Couleur

## Quatre couches

1. **Marque** : valeurs CMDR et Command, sources uniques dans `02-brand/`.
2. **Produit** : aliases ; Investigate, Govern et Studio restent non résolus.
3. **Sémantique** : statut et interaction, indépendants de la marque.
4. **Visualisation** : séries et encodages, indépendants du produit et du statut.

## Règles

Un composant consomme `color.surface.*`, `color.text.*`, `color.status.*` ou un token composant. Moss n'est pas succès ; Ember n'est pas critique ; Juniper n'est pas succès. Une couleur est accompagnée d'un libellé, d'une icône ou d'une forme.

## Produit

Command peut résoudre ses aliases vers sa palette canonique. Investigate, Govern et Studio exposent seulement des slots `unresolved`. Settings utilise les neutres CMDR ; Endpoint adopte l'identité du workflow propriétaire.

## Validation

Chaque paire réelle est testée en clair/sombre pour texte, icône, bordure, focus, contrôle et graphique. Les contrastes de Phase 2 restent des preuves de marque et ne sont pas recopiés comme tokens.
