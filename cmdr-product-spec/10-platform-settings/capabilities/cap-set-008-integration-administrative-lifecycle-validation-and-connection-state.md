---
id: CAP-SET-008
title: Integration Administrative Lifecycle, Validation and Connection State
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
delivery_status: defined
delivery_mode: planned
updated: 2026-08-13
requirement_ids: [REQ-PROD-006, REQ-PROD-009, REQ-PROD-010, REQ-PROD-011, REQ-PROD-012, REQ-AI-001, REQ-AI-003, REQ-AI-004, REQ-AI-009, REQ-AI-010]
open_decisions: [OPEN-008, OPEN-013]
source-of-truth: canonical
---
# CAP-SET-008 — Integration Administrative Lifecycle, Validation and Connection State

## 1. Définition
Capability Platform Settings qui définit l’administration du lifecycle du `Integration` canonique, son état de configuration et les sémantiques administratives de validation/connexion sans créer de `Connection`, `Connector`, `Credential` ni moteur d’exécution externe.

## 2. Problème utilisateur
Un Platform Administrator doit pouvoir configurer, inspecter, valider et désactiver une Integration tenant-scoped, comprendre un état `degraded` ou `error`, demander un test de connexion et exploiter un résultat observé sans confondre administration Settings et exécution technique d’un provider ou connector.

## 3. Objectifs
Fournir un lifecycle administratif unique pour Integration; préserver les états canoniques; appliquer Tenant, permissions et provenance; distinguer validation déterministe de probe externe; conserver un chemin manuel/déterministe sans IA; permettre un handoff vers le véritable runtime owner lorsqu’un test externe doit être exécuté.

## 4. Non-objectifs
Ne pas définir de Connection/Connector/Credential, Model Provider, Data Source, Parser, protocole, API, client réseau, moteur connector, credential-retrieval runtime, provider support matrix, schéma physique, Screen ID, Permission ID, implémentation d’ingestion ni runtime externe. Ne pas rendre Environment obligatoire sans source canonique.

## 5. Propriétaire
`Platform Settings Product Lead` est l’unique capability owner. Platform Settings reste owner de `Integration`. Security reste owner du Permission Model, RBAC/ABAC, tenant isolation, SoD, step-up et autorisation. L’exécuteur technique d’un probe externe n’est pas attribué par les sources canoniques actuelles et n’est donc pas revendiqué ici.

## 6. Utilisateurs
Platform Administrator principal; Security Administrator et consommateurs autorisés secondaires pour lecture/provenance. Studio, Endpoint, Command, Investigate et Govern peuvent consommer une projection Integration sans obtenir son ownership administratif.

## 7. Conditions d’entrée
Tenant résolu avant toute cible; acteur administratif identifiable; permission serveur actuelle; Integration stable lorsqu’existante; configuration suffisamment présente pour l’opération demandée; sources Security disponibles. Environment est accepté seulement lorsqu’une relation canonique le fournit.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Integration cible ou intention administrative | `SET-SEC-001` / administrateur | référence ou intention | oui | courante | aucune mutation |
| Tenant | contexte Settings | référence canonique | oui | courante | refus sans fallback |
| version/état Integration | objet Integration | état + version | oui pour mutation | courante | relecture requise |
| endpoint et metadata `capabilities` de l’Integration | objet/module Integration | configuration administrative | selon opération | courante | validation/test non éligible |
| permission context | Security | décision RBAC/ABAC/SoD/step-up | oui pour mutation | courante | refus serveur |
| résultat de test réellement observé | runtime owner externe à cette capability | résultat/error + provenance | seulement pour projection de résultat | observé | aucun succès déduit |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Integration | Platform Settings | id, tenant-id, version, état, endpoint, metadata `capabilities`, timestamps/provenance | read/manage selon action |
| Tenant | Platform Settings | id et scope | read |
| Environment | Platform Settings | référence uniquement si source canonique | read si applicable |
| Permission context | Security | autorisation, SoD, step-up, tenant outcome | consume |

