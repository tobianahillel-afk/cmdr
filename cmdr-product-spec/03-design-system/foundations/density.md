---
id: foundation-density
domain: 03-design-system
status: draft
owner: Design System Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-007
  - REQ-PROD-057
  - REQ-UX-004
---

# Densité

| Mode | Contrôle | Ligne | Padding principal | Usage |
|---|---:|---:|---:|---|
| compact | 28 px | 32 px | 8 px | queues, tables, workbenches |
| standard | 36 px | 40 px | 12 px | Case, Run, Settings |
| comfortable | 44 px | 48 px | 16 px | formulaires, lecture, touch |

La cible minimum interactive reste 44×44 CSS px lorsque nécessaire via zone de hit invisible. Densité n'altère pas libellés, statuts, permission ou contenu. Le défaut suit l'activité ; l'override utilisateur est permis. `OPEN-010` reste ouverte pour les préférences finales par rôle.
