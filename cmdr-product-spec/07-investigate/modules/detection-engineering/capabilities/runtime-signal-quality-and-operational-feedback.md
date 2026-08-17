---
id: CAP-INV-427
title: Runtime Signal Quality and Operational Feedback
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
# CAP-INV-427 — Runtime Signal Quality and Operational Feedback

## 1. Définition
Consommer les projections Command de Signals, Alerts, Incidents et dispositions pour évaluer la qualité opérationnelle par version et environnement sans reprendre l’ownership des objets runtime.

## 2. Problème utilisateur
Le volume de Signals ou la clôture d’un Incident peut être utilisé comme mesure absolue de qualité alors que les dispositions et la ground truth restent incomplètes.

## 3. Objectifs
- consommer Signals, Alerts, Incidents, dispositions, comments, escalations and closures
- voir priorities, handling time, reported noise, duplicates, missing context and operational value
- regrouper, filtrer and compare by environment/version
- documenter ground-truth limitations
- préparer Runtime Quality Assessment and tuning proposal

## 4. Non-objectifs
Aucune API, protocole, syntaxe vendor, moteur, langage, compilateur, parser, AST, format technique, commande, code, modèle ML, mutation runtime directe, capability CAP-INV-5xx, Threat Intelligence, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède **Runtime Quality Assessment** comme concept fonctionnel. Command garde Detection/Signal/Alert/Incident ; Settings les runtimes, targets, environments et health ; Endpoint ses capacités et exécutions locales ; Govern Action Request/Decision/Approval/Response Run/Result ; Studio Tools et Automation Runs ; Shared les mécanismes transverses.

## 6. Utilisateurs
Principal : **Detection Reviewer**. Secondaires : Detection Engineer, Detection Reviewer, Detection Owner, Incident Commander, Platform Operator et Approver autorisés.

## 7. Conditions d’entrée
Tenant, environnement, target, version, période, owner, permissions, restrictions et return origin sont explicites. Les preuves 4B.3A.1 sont référencées sans duplication. Toute absence produit un état incomplete, partial, blocked ou unknown.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | --- | --- | --- |
| Runtime Detection and version | Command / CAP-INV-425 | runtime identity, version and targets | oui | current selected scope | quality unknown |
| Signals/Alerts/Incidents projections | Command | objects, dispositions, comments, escalations and closures | oui | time-bounded | insufficient-data |
| Health and source context | CAP-INV-426 / Settings | technical health, gaps and limitations | oui | same period | interpretation limited |
| Operational feedback | Command analysts / Incident Commander | noise, duplicates, value and missing context | non | attributed/current | quantitative-only assessment |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Runtime Detection / Signal / Alert / Incident | Command | version, relationships, disposition and operational feedback | lecture uniquement |
| Detection Health Assessment | Investigate | technical context | lecture |
| Telemetry Event / source context | Shared | source evidence where permitted | lecture |
| Comments / Assignments / Reporting | Shared | attributed feedback | consommation |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Runtime Quality Assessment | créer, comparer, contester, superseder | Investigate concept | does not mutate Command objects |
| Operational feedback annotation | enregistrer/link | Investigate / Shared | Command source and author preserved |
| Tuning Proposal context | préparer | CAP-INV-429 | no active change |

## 11. Fonctionnalités
- consommer Signals, Alerts, Incidents, dispositions, comments, escalations and closures
- voir priorities, handling time, reported noise, duplicates, missing context and operational value
- regrouper, filtrer and compare by environment/version
- documenter ground-truth limitations
- préparer Runtime Quality Assessment and tuning proposal
- conserver sources, versions, erreurs, partialité, attribution et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | --- | --- | --- | --- |
| consulter/comparer | Detection Reviewer | Runtime Quality Assessment | 0 | lecture autorisée | projection sourcée | non |
| assessment bornée | Detection Reviewer | Tool Call / assessment | 1 | lancement explicite | résultat attribué | policy |
| créer/modifier/contester | Detection Reviewer | Runtime Quality Assessment | 2 | mutation réversible | nouvelle version | OPEN-013 |
| préparer demande | Detection Reviewer | Action Request context | 2 | risque, cible et rollback visibles | package non effectif | Govern |
| changement réel | Govern/runtime owner | runtime target | 3 | Decision/Approval | Result projeté | owner |
| détruire historique | aucun rôle local | provenance | 4 | interdit | refus | strict |

Investigate n’exécute jamais les classes 3/4.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | --- | --- | --- | --- | --- |
| construire/revoir Runtime Quality Assessment | oui | formulaires et contrôles | oui | proposition | matrices et revue humaine |
| comparer/valider | oui | diff et règles explicables | oui | explication | diagnostics et checklist |
| résumer risques/erreurs | oui | agrégations | oui | résumé sourcé | tables et timeline |
| approuver/exécuter | non localement | non | non | interdit | Govern/runtime owner |

