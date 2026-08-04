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

**Situation :** les coordinateurs doivent comprendre la séquence et l’origine des changements malgré des horloges, sources et produits différents. **Utilisateurs :** Incident Commander, SOC Analyst L2, Auditor, Business Owner. **Sans la capacité :** une action peut être interprétée avant sa cause ou un résultat comme une décision.

## 3. Objectifs

- séparer événements observés, inférés, actions, Decisions et Results ;
- afficher event time, ingestion time, source et incertitude ;
- permettre filtrage et retour à l’objet source ;
- fournir une alternative tabulaire accessible.

## 4. Non-objectifs

- posséder le Timeline Engine ;
- réécrire les événements source ;
- corriger silencieusement le clock skew ;
- produire une conclusion analytique.

## 5. Propriétaire

Command / Mission Control / Command Product Lead. Command possède le contenu métier de la projection ; Shared Capabilities possède le moteur temporel.

## 6. Utilisateurs

Rôle principal : Incident Commander. Rôles secondaires : SOC Analyst L2, Auditor, Business Owner. Les permissions des objets sources sont toujours réévaluées.

## 7. Conditions d’entrée

Contexte tenant, période définie et au moins une source autorisée.

## 8. Entrées fonctionnelles

| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Timeline entries | Timeline Engine | événements normalisés | oui | event/ingestion time | afficher les sources disponibles et la lacune |
| Object links | Object Linking Service | références stables | oui | résolution courante | conserver un tombstone autorisé |
| Clock quality | sources | métadonnée de temps | non | par événement | marquer l’ordre incertain |

## 9. Objets lus

| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Timeline Entry | Shared | type, timestamps, source, objet | projection |
| Incident | Command | événements de coordination | lecture |
| Decision / Response Run / Result | Govern | événements gouvernés | projection |

## 10. Objets créés ou modifiés

Aucun objet. Les interactions métier sont de classe 0 ; un export crée seulement un job via Shared Export Engine.

## 11. Fonctionnalités

- filtrer par période, type, source, objet et acteur ;
- distinguer les catégories d’événements ;
- inspecter provenance et correlation ID ;
- ouvrir l’objet source puis restaurer la position ;
- exporter une projection autorisée.

## 12. Actions utilisateur

| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Filtrer la timeline | lecteur autorisé | timeline | 0 | source chargée | projection réduite | non |
| Inspecter la provenance | lecteur autorisé | Timeline Entry | 0 | permission source | Trace ouverte | non |
| Exporter la sélection | rôle export autorisé | projection | 0 | classification/redaction vérifiées | job d’export | non |

## 13. Automatisation et IA

| Fonction | Humain | Règle | Moteur déterministe | Workflow | Agent | Govern | Alternative sans IA |
|---|---|---|---|---|---|---|---|
| Ordonner et filtrer | oui | possible | oui, source principale | possible | résumé facultatif | non | moteur temporel et filtres |
| Regrouper une rafale | correction possible | règle versionnée | regroupement explicable | possible | proposition uniquement | non | règles et liste brute |

Aucune fonction essentielle ne dépend d’un modèle.

## 14. États fonctionnels

`current`, `delayed-source`, `incomplete`, `clock-conflict`, `source-unavailable`, `empty-period`.

## 15. États d’interface

Loading conserve période et filtres ; Empty indique une période vide ; Partial nomme les sources manquantes ; Error conserve les entrées valides ; Offline utilise la dernière projection garantie ; Permission denied masque l’entrée protégée ; Stale montre la source et le retard. Le Design System possède le rendu.

## 16. Sorties

| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Situation chronology | projection ordonnée | Mission Control, handover, reporting | chaque entrée conserve type et provenance |
| Selected event context | return origin | objet source | filtre, période et position préservés |

## 17. Transitions

| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Timeline | sélection d’un objet | produit propriétaire | tenant, objet, timestamp, source | retour à l’entrée sélectionnée |

## 18. Dépendances

Timeline Engine, Object Linking Service, Trace, Export Engine et CAP-CMD-006. Aucune dépendance ne transfère l’ownership.

## 19. Source de vérité

Les événements restent propriétaires de leur source. Timeline Engine normalise l’ordre ; Command choisit les catégories et filtres métier. Clock quality, source, fraîcheur et permissions restent visibles.

## 20. Provenance et audit

Chaque entrée expose source, producer type, event time, ingestion time, correlation ID et objet lié. Les exports enregistrent acteur, scope et classification.

## 21. Permissions fonctionnelles

`perm.command.read`, permissions de chaque objet source et permission d’export distincte. Les namespaces finaux sont reportés.

## 22. Limites et erreurs

Horloge incertaine, source absente, lien inaccessible, données stale, changement de tenant ou permission refusée ne doivent pas être masqués. Une navigation échouée conserve la timeline, les filtres et la position.

## 23. Métriques

- part des événements avec source et timestamps complets ;
- nombre de conflits d’horloge visibles ;
- taux de retours restaurant la position.

## 24. Classification de livraison

`delivery_status: defined`, `delivery_mode: planned`, cible native. Aucun moteur ou logiciel livré n’est prouvé.

## 25. Critères d’acceptation

**Given** plusieurs événements Command et Govern, **When** l’utilisateur filtre une période, **Then** l’ordre, la catégorie, la source et la fraîcheur de chaque entrée restent visibles.

**Given** un clock conflict, **When** les entrées sont affichées, **Then** l’incertitude est signalée et aucun ordre certain n’est inventé.

**Given** aucun modèle IA, **When** la timeline est utilisée, **Then** le moteur déterministe, les filtres et l’alternative tabulaire fournissent le résultat complet.

## 26. Questions ouvertes

- Quelle politique de regroupement évite une timeline illisible lors des rafales ? — REQ-PROD-005, REQ-PROD-008, REQ-PROD-013, REQ-UX-007.
- Quels événements Command sont obligatoires dans un handover ? — mêmes Requirement IDs.

## 27. Consommateurs documentaires

Mission Control Situation, Incident Detail, Handover, reports opérationnels, parcours Phase 5, écrans Phase 6, objets Phase 7 et contrats ultérieurs.
