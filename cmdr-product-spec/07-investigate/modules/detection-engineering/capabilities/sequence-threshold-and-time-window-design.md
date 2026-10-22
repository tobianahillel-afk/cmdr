---
id: CAP-INV-408
title: Sequence, Threshold and Time-Window Design
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
open_decisions:
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-408 — Sequence, Threshold and Time-Window Design
## 1. Définition
Définir fonctionnellement séquences, ordre, événements optionnels ou interdits, seuils, compteurs, groupes, fenêtres, retard, données hors ordre, tolérance aux manques et limites.

## 2. Problème utilisateur
Un ordre observé, un seuil dépassé ou une absence d’événement peut être interprété à tort comme chemin d’attaque, malveillance ou preuve négative.

## 3. Objectifs
- définir événements attendus, ordre, options et interdictions
- définir seuils, compteurs et groupes
- définir fenêtres, late data et out-of-order handling fonctionnels
- définir tolérance aux données manquantes et limites
- comparer configurations, versionner et préparer tests

## 4. Non-objectifs
- aucun moteur, langage, syntaxe vendor, API, protocole, parser, compilateur, AST, modèle ML, commande ou code
- aucune promotion, deployment, activation, deactivation, rollback, exception active ou mutation Signal/Alert
- aucune capability CAP-INV-5xx, Intelligence, Cloud/Mobile ou réécriture détaillée d’écran

## 5. Propriétaire
Investigate possède Sequence and Threshold Definition et sa disposition humaine. Command conserve runtime Detection/Signal/Alert/Incident ; Settings les sources/parsers/schemas/health/retention ; Endpoint Agent ses capacités et résultats locaux ; Studio Tool/Tool Call/Workflow/Automation Run ; Govern l’autorité future ; Shared les mécanismes génériques.

## 6. Utilisateurs
Principal : **Detection Engineer**. Secondaires : Detection Reviewer; Threat Hunter.

## 7. Conditions d’entrée
Tenant, environnement, objective, scope, versions, sources, permissions, restrictions et return origin sont explicites. Une dépendance absente produit un état incomplet, partiel ou bloqué ; aucune donnée ou autorité n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | ---: | --- | --- |
| Condition and correlation design | CAP-INV-407 | events, keys, exclusions and relations | oui | version liée | definition `incomplete` |
| Telemetry readiness and time quality | CAP-INV-404 | coverage, gaps, freshness et time semantics | oui | assessment courant | warnings/blocked |
| Field mapping review | CAP-INV-405 | timestamp, grouping et identity fields | oui | version liée | `field-missing` |
| Detection Hypothesis | CAP-INV-403 | behavior and success criteria | oui | version liée | aucune configuration utile |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Detection Content Draft / Conditions | Investigate concepts | version et logic base | lecture |
| Telemetry Event time projections | Shared | event/ingestion time et quality | lecture |
| Readiness Assessment | Investigate concept | gaps et coverage | lecture |
| Field Reference | Investigate/Settings | grouping and timestamp semantics | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Sequence Definition | créer, modifier, comparer, superseder | Investigate concept | matched sequence ≠ confirmed path |
| Threshold Definition | créer, modifier, comparer, superseder | Investigate concept | threshold ≠ maliciousness |
| Time-window specification | versionner | Investigate concept | fonctionnelle, pas moteur |

## 11. Fonctionnalités
- définir sequence, optional/forbidden events et order
- définir thresholds, counters and groups
- définir time window, delay and out-of-order policies functionally
- définir missing-data tolerance
- comparer variants et préparer scenarios
- conserver versions, erreurs, partialité, restrictions et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | ---: | --- | --- | --- |
| Consulter/Comparer | Detection Engineer | Sequence and Threshold Definition | 0 | lecture autorisée | projection sourcée | non |
| Exécuter traitement borné | Detection Engineer | Tool Call / Result | 1 | déclenchement explicite et permission | résultat attribué | policy |
| Créer/Modifier/Contester | Detection Engineer | Sequence and Threshold Definition | 2 | mutation réversible | nouvelle version | OPEN-013 |
| Préparer handoff | Detection Engineer | candidate package | 2 | sources et limites visibles | package non effectif | destination |

