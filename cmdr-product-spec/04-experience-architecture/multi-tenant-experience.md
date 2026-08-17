---
id: experience-multi-tenant
domain: 04-experience-architecture
status: draft
owner: UX Architecture Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-008
  - REQ-SEC-001
  - REQ-UX-006
---
# Expérience multi-tenant


Le tenant est visible dans le Global Header ou le Context Bar avant toute donnée. Une vue agrégée affiche le tenant par objet et interdit les mutations multi-tenant.

Changer de tenant exige confirmation si contexte ou travail actif ; filtres, sélection, panels et drafts incompatibles sont nettoyés. Aucun fallback vers le premier tenant. Les liens sont tenant-scoped et réévalués.

**Given** un draft dans Tenant A, **When** l'utilisateur choisit Tenant B, **Then** il doit enregistrer/abandonner/rester, le contexte A est supprimé et aucune donnée A n'apparaît dans B.
