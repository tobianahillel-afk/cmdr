---
id: CAP-INV-417
title: Detection Authoring Provenance and Review Handoff
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
  - OPEN-014
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-417 — Detection Authoring Provenance and Review Handoff
## 1. Définition
Retracer toutes les sources, versions, Tools, Runs, datasets, résultats et décisions humaines de l’authoring, puis préparer un package versionné destiné à la future revue et promotion sans créer Approval, Decision, deployment ou runtime Detection.

## 2. Problème utilisateur
Un draft transmis sans lineage complet peut masquer sources, gaps, résultats partiels, propositions automatisées ou éléments non résolus et être pris pour une règle approuvée.

## 3. Objectifs
- trace Case, Incident, Finding, Hypothesis, Hunt and technical handoffs
- trace project, draft versions, sources, schemas, fields, mappings, Tools and Runs
- trace enrichments, scenarios, datasets, outcomes, validations, replay, match reviews, coverage and gaps
- trace authors, reviewers, timestamps, errors, interruptions and dispositions
- prepare review package with version, objective, dependencies, results, limits, risks and unresolved items

## 4. Non-objectifs
- aucun moteur, langage, syntaxe vendor, API, protocole, parser, compilateur, AST, modèle ML, commande ou code
- aucune promotion, deployment, activation, deactivation, rollback, exception active ou mutation Signal/Alert
- aucune capability CAP-INV-5xx, Intelligence, Cloud/Mobile ou réécriture détaillée d’écran

## 5. Propriétaire
Investigate possède Detection Authoring Provenance and Review Package et sa disposition humaine. Command conserve runtime Detection/Signal/Alert/Incident ; Settings les sources/parsers/schemas/health/retention ; Endpoint Agent ses capacités et résultats locaux ; Studio Tool/Tool Call/Workflow/Automation Run ; Govern l’autorité future ; Shared les mécanismes génériques.

## 6. Utilisateurs
Principal : **Detection Engineering Lead**. Secondaires : Detection Reviewer; Audit Analyst; Future Promotion Reviewer.

## 7. Conditions d’entrée
Tenant, environnement, objective, scope, versions, sources, permissions, restrictions et return origin sont explicites. Une dépendance absente produit un état incomplet, partiel ou bloqué ; aucune donnée ou autorité n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | ---: | --- | --- |
| Detection Engineering Project and candidate version | CAP-INV-402/406/410 | owner, objective, draft and documentation | oui | selected immutable version | package `incomplete` |
| Validation, tests, replay and reviews | CAP-INV-411..415 | results, errors, partiality and dispositions | oui | selected versions | package with missing evidence |
| Coverage and gaps | CAP-INV-416 | assessment, limitations and unresolved items | oui | current selected assessment | review not ready |
| Source and automation provenance | Cases/Hunts/analyses/Studio/Shared | sources, Tools, Calls, Runs and timestamps | oui | resolvable chain | `not-reproducible` or blocked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Case / Incident / Finding / Hypothesis / Hunt / technical packages | owners | origin and supporting context | lecture/lien |
| Project / Detection Content / test and result concepts | Investigate | authoring chain | lecture |
| Tool / Tool Call / Automation Run | Studio | producer, version, parameters and outputs | lecture |
| Trace / Activity / Version | Shared | events and lineage | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Provenance Record relation | assembler/annoter | Investigate using Shared records | no Trace redefinition |
| Review Package | créer, versionner, submit-to-future, withdraw, superseder | Investigate concept | not Approval/Decision/deployment |
| Handoff disposition | enregistrer | Investigate / Shared | destination response and return origin retained |

## 11. Fonctionnalités
- assemble complete source and authoring lineage
- expose versions, parameters, errors, interruptions and human decisions
- assess missing or disputed provenance
- prepare candidate review package
- handoff only to future 4B.3A.2 and preserve return
- conserver versions, erreurs, partialité, restrictions et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | ---: | --- | --- | --- |
| Consulter/Comparer | Detection Engineering Lead | Detection Authoring Provenance and Review Package | 0 | lecture autorisée | projection sourcée | non |
| Exécuter traitement borné | Detection Engineering Lead | Tool Call / Result | 1 | déclenchement explicite et permission | résultat attribué | policy |
| Créer/Modifier/Contester | Detection Engineering Lead | Detection Authoring Provenance and Review Package | 2 | mutation réversible | nouvelle version | OPEN-013 |
| Préparer handoff | Detection Engineering Lead | candidate package | 2 | sources et limites visibles | package non effectif | destination |

