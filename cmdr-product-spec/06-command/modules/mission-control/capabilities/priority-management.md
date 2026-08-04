---
id: CAP-CMD-002
title: Priority Management
product: command
module: mission-control
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-003
  - REQ-PROD-010
  - REQ-PROD-013
  - REQ-PROD-021
open_decisions:
  - OPEN-013
source-of-truth: canonical
---

# CAP-CMD-002 — Priority Management

## 1. Définition

Établit et explique la priorité opérationnelle effective d’un Incident ou d’une Task en séparant impact, urgence, SLA, criticité de service, dépendances et recommandations.

## 2. Problème utilisateur

**Situation.** Les équipes doivent arbitrer rapidement entre plusieurs travaux sans confondre sévérité technique, urgence, impact et priorité opérationnelle.

**Utilisateurs concernés.** Incident Commander en premier lieu ; SOC Analyst L2, Business Owner, Team Lead.

**Conséquence sans la capacité.** Sans facteurs visibles et historique, la priorité devient arbitraire, difficile à contester et vulnérable aux suggestions automatisées opaques.

## 3. Objectifs

- montrer la priorité effective et ses facteurs
- permettre une modification autorisée avec justification
- distinguer décision humaine, règle déterministe et proposition IA
- conserver l’historique et les valeurs superseded

## 4. Non-objectifs

- définir un algorithme universel de score
- modifier la severity de la source
- appliquer silencieusement une recommandation
- créer une Decision Govern pour chaque priorité

## 5. Propriétaire

- **Produit :** Command.
- **Module :** Mission Control.
- **Rôle responsable :** Command Product Lead.
- **Raison :** Priority Management contribue à maintenir une conscience partagée de la situation et des priorités. Command possède uniquement les mutations listées ; les projections externes gardent leur owner.

## 6. Utilisateurs

- **Rôle principal :** Incident Commander.
- **Rôles secondaires :** SOC Analyst L2, Business Owner, Team Lead.
- Les rôles secondaires consultent ou contribuent uniquement dans leur tenant, environnement, scope et permissions.

## 7. Conditions d’entrée

Incident ou Task accessible ; facteurs disponibles ou explicitement incomplets ; rôle autorisé à coordonner.

## 8. Entrées fonctionnelles

| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Priority factors | Incident, Task, SLA, Service | facteurs structurés | oui | valeurs courantes | afficher les facteurs manquants |
| Current priority | Command | valeur effective | oui | version courante | refuser la mutation |
| Recommendation | règle, moteur, workflow ou agent | proposition attribuée | non | avec run/version | aucun effet sur la priorité effective |

## 9. Objets lus

| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Incident / Task | Command | priority, impact, urgency, owner, state | lecture et modification autorisée |
| Service | Shared | criticité et dépendances | projection |
| SLA context | Command/policy source | risque et échéance | projection explicable |

## 10. Objets créés ou modifiés

| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Incident / Task | modifier priority et justification | Command | classe 2, concurrence optimiste et audit |
| Audit trail | ajouter l’historique | source audit partagée | append-only conceptuel |

Les objets externes restent des projections ; aucune relation ou suggestion ne transfère leur ownership à Command.

## 11. Fonctionnalités

- comparer les facteurs sans les réduire à un score unique ;
- afficher l’origine de la priorité courante ;
- prévisualiser l’effet d’une nouvelle priorité sur la file ;
- accepter ou rejeter une proposition attribuée ;
- restaurer une valeur antérieure par une nouvelle mutation autorisée.

## 12. Actions utilisateur

| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter les facteurs | rôles lecteurs | Incident/Task | 0 | objet accessible | explication visible | non |
| Modifier la priorité | coordinateur autorisé | Incident/Task | 2 | justification et version courante | nouvelle priorité effective | OPEN-013 |
| Accepter une recommandation | coordinateur autorisé | proposition | 2 | facteurs visibles | mutation humaine auditée | OPEN-013 |
| Rejeter une recommandation | coordinateur autorisé | proposition | 2 | proposition active | rejet et motif audités | non |

Une demande dont l’effet cible est de classe 3 ou 4 reste une préparation ou transition de classe 2 dans Command ; l’autorité et l’exécution appartiennent à Govern.

## 13. Automatisation et IA

