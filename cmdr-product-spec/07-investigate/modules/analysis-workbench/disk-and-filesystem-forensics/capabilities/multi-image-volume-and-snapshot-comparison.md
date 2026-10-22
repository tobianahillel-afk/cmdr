---
id: CAP-INV-377
title: Multi-Image, Volume and Snapshot Comparison
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
  - REQ-AI-002
  - REQ-UX-006
  - REQ-SEC-001
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-377 — Multi-Image, Volume and Snapshot Comparison

## 1. Définition
Comparer plusieurs Disk Images, volumes ou snapshots et leurs structures, filesystems, metadata, files, journals, configurations, user artifacts, startup artifacts and timelines, without merging or modifying sources and without treating `not observed` as absent.

## 2. Problème utilisateur
Sources partielles, incompatibles ou acquises à des moments différents peuvent produire des faux ajouts/suppressions. Une comparaison doit exposer ses préconditions et limites.

## 3. Objectifs
Select sources and verify relationships; compare structures/filesystems/metadata/files/journals/config/user/startup/timeline; distinguish added/deleted/modified/moved/inaccessible/not-observed; show precondition differences and partial results; save comparison; link Hypothesis; prepare Finding Draft.

## 4. Non-objectifs
Ne pas fusionner sources, affirmer absence depuis `not observed`, définir diff algorithm, normalize final schemas, modify images, auto-confirm Finding or perform Network capture comparison.

## 5. Propriétaire
Investigate owns comparison context/result and interpretation; each source retains ownership/version; Shared owns comparison UI mechanism/versioning where applicable.

## 6. Utilisateurs
Principal : DFIR Comparison Analyst. Secondaires : Investigation Lead, Timeline Analyst, Artifact Analyst, Evidence Reviewer.

## 7. Conditions d’entrée
At least two authorized sources; relationships and acquisition times known or explicitly unknown; comparable scopes identified; source limitations and permissions visible.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Images/volumes/snapshots | Investigate | sources/versions | yes | selected versions | blocked |
| Source relationship/acquisition | provenance | temporal/context relation | yes or explicit unknown | visible | ambiguous |
| Structure/filesystem results | CAP-INV-366 | compatibility | yes | per source | incompatible |
| Observation sets | CAP-INV-367..374 | domains to compare | according scope | per source | partial |
| Comparison permissions | Security | cross-source access | yes | current | restricted |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Disk Images/Artifacts | Investigate | source/version/restrictions | read |
| Volumes/Filesystems/Observations | Investigate concepts | comparable projections | read |
| Tool/Tool Call/Run | Studio | producer/version | read |
| Hypothesis/Finding | Investigate | reasoning/draft target | read/link |
| Timeline | Shared/Investigate | temporal comparison | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Disk Comparison Result | create/update/supersede | Investigate concept | source set and preconditions required |
| Difference relation | create/dispute | Investigate | category and uncertainty visible |
| Saved comparison scope | create/version | Investigate/Shared | no source mutation |
| Finding Draft relation | prepare | Investigate | owner review required |

## 11. Fonctionnalités
Select sources; verify relations; compare structures/filesystems/metadata/files/journals/config/user/startup/timelines; classify differences; expose incompatible/partial preconditions; save/filter/annotate; link Hypothesis; prepare Finding Draft. `Not observed` remains separate from absent.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Select/inspect sources | Analyst | comparison scope | 0/2 | access | versioned scope | OPEN-013 for change |
| Run comparison | Analyst | Tool Call/result | 1 | compatible or limits accepted | attributed result | no |
| Annotate/dispute difference | Analyst | relation | 2 | permission | versioned interpretation | OPEN-013 |
| Prepare Finding/handoff | Reviewer | candidate package | 2 | provenance | CAP-INV-379 | no |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| compare structures/content | yes | yes | yes | no need | deterministic comparators |
| classify differences | yes | rules | yes | suggestion | explicit categories |
| summarize significant changes | yes | aggregation | yes | yes | filter/sort/table |
| conclude cause/malice | human | no | no | assistance | Hypothesis/Evidence review |

## 14. États fonctionnels
`draft`, `ready`, `processing`, `available`, `partial`, `incompatible`, `conflicting`, `failed`, `disputed`, `superseded`.

## 15. États d’interface
Loading preserves source order; Empty explains no comparable data; Partial shows missing domains; Error keeps valid differences; Offline read-only; Permission denied masks source differences; Stale shows versions.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Disk Comparison Result | result | Workbench/Hypothesis | sources/preconditions/limits visible |
| Difference relation | relation | CAP-INV-373/374/379 | category and uncertainty retained |
| Finding Draft selection | candidate | CAP-INV-379/109 | no automatic confirmation |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Sessions/images/volumes | compare | CAP-INV-377 | sources, relations, scope, limits | source context |
| Difference | correlate | CAP-INV-373/374 | category, source/time, supporting/contradicting data | Comparison |
| Comparison | handoff | CAP-INV-379 | selected differences, preconditions, uncertainty | Comparison |

## 18. Dépendances
CAP-INV-364/366..374/378/379, Shared Comparison/Versioning/Trace, Studio Tools, Security, OPEN-005/008/013/015.

## 19. Source de vérité
Each source remains immutable truth for its recorded content; Investigate owns comparison interpretation. No source is merged or silently normalized.

## 20. Provenance et audit
Source IDs/versions/relations/acquisition times, structures/profiles, Tool/version/parameters, compared domains, filters, difference categories, partial/incompatible states, annotations and handoff.

## 21. Permissions fonctionnelles
Comparison create/read/export, multi-image/raw source read, sensitive difference view, annotate/dispute, Hypothesis/Finding Draft prepare and cross-tenant restrictions.

## 22. Limites et erreurs
Not observed ≠ absent; partial/incompatible sources may prevent conclusions. Missing source, different scope/profile/timezone, Tool mismatch, permission denial and large-volume limits remain explicit.

## 23. Métriques
Comparisons by state, incompatible/partial domains, difference categories, disputed results, `not observed` usage and handoffs.

## 24. Classification de livraison
`defined` / `planned`; no diff algorithm, source merge or implementation.

## 25. Critères d’acceptation
**Given** one partial and one complete image **When** an entry is not observed in the partial source **Then** result is `not observed`, not absent, and limitations remain visible.

**Given** incompatible filesystems **When** compared **Then** incompatible domains are excluded/marked and no normalized equality is invented.

**Given** no AI **When** comparison runs **Then** deterministic comparators, categories, filters and human review work.

## 26. Questions ouvertes
OPEN-005/008/013/015 remain open; comparison objects/algorithms and permissions final are future.

## 27. Consommateurs documentaires
INV-DSK-001, CAP-INV-373/374/378/379, Hypothesis/Finding, Shared Comparison and Objects/Permissions/Technique.
