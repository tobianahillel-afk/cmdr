---
id: CAP-INV-370
title: Filesystem Journal and Change History Analysis
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
# CAP-INV-370 — Filesystem Journal and Change History Analysis

## 1. Définition
Examiner les journaux et historiques filesystem disponibles, leurs records, timestamps, relations, gaps et conflicts, sans présenter un record comme action utilisateur certaine ni remplacer Disk ou Case Timeline.

## 2. Problème utilisateur
Un journal peut être partiel, réordonné, sans auteur ou lié à un compte sans prouver une interaction humaine. Une lecture directe peut produire de fausses attributions.

## 3. Objectifs
Voir journals/history, records, timestamps, change types, associated entries/identities/sequences, gaps/conflicts/partials; filter/compare/correlate/annotate; link Disk/Case Timeline; prepare Evidence candidate.

## 4. Non-objectifs
Ne pas définir structures internes, journal formats, parsing algorithms, authorship certainty, Case Timeline, engine, command, API or full Network Forensics.

## 5. Propriétaire
Investigate owns journal interpretation; Disk Image/Tool result remain source; Shared owns Timeline mechanism; Case Timeline remains Case-owned.

## 6. Utilisateurs
Principal : Timeline/Filesystem Analyst. Secondaires : Investigation Lead, Evidence Reviewer, User Activity Analyst and Audit Analyst.

## 7. Conditions d’entrée
Selected filesystem; journal projection attributed to Tool/version; timestamp quality/timezone visible; partiality and permissions inherited.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Journal/change-history result | Tool Call | records/types/sequences | yes | source version | unavailable |
| File/entry relations | CAP-INV-367/368 | associated entries | no | same image | unlinked |
| Timestamp quality/timezone | source/analyst | temporal context | yes | per record | ambiguous |
| Image/filesystem limitations | CAP-INV-365/366 | gaps/partiality | yes | inherited | partial |
| Case/Disk timeline context | owners | correlation targets | no | current selection | local-only |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Disk Image/Filesystem | Investigate | source/limits | read |
| Journal Record | Investigate concept | timestamp/type/relation | read |
| Filesystem Entry/File Observation | Investigate concepts | associated context | read/link |
| Timeline mechanism/Case Timeline | Shared/Investigate | ordering/correlation | consume/read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Journal observation | create/annotate/dispute | Investigate | record ≠ user action |
| Change-history relation | create/supersede | Investigate | source/sequence/uncertainty required |
| Disk Timeline selection | prepare/link | Investigate | CAP-INV-374 owns Disk Timeline projection |
| Evidence candidate | prepare | Investigate | qualification remains future |

## 11. Fonctionnalités
Display records/timestamps/change types/associated files/identities/sequences/gaps/conflicts/partial records; filter, compare, correlate, annotate; link Disk Timeline and Case Timeline; prepare Evidence candidate with uncertainty.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspect/filter | Analyst | journal records | 0 | read | sourced view | no |
| Compare/correlate | Analyst | history relation | 0/1 | sources readable | gaps/conflicts | no |
| Annotate/dispute | Analyst | observation | 2 | permission | versioned interpretation | OPEN-013 |
| Link timeline/prepare Evidence | Analyst | selection | 2 | provenance | package | no |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| parse/order records | yes | yes | yes | no need | deterministic parser/table |
| relate files/records | yes | rules | yes | suggestion | IDs/timestamps comparison |
| explain gaps/conflicts | yes | yes | yes | summary | explicit diagnostics |
| infer human action | human only | no | no | assistance | corroboration review |

## 14. États fonctionnels
`queued`, `processing`, `available`, `partial`, `conflicting`, `unlinked`, `failed`, `disputed`, `superseded`.

## 15. États d’interface
Loading preserves filters; Empty means no available record, not no activity; Partial exposes gaps; Error keeps valid records; Offline read-only; Permission denied masks protected data; Stale shows source version.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Journal observation set | observations | analyst/Case | source/quality/gaps visible |
| Change-history relation | relation | CAP-INV-374 | no certain authorship |
| Evidence candidate selection | package | CAP-INV-379/107/108 | uncertainty and provenance |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Filesystem | open journal | CAP-INV-370 | image, filesystem, Tool result, limits | Navigation |
| Journal Record | correlate | CAP-INV-374 | timestamp quality, file relation, gaps | Journal |
| Journal observation | handoff | CAP-INV-379 | records, contradictions, provenance | Journal |

## 18. Dépendances
CAP-INV-365/366/367/368/374/379, Shared Timeline/Linking, Case Timeline, Studio Tool Calls, OPEN-005/008/013/015.

## 19. Source de vérité
Raw records stay sourced to Disk Image/Tool result; Investigate owns interpretation; Shared orders but does not own meaning; Case Timeline is not replaced.

## 20. Provenance et audit
Image/session/filesystem/journal source, Tool/version, record as displayed, timestamps/quality/timezone, relationships, filters, annotations, disputes and handoffs.

## 21. Permissions fonctionnelles
Journal read, sensitive metadata read, compare/correlate, annotate/dispute, timeline link and Evidence candidate prepare; final access model future.

## 22. Limites et erreurs
Record ≠ certain user action; account ≠ author; journal ≠ Disk/Case Timeline. Missing/corrupt records, clock issues, unsupported journal, relation conflict and permission denial remain explicit.

## 23. Métriques
Records by quality, gaps/conflicts, unlinked records, disputed interpretations and handoffs with complete source.

## 24. Classification de livraison
`defined` / `planned`; no format/parser implementation or engine choice.

## 25. Critères d’acceptation
**Given** a journal record with associated account but no human interaction proof **When** correlated **Then** account is context, uncertainty is visible and no human action is asserted.

**Given** partial journal with gaps **When** added to Disk Timeline **Then** gaps and source quality remain visible and no missing event is invented.

**Given** no AI **When** analyzed **Then** deterministic parser, filters, comparison and human review work.

## 26. Questions ouvertes
OPEN-005/008/013/015 remain open; journal support, objects and permissions final are future.

## 27. Consommateurs documentaires
INV-DSK-001, CAP-INV-374/379, Case Timeline/Evidence/Replay and Objects/Permissions/Technique phases.
