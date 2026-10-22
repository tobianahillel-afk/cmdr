---
id: CAP-INV-379
title: Disk Forensics Handoff to Evidence, Findings and Future Analysis
product: investigate
module: analysis-workbench
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-05
requirement_ids:
  - REQ-INV-001
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-OBJ-004
  - REQ-INV-006
  - REQ-AI-002
open_decisions:
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-379 — Disk Forensics Handoff to Evidence, Findings and Future Analysis

## 1. Définition
Sélectionner des observations et Artifacts Disk/Filesystem, expliquer leur relation à l’Hypothesis et préparer des packages vers Evidence, Finding, Static, Reverse, Memory et futures Detection/Network phases sans qualification ou exécution automatique.

## 2. Problème utilisateur
Un fichier, journal record, user artifact, persistence candidate ou carved content peut être promu trop tôt en Evidence, Finding, IOC ou règle si contradictions, partialité and provenance are not retained.

## 3. Objectifs
Select files/metadata/journal records/system/user/application/persistence/timeline/recovered/Derived Artifacts and contradictions; explain Hypothesis relation; prepare Evidence candidate, Finding Draft, Static/Reverse/Memory handoff, future Detection package and future Network package; send to CAP-INV-107/108/109; preserve provenance and return origin.

## 4. Non-objectifs
Ne pas qualify Evidence, confirm Finding/persistence/intent, deploy detection rule, conduct Network Forensics, modify source, execute file, define API/protocol or create Report canonical.

## 5. Propriétaire
Investigate owns selection and candidate/draft package. CAP-INV-107/108 own Evidence qualification, CAP-INV-109 Finding management, Static/Reverse/Memory own destination analysis, future phases own Detection/Network work.

## 6. Utilisateurs
Principal : Investigation Lead or Evidence Reviewer. Secondaires : Disk Analyst, Malware/Reverse/Memory analysts and future Detection/Network analysts.

## 7. Conditions d’entrée
Selected observations/results with source, uncertainty, restrictions and provenance; Case/Hypothesis accessible; destination permission; no hidden sensitive content.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Selected observations/Artifacts | CAP-INV-365..378 | evidence basis candidates | yes | source versions | incomplete |
| Case/Hypothesis/existing Evidence | Investigate | reasoning context | yes | current | draft-only |
| Contradictions/uncertainty | analyses | limitations | yes if present | linked | invalid package if hidden |
| Provenance/reproducibility | CAP-INV-378 | source/producer/status | yes | current assessment | blocked/partial |
| Destination permissions/restrictions | owners/Security | allowed handoff | yes | current | denied |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case/Hypothesis/Artifact/Evidence/Finding | Investigate | context/current state | read/link |
| Disk Image/observations/results | Investigate | selected sources | read |
| Tool/Tool Call/Run | Studio | producer/provenance | read |
| Provenance/Trace | Investigate/Shared | chain/gaps | read |
| Destination capabilities | owners | intake requirements | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Evidence Candidate Package | create/update/withdraw | Investigate | observation ≠ Evidence |
| Finding Draft | create/update/withdraw | Investigate | candidate ≠ confirmed Finding |
| Analysis Handoff Package | create/version | Investigate | destination retains ownership |
| Detection/Network future package | prepare only | future owners | no rule/network analysis |

## 11. Fonctionnalités
Select all relevant categories and contradictions; explain relation to Hypothesis; preview package and restrictions; prepare/send to Evidence/Finding/Static/Reverse/Memory; prepare future Detection/Network knowledge with explicit limitations; preserve source/session/provenance and return to Workbench.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Select/preview | Analyst | package selection | 0 | read | included/excluded visible | no |
| Create/update/withdraw package | Analyst | candidate/draft | 2 | provenance | versioned package | OPEN-013 |
| Submit to owner capability | Reviewer | handoff | 2 | destination permission | owner review starts | no |
| Export authorized package | Reviewer | export | 1/2 | policy/redaction | controlled output | policy |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| collect selected context | yes | relations | yes | summary | checklist/package builder |
| propose Evidence/Finding/Detection knowledge | yes | rules | yes as draft | yes | manual selection/template |
| validate contradictions/provenance | yes | validators | yes | explanation | explicit review |
| qualify/confirm/deploy | owner human process | no | prohibited here | prohibited | destination capability review |

