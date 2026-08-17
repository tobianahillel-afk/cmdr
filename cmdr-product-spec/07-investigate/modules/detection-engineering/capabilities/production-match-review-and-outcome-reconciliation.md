---
id: CAP-INV-428
title: Production Match Review and Outcome Reconciliation
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
  - OPEN-013
  - OPEN-015
  - OPEN-017
source-of-truth: canonical
---
# CAP-INV-428 — Production Match Review and Outcome Reconciliation

## 1. Définition
Examiner les matches runtime, Signals, événements sources, dispositions Command, Findings, Evidence et Incidents afin de réconcilier les observations préproduction et production avec incertitude.

## 2. Problème utilisateur
Une disposition opérationnelle peut être traitée comme ground truth absolue et transformer automatiquement un match en TP ou FP certain.

## 3. Objectifs
- voir runtime matches, Signals, source events and Command dispositions
- voir Findings, Evidence and confirmed/non-confirmed Incidents when available
- revoir candidate TP/FP and identify candidate FN from Incidents or Hunts
- classer inconclusive and document reviewer/justification
- compare versions and preproduction results; prepare new Draft or Gap

## 4. Non-objectifs
Aucune API, protocole, syntaxe vendor, moteur, langage, compilateur, parser, AST, format technique, commande, code, modèle ML, mutation runtime directe, capability CAP-INV-5xx, Threat Intelligence, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède **Production Match Review** comme concept fonctionnel. Command garde Detection/Signal/Alert/Incident ; Settings les runtimes, targets, environments et health ; Endpoint ses capacités et exécutions locales ; Govern Action Request/Decision/Approval/Response Run/Result ; Studio Tools et Automation Runs ; Shared les mécanismes transverses.

## 6. Utilisateurs
Principal : **Detection Reviewer**. Secondaires : Detection Engineer, Detection Reviewer, Detection Owner, Incident Commander, Platform Operator et Approver autorisés.

## 7. Conditions d’entrée
Tenant, environnement, target, version, période, owner, permissions, restrictions et return origin sont explicites. Les preuves 4B.3A.1 sont référencées sans duplication. Toute absence produit un état incomplete, partial, blocked ou unknown.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | --- | --- | --- |
| Runtime matches and Signals | Command | matches, Signal links, version and target | oui | time-bounded | unreviewed |
| Command dispositions / Incident outcomes | Command | triage, closure, impact and analyst comments | non | current linked objects | ground truth limited |
| Source events / Findings / Evidence | Shared / Investigate | technical and investigative context | oui when permitted | insufficient-context | insufficient-context |
| Preproduction evidence | CAP-INV-411..415 | tests, Expected Outcomes, replay and candidate dispositions | oui | selected versions | comparison unavailable |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Runtime Detection / Signal / Alert / Incident | Command | runtime chain and dispositions | lecture |
| Telemetry Event | Shared | source events and time context | lecture |
| Finding / Evidence / Hunt | Investigate | qualified context and FN candidates | lecture/lien |
| Validation/Replay/Match Review | Investigate | preproduction baseline | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Production Match Review | créer, classifier, contester, superseder | Investigate concept | candidate labels with uncertainty |
| Outcome reconciliation relation | lier/versionner | Investigate | preproduction vs runtime evidence preserved |
| New Draft / Gap context | préparer | CAP-INV-406/416 | never active mutation |

## 11. Fonctionnalités
- voir runtime matches, Signals, source events and Command dispositions
- voir Findings, Evidence and confirmed/non-confirmed Incidents when available
- revoir candidate TP/FP and identify candidate FN from Incidents or Hunts
- classer inconclusive and document reviewer/justification
- compare versions and preproduction results; prepare new Draft or Gap
- conserver sources, versions, erreurs, partialité, attribution et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | --- | --- | --- | --- |
| consulter/comparer | Detection Reviewer | Production Match Review | 0 | lecture autorisée | projection sourcée | non |
| assessment bornée | Detection Reviewer | Tool Call / assessment | 1 | lancement explicite | résultat attribué | policy |
| créer/modifier/contester | Detection Reviewer | Production Match Review | 2 | mutation réversible | nouvelle version | OPEN-013 |
| préparer demande | Detection Reviewer | Action Request context | 2 | risque, cible et rollback visibles | package non effectif | Govern |
| changement réel | Govern/runtime owner | runtime target | 3 | Decision/Approval | Result projeté | owner |
| détruire historique | aucun rôle local | provenance | 4 | interdit | refus | strict |

Investigate n’exécute jamais les classes 3/4.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | --- | --- | --- | --- | --- |
| construire/revoir Production Match Review | oui | formulaires et contrôles | oui | proposition | matrices et revue humaine |
| comparer/valider | oui | diff et règles explicables | oui | explication | diagnostics et checklist |
| résumer risques/erreurs | oui | agrégations | oui | résumé sourcé | tables et timeline |
| approuver/exécuter | non localement | non | non | interdit | Govern/runtime owner |

