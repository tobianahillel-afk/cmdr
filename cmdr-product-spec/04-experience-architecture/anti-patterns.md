---
id: experience-anti-patterns
domain: 04-experience-architecture
status: draft
owner: UX Architecture Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-001
  - REQ-UX-002
  - REQ-UX-004
  - REQ-AI-001
---

# Anti-patterns d’expérience

| Anti-pattern | Problème | Alternative |
|---|---|---|
| page par filtre | navigation et états dupliqués | View + URL |
| dashboard de cartes | faible scan et densité artificielle | sections, table, métrique justifiée |
| deux Inspectors | sélection et focus ambigus | Inspector canonique |
| drawer empilé | retour imprévisible | fermer/remplacer ou workspace |
| modal longue | focus prisonnier | page/workspace |
| chat central | dépendance IA et contexte opaque | assistance contextuelle |
| refresh complet | perte de sélection | patch local et fraîcheur |
| changement de tenant silencieux | risque de fuite | confirmation et nettoyage |
| Back vers accueil | perte de travail | historique réel |
| panneau sans limites | canvas écrasé | min/max et collapse |
| couleur seule | inaccessible | label, forme, icône |
| mobile simulant le workbench complet | interaction impraticable | lecture/stack explicite |
