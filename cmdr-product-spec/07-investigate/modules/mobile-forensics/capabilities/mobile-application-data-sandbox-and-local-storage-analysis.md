---
id: CAP-INV-708
title: Mobile Application Data, Sandbox and Local Storage Analysis
product: investigate
module: mobile-forensics
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-07
requirement_ids: [REQ-INV-001, REQ-PROD-014, REQ-PROD-020, REQ-OBJ-003, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-005, OPEN-011, OPEN-013, OPEN-014, OPEN-015]
source-of-truth: canonical
---
# CAP-INV-708 — Mobile Application Data, Sandbox and Local Storage Analysis

## 1. Définition
Analyser les données applicatives représentées : sandbox, local databases/files/caches/preferences/storage, session records, account references, downloads, browser data, application-specific records, shared containers, protected data, timestamps, deleted candidates, restrictions et provenance.

## 2. Problème utilisateur
Les données d’application sont fortement contextuelles et privées. Un record local ne prouve ni intention utilisateur, ni authorship, ni état courant; une base partielle ou synchronisée peut aussi contenir des données historiques ou de tiers.

## 3. Objectifs
- parcourir/search/filter des records applicatifs autorisés avec app/source context;
- distinguer raw/protected/masked/deleted candidate states;
- comparer versions/backups and relate account/download/browser/shared-container records;
- préparer communications/media/sensitive/deleted/timeline/artifact analysis without losing lineage.

## 4. Non-objectifs
No application execution, login, credential use, network request, app-specific offensive procedure, bypass protected data, final database schema, parser implementation, cloud-service interaction or user-intent inference.

## 5. Propriétaire
Investigate owns Application Data Observations. Raw records remain package/source material; Security governs private access; Settings owns credentials/connections; destination analysis families own their specialized workflows.

## 6. Utilisateurs
Principal : Mobile Forensics Analyst. Secondaires : DFIR Analyst, Evidence Reviewer, Sensitive Data Reviewer and Investigation Lead.

## 7. Conditions d’entrée
Session, application observation, represented app data area, integrity/completeness/accessibility context and content-level permission where necessary.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| App identifier/version/context | CAP-INV-707 | interpretation scope | oui | selected source version | generic/unknown app context |
| App data records/files/databases | package/extraction | analytical source | oui | source version | no data claim |
| Filesystem/source path | CAP-INV-706 | lineage | non if source already links | current observation | provenance gap |
| Integrity/completeness/accessibility | CAP-INV-705 | limitations | oui | reviewed version | partial/unverified |
| Sensitive/private classification | source/Security | handling | selon data | access time | mask/deny |
| Deleted/sync indicators | source | record state candidates | non | source version | unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Session / App Observation | Investigate | scope/app/source | read |
| Package/filesystem/source records | source | databases/files/cache/preferences | read according to permission |
| Account/Secret Reference projection | Settings/Security | context only | metadata/masked read only |
| Artifact | Investigate | derivation target/ref | read/link |
| Shared Search/Comparison | Shared | mechanisms | consume |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Application Data Observation | create/review/dispute/supersede | Investigate concept | app/source/record version required |
| Record relation candidate | create/dispute | Investigate | data ≠ user intent/authorship |
| Sensitive-access annotation | create | Investigate/Security audit | content not copied into audit |
| Derived Artifact proposal | prepare | CAP-INV-718 | bounded record/file + lineage |

## 11. Fonctionnalités
Tree/table/database-like viewers where source permits, search/filter, field/record inspection, timestamp/source display, protected-data masking, compare versions/backups, browser/download/account/shared-container relations, deleted-candidate tagging and deep links to communications/media/sensitive/deleted analysis.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Browse/search/filter | analyste | app records | 0 | read | scoped results | non |
| Reveal/read restricted field | authorized reviewer | protected data | 0/2 | explicit permission/policy | audited view | step-up possible |
| Compare records/versions | analyste | observations | 1 | both source permissions | diff | non |
| Annotate/dispute relation | analyste | observation | 2 | source | versioned disposition | OPEN-013 |
| Prepare Derived Artifact | analyste | selected file/record | 1/2 | purpose + permission | CAP-INV-718 input | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Parse known represented structure | oui | oui where parser exists | oui | non | deterministic parser/viewer |
| Search/filter/compare | oui | oui | oui | non necessary | query/table/diff |
| Group related records | oui | rules | oui | suggestion | keys/tags |
| Summarize authorized records | oui | structured counts | oui | yes with permission | table/exported summary |
| Infer user intent/authorship | non | non | non | interdit | human evidence-based review |

