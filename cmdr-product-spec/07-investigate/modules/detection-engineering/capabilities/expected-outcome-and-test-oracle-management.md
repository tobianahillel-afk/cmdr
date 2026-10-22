---
id: CAP-INV-413
title: Expected Outcome and Test Oracle Management
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
# CAP-INV-413 — Expected Outcome and Test Oracle Management
## 1. Définition
Définir, attribuer, versionner et faire revoir les matches, non-matches, options, fields, relations, counts, windows, warnings, blocking errors et partial states attendus, sans confondre attendu, résultat réel ou ground truth.

## 2. Problème utilisateur
Un Expected Outcome produit sans source, confiance ou reviewer peut devenir un test oracle autoritaire, notamment lorsqu’il est généré par IA.

## 3. Objectifs
- définir match/non-match/optional match expectations
- définir expected fields, relations, counts and windows
- définir acceptable warnings, blocking errors and partial status
- attribuer source, author, confidence and contradictions
- versionner, contester, confirmer par revue and superseder

## 4. Non-objectifs
- aucun moteur, langage, syntaxe vendor, API, protocole, parser, compilateur, AST, modèle ML, commande ou code
- aucune promotion, deployment, activation, deactivation, rollback, exception active ou mutation Signal/Alert
- aucune capability CAP-INV-5xx, Intelligence, Cloud/Mobile ou réécriture détaillée d’écran

## 5. Propriétaire
Investigate possède Expected Outcome et sa disposition humaine. Command conserve runtime Detection/Signal/Alert/Incident ; Settings les sources/parsers/schemas/health/retention ; Endpoint Agent ses capacités et résultats locaux ; Studio Tool/Tool Call/Workflow/Automation Run ; Govern l’autorité future ; Shared les mécanismes génériques.

## 6. Utilisateurs
Principal : **Detection Test Reviewer**. Secondaires : Detection Engineer; Dataset Curator; Quality Reviewer.

## 7. Conditions d’entrée
Tenant, environnement, objective, scope, versions, sources, permissions, restrictions et return origin sont explicites. Une dépendance absente produit un état incomplet, partiel ou bloqué ; aucune donnée ou autorité n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | ---: | --- | --- |
| Detection Test Scenario | CAP-INV-412 | scenario, variants, dataset ref and restrictions | oui | version selected | outcome `draft` |
| Detection Hypothesis / Content logic | CAP-INV-403,406..408 | behavior and functional expectations | oui | linked versions | incomplete oracle |
| Ground-truth evidence projection | Evidence / curated labels / reviewer knowledge | support and uncertainty | non | source visible | confidence limited |
| Automated proposal context | Studio Automation Run | agent, sources and proposal | non | run resolvable | manual authoring |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Test Scenario / Dataset relation | Investigate concepts | target and variants | lecture |
| Detection Content Draft / Hypothesis | Investigate concepts | logic and success criteria | lecture |
| Evidence / labels | Investigate / source owner | support and limitations | lecture |
| Automation Run / Tool Calls | Studio | proposal provenance | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Expected Outcome | créer, modifier, contester, confirmer-by-review, superseder | Investigate concept | expected ≠ actual and ≠ certain ground truth |
| Oracle review disposition | enregistrer | Investigate concept | reviewer and confidence mandatory |
| Trace / Version event | émettre | Shared | proposal and human disposition retained |

## 11. Fonctionnalités
- define required/forbidden/optional matches
- define expected fields, relations, count and window
- define warnings, blocking errors and partial outcomes
- record source, confidence and contradictions
- review, dispute, version and supersede
- conserver versions, erreurs, partialité, restrictions et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | ---: | --- | --- | --- |
| Consulter/Comparer | Detection Test Reviewer | Expected Outcome | 0 | lecture autorisée | projection sourcée | non |
| Exécuter traitement borné | Detection Test Reviewer | Tool Call / Result | 1 | déclenchement explicite et permission | résultat attribué | policy |
| Créer/Modifier/Contester | Detection Test Reviewer | Expected Outcome | 2 | mutation réversible | nouvelle version | OPEN-013 |
| Préparer handoff | Detection Test Reviewer | candidate package | 2 | sources et limites visibles | package non effectif | destination |

