---
id: CAP-INV-369
title: Deleted and Unallocated Data Analysis
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
  - REQ-OBJ-003
  - REQ-OBJ-004
  - REQ-SEC-001
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-369 — Deleted and Unallocated Data Analysis

## 1. Définition
Analyser entrées supprimées candidates, metadata résiduelles, données non allouées et espaces résiduels en exposant qualité, conflits et attribution incertaine, sans confondre entrée, contenu récupéré ou résultat de carving.

## 2. Problème utilisateur
Des fragments peuvent être attribués à tort à un fichier, un path ou un utilisateur. Une entrée supprimée ne garantit ni contenu, ni complétude, ni récupération.

## 3. Objectifs
Voir deleted candidates, residual metadata, unallocated/residual areas, possible relations, reconstruction quality, missing parts and conflicts; compare methods; annotate; select for carving/extraction; link Hypothesis; prepare Evidence candidate.

## 4. Non-objectifs
Ne pas garantir attribution/path/name/completeness, définir carving algorithm, modifier source, fournir récupération offensive, créer Evidence automatiquement ou analyser full Network Forensics.

## 5. Propriétaire
Investigate possède observations, selections et Derived Artifact relations; Studio Tools; Shared Jobs/Trace; Evidence owner retains qualification.

## 6. Utilisateurs
Principal : Disk Forensics Analyst. Secondaires : Recovery Analyst, Evidence Reviewer, Privacy Reviewer et Investigation Lead.

## 7. Conditions d’entrée
Image/session/filesystem disponibles; partiality visible; deleted/unallocated projection attributed to Tool/version; permissions for recovered/private data.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Deleted entry candidates | filesystem Tool Call | residual metadata/state | non | same image/version | empty |
| Unallocated/residual sources | structure result | source ranges/projections | oui pour analyse | immutable image | blocked |
| Reconstruction results | Tool Calls | content fragments/quality/errors | non | attributed | unavailable |
| Disk limitations | CAP-INV-365 | missing ranges | oui | inherited | partial |
| Permissions/purpose | Security/Case | sensitive recovery scope | oui | current | restricted |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Disk Image/Filesystem Entry | Investigate | source/metadata | lire |
| Deleted Entry Candidate / Unallocated Region | Investigate concepts | state/quality/source | lire |
| Tool/Tool Call/Job | Studio/Shared | producer/progress/errors | lire |
| Hypothesis/Evidence | Investigate | context/qualification target | lire/lier |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Deleted/Unallocated Observation | créer/annoter/dispute | Investigate concept | uncertainty mandatory |
| Recovery selection | créer/update | Investigate | bounded source required |
| Derived Artifact relation | créer via CAP-INV-311/375 | Investigate | parent/source/quality required |
| Evidence candidate relation | préparer | Investigate | no auto-qualification |

## 11. Fonctionnalités
Display deleted candidates, residual metadata, unallocated/residual sources, possible relations, quality, missing parts and conflicts; compare methods; annotate; select for carving; extract Derived Artifact; link Hypothesis; prepare Evidence candidate. Keep deleted entry, recovered content, partial content, unallocated data and carving result distinct.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter/compare | Analyst | candidates/sources | 0/1 | read | quality/conflicts visible | non |
| Annoter/link | Analyst | observation | 2 | session modifiable | version relation | OPEN-013 |
| Sélectionner pour carving | Analyst | bounded source | 2 | scope/permission | recovery plan | OPEN-013 |
| Extraire Derived Artifact | Analyst | result | 1 | authorized Tool | partiality/lineage | non |
| Préparer Evidence | Reviewer | candidate | 2 | provenance | package | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| détecter candidates | oui | Tool/rules | oui | suggestion | deterministic scans |
| comparer méthodes | oui | oui | oui | summary | comparison table |
| proposer relation | oui | rules | oui | yes | metadata/source review |
| confirmer attribution | humain | non | non | assistance | evidence review |

## 14. États fonctionnels
`candidate`, `processing`, `partial`, `available`, `inconsistent`, `unattributed`, `restricted`, `invalid`, `superseded`.

## 15. États d’interface
Loading preserves source; Empty means no candidate, not absence of deleted data; Partial exposes missing ranges; Error keeps valid fragments; Offline read-only; Permission denied masks content; Stale shows source version.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Deleted/Unallocated observation | observation | Workbench/Hypothesis | state/quality/uncertainty visible |
| Recovery selection | plan | CAP-INV-375 | bounded source and restrictions |
| Derived Artifact | Artifact relation | Static/Reverse/CAP-INV-379 | partiality and lineage preserved |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Deleted Entry | analyser | CAP-INV-369 | entry, residual metadata, image limits | Navigation |
| Unallocated source | carve | CAP-INV-375 | bounded source, objective, permissions | Analysis |
| Observation/Artifact | handoff | CAP-INV-379 | quality, conflicts, source, uncertainty | Analysis |

## 18. Dépendances
CAP-INV-311/365/367/368/375/379, Studio Tools, Shared Jobs/Trace, Security, OPEN-005/008/013/014/015.

## 19. Source de vérité
Disk Image and Tool results remain sources; Investigate owns interpretation and relations; no inferred path/name becomes source truth.

## 20. Provenance et audit
Image/session/filesystem/source region, entry metadata, Tool/version/Call, method comparison, quality, partial ranges, analyst, extraction/handoff and access events.

## 21. Permissions fonctionnelles
Deleted-data read, unallocated-data read, recovery prepare, Derived Artifact create/read/export and Evidence candidate prepare; minimization, step-up and SoD future.

## 22. Limites et erreurs
Deleted entry ≠ recovered content; unallocated data ≠ identified file; residual space ≠ attributed content. Errors include overlap, missing ranges, corruption, false relation, restricted data and Tool failure.

## 23. Métriques
Candidates by quality/state, unattributed/partial results, conflicting methods, recoveries with lineage and automatic qualifications prevented.

## 24. Classification de livraison
`defined` / `planned`; no carving algorithm, engine, format or implementation chosen.

## 25. Critères d’acceptation
**Given** a deleted candidate with partial content and no confirmed full path **When** reviewed **Then** partial status and missing path remain explicit, entry and content remain distinct and no Evidence is created.

**Given** unallocated data with conflicting reconstructions **When** compared **Then** all methods, sources and conflicts remain visible without certain attribution.

**Given** no AI **When** analysis runs **Then** deterministic Tools, tables, comparisons and human review cover the workflow.

## 26. Questions ouvertes
OPEN-005/008/013/014/015 remain open; algorithms, objects and permissions final are future.

## 27. Consommateurs documentaires
INV-DSK-001, CAP-INV-375/379, Static/Reverse, Hypothesis/Evidence and Objects/Permissions/Technique phases.
