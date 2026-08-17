---
id: CAP-GOV-034
title: Govern Audit Trail Intake and Event Semantics
product: govern
module: audit-trail
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-002, REQ-PROD-004, REQ-PROD-005, REQ-PROD-008, REQ-PROD-015, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-015, OPEN-019]
source-of-truth: canonical
---
# CAP-GOV-034 — Govern Audit Trail Intake and Event Semantics

## 1. Définition
Interpréter et corréler, pour le domaine Govern, les événements et changements relatifs à Action Request, Policy, authority, Approval, Decision, Playbook, Response Run, verification, rollback/recovery et Result afin de former des Govern Audit Events attribués, versionnés et permission-aware, sans créer un format technique d’événement ni un moteur de collecte.

## 2. Problème utilisateur
Les faits de gouvernance proviennent de plusieurs objets et runtimes. Sans sémantique Govern commune, un auditeur peut confondre un log technique, une Activity générique, un changement métier et une conclusion d’audit, perdre l’acteur ou la version exacte, ou traiter l’absence d’un event comme preuve qu’une action n’a jamais eu lieu.

## 3. Objectifs
- définir les classes d’événements Govern attendues sur GOV-1/GOV-2 ;
- préserver source, owner, actor, action, object/version, tenant/environment, timestamps et correlation identifiers ;
- distinguer event source, Govern Audit Event interprété, Trace et Activity ;
- conserver raisons, before/after references, limitations et disponibilité source ;
- rendre les événements exploitables par reconstruction, completeness review et metrics ;
- signaler explicitement unknown/partial/inaccessible au lieu d’inventer des faits.

## 4. Non-objectifs
Ne pas créer SIEM, log pipeline, event bus, audit storage, immutable ledger technique, cryptographic proof, event schema/API, ingestion protocol, retention policy, Evidence Investigate ou conclusion de fraude/violation.

## 5. Propriétaire
Govern / Audit Trail possède la sémantique métier du Govern Audit Event et sa relation au lifecycle Govern. Shared possède Trace/Activity/audit plumbing; source products et runtimes restent propriétaires des faits sources; Settings/Security possèdent retention, stockage, intégrité et accès.

## 6. Utilisateurs
Principal : Govern Auditor. Secondaires : Govern Reviewer, Decision Maker, Response Operator, Security Reviewer, Compliance Reviewer, Incident Commander, Investigator et product owners autorisés qui inspectent une projection.

## 7. Conditions d’entrée
Une source Govern ou source propriétaire expose un event/changement ou une absence explicite de source, avec tenant/environment et permission suffisants. Les références d’objet doivent rester résolubles ou être marquées indisponibles/retention-limited.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| object/change event | GOV-1/GOV-2 or source owner | lifecycle fact candidate | oui pour event interprété | source timestamp/version | no event invented |
| actor/service identity | source + Settings/Security projection | attribution | lorsque disponible | event-time context | actor-unknown |
| object/version reference | canonical owner | subject/version | oui | historical version | orphan/gap candidate |
| tenant/environment | Settings/source context | isolation scope | oui | event-time context | block cross-scope interpretation |
| timestamps | source | event/record time | oui | source supplied | ordering unknown |
| correlation/request/run ids | source/Shared Trace | linkage | selon lifecycle | source supplied | reconstruction partial |
| rationale/result/before-after refs | GOV source | explanatory context | lorsque applicable | event version | explicitly unavailable |
| source availability/retention state | source/Settings | audit limitation | oui | current knowledge | unknown limitation |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Action Request / Approval / Decision | Govern | lifecycle versions/actions | read/audit |
| Playbook / Response Run / Result | Govern | execution lifecycle refs | read/audit |
| Policy / authority context | Govern/Security source | governance outcomes | restricted read |
| Workflow/Automation Run/Tool Call | Studio | correlated technical provenance | read/link only |
| local-audit-event / technical output | Endpoint/source owner | raw technical reference | restricted read/link |
| Trace / Activity | Shared | correlation/presentation mechanism | consume only |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Govern Audit Event | create/interpret/version/supersede | Govern local concept | must cite source/ref/version; no raw-event ownership transfer |
| audit correlation relation | create/update | Govern local concept | correlation ≠ causation |
| source object/event | no mutation | original owner | audit is read/reconstruction only |
| Trace/Activity | no redefinition | Shared | generic infrastructure remains Shared |

## 11. Fonctionnalités
Map source events to Govern lifecycle class; attach exact object/version and actor; preserve tenant/environment; distinguish event time/record time; link correlation IDs; capture action/result/rationale/before-after refs; expose source availability and masking; classify unknown/orphan candidate; compare duplicate candidates without deleting either; feed reconstruction and completeness review.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect source/audit event | auditor | Govern Audit Event/source ref | 0 | read permission | sourced event view | non |
| correlate/compare events | auditor | audit relations | 1 | visible refs | bounded correlation | non |
| interpret/annotate event | auditor | Govern Audit Event | 2 | source retained + rationale | versioned interpretation | OPEN-013 |
| mark source unavailable/disputed | auditor | audit event | 2 | evidence of limitation | limitation preserved | OPEN-013 |
| delete/alter source history | none | source/audit history | 4 | forbidden | no action | prohibited |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| classify known event classes | oui | mapping/rules | oui | explanation | event-class table |
| correlate exact IDs/versions | oui | oui | oui | candidate links only | deterministic joins |
| flag missing attribution/ref | oui | oui | oui | summary | completeness checklist |
| summarize event sequence | oui | ordering rules | oui | sourced summary | timeline/table |
| invent missing event/declare fraud | no | no | no | prohibited | explicit unknown/review |

