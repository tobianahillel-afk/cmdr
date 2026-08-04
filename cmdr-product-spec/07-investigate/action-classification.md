---
id: investigate-action-classification
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-003
  - REQ-PROD-004
  - REQ-SEC-001
  - REQ-SEC-002
open_decisions:
  - OPEN-007
  - OPEN-013
---
# Action classification — Investigate through Phase 4B.2A

## Classes
| Classe | Sens fonctionnel | Autorité Investigate |
|---:|---|---|
| 0 | observation, recherche, inspection, filtre, pivot ou comparaison | selon permission |
| 1 | collecte ou export contrôlé et borné | selon capability, policy et impact |
| 2 | modification réversible, session, transfert ou interruption | OPEN-013; step-up/Govern selon policy |
| 3 | containment | préparation seulement; Govern obligatoire |
| 4 | destructive ou irréversible | contexte seulement; Govern obligatoire |

## Phase 4B.1
Les classifications CAP-INV-001..114 restent celles de leurs specifications et du registre antérieur.

## Phase 4B.2A
| Capability | Action | Classe | Acteur | Cible | Risque | Govern | OPEN |
|---|---|---:|---|---|---|---|---|
| CAP-INV-201 | read/filter/link Endpoint context | 0/2 | Analyst | relation Case–Endpoint | wrong/stale target | policy dependent | OPEN-008/013 |
| CAP-INV-202 | create/modify/submit Collection Request | 1/2 | Analyst | request/scope | unbounded or disruptive collection | by class/policy | OPEN-008/013 |
| CAP-INV-203 | cancel/retry Collection Job | 1/2 | Analyst | job/elements | duplicate effects or lost partials | by class | OPEN-008/013 |
| CAP-INV-204 | request bounded triage package | 1 | Analyst | endpoint/categories | resource impact | by profile/policy | OPEN-008/013 |
| CAP-INV-205 | acquire files/directories | 1/2 | DFIR Analyst | bounded paths/patterns | sensitive/unbounded data | by impact | OPEN-008/013 |
| CAP-INV-206 | inspect process/system snapshot | 0/1 | Analyst | host snapshot | staleness, sensitive data | no normally | OPEN-008 |
| CAP-INV-207 | request memory acquisition | 1/2 | DFIR Analyst | Endpoint | high resource/size/platform | by impact | OPEN-005/008/013 |
| CAP-INV-208 | request/stop network capture | 1/2 | DFIR Analyst | bounded capture | privacy/resource/loss | by impact | OPEN-008/013 |
| CAP-INV-209 | request/open/extend/close Live Session | 2 | Response Operator | session | persistent access/session conflict | policy/gate possible | OPEN-007/008/013/015 |
| CAP-INV-210 | execute/interrupt endpoint operation | 0..4 by effect | Response Operator | Endpoint Operation | mutation/containment/destruction | mandatory class 3/4 | OPEN-007/008/013/015 |
| CAP-INV-211 | upload/download/cancel transfer | 1/2 | Response Operator | file/Endpoint | collision, sensitive data, deployment misuse | by class | OPEN-008/013 |
| CAP-INV-212 | verify/dispute/link Operation Result | 0/2 | Analyst/Reviewer | result/relations | local output mistaken for Govern Result | no normally | OPEN-013/015 |
| CAP-INV-213 | annotate/verify/dispute custody | 0/2 | Evidence Reviewer | custody events | history loss or false integrity | separation possible | OPEN-013/014 |
| CAP-INV-214 | annotate/correct/export provenance | 0/1/2 | Reviewer/Auditor | provenance chain | source attribution/redaction | by sensitivity | OPEN-007/015 |
| CAP-INV-215 | prepare containment package | 2/3/4 | Investigation Lead | Action Request content | unauthorized containment | mandatory for 3/4 | OPEN-007/008/013/015 |

## Invariants
No class 3/4 execution occurs in Investigate. No automation self-grants permission. No command syntax or implementation technique is specified. Artifact qualification and historical trace are never silently rewritten.