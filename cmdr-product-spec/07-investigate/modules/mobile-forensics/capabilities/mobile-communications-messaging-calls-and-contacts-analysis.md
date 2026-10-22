---
id: CAP-INV-709
title: Mobile Communications, Messaging, Calls and Contacts Analysis
product: investigate
module: mobile-forensics
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-07
requirement_ids: [REQ-INV-001, REQ-PROD-014, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-011, OPEN-013, OPEN-014, OPEN-015]
source-of-truth: canonical
---
# CAP-INV-709 — Mobile Communications, Messaging, Calls and Contacts Analysis

## 1. Définition
Analyser des messaging records, conversation/participant candidates, call records, voicemail metadata disponible, contacts/address books, stored emails/communications, attachments, timestamps, declared delivery/read states, deleted candidates, account/application context, source, contradictions et restrictions.

## 2. Problème utilisateur
Les communications mobiles contiennent des données privées de multiples personnes. Un message stocké ne prouve pas son auteur, un message reçu ne prouve pas sa lecture, un call record ne prouve pas le contenu d’une conversation et un contact ne prouve pas une relation personnelle.

## 3. Objectifs
- présenter records et participants candidats avec source/app/account/time context;
- distinguer sent/received/delivery/read states comme déclarés par la source;
- préserver deleted candidates, contradictions and privacy restrictions;
- permettre search/filter/compare and source-aware timeline/handoff without automatic attribution.

## 4. Non-objectifs
No account login, message sending, remote lookup, interception, decryption bypass, contact-person identity merge, authorship inference, conversation-content inference from call logs, legal conclusion, API/protocol or provider integration.

## 5. Propriétaire
Investigate owns Communication/Call/Contact Observations. Raw app records remain source material; Shared owns generic graph/search/timeline; Security controls private access; TI/Detection receive only explicit packages.

## 6. Utilisateurs
Principal : Mobile Forensics Analyst. Secondaires : Sensitive Data Reviewer, DFIR Analyst, Evidence Reviewer, Investigation Lead et SOC Analyst.

## 7. Conditions d’entrée
Session, authorized application/communication records, source/app/account context, integrity/completeness limits, private-data permission and explicit analytical purpose.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Messaging/communication records | package/app data | source records | oui pour message analysis | source version | no communication claim |
| Call/contact/address-book records | package/app data | source records | selon scope | source version | unknown/not represented |
| App/account context | CAP-INV-707/708/712 | context candidates | non | Session version | attribution more limited |
| Timestamps/status fields | source | time/state declarations | non | source version | unknown, not inferred |
| Integrity/completeness limits | CAP-INV-705 | coverage boundary | oui | reviewed version | partial/unverified |
| Privacy/permission/classification | Security/source | access boundary | oui | access time | masked/denied |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Session/App Data Observation | Investigate | source/app/scope | read |
| Communication/call/contact source records | source | content/metadata/status | read by permission |
| Account Observation | CAP-INV-712 concept | account context | read/link candidate |
| Entity/Graph/Timeline | Shared | candidate relation/time mechanisms | consume, no identity merge |
| Artifact/Evidence/Finding | Investigate | attachments/handoff context | read/link only |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Communication Observation | create/review/dispute/supersede | Investigate concept | record/source/app context required |
| Call Observation | create/review/dispute/supersede | Investigate concept | call record ≠ conversation content |
| Contact Observation | create/review/dispute/supersede | Investigate concept | contact/number ≠ certain person/relationship |
| Participant/conversation relation candidate | create/dispute | Investigate | source/confidence/contradiction required |

## 11. Fonctionnalités
Search/filter conversations/records, inspect participants candidates and attachments, compare status/timestamps across sources, display delivery/read flags exactly as source declarations, mask private content, link contacts/accounts as candidates, identify deleted candidates and prepare Timeline/Artifact/Evidence/Finding/TI handoffs.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Browse metadata/content | analyste/reviewer | records | 0 | proper read permission | scoped view | non |
| Reveal restricted content | sensitive reviewer | protected record | 0/2 | explicit permission/purpose | audited access | step-up possible |
| Search/compare records | analyste | observations | 1 | source permissions | result/diff | non |
| Annotate/dispute participant/state | analyste | observation | 2 | source | versioned disposition | OPEN-013 |
| Prepare attachment/record handoff | analyste | source selection | 1/2 | purpose/restrictions | CAP-INV-718 package input | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Parse/normalize declared records | oui | oui where supported | oui | non | deterministic parser/table |
| Search/filter/compare | oui | oui | oui | non necessary | query/diff |
| Group conversation candidates | oui | relation rules | oui | suggestion | participant/time grouping |
| Summarize authorized conversation | oui | structured counts | oui | yes with permission | viewer/exported summary |
| Attribute message/person/intent | non | non | non | interdit | human sourced assessment only |