La metadata `capabilities` d’une Integration décrit la connexion externe selon son contrat; elle n’est jamais le canonical object `Capability`, un `CAP-*` ID ou une entrée du Capability Register.

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Integration | création/configuration ou transition administrative sourcée | Platform Settings | états canoniques seulement, Tenant obligatoire, version/provenance conservées |
| Administrative audit event | émission/référence | Platform Settings / Administrative Audit | résultat administratif et corrélation, jamais de secret |

Aucun objet Connection, Connector, Credential ou Test Request canonique n’est créé. Un test externe reste un handoff fonctionnel lorsque son exécuteur n’est pas canonique dans Settings.

## 11. Fonctionnalités
Créer/configurer une Integration, inspecter endpoint et metadata, exécuter une validation locale sans effet, préparer/demander un test de connexion, projeter un résultat/error réellement observé, désactiver de façon administrative, exposer `degraded`/`error`, refuser les opérations cross-tenant et conserver provenance/audit.

## 12. Actions utilisateur
Class 0: inspecter Integration, état, configuration non secrète et provenance. Class 1: uniquement validation déterministe sans effet externe. Class 2: création/mise à jour/désactivation administrative réversible/versionnée lorsque les règles existantes l’autorisent; `OPEN-013` reste ouvert pour la politique par défaut. Un vrai probe sortant n’est classé ni C1 ni C2 ici car son exécuteur, son effet et sa réversibilité ne sont pas attribués par les sources; cette exécution est exclue plutôt qu’inventée.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| expliquer configuration/état non secret | oui | oui | oui | oui | lecture des champs/états canoniques |
| expliquer une validation locale refusée | oui | oui | oui | oui | règles déterministes + cause |
| préparer un test administratif | oui | oui | oui selon orchestration autorisée | oui, suggestion seulement | action explicite + préconditions |
| interpréter un résultat réellement observé | oui | oui pour faits bruts | oui | oui avec provenance/incertitude | résultat/error brut + provenance |
| mutation Integration | oui | oui | non autonome par IA | suggestion seulement | action explicite autorisée |

L’IA ne fabrique jamais un succès de connexion, ne contourne pas permissions/SoD/step-up et ne sélectionne pas un provider ou runtime.

## 14. États fonctionnels
États `Integration` exactement: `draft`, `validating`, `active`, `degraded`, `disabled`, `error`. `validating` ne prouve pas qu’un probe externe est en cours; `active` ne prouve ni disponibilité universelle ni couverture complète; `degraded`/`error` sont des états administratifs sourcés, pas des diagnostics inventés.

## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied et Stale exposent cause, fraîcheur et conséquence. Offline ou stale interdit une mutation exigeant garantie serveur. Un test sans résultat observé reste pending/unknown dans la projection; il n’est jamais présenté comme réussi.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Integration state projection | référence Integration | Settings et consommateurs autorisés | état canonique + Tenant + version |
| validation outcome | résultat administratif | administrateur/audit | succès/refus local + cause, sans prétendre à un probe externe |
| connection-test administrative handoff | intention + contexte non secret | runtime owner lorsqu’il existe | préconditions/provenance seulement, aucun moteur inventé |
| observed test result projection | résultat/error sourcé | administrateur/audit | reflète uniquement un résultat réellement reçu |
| lifecycle outcome | résultat administratif | Settings/audit | transition/refus + corrélation |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| `SET-SEC-001` | inspect/configure/disable | Integration owner surface | integration/tenant/state/version | projection mise à jour |
| Settings | validation locale | Integration/admin UI | configuration non secrète + règles | validation outcome |
| Settings | test externe requis | runtime owner non attribué par cette capability | integration ref, tenant, préconditions, correlation id; aucun secret brut | observed result/error si produit |
| Settings | mutation/outcome | `SET-AUD-001` | acteur/cible/action/résultat/corrélation | audit immutable |

Aucun handoff ne transfère l’ownership de l’Integration.

## 18. Dépendances
Objet `Integration`, Tenant, Environment seulement si sourcé, `SET-SEC-001`, Administrative Audit, Security Permission Model, ABAC, tenant isolation, SoD, step-up, audit/immutability, provenance/integrity, privacy/minimization. Models & Providers et Sources & Parsers sont adjacents mais ne sont pas absorbés. Studio/Endpoint/Investigate/Govern restent propriétaires de leurs runtimes/objets.

