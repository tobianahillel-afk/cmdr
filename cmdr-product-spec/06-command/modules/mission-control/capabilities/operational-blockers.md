---
id: CAP-CMD-005
title: Operational Blockers
product: command
module: mission-control
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-006
  - REQ-PROD-009
  - REQ-PROD-013
  - REQ-PROD-021
open_decisions:
  - OPEN-013
source-of-truth: canonical
---

# CAP-CMD-005 — Operational Blockers

## 1. Définition

Identifie et suit un empêchement opérationnel en le représentant comme état ou relation d’un Incident/Task, ou comme Task de suivi, sans créer un nouvel objet canonique autonome.

## 2. Problème utilisateur

**Situation.** Les équipes voient qu’un travail n’avance plus, mais ne savent pas toujours ce qui est bloqué, par qui, jusqu’à quand et quelle action peut le débloquer.

**Utilisateurs concernés.** Incident Commander en premier lieu ; SOC Analyst L2, Task owner, Team Lead.

**Conséquence sans la capacité.** Le blocage devient une note informelle, disparaît des priorités et n’est pas pris en compte dans le handover ou l’escalade.

## 3. Objectifs

- nommer le type, la cause, l’objet bloqué et l’owner
- lier une prochaine action et une échéance
- permettre escalade et résolution auditées
- réutiliser Incident/Task plutôt que créer un objet concurrent

## 4. Non-objectifs

- créer un objet Blocker définitif
- modifier l’objet source d’un autre produit
- résoudre automatiquement une dépendance
- assimiler tout état pending à un blocage

## 5. Propriétaire

- **Produit :** Command.
- **Module :** Mission Control.
- **Rôle responsable :** Command Product Lead.
- **Raison :** Operational Blockers contribue à « Maintenir une conscience partagée de la situation, des priorités, des blocages, des handovers et des résultats sans devenir une seconde Work Queue. ». Command possède uniquement les mutations listées pour cette capability ; les projections externes gardent leur owner.

## 6. Utilisateurs

- **Rôle principal :** Incident Commander — responsable du résultat opérationnel principal.
- **Rôles secondaires :** SOC Analyst L2, Task owner, Team Lead.
- Pour Operational Blockers, le rôle principal répond du résultat décrit en section 3 ; les rôles secondaires consultent ou contribuent uniquement dans leur tenant, environnement, scope et permissions.

## 7. Conditions d’entrée

Incident ou Task accessible; raison de blocage connue ou déclarée inconnue; owner ou équipe de reprise.

## 8. Entrées fonctionnelles

| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Blocked work | Incident ou Task | référence de travail | oui | version courante | aucune création de blocage sans objet source |
| Blocker reason | utilisateur ou projection externe | cause structurée | oui | au moment de la déclaration | utiliser unknown avec prochaine action de qualification |
| Dependency | objet ou acteur lié | relation | non | fraîcheur déclarée | conserver le blocage sans inventer la dépendance |

## 9. Objets lus

| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Incident / Task | Command | état, owner, dépendances, prochaine action | lecture/modification |
| Decision / Case / external dependency | produit propriétaire | statut pertinent | projection |

## 10. Objets créés ou modifiés

| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Incident / Task | marquer bloqué, raison, owner et next action | Command | classe 2 |
| Task | créer une action de déblocage si nécessaire | Command | éviter un objet Blocker distinct |

Dans Operational Blockers, les objets externes listés en section 9 sont exclusivement des projections : aucune relation, suggestion ou transition ne transfère leur ownership à Command.

## 11. Fonctionnalités

- déclarer et catégoriser un blocage
- lier une dépendance ou une personne attendue
- affecter la résolution et une échéance
- escalader sans changer l’owner avant acceptation
- résoudre avec motif et preuve de reprise

## 12. Actions utilisateur

| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Déclarer bloqué | owner/coordinateur | Incident/Task | 2 | raison et prochaine action | état de travail bloqué | OPEN-013 |
| Créer une Task de déblocage | coordinateur | Task | 2 | aucune Task équivalente active | Task liée | OPEN-013 |
| Escalader | coordinateur | Incident/Task | 2 | destinataire et délai | escalade tracée | selon destination |
| Résoudre le blocage | owner | Incident/Task | 2 | cause traitée | état repris et historique | OPEN-013 |

