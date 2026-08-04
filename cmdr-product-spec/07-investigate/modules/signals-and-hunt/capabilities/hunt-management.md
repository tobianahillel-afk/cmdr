---
id: CAP-INV-005
title: Hunt Management
product: investigate
module: signals-and-hunt
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
# CAP-INV-005 — Hunt Management

## 1. Définition
Organiser une activité de recherche orientée question, scope et Hypotheses, sans imposer encore un objet Hunt canonique.

## 2. Problème utilisateur
Les recherches liées à une même question se dispersent entre Queries, résultats, notes et Cases. Sans cadre Hunt, les conclusions, contributeurs et limites sont difficiles à rejouer ou transmettre.

## 3. Objectifs
- définir question, scope, période, Hypotheses, recherches, résultats, contributeurs et conclusion ;
- promouvoir sélectivement vers un Case ;
- préserver provenance et réutilisation.

## 4. Non-objectifs
- ne pas définir un schéma Hunt final ;
- ne pas remplacer Case ;
- ne pas créer une Work Queue générale ;
- ne pas définir Detection Engineering.

## 5. Propriétaire
Investigate / Signals and Hunt / Investigate Product Lead. Le workspace Hunt est Investigate ; son éventuel statut d’objet canonique reste ouvert.

## 6. Utilisateurs
Principal : Threat Hunter. Secondaires : Investigation Lead, SOC Analyst et Case Analyst consommateur.

## 7. Conditions d’entrée
Question explicite, tenant/environnement, période, owner, sources ou Query disponibles et permissions valides.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Question de Hunt | utilisateur ou Incident | objectif analytique | oui | versionnée | rester draft |
| Scope et période | utilisateur ou Case | périmètre | oui | timezone et sources visibles | bloquer activation |
| Queries et Search Jobs | Shared Capabilities | travail de recherche | non | versions et runs visibles | permettre planification |
| Hypotheses | Investigate | raisonnement | non | statut courant | signaler la lacune |
| Contributeurs | Collaboration/Presence | responsabilités | oui pour active | disponibilité visible | rester unowned |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Query | Shared Capabilities | version et paramètres | consulter et exécuter via Event Search |
| Search Job | Shared Capabilities | runs, résultats et erreurs | consulter et comparer |
| Hypothesis | Investigate | question, statut et Evidence | consulter et relier |
| Case | Investigate | scope et objets liés | consulter ou promouvoir vers |
| Incident | Command | contexte opérationnel | consulter sans coordonner |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Hunt workspace state | créer, modifier, suspendre ou clôturer | Investigate | reste workspace tant que l’objet Hunt n’est pas décidé |
| Hypothesis | créer ou lier | Investigate | distincte d’un Finding |
| Case relation | promouvoir résultats et Queries | Investigate | sources et provenance conservées |
| Query Asset relation | lier et versionner | Investigate/Shared | la Query reste Shared |

## 11. Fonctionnalités
Créer un Hunt, lier Queries/runs/résultats/Hypotheses, gérer contributeurs, états, conclusion, limites, promotion Case et réutilisation.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Créer Hunt | Threat Hunter | workspace | 2 | question et scope | draft | OPEN-013 |
| Activer ou suspendre | Hunt Lead | état | 2 | owner et sources | transition auditée | OPEN-013 |
| Ajouter Query ou résultat | Contributor | relation | 2 | objet accessible | relation sourcée | OPEN-013 |
| Promouvoir vers Case | Hunt Lead | Case | 2 | sélection | liens ou Case | OPEN-013 |
| Clôturer | Hunt Lead | état | 2 | conclusion ou limites | completed/inconclusive | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Agréger les runs | oui | oui | oui | résumé | timeline et statistiques déterministes |
| Proposer Hypothesis | oui | oui | oui | oui | création manuelle et règles |
| Proposer pivot | oui | oui | oui | oui | Query Assets et templates |
| Préparer conclusion | oui | faits seulement | workflow de revue | brouillon | synthèse humaine |

## 14. États fonctionnels
`draft`, `active`, `paused`, `blocked`, `completed`, `inconclusive`, `archived`. États Draft ; machine finale reportée à la phase Objets.

## 15. États d’interface
Loading conserve le draft ; Partial nomme les sources manquantes ; Error conserve correlation ID ; Offline reste en lecture ; Permission denied ne divulgue rien ; Stale affiche la date.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Hunt workspace | état de workspace | Threat Hunters | question, scope, owner et version |
| Conclusion | conclusion record | Case/Reporting | sources, limites et auteur |
| Promotion package | relations Query/résultat/Hypothesis | Case Lifecycle | sélectif et sourcé |
| Query set réutilisable | liens Query Assets | futurs Hunts | versions et prérequis visibles |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Signal ou Incident | démarrer Hunt | Hunt workspace | tenant, environnement, source, question, période | objet source |
| Hunt | promouvoir | Case Lifecycle | question, scope, Queries, résultats, Hypotheses, auteurs | Hunt |
| Hunt | ouvrir Query | Event Search | Query/version, paramètres, contexte | Hunt avec run lié |
| Hunt archivé | réutiliser | nouveau Hunt ou Query Asset | conclusion, sources, assets | archive inchangée |

## 18. Dépendances
CAP-INV-002/003/006/008, Case Lifecycle, Hypothesis Management, Shared Query/Search Job/Collaboration/Trace et future décision objet Hunt.

## 19. Source de vérité
Queries et Search Jobs restent Shared ; Hypotheses et Cases restent Investigate ; le workspace Hunt est Investigate ; aucun objet Hunt n’est créé implicitement.

## 20. Provenance et audit
Question/scope versions, owner/contributeurs, Query versions, Search Job IDs, inclusions/exclusions, assistance automatisée et promotion Case.

## 21. Permissions fonctionnelles
Hunt create/manage/archive, contributor add/remove, Query execute, Case create/link, share/export. Matrice atomique reportée.

## 22. Limites et erreurs
Source indisponible, Query expirée, contributor sans permission, Case inaccessible, scope conflict, stale results, clôture sans conclusion, tenant change.

## 23. Métriques
Temps création→première Query, complétude question/scope, taux de promotion Case, causes inconclusive et réutilisation d’assets.

## 24. Classification de livraison
`defined` / `planned`, cible native. Promotion conditionnée par versioning, relations traçables et décision future sur l’objet Hunt.

## 25. Critères d’acceptation
**Given** un Hunt actif avec résultats et Hypotheses **When** le lead promeut **Then** le package est sélectif, le Case appartient à Investigate, les runs restent sourcés et le retour est conservé.

**Given** des sources partielles **When** le lead clôture **Then** l’état est `inconclusive`, les limites sont visibles et aucun Finding n’est confirmé.

**Given** un Hunt archivé **When** un asset est réutilisé **Then** les versions sont référencées et l’archive reste inchangée.

## 26. Questions ouvertes
Le statut objet/workspace/activity de Hunt, la rétention et les permissions de partage restent ouverts.

## 27. Consommateurs documentaires
Signals and Hunt, Case Lifecycle, Query Assets, Case Replay, futures phases Detection Engineering et Readiness.
