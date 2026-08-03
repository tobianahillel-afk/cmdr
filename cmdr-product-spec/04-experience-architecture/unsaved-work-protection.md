---
id: experience-unsaved-work
domain: 04-experience-architecture
status: draft
owner: UX Architecture Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-007
  - REQ-UX-009
---
# Protection du travail non enregistré


Autosave affiche `Saving`, `Saved`, `Offline queued`, `Conflict` ou `Failed`. Il ne promet jamais un enregistrement avant confirmation du serveur.

Avant navigation destructive : enregistrer, abandonner ou rester. Après crash/session expirée, une récupération scoped au tenant propose diff et restauration. Les secrets temporaires ne sont pas persistés.

**Given** un workflow modifié et offline, **When** l'utilisateur change de produit, **Then** le draft est conservé localement de façon sûre ou la navigation est bloquée avec explication.
