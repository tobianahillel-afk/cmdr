---
id: CAP-CMD-110
title: Escalation
product: command
module: incidents-and-work-queue
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-003
  - REQ-PROD-004
  - REQ-PROD-008
  - REQ-PROD-013
  - REQ-SEC-001
open_decisions:
  - OPEN-013
source-of-truth: canonical
---
# CAP-CMD-110 — Escalation

## 1. Définition
Formalise une escalade fonctionnelle, hiérarchique, technique, vers Investigate ou vers Govern avec justification, destinataire, délai, contexte transmis et acceptation explicite.

## 2. Problème utilisateur
Une escalade informelle perd l’owner courant et confond demande d’aide, investigation et autorité. Principal : Incident Commander ; secondaires : SOC Analyst L2, Team Lead, Business Owner.

## 3. Objectifs
Qualifier type/raison ; maintenir owner jusqu’à acceptation ; transmettre contexte minimum/return origin ; distinguer Investigate et Action Request Govern.

## 4. Non-objectifs
Ne pas créer Decision automatiquement, transférer ownership canonique, exécuter containment ou considérer une notification comme acceptation.

## 5. Propriétaire
Command possède le record/état de coordination ; Investigate/Govern possèdent les objets destination.

## 6. Utilisateurs
Principal : Incident Commander. Secondaires : SOC Analyst L2, Team Lead, Business Owner.

## 7. Conditions d’entrée
Incident/Task accessible ; motif, destinataire, délai, contexte autorisé et type d’escalade.

## 8. Entrées fonctionnelles
| Entrée | Source | Type | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Source work | Incident/Task | objet/next action | oui | version courante | refuser |
| Escalation type | utilisateur/rule | functional/hierarchical/technical/investigate/govern | oui | envoi | sélection explicite |
| Context package | Command+projections | impact, urgence, blockers, refs | oui | visible | partial déclaré si autorisé |

## 9. Objets lus
| Objet | Owner | Projection | Droit local |
|---|---|---|---|
| Incident / Task | Command | owner, reason, deadline, context | lecture/modification |
| Case / Finding / Decision | produit propriétaire | liens/statut | projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Owner | Règle |
|---|---|---|---|
| Escalation record/state | créer, envoyer, accepter, décliner, retourner | Command coordination | classe 2 |
| Case or Action Request | demander création/liaison | Investigate/Govern | ownership destination |

## 11. Fonctionnalités
Créer escalade typée ; notifier sans perdre owner ; accepter/décliner/retourner ; ouvrir/lier Case ; préparer Action Request.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Escalade fonctionnelle/hiérarchique | coordinateur | Incident/Task | 2 | destinataire | sent | OPEN-013 |
| Vers Investigate | analyste autorisé | Case request | 2 | Incident/contexte | Case créé/lié par Investigate | non |
| Vers Govern | requester | Action Request | 2 | action/impact/targets | draft/submitted | oui |
| Accepter/décliner | destinataire | escalation | 2 | sent | accepted/declined | non |

## 13. Automatisation et IA
| Fonction | Humain | Règle | Moteur | Workflow | Agent | Govern | Sans IA |
|---|---|---|---|---|---|---|---|
| Préparer package | correction | complétude | agrégation | possible | brouillon | selon destination | formulaire manuel |
| Envoyer/accepter | décision humaine | possible | validation/version | possible | jamais autorité | OPEN-013 | actions manuelles |

## 14. États fonctionnels
`drafted`, `sent`, `accepted`, `declined`, `returned-for-information`, `expired`, `cancelled`.

## 15. États d’interface
Partial montre contexte manquant ; Error garde draft ; Offline bloque envoi ; Permission denied ne révèle pas destination ; Stale impose validation. Rendu DS.

## 16. Sorties
| Sortie | Objet/événement | Consommateur | Garantie |
|---|---|---|---|
| Escalation | record coordination | source/destination | type, raison, owner, délai, acknowledgement |
| Case/Action Request link | relation | Command/destination | ownership destination préservé |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte | Retour |
|---|---|---|---|---|
| Command | technical escalation | Investigate | tenant, env, Incident, time range, entities, blockers | Incident Detail |
| Command | authority/action | Govern | Incident, refs, impact, urgence, target, action, alternatives | Incident Detail |
| Destination | accept/decline | Command | result/ref | owner change séparé |

## 18. Dépendances
CAP-CMD-005, CAP-CMD-106, Notification Center, Collaboration Service, Object Linking Service, Context preservation.

## 19. Source de vérité
Escalation record : Command ; Case : Investigate ; Action Request : Govern. Les liens ne transfèrent aucun droit.

## 20. Provenance et audit
Type, sender, destination, package, freshness, disposition, owner, timestamps et correlation ID sont enregistrés.

## 21. Permissions fonctionnelles
`perm.command.coordinate`, `perm.investigate.case.create`, `perm.govern.action-request.create`, permissions destination. Atomisation reportée.

## 22. Limites et erreurs
Destination inaccessible, contexte incomplet, expiration, refus, conflit, changement tenant/env ou service indisponible conservent draft/source et owner courant.

## 23. Métriques
Délai envoi→acceptation ; escalades retournées ; part conservant owner/return origin. Aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire seulement.

## 25. Critères d’acceptation
**Given** un Incident sans Case, **When** escalade technique est envoyée, **Then** Investigate reçoit contexte, possède le Case et le retour restaure Incident Detail.

**Given** une action classe 3, **When** l’escalade Govern est envoyée, **Then** Action Request est préparée, aucune exécution locale ne se produit et l’origine Command reste tracée.

**Given** aucun modèle IA, **When** une escalade est préparée, **Then** toutes les sections sont saisissables manuellement.

## 26. Questions ouvertes
Quels types exigent acceptation avant changement d’owner ? Quand créer une Task de suivi ? — REQ-PROD-003, REQ-PROD-004, REQ-PROD-008, REQ-PROD-013, REQ-SEC-001. `OPEN-013` reste ouverte.

## 27. Consommateurs documentaires
Incident Detail, Work Queue, Investigate/Govern transitions, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions ultérieures.
