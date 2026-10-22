---
id: CAP-SET-003
title: Tenant and Environment Administrative Change Validation and Provenance
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
# CAP-SET-003 — Tenant and Environment Administrative Change Validation and Provenance
## 1. Définition
Capability Settings qui valide les changements administratifs Tenant/Environment et produit leur provenance sans créer de moteur générique de configuration.
## 2. Problème utilisateur
Un administrateur doit savoir si un changement est permis, ce qu'il affecte et comment son résultat est auditable avant et après mutation.
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
| cible | Settings | Tenant/Environment ref | oui | courant | validation impossible |
| intent | utilisateur | changement demandé | oui | requête courante | no-op |
| état/version courants | owner object | preconditions | oui | frais | stale/Partial |
| autorisation | Security | permission + ABAC/step-up/SoD | oui | courant | refus |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Tenant | Platform Settings | state/version/scope | read |
| Environment | Platform Settings | state/version/tenant-id | read |
| Permission context | Security | permission, ABAC, step-up/SoD result | consume |
## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Tenant/Environment | mutation owner-defined | Platform Settings | seulement après validation courante |
| Administrative audit event | emit/reference | Platform Settings / Audit | actor/target/action/result/justification/correlation |
## 11. Fonctionnalités
Validation de préconditions, transition state/scope, preview déterministe quand possible, refus cross-tenant, revalidation de permission, distinction no-effect/failure/success et handoff audit.
## 12. Actions utilisateur
Class 0: inspecter provenance. Class 1: valider/preview sans effet. Class 2: mutation Settings sourcée après revalidation. Class 3/4 non introduites.
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| expliquer refus | oui | oui | oui | oui | code/raison déterministe |
| résumer changement | oui | oui | oui | oui | diff fonctionnel sourcé |
## 14. États fonctionnels
`valid`, `invalid`, `permission-denied`, `scope-conflict`, `stale-precondition`, `applied`, `failed`; ces états sont fonctionnels et ne deviennent pas un nouvel objet canonique.
## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied et Stale exposent la cause, la fraîcheur et les conséquences sans inventer de données ni élargir le scope.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| validation outcome | résultat fonctionnel | administrateur | reason/preconditions/current version |
| administrative change outcome | événement | Administrative Audit | actor/tenant/object/action/result/correlation |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| SET-TEN-001 | validate change | validation fonctionnelle | target/current version/intent | Settings |
| validation | apply allowed change | owner lifecycle | validated intent + fresh authorization | Settings |
| Settings | emit audit | Administrative Audit | event reference | Settings |
## 18. Dépendances
Tenant/Environment objects, Security Permission Model/SoD/step-up, Administrative Audit. Shared jobs/notifications ne sont que mécanismes éventuels.
## 19. Source de vérité
Les objets Tenant/Environment restent canoniques; la capability ne crée ni Configuration object ni Effective Configuration engine.
## 20. Provenance et audit
Audit minimal: actor, tenant, environment si applicable, objet, action, precondition/version, résultat, justification et correlation id. Refus et cross-scope attempts sont aussi audités.
## 21. Permissions fonctionnelles
Permissions existantes read/manage uniquement. Security reste owner du modèle, du step-up et de la SoD. Nouvelle permission requise => BLOCKED et run séparé.
## 22. Limites et erreurs
Pas de default/inheritance/override/effective-config générique; pas de rollback Govern automatique; pas de secret brut. Stale precondition impose revalidation.
## 23. Métriques
Validation pass/fail, stale-precondition, permission-denied, cross-scope denial, audit handoff completeness.
## 24. Classification de livraison
`defined / planned`; aucune disponibilité runtime, implémentation, intégration ou support plateforme n'est revendiqué.
## 25. Critères d’acceptation
**Given** version de cible périmée, **When** une mutation est soumise, **Then** elle est refusée ou revalidée avant effet.  
**Given** permission retirée entre preview et apply, **When** apply est demandé, **Then** l'autorisation est réévaluée et la mutation est bloquée.  
**Given** changement réussi, **When** l'audit administratif est consulté, **Then** actor, scope, target, action, result, justification, version et correlation sont reconstruisibles.
## 26. Questions ouvertes
OPEN-013 reste ouvert; aucun choix de politique d'approbation par défaut n'est fait.
## 27. Consommateurs documentaires
Tenants & Environments, Administrative Audit, Security, Govern conditionnel, Shared mécanismes, Quality et Roadmap.
