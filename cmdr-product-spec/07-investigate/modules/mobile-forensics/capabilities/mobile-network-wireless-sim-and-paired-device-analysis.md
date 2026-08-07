---
id: CAP-INV-713
title: Mobile Network, Wireless, SIM and Paired Device Analysis
product: investigate
module: mobile-forensics
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-07
requirement_ids: [REQ-INV-001, REQ-PROD-014, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-005, OPEN-008, OPEN-011, OPEN-013, OPEN-014, OPEN-015]
source-of-truth: canonical
---
# CAP-INV-713 — Mobile Network, Wireless, SIM and Paired Device Analysis

## 1. Définition
Analyser network configuration observations, available connection records, Wi-Fi network candidates, Bluetooth records, paired devices, NFC-related records when represented, SIM/eSIM/carrier metadata, phone-number candidates, interfaces, local/remote endpoint candidates, application associations, timestamps and contradictions, puis préparer un handoff vers Network Forensics si nécessaire.

## 2. Problème utilisateur
Un Wi-Fi record ne prouve pas une connexion réussie; Bluetooth pairing ne prouve pas une interaction malveillante; paired device ne prouve pas le même owner; SIM/eSIM ou numéro ne prouve pas une personne. Les états réseau stockés ne remplacent pas une capture Network Forensics.

## 3. Objectifs
- afficher chaque observation avec source/device/app/time context;
- distinguer configuration, historical record, observed connection candidate and pairing;
- corréler Wi-Fi/Bluetooth/SIM/eSIM/endpoint candidates sans interaction active;
- préparer Network Forensics/Timeline/Hypothesis handoffs avec limites explicites.

## 4. Non-objectifs
No active scan, packet generation, wireless probing, Bluetooth/NFC connection, SIM/eSIM action, carrier lookup, phone-number ownership resolution, traffic capture, credential use, Network Forensics replacement, API/protocol or command.

## 5. Propriétaire
Investigate owns Mobile Connectivity/SIM/eSIM/Paired Device Observations. Network Forensics owns full captured-network analysis. Settings/Endpoint own configured sources/capabilities. Shared owns Graph/Timeline.

## 6. Utilisateurs
Principal : DFIR Analyst. Secondaires : Mobile Forensics Analyst, Evidence Reviewer, Investigation Lead and SOC Analyst.

## 7. Conditions d’entrée
Session, device/source context, represented network/wireless/SIM/pairing records, integrity/completeness limits and permission to related identifiers.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Network/interface/config records | package/app data | source records | non | source version | unknown/no claim |
| Wi-Fi/Bluetooth/NFC records | package/app data | wireless candidates | non | source version | not represented |
| SIM/eSIM/carrier/number metadata | package/source | identity/connectivity metadata | non | source version | unknown |
| Endpoint/application associations | source/app data | relation candidates | non | source version | unknown |
| Device/platform context | CAP-INV-703 | scope | oui | Session version | limited interpretation |
| Integrity/permission restrictions | CAP-INV-705/Security | trust/access | oui | access time | partial/masked/denied |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Session / Device Context | Investigate | scope/device/source | read |
| Connectivity/wireless/SIM source records | source | config/history/metadata | read by permission |
| App Data Observation | Investigate | app/endpoint relation | read/link |
| Network Forensics artifacts | Investigate Network module | captured-flow context if already available | read/link only |
| Entity/Graph/Timeline | Shared | relation/time mechanism | consume |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Connectivity Observation | create/review/dispute/supersede | Investigate concept | record ≠ successful connection unless source proves state |
| SIM/eSIM Observation | create/review/dispute/supersede | Investigate concept | SIM/number ≠ person |
| Paired Device Observation | create/review/dispute/supersede | Investigate concept | pairing ≠ same owner/malicious interaction |
| Network Forensics handoff package | prepare | Network module | no scan/capture execution |

## 11. Fonctionnalités
Tables/graph candidates for interfaces, network configs, Wi-Fi identifiers/records, Bluetooth/NFC/paired-device records, SIM/eSIM/carrier/phone-number candidates, local/remote endpoints and app associations; filters by time/source/type; comparison across backups; contradictions; deep links to Network Forensics.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Browse/search/filter records | analyste | observations | 0 | read | scoped view | non |
| Compare sources/backups | analyste | observations | 1 | both accessible | diff | non |
| Annotate/dispute relation | analyste | observation | 2 | source | versioned disposition | OPEN-013 |
| Correlate endpoint/app/device candidates | analyste | relations | 1/2 | source permission | candidate relation | non |
| Prepare Network Forensics handoff | analyste | package | 2 | relevant source gap/context | destination package | no active collection |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Normalize identifiers/endpoints | oui | oui | oui | non | deterministic normalization |
| Group repeated networks/pairings | oui | rules | oui | suggestion | table/graph filters |
| Compare configs/backups | oui | oui | oui | summary | diff |
| Propose anomaly/correlation | oui | rules | oui | suggestion | explicit rules/timeline |
| Probe/connect/resolve ownership | non | non | non | interdit | no active alternative in Mobile |

