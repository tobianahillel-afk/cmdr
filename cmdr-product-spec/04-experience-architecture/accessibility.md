---
id: experience-accessibility
domain: 04-experience-architecture
status: draft
owner: Accessibility Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-004
  - REQ-UX-005
  - REQ-UX-010
  - REQ-PROD-007
---

# Accessibilité de l’expérience

CMDR vise WCAG 2.2 AA. Tous les shells possèdent landmarks, skip links vers navigation locale, workspace, Inspector et console lorsqu'ils existent.

Le clavier suit l'ordre visuel ; chaque composant complexe documente ses flèches, Home/End, Enter/Space et Escape. Le focus est visible, distinct de la sélection et restauré après fermeture.

Zoom 200 %, reflow 400 %, text spacing et longues chaînes ne doivent pas masquer actions ou statuts. Tables, graphes, canvases, éditeurs, consoles et drag-and-drop ont une alternative structurée. Les mises à jour, erreurs, jobs et changements de statut utilisent des live regions proportionnées.

Timeout et expiration annoncent la durée, proposent extension et récupération sûre. Permission denied ne révèle aucune donnée protégée. Reduced motion supprime les transitions non nécessaires.

**Given** un Technical Workbench à 200 % de zoom et reduced motion,  
**When** l'utilisateur parcourt explorer, canvas alternatif, Inspector et console,  
**Then** chaque région est atteignable, nommée, non recouverte et aucune information n'est transmise uniquement par mouvement ou couleur.
