---
id: CAP-INV-339
title: Debugger Session Management
product: investigate
module: analysis-workbench
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-05
requirement_ids:
  - REQ-INV-004
  - REQ-PROD-014
  - REQ-PROD-017
  - REQ-PROD-052
  - REQ-SEC-001
  - REQ-SEC-002
open_decisions:
  - OPEN-005
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---

# CAP-INV-339 — Debugger Session Management

## 1. Définition
Debugger Session Management définit le comportement produit permettant de créer, reprendre et clôturer une Debugger Session isolée, attribuée et liée à une Reverse Analysis Session sans cibler un Endpoint réel.

## 2. Problème utilisateur
de créer, reprendre et clôturer une Debugger Session isolée, attribuée et liée à une Reverse Analysis Session sans cibler un Endpoint réel.

## 3. Objectifs
- sélectionner Artifact ou copie isolée, environnement autorisé, Tool/version, objectif et owner
- conserver breakpoints, snapshots, événements, exceptions et traces
- gérer pause, reprise, clôture, réouverture, crash et timeout
- comparer plusieurs sessions et revenir à la Reverse Analysis Session

## 4. Non-objectifs
- déboguer un processus sur un Endpoint réel
- posséder ou administrer l’environnement ou le Debugger Tool
- confondre Debugger Session avec Sandbox Run, Automation Run ou Response Run

## 5. Propriétaire
Investigate possède le contexte analytique, les annotations et les relations au Case. Studio conserve les Tools, Tool Calls et Automation Runs; Platform Settings administre les environnements; Govern conserve l’autorité sur toute cible réelle; Shared conserve les mécanismes transversaux.

## 6. Utilisateurs
Reverse Engineer principal; Malware Analyst et Debug Reviewer secondaires.

## 7. Conditions d’entrée
- Reverse Analysis Session accessible
- Artifact ou copie isolée autorisée
- environnement healthy et policy-compatible
- Debugger Tool/version disponible
- permissions évaluées

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Reverse Analysis Session | CAP-INV-330 | Artifact, objective, locations and restrictions | Oui | session active | debugger session not created |
| Execution environment | Platform Settings | health, isolation and allowed use | Oui | courante | environment-unavailable |
| Debugger Tool/version | CMDR Studio | capabilities and version | Oui | courante | tool-unavailable |
| Existing breakpoints/snapshots | CAP-INV-340/341 | session context | Non | session selected | start empty |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | contexte, Hypothesis et return origin | lecture et liaison uniquement |
| Artifact / Derived Artifact / Runtime Artifact | Investigate | source, versions, restrictions et provenance | lecture; aucun original modifié |
| Tool / Tool Call / Automation Run | CMDR Studio | outil, version, exécution et attribution | lecture et sélection; lifecycle Studio |
| Reverse Analysis Session | Investigate | Artifact/copie isolée, objectif, emplacements et restrictions | lecture/lien |
| Authorized execution environment | Platform Settings | availability, health, isolation and policy | read/select only |
| Debugger Tool/version | CMDR Studio | capability and lifecycle projection | read/select only |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Trace / Activity event | Shared | émission de l’action et de son résultat | append-only; aucune suppression locale |
| Debugger Session | Investigate concept | create/update/pause/close/reopen | distinct from Sandbox/Automation/Response Run |
| Debugger session linkage | Investigate | link to Reverse Session/Case | no real Endpoint target |

## 11. Fonctionnalités
- créer/reprendre une session explicitement
- lier Reverse Session, Artifact/copie, Case, Tool et environnement
- conserver breakpoints, snapshots, events, exceptions and traces
- suspendre, reprendre, clôturer, rouvrir et comparer
- gérer crash, timeout, revocation and environment failure sans inventer d’état

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Préparer une session | Reverse Engineer | Debugger Session | 2 | Reverse Session et environnement valides | requested/validating | OPEN-013 |
| Ouvrir explicitement | Authorized Analyst | Debugger Session | 1 | ready et confirmation | opening puis active/paused | Non |
| Suspendre/reprendre | Session owner | Debugger Session | 2 | active/paused | context preserved | OPEN-013 |
| Clôturer/rouvrir | Investigation Lead | Debugger Session | 2 | results preserved | auditable disposition | OPEN-013 |
| Comparer sessions | Reviewer | Debugger comparison | 0 | access to sessions | differences and preconditions visible | Non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Inspecter la Debugger Session | Oui | Oui | Oui | Oui | viewer déterministe et navigation manuelle de la Debugger Session |
| Proposer un objectif, un breakpoint ou une reprise | Oui | Oui | Oui | Oui | règles, heuristiques visibles, comparaison et revue humaine |
| Préparer un handoff | Oui | Oui | Oui | Oui | sélection manuelle, checklist et validation humaine |

