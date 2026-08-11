---
id: CAP-EPT-042
title: Endpoint Timeline and Activity Correlation Investigation
product: endpoint-agent
module: investigation
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-012, REQ-PROD-018, REQ-PROD-019, REQ-INV-001, REQ-INV-006, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-017]
source-of-truth: canonical
---
# CAP-EPT-042 — Endpoint Timeline and Activity Correlation Investigation

## 1. Définition
Définir une contextual timeline locale et la corrélation d’activité à partir des observations/contextes EPT-2/EPT-3 existants, en conservant timestamps, time skew, late events, gaps, entity/signal relations et provenance, sans revendiquer le generic Timeline/Correlation engine Shared.

## 2. Problème utilisateur
Les événements endpoint arrivent de sources différentes et parfois en retard. L’analyste doit reconstruire un ordre contextualisé tout en sachant qu’ordre et corrélation ne prouvent pas la causalité.

## 3. Objectifs
Assembler event/context refs ; distinguer observed/inferred link ; représenter source/observation time, late/gap/time skew ; relier process/session/files/network/system et candidates ; produire une timeline locale permission-aware.

## 4. Non-objectifs
Aucun causal graph proof, aucun generic Timeline/Linking/Correlation ownership, aucune central Case Timeline, aucune acquisition, aucun SIEM/correlation engine final.

## 5. Propriétaire
Endpoint Agent possède la contextual timeline locale pour ses facts. Shared conserve les generic Timeline/Linking/Correlation primitives ; Investigate conserve Case Timeline et raisonnement central.

## 6. Utilisateurs
SOC/Investigate Analyst ; Endpoint Operator ; Detection Engineer ; Auditor ; Security Reviewer ; Command consumer.

## 7. Conditions d’entrée
Observations/context refs permissionnés avec timestamps/provenance ; clock health/time skew connu si possible ; tenant scope ; gaps/late state EPT-2 conservés.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| observations + time semantics | CAP-EPT-016..024 | events/facts | oui | source freshness | timeline partial |
| detection candidates/context | CAP-EPT-032..036 | local refs | non | candidate freshness | timeline sans detection links |
| process/file/network/user/system contexts | CAP-EPT-037..041 | contextual refs | non | snapshot freshness | relation omitted |
| clock/gap/late state | EPT-2 quality + host timeline | quality | non | period-scoped | limitation explicit |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| telemetry-event / timeline-entry refs | Shared | generic refs/time semantics | read/reference |
| Local Observation/Contexts | Endpoint Agent | facts/relations | read |
| Local Detection Candidate | Endpoint Agent | candidate refs | read |
| Case Timeline | Investigate | destination/handoff only | no local ownership |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Endpoint Contextual Timeline | dériver/rafraîchir | Endpoint Agent | local projection, not Shared engine |
| Local Correlation Relation | dériver | Endpoint Agent | correlation ≠ causation |
| Timeline Gap/Late Marker | dériver | Endpoint Agent | quality state retained |

## 11. Fonctionnalités
Ordonnancer les refs selon timestamps disponibles, annoter clock skew/late/gaps, relier entities/process/session/file/network/system/candidates par clés source-backed, distinguer relation observée vs dérivée et préparer handoff vers Investigate.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect local timeline | Analyst | Endpoint Timeline | 0 | read | chronology + gaps | non |
| correlate existing facts | deterministic service | local refs | 1 | source refs | relations contextualisées | non |
| request bounded timeline refresh | Operator | existing local data | 2 | aucune nouvelle acquisition | re-evaluation only | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| ordering | oui | oui | oui | non nécessaire | timestamps/time rules |
| entity/context linking | oui | oui | oui | suggestion possible | stable refs/keys |
| timeline summary | oui | oui | oui | oui | ordered event list |
| déclarer causalité | non | non | non | interdit | relation type + uncertainty |

## 14. États fonctionnels
`available`, `partial`, `gap-present`, `late-events`, `clock-uncertain`, `stale`, `restricted`, `unknown`.

## 15. États d’interface
Aucun Screen ID. Late/gap/clock uncertainty doivent être visibles ; la ligne temporelle ne doit pas suggérer une causalité non sourcée.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Endpoint Contextual Timeline | projection Endpoint | CAP-EPT-043..045/Investigate | source/time/limits |
| correlation relations | projections Endpoint | Analyst | correlation ≠ causation |
| gap/late markers | quality context | Analyst/Detection Engineering | visibility limitations explicites |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| observations/contexts | timeline build | Endpoint Timeline | refs/timestamps/quality | origin retained |
| Timeline event/entity | pivot | CAP-EPT-043 | target ref + time context | return origin retained |
| Endpoint Timeline | handoff | Investigate Case Timeline context | stable refs/gaps/uncertainty | no ownership transfer |

## 18. Dépendances
CAP-EPT-016..041 ; `investigation/host-timeline.md`; `detection/local-correlation.md`; Shared Timeline Engine/Object Linking/Correlation Engine ; Investigate Case Timeline ; `OPEN-008`, `OPEN-017`.

## 19. Source de vérité
Endpoint est SOT de cette projection locale ; chaque observation reste owner de ses facts ; Shared reste SOT des primitives génériques ; Investigate reste SOT de la timeline de Case.

## 20. Provenance et audit
Event/context ids, source/observation/ingestion times, ordering rule/version, clock skew, late/gap markers, relation type/source, tenant et actor/service sont conservés.

## 21. Permissions fonctionnelles
Local timeline read, sensitive event/context read selon source, correlation relation read, provenance, cross-tenant deny. Aucune permission Shared/Investigate n’est héritée par navigation.

## 22. Limites et erreurs
Timeline ≠ causal proof ; correlation ≠ causation ; late events peuvent réordonner la vue ; gaps réduisent la confiance ; no event ≠ no activity ; local timeline ≠ central Case Timeline.

## 23. Métriques
Timeline completeness, late/gap counts, clock uncertainty, correlation relation coverage, restricted events, refresh/pivot resolution.

## 24. Classification de livraison
`draft / defined / planned`; aucun timeline/correlation engine technique n’est sélectionné.

## 25. Critères d’acceptation
**Given** un événement arrive en retard, **When** la timeline est reconstruite, **Then** son event time et late state sont conservés et l’ordre peut être réévalué sans falsifier l’historique.

**Given** plusieurs événements sont corrélés, **When** l’analyste les inspecte, **Then** la relation est présentée comme corrélation et non causalité prouvée.

**Given** un gap de télémétrie existe, **When** la timeline est résumée, **Then** le gap reste visible et le résumé n’invente aucun événement.

## 26. Questions ouvertes
`OPEN-008` et `OPEN-017` restent ouvertes. Le choix d’un engine générique ou SIEM est hors EPT-3.

## 27. Consommateurs documentaires
EPT-3 pivots/handoff/summary ; Investigate Case Timeline ; Shared capability contracts ; Detection Engineering ; Security ; Quality ; Roadmap.