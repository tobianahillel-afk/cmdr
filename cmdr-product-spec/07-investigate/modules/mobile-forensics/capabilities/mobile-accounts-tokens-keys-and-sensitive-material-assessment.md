---
id: CAP-INV-712
title: Mobile Accounts, Tokens, Keys and Sensitive Material Assessment
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
# CAP-INV-712 — Mobile Accounts, Tokens, Keys and Sensitive Material Assessment

## 1. Définition
Identifier et évaluer accounts, account identifiers, token/key candidates, certificates, session material, authentication references, application context, expiry, restrictions, source, access events, masking and audit without validating, testing or using credentials.

## 2. Problème utilisateur
Mobile data can contain highly sensitive authentication material. Presence does not prove validity, ownership, current usability or authority to reveal/copy/export/use. Mishandling can expose third-party or cross-tenant secrets.

## 3. Objectifs
- distinguish presence candidate, metadata, masked preview, reveal, copy, export and use;
- keep account ≠ person and token/key/secret candidate ≠ valid/usable credential;
- preserve app/source/expiry/access/restriction context;
- provide Evidence/Detection/TI/Govern handoff context without secret use.

## 4. Non-objectifs
No credential dumping, password/code testing, token validation, login, replay, authentication attempt, cryptographic operation using discovered material, secret use, revocation, rotation, external sharing, bypass, API or command.

## 5. Propriétaire
Investigate owns Account Observations and Sensitive Material Candidates. Platform Settings/Security own secret administration and references; Govern owns revocation/response authority; source owns raw value; Shared owns audit/export mechanisms.

## 6. Utilisateurs
Principal : Sensitive Data Reviewer. Secondaires : Mobile Forensics Analyst, DFIR Analyst, Evidence Reviewer and Investigation Lead.

## 7. Conditions d’entrée
Session, authorized source/app records, explicit purpose, sensitivity classification, permission for requested access level and integrity/completeness limits.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Account identifiers/context | app/package source | account observation | non | source version | unknown |
| Token/key/certificate/session candidates | source | sensitive candidate | non | source version | no candidate invented |
| Application/source relation | CAP-INV-708/707 | context | oui for candidate assessment | Session version | context gap |
| Expiry/status metadata | source | descriptive fields | non | source version | unknown, no validity inference |
| Classification/restriction/permission | Security/source | access boundary | oui | access time | masked/denied |
| Access/audit events | Shared/Security | provenance | according to operation | current | operation blocked if audit required/unavailable |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Session/App Data Observation | Investigate | source/app/scope | read |
| Account/sensitive source records | source | metadata/value only by explicit permission | read per access level |
| Secret Reference | Settings/Security | administrative reference context | metadata read only |
| Action Request/Decision | Govern | future revocation/request context | read only |
| Audit/Trace | Shared/Security | access history | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Account Observation | create/review/dispute/supersede | Investigate concept | account ≠ person |
| Sensitive Material Candidate | create/review/dispute/supersede | Investigate concept | no validity/use claim |
| Sensitive access record | create | Security/Shared audit projection | reason/permission/result, never raw secret value |
| Revocation/response context package | prepare | Govern | no revocation executed |

## 11. Fonctionnalités
Show candidate type, app/source, masked identifier/value preview according to policy, expiry metadata, account relation, classification, restrictions and prior accesses; request elevated reveal/copy/export separately; compare candidates; dispute false candidate; prepare downstream package without exposing raw value where unnecessary.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| View presence/metadata | reviewer | candidate | 0 | metadata permission | masked context | non |
| Reveal value | sensitive reviewer | candidate | 0/2 | explicit reveal permission + reason | audited temporary view | step-up possible |
| Copy/export value | specially authorized reviewer | candidate | 1/2 | separate permission/policy | audited bounded output | step-up/SoD possible |
| Annotate/dispute candidate | reviewer | candidate | 2 | source/reason | versioned disposition | OPEN-013 |
| Prepare revocation/response request | Investigation Lead | context package | 2→3 | sourced risk + authority need | Govern handoff | Govern executes |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Detect candidate patterns | oui | rules/parsers | oui | suggestion | deterministic pattern catalogue |
| Mask values | oui | policy | oui | non | deterministic masking |
| Compare metadata/expiry | oui | oui | oui | summary | table/diff |
| Suggest related account/app | oui | source relations | oui | suggestion | relation table |
| Validate/use/test secret | non | non | non | interdit | no alternative inside Investigate |