| Fonction | Humain | Règle | Moteur déterministe | Workflow | Agent | Govern | Alternative sans IA |
|---|---|---|---|---|---|---|---|
| Modifier la priorité | décision explicite | possible si policy versionnée | validation et calcul sourcé | possible | proposition uniquement | OPEN-013 | mutation manuelle complète |
| Comparer les facteurs | consultation/correction | sélection explicable | agrégation et calcul | possible | résumé attribué | non | données sources et règles |
| Restaurer la valeur précédente | décision humaine | possible | contrôle de version | possible | facultatif | OPEN-013 | historique et action manuelle |

Priority Management reste utilisable sans fournisseur de modèle.

## 14. États fonctionnels

- `effective`
- `proposed`
- `disputed`
- `awaiting-owner`
- `factors-incomplete`
- `superseded`

Ces états décrivent le travail de la capability, pas la machine d’état finale des objets.

## 15. États d’interface

Loading conserve le contexte ; Empty explique l’absence ; Partial nomme les facteurs manquants ; Error garde les données valides ; Offline bloque les mutations non vérifiables ; Permission denied ne révèle rien ; Stale affiche source, date et conséquence. Le rendu reste possédé par le Design System.

## 16. Sorties

| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Priority change | Incident/Task event | Work Queue et Mission Control | acteur, facteurs, avant/après et justification |
| Recommendation disposition | audit event | moteur/workflow et analystes | acceptation ou rejet distinct de la proposition |

## 17. Transitions

| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Priority Management | besoin d’investigation des facteurs | Investigate | Incident et facteurs incomplets | retour à l’Incident |
| Priority Management | impact nécessitant autorité | Govern | Incident, impact, urgence et action proposée | aucune Decision créée localement |

## 18. Dépendances

CAP-CMD-104, CAP-CMD-105, CAP-CMD-204, CAP-CMD-205, Metrics Engine et audit hooks. `OPEN-013` reste ouverte.

## 19. Source de vérité

La priorité effective et sa justification sont Command ; Service et SLA sont des projections sourcées ; recommandations et calculs exposent producteur, version, facteurs et limites.

## 20. Provenance et audit

Toute mutation ou proposition enregistre acteur/producteur, source, version/run, objet, avant/après, justification, résultat, tenant, environnement et correlation ID.

## 21. Permissions fonctionnelles

`perm.command.coordinate`, `perm.command.incident.manage`, `perm.command.task.manage` et ABAC ownership/tenant. L’atomisation, le step-up et les namespaces finaux sont reportés.

## 22. Limites et erreurs

Données absentes, projection stale, dépendance indisponible, conflit de version, changement de tenant/environnement ou refus d’autorisation ne doivent jamais produire une priorité présentée comme complète. Une transition échouée conserve le workspace source.

## 23. Métriques

- part des changements avec justification complète ;
- délai de résolution des priorités disputées ;
- taux de recommandations acceptées, rejetées ou ignorées.

Aucune cible chiffrée définitive n’est fixée.

## 24. Classification de livraison

- **Delivery status :** `defined`.
- **Delivery mode courant :** `planned`.
- **Cible :** native.
- **Preuve actuelle :** specification fonctionnelle uniquement, sans implémentation.
- **Promotion :** objets, permissions, parcours, écrans, contrats, implémentation et validation.

## 25. Critères d’acceptation

**Given** un utilisateur autorisé et des facteurs courants, **When** il modifie la priorité, **Then** la nouvelle valeur, les facteurs, l’acteur et la justification sont auditables.

**Given** une recommandation automatisée différente, **When** elle est consultée, **Then** la priorité effective reste inchangée jusqu’à acceptation humaine autorisée.

**Given** aucun modèle IA, **When** la capability est utilisée, **Then** règles, moteurs déterministes et actions manuelles permettent le résultat essentiel.

## 26. Questions ouvertes

- Quels facteurs sont obligatoires avant une modification ? — REQ-PROD-003, REQ-PROD-010, REQ-PROD-013, REQ-PROD-021.
- Quelles mutations de classe 2 exigent Govern ? — `OPEN-013`.

## 27. Consommateurs documentaires

Mission Control Priorities, Unified Work Queue, Incident Detail, parcours Phase 5, écrans Phase 6, objets/relations Phase 7 et permissions/contrats ultérieurs.
