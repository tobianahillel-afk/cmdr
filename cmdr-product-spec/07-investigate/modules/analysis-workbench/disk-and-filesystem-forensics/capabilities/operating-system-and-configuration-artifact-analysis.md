---
id: CAP-INV-371
title: Operating System and Configuration Artifact Analysis
product: investigate
module: analysis-workbench
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-05
requirement_ids:
  - REQ-INV-001
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-AI-002
  - REQ-SEC-001
  - REQ-PROD-055
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-371 — Operating System and Configuration Artifact Analysis

## 1. Définition
Examiner les artefacts persistants de système et de configuration comme observations historiques, sans modifier une configuration, administrer la plateforme ni présenter un état ancien comme état actuel certain.

## 2. Problème utilisateur
Des fichiers de configuration, comptes, services, tâches, paramètres ou logs persistants peuvent être incomplets, anciens ou contradictoires. Leur lecture directe peut produire une fausse représentation de l’Endpoint actuel.

## 3. Objectifs
Examiner configuration système persistante, comptes/profils, services ou mécanismes équivalents, tâches planifiées candidates, paramètres de sécurité, configurations réseau persistantes, logs locaux, installation/désinstallation, état de composants et changements; montrer missing/conflicts; compare snapshots; annotate/link Endpoint/Entity/Case.

## 4. Non-objectifs
Ne pas modifier ou administrer une configuration, devenir Platform Settings, affirmer un état actuel, fournir commandes, moteur, format, bypass ou capability Network Forensics complète.

## 5. Propriétaire
Investigate possède l’observation et l’interprétation. Platform Settings possède configuration administrative/Fleet/Policies; Endpoint Agent l’état actuel déclaré; Studio Tools; Shared linking/timeline.

## 6. Utilisateurs
Principal : DFIR/System Artifact Analyst. Secondaires : Investigation Lead, Endpoint Analyst, Evidence Reviewer et Security Reviewer.

## 7. Conditions d’entrée
Session/image/filesystem accessibles; sources et age visibles; Tool/version attribués; permissions pour données système/configuration; limites de l’image héritées.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| System/config artifacts | Filesystem/Tool Calls | persistent observations | oui | image version | unavailable |
| Platform context | Endpoint/Settings projection | support/context | non | freshness visible | artifact-only |
| Change/journal/timeline context | CAP-INV-370/374 | temporal relations | non | same/correlated source | uncorrelated |
| Comparison sources | CAP-INV-377 | snapshots/images | non | versions visible | single-source |
| Permissions/privacy | Security | scope/masking | oui | current | restricted |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Disk Image/File Observation | Investigate | source/artifact/metadata | lire |
| Endpoint/Entity/Case | owners | contextual projection | lire/lier |
| Settings policy/config projection | Platform Settings | current/admin context | lire uniquement |
| Tool/Tool Call | Studio | extraction/result | lire |
| Journal/Timeline | Investigate/Shared | change context | lire/lier |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| System Artifact Observation | créer/annoter/dispute | Investigate concept | persisted ≠ current state |
| Configuration relation | créer/supersede | Investigate | source/time/uncertainty required |
| Endpoint/Entity/Case link | créer | owner relation mechanism | no ownership transfer |
| Platform configuration | aucune mutation | Platform Settings | boundary invariant |

## 11. Fonctionnalités
Show persistent accounts/profiles, services, tasks, security/network configuration, local logs, installation/removal candidates, component states and configuration changes; show gaps/conflicts/age; compare snapshots; annotate and link context. Persistent network configuration is not a connection or network analysis.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspect/filter | Analyst | observations | 0 | read | sourced view | non |
| Compare snapshots | Analyst | comparison | 1 | sources compatible | differences/limits | non |
| Annotate/dispute/link | Analyst | observation/relation | 2 | permission | versioned interpretation | OPEN-013 |
| Prepare handoff | Analyst | selected package | 2 | provenance | CAP-INV-379 | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| identify artifact categories | oui | parser/rules | oui | suggestion | deterministic parsers |
| compare configuration | oui | oui | oui | summary | diff/table |
| explain conflicts | oui | diagnostics | oui | yes | source/age view |
| assert current state | human/corroborated | no | no | prohibited certainty | Endpoint current projection review |

## 14. États fonctionnels
`available`, `partial`, `stale-context`, `conflicting`, `missing`, `restricted`, `unsupported`, `disputed`, `superseded`.

## 15. États d’interface
Loading preserves selection; Empty is no observable artifact, not no configuration; Partial shows gaps; Error keeps valid observations; Offline read-only; Permission denied masks sensitive configuration; Stale highlights age.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| System/Configuration Observation | observation | Case/Hypothesis | persistent source and age visible |
| Comparison/change relation | relation/result | CAP-INV-374/377 | preconditions and conflicts visible |
| Handoff package | candidate | CAP-INV-379 | no current-state or network claim |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Filesystem/System artifact | inspect | CAP-INV-371 | file/source, metadata, age, permissions | Navigation |
| Observation | correlate timeline | CAP-INV-374 | timestamp quality, source, conflicts | System view |
| Persistent network config | future handoff | CAP-INV-379/future Network | source/config candidate/limitations | System view |

## 18. Dépendances
CAP-INV-367/368/370/374/377/379, Endpoint Agent/Settings projections, Studio Tools, Shared Linking/Timeline, OPEN-005/008/013/015.

## 19. Source de vérité
Disk artifacts remain historical sources; Investigate owns interpretation; Settings/Endpoint sources own current/admin state. No old artifact overwrites current configuration.

## 20. Provenance et audit
Image/session/filesystem/file, category, Tool/version, timestamps/age, parsed fields as displayed, current projection used, comparison, annotations, privacy access, links and disposition.

## 21. Permissions fonctionnelles
System artifact read, configuration artifact read, sensitive settings/path read, comparison, annotation, Endpoint/Entity/Case linking and handoff prepare. Final access rules future.

## 22. Limites et erreurs
Persistent artifact ≠ current state; network config ≠ network connection. Missing/corrupt/stale artifacts, unsupported parser, conflicting snapshots, permission denial and incomplete image remain explicit.

## 23. Métriques
Artifacts by category/age/state, conflicts, comparisons, links with provenance and current-state claims prevented.

## 24. Classification de livraison
`defined` / `planned`; no platform administration, engine or implementation.

## 25. Critères d’acceptation
**Given** an old configuration artifact differing from current Endpoint projection **When** compared **Then** both sources and times remain visible and neither silently replaces the other.

**Given** a persistent network address without capture **When** observed **Then** it remains configuration context, not a confirmed connection or IOC.

**Given** no AI **When** analyzed **Then** parsers, tables, diff and human review work.

## 26. Questions ouvertes
OPEN-005/008/013/015 remain open; platform support and final observation/permission models are future.

## 27. Consommateurs documentaires
INV-DSK-001, Case/Endpoint/Entity, CAP-INV-374/377/379, future Network handoff and Objects/Permissions/Technique.
