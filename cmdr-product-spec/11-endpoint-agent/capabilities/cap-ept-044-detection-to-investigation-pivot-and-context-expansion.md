---
id: CAP-EPT-044
title: Detection-to-Investigation Pivot and Context Expansion
product: endpoint-agent
module: investigation
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-018, REQ-PROD-019, REQ-INV-001, REQ-INV-006, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-017]
source-of-truth: canonical
---
# CAP-EPT-044 — Detection-to-Investigation Pivot and Context Expansion

## 1. Définition
Formaliser le passage d’un Local Detection Signal Candidate vers ses observations pertinentes, contextes process/file/network/session/system, timeline et related entities, puis l’expansion locale jusqu’aux limites des données disponibles.

## 2. Problème utilisateur
Après une détection locale, l’analyste doit comprendre “quoi regarder ensuite” sans création automatique de Case/Finding/Evidence et sans qu’une expansion locale déclenche de Collection ou Live Response.

## 3. Objectifs
Fournir initial pivot, context expansion, missing/unsupported/stale states, related refs, timeline et provenance ; conserver uncertainty ; préparer un résumé local et un handoff vers Investigate.

## 4. Non-objectifs
Aucun Case automatique, aucune Evidence/Finding, aucune acquisition, aucun remote query/command, aucun response, aucune causalité automatique.

## 5. Propriétaire
Endpoint Agent possède le pivot technique local. Investigate conserve central Case investigation et qualification ; Shared conserve generic linking/search ; Govern conserve toute réponse.

## 6. Utilisateurs
SOC/Investigate Analyst ; Endpoint Operator ; Detection Engineer ; Security Reviewer ; Auditor.

## 7. Conditions d’entrée
Local candidate/match/context disponible ; related observation refs traçables ; tenant/permissions ; EPT-3 context capabilities disponibles ou limitations explicites.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| local candidate + match | CAP-EPT-032..034 | detection context | oui | candidate freshness | no pivot |
| related observation refs | CAP-EPT-017..024/034 | facts | oui au moins partiellement | source freshness | missing context |
| investigation contexts | CAP-EPT-037..043 | context projections | non | snapshot freshness | partial expansion |
| coverage/gap state | CAP-EPT-036 | visibility context | non | period scoped | uncertainty explicit |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Local Detection Candidate/Context | Endpoint Agent | source/related refs | read |
| Endpoint investigation contexts | Endpoint Agent | stable local refs | read |
| Case/Finding/Evidence | Investigate | destination concepts only | no auto creation |
| Shared links/timeline refs | Shared | generic relation mechanism | read/reference |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Detection-to-Investigation Pivot | créer/rafraîchir | Endpoint Agent | no Case creation |
| Context Expansion Set | dériver | Endpoint Agent | existing data only |
| Missing/Unsupported Context Marker | dériver | Endpoint Agent | no acquisition hidden |

## 11. Fonctionnalités
Partir du match/candidate, résoudre relevant observations, ouvrir context capabilities, construire local timeline slice, lier related entities, marquer gaps/stale/unsupported, et arrêter l’expansion dès qu’une donnée nécessite Collection.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| open detection investigation | Analyst | pivot | 0 | read | initial context | non |
| expand existing local context | Analyst/service | context set | 1 | refs disponibles | expanded set | non |
| request bounded refresh | Operator | existing context | 2 | no acquisition | refreshed local facts only | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| initial pivot | oui | oui | oui | suggestion possible | related refs |
| context expansion | oui | oui | oui | suggestions possibles | deterministic link traversal |
| proposer hypothèse | oui | non requise | oui | oui, clairement hypothesis | analyst/manual hypothesis |
| créer Finding/Evidence/Case | non | non | non | interdit | Investigate workflow explicite |

## 14. États fonctionnels
`initial`, `expanded`, `partial`, `missing-context`, `stale`, `unsupported-pivot`, `permission-denied`, `collection-required`, `ready-for-handoff`.

## 15. États d’interface
Aucun Screen ID. Future UI doit montrer missing visibility/coverage gaps, source limitations et différence entre hypothesis AI et fact.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| expanded local context | concept Endpoint | CAP-EPT-045 | source refs/uncertainty |
| related timeline/entities | projections | Analyst/Investigate | correlation ≠ causation |
| boundary/missing markers | diagnostics | Analyst/EPT-4 future | no hidden collection |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Local Detection Candidate | analyst opens | initial investigation | match/content/obs refs | detection origin retained |
| initial context | expand | process/file/network/session/system/timeline contexts | stable refs/limits | return to detection |
| expanded context | summarize/handoff | CAP-EPT-045/Investigate | context + uncertainty | no Case/Finding auto |

## 18. Dépendances
CAP-EPT-017..043 ; Shared Linking/Timeline/Search ; Investigate Case/Evidence/Finding ; Detection Engineering ; Security permissions ; `OPEN-008`, `OPEN-017`.

## 19. Source de vérité
Endpoint est SOT du pivot/expanded local context. Les facts restent leurs sources EPT-2 ; Investigate est SOT de toute central investigation/qualification.

## 20. Provenance et audit
Detection/content/match refs, initial/expanded targets, included/excluded refs, missing context, coverage/gaps, user/service, time, tenant, AI hypothesis attribution et return origin.

## 21. Permissions fonctionnelles
Local signal/context reads, sensitive process/file/network/user/system fields, timeline/provenance, cross-tenant deny ; aucun create Case/Evidence/Finding implicite.

## 22. Limites et erreurs
Context expansion ≠ acquisition ; pivot ≠ Collection ; local investigation ≠ Case ; related events ≠ causality ; missing data reste missing ; unsupported pivot ne déclenche aucun tool/run.

## 23. Métriques
Expansion depth, context completeness, missing/stale/unsupported/denied reasons, collection-required count, handoff readiness completeness.

## 24. Classification de livraison
`draft / defined / planned`; aucun investigation engine, collector ou Case automation implémenté.

## 25. Critères d’acceptation
**Given** un signal local est forwardé vers Investigation context, **When** l’analyste l’ouvre, **Then** les observations/contextes disponibles sont liés sans créer de Finding.

**Given** une expansion nécessite une donnée non locale, **When** le pivot est tenté, **Then** `collection-required` est retourné et EPT-3 s’arrête.

**Given** l’IA propose une hypothèse, **When** le contexte est présenté, **Then** l’hypothèse est clairement attribuée et aucune qualification analytique n’est créée.

## 26. Questions ouvertes
`OPEN-008`/`OPEN-017` restent ouvertes. Le bridge exact vers future Collection/Response ne fait pas partie d’EPT-3.

## 27. Consommateurs documentaires
CAP-EPT-045/046 ; Investigate Case workflows ; Detection Engineering ; Shared ; Security ; Quality ; Roadmap ; EPT-4 futur seulement.