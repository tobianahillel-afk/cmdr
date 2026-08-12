---
id: CAP-SET-002
title: Environment Administrative Lifecycle and Tenant Scope
product: platform-settings
module: tenants-and-environments
owner: Platform Settings Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-12
requirement_ids: [REQ-PROD-006, REQ-PROD-008, REQ-SEC-001]
open_decisions: [OPEN-013]
source-of-truth: canonical
---
# CAP-SET-002 — Environment Administrative Lifecycle and Tenant Scope
## 1. Définition
Capability Settings qui définit l'administration d'un Environment dans le scope obligatoire d'un Tenant, avec type, labels, lifecycle et isolation.
## 2. Problème utilisateur
Un administrateur doit créer, classifier, administrer et retirer un Environment sans confondre Environment, Tenant, Deployment, Endpoint ou health.
## 3. Objectifs
Définir une responsabilité administrative autonome, permission-aware, tenant-scoped et auditable, sans redéfinir les objets, permissions ou mécanismes partagés.
## 4. Non-objectifs
Ne pas définir d'API, protocole, schéma physique, moteur générique de configuration, nouvelle permission, nouvel Screen ID, implémentation, ni autorité Govern parallèle.
## 5. Propriétaire
Platform Settings Product Lead est l'unique capability owner. Les mécanismes Security, Shared, Design System et Experience Architecture restent des dépendances possédées par leurs domaines.
## 6. Utilisateurs
Platform Administrator principal; Security Administrator, SOC Lead et consommateurs CMDR autorisés secondaires selon permission.
## 7. Conditions d’entrée
Principal authentifié, Tenant résolu avant l'objet, contexte Environment compatible lorsque requis, permissions serveur courantes et sources canoniques disponibles.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Tenant parent | Platform Settings | référence canonique | oui | courant | création impossible |
| Environment/intent | utilisateur | ref ou création | oui | courant | aucune mutation |
| autorisation | Security | permission + tenant/env scope | oui | courant | refus explicite |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Environment | Platform Settings | id, tenant-id, version, type/state/labels | read/manage |
| Tenant | Platform Settings | tenant reference et scope | read |
| Permission context | Security | tenant/environment ABAC | consume |
## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Environment | lifecycle/metadata mutation | Platform Settings | tenant obligatoire; transitions invalides refusées |
| Administrative audit event | emit/reference | Platform Settings / audit consumer | provenance seulement |
## 11. Fonctionnalités
Types fonctionnels Production/Preproduction/Test/Development/Lab, labels, tenant scope, isolation, promotion context comme contexte uniquement, lifecycle provisioning/active/frozen/retired.
## 12. Actions utilisateur
Class 0: inspecter. Class 1: valider type/scope/transition. Class 2: créer/provisionner, modifier labels, freeze/reactiver/retire lorsqu'autorisé. Aucun déploiement ou promotion automatique.
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| expliquer type/état | oui | oui | oui | oui | champs canoniques |
| suggérer configuration contextuelle | oui | oui | oui | oui | types/labels/actions déterministes |
## 14. États fonctionnels
`provisioning`, `active`, `frozen`, `retired`; le type d'Environment est distinct de son état. `retired` n'est pas une preuve d'effacement physique.
## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied et Stale exposent la cause, la fraîcheur et les conséquences sans inventer de données ni élargir le scope.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| environment projection | Environment reference | Settings/consumers | tenant-id + type/state/version |
| scope validation outcome | résultat fonctionnel | administrateur | autorisé/refusé + raison |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| SET-TEN-001 | create/configure/freeze/retire | Environment owner surface | tenant/env/ref/justification | même contexte |
| consumer product | open environment context | destination authorized | tenant/env references | return-origin |
## 18. Dépendances
Tenant object, Security tenant isolation/ABAC, Experience context preservation, Administrative Audit.
## 19. Source de vérité
`05-domain-model/objects/environment.md` et `environment-management.md`; Tenant reste une référence owner-owned.
## 20. Provenance et audit
Chaque création et transition conserve actor, tenant, environment, action, résultat, justification, version et correlation id.
## 21. Permissions fonctionnelles
`perm.platform-settings.environment.read` et `perm.platform-settings.environment.manage`; aucun nouvel ID n'est permis. Tenant scope et ABAC s'ajoutent à la permission ressource.
## 22. Limites et erreurs
Ne pas inventer cardinalité exacte, nesting, shared Environment cross-tenant, migration entre tenants ou promotion automatique. Si un Screen/permission nouveau devient nécessaire: BLOCKED.
## 23. Métriques
Environments par état/type, transitions refusées, scope conflicts, stale context; aucune cible KPI/SLO.
## 24. Classification de livraison
`defined / planned`; aucune disponibilité runtime, implémentation, intégration ou support plateforme n'est revendiqué.
## 25. Critères d’acceptation
**Given** Tenant autorisé et type `Production`, **When** un Environment est provisionné, **Then** son tenant-id et sa provenance sont obligatoires.  
**Given** Environment `retired`, **When** une mutation incompatible est demandée, **Then** la transition est refusée côté serveur.  
**Given** un Environment d'un autre Tenant, **When** il est projeté dans le contexte courant, **Then** l'accès est refusé et aucun fallback tenant n'est appliqué.
## 26. Questions ouvertes
Cardinalité, nesting et migration inter-tenant restent non spécifiés tant qu'une source canonique ne les tranche pas; aucune nouvelle OPEN ici.
## 27. Consommateurs documentaires
Tenants & Environments, Security, Experience Architecture, Command, Investigate, Govern, Studio, Endpoint, Shared, Quality et Roadmap.
