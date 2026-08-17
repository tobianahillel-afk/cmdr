---
id: CAP-INV-435
title: Detection Lifecycle Provenance and Continuous Improvement Handoff
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
  - OPEN-014
  - OPEN-015
  - OPEN-017
source-of-truth: canonical
---
# CAP-INV-435 — Detection Lifecycle Provenance and Continuous Improvement Handoff

## 1. Définition
Retracer de bout en bout le besoin, l’authoring, la revue, l’autorité Govern, le changement runtime, les observations, le feedback, les propositions, le rollback et le retirement, puis préparer un Continuous Improvement Package non effectif vers un nouveau cycle.

## 2. Problème utilisateur
Sans lineage unique, un match, un Result, un tuning ou un rollback peut être dissocié de la version, des targets, des décisions et des limites qui lui donnent son sens.

## 3. Objectifs
- retracer initial need, Findings, Cases, Hunts and technical analyses
- retracer Project, Hypothesis, Drafts, Versions, Review Packages, Release Candidates and Reviews
- retracer readiness, plans, requests, Decisions, Approvals, Runs, Results and targets
- retracer runtime versions, health, shadow/canary, Signals/Alerts/Incidents and Command feedback
- retracer reviews, tuning, suppressions, exceptions, drift, performance, rollback and retirement
- préparer Continuous Improvement Package vers new Draft/Hypothesis/Data Source Gap/Settings/reporting/future Intelligence handoff

## 4. Non-objectifs
Aucune API, protocole, syntaxe vendor, moteur, langage, compilateur, parser, AST, format technique, commande, code, modèle ML, mutation runtime directe, capability CAP-INV-5xx, Threat Intelligence, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède **Lifecycle Provenance Assessment** et **Continuous Improvement Package** comme concepts fonctionnels. Command garde Detection/Signal/Alert/Incident ; Settings les runtimes, targets, environments et health ; Endpoint ses capacités et exécutions locales ; Govern Action Request/Decision/Approval/Response Run/Result ; Studio Tools et Automation Runs ; Shared les mécanismes transverses.

## 6. Utilisateurs
Principal : **Detection Owner**. Secondaires : Detection Engineer, Detection Reviewer, Incident Commander, Platform Operator, Approver, Auditor et Product Lead autorisés.

## 7. Conditions d’entrée
Les identifiants, versions, tenants, environnements, targets, auteurs, timestamps, permissions, restrictions, sources et return origins de chaque étape sont résolubles ou explicitement manquants. Les traces partielles restent visibles.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | --- | --- | --- |
| Authoring lineage | CAP-INV-401..417 | need, Project, Hypothesis, Drafts, validations, replay, reviews, coverage and package | oui | canonical versions | lifecycle partial |
| Lifecycle lineage | CAP-INV-418..434 | candidates, reviews, readiness, plans, runtime observations and proposals | oui | canonical versions | lifecycle partial |
| Authority/execution lineage | Govern | Action Requests, Decisions, Approvals, Runs and Results | selon changement | immutable links | execution unverified |
| Runtime/operational lineage | Command / Settings / Endpoint | Detection versions, health, Signals, Alerts, Incidents and target state | selon période | source projections | runtime evidence missing |
| Automation lineage | Studio / Shared | Tools, Tool Calls, Automation Runs, Trace and Activity | selon usage | immutable/versioned | automation unreproducible |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Project/Drafts/Reviews/Assessments/Proposals | Investigate | full lifecycle relations | lecture/lien |
| Runtime Detection / Signal / Alert / Incident | Command | operational lineage and feedback | lecture/lien |
| Environment / Target / Runtime observations | Settings / Endpoint | version, health and execution context | lecture |
| Action Request / Decision / Approval / Run / Result | Govern | authority and effect chain | lecture/lien |
| Tool / Tool Call / Automation Run | Studio | automated contributor lineage | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Lifecycle Provenance Assessment | créer, compléter, contester, superseder | Investigate concept | missing links visible |
| Continuous Improvement Package | créer, versionner, retirer, superseder | Investigate concept | package ≠ active change |
| Cross-product lineage relation | lier/versionner | Shared/owner products | no ownership transfer |
| Future Intelligence handoff candidate | préparer | future 4B.3B | no CAP-INV-5xx or canonical Intelligence object |

