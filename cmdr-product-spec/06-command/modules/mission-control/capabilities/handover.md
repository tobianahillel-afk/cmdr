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
Lors d’un changement d’équipe, contexte et responsabilité se perdent dans des messages libres. Sans handover structuré, Decisions en attente et actions sans owner sont oubliées.

## 3. Objectifs
Permettre un handover complet sans IA, lier chaque section aux sources, rendre visibles expéditeur/destinataire/acknowledgement et conserver corrections/supersessions.

## 4. Non-objectifs
Ne transfère pas implicitement l’ownership, ne remplace pas les objets sources, n’expose pas d’Evidence non autorisée et ne dépend pas d’un modèle.

## 5. Propriétaire
Command / Mission Control / Command Product Lead pour le record de coordination ; les objets liés gardent leur owner.

## 6. Utilisateurs
Principal : Incident Commander. Secondaires : SOC Team Lead, SOC Analyst L2, Business Owner.

## 7. Conditions d’entrée
Équipe source et destinataire identifiés, tenant courant et objets autorisés ou lacunes déclarées.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Situation selection | CAP-CMD-001 | snapshot et objets liés | oui | au moment du brouillon | saisie manuelle permise avec lacune visible |
| Owners and next actions | Incident et Task | coordination courante | oui | version courante | état `ready` bloqué pour les éléments critiques |
| Pending governance | Govern | Decisions et Runs en attente | non | état courant | Govern indisponible signalé, références conservées |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Incident / Task | Command | situation, owner, blocage et prochaine action | consulter et relier |
| Decision / Response Run / Result | Govern | statut résumé et conditions | consulter en projection |
| Case / Finding | Investigate | référence et résumé autorisé | consulter et naviguer |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Handover record | créer, mettre à jour, envoyer, reconnaître et superséder | Command coordination | classe 2, versionné et audité |
| Incident / Task ownership | transfert séparé après acceptation | Command | jamais implicite dans l’envoi du handover |
| Objets externes | aucune mutation | Investigate ou Govern | projections en lecture seule |

## 11. Fonctionnalités
Composer toutes les sections manuellement, agréger les objets sélectionnés, vérifier complétude/fraîcheur/owners, générer un brouillon attribué et gérer envoi/acceptation/correction/supersession.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Modifier le brouillon | Incident Commander | Handover | 2 | contexte actif | draft versionné | non |
| Marquer prêt | Incident Commander | Handover | 2 | sections critiques complètes | ready | non |
| Envoyer | Incident Commander | Handover | 2 | destinataire et périmètre confirmés | sent | non |
| Accepter ou rejeter | destinataire | Handover | 2 | version courante | acknowledged ou correction demandée | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Composer le handover | oui | agrégation possible | oui | brouillon attribué | formulaire et liens sources manuels |
| Vérifier la complétude | oui | oui | oui | explication facultative | règles de complétude versionnées |
| Envoyer et reconnaître | oui | validation/version | workflow possible | non décisionnelle | actions humaines explicites |
| Préparer le résumé | oui | agrégation sourcée | oui | oui, modifiable | synthèse manuelle des sections |

## 14. États fonctionnels
`draft`, `ready`, `sent`, `acknowledged`, `rejected-for-correction`, `superseded`.

## 15. États d’interface
Partial nomme les sections/sources manquantes ; Offline bloque envoi et acknowledgement ; Permission denied masque l’objet protégé ; Stale affiche source/date et peut bloquer `ready`.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Handover package | record Command | équipe destinataire | versionné, attribué, sourcé et permission-aware |
| Acknowledgement | événement de coordination | équipe source et audit | acteur, date et périmètre confirmés |
| Correction request | événement | Incident Commander | motif, version et sections concernées conservés |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Handover | ouverture d’un objet lié | produit propriétaire | tenant, objet, handover origin | handover et section restaurés |
| Handover acknowledged | transfert d’ownership demandé | Work Assignment | items sélectionnés et destinataire | résultats élément par élément |
| Handover | Decision en attente sélectionnée | Govern | Incident, Decision et return origin | même handover restauré |

## 18. Dépendances
CAP-CMD-001, CAP-CMD-005, CAP-CMD-006, Notification Center, Collaboration Service et Object Linking Service.

## 19. Source de vérité
Handover record et coordination : Command. Objets Investigate/Govern : projections stables et permission-aware.

## 20. Provenance et audit
Auteur, producteur de brouillon, sources, version, before/after, envoi, acknowledgement, tenant et correlation ID.

## 21. Permissions fonctionnelles
`perm.command.read`, `perm.command.coordinate`, permissions objet par référence et interdiction inter-tenant ; atomisation reportée.

## 22. Limites et erreurs
Objets absents, stale ou interdits, conflit de version, destinataire inaccessible et changement de tenant empêchent un handover présenté comme complet.

## 23. Métriques
Handovers avec owner/prochaine action, délai d’acknowledgement et taux de correction/supersession ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire uniquement.

## 25. Critères d’acceptation
**Given** aucun fournisseur de modèle, **When** un Incident Commander prépare un handover, **Then** toutes les sections restent saisissables et envoyables manuellement.

**Given** des items critiques sans owner, **When** le handover est marqué prêt, **Then** la complétude échoue avec les lacunes listées.

**Given** un handover envoyé, **When** le destinataire l’accepte, **Then** acknowledgement et éventuel transfert d’ownership restent deux opérations distinctes.

## 26. Questions ouvertes
Le handover devient-il un objet canonique et quels champs bloquent `ready` par tenant ? — Requirement IDs ci-dessus. Aucun nouvel OPEN.

## 27. Consommateurs documentaires
Mission Control Handover, parcours de relève, Incident Detail, audit, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions ultérieures.