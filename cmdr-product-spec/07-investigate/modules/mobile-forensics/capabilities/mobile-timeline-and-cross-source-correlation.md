---
id: CAP-INV-716
title: Mobile Timeline and Cross-Source Correlation
product: investigate
module: mobile-forensics
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-07
requirement_ids: [REQ-INV-001, REQ-PROD-014, REQ-PROD-019, REQ-PROD-020, REQ-AI-002, REQ-SEC-001]
open_decisions: [OPEN-011, OPEN-013, OPEN-014, OPEN-015]
source-of-truth: canonical
---
# CAP-INV-716 — Mobile Timeline and Cross-Source Correlation

## 1. Définition
Construire une Mobile Timeline et des Cross-Source Correlation Candidates en distinguant observed, recorded, reconstructed, estimated, synchronized and absent timestamps, timezone, source, confidence, gaps and conflicts, puis corréler avec Case Timeline, Cloud, Endpoint, Network, Disk, Memory, Detection Engineering, Threat Intelligence et plusieurs extractions/backups.

## 2. Problème utilisateur
Les artefacts mobiles utilisent plusieurs horloges et copies. Un timestamp synchronisé n’est pas un événement local certain; deux événements proches ne prouvent pas une causalité. Sans qualité temporelle, une timeline peut transformer une reconstruction en vérité.

## 3. Objectifs
- normaliser l’affichage sans perdre timestamp original/source/timezone;
- produire chronology and correlation candidates with support/conflicts;
- compare multiple packages/backups and owner-domain observations;
- preserve Mobile Timeline ≠ Case Timeline and correlation ≠ causality.

## 4. Non-objectifs
No final time-normalization algorithm, clock correction engine, causality inference, Case Timeline ownership, active query of external systems, provider-specific format, automatic Finding or response.

## 5. Propriétaire
Investigate owns Mobile Timeline semantics and correlation candidates. Shared owns generic Timeline/Graph/Linking. Case owns Case context; Cloud/Network/Disk/Memory/Detection/TI retain their source observations and workflows.

## 6. Utilisateurs
Principal : DFIR Analyst. Secondaires : Mobile Forensics Analyst, Investigation Lead, Evidence Reviewer and SOC Analyst.

## 7. Conditions d’entrée
Session, sourced observations from one or more Mobile capabilities, timestamp/source metadata where available, scope/permissions for every linked source and integrity/completeness limits.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Mobile observations | CAP-INV-706..715 | sourced events/records | oui | source versions | no timeline entry invented |
| Timestamp/timezone/quality | source/observations | temporal metadata | non | source version | timestamp `absent`/unknown |
| Multiple packages/backups | Session/CAP-INV-714 | comparison sources | non | selected versions | single-source timeline |
| Case/Cloud/Endpoint/Network/Disk/Memory observations | owner modules | cross-source context | non | owner versions | no external correlation claim |
| Detection/TI context | owner modules | analytical context | non | versioned | absent, not invented |
| Permissions/restrictions | Security/source owners | cross-source boundary | oui | correlation time | source excluded/masked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Mobile Observations / Session | Investigate | events, sources, scope | read |
| Case Timeline / Case | Investigate/Shared mechanism | contextual chronology | read/link, not modify as Mobile truth |
| Cloud/Network/Disk/Memory observations | owner modules | candidate correlated context | read/link by permission |
| Detection/TI objects | owner modules | hypothesis/intelligence context | read/link |
| Shared Timeline/Graph/Linking | Shared | rendering/relation mechanisms | consume |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Mobile Timeline entry semantics | create/reconstruct/review/supersede | Investigate using Shared | original time/source/quality retained |
| Cross-source Correlation Candidate | create/review/dispute/supersede | Investigate concept | correlation ≠ causality |
| Time quality assessment | create/update | Investigate | observed/recorded/reconstructed/estimated/synchronized/absent explicit |
| Case/owner source object | none | owner | relation only, no source mutation |

