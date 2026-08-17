---
id: CAP-INV-707
title: Mobile Application Inventory and Package Analysis
product: investigate
module: mobile-forensics
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-07
requirement_ids: [REQ-INV-001, REQ-PROD-014, REQ-PROD-020, REQ-OBJ-003, REQ-SEC-001]
open_decisions: [OPEN-005, OPEN-011, OPEN-013, OPEN-014, OPEN-015]
source-of-truth: canonical
---
# CAP-INV-707 — Mobile Application Inventory and Package Analysis

## 1. Définition
Analyser les applications et packages représentés par la source : identifiers, versions, install sources déclarées, signatures/metadata disponibles, install/update observations, removal candidates, permissions déclarées, associated accounts, shared components, extensions, configuration, support status et contradictions.

## 2. Problème utilisateur
La présence d’un package ne prouve ni son usage, ni son activité, ni son caractère malveillant. Les inventaires incomplets ou synchronisés peuvent aussi présenter une vue historique plutôt qu’un état device courant.

## 3. Objectifs
- produire des Application Observations sourcées et versionnées;
- comparer packages/versions/configurations entre représentations;
- lier les zones de données associées sans inférer user intent;
- préparer app-data, static/reverse/dynamic ou anomaly handoffs lorsque justifiés.

## 4. Non-objectifs
Aucune exécution/install/uninstall, réputation réseau, marketplace query, signature validation implementation, malware confirmation, package format final, platform support list, engine, connector, API ou command.

## 5. Propriétaire
Investigate possède Application Observations. Source/package owner keeps raw inventory; Settings owns supported-platform/source administration; Static/Reverse/Dynamic owners analyze Derived Artifacts; TI owns canonical intelligence knowledge.

## 6. Utilisateurs
Principal : Mobile Forensics Analyst. Secondaires : DFIR Analyst, Investigation Lead, Evidence Reviewer et SOC Analyst.

## 7. Conditions d’entrée
Session, represented filesystem/package inventory, integrity/completeness/accessibility context, device/platform candidate context and permission to app metadata/source areas.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Application/package records | extraction/package | inventory source | oui | source version | no inventory claim |
| Filesystem/application areas | CAP-INV-706 | source relations | non | current analysis | limited package-only view |
| Platform candidate | CAP-INV-703 | interpretation context | non | Session version | generic app semantics only |
| Integrity/completeness limits | CAP-INV-705 | confidence boundary | oui | reviewed version | partial/unverified banner |
| Signature/install/update metadata | source | descriptive fields | non | source version | unknown |
| Permissions/accounts/shared components | source | app relations | non | source version | unknown, never inferred |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Session/Device Context | Investigate | source/platform/scope | read |
| Filesystem Observation | Investigate | app/package area relation | read |
| Package/source records | source | identifiers/version/configuration | read |
| Artifact / Malware Sample | Investigate | potential derived handoff refs | read/link |
| TI knowledge / Settings source projection | owner | context only | read, no promotion |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Application Observation | create/review/dispute/supersede | Investigate concept | package source/version required |
| Package relation candidate | create/dispute | Investigate | installed/present ≠ used/active |
| Associated account/component relation | create/review | Investigate | candidate, not person/ownership proof |
| Derived Artifact proposal | prepare | CAP-INV-718 | source lineage mandatory |

## 11. Fonctionnalités
Inventory table/tree, identifier/version filters, install/update/removal observations, permissions/configuration fields, signature metadata as observed, associated accounts/components/extensions, version comparison, source contradiction display, open related data and prepare technical handoffs.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Browse/search/filter inventory | analyste | app records | 0 | read | scoped observations | non |
| Compare versions/packages | analyste | observations | 1 | sources accessible | diff | non |
| Annotate/dispute relation | analyste | observation | 2 | source | versioned disposition | OPEN-013 |
| Open app data | analyste | relation | 0 | data permission | CAP-INV-708 | non |
| Prepare package Artifact | analyste | package file/record | 1/2 | export/source permission | CAP-INV-718 | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Parse inventory | oui | oui | oui | non | parser/table |
| Normalize identifiers/versions | oui | oui | oui | non | deterministic normalization |
| Group related packages/components | oui | rules | oui | suggestion | tags/relations |
| Flag version/config anomaly candidate | oui | rules | oui | suggestion | filters/diff |
| Confirm malicious application | humain downstream | non | non | interdit | Static/Reverse/Dynamic + review |

