---
id: CAP-SET-005
title: Principal Administrative Lifecycle and Identity State
product: platform-settings
module: users-and-roles
owner: Platform Settings Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-12
requirement_ids: [REQ-PROD-006, REQ-OBJ-001, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-013]
source-of-truth: canonical
---
# CAP-SET-005 — Principal Administrative Lifecycle and Identity State
## 1. Définition
Capability Settings qui définit l’administration du lifecycle du `Principal` canonique et de son état d’identité, pour identités humaines ou de service, sans créer d’objet User ou ServiceIdentity.
## 2. Problème utilisateur
Un administrateur doit initialiser, consulter, suspendre ou révoquer un Principal dans son Tenant tout en distinguant identité, authentification et autorisation et en conservant une provenance exploitable.
## 3. Objectifs
Fournir une responsabilité administrative autonome, tenant-scoped, permission-aware et auditable autour des états `pending`, `active`, `suspended`, `revoked`, avec chemins manuels/déterministes complets.
## 4. Non-objectifs
Ne pas définir Group, Group Membership, Role, AccessAssignment, EffectiveAccess, protocole IdP/SSO, authentification, API, schéma physique, Permission ID, Screen ID, suppression physique ni autorité Govern.
## 5. Propriétaire
Platform Settings Product Lead est l’unique capability owner. Platform Settings reste object owner de Principal; Security reste owner des permissions, RBAC/ABAC, SoD, step-up et autorisation.
## 6. Utilisateurs
Platform Administrator principal; Security Administrator et consommateurs CMDR autorisés secondaires. Un Service Principal reste une identité de service, jamais un humain implicite.
## 7. Conditions d’entrée
Principal acteur authentifié, Tenant résolu avant la cible, cible Principal stable lorsqu’existante, permissions serveur actuelles, règles Security disponibles et source canonique résoluble.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Principal cible ou intention d’initialisation | utilisateur/SET-IAM-001 | référence ou intention | oui | courant | aucune mutation |
| état/intention de transition | utilisateur | lifecycle intent | oui pour mutation | courant | lecture seule |
| contexte Tenant | Settings/Security | tenant scope | oui | courant | refus sans fallback |
| mapping MFA/SSO | source d’identité configurée | référence administrative | non | sourcée | ne pas inventer |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Principal | Platform Settings | id, tenant-id, version, timestamps, state, relations sourcées | read/manage |
| Tenant | Platform Settings | id et scope | read |
| Permission context | Security | RBAC/ABAC/tenant/SoD/step-up outcome | consume |
## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Principal | initialisation ou transition sourcée | Platform Settings | états canoniques seulement; mutation auditée/versionnée |
| Administrative audit event | émission/référence | Platform Settings / audit consumer | événement de provenance, pas nouvel objet canonique |
## 11. Fonctionnalités
Lecture d’état, distinction human/service, initialisation correspondant à l’action UI `Inviter` sans état `invited`, validation de transition, suspension, révocation, référence MFA/SSO et refus tenant/permission explicite.
## 12. Actions utilisateur
Class 0: inspecter Principal/état/provenance. Class 1: valider sans effet une transition ou contrainte. Class 2: initialiser ou appliquer une mutation sourcée et autorisée. Class 3/4: aucune nouvelle autorité Settings.
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| expliquer état ou refus | oui | oui | oui | oui | règles et état sourcés |
| suggérer prochaine action sûre | oui | oui | oui | oui | transitions canoniques + permissions |
| mutation Principal | oui | oui | non autonome | non autonome | action explicite autorisée |
## 14. États fonctionnels
`pending`, `active`, `suspended`, `revoked` exactement. `suspended` ≠ `revoked`; `revoked` ≠ suppression physique; `Inviter` ne crée jamais un état canonique `invited`.
## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied et Stale exposent cause, fraîcheur et conséquence. Offline interdit toute mutation dont la garantie serveur n’est pas disponible.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Principal state projection | Principal reference | Settings/consommateurs autorisés | état canonique + version |
| lifecycle outcome | résultat administratif | administrateur/audit | succès/refus + cause + corrélation |
| MFA/SSO mapping reference | référence | Settings | contexte seulement, pas preuve d’autorisation |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| SET-IAM-001 | inspect/transition | Principal owner surface | principal/tenant/state/version | même contexte |
| Settings | identity mutation outcome | SET-AUD-001 | acteur/cible/action/résultat/corrélation | audit immutable |
| Settings | autorité externe réellement requise | Govern/Security owner | référence seulement | Settings |
## 18. Dépendances
`OBJ-PRINCIPAL`, Tenant, Users & Roles, Security Permission Model, ABAC, tenant isolation, SoD, step-up, Administrative Audit et mécanismes Shared consommés sans transfert d’ownership.
## 19. Source de vérité
`05-domain-model/objects/principal.md` pour objet/états; `users-and-roles/user-management.md` pour lifecycle/MFA-SSO/suspension/audit; `SET-IAM-001` pour actions UI; Security pour autorisation.
## 20. Provenance et audit
Toute création, transition, relation ou suppression logique sourcée conserve acteur, Tenant, justification, version, timestamps et correlation id. Aucun secret ni donnée non nécessaire n’est enregistré.
## 21. Permissions fonctionnelles
Consomme `perm.settings.identity.read/manage` côté surface et `perm.platform-settings.principal.read/manage` côté objet. Leur coexistence est dette d’alias/normalisation; aucun nouvel ID ni bulk rename.
## 22. Limites et erreurs
Principal ≠ Person/Customer/Tenant; Human Principal ≠ Service Principal; mapping d’authentification ≠ autorisation. Toute cible cross-tenant, transition invalide ou permission insuffisante est refusée côté serveur sans fallback.
## 23. Métriques
Transitions acceptées/refusées, suspensions/révocations, refus cross-tenant, erreurs de provenance et fraîcheur. Aucune cible SLO/KPI ni implication d’implémentation.
## 24. Classification de livraison
`defined / planned`; aucune disponibilité runtime, intégration IdP, protocole, support plateforme ou implémentation n’est revendiqué.
## 25. Critères d’acceptation
**Given** un Principal `active` tenant-scoped et une permission manage valide, **When** une suspension sourcée est demandée, **Then** la transition vers `suspended` est validée, versionnée et auditée.  
**Given** l’action UI `Inviter`, **When** un Principal est initialisé selon la source, **Then** aucun état canonique `invited` n’est créé et seuls les états de `OBJ-PRINCIPAL` sont utilisés.  
**Given** une cible d’un autre Tenant ou une autorisation absente, **When** une mutation est demandée, **Then** elle est refusée sans fallback ni fuite de données.  
**Given** l’IA indisponible, **When** le lifecycle est administré, **Then** lecture, validation et mutation autorisée restent possibles par chemin déterministe/manual.
## 26. Questions ouvertes
OPEN-013 reste ouvert pour la politique par défaut des mutations Class 2. Aucune décision sur protocole d’authentification, Group ou assignment n’est créée par cette capability.
## 27. Consommateurs documentaires
Users & Roles, Administrative Audit, Security, Command, Investigate, Govern, Studio, Endpoint, Shared, Quality et Roadmap; toute projection conserve ownership et autorisation du consommateur.