## 14. États fonctionnels
`draft`, `incomplete`, `ready-for-review`, `submitted`, `accepted-by-destination`, `returned`, `withdrawn`, `superseded`, `blocked`, `partial`.

## 15. États d’interface
Loading preserves selection; Empty requires explicit source selection; Partial shows missing sources; Error keeps draft; Offline blocks submit; Permission denied masks protected items; Stale requires source revalidation.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Evidence Candidate Package | candidate | CAP-INV-107/108 | not qualified; provenance/uncertainty retained |
| Finding Draft | draft | CAP-INV-109 | not confirmed; contradictions retained |
| Static/Reverse/Memory package | handoff | owner capability | source/goal/restrictions/return origin |
| Future Detection/Network package | knowledge package | Phase 4B.3 / 4B.2B.3B.2 | no rule, IOC or network analysis created |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Disk observations/results | prepare package | CAP-INV-379 | selected items, Hypothesis, contradictions, provenance | source view |
| Candidate package | submit | CAP-INV-107/108/109 | sources, uncertainty, restrictions, analyst | Disk Workbench |
| File/Derived Artifact | analyze | Static/Reverse/Memory | parent/source/context/objective | Disk Workbench |
| Persistent network artifact | future handoff | Phase 4B.2B.3B.2 | config/address candidate, source, no connection claim | Disk Workbench |

## 18. Dépendances
CAP-INV-103/105/107/108/109/301/329/347..378, Shared Linking/Export/Trace, Security, future 4B.2B.3B.2/4B.3, OPEN-008/013/014/015.

## 19. Source de vérité
Each observation/Artifact remains source; Investigate owns package; destination owners decide qualification/analysis. Local report/package is not canonical Report.

## 20. Provenance et audit
Case/Hypothesis, selected/excluded items, source image/session/filesystem, Tool/version, uncertainty/contradictions, provenance status, author/reviewer, package versions, submissions, destination dispositions and return origin.

## 21. Permissions fonctionnelles
Evidence candidate prepare/submit, Finding Draft prepare, Static/Reverse/Memory/Detection/Network handoff prepare, sensitive item inclusion, package export and withdrawal. Atomic rules future.

## 22. Limites et erreurs
File/journal/user artifact/persistence candidate/carved content ≠ Evidence/Finding/original file. Destination unavailable, stale source, hidden contradiction, missing provenance, permission denial or restricted content block/partial the package.

## 23. Métriques
Packages by destination/state, returned/withdrawn, missing provenance, contradictions included, candidates accepted/rejected and automatic qualification/deployment prevented.

## 24. Classification de livraison
`defined` / `planned`; no Detection rule, Network analysis, API, protocol or implementation.

## 25. Critères d’acceptation
**Given** partial carved content and uncertain path **When** Evidence candidate prepared **Then** partiality, uncertainty, source and lineage remain visible and no Evidence is automatically created.

**Given** persistent network configuration without capture **When** future Network handoff prepared **Then** package states no confirmed connection/IOC and no Network capability is treated as complete.

**Given** AI-generated draft **When** reviewed **Then** producer/Tool Calls/sources are visible and analyst can accept, modify or reject; no qualification occurs.

## 26. Questions ouvertes
OPEN-008/013/014/015 remain open; destination contracts and permission model final are future.

## 27. Consommateurs documentaires
Evidence Creation/Review, Finding Management, Static/Reverse/Memory, future Network and Detection phases, Case Replay/Reporting preparation and Objects/Permissions/Journeys/Screens.
