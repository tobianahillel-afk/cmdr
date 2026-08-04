---
id: CAP-CMD-201
title: Service Context
product: command
module: risk-and-coverage
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids: [REQ-PROD-005, REQ-PROD-013, REQ-PROD-021, REQ-PROD-032]
open_decisions: [OPEN-013]
source-of-truth: canonical
---
# CAP-CMD-201 — Service Context
## 1. Définition
Consomme le Business Service Catalog pour relier un Incident/Task au service, owner, criticité, environnement et dépendances.
## 2. Problème utilisateur
La coordination technique ne suffit pas sans service métier ni owner. Principal : Incident Commander ; secondaires : Business Owner, SOC L2, Service Delivery Manager. Sans capacité, l’impact et l’escalade sont incomplets.
## 3. Objectifs
Présenter identité/owner/criticité/dépendances ; lier work items sans dupliquer le catalogue ; montrer source/fraîcheur ; signaler les lacunes au propriétaire.
## 4. Non-objectifs
Ne pas devenir CMDB, créer un Service local, modifier les dépendances ou inférer un owner.
## 5. Propriétaire
Command possède la relation opérationnelle ; Shared Business Service Catalog possède Service.
## 6. Utilisateurs
Incident Commander principal ; Business Owner, SOC L2, Service Delivery Manager secondaires. Scope et permissions obligatoires.
## 7. Conditions d’entrée
Incident/Task ou contexte risque ; référence explicite ou candidats à confirmer ; permission catalogue.
## 8. Entrées fonctionnelles
| Entrée | Source | Type | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Service reference | Incident/Task/utilisateur | ID/candidat | oui ou unknown | courante | impact conservé sans service confirmé |
| Service metadata | Catalog | owner/criticality/dependencies/env | non | déclarée | Partial |
| Operational links | Linking Service | incidents/exposures/coverage | non | courante | relation indisponible |
## 9. Objets lus
| Objet | Owner | Projection | Droit local |
|---|---|---|---|
| Service | Shared Catalog | identity/owner/criticality/dependencies | projection |
| Incident / Task | Command | service links/impact | lecture/modification relation |
| Exposure / coverage | sources propriétaires | relation service | projection |
## 10. Objets créés ou modifiés
| Objet | Opération | Owner | Règle |
|---|---|---|---|
| Incident / Task | ajouter/corriger lien Service | Command | classe 2 ; Service inchangé |
| Task | correction de données | Command | follow-up catalogue |
## 11. Fonctionnalités
Chercher/confirmer service ; afficher owner/criticité/env/dépendances ; montrer incidents/expositions ; marquer assumed/confirmed ; créer correction.
## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter | lecteur | Service projection | 0 | permission | contexte | non |
| Lier | coordinateur | Incident/Task | 2 | confirmed/assumed | relation sourcée | OPEN-013 |
| Confirmer/contester | Business Owner/coordinateur | lien | 2 | source visible | statut lien | OPEN-013 |
| Créer correction | coordinateur | Task | 2 | lacune | Task | non |
## 13. Automatisation et IA
| Fonction | Humain | Règle | Moteur | Workflow | Agent | Govern | Sans IA |
|---|---|---|---|---|---|---|---|
| Candidats service | confirmation | matching rule | deterministic match | possible | proposition | non | recherche manuelle |
| Lien effectif | décision | policy possible | validation | possible | jamais silencieux | OPEN-013 | action manuelle |
## 14. États fonctionnels
`confirmed`, `assumed`, `partial`, `stale`, `owner-unknown`, `dependency-unknown`, `not-linked`.
## 15. États d’interface
Partial/Stale nomment champs/source ; Offline bloque mutation ; Permission denied ne révèle pas Service. Rendu DS.
## 16. Sorties
| Sortie | Objet/événement | Consommateur | Garantie |
|---|---|---|---|
| Service context | projection | Mission Control/Incident/Risk | source/freshness/status |
| Service link | relation typée | work item/catalog consumers | aucun ownership transfer |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte | Retour |
|---|---|---|---|---|
| Service Context | correction | Catalog owner | service ref/field/evidence | work item |
| Service Context | impact/autorité | Govern | Incident/service/impact/owner | Command |
## 18. Dépendances
Catalog, Linking Service, CAP-CMD-204/106/107.
## 19. Source de vérité
Service et metadata : Catalog ; relation/Task : Command ; projections externes sourcées.
## 20. Provenance et audit
Acteur, candidat, source, statut assumed/confirmed, before/after, tenant et correlation ID.
## 21. Permissions fonctionnelles
Lecture Command/catalog ; manage Incident/Task pour relation. Atomisation reportée.
## 22. Limites et erreurs
Service absent, multiple, stale, inaccessible ou conflit ne doit jamais produire un owner/impact inventé.
## 23. Métriques
Incidents avec service confirmed/assumed ; services sans owner/criticité ; délai correction. Pas de cible définitive.
## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire uniquement.
## 25. Critères d’acceptation
**Given** deux services candidats, **When** le coordinateur confirme l’un, **Then** source et statut confirmed sont audités sans modifier le catalogue.

**Given** aucun modèle, **When** la recherche est utilisée, **Then** catalogue, règles et sélection manuelle suffisent.

**Given** catalog inaccessible, **When** l’Incident s’ouvre, **Then** le lien est Partial/unknown et les données Command restent utilisables.
## 26. Questions ouvertes
Objet Service canonique ou catalogue fonctionnel ? Qui confirme service-impact ? — requirements ci-dessus. `OPEN-013` reste ouverte.
## 27. Consommateurs documentaires
Mission Control, Incident Detail, Risk, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions ultérieures.