## 14. États fonctionnels
`observed`, `interpreted`, `partially-attributed`, `actor-unknown`, `object-unresolved`, `source-unavailable`, `retention-limited`, `duplicate-candidate`, `contradictory-candidate`, `disputed`, `superseded`. Ces états décrivent l’interprétation d’audit, pas l’état de l’objet source.

## 15. États d’interface
Loading conserve filtres et contexte ; Empty signifie aucun event visible pour le scope courant, pas absence d’action ; Partial expose sources manquantes ; Error garde les dernières données valides ; Offline limite les nouvelles résolutions ; Permission denied masque contenu/identité protégés ; Stale expose la date de dernière source connue.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Govern Audit Event | interpreted audit record | CAP-GOV-035..038/metrics | source/version/actor/limitations attributable |
| correlation relation | audit linkage | reconstruction | IDs and rationale explicit |
| source limitation | audit event state | completeness review | absence not converted to fact |
| event-class inventory | audit semantics | GOV-3 metrics/closure | expected source classes named |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| GOV-1 lifecycle event | event available | CAP-GOV-034 | request/policy/approval/decision ref/version | source object |
| GOV-2 lifecycle event | event available | CAP-GOV-034 | run/step/verification/rollback/result refs | source object |
| CAP-GOV-034 | decision-chain event set | CAP-GOV-035 | attributed events/limitations | Audit Trail |
| CAP-GOV-034 | run-chain event set | CAP-GOV-036 | attributed events/limitations | Audit Trail |
| CAP-GOV-034 | expected vs observed inventory | CAP-GOV-037 | classes/sources/availability | Audit Trail |

## 18. Dépendances
CAP-GOV-001..033, canonical objects, Shared Trace/Activity/Linking/Search, Settings tenant/retention sources, Security audit/privacy/export policy, source runtime records, OPEN-007/008/013/015/019.

## 19. Source de vérité
Source objects/events remain authoritative for what they recorded. Govern owns only its audit interpretation and correlation context. Shared Trace/Activity remain mechanisms. Absence, redaction or retention loss remains an explicit limitation rather than a fabricated event.

## 20. Provenance et audit
Record interpretation author/method, source/ref/version, actor/service identity ref, action, before/after refs, timestamps, tenant/environment, correlation IDs, rationale/result, masking, source availability, deterministic/AI assistance and supersession history.

## 21. Permissions fonctionnelles
Audit Trail read; restricted audit read; actor identity read; audit search; Govern Audit Event interpretation/annotation; cross-tenant access only when separately authorized; provenance export preparation. Final atomic permissions remain Security-owned future work.

## 22. Limites et erreurs
Missing source, retention gap, duplicate IDs, conflicting timestamps, unresolved actor, inaccessible object, tenant mismatch, stale correlation data and source disagreement remain visible. Timestamp order never proves causality; integrity is not claimed cryptographically without implementation evidence.

## 23. Métriques
Conceptual quality metrics: events by class/source; unresolved actors; orphan candidates; source-unavailable/retention-limited rates; duplicate/contradiction candidates; percentage carrying object/version/correlation metadata. No universal quality threshold is set.

## 24. Classification de livraison
`defined` / `planned`. Current evidence is documentary semantics only; no audit pipeline, SIEM, storage, immutable ledger, cryptographic proof or event format is delivered. Promotion requires future object/permission/technical/release evidence.

## 25. Critères d’acceptation
**Given** an Approval exists but its Decision event is absent, **When** audit events are interpreted, **Then** the Approval remains visible, the Decision event is not invented and a reconstruction gap can be assessed.

**Given** two source records have conflicting timestamps, **When** they are correlated, **Then** both timestamps/sources remain visible and chronology is not presented as proven causality.

**Given** the actor identity is permission-restricted, **When** an auditor lacks access, **Then** the event remains structurally visible with masking and no identity leakage.

**Given** no AI provider, **When** events are classified/correlated, **Then** mappings, IDs, tables and human review provide the complete essential function.

## 26. Questions ouvertes
OPEN-008 covers actual source/runtime support; OPEN-013 covers governance of reversible audit annotations where applicable; OPEN-015 covers Automation Run/Response Run linkage; OPEN-019 covers external dissemination. No new OPEN is required. Final event/storage/retention/integrity implementation remains future.

## 27. Consommateurs documentaires
Audit Trail, CAP-GOV-035..038, CAP-GOV-039..047, Govern closure reports, future object/permission/screen/technical phases, Shared Trace/Reporting/Export integrations and Delivery Roadmap validation.