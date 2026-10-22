---
id: 07-investigate-navigation
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-008
  - REQ-UX-001
  - REQ-UX-006
  - REQ-UX-007
---

# Navigation Investigate

## Navigation produit

La navigation locale suit les six modules canoniques, mais seules les activités disponibles et autorisées sont exposées. Une capability planned ne crée pas automatiquement une destination active.

## Entrées principales

- Command Incident → Case intake ou Case existant ;
- Signal → triage, recherche ou Case ;
- Global Search → objet source ;
- Case Queue → Case Workspace ;
- Case → Event Search, Evidence, Finding, Timeline ou Report ;
- Settings Endpoint projection → inspection ou demande de collecte future.

## Retours

Chaque transition transmet un `return origin`. Back restaure la vraie vue précédente, y compris filtres, période, sélection, scroll, tabs, Inspector et Query draft lorsque cela est sûr.

## Règles

- aucune redirection silencieuse vers un Case différent ;
- une destination interdite conserve le workspace source ;
- un lien profond réévalue les permissions ;
- une transition vers Govern ne présente jamais une Decision déjà créée ;
- une transition vers 4B.2 ou 4B.3 reste documentaire tant que la capability n’est pas traitée.
