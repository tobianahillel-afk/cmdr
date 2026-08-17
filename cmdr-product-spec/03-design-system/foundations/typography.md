---
id: foundation-typography
domain: 03-design-system
status: draft
owner: Design System Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-BRAND-007
  - REQ-UX-004
  - REQ-PROD-051
---

# Système typographique Draft

`OPEN-004` reste ouverte. `font.family.ui` et `font.family.mono` sont des aliases non résolus ; aucune police n'est distribuée.

| Rôle | Taille / ligne | Poids Draft | Règle |
|---|---|---|---|
| Display | 40/48 | 600 | marketing rare |
| Page title | 28/36 | 600 | objectif précis |
| Section title | 20/28 | 600 | section majeure |
| Subsection title | 16/24 | 600 | groupe local |
| Body | 14/22 | 400 | lecture continue |
| Compact body | 13/18 | 400 | workbench/table |
| UI label | 13/18 | 600 | action et contrôle |
| Metadata | 12/16 | 400 | source, date, owner |
| Table | 13/18 | 400/600 | chiffres tabulaires |
| Badge | 11/16 | 600 | texte court |
| Code | 12/18 | 400 | monospace |
| Query | 12/18 | 400 | monospace |
| Identifier | 12/18 | 400 | monospace, wrap/copy |

Letter spacing : normal pour body ; `-0.01em` seulement grands titres ; `0.02em` labels courts. Capitales longues interdites. Hashes et chemins utilisent wrap opportuniste, copie explicite et valeur complète accessible. Les tables activent les chiffres tabulaires.

La pile candidate reste définie dans `02-brand/cmdr/typography.md`. La décision finale exige tests Windows/macOS/Linux, langues, performance et licence.
