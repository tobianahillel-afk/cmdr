---
id: CAP-GOV-035
title: Decision, Approval and Authority Audit Reconstruction
product: govern
module: audit-trail
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-008, REQ-PROD-015, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-GOV-035 — Decision, Approval and Authority Audit Reconstruction

## 1. Définition
Reconstruire, sans réexécution, la chaîne historique `Action Request → Policy → Authority → Approver Eligibility → Approval Request → Approval → Decision` en conservant versions, acteurs, SoD, délégations, emergency context, conditions, expiration, supersession, liens manquants et contradictions.

## 2. Problème utilisateur
Une Decision peut être historiquement valide tout en ayant traversé plusieurs versions, approvers, exceptions ou délégations. Sans reconstruction sourcée, une revue peut confondre Approval et Decision, ignorer un changement de request/version ou interpréter un événement manquant comme preuve de non-action.

## 3. Objectifs
- pinner request/policy/authority/approval/decision versions ;
- reconstruire acteurs, eligibility/SoD, delegation et emergency path ;
- exposer Policy outcomes/conflicts/exceptions ;
- montrer Decision disposition, conditions, expiry et supersession ;
- signaler links/gaps/contradictions sans fabriquer un fait ;
- fournir une séquence navigable et reviewable.

## 4. Non-objectifs
Ne pas refaire la Decision, modifier Approval/Decision, recalculer authority comme vérité actuelle, réexécuter une policy, conclure fraude/violation, créer Evidence, définir un audit engine ou une state machine finale.

## 5. Propriétaire
Govern / Audit Trail possède la reconstruction métier. GOV-1 reste owner des objets/reviews d’origine ; Shared fournit Trace/Activity/Search ; Settings/Security fournissent identité/permission/authority policy projections.

## 6. Utilisateurs
Principal : Govern Auditor. Secondaires : Decision Reviewer, Security Reviewer, Authority Reviewer, Compliance Reviewer, Decision Maker et source requester avec accès limité.

## 7. Conditions d’entrée
Un Action Request, Approval ou Decision est sélectionné avec suffisamment de références historiques pour démarrer une reconstruction ; les absences/inaccessibilités sont explicitement représentées.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Action Request/version lineage | CAP-GOV-001..006 | audit subject | oui | historical pinned versions | reconstruction partial |
| Policy evaluations/conflicts/exceptions | CAP-GOV-007/008 | governance events | selon applicability | historical versions | explicit gap/unknown |
| authority/eligibility/SoD | CAP-GOV-009/010 | authority context | selon action | event-time context | authority gap |
| Approval Requests/Approvals | CAP-GOV-011..013 | approval history | selon requirement | exact historical records | approval gap |
| Decision/version | CAP-GOV-014..016 | disposition history | oui for Decision reconstruction | exact version | no Decision invented |
| identities/delegations | Settings/Security/GOV-1 | actor context | when available | event-time | masked/unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Action Request | Govern | versions/lineage/requester | audit read |
| Policy/Policy Evaluation | Govern/Security source | outcome/version | audit read |
| Approval | Govern | approver/disposition/authority refs | audit read |
| Decision | Govern | disposition/conditions/expiry | audit read |
| Principal/Role | Settings | historical identity projection | restricted read |
| Govern Audit Event | Govern | event sequence | read/correlate |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Decision Audit Reconstruction | create/update/version | Govern local concept | reconstruction ≠ source mutation |
| missing-link/contradiction candidate | create/review/supersede | Govern local concept | candidate ≠ falsification |
| source objects | no mutation | original owner | historical read only |

## 11. Fonctionnalités
Select reconstruction anchor; traverse request versions; align Policy evaluations; display authority/eligibility/SoD; map Approval Requests to actual Approvals; expose delegation/emergency context; align Decision Draft/final Decision; show conditions/expiry/supersession; identify missing/duplicate/conflicting links; compare two Decisions or versions; return to exact source.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect reconstruction | auditor | reconstruction | 0 | read permission | sourced chain | non |
| compare versions/actors | auditor | reconstruction | 1 | ≥2 refs | deterministic diff | non |
| flag missing/contradictory link | auditor | audit assessment | 2 | rationale/source | attributed candidate | OPEN-013 |
| annotate review | auditor | reconstruction | 2 | write review permission | versioned note | OPEN-013 |
| change Decision/Approval | none | source object | 3/4 | outside audit | prohibited | GOV-1 owner only |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| resolve exact IDs/versions | oui | oui | oui | no authority | relation table |
| compare Policy/Approval/Decision versions | oui | diff | oui | explanation | version diff |
| flag missing links | oui | expected-link rules | oui | candidate summary | completeness matrix |
| summarize rationale/dissent | oui | source aggregation | oui | sourced summary | source excerpts/fields |
| declare improper/fraudulent Decision | human external review | no | no | prohibited | explicit review process |

