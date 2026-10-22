---
id: CAP-INV-718
title: Mobile Artifact Extraction and Handoff to Evidence, Findings and Detection
product: investigate
module: mobile-forensics
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-07
requirement_ids: [REQ-INV-001, REQ-PROD-014, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-005, OPEN-011, OPEN-013, OPEN-014, OPEN-015, OPEN-017, OPEN-018]
source-of-truth: canonical
---
# CAP-INV-718 — Mobile Artifact Extraction and Handoff to Evidence, Findings and Detection

## 1. Définition
Sélectionner un fichier ou record et son contexte source, définir objectif/restrictions, exécuter une extraction analytique bornée autorisée, produire un Derived Artifact avec lineage/Tool/version/paramètres/erreurs/partialité, puis préparer Evidence Candidate Package, Finding Draft, Static/Reverse/Dynamic, Detection Engineering, Threat Intelligence ou complementary Collection handoffs.

## 2. Problème utilisateur
Un résultat Mobile utile doit pouvoir quitter son viewer sans perdre source, privacy, transformations ou uncertainty. Sans handoff explicite, un Derived Artifact peut être confondu avec Evidence ou un Finding, ou entraîner une règle/réponse automatique.

## 3. Objectifs
- bounded extraction/recovery/derivation from authorized represented data;
- preserve source device/package/app/record/region and all transformations;
- create Derived Artifact distinct from Evidence;
- prepare destination packages/drafts with restrictions and return origin;
- never auto-qualify Evidence/Finding, deploy rule or execute collection/response.

## 4. Non-objectifs
No acquisition from original device, unlock/bypass/root/jailbreak, external transmission, automatic Evidence/Finding, rule creation/deployment, Indicator activation, response, target mutation, secret use, final artifact format or engine.

## 5. Propriétaire
Investigate owns bounded Mobile derived-selection semantics and candidate/draft packages. Artifact/Evidence/Finding owner workflows remain Investigate canonical; Static/Reverse/Dynamic, Detection, TI, Collection and Govern own destination processing/actions; Studio owns Tool/Run.

## 6. Utilisateurs
Principal : Investigation Lead or Mobile Forensics Analyst. Secondaires : DFIR Analyst, Evidence Reviewer, Sensitive Data Reviewer and SOC Analyst.

## 7. Conditions d’entrée
Selected authorized source observation/file/record, purpose, source permission, export/extraction permission where applicable, source lineage, integrity/partiality context, restrictions and destination question.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Selected source file/record/observation | CAP-INV-706..717 | derivation source | oui | source version | no artifact/handoff |
| Device/package/app/region context | Mobile observations | provenance | oui where known | Session version | provenance gap/blocked |
| Integrity/completeness/accessibility | CAP-INV-705 | trust limits | oui | reviewed version | partial/unverified retained |
| Purpose/destination/restrictions | analyst/Security | bounded use | oui | handoff time | blocked |
| Tool/Tool Call/Run/parameters | Studio | derivation execution provenance | if Tool used | run version | no Tool result invented |
| Permissions/classification | Security/source | extract/export boundary | oui | action time | denied/masked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Mobile observations/Session | Investigate | source/context/selection | read |
| Artifact / Evidence / Finding | Investigate canonical workflows | destination/status context | read/link, no auto creation beyond owner workflow |
| Tool/Tool Call/Automation Run | Studio | execution/provenance | select/invoke/read by permission |
| Detection/TI objects | owner modules | destination context | read/link |
| Collection Request / Govern Action Request | owners | complementary action context | prepare/link only |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Derived Artifact | create/read/version/supersede/withdraw | Investigate | source lineage, transformations, Tool/params, restrictions mandatory |
| Evidence Candidate Package | prepare/version/withdraw | Evidence workflow | candidate ≠ qualified Evidence |
| Finding Draft | prepare/version/withdraw | Finding workflow | draft ≠ confirmed Finding |
| Detection Engineering / Intelligence Handoff Package | prepare/version/withdraw | destination owner | no rule/Indicator/knowledge auto-created |
| Complementary Collection proposal | prepare | Collection | no acquisition executed |

## 11. Fonctionnalités
Select source region/record/file, preview bounded derivation, choose purpose/destination, run deterministic Tool where permitted, show transformations/Tool/params/status/errors, preserve masking/classification, compare output with source, create Derived Artifact, build destination package with source/uncertainty/limitations and track return status.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Select source | analyste | record/file | 0 | source read | bounded selection | non |
| Run bounded extraction/recovery | authorized analyst | derivation | 1 | extraction permission + parameters | Derived Artifact candidate/result | no real device action |
| Create/version Derived Artifact | analyste | Artifact | 2 | successful/partial result + lineage | versioned derived material | OPEN-013 |
| Prepare Evidence/Finding/Detection/TI package | lead/reviewer | handoff | 2 | destination purpose | candidate/draft/package | destination re-evaluates |
| Prepare collection/response need | lead | request context | 2→3 | gap/risk | Collection/Govern package | destination executes if approved |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Bounded deterministic extraction | oui | oui where Tool exists | oui | non | explicit Tool/manual selection |
| Generate manifest/lineage | oui | oui | oui | non | deterministic metadata |
| Suggest destination/package fields | oui | templates/rules | oui | suggestion | forms/checklists |
| Summarize observations/limitations | oui | structured | oui | yes | table/manual summary |
| Qualify Evidence/Finding/deploy rule/respond | non | destination human workflow | non | interdit | destination review/authority |