## 11. Fonctionnalités
- retracer initial need, Findings, Cases, Hunts and technical analyses
- retracer Project, Hypothesis, Drafts, Versions, Review Packages, Release Candidates and Reviews
- retracer readiness, plans, requests, Decisions, Approvals, Runs, Results and targets
- retracer runtime versions, health, shadow/canary, Signals/Alerts/Incidents and Command feedback
- retracer reviews, tuning, suppressions, exceptions, drift, performance, rollback and retirement
- préparer Continuous Improvement Package vers new Draft/Hypothesis/Data Source Gap/Settings/reporting/future Intelligence handoff
- conserver sources, versions, erreurs, partialité, attribution et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | --- | --- | --- | --- |
| consulter/naviguer/comparer | Detection Owner | lifecycle provenance | 0 | lecture autorisée | lineage sourcé | non |
| vérifier/reproduire une relation | Detection Owner | provenance assessment | 1 | sources disponibles | résultat attribué | policy |
| créer/compléter/contester package | Detection Owner | Continuous Improvement Package | 2 | mutation réversible | nouvelle version | OPEN-013 |
| préparer nouveau cycle/request | Detection Owner | CAP-INV-401/402/406/Settings context | 2 | objectifs, risques et sources visibles | package non effectif | owner |
| changement réel | Govern/runtime owner | runtime target | 3 | Decision/Approval | Result projeté | obligatoire |
| détruire trace/historique | aucun rôle local | provenance | 4 | interdit | refus | strict |

Investigate n’exécute jamais les classes 3/4.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | --- | --- | --- | --- | --- |
| construire/revoir lifecycle provenance | oui | link validation and diff | oui | gap suggestion | trace graph plus table alternative |
| préparer improvement package | oui | templates/checklists | oui | draft | structured form and review |
| expliquer incohérences | oui | deterministic diagnostics | oui | sourced explanation | inspector and comparison |
| approuver/exécuter change | non localement | non | non | interdit | Govern/runtime owner |

Initiateur, agent/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner et disposition humaine sont visibles. Aucun choix silencieux.

## 14. États fonctionnels
`draft`, `collecting`, `complete-for-declared-scope`, `partial`, `missing-authority-link`, `missing-runtime-observation`, `inconsistent`, `disputed`, `reviewed`, `handoff-ready`, `superseded`, `withdrawn`. Projections fonctionnelles, pas machine d’état objet définitive.

## 15. États d’interface
Loading conserve le contexte ; Empty distingue absence et interdiction ; Partial détaille les links manquants ; Error conserve le valide ; Offline bloque les mutations ; Permission denied masque ; Stale distingue ancien/courant ; les conflits offrent timeline, diff, alternative tabulaire et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Lifecycle Provenance Assessment | assessment | Audit/Reporting/owners | sources, gaps and contradictions visible |
| Continuous Improvement Package | versioned package | CAP-INV-401/402/403/406 | no active mutation; return origin retained |
| Data Source/Settings request context | request | CAP-INV-404 / Settings | gap, impact and evidence linked |
| Future Intelligence handoff candidate | non-canonical package | future 4B.3B | candidate only, no Intelligence object created |
| Canonical reporting preparation | report input | Shared Reporting | local package ≠ canonical Report |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-401..434 | lifecycle event/disposition | CAP-INV-435 | IDs, versions, sources, owners, timestamps, errors and authority | source object |
| CAP-INV-435 | improvement identified | CAP-INV-401/402/403/406 | problem, evidence, runtime feedback, risks and unresolved gaps | Lifecycle |
| CAP-INV-435 | data/platform gap identified | CAP-INV-404 / Settings request | source, field, target, impact and desired outcome | Lifecycle |
| CAP-INV-435 | reporting requested | Shared Reporting | lineage, assessments and restrictions | Lifecycle |
| CAP-INV-435 | future Intelligence relevance | future 4B.3B | candidate observables/context and provenance | Lifecycle |

