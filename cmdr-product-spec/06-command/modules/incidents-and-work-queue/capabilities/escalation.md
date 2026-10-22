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
Formalise une escalade fonctionnelle, hiérarchique, technique, vers Investigate ou Govern avec justification, destinataire, délai, contexte transmis et acceptation explicite.

## 2. Problème utilisateur
Une escalade informelle perd l’owner courant et confond demande d’aide, investigation et autorité. Sans record, notification et acceptation sont prises pour un transfert.

## 3. Objectifs
Qualifier type/raison, maintenir owner jusqu’à acceptation, transmettre contexte minimum/return origin et distinguer Investigate d’Action Request Govern.

## 4. Non-objectifs
Ne crée pas Decision automatiquement, ne transfère pas ownership canonique, n’exécute pas containment et ne considère pas une notification comme acceptation.

## 5. Propriétaire
Command possède le record/état de coordination ; Investigate/Govern possèdent les objets destination.

## 6. Utilisateurs
Principal : Incident Commander. Secondaires : SOC Analyst L2, Team Lead, Business Owner.

## 7. Conditions d’entrée
Incident/Task accessible, motif, destinataire, délai, contexte autorisé et type d’escalade.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Source work | Incident ou Task | objet et prochaine action | oui | version courante | escalade refusée |
| Escalation type | utilisateur ou règle | fonctionnelle, hiérarchique, technique, Investigate ou Govern | oui | validé à l’envoi | sélection explicite requise |
| Context package | Command et projections | impact, urgence, blockers et refs | oui | fraîcheur visible | `partial` déclaré seulement si autorisé |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Incident / Task | Command | owner, raison, délai et contexte | consulter et modifier |
| Case / Finding | Investigate | lien et statut | consulter et relier en projection |
| Decision | Govern | statut et relation | consulter en lecture seule |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Escalation record | créer, envoyer, accepter, décliner ou retourner | Command coordination | classe 2, versionné et audité |
| Case | demander création ou liaison | Investigate | ownership destination préservé |
| Action Request | demander création ou soumettre | Govern | Command fournit le package, aucune Decision locale |

## 11. Fonctionnalités
Créer une escalade typée, notifier sans perdre owner, accepter/décliner/retourner, ouvrir/lier Case et préparer Action Request.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Escalade interne | coordinateur | Incident/Task | 2 | destinataire et délai | sent | OPEN-013 |
| Vers Investigate | analyste autorisé | Case request | 2 | Incident et contexte | Case créé/lié par Investigate | non |
| Vers Govern | requester | Action Request | 2 | action, impact et targets | draft/submitted | oui |
| Accepter ou décliner | destinataire | Escalation | 2 | sent et version courante | accepted/declined | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Préparer le package | oui | agrégation/complétude | oui | brouillon attribué | formulaire et références manuelles |
| Choisir le type/destinataire | oui | règles possibles | oui | suggestion seulement | catalogue de types et sélection humaine |
| Envoyer l’escalade | oui | validation/version | workflow possible | jamais autonome | action humaine explicite |
| Accepter ou décliner | oui | contrôle de version | notification possible | non décisionnelle | action humaine avec motif |

## 14. États fonctionnels
`drafted`, `sent`, `accepted`, `declined`, `returned-for-information`, `expired`, `cancelled`.

## 15. États d’interface
Partial montre le contexte manquant ; Error garde le draft ; Offline bloque envoi ; Permission denied ne révèle pas la destination ; Stale impose validation.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Escalation | record Command | source et destination | type, raison, owner, délai et acknowledgement visibles |
| Case link | relation | Command et Investigate | owner Investigate et permissions conservés |
| Action Request context | objet ou événement Govern | Govern | origine Command, impact et alternatives attribués |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Command | escalade technique | Investigate | tenant, environnement, Incident, période, entities et blockers | Incident Detail restauré |
| Command | besoin d’autorité/action | Govern | Incident, refs, impact, urgence, target, action et alternatives | Incident Detail restauré |
| Destination | acceptation ou refus | Command | résultat, référence et motif | owner change séparé |

## 18. Dépendances
CAP-CMD-005, CAP-CMD-106, Notification Center, Collaboration Service, Object Linking Service et Context preservation.

## 19. Source de vérité
Escalation record : Command. Case : Investigate. Action Request : Govern. Les liens ne transfèrent aucun droit.

## 20. Provenance et audit
Type, sender, destination, package, fraîcheur, disposition, owner, timestamps et correlation ID.

## 21. Permissions fonctionnelles
`perm.command.coordinate`, `perm.investigate.case.create`, `perm.govern.action-request.create` et permissions destination ; atomisation reportée.

## 22. Limites et erreurs
Destination inaccessible, contexte incomplet, expiration, refus, conflit, tenant/env changé ou service indisponible conservent draft/source et owner courant.

## 23. Métriques
Délai envoi→acceptation, escalades retournées et part conservant owner/return origin ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire seulement.

## 25. Critères d’acceptation
**Given** un Incident sans Case, **When** une escalade technique est envoyée, **Then** Investigate reçoit le contexte, possède le Case et le retour restaure Incident Detail.

**Given** une action classe 3, **When** l’escalade Govern est envoyée, **Then** Action Request est préparée et aucune exécution locale n’a lieu.

**Given** aucun modèle IA, **When** une escalade est préparée, **Then** toutes les sections restent saisissables manuellement.

## 26. Questions ouvertes
Quels types exigent acceptation avant changement d’owner et quand créer une Task de suivi ? — Requirement IDs ci-dessus ; `OPEN-013` reste ouverte.

## 27. Consommateurs documentaires
Incident Detail, Work Queue, transitions Investigate/Govern, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions ultérieures.