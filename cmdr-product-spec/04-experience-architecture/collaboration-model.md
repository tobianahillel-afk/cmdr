---
id: experience-collaboration
domain: 04-experience-architecture
status: draft
owner: UX Architecture Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-006
  - REQ-PROD-003
---
# Collaboration


Présence, commentaires, mentions et conflits sont des projections de la capacité partagée. La présence n'accorde aucune permission. Un commentaire reste distinct d'un fait, d'une Evidence, d'un Finding ou d'une Decision.

Les modifications concurrentes affichent version source, version locale et options `Comparer`, `Fusionner`, `Recharger`; aucun last-write-wins silencieux. Attribution, timestamp, modification et suppression restent auditables.

**Given** deux utilisateurs modifiant une note, **When** un conflit est détecté, **Then** les deux versions restent accessibles, le focus ne saute pas, et aucune Evidence n'est modifiée par la résolution d'un commentaire.
