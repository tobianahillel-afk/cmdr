---
id: CAP-INV-340
title: Breakpoint and Execution Control
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
  - REQ-PROD-020
  - REQ-AI-002
  - REQ-SEC-002
open_decisions:
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---

# CAP-INV-340 — Breakpoint and Execution Control

## 1. Définition
Breakpoint and Execution Control définit le comportement produit permettant de gérer des breakpoints et une progression fonctionnelle explicite dans une Debugger Session isolée sans exposer de commandes réelles ni lancer une exécution silencieuse.

## 2. Problème utilisateur
de gérer des breakpoints et une progression fonctionnelle explicite dans une Debugger Session isolée sans exposer de commandes réelles ni lancer une exécution silencieuse.

## 3. Objectifs
- créer, activer, désactiver, modifier et retirer un breakpoint de la session
- rendre emplacement, condition fonctionnelle, source, état et atteinte visibles
- contrôler pause, reprise, étape fonctionnelle, poursuite conditionnelle et arrêt
- gérer breakpoint invalide ou emplacement non mappé
- conserver une trace et revenir à un état antérieur lorsque possible

## 4. Non-objectifs
- définir l’implémentation technique d’un breakpoint
- fournir des commandes de debugger, scripts ou patching
- contrôler un Endpoint réel ou contourner une protection

## 5. Propriétaire
Investigate possède le contexte analytique, les annotations et les relations au Case. Studio conserve les Tools, Tool Calls et Automation Runs; Platform Settings administre les environnements; Govern conserve l’autorité sur toute cible réelle; Shared conserve les mécanismes transversaux.

## 6. Utilisateurs
Reverse Engineer principal; Debug Operator autorisé et Reviewer secondaires.

## 7. Conditions d’entrée
- Debugger Session ready/active/paused
- Code Location résolue ou état unresolved visible
- permission d’exécution contrôlée
- environnement autorisé

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Debugger Session state | CAP-INV-339 | current lifecycle and environment | Oui | courante | controls disabled |
| Breakpoint definition | analyst or suggestion | location, condition, source | Non | session version | session may run without breakpoint after explicit confirm |
| Location mapping | CAP-INV-331 | address/offset mapping | Oui pour enable | current session image | invalid/unresolved |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | contexte, Hypothesis et return origin | lecture et liaison uniquement |
| Artifact / Derived Artifact / Runtime Artifact | Investigate | source, versions, restrictions et provenance | lecture; aucun original modifié |
| Tool / Tool Call / Automation Run | CMDR Studio | outil, version, exécution et attribution | lecture et sélection; lifecycle Studio |
| Debugger Session | CAP-INV-339 | state, Artifact/copy, environment and permissions | read/control within session |
| Code Location / Function | CAP-INV-331/334 | target and mapping status | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Trace / Activity event | Shared | émission de l’action et de son résultat | append-only; aucune suppression locale |
| Breakpoint | Investigate concept | create/update/enable/disable/remove-from-session | session-scoped and reversible |
| Execution control event | Investigate/Shared trace | pause/resume/step/continue/stop record | no real command text persisted as product spec |

## 11. Fonctionnalités
- create/disable/enable/update/remove session breakpoint
- show source, location, condition and reached state
- explicit pause/resume/functional step/continue-until/stop/cancel
- validate mapping before control
- preserve all control events and reversible state where supported

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Créer/modifier breakpoint | Reverse Engineer | Breakpoint | 2 | session prepared and location known/explicit unresolved | versioned breakpoint | OPEN-013 |
| Activer/désactiver/retirer | Session owner | Breakpoint | 2 | breakpoint exists | session state updated | OPEN-013 |
| Pause/reprendre/step/continue | Authorized Operator | Execution control | 2 | active/paused and policy allows | visible state transition | OPEN-013 |
| Ouvrir la session | Authorized Operator | Debugger Session | 1 | ready and confirmed | execution begins explicitly | Non |
| Stop/cancel | Authorized Operator | Debugger Session | 2 | session not closed | stopping/closed or failure visible | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Inspecter les breakpoints et contrôles | Oui | Oui | Oui | Oui | viewer déterministe et navigation manuelle de les breakpoints et contrôles |
| Proposer un breakpoint ou une prochaine étape fonctionnelle | Oui | Oui | Oui | Oui | règles, heuristiques visibles, comparaison et revue humaine |
| Préparer un handoff | Oui | Oui | Oui | Oui | sélection manuelle, checklist et validation humaine |

