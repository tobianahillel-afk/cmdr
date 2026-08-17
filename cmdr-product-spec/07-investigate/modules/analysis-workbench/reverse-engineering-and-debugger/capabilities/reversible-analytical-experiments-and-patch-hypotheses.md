---
id: CAP-INV-344
title: Reversible Analytical Experiments and Patch Hypotheses
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
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-PROD-060
  - REQ-AI-002
open_decisions:
  - OPEN-005
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---

# CAP-INV-344 — Reversible Analytical Experiments and Patch Hypotheses

## 1. Définition
Reversible Analytical Experiments and Patch Hypotheses définit le comportement produit permettant de tester une Patch Hypothesis sur une copie isolée, avec confirmation, comparaison avant/après et rollback, sans modifier l’original ni déployer une modification.

## 2. Problème utilisateur
de tester une Patch Hypothesis sur une copie isolée, avec confirmation, comparaison avant/après et rollback, sans modifier l’original ni déployer une modification.

## 3. Objectifs
- formuler objectif, zone, session, modification fonctionnelle, effets attendus, risques et limites
- créer ou sélectionner une copie isolée et exiger confirmation
- appliquer l’expérience uniquement dans un environnement autorisé
- comparer avant/après, observer, revert et produire un Derived Artifact
- accepter, rejeter ou superseder l’hypothèse avec provenance

## 4. Non-objectifs
- patcher ou déployer en production
- modifier un Endpoint réel ou l’Artifact source
- fournir commandes de patching, bypass, exploit ou payload
- présenter l’hypothèse comme correction validée

## 5. Propriétaire
Investigate possède le contexte analytique, les annotations et les relations au Case. Studio conserve les Tools, Tool Calls et Automation Runs; Platform Settings administre les environnements; Govern conserve l’autorité sur toute cible réelle; Shared conserve les mécanismes transversaux.

## 6. Utilisateurs
Reverse Engineer principal; Debug Analyst, Reviewer et Investigation Lead secondaires.

## 7. Conditions d’entrée
- Artifact source and lineage available
- isolated copy permitted
- authorized environment healthy
- Patch Hypothesis ready-for-experiment
- class-2 policy/confirmation satisfied

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Artifact source/copy | CAP-INV-105/311 | immutable source and isolated derivative | Oui | versions fixed | experiment blocked |
| Patch Hypothesis | analyst | objective, functional change, risk and expected effect | Oui | current version | draft only |
| Environment/policy | Platform Settings | authorization, isolation and restrictions | Oui | current | policy-blocked |
| Runtime baseline | CAP-INV-341/343 | before state and events | Non | experiment session | comparison marked partial |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | contexte, Hypothesis et return origin | lecture et liaison uniquement |
| Artifact / Derived Artifact / Runtime Artifact | Investigate | source, versions, restrictions et provenance | lecture; aucun original modifié |
| Tool / Tool Call / Automation Run | CMDR Studio | outil, version, exécution et attribution | lecture et sélection; lifecycle Studio |
| Artifact source / isolated copy | Artifact Management | original, copy lineage and restrictions | read; original immutable |
| Reverse/Debugger context | CAP-INV-330/339 | locations, hypothesis and environment | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Trace / Activity event | Shared | émission de l’action et de son résultat | append-only; aucune suppression locale |
| Patch Hypothesis | Investigate concept | create/update/review/supersede/reject | not validated fix |
| Isolated Derived Artifact | Investigate | create/revert/withdraw-from-use | source lineage and no deployment |
| Experiment observation | Investigate | record/compare | no production effect |

## 11. Fonctionnalités
- create and review Patch Hypothesis
- identify source, isolated copy, location, purpose, expected effect, risks and limits
- require visible confirmation and authorized environment
- apply only to isolated representation without command specification
- compare before/after, observe, revert and preserve source
- produce Derived Artifact and disposition supported/contradicted/inconclusive

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Créer/modifier hypothesis | Reverse Engineer | Patch Hypothesis | 2 | source and location known | versioned draft | OPEN-013 |
| Appliquer expérience isolée | Authorized Analyst | Isolated copy | 2 | confirmation and policy | applied-in-isolated-copy | OPEN-013 |
| Observer/comparer | Reviewer | Experiment result | 0 | before/after available | differences and limits visible | Non |
| Revert/restaurer | Session owner | Isolated copy | 2 | reversible state available | reverted with trace | OPEN-013 |
| Créer Derived Artifact | Authorized Analyst | Derived Artifact | 1 | lineage complete | artifact created, no deployment | Non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Inspecter l’expérience et ses effets | Oui | Oui | Oui | Oui | viewer déterministe et navigation manuelle de l’expérience et ses effets |
| Proposer une Patch Hypothesis ou résumer avant/après | Oui | Oui | Oui | Oui | règles, heuristiques visibles, comparaison et revue humaine |
| Préparer un handoff | Oui | Oui | Oui | Oui | sélection manuelle, checklist et validation humaine |

