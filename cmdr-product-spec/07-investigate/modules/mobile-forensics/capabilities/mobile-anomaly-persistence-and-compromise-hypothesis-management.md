---
id: CAP-INV-717
title: Mobile Anomaly, Persistence and Compromise Hypothesis Management
product: investigate
module: mobile-forensics
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-07
requirement_ids: [REQ-INV-001, REQ-PROD-014, REQ-PROD-020, REQ-AI-002, REQ-SEC-001]
open_decisions: [OPEN-005, OPEN-011, OPEN-013, OPEN-014, OPEN-015]
source-of-truth: canonical
---
# CAP-INV-717 — Mobile Anomaly, Persistence and Compromise Hypothesis Management

## 1. Définition
Gérer des Mobile Anomalies et Mobile Hypotheses pour application/configuration/identity/permission/communication/network/storage anomalies, persistence candidates, suspicious-application candidates et compromise hypotheses, avec supporting/contradicting elements, confidence, dispute, withdrawal, supersession et handoff.

## 2. Problème utilisateur
Une anomalie ou application suspecte peut provenir d’une configuration légitime, d’une source partielle ou d’une reconstruction erronée. Sans hypothèses concurrentes et contradictions visibles, l’analyse peut transformer un candidat en compromission ou malware confirmé.

## 3. Objectifs
- créer anomalies/hypotheses explicitement candidates;
- relier support, contradictions, missing data and source quality;
- séparer persistence candidate, suspicious application and compromise assessment;
- permettre review/dispute/withdraw/supersede and bounded handoff.

## 4. Non-objectifs
No malware confirmation, exploit/persistence/evasion procedure, automatic attribution, automatic Finding/Evidence, response, rule deployment, target interaction, vulnerability confirmation, engine or code.

## 5. Propriétaire
Investigate owns Mobile Anomaly/Hypothesis concepts. Evidence/Finding, Detection Engineering, Threat Intelligence and Govern retain their destination lifecycles. Shared owns generic Graph/Timeline/Versioning.

## 6. Utilisateurs
Principal : Investigation Lead. Secondaires : Mobile Forensics Analyst, DFIR Analyst, Evidence Reviewer and SOC Analyst.

## 7. Conditions d’entrée
Session plus sourced observations/correlations, source quality and restrictions, analyst purpose, permission to supporting data and owner assigned to review.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Mobile observations | CAP-INV-706..715 | supporting/contradicting facts | oui | source versions | no anomaly/hypothesis claim |
| Timeline/correlations | CAP-INV-716 | temporal/cross-source candidates | non | current assessment | limited interpretation |
| Integrity/completeness/accessibility | CAP-INV-705 | source limits | oui | reviewed version | confidence constrained |
| External owner context | Detection/TI/Case | contextual candidate knowledge | non | owner version | not invented |
| Missing-data/collection gaps | Intake/Collection | uncertainty | non | current | explicitly unknown |
| Permissions/restrictions | Security/source | access boundary | oui | review time | protected support hidden/blocked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Mobile Observations/Timeline | Investigate | facts/candidates/time | read/link |
| Case Hypothesis / Finding / Evidence | Investigate owner workflows | contextual relation | read/link; no local confirmation |
| Detection/TI context | owner modules | relevant behavior/knowledge | read/link |
| Tool/Run results | Studio | analytical outputs | read, never conclusion |
| Shared Graph/Versioning | Shared | relation/version mechanisms | consume |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Mobile Anomaly | create/review/dispute/withdraw/supersede | Investigate concept | anomaly ≠ compromise/Finding |
| Persistence Candidate | create/review/dispute/withdraw | Investigate concept | no procedure or confirmation |
| Mobile Hypothesis | create/update/review/dispute/withdraw/supersede | Investigate concept | support/contradictions/missing data explicit |
| Handoff package | prepare/version/withdraw | destination owner receives | no destination object auto-created |

## 11. Fonctionnalités
Create anomaly from observation set, classify analytical category, record rationale, support/against/missing evidence, compare alternative hypotheses, link timeline/correlation, confidence with textual basis, assign reviewer, comment/dispute, withdraw/supersede and prepare technical/Evidence/Detection/TI handoffs.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Create anomaly/hypothesis | analyste | candidate | 2 | sourced support | proposed object | OPEN-013 |
| Add support/contradiction | analyste | relation | 2 | source access | versioned relation | non |
| Compare alternatives | analyste/reviewer | hypotheses | 1 | read | comparison matrix | non |
| Dispute/withdraw/supersede | reviewer/owner | candidate | 2 | reason | history preserved | non |
| Prepare handoff | lead | package | 2 | sufficient context | destination package | destination re-evaluates |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Flag rule-based anomaly candidate | oui | oui | oui | suggestion | explicit rules |
| Group supporting/contradicting observations | oui | relation rules | oui | suggestion | hypothesis matrix |
| Propose competing hypothesis | oui | templates/rules | oui | suggestion | manual hypothesis template |
| Summarize gaps/confidence basis | oui | structured fields | oui | summary | checklist/table |
| Confirm compromise/malware/persistence | humain downstream | non | non | interdit | Evidence/Finding/review workflows |