## 14. États fonctionnels
Application observations use `proposed`, `under-review`, `supported`, `weakly-supported`, `contradicted`, `inconclusive`, `disputed`, `superseded`, `withdrawn`. Removal is always candidate unless source semantics establish it.

## 15. États d’interface
Loading keeps filters; Empty means no represented app records, not no apps on original device; Partial shows coverage limits; Error keeps valid rows; Offline read-only; Permission denied masks restricted details; Stale shows source/version age.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Application Observations | concepts | CAP-INV-708/717/716 | source/version/limitations |
| Package/version comparison | result | analyst/reviewer | no usage/malware inference |
| App-data relations | links | CAP-INV-708 | source path and permissions preserved |
| Derived Artifact proposal | package input | CAP-INV-718 | bounded source selection |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-706 | app/package areas | CAP-INV-707 | source paths, package metadata, limitations | filesystem |
| CAP-INV-707 | analyze app records | CAP-INV-708 | app identifier/version/source areas | inventory |
| CAP-INV-707 | suspicious candidate | CAP-INV-717 | observation, source, contradiction | inventory |
| Package/file selected | technical analysis | CAP-INV-718 → Static/Reverse/Dynamic | Derived Artifact proposal + lineage | inventory |
| App timestamps | correlation | CAP-INV-716 | install/update/removal candidates | inventory |

## 18. Dépendances
CAP-INV-703/705/706/708/716/717/718/719, Static/Reverse/Dynamic, Settings source projection, TI handoff, Shared Search/Comparison, OPEN-005/011/013/014/015.

## 19. Source de vérité
Raw package/inventory fields remain source facts. Application Observation is Investigate analysis. Platform support and integrations remain Settings/Technique. Maliciousness and Finding status require downstream analytical/review workflows.

## 20. Provenance et audit
Package/source version, identifier/version, install source, signature metadata, timestamps, permissions/configuration, accounts/components, parser/Tool/Run, comparisons, annotations, disputes, derived selection and human disposition.

## 21. Permissions fonctionnelles
Application inventory read, restricted package metadata read, comparison, annotation/dispute, associated account read, package content read separately, Derived Artifact prepare/export and provenance read.

## 22. Limites et erreurs
Missing inventory areas, inconsistent identifiers, partial package metadata, stale backup, unsupported record, permission restriction or parser failure remain explicit; no absent app or usage conclusion is inferred.

## 23. Métriques
Apps/packages observed, unknown versions, conflicting metadata, comparison changes, partial inventory rate, derived-package handoffs, suspicious candidates later rejected/confirmed and manual override of suggestions.

## 24. Classification de livraison
`defined` / `planned`. No app store integration, reputation service, package parser implementation, platform list, malware engine, API or code delivered.

## 25. Critères d’acceptation
**Given** a package record in a backup **When** inventory is reviewed **Then** it is an observed package in that representation, not proof the app is currently installed or used.

**Given** a suspicious app candidate **When** selected **Then** it may be handed to analysis but is not automatically malware or a Finding.

**Given** no AI **When** inventory is analyzed **Then** parser, tables, filters, comparisons and annotations provide full functionality.

## 26. Questions ouvertes
OPEN-011 future supported platforms/package coverage; OPEN-005 engines; OPEN-013 dispositions; OPEN-014 final Artifact/package relations; OPEN-015 Tool/Run contracts.

## 27. Consommateurs documentaires
Application Data, Timeline, Anomaly/Hypothesis, Artifact Handoff, Static/Reverse/Dynamic, Detection/TI, Objects, Permissions, Screens, Quality and Technique.
