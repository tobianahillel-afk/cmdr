---
id: CAP-CMD-004
title: Handover
product: command
module: mission-control
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-008
  - REQ-PROD-013
  - REQ-PROD-021
  - REQ-PROD-010
open_decisions:
  - none
source-of-truth: canonical
---

# CAP-CMD-004 — Handover

## 1. Définition

Prépare, transmet, accuse réception et supersède une relève opérationnelle structurée reliant situation, Incidents, Tasks, risques, blocages, Decisions, Runs, owners et prochaines actions.

## 2. Problème utilisateur

**Situation.** Lors d’un changement d’équipe ou de rôle, le contexte et la responsabilité se perdent dans des messages libres et des listes non synchronisées.

**Utilisateurs concernés.** Incident Commander en premier lieu ; SOC Team Lead, SOC Analyst L2, Business Owner.

**Conséquence sans la capacité.** Un handover incomplet provoque duplication, délais, oubli de Decision en attente et actions sans owner.

## 3. Objectifs

- permettre un handover complet sans IA
- lier chaque section aux objets sources
- rendre visible l’owner sortant, le destinataire et l’accusé de réception
- conserver les corrections et supersessions

## 4. Non-objectifs

- transférer automatiquement l’ownership canonique des objets
- remplacer les objets par un document figé
- envoyer des Evidence brutes non autorisées
- faire dépendre l’envoi d’un fournisseur de modèle

## 5. Propriétaire

- **Produit :** Command.
- **Module :** Mission Control.
- **Rôle responsable :** Command Product Lead.
- **Raison :** Handover contribue à « Maintenir une conscience partagée de la situation, des priorités, des blocages, des handovers et des résultats sans devenir une seconde Work Queue. ». Command possède uniquement les mutations listées pour cette capability ; les projections externes gardent leur owner.

## 6. Utilisateurs

- **Rôle principal :** Incident Commander — responsable du résultat opérationnel principal.
- **Rôles secondaires :** SOC Team Lead, SOC Analyst L2, Business Owner.
- Pour Handover, le rôle principal répond du résultat décrit en section 3 ; les rôles secondaires consultent ou contribuent uniquement dans leur tenant, environnement, scope et permissions.

## 7. Conditions d’entrée

équipe source et destinataire identifiés; contexte tenant; objets autorisés disponibles ou lacunes déclarées.

## 8. Entrées fonctionnelles

| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Situation selection | CAP-CMD-001 | snapshot et objets liés | oui | au moment du brouillon | autoriser la saisie manuelle avec lacunes visibles |
| Owners and next actions | Incident/Task | coordination courante | oui | version courante | bloquer l’état ready si les éléments critiques n’ont pas d’owner |
| Pending governance | Govern projections | Decisions/Runs en attente | non | état courant | indiquer Govern indisponible et conserver les références |

## 9. Objets lus

| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Incident / Task | Command | situation, owner, blocage, next action | lecture |
| Decision / Response Run / Result | Govern | statut résumé | projection |
| Case / Finding | Investigate | références et résumé autorisé | projection |

## 10. Objets créés ou modifiés

| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Handover record | créer, mettre à jour, envoyer, reconnaître, superséder | Command coordination | classe 2 et audit |
| Incident / Task ownership | transfert séparé après acceptation | Command | jamais implicite dans l’envoi |

Dans Handover, les objets externes listés en section 9 sont exclusivement des projections : aucune relation, suggestion ou transition ne transfère leur ownership à Command.

## 11. Fonctionnalités

- composer manuellement toutes les sections obligatoires
- agréger déterministement les objets sélectionnés
- générer un brouillon IA attribué et modifiable
- valider complétude, fraîcheur et owners
- envoyer, accepter, rejeter pour correction et superséder

## 12. Actions utilisateur

| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Créer ou modifier le brouillon | Incident Commander | handover | 2 | contexte actif | draft versionné | non |
| Marquer prêt | Incident Commander | handover | 2 | sections critiques complètes | ready | non |
| Envoyer | Incident Commander | handover | 2 | destinataire et périmètre confirmés | sent | non |
| Accepter | destinataire autorisé | handover | 2 | handover sent et courant | acknowledged | non |
| Rejeter pour correction | destinataire autorisé | handover | 2 | motif fourni | rejected-for-correction | non |

