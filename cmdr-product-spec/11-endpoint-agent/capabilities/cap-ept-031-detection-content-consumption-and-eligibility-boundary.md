---
id: CAP-EPT-031
title: Detection Content Consumption and Eligibility Boundary
product: endpoint-agent
module: detection
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-018, REQ-PROD-019, REQ-INV-006, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-017]
source-of-truth: canonical
---
# CAP-EPT-031 — Detection Content Consumption and Eligibility Boundary

## 1. Définition
Définir comment un Endpoint Agent référence et évalue l’éligibilité locale d’un Detection Content/version produit par les sources canoniques, sans devenir propriétaire de l’authoring ni choisir un runtime, langage ou format exécutable final.

## 2. Problème utilisateur
Un contenu de détection peut exister sans être applicable à un endpoint précis. Les opérateurs doivent comprendre version, origine, dépendances, observations requises, conditions de plateforme et raisons d’inéligibilité sans confondre contenu disponible, contenu activé et contenu exécutable localement.

## 3. Objectifs
- conserver une référence stable vers Detection Content et sa version ;
- évaluer localement applicabilité et préconditions à partir des faits EPT-1/EPT-2 ;
- exposer explicitement dépendances, limitations et raisons d’inéligibilité ;
- préserver les owners Investigate, Settings et Govern.

## 4. Non-objectifs
Aucun authoring, publication, déploiement, activation administrative, compilation, moteur, langage final, package technique, API ou choix Sigma/YARA. Cette capability ne crée ni canonical Detection ni canonical Signal.

## 5. Propriétaire
Endpoint Agent possède uniquement la sémantique de consommation et d’éligibilité technique locale. Investigate conserve Detection Engineering/authoring/validation ; Settings conserve configuration, runtimes, targets, assignments et policies ; Govern conserve l’autorité de changement de production.

## 6. Utilisateurs
Endpoint Operator ; SOC/Investigate Analyst ; Detection Engineer comme consommateur de feedback ; Platform Administrator pour la projection administrative ; Security/Privacy Reviewer ; Auditor.

## 7. Conditions d’entrée
Endpoint Agent identifié et tenant-bound ; référence de contenu/version traçable ; état de plateforme/version/capabilities connu ou explicitement inconnu ; observations/dépendances déclarées ; permissions de lecture ; `OPEN-008` et `OPEN-017` restent ouvertes.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Detection Content + version reference | Investigate Detection Engineering | référence/version | oui | version courante connue | eligibility `unknown` |
| platform/version/capability facts | CAP-EPT-004/005/027/028 | faits techniques | oui | source freshness | `unsupported`/`unknown` explicite |
| required observation/source state | CAP-EPT-015..026 | dépendance télémétrique | oui selon contenu | source freshness | `dependency-unavailable` |
| enabled/assignment projection | Platform Settings si disponible | contexte administratif | non | source version | ne pas inférer `enabled` |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Detection Content / version concept | Investigate Detection Engineering | id, version, requirements, limitations | lecture référencée |
| endpoint-agent | Endpoint Agent | platform/version/capability state | lecture locale |
| endpoint-policy / assignment projection | Platform Settings | activation/assignment si sourcée | lecture seulement |
| telemetry-event / Local Observation refs | Shared / Endpoint | disponibilité des inputs | lecture permissionnée |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Detection Content Runtime Reference | dériver/rafraîchir | Endpoint Agent | référence locale, pas copie d’ownership |
| Detection Eligibility Assessment | dériver | Endpoint Agent | conditions et raisons explicites |
| Dependency Requirement Projection | dériver | Endpoint Agent | aucun changement de source |

## 11. Fonctionnalités
Résoudre une version de contenu, vérifier les conditions déclarées, comparer observations/capabilities disponibles aux dépendances requises, exposer eligibility et enabled-state seulement si sourcé, et conserver toute limitation de plateforme/version/source.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect content eligibility | Endpoint Operator / Analyst | content reference | 0 | read + tenant scope | raisons visibles | non |
| recompute eligibility | Endpoint deterministic service | eligibility assessment | 1 | facts disponibles | nouvel état sourcé | non |
| inspect enabled projection | Platform Admin / Analyst | assignment projection | 0 | source Settings accessible | projection seulement | non |

