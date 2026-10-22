---
id: CAP-INV-715
title: Mobile Deleted, Residual and Recovered Data Analysis
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
# CAP-INV-715 — Mobile Deleted, Residual and Recovered Data Analysis

## 1. Définition
Analyser deleted-entry candidates, residual metadata/content, recovered fragments, carved content, partial/orphaned records, application/filesystem deleted-data candidates, source region, recovery quality, attribution limits, contradictions, Derived Artifacts and withdrawal/supersession.

## 2. Problème utilisateur
Un record marqué supprimé ne prouve pas l’intention de suppression; un fragment récupéré ou carved ne prouve ni contenu original complet ni attribution certaine. Recovery quality and source-region context are critical to avoid overclaiming.

## 3. Objectifs
- preserve candidate/recovery state and source region;
- distinguish deleted, residual, recovered, carved, partial and orphaned material;
- record recovery quality, gaps, overlaps and contradictions;
- allow bounded deterministic recovery where authorized without choosing a final engine;
- prepare Derived Artifact/Timeline/Hypothesis handoffs with attribution limits.

## 4. Non-objectifs
No destructive recovery, source modification, undelete on original device, final carving algorithm, filesystem-specific exploit, bypass, acquisition, malware confirmation, user-intent inference, engine/tool selection or code.

## 5. Propriétaire
Investigate owns Deleted Entry Candidates, Recovery Results and analytical observations. Source remains package/extraction; Shared may provide Job/Trace mechanisms; CAP-INV-718 owns bounded derived handoff preparation.

## 6. Utilisateurs
Principal : DFIR Analyst. Secondaires : Mobile Forensics Analyst, Evidence Reviewer and Investigation Lead.

## 7. Conditions d’entrée
Session, represented filesystem/app/backup area, source region and limitations, permission to deleted/residual content, integrity/completeness context and explicit analytical purpose.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Deleted/residual candidate records | filesystem/app/backup source | candidate material | oui for candidate analysis | source version | no deleted claim |
| Source region/path/app | CAP-INV-706/708/714 | provenance context | oui where available | Session version | attribution limited |
| Recovery/parsing result | Tool/Studio or deterministic analyzer | bounded derived result | non | run/version | no recovered content invented |
| Integrity/completeness limits | CAP-INV-705 | trust boundary | oui | reviewed version | partial/unverified |
| Permissions/classification | Security/source | access boundary | oui | access time | masked/denied |
| Contradicting live/current records | other observations | comparison context | non | source versions | unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Session / Filesystem / App / Backup observations | Investigate | source region/context | read |
| Deleted/residual source material | source | metadata/content | read by permission |
| Tool / Tool Call / Automation Run | Studio | recovery provenance | read/link |
| Artifact | Investigate | Derived Artifact relation | read/link |
| Shared Jobs/Trace/Versioning | Shared | mechanisms | consume |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Deleted Entry Candidate | create/review/dispute/supersede/withdraw | Investigate concept | deletion ≠ user intent |
| Recovery Result | create/review/invalidate/supersede | Investigate concept | quality/source region/Tool required |
| Residual/Carved Observation | create/review/dispute | Investigate | fragment ≠ complete original or certain attribution |
| Derived Artifact proposal | prepare/withdraw | CAP-INV-718 | source/recovery lineage mandatory |

## 11. Fonctionnalités
List/search/filter deleted/residual candidates, inspect source region and raw/metadata according to permission, run/request bounded deterministic recovery if authorized, compare recovered/live records, score recovery quality conceptually, flag partial/orphaned/overlap, annotate attribution limits and select Derived Artifacts.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Browse candidates | analyste | source/candidates | 0 | read | scoped view | non |
| Run bounded recovery | authorized analyst | source region | 1 | source permission + bounded parameters | Recovery Result | no device mutation |
| Compare recovered/live | analyste | records | 1 | both accessible | diff/relations | non |
| Dispute/invalidate candidate/result | reviewer | observation | 2 | reason/source | versioned status | OPEN-013 |
| Prepare Derived Artifact | analyste | recovery/result | 1/2 | purpose/restrictions | CAP-INV-718 input | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Identify deletion markers/candidates | oui | rules | oui | suggestion | deterministic parser |
| Recover bounded fragments | oui | deterministic tool where available | oui | non | explicit recovery Tool |
| Compare recovered/live | oui | oui | oui | summary | diff/table |
| Group residual fragments | oui | rules | oui | suggestion | source-region grouping |
| Infer user deletion intent/attribution | non | non | non | interdit | human sourced assessment |

