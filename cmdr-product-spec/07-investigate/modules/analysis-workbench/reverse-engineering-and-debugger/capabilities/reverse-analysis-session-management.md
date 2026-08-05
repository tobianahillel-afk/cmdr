---
id: CAP-INV-330
title: Reverse Analysis Session Management
product: investigate
module: analysis-workbench
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-05
requirement_ids:
  - REQ-INV-003
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-OBJ-009
  - REQ-AI-002
open_decisions:
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---

# CAP-INV-330 — Reverse Analysis Session Management

## 1. Définition
Reverse Analysis Session Management définit le comportement produit permettant de conserver l’objectif, les Artifacts, Tools, vues, annotations, interprétations et Debugger Sessions d’un travail Reverse durable et révisable.

## 2. Problème utilisateur
de conserver l’objectif, les Artifacts, Tools, vues, annotations, interprétations et Debugger Sessions d’un travail Reverse durable et révisable.

## 3. Objectifs
- créer, reprendre, suspendre, clôturer, rouvrir et superseder une session
- lier Case, Artifacts, contributors, Tools, Tool Calls, versions et paramètres
- conserver vues ouvertes, emplacements, annotations, renommages et relations
- comparer plusieurs sessions et transmettre leurs résultats

## 4. Non-objectifs
- remplacer l’Analysis Session générale
- posséder Tool, Tool Call ou Automation Run
- administrer un environnement ou lancer un debugger silencieusement

## 5. Propriétaire
Investigate possède le contexte analytique, les annotations et les relations au Case. Studio conserve les Tools, Tool Calls et Automation Runs; Platform Settings administre les environnements; Govern conserve l’autorité sur toute cible réelle; Shared conserve les mécanismes transversaux.

## 6. Utilisateurs
Reverse Engineer principal; Malware Analyst, Reviewer et Investigation Lead secondaires.

## 7. Conditions d’entrée
- Reverse Intake ready
- Artifact accessible
- owner et objectif définis
- droits de contribution évalués

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Reverse Intake | CAP-INV-329 | préconditions et scope | Oui | courante | session reste draft |
| Artifacts et versions | Artifact Management | sources de travail | Oui | versions sélectionnées | bloquer la branche de travail manquante |
| Tools/Tool Calls | CMDR Studio | catalogue et historique d’exécution | Non | version enregistrée | fonction déterministe manuelle disponible |
| Debugger Sessions liées | CAP-INV-339 | résumés, états et retours | Non | dernière synchronisation | marquer relation stale |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | contexte, Hypothesis et return origin | lecture et liaison uniquement |
| Artifact / Derived Artifact / Runtime Artifact | Investigate | source, versions, restrictions et provenance | lecture; aucun original modifié |
| Tool / Tool Call / Automation Run | CMDR Studio | outil, version, exécution et attribution | lecture et sélection; lifecycle Studio |
| Reverse Intake | Investigate | préconditions, objectif et return origin | lecture |
| Debugger Session | Investigate concept | sessions liées et résumés | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Trace / Activity event | Shared | émission de l’action et de son résultat | append-only; aucune suppression locale |
| Reverse Analysis Session | Investigate | create/update/pause/close/reopen/supersede | distincte d’Analysis Session et Automation Run |
| Session membership and view state | Investigate | versioned update | aucune mutation d’Artifact |

## 11. Fonctionnalités
- gérer lifecycle fonctionnel de session
- enregistrer owner, contributeurs, objectif, paramètres et sources
- restaurer vues, emplacements, annotations et sélection
- lier plusieurs Debugger Sessions sans les absorber
- comparer et superseder des sessions avec historique

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Créer une session | Reverse Engineer | Reverse Analysis Session | 2 | intake ready | session draft versionnée | Non |
| Suspendre ou reprendre | Owner | Reverse Analysis Session | 2 | session active/paused | état et contexte conservés | OPEN-013 |
| Clôturer ou rouvrir | Investigation Lead | Reverse Analysis Session | 2 | résultats enregistrés | disposition auditable | OPEN-013 |
| Comparer deux sessions | Reviewer | Session comparison | 0 | accès aux deux sessions | différences et limites visibles | Non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Inspecter la session Reverse | Oui | Oui | Oui | Oui | viewer déterministe et navigation manuelle de la session Reverse |
| Résumer la progression ou proposer une prochaine vue | Oui | Oui | Oui | Oui | règles, heuristiques visibles, comparaison et revue humaine |
| Préparer un handoff | Oui | Oui | Oui | Oui | sélection manuelle, checklist et validation humaine |

