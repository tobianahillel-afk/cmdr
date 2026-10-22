---
id: CAP-INV-432
title: Detection Performance and Resource Impact Assessment
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
# CAP-INV-432 — Detection Performance and Resource Impact Assessment

## 1. Définition
Comparer fonctionnellement volume, fréquence, latency, processing delay, throughput, errors, timeouts, declared resources, cost, throttling, dropped evaluations, skipped windows et impacts sans imposer de seuil universel.

## 2. Problème utilisateur
Une hausse de coût ou de latency peut être confondue avec une baisse d’accuracy, tandis qu’une performance rapide peut être présentée comme qualité élevée.

## 3. Objectifs
- voir and compare volume, frequency, latency, delay, throughput, errors and timeouts
- voir declared memory/resource/cost when available
- voir throttling, dropped evaluations and skipped windows
- compare targets, versions, periods, sources and relative load
- document operational impact, limitations and observed benefit

## 4. Non-objectifs
Aucune API, protocole, syntaxe vendor, moteur, langage, compilateur, parser, AST, format technique, commande, code, modèle ML, mutation runtime directe, capability CAP-INV-5xx, Threat Intelligence, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède **Performance Assessment** comme concept fonctionnel. Command garde Detection/Signal/Alert/Incident ; Settings les runtimes, targets, environments et health ; Endpoint ses capacités et exécutions locales ; Govern Action Request/Decision/Approval/Response Run/Result ; Studio Tools et Automation Runs ; Shared les mécanismes transverses.

## 6. Utilisateurs
Principal : **Platform Operator**. Secondaires : Detection Engineer, Detection Reviewer, Detection Owner, Incident Commander, Platform Operator et Approver autorisés.

## 7. Conditions d’entrée
Tenant, environnement, target, version, période, owner, permissions, restrictions et return origin sont explicites. Les preuves 4B.3A.1 sont référencées sans duplication. Toute absence produit un état incomplete, partial, blocked ou unknown.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | --- | --- | --- |
| Runtime/version/target context | CAP-INV-425/426 | version, state, health and target scope | oui | current/time-bounded | not-assessed |
| Metrics/resource projections | Settings / Endpoint / Shared Metrics | volume, latency, throughput, errors, resources and cost | oui where available | same period | insufficient-data |
| Quality/benefit context | CAP-INV-427/428 | Signals, operational value and match review | oui | same version/period | cost without value context |
| Source/drift context | CAP-INV-431 / Settings | data volume and dependency changes | non | linked assessment | cause uncertain |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Runtime Detection / Version | Command | runtime scope and results | lecture |
| Platform/Endpoint metrics | Settings / Endpoint | resource and execution projections | lecture |
| Metrics Engine / Data Quality | Shared | definitions, freshness and completeness | lecture |
| Runtime Quality / Production Review | Investigate | operational benefit and limitations | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Performance Assessment | créer, comparer, contester, superseder | Investigate concept | performance ≠ quality |
| Resource/capacity request context | préparer | Settings/Govern | no resource mutation |
| Tuning Proposal context | préparer | CAP-INV-429 | no active tuning |

## 11. Fonctionnalités
- voir and compare volume, frequency, latency, delay, throughput, errors and timeouts
- voir declared memory/resource/cost when available
- voir throttling, dropped evaluations and skipped windows
- compare targets, versions, periods, sources and relative load
- document operational impact, limitations and observed benefit
- conserver sources, versions, erreurs, partialité, attribution et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | --- | --- | --- | --- |
| consulter/comparer | Platform Operator | Performance Assessment | 0 | lecture autorisée | projection sourcée | non |
| assessment bornée | Platform Operator | Tool Call / assessment | 1 | lancement explicite | résultat attribué | policy |
| créer/modifier/contester | Platform Operator | Performance Assessment | 2 | mutation réversible | nouvelle version | OPEN-013 |
| préparer demande | Platform Operator | Action Request context | 2 | risque, cible et rollback visibles | package non effectif | Govern |
| changement réel | Govern/runtime owner | runtime target | 3 | Decision/Approval | Result projeté | owner |
| détruire historique | aucun rôle local | provenance | 4 | interdit | refus | strict |

Investigate n’exécute jamais les classes 3/4.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | --- | --- | --- | --- | --- |
| construire/revoir Performance Assessment | oui | formulaires et contrôles | oui | proposition | matrices et revue humaine |
| comparer/valider | oui | diff et règles explicables | oui | explication | diagnostics et checklist |
| résumer risques/erreurs | oui | agrégations | oui | résumé sourcé | tables et timeline |
| approuver/exécuter | non localement | non | non | interdit | Govern/runtime owner |

Initiateur, agent/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner et disposition humaine sont visibles. Aucun choix silencieux.