Toute sortie automatisée expose initiateur, producteur et version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres, timestamp, statut, incertitude, owner humain et disposition acceptée, modifiée ou rejetée.

## 14. États fonctionnels
- requested
- validating
- ready
- opening
- active
- paused
- stepping
- running
- stopping
- closed
- failed
- crashed
- timed-out
- environment-unavailable
- revoked

Ces états sont fonctionnels et ne constituent pas une machine d’état objet définitive.

## 15. États d’interface
- **Loading** conserve le Workbench, l’Artifact, la sélection et le return origin.
- **Empty** explique l’absence de résultat sans simuler une analyse.
- **Partial** identifie les sources, vues ou événements manquants et leurs conséquences.
- **Error** conserve les résultats valides, l’erreur et une reprise sûre.
- **Offline** limite les mutations et affiche la dernière synchronisation.
- **Permission denied** ne révèle aucune donnée protégée.
- **Stale** distingue la dernière observation connue de l’état courant.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Debugger Session | Session concept | CAP-INV-340..346 | Tool/environment, objective and restrictions attributed |
| Session failure record | Error/interruption event | analyst and audit | last known state distinct from current unknown state |
| Return context | Context reference | Reverse Analysis Session | selection, Artifact and return origin preserved |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Reverse Analysis Session | préparer debugger | Debugger Session | Artifact/copie, environment, objective, locations, proposed breakpoints, restrictions | return Reverse Session |
| Debugger Session | ouvrir/continue | Execution Control | session, state, permissions and initiator | return session |
| Debugger Session | close or failure | Reverse Analysis Session | snapshots, events, errors, disposition | return same location |

Chaque transition conserve tenant, environnement, Case, Artifact, sélection, permissions et return origin. Une erreur ne détruit pas la source ni les résultats déjà valides.

## 18. Dépendances
- CAP-INV-330
- CAP-INV-315/326 environment safety principles
- CMDR Studio Tool/Tool Call
- Platform Settings environments/health/policies
- Shared Background Jobs/Recovery
- OPEN-005/013/015

## 19. Source de vérité
Debugger Session and analytical context belong to Investigate; Tool/version belongs to Studio; environment administration belongs to Settings; no Response Run is created.

## 20. Provenance et audit
Trace actor, Artifact/copy, Reverse Session, environment/version, Tool/version, objective, permissions, lifecycle events, breakpoints, snapshots, errors, crash/timeout and disposition.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up éventuel | Séparation des tâches | Owner | Phase propriétaire |
|---|---|---:|---|---|---|---|
| Debugger Session prepare | configuration de traitement isolé | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| Debugger Session open | exécution isolée | 1 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| close/reopen | mutation de session | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| restricted environment use | risque environnemental | 1 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| cross-tenant analysis | isolation | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |

La matrice atomique, les namespaces et le modèle RBAC/ABAC restent reportés à la phase Permissions.

## 22. Limites et erreurs
- real Endpoint target blocked
- tool/environment unavailable visible
- crash preserves existing snapshots and last known state
- revocation stops new controls
- no silent open

## 23. Métriques
- sessions by state
- crashes/timeouts/environment failures
- recovery/new session rate
- return-to-Reverse context success

## 24. Classification de livraison
Delivery status `defined`; delivery mode `planned`. La promotion exige objets et permissions approuvés, Tools et environnements évalués, contrats techniques, tests de sécurité, écrans et release evidence. Aucun moteur, debugger, produit tiers ou implémentation n’est choisi.

## 25. Critères d’acceptation
### Scénario 1
**Given** une session active dont le Tool ou environnement crash
**When** l’incident survient
**Then** session becomes crashed/failed; snapshots remain; last known state distinguished; new session can be prepared

### Scénario 2
**Given** un Endpoint réel lié au Case
**When** l’analyste tente un direct debugger
**Then** action unavailable; no command; boundary to Collection/Live Response and Govern visible

### Scénario 3
**Given** aucun modèle IA
**When** l’analyste prépare et ouvre
**Then** manual Tool selection, permissions, breakpoints and controls remain available

## 26. Questions ouvertes
- OPEN-005 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-013 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-015 reste ouverte; aucune décision n’est fermée par cette spécification.

Les schémas, cardinalités, machines d’état, formats d’adresse et permissions atomiques sont reportés aux phases propriétaires.

## 27. Consommateurs documentaires
- INV-DBG-001
- INV-REV-001
- Technical Workbench/Run Shell
- Platform Settings Sandbox Environments
- future Debugger Session object and journey
