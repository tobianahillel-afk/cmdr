---
id: foundation-focus-selection
domain: 03-design-system
status: draft
owner: Accessibility Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-002
  - REQ-UX-004
  - REQ-UX-005
---

# Focus et sélection

Focus = position d'interaction ; sélection = objet actif. Ils peuvent coexister et sont visuellement distincts.

- focus ring `2px` + offset `2px`, contraste non-textuel ≥3:1 ;
- sélection : fond + bordure/rail + `aria-selected` ;
- hover n'est jamais persistant ;
- focus n'est pas supprimé au clic clavier ;
- fermeture d'une surface restaure le déclencheur ou la ligne source ;
- sélection supprimée : focus revient au voisin logique et annonce le changement.

Dans une table, les flèches déplacent le focus selon le modèle choisi ; Space sélectionne, Enter ouvre, `]` ouvre l'Inspector. Une couleur seule ne suffit jamais.
