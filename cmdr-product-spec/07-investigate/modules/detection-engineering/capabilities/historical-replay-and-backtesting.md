---
id: CAP-INV-414
title: Historical Replay and Backtesting
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
  - REQ-AI-010
open_decisions:
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-414 — Historical Replay and Backtesting
## 1. Définition
Préparer et exécuter explicitement un replay historique autorisé d’un Detection Content Draft sur une période, un environnement, des sources ou un dataset, avec rétention, gaps, paramètres, progression, interruption, errors, matches, partiality and provenance visibles.

## 2. Problème utilisateur
Un backtest peut être confondu avec un déploiement, un runtime Detection ou une garantie de performance future, surtout lorsque les données historiques sont incomplètes.

## 3. Objectifs
- sélectionner draft, environment, period, sources and dataset
- voir retention, gaps, limitations and functional estimated cost
- launch, track, cancel and inspect replay explicitly
- compare versions and periods and link test scenarios
- preserve parameters, errors, partial results and provenance

## 4. Non-objectifs
- aucun moteur, langage, syntaxe vendor, API, protocole, parser, compilateur, AST, modèle ML, commande ou code
- aucune promotion, deployment, activation, deactivation, rollback, exception active ou mutation Signal/Alert
- aucune capability CAP-INV-5xx, Intelligence, Cloud/Mobile ou réécriture détaillée d’écran

## 5. Propriétaire
Investigate possède Replay Result et sa disposition humaine. Command conserve runtime Detection/Signal/Alert/Incident ; Settings les sources/parsers/schemas/health/retention ; Endpoint Agent ses capacités et résultats locaux ; Studio Tool/Tool Call/Workflow/Automation Run ; Govern l’autorité future ; Shared les mécanismes génériques.

## 6. Utilisateurs
Principal : **Detection Engineer**. Secondaires : Detection Test Reviewer; Data Access Reviewer.

## 7. Conditions d’entrée
Tenant, environnement, objective, scope, versions, sources, permissions, restrictions et return origin sont explicites. Une dépendance absente produit un état incomplet, partiel ou bloqué ; aucune donnée ou autorité n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | ---: | --- | --- |
| Detection Content Draft and validation | CAP-INV-406/411 | target version, logic and diagnostics | oui | immutable version | replay `blocked` |
| Historical data scope | Event Search / Shared data access | environment, period, sources and retention | oui | freshness and retention snapshot | `retention-limited` |
| Test scenarios and expected outcomes | CAP-INV-412/413 | controlled expectations | non | linked versions | review without oracle |
| Permission and Tool context | Security / Studio | data access, Tool/version and run controls | oui | re-evaluated before run | `permission-blocked` |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Detection Content Draft / Validation | Investigate concepts | selected version | lecture |
| Telemetry Event / Query / Search Job | Shared | historical authorized data and execution context | lecture/search |
| Test Scenario / Expected Outcome | Investigate concepts | comparison context | lecture |
| Tool / Tool Call / Automation Run | Studio | execution and attribution | invoke/read/cancel |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Replay Result | créer, interrompre, comparer, superseder | Investigate concept | not runtime Detection and no Signal |
| Replay parameter record | enregistrer | Investigate / Shared Trace | period, sources and version immutable |
| Background Job / Trace event | request/emit | Shared | progress and errors visible |

## 11. Fonctionnalités
- prepare scoped historical replay
- estimate functional cost without engine choice
- explicitly launch, monitor and cancel
- inspect matches, relevant non-matches, partial and errors
- compare versions/periods and link scenarios
- conserver versions, erreurs, partialité, restrictions et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | ---: | --- | --- | --- |
| Consulter/Comparer | Detection Engineer | Replay Result | 0 | lecture autorisée | projection sourcée | non |
| Exécuter traitement borné | Detection Engineer | Tool Call / Result | 1 | déclenchement explicite et permission | résultat attribué | policy |
| Créer/Modifier/Contester | Detection Engineer | Replay Result | 2 | mutation réversible | nouvelle version | OPEN-013 |
| Préparer handoff | Detection Engineer | candidate package | 2 | sources et limites visibles | package non effectif | destination |

