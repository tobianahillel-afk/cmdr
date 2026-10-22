---
id: CAP-EPT-032
title: Local Detection Evaluation and Match Semantics
product: endpoint-agent
module: detection
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-012, REQ-PROD-018, REQ-PROD-019, REQ-INV-006, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-017]
source-of-truth: canonical
---
# CAP-EPT-032 — Local Detection Evaluation and Match Semantics

## 1. Définition
Définir l’évaluation locale d’observations EPT-2 contre un Detection Content éligible et la sémantique conceptuelle `matched / not-matched / inconclusive / error`, sans sélectionner moteur, langage ou représentation physique.

## 2. Problème utilisateur
Les consommateurs doivent distinguer une évaluation réellement exécutée, ses entrées, son contenu/version et ses limitations d’un simple contenu chargé, d’une conclusion malveillante ou d’un Finding.

## 3. Objectifs
Tracer observations évaluées, content/version, contexte, temps, résultat conceptuel, rationale technique, inputs incomplets, dépendances et provenance ; conserver l’évaluation sans side effect.

## 4. Non-objectifs
Aucun canonical Command Detection, Finding, Evidence, Incident, response action, compiler, AST, runtime final, ML engine final, API ou schéma physique.

## 5. Propriétaire
Endpoint Agent possède la sémantique technique de l’évaluation locale. Command conserve l’objet canonique Detection ; Investigate conserve Detection Engineering et Finding/Evidence/Case.

## 6. Utilisateurs
Endpoint Operator ; SOC/Investigate Analyst ; Detection Engineer pour feedback ; Security Reviewer ; Auditor ; Command consumer autorisé des résultats projetés.

## 7. Conditions d’entrée
Content/version référencé et éligible via CAP-EPT-031 ; observations EPT-2 disponibles ; dépendances explicites ; tenant/permission valides ; runtime concret non requis par le contrat documentaire.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| eligibility + content/version | CAP-EPT-031 | contexte d’évaluation | oui | courant | `inconclusive/not-run` |
| observation refs | CAP-EPT-015..024 | faits locaux | oui selon contenu | source freshness | input missing explicite |
| normalization/quality state | CAP-EPT-023/024 | qualité | non | liée aux observations | limitation visible |
| capability/dependency state | CAP-EPT-027/028 | availability | oui si requis | courant | evaluation blocked/inconclusive |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Local Observation | Endpoint Agent | refs/facts/time/source | read |
| telemetry-event | Shared | envelope projection si utilisée | read |
| Detection Content/version | Investigate | référence/version seulement | read |
| endpoint-agent | Endpoint Agent | runtime context facts | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Local Detection Evaluation | créer/mettre à jour état technique | Endpoint Agent | non-canonical Command Detection |
| Detection Match | dériver | Endpoint Agent | match technique, pas verdict malveillant |
| Evaluation Error/Reason | dériver | Endpoint Agent | provenance et input refs conservés |

## 11. Fonctionnalités
Valider préconditions, fixer le snapshot logique des entrées, lancer une évaluation conceptuelle déterministe/behavioral/model-eligible selon content contract, produire outcome et rationale technique, et conserver toute absence/erreur sans inférence abusive.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect evaluation | Operator/Analyst | Local Detection Evaluation | 0 | read permission | contexte/résultat visible | non |
| evaluate eligible content | deterministic local service | observations/content | 1 | eligibility + inputs | outcome technique | non |
| recompute after new local context | authorized operator/service | evaluation | 1 | mêmes règles + nouvelles refs | nouvelle assessment liée | non |

Evaluation n’exécute aucun side effect et ne constitue aucune autorisation de réponse.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| préparer inputs | oui | oui | oui | non nécessaire | résolution structurée |
| évaluer règle/condition | oui | oui | oui | seulement si content autorise un model hook | voie déterministe/content-defined |
| expliquer rationale | oui | oui | oui | oui | facts/reasons structurés |
| déclarer malicious avec autorité | non | non | non | interdit | analyst qualification externe |

## 14. États fonctionnels
`not-run`, `ready`, `evaluating`, `matched`, `not-matched`, `inconclusive`, `error`, `dependency-unavailable`, `stale`. États conceptuels, pas machine runtime finale.

## 15. États d’interface
Aucun Screen ID. Une future projection doit distinguer No match, Inconclusive, Error, Missing input, Stale et Permission denied ; `no match` ne signifie jamais benign.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| evaluation outcome | état Endpoint | CAP-EPT-033/034 | inputs/content/version/time liés |
| Detection Match technique | concept Endpoint | CAP-EPT-033/044 | pas Finding/Evidence/Incident |
| evaluation reason/error | diagnostic Endpoint | Operator/Detection Engineering | aucune donnée inventée |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-031 | eligible content | Local Evaluation | content/version/dependencies | provenance conservée |
| EPT-2 observations | evaluation | Match/outcome | exact observation refs | source refs conservées |
| Match/outcome | match ou besoin de contexte | CAP-EPT-033/034/044 | rationale + refs | aucune promotion automatique |

## 18. Dépendances
CAP-EPT-015..031 ; Endpoint `local-detection-engine.md`, `rule-model.md`, `behavioral-detection.md`, `model-inference.md` ; Investigate Detection Engineering ; `OPEN-008` ; `OPEN-017`.

## 19. Source de vérité
Endpoint possède l’évaluation locale et ses facts techniques. Command reste owner de l’objet canonical Detection et Signal. Investigate reste owner du contenu engineering et des conclusions Finding/Evidence/Case.

## 20. Provenance et audit
Content id/version/origin, eligibility assessment, Agent/version, observation ids/source/time, normalization/quality state, evaluation time, outcome, rationale, dependency errors et actor/service sont conservés.

## 21. Permissions fonctionnelles
Lire observations, contenu/rationale restreints selon politique, exécuter évaluation locale autorisée, lire provenance ; aucune permission ne permet de créer Finding/Evidence ou response. Pas de RBAC final.

## 22. Limites et erreurs
Inputs partiels, source unavailable, content stale, dependency error, platform unsupported, permission denied, model hook disabled produisent des outcomes explicites. Match ≠ malicious verdict ; no-match ≠ benign.

## 23. Métriques
Evaluations par outcome ; missing-input rate ; inconclusive/error reasons ; latency conceptuelle ; provenance completeness ; aucun objectif SLO numérique.

## 24. Classification de livraison
`draft / defined / planned` ; contrat fonctionnel documentaire, aucune preuve de runtime/engine implémenté.

## 25. Critères d’acceptation
**Given** un contenu est enabled mais l’observation ne satisfait pas ses conditions, **When** l’évaluation s’exécute, **Then** outcome = `not-matched` sans conclusion benign.

**Given** une observation satisfait le contenu éligible, **When** l’évaluation est effectuée, **Then** un Detection Match technique est produit avec content/version et observation refs, sans Finding.

**Given** une dépendance requise disparaît, **When** l’évaluation est demandée, **Then** outcome = `inconclusive/dependency-unavailable` et aucune donnée n’est inventée.

## 26. Questions ouvertes
`OPEN-017` garde runtime/language/portability/model execution indécis. `OPEN-008` conserve les limites de plateforme/source.

## 27. Consommateurs documentaires
EPT-3 detection/signals/context ; Investigate Detection Engineering ; Command projections futures ; Security/Trust ; Quality ; Roadmap ; EPT-4+ uniquement comme futurs consommateurs, non démarrés.