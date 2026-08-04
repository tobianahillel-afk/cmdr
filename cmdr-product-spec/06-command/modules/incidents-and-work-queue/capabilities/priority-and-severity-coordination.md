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
Une severity élevée peut être confondue avec une priorité automatique. Rôles : SOC Analyst L2, Incident Commander, Business Owner, Detection Engineer en consultation. Sans distinction, la file est mal ordonnée.

## 3. Objectifs
Afficher chaque dimension/source ; permettre mutation de priority/impact seulement ; conserver severity/confidence comme projections ; expliquer divergences.

## 4. Non-objectifs
Ne pas créer un score universel, modifier une severity externe, fusionner les unités ou utiliser la couleur seule.

## 5. Propriétaire
Command coordonne priority/impact sur Incident/Task ; les sources propriétaires gardent severity/confidence.

## 6. Utilisateurs
Principal : SOC Analyst L2. Secondaires : Incident Commander, Business Owner, Detection Engineer en lecture.

## 7. Conditions d’entrée
Work item accessible ; facteurs autorisés ; définition/version disponible ou lacune déclarée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Severity | Alert/Signal/Detection source | technique | non | source courante | unknown |
| Priority | Command | opérationnelle | oui | version courante | triage demandé |
| Impact/Urgency/Confidence/SLA risk | sources propriétaires | facteurs | non | par facteur | aucun calcul silencieux |

## 9. Objets lus
| Objet | Owner | Projection | Droit local |
|---|---|---|---|
| Incident / Task | Command | priority, urgency, impact | lecture/modification |
| Signal / Alert | source/Command | severity, confidence | projection |
| Service / SLA | Shared/policy | criticité, risk | projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Owner | Règle |
|---|---|---|---|
| Incident / Task | priority, urgency ou impact autorisé | Command | classe 2, historique |
| Recommendation disposition | audit | Shared audit | aucune mutation implicite |

## 11. Fonctionnalités
Présenter six dimensions ; comparer sans fusion ; afficher source/calcul/fraîcheur/owner ; expliquer divergence ; diriger vers Priority Management.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter dimensions | lecteur | work item | 0 | données autorisées | facteurs visibles | non |
| Modifier priority | coordinateur | Incident/Task | 2 | justification | priority effective | OPEN-013 |
| Confirmer impact | IC/Business Owner | Incident | 2 | source/certitude | impact mis à jour | OPEN-013 |
| Contester projection | analyste | severity/confidence | 2 | motif | retour owner source | non |

## 13. Automatisation et IA
| Fonction | Humain | Règle | Moteur | Workflow | Agent | Govern | Sans IA |
|---|---|---|---|---|---|---|---|
| Calculer facteurs | correction | possible | oui, explicable | possible | proposition | non | sources/règles |
| Modifier priority/impact | décision humaine | policy possible | validation | possible | jamais silencieux | OPEN-013 | mutation manuelle |

## 14. États fonctionnels
`dimensions-complete`, `dimensions-partial`, `priority-proposed`, `priority-effective`, `source-conflict`, `unknown`.

## 15. États d’interface
Partial nomme la dimension absente ; Error conserve les facteurs valides ; Offline bloque mutation ; Permission denied masque seulement le facteur protégé ; Stale montre source/date. Rendu DS.

## 16. Sorties
| Sortie | Objet/événement | Consommateur | Garantie |
|---|---|---|---|
| Coordination factors | projection | Queue/Mission Control | dimensions non confondues |
| Priority/impact event | update | audit/consumers | source, before/after, justification |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte | Retour |
|---|---|---|---|---|
| Coordination | severity/confidence contestée | source/Investigate | objet, facteur, raison | file |
| Coordination | impact exige autorité | Govern | Incident, facteurs, action | aucune Decision locale |

## 18. Dépendances
CAP-CMD-002, 105, 204, 205 et Metrics Engine.

## 19. Source de vérité
Priority/impact autorisés : Command. Severity/confidence/service/SLA : sources propriétaires. Calculs exposent définitions/version/facteurs.

## 20. Provenance et audit
Chaque dimension expose source/owner/fraîcheur ; mutations et dispositions enregistrent acteur et before/after.

## 21. Permissions fonctionnelles
`perm.command.read`, `perm.command.coordinate`, `perm.command.incident.manage`, permissions source. Atomisation reportée.

## 22. Limites et erreurs
Facteur absent/stale, définition incompatible, source refusée ou conflit ne doivent jamais produire un score total opaque.

## 23. Métriques
Items avec dimensions sourcées ; divergences expliquées ; impacts confirmed/assumed. Aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire seulement.

## 25. Critères d’acceptation
**Given** severity élevée et impact faible, **When** l’item est inspecté, **Then** severity, priority, impact, urgence, confiance et SLA risk sont séparés et sourcés.

**Given** une proposition automatisée, **When** elle est consultée, **Then** priority effective reste inchangée avant acceptation autorisée.

**Given** aucun modèle IA, **When** la capability est utilisée, **Then** sources, règles et actions manuelles suffisent.

## 26. Questions ouvertes
Quelles dimensions sont obligatoires et qui valide confidence entre produits ? — REQ-PROD-005, REQ-PROD-013, REQ-PROD-021, REQ-UX-005. `OPEN-013` reste ouverte.

## 27. Consommateurs documentaires
Work Queue, Mission Control, Incident Detail, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions ultérieures.
