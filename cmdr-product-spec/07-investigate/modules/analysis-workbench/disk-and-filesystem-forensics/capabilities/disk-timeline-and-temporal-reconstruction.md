---
id: CAP-INV-374
title: Disk Timeline and Temporal Reconstruction
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
  - REQ-UX-007
  - REQ-SEC-001
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-374 — Disk Timeline and Temporal Reconstruction

## 1. Définition
Agréger et corréler les informations temporelles Disk/Filesystem en distinguant timestamp observé, enregistré, reconstruit, estimé ou absent, sans fabriquer de temps ni remplacer Memory Timeline ou Case Timeline.

## 2. Problème utilisateur
Des timestamps de metadata, journals, applications ou configurations peuvent diverger, manquer de timezone ou représenter des sémantiques différentes. Une fusion naïve crée une chronologie artificiellement certaine.

## 3. Objectifs
Afficher source, filesystem, file/Artifact, journal, timezone and timestamp quality; gaps/conflicts; filter/group/compare; correlate Memory Timeline, Case Timeline and Endpoint telemetry; compare images; annotate; prepare Evidence candidate.

## 4. Non-objectifs
Ne pas fabriquer timestamps, définir moteur/clock algorithm, remplacer timelines propriétaires, affirmer user action, créer Evidence automatiquement ou couvrir packet/network timeline.

## 5. Propriétaire
Investigate owns Disk Timeline projection and temporal interpretation; Shared owns Timeline mechanism; Case/Memory timelines remain owner projections.

## 6. Utilisateurs
Principal : Timeline Analyst. Secondaires : Disk Analyst, Memory Analyst, Investigation Lead, Evidence Reviewer and Audit Analyst.

## 7. Conditions d’entrée
Sourced observations with timestamp value/type/quality/timezone or explicit absence; image/session and limitations visible; permission across correlated sources.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Disk/Filesystem observations | CAP-INV-367..373 | events/records/relations | yes | source versions | empty/partial |
| Timestamp quality/timezone | source/analyst | observed/recorded/reconstructed/estimated/absent | yes | per entry | ambiguous |
| Memory/Case/Endpoint events | owners | correlation context | no | freshness visible | disk-only |
| Comparison scope | analyst/CAP-INV-377 | images/period/filters | no | versioned | default scope |
| Permissions | Security | source-level access | yes | current | restricted |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Disk observations/Journal Records | Investigate | source/time quality | read |
| Disk Image/File/Artifact | Investigate | source relation | read/link |
| Memory Timeline | Investigate projection | correlation | read/link |
| Case Timeline/Timeline Entry | Investigate/Shared | correlation/order | read/link |
| Endpoint telemetry | owner | external projection | read if authorized |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Disk Timeline Entry/projection | create/annotate/dispute | Investigate concept | source and time quality required |
| Temporal correlation relation | create/supersede | Investigate | no invented timestamp |
| Evidence candidate selection | prepare | Investigate | qualification future |
| Trace event | emit | Shared | filters/correlation/handoff attributed |

## 11. Fonctionnalités
Aggregate timestamps; distinguish quality/type/absence; show timezone/source/filesystem/file/journal; gaps/conflicts; filter/group/compare; correlate Memory/Case/Endpoint; compare images; annotate; prepare Evidence. Metadata timestamp remains not a certain event.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| View/filter/group | Analyst | Disk Timeline | 0 | read | sourced view | no |
| Compare/correlate | Analyst | temporal relation | 0/1 | sources authorized | gaps/conflicts | no |
| Annotate/dispute | Analyst | entry/relation | 2 | permission | versioned interpretation | OPEN-013 |
| Prepare Evidence/handoff | Reviewer | selection | 2 | provenance | CAP-INV-379 | no |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| order entries | yes | yes | yes | no need | deterministic Timeline |
| label timestamp quality | yes | rules | yes | suggestion | explicit labels/checklist |
| propose correlations | yes | matching rules | yes | suggestion | filters/comparator |
| conclude sequence/intent | human review | no | no | assistance | source-by-source review |

## 14. États fonctionnels
`queued`, `processing`, `available`, `partial`, `conflicting`, `gapped`, `incompatible`, `failed`, `disputed`, `superseded`.

## 15. États d’interface
Loading preserves filters; Empty means no available temporal data; Partial shows missing sources; Error keeps valid entries; Offline read-only; Permission denied masks source details; Stale shows source versions.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Disk Timeline | projection | Workbench/Case | source/type/quality/timezone visible |
| Temporal correlation | relation | Hypothesis/Memory/Case | gaps/conflicts and uncertainty retained |
| Evidence selection | candidate | CAP-INV-379 | no timestamp invented or qualification |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Journal/System/User/Persistence observation | add/correlate | CAP-INV-374 | source, timestamp quality, relations, gaps | source view |
| Disk Timeline | correlate | Memory/Case Timeline | selected entries, quality, conflicts, return origin | Disk Timeline |
| Timeline selection | handoff | CAP-INV-379 | entries, sources, uncertainty | Disk Timeline |

## 18. Dépendances
CAP-INV-359/367..373/377/378/379, Shared Timeline/Linking, Case Timeline, Endpoint telemetry, OPEN-005/008/013/015.

## 19. Source de vérité
Underlying artifacts/records remain sources; Investigate owns Disk temporal interpretation; Shared owns ordering mechanism; owner timelines remain distinct.

## 20. Provenance et audit
Image/session/source artifact, Tool/version, timestamp value/type/quality/timezone, transforms, filters, correlation sources, gaps/conflicts, analyst annotation and handoff.

## 21. Permissions fonctionnelles
Disk timeline read, cross-source correlate, comparison, annotation/dispute, Evidence candidate and provenance export. Source-level permissions and final step-up/SoD future.

## 22. Limites et erreurs
Disk Timeline ≠ Memory/Case Timeline; metadata/journal record ≠ certain action; absent timestamp never fabricated. Clock/timezone conflict, partial image, inaccessible source and permission denial remain visible.

## 23. Métriques
Entries by timestamp quality, absent/estimated/conflicting/gapped, correlations accepted/disputed and handoffs with complete provenance.

## 24. Classification de livraison
`defined` / `planned`; no implementation, algorithm, engine or packet timeline.

## 25. Critères d’acceptation
**Given** an observation without timestamp **When** added **Then** `absent` remains visible and no time is invented.

**Given** Disk and Memory timestamps conflict **When** correlated **Then** both values, sources, quality and conflict remain visible without replacement.

**Given** no AI **When** timeline built **Then** deterministic ordering, filters, comparison and annotations work.

## 26. Questions ouvertes
OPEN-005/008/013/015 remain open; temporal objects/contracts and permissions final are future.

## 27. Consommateurs documentaires
INV-DSK-001, Case/Memory Timeline, CAP-INV-370..379, Evidence/Replay and Objects/Permissions/Technique.
