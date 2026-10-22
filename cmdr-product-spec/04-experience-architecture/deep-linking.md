---
id: experience-deep-linking
domain: 04-experience-architecture
status: draft
owner: UX Architecture Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-008
  - REQ-UX-006
  - REQ-UX-007
---
# Deep linking


Un deep link encode route, tenant, environnement, objet, vue/mode/filtres sûrs, sélection et return origin lorsque pertinent. Secrets, payloads, tokens et valeurs sensibles sont exclus.

À l'ouverture, session, tenant et permissions sont vérifiés avant chargement. Un lien expiré ou interdit explique l'état sans révéler l'objet et propose retour sûr. Les liens de contenu sensible peuvent être expirables et audités.

**Given** un lien vers une Evidence après expiration de session, **When** l'utilisateur se reconnecte, **Then** la route est restaurée après autorisation, ou une Permission denied sûre est affichée.