Classes 3/4 exclues ; production et runtime appartiennent à 4B.3A.2/owners.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | ---: | ---: | ---: | ---: | --- |
| Construire Expected Outcome | oui | éditeur/contrôles explicables | oui | proposition | formulaire/table/revue |
| Valider ou comparer | oui | validateur/comparateur | oui | explication | diagnostics/diff |
| Expliquer erreurs | oui | catalogue | oui | résumé sourcé | erreurs brutes/checklist |
| Promouvoir/déployer/qualifier runtime | non | non | non | interdit | future phase/owner |

Toute sortie expose initiateur, producteur/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun choix silencieux.

## 14. États fonctionnels
`draft`, `incomplete`, `proposed`, `under-review`, `reviewed`, `disputed`, `low-confidence`, `superseded`, `withdrawn`. États fonctionnels, pas machine objet finale.

## 15. États d’interface
Loading conserve context/version ; Empty distingue absence et interdiction ; Partial expose gaps ; Error conserve le valide ; Offline bloque les nouveaux runs ; Permission denied masque ; Stale distingue ancien/courant ; conflits fournissent diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Expected Outcome | test oracle concept | CAP-INV-411,414,415 | source, confidence and reviewer visible |
| Oracle review disposition | review event | Project / Audit | AI proposal not auto-accepted |
| Outcome diff | comparison | Reviewer | versions and contradictions visible |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-412 | définir attendu | CAP-INV-413 | scenario, dataset source, target version and restrictions | tests |
| CAP-INV-413 | execute test/replay | CAP-INV-411/414 | reviewed expected outcome, confidence and version | tests |
| CAP-INV-413 | compare actual result | CAP-INV-415 | expected vs actual, warnings, gaps and provenance | tests |

Transitions conservent ownership, tenant/env, permissions, restrictions, versions, erreurs, provenance et return origin.

## 18. Dépendances
CAP-INV-403,406..408,412; Studio Automation provenance; Shared Versioning/Comparison; OPEN-013/015. Shared Jobs/Trace/Versioning/Linking/Search/Export/Reporting/Collaboration/Comparison/Recovery consommés sans redéfinition.

## 19. Source de vérité
Investigate est source de Expected Outcome ; tous les objets consommés restent chez leurs owners. Draft ≠ runtime Detection.

## 20. Provenance et audit
Source du besoin, Project, versions, sources/schemas/fields/mappings/logic, Tool/Calls/Runs, paramètres, datasets, résultats, erreurs, interruptions, auteurs, reviewers, dispositions, exports et correlation ID applicables à Expected Outcome.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | ---: | --- | --- | --- | --- | --- |
| Expected Outcome read/create/update | test authority | 0/2 | sensitive fields masked | OPEN-013 | author/reviewer | Investigate | Permissions |
| Expected Outcome review | ground-truth risk | 2 | confidence and sources mandatory | possible | reviewer distinct | Investigate | Permissions |
| Automated oracle proposal | hallucinated truth | 1/2 | agent/run/source visible | possible | human accept/modify/reject | Studio/Investigate | Permissions |

Matrice atomique, namespaces, RBAC/ABAC, step-up et SoD finaux reportés ; permissions production exclues.

## 22. Limites et erreurs
- Expected Outcome ≠ actual result.
- Expected Outcome ≠ certain ground truth.
- AI-generated output cannot become oracle automatically.
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
### 1. AI proposal
**Given** un Expected Outcome proposed by an agent  
**When** reviewer opens it  
**Then** agent/version/sources are visible, proposal is not accepted automatically and can be modified or rejected

### 2. Contradictory labels
**Given** two curated labels disagree  
**When** outcome is reviewed  
**Then** contradiction and confidence remain visible and status may be disputed

### 3. Sans IA
**Given** aucun modèle  
**When** outcome is authored  
**Then** manual forms, examples and reviewer workflow remain available

## 26. Questions ouvertes
- OPEN-013 reste ouverte.
- OPEN-015 reste ouverte.
- Le moteur/langage Detection est une lacune future non couverte ; OPEN-005 reste forensic-only.
- Schémas, formats, permissions et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering module, Event Search/Hunt/Case/Evidence/technical handoffs, Command boundaries, Settings/Endpoint, Studio, Govern future review, Shared, Objects/Permissions/Screens/Journeys/Technique/4B.3A.2/Validation.
