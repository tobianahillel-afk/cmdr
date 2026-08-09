---
id: CAP-GOV-037
title: Audit Completeness, Integrity, Gap and Contradiction Assessment
product: govern
module: audit-trail
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-002, REQ-PROD-005, REQ-PROD-008, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-019]
source-of-truth: canonical
---
# CAP-GOV-037 — Audit Completeness, Integrity, Gap and Contradiction Assessment

## 1. Définition
Évaluer si une reconstruction Govern contient suffisamment d’événements et de provenance pour une revue donnée, en comparant classes attendues/observées, gaps, orphans, duplicates, conflicting states/timestamps, actor ambiguity, retention limitations et source availability sans assimiler completeness à vérité ni integrity à preuve cryptographique.

## 2. Problème utilisateur
Un historique peut sembler complet alors que certaines sources sont inaccessibles, retenues moins longtemps ou contradictoires. À l’inverse, une absence ou un doublon peut être bénin. Sans assessment explicite, l’utilisateur risque de sur-interpréter un trou comme fraude ou un ledger présent comme preuve absolue.

## 3. Objectifs
- comparer event classes attendues et observées ;
- identifier missing-event candidates, orphans, version/timestamp/provenance gaps ;
- détecter duplicates, conflicting states et actor ambiguity ;
- distinguer retention/source/access limitations ;
- produire un statut reviewable avec raisons ;
- conserver dispute et uncertainty.

## 4. Non-objectifs
Ne pas prouver cryptographiquement l’intégrité, déclarer falsification/malice, réparer un historique, inventer un event, supprimer un doublon, choisir retention/storage, modifier un objet source ou créer une alerte SIEM.

## 5. Propriétaire
Govern / Audit Trail owns Govern completeness/gap/contradiction assessment. Security owns integrity requirements; Settings owns retention/storage; Shared owns audit plumbing; source owners retain source records.

## 6. Utilisateurs
Principal : Govern Auditor. Secondaires : Security/Compliance Reviewer, Govern Reviewer, Decision/Response Reviewer, Platform Administrator as source owner and Investigate reviewer on explicit handoff.

## 7. Conditions d’entrée
At least one reconstruction or scoped audit event set exists, along with an expected-event-class profile derived from relevant GOV-1/GOV-2 lifecycle semantics and source-availability metadata.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| expected event classes | CAP-GOV-034 + capability semantics | completeness expectation | oui | versioned review profile | assessment impossible |
| observed audit events | CAP-GOV-034..036 | observed records | oui | scoped snapshot | empty/partial explicit |
| source availability | source owners/Settings | limitation context | oui | current/review-time | unknown limitation |
| retention/legal-hold context | Settings/Security | availability constraint | when applicable | policy/version shown | retention unknown |
| actor/object/version refs | audit events | integrity/context | when available | historical | ambiguity/gap |
| duplicate/conflict signals | deterministic comparison | consistency context | non | snapshot | none invented |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Govern Audit Event | Govern | classes/refs/limitations | read/assess |
| Decision/Run Reconstructions | Govern | expected/observed chain | read |
| retention config | Settings | retention/hold context | restricted read |
| audit/integrity policy | Security | requirements only | read |
| source event refs | original owner | existence/version metadata | restricted read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Audit Completeness Assessment | create/update/version/supersede | Govern local concept | scoped purpose and expectations required |
| Audit Gap | create/review/resolve/supersede | Govern local concept | gap ≠ malicious/failure claim |
| Audit Contradiction | create/review/dispute/supersede | Govern local concept | contradiction ≠ falsity automatically |
| source records | no mutation | source owner | assessment only |

## 11. Fonctionnalités
Define review scope; derive expected classes; compare observed classes; detect missing-event candidates/orphan refs/version gaps/timestamp gaps/actor ambiguity/duplicates/conflicting states/provenance gaps; attach source/retention limitations; distinguish inaccessible from absent; allow dispute/review; compare assessments across versions; hand off unresolved issues.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect expected/observed matrix | auditor | assessment | 0 | read | scoped matrix | non |
| run deterministic completeness/consistency check | auditor | assessment | 1 | profile+snapshot | candidate gaps/conflicts | non |
| create/review gap or contradiction | auditor | Gap/Contradiction | 2 | source/rationale | versioned assessment | OPEN-013 |
| dispute/supersede assessment | authorized reviewer | assessment | 2 | reason | preserved history | OPEN-013 |
| delete/fabricate history | none | source/audit | 4 | forbidden | no action | prohibited |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| expected-vs-observed comparison | oui | oui | oui | explain | matrix |
| duplicate/orphan detection | oui | IDs/rules | oui | candidate explanation | deterministic list |
| contradiction candidate detection | oui | explicit comparisons | oui | candidate only | diff table |
| summarize limitations | oui | aggregation | oui | sourced summary | limitation table |
| declare fraud/tampering | human external process | no | no | prohibited | evidence/review escalation |