Initiateur, agent/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner et disposition humaine sont visibles. Aucun choix silencieux.

## 14. États fonctionnels
`unknown`, `insufficient-data`, `under-review`, `acceptable`, `noisy`, `low-value`, `coverage-limited`, `candidate-regression`, `candidate-improvement`, `disputed`, `superseded`. Projections fonctionnelles, pas machine d’état objet définitive.

## 15. États d’interface
Loading conserve le contexte ; Empty distingue absence et interdiction ; Partial détaille les targets ; Error conserve le valide ; Offline bloque les mutations ; Permission denied masque ; Stale distingue ancien/courant ; les conflits offrent diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Runtime Quality Assessment | assessment | CAP-INV-428/429/435 | health and quality separate, limits visible |
| Feedback grouping | analysis projection | Detection Owner / Command | source objects remain Command |
| Tuning candidate | proposal context | CAP-INV-429 | problem and evidence linked |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| Runtime Detection | produces Signal | Command Signal | runtime context only; ownership Command | Command |
| Signal | promoted/linked | Alert / Incident | Command chain and disposition | Command |
| Command feedback | available for lifecycle review | CAP-INV-427 | objects, dispositions, comments, timing, version and environment | Command |
| CAP-INV-427 | quality assessment prepared | CAP-INV-428/429 | grouped feedback, limits, health and version | Runtime Quality |

Chaque transition conserve ownership, tenant/env, target, version, permissions, restrictions, erreurs, autorité, provenance et return origin.

## 18. Dépendances
CAP-INV-001,425,426,428,429,435; Command Detection/Signal/Alert/Incident; Shared Metrics/Reporting. OPEN-017 couvre la future stratégie runtime/langage/portabilité sans option sélectionnée. Shared Jobs, Notifications, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Collaboration, Comparison, Inspector, Audit Hooks et Recovery sont consommés sans redéfinition.

## 19. Source de vérité
Investigate possède l’assessment/proposition locale. Les objets Command, Settings, Endpoint, Govern, Studio et Shared restent canoniques chez leurs owners. Une projection ne remplace jamais sa source.

## 20. Provenance et audit
Conserver besoin, Project, Hypothesis, Drafts/Versions, Review Package, Release Candidate, reviews, readiness, targets, plans, Action Requests, Decisions, Approvals, Runs, Results, runtime observations, health, Signals/Alerts/Incidents, assessments, propositions, erreurs, auteurs, timestamps, Tools/Runs, paramètres et dispositions. Aucune trace supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Signal quality / Alert / Incident feedback read | operational and sensitive context | 0 | fields scoped by Command permissions | possible | viewer/Command owner | Command | Permissions |
| Runtime Quality Assessment create/update | quality interpretation | 2 | no source mutation | OPEN-013 | reviewer/content owner | Investigate | Permissions |
| Automated feedback grouping | bias/overgeneralization | 1/2 | sources and sample visible | possible | human disposition | Studio/Investigate | Permissions |

Permissions atomiques, namespaces, RBAC/ABAC, step-up et séparation finale restent reportés.

## 22. Limites et erreurs
- Command retains all runtime object ownership.
- Signal volume ≠ detection quality; alert volume ≠ value.
- Low volume ≠ effective detection; high volume ≠ FP certainty.
- Command disposition and incident closure are not absolute ground truth.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, target offline et version superseded restent visibles.
- Result, match, health, disposition Command ou sortie IA ne vaut pas conclusion universelle.

## 23. Métriques conceptuelles
- signals/alerts/incidents by version/environment
- noise/duplicate/value feedback distribution
- insufficient-context rate
- runtime object mutations by Investigate — target zero
- provenance et dispositions humaines complètes
- exécutions silencieuses, auto-approbations et suppressions de trace : cible zéro

Aucun seuil technique universel n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; documentation only. Aucun `validated`, `implemented`, `deployed`, `active`, `native` ou `integrated`.

## 25. Critères d’acceptation
### 1. Runtime healthy mais noisy
**Given** un runtime healthy avec beaucoup de Signals et noise feedback  
**When** quality is assessed  
**Then** health and quality remain distinct, quality may be noisy and tuning may be proposed

### 2. Signal disposition incertaine
**Given** un Signal closed without sufficient Evidence  
**When** production feedback is reviewed  
**Then** the disposition is not absolute truth and may remain inconclusive

### 3. Sans IA
**Given** aucun modèle  
**When** feedback is analyzed  
**Then** filters, groupings, tables and human review remain available

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-015 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-017 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-005 reste forensic-only.
- Schémas, formats, cardinalités, permissions atomiques, contrats techniques et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering lifecycle, Command, Settings, Endpoint, Govern, Studio, Shared, Objects, Permissions, Screens, Journeys, Technique, validation et futur handoff 4B.3B non canonique.