Initiateur, agent/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner et disposition humaine sont visibles. Aucun choix silencieux.

## 14. États fonctionnels
`unreviewed`, `candidate-true-positive`, `candidate-false-positive`, `candidate-false-negative`, `expected-non-match`, `ambiguous`, `insufficient-context`, `disputed`, `reconciled`, `superseded`. Projections fonctionnelles, pas machine d’état objet définitive.

## 15. États d’interface
Loading conserve le contexte ; Empty distingue absence et interdiction ; Partial détaille les targets ; Error conserve le valide ; Offline bloque les mutations ; Permission denied masque ; Stale distingue ancien/courant ; les conflits offrent diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Production Match Review | review disposition | CAP-INV-429/416/435 | candidate status, evidence and limits |
| New Draft proposal | change context | CAP-INV-406 | source active version and rationale retained |
| Detection Gap candidate | gap context | CAP-INV-416 | data/logic/coverage cause separated |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-427 | quality review requires sample | CAP-INV-428 | runtime matches, Command feedback, version and period | Quality |
| CAP-INV-428 | tuning indicated | CAP-INV-429 | candidate FP/FN, evidence, risks and source version | Production Review |
| CAP-INV-428 | gap indicated | CAP-INV-416 | missing behavior/data/context and uncertainty | Production Review |
| CAP-INV-428 | new investigation needed | CAP-INV-401/005 | source events, disposition limits and return origin | Production Review |

Chaque transition conserve ownership, tenant/env, target, version, permissions, restrictions, erreurs, autorité, provenance et return origin.

## 18. Dépendances
CAP-INV-401,005,406,411..416,427,429,435; Command runtime chain; Shared events; Case/Evidence. OPEN-017 couvre la future stratégie runtime/langage/portabilité sans option sélectionnée. Shared Jobs, Notifications, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Collaboration, Comparison, Inspector, Audit Hooks et Recovery sont consommés sans redéfinition.

## 19. Source de vérité
Investigate possède l’assessment/proposition locale. Les objets Command, Settings, Endpoint, Govern, Studio et Shared restent canoniques chez leurs owners. Une projection ne remplace jamais sa source.

## 20. Provenance et audit
Conserver besoin, Project, Hypothesis, Drafts/Versions, Review Package, Release Candidate, reviews, readiness, targets, plans, Action Requests, Decisions, Approvals, Runs, Results, runtime observations, health, Signals/Alerts/Incidents, assessments, propositions, erreurs, auteurs, timestamps, Tools/Runs, paramètres et dispositions. Aucune trace supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Production match review read/create/update | sensitive runtime context | 0/2 | raw data permission-bound | OPEN-013 | reviewer/content owner | Investigate | Permissions |
| Command disposition read | operational context | 0 | Command policies | possible | read-only consumer | Command | Permissions |
| Candidate TP/FP/FN classify | ground-truth risk | 2 | uncertainty mandatory | possible | independent reviewer | Investigate | Permissions |

Permissions atomiques, namespaces, RBAC/ABAC, step-up et séparation finale restent reportés.

## 22. Limites et erreurs
- Runtime match ≠ confirmed malicious activity.
- Command disposition ≠ absolute ground truth.
- Incident closure ≠ rule correctness.
- No active rule or Command object is modified.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, target offline et version superseded restent visibles.
- Result, match, health, disposition Command ou sortie IA ne vaut pas conclusion universelle.

## 23. Métriques conceptuelles
- reviews by candidate disposition/version
- inconclusive and insufficient-context rates
- preproduction/runtime disagreement
- certain auto-classifications — target zero
- provenance et dispositions humaines complètes
- exécutions silencieuses, auto-approbations et suppressions de trace : cible zéro

Aucun seuil technique universel n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; documentation only. Aucun `validated`, `implemented`, `deployed`, `active`, `native` ou `integrated`.

## 25. Critères d’acceptation
### 1. Disposition incertaine
**Given** un Signal closed sans Evidence suffisante  
**When** production review is performed  
**Then** the match may remain inconclusive, sources/limits remain visible and no certain FP is imposed

### 2. Candidate FN
**Given** an Incident or Hunt indicates behavior with no runtime match  
**When** it is reviewed  
**Then** data availability and coverage are checked before candidate-FN disposition

### 3. Sans IA
**Given** aucun modèle  
**When** production matches are reviewed  
**Then** source pivots, comparison, tables and human classification suffice

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-015 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-017 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-005 reste forensic-only.
- Schémas, formats, cardinalités, permissions atomiques, contrats techniques et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering lifecycle, Command, Settings, Endpoint, Govern, Studio, Shared, Objects, Permissions, Screens, Journeys, Technique, validation et futur handoff 4B.3B non canonique.
