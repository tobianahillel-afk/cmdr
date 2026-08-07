---
id: CAP-INV-706
title: Mobile Filesystem, Partition and Storage Analysis
product: investigate
module: mobile-forensics
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-07
requirement_ids: [REQ-INV-001, REQ-PROD-014, REQ-PROD-020, REQ-OBJ-003, REQ-SEC-001]
open_decisions: [OPEN-005, OPEN-011, OPEN-013, OPEN-014, OPEN-015]
source-of-truth: canonical
---
# CAP-INV-706 — Mobile Filesystem, Partition and Storage Analysis

## 1. Définition
Analyser en lecture seule les partitions candidates, volumes, filesystems, storage areas, directories, entries, links, metadata, permissions disponibles, timestamps, protected/shared/application areas, removable-storage candidates et zones residual/unallocated lorsqu’elles sont réellement représentées.

## 2. Problème utilisateur
Une extraction Mobile peut ne représenter qu’une partie du stockage. Sans modèle fonctionnel explicite, un analyste peut traiter une structure partielle comme complète, confondre metadata et contenu, ou réutiliser des hypothèses Disk Forensics non supportées par la source.

## 3. Objectifs
- parcourir et rechercher la structure représentée sans liste finale de filesystems supportés;
- afficher source, path/identifier, type, size/metadata, timestamps, permission/accessibility et zone;
- conserver protected/missing/deleted/residual states et liens vers package/acquisition;
- permettre filtres, bookmarks, comparaison et handoffs vers app/deleted/artifact analysis.

## 4. Non-objectifs
Aucune acquisition, mount réel, écriture source, filesystem support matrix finale, carving algorithm, bypass protected area, exploit, command, parser implementation, format ou Disk Forensics complète.

## 5. Propriétaire
Investigate possède Mobile Filesystem/Storage Observations. Collection/source owns acquisition representation; generic Artifact/Evidence remain canonical Investigate objects; Shared owns search/version mechanisms.

## 6. Utilisateurs
Principal : Mobile Forensics Analyst. Secondaires : DFIR Analyst, Evidence Reviewer et Investigation Lead.

## 7. Conditions d’entrée
Session, device/platform scope, accepted or explicit-limited acquisition context, integrity/completeness/accessibility assessments and permission to represented storage.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Filesystem/storage representation | package/extraction | analytical source | oui | selected version | no storage analysis |
| Device/platform context | CAP-INV-703 | interpretation context | oui | Session version | generic/unknown labels only |
| Integrity/completeness/accessibility | CAP-INV-705 | limitations | oui | reviewed version | partial/unverified banner |
| Entry metadata/timestamps | source/parser output | observation fields | selon source | source version | field unknown |
| Permissions/classification | Security/source | content access | oui | access time | metadata-only/denied |
| Previous observations/bookmarks | Session | analyst context | non | versioned | none invented |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Mobile Session / Device Context | Investigate | scope/source/platform | read |
| Package / filesystem extraction | source | represented structure | read only |
| Integrity/Completeness Assessment | Investigate | limitations | read |
| Artifact | Investigate | referenced file/fragment relation | read/link |
| Shared Search/Versioning | Shared | mechanisms | consume |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Mobile Filesystem Observation | create/review/dispute/supersede | Investigate concept | source path/region/version required |
| Storage/partition relation | create/dispute | Investigate | candidate if source is ambiguous |
| Bookmark/annotation | create/update | Investigate/Shared | source-scoped |
| Derived Artifact proposal | prepare | CAP-INV-718 | no Evidence creation |

## 11. Fonctionnalités
Browse trees/tables, search/filter by path/type/time/area/accessibility, inspect links/metadata/permissions/timestamps, compare versions, bookmark, identify app/shared/protected/residual areas, show source gaps and open app/deleted analysis or bounded extraction handoff.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Browse/search/filter | analyste | represented entries | 0 | read | scoped results | non |
| Compare metadata/versions | analyste | entries | 1 | both accessible | diff | non |
| Bookmark/annotate/dispute | analyste | observation | 2 | source | versioned note | OPEN-013 |
| Prepare Derived Artifact | analyste | selected entry | 1/2 | source permission | CAP-INV-718 proposal | non |
| Open application/deleted analysis | analyste | relation | 0 | relevant source | contextual transition | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Parse represented structure | oui | oui | oui | non | deterministic parser |
| Search/filter/compare | oui | oui | oui | non necessary | query/table/diff |
| Group directories/areas | oui | rules | oui | suggestion | tree/tags |
| Flag unusual metadata | oui | rules | oui | candidate | explicit filters |
| Infer unsupported missing filesystem | non | non | non | interdit | mark unknown/missing |

