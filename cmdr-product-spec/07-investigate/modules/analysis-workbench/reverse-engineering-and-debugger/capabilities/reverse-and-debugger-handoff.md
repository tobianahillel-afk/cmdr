---
id: CAP-INV-346
title: Reverse and Debugger Handoff to Evidence, Findings and Detection Engineering
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
  - REQ-OBJ-004
  - REQ-AI-002
open_decisions:
  - OPEN-013
  - OPEN-014
  - OPEN-015
source-of-truth: canonical
---

# CAP-INV-346 — Reverse and Debugger Handoff to Evidence, Findings and Detection Engineering

## 1. Définition
Reverse and Debugger Handoff to Evidence, Findings and Detection Engineering définit le comportement produit permettant de sélectionner et transmettre des résultats Reverse/Debugger avec provenance vers Evidence, Finding Draft et future Detection Engineering sans qualification, confirmation ou déploiement automatique.

## 2. Problème utilisateur
de sélectionner et transmettre des résultats Reverse/Debugger avec provenance vers Evidence, Finding Draft et future Detection Engineering sans qualification, confirmation ou déploiement automatique.

## 3. Objectifs
- sélectionner observations, functions, refs, graphs, annotations, snapshots, events, exceptions and Artifacts
- relier la sélection à une Hypothesis et aux contradictions
- préparer Evidence candidate et Finding Draft
- préparer une connaissance réutilisable pour future Detection Engineering
- préserver sessions, provenance et return origin

## 4. Non-objectifs
- qualifier automatiquement Evidence ou confirmer Finding
- créer ou déployer une règle Detection Engineering
- transformer annotation, decompilation, snapshot, exception ou Patch Hypothesis en conclusion

## 5. Propriétaire
Investigate possède le contexte analytique, les annotations et les relations au Case. Studio conserve les Tools, Tool Calls et Automation Runs; Platform Settings administre les environnements; Govern conserve l’autorité sur toute cible réelle; Shared conserve les mécanismes transversaux.

## 6. Utilisateurs
Investigation Lead ou Evidence Reviewer principal; Reverse Engineer, Debug Analyst et Detection Engineer receiver secondaires.

## 7. Conditions d’entrée
- Case actif
- sources sélectionnées et accessibles
- provenance sufficient or gaps explicit
- permissions de handoff

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Reverse/Debugger selections | CAP-INV-329..345 | observations and results | Oui | selected versions | handoff blocked or partial |
| Hypothesis and existing Evidence | Investigate | reasoning context | Non | Case current | candidate without relation |
| Provenance assessment | CAP-INV-345 | sources, gaps and disputes | Oui | latest selected | uncertainty mandatory |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | contexte, Hypothesis et return origin | lecture et liaison uniquement |
| Artifact / Derived Artifact / Runtime Artifact | Investigate | source, versions, restrictions et provenance | lecture; aucun original modifié |
| Tool / Tool Call / Automation Run | CMDR Studio | outil, version, exécution et attribution | lecture et sélection; lifecycle Studio |
| Reverse/Debugger observations and results | CAP-INV-329..345 | functions, refs, graphs, annotations, snapshots, events, exceptions, artifacts and assessments | read/select |
| Hypothesis/Evidence/Finding | Investigate | reasoning and qualified context | read/link |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Trace / Activity event | Shared | émission de l’action et de son résultat | append-only; aucune suppression locale |
| Evidence candidate package | Investigate | prepare/submit to CAP-INV-107/108 | not Evidence until qualification |
| Finding Draft | Investigate | prepare/submit to CAP-INV-109 | not confirmed Finding |
| Detection Engineering handoff package | Investigate | prepare only | no rule creation or deployment |

## 11. Fonctionnalités
- select heterogeneous results and contradictions
- explain relation to Hypothesis
- prepare Evidence candidate with source links
- prepare Finding Draft with uncertainty and existing Evidence
- prepare future Detection Engineering knowledge with behavior, conditions, limits and sources
- submit to owner capabilities and return to Workbench

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Sélectionner résultats | Analyst | Handoff selection | 0 | sources readable | selection visible | Non |
| Préparer Evidence candidate | Evidence Reviewer | Candidate package | 2 | provenance reviewed | submitted to 107/108 | OPEN-013 |
| Préparer Finding Draft | Investigation Lead | Finding Draft | 2 | Hypothesis and evidence context | submitted to 109 | OPEN-013 |
| Préparer Detection handoff | Detection Engineer/Analyst | Knowledge package | 2 | scope and limits explicit | future intake only, no rule | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Inspecter le handoff et ses contradictions | Oui | Oui | Oui | Oui | viewer déterministe et navigation manuelle de le handoff et ses contradictions |
| Proposer une Evidence candidate, un Finding Draft ou une connaissance Detection | Oui | Oui | Oui | Oui | règles, heuristiques visibles, comparaison et revue humaine |
| Préparer un handoff | Oui | Oui | Oui | Oui | sélection manuelle, checklist et validation humaine |

