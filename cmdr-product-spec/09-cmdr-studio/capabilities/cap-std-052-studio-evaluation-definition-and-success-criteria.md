---
id: CAP-STD-052
title: Studio Evaluation Definition and Success Criteria
product: cmdr-studio
module: assurance
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-10
requirement_ids: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-016, REQ-PROD-019, REQ-OBJ-009, REQ-AI-002, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-013]
source-of-truth: canonical
---
# CAP-STD-052 — Studio Evaluation Definition and Success Criteria

## 1. Définition
Define a version-pinned Studio Evaluation definition with purpose, target asset/version, scope, criteria, expected characteristics, required observations, failure and inconclusive conditions, tenant/environment and source restrictions without turning criteria into universal policy.

## 2. Problème utilisateur
Sans cette capability, les opérateurs Studio peuvent confondre état documentaire, qualité observée, autorité, publication, déploiement ou résultat technique, perdre la version exacte évaluée et prendre une projection issue d’un autre owner pour une vérité Studio. Le comportement doit rester explicable même avec dépendance indisponible, donnée restreinte ou IA absente.

## 3. Objectifs
- fournir une sémantique Studio versionnée et tenant-scoped pour **Evaluation Definition** ;
- pinner les versions, sources, critères/dépendances et limites nécessaires ;
- préserver les owners Govern, Settings, Shared et Endpoint ;
- distinguer état de lifecycle, résultat observé, permission et autorité ;
- rendre warnings, partial/inconclusive/failure et supersession visibles ;
- conserver une voie manuelle ou déterministe sans IA.

## 4. Non-objectifs
Ne pas définir d’API/protocole, moteur d’évaluation/simulation/déploiement, package registry, scheduler, schéma JSON final, stockage physique, RBAC/ABAC final, provider/runtime obligatoire, Screen ID, detailed screen rewrite, commande Endpoint ou autorité Govern parallèle. Un PASS documentaire ne prouve aucune implémentation.

## 5. Propriétaire
CMDR Studio / assurance owns la sémantique fonctionnelle de **Evaluation Definition**. Les objets source restent chez leur owner canonique; Govern conserve Approval/Decision/Response Run/Result et rollback de réponse, Settings l’administration provider/integration/secret/environment, Shared les moteurs génériques et Endpoint les primitives endpoint.

## 6. Utilisateurs
Principal : Studio Assurance / Release Operator selon l’étape. Secondaires : Studio Asset Owner, Workflow/Agent/Skill Author, Studio Reviewer, Platform Administrator comme source Settings, Govern Reviewer lorsque l’autorité est requise, Auditor et consommateurs produit autorisés.

## 7. Conditions d’entrée
Asset et version exacts, tenant/environment explicites, références source résolubles, permission de lecture puis de mutation locale quand nécessaire, dépendances et restrictions connues, et absence de tentative de transformer une projection technique en Approval/Decision/Result. Les valeurs de secret restent des Secret References.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| target Studio asset/version | Studio or canonical source owner | functional reference/context | oui | exact version/current assessment | blocked/incomplete |
| evaluation purpose and scope | Studio or canonical source owner | functional reference/context | oui | exact version/current assessment | blocked/incomplete |
| criterion set | Studio or canonical source owner | functional reference/context | oui | exact version/current assessment | blocked/incomplete |
| expected characteristics | Studio or canonical source owner | functional reference/context | selon scope | exact version/current assessment | explicit missing/unknown |
| required observations | Studio or canonical source owner | functional reference/context | selon scope | exact version/current assessment | explicit missing/unknown |
| failure/inconclusive conditions | Studio or canonical source owner | functional reference/context | selon scope | exact version/current assessment | explicit missing/unknown |
| tenant/environment | Studio or canonical source owner | functional reference/context | selon scope | exact version/current assessment | explicit missing/unknown |
| data/source restrictions | Studio or canonical source owner | functional reference/context | selon scope | exact version/current assessment | explicit missing/unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| version | CMDR Studio | exact version/scope or authorized projection needed by Evaluation Definition | read/link only unless this capability owns the derived Studio record |
| skill | CMDR Studio | exact version/scope or authorized projection needed by Evaluation Definition | read/link only unless this capability owns the derived Studio record |
| workflow | CMDR Studio | exact version/scope or authorized projection needed by Evaluation Definition | read/link only unless this capability owns the derived Studio record |
| automation-agent | CMDR Studio | exact version/scope or authorized projection needed by Evaluation Definition | read/link only unless this capability owns the derived Studio record |
| tenant | Platform Settings | exact version/scope or authorized projection needed by Evaluation Definition | read/link only unless this capability owns the derived Studio record |
| environment | Platform Settings | exact version/scope or authorized projection needed by Evaluation Definition | read/link only unless this capability owns the derived Studio record |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Evaluation Definition | create/update-review/supersede as applicable | CMDR Studio | derived Studio record; source objects and external owners are not mutated implicitly |
| Evaluation Criterion Set | create/update-review/supersede as applicable | CMDR Studio | derived Studio record; source objects and external owners are not mutated implicitly |

## 11. Fonctionnalités
target Studio asset/version; evaluation purpose and scope; criterion set; expected characteristics; required observations; failure/inconclusive conditions; explicit missing/stale/conflict handling; source-version pinning; tenant/environment preservation; auditable supersession.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect Evaluation Definition and source lineage | authorized Studio user | Evaluation Definition | 0 | read permission | current context without mutation | non |
| validate Evaluation Definition prerequisites | Studio reviewer | Evaluation Definition | 1 | source refs available | explainable eligible/blocked/warning result | non |
| create or update Evaluation Definition | Studio owner/reviewer | Evaluation Definition | 2 | functional manage permission | versioned Studio record | OPEN-013 |
| advance effect-bearing lifecycle where applicable | authorized Studio operator | Evaluation Definition | 3 | current permission plus applicable Govern/Settings dependencies | explicit request/transition only; never hidden authority | governed when effectful |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/compare inputs for Evaluation Definition | oui | oui when rule-defined | oui | sourced explanation possible | structured source tables and manual review |
| draft observations/summary for Evaluation Definition | oui | templates/checks | oui | oui, explicitly attributed | human-authored structured record |
| decide authority, invent evidence or bypass a gate | accountable owner only | validation only | non autonome | interdit | explicit human/deterministic governed path |