## 14. États fonctionnels
Observations: `proposed`, `under-review`, `supported`, `weakly-supported`, `contradicted`, `inconclusive`, `disputed`, `superseded`, `withdrawn`. Deleted/read/delivery fields remain source-declared candidates, not certainty beyond source semantics.

## 15. États d’interface
Loading preserves conversation/filter; Empty means no represented records; Partial shows missing sources/participants/status fields; Error keeps valid records; Offline read-only; Permission denied hides content and identifiers as required; Stale shows package/app version.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Communication/Call/Contact Observations | concepts | CAP-INV-716/717/718 | source/app/time/status/limits |
| Participant/conversation candidates | relations | Shared Graph/analyst | no certain person/authorship |
| Attachment/record selections | package inputs | CAP-INV-718 | privacy/classification preserved |
| Deleted candidates | candidates | CAP-INV-715 | deletion ≠ user intent |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-708 | communication records | CAP-INV-709 | app/source/records/restrictions | app data |
| CAP-INV-709 | temporal correlation | CAP-INV-716 | observations/status/timestamps/quality | communications |
| CAP-INV-709 | anomaly/hypothesis | CAP-INV-717 | patterns, support, contradictions | communications |
| CAP-INV-709 | selected attachment/record | CAP-INV-718 | source, app, record, privacy limits | communications |
| Deleted record candidate | recovery review | CAP-INV-715 | source record/state | communications |

## 18. Dépendances
CAP-INV-705/708/712/715..719, Shared Graph/Search/Timeline, Security privacy, Evidence/Artifact workflows, OPEN-011/013/014/015.

## 19. Source de vérité
Raw communication/status/content fields remain source records. Mobile observations/relations are Investigate assessments. Entity identity and generic graph mechanisms remain Shared/owner controlled. No source field is silently upgraded to person/authorship truth.

## 20. Provenance et audit
Package/app/account source, record IDs, timestamps/timezone, status fields, content access level, participant candidates, attachments, deleted indicators, parsers/Tools/Runs, queries, comparisons, sensitive accesses, annotations, disputes, handoffs and human dispositions.

## 21. Permissions fonctionnelles
Communications metadata/content read separately, contacts read, call records read, attachments read/export, participant/account relation, restricted/private content reveal, search/compare, annotation/dispute, handoff prepare, cross-tenant source correlation and provenance export.

## 22. Limites et erreurs
Unsupported app format, malformed record, missing participant/account, timezone ambiguity, partial extraction, deleted ambiguity, sync duplication, encrypted content, permission denial or cross-tenant restriction remain explicit and do not create authorship/intent conclusions.

## 23. Métriques
Records/conversations/calls/contacts analyzed, private accesses, masked/denied fields, missing status/time data, participant disputes, duplicates, deleted candidates and downstream handoffs.

## 24. Classification de livraison
`defined` / `planned`. No messaging provider, app schema library, API, decryption method, identity resolution engine or code delivered.

## 25. Critères d’acceptation
**Given** a message stored under a device-owner account **When** reviewed **Then** the source/account context is visible but authorship is not automatically assigned to the owner.

**Given** a received message with no read-state field **When** displayed **Then** CMDR does not infer that it was read.

**Given** no AI **When** communications are analyzed **Then** deterministic viewers/parsers/search/filter/diff and human review provide full functionality.

## 26. Questions ouvertes
OPEN-011 supported platforms/apps; OPEN-013 sensitive/analytical mutations; OPEN-014 final record/material relations; OPEN-015 Tool/Run provenance remain open.

## 27. Consommateurs documentaires
Timeline, Anomaly/Hypothesis, Evidence/Finding handoffs, TI, Privacy, Shared Graph, Objects, Permissions, Screens, Quality and Technique.