## 14. États fonctionnels
Observations: `proposed`, `under-review`, `supported`, `weakly-supported`, `contradicted`, `inconclusive`, `disputed`, `superseded`, `withdrawn`. Source records may be `available`, `partial`, `protected`, `restricted`, `deleted-candidate`, `inaccessible`.

## 15. États d’interface
Loading preserves app/source; Empty distinguishes no represented records from no activity; Partial lists missing databases/areas; Error keeps valid records; Offline read-only; Permission denied reveals no protected value; Stale shows package/version.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Application Data Observations | concepts | CAP-INV-709..717 | app/source/record/version/limitations |
| Communication/media/account relations | candidate links | 709/710/712 | no authorship/person intent inferred |
| Deleted record candidates | candidates | CAP-INV-715 | source state + uncertainty |
| Derived Artifact proposal | package input | CAP-INV-718 | bounded source + privacy restrictions |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-707 | app data available | CAP-INV-708 | app/version/source area/limits | inventory |
| CAP-INV-708 | messaging/call/contact record | CAP-INV-709 | record/source/app context/restrictions | app data |
| CAP-INV-708 | media/document | CAP-INV-710 | file/record/metadata/source | app data |
| CAP-INV-708 | token/account/secret candidate | CAP-INV-712 | masked candidate context | app data |
| CAP-INV-708 | deleted/residual record | CAP-INV-715 | source region/record/recovery hints | app data |
| Records/timestamps | correlate | CAP-INV-716/717 | observations, source quality, contradictions | app data |
| Selected source | derive | CAP-INV-718 | bounded record/file, purpose, restrictions | app data |

## 18. Dépendances
CAP-INV-705..707/709..719, Security privacy/secrets, Settings secret references, Shared Search/Comparison/Trace, OPEN-005/011/013/014/015.

## 19. Source de vérité
Raw application records remain source facts. Mobile data observations and candidate relations remain Investigate. Credentials/secrets remain Settings/Security. User intent/authorship remains an analyst assessment, never a source field invented by CMDR.

## 20. Provenance et audit
App/package/version, source path/record ID, parser/Tool/Run, query/filter, field access level, sensitive access reason, timestamps, relations, deleted candidate state, comparison, annotations, Derived Artifact selection, errors and human disposition.

## 21. Permissions fonctionnelles
Application data read, restricted data read, communication/media/location/account field read separately, sensitive presence/metadata/reveal permissions, search/compare, annotation, Derived Artifact prepare/export and provenance read.

## 22. Limites et erreurs
Malformed DB/file, unsupported record, missing shared container, encrypted/protected field, partial extraction, deleted ambiguity, sync ambiguity, permission denial or parser failure remain explicit; no missing activity or user intent inferred.

## 23. Métriques
Records/databases parsed, restricted fields accessed/denied, partial areas, parser errors, comparisons, deleted candidates, downstream handoffs and manual rejection/modification of suggestions.

## 24. Classification de livraison
`defined` / `planned`. No app-specific schema library, cloud login, credential use, parser engine, API, protocol or code delivered.

## 25. Critères d’acceptation
**Given** a stored message record in an app database **When** viewed **Then** it remains a stored record with app/source context and is not automatically attributed to the device owner.

**Given** a protected token candidate **When** permission is absent **Then** presence/metadata may be shown only as allowed, value remains masked and use is impossible.

**Given** no AI **When** app data is analyzed **Then** deterministic viewers/parsers/search/filter/diff and human review provide complete functionality.

## 26. Questions ouvertes
OPEN-011 platform/app coverage; OPEN-005 parsers/engines; OPEN-013 sensitive/class-2 access; OPEN-014 final material relations; OPEN-015 Tool/Run contracts.

## 27. Consommateurs documentaires
Communications, Media/Documents, Location/Sensors, Sensitive Material, Backup/Sync, Deleted/Recovery, Timeline, Hypothesis, Artifact Handoff, Privacy, Objects, Permissions and Quality.
