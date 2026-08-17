---
id: investigate-memory-forensics-user-questions
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-INV-001
  - REQ-PROD-014
  - REQ-PROD-020
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# User questions

- Quelle Memory Image suis-je en train d’analyser et de quel Endpoint provient-elle ?
- L’image est-elle complète, intègre, exploitable, partielle, corrupted ou disputed ?
- Quel est le contexte d’acquisition, l’autorisation et la custody ?
- Quelle plateforme est déclarée, détectée ou seulement proposée ?
- Quel profil est sélectionné, avec quelle confiance et quelles contradictions ?
- Quels processus, threads, régions, modules, drivers, handles, objets et connexions sont reconstruits ?
- Quelles vues divergent et quel effet la partialité de l’image a-t-elle ?
- Quelles données sensibles candidates existent sans révéler leur valeur ?
- Quelles anomalies ou incohérences restent candidates ?
- Quels timestamps sont observés, reconstruits, estimés ou absents ?
- Quels Derived Artifacts ont été extraits et avec quelle parenté ?
- L’analyse est-elle reproductible et quelles préconditions manquent ?
- Quels éléments peuvent être transmis à Evidence, Finding ou future Detection Engineering ?
- Comment revenir au Case ou à la vue source sans perdre la sélection ?