Toute sortie automatisée expose initiateur, producteur et version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres, timestamp, statut, incertitude, owner humain et disposition acceptée, modifiée ou rejetée.

## 14. États fonctionnels
- draft
- ready
- active
- paused
- blocked
- partial
- completed
- failed
- archived
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
| Reverse Analysis Session | Session concept | CAP-INV-331..346 | objectif, sources et contexte versionnés |
| Session summary | Analysis Result concept | Case / reviewers | erreurs et incertitudes conservées |
| Supersession relation | Version relation | Audit and navigation | ancienne session toujours consultable |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Reverse Intake | create/resume | Reverse Analysis Session | Artifact, objectif, owner, scope, restrictions | retour intake |
| Session | ouvrir une vue | Disassembly/Decompilation/Functions/Debugger | session, Artifact, emplacement, sélection | retour session |
| Session | transmettre les résultats | CAP-INV-346 | observations, annotations, contradictions, provenance | retour Workbench |

Chaque transition conserve tenant, environnement, Case, Artifact, sélection, permissions et return origin. Une erreur ne détruit pas la source ni les résultats déjà valides.

## 18. Dépendances
- CAP-INV-329
- CAP-INV-301/302 session foundations
- CAP-INV-337 annotations
- CAP-INV-339 debugger sessions
- Shared Versioning/Activity/Recovery
- OPEN-013/015

## 19. Source de vérité
La session et ses interprétations appartiennent à Investigate; Tools/Tool Calls/Automation Runs restent Studio; vues et historique consomment Shared Versioning et Activity.

## 20. Provenance et audit
Tracer création, reprise, changements d’owner/contributeurs, paramètres, vues, annotations, Tools/versions, erreurs, clôture, réouverture et supersession.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up éventuel | Séparation des tâches | Owner | Phase propriétaire |
|---|---|---:|---|---|---|---|
| Reverse Analysis Session create/update | modification du contexte partagé | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| close/reopen/supersede | impact sur le statut d’usage | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| session comparison | exposition de plusieurs analyses | 0 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| cross-tenant analysis | risque d’isolation | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |

La matrice atomique, les namespaces et le modèle RBAC/ABAC restent reportés à la phase Permissions.

## 22. Limites et erreurs
- conflit de modification visible
- session source manquante ou stale
- Tool indisponible n’efface pas le contexte
- archived reste lisible selon policy
- aucune fusion silencieuse de sessions

## 23. Métriques
- sessions actives/paused/blocked
- taux de reprise avec contexte restauré
- conflits et résolutions
- sessions superseded avec lien valide

## 24. Classification de livraison
Delivery status `defined`; delivery mode `planned`. La promotion exige objets et permissions approuvés, Tools et environnements évalués, contrats techniques, tests de sécurité, écrans et release evidence. Aucun moteur, debugger, produit tiers ou implémentation n’est choisi.

## 25. Critères d’acceptation
### Scénario 1
**Given** une session active avec vues et annotations
**When** l’owner la suspend puis la reprend
**Then** Artifact, emplacement, vues, annotations et historique sont restaurés

### Scénario 2
**Given** deux contributeurs modifient la session
**When** un conflit est détecté
**Then** aucune modification n’est écrasée silencieusement et une résolution est exigée

### Scénario 3
**Given** aucun modèle IA n’est configuré
**When** la session est créée
**Then** toutes les fonctions essentielles restent disponibles

## 26. Questions ouvertes
- OPEN-013 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-015 reste ouverte; aucune décision n’est fermée par cette spécification.

Les schémas, cardinalités, machines d’état, formats d’adresse et permissions atomiques sont reportés aux phases propriétaires.

## 27. Consommateurs documentaires
- Analysis Workbench
- INV-REV-001
- INV-DBG-001
- Case Workspace
- future objects, journeys and permissions
