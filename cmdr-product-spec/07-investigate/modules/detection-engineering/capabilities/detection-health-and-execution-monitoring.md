---
id: CAP-INV-426
title: Detection Health and Execution Monitoring
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
  - REQ-PROD-055
  - REQ-AI-002
  - REQ-SEC-001
  - REQ-SEC-002
  - REQ-UX-010
open_decisions:
  - OPEN-008
  - OPEN-013
  - OPEN-015
  - OPEN-017
source-of-truth: canonical
---
# CAP-INV-426 — Detection Health and Execution Monitoring

## 1. Définition
Construire une projection fonctionnelle de santé d’un runtime Detection par version et target à partir de l’exécution, des sources, schémas, dépendances, erreurs, délais et périodes manquées, sans confondre santé et efficacité.

## 2. Problème utilisateur
Une règle active peut être techniquement healthy mais inutile ou noisy, tandis qu’un warning peut être présenté à tort comme outage ou perte de détection.

## 3. Objectifs
- voir runtime, version, targets, last execution and status
- voir source availability, schema compatibility and dependency health
- voir errors, warnings, latency, processing delay, partial results and skipped periods
- voir disabled/suspended/throttling states and health history
- relier incidents techniques, owner and support path

## 4. Non-objectifs
Aucune API, protocole, syntaxe vendor, moteur, langage, compilateur, parser, AST, format technique, commande, code, modèle ML, mutation runtime directe, capability CAP-INV-5xx, Threat Intelligence, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède **Detection Health Assessment** comme concept fonctionnel. Command garde Detection/Signal/Alert/Incident ; Settings les runtimes, targets, environments et health ; Endpoint ses capacités et exécutions locales ; Govern Action Request/Decision/Approval/Response Run/Result ; Studio Tools et Automation Runs ; Shared les mécanismes transverses.

## 6. Utilisateurs
Principal : **Platform Operator**. Secondaires : Detection Engineer, Detection Reviewer, Detection Owner, Incident Commander, Platform Operator et Approver autorisés.

## 7. Conditions d’entrée
Tenant, environnement, target, version, période, owner, permissions, restrictions et return origin sont explicites. Les preuves 4B.3A.1 sont référencées sans duplication. Toute absence produit un état incomplete, partial, blocked ou unknown.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | --- | --- | --- |
| Reconciled runtime state | CAP-INV-425 | version, target and observed state | oui | current reconciliation | health unknown |
| Execution/health projections | Settings / Endpoint / runtime owner | last execution, status, errors, delay and throttling | oui | current/time series | unavailable |
| Source/schema/dependency health | Settings / Shared Data Quality | availability, compatibility, gaps and failures | oui | current snapshot | partial assessment |
| Technical incident/support context | Command / Settings | linked incidents, owner and support path | non | current refs | no incident link |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Runtime Detection | Command | runtime identity/version/state | lecture |
| Platform/target health | Settings / Endpoint | execution and dependency health | lecture |
| Data Source / Parser / Schema quality | Settings / Shared | availability and compatibility | lecture |
| Incident / Result | Command / Govern | technical effect and reported outcome | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Detection Health Assessment | créer, comparer, contester, superseder | Investigate concept | health ≠ effectiveness |
| Health warning link | annoter/link | Investigate / Shared | source and timestamp retained |
| Support/remediation request context | préparer | Settings/Govern | no runtime mutation |

## 11. Fonctionnalités
- voir runtime, version, targets, last execution and status
- voir source availability, schema compatibility and dependency health
- voir errors, warnings, latency, processing delay, partial results and skipped periods
- voir disabled/suspended/throttling states and health history
- relier incidents techniques, owner and support path
- conserver sources, versions, erreurs, partialité, attribution et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | --- | --- | --- | --- |
| consulter/comparer | Platform Operator | Detection Health Assessment | 0 | lecture autorisée | projection sourcée | non |
| assessment bornée | Platform Operator | Tool Call / assessment | 1 | lancement explicite | résultat attribué | policy |
| créer/modifier/contester | Platform Operator | Detection Health Assessment | 2 | mutation réversible | nouvelle version | OPEN-013 |
| préparer demande | Platform Operator | Action Request context | 2 | risque, cible et rollback visibles | package non effectif | Govern |
| changement réel | Govern/runtime owner | runtime target | 3 | Decision/Approval | Result projeté | owner |
| détruire historique | aucun rôle local | provenance | 4 | interdit | refus | strict |

Investigate n’exécute jamais les classes 3/4.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | --- | --- | --- | --- | --- |
| construire/revoir Detection Health Assessment | oui | formulaires et contrôles | oui | proposition | matrices et revue humaine |
| comparer/valider | oui | diff et règles explicables | oui | explication | diagnostics et checklist |
| résumer risques/erreurs | oui | agrégations | oui | résumé sourcé | tables et timeline |
| approuver/exécuter | non localement | non | non | interdit | Govern/runtime owner |

