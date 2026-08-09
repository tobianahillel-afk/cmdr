---
id: CAP-GOV-038
title: Govern Audit Review, Search and Evidence Package Preparation
product: govern
module: audit-trail
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-002, REQ-PROD-006, REQ-PROD-008, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-013, OPEN-014, OPEN-019]
source-of-truth: canonical
---
# CAP-GOV-038 — Govern Audit Review, Search and Evidence Package Preparation

## 1. Définition
Permettre une revue Govern permission-aware des Audit Events, reconstructions et completeness assessments avec recherche/pivots/comparaisons/annotations, puis préparer un Audit Evidence Package sourcé et masqué pour revue, Reporting/Export ou handoff Investigate, sans redéfinir Shared Search/Export/Report ni créer automatiquement Evidence.

## 2. Problème utilisateur
Un auditeur doit explorer une grande chaîne sans perdre tenant, restrictions ou provenance. Un export brut peut exposer des identités ou exceptions sensibles, et un « evidence package » d’audit ne doit pas être pris pour Evidence Investigate sans qualification explicite.

## 3. Objectifs
- rechercher/filtrer/pivoter sur objets, actors, versions, Decision/Run et périodes ;
- comparer Decisions/Runs/reconstructions ;
- annoter, assigner review, bookmarker et produire review summaries ;
- préparer un package avec exact sources/restrictions/masking ;
- déléguer rendu/export à Shared ;
- permettre handoff Investigate sans transfert implicite d’ownership.

## 4. Non-objectifs
Ne pas créer un moteur Search/Reporting/Export, publier à l’externe automatiquement, créer Evidence canonique, contourner masking, modifier les objets historiques, finaliser un format d’export/API ou définir un dashboard détaillé.

## 5. Propriétaire
Govern / Audit Trail owns audit-review semantics and Audit Evidence Package composition. Shared owns Search, Inspector, Reporting and Export mechanics. Investigate owns Evidence qualification/lifecycle. Security/Settings own permission, sharing, retention and export destinations.

## 6. Utilisateurs
Principal : Govern Auditor. Secondaires : Compliance/Security Reviewer, Govern Reviewer, Decision/Response Reviewer, Investigator receiving an explicit handoff, authorized Report Reviewer.

## 7. Conditions d’entrée
At least one Govern Audit Event, reconstruction or assessment is visible within a permitted tenant/environment scope. Restricted fields have masking/classification metadata and export/sharing rights are independently evaluated.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| audit events/reconstructions | CAP-GOV-034..037 | review content | oui | selected snapshot | empty/partial |
| search/filter capability | Shared | query mechanism | oui for search | service freshness shown | browse fallback/limited |
| classification/masking | Security/source | sensitivity controls | oui when sensitive | current policy | export blocked/masked |
| tenant/environment | Settings/context | scope | oui | current context | block operation |
| Report/Export capability | Shared | rendering/export | when requested | service state | package retained, export unavailable |
| destination/handoff context | destination owner | consumer restrictions | when handoff | current | no implicit handoff |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Govern Audit Event/Reconstruction/Assessment | Govern | review content | read/review |
| Report | Shared Reporting | rendered package/report state | read/link only |
| Search/Export Job | Shared | query/export status | read/link |
| Evidence/Case | Investigate | target handoff context | read/link, no creation as Evidence |
| classification/tenant | Security/Settings | masking/scope | restricted read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Audit Review | create/update/assign/annotate/close/supersede | Govern local concept | source refs retained |
| Audit Evidence Package | create/update/version/supersede | Govern local concept | package ≠ Evidence; restrictions preserved |
| Report/Export Job | request via Shared | Shared | no local engine/ownership |
| Evidence/Case | handoff request only | Investigate | destination decides qualification |

## 11. Fonctionnalités
Search/filter/pivot by lifecycle object/actor/action/tenant/environment/time; traverse correlation; compare Decision/Run versions; annotate and assign review; bookmark stable refs; collect scoped source refs; apply masking/classification; create package manifest/summary/limitations; preview package; request Shared export/report; prepare Investigate handoff; preserve review history and return origin.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| search/filter/inspect | auditor | audit corpus | 0 | read permission | scoped review view | non |
| deterministic compare/package completeness check | auditor | review/package | 1 | source set | diff/missing list | non |
| annotate/assign/bookmark | auditor | Audit Review | 2 | review permission | versioned review state | OPEN-013 |
| create Audit Evidence Package | auditor | package | 2 | scoped sources + classification | versioned package | OPEN-013/014 |
| request export/handoff | authorized auditor | package | 2 | destination/export permission | Shared job/handoff request | OPEN-019 |
| external disclosure/create Evidence automatically | none | package | 3/4 | prohibited here | no action | destination governance |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| search/filter/pivot | oui | query/filter | oui | query suggestion | filter builder |
| compare versions/runs | oui | diff | oui | explanation | diff viewer |
| assemble package manifest | oui | selected refs | oui | draft summary | manifest table |
| summarize review | oui | templates | oui | sourced draft | structured summary |
| decide external sharing/Evidence qualification | destination human/governance | permission checks | no autonomous | prohibited | explicit approval/handoff |

