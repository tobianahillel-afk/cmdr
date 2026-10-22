---
id: CAP-INV-416
title: Detection Coverage Mapping and Gap Analysis
product: investigate
module: detection-engineering
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-06
requirement_ids:
  - REQ-INV-006
  - REQ-PROD-014
  - REQ-PROD-019
  - REQ-PROD-020
  - REQ-AI-002
  - REQ-SEC-001
  - REQ-PROD-055
open_decisions:
  - OPEN-008
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-416 — Detection Coverage Mapping and Gap Analysis
## 1. Définition
Cartographier fonctionnellement comportements, sources, plateformes, environnements, taxonomies, étapes, exclusions, tests, observations historiques et dépendances, puis identifier et versionner les gaps sans présenter un mapping comme couverture effective.

## 2. Problème utilisateur
Un mapping technique ou des tests synthétiques peuvent être présentés à tort comme couverture réelle, protection ou performance en production.

## 3. Objectifs
- définir behavior, sources, platforms, environments and mappings
- définir steps/conditions covered, exclusions, limits and dependencies
- voir linked Detection Contents, Hunts, Findings, missing sources/fields/scenarios
- classify coverage level and gaps
- compare environments and prepare backlog/future handoff

## 4. Non-objectifs
- aucun moteur, langage, syntaxe vendor, API, protocole, parser, compilateur, AST, modèle ML, commande ou code
- aucune promotion, deployment, activation, deactivation, rollback, exception active ou mutation Signal/Alert
- aucune capability CAP-INV-5xx, Intelligence, Cloud/Mobile ou réécriture détaillée d’écran

## 5. Propriétaire
Investigate possède Detection Coverage Assessment and Detection Gap et sa disposition humaine. Command conserve runtime Detection/Signal/Alert/Incident ; Settings les sources/parsers/schemas/health/retention ; Endpoint Agent ses capacités et résultats locaux ; Studio Tool/Tool Call/Workflow/Automation Run ; Govern l’autorité future ; Shared les mécanismes génériques.

## 6. Utilisateurs
Principal : **Detection Coverage Analyst**. Secondaires : Detection Engineer; Threat Hunter; Platform Owner.

## 7. Conditions d’entrée
Tenant, environnement, objective, scope, versions, sources, permissions, restrictions et return origin sont explicites. Une dépendance absente produit un état incomplet, partiel ou bloqué ; aucune donnée ou autorité n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | ---: | --- | --- |
| Detection Content Draft and Hypothesis | CAP-INV-403/406 | behavior, conditions, sources and limits | oui | selected versions | coverage `unknown` |
| Readiness / Mapping / Test / Replay / Review | CAP-INV-404,405,412..415 | data and evidence of validation | oui | linked snapshots | coverage `planned` or partial |
| Taxonomy projections | source owner / Shared Mapping | candidate technique or behavior mappings | non | version visible | no taxonomy claim |
| Environment/platform context | Settings / Endpoint Agent | scope and support projections | oui | current assessment | coverage-limited |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Detection Content / Hypothesis | Investigate concepts | scope and claimed behavior | lecture |
| Readiness, tests, replay and reviews | Investigate concepts | support and limitations | lecture |
| Hunt / Finding | Investigate | related observations and gaps | lecture/lien |
| Platform / taxonomy projections | Settings / Shared | scope and mappings | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Detection Coverage Assessment | créer, comparer, contester, superseder | Investigate concept | mapping ≠ effective coverage |
| Detection Gap | créer, update, prioritize-propose, superseder | Investigate concept | approval/backlog owner remains future |
| Backlog/review context | préparer | future owner | no automatic requirement approval |

## 11. Fonctionnalités
- map behavior, data, platform and taxonomy context
- evaluate tests/historical observations and limitations
- identify missing sources, fields, scenarios and logic
- classify coverage states
- compare environments and prepare gaps/backlog
- conserver versions, erreurs, partialité, restrictions et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | ---: | --- | --- | --- |
| Consulter/Comparer | Detection Coverage Analyst | Detection Coverage Assessment and Detection Gap | 0 | lecture autorisée | projection sourcée | non |
| Exécuter traitement borné | Detection Coverage Analyst | Tool Call / Result | 1 | déclenchement explicite et permission | résultat attribué | policy |
| Créer/Modifier/Contester | Detection Coverage Analyst | Detection Coverage Assessment and Detection Gap | 2 | mutation réversible | nouvelle version | OPEN-013 |
| Préparer handoff | Detection Coverage Analyst | candidate package | 2 | sources et limites visibles | package non effectif | destination |

Classes 3/4 exclues ; production et runtime appartiennent à 4B.3A.2/owners.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | ---: | ---: | ---: | ---: | --- |
| Construire Detection Coverage Assessment and Detection Gap | oui | éditeur/contrôles explicables | oui | proposition | formulaire/table/revue |
| Valider ou comparer | oui | validateur/comparateur | oui | explication | diagnostics/diff |
| Expliquer erreurs | oui | catalogue | oui | résumé sourcé | erreurs brutes/checklist |
| Promouvoir/déployer/qualifier runtime | non | non | non | interdit | future phase/owner |