Les classes suivent le modèle 0–4. Une demande dont l’effet cible est de classe 3 ou 4 reste une préparation ou une transition de classe 2 dans Command ; l’autorité et l’exécution appartiennent à Govern.

## 13. Automatisation et IA

| Fonction | Humain | Règle | Moteur déterministe | Workflow | Agent | Govern | Alternative sans IA |
|---|---|---|---|---|---|---|---|
| Créer ou modifier le brouillon | oui; action explicite | possible si policy versionnée | validation, déduplication et contrôle de version | possible avec étapes visibles | proposition uniquement, jamais autorité | non par défaut; OPEN-013 peut s’appliquer à une mutation de classe 2 | action manuelle complète |
| composer manuellement toutes les sections obligatoires | oui; consultation et correction | oui pour sélection/alerte explicable | oui pour agrégation et calcul sourcé | possible | résumé ou proposition attribuée | non pour l’observation | données sources, filtres et règles |
| envoyer, accepter, rejeter pour correction et superséder | oui; décision finale humaine | possible | source de repli obligatoire | possible | facultatif; run, sources et incertitude visibles | non par défaut; OPEN-013 peut s’appliquer à une mutation de classe 2 | agréger déterministement les objets sélectionnés |

Handover reste entièrement utilisable sans fournisseur de modèle : « Créer ou modifier le brouillon » et les sorties essentielles disposent d’une voie humaine ou déterministe.

## 14. États fonctionnels

Les états ci-dessous décrivent exclusivement le travail de Handover ; ils ne valident ni ne remplacent les machines d’état futures de Incident / Task, Decision / Response Run / Result, Case / Finding :

- `draft`
- `ready`
- `sent`
- `acknowledged`
- `rejected-for-correction`
- `superseded`

## 15. États d’interface

- **Loading :** conserver l’anatomie du module et indiquer quelles entrées de Handover sont en cours de résolution.
- **Empty :** expliquer si aucun objet ne correspond, si le scope est vide ou si la capability n’est pas configurée.
- **Partial :** nommer les sources manquantes et leur effet sur Handover.
- **Error :** conserver les données valides et fournir une reprise sûre avec correlation ID lorsque disponible.
- **Offline :** rester en lecture sur la dernière donnée garantie ; bloquer toute mutation dont la version ne peut pas être validée.
- **Permission denied :** ne révéler ni objet ni valeur protégée ; indiquer la famille de permission et le chemin de demande.
- **Stale :** afficher source, dernière mise à jour et conséquence fonctionnelle avant toute action.

Le rendu détaillé de ces états reste dans le Design System et les futurs écrans.

## 16. Sorties

| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Handover package | record Command | équipe destinataire | liens sources, fraîcheur et owner visibles |
| Acknowledgement | event | équipe source et audit | acteur, date et périmètre confirmés |

## 17. Transitions

| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Handover | ouverture d’un objet lié | produit propriétaire | tenant, objet, handover origin | retour au handover |
| Handover acknowledged | transfert d’ownership demandé | Work Assignment | objets sélectionnés et destinataire | résultats élément par élément |

## 18. Dépendances

- CAP-CMD-001
- CAP-CMD-005
- CAP-CMD-006
- Notification Center
- Collaboration Service
- Object Linking Service

Pour Handover, la dépendance primaire est « CAP-CMD-001 ». Permission Model, audit/provenance, tenant isolation et Object Linking Service restent transversaux lorsque des références sont utilisées.

## 19. Source de vérité

- **Données propriétaires :** pour Handover, seules les opérations listées en section 10 sur Handover record, Incident / Task ownership sont autorisées.
- **Projections :** Decision / Response Run / Result, Case / Finding restent résolues par référence stable et permission-aware.
- **Sources externes :** chaque entrée issue de CAP-CMD-001, Incident/Task, Govern projections conserve source, version ou timestamp, fraîcheur et classification.
- **Données dérivées :** les calculs propres à Handover exposent leurs facteurs, leur version et leurs limites ; ils ne remplacent pas les objets sources.

## 20. Provenance et audit