## 14. États fonctionnels
Extraction/derivation: `proposed`, `queued`, `processing`, `available`, `partial`, `invalid`, `restricted`, `failed`, `cancelled`, `superseded`, `withdrawn-from-use`. Handoff packages add `draft`, `ready`, `submitted`, `returned`, `superseded`, `withdrawn` as functional vocabulary only.

## 15. États d’interface
Loading preserves source selection; Empty requires source/purpose; Partial exposes partial output; Error keeps source and prior valid result; Offline disallows new Tool execution/export not guaranteed; Permission denied reveals no protected content; Stale requires source/version revalidation.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Derived Artifact | Investigate material | Static/Reverse/Dynamic/Evidence | bounded source lineage + restrictions; not Evidence |
| Evidence Candidate Package | candidate package | Evidence Board/workflow | qualification remains destination-owned |
| Finding Draft | draft | Finding workflow | confirmation remains destination-owned |
| Detection Engineering Package | package | CAP-INV-401..435 | no rule/content/deployment automatic |
| Intelligence Handoff Package | package | CAP-INV-501..537 | no attribution/knowledge promotion automatic |
| Collection gap package | request input | Collection | no acquisition execution |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-706..717 | selected material/result | CAP-INV-718 | source/device/package/app/record, purpose, restrictions | source analysis |
| Derived Artifact | static need | Static Analysis | Artifact + lineage + question | Mobile Session |
| Derived Artifact | reverse/dynamic need | Reverse/Dynamic | Artifact + lineage + question | Mobile Session |
| Mobile result package | evidence/finding | Evidence/Finding owner | candidate/draft + source/limits | Mobile Session |
| Mobile result package | detection/TI | Detection/TI | sourced package + uncertainty | Mobile Session |
| Missing source | collection | CAP-INV-202 | bounded need/source/authority | Mobile Session |
| All outputs | closure | CAP-INV-719 | lineage, Tools/Runs, errors, dispositions | Session |

## 18. Dépendances
CAP-INV-202, CAP-INV-301..346, CAP-INV-401..537, CAP-INV-705..717/719, Artifact/Evidence/Finding, Studio Tools/Runs, Shared Export/Trace, Govern, OPEN-005/011/013/014/015/017/018.

## 19. Source de vérité
Raw selected material remains source. Derived Artifact/selection lineage is Investigate. Tool/Run is Studio. Evidence/Finding, Detection, TI, Collection and Govern destinations own qualification/action. Handoff never transfers source permissions.

## 20. Provenance et audit
Source package/device/app/record/region/version, selection, purpose, restrictions, integrity/partiality, Tool/Call/Run/version/parameters, transformations, output hashes/identifiers if provided by implementation later, errors, sensitive accesses, destination package/version, reviewer and destination return.

## 21. Permissions fonctionnelles
Source read, bounded extraction/recovery run, Derived Artifact create/read/export, restricted source/content read, Evidence Candidate prepare, Finding Draft prepare, Detection/TI handoff prepare, Collection gap prepare, automated Tool request, provenance export. Export never implies source access for recipient.

## 22. Limites et erreurs
Tool failure, partial output, unsupported source, transformed/ambiguous record, permission denial, source version change, privacy restriction or destination rejection preserves lineage and allows invalidate/supersede/withdraw without deleting history.

## 23. Métriques
Derived Artifacts by source/status, partial/failed derivations, Evidence/Finding/Detection/TI/Collection handoffs, destination accept/modify/reject returns, permission blocks, source-to-output lineage completeness and withdrawn artifacts.

## 24. Classification de livraison
`defined` / `planned`. No artifact format, extractor engine, decompiler/sandbox implementation, detection runtime, TI protocol, collection mechanism, API or code delivered.

## 25. Critères d’acceptation
**Given** a selected recovered fragment **When** a Derived Artifact is produced **Then** recovery quality, source region, Tool/parameters and partiality remain linked and it is not qualified Evidence.

**Given** a Detection Engineering handoff **When** submitted **Then** no Detection Content or rule is created/deployed automatically.

**Given** no AI **When** an artifact/handoff is prepared **Then** deterministic extraction, manifests, forms, templates and human review provide full functionality.

## 26. Questions ouvertes
OPEN-005 Tool/engine selection; OPEN-011 Mobile delivery; OPEN-013 class-2 mutations; OPEN-014 Artifact/material relations; OPEN-015 Run lineage; OPEN-017 Detection runtime; OPEN-018 TI interoperability remain open.

## 27. Consommateurs documentaires
Evidence/Finding, Static/Reverse/Dynamic, Detection Engineering, Threat Intelligence, Collection, Govern, Studio, Shared Export/Trace, Objects, Permissions, Quality and Technique.