Toute sortie expose initiateur, producteur/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun choix silencieux.

## 14. États fonctionnels
`unknown`, `no-coverage`, `planned`, `partial`, `conditionally-covered`, `test-covered`, `historically-observed`, `gap-identified`, `disputed`, `superseded`. États fonctionnels, pas machine objet finale.

## 15. États d’interface
Loading conserve context/version ; Empty distingue absence et interdiction ; Partial expose gaps ; Error conserve le valide ; Offline bloque les nouveaux runs ; Permission denied masque ; Stale distingue ancien/courant ; conflits fournissent diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Detection Coverage Assessment | coverage assessment | Project / CAP-INV-417 | evidence level and limitations visible |
| Detection Gap | gap | Settings backlog / future review | impact and owner probable visible |
| Coverage comparison | comparison result | Reviewer | environment differences explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-403/406 | assess coverage | CAP-INV-416 | behavior, sources, platforms, mappings and limitations | authoring |
| CAP-INV-412..415 | update evidence | CAP-INV-416 | tests, replay, reviews, gaps and confidence | tests/review |
| CAP-INV-416 | prepare review | CAP-INV-417 | assessment, gaps, sources, risks and unresolved items | coverage |

Transitions conservent ownership, tenant/env, permissions, restrictions, versions, erreurs, provenance et return origin.

## 18. Dépendances
CAP-INV-005,109,403..415; Settings/Endpoint projections; Shared Mapping/Comparison; OPEN-008/013/015. Shared Jobs/Trace/Versioning/Linking/Search/Export/Reporting/Collaboration/Comparison/Recovery consommés sans redéfinition.

## 19. Source de vérité
Investigate est source de Detection Coverage Assessment and Detection Gap ; tous les objets consommés restent chez leurs owners. Draft ≠ runtime Detection.

## 20. Provenance et audit
Source du besoin, Project, versions, sources/schemas/fields/mappings/logic, Tool/Calls/Runs, paramètres, datasets, résultats, erreurs, interruptions, auteurs, reviewers, dispositions, exports et correlation ID applicables à Detection Coverage Assessment and Detection Gap.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | ---: | --- | --- | --- | --- | --- |
| Coverage assessment read/create/update | overclaim and business impact | 0/2 | sensitive scope masked | OPEN-013 | author/reviewer | Investigate | Permissions |
| Detection Gap create/update | backlog impact | 2 | no automatic approval | possible | analyst/owner | Investigate/future owner | Permissions |
| Cross-environment comparison | tenant/platform exposure | 1/2 | aggregate/minimize | step-up possible | reviewer independent | Security/Settings | Permissions |

Matrice atomique, namespaces, RBAC/ABAC, step-up et SoD finaux reportés ; permissions production exclues.

## 22. Limites et erreurs
- Coverage mapping ≠ effective coverage; technique mapping ≠ proof.
- Data coverage ≠ detection coverage; detection coverage ≠ prevention.
- `test-covered` ≠ effective in production; gap identified ≠ requirement approved.
- source stale/restricted/partial, tenant mismatch, permission revoked, Tool/version unavailable, timeout or cancellation
- Tool result, score, match, non-match, AI output or mapping is not a conclusion by itself

## 23. Métriques conceptuelles
- volume par état/version
- partial/blocked/disputed/failed
- provenance et dispositions humaines complètes
- silent promotion/deployment count — cible zéro

Aucune cible runtime, precision/recall garantie, drift, health ou coût production.

## 24. Classification de livraison
`defined` / `planned` ; documentation only. Aucun `validated`, `implemented`, `deployed`, `active`, `native` ou `integrated`.

## 25. Critères d’acceptation
### 1. Synthetic-only
**Given** a draft mapped to a behavior using only synthetic tests  
**When** coverage is assessed  
**Then** state may be `test-covered`, not production-effective, with platforms and limitations visible

### 2. Different environments
**Given** one environment lacks a required source  
**When** comparison is opened  
**Then** coverage differs, gap is explicit and no global coverage is claimed

### 3. Sans IA
**Given** aucun modèle  
**When** coverage is mapped  
**Then** matrices, catalogs, tests, replay evidence and human review suffice

## 26. Questions ouvertes
- OPEN-008 reste ouverte.
- OPEN-013 reste ouverte.
- OPEN-015 reste ouverte.
- Le moteur/langage Detection est une lacune future non couverte ; OPEN-005 reste forensic-only.
- Schémas, formats, permissions et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering module, Event Search/Hunt/Case/Evidence/technical handoffs, Command boundaries, Settings/Endpoint, Studio, Govern future review, Shared, Objects/Permissions/Screens/Journeys/Technique/4B.3A.2/Validation.