Initiateur, agent/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner et disposition humaine sont visibles. Aucun choix silencieux.

## 14. États fonctionnels
`unknown`, `healthy`, `degraded`, `delayed`, `partially-executing`, `failing`, `source-unavailable`, `schema-incompatible`, `dependency-failed`, `disabled`, `suspended`, `stale`, `disputed`. Projections fonctionnelles, pas machine d’état objet définitive.

## 15. États d’interface
Loading conserve le contexte ; Empty distingue absence et interdiction ; Partial détaille les targets ; Error conserve le valide ; Offline bloque les mutations ; Permission denied masque ; Stale distingue ancien/courant ; les conflits offrent diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Detection Health Assessment | assessment | CAP-INV-427/431/432/433 | health dimensions and limits visible |
| Health incident context | link/request | Settings / Command | technical impact and source retained |
| Health history comparison | comparison result | Lifecycle / Audit | versions and periods attributed |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-425 | runtime reconciled | CAP-INV-426 | version, target, state and observation freshness | Reconciliation |
| Settings/Endpoint/runtime owner | health update | CAP-INV-426 | execution, errors, dependencies, delay and timestamps | source owner |
| CAP-INV-426 | degraded/failing | CAP-INV-431/432/433 / support request | health facts, targets, impact and limits | Health |
| CAP-INV-426 | assessment available | CAP-INV-427 | runtime health context and period | Health |

Chaque transition conserve ownership, tenant/env, target, version, permissions, restrictions, erreurs, autorité, provenance et return origin.

## 18. Dépendances
CAP-INV-425,427,431..433; Command Detection/Incident; Settings Health/Sources; Endpoint Health; Shared Metrics/Data Quality. OPEN-017 couvre la future stratégie runtime/langage/portabilité sans option sélectionnée. Shared Jobs, Notifications, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Collaboration, Comparison, Inspector, Audit Hooks et Recovery sont consommés sans redéfinition.

## 19. Source de vérité
Investigate possède l’assessment/proposition locale. Les objets Command, Settings, Endpoint, Govern, Studio et Shared restent canoniques chez leurs owners. Une projection ne remplace jamais sa source.

## 20. Provenance et audit
Conserver besoin, Project, Hypothesis, Drafts/Versions, Review Package, Release Candidate, reviews, readiness, targets, plans, Action Requests, Decisions, Approvals, Runs, Results, runtime observations, health, Signals/Alerts/Incidents, assessments, propositions, erreurs, auteurs, timestamps, Tools/Runs, paramètres et dispositions. Aucune trace supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Detection health read | production telemetry | 0 | target details scoped | possible | viewer/operator | Settings/Command | Permissions |
| Health assessment create/update | quality interpretation | 1/2 | raw errors masked as required | OPEN-013 | assessor/owner | Investigate | Permissions |
| Support/remediation request | operational impact | 2 | no direct action | step-up possible | requester/operator | Settings/Govern | Permissions |

Permissions atomiques, namespaces, RBAC/ABAC, step-up et séparation finale restent reportés.

## 22. Limites et erreurs
- Active ≠ healthy; healthy ≠ useful or effective.
- Health warning ≠ outage.
- Latency ≠ accuracy.
- No automatic disable, rollback or tuning.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, target offline et version superseded restent visibles.
- Result, match, health, disposition Command ou sortie IA ne vaut pas conclusion universelle.

## 23. Métriques conceptuelles
- health states by version/target
- execution gaps and delay periods
- source/schema/dependency failure categories
- automatic deactivations — target zero
- provenance et dispositions humaines complètes
- exécutions silencieuses, auto-approbations et suppressions de trace : cible zéro

Aucun seuil technique universel n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; documentation only. Aucun `validated`, `implemented`, `deployed`, `active`, `native` ou `integrated`.

## 25. Critères d’acceptation
### 1. Healthy mais noisy
**Given** un runtime technically healthy et des feedbacks de bruit  
**When** la qualité est évaluée  
**Then** health remains healthy while quality may be noisy and a tuning proposal can be prepared

### 2. Performance dégradée
**Given** latency increases without evidence of accuracy loss  
**When** health is assessed  
**Then** health may be degraded, latency and accuracy remain distinct and no rule is disabled automatically

### 3. Sans IA
**Given** aucun modèle  
**When** health is monitored  
**Then** deterministic projections, tables, timelines and human review suffice

## 26. Questions ouvertes
- OPEN-008 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-015 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-017 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-005 reste forensic-only.
- Schémas, formats, cardinalités, permissions atomiques, contrats techniques et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering lifecycle, Command, Settings, Endpoint, Govern, Studio, Shared, Objects, Permissions, Screens, Journeys, Technique, validation et futur handoff 4B.3B non canonique.
