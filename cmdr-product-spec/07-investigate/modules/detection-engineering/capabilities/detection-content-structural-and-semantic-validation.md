---
id: CAP-INV-411
title: Detection Content Structural and Semantic Validation
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
# CAP-INV-411 — Detection Content Structural and Semantic Validation
## 1. Définition
Vérifier déterministement complétude, références, sources, fields, types, conditions, relations, fenêtres, paramètres, exclusions, enrichissements, permissions, dépendances, obsolescence, ambiguïtés et contradictions, en séparant validation structurelle et revue sémantique.

## 2. Problème utilisateur
Une syntaxe ou structure valide peut porter une logique trop large, incohérente ou inutile ; une validation fonctionnelle ne prédit pas la performance de production.

## 3. Objectifs
- vérifier complétude documentaire et références
- vérifier sources, fields, types, logic, windows, exclusions and enrichments
- détecter missing dependencies, stale refs, ambiguities and contradictions
- classer et attribuer errors/warnings
- relancer, comparer validations et préparer tests/replay

## 4. Non-objectifs
- aucun moteur, langage, syntaxe vendor, API, protocole, parser, compilateur, AST, modèle ML, commande ou code
- aucune promotion, deployment, activation, deactivation, rollback, exception active ou mutation Signal/Alert
- aucune capability CAP-INV-5xx, Intelligence, Cloud/Mobile ou réécriture détaillée d’écran

## 5. Propriétaire
Investigate possède Validation Result et sa disposition humaine. Command conserve runtime Detection/Signal/Alert/Incident ; Settings les sources/parsers/schemas/health/retention ; Endpoint Agent ses capacités et résultats locaux ; Studio Tool/Tool Call/Workflow/Automation Run ; Govern l’autorité future ; Shared les mécanismes génériques.

## 6. Utilisateurs
Principal : **Detection Validator**. Secondaires : Detection Engineer; Data Source Reviewer; Content Owner.

## 7. Conditions d’entrée
Tenant, environnement, objective, scope, versions, sources, permissions, restrictions et return origin sont explicites. Une dépendance absente produit un état incomplet, partiel ou bloqué ; aucune donnée ou autorité n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | ---: | --- | --- |
| Detection Content Draft and metadata | CAP-INV-406..410 | version and full functional specification | oui | immutable selected version | validation `blocked` |
| Readiness and schema reviews | CAP-INV-404/405 | sources, gaps, fields and mappings | oui | linked snapshots | validation partial |
| Permission context | Security / owners | data and execution rights | oui | run time | permission-blocked |
| Validator / Tool context | Studio Tool / Tool Call | deterministic checks and version | oui for execution | re-evaluated | tool-unavailable |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Detection Content Draft / Metadata | Investigate concepts | selected version | lecture |
| Readiness / Mapping / Enrichment requirements | Investigate concepts | dependencies and limits | lecture |
| Tool / Tool Call / Automation Run | Studio | validator identity, version and outputs | invoke/read |
| Data Source / Schema projections | Settings / Shared | current refs and quality | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Validation Result | créer, relancer, comparer, superseder | Investigate concept | structure and semantic review dimensions separated |
| Validation diagnostics | enregistrer and disposition | Investigate concept | source location and severity visible |
| Trace / Job event | émettre | Shared | run, parameters and errors retained |

## 11. Fonctionnalités
- run completeness/reference checks
- run source/field/type/condition/window checks
- detect stale/missing/ambiguous/contradictory dependencies
- record structural and semantic-review dimensions separately
- retry and compare results
- conserver versions, erreurs, partialité, restrictions et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | ---: | --- | --- | --- |
| Consulter/Comparer | Detection Validator | Validation Result | 0 | lecture autorisée | projection sourcée | non |
| Exécuter traitement borné | Detection Validator | Tool Call / Result | 1 | déclenchement explicite et permission | résultat attribué | policy |
| Créer/Modifier/Contester | Detection Validator | Validation Result | 2 | mutation réversible | nouvelle version | OPEN-013 |
| Préparer handoff | Detection Validator | candidate package | 2 | sources et limites visibles | package non effectif | destination |

