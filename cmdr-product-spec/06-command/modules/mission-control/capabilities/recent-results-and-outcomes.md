---
id: CAP-CMD-006
title: Recent Results and Outcomes
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
  - REQ-PROD-021
open_decisions:
  - OPEN-013
source-of-truth: canonical
---

# CAP-CMD-006 — Recent Results and Outcomes

## 1. Définition

Consomme les Results Govern récents et les relie aux Response Runs, Decisions et Incidents afin d’actualiser la situation, le risque résiduel et la prochaine action sans modifier le Result.

## 2. Problème utilisateur

**Situation :** une exécution peut réussir techniquement sans que sa vérification ou son effet métier soit compris dans Command. **Utilisateurs :** Incident Commander, Response Operator, SOC Analyst L2, Business Owner. **Sans la capacité :** l’Incident peut être fermé trop tôt ou le même travail relancé.

## 3. Objectifs

- distinguer outcome technique, vérification et impact métier ;
- lier Result, Run, Decision et Incident ;
- réinjecter le résultat dans la situation ;
- permettre une contestation sans réécrire le Result.

## 4. Non-objectifs

- créer ou valider un Result ;
- déterminer la réussite du Response Run ;
- modifier une Decision ;
- masquer un résultat disputé.

## 5. Propriétaire

Command / Mission Control / Command Product Lead pour la consommation opérationnelle. Govern reste propriétaire de Decision, Response Run et Result.

## 6. Utilisateurs

Rôle principal : Incident Commander. Rôles secondaires : Response Operator, SOC Analyst L2, Business Owner.

## 7. Conditions d’entrée

Result autorisé disponible ou projection partielle ; lien Run/Decision ; Incident cible identifiable ou association à confirmer.

## 8. Entrées fonctionnelles

| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Result projection | Govern | outcome, verification, residual risk | oui | état courant | afficher pending verification ou association manuelle |
| Incident context | Command | statut, impact, prochaine action | oui | version courante | ne pas muter si inaccessible |
| Business validation | Business Owner/Incident Commander | confirmation opérationnelle | non | datée | distinguer technique et métier |

## 9. Objets lus

| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Result / Response Run / Decision | Govern | résultat, vérification, portée, conditions | projection |
| Incident | Command | état, impact, next action | lecture/modification |

## 10. Objets créés ou modifiés

| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Incident | lier Result, ajuster next action ou état de travail | Command | classe 2, jamais automatique par défaut |
| Timeline | ajouter un événement de consommation | Shared/Command content | référence vers le Result source |

## 11. Fonctionnalités

- présenter Results récents et fraîcheur ;
- différencier succeeded, verified et business-accepted ;
- associer un Result à l’Incident pertinent ;
- marquer contestation ou vérification manquante ;
- créer une prochaine action d’amélioration.

## 12. Actions utilisateur

| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter le Result | rôle autorisé | Result | 0 | permission Govern | projection visible | non |
| Associer à un Incident | coordinateur | Incident | 2 | objets compatibles | relation auditée | OPEN-013 |
| Mettre à jour la prochaine action | Incident Commander | Incident | 2 | effet compris | Incident mis à jour | OPEN-013 |
| Demander vérification | Incident Commander | Result/Task | 2 | résultat absent ou contesté | Task ou retour Govern | selon destination |

## 13. Automatisation et IA

| Fonction | Humain | Règle | Moteur déterministe | Workflow | Agent | Govern | Alternative sans IA |
|---|---|---|---|---|---|---|---|
| Classer les résultats | correction humaine | possible | agrégation sourcée | possible | résumé attribué | propriétaire source | filtres et règles |
| Mettre à jour Incident | décision humaine | possible si policy | validation/version | possible | proposition uniquement | OPEN-013 | mutation manuelle |

Aucune suggestion ne modifie Result ou Incident silencieusement.

## 14. États fonctionnels

`received`, `pending-verification`, `verified`, `disputed`, `applied-to-situation`, `superseded`.

## 15. États d’interface

Loading conserve les liens ; Empty explique l’absence ; Partial montre Result sans certaines projections ; Error conserve les données valides ; Offline bloque la mutation ; Permission denied ne révèle pas le Result ; Stale affiche source et date. Le rendu reste au Design System.

## 16. Sorties

| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Outcome projection | linked Result | Mission Control, Incident, reporting | source Govern et niveau de vérification visibles |
| Operational follow-up | Incident update ou Task | Work Queue / Readiness | ne modifie pas le Result |

## 17. Transitions

| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Govern Result | publication autorisée | Command Incident | Result, Run, Decision, residual risk | retour Govern possible |
| Command | contestation | Govern / Investigate | Result, motif, Incident, Evidence requise | retour à l’Incident |

## 18. Dépendances

CAP-CMD-001, CAP-CMD-003, CAP-CMD-106, Object Linking Service, Timeline Engine et Reporting Engine.

## 19. Source de vérité

Govern possède le Result ; Command possède seulement les relations et mutations Incident explicites. Les projections conservent source, version, vérification, fraîcheur et permissions.

## 20. Provenance et audit

Toute association, contestation, mutation ou proposition enregistre acteur/producteur, source, avant/après, justification, tenant, environnement et correlation ID.

## 21. Permissions fonctionnelles

`perm.govern.result.read`, `perm.command.incident.manage`, `perm.command.task.manage`. L’atomisation et `OPEN-013` restent reportées.

## 22. Limites et erreurs

Result absent, stale, inaccessible, non associé, conflit d’Incident ou dépendance indisponible ne doivent pas être présentés comme un outcome complet. Une transition échouée conserve le workspace source.

## 23. Métriques

- part des Results associés à un Incident ;
- part distinguant vérification technique et métier ;
- délai Result → prochaine action Command.

## 24. Classification de livraison

`delivery_status: defined`, `delivery_mode: planned`, cible native. La preuve est documentaire uniquement.

## 25. Critères d’acceptation

**Given** un Result vérifié lié à un Incident, **When** Command le consomme, **Then** le Run, la Decision, le niveau de vérification et le risque résiduel restent visibles.

**Given** un Result disputé, **When** une vérification est demandée, **Then** Command crée seulement la Task ou le contexte de retour et ne modifie pas le Result.

**Given** aucun modèle IA, **When** la capability est utilisée, **Then** filtres, règles, moteurs déterministes et actions manuelles fournissent le résultat essentiel.

## 26. Questions ouvertes

- Quelle validation métier est nécessaire avant `applied-to-situation` ? — REQ-PROD-005, REQ-PROD-008, REQ-PROD-013, REQ-PROD-021.
- Quand une contestation repart-elle vers Govern ou Investigate ? — mêmes Requirement IDs.

Décision liée : `OPEN-013` pour les mutations de classe 2.

## 27. Consommateurs documentaires

Mission Control, Incident Detail, Readiness improvement actions, reports, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions/contrats ultérieurs.
