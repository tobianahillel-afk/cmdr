---
id: CAP-INV-714
title: Mobile Backup, Synchronization and Cross-Device Artifact Analysis
product: investigate
module: mobile-forensics
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-07
requirement_ids: [REQ-INV-001, REQ-PROD-014, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-011, OPEN-013, OPEN-014, OPEN-015]
source-of-truth: canonical
---
# CAP-INV-714 — Mobile Backup, Synchronization and Cross-Device Artifact Analysis

## 1. Définition
Analyser local backups, synchronized backups, backup versions/sources/destinations, sync records, cloud-backed records, paired-device/shared-account relations, application/media/contact/message sync candidates, gaps/conflicts, timestamps, provenance, scope and restrictions.

## 2. Problème utilisateur
Un backup timestamp ne représente pas forcément l’état actuel du device. Un record synchronisé n’est pas une preuve qu’il était local, et un cloud-backed artifact ne constitue pas une Cloud Analysis complète. Cross-device/shared-account relations can also cross privacy/tenant boundaries.

## 3. Objectifs
- compare backup versions and source/destination context;
- identify sync/cross-device candidates while preserving local/cloud ambiguity;
- expose gaps, conflicts, duplicate/copied records and timestamp quality;
- prepare Cloud/Timeline/Deleted/Anomaly handoffs without source-permission expansion.

## 4. Non-objectifs
No cloud login/API, backup creation/restore, synchronization command, paired-device interaction, cross-tenant access grant, provider-specific backup format, Cloud Analysis replacement, credential use or device mutation.

## 5. Propriétaire
Investigate owns Backup/Synchronization Observations. Collection/source owns acquired backup representation; Cloud Analysis owns full Cloud analytical scope; Settings owns providers/connectors/credentials; Shared owns comparison/linking.

## 6. Utilisateurs
Principal : Mobile Forensics Analyst. Secondaires : DFIR Analyst, Evidence Reviewer, Investigation Lead and Sensitive Data Reviewer.

## 7. Conditions d’entrée
Session, at least one authorized backup/sync representation, device/source context, integrity/completeness limits, permissions for each cross-device/cloud source and explicit purpose.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Local/synchronized backup | source/package | backup representation | oui for backup analysis | selected version | no backup claim |
| Backup source/destination/version | source | provenance metadata | oui if available | backup version | unknown/gap |
| Sync records/cloud-backed flags | app/package/source | synchronization candidates | non | source version | unknown |
| Paired/shared-account/device relations | CAP-INV-712/713/source | cross-device candidates | non | Session/source version | unknown |
| Integrity/completeness | CAP-INV-705 | trust boundary | oui | reviewed version | partial/unverified |
| Permissions/classification | Security/Settings/source | scope/access | oui for each source | access time | masked/blocked; no access propagation |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Session / Device Context | Investigate | scope/device/source | read |
| Device Backup / synchronized representation | source | data/version/provenance | read by permission |
| Account/Paired Device Observations | Investigate concepts | relation candidates | read/link |
| Cloud source/Cloud Analysis context | Settings/Cloud Analysis | source relation only | read if permitted; no Cloud ownership |
| Shared Versioning/Comparison/Linking | Shared | mechanisms | consume |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Synchronization Observation | create/review/dispute/supersede | Investigate concept | synchronized ≠ certain local state |
| Backup version relation | create/compare/dispute | Investigate | source/destination/version explicit |
| Cross-device/shared-account relation candidate | create/dispute | Investigate | no permission or owner inference |
| Cloud Analysis handoff context | prepare | Cloud Analysis | cloud-backed artifact ≠ Cloud analysis |

## 11. Fonctionnalités
List/compare backup versions, show source/destination/timestamps, group sync candidates, compare local/synchronized records, mark conflicts/gaps/duplicates, link paired/shared-account candidates, preserve cloud-backed ambiguity, filter by app/data category/device and prepare downstream handoffs.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Browse backup/sync metadata | analyste | observations | 0 | read | scoped view | non |
| Compare versions/records | analyste | backups/records | 1 | permissions to both | diff | non |
| Annotate/dispute sync relation | analyste | observation | 2 | source | versioned disposition | OPEN-013 |
| Correlate cross-device | analyste | relation candidate | 1/2 | access each source | candidate relation only | non |
| Prepare Cloud/Collection handoff | lead | package | 2 | question/gap | destination package | destination reauthorizes |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Compare backup versions | oui | oui | oui | summary | diff/table |
| Match duplicate/sync records | oui | keys/rules | oui | suggestion | deterministic matching |
| Detect conflicts/gaps | oui | rules | oui | explanation | checklist/diff |
| Propose cross-device relation | oui | sourced rules | oui | suggestion | relation table |
| Fetch remote cloud/device data | non | non | non | interdit | destination-owned collection/Cloud workflow |