## 14. États fonctionnels
Extraction/recovery: `proposed`, `queued`, `processing`, `available`, `partial`, `invalid`, `restricted`, `failed`, `cancelled`, `superseded`, `withdrawn-from-use`. Candidate review also uses supported/contradicted/inconclusive/disputed states.

## 15. États d’interface
Loading preserves source region; Empty means no represented candidates, not no prior data; Partial marks incomplete recovery; Error preserves valid fragments; Offline only on existing data; Permission denied masks residual content; Stale shows source/Tool version.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Deleted/Residual Observations | concepts | CAP-INV-716/717/718 | candidate state/source/attribution limits |
| Recovery Results | concepts | analyst/QA/CAP-INV-718 | Tool/version/parameters/quality/errors |
| Comparison relations | candidates | Timeline/Hypothesis | no original completeness claim |
| Derived Artifact proposal | package input | CAP-INV-718 | recovery lineage/restrictions |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-706/708/709/710/714 | deleted/residual candidate | CAP-INV-715 | source region/record/state/limits | source analysis |
| CAP-INV-715 | temporal/correlation | CAP-INV-716 | candidate/result/timestamps/quality | deleted analysis |
| CAP-INV-715 | anomaly/hypothesis | CAP-INV-717 | recovery observations + contradictions | deleted analysis |
| Recovery result selected | derive/handoff | CAP-INV-718 | source region, Tool/params, quality, restrictions | deleted analysis |
| Closure | provenance | CAP-INV-719 | recovery runs/results/decisions | Session |

## 18. Dépendances
CAP-INV-705/706/708..710/714/716..719, Studio Tools/Runs, Shared Jobs/Trace, Security, OPEN-005/011/013/014/015.

## 19. Source de vérité
Raw residual/deleted markers remain source facts. Recovery Result is an Investigate analytical concept derived by attributed Tool/Run. Source is never rewritten. Evidence qualification and Finding confirmation remain downstream workflows.

## 20. Provenance et audit
Source package/region/path/app, candidate marker, Tool/Run/version/parameters, recovery timestamps/status/errors, quality/partiality, comparison inputs, privacy access, reviewer/dispute/invalidation, Derived Artifact selection and supersession.

## 21. Permissions fonctionnelles
Deleted-data read, residual-content read, recovery prepare/run/read, restricted recovery read, compare, invalidate/dispute, Derived Artifact create/export, Tool result read and provenance export.

## 22. Limites et erreurs
Overwritten/fragmented/malformed regions, parser/recovery failure, false deletion marker, partial extraction, source ambiguity, permission denial, encrypted/protected area or conflicting live record remain explicit; no original completeness/user intent/attribution is inferred.

## 23. Métriques
Candidates by source/type, recovery success/partial/invalid rates, fragments recovered, disputed/withdrawn results, permission blocks, derived handoffs and false-positive candidate rate after human review.

## 24. Classification de livraison
`defined` / `planned`. No carving engine, filesystem-specific recovery algorithm, tool vendor, source-modifying recovery, API or code delivered.

## 25. Critères d’acceptation
**Given** a carved fragment **When** reviewed **Then** recovery quality, source region and attribution limits are visible and it is not presented as complete original content.

**Given** a deleted-record marker **When** no user-action evidence exists **Then** it remains a deletion candidate and does not establish user intent.

**Given** no AI **When** recovery is performed **Then** deterministic Tools/rules, tables, diff and human review cover the full workflow.

## 26. Questions ouvertes
OPEN-005 recovery engines; OPEN-011 platform/source methods; OPEN-013 recovery/dispute mutations; OPEN-014 material relations; OPEN-015 Tool/Run contracts remain open.

## 27. Consommateurs documentaires
Timeline, Anomaly/Hypothesis, Artifact/Evidence handoff, Disk/Filesystem boundaries, Privacy/Trust, Objects, Permissions, Quality and Technique.