## 19. Source de vérité
`05-domain-model/objects/integration.md` pour objet/états; `10-platform-settings/secrets-and-connections/connection-management.md` pour Endpoint/Capabilities/Test/Disable; `secrets-and-connections/README.md` et `SET-SEC-001` pour module/actions; Security pour permissions et trust; Product Boundaries/Object Register pour ownership.

## 20. Provenance et audit
Toute création, mutation, validation, demande de test, résultat observé, désactivation ou refus conserve acteur, Tenant, Integration ref, version/état, justification pertinente, timestamps et correlation id selon les mécanismes canoniques. Aucune valeur secrète, token, mot de passe ou clé n’est enregistrée.

## 21. Permissions fonctionnelles
Consomme `perm.platform-settings.integration.read` et `perm.platform-settings.integration.manage` côté objet. `SET-SEC-001` conserve ses permissions de surface existantes; aucun nouvel ID et aucun bulk rename. `.manage` ne doit pas être interprété comme droit d’exécuter un probe externe si cette sémantique n’est pas canonique.

## 22. Limites et erreurs
Integration ≠ Connection object ≠ Connector object ≠ Model Provider ≠ Data Source. Metadata `capabilities` ≠ CMDR Capability/CAP-* ID. Une cible cross-tenant, une version stale, une transition invalide ou une permission insuffisante est refusée sans fallback. L’absence d’exécuteur de probe canonique bloque l’exécution technique, pas l’administration et le handoff sourcés.

## 23. Métriques
Transitions administratives acceptées/refusées, validations locales, demandes de test, résultats observés par catégorie success/error/unknown lorsque réellement fournis, passages `degraded`/`error`, désactivations, refus cross-tenant et erreurs de provenance. Aucune cible SLO/KPI, disponibilité provider ou implémentation n’est déduite.

## 24. Classification de livraison
`defined / planned`. Cette capability est documentaire et provider-neutral. Elle ne prouve aucun connector/runtime, API, protocole, provider support, probe engine, ingestion, déploiement ou disponibilité produit. `OPEN-008` reste ouvert sur disponibilité/support.

## 25. Critères d’acceptation
**Given** une Integration tenant-scoped et une permission manage valide, **When** une configuration sourcée est modifiée, **Then** la version, le lifecycle canonique et la provenance sont préservés sans créer de Connection/Connector object.  
**Given** une validation strictement locale et sans effet, **When** elle est exécutée, **Then** elle est classée Class 1 et retourne une cause déterministe sans déclarer la connexion externe réussie.  
**Given** l’action `Tester` et aucun exécuteur technique canonique attribué à Settings, **When** un test externe est requis, **Then** Settings limite son contrat à la demande/préconditions/handoff et ne revendique ni moteur réseau ni classe d’effet inventée.  
**Given** un résultat de test externe réellement reçu, **When** il est projeté, **Then** le résultat/error et sa provenance sont exposés sans fabriquer de succès ni secret.  
**Given** une cible d’un autre Tenant ou une permission absente, **When** une mutation est demandée, **Then** elle est refusée sans fallback ni fuite de données.  
**Given** l’IA indisponible, **When** l’Integration est administrée, **Then** lecture, validation et mutation autorisée restent accessibles manuellement/déterministiquement.

## 26. Questions ouvertes
`OPEN-008` reste ouvert pour disponibilité/support des plateformes/sources. `OPEN-013` reste ouvert pour la politique d’autorité par défaut des mutations Class 2. L’owner et la classification d’effet d’un éventuel probe externe ne sont pas inventés; ils devront être apportés par une source canonique avant de faire de Settings l’exécuteur.

## 27. Consommateurs documentaires
Secrets & Connections, Models & Providers, Sources & Parsers, Administrative Audit, Security, Studio, Endpoint, Command, Investigate, Govern, Shared, Quality et Roadmap. Les consommateurs utilisent une référence/projection Integration selon permissions et Tenant sans transfert d’ownership.
