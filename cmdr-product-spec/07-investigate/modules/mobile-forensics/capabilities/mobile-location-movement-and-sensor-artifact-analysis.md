---
id: CAP-INV-711
title: Mobile Location, Movement and Sensor Artifact Analysis
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
# CAP-INV-711 — Mobile Location, Movement and Sensor Artifact Analysis

## 1. Définition
Analyser des location records, GPS-like observations, cell/network location candidates, movement/routes/geofence candidates, map searches, navigation records, sensor records et health/activity records lorsqu’ils sont autorisés, avec timestamps/timezone, confidence, gaps, contradictions, source application et device context.

## 2. Problème utilisateur
La localisation d’un appareil ne prouve pas la présence d’une personne, et un sensor record ne prouve pas une action humaine. Les données de santé/activité et localisation sont particulièrement sensibles et peuvent être synchronisées, estimées, copiées ou partielles.

## 3. Objectifs
- présenter chaque observation avec source, device/app context, timestamp quality and precision/confidence where declared;
- distinguer device location, user presence, route/movement and network-derived candidates;
- comparer sources, gaps and contradictions;
- produire timeline/correlation and anomaly candidates without automatic attribution.

## 4. Non-objectifs
No live tracking, geofencing action, device query, cell/network triangulation implementation, external map/provider call, certain human-presence inference, health diagnosis, biometric identification, legal conclusion, API or code.

## 5. Propriétaire
Investigate owns Location/Sensor Observations. Raw records remain source material; Shared owns generic map-independent Graph/Timeline mechanisms; Security controls sensitive access; Settings owns configured sources.

## 6. Utilisateurs
Principal : Sensitive Data Reviewer or Mobile Forensics Analyst. Secondaires : DFIR Analyst, Evidence Reviewer and Investigation Lead.

## 7. Conditions d’entrée
Session, device/source context, authorized location/sensor/health records, explicit purpose, privacy permission and integrity/completeness limits.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Location/navigation records | package/app data/media metadata | source observations | selon scope | source version | no location claim |
| Movement/sensor/health records | package/app data | sensitive source observations | non | source version | unknown/not represented |
| Device/app context | CAP-INV-703/707/708 | interpretation context | oui where available | Session version | attribution limited |
| Timestamps/timezone/precision | source | time/quality fields | non | source version | unknown/estimated explicit |
| Integrity/completeness | CAP-INV-705 | trust boundary | oui | reviewed version | partial/unverified |
| Privacy/permission | Security/source | access boundary | oui | access time | masked/denied |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Session / Device Context | Investigate | scope/device/source | read |
| Location/Sensor source records | source | coordinates/route/activity metadata | read by permission |
| Media/App observations | Investigate | related metadata/app context | read/link candidate |
| Entity/Timeline | Shared | generic relation/time mechanisms | consume |
| Tenant/Environment | Settings | boundary | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Location Observation | create/review/dispute/supersede | Investigate concept | source/time/quality required; device location ≠ person location |
| Sensor Observation | create/review/dispute/supersede | Investigate concept | sensor record ≠ human action |
| Route/geofence/presence relation candidate | create/dispute | Investigate | candidate only; evidence support visible |
| Timeline correlation proposal | create | CAP-INV-716 | source/time quality preserved |

## 11. Fonctionnalités
Table/map-independent list of records, time/precision/source filters, route/movement grouping, source-app pivots, media metadata comparison, timezone normalization, gaps/contradictions, sensitive masking, multi-source comparison and bounded timeline/hypothesis handoff.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Browse metadata/records | reviewer | observations | 0 | read | scoped view | non |
| Reveal sensitive location/health record | sensitive reviewer | record | 0/2 | explicit permission/purpose | audited access | step-up possible |
| Compare sources/timezones | analyste | observations | 1 | source access | comparison | non |
| Annotate/dispute relation | analyste | observation | 2 | source | versioned disposition | OPEN-013 |
| Prepare correlation/handoff | analyste | selected records | 1/2 | purpose | CAP-INV-716/718 input | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Normalize timestamps/timezones | oui | oui | oui | non | deterministic conversion |
| Compare location sources | oui | rules | oui | summary | table/diff |
| Group route/movement candidates | oui | rules | oui | suggestion | time/distance filters |
| Propose anomaly/correlation | oui | rules | oui | suggestion | explicit rules/timeline |
| Infer certain human presence/health meaning | non | non | non | interdit | human evidence-based assessment |