Toute contribution IA conserve provider/model/version quand disponibles, sources, limites et auteur/reviewer. L’IA ne déclare jamais seule production-safe, n’accorde aucune permission/autorité, ne publie/déploie pas seule, n’invente aucun résultat et ne masque ni échec ni régression.

## 14. États fonctionnels
`draft`, `complete`, `blocked`, `superseded`. Ces états sont fonctionnels et ne remplacent pas les machines d’état canoniques des objets lus.

## 15. États d’interface
Loading conserve asset/version et return origin; Empty explique quelle pièce manque; Partial nomme exactement observations/dépendances absentes; Error conserve le dernier état valide et correlation id; Offline interdit une mutation dont l’écriture/autorité ne peut être garantie; Permission denied masque les données protégées; Stale affiche la dernière version/fraîcheur et force revalidation lorsque nécessaire. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Evaluation Definition | Studio documentary/runtime-lifecycle record | Assurance, Library, Versions & Deployment and authorized consumers | source/version/scope/provenance and limitations remain explicit |
| Evaluation Criterion Set | Studio documentary/runtime-lifecycle record | Assurance, Library, Versions & Deployment and authorized consumers | source/version/scope/provenance and limitations remain explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| canonical Studio/source objects | prerequisites satisfied | Evaluation Definition | exact versions, tenant/environment, permission and provenance refs | source context |
| Evaluation Definition | review identifies gap | source capability / Assurance | missing/contradictory refs and limitations | same Evaluation Definition lineage |
| Evaluation Definition | accepted functional outcome | downstream Studio lifecycle | versioned outcome, warnings, dependencies and provenance | return origin retained |

## 18. Dépendances
CAP-STD-049, CAP-STD-050, CAP-STD-051, Studio canonical objects, Govern/Settings/Shared/Endpoint boundaries, Security permission model, OPEN-013. Les dépendances techniques restent des références et aucune n’est déclarée implémentée par cette capability.

## 19. Source de vérité
La source de vérité pour **Evaluation Definition** est le record Studio créé par cette capability et ses références versionnées. Les faits externes restent chez leurs owners. Une observation, simulation, score, health projection ou technical outcome ne devient pas automatiquement Govern Result, Approval, Decision, Evidence ou état business.

## 20. Provenance et audit
Enregistrer id/version de l’asset et du record, tenant/environment, sources et versions, critères/dépendances, permissions/authority refs, auteur/reviewer, règles ou assistance IA, observations, warnings/failures/partials, décisions humaines locales, timestamps, supersession, return origin et correlation ids. L’historique n’est jamais effacé.

## 21. Permissions fonctionnelles
Besoins : read; create/update du record Studio; exécution no-effect/check lorsqu’applicable; lecture sensible masquée; review; lifecycle transition; provenance/export preparation. Les références `perm.studio.*` et `perm.cmdr-studio.*` restent une anomalie documentée : aucun bulk rename ni choix RBAC final. Toute action effectful réévalue les permissions et dépendances Govern/Settings applicables.

## 22. Limites et erreurs
Version/source stale, tenant mismatch, accès refusé, Secret Reference manquante, provider/runtime indisponible, dependency incompatible, conflicting observations, incomplete evidence, concurrent supersession ou trace indisponible donnent un état explicite blocked/partial/inconclusive/error. Aucun fallback ne peut étendre le scope, inventer une source ou transformer warning en succès.

## 23. Métriques
Comptabiliser records complets/incomplets, stale/version mismatches, blocked permission/dependency cases, partial/inconclusive outcomes, supersessions, manual reviews, AI-assisted reviews identifiables et violations de frontières détectées. Aucune cible numérique/SLO n’est inventée.

## 24. Classification de livraison
`defined` / `planned`. Cette capability ferme une lacune fonctionnelle documentaire STD-4 mais ne prouve aucun engine, API, UI détaillée, stockage, permission atomique, provider/runtime, déploiement réel ou Endpoint capability.

## 25. Critères d’acceptation
**Given** an Evaluation Case lacks an expected outcome, **When** the Evaluation Definition is checked for executability, **Then** the missing oracle is reported and the affected criterion cannot be treated as evaluable.

**Given** a baseline/reference result exists, **When** a new criterion is added, **Then** the baseline remains contextual evidence rather than universal truth.

**Given** restricted source material is referenced, **When** an evaluator lacks source permission, **Then** the evaluation is blocked or masked without expanding access.

**Given** AI is unavailable, **When** the evaluation definition is prepared, **Then** criteria and expected characteristics can be authored and reviewed manually.

## 26. Questions ouvertes
OPEN-013 remain open where consumed. Les 18 décisions OPEN du programme restent inchangées; aucune décision de framework, provider, permission finale, bridge Govern ou support Endpoint n’est prise par rédaction.

## 27. Consommateurs documentaires
Studio Assurance, Evaluations, Simulations, Versions & Deployment, Library/Builder/Skills/Workflows/Agents/Control Room selon besoin, Capability/Object/Dependency/Permission registers, Requirements/baseline, Studio screens as conceptual mappings, Quality/Validation, Roadmap Phase 5, Govern/Settings/Shared/Endpoint boundaries et futurs contrats d’implémentation.