## 14. États fonctionnels
Candidate states: `proposed`, `under-review`, `supported`, `weakly-supported`, `contradicted`, `inconclusive`, `disputed`, `superseded`, `withdrawn`. Access state is independently `masked`, `reveal-authorized`, `restricted`, `denied`; none means valid credential.

## 15. États d’interface
Loading never flashes raw values; Empty means no represented candidates; Partial shows masked/unknown metadata; Error preserves masking and logs no value; Offline never enables reveal/copy not guaranteed; Permission denied reveals no secret hint beyond allowed presence; Stale shows source/version.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Account Observations | concepts | CAP-INV-716/717/718 | account context without person attribution |
| Sensitive Material Candidates | concepts | CAP-INV-717/718 | masked by default; validity/use not asserted |
| Sensitive access events | audit events | Security/Trust | actor/reason/permission/result, no secret value |
| Response/revocation context package | handoff | Govern | no action or secret use by Mobile |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-708/709 | account/token/key candidate | CAP-INV-712 | app/source/record/restrictions | source view |
| CAP-INV-712 | correlation | CAP-INV-716/717 | masked candidate metadata, relations, uncertainty | sensitive view |
| CAP-INV-712 | Evidence/Detection/TI need | CAP-INV-718 | minimized candidate context | sensitive view |
| Candidate requires real response | request | Govern | sourced risk, masked reference, authority need | Session |
| Closure | provenance | CAP-INV-719 | access events, decisions, versions | Session |

## 18. Dépendances
CAP-INV-708/709/716..719, Settings Secrets & Connections, Security permission/privacy/secrets, Govern authority, Shared Trace/Export, OPEN-011/013/014/015.

## 19. Source de vérité
Raw sensitive material remains source-controlled. Mobile candidate assessment is Investigate. Secret administration is Settings/Security. Revocation/response is Govern. No Mobile state can convert a candidate into an approved usable secret.

## 20. Provenance et audit
Source package/app/record, candidate type, masked fingerprint/reference, metadata/expiry, classification, access-level requests/results, permission/policy version, Tool/Run that detected candidate, reviewer/dispute, handoff and supersession. Never log raw secret values.

## 21. Permissions fonctionnelles
Account metadata read; sensitive presence/metadata/masked preview/reveal/copy/export as distinct permissions; candidate review/dispute; cross-tenant restriction; handoff prepare; provenance export. **Secret use permission does not exist in Investigate.**

## 22. Limites et erreurs
False pattern, stale/expired-looking metadata, unknown encoding, encrypted/protected source, permission denial, audit failure, cross-tenant source or duplicate candidate keeps values masked and status uncertain; no validation or use occurs.

## 23. Métriques
Candidates by type/state, masked views, reveal/copy/export requests granted/denied, disputes/false candidates, access-policy failures, handoffs and zero secret-use operations.

## 24. Classification de livraison
`defined` / `planned`. No secret engine, credential validator, cracking tool, authentication connector, revocation implementation, API or command delivered.

## 25. Critères d’acceptation
**Given** a token-looking record **When** detected **Then** it is a masked candidate with source/app/metadata and is not tested or declared valid.

**Given** an analyst with Session access but no reveal permission **When** opening the candidate **Then** raw value remains hidden and Session permission grants no additional source access.

**Given** no AI **When** sensitive material is assessed **Then** deterministic patterns, masking, metadata tables and human review provide the complete workflow.

## 26. Questions ouvertes
OPEN-011 future platform/source coverage; OPEN-013 sensitive access class/step-up; OPEN-014 final material relations; OPEN-015 Tool/Run contracts remain open.

## 27. Consommateurs documentaires
Privacy/Trust, Evidence/Finding handoff, Detection Engineering, Threat Intelligence, Govern, Settings Secrets, Objects, Permissions, Screens, Quality and Technique.
