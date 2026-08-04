---
id: CAP-CMD-401
title: Customers and Delivery Context
product: command
module: customers-and-delivery
owner: Command Product Lead
status: draft
delivery_status: proposed
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids: [REQ-PROD-053, REQ-PROD-013, REQ-PROD-019, REQ-PROD-033]
open_decisions: [OPEN-006, OPEN-013]
source-of-truth: proposal
---
# CAP-CMD-401 — Customers and Delivery Context
## 1. Définition
Proposition deployment-dependent présentant contexte client, portefeuille, engagements, delivery et reporting depuis Command/Reporting Engine sans imposer un modèle MSP/MSSP.
## 2. Problème utilisateur
Certains déploiements coordonnent plusieurs clients/engagements ; d’autres sont internes. Principal : Service Delivery Manager ; secondaires : Customer Success, Incident Commander, Business Owner, Internal Security Lead.
## 3. Objectifs
Réutiliser les fonctions indépendantes du modèle économique ; documenter internal/enterprise/MSP ; consommer Reporting Engine/métriques ; rendre contractual features explicitement activées.
## 4. Non-objectifs
Ne pas imposer multi-client, billing, portail, dupliquer Reporting Engine ou fermer `OPEN-006`.
## 5. Propriétaire
Command possède seulement le contexte local et les Tasks ; Reporting Engine et sources client/contrat gardent ownership.
## 6. Utilisateurs
Service Delivery Manager principal ; rôles secondaires selon deployment et permissions.
## 7. Conditions d’entrée
Deployment model déclaré, capability activée, audience/scope tenant-client, permissions, `OPEN-006` ouverte.
## 8. Entrées fonctionnelles
| Entrée | Source | Type | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Deployment model | architecture/deployment owner | internal/enterprise/MSP | oui | courante | disabled/proposed |
| Command projections | Incident/Task/Result/Readiness | delivery context | non | visible | partial |
| Engagement context | contract/customer source | SLA/obligations/audience | non | effective version | aucune obligation inférée |
## 9. Objets lus
Incident/Task (Command), Result (Govern), Report (Reporting Engine), Customer/engagement context (source deployment).
## 10. Objets créés ou modifiés
Report request/draft via Reporting Engine ; delivery follow-up Task (Command, C2) ; aucune mutation source client/contrat.
## 11. Fonctionnalités
Internal reporting sans Customer ; portfolio/engagement optionnels ; governed report request avec citations ; Tasks delivery ; séparer SLA opérationnel/contractuel.
## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter | autorisé | projections | 0 | enabled | scope visible | non |
| Créer report draft | autorisé | Report | 2 | audience/scope/template | Shared draft | workflow report |
| Créer follow-up | delivery manager | Task | 2 | obligation/gap | Task | OPEN-013 |
| Publier/exporter | report role | Report | 2 | review/permission | Shared result | pas autorité Command |
## 13. Automatisation et IA
Moteurs agrègent citations/snapshots ; workflows gèrent review ; IA peut rédiger un brouillon attribué. Sans IA : reporting déterministe et édition humaine.
## 14. États fonctionnels
`proposed`, `disabled-for-deployment`, `configured`, `partial`, `data-stale`, `audience-restricted`, `suspended`.
## 15. États d’interface
Disabled explique deployment ; Partial/Stale montrent sources ; Permission denied ne révèle pas client/contrat. Rendu DS.
## 16. Sorties
Delivery context vers rôles autorisés ; Report draft/request vers Reporting Engine ; Delivery Task vers Work Queue.
## 17. Transitions
Report creation vers Reporting Engine avec tenant/client scope/citations/audience ; obligation vers Work Queue avec engagement ref.
## 18. Dépendances
`OPEN-006`, Reporting, Metrics, Export, Catalog, CAP-CMD-105/303.
## 19. Source de vérité
Command objects, Govern Result, Shared Report et deployment customer/contract sources restent séparés.
## 20. Provenance et audit
Deployment model, audience, citations, snapshot, redaction, reviewer, publication/export et correlation ID.
## 21. Permissions fonctionnelles
Command read ; Shared report create/review/publish/export ; customer/contract source permissions ; tenant isolation.
## 22. Limites et erreurs
Module disabled, scope client ambigu, contractual data absent/stale, audience refusée ou cross-tenant risk empêchent toute promesse/donnée implicite.
## 23. Métriques
Déploiements utilisant generic reporting vs client context ; reports avec citations/snapshots ; Tasks liées à une source d’engagement.
## 24. Classification de livraison
`proposed` / `planned`, deployment-dependent. Aucune promotion ni inclusion universelle sans décision `OPEN-006` et preuve.
## 25. Critères d’acceptation
**Given** déploiement interne sans Customer, **When** report est créé, **Then** Reporting Engine fonctionne et le module client reste disabled/proposed.

**Given** MSP configuré, **When** delivery context est ouvert, **Then** tenant/client/audience/source sont explicites et isolés.

**Given** aucun modèle, **When** report est préparé, **Then** citations, templates et édition humaine suffisent.
## 26. Questions ouvertes
Quels scénarios activent le module ? Concepts génériques vs économiques ? SLA contractuel distinct ? — requirements ci-dessus. `OPEN-006` et `OPEN-013` restent ouvertes.
## 27. Consommateurs documentaires
Customer Overview seulement si retenu, Reporting, Work Queue, parcours Phase 5, écrans Phase 6 conditionnels, objets/permissions ultérieurs.
