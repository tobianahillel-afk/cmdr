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
Identifie et suit un empêchement opérationnel comme état ou relation d’un Incident/Task, ou comme Task de suivi, sans créer un objet Blocker concurrent.

## 2. Problème utilisateur
Les équipes voient qu’un travail n’avance plus sans toujours connaître cause, owner, délai et prochaine action. Le blocage disparaît alors des priorités et handovers.

## 3. Objectifs
Nommer type/cause/objet/owner, lier prochaine action et échéance, escalader/résoudre avec audit et réutiliser Incident/Task.

## 4. Non-objectifs
Ne crée pas d’objet Blocker final, ne modifie pas un objet source externe, ne résout pas automatiquement une dépendance et ne confond pas tout état pending avec un blocage.

## 5. Propriétaire
Command / Mission Control / Command Product Lead pour le suivi opérationnel.

## 6. Utilisateurs
Principal : Incident Commander. Secondaires : SOC Analyst L2, Task owner, Team Lead.

## 7. Conditions d’entrée
Incident ou Task accessible, raison connue ou explicitement inconnue, owner ou équipe de reprise.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Blocked work | Incident ou Task | référence de travail | oui | version courante | aucun blocage créé sans objet source |
| Blocker reason | utilisateur ou projection externe | cause structurée | oui | au moment de la déclaration | `unknown-cause` et Task de qualification |
| Dependency | objet ou acteur lié | relation | non | fraîcheur déclarée | blocage conservé sans dépendance inventée |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Incident / Task | Command | état, owner, dépendances et prochaine action | consulter et modifier |
| Decision / Case / dépendance externe | produit propriétaire | statut pertinent | consulter et relier en lecture seule |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Incident / Task | marquer bloqué, cause, owner et prochaine action | Command | classe 2, version et audit |
| Task | créer une action de déblocage | Command | réutiliser Task, sans objet concurrent |
| Dépendance externe | aucune mutation | produit propriétaire | projection uniquement |

## 11. Fonctionnalités
Déclarer/catégoriser, lier une dépendance, affecter résolution/échéance, escalader sans changer implicitement l’owner et résoudre avec preuve de reprise.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Déclarer bloqué | owner ou coordinateur | Incident/Task | 2 | raison et prochaine action | état bloqué audité | OPEN-013 |
| Créer Task de déblocage | coordinateur | Task | 2 | aucune Task équivalente | Task liée | OPEN-013 |
| Escalader | coordinateur | Incident/Task | 2 | destinataire et délai | escalade tracée | selon destination |
| Résoudre | owner | Incident/Task | 2 | cause traitée | reprise et historique | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Détecter un blocage potentiel | oui | règle possible | oui | suggestion attribuée | observation humaine et règles d’état |
| Déclarer ou catégoriser | oui | validation possible | workflow possible | proposition seulement | formulaire manuel et types explicites |
| Créer une action de déblocage | oui | déduplication | oui | brouillon Task | création manuelle de Task |
| Résoudre le blocage | oui | contrôle de version | workflow possible | non décisionnelle | action humaine avec preuve |

## 14. États fonctionnels
`identified`, `owned`, `escalated`, `mitigation-in-progress`, `resolved`, `unknown-cause`.

## 15. États d’interface
Partial nomme cause/dépendance manquante ; Offline bloque mutation ; Permission denied masque la dépendance ; Stale impose validation ; Error conserve les données Command valides.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Blocker state | mise à jour Incident/Task | Work Queue, Mission Control et Handover | cause, owner, prochaine action et fraîcheur visibles |
| Deblocking Task | Task Command | Work Queue et Readiness | source, résultat attendu et owner conservés |
| Escalation context | événement de transition | Investigate ou Govern | attribué, permission-aware et non destructif |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Operational Blocker | besoin d’expertise | Investigate | Incident, cause, délai et objets liés | objet bloqué restauré |
| Operational Blocker | besoin d’autorité | Govern | contexte, impact, urgence et action proposée | owner Command conservé jusqu’à acceptation |
| Operational Blocker | action interne requise | Task Coordination | source, cause et résultat attendu | Incident/Task source |

## 18. Dépendances
CAP-CMD-110, CAP-CMD-107, Object Linking Service et Notification Center ; `OPEN-013` reste ouverte.

## 19. Source de vérité
État de coordination et Tasks : Command. Decision/Case/dépendance externe : propriétaires sources.

## 20. Provenance et audit
Acteur/producteur, cause, objet, dépendance, before/after, justification, résultat, tenant et correlation ID.

## 21. Permissions fonctionnelles
`perm.command.coordinate`, `perm.command.incident.manage`, `perm.command.task.manage` ; atomisation et step-up reportés.

## 22. Limites et erreurs
Cause inconnue, dépendance stale/inaccessible, conflit de version, tenant incompatible ou refus ne produisent jamais une résolution fictive.

## 23. Métriques
Travaux bloqués sans owner/prochaine action, âge par cause et taux de reprise vérifiée ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; aucune implémentation prouvée.

## 25. Critères d’acceptation
**Given** un Incident bloqué, **When** un coordinateur déclare cause/owner/prochaine action, **Then** la mutation et l’historique sont audités.

**Given** une cause inconnue, **When** le blocage est créé, **Then** `unknown-cause` et une action de qualification restent visibles.

**Given** aucun modèle IA, **When** la capability est utilisée, **Then** types, formulaire, règles et Tasks manuelles suffisent.

## 26. Questions ouvertes
Quels attributs appartiennent à Incident/Task et quand une dépendance devient-elle une Task ? — Requirement IDs ci-dessus ; `OPEN-013` reste ouverte.

## 27. Consommateurs documentaires
Mission Control, Unified Work Queue, Handover, Readiness, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions ultérieures.