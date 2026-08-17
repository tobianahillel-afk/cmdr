---
id: CAP-INV-375
title: File Carving and Recovery
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
# CAP-INV-375 — File Carving and Recovery

## 1. Définition
Lancer explicitement un traitement borné de carving ou recovery sur une source sélectionnée, conserver progression, erreurs, partialité et lineage, et produire des Derived Artifacts sans modifier l’image ni garantir nom, chemin, complétude ou attribution.

## 2. Problème utilisateur
Une récupération peut produire des fragments ou faux rapprochements. Sans bornes et provenance, un résultat incomplet peut être présenté comme fichier original.

## 3. Objectifs
Sélectionner source/objective/limits/restrictions/volume risk; launch/interrupt; show progress/partial results; create Derived Artifacts; preserve source/result, Tool/version/parameters/errors; compare metadata; mark incomplete/invalid; route Static/Reverse; withdraw without deleting trace.

## 4. Non-objectifs
Ne pas imposer algorithme, signature, moteur, commande ou format; ne pas promettre original name/path/completeness/attribution; ne pas modify source, bypass protection or auto-create Evidence.

## 5. Propriétaire
Investigate owns recovery plan/result interpretation and Derived Artifact relation; Studio owns Tool/Call/Run; Shared Jobs/Export; Security controls sensitive scope.

## 6. Utilisateurs
Principal : Recovery Analyst. Secondaires : Disk Analyst, Artifact Analyst, Privacy Reviewer, Evidence Reviewer.

## 7. Conditions d’entrée
Bounded deleted/unallocated/source region; immutable image; objective/limits/volume risk shown; Tool and permission explicit; storage/retention projection available where required.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Bounded source selection | CAP-INV-369 or analyst | source/scope | yes | image version | blocked |
| Objective/limits/volume risk | analyst/policy | processing bounds | yes | launch snapshot | incomplete |
| Tool/version/parameters | Studio | execution context | yes | selected version | tool-unavailable |
| Restrictions/permissions | Security/Settings | access/export | yes | current | policy-blocked |
| Existing metadata/relations | CAP-INV-367/368 | comparison context | no | same image | unattributed |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Disk Image / source region | Investigate | immutable source | read |
| Deleted/Unallocated Observation | Investigate concept | quality/selection | read |
| Tool/Tool Call/Automation Run | Studio | producer/version/parameters | select/read |
| Background Job | Shared | progress/cancel/partial | consume |
| Artifact/Derived Artifact | Investigate | lineage/versions | read/link |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Recovery Plan | create/update/supersede | Investigate concept | bounded scope required |
| Recovery Result | create/mark partial-invalid-withdraw | Investigate concept | uncertainty and errors retained |
| Derived Artifact | create/link | Investigate | parent, source range, Tool/version required |
| Source image | no mutation | Investigate | immutable |

## 11. Fonctionnalités
Define source/objective/limits/restrictions/volume risk; explicit start and interrupt; progress and partial results; Derived Artifacts and source/result relation; Tool/version/parameters/errors; metadata comparison; incomplete/invalid/unattributed/restricted; route Static/Reverse; withdraw from active use without trace deletion.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Prepare bounded recovery | Analyst | Recovery Plan | 2 | scope/permission | versioned plan | OPEN-013 |
| Start/reproduce | Analyst | Tool Call | 1 | explicit confirmation | Background Job | no |
| Interrupt | Analyst | Job | 2 | cancellable | partials retained | OPEN-013 |
| Review/mark/withdraw | Analyst | Recovery Result | 2 | result available | disposition retained | OPEN-013 |
| Export/route Artifact | Analyst | Derived Artifact | 1/2 | policy | handoff/export | policy |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| scan bounded source | yes | Tool | yes | no need | deterministic recovery Tool |
| group candidate results | yes | rules | yes | suggestion | metadata/content indicators |
| explain partiality | yes | diagnostics | yes | summary | error/coverage table |
| confirm original identity | human | no | no | assistance | provenance review |

No silent Tool, no automatic certainty.

## 14. États fonctionnels
`proposed`, `processing`, `available`, `partial`, `invalid`, `unattributed`, `restricted`, `superseded`, `withdrawn-from-use`.

## 15. États d’interface
Loading shows source/bounds; Empty means no result, not no data; Partial shows missing fragments; Error preserves outputs; Offline blocks new run; Permission denied masks content; Stale shows Tool/source version.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Recovery Result | analysis concept | analyst/CAP-INV-378 | scope, errors and quality visible |
| Derived Artifact | Artifact relation | Static/Reverse/CAP-INV-379 | parent/source/Tool lineage |
| Progress/disposition event | Job/Trace | Session/Notifications | interruption and partials visible |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Unallocated/deleted source | prepare | CAP-INV-375 | image, range, objective, restrictions | CAP-INV-369 |
| Recovery result | route | Static/Reverse | Derived Artifact, source, quality, objective | Recovery |
| Result/Artifact | handoff | CAP-INV-379 | partiality, contradictions, lineage | Recovery |

## 18. Dépendances
CAP-INV-311/369/368/376/378/379, Studio Tools/Runs, Shared Jobs/Trace/Export, Settings retention/storage, OPEN-005/008/013/014/015.

## 19. Source de vérité
Disk Image remains source; Investigate owns Recovery Plan/Result and lineage; Studio owns execution records; Shared owns Job/Export mechanisms.

## 20. Provenance et audit
Case/session/image/source region, objective/bounds/risk, Tool/version/Call/Run, parameters, progress, interruption, result quality, errors, Derived Artifact, access/export and human disposition.

## 21. Permissions fonctionnelles
Recovery prepare/run/cancel/read/export, unallocated/deleted read, Derived Artifact create/read/export, restricted source access and reproducibility review. Final step-up/SoD future.

## 22. Limites et erreurs
No guaranteed name, path, completeness or attribution. Volume limit, overlap, corruption, unsupported content, policy block, permission denial, Tool failure and interruption remain explicit. Source never modified.

## 23. Métriques
Runs by status, bytes/regions in authorized scope, partial/unattributed/invalid results, interruptions, outputs with lineage and withdrawals.

## 24. Classification de livraison
`defined` / `planned`; no algorithm, engine, format, command or implementation chosen.

## 25. Critères d’acceptation
**Given** unallocated source and bounded treatment producing incomplete output **When** result completes **Then** source stays unchanged, result is partial, no original name/path is invented and Tool/version/lineage are visible.

**Given** volume limit reached **When** run stops **Then** completed outputs and excluded scope are visible and no hidden continuation occurs.

**Given** no AI **When** recovery runs **Then** deterministic Tool, progress, filters, comparison and human review work.

## 26. Questions ouvertes
OPEN-005/008/013/014/015 remain open; algorithms, formats, objects and permission rules are future.

## 27. Consommateurs documentaires
INV-DSK-001, CAP-INV-368/369/376/378/379, Static/Reverse, Artifact/Evidence and Objects/Permissions/Technique.