Classes 3/4 exclues ; production et runtime appartiennent à 4B.3A.2/owners.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | ---: | ---: | ---: | ---: | --- |
| Construire Sequence and Threshold Definition | oui | éditeur/contrôles explicables | oui | proposition | formulaire/table/revue |
| Valider ou comparer | oui | validateur/comparateur | oui | explication | diagnostics/diff |
| Expliquer erreurs | oui | catalogue | oui | résumé sourcé | erreurs brutes/checklist |
| Promouvoir/déployer/qualifier runtime | non | non | non | interdit | future phase/owner |

Toute sortie expose initiateur, producteur/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun choix silencieux.

## 14. États fonctionnels
`draft`, `incomplete`, `reviewable`, `coverage-limited`, `time-semantics-uncertain`, `valid-with-warnings`, `blocked`, `superseded`, `withdrawn`. États fonctionnels, pas machine objet finale.

## 15. États d’interface
Loading conserve context/version ; Empty distingue absence et interdiction ; Partial expose gaps ; Error conserve le valide ; Offline bloque les nouveaux runs ; Permission denied masque ; Stale distingue ancien/courant ; conflits fournissent diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Sequence Definition | functional definition | CAP-INV-411..414 | order and gaps explicit |
| Threshold Definition | functional definition | CAP-INV-411..415 | counting scope visible |
| Time-window context | functional parameter set | Tests / Replay | time semantics and limits visible |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-407 | ajouter sequence/seuil/window | CAP-INV-408 | conditions, fields, keys and readiness | logic |
| CAP-INV-408 | préparer scenario | CAP-INV-412/413 | events, order, thresholds, window and expected behavior | logic |
| CAP-INV-408 | valider/replay | CAP-INV-411/414 | version, parameters, gaps and permissions | logic |

Transitions conservent ownership, tenant/env, permissions, restrictions, versions, erreurs, provenance et return origin.

## 18. Dépendances
CAP-INV-403..407; Shared time/normalization; OPEN-013/015. Shared Jobs/Trace/Versioning/Linking/Search/Export/Reporting/Collaboration/Comparison/Recovery consommés sans redéfinition.

## 19. Source de vérité
Investigate est source de Sequence and Threshold Definition ; tous les objets consommés restent chez leurs owners. Draft ≠ runtime Detection.

## 20. Provenance et audit
Source du besoin, Project, versions, sources/schemas/fields/mappings/logic, Tool/Calls/Runs, paramètres, datasets, résultats, erreurs, interruptions, auteurs, reviewers, dispositions, exports et correlation ID applicables à Sequence and Threshold Definition.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | ---: | --- | --- | --- | --- | --- |
| Sequence/threshold/window read | logique sensible | 0 | scope masqué si requis | possible | viewer/reviewer | Investigate | Permissions |
| Sequence/threshold/window update | impact logique | 2 | aucune valeur sensible copiée | OPEN-013 | author/reviewer | Investigate | Permissions |
| Automated parameter proposal | sur-ajustement | 1/2 | sources et uncertainty visibles | possible | human review | Studio/Investigate | Permissions |

Matrice atomique, namespaces, RBAC/ABAC, step-up et SoD finaux reportés ; permissions production exclues.

## 22. Limites et erreurs
- Sequence matched ≠ attack path confirmed.
- Threshold exceeded ≠ maliciousness; absence of event ≠ proof it did not occur.
- Time window ≠ universal truth; no runtime scheduler or algorithm defined.
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
### 1. Données hors ordre
**Given** une sequence attendue et des événements hors ordre  
**When** la configuration est revue  
**Then** late/out-of-order assumptions et uncertainty sont visibles

### 2. Gap
**Given** une source manque pendant la fenêtre  
**When** un non-match est observé  
**Then** aucune absence d’attaque ou FN certain n’est affirmée

### 3. Sans IA
**Given** aucun modèle  
**When** la sequence et le seuil sont définis  
**Then** éditeur structuré, tableaux et tests manuels restent disponibles

## 26. Questions ouvertes
- OPEN-013 reste ouverte.
- OPEN-015 reste ouverte.
- Le moteur/langage Detection est une lacune future non couverte ; OPEN-005 reste forensic-only.
- Schémas, formats, permissions et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering module, Event Search/Hunt/Case/Evidence/technical handoffs, Command boundaries, Settings/Endpoint, Studio, Govern future review, Shared, Objects/Permissions/Screens/Journeys/Technique/4B.3A.2/Validation.