Toute sortie automatisée expose initiateur, producteur et version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres, timestamp, statut, incertitude, owner humain et disposition acceptée, modifiée ou rejetée.

## 14. États fonctionnels
- draft
- incomplete
- ready-for-review
- submitted
- accepted
- returned
- rejected
- superseded

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
| Evidence candidate | Candidate package | CAP-INV-107/108 | annotation/decompilation/snapshot remain sources, qualification required |
| Finding Draft | Draft | CAP-INV-109 | exception/crash/Patch Hypothesis not confirmation |
| Detection Engineering knowledge | Handoff package | Future Phase 4B.3 | no rule created or deployed |
| Return context | Context reference | Reverse/Debugger Workbench | session, selection and origin preserved |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Reverse/Debug result | prepare Evidence | CAP-INV-107/108 | sources, sessions, Artifacts, uncertainty, qualification need | return Workbench |
| Reverse/Debug result | prepare Finding | CAP-INV-109 | Hypothesis, existing Evidence, observations, contradictions, author, provenance | return Workbench |
| Reverse knowledge | prepare future detection | Future Detection Engineering | behavior, functions, strings, structures, conditions, limits, sources | return Workbench |

Chaque transition conserve tenant, environnement, Case, Artifact, sélection, permissions et return origin. Une erreur ne détruit pas la source ni les résultats déjà valides.

## 18. Dépendances
- CAP-INV-103/107/108/109
- CAP-INV-313/328 existing handoff patterns
- CAP-INV-329..345
- Future Phase 4B.3
- Shared Object Linking/Export/Notifications
- OPEN-013/014/015

## 19. Source de vérité
Sources remain at their capabilities; Evidence and Finding ownership remains CAP-INV-107..109; Detection package is only an Investigate handoff record for future owner intake.

## 20. Provenance et audit
Trace selected sources and versions, sessions, Artifact lineage, uncertainty, contradictions, author, automated suggestions, human disposition, submission result and return origin.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up éventuel | Séparation des tâches | Owner | Phase propriétaire |
|---|---|---:|---|---|---|---|
| Evidence candidate prepare | qualification risk | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| Finding Draft prepare | conclusion risk | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| Detection Engineering handoff prepare | future rule risk | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| handoff package export | diffusion | 1 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |

La matrice atomique, les namespaces et le modèle RBAC/ABAC restent reportés à la phase Permissions.

## 22. Limites et erreurs
- annotation not Evidence
- decompilation not source truth/Evidence
- snapshot/value not Evidence
- exception not Finding
- Patch Hypothesis not confirmed Finding
- no automatic rule deployment

## 23. Métriques
- candidate packages ready/submitted/returned
- Finding Drafts accepted/rejected
- Detection handoffs prepared
- packages missing provenance

## 24. Classification de livraison
Delivery status `defined`; delivery mode `planned`. La promotion exige objets et permissions approuvés, Tools et environnements évalués, contrats techniques, tests de sécurité, écrans et release evidence. Aucun moteur, debugger, produit tiers ou implémentation n’est choisi.

## 25. Critères d’acceptation
### Scénario 1
**Given** une exception, snapshot et trace sélectionnés
**When** l’analyste prépare Evidence candidate
**Then** each remains source; provenance referenced; qualification required; no Finding automatic

### Scénario 2
**Given** une Patch Hypothesis supported in isolated copy
**When** a Finding Draft is prepared
**Then** hypothesis remains experimental and contradictions/limits are included

### Scénario 3
**Given** une connaissance est envoyée vers future Detection Engineering
**When** handoff occurs
**Then** no rule is created/deployed and sources/conditions/limits are preserved

## 26. Questions ouvertes
- OPEN-013 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-014 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-015 reste ouverte; aucune décision n’est fermée par cette spécification.

Les schémas, cardinalités, machines d’état, formats d’adresse et permissions atomiques sont reportés aux phases propriétaires.

## 27. Consommateurs documentaires
- Evidence Creation/Review
- Finding Management
- Future Detection Engineering Phase 4B.3
- INV-REV-001/INV-DBG-001
- Case Workspace