## 14. États fonctionnels
Observations use `proposed`, `under-review`, `supported`, `weakly-supported`, `contradicted`, `inconclusive`, `disputed`, `superseded`, `withdrawn`. Precision/time source may be `observed`, `recorded`, `reconstructed`, `estimated`, `synchronized`, `absent`.

## 15. États d’interface
Loading preserves time/source filters; Empty means no represented authorized records; Partial shows gaps/precision missing; Error keeps valid records; Offline read-only; Permission denied masks location/health values; Stale shows source version/age.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Location/Sensor Observations | concepts | CAP-INV-716/717/718 | source/device/app/time/quality/limits |
| Route/movement candidates | relations | analyst/timeline | no human-presence certainty |
| Sensitive-access events | audit | Security/Trust | no sensitive content in log |
| Correlation proposal | package input | CAP-INV-716 | source/time quality retained |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-703/708/710 | location/sensor records | CAP-INV-711 | device/app/media/source context | source analysis |
| CAP-INV-711 | temporal correlation | CAP-INV-716 | observations, time quality, gaps | location |
| CAP-INV-711 | anomaly/hypothesis | CAP-INV-717 | candidate movement/sensor pattern, contradictions | location |
| Selected record | derive/handoff | CAP-INV-718 | source, purpose, privacy restriction | location |
| Sync ambiguity | cross-device review | CAP-INV-714 | source/device/sync relation | location |

## 18. Dépendances
CAP-INV-703/705/708/710/714/716..719, Security privacy, Shared Timeline/Graph, Settings source projection, OPEN-011/013/014/015.

## 19. Source de vérité
Raw location/sensor/health records remain source material. Mobile observations are Investigate assessments. Generic Timeline/Graph are Shared. Human presence or identity is never created automatically from device records.

## 20. Provenance et audit
Source package/app/device, record ID, time/timezone/precision, location/movement/sensor fields, privacy classification, access level/reason, parser/Tool/Run, comparisons, candidate relations, disputes, handoffs and human disposition.

## 21. Permissions fonctionnelles
Location metadata/content read, health/sensor read, map/navigation record read, cross-device relation, sensitive reveal, compare/correlate, annotation/dispute, Derived Artifact prepare/export and provenance read.

## 22. Limites et erreurs
Missing precision/timezone, spoofable/copied/synchronized records, partial extraction, stale records, conflicting sources, permission denial or unsupported format remain explicit; no certain user presence or activity is inferred.

## 23. Métriques
Records by source/quality, restricted accesses, timezone/precision gaps, contradictions, route candidates, correlation proposals, disputes and downstream handoffs.

## 24. Classification de livraison
`defined` / `planned`. No GPS/network location engine, map provider, health model, tracking integration, API or code delivered.

## 25. Critères d’acceptation
**Given** a location record from a synchronized app **When** viewed **Then** device/app/source/sync context is visible and it is not treated as certain local user presence.

**Given** protected health/activity data **When** permission is absent **Then** no value is revealed and denial exposes no protected detail.

**Given** no AI **When** records are analyzed **Then** deterministic time normalization, tables, filters, comparison and human review cover the workflow.

## 26. Questions ouvertes
OPEN-011 platform/sensor coverage; OPEN-013 sensitive/class-2 access; OPEN-014 final relation models; OPEN-015 Tool/Run contracts remain open.

## 27. Consommateurs documentaires
Timeline, Anomaly/Hypothesis, Evidence/Finding handoff, Privacy/Trust, Shared Graph, Objects, Permissions, Screens, Quality and Technique.