Classes 3/4 exclues ; production et runtime appartiennent à 4B.3A.2/owners.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | ---: | ---: | ---: | ---: | --- |
| Construire Replay Result | oui | éditeur/contrôles explicables | oui | proposition | formulaire/table/revue |
| Valider ou comparer | oui | validateur/comparateur | oui | explication | diagnostics/diff |
| Expliquer erreurs | oui | catalogue | oui | résumé sourcé | erreurs brutes/checklist |
| Promouvoir/déployer/qualifier runtime | non | non | non | interdit | future phase/owner |

Toute sortie expose initiateur, producteur/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun choix silencieux.

## 14. États fonctionnels
`draft`, `queued`, `running`, `partial`, `completed`, `failed`, `cancelled`, `retention-limited`, `data-gap`, `permission-blocked`, `superseded`. États fonctionnels, pas machine objet finale.

## 15. États d’interface
Loading conserve context/version ; Empty distingue absence et interdiction ; Partial expose gaps ; Error conserve le valide ; Offline bloque les nouveaux runs ; Permission denied masque ; Stale distingue ancien/courant ; conflits fournissent diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Replay Result | historical result | CAP-INV-415,416,417 | period, data gaps and version visible |
| Match/non-match projection | result set | Match Review | not Signal or runtime Detection |
| Replay comparison | comparison result | Project / Reviewer | versions and periods attributed |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-406/411 | préparer replay | CAP-INV-414 | version, sources, period, environment, permissions | authoring |
| CAP-INV-414 | replay completed/partial | CAP-INV-415 | matches, relevant non-matches, source events, gaps and parameters | replay |
| CAP-INV-414 | evaluate coverage | CAP-INV-416 | period, platforms, sources, scenarios and limitations | replay |

Transitions conservent ownership, tenant/env, permissions, restrictions, versions, erreurs, provenance et return origin.

## 18. Dépendances
CAP-INV-002/003/008,406,411..413; Shared Query/Jobs/Trace; Studio Tool/Run; OPEN-013/015. Shared Jobs/Trace/Versioning/Linking/Search/Export/Reporting/Collaboration/Comparison/Recovery consommés sans redéfinition.

## 19. Source de vérité
Investigate est source de Replay Result ; tous les objets consommés restent chez leurs owners. Draft ≠ runtime Detection.

## 20. Provenance et audit
Source du besoin, Project, versions, sources/schemas/fields/mappings/logic, Tool/Calls/Runs, paramètres, datasets, résultats, erreurs, interruptions, auteurs, reviewers, dispositions, exports et correlation ID applicables à Replay Result.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | ---: | --- | --- | --- | --- | --- |
| Historical replay prepare/run/cancel/read | resource and sensitive data | 1 | raw values masked by permission | step-up possible | runner/reviewer | Investigate/Shared/Studio | Permissions |
| Cross-tenant replay | tenant isolation | 1/2 | aggregate/minimize | step-up required | independent approval future | Security | Permissions |
| Replay result annotate/link | analytical mutation | 2 | source access rechecked | OPEN-013 | author/reviewer | Investigate | Permissions |

Matrice atomique, namespaces, RBAC/ABAC, step-up et SoD finaux reportés ; permissions production exclues.

## 22. Limites et erreurs
- Historical replay ≠ future performance.
- Backtest match ≠ production Signal; Replay Result ≠ runtime Detection.
- No deployment, activation, shadow, canary, production tuning or runtime performance measurement.
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
### 1. Data gaps
**Given** a historical period contains known gaps  
**When** replay runs  
**Then** result is `partial` or `data-gap`, gaps are visible and non-matches are not certain FNs

### 2. Cancellation
**Given** a long replay is explicitly cancelled  
**When** results are opened  
**Then** completed partial outputs and cancellation provenance remain visible

### 3. Sans IA
**Given** aucun modèle  
**When** replay runs  
**Then** deterministic execution, progress, comparison and human review remain available

## 26. Questions ouvertes
- OPEN-013 reste ouverte.
- OPEN-015 reste ouverte.
- Le moteur/langage Detection est une lacune future non couverte ; OPEN-005 reste forensic-only.
- Schémas, formats, permissions et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering module, Event Search/Hunt/Case/Evidence/technical handoffs, Command boundaries, Settings/Endpoint, Studio, Govern future review, Shared, Objects/Permissions/Screens/Journeys/Technique/4B.3A.2/Validation.
