---
id: experience-screen-states
domain: 04-experience-architecture
status: draft
owner: UX Architecture Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-004
  - REQ-UX-005
  - REQ-UX-010
---
# États obligatoires


Tout écran documente Loading, Empty, Partial, Error, Offline et Permission denied. Chaque état précise données conservées, action sûre, focus initial, annonce, responsive et audit/correlation ID si applicable.

Un état non applicable doit être justifié par le comportement réel ; un texte générique n'est pas une preuve. Partial ne devient pas success et Offline n'autorise aucune mutation non garantie.