## 14. États fonctionnels
Observations use `proposed`, `under-review`, `supported`, `weakly-supported`, `contradicted`, `inconclusive`, `disputed`, `superseded`, `withdrawn`. Backups may be `available`, `partial`, `restricted`, `inaccessible`, `superseded`.

## 15. États d’interface
Loading preserves compared versions; Empty means no represented backup/sync data; Partial shows gaps; Error preserves valid matches; Offline read-only; Permission denied blocks cross-source values; Stale marks source/version age.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Backup/Synchronization Observations | concepts | CAP-INV-715..718 | source/version/device/sync ambiguity |
| Version comparison | result | analyst/QA | local/sync/cloud states remain distinct |
| Cross-device relation candidates | relations | Shared Graph/Timeline | no permission/ownership inference |
| Cloud/Collection handoff | package | CAP-INV-601 or Collection | question/gap only; no remote action |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Backup/package | sync/version analysis | CAP-INV-714 | source/version/device/restrictions | Session |
| CAP-INV-713/712 | paired/shared account relation | CAP-INV-714 | candidates/source/confidence | connectivity/account |
| CAP-INV-714 | deleted/conflicting record | CAP-INV-715 | record/source/version/conflict | backup |
| CAP-INV-714 | temporal/correlation | CAP-INV-716/717 | observations/gaps/time quality | backup |
| Cloud-backed question | full cloud analysis | CAP-INV-601 | minimized source relation + question | Mobile Session |
| Selected record | derive/handoff | CAP-INV-718 | source/version/restrictions | backup |

## 18. Dépendances
CAP-INV-601..618 Cloud, CAP-INV-705/712/713/715..719, Settings sources/connectors, Shared Comparison/Linking, Security, OPEN-008/011/013/014/015.

## 19. Source de vérité
Backup/sync raw records remain source. Mobile observations are Investigate. Cloud configuration/source remains Settings; Cloud Analysis owns full Cloud analytical Session; cross-device link never grants access.

## 20. Provenance et audit
Backup/source/destination/version, sync fields, device/account/pairing refs, timestamps, comparison inputs, parsers/Tools/Runs, conflicts/gaps, permissions to each source, cross-device relation, handoffs and human disposition.

## 21. Permissions fonctionnelles
Backup read, synchronized data read, cross-device relation read/create, cloud-backed metadata read, compare, restricted source access, annotation/dispute, handoff prepare, Derived Artifact export and provenance export.

## 22. Limites et erreurs
Missing backup manifest/version, stale timestamp, partial sync, duplicate/conflicting records, cloud permission denial, cross-tenant restriction, unknown source device, parser failure or inaccessible backup remain explicit; no local/current-state inference is forced.

## 23. Métriques
Backups/versions compared, sync candidates/conflicts/gaps, cross-device relations, cloud handoffs, permission blocks, duplicates and manually disputed automated matches.

## 24. Classification de livraison
`defined` / `planned`. No backup format, cloud provider, sync engine, connector, restore method, API or code delivered.

## 25. Critères d’acceptation
**Given** a synchronized message present in a backup **When** reviewed **Then** it remains synchronized/source-scoped and is not presented as certain local device data.

**Given** a cloud-backed artifact **When** more Cloud context is needed **Then** a handoff to Cloud Analysis is prepared rather than expanding Mobile scope.

**Given** no AI **When** backups are compared **Then** deterministic version diff, matching rules, tables and human review provide full functionality.

## 26. Questions ouvertes
OPEN-011 backup/platform/tool strategy; OPEN-008 provider/source support; OPEN-013 relations; OPEN-014 final backup/material model; OPEN-015 Tool/Run contracts remain open.

## 27. Consommateurs documentaires
Deleted/Recovery, Timeline, Anomaly/Hypothesis, Cloud Analysis, Collection, Artifact/Evidence handoff, Privacy, Objects, Permissions, Quality and Technique.