## 11. Fonctionnalités
Timeline table/visualization, source/time-quality filters, timezone view conversion with originals preserved, gaps/conflict highlighting, multi-backup comparison, cross-domain pivots, candidate relation graph with tabular alternative, provenance drill-down, contradiction annotations and selected correlation handoffs.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Browse/filter timeline | analyste | timeline | 0 | read | scoped chronology | non |
| Reconstruct/normalize view | analyste | time projection | 1 | sources | attributed derived view | non |
| Create/dispute correlation | analyste | candidate | 2 | source permissions | versioned candidate | OPEN-013 |
| Compare backups/sources | analyste | observations | 1 | read all | differences/gaps | non |
| Prepare hypothesis/handoff | analyste | selected correlation | 2 | support/context | CAP-INV-717/718 input | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Normalize display timezone | oui | oui | oui | non | deterministic conversion |
| Sort/group events | oui | oui | oui | non | timeline/table |
| Detect temporal proximity/candidate relation | oui | rules | oui | suggestion | explicit correlation rules |
| Summarize conflicts/gaps | oui | structured | oui | yes | conflict table |
| Infer causality/compromise | non | non | non | interdit | human Hypothesis workflow |

## 14. États fonctionnels
Correlation/entries use `proposed`, `under-review`, `supported`, `weakly-supported`, `contradicted`, `inconclusive`, `disputed`, `superseded`, `withdrawn`. Time quality separately uses `observed`, `recorded`, `reconstructed`, `estimated`, `synchronized`, `absent`.

## 15. États d’interface
Loading preserves time window/source filters; Empty means no in-scope timed observations; Partial shows missing sources/timestamps; Error keeps valid timeline entries; Offline read-only; Permission denied hides restricted source details; Stale marks source/version age.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Mobile Timeline | analytical chronology | analyst/CAP-INV-717/718 | source/time quality visible; not Case Timeline |
| Correlation Candidates | concepts | CAP-INV-717/718 | support/conflicts/permissions; no causality claim |
| Gap/conflict assessment | result | analyst/Collection feedback | missing/contradicting sources explicit |
| Timeline provenance | trace context | CAP-INV-719 | original + derived times retained |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-706..715 | timed observations | CAP-INV-716 | records, source, time quality, restrictions | source analysis |
| Case/Cloud/Network/Disk/Memory | authorized correlation | CAP-INV-716 | owner refs + selected context | owner source |
| CAP-INV-716 | interpretation needed | CAP-INV-717 | correlations, gaps, conflicts, support | timeline |
| CAP-INV-716 | Evidence/Finding/Detection/TI handoff | CAP-INV-718 | selected events/correlations/provenance | timeline |
| Closure | reproducibility | CAP-INV-719 | source versions, transformations, corrections | Session |

## 18. Dépendances
CAP-INV-706..715/717..719, Case Timeline, Cloud, Network, Disk, Memory, Detection, TI, Shared Timeline/Graph/Linking, Security, OPEN-011/013/014/015.

## 19. Source de vérité
Original timestamps/records remain source-owner facts. Mobile owns time-quality assessments and correlation candidates. Shared owns generic Timeline/Graph mechanisms. Case Timeline remains Case-owned semantic context.

## 20. Provenance et audit
Source record/version, original timestamp/timezone, normalized/reconstructed value and method, time quality, parser/Tool/Run, correlation rule/parameters, linked owner objects, permissions, conflicts/gaps, reviewer/dispute and supersession.

## 21. Permissions fonctionnelles
Timeline read, cross-source source read, time reconstruction/correlation create, cross-tenant correlation, restricted event read, correlation dispute, export and provenance read. Correlation permission never grants source permission.

## 22. Limites et erreurs
Missing/ambiguous timezone, clock drift, duplicated synchronized records, absent timestamps, conflicting sources, stale owner object, permission denial or partial extraction remains visible; no gap is filled by fabricated time.

## 23. Métriques
Entries by time-quality state, missing timestamps, conflicts, cross-source candidate count, disputed correlations, source-permission blocks, multi-backup comparisons and accepted/modified/rejected suggestions.

## 24. Classification de livraison
`defined` / `planned`. No time engine, causal model, external query, provider format, API or code delivered.

## 25. Critères d’acceptation
**Given** one synchronized timestamp and one device-recorded timestamp **When** correlated **Then** both qualities remain distinct and the correlation does not imply causality.

**Given** a Case Timeline **When** Mobile events are linked **Then** Mobile Timeline remains a separate analytical context and no Case entry is silently rewritten.

**Given** no AI **When** timeline is built **Then** deterministic sorting/time conversion, filters, rules and human review provide full functionality.

## 26. Questions ouvertes
OPEN-011 source/platform timing semantics; OPEN-013 correlation mutations; OPEN-014 final event relations; OPEN-015 Tool/Run/trace contracts remain open.

## 27. Consommateurs documentaires
Mobile Hypotheses, Artifact/Evidence handoff, Case Timeline, Cloud/Network/Disk/Memory, Detection/TI, Shared Timeline/Graph, Objects, Permissions, Quality and Technique.