Toute sortie automatisée expose initiateur, producteur et version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres, timestamp, statut, incertitude, owner humain et disposition acceptée, modifiée ou rejetée.

## 14. États fonctionnels
- draft
- ready-for-experiment
- applied-in-isolated-copy
- observed
- supported
- contradicted
- inconclusive
- reverted
- superseded
- rejected

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
| Patch Hypothesis | Hypothesis concept | Reverse Session/Reviewer | never a validated fix |
| Experiment result | Analysis Result concept | CAP-INV-345/346 | before/after, risks and rollback visible |
| Derived Artifact | Artifact | CAP-INV-311 | source immutable, lineage complete, no deployment |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Reverse/Debugger finding | formulate experiment | Patch Hypothesis | source, copy, location, objective, risks, authorization | return analysis |
| Patch Hypothesis | explicit confirmation | Isolated Experiment | copy, functional change, limits, rollback, provenance | return hypothesis |
| Experiment | review result | Handoff/Derived Artifact | before/after, observations, contradictions, disposition | return Workbench |

Chaque transition conserve tenant, environnement, Case, Artifact, sélection, permissions et return origin. Une erreur ne détruit pas la source ni les résultats déjà valides.

## 18. Dépendances
- CAP-INV-311
- CAP-INV-330/339/341/343
- Platform Settings environment policy
- Shared Undo/Rollback/Versioning
- OPEN-005/013/015

## 19. Source de vérité
Original Artifact remains canonical and immutable; Patch Hypothesis and experiment records belong to Investigate; Derived Artifact lineage follows CAP-INV-311; environment administration remains Settings.

## 20. Provenance et audit
Trace author, source/copy/version, location, functional change description, risks, confirmation, environment/version, Tool Calls, before/after, rollback, disposition and human reviewer.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up éventuel | Séparation des tâches | Owner | Phase propriétaire |
|---|---|---:|---|---|---|---|
| Patch Hypothesis create/review | reversible analytical modification | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| isolated experiment apply/revert | state change in isolated copy | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| Derived Artifact create/export | new sensitive artifact | 1 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| restricted environment use | environment risk | 1 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |

La matrice atomique, les namespaces et le modèle RBAC/ABAC restent reportés à la phase Permissions.

## 22. Limites et erreurs
- no original mutation
- no production/Endpoint deployment
- no bypass/exploit instructions
- rollback failure visible
- hypothesis not validated fix
- class-2 governance unresolved

## 23. Métriques
- hypotheses by disposition
- experiments reverted
- source immutability violations (must remain zero)
- results partial/inconclusive
- Derived Artifacts with complete lineage

## 24. Classification de livraison
Delivery status `defined`; delivery mode `planned`. La promotion exige objets et permissions approuvés, Tools et environnements évalués, contrats techniques, tests de sécurité, écrans et release evidence. Aucun moteur, debugger, produit tiers ou implémentation n’est choisi.

## 25. Critères d’acceptation
### Scénario 1
**Given** source Artifact, isolated copy and ready hypothesis
**When** analyst executes experiment
**Then** source unchanged; change attributed; Tool/version visible; before/after comparable; revert possible; Derived Artifact optional; no deployment

### Scénario 2
**Given** rollback cannot restore state
**When** the experiment stops
**Then** failure visible; source remains intact; isolated copy marked invalid/restricted

### Scénario 3
**Given** AI proposes a hypothesis
**When** reviewer examines
**Then** proposal attributed and no application without explicit human confirmation

## 26. Questions ouvertes
- OPEN-005 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-013 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-015 reste ouverte; aucune décision n’est fermée par cette spécification.

Les schémas, cardinalités, machines d’état, formats d’adresse et permissions atomiques sont reportés aux phases propriétaires.

## 27. Consommateurs documentaires
- INV-REV-001
- INV-DBG-001
- CAP-INV-311/345/346
- future Patch Hypothesis concept and permission model