Chaque transition conserve ownership, tenant/env, target, version, permissions, restrictions, erreurs, autorité, provenance et return origin.

## 18. Dépendances
CAP-INV-401..434; Command; Settings; Endpoint; Govern; Studio; Shared Trace, Activity, Timeline, Linking, Versioning, Reporting and Audit. OPEN-017 couvre la future stratégie runtime/langage/portabilité sans option sélectionnée. Le handoff Intelligence reste futur et non canonique.

## 19. Source de vérité
Chaque objet reste canonique chez son owner. CAP-INV-435 possède uniquement la relation et l’assessment fonctionnelle locale ; elle ne remplace ni Trace, ni Activity, ni Audit, ni Reporting.

## 20. Provenance et audit
Conserver explicitement besoin initial, Finding, Case, Hunt, analyses, Project, Hypothesis, Drafts, Versions, Review Packages, Release Candidates, Reviews, Readiness, Plans, Requests, Decisions, Approvals, Runs, Results, targets, runtime versions, health, shadow/canary assessments, Signals, Alerts, Incidents, Command feedback, reviews, tuning, suppressions, exceptions, drift, performance, rollback, retirement, décisions humaines, Automation Runs et Tool Calls. Aucune trace supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Lifecycle provenance read/export | cross-product sensitive lineage | 0/1 | evidence, secrets and tenant data scoped | possible | viewer/auditor | owners/Shared | Permissions |
| Continuous Improvement Package create/update | future change influence | 2 | no active secret/action | OPEN-013 | author/reviewer | Investigate | Permissions |
| Future handoff/request prepare | scope expansion | 2 | candidate data minimized | step-up possible | requester/owner | future owner/Govern | Permissions |
| Trace deletion | irrecoverable loss | 4 | not available | strict | prohibited | Shared/Govern | Permissions |

Permissions atomiques, namespaces, RBAC/ABAC, step-up et séparation finale restent reportés.

## 22. Limites et erreurs
- Local report/package ≠ canonical Report.
- AI recommendation ≠ review or Govern decision.
- Missing link remains visible; no relation is invented.
- Continuous Improvement Package applies no active modification.
- Future Intelligence handoff creates no Indicator, Campaign, Threat Actor or CAP-INV-5xx.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, target offline et version superseded restent visibles.
- Result, match, health, disposition Command ou sortie IA ne vaut pas conclusion universelle.

## 23. Métriques conceptuelles
- lifecycle links resolved/missing by owner
- time from runtime feedback to new Draft
- packages by destination and disposition
- active changes caused by package — target zero
- provenance et dispositions humaines complètes
- exécutions silencieuses, auto-approbations et suppressions de trace : cible zéro

Aucun seuil technique universel n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; documentation only. Aucun `validated`, `implemented`, `deployed`, `active`, `native` ou `integrated`.

## 25. Critères d’acceptation
### 1. Provenance partielle
**Given** a deployment Result without runtime confirmation  
**When** lifecycle provenance is reviewed  
**Then** the missing runtime observation remains visible and no active state is invented

### 2. Continuous improvement
**Given** noisy runtime feedback and a tuning proposal  
**When** an improvement package is prepared  
**Then** sources, limitations and lineage are included and the package returns to a new Draft without changing active content

### 3. Sans IA
**Given** aucun modèle  
**When** lifecycle provenance and handoff are prepared  
**Then** trace links, tables, timelines, diff and human review suffice

## 26. Questions ouvertes
- OPEN-008 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-014 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-015 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-017 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-005 reste forensic-only.
- Schémas, formats, cardinalités, permissions atomiques, contrats techniques et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering lifecycle, Command, Settings, Endpoint, Govern, Studio, Shared, Objects, Permissions, Screens, Journeys, Technique, validation, Reporting and future 4B.3B handoff.