Classes 3/4 exclues ; production et runtime appartiennent à 4B.3A.2/owners.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | ---: | ---: | ---: | ---: | --- |
| Construire Validation Result | oui | éditeur/contrôles explicables | oui | proposition | formulaire/table/revue |
| Valider ou comparer | oui | validateur/comparateur | oui | explication | diagnostics/diff |
| Expliquer erreurs | oui | catalogue | oui | résumé sourcé | erreurs brutes/checklist |
| Promouvoir/déployer/qualifier runtime | non | non | non | interdit | future phase/owner |

Toute sortie expose initiateur, producteur/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun choix silencieux.

## 14. États fonctionnels
`not-run`, `queued`, `validating`, `valid`, `valid-with-warnings`, `invalid`, `blocked`, `partial`, `failed`, `superseded`. États fonctionnels, pas machine objet finale.

## 15. États d’interface
Loading conserve context/version ; Empty distingue absence et interdiction ; Partial expose gaps ; Error conserve le valide ; Offline bloque les nouveaux runs ; Permission denied masque ; Stale distingue ancien/courant ; conflits fournissent diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Validation Result | validation assessment | CAP-INV-412,414,417 | structural and semantic dimensions explicit |
| Diagnostics | versioned issues | Authoring / Reviewer | source and disposition visible |
| Validation comparison | comparison result | Project | different versions/runs attributed |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-406/410 | lancer validation | CAP-INV-411 | version, sources, fields, logic, metadata, dependencies | authoring |
| CAP-INV-411 | valid or warnings | CAP-INV-412/414 | result, warnings, version and limitations | validation |
| CAP-INV-411 | invalid or blocked | CAP-INV-406 | diagnostics, missing dependencies and return origin | validation |

Transitions conservent ownership, tenant/env, permissions, restrictions, versions, erreurs, provenance et return origin.

## 18. Dépendances
CAP-INV-404..410; Studio Tool/Tool Call/Automation Run; Shared Jobs/Trace; OPEN-013/015. Shared Jobs/Trace/Versioning/Linking/Search/Export/Reporting/Collaboration/Comparison/Recovery consommés sans redéfinition.

## 19. Source de vérité
Investigate est source de Validation Result ; tous les objets consommés restent chez leurs owners. Draft ≠ runtime Detection.

## 20. Provenance et audit
Source du besoin, Project, versions, sources/schemas/fields/mappings/logic, Tool/Calls/Runs, paramètres, datasets, résultats, erreurs, interruptions, auteurs, reviewers, dispositions, exports et correlation ID applicables à Validation Result.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | ---: | --- | --- | --- | --- | --- |
| Validation run/read | processing and sensitive metadata | 1 | values masked | possible | runner/reviewer | Investigate/Studio | Permissions |
| Validation result disposition | quality decision | 2 | no hidden diagnostics | OPEN-013 | validator/content owner | Investigate | Permissions |
| Automated validation request | tool risk | 1 | Tool/version visible | possible | requester/reviewer | Studio | Permissions |

Matrice atomique, namespaces, RBAC/ABAC, step-up et SoD finaux reportés ; permissions production exclues.

## 22. Limites et erreurs
- Structural validation ≠ functional usefulness.
- Semantic review ≠ production performance.
- No compiler, parser, engine, language, runtime or production approval is defined.
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
### 1. Structure valid, logic too broad
**Given** un draft structurally valid avec condition sémantiquement large  
**When** la validation se termine  
**Then** structure may be valid, warnings remain visible and no `review-ready` is automatic

### 2. Dependency stale
**Given** un field mapping superseded  
**When** validation runs  
**Then** the stale reference is named and result may be invalid/partial

### 3. Sans IA
**Given** aucun modèle  
**When** validation runs  
**Then** deterministic validators, diagnostics and human semantic review remain available

## 26. Questions ouvertes
- OPEN-013 reste ouverte.
- OPEN-015 reste ouverte.
- Le moteur/langage Detection est une lacune future non couverte ; OPEN-005 reste forensic-only.
- Schémas, formats, permissions et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering module, Event Search/Hunt/Case/Evidence/technical handoffs, Command boundaries, Settings/Endpoint, Studio, Govern future review, Shared, Objects/Permissions/Screens/Journeys/Technique/4B.3A.2/Validation.
