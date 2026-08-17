---
id: CAP-EPT-043
title: Endpoint Investigation Context Retrieval and Pivot Semantics
product: endpoint-agent
module: investigation
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-018, REQ-PROD-019, REQ-INV-001, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008]
source-of-truth: canonical
---
# CAP-EPT-043 — Endpoint Investigation Context Retrieval and Pivot Semantics

## 1. Définition
Définir les pivots read-only entre process, parent/child, file, network, user/session, system component, local detection candidate, source et timeline en utilisant uniquement les données déjà disponibles dans EPT-3.

## 2. Problème utilisateur
L’analyste doit naviguer rapidement entre contextes liés sans que “pivot”, “search” ou “retrieve” déclenche par erreur une collecte, un Live Response ou un accès cross-tenant.

## 3. Objectifs
Résoudre stable refs, préserver origin/tenant/permissions, expliciter unsupported/missing/stale pivots, garder la provenance et identifier clairement quand une expansion exigerait future Collection.

## 4. Non-objectifs
Aucune generic Shared Search ownership, aucun remote query shell, aucune Collection, aucun Live Response, aucun deep scan, aucun cross-tenant privilege expansion.

## 5. Propriétaire
Endpoint Agent possède les semantics des pivots au sein de son contexte local. Shared conserve Global Search/Object Linking ; Investigate conserve central investigation/navigation de Case.

## 6. Utilisateurs
SOC/Investigate Analyst ; Endpoint Operator ; Detection Engineer ; Auditor ; Security Reviewer.

## 7. Conditions d’entrée
Source context et target ref existants ; target data déjà localement disponible ; tenant/permission valides ; freshness visible ; return origin conservable.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| source context/ref | CAP-EPT-037..042 | context | oui | source freshness | pivot impossible |
| target stable ref | Endpoint/Shared linking projection | reference | oui | ref validity | unsupported/not-found |
| permission/tenant scope | Security/Settings | access context | oui | current | denied |
| local availability | CAP-EPT-027/028 | capability/data state | oui | current | boundary-to-Collection reason |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Endpoint local contexts | Endpoint Agent | stable refs/summary | read |
| Object Link/Timeline refs | Shared | navigation mechanism | permission-aware reference |
| Case/Evidence/Finding | Investigate | external destination refs | source permission required |
| endpoint-agent | Endpoint Agent | tenant/availability | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Investigation Pivot | créer état/navigation contextuel | Endpoint Agent | no source mutation |
| Pivot Resolution Result | dériver | Endpoint Agent | available/missing/denied reason |
| Return-Origin Context | préserver | Endpoint Agent | navigation ne transfère pas permission |

## 11. Fonctionnalités
Résoudre localement target refs, appliquer tenant/permission, charger uniquement context déjà disponible, transmettre time/source context, retourner explicitement missing/unsupported/stale/denied, et marquer “Collection required” comme frontière non exécutée.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| pivot to related local context | Analyst | Investigation Pivot | 0 | target available/read | context opened | non |
| resolve contextual links | deterministic service | refs | 1 | stable keys/permissions | result list | non |
| request bounded local refresh | Operator | existing context | 2 | no acquisition/query execution | freshness only | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| resolve direct pivot | oui | oui | oui | non nécessaire | stable refs |
| suggest related pivots | oui | oui | oui | oui, suggestion | deterministic relation list |
| explain unavailable pivot | oui | oui | oui | oui | reason codes |
| turn pivot into collection silently | non | non | non | interdit | explicit future boundary |

## 14. États fonctionnels
`available`, `missing`, `stale`, `unsupported`, `permission-denied`, `cross-tenant-denied`, `collection-required`, `unknown`.

## 15. États d’interface
Aucun Screen ID. Future UI doit distinguer unavailable/denied/collection-required et conserver return origin ; navigation ne révèle aucun summary non autorisé.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| resolved local context | Endpoint projection | Analyst | only existing data |
| pivot reason/result | state Endpoint | Analyst/Quality | no silent acquisition |
| future-collection boundary marker | diagnostic | EPT-4 future | request non exécutée |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| process/file/network/session/system/timeline | user pivot | related EPT-3 context | stable ref/time/tenant | return origin |
| local candidate | inspect context | CAP-EPT-044 | candidate + related refs | no Case creation |
| missing data | expansion requested | EPT-3 stop | missing requirement | future Collection only |

## 18. Dépendances
CAP-EPT-027/028/037..042 ; Shared Global Search/Object Linking/Timeline ; Investigate navigation boundaries ; Security permission model ; `OPEN-008`.

## 19. Source de vérité
Endpoint est SOT du pivot contextuel local ; target source conserve ses facts/permissions. Shared possède generic Search/Linking et Investigate possède Case navigation.

## 20. Provenance et audit
Source/target refs, tenant, permission outcome, resolver rule, time, freshness, return origin, missing/denied reason et actor/service sont conservés.

## 21. Permissions fonctionnelles
Read par context family, sensitive command/file/network/user/system fields, local timeline, provenance ; cross-tenant denial obligatoire ; aucune permission de collection/response implicite.

## 22. Limites et erreurs
Pivot ≠ Collection ; context expansion ≠ forensic acquisition ; investigation query ≠ Live Response ; local retrieval ≠ Shared Search ownership ; denied source ne révèle pas ses détails.

## 23. Métriques
Pivot resolution success, missing/stale/denied/collection-required reasons, cross-tenant denials, return-origin completeness.

## 24. Classification de livraison
`draft / defined / planned`; pas de query engine, collector ou UX finale implémentée.

## 25. Critères d’acceptation
**Given** un analyste pivote vers un parent process déjà observé, **When** la ref est résolue, **Then** le contexte local existant est retourné avec origin/permissions.

**Given** le pivot nécessite des bytes fichier non disponibles, **When** l’expansion est demandée, **Then** résultat = `collection-required` et EPT-3 ne lance rien.

**Given** un target appartient à un autre tenant, **When** le pivot est résolu, **Then** l’accès est refusé sans fuite de contexte.

## 26. Questions ouvertes
`OPEN-008` reste ouverte pour disponibilité des contextes. Les opérations de future Collection appartiennent à EPT-4.

## 27. Consommateurs documentaires
EPT-3 detection-to-investigation/summary ; Investigate ; Shared navigation contracts ; Security ; Quality ; future EPT-4 boundary.