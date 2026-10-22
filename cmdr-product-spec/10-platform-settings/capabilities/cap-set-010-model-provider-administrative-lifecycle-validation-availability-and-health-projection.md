---
id: CAP-SET-010
title: Model Provider Administrative Lifecycle, Validation, Model Availability and Health Projection
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
delivery_status: defined
delivery_mode: planned
updated: 2026-08-13
requirement_ids: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-006, REQ-PROD-009, REQ-PROD-010, REQ-PROD-011, REQ-PROD-012, REQ-AI-001, REQ-AI-003, REQ-AI-004, REQ-AI-009, REQ-AI-010, REQ-AI-011]
open_decisions: [OPEN-008, OPEN-012, OPEN-013]
source-of-truth: canonical
---
# CAP-SET-010 — Model Provider Administrative Lifecycle, Validation, Model Availability and Health Projection

## 1. Définition
Capability Platform Settings qui définit l'administration provider-neutral du canonical `Model Provider`: lifecycle, configuration, validation locale déterministe, handoff de validation externe et projection de model availability/health réellement sourcées. Elle ne crée aucun runtime provider.

## 2. Problème utilisateur
Un Platform Administrator doit pouvoir configurer un Model Provider tenant-scoped, comprendre son état, vérifier localement sa configuration, demander une validation externe sans confondre demande et exécution, et voir des observations de disponibilité/health sans qu'elles soient inventées par Settings.

## 3. Objectifs
Préserver un owner unique; utiliser le lifecycle canonique; appliquer Tenant/permissions/provenance; distinguer validation locale d'une opération externe; exposer model availability et provider-local health comme projections sourcées; préserver un chemin manuel/déterministe; conserver les décisions de support/runtime ouvertes.

## 4. Non-objectifs
Ne pas définir de canonical `Model`, Model Catalog object, generic Provider, Provider Instance, Provider Endpoint, Routing Policy, Connection, Connector ou Credential. Ne pas définir de provider SDK, inference gateway, model-discovery engine, health-probe engine, protocole, API, schéma physique, provider support matrix, nouvel écran, nouvelle permission ou implémentation.

## 5. Propriétaire
`Platform Settings Product Lead` est l'unique capability owner. Platform Settings reste owner de `Model Provider`. Security conserve authorization/RBAC/ABAC/SoD/step-up. Studio conserve les objets et exécutions agentiques. Govern conserve Policy/Decision/Approval. Un exécuteur technique externe n'est pas attribué à Settings par cette capability.

## 6. Utilisateurs
Platform Administrator principal; Security Administrator secondaire; Studio et autres consommateurs autorisés peuvent lire une projection Model Provider sans obtenir l'ownership administratif ni le runtime.

## 7. Conditions d'entrée
Tenant résolu; acteur identifiable; permissions actuelles; Model Provider stable lorsqu'existant; version/état disponibles pour mutation; configuration minimale sourcée pour la validation demandée. Environment n'est requis que lorsqu'une relation canonique le fournit.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Model Provider cible ou intention administrative | `SET-MDL-001` / administrateur | référence/intention | oui | courante | aucune mutation |
| Tenant | contexte Settings | référence canonique | oui | courante | refus sans fallback |
| état/version/configuration provider | Model Provider | état + version + metadata | oui selon action | courante | relecture requise |
| Secret Reference liée | source configurée | référence seulement | seulement si sourcée | courante | dépendance non affirmée |
| contraintes de données/usage | owner canonique | contraintes | selon usage | courante | opération non éligible si requise |
| résultat de validation externe observé | runtime owner/source externe | outcome + provenance | seulement pour projection | observé | aucun succès déduit |
| model availability observée | source autorisée | liste/metadata + fraîcheur | seulement pour projection | sourcée | état unknown/partial |
| provider health observé | source autorisée | health + fraîcheur/provenance | seulement pour projection | sourcée | état unknown/partial |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Model Provider | Platform Settings | id, tenant-id, version, état, configuration administrative, timestamps/provenance | read/manage selon action |
| Tenant | Platform Settings | id/scope | read |
| Environment | Platform Settings | référence seulement si sourcée | read si applicable |
| Secret Reference | Platform Settings | identifiant/metadata non sensible seulement si relation sourcée | read selon dépendance |
| Permission context | Security | authorization/SoD/step-up outcome | consume |

