---
id: CAP-SET-001
title: Tenant Administrative Lifecycle and Isolation Boundary
product: platform-settings
module: tenants-and-environments
owner: Platform Settings Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-12
requirement_ids: [REQ-PROD-006, REQ-SEC-001]
open_decisions: [OPEN-013]
source-of-truth: canonical
---
# CAP-SET-001 — Tenant Administrative Lifecycle and Isolation Boundary
## 1. Définition
Capability Settings qui définit l'administration du lifecycle Tenant et sa frontière d'isolation. Tenant n'est ni Customer, ni Environment, ni état de santé.
## 2. Problème utilisateur
Un administrateur doit provisionner, suspendre, conduire l'offboarding et archiver logiquement un Tenant tout en conservant isolation, provenance et limites d'autorité.
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
| Tenant ciblé | utilisateur/route | référence canonique | oui | courant | action impossible |
| transition demandée | utilisateur | lifecycle intent | oui | requête courante | aucune mutation |
| autorisation | Security | permission + tenant scope | oui | courant | refus explicite |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Tenant | Platform Settings | id, tenant-id, version, timestamps, state | read/manage |
| Permission context | Security | tenant scope, ABAC/RBAC result | consume |
## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Tenant | lifecycle transition | Platform Settings | transitions invalides refusées; version/provenance conservées |
| Administrative audit event | emit/reference | Platform Settings / audit consumer | événement, pas nouvel objet canonique |
## 11. Fonctionnalités
Provisioning, lecture d'état, transition vers active/suspended/offboarding/archived selon l'objet canonique, contexte de residency par référence, visibilité consumer et refus cross-tenant.
## 12. Actions utilisateur
Class 0: inspecter. Class 1: valider une transition sans effet. Class 2: appliquer une transition Settings réversible/versionnée lorsque la source l'autorise. Class 3/4: aucune autorité créée par cette capability.
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| expliquer état/transition | oui | oui | oui | oui | états + règles déterministes |
| proposer prochaine étape | oui | oui | oui | oui | actions sourcées et permissions |
## 14. États fonctionnels
`provisioning`, `active`, `suspended`, `offboarding`, `archived`; `suspended` ≠ deleted et `offboarding` ≠ effacement irréversible.
## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied et Stale exposent la cause, la fraîcheur et les conséquences sans inventer de données ni élargir le scope.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Tenant state projection | Tenant reference | Settings/consumers | état canonique + version |
| administrative change outcome | événement/résultat fonctionnel | Audit/administrateur | succès/refus + corrélation |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| SET-TEN-001 | lifecycle action | Tenant owner surface | tenant/ref/state/justification | même contexte |
| Settings | governed exception if separately required | Govern | Action Request reference only | Settings |
## 18. Dépendances
Security Permission Model, tenant-isolation, data-residency, Administrative Audit, Shared notification/job mechanisms seulement comme mécanismes consommés.
## 19. Source de vérité
`05-domain-model/objects/tenant.md` pour l'objet; `tenant-management.md` pour l'administration; Security pour isolation/permission.
## 20. Provenance et audit
Chaque mutation conserve actor, tenant, cible, action, résultat, justification, version, timestamps et correlation id; aucun secret brut.
## 21. Permissions fonctionnelles
`perm.platform-settings.tenant.read` et `perm.platform-settings.tenant.manage`; les aliases `perm.settings.tenant.*` restent dette de migration. Aucun nouvel ID de permission n'est permis dans ce lot.
## 22. Limites et erreurs
Pas de suppression physique, pas de Customer lifecycle, pas de fallback tenant, pas de cardinalité Environment inventée. Une nouvelle permission ou un nouvel Screen ID requis bloque le lot.
## 23. Métriques
Transitions acceptées/refusées, échecs cross-tenant, freshness/provenance gaps; aucune cible KPI/SLO.
## 24. Classification de livraison
`defined / planned`; aucune disponibilité runtime, implémentation, intégration ou support plateforme n'est revendiqué.
## 25. Critères d’acceptation
**Given** Tenant `active` et permission manage valide, **When** une suspension sourcée est demandée, **Then** la transition est validée, auditée et versionnée.  
**Given** une cible d'un autre Tenant, **When** une mutation est demandée, **Then** elle est refusée sans fallback.  
**Given** IA indisponible, **When** le lifecycle est administré, **Then** toutes les validations et mutations autorisées restent réalisables sans IA.
## 26. Questions ouvertes
OPEN-013 reste ouvert pour la politique par défaut des mutations Class 2; aucune nouvelle OPEN n'est créée.
## 27. Consommateurs documentaires
Tenants & Environments, Administrative Audit, Security, Command, Investigate, Govern, Studio, Endpoint, Shared, Quality et Roadmap.