Les classes suivent le modèle 0–4. Une demande dont l’effet cible est de classe 3 ou 4 reste une préparation ou une transition de classe 2 dans Command ; l’autorité et l’exécution appartiennent à Govern.

## 13. Automatisation et IA

| Fonction | Humain | Règle | Moteur déterministe | Workflow | Agent | Govern | Alternative sans IA |
|---|---|---|---|---|---|---|---|
| Déclarer bloqué | oui; action explicite | possible si policy versionnée | validation, déduplication et contrôle de version | possible avec étapes visibles | proposition uniquement, jamais autorité | non par défaut; OPEN-013 peut s’appliquer à une mutation de classe 2 | action manuelle complète |
| déclarer et catégoriser un blocage | oui; consultation et correction | oui pour sélection/alerte explicable | oui pour agrégation et calcul sourcé | possible | résumé ou proposition attribuée | non pour l’observation | données sources, filtres et règles |
| résoudre avec motif et preuve de reprise | oui; décision finale humaine | possible | source de repli obligatoire | possible | facultatif; run, sources et incertitude visibles | non par défaut; OPEN-013 peut s’appliquer à une mutation de classe 2 | lier une dépendance ou une personne attendue |

Operational Blockers reste entièrement utilisable sans fournisseur de modèle : « Déclarer bloqué » et les sorties essentielles disposent d’une voie humaine ou déterministe.

## 14. États fonctionnels

Les états ci-dessous décrivent exclusivement le travail de Operational Blockers ; ils ne valident ni ne remplacent les machines d’état futures de Incident / Task, Decision / Case / external dependency, Task :

- `identified`
- `owned`
- `escalated`
- `mitigation-in-progress`
- `resolved`
- `unknown-cause`

## 15. États d’interface

- **Loading :** conserver l’anatomie du module et indiquer quelles entrées de Operational Blockers sont en cours de résolution.
- **Empty :** expliquer si aucun objet ne correspond, si le scope est vide ou si la capability n’est pas configurée.
- **Partial :** nommer les sources manquantes et leur effet sur Operational Blockers.
- **Error :** conserver les données valides et fournir une reprise sûre avec correlation ID lorsque disponible.
- **Offline :** rester en lecture sur la dernière donnée garantie ; bloquer toute mutation dont la version ne peut pas être validée.
- **Permission denied :** ne révéler ni objet ni valeur protégée ; indiquer la famille de permission et le chemin de demande.
- **Stale :** afficher source, dernière mise à jour et conséquence fonctionnelle avant toute action.

Le rendu détaillé de ces états reste dans le Design System et les futurs écrans.

## 16. Sorties

| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Blocker relation/state | Incident/Task update | Work Queue, Mission Control, handover | cause, owner, next action et fraîcheur |
| Improvement task | Task | Readiness and Operations | cause et résultat liés |

## 17. Transitions

| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Operational Blocker | besoin d’expertise | Investigate | Incident, cause, délai, objets liés | retour à l’objet bloqué |
| Operational Blocker | besoin d’autorité | Govern | contexte, impact, urgence, action proposée | owner Command conservé jusqu’à acceptation |

## 18. Dépendances

- CAP-CMD-110
- CAP-CMD-107
- Object Linking Service
- Notification Center

Pour Operational Blockers, la dépendance primaire est « CAP-CMD-110 ». Permission Model, audit/provenance, tenant isolation et Object Linking Service restent transversaux lorsque des références sont utilisées.

## 19. Source de vérité

- **Données propriétaires :** pour Operational Blockers, seules les opérations listées en section 10 sur Incident / Task, Task sont autorisées.
- **Projections :** Decision / Case / external dependency restent résolues par référence stable et permission-aware.
- **Sources externes :** chaque entrée issue de Incident ou Task, utilisateur ou projection externe, objet ou acteur lié conserve source, version ou timestamp, fraîcheur et classification.
- **Données dérivées :** les calculs propres à Operational Blockers exposent leurs facteurs, leur version et leurs limites ; ils ne remplacent pas les objets sources.

## 20. Provenance et audit

