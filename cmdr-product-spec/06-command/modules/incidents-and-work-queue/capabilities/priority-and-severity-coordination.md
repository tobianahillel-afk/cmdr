---
id: CAP-CMD-104
title: Priority and Severity Coordination
product: command
module: incidents-and-work-queue
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-005
  - REQ-PROD-013
  - REQ-PROD-021
  - REQ-UX-005
open_decisions:
  - OPEN-013
source-of-truth: canonical
---

# CAP-CMD-104 — Priority and Severity Coordination

## 1. Définition
Présente severity, priority, urgency, impact, confidence et SLA risk comme dimensions distinctes, avec sources et droits propres.

## 2. Problème utilisateur
Une severity élevée peut être confondue avec une priorité automatique. Sans distinction, la file est mal ordonnée et les facteurs ne sont pas contestables.

## 3. Objectifs
Afficher chaque dimension/source, permettre mutation de priority/impact seulement, conserver severity/confidence comme projections et expliquer les divergences.

## 4. Non-objectifs
Ne crée pas de score universel, ne modifie pas une severity externe, ne fusionne pas les unités et n’utilise pas la couleur seule.

## 5. Propriétaire
Command coordonne priority/impact sur Incident/Task ; les sources propriétaires gardent severity/confidence.

## 6. Utilisateurs
Principal : SOC Analyst L2. Secondaires : Incident Commander, Business Owner, Detection Engineer en lecture.

## 7. Conditions d’entrée
Work item accessible, facteurs autorisés et définition/version disponible ou lacune déclarée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Severity | Alert, Signal ou Detection source | dimension technique | non | timestamp source | `unknown`, sans priorité déduite |
| Priority | Command | dimension opérationnelle | oui | version courante | triage demandé |
| Impact / urgency / confidence / SLA risk | propriétaires respectifs | facteurs séparés | non | par facteur | facteur absent explicitement |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Incident / Task | Command | priority, urgency et impact | consulter et modifier si autorisé |
| Signal / Alert | Command ou source | severity et confidence | consulter et contester via retour source |
| Service / SLA | Shared ou policy source | criticité et risque | consulter et comparer |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Incident / Task | modifier priority, urgency ou impact autorisé | Command | classe 2 et historique |
| Recommendation disposition | créer acceptation/rejet | Command audit | aucune mutation implicite |
| Severity / confidence source | aucune mutation | propriétaire source | projection en lecture seule |

## 11. Fonctionnalités
Présenter six dimensions, comparer sans fusion, afficher source/calcul/fraîcheur/owner, expliquer divergence et diriger vers Priority Management.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter dimensions | lecteur | work item | 0 | données autorisées | facteurs visibles | non |
| Modifier priority | coordinateur | Incident/Task | 2 | justification | priority effective | OPEN-013 |
| Confirmer impact | IC ou Business Owner | Incident | 2 | source/certitude | impact mis à jour | OPEN-013 |
| Contester projection | analyste | severity/confidence | 2 | motif | retour owner source | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Résoudre les dimensions | oui | oui | oui | résumé facultatif | lecture des sources et règles |
| Calculer un contexte | oui | facteurs explicables | oui | proposition attribuée | moteur déterministe ou comparaison humaine |
| Modifier priority/impact | oui | validation/version | workflow possible | jamais silencieusement | mutation humaine complète |
| Contester une projection | oui | routage source | oui | brouillon de motif | action manuelle vers owner source |

## 14. États fonctionnels
`dimensions-complete`, `dimensions-partial`, `priority-proposed`, `priority-effective`, `source-conflict`, `unknown`.

## 15. États d’interface
Partial nomme la dimension absente ; Error conserve les facteurs valides ; Offline bloque mutation ; Permission denied masque seulement le facteur protégé ; Stale montre source/date.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Coordination factors | projection | Queue et Mission Control | dimensions non confondues et sourcées |
| Priority/impact event | mise à jour | audit et consommateurs | source, before/after et justification |
| Projection dispute | événement | propriétaire source | motif attribué et objet référencé |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Coordination | severity/confidence contestée | source ou Investigate | objet, facteur et raison | work item restauré |
| Coordination | impact exige autorité | Govern | Incident, facteurs et action | aucune Decision locale |
| Coordination | modification priority | Priority Management | objet, facteurs et proposition | file ou détail restauré |

## 18. Dépendances
CAP-CMD-002, CAP-CMD-105, CAP-CMD-204, CAP-CMD-205 et Metrics Engine.

## 19. Source de vérité
Priority/impact autorisés : Command. Severity/confidence/service/SLA : propriétaires sources. Calculs : définition/version/facteurs visibles.

## 20. Provenance et audit
Chaque dimension expose source/owner/fraîcheur ; mutations et dispositions enregistrent acteur et before/after.

## 21. Permissions fonctionnelles
`perm.command.read`, `perm.command.coordinate`, `perm.command.incident.manage` et permissions sources ; atomisation reportée.

## 22. Limites et erreurs
Facteur absent/stale, définition incompatible, refus ou conflit ne produisent jamais un score total opaque.

## 23. Métriques
Items avec dimensions sourcées, divergences expliquées et impacts confirmed/assumed ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire seulement.

## 25. Critères d’acceptation
**Given** severity élevée et impact faible, **When** l’item est inspecté, **Then** les six dimensions sont séparées et sourcées.

**Given** une proposition automatisée, **When** elle est consultée, **Then** la priorité effective reste inchangée avant acceptation.

**Given** aucun modèle IA, **When** la capability est utilisée, **Then** sources, règles et actions manuelles suffisent.

## 26. Questions ouvertes
Quelles dimensions sont obligatoires et qui valide confidence entre produits ? — Requirement IDs ci-dessus ; `OPEN-013` reste ouverte.

## 27. Consommateurs documentaires
Work Queue, Mission Control, Incident Detail, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions ultérieures.