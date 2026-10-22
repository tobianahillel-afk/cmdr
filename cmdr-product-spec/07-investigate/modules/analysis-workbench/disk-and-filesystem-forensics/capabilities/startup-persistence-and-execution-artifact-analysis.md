---
id: CAP-INV-373
title: Startup, Persistence and Execution Artifact Analysis
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
  - REQ-INV-006
  - REQ-AI-002
  - REQ-SEC-001
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-373 — Startup, Persistence and Execution Artifact Analysis

## 1. Définition
Regrouper et évaluer défensivement les artefacts persistants candidats liés au démarrage, à l’exécution et à la persistance, avec éléments pour/contre et corrélations, sans fournir de procédure offensive ni confirmation automatique.

## 2. Problème utilisateur
Un mécanisme de démarrage légitime, ancien, désactivé ou incomplet peut être qualifié à tort de persistance malveillante. Les contradictions et la source doivent rester visibles.

## 3. Objectifs
Group startup/execution/persistence candidates; show associated files/configurations/identities/timestamps/sources and supporting/contradicting elements; compare images; correlate Memory/Dynamic; annotate; prepare Hypothesis, Evidence candidate and future Detection package.

## 4. Non-objectifs
No persistence creation procedure, command, code, mechanism implementation, bypass/evasion, automatic confirmation, Finding confirmation or detection rule deployment.

## 5. Propriétaire
Investigate owns Persistence Candidate interpretation; Static/Reverse/Memory/Dynamic retain their contexts; Studio Tools; Govern real-target actions; future Detection owner receives package only.

## 6. Utilisateurs
Principal : DFIR/Persistence Analyst. Secondaires : Malware Analyst, Memory Analyst, Investigation Lead, Evidence Reviewer and Detection Engineer future consumer.

## 7. Conditions d’entrée
System/user/file observations available; source/time/limits visible; permissions; no content execution; related Memory/Dynamic results optional and sourced.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Startup/execution artifacts | CAP-INV-368/371/372 | candidate sources | yes | image version | unavailable |
| Associated files/config/identities | Disk observations | context | no | same/correlated source | partial |
| Timeline/journal | CAP-INV-370/374 | temporal evidence | no | quality visible | uncorrelated |
| Memory/Dynamic observations | CAP-INV-347..362/314..328 | corroboration/contradiction | no | source/version visible | disk-only |
| Rules/Tool results | Studio | candidate grouping | no | attributed | manual review |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| File/System/User Observation | Investigate | sources/context | read |
| Memory/Dynamic observation | Investigate | corroboration only | read/link |
| Hypothesis/Evidence/Finding | Investigate | reasoning/handoff target | read/link |
| Tool/Tool Call/Run | Studio | producer/version | read |
| Timeline | Shared/Investigate | temporal context | read/link |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Persistence Candidate | create/classify/dispute/supersede/withdraw | Investigate concept | candidate until owner review |
| Supporting/contradicting relation | create/link | Investigate | source and uncertainty required |
| Hypothesis/Evidence candidate | prepare | Investigate | no automatic qualification |
| Detection knowledge package | prepare | future owner | no rule creation/deployment |

## 11. Fonctionnalités
Group candidate artifacts; show mechanisms as functional categories, associated files/config/identities/timestamps/sources, for/against; compare images; correlate Memory/Dynamic; annotate/classify/dispute; prepare Hypothesis/Evidence/Future Detection package. No offensive instructions.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspect/group/compare | Analyst | candidates | 0/1 | sources read | evidence matrix | no |
| Annotate/classify/dispute | Analyst | Persistence Candidate | 2 | permission | versioned status | OPEN-013 |
| Link Memory/Dynamic/Hypothesis | Analyst | relation | 2 | objects readable | sourced link | OPEN-013 |
| Prepare handoff/package | Reviewer | candidate package | 2 | provenance | CAP-INV-379/future Detection | no |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| group candidate artifacts | yes | rules/Tools | yes | suggestion | explicit rule sets |
| compare/correlate | yes | yes | yes | summary | matrices/timeline |
| propose Hypothesis | yes | no | yes as draft | yes | analyst-authored Hypothesis |
| confirm persistence/Finding/rule | human owner | no | prohibited | prohibited | Evidence/Finding/Detection review |

## 14. États fonctionnels
`candidate`, `weakly-supported`, `supported`, `contradicted`, `inconclusive`, `disputed`, `superseded`, `withdrawn`.

## 15. États d’interface
Loading preserves candidate; Empty means none observed, not none exists; Partial shows missing sources; Error keeps evidence matrix; Offline read-only; Permission denied masks protected artifacts; Stale shows source age.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Persistence Candidate | observation | Case/Hypothesis | status/elements for/against visible |
| Correlation/comparison | relation/result | Memory/Dynamic/Timeline | source distinction preserved |
| Evidence/Detection package | candidate | CAP-INV-379/future 4B.3 | no confirmation or deployment |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Startup/system/user artifact | evaluate | CAP-INV-373 | source, file/config, identity, timestamps | source view |
| Candidate | correlate | Memory/Dynamic | candidate, source, time, contradictions | Persistence view |
| Candidate | handoff | CAP-INV-379 | Hypothesis, selected observations, uncertainty | Persistence view |

## 18. Dépendances
CAP-INV-368/370/371/372/374/377/379, CAP-INV-314..328/347..362, Studio Tools, Shared Timeline/Linking, future 4B.3, OPEN-005/008/013/015.

## 19. Source de vérité
Artifacts and observations remain sources; Investigate owns candidate interpretation; future Detection owner creates/evaluates rules; no candidate becomes Finding automatically.

## 20. Provenance et audit
Image/session/source artifacts, Tool/version, grouping rules, timestamps, related files/config/identity, Memory/Dynamic links, evidence for/against, analyst status, handoff and disposition.

## 21. Permissions fonctionnelles
Persistence candidate read/review/classify, system/user artifact read, cross-analysis link, Hypothesis/Evidence/Detection handoff prepare. Final step-up/SoD future.

## 22. Limites et erreurs
Startup artifact ≠ confirmed persistence; candidate ≠ Finding. Missing sources, stale config, disabled mechanism, conflicting images, permission denial and Tool failure remain explicit. No procedure or command.

## 23. Métriques
Candidates by state, contradictions, cross-source correlations, withdrawn/superseded, Evidence/Detection packages and automatic confirmations prevented.

## 24. Classification de livraison
`defined` / `planned`; no offensive method, engine or implementation.

## 25. Critères d’acceptation
**Given** startup artifact, associated file and contradictory evidence **When** Hypothesis prepared **Then** candidate status and elements for/against remain visible, no persistence procedure is provided and no Finding is confirmed.

**Given** Memory contradicts Disk candidate **When** correlated **Then** both sources/times remain distinct and status may become contradicted/inconclusive.

**Given** no AI **When** analyzed **Then** explicit rules, tables, comparisons and human review work.

## 26. Questions ouvertes
OPEN-005/008/013/015 remain open; objects, confidence and Detection handoff contracts are future.

## 27. Consommateurs documentaires
INV-DSK-001, Case/Hypothesis/Evidence, Memory/Dynamic, CAP-INV-374/377/379, future Detection Engineering and Objects/Permissions/Technique.