Pour Operational Blockers, toute mutation ou proposition enregistre acteur humain ou service, producteur (`human`, `rule`, `engine`, `workflow`, `agent`, `external`), source, version/run, objet, avant/après, justification, résultat, tenant, environnement et correlation ID. La conservation détaillée sera possédée par les contrats d’audit ultérieurs ; la Phase 4A fixe seulement la sémantique fonctionnelle.

## 21. Permissions fonctionnelles

- perm.command.coordinate
- perm.command.incident.manage
- perm.command.task.manage

Les familles utilisées par Operational Blockers incluent perm.command.coordinate, perm.command.incident.manage, perm.command.task.manage. Leur atomisation, les namespaces manquants, le step-up et les règles ABAC finales restent la responsabilité de Security/Permissions.

## 22. Limites et erreurs

- données requises absentes ou non autorisées; impact spécifique : Operational Blockers ne doit pas produire un résultat présenté comme complet.
- projection stale ou dépendance indisponible; impact spécifique : Operational Blockers ne doit pas produire un résultat présenté comme complet.
- conflit de version lors d'une mutation; impact spécifique : Operational Blockers ne doit pas produire un résultat présenté comme complet.
- tenant ou environnement devenu incompatible; impact spécifique : Operational Blockers ne doit pas produire un résultat présenté comme complet.
- autorisation refusée sans divulgation de données; impact spécifique : Operational Blockers ne doit pas produire un résultat présenté comme complet.

Une transition de Operational Blockers vers Investigate qui échoue conserve le workspace source, le contexte sûr et le travail non enregistré récupérable.

## 23. Métriques

- nombre de travaux bloqués sans owner ou prochaine action
- âge des blocages par cause
- taux de résolution avec reprise vérifiée

Ces métriques sont conceptuelles ; aucune cible chiffrée définitive n’est fixée en Phase 4A.

## 24. Classification de livraison

- **Delivery status fonctionnel :** `defined`.
- **Delivery mode courant :** `planned` — aucun logiciel ou runtime n’est prouvé par ce document.
- **Cible produit :** native lorsque la capability fait partie du cœur Command ; deployment-dependent pour Customers and Delivery.
- **Preuve actuelle :** le document CAP-CMD-005 définit ownership, entrées/sorties, actions, états, dépendances et critères de Operational Blockers ; aucune implémentation.
- **Conditions de promotion :** objets suffisants, permissions approuvées, parcours et écrans finalisés, contrats techniques, implémentation et validation.

## 25. Critères d’acceptation

### Scénario A — comportement principal

**Given** un Incident Commander autorisé, le tenant et l’environnement corrects, des entrées valides et une version courante,  
**When** il exécute « Déclarer bloqué »,  
**Then** le résultat de Operational Blockers est observable, la classe d’action et l’owner sont explicites, la mutation éventuelle est auditée et aucun objet d’un autre produit n’est redéfini.

### Scénario B — fonctionnement sans IA

**Given** aucun fournisseur de modèle et aucun Automation Agent disponible,  
**When** l’utilisateur utilise Operational Blockers,  
**Then** les données sources, règles, moteurs déterministes et actions manuelles permettent le résultat essentiel ; aucune suggestion n’est requise.

### Scénario C — dépendance partielle ou permission refusée

**Given** le cas limite « données requises absentes ou non autorisées »,  
**When** Operational Blockers est ouverte ou qu’une mutation est tentée,  
**Then** les données encore valides sont conservées, l’impact précis sur les sorties de la section 16 est visible, aucune donnée protégée n’est révélée et le workspace source reste récupérable.

## 26. Questions ouvertes

- Quels attributs minimaux du blocage doivent être portés par Incident et Task en Phase Objets ? — Requirement IDs : REQ-PROD-006, REQ-PROD-009, REQ-PROD-013, REQ-PROD-021.
- Quand une dépendance externe devient-elle une Task de suivi ? — Requirement IDs : REQ-PROD-006, REQ-PROD-009, REQ-PROD-013, REQ-PROD-021.

Décisions structurées liées : `OPEN-013`.

## 27. Consommateurs documentaires

- Mission Control
- Unified Work Queue
- Handover
- Readiness improvement actions
- futurs parcours de Phase 5
- futurs écrans et leur front matter de Phase 6
- objets et relations de Phase 7
- familles de permissions et contrats techniques ultérieurs