## 14. États fonctionnels
Observations use `proposed`, `under-review`, `supported`, `weakly-supported`, `contradicted`, `inconclusive`, `disputed`, `superseded`, `withdrawn`. Connection/pairing states remain exactly source-declared; unknown remains unknown.

## 15. États d’interface
Loading preserves source/type filters; Empty means no represented records; Partial lists missing interfaces/areas; Error keeps valid records; Offline read-only; Permission denied masks identifiers; Stale shows source/back-up version.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Connectivity/SIM/Paired Device Observations | concepts | CAP-INV-714/716/717/718 | source/time/device/app/limits |
| Endpoint/network relation candidates | relations | Shared Graph/Network Forensics | no successful/malicious inference |
| Network Forensics handoff | package | CAP-INV-380..397 | no active scan/capture |
| Timeline inputs | events | CAP-INV-716 | timestamps/source quality preserved |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-703/708 | connectivity records | CAP-INV-713 | device/app/source/restrictions | source view |
| CAP-INV-713 | cross-device/pairing/sync | CAP-INV-714 | paired candidates, source/version | connectivity |
| CAP-INV-713 | temporal/correlation need | CAP-INV-716/717 | observations, endpoints, contradictions | connectivity |
| CAP-INV-713 | full network analysis need | Network Forensics | source refs, captured-data refs if any, question | connectivity |
| Selected record | derive/handoff | CAP-INV-718 | source/context/restrictions | connectivity |

## 18. Dépendances
CAP-INV-703/705/708/714/716..719, CAP-INV-380..397 Network Forensics, Settings/Endpoint, Shared Graph/Timeline, OPEN-005/008/011/013/014/015.

## 19. Source de vérité
Raw connectivity records remain source. Mobile observations are Investigate analysis. Network Forensics remains canonical for captured traffic. Settings/Endpoint retain configuration/capabilities. No relation grants active network/device access.

## 20. Provenance et audit
Package/device/app/source record, identifiers/endpoints, SIM/eSIM/carrier fields, timestamp/timezone, parser/Tool/Run, comparison inputs, relation candidates, disputes, network handoff and human disposition.

## 21. Permissions fonctionnelles
Network/connectivity metadata read, Wi-Fi/Bluetooth/NFC record read, SIM/eSIM/carrier/number metadata read, paired-device relation read, endpoint/app relation, cross-device correlation, Network handoff prepare, Derived Artifact export and provenance read.

## 22. Limites et erreurs
Rotating/random identifiers, stale configs, partial extraction, missing success state, duplicated sync record, paired-device ambiguity, restricted phone/carrier data, parser failure or permission denial remain explicit; no active verification is performed.

## 23. Métriques
Observations by type, unknown/contradicted connection states, paired-device candidates, SIM/eSIM records, identifier masks, network handoffs, cross-device correlations and active-interaction count fixed at zero.

## 24. Classification de livraison
`defined` / `planned`. No wireless stack, scanner, carrier service, packet capture, protocol, API, connector or code delivered.

## 25. Critères d’acceptation
**Given** a stored Wi-Fi network record **When** analyzed **Then** it remains a stored record and is not presented as proof of successful connection.

**Given** a Bluetooth paired-device record **When** correlated **Then** pairing does not imply same owner or malicious interaction.

**Given** no AI **When** connectivity is analyzed **Then** deterministic normalization, tables/graphs, filters, diff and human review provide the workflow.

## 26. Questions ouvertes
OPEN-011 platform/wireless/SIM coverage; OPEN-008 source/Endpoint support; OPEN-013 relations; OPEN-014 final objects; OPEN-015 Tools/Runs remain open.

## 27. Consommateurs documentaires
Backup/Sync, Timeline, Anomaly/Hypothesis, Network Forensics, Artifact/Evidence handoff, Privacy, Objects, Permissions, Screens, Quality and Technique.