## 14. États fonctionnels
Audit Review: `draft`, `assigned`, `in-review`, `information-required`, `reviewed`, `closed`, `superseded`. Package: `draft`, `incomplete`, `ready-for-review`, `restricted`, `export-prepared`, `handoff-prepared`, `superseded`. These are not Report/Evidence states.

## 15. États d’interface
Loading preserves query/review selection ; Empty means no visible results ; Partial exposes inaccessible sources ; Error retains stable refs ; Offline supports safe local reading only where available ; Permission denied masks restricted data ; Stale identifies snapshot/index freshness.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Audit Review | Govern review record | auditor/CAP-GOV-047 | sources/annotations/limitations attributable |
| Audit Evidence Package | Govern package | Shared Reporting/Export/Investigate handoff | exact refs + masking/classification retained |
| review summary | derived review output | authorized consumers | summary ≠ audit finding/fact |
| export/handoff request | request event | Shared/Investigate | destination permission/ownership retained |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-035..037 | review requested | CAP-GOV-038 | reconstruction/assessment refs | source review |
| CAP-GOV-038 | export/report requested | Shared Reporting/Export | package manifest/restrictions | Audit Review |
| CAP-GOV-038 | Investigate qualification requested | Investigate | package/provenance/classification | Audit Review |
| CAP-GOV-038 | improvement issue identified | CAP-GOV-047 | review summary/gaps | Audit Review |

## 18. Dépendances
CAP-GOV-034..037/047, Shared Search/Inspector/Reporting/Export/Jobs/Versioning, Security secure-export/privacy, Settings tenant/export destinations/retention, Investigate Evidence/Case, OPEN-013/014/019.

## 19. Source de vérité
Govern owns review/package composition; included objects/events remain source-owned. Shared owns rendered Report/export execution. Investigate alone determines whether a received package/material becomes canonical Evidence.

## 20. Provenance et audit
Record query/filter snapshot, selected refs/versions, reviewer/assignee, annotations, package manifest, classification/masking decisions, excluded items/reasons, Report/Export Job refs, handoff destination, AI draft provenance and supersession.

## 21. Permissions fonctionnelles
Audit Trail read/search; restricted audit read; actor identity read; Audit Review create/assign; package create; provenance export preparation; cross-tenant metric/audit reads separately authorized; external-sharing preparation does not grant disclosure permission.

## 22. Limites et erreurs
Search unavailable, index stale, source permission lost, export engine unavailable, redaction policy missing, destination inaccessible, cross-tenant selection or oversized/partial package keeps review/package partial/blocked. Export cannot widen visibility.

## 23. Métriques
Review volume/age; package completeness; masked/restricted item count; export/handoff requests; failed exports due to permission/classification; packages promoted automatically to Evidence — target zero.

## 24. Classification de livraison
`defined` / `planned`; no Search/Reporting/Export implementation or external-sharing contract selected.

## 25. Critères d’acceptation
**Given** an audit export includes restricted content, **When** export is requested, **Then** classification/masking and export permission are re-evaluated and the export cannot widen visibility.

**Given** an Audit Evidence Package is handed to Investigate, **When** it arrives, **Then** it remains a sourced package and does not become canonical Evidence until Investigate explicitly qualifies it.

**Given** cross-tenant records are selected without permission, **When** review/package is prepared, **Then** those records remain inaccessible/excluded and the scope does not silently expand.

**Given** no AI, **When** an audit review/package is produced, **Then** filters, source tables, diffs, manifests and human review provide full functionality.

## 26. Questions ouvertes
OPEN-013/014/019 remain open. Final export formats, attachment/material relation and external release policy remain future; no new OPEN.

## 27. Consommateurs documentaires
Audit Trail, Shared Reporting/Export, Investigate handoff, CAP-GOV-047, Govern closure, future screen/object/permission/technical phases and compliance/quality reports.