Aucun canonical `Model` n'est lu parce qu'aucun objet de ce type n'est créé par le corpus actuel; les noms/capacités de modèles restent des metadata/projections sourcées.

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Model Provider | création/configuration/transition administrative sourcée | Platform Settings | Tenant obligatoire, états canoniques, version/provenance |
| Administrative audit event | émission/référence | Platform Settings / Administrative Audit | acteur/cible/action/outcome/corrélation sans valeur sensible |

Aucun Model, Model Availability, Provider Health, Validation Request, Credential, Connection ou Routing Policy canonique n'est créé.

## 11. Fonctionnalités
Créer/configurer un Model Provider; inspecter son état et metadata; valider localement la cohérence; préparer une demande de test/validation externe; projeter un résultat observé; projeter model availability et provider-local health avec fraîcheur; activer/désactiver selon lifecycle sourcé; refuser cross-tenant/stale/invalid; conserver audit et provenance.

## 12. Actions utilisateur
Class 0: inspection du Model Provider, projections et provenance. Class 1: validation locale strictement déterministe et sans effet externe. Class 2: création/mise à jour/activation/désactivation administrative réversible/versionnée lorsque sourcée, avec `OPEN-013` préservé. Une opération externe de provider n'est pas reclassifiée comme action Settings par cette capability.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| expliquer état/configuration | oui | oui | oui | oui | lecture canonique |
| validation locale | oui | oui | oui | explication seulement | règles déterministes |
| préparer demande externe | oui | oui pour préconditions | selon orchestration autorisée | suggestion seulement | action explicite |
| interpréter availability/health sourcées | oui | oui pour faits | oui | oui avec provenance/incertitude | projection brute + fraîcheur |
| mutation administrative | oui | oui | non autonome par IA | suggestion seulement | action explicite autorisée |

L'IA ne fabrique aucun modèle disponible, health, succès de test ou support provider et ne modifie pas le Model Provider de manière autonome.

## 14. États fonctionnels
États `Model Provider` exactement: `configured`, `validating`, `active`, `degraded`, `disabled`. `validating` ne prouve pas qu'une opération externe est en cours. `active` ne prouve pas qu'un provider/model est supporté ou disponible universellement. `degraded` doit rester sourcé.

## 15. États d'interface
Loading, Empty, Partial, Error, Offline, Permission denied et Stale. Partial/Stale expose la source manquante, la fraîcheur et la conséquence. Une disponibilité/health absente reste unknown/partial; elle n'est jamais synthétisée en succès.

## 16. Sorties
| Sortie | Objet/événement | Consommateur | Garantie |
|---|---|---|---|
| Model Provider state projection | référence Model Provider | Settings/consommateurs autorisés | état canonique + Tenant + version |
| local validation outcome | résultat administratif | administrateur/audit | succès/refus local + cause |
| external validation handoff | intention + contexte non sensible | runtime owner lorsqu'il existe | préconditions/provenance seulement |
| observed validation result | outcome sourcé | administrateur/audit | reflète uniquement le résultat reçu |
| model availability projection | metadata sourcée | Settings/Studio autorisés | source + fraîcheur, aucun Model object |
| provider health projection | observation sourcée | Settings/Health/Studio autorisés | source + fraîcheur, aucun probe revendiqué |

## 17. Transitions / handoffs
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| `SET-MDL-001` | inspect/configure/disable | Model Provider owner surface | provider/tenant/state/version | projection mise à jour |
| Settings | validation locale | admin UI/audit | configuration + règles | deterministic outcome |
| Settings | validation externe requise | runtime owner séparément sourcé | provider ref, tenant, préconditions, correlation id | observed outcome si produit |
| source autorisée | availability/health observée | Settings/Health/Studio | provider ref, observation, fraîcheur/provenance | projection |
| Settings | mutation/outcome | `SET-AUD-001` | acteur/cible/action/outcome/corrélation | trace administrative |

