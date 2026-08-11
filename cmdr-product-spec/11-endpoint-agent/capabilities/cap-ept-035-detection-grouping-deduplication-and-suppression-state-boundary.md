---
id: CAP-EPT-035
title: Detection Grouping, Deduplication and Suppression-State Boundary
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
# CAP-EPT-035 — Detection Grouping, Deduplication and Suppression-State Boundary

## 1. Définition
Définir la sémantique locale de répétition, grouping, duplicate-candidate et la projection d’un état de suppression sourcé, tout en conservant chaque match/provenance et l’owner administratif de la suppression.

## 2. Problème utilisateur
Les détections répétées peuvent produire du bruit ; sans frontière normative, déduplication ou suppression pourrait effacer des faits, masquer l’historique ou laisser Endpoint inventer une policy qu’il ne possède pas.

## 3. Objectifs
Identifier repeats/duplicate candidates ; appliquer des grouping keys conceptuelles ; conserver recurrence et tous les refs ; projeter suppression scope/owner/reason/expiry ; interdire tenant leakage et suppression destructive.

## 4. Non-objectifs
Aucune policy de suppression finale, aucun delete de Match/Signal, aucune activation/deactivation de content, aucune mutation Settings/Investigate/Govern, aucun moteur de grouping final.

## 5. Propriétaire
Endpoint Agent possède uniquement le grouping/dedup technique local et l’état projeté. L’owner de suppression administrative reste celui de la source canonique ; Investigate Detection Engineering conserve les propositions de suppression/exception et Settings les configurations administratives applicables.

## 6. Utilisateurs
Endpoint Operator ; SOC/Investigate Analyst ; Detection Engineer ; Platform Administrator ; Auditor ; Security Reviewer.

## 7. Conditions d’entrée
Matches/candidates traçables ; tenant/time/content refs connus ; grouping inputs suffisants ; suppression-state uniquement si une source autorisée fournit owner/scope/reason/expiry.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| matches/candidates | CAP-EPT-032/033 | faits techniques | oui | event/evaluation time | aucun grouping |
| context/entity keys | CAP-EPT-034 + local-correlation | contexte | non | source freshness | grouping limité |
| suppression projection | source administrative/engineering autorisée | policy/state ref | non | source version/expiry | état non suppressed |
| tenant/time window facts | Endpoint/Settings | scope | oui | courant | grouping refusé |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Detection Match / local candidate | Endpoint Agent | ids/context/time | read |
| suppression/exception proposal or policy | Investigate/Settings selon source | scope/reason/expiry | read only |
| endpoint-agent | Endpoint Agent | tenant/Agent scope | read |
| canonical Signal | Command | destination/reference | no local mutation |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Grouping Relation | créer/mettre à jour | Endpoint Agent | membres conservés |
| Duplicate Candidate | dériver | Endpoint Agent | duplicate ≠ erased |
| Suppression-State Projection | dériver/rafraîchir | Endpoint Agent | owner/reason/scope/expiry obligatoires si suppressed |

## 11. Fonctionnalités
Comparer des matches selon content/version, entity/context/time keys ; créer une relation de groupe ou candidate duplicate ; conserver recurrence ; appliquer uniquement la visibilité d’une suppression autorisée et auditable ; expirer la projection sans supprimer les sources.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect group/recurrence | Analyst | Grouping Relation | 0 | read | membres et raisons visibles | non |
| compute grouping candidate | deterministic service | matches | 1 | keys/time scope | relation candidate | non |
| acknowledge local grouped state | Operator | local candidate | 2 | workflow state permis | état local seulement | non |
| request suppression change | authorized user | suppression source | 2 request | owner connu | aucune mutation implicite | selon owner |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| détecter repeats | oui | oui | oui | suggestion possible | exact/grouping keys |
| calculer group candidate | oui | oui | oui | suggestion possible | règles déterministes |
| expliquer groupe | oui | oui | oui | oui | member/reason list |
| supprimer provenance | non | non | non | interdit | all source refs retained |

## 14. États fonctionnels
`ungrouped`, `group-candidate`, `grouped`, `duplicate-candidate`, `repeated`, `suppressed-projection`, `suppression-expired`, `partial`, `unknown`.

## 15. États d’interface
Aucun Screen ID. Suppressed doit rester inspectable selon permission ; group count ne remplace pas la visibilité des membres ; stale suppression est explicitement marquée.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| grouping relation | concept Endpoint | Analyst/CAP-EPT-044 | tous membres référencés |
| duplicate/recurrence assessment | état Endpoint | Detection Engineering | aucune suppression destructive |
| suppression-state projection | état Endpoint | consumers autorisés | owner/reason/expiry/provenance |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| matches/candidates | recurrence | grouping assessment | ids/keys/time | sources intactes |
| suppression source | policy/state change | projection refresh | owner/scope/reason/expiry | history retained |
| group | investigation pivot | CAP-EPT-042/044 | member refs/context | grouped ≠ identical |

## 18. Dépendances
CAP-EPT-032..034 ; `detection-suppression.md` ; `local-correlation.md` ; Investigate suppression/exception proposals ; Settings policy/configuration ; tenant isolation ; `OPEN-008`, `OPEN-017`.

## 19. Source de vérité
Endpoint est source du group local et de sa projection technique. Il n’est pas source de vérité de la policy administrative de suppression ni du canonical Command Signal.

## 20. Provenance et audit
Conserver member ids, grouping keys/version, window/time skew, actor/service, suppression owner/reason/scope/expiry, previous/new state, tenant et correlation id. Aucun member n’est effacé.

## 21. Permissions fonctionnelles
Lire matches/context, groupes, rationale et suppression state ; interaction de suppression uniquement via owner autorisé ; deny cross-tenant ; pas de RBAC final.

## 22. Limites et erreurs
Partial telemetry/time skew peut rendre grouping inconclusif ; grouped ≠ identical ; deduplicated ≠ erased ; suppressed ≠ deleted ; multiple matches ≠ Incident ; no suppression source = aucune suppression inférée.

## 23. Métriques
Repeat/group candidate counts, recurrence depth, suppression projection coverage/expiry, grouping uncertainty, provenance completeness.

## 24. Classification de livraison
`draft / defined / planned` ; pas de moteur/policy runtime déployé.

## 25. Critères d’acceptation
**Given** un match équivalent se répète, **When** le grouping est calculé, **Then** un group/recurrence est créé en conservant chaque match original.

**Given** une suppression est active, **When** le signal candidate est consulté, **Then** il reste visible avec reason/owner/expiry et n’est pas supprimé.

**Given** deux matches ont une clé proche mais un contexte différent, **When** ils sont groupés comme candidate, **Then** `grouped` ne les déclare pas identiques.

## 26. Questions ouvertes
`OPEN-017` conserve runtime/portability ; `OPEN-008` conserve plateformes/sources. Le default owner/policy exact d’une suppression dépend de la source existante et n’est pas réinventé ici.

## 27. Consommateurs documentaires
EPT-3 signal/context/timeline/handoff ; Investigate Detection Engineering ; Settings ; Command ; Security ; Quality ; Roadmap.