## 14. États fonctionnels
`not-assessed`, `complete-enough-for-review`, `partial`, `gap-detected`, `contradictory`, `source-unavailable`, `retention-limited`, `integrity-unverified`, `disputed`, `reviewed`, `superseded`.

## 15. États d’interface
Loading preserves review scope ; Empty distinguishes no expected events from no visible events ; Partial lists missing sources ; Error keeps prior assessment ; Offline blocks fresh source checks ; Permission denied hides restricted records but counts limitation ; Stale marks snapshot/date drift.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Completeness Assessment | audit review object | CAP-GOV-038/046/047 | review scope/expectations/limitations explicit |
| Audit Gap | candidate issue | reviewer/improvement | no malicious conclusion |
| Audit Contradiction | candidate issue | reviewer/improvement | conflicting sources preserved |
| source/retention limitation | assessment evidence | reporting/closure | absence reason distinguished |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-035/036 | reconstruction ready | CAP-GOV-037 | expected/observed chains | reconstruction |
| CAP-GOV-037 | review/package needed | CAP-GOV-038 | assessment/gaps/restrictions | assessment |
| CAP-GOV-037 | trend/control review | CAP-GOV-046 | gap/contradiction observations | assessment |
| CAP-GOV-037 | improvement candidate | CAP-GOV-047 | unresolved issue/source | assessment |

## 18. Dépendances
CAP-GOV-034..036/038/046/047, Settings retention/storage, Security audit/integrity/legal-hold/privacy, Shared Trace/Activity/Search, source availability, OPEN-008/013/019.

## 19. Source de vérité
Govern owns the assessment, not source truth. Source owners and retention/security sources remain authoritative for their records/configuration. Completeness is scoped to expected data and available access, never global truth completeness.

## 20. Provenance et audit
Record review purpose/scope, expected profile version, snapshot time, observed refs, source availability, retention context, deterministic rules, gap/conflict classifications, reviewers, disputes, AI candidates and supersession.

## 21. Permissions fonctionnelles
Gap assessment create/review; contradiction review; restricted audit read; retention metadata read; actor identity read where authorized; cross-tenant comparison only with explicit permission; provenance export preparation. No final RBAC/ABAC.

## 22. Limites et erreurs
Incomplete indexes, retention expiry, legal hold restrictions, permission masking, clock skew, duplicated sources, unresolved identities and unavailable providers can prevent a complete assessment. An integrity requirement does not prove implementation of immutable/cryptographic storage.

## 23. Métriques
Assessments by outcome; gaps by class/source; retention-limited reviews; disputed contradictions; orphan/duplicate candidate counts; inaccessible-source frequency. Counts are quality signals, not conclusions.

## 24. Classification de livraison
`defined` / `planned`; no audit-integrity engine, anomaly detector, warehouse, storage or proof mechanism selected.

## 25. Critères d’acceptation
**Given** source data is unavailable, **When** completeness is assessed, **Then** outcome is partial/source-unavailable rather than complete or fabricated.

**Given** retention removed an expected event class, **When** review runs, **Then** the gap is marked retention-limited and absence is not treated as proof that the action did not occur.

**Given** duplicate records disagree on state, **When** contradiction review runs, **Then** both records remain visible and falsity is not automatically assigned.

**Given** no AI, **When** assessment is performed, **Then** expected/observed matrices and deterministic comparisons provide full functionality.

## 26. Questions ouvertes
OPEN-008/013/019 remain open. Technical integrity proof, retention implementation and external-disclosure policy remain future. No new OPEN.

## 27. Consommateurs documentaires
Audit Trail, CAP-GOV-038/046/047, Govern closure reports, Security/Settings future phases, future screens/objects/permissions/technical implementation and quality gates.