Aucun handoff ne transfère l'ownership du Model Provider ou du runtime.

## 18. Dépendances
`Model Provider`, Tenant, Environment seulement si sourcé, Secret Reference seulement si relation sourcée, `SET-MDL-001`, `SET-HLT-001`, `SET-AUD-001`, Security, data-policy/residency owners, Studio consommateurs. `OPEN-008`, `OPEN-012`, `OPEN-013` restent ouverts.

## 19. Source de vérité
`05-domain-model/objects/model-provider.md` pour objet/owner/états; `10-platform-settings/models-and-providers/README.md` et `provider-management.md` pour frontière administrative; Screen Register/`SET-MDL-001` pour surface; Permission Register pour IDs existants; Security, Product Boundaries, Requirements et OPEN pour trust/ownership/maturité.

## 20. Provenance et audit
Toute mutation, validation locale, demande externe, résultat observé, availability/health projetée et refus conserve acteur/source, Tenant, Model Provider ref, version/état, timestamps/fraîcheur et correlation id selon les mécanismes canoniques. Aucune donnée sensible n'est requise dans cette trace.

## 21. Permissions fonctionnelles
Consomme `perm.platform-settings.model-provider.read` et `perm.platform-settings.model-provider.manage`; `SET-MDL-001` conserve les aliases UI existants `perm.settings.model.read/manage`. Aucun nouvel ID. `.manage` reste une permission administrative et n'attribue pas un runtime provider.

## 22. Limites et erreurs
Model Provider ≠ Integration ≠ generic Provider. Model metadata ≠ canonical Model object. Cross-tenant, version stale, transition invalide, configuration insuffisante ou permission absente produit un refus explicite. Une source de disponibilité/health absente bloque la conclusion correspondante, pas l'administration provider-neutral.

## 23. Métriques
Transitions administratives acceptées/refusées, validations locales, demandes externes, résultats observés par catégorie, fraîcheur des projections availability/health, passages `degraded`, désactivations et refus cross-tenant. Aucun SLO, support provider ou runtime performance n'est inventé.

## 24. Classification de livraison
`defined / planned`. La capability est documentaire/provider-neutral. Elle ne prouve aucune implémentation, integration déployée, provider support, modèle supporté, runtime, protocole ou disponibilité produit. `OPEN-008`/`OPEN-012` restent ouverts.

## 25. Critères d'acceptation
**Given** un Model Provider tenant-scoped et une permission manage valide, **When** une configuration sourcée est modifiée, **Then** version, lifecycle canonique et provenance sont conservés sans créer de generic Provider/Model object.  
**Given** une validation strictement locale, **When** elle est lancée, **Then** elle retourne un outcome déterministe sans déclarer la validation externe réussie.  
**Given** une validation externe nécessaire et aucun exécuteur Settings sourcé, **When** l'administrateur choisit `Tester`, **Then** Settings produit seulement préconditions/handoff/corrélation et attend un outcome observé.  
**Given** une source autorisée de model availability, **When** sa projection est affichée, **Then** source et fraîcheur sont visibles et aucun canonical Model n'est créé.  
**Given** une observation provider health, **When** elle est projetée, **Then** elle reste attribuée à sa source et n'est pas présentée comme un probe Settings.  
**Given** une cible cross-tenant ou une version stale, **When** une mutation est demandée, **Then** elle est refusée sans fallback.  
**Given** l'IA indisponible, **When** le provider est administré, **Then** inspection, validation locale et mutations autorisées restent manuelles/déterministes.

## 26. Questions ouvertes
`OPEN-008` conserve la question de disponibilité/support; `OPEN-012` conserve provider/delivery scope; `OPEN-013` conserve l'autorité par défaut des mutations Class 2. Aucune n'est résolue par ce texte.

## 27. Consommateurs documentaires
Models & Providers, Administrative Audit, Health, Security, Studio, Shared, Quality et Roadmap; les autres produits peuvent consommer une projection autorisée. Les consommateurs ne deviennent ni owner du Model Provider ni owner de son administration.