Pour Handover, toute mutation ou proposition enregistre acteur humain ou service, producteur (`human`, `rule`, `engine`, `workflow`, `agent`, `external`), source, version/run, objet, avant/après, justification, résultat, tenant, environnement et correlation ID. La conservation détaillée sera possédée par les contrats d’audit ultérieurs ; la Phase 4A fixe seulement la sémantique fonctionnelle.

## 21. Permissions fonctionnelles

- perm.command.read
- perm.command.coordinate
- permissions objet par référence
- accès inter-tenant interdit

Les familles utilisées par Handover incluent perm.command.read, perm.command.coordinate, permissions objet par référence. Leur atomisation, les namespaces manquants, le step-up et les règles ABAC finales restent la responsabilité de Security/Permissions.

## 22. Limites et erreurs

- données requises absentes ou non autorisées; impact spécifique : Handover ne doit pas produire un résultat présenté comme complet.
- projection stale ou dépendance indisponible; impact spécifique : Handover ne doit pas produire un résultat présenté comme complet.
- conflit de version lors d'une mutation; impact spécifique : Handover ne doit pas produire un résultat présenté comme complet.
- tenant ou environnement devenu incompatible; impact spécifique : Handover ne doit pas produire un résultat présenté comme complet.
- autorisation refusée sans divulgation de données; impact spécifique : Handover ne doit pas produire un résultat présenté comme complet.

Une transition de Handover vers produit propriétaire qui échoue conserve le workspace source, le contexte sûr et le travail non enregistré récupérable.

## 23. Métriques

- part des handovers avec owner et prochaine action sur chaque élément critique
- délai jusqu’à acknowledgement
- taux de handovers rejetés ou superseded

Ces métriques sont conceptuelles ; aucune cible chiffrée définitive n’est fixée en Phase 4A.

## 24. Classification de livraison

- **Delivery status fonctionnel :** `defined`.
- **Delivery mode courant :** `planned` — aucun logiciel ou runtime n’est prouvé par ce document.
- **Cible produit :** native lorsque la capability fait partie du cœur Command ; deployment-dependent pour Customers and Delivery.
- **Preuve actuelle :** le document CAP-CMD-004 définit ownership, entrées/sorties, actions, états, dépendances et critères de Handover ; aucune implémentation.
- **Conditions de promotion :** objets suffisants, permissions approuvées, parcours et écrans finalisés, contrats techniques, implémentation et validation.

## 25. Critères d’acceptation

### Scénario A — comportement principal

**Given** un Incident Commander autorisé, le tenant et l’environnement corrects, des entrées valides et une version courante,  
**When** il exécute « Créer ou modifier le brouillon »,  
**Then** le résultat de Handover est observable, la classe d’action et l’owner sont explicites, la mutation éventuelle est auditée et aucun objet d’un autre produit n’est redéfini.

### Scénario B — fonctionnement sans IA

**Given** aucun fournisseur de modèle et aucun Automation Agent disponible,  
**When** l’utilisateur utilise Handover,  
**Then** les données sources, règles, moteurs déterministes et actions manuelles permettent le résultat essentiel ; aucune suggestion n’est requise.

### Scénario C — dépendance partielle ou permission refusée

**Given** le cas limite « données requises absentes ou non autorisées »,  
**When** Handover est ouverte ou qu’une mutation est tentée,  
**Then** les données encore valides sont conservées, l’impact précis sur les sorties de la section 16 est visible, aucune donnée protégée n’est révélée et le workspace source reste récupérable.

## 26. Questions ouvertes

- Le handover devient-il un objet canonique Command en Phase Objets ou un record de coordination versionné ? — Requirement IDs : REQ-PROD-008, REQ-PROD-013, REQ-PROD-021, REQ-PROD-010.
- Quels champs sont bloquants pour l’état ready selon le tenant ? — Requirement IDs : REQ-PROD-008, REQ-PROD-013, REQ-PROD-021, REQ-PROD-010.

Aucune nouvelle décision ouverte n’est créée par cette capability.

## 27. Consommateurs documentaires

- Mission Control — Handover
- parcours de relève
- Incident Detail
- audit opérationnel
- futurs parcours de Phase 5
- futurs écrans et leur front matter de Phase 6
- objets et relations de Phase 7
- familles de permissions et contrats techniques ultérieurs
