---
id: CAP-INV-415
title: Match Review and False Positive/False Negative Analysis
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
# CAP-INV-415 — Match Review and False Positive/False Negative Analysis
## 1. Définition
Examiner matches et non-matches attendus en les reliant aux événements sources, scénarios, Expected Outcomes, contradictions et limites de ground truth, puis enregistrer des dispositions candidates TP, FP, FN ou inconclusive.

## 2. Problème utilisateur
Un match n’est pas une activité malveillante confirmée et un non-match n’est pas un faux négatif certain lorsque données, labels ou contextes sont incomplets.

## 3. Objectifs
- review matches and expected non-matches
- link source events, scenarios and Expected Outcomes
- classify candidate TP/FP/FN, ambiguous or insufficient context
- document justification, reviewer and ground-truth limits
- identify noise/logic gaps, compare versions/periods and propose new draft

## 4. Non-objectifs
- aucun moteur, langage, syntaxe vendor, API, protocole, parser, compilateur, AST, modèle ML, commande ou code
- aucune promotion, deployment, activation, deactivation, rollback, exception active ou mutation Signal/Alert
- aucune capability CAP-INV-5xx, Intelligence, Cloud/Mobile ou réécriture détaillée d’écran

## 5. Propriétaire
Investigate possède Match Review et sa disposition humaine. Command conserve runtime Detection/Signal/Alert/Incident ; Settings les sources/parsers/schemas/health/retention ; Endpoint Agent ses capacités et résultats locaux ; Studio Tool/Tool Call/Workflow/Automation Run ; Govern l’autorité future ; Shared les mécanismes génériques.

## 6. Utilisateurs
Principal : **Detection Reviewer**. Secondaires : Detection Engineer; Threat Hunter; Investigation Analyst.

## 7. Conditions d’entrée
Tenant, environnement, objective, scope, versions, sources, permissions, restrictions et return origin sont explicites. Une dépendance absente produit un état incomplet, partiel ou bloqué ; aucune donnée ou autorité n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | ---: | --- | --- |
| Replay or validation results | CAP-INV-411/414 | matches, non-matches, errors and parameters | oui | selected result version | review `unreviewed` |
| Test Scenario / Expected Outcome | CAP-INV-412/413 | expected behavior and confidence | non | linked versions | review without oracle |
| Source events and investigation context | Shared / Case / Evidence / Finding | context and provenance | oui where accessible | `insufficient-context` |
| Data readiness and gaps | CAP-INV-404 | coverage, retention and quality limits | oui | assessment for period | classification uncertainty |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Replay / Validation Result | Investigate concepts | actual results and errors | lecture |
| Telemetry Event | Shared | source event projection | lecture |
| Test Scenario / Expected Outcome | Investigate concepts | expected behavior and confidence | lecture |
| Case / Evidence / Finding | Investigate | qualified context when available | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Match Review | créer, classifier, contester, superseder | Investigate concept | candidate labels only |
| FP/FN candidate disposition | enregistrer | Investigate concept | justification and uncertainty mandatory |
| Draft change proposal | préparer | CAP-INV-406 | new version, never active mutation |

## 11. Fonctionnalités
- inspect matches and expected non-matches
- compare actual vs expected and source context
- classify candidate TP/FP/FN or inconclusive
- record justification/reviewer/ground-truth limits
- identify noise/gaps and propose new draft
- conserver versions, erreurs, partialité, restrictions et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | ---: | --- | --- | --- |
| Consulter/Comparer | Detection Reviewer | Match Review | 0 | lecture autorisée | projection sourcée | non |
| Exécuter traitement borné | Detection Reviewer | Tool Call / Result | 1 | déclenchement explicite et permission | résultat attribué | policy |
| Créer/Modifier/Contester | Detection Reviewer | Match Review | 2 | mutation réversible | nouvelle version | OPEN-013 |
| Préparer handoff | Detection Reviewer | candidate package | 2 | sources et limites visibles | package non effectif | destination |