## 14. États fonctionnels
`proposed`, `under-review`, `supported`, `weakly-supported`, `contradicted`, `inconclusive`, `disputed`, `superseded`, `withdrawn`. `supported` remains an analytical assessment and is not qualified Evidence or confirmed Finding/compromise.

## 15. États d’interface
Loading preserves selected hypothesis; Empty suggests creating only from sourced observations; Partial highlights unavailable support; Error preserves valid relations; Offline read-only; Permission denied masks protected support without leaking it; Stale marks superseded sources.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Mobile Anomaly / Hypothesis | concepts | analyst/Case/CAP-INV-718 | sources, alternatives, contradictions, uncertainty |
| Persistence/Suspicious App candidates | concepts | Static/Reverse/Dynamic/TI | candidate only, no technique procedure |
| Gap assessment | result | Collection/Investigation Lead | missing data explicit |
| Handoff package | package | CAP-INV-718 / owner modules | destination ownership preserved |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-706..716 | unusual/correlated observations | CAP-INV-717 | sources, quality, contradictions, gaps | source analysis |
| CAP-INV-717 | technical artifact needed | CAP-INV-718 | selected source refs, question, restrictions | hypothesis |
| CAP-INV-717 | Detection/TI relevance | CAP-INV-718 | hypothesis, support, uncertainty | hypothesis |
| Missing evidence | collection need | CAP-INV-202 | gap, source, purpose, authority | hypothesis |
| Closure/review | provenance | CAP-INV-719 | versions, reviewers, decisions | Session |

## 18. Dépendances
CAP-INV-202, CAP-INV-705..716/718/719, Static/Reverse/Dynamic, Detection Engineering, Threat Intelligence, Shared Graph/Versioning, OPEN-005/011/013/014/015.

## 19. Source de vérité
Mobile anomaly/hypothesis assessments are Investigate. Source observations remain their owner/version. Evidence/Finding qualification, Detection Content, TI knowledge and response remain destination-owned.

## 20. Provenance et audit
Session, candidate ID/version, author/reviewer, source observations/versions, timeline/correlations, rules/Tools/Runs, support/against/missing elements, confidence rationale, comments/disputes, withdrawals/supersession and handoffs.

## 21. Permissions fonctionnelles
Anomaly/Hypothesis create/read/update/review/dispute/withdraw, supporting source read, restricted support read, cross-source relation create, handoff prepare, collection-gap prepare and provenance export.

## 22. Limites et erreurs
Source gaps, partial extraction, inaccessible support, contradictory evidence, stale TI/Detection context, Tool/model failure or permission denial lower confidence or block review; no automatic escalation to compromise/malware/Finding.

## 23. Métriques
Anomalies/hypotheses by state/category, supporting vs contradicting elements, time to review, disputes/withdrawals, downstream handoffs, rejected automated suggestions and hypotheses reopened after new data.

## 24. Classification de livraison
`defined` / `planned`. No anomaly model, malware classifier, persistence detector engine, scoring algorithm, rule deployment, API or code delivered.

## 25. Critères d’acceptation
**Given** a suspicious application and network anomaly **When** a compromise Hypothesis is created **Then** each source, contradiction and missing element is visible and compromise is not automatically confirmed.

**Given** new evidence contradicting a persistence candidate **When** reviewed **Then** the candidate can be disputed/withdrawn while prior versions remain traceable.

**Given** no AI **When** hypotheses are managed **Then** explicit rules, hypothesis matrices, tables and human review cover the full workflow.

## 26. Questions ouvertes
OPEN-005 future engines; OPEN-011 platform/method coverage; OPEN-013 hypothesis mutations; OPEN-014 final relations; OPEN-015 Tool/Run contracts remain open.

## 27. Consommateurs documentaires
Case/Hypotheses, Evidence/Finding, Static/Reverse/Dynamic, Detection Engineering, Threat Intelligence, Collection, Provenance, Objects, Permissions, Quality and Technique.
