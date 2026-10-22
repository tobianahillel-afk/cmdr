---
id: CAP-STD-015
title: Provider, Runtime, Integration and Secret Reference Boundaries
product: cmdr-studio
module: studio-foundations
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-PROD-017, REQ-SEC-001]
open_decisions: []
source-of-truth: canonical
---
# CAP-STD-015 — Provider, Runtime, Integration and Secret Reference Boundaries
## 1. Définition
Studio consomme Provider, Runtime, Integration, Secret Reference, health et tenant/environment projections; Platform Settings conserve configuration, credentials, secrets et administrative lifecycle.
## 2. Problème utilisateur
Une dépendance disponible ne doit être confondue ni avec permission, compatibility ni execution eligibility.
## 3. Objectifs
Référencer dependencies sans copier les valeurs sensibles ni déplacer l'ownership Settings.
## 4. Non-objectifs
Aucun provider imposé, credential store, secret reveal, connection setup ou runtime implementation.
## 5. Propriétaire
Studio possède seulement les bindings fonctionnels; Platform Settings possède provider/integration/secret/admin runtime configuration.
## 6. Utilisateurs
Automation Designer, Studio Operator et Platform Administrator.
## 7. Conditions d’entrée
Tenant/env, dependency refs et authorization context connus.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Provider/Runtime refs | Platform Settings | reference metadata | conditionnel | health visible | unavailable |
| Integration ref | Platform Settings | configured ref | conditionnel | courant | blocked |
| Secret Reference | Platform Settings | opaque ref | conditionnel | authorization courant | masked/blocked |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Model Provider/Integration | Platform Settings | identity/health | read |
| Secret Reference | Platform Settings | opaque metadata | read metadata |
| Tenant/Environment | Platform Settings | scope/config | read |
## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Studio dependency binding | bind reference | CMDR Studio | jamais credential/secret brut |
## 11. Fonctionnalités
Dependency discovery, reference binding, health/freshness projection, compatibility inputs et secret-safe display.
## 12. Actions utilisateur
Class 0 inspect refs/health; Class 1 validate availability; Class 2 bind an authorized reference. No secret reveal/use is implied.
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| availability check | oui | oui | oui | oui | health/ref check |
| suggest dependency | oui | oui | oui | oui | manual selection |
## 14. États fonctionnels
available, degraded, unavailable, stale, unauthorized, incompatible; these do not modify Settings lifecycle.
## 15. États d’interface
Restricted refs expose only authorized metadata; raw secret is never rendered/logged.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| dependency refs | reference set | Tool/Skill | Settings ownership retained |
| health projection | status | Studio | availability != authorization |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Studio | select ref | Settings projection | tenant/env/ref | Studio |
| Settings health | refresh | Studio | ref/health/freshness | no rebinding |
## 18. Dépendances
Platform Settings, Security tenant isolation/secrets, Shared health display where applicable.
## 19. Source de vérité
Settings is source for configured provider/integration/Secret Reference/admin runtime state; Studio stores references only.
## 20. Provenance et audit
Record ref ID, actor, binding purpose, tenant/env and health snapshot; never raw secret.
## 21. Permissions fonctionnelles
Reference read, restricted metadata read, binding prepare and use eligibility; raw secret operations excluded.
## 22. Limites et erreurs
Missing/stale ref, revoked integration, unavailable runtime, cross-tenant mismatch and denied metadata remain explicit.
## 23. Métriques
Unavailable dependencies, stale health, invalid bindings, secret-safe violations target zero conceptually without KPI threshold.
## 24. Classification de livraison
`defined / planned`; no provider/runtime selected or integrated.
## 25. Critères d’acceptation
**Given** provider healthy, **When** Tool eligibility is checked, **Then** health alone does not grant authorization.  
**Given** Secret Reference, **When** Studio displays it, **Then** raw value stays hidden and Settings owner visible.  
**Given** runtime unavailable, **When** dependency validated, **Then** Tool/Skill remains defined but execution unavailable.
## 26. Questions ouvertes
No OPEN-008 resolution: Endpoint support remains future. No provider/runtime choice is made.
## 27. Consommateurs documentaires
Tools, Skills, future STD lots, Platform Settings, Security, Quality and Roadmap.
