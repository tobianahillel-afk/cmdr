---
id: CAP-GOV-017
title: Response Playbook Catalog and Selection
product: govern
module: playbooks
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-006, REQ-PROD-008, REQ-PROD-015, REQ-PROD-016, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-GOV-017 — Response Playbook Catalog and Selection

## 1. Définition
Permettre à Govern d’identifier, rechercher, comparer et sélectionner un **Response Playbook** candidat pour une Decision et son Execution Handoff Package, en évaluant version, action types, target types, environnements, dépendances runtime, reversibility, rollback/verification support, restrictions et disponibilité sans confondre Playbook, Workflow ou autorisation d’exécution.

## 2. Problème utilisateur
Une Decision autorisant une action ne détermine pas automatiquement la procédure d’exécution. Choisir un Playbook incompatible, déprécié, indisponible ou destiné à un autre target type peut produire un plan incorrect. Une simple sélection de catalogue ne doit jamais devenir une autorisation ni un démarrage de Run.

## 3. Objectifs
- fournir un catalogue Govern de Response Playbooks versionnés ;
- filtrer les candidats par action type, target type, tenant/environment et contraintes déclarées ;
- afficher owner, version, availability, risk class, reversibility, rollback et verification support ;
- conserver les références Workflow/Tool/runtime sans transfert d’ownership ;
- permettre compare/select/reject/annotate/request-alternative avec rationale ;
- produire une sélection candidate traçable pour CAP-GOV-018.

## 4. Non-objectifs
Ne pas authorer un Studio Workflow, déployer un Tool, choisir un provider/runtime global, résoudre les targets, lier des secrets, autoriser l’exécution, créer un Response Run, exécuter une action ou finaliser le schéma physique du Playbook.

## 5. Propriétaire
Govern / Playbooks possède la sémantique du Response Playbook et sa sélection pour réponse. CMDR Studio conserve Workflow, Workflow Version, Tool, Tool Call et Automation Run. Settings conserve providers/integrations/runtime configuration ; Endpoint conserve ses capabilities techniques.

## 6. Utilisateurs
Principal : Response Operator / Govern Reviewer. Secondaires : Decision Maker, Playbook Owner, Runtime/Endpoint Operator, Studio Operator comme owner de références, Security Reviewer, Auditor.

## 7. Conditions d’entrée
Execution Handoff Package de CAP-GOV-016, Decision autorisante non expirée comme contexte, tenant/environment, action et target types déclarés, accès au catalogue Playbook et aux projections de disponibilité/dépendances autorisées.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Execution Handoff Package | CAP-GOV-016 | authorized bounds package | oui | current package/Decision | aucune sélection exploitable |
| Decision action/scope/conditions | CAP-GOV-015/016 | selection constraints | oui | Decision current | selection blocked |
| Playbook catalog/version metadata | Govern | catalog candidates | oui | current catalog snapshot | no candidate invented |
| target/environment types | Decision + source projections | compatibility dimensions | oui | handoff version | candidate marked unknown/incompatible |
| runtime/provider dependencies | Settings/Studio/Endpoint metadata | dependency context | selon Playbook | latest known | availability unknown |
| rollback/verification metadata | Playbook/version | safety capability metadata | selon action/Decision | exact Playbook version | explicit unsupported/unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Decision | Govern | action, scope, conditions, expiry | read |
| Execution Handoff Package | Govern | exact target/scope/requirements | read |
| Playbook | Govern | candidate metadata/version/restrictions | read/select candidate |
| Workflow / Workflow Version | CMDR Studio | optional implementation reference | read/link only |
| Tool | CMDR Studio | required tool reference/capability | read/link only |
| Integration/Environment/Secret Reference metadata | Platform Settings | availability/dependency metadata | restricted read, no value |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Playbook Selection Context | create/update/version/supersede | Govern local concept | candidate selection ≠ authorization |
| Playbook candidate relation | add/reject/annotate | Govern | exact version and rationale retained |
| Playbook | no authoring mutation by selection | Govern | catalog source remains versioned |
| Workflow/Tool/Settings objects | no mutation | source owners | references only |

## 11. Fonctionnalités
Search/filter catalog; inspect owner/version/intended action types/target types/environments; compare risk/reversibility/rollback/verification support; expose required Workflow/Tool/runtime dependencies; show deprecation/availability/restrictions; reject incompatible candidate; annotate rationale; request alternative; select one or more candidates for compatibility review without creating execution authority.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| search/inspect Playbooks | Response Operator | Playbook catalog | 0 | catalog read | candidate list | non |
| compare exact versions | Govern Reviewer | candidates | 1 | metadata accessible | bounded comparison | non |
| select candidate | Response Operator | Selection Context | 2 | rationale + exact version | candidate selected | OPEN-013 |
| reject/annotate/request alternative | reviewer | Selection Context | 2 | reason | versioned disposition | OPEN-013 |
| publish/edit Workflow or execute | none via this capability | Studio/runtime objects | — | outside owner | no mutation | source owner only |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| filter candidates by action/target | oui | oui | oui | candidate suggestion | catalog filters |
| compare version metadata | oui | deterministic diff | oui | sourced summary | comparison table |
| rank by declared compatibility | oui | explicit criteria | oui | recommendation only | matrix/filter |
| explain missing dependencies | oui | oui | oui | explanation | dependency checklist |
| authorize/start execution | authorized Govern path later | validation only | no autonomous | prohibited | CAP-GOV-021..023 |