## 14. États fonctionnels
Observations: `proposed`, `under-review`, `supported`, `weakly-supported`, `contradicted`, `inconclusive`, `disputed`, `superseded`, `withdrawn`. Source regions may be `available`, `partial`, `restricted`, `inaccessible`.

## 15. États d’interface
Loading keeps tree path; Empty distinguishes no entries from unavailable region; Partial marks missing/protected regions; Error preserves parsed entries; Offline read-only; Permission denied shows metadata only if allowed; Stale identifies source version.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Filesystem/Storage Observations | concepts | CAP-INV-707/708/710/715/716/717 | source/version/region/limitations |
| Bookmarks/annotations | session data | analyst/reviewer | scoped and versioned |
| Candidate app/shared/protected areas | relations | CAP-INV-707/708 | candidates, not platform truth |
| Derived Artifact proposal | package input | CAP-INV-718 | bounded source selection |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-703/705 | represented storage ready | CAP-INV-706 | Session, device, package, limitations | Session |
| CAP-INV-706 | app/package area identified | CAP-INV-707/708 | paths, entries, metadata, source | Filesystem |
| CAP-INV-706 | deleted/residual candidate | CAP-INV-715 | source region, metadata, quality | Filesystem |
| CAP-INV-706 | select record/file | CAP-INV-718 | entry/version/path, restrictions | Filesystem |
| Observations | temporal review | CAP-INV-716 | timestamps/source quality | Filesystem |

## 18. Dépendances
CAP-INV-703/705/707/708/710/715/716/718/719, Disk/Artifact concepts, Shared Search/Versioning/Inspector, OPEN-005/011/013/014/015.

## 19. Source de vérité
Raw structure remains the package/extraction source. Mobile observations are Investigate analytical concepts. Artifact remains canonical Investigate object. No observation rewrites source filesystem metadata.

## 20. Provenance et audit
Package/extraction version, partition/volume/region/path, parser/Tool/Run, fields/timestamps, access state, queries/filters, bookmarks, annotations, comparisons, derived selections, errors and human disposition.

## 21. Permissions fonctionnelles
Filesystem read, protected/restricted area metadata/content read separately, search/compare, bookmark/annotation, deleted/residual read, Derived Artifact prepare/export and provenance read. No source write permission.

## 22. Limites et erreurs
Unknown filesystem, malformed entry, missing region, unsupported metadata, encrypted/protected area, broken link, timestamp ambiguity, parser failure or permission denial are explicit and do not imply absence or corruption of the entire source.

## 23. Métriques
Regions represented/missing/restricted, parsed entries, parser failures, searches, bookmarks, derived selections, disputed observations and comparisons across versions.

## 24. Classification de livraison
`defined` / `planned`. No filesystem list, parser implementation, carving engine, mount mechanism, format, API or product code delivered.

## 25. Critères d’acceptation
**Given** a filesystem extraction missing protected areas **When** browsed **Then** represented entries remain usable while completeness and inaccessible regions stay explicit.

**Given** a recovered-looking residual entry **When** selected **Then** it remains a candidate and is handed to CAP-INV-715 rather than treated as certain original content.

**Given** no AI **When** storage is analyzed **Then** tree/table/search/filter/diff and deterministic parsers provide full workflow.

## 26. Questions ouvertes
OPEN-005 future engines/parsers; OPEN-011 platform/filesystem delivery; OPEN-013 annotations; OPEN-014 final Artifact/material relations; OPEN-015 Tool/Run contracts.

## 27. Consommateurs documentaires
Application inventory/data, media/documents, deleted/recovery, timeline, artifact handoff, Disk Forensics boundaries, Objects, Permissions, Screens, Quality and Technique.
