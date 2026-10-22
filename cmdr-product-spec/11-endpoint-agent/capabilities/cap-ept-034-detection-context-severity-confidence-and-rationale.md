---
id: CAP-EPT-034
title: Detection Context, Severity, Confidence and Rationale
product: endpoint-agent
module: detection
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-012, REQ-PROD-018, REQ-PROD-019, REQ-INV-006, REQ-AI-002, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-017]
source-of-truth: canonical
---
# CAP-EPT-034 — Detection Context, Severity, Confidence and Rationale

## 1. Définition
Définir le contexte technique attaché à une évaluation/match/candidate local : observations pertinentes, entité locale, sévérité contextuelle, confiance, rationale, incertitude et limitations, sans transformer ces éléments en impact métier, preuve ou verdict malveillant.

## 2. Problème utilisateur
Un match brut sans contexte est difficile à interpréter ; à l’inverse, une sévérité ou confiance affichée sans provenance peut être prise à tort pour un impact, une certitude ou une preuve.

## 3. Objectifs
Relier process/file/network/user/system refs déjà disponibles ; exprimer severity/confidence avec origine et méthode ; exposer rationale, missing context et uncertainty ; conserver les limites de source/plaqueforme.

## 4. Non-objectifs
Aucune qualification Finding/Evidence/Incident, aucun business impact automatique, aucune probabilité obligatoire, aucune causalité, aucun scoring engine final ni modèle ML final.

## 5. Propriétaire
Endpoint Agent possède le contexte technique local dérivé. Investigate conserve les conclusions et confidence analytique de Findings ; Command conserve la priorité/Signal/Detection canonique ; Govern conserve l’impact/autorité de réponse.

## 6. Utilisateurs
Endpoint Operator ; SOC/Investigate Analyst ; Detection Engineer ; Command consumer ; Security/Privacy Reviewer ; Auditor.

## 7. Conditions d’entrée
Evaluation/match ou Local Detection Signal Candidate sourcé ; observations/context refs accessibles ; tenant/permissions valides ; limitations et source quality disponibles ou explicitement unknown.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| match/evaluation | CAP-EPT-032 | résultat technique | oui | evaluation time | contexte non calculable |
| local candidate | CAP-EPT-033 | relation technique | non | candidate freshness | contexte attaché au match seulement |
| process/file/network/auth/system refs | CAP-EPT-017..021 | observations | non selon contenu | source freshness | missing context explicite |
| content severity/confidence hints | Detection Content | métadonnée référencée | non | version liée | pas d’inférence |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Local Detection Evaluation/Match | Endpoint Agent | outcome/rationale refs | read |
| Local Observation | Endpoint Agent | process/file/network/session/system facts | read selon permission |
| Detection Content/version | Investigate | severity/confidence intent si sourcé | read |
| canonical Finding/Signal | Investigate / Command | destination/reference seulement | aucune mutation |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Detection Context | dériver/rafraîchir | Endpoint Agent | facts et limites explicites |
| Contextual Severity | dériver/projeter | Endpoint Agent | severity ≠ impact |
| Detection Confidence | dériver/projeter | Endpoint Agent | confidence ≠ certainty/probability obligatoire |
| Detection Rationale | dériver | Endpoint Agent | rationale ≠ proof |

## 11. Fonctionnalités
Sélectionner les observations pertinentes, appliquer une méthode de contexte documentée, lier entités locales, conserver les éléments manquants, représenter severity/confidence/rationale et leur provenance sans dépasser la portée du Detection Content ou des données disponibles.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect detection context | Analyst/Operator | Detection Context | 0 | read permission | contexte et limites visibles | non |
| recompute context | deterministic local service | context | 1 | nouvelles observations disponibles | nouvelle projection liée | non |
| inspect restricted rationale | authorized reviewer | rationale | 0 | permission spécifique | champs autorisés seulement | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| lier observations | oui | oui | oui | suggestion possible | keys/refs déterministes |
| calculer contexte/sévérité selon contenu | oui | oui | oui | explication seulement | règles/métadonnées sourcées |
| résumer rationale | oui | oui | oui | oui, attribuée | rationale structurée |
| déclarer causalité/malicious | non | non | non | interdit | analyst qualification externe |

## 14. États fonctionnels
`complete`, `partial`, `restricted`, `uncertain`, `stale`, `unsupported`, `unknown`. Severity et confidence sont valeurs contextualisées, pas des états d’autorité.

## 15. États d’interface
Aucun Screen ID. Les projections futures doivent rendre uncertainty/missing context/restricted rationale visibles ; un badge severity ne doit pas être présenté comme impact métier.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Detection Context | concept Endpoint | CAP-EPT-035..046 | source refs + limitations |
| contextual severity/confidence | projection Endpoint | Analyst/Command/Investigate | méthode/origine visibles |
| rationale/uncertainty | projection Endpoint | analyst/reviewer | aucune preuve inventée |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-032/033 | match/candidate | Detection Context | content + observation refs | provenance conservée |
| EPT-2 observation | contexte additionnel | context refresh | old/new refs | limitation history retained |
| Detection Context | pivot/handoff | EPT-3 investigation/consumer | severity/confidence/rationale/uncertainty | aucune qualification auto |

## 18. Dépendances
CAP-EPT-017..033 ; `rule-model.md`, `behavioral-detection.md`, `local-correlation.md` ; Investigate Detection Engineering ; Security privacy/minimization ; `OPEN-008`, `OPEN-017`.

## 19. Source de vérité
Endpoint est source des facts de contexte technique local. Les source observations restent EPT-2/Shared-owned selon leur contrat. Investigate/Command/Govern conservent leurs propres sévérités, impacts et conclusions canoniques.

## 20. Provenance et audit
Conserver match/candidate refs, observations retenues/écartées, method/version de contexte, severity/confidence origin, rationale, missing fields, masking, actor/service, time et tenant.

## 21. Permissions fonctionnelles
Lire signal candidate, observation refs, command metadata restreinte, user/session data, rationale et provenance selon moindre privilège ; deny cross-tenant ; aucun droit implicite Finding/Evidence/response.

## 22. Limites et erreurs
Contexte incomplet, données stale, champ masqué, source unsupported ou permission denied restent explicites. Severity ≠ impact ; confidence ≠ certainty ; rationale ≠ proof ; related observation ≠ malicious relationship.

## 23. Métriques
Context completeness, missing/refused fields, confidence-source distribution, rationale completeness, stale/unsupported reasons. Pas de seuil universel non approuvé.

## 24. Classification de livraison
`draft / defined / planned` ; aucune preuve d’implémentation, scoring engine ou modèle livré.

## 25. Critères d’acceptation
**Given** un match manque de contexte requis, **When** le contexte est construit, **Then** il reste `partial` et les champs manquants sont nommés.

**Given** une confidence faible est sourcée, **When** l’analyste consulte le match, **Then** elle est affichée comme confidence/uncertainty et non comme certitude ou probabilité obligatoire.

**Given** une relation process/file existe, **When** le rationale est produit, **Then** la relation est décrite sans être présentée comme preuve malveillante.

## 26. Questions ouvertes
`OPEN-017` garde modèle/runtime/langage indécis ; `OPEN-008` garde source/platform support indécis. Aucune échelle universelle d’impact n’est choisie.

## 27. Consommateurs documentaires
EPT-3 grouping/coverage/investigation/handoff ; Investigate Detection Engineering et Case workflows comme consommateurs ; Command ; Security ; Quality ; Roadmap.