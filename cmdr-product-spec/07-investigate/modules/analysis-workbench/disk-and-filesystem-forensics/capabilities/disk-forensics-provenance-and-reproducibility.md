---
id: CAP-INV-378
title: Disk Forensics Provenance and Reproducibility
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
  - REQ-OBJ-009
  - REQ-AI-002
  - REQ-SEC-001
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-378 — Disk Forensics Provenance and Reproducibility

## 1. Définition
Retracer le Case, la source, l’acquisition, la Disk Image, la session, les structures, Tools, paramètres, observations, Derived Artifacts, accès restreints, erreurs et décisions, puis évaluer la reproductibilité sans redéfinir Trace, Activity ou Audit.

## 2. Problème utilisateur
Une conclusion ou un Artifact dérivé perd sa valeur si l’image, le filesystem, la version du Tool, les paramètres ou les limitations ne peuvent être retrouvés.

## 3. Objectifs
Trace Case, Endpoint/source, Collection Request/Job, Disk Image, acquisition/custody, session, partitions/volumes/filesystems, Tools/versions/Calls/Runs, parameters/queries/filters, observations, Derived Artifacts, restricted access, errors/interruptions, human decisions and final dispositions; classify reproducibility.

## 4. Non-objectifs
Ne pas définir infrastructure, manifests, hashes/PKI, storage, replay engine, API, protocol, schema final, audit service or automatic conclusion.

## 5. Propriétaire
Investigate owns forensic provenance assessment and local relations; Shared owns Trace/Activity/Audit mechanisms; Studio owns Tool/Call/Run; Settings administrative projections; source object owners remain unchanged.

## 6. Utilisateurs
Principal : Audit/Forensics Reviewer. Secondaires : Disk Analyst, Investigation Lead, Evidence Reviewer, Security Reviewer and Release/Quality reviewer.

## 7. Conditions d’entrée
Session or result exists; source/version and Tool records accessible or explicitly missing; permissions for provenance and restricted access events; no trace deletion.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Case/source/acquisition/custody | owner capabilities | provenance chain | yes | versions visible | not-reproducible/disputed |
| Disk session and structures | CAP-INV-364/366 | analytical context | yes | session version | missing-volume/filesystem-support |
| Tools/Calls/Runs/parameters | Studio | execution attribution | yes for executed work | exact versions | missing-tool/version |
| Observations/Derived Artifacts | CAP-INV-367..377 | outputs/lineage | according scope | linked versions | partial |
| Access/errors/decisions | Security/Shared/Investigate | disposition context | yes when applicable | append-only | disputed |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case/Disk Image/Artifact | Investigate | source/version/lineage | read |
| Collection Request/Job/Endpoint | owners | acquisition context | read/link |
| Disk Session/observations/results | Investigate concepts | analytical history | read |
| Tool/Tool Call/Automation Run | Studio | execution/version | read |
| Trace/Activity/Audit events | Shared/Security | mechanism records | read/export if authorized |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Reproducibility Assessment | create/dispute/supersede | Investigate concept | missing prerequisites explicit |
| Provenance relation | create/link/version | Investigate/Shared | source/producer/time required |
| Reproduction request/result | prepare/record | Investigate + Studio | explicit Tool and bounds |
| Source trace | no deletion/rewrite | owner mechanism | append/supersede only |

## 11. Fonctionnalités
Navigate end-to-end lineage; display all sources, Tools, versions, parameters, filters, decisions and errors; distinguish reproducible, partially-reproducible, not-reproducible, missing-image, incomplete-image, missing-volume, missing-filesystem-support, missing-tool, missing-version, restricted-content, policy-blocked and disputed; request bounded reproduction; compare result and original assessment.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspect/export provenance | Reviewer | provenance | 0/1 | permissions | sourced chain/export | policy |
| Assess/dispute | Reviewer | assessment | 2 | reason/sources | versioned status | OPEN-013 |
| Reproduce analysis | Analyst | Tool Call request | 1 | prerequisites/authorization | linked result | no |
| Supersede assessment | Reviewer | assessment | 2 | new evidence | old retained | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| assemble lineage | yes | IDs/relations/timestamps | yes | summary | deterministic trace view |
| detect missing prerequisite | yes | rules | yes | explanation | checklist/status matrix |
| reproduce bounded analysis | explicit | Tool/workflow | yes | no silent agent | manual Tool run |
| assess equivalence/conclusion | human | comparison aids | no autonomous | assistance | reviewer comparison |

## 14. États fonctionnels
`reproducible`, `partially-reproducible`, `not-reproducible`, `missing-image`, `incomplete-image`, `missing-volume`, `missing-filesystem-support`, `missing-tool`, `missing-version`, `restricted-content`, `policy-blocked`, `disputed`.

## 15. États d’interface
Loading preserves chain; Empty names missing records; Partial highlights gaps; Error retains valid lineage; Offline read-only; Permission denied masks protected access; Stale shows superseded versions.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Provenance chain | relations/projection | Case/Replay/Evidence/Audit | sources, producers, errors and gaps visible |
| Reproducibility Assessment | assessment | CAP-INV-379/Quality | explicit status/reasons |
| Reproduction result relation | relation | analyst/Trace | original and reproduced contexts distinct |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Session/result | review provenance | CAP-INV-378 | Case, image, structures, Tools, outputs, errors | source view |
| Assessment | reproduce | Studio Tool/Workflow | exact source/version/parameters/bounds | Provenance |
| Provenance selection | handoff | CAP-INV-379 | chain, gaps, status, restrictions | Provenance |

## 18. Dépendances
CAP-INV-105/203/213/214/311/364..377/379, Studio Tool/Call/Run, Shared Trace/Activity/Versioning/Export/Recovery, Security, OPEN-005/008/013/014/015.

## 19. Source de vérité
Every owner remains source for its records; Investigate owns assessment and local provenance relations; Shared mechanisms are consumed, not redefined.

## 20. Provenance et audit
This capability traces all phase elements plus reviewer, reason, assessment version, reproduction request/result, comparison and disposition. Trace removal is prohibited.

## 21. Permissions fonctionnelles
Custody/provenance read/review/export, reproducibility assess/dispute, automated reproduction request, restricted access-event read and cross-tenant constraints. Atomic rules future.

## 22. Limites et erreurs
Missing image/volume/filesystem support/Tool/version, restricted content, policy block, partial source, unavailable provider or inconsistent parameters yield explicit non/partial reproducibility; no infrastructure defined.

## 23. Métriques
Assessments by state, missing prerequisites, successful/partial reproductions, lineage gaps, disputed/superseded assessments and provenance exports.

## 24. Classification de livraison
`defined` / `planned`; no replay infrastructure, schema, engine, API, protocol or implementation.

## 25. Critères d’acceptation
**Given** the Tool version is unavailable **When** reproducibility reviewed **Then** status is `missing-version` or partially/not reproducible, with original result preserved.

**Given** restricted content blocks replay **When** reviewed **Then** `restricted-content` is visible and no bypass is attempted.

**Given** no AI **When** provenance assembled **Then** IDs, relations, versions, parameters, trace and checklist provide the essential workflow.

## 26. Questions ouvertes
OPEN-005/008/013/014/015 remain open; final provenance contracts and infrastructure are future.

## 27. Consommateurs documentaires
INV-DSK-001, Case Replay/Evidence/Audit, CAP-INV-379, Quality/Release and Objects/Permissions/Trust/Technique.
