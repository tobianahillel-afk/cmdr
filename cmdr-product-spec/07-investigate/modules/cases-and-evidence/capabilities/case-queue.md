---
id: CAP-INV-101
title: Case Queue
product: investigate
module: cases-and-evidence
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-PROD-020
open_decisions: []
---
# CAP-INV-101 — Case Queue

## 1. Définition
Présenter et organiser les Cases Investigate sans reproduire la Work Queue générale de Command.

## 2. Problème utilisateur
Les analystes doivent trouver les Cases actifs, bloqués ou en revue. Une file mal délimitée dupliquerait Tasks, Decisions et Response Runs.

## 3. Objectifs
Recherche, filtres, tri, Saved Views, statut, owner, Incident parent, priorité contextuelle, fraîcheur, prochaine action, Findings, activité et ouverture du Case Workspace.

## 4. Non-objectifs
Ne pas remplacer Work Queue ; ne pas posséder Tasks générales, Decisions ou Runs ; ne pas fixer les colonnes finales.

## 5. Propriétaire
Investigate / Cases and Evidence / Investigate Product Lead. Saved Views restent Shared ; Incident reste Command.

## 6. Utilisateurs
Principal : Case Analyst. Secondaires : Investigation Lead, SOC Analyst et Incident Commander en projection.

## 7. Conditions d’entrée
Tenant/environnement autorisés, permission Case read, projections disponibles ou lacunes visibles et Saved View réévaluée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Cases | Investigate | objets de file | oui | updated-at par Case | empty ou partial explicite |
| Incident parent | Command | contexte | non | projection datée | Case sans contexte parent |
| Saved View | Shared Capabilities | configuration | non | version/permissions | vue par défaut |
| Findings/activity | Investigate/Shared | indicateurs | non | timestamp | unavailable explicite |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | statut, owner, next action, fraîcheur | consulter, filtrer, trier, ouvrir |
| Incident | Command | priorité, impact, owner | consulter comme contexte |
| Finding | Investigate | count, statut, dernier changement | consulter |
| Saved View | Shared Capabilities | filtres, tri, colonnes | appliquer/personnaliser |
| Task | Command | relation éventuelle | consulter sans agréger |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Saved View reference | sélectionner/personnaliser | Shared Capabilities | ne possède pas les Cases |
| Case assignment/status | demander mutation | Investigate | via CAP-INV-102 |
| Queue selection state | modifier | workspace Investigate | pas d’objet métier |
| Task/Decision/Run | aucune création | Command/Govern | Case Queue ≠ Work Queue |

## 11. Fonctionnalités
Lister les Cases accessibles, rechercher, filtrer, trier, appliquer Saved Views, afficher contexte et fraîcheur, ouvrir le workspace et signaler partial/stale.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter/filtrer | Case Analyst | Case list | 0 | read | vue | non |
| Ouvrir Case | Case Analyst | Case | 0 | accessible | workspace | non |
| Appliquer Saved View | Analyst | Saved View | 0 | vue accessible | configuration | non |
| Affecter | Lead | Case | 2 | assignment | owner modifié | OPEN-013 |
| Exporter liste | Lead | Export Job | 1 | permission/redaction | job | selon données |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Tri/filtres | oui | oui | oui | suggestion | Saved Views |
| Priorité de revue | oui | oui | oui | proposition | facteurs visibles |
| Assignment | oui | oui | oui | proposition | règles/availability |
| Résumé Case | oui | agrégation | oui | oui | champs et timeline |

## 14. États fonctionnels
`available`, `empty`, `partial`, `stale`, `permission-filtered`, `selection-active`, `view-dirty`. États Draft.

## 15. États d’interface
Loading préserve vue ; Empty distingue absence de Case ; Partial nomme les projections ; Error conserve filtres ; Offline lecture stale ; Permission denied sans fuite.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Case list projection | queue view | Case Analysts | permission-aware et fraîcheur visible |
| Selected Case context | transition | Case Workspace | tenant, vue, sélection et retour |
| Assignment request | mutation Case | Case Lifecycle | auditée et réversible |
| Export request | Export Job | Lead/Reporting | redacted |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Global navigation | ouvrir Cases | Case Queue | tenant, env, vue | destination précédente |
| Case Queue | sélectionner | Inspector | Case ID/projections | focus restauré |
| Case Queue | ouvrir | Case Workspace | Case, Incident, vue, scroll | retour exact |
| Command Incident | ouvrir Cases liés | Queue/Workspace | Incident, tenant, filtre | Incident Detail |

## 18. Dépendances
CAP-INV-102, Shared Saved Views/Search/Inspector/Activity/Export, Command Incident projection, Queue/Case Shell et permissions Case.

## 19. Source de vérité
Case reste Investigate ; Incident reste Command ; Saved View reste Shared ; l’état de file est un workspace state.

## 20. Provenance et audit
Case version/fraîcheur, Saved View/version, assignment actor/reason, origin/return context et facteurs automatisés.

## 21. Permissions fonctionnelles
Case read/assign, Saved View use/share, export/redaction et cross-tenant interdit.

## 22. Limites et erreurs
Case inaccessible, projection Incident indisponible, Saved View interdite, assignment conflict, stale list, export partial ou tenant change.

## 23. Métriques
Temps vers Case utile, Cases avec owner/next action, usage Saved Views, conflits d’affectation et contexte stale.

## 24. Classification de livraison
`defined` / `planned`, cible native. Promotion conditionnée par contrats Case, Saved Views permission-aware et restauration de contexte.

## 25. Critères d’acceptation
**Given** Cases, Tasks, Decisions et Runs **When** Case Queue s’ouvre **Then** seuls Cases sont work items propriétaires.

**Given** une vue filtrée et un Case sélectionné **When** l’utilisateur revient **Then** vue, filtres, scroll et sélection sont restaurés.

**Given** une Saved View avec champ interdit **When** elle est appliquée **Then** aucune fuite ne se produit.

## 26. Questions ouvertes
Colonnes finales, machine d’état Case et permissions d’affectation sont reportées.

## 27. Consommateurs documentaires
Case Workspace, Command projection, Screen Capability Map, phases Objets, Permissions et Écrans.
