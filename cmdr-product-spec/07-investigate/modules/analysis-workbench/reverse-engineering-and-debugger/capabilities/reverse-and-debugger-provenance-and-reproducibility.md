---
id: CAP-INV-345
title: Reverse and Debugger Provenance and Reproducibility
product: investigate
module: analysis-workbench
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-05
requirement_ids:
  - REQ-INV-003
  - REQ-INV-004
  - REQ-PROD-020
  - REQ-OBJ-009
  - REQ-AI-002
open_decisions:
  - OPEN-005
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---

# CAP-INV-345 — Reverse and Debugger Provenance and Reproducibility

## 1. Définition
Reverse and Debugger Provenance and Reproducibility définit le comportement produit permettant de retracer et évaluer la reproductibilité du travail Reverse et Debugger sans dupliquer Trace, Activity ou Audit ni définir l’infrastructure technique.

## 2. Problème utilisateur
de retracer et évaluer la reproductibilité du travail Reverse et Debugger sans dupliquer Trace, Activity ou Audit ni définir l’infrastructure technique.

## 3. Objectifs
- relier Case, Artifacts, sessions, Tools, versions, Tool Calls, Automation Runs, environnement et paramètres
- retracer fonctions, symboles, annotations, breakpoints, controls, snapshots, events, exceptions and experiments
- distinguer reproducible, partial, not reproducible et causes manquantes
- conserver erreurs, interruptions, décisions humaines et dispositions finales
- permettre contestation et supersession

## 4. Non-objectifs
- définir stockage, schéma de trace ou infrastructure de replay
- absorber les objets Shared/Studio/Settings
- présenter non-reproduction comme absence de comportement ou menace

## 5. Propriétaire
Investigate possède le contexte analytique, les annotations et les relations au Case. Studio conserve les Tools, Tool Calls et Automation Runs; Platform Settings administre les environnements; Govern conserve l’autorité sur toute cible réelle; Shared conserve les mécanismes transversaux.

## 6. Utilisateurs
Reviewer et Audit Analyst principaux; Reverse Engineer, Debug Analyst et Investigation Lead secondaires.

## 7. Conditions d’entrée
- au moins une session ou result reference
- identity/version sources available or gaps explicit
- permission to view provenance

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Sessions/results | CAP-INV-329..344 | contexts, events, knowledge and outputs | Oui | selected versions | assessment partial |
| Tool/automation provenance | CMDR Studio | Tool versions, Tool Calls, Automation Runs | Non | recorded execution | missing-tool/version state |
| Environment provenance | Platform Settings | environment/version/policy | Non | recorded session | missing-environment |
| Human decisions | Investigate/Shared Activity | accept/modify/reject/revert/disposition | Oui | event time | assessment disputed |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | contexte, Hypothesis et return origin | lecture et liaison uniquement |
| Artifact / Derived Artifact / Runtime Artifact | Investigate | source, versions, restrictions et provenance | lecture; aucun original modifié |
| Tool / Tool Call / Automation Run | CMDR Studio | outil, version, exécution et attribution | lecture et sélection; lifecycle Studio |
| Reverse/Debugger sessions and results | CAP-INV-329..344 | all analytical context and outputs | read |
| Tool/Tool Call/Automation Run | CMDR Studio | versions, parameters, status and sources | read |
| Environment projection | Platform Settings | environment/version and availability | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Trace / Activity event | Shared | émission de l’action et de son résultat | append-only; aucune suppression locale |
| Provenance relation / Reproducibility Assessment | Investigate concepts + Shared trace | create/review/dispute/supersede | does not duplicate Trace/Activity/Audit |

## 11. Fonctionnalités
- assemble a provenance chain without copying canonical records
- show initiators, Tools, versions, parameters, timestamps, errors and interruptions
- show human decisions and dispositions
- evaluate reproducibility and causes
- compare attempts, dispute and supersede an assessment

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter provenance | Reviewer | Provenance chain | 0 | references accessible | source links and gaps visible | Non |
| Reproduire session | Authorized Analyst | Reproduction attempt | 1 | inputs/tools/environment authorized | new attempt linked | Non |
| Qualifier/disputer assessment | Reviewer | Reproducibility Assessment | 2 | evidence reviewed | versioned disposition | OPEN-013 |
| Exporter provenance | Auditor | Provenance package | 1 | policy allows | redacted attributed package | Non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Inspecter la provenance et reproductibilité | Oui | Oui | Oui | Oui | viewer déterministe et navigation manuelle de la provenance et reproductibilité |
| Résumer les écarts ou proposer une cause de non-reproduction | Oui | Oui | Oui | Oui | règles, heuristiques visibles, comparaison et revue humaine |
| Préparer un handoff | Oui | Oui | Oui | Oui | sélection manuelle, checklist et validation humaine |