Classes 3/4 exclues ; production et runtime appartiennent à 4B.3A.2/owners.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | ---: | ---: | ---: | ---: | --- |
| Construire Match Review | oui | éditeur/contrôles explicables | oui | proposition | formulaire/table/revue |
| Valider ou comparer | oui | validateur/comparateur | oui | explication | diagnostics/diff |
| Expliquer erreurs | oui | catalogue | oui | résumé sourcé | erreurs brutes/checklist |
| Promouvoir/déployer/qualifier runtime | non | non | non | interdit | future phase/owner |

Toute sortie expose initiateur, producteur/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun choix silencieux.

## 14. États fonctionnels
`unreviewed`, `candidate-true-positive`, `candidate-false-positive`, `candidate-false-negative`, `expected-non-match`, `ambiguous`, `insufficient-context`, `disputed`, `superseded`. États fonctionnels, pas machine objet finale.

## 15. États d’interface
Loading conserve context/version ; Empty distingue absence et interdiction ; Partial expose gaps ; Error conserve le valide ; Offline bloque les nouveaux runs ; Permission denied masque ; Stale distingue ancien/courant ; conflits fournissent diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Match Review | review disposition | CAP-INV-406,416,417 | candidate status and justification visible |
| Draft change proposal | change context | CAP-INV-406 | source version and provenance retained |
| Noise or logic gap | gap candidate | CAP-INV-416 | data vs logic cause not conflated |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-411/414 | review results | CAP-INV-415 | actual results, events, scenarios, expected outcomes and limits | validation/replay |
| CAP-INV-415 | propose modification | CAP-INV-406 | candidate FP/FN, rationale, source version and provenance | review |
| CAP-INV-415 | record gap | CAP-INV-416 | gap type, impact, data/logic evidence and uncertainty | review |

Transitions conservent ownership, tenant/env, permissions, restrictions, versions, erreurs, provenance et return origin.

## 18. Dépendances
CAP-INV-404,406,411..414; Case/Evidence/Finding; Shared event access; OPEN-013/015. Shared Jobs/Trace/Versioning/Linking/Search/Export/Reporting/Collaboration/Comparison/Recovery consommés sans redéfinition.

## 19. Source de vérité
Investigate est source de Match Review ; tous les objets consommés restent chez leurs owners. Draft ≠ runtime Detection.

## 20. Provenance et audit
Source du besoin, Project, versions, sources/schemas/fields/mappings/logic, Tool/Calls/Runs, paramètres, datasets, résultats, erreurs, interruptions, auteurs, reviewers, dispositions, exports et correlation ID applicables à Match Review.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | ---: | --- | --- | --- | --- | --- |
| Match review read/create/update | potentially sensitive event context | 0/2 | raw fields masked | OPEN-013 | reviewer/content owner | Investigate | Permissions |
| FP/FN candidate classify | quality label risk | 2 | uncertainty mandatory | possible | independent reviewer if required | Investigate | Permissions |
| Source event read | production data | 0 | field-level controls | step-up possible | reviewer/data owner | Shared/Settings | Permissions |

Matrice atomique, namespaces, RBAC/ABAC, step-up et SoD finaux reportés ; permissions production exclues.

## 22. Limites et erreurs
- Match ≠ confirmed malicious activity; suspected TP ≠ confirmed TP.
- Non-match ≠ certain FN; suspected FP/FN remain candidates.
- Estimated precision/recall are not guarantees; no observed FP ≠ zero FP.
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
### 1. Candidate FP
**Given** a match with possible benign context and incomplete truth  
**When** reviewer examines it  
**Then** status may be candidate-FP, justification/uncertainty/source remain visible and no certainty is imposed

### 2. Candidate FN
**Given** a positive scenario has no match and data may be missing  
**When** result is reviewed  
**Then** data and Expected Outcome are checked and status may remain insufficient-context

### 3. Sans IA
**Given** aucun modèle  
**When** match review is performed  
**Then** tables, source pivots, comparison and human classification remain available

## 26. Questions ouvertes
- OPEN-013 reste ouverte.
- OPEN-015 reste ouverte.
- Le moteur/langage Detection est une lacune future non couverte ; OPEN-005 reste forensic-only.
- Schémas, formats, permissions et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering module, Event Search/Hunt/Case/Evidence/technical handoffs, Command boundaries, Settings/Endpoint, Studio, Govern future review, Shared, Objects/Permissions/Screens/Journeys/Technique/4B.3A.2/Validation.