## 14. États fonctionnels
`not-assessed`, `assessing`, `acceptable`, `degraded`, `resource-constrained`, `throttled`, `partial`, `insufficient-data`, `candidate-regression`, `candidate-improvement`, `disputed`, `superseded`. Projections fonctionnelles, pas machine d’état objet définitive.

## 15. États d’interface
Loading conserve le contexte ; Empty distingue absence et interdiction ; Partial détaille les targets ; Error conserve le valide ; Offline bloque les mutations ; Permission denied masque ; Stale distingue ancien/courant ; les conflits offrent diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Performance Assessment | assessment | CAP-INV-429/433/435 | metrics, definitions, freshness and limitations |
| Capacity request | request context | Platform Settings / Govern | targets and impact visible |
| Tuning candidate | proposal context | CAP-INV-429 | performance and quality evidence separated |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-425/426 | runtime performance review | CAP-INV-432 | version, targets, health, period and metric definitions | Runtime |
| CAP-INV-432 | tuning indicated | CAP-INV-429 | performance issue, benefit context, risks and required tests | Performance |
| CAP-INV-432 | capacity/remediation indicated | Settings/Govern request | targets, resource impact and operational risk | Performance |
| CAP-INV-432 | rollback risk indicated | CAP-INV-433 | regression, stop criteria and affected targets | Performance |

Chaque transition conserve ownership, tenant/env, target, version, permissions, restrictions, erreurs, autorité, provenance et return origin.

## 18. Dépendances
CAP-INV-425..431,433,435; Settings/Endpoint metrics; Shared Metrics/Data Quality; Command. OPEN-017 couvre la future stratégie runtime/langage/portabilité sans option sélectionnée. Shared Jobs, Notifications, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Collaboration, Comparison, Inspector, Audit Hooks et Recovery sont consommés sans redéfinition.

## 19. Source de vérité
Investigate possède l’assessment/proposition locale. Les objets Command, Settings, Endpoint, Govern, Studio et Shared restent canoniques chez leurs owners. Une projection ne remplace jamais sa source.

## 20. Provenance et audit
Conserver besoin, Project, Hypothesis, Drafts/Versions, Review Package, Release Candidate, reviews, readiness, targets, plans, Action Requests, Decisions, Approvals, Runs, Results, runtime observations, health, Signals/Alerts/Incidents, assessments, propositions, erreurs, auteurs, timestamps, Tools/Runs, paramètres et dispositions. Aucune trace supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Performance assessment read/create/update | production metrics and cost | 0/2 | sensitive capacity data scoped | OPEN-013 | assessor/platform owner | Investigate | Permissions |
| Resource metrics read | infrastructure exposure | 0 | aggregate/mask | possible | viewer/admin | Settings/Endpoint | Permissions |
| Capacity/tuning/rollback request prepare | operational impact | 2 | evidence visible | step-up possible | requester/approver future | Settings/Govern | Permissions |

Permissions atomiques, namespaces, RBAC/ABAC, step-up et séparation finale restent reportés.

## 22. Limites et erreurs
- Latency ≠ detection accuracy.
- High performance ≠ high quality; high cost ≠ low value automatically.
- No universal technical threshold is imposed.
- No automatic disable, capacity change or active tuning.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, target offline et version superseded restent visibles.
- Result, match, health, disposition Command ou sortie IA ne vaut pas conclusion universelle.

## 23. Métriques conceptuelles
- metrics completeness/freshness
- performance by version/target/period
- candidate regressions/improvements
- automatic deactivations — target zero
- provenance et dispositions humaines complètes
- exécutions silencieuses, auto-approbations et suppressions de trace : cible zéro

Aucun seuil technique universel n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; documentation only. Aucun `validated`, `implemented`, `deployed`, `active`, `native` ou `integrated`.

## 25. Critères d’acceptation
### 1. Performance dégradée
**Given** latency rises while runtime remains active and no accuracy loss is proven  
**When** performance is assessed  
**Then** latency and accuracy remain distinct, health may be degraded and tuning/capacity requests may be prepared

### 2. Cost high, value high
**Given** resource cost increases with observed operational benefit  
**When** assessment completes  
**Then** cost and value remain separately visible and no automatic low-value disposition is made

### 3. Sans IA
**Given** aucun modèle  
**When** performance is assessed  
**Then** deterministic metrics, comparisons and human analysis suffice

## 26. Questions ouvertes
- OPEN-008 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-015 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-017 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-005 reste forensic-only.
- Schémas, formats, cardinalités, permissions atomiques, contrats techniques et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering lifecycle, Command, Settings, Endpoint, Govern, Studio, Shared, Objects, Permissions, Screens, Journeys, Technique, validation et futur handoff 4B.3B non canonique.