Toute sortie automatisée expose initiateur, producteur et version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres, timestamp, statut, incertitude, owner humain et disposition acceptée, modifiée ou rejetée.

## 14. États fonctionnels
- reproducible
- partially-reproducible
- not-reproducible
- missing-tool
- missing-version
- missing-environment
- missing-input
- state-not-restorable
- behavior-not-reproduced
- disputed

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
| Provenance chain | Linked references | Case/Review/Audit | canonical records remain at their owners |
| Reproducibility Assessment | Assessment concept | CAP-INV-346 | cause, attempts and uncertainty visible |
| Reproduction attempt | Session/run relation | Reverse/Debugger Session | new attempt never overwrites prior one |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Any Reverse/Debugger result | review provenance | Provenance view | sources, sessions, Tools, parameters, timestamps | return source |
| Assessment | request reproduction | New authorized session/attempt | inputs, versions, environment requirements, restrictions | return assessment |
| Provenance result | prepare handoff | CAP-INV-346 | selected sources, gaps, disputes and assessment | return Workbench |

Chaque transition conserve tenant, environnement, Case, Artifact, sélection, permissions et return origin. Une erreur ne détruit pas la source ni les résultats déjà valides.

## 18. Dépendances
- CAP-INV-312/327 existing provenance
- CAP-INV-329..344
- CMDR Studio provenance
- Platform Settings environment metadata
- Shared Trace/Activity/Audit/Versioning
- OPEN-005/013/015

## 19. Source de vérité
Canonical source records remain owned by Investigate, Studio, Settings or Shared; this capability links them and owns only the analytical reproducibility assessment.

## 20. Provenance et audit
Trace every reference, initiator, producer/version, parameters, environment, timestamp, status/error, interruption, decision, attempt, comparison, dispute and final disposition.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up éventuel | Séparation des tâches | Owner | Phase propriétaire |
|---|---|---:|---|---|---|---|
| provenance/reproducibility read | exposition de sources sensibles | 0 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| reproduction attempt | isolated execution | 1 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| assessment/dispute | qualification réversible | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| provenance export | diffusion | 1 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |

La matrice atomique, les namespaces et le modèle RBAC/ABAC restent reportés à la phase Permissions.

## 22. Limites et erreurs
- missing source visible
- state-not-restorable not silently reconstructed
- behavior-not-reproduced not exculpatory
- disputed assessment preserved
- no duplicate audit record

## 23. Métriques
- assessments by state
- missing tool/version/input/environment
- reproduction attempts/outcomes
- disputes and supersessions

## 24. Classification de livraison
Delivery status `defined`; delivery mode `planned`. La promotion exige objets et permissions approuvés, Tools et environnements évalués, contrats techniques, tests de sécurité, écrans et release evidence. Aucun moteur, debugger, produit tiers ou implémentation n’est choisi.

## 25. Critères d’acceptation
### Scénario 1
**Given** une session manque Tool version
**When** reviewer evaluates reproducibility
**Then** missing-version visible and no reproducible claim

### Scénario 2
**Given** un comportement n’est pas reproduit
**When** assessment is recorded
**Then** behavior-not-reproduced remains distinct from absence of threat

### Scénario 3
**Given** aucun modèle IA
**When** provenance is reviewed
**Then** linked records, deterministic comparisons and manual assessment remain available

## 26. Questions ouvertes
- OPEN-005 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-013 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-015 reste ouverte; aucune décision n’est fermée par cette spécification.

Les schémas, cardinalités, machines d’état, formats d’adresse et permissions atomiques sont reportés aux phases propriétaires.

## 27. Consommateurs documentaires
- Case Replay
- INV-REV-001/INV-DBG-001
- CAP-INV-346
- Shared Trace/Activity/Audit
- future Provenance/Reproducibility concepts