## 14. États fonctionnels
`building`, `complete-enough-for-review`, `partial`, `authority-gap`, `approval-gap`, `decision-gap`, `contradictory`, `source-unavailable`, `retention-limited`, `disputed`, `reviewed`, `superseded`.

## 15. États d’interface
Loading preserves anchor ; Empty means no visible chain ; Partial lists missing nodes ; Error retains last valid reconstruction ; Offline read-only ; Permission denied masks restricted identity/rationale ; Stale indicates later source supersession or incomplete indexing.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Decision Audit Reconstruction | audit review object | CAP-GOV-037/038/metrics | exact refs/versions/actors retained |
| missing-link candidates | audit findings candidate | CAP-GOV-037 | absence not equated to wrongdoing |
| Decision/Approval comparison | review artifact | auditor/control review | differences sourced |
| review annotation | audit event | Audit Trail | author/rationale preserved |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-034 | decision-chain events available | CAP-GOV-035 | attributed events/refs | Audit Trail |
| CAP-GOV-035 | gaps/contradictions detected | CAP-GOV-037 | expected/observed relations | reconstruction |
| CAP-GOV-035 | review/export requested | CAP-GOV-038 | scoped reconstruction/restrictions | reconstruction |
| CAP-GOV-035 | metric aggregation | CAP-GOV-039..041 | event-derived observations | audit source |

## 18. Dépendances
CAP-GOV-001..016/034/037/038, Shared Trace/Activity/Search/Versioning, Settings identity history, Security SoD/authority/privacy policy, OPEN-007/013/015.

## 19. Source de vérité
GOV-1 objects remain authoritative for their recorded state. The reconstruction is a Govern audit interpretation over immutable/supersedable references; it cannot retroactively alter the request, Approval or Decision.

## 20. Provenance et audit
Store reconstruction anchor, traversed IDs/versions, actor refs, Policy/authority/SoD/Approval/Decision relations, missing/duplicate/conflicting candidates, reviewer annotations, deterministic/AI assistance and every reconstruction version.

## 21. Permissions fonctionnelles
Audit read; restricted actor/authority read; reconstruction create/read; contradiction review; compare Decisions; provenance export preparation. Cross-tenant or sensitive authority data requires explicit scope; final RBAC remains future.

## 22. Limites et erreurs
Historical identity unavailable, retention gap, inaccessible Policy, missing Approval event, duplicated event, conflicting timestamps, superseded Decision and tenant mismatch keep reconstruction partial/disputed. A present audit record does not certify action correctness.

## 23. Métriques
Reconstructions complete-enough; authority/Approval/Decision gap candidates; superseded Decision chains; emergency chains lacking expected retrospective-review reference; average unresolved-link count. No universal quality target.

## 24. Classification de livraison
`defined` / `planned`; no audit/reconstruction engine, storage schema, causal model or compliance certification is delivered.

## 25. Critères d’acceptation
**Given** an Approval is present but no Decision event is visible, **When** reconstruction runs, **Then** Approval remains distinct, Decision is not invented and the gap is explicit.

**Given** requester equals an approver in a chain where SoD applied, **When** the historical chain is reviewed, **Then** eligibility/SoD evidence and actual approver remain independently visible rather than inferred from role alone.

**Given** a superseded Decision, **When** reconstruction opens the latest Decision, **Then** the superseded version and reasons remain resolvable.

**Given** no AI, **When** reconstruction is performed, **Then** ID/version traversal, diffs and human review provide complete functionality.

## 26. Questions ouvertes
OPEN-007/013/015 remain open. Final historical identity retention, technical event model and proof mechanisms remain future; no new OPEN is created.

## 27. Consommateurs documentaires
Audit Trail, CAP-GOV-037/038/039..041, Govern closure, future screens/objects/permissions/technical audit implementation and quality reports.