## 14. États fonctionnels
`not-selected`, `candidates-available`, `no-candidate`, `candidate-selected`, `candidate-rejected`, `alternative-requested`, `dependency-unknown`, `deprecated-candidate`, `unavailable-candidate`, `compatibility-review-required`, `superseded`.

## 15. États d’interface
Loading conserve handoff/Decision context ; Empty distingue aucun candidat d’un accès refusé ; Partial nomme metadata/dépendances absentes ; Error conserve la dernière comparaison valide ; Offline est read-only ; Permission denied masque les dépendances protégées ; Stale exige refresh avant sélection autoritative.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| candidate selection | Playbook Selection Context | CAP-GOV-018 | exact Playbook/version + rationale |
| candidate comparison | review result | Response Operator/Decision reviewer | sources and unknowns explicit |
| no-compatible-candidate condition | blocker | CAP-GOV-018/Response Inbox | no silent fallback/substitution |
| alternative request | review event | Playbook owner/operator | requirements and return origin retained |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-016 | handoff accepted for planning | CAP-GOV-017 | Decision, action, targets, scope, conditions | Decision Register |
| CAP-GOV-017 | candidate selected | CAP-GOV-018 | exact Playbook/version/dependencies/rationale | Playbooks |
| CAP-GOV-017 | no candidate | source/Playbook owner | action/target requirements and blocker | Playbooks |
| Playbook metadata changes | candidate stale | CAP-GOV-017/018 | previous/new version refs | same selection lineage |

## 18. Dépendances
CAP-GOV-015/016/018/019, Playbook object, Studio Workflow/Tool metadata, Settings integrations/environments/health/Secret Reference metadata, Endpoint capability projections, Shared Search/Versioning/Trace/Linking, OPEN-008/013/015.

## 19. Source de vérité
Govern Playbook/version is source of response-procedure semantics. Studio Workflow remains source of Studio orchestration. Settings/Endpoint sources determine current dependency availability. Selection Context is source only for the selected candidate, never for execution authority.

## 20. Provenance et audit
Record Decision/handoff versions, search/filter context, candidate ids/versions, metadata sources/timestamps, dependency availability, deprecation/restrictions, comparisons, rejected alternatives, selection rationale, reviewer, automation/AI recommendation producer and timestamps.

## 21. Permissions fonctionnelles
Playbook read/select, version compare, restricted dependency metadata read, automated recommendation request, provenance export. Playbook management may exist separately; Workflow/Tool/secret/runtime administration remains source-owned. Selection may require class-2 governance under OPEN-013.

## 22. Limites et erreurs
No candidate, stale catalog, deprecated version, unavailable dependency, unknown rollback/verification support, tenant mismatch, inaccessible metadata or conflicting target type must remain explicit. No candidate is fabricated and no version is silently substituted.

## 23. Métriques
Candidates per handoff; no-compatible-candidate rate; deprecated/unavailable candidates encountered; selection changes after compatibility review; selection-to-compatibility time; executions started by selection alone — target zero.

## 24. Classification de livraison
`defined` / `planned`; provider/runtime-neutral catalog semantics only. No execution engine, provider, Workflow runtime, API or protocol selected.

## 25. Critères d’acceptation
**Given** a Decision authorizes an endpoint action but the only candidate Playbook targets a different resource type, **When** selection is reviewed, **Then** the mismatch is explicit and the candidate cannot become an authorized execution path.

**Given** a selected Playbook version is deprecated before planning, **When** the selection is reopened, **Then** its exact version remains visible and a replacement requires explicit selection plus compatibility review.

**Given** no AI, **When** a Playbook is selected, **Then** catalog filters, comparison tables, dependency metadata and human rationale provide the full function.

## 26. Questions ouvertes
OPEN-008 remains open for actual platform/runtime support, OPEN-013 for class-2 selection mutations, and OPEN-015 for Workflow/Automation Run/Response Run bridging. No new OPEN is created.

## 27. Consommateurs documentaires
Playbooks module, CAP-GOV-018/019, Runs & Rollback, future Objects/Permissions/Screens/Journeys/Technique, GOV-2 conformance report and GOV-3 audit/metrics consumers.