Toute sortie automatisée expose initiateur, producteur et version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres, timestamp, statut, incertitude, owner humain et disposition acceptée, modifiée ou rejetée.

## 14. États fonctionnels
- proposed
- enabled
- disabled
- reached
- not-reached
- invalid
- unresolved
- superseded
- removed-from-session

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
| Breakpoint set | Breakpoint concepts | Debugger Session | source, mapping and condition visible |
| Execution control event | Debug event | CAP-INV-341/343/345 | initiator, prior/new state and result |
| Control error | Error event | analyst/recovery | other breakpoints unchanged |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Reverse location | prepare breakpoint | Breakpoint configuration | session, location, source, condition, permission | return code view |
| Breakpoint | explicit control | Runtime State | session, thread, location, state, timestamp | return debugger |
| Control failure | error/revocation | Recovery or stop | last known state, error, trace | return session |

Chaque transition conserve tenant, environnement, Case, Artifact, sélection, permissions et return origin. Une erreur ne détruit pas la source ni les résultats déjà valides.

## 18. Dépendances
- CAP-INV-331/334
- CAP-INV-339
- CAP-INV-341/343
- Shared Trace/Recovery
- OPEN-013
- Platform Settings policy projection

## 19. Source de vérité
Breakpoint definitions and controls belong to Investigate session context; environment enforcement remains Settings/technical implementation; no real command syntax is defined.

## 20. Provenance et audit
Trace creator, location system/value, mapping source, condition semantics, status, enable/disable/remove, control initiator, prior/new state, result, error and timestamp.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up éventuel | Séparation des tâches | Owner | Phase propriétaire |
|---|---|---:|---|---|---|---|
| breakpoint create/update/disable/remove | reversible execution control | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| execution pause/resume/step/stop | runtime state change | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| Debugger Session open | isolated execution | 1 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| restricted environment use | environment risk | 1 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |

La matrice atomique, les namespaces et le modèle RBAC/ABAC restent reportés à la phase Permissions.

## 22. Limites et erreurs
- unmapped location remains unresolved/invalid
- no execution claimed controlled by invalid breakpoint
- control denied preserves session
- stop failure visible
- no invisible breakpoint or control

## 23. Métriques
- breakpoints by state
- reached/not-reached/unresolved
- control failures/revocations
- stop and recovery outcomes

## 24. Classification de livraison
Delivery status `defined`; delivery mode `planned`. La promotion exige objets et permissions approuvés, Tools et environnements évalués, contrats techniques, tests de sécurité, écrans et release evidence. Aucun moteur, debugger, produit tiers ou implémentation n’est choisi.

## 25. Critères d’acceptation
### Scénario 1
**Given** un breakpoint lié à un emplacement non mappé
**When** l’analyste prépare l’exécution
**Then** breakpoint marked unresolved/invalid; no claim of control; may modify/disable; trace retains error; others unchanged

### Scénario 2
**Given** une session ready
**When** l’utilisateur ne confirme pas l’ouverture
**Then** aucune exécution ne commence

### Scénario 3
**Given** une suggestion IA de breakpoint
**When** l’analyste l’examine
**Then** producer/version/sources visible and no automatic creation or execution

## 26. Questions ouvertes
- OPEN-013 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-015 reste ouverte; aucune décision n’est fermée par cette spécification.

Les schémas, cardinalités, machines d’état, formats d’adresse et permissions atomiques sont reportés aux phases propriétaires.

## 27. Consommateurs documentaires
- INV-DBG-001
- CAP-INV-331/339/341/343
- Run Shell and Trace component
- future Breakpoint object and permission model