Classes 3/4 exclues ; production et runtime appartiennent à 4B.3A.2/owners.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | ---: | ---: | ---: | ---: | --- |
| Construire Detection Authoring Provenance and Review Package | oui | éditeur/contrôles explicables | oui | proposition | formulaire/table/revue |
| Valider ou comparer | oui | validateur/comparateur | oui | explication | diagnostics/diff |
| Expliquer erreurs | oui | catalogue | oui | résumé sourcé | erreurs brutes/checklist |
| Promouvoir/déployer/qualifier runtime | non | non | non | interdit | future phase/owner |

Toute sortie expose initiateur, producteur/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun choix silencieux.

## 14. États fonctionnels
`draft`, `incomplete`, `ready-for-review`, `submitted-to-future-review`, `returned`, `withdrawn`, `superseded`, `blocked`, `provenance-disputed`. États fonctionnels, pas machine objet finale.

## 15. États d’interface
Loading conserve context/version ; Empty distingue absence et interdiction ; Partial expose gaps ; Error conserve le valide ; Offline bloque les nouveaux runs ; Permission denied masque ; Stale distingue ancien/courant ; conflits fournissent diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Detection Review Package | review package | Future Phase 4B.3A.2 | not Approval, deployment or activation |
| Authoring provenance assessment | assessment | Audit / Project | sources, gaps and reproducibility visible |
| Handoff disposition | event | Project / Trace | submitted, returned, withdrawn or superseded |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-402/406/410..416 | prepare review package | CAP-INV-417 | candidate version, sources, validation, tests, replay, FP/FN, coverage, risks | project |
| CAP-INV-417 | submit for future review | Future Phase 4B.3A.2 | review package, unresolved items, provenance and return origin | project |
| Future owner | return or reject | CAP-INV-417/406 | disposition, requested changes and source version | future review |

Transitions conservent ownership, tenant/env, permissions, restrictions, versions, erreurs, provenance et return origin.

## 18. Dépendances
CAP-INV-401..416; Shared Trace/Activity/Versioning/Export; Studio Runs; future 4B.3A.2; OPEN-013/014/015. Shared Jobs/Trace/Versioning/Linking/Search/Export/Reporting/Collaboration/Comparison/Recovery consommés sans redéfinition.

## 19. Source de vérité
Investigate est source de Detection Authoring Provenance and Review Package ; tous les objets consommés restent chez leurs owners. Draft ≠ runtime Detection.

## 20. Provenance et audit
Source du besoin, Project, versions, sources/schemas/fields/mappings/logic, Tool/Calls/Runs, paramètres, datasets, résultats, erreurs, interruptions, auteurs, reviewers, dispositions, exports et correlation ID applicables à Detection Authoring Provenance and Review Package.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | ---: | --- | --- | --- | --- | --- |
| Review package prepare | future production impact | 2 | sensitive datasets/references masked | OPEN-013 | author/reviewer | Investigate | Permissions |
| Provenance export | data diffusion | 1 | redaction and export policy | step-up possible | exporter/reviewer | Shared/Security | Permissions |
| Future review submit | governance boundary | 2 | no approval implied | future SoD | Investigate/future owner | 4B.3A.2 |

Matrice atomique, namespaces, RBAC/ABAC, step-up et SoD finaux reportés ; permissions production exclues.

## 22. Limites et erreurs
- Review Package ≠ Approval or Decision.
- Approval ≠ deployment; version created ≠ version active.
- No promotion, deployment, activation, deactivation, rollback, production exception, Signal or Alert mutation.
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
### 1. Future handoff
**Given** a validated draft with tests, replay, FP/FN review and documented coverage  
**When** review package is prepared  
**Then** all items are versioned, limitations/unresolved items visible, no Approval or deployment occurs, destination is future 4B.3A.2

### 2. Missing provenance
**Given** a Tool Call version is unavailable  
**When** package is prepared  
**Then** gap is visible and status remains incomplete/blocked rather than reproducible

### 3. Sans IA
**Given** aucun model  
**When** package is assembled  
**Then** deterministic relations, checklists, versioning and human review cover the workflow

## 26. Questions ouvertes
- OPEN-013 reste ouverte.
- OPEN-014 reste ouverte.
- OPEN-015 reste ouverte.
- Le moteur/langage Detection est une lacune future non couverte ; OPEN-005 reste forensic-only.
- Schémas, formats, permissions et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering module, Event Search/Hunt/Case/Evidence/technical handoffs, Command boundaries, Settings/Endpoint, Studio, Govern future review, Shared, Objects/Permissions/Screens/Journeys/Technique/4B.3A.2/Validation.
