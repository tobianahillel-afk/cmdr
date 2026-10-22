---
id: CAP-EPT-045
title: Endpoint Investigation Summary and Consumer Handoff
product: endpoint-agent
module: investigation
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-018, REQ-PROD-019, REQ-INV-001, REQ-INV-006, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-015, OPEN-017]
source-of-truth: canonical
---
# CAP-EPT-045 — Endpoint Investigation Summary and Consumer Handoff

## 1. Définition
Définir un résumé technique Endpoint limité et un handoff permission-aware contenant contexte local, observations/detections/timeline/process/file/network/session/system refs, incertitude, visibilité manquante, limitations et provenance, sans créer Evidence, Finding, Case, Decision ou Result.

## 2. Problème utilisateur
L’analyste doit transmettre le résultat d’une investigation locale à Investigate/Command sans perdre les sources ni transformer un résumé machine en conclusion canonique.

## 3. Objectifs
Produire un summary reconstructible ; inclure related refs et uncertainty ; indiquer missing visibility/collection-required ; préserver source owners et permissions ; proposer un consumer pivot conceptuel.

## 4. Non-objectifs
Aucune Evidence/Finding/Case/Incident/Decision/Result creation, aucune Action Request implicite, aucune response, aucune external publication, aucun report engine concurrent.

## 5. Propriétaire
Endpoint Agent possède le résumé local et son handoff technique. Investigate qualifie Evidence/Finding/Case ; Command qualifie les objets opérationnels ; Govern qualifie Decision/Result ; Shared conserve Reporting/Export génériques.

## 6. Utilisateurs
SOC/Investigate Analyst ; Endpoint Operator ; Command Analyst ; Detection Engineer ; Auditor ; Security/Privacy Reviewer ; Govern consumer autorisé.

## 7. Conditions d’entrée
Context Expansion EPT-3 présent ou explicitement partiel ; stable refs ; tenant/permissions ; limitations/gaps et provenance disponibles ; consumer destination identifiable.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| local detection refs/context | CAP-EPT-031..036 | detection context | oui si handoff detection | candidate freshness | summary detection partiel |
| investigation contexts | CAP-EPT-037..044 | local context | oui au moins partiellement | snapshot freshness | summary partial |
| timeline/coverage/gaps | CAP-EPT-036/042 | visibility context | non | period scoped | uncertainty explicit |
| destination/permission context | Investigate/Command/Govern + Security | consumer scope | oui | current | handoff denied/not-ready |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Endpoint local detection/investigation contexts | Endpoint Agent | stable refs/summaries | read |
| telemetry-event / timeline refs | Shared | references only | source permission |
| Case/Evidence/Finding | Investigate | destination/reference only | no creation |
| Signal/Detection/Incident | Command | destination/reference only | no creation |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Endpoint Investigation Summary | créer/rafraîchir | Endpoint Agent | summary ≠ Finding |
| Consumer Handoff Package | créer | Endpoint Agent | references, not foreign objects |
| Handoff Readiness/Restriction State | dériver | Endpoint Agent | permission/visibility explicit |

## 11. Fonctionnalités
Assembler scope/time, observations, detection refs, contextual timeline, entity/context refs, uncertainty/gaps/limitations et provenance ; minimiser données sensibles ; calculer readiness ; transmettre stable refs sans permission expansion.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect local summary | Analyst | Endpoint Investigation Summary | 0 | read | summary + limits | non |
| build deterministic handoff | service/operator | Handoff Package | 1 | contexts disponibles | package référencé | non |
| forward to authorized consumer | Analyst/Operator | handoff | 2 | destination/read scope valide | consumer ref created | non, sauf action future séparée |

Forwarding n’accorde aucun droit source et ne crée aucune autorisation de réponse.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| assembler refs/context | oui | oui | oui | non nécessaire | structured package |
| résumer contexte | oui | oui | oui | oui, attribuée | deterministic structured summary |
| proposer consumer pivot | oui | oui | oui | oui | configured destination list |
| qualifier Finding/Evidence/Decision | non | non | non | interdit | owner workflow externe |

## 14. États fonctionnels
`draft-local`, `partial`, `ready-for-handoff`, `restricted`, `stale`, `forwarded`, `destination-denied`, `collection-required`, `superseded`.

## 15. États d’interface
Aucun Screen ID. Future UI doit séparer facts, uncertainty, AI summary, visibility gaps, permission restrictions et destination ownership.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Endpoint Investigation Summary | concept Endpoint | Investigate/Command/authorized consumers | summary ≠ Finding/Evidence |
| Consumer Handoff Package | référence/projection | destination owner | stable refs + provenance |
| missing visibility/limitations | diagnostic | analyst/future Collection | no implicit acquisition |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-044 expanded context | summarize | Endpoint Summary | all relevant refs/limits | detection origin retained |
| Endpoint Summary | analyst forwards | Investigate/Command consumer | permission-aware refs | owner destination qualifies |
| destination needs new data/action | consumer request | external/future lot | explicit missing need | no execution by EPT-3 |

## 18. Dépendances
CAP-EPT-031..044 ; Investigate Case/Evidence/Finding ; Command Detection/Signal/Incident ; Govern Decision/Result ; Shared Reporting/Linking ; Security permissions ; `OPEN-008`, `OPEN-015`, `OPEN-017`.

## 19. Source de vérité
Endpoint est SOT du summary/handoff package technique. Chaque destination conserve l’ownership de ses objets et qualifications ; source permissions ne sont jamais héritées.

## 20. Provenance et audit
Summary version/time, included/excluded refs, detection/content versions, coverage/gaps, uncertainty, masking, AI attribution, destination, permission outcome, tenant, actor et correlation id.

## 21. Permissions fonctionnelles
Summary read, sensitive context read, detection rationale/content read, provenance, destination handoff permission et cross-tenant deny ; aucun manage sur Evidence/Finding/Case/Decision/Result implicite.

## 22. Limites et erreurs
Local summary ≠ Finding ; context ≠ Evidence ; forwarded ≠ accepted/qualified ; destination denial n’autorise aucune redaction bypass ; incomplete context reste incomplete.

## 23. Métriques
Summary completeness, visibility-gap count, restricted fields, handoff readiness, destination-denied rate, provenance completeness ; aucune target finale.

## 24. Classification de livraison
`draft / defined / planned`; aucun reporting/handoff runtime ou integration déployé.

## 25. Critères d’acceptation
**Given** un local investigation context est partiel, **When** le summary est construit, **Then** missing visibility et uncertainty sont inclus au lieu d’être masqués.

**Given** un summary est forwardé à Investigate, **When** il est reçu, **Then** aucun Finding/Evidence/Case n’est créé automatiquement et Investigate conserve la qualification.

**Given** l’IA est indisponible, **When** le handoff est préparé, **Then** toutes les refs, limitations et provenance sont produites par la voie structurée déterministe.

## 26. Questions ouvertes
`OPEN-015` conserve les bridges cross-product run/provenance futurs ; `OPEN-008`/`OPEN-017` restent ouvertes. Aucun bridge de response n’est finalisé.

## 27. Consommateurs documentaires
Investigate Case/Detection Engineering ; Command ; Govern comme consumer contextuel ; Shared Reporting/Linking ; Security ; Quality ; EPT-4+ futurs seulement.