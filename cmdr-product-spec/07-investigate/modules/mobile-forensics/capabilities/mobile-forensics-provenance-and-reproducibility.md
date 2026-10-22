---
id: CAP-INV-719
title: Mobile Forensics Provenance and Reproducibility
product: investigate
module: mobile-forensics
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-07
requirement_ids: [REQ-INV-001, REQ-PROD-014, REQ-PROD-020, REQ-OBJ-009, REQ-AI-002, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-011, OPEN-013, OPEN-014, OPEN-015]
source-of-truth: canonical
---
# CAP-INV-719 — Mobile Forensics Provenance and Reproducibility

## 1. Définition
Retracer de bout en bout Case/Incident/origine, device and declared owner context, platform candidate, acquisition authorization, Collection Request/Job, Mobile Evidence Package, backup/extraction representation, custody, integrity/completeness/accessibility, Session, Tools/Tool Calls/Automation Runs/parameters/queries, observations, Derived Artifacts, sensitive accesses, timelines/correlations, handoffs, errors/interruptions, human decisions, versions and Reproducibility Assessment.

## 2. Problème utilisateur
Une analyse Mobile peut combiner plusieurs représentations, données privées, transformations et Tools. Sans lineage complet, un reviewer ne peut distinguer source, dérivation, automatisation et jugement humain ni reproduire raisonnablement l’analyse.

## 3. Objectifs
- link every analytical result to source/version/scope/permissions;
- preserve Tool/Run/parameters/errors and human accept/modify/reject decisions;
- record sensitive-access events without raw secret/content leakage;
- assess reproducibility conditions, missing dependencies and limitations;
- preserve superseded/withdrawn history and destination handoff returns.

## 4. Non-objectifs
No guaranteed byte-for-byte replay, final provenance schema, immutable-ledger technology, evidence notarization, Tool packaging format, secret/content logging, automatic legal opinion, audit deletion or product implementation.

## 5. Propriétaire
Investigate owns Mobile provenance semantics and Reproducibility Assessment. Shared owns generic Trace/Activity/Versioning/Export. Studio owns Tool/Run execution records. Collection/Trust/source owners retain acquisition/custody facts.

## 6. Utilisateurs
Principal : Evidence Reviewer. Secondaires : Investigation Lead, Mobile Forensics Analyst, DFIR Analyst, Sensitive Data Reviewer, QA/Traceability Lead and auditors with permission.

## 7. Conditions d’entrée
A Mobile Session or result with at least one source reference, owner identity, current permission and accessible trace/version metadata. Missing provenance becomes an explicit gap, never reconstructed silently.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Origin/Session/device/platform/scope | CAP-INV-701..703 | analytical context | oui | final reviewed versions | reproducibility incomplete |
| Acquisition/Collection/custody | CAP-INV-704/Collection | source lineage | oui if applicable | recorded versions | gap explicit |
| Integrity/completeness/accessibility | CAP-INV-705 | trust context | oui | final assessment versions | gap explicit |
| Observations/timeline/hypotheses/artifacts | CAP-INV-706..718 | analytical results | selon Session | versioned | omitted result not reproducible |
| Tools/Calls/Runs/parameters/queries | Studio/Session | execution lineage | if used | immutable/versioned refs | Tool step marked missing |
| Sensitive access + human decisions | Security/Shared/Session | audit/disposition | if applicable | event versions | compliance gap explicit |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case/Incident/Session/Mobile concepts | owner/Investigate | origin and analysis lineage | read/link |
| Collection Request/Job/custody | Collection/Investigate | acquisition lineage | read/link |
| Tool/Tool Call/Automation Run | Studio | execution provenance | read/link |
| Trace/Activity/Versioning | Shared | generic audit mechanisms | consume |
| Evidence/Finding/Detection/TI handoffs | destination owners | return/disposition refs | read/link |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Mobile Analysis Provenance Package | create/version/supersede/export | Investigate concept using Shared | refs/versions/permissions/errors/human decisions mandatory |
| Reproducibility Assessment | create/review/dispute/supersede | Investigate concept | conditions/gaps explicit; replay not guaranteed |
| Provenance gap/correction proposal | create/resolve/supersede | Investigate | never edits historical source event |
| Shared Trace/Activity | append via owner mechanism | Shared | no trace deletion |

## 11. Fonctionnalités
Lineage graph/table with tabular alternative, origin-to-result trace, source/version filters, Tool/Run drill-down, parameter/query references, sensitive-access audit summaries, human disposition timeline, error/interruption chain, handoff return links, version diff and reproducibility checklist/export with masking.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspect lineage | reviewer | provenance | 0 | read | trace view | non |
| Compare versions | reviewer | provenance/results | 1 | read | diff | non |
| Create/update reproducibility assessment | reviewer | assessment | 2 | source refs | versioned assessment | OPEN-013 |
| Flag/correct provenance gap | reviewer/owner | proposal | 2 | reason | gap proposal/history preserved | non |
| Export provenance package | authorized reviewer | package | 1/2 | export permission | minimized classified export | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Assemble stable references/versions | oui | oui | oui | non | deterministic trace |
| Check missing lineage fields | oui | rules | oui | explanation | checklist |
| Compare parameter/source versions | oui | oui | oui | summary | diff/table |
| Summarize reproducibility gaps | oui | structured | oui | yes | checklist/report |
| Invent missing provenance/delete trace | non | non | non | interdit | explicit gap/correction proposal |

## 14. États fonctionnels
Reproducibility: `draft`, `under-review`, `reproducible-with-current-evidence`, `reproducible-with-limitations`, `not-reproducible`, `blocked`, `disputed`, `superseded`. These are functional assessments, not final object states.

## 15. États d’interface
Loading preserves lineage focus; Empty indicates missing trace rather than success; Partial highlights missing links; Error keeps resolved links; Offline read-only; Permission denied masks restricted sources; Stale identifies changed/superseded versions requiring re-review.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Provenance Package | concept/export | reviewer/QA/Evidence/Technique | complete known refs, versions, restrictions and gaps |
| Reproducibility Assessment | concept | QA/Investigation Lead | conditions/gaps, not guaranteed replay |
| Provenance gap/correction proposal | package | source/Studio/Shared owner | no historical mutation |
| Closure activity | event | Shared Trace/Status | human disposition and final Mobile state |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-701..718 | review/closure | CAP-INV-719 | source refs, versions, Tools/Runs, decisions, restrictions | Session/source capability |
| CAP-INV-719 | source/custody gap | Collection/source owner | gap ref, needed correction context | provenance |
| CAP-INV-719 | Tool/Run gap | Studio | Tool/Run refs and missing lineage | provenance |
| CAP-INV-719 | generic trace/version issue | Shared | trace/version refs | provenance |
| CAP-INV-719 | Evidence/Finding/Detection/TI review | destination owner | minimized provenance package | Mobile Session |

## 18. Dépendances
CAP-INV-701..718, Collection/Trust, Studio Tools/Runs, Shared Trace/Activity/Versioning/Export, Security privacy/audit, destination handoffs, OPEN-011/013/014/015.

## 19. Source de vérité
Each source owner remains authoritative for its raw event/object. Mobile provenance links those immutable/stable references and adds analytical semantics. Shared remains canonical for generic trace/version mechanisms; Studio remains canonical for execution records.

## 20. Provenance et audit
This capability itself records actor/reviewer, Session/version, all source refs/versions, tenant/environment/scope, acquisition/custody, assessments, Tools/Runs/parameters/queries, observation and artifact versions, sensitive-access event refs, timeline/correlations, errors/interruptions, handoffs/returns, human decisions, corrections and assessment version.

## 21. Permissions fonctionnelles
Provenance read, restricted-source metadata read, Tool/Run trace read, sensitive-access audit read, reproducibility create/review/dispute, provenance gap prepare, classified export and cross-tenant correlation review. Export never grants underlying source permission.

## 22. Limites et erreurs
Missing historical source, deleted external record, unavailable Tool version, redacted sensitive source, cross-tenant denial, incomplete custody, unsupported parser, failed Run or inconsistent versions results in explicit reproducibility limitation/block, not fabricated lineage.

## 23. Métriques
Lineage completeness, missing source/version/Tool refs, reproducible-with-limitations rate, correction proposals, unresolved provenance gaps, sensitive-access trace coverage, handoff return coverage and export attempts denied/allowed.

## 24. Classification de livraison
`defined` / `planned`. No provenance storage technology, immutable ledger, replay engine, artifact package format, API or code delivered.

## 25. Critères d’acceptation
**Given** a Finding Draft prepared from a recovered fragment **When** provenance is reviewed **Then** source region, recovery Tool/parameters/quality, Derived Artifact, draft version and human disposition are traversable without treating the fragment as complete original Evidence.

**Given** a Tool version no longer available **When** reproducibility is assessed **Then** the limitation is explicit and missing execution is not fabricated.

**Given** no AI **When** provenance is reviewed **Then** deterministic references, version diffs, trace tables and checklists provide full functionality.

## 26. Questions ouvertes
OPEN-011 Mobile implementation delivery; OPEN-013 review/correction mutations; OPEN-014 final material relations; OPEN-015 Tool/Automation/Response run bridging and final trace contracts remain open.

## 27. Consommateurs documentaires
All Mobile capabilities, Evidence/Finding, Collection, Studio, Shared Trace/Versioning, Security/Trust, Quality, Requirements, Roadmap, Technique and future implementation contracts.
