---
id: CAP-CMD-003
title: Situation Timeline
product: command
module: mission-control
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-005
  - REQ-PROD-008
  - REQ-PROD-013
  - REQ-UX-007
open_decisions:
  - none
source-of-truth: canonical
---

# CAP-CMD-003 — Situation Timeline

## 1. Définition
Projette dans un ordre temporel explicite les changements Command et les références autorisées vers Incidents, Decisions, Response Runs et Results, sans dupliquer le Timeline Engine.

## 2. Problème utilisateur
Les coordinateurs doivent comprendre une séquence issue de sources et d’horloges différentes. Sans cette capacité, une action peut être interprétée avant sa cause ou un Result comme une Decision.

## 3. Objectifs
- séparer événements observés, inférés, actions, Decisions et Results ;
- afficher event time, ingestion time, source et incertitude ;
- permettre filtrage et retour exact à l’objet source ;
- fournir une alternative tabulaire accessible.

## 4. Non-objectifs
Ne possède pas le Timeline Engine, ne réécrit pas les événements sources, ne corrige pas silencieusement le clock skew et ne produit pas de Finding.

## 5. Propriétaire
Command / Mission Control / Command Product Lead pour le contenu métier ; Shared Capabilities possède le moteur temporel.

## 6. Utilisateurs
Principal : Incident Commander. Secondaires : SOC Analyst L2, Auditor, Business Owner.

## 7. Conditions d’entrée
Tenant, période et au moins une source autorisée ; permissions de chaque objet réévaluées.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Timeline entries | Timeline Engine | événements normalisés | oui | event time et ingestion time | afficher les sources disponibles et nommer la lacune |
| Object links | Object Linking Service | références stables | oui | résolution à l’ouverture | conserver un tombstone autorisé sans données protégées |
| Clock quality | producteurs sources | métadonnée temporelle | non | par événement | marquer l’ordre comme incertain |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Timeline Entry | Shared Capabilities | type, timestamps, source et objet | consulter, filtrer et inspecter |
| Incident | Command | événements de coordination | consulter et naviguer |
| Decision / Response Run / Result | Govern | événements gouvernés | consulter et relier en lecture seule |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Objets sources | aucune mutation | produits propriétaires | projection en lecture seule ; la timeline ne réécrit aucun événement |
| Export job | créer sur demande | Shared Export Engine | classe 0, scope et permissions de la sélection conservés |

## 11. Fonctionnalités
Filtrer par période/type/source/objet/acteur, distinguer les catégories, inspecter provenance et correlation ID, ouvrir l’objet source, restaurer la position et exporter une projection autorisée.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Filtrer | lecteur autorisé | timeline | 0 | source chargée | projection réduite | non |
| Inspecter provenance | lecteur autorisé | Timeline Entry | 0 | permission source | Trace ouverte | non |
| Exporter la sélection | rôle export | projection | 0 | classification et redaction vérifiées | job d’export | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Ordonner les événements | correction humaine possible | oui | oui | non nécessaire | Timeline Engine et ordre explicite |
| Filtrer et regrouper | oui | oui | oui | résumé facultatif | filtres, règles versionnées et liste brute |
| Résumer une période | oui | agrégation possible | oui | brouillon attribué | lecture chronologique et export tabulaire |

## 14. États fonctionnels
`current`, `delayed-source`, `incomplete`, `clock-conflict`, `source-unavailable`, `empty-period`.

## 15. États d’interface
Loading conserve période/filtres ; Empty indique une période vide ; Partial nomme les sources ; Error conserve les entrées valides ; Offline utilise la dernière projection garantie ; Permission denied masque l’entrée ; Stale expose le retard.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Situation chronology | projection ordonnée | Mission Control, Handover et Reporting | type, source, timestamps et incertitude conservés |
| Selected event context | événement de navigation | produit propriétaire | filtre, période, position et return origin préservés |
| Timeline export | job d’export | utilisateur autorisé | permission-aware, redacted et non destructif |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Situation Timeline | sélection d’un Incident | Incident Detail | tenant, Incident, timestamp et filtres | entrée et position restaurées |
| Situation Timeline | sélection d’une Decision ou d’un Run | Govern | référence stable, Incident et return origin | même période restaurée |
| Situation Timeline | export utilisateur | Shared Export Engine | sélection, classification et colonnes | notification puis timeline inchangée |

## 18. Dépendances
Timeline Engine, Object Linking Service, Trace, Export Engine et CAP-CMD-006.

## 19. Source de vérité
Chaque événement reste propriétaire de sa source ; Timeline Engine normalise l’ordre ; Command définit catégories et filtres métier.

## 20. Provenance et audit
Source, producer type, event time, ingestion time, correlation ID, objet lié et scope d’export sont conservés.

## 21. Permissions fonctionnelles
`perm.command.read`, permissions des objets sources et permission d’export distincte ; atomisation reportée.

## 22. Limites et erreurs
Clock skew, source absente, lien inaccessible, données stale, tenant incompatible ou refus ne sont jamais masqués ; une navigation échouée conserve filtres et position.

## 23. Métriques
Part des événements avec timestamps complets, conflits d’horloge visibles et taux de retour exact ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; aucun moteur ni logiciel livré n’est prouvé.

## 25. Critères d’acceptation
**Given** des événements Command et Govern, **When** une période est filtrée, **Then** ordre, catégorie, source et fraîcheur restent visibles.

**Given** un clock conflict, **When** les entrées sont affichées, **Then** l’incertitude est signalée sans ordre certain inventé.

**Given** aucun modèle IA, **When** la timeline est utilisée, **Then** moteur déterministe, filtres et alternative tabulaire fournissent le résultat complet.

## 26. Questions ouvertes
Quelle politique de regroupement et quels événements sont obligatoires dans un Handover ? — Requirement IDs ci-dessus. Aucun nouvel OPEN.

## 27. Consommateurs documentaires
Mission Control Situation, Incident Detail, Handover, reports, parcours Phase 5, écrans Phase 6, objets Phase 7 et contrats ultérieurs.