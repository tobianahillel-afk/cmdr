---
id: CAP-SET-006
title: Role Administrative Lifecycle, Constraints and Principal-Relation Boundary
product: platform-settings
module: users-and-roles
owner: Platform Settings Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-12
requirement_ids: [REQ-PROD-006, REQ-OBJ-001, REQ-SEC-001, REQ-SEC-002, REQ-SEC-006]
open_decisions: [OPEN-013]
source-of-truth: canonical
---
# CAP-SET-006 — Role Administrative Lifecycle, Constraints and Principal-Relation Boundary
## 1. Définition
Capability Settings qui définit l’administration du `Role` canonique, ses contraintes sourcées et la frontière de relation typée Principal/Role, sans devenir moteur d’autorisation ni d’assignation générique.
## 2. Problème utilisateur
Un administrateur doit créer, activer, déprécier et contraindre un Role tenant-scoped et comprendre sa relation à un Principal sans confondre Role, Permission, Group ou Decision Authority.
## 3. Objectifs
Définir lifecycle `draft/active/deprecated`, baselines, contraintes personnalisées, condition d’expiry, provenance et validation de relation tenant-scoped, avec sécurité et absence d’autorité implicite.
## 4. Non-objectifs
Ne pas définir Group, Assignment/RoleAssignment/PermissionAssignment, Permission, effective-access computation, inheritance, precedence, moteur RBAC/ABAC, Decision Authority, API, schéma physique, Permission ID ou Screen ID.
## 5. Propriétaire
Platform Settings Product Lead est capability owner; Platform Settings reste object owner de Role et Principal. Security reste owner du Permission Model, de l’autorisation et de l’évaluation effective.
## 6. Utilisateurs
Platform Administrator principal; Security Administrator et reviewers autorisés secondaires. Une relation Principal/Role n’est jamais une preuve autonome de permission ou d’autorité Govern.
## 7. Conditions d’entrée
Tenant résolu, Role ou intention de création connue, Principal de référence résoluble lorsqu’une relation est inspectée, permissions serveur valides, contraintes Security et versions courantes disponibles.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Role cible/draft | utilisateur/SET-IAM-001 | référence ou intention | oui | courant | aucune mutation |
| lifecycle/constraint intent | utilisateur | intention administrative | oui pour mutation | courant | lecture seule |
| Principal relation reference | objet canonique | relation typée sourcée | non | versionnée | aucune relation inventée |
| contexte Security | Security | tenant/RBAC/ABAC/SoD/step-up | oui | courant | refus explicite |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Role | Platform Settings | id, tenant-id, version, state, relations, contraintes sourcées | read/manage |
| Principal | Platform Settings | id, tenant-id, version, relation typée uniquement | read |
| Permission context | Security | décision d’autorisation et contraintes | consume |
## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Role | création/lifecycle/contrainte sourcée | Platform Settings | `draft/active/deprecated` seulement; version/provenance conservées |
| Principal/Role typed relation | mutation uniquement si opération déjà sourcée | Platform Settings | tenant-scoped; aucune sémantique de Permission inventée |
| Administrative audit event | émission/référence | Platform Settings / audit consumer | trace, pas nouvel objet canonique |
## 11. Fonctionnalités
Création de Role, inspection, lifecycle, baselines, custom constraints, condition/contrainte d’expiry, validation tenant scope, inspection de relation typée Principal/Role et refus de toute autorité implicite.
## 12. Actions utilisateur
Class 0: inspecter Role/relation. Class 1: valider transition, tenant scope ou condition d’expiry. Class 2: mutation Role ou relation uniquement si explicitement sourcée et autorisée. Class 3/4: aucune autorité nouvelle.
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| expliquer Role/contrainte | oui | oui | oui | oui | source + règles |
| valider lifecycle/expiry | oui | oui | oui | oui | moteur déterministe propriétaire |
| suggérer relation sûre | oui | oui | non autonome | oui | inspection + validation humaine |
| muter Role/relation | oui | oui | non autonome | non autonome | action explicite autorisée |
## 14. États fonctionnels
Role utilise exactement `draft`, `active`, `deprecated`. L’expiry est une condition/contrainte temporelle, jamais un état `expired`. Principal conserve ses propres états et lifecycle séparés.
## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied et Stale. Une donnée de relation ou contrainte partielle reste explicitement partielle et n’est pas complétée par inférence.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Role state projection | Role reference | Settings/consommateurs | état canonique + version |
| constraint/expiry assessment | assessment fonctionnel | administrateur | condition sourcée, pas autorisation effective |
| relation boundary outcome | validation/référence | Settings/Security consumer | relation typée tenant-scoped, pas Permission |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| SET-IAM-001 | create/inspect/lifecycle | Role owner surface | role/tenant/state/version | même contexte |
| Settings | constraint/relation validation | Security context | refs + scope uniquement | allow/deny/constraints |
| Settings | mutation administrative | SET-AUD-001 | acteur/cible/action/résultat | audit immutable |
## 18. Dépendances
`OBJ-ROLE`, `OBJ-PRINCIPAL`, Tenant, Users & Roles, Security Permission Model/catalog, ABAC, SoD, step-up, Decision Authority boundary et Administrative Audit.
## 19. Source de vérité
`05-domain-model/objects/role.md` pour états/relations; `role-management.md` pour baselines/custom constraints/expiry/no implicit authority; Principal object pour référence; Security pour permission/autorisation.
## 20. Provenance et audit
Chaque mutation Role ou relation réellement sourcée conserve actor, Tenant, cible, justification, version, timestamps et correlation id. L’historique de Role déprécié n’est pas effacé.
## 21. Permissions fonctionnelles
Consomme `perm.settings.identity.read/manage` côté écran et `perm.platform-settings.role.read/manage` côté objet, avec Principal read lorsque nécessaire. Dette d’alias conservée; aucun nouvel ID.
## 22. Limites et erreurs
Role ≠ Permission/Group/job title/Decision Authority. Relation ≠ Permission grant. Aucune assignment inheritance/precedence/effective-access. Cross-tenant et transitions invalides sont refusés sans fallback.
## 23. Métriques
Créations/transitions Role, refus de contraintes, conditions d’expiry rencontrées, validations/refus de relation et erreurs cross-tenant/provenance. Aucun SLO/KPI imposé.
## 24. Classification de livraison
`defined / planned`; aucune implémentation RBAC/ABAC, moteur d’assignation, support runtime ou disponibilité effective n’est revendiqué.
## 25. Critères d’acceptation
**Given** un Role `draft` tenant-scoped et une permission manage valide, **When** une activation sourcée est demandée, **Then** la transition vers `active` est validée, versionnée et auditée.  
**Given** une condition d’expiry atteinte, **When** le Role est évalué administrativement, **Then** l’expiry est exposée comme contrainte/condition et aucun état `expired` n’est créé.  
**Given** un Principal et un Role de Tenants incompatibles, **When** leur relation est inspectée ou mutée, **Then** l’opération est refusée sans créer de sémantique d’assignation.  
**Given** une relation Principal/Role visible, **When** l’administrateur l’inspecte, **Then** elle n’est jamais présentée comme permission effective ou Decision Authority.
## 26. Questions ouvertes
OPEN-013 reste ouvert pour les mutations Class 2. Les sémantiques génériques d’assignment, d’inheritance et d’effective access restent explicitement hors scope plutôt que résolues localement.
## 27. Consommateurs documentaires
Users & Roles, Administrative Audit, Security, Govern et tous produits consommateurs autorisés; chaque projection conserve l’ownership Security de l’autorisation et l’ownership Govern de Decision Authority.
