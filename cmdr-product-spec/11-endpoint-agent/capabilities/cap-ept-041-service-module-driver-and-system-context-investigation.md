---
id: CAP-EPT-041
title: Service, Module, Driver and System Context Investigation
product: endpoint-agent
module: investigation
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-012, REQ-PROD-018, REQ-PROD-019, REQ-INV-001, REQ-INV-006, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008]
source-of-truth: canonical
---
# CAP-EPT-041 — Service, Module, Driver and System Context Investigation

## 1. Définition
Définir l’investigation locale du contexte services, modules, drivers et activité système déjà observés, y compris les relations process/file et les persistence-like observations, sans conclure automatiquement à une persistance malveillante.

## 2. Problème utilisateur
Des services, modules, drivers, tasks ou autoruns peuvent être légitimes ou malveillants. L’analyste doit voir origin/signature/chronology/relations et limitations de plateforme sans qu’une observation soit promue en conclusion.

## 3. Objectifs
Relier service/module/driver/system refs, origin/signature facts, process/file relations, related candidates, chronology, persistence-like context et platform-specific coverage ; exposer kernel restrictions et gaps.

## 4. Non-objectifs
Aucun persistence Finding, aucune modification service/driver/task/registry, aucun kernel inspection supplémentaire, aucun module dump, aucun response, aucun support Windows/Linux/macOS universel.

## 5. Propriétaire
Endpoint Agent possède le System Context local dérivé. Investigate conserve Findings/Evidence et toute conclusion persistence ; Settings conserve configuration/policies ; Govern conserve response authority.

## 6. Utilisateurs
SOC/Investigate Analyst ; Endpoint Operator ; Detection Engineer ; Security Reviewer ; Auditor ; Command consumer autorisé.

## 7. Conditions d’entrée
System/service/module/driver observations disponibles ; permissions ; platform/source support déclaré ; related process/file/detection refs optionnels ; aucun fetch supplémentaire équivalent à Collection.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| system/service/module/driver observations | CAP-EPT-021 | faits locaux | oui | source freshness | context unavailable |
| process/file context | CAP-EPT-017/018/037/038 | relations | non | snapshot/source time | relation unknown |
| local detection refs | CAP-EPT-032..035 | signal context | non | candidate freshness | system context only |
| platform/source limitation | CAP-EPT-027/028 | availability | oui | current | unknown/unsupported |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| System/Service/Module/Driver Observation | Endpoint Agent | refs/state/origin/signature/time | read |
| Process/File Context | Endpoint Agent | relations | read |
| Local Detection Candidate | Endpoint Agent | related refs | read |
| Finding/Evidence | Investigate | destination only | no local creation |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| System Investigation Context | dériver/rafraîchir | Endpoint Agent | observed facts only |
| Persistence-Like Context Reference | dériver | Endpoint Agent | observation ≠ persistence conclusion |
| System Chronology Projection | dériver | Endpoint Agent | source/platform limits visibles |

## 11. Fonctionnalités
Associer services/modules/drivers/tasks/registry-like refs existantes aux process/file observations et related signals, afficher signature/origin quand disponibles, construire chronology et qualifier uniquement les limitations techniques.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect system context | Analyst | System Context | 0 | read | facts/relations visibles | non |
| correlate existing system refs | deterministic service | context | 1 | observations locales | liens sourcés | non |
| request bounded status refresh | Operator | existing component state | 2 | sans mutation/acquisition | freshness update | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| lier component/process/file | oui | oui | oui | suggestion possible | stable refs/time |
| construire chronology | oui | oui | oui | non nécessaire | timestamps |
| résumer persistence-like context | oui | oui | oui | oui, comme hypothèse/context | structured facts |
| conclure persistence/malicious | non | non | non | interdit | Investigate qualification |

## 14. États fonctionnels
`available`, `partial`, `platform-dependent`, `restricted`, `stale`, `unsupported`, `unavailable`, `unknown`.

## 15. États d’interface
Aucun Screen ID. Platform-specific sources et kernel restrictions restent visibles ; persistence-like label doit être présenté comme contexte, jamais Finding.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| System Investigation Context | concept Endpoint | CAP-EPT-042..045 | source refs/limits |
| component/process/file relations | projections Endpoint | Analyst | observation ≠ malicious relation |
| persistence-like context | contexte Endpoint | Investigate | no conclusion/promotion |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| system observation/candidate | pivot | System Context | component/process/file refs | origin retained |
| System Context | process/file/timeline pivot | CAP-EPT-037/038/042 | stable refs/time | no collection |
| deeper artifact/kernel data required | unavailable local data | EPT-3 stop | missing requirement | future Collection only |

## 18. Dépendances
CAP-EPT-017/018/021/027/028/032..038 ; `investigation/modules-and-drivers.md`; `persistence.md`; EPT-2 system/registry sources ; Shared Timeline/Linking ; `OPEN-008`.

## 19. Source de vérité
Endpoint est SOT du contexte local source-backed ; Investigate reste SOT des Findings/Evidence et conclusions persistence ; Settings reste SOT de l’administration.

## 20. Provenance et audit
Component/service/module/driver ids, origin/signature source, platform/source capability, process/file relations, candidate refs, timestamps, limitations, tenant et actor/service sont conservés.

## 21. Permissions fonctionnelles
System/service/module/driver context read, restricted kernel/signature/path metadata read, provenance, cross-tenant deny. Aucun manage/stop/remove/modify implicite.

## 22. Limites et erreurs
Service/driver observation ≠ persistence proof ; module relation ≠ malicious relationship ; platform-specific coverage peut être partial ; absent source ne signifie pas absent component.

## 23. Métriques
Context completeness, signature/origin availability, platform limitations, related process/file coverage, stale/unsupported reasons, pivot success.

## 24. Classification de livraison
`draft / defined / planned`; aucune mutation système, collector ou kernel engine implémenté.

## 25. Critères d’acceptation
**Given** un driver load est lié à un signal, **When** System Context est construit, **Then** origin/signature/process refs disponibles sont montrés sans conclure rootkit/persistence.

**Given** une source registry est unsupported sur la plateforme, **When** le contexte est consulté, **Then** la limitation est explicite et aucune visibilité universelle n’est revendiquée.

**Given** des données kernel plus profondes seraient nécessaires, **When** l’analyste demande expansion, **Then** EPT-3 s’arrête avant acquisition.

## 26. Questions ouvertes
`OPEN-008` reste ouverte ; les détails de support/kernel sources et toute réponse future restent hors EPT-3.

## 27. Consommateurs documentaires
EPT-3 timeline/pivots/summary ; Investigate Case/Finding ; Command ; Security ; Quality ; future EPT-4/EPT-5 boundaries.