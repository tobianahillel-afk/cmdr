---
id: foundation-accessibility
domain: 03-design-system
status: draft
owner: Accessibility Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-004
  - REQ-UX-005
  - REQ-UX-010
---

# Accessibilité du Design System

Objectif WCAG 2.2 AA. Obligations : landmarks, skip links, focus visible, ordre logique, noms accessibles, messages associés, status announcements, zoom/reflow, text spacing, reduced motion et contraste.

## Composants complexes

- table/grid : modèle clavier documenté et alternative HTML ;
- graph/canvas : liste ou table synchronisée ;
- code editor/console : mode lecture, recherche, line numbers et sortie annoncée ;
- drag-and-drop : actions déplacer avant/après/vers ;
- resizer : séparateur focusable avec flèches et Home/End ;
- modal/drawer : focus initial, piège, Escape et restauration ;
- timeout : avertissement, prolongation et récupération.

Les identifiants longs se replient sans perte, restent copiables et lisibles. Permission denied ne révèle pas le nom ou la valeur d'un objet interdit. Les six états Loading, Empty, Partial, Error, Offline et Permission denied sont inclus dans chaque composant pertinent.
