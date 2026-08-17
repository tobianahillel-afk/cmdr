---
id: CAP-SET-004
title: Cross-Product Tenant and Environment Context Preservation
product: platform-settings
module: tenants-and-environments
owner: Platform Settings Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-12
requirement_ids: [REQ-PROD-008, REQ-UX-006, REQ-UX-007, REQ-UX-009, REQ-SEC-001]
open_decisions: []
source-of-truth: canonical
---
# CAP-SET-004 — Cross-Product Tenant and Environment Context Preservation
## 1. Définition
Capability Settings qui définit la sémantique des références Tenant/Environment projetées entre produits. **Platform Settings Product Lead est l'unique capability owner**; Experience Architecture reste owner du mécanisme de propagation et de conservation du contexte.
## 2. Problème utilisateur
Un utilisateur qui traverse Command, Investigate, Govern, Studio, Endpoint ou Settings doit conserver un contexte Tenant/Environment sûr sans gain de permission ni transfert d'ownership.
## 3. Objectifs
Définir une responsabilité administrative autonome, permission-aware, tenant-scoped et auditable, sans redéfinir les objets, permissions ou mécanismes partagés.
## 4. Non-objectifs
Ne pas définir d'API, protocole, schéma physique, moteur générique de configuration, nouvelle permission, nouvel Screen ID, implémentation, ni autorité Govern parallèle.
## 5. Propriétaire
**Platform Settings Product Lead est l'unique capability owner.** Experience Architecture est uniquement dependency/mechanism owner pour `context-preservation.md`; Design System possède l'affichage Context Bar/Inspector; Security possède l'autorisation.
## 6. Utilisateurs
Platform Administrator principal; utilisateurs Command, Investigate, Govern, Studio et Endpoint autorisés comme consommateurs secondaires.
## 7. Conditions d’entrée
Principal authentifié, Tenant source résolu, Environment compatible lorsque présent, return-origin non sensible et permission destination réévaluée.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| tenant ref | Platform Settings | contexte canonique | oui | courant | destination bloquée |
| environment ref | Platform Settings | contexte compatible | conditionnel | courant | choix explicite requis |
| return-origin | Experience Architecture | route/état non sensible | oui | session | retour minimal |
| autorisation destination | Security | permission actuelle | oui | courant | permission denied |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Tenant context | Platform Settings | tenant reference | read |
| Environment context | Platform Settings | environment reference | read |
| return-origin | Experience Architecture | route + état non sensible | consume |
| Permission context | Security | destination authorization | consume |
## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Context projection | compose/clear | Platform Settings semantics / Experience mechanism | aucune mutation de Tenant/Environment |
| return-origin | preserve | Experience Architecture | non sensible uniquement |
## 11. Fonctionnalités
Projection de Tenant, projection d'Environment compatible, inherited context visible, destination permission re-evaluation, clear on tenant change, no fallback, deep-link sûr et return-origin.
## 12. Actions utilisateur
Class 0: inspecter/naviguer. Class 1: valider la compatibilité du contexte destination. Aucune Class 2/3/4 mutation de Tenant/Environment n'est créée par la propagation.
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| expliquer contexte perdu | oui | oui | oui | oui | règles de propagation déterministes |
| suggérer destination compatible | oui | oui | oui | oui | navigation/permissions existantes |
## 14. États fonctionnels
`context-valid`, `context-incompatible`, `permission-denied`, `expired`, `source-missing`; les objets source gardent leurs propres lifecycles.
## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied et Stale exposent la cause, la fraîcheur et les conséquences sans inventer de données ni élargir le scope.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| context envelope | références | destination CMDR | tenant/env refs autorisées seulement |
| cleared context notice | événement UX | utilisateur | incompatibilités explicitement signalées |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Settings | open consumer | Command/Investigate/Govern/Studio/Endpoint | tenant/env/ref/return-origin | Settings |
| consumer | return | source product | refs encore autorisées | origin |
| tenant change | incompatibility | local context reset | nouveau tenant seulement | destination |
## 18. Dépendances
Experience Architecture `context-preservation.md` pour propagation, Design System Context Bar/Inspector pour présentation, Security pour permission, owners consommateurs pour données.
## 19. Source de vérité
Platform Settings reste owner sémantique Tenant/Environment. Experience Architecture possède le mécanisme de conservation; Design System affiche; la destination réévalue l'accès. Aucune projection ne transfère l'ownership.
## 20. Provenance et audit
Tracer source, destination, tenant/env refs, permission outcome, context clear/expiry et correlation id sans payload brut ni secret.
## 21. Permissions fonctionnelles
Permissions de lecture source existantes seulement; aucune nouvelle permission. La navigation ne crée jamais manage/export/execute/approve. **Nouveau Permission ID ou Screen ID requis => BLOCKED et run séparé.**
## 22. Limites et erreurs
Pas de cross-tenant mutation, pas de fallback silencieux, pas de payload brut/secret, pas de duplication de Context Bar, pas de nouvelle navigation globale.
## 23. Métriques
Context restore success/failure, incompatible clears, permission-denied transitions, stale/expired context; aucune cible KPI/SLO.
## 24. Classification de livraison
`defined / planned`; aucune disponibilité runtime, implémentation, intégration ou support plateforme n'est revendiqué.
## 25. Critères d’acceptation
**Given** Tenant/Environment autorisés, **When** l'utilisateur ouvre un produit consommateur, **Then** les références sont conservées et la destination réévalue la permission.  
**Given** changement de Tenant, **When** le contexte précédent devient incompatible, **Then** il est effacé selon le mécanisme Experience Architecture et aucun fallback n'est choisi.  
**Given** permission retirée, **When** un deep link est restauré, **Then** le contenu protégé est masqué et seule une trace non sensible est conservée.
## 26. Questions ouvertes
Aucune nouvelle OPEN. La mécanique de propagation reste Experience Architecture; aucun co-ownership de capability n'est créé.
## 27. Consommateurs documentaires
Command, Investigate, Govern, Studio, Endpoint, Platform Settings, Shared, Experience Architecture, Design System, Security, Quality et Roadmap.