Aucune action de classe 3/4. Eligible ≠ enabled ; enabled ≠ applicable à chaque observation.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| résoudre référence/version | oui | oui | oui | non nécessaire | résolution déterministe |
| calculer eligibility | oui | oui | oui | explication seulement | règles de préconditions |
| expliquer limitation | oui | oui | oui | oui, attribuée | codes de raisons structurés |
| choisir runtime/langage | non | non | non | interdit | `OPEN-017` ouvert |

## 14. États fonctionnels
`referenced`, `eligible`, `ineligible`, `dependency-unavailable`, `platform-unsupported`, `version-incompatible`, `unknown`, `stale`; `enabled/disabled` n’est qu’une projection administrative lorsqu’une source canonique la fournit.

## 15. États d’interface
Aucun Screen ID Endpoint. Toute future surface doit montrer Loading/Partial/Stale/Permission denied/Unavailable et les raisons sans masquer la différence eligibility/enablement/applicability.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| content runtime reference | projection Endpoint | CAP-EPT-032, opérateur | owner/version/origin préservés |
| eligibility assessment | état technique Endpoint | CAP-EPT-032/036, Investigate | raisons et dépendances explicites |
| limitation/dependency projection | état technique Endpoint | Settings/Investigate | aucune promesse de support |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Detection Content/version | référence reçue | eligibility assessment | owner/version/requirements | provenance conservée |
| EPT-2 source/capability facts | changement de dépendance | reassessment | source/platform/reason | ancien/nouvel état |
| eligibility assessment | consommation locale | CAP-EPT-032 | content ref + reasons | aucun ownership transfer |

## 18. Dépendances
CAP-EPT-004/005/011/015..030 ; Investigate Detection Engineering `CAP-INV-401..435` ; Settings Endpoint Policy/targets/configured runtime projections ; `OPEN-008` ; `OPEN-017`.

## 19. Source de vérité
Investigate est source de vérité du Detection Content et de son cycle d’authoring/validation. Settings est source administrative des assignments/configuration/runtime targets. Endpoint est source de vérité uniquement de l’éligibilité technique locale dérivée.

## 20. Provenance et audit
Conserver content id/version/owner/origin, Agent id/version, tenant/environment, source/capability states, assessment time, reasons, limitations, actor/service et correlation reference. Local audit reste distinct de Shared Trace.

## 21. Permissions fonctionnelles
Besoins : lecture référence Detection Content, lecture éventuelle de contenu restreint, lecture platform/capability/source state, lecture assignment/policy projection, lecture provenance, rejet cross-tenant. Aucun RBAC/ABAC final n’est créé.

## 22. Limites et erreurs
Référence inconnue, version retirée, source manquante, plateforme inconnue, version incompatible, permission refusée ou données stale produisent un état explicite. `OPEN-017` non résolue n’empêche pas le contrat fonctionnel.

## 23. Métriques
Distribution eligible/ineligible/unknown ; raisons de dépendance ; fraîcheur ; proportion de contenu avec exigences résolues ; completeness de provenance. Aucun SLO numérique non approuvé.

## 24. Classification de livraison
`draft / defined / planned`. Spécification documentaire seulement ; aucune preuve d’implémentation, de runtime ou de support plateforme livré.

## 25. Critères d’acceptation
**Given** un Detection Content n’est pas applicable à la plateforme observée, **When** eligibility est calculée, **Then** il est `platform-unsupported/ineligible` sans exécution.

**Given** un contenu exige une télémétrie indisponible, **When** les dépendances sont évaluées, **Then** l’état indique `dependency-unavailable` et ne déclare pas le contenu matched.

**Given** l’IA est indisponible, **When** eligibility est calculée, **Then** le résultat complet est produit déterministement avec ses reasons.

## 26. Questions ouvertes
`OPEN-008` reste ouverte pour support plateforme/source. `OPEN-017` reste ouverte pour runtime, target language, portability, executable representation et moteur final.

## 27. Consommateurs documentaires
EPT-3 ; Endpoint Detection/Investigation ; Investigate Detection Engineering ; Platform Settings ; Command comme futur consommateur de projections ; Security/Trust ; Quality ; Roadmap Phase 5.