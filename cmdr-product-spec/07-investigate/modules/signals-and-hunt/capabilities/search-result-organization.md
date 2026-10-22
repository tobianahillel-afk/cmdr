---
id: CAP-INV-007
title: Search Result Organization
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
# CAP-INV-007 — Search Result Organization

## 1. Définition
Organiser les résultats de recherche comme état de travail traçable, sans créer automatiquement un objet canonique ou une Evidence.

## 2. Problème utilisateur
Les analystes doivent comparer, annoter, regrouper et exclure des résultats avant de décider lesquels méritent un Case. Sans état explicite, sélection et raisons disparaissent.

## 3. Objectifs
Permettre sélection, annotation, groupement, comparaison, épinglage, exclusion justifiée, persistance, transmission Case et export permission-aware.

## 4. Non-objectifs
Ne pas créer un objet Collection sans décision ; ne pas qualifier Evidence ; ne pas remplacer Saved Views ; ne pas fixer le format d’export.

## 5. Propriétaire
Investigate / Signals and Hunt / Investigate Product Lead. Les résultats restent Search Job/Telemetry Event et les Saved Views restent Shared.

## 6. Utilisateurs
Principal : Threat Hunter. Secondaires : SOC Analyst, Case Analyst et Investigation Lead reviewer.

## 7. Conditions d’entrée
Search Job terminé ou partial, résultats accessibles, scope conservé, workspace ou Case cible identifiable et permissions valides.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Result set | Search Job | population de travail | oui | run et expiration visibles | empty/expired explicite |
| Saved View | Shared Capabilities | présentation | non | version courante | vue par défaut |
| Annotations | utilisateurs | raisonnement | non | versionnées | organisation sans note |
| Case cible | Investigate | destination | non | permission réévaluée | workspace standalone |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Search Job | Shared Capabilities | résultats, erreurs, limites | consulter et sélectionner |
| Telemetry Event | Shared Capabilities | champs et raw refs | inspecter et pivoter |
| Entity | Shared Capabilities | groupements et relations | consulter et comparer |
| Case | Investigate | scope et objets liés | consulter et transmettre |
| Saved View | Shared Capabilities | colonnes/filtres | appliquer sans posséder |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Workspace state | sélectionner, grouper, épingler, exclure | Investigate | versionné et lié au Search Job |
| Annotation | créer/modifier | Investigate/Collaboration | auteur et raison visibles |
| Case relation | transmettre références | Investigate | aucune Evidence automatique |
| Export Job | demander | Shared Capabilities | permission, redaction et source |

## 11. Fonctionnalités
Sélection multi-résultats, groupement, comparaison, épinglage, exclusion justifiée, annotations, persistance, transmission Case et export.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Sélectionner/grouper | Analyst | workspace | 0 | résultats disponibles | état local | non |
| Annoter | Analyst | annotation | 2 | collaboration autorisée | version | OPEN-013 |
| Exclure avec raison | Threat Hunter | workspace | 2 | sélection | exclusion auditée | OPEN-013 |
| Transmettre au Case | Case Analyst | relation | 2 | Case accessible | liens sourcés | OPEN-013 |
| Exporter | Lead | Export Job | 1 | permission/redaction | job | selon données |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Groupement | oui | oui | oui | oui | agrégations déterministes |
| Détection de doublons | oui | oui | oui | oui | clés et règles explicites |
| Suggestion d’exclusion | oui | oui | oui | oui | filtres et revue humaine |
| Transmission Case | oui | oui | workflow | candidate | sélection manuelle |

## 14. États fonctionnels
`unorganized`, `selected`, `grouped`, `annotated`, `excluded`, `pinned`, `linked-to-case`, `expired`. États Draft.

## 15. États d’interface
Loading conserve la sélection ; Partial nomme les partitions manquantes ; Error garde le workspace ; Offline reste en lecture ; Permission denied masque les champs ; Stale expose l’expiration.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Workspace organization | state record | Threat Hunter | versionné et lié au run |
| Selected package | références | Case Workspace | sélectif, attribué, non probatoire |
| Exclusion record | analysis event | Replay | raison et auteur |
| Export request | Export Job | analyst/reporting | permission-aware et redacted |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Event Search | organiser | result workspace | run, Query, résultats, filtres | Search |
| Result workspace | inspecter | Event Inspector | sélection et position | même sélection |
| Result workspace | transmettre | Case Workspace | références, annotations, exclusions, provenance | workspace |
| Result workspace | exporter | Shared Export | scope, champs, redactions, run | job/status |

## 18. Dépendances
Shared Search Job, Saved Views, Export, Object Linking, CAP-INV-004/008/102, Collaboration et future décision éventuelle sur un objet Collection.

## 19. Source de vérité
Résultats restent Shared ; organisation et annotations de travail restent Investigate ; Saved View et Export restent Shared ; aucun objet Collection par défaut.

## 20. Provenance et audit
Search Job/Query version, sélections, groupes, exclusions, auteur/raisons, liens Case, exports et disposition des suggestions.

## 21. Permissions fonctionnelles
Result read, annotation/link, bulk linking, export/redaction, sensitive fields et cross-tenant interdit.

## 22. Limites et erreurs
Run expiré, résultat supprimé, Case inaccessible, conflit de workspace, export partial, redaction requise ou permission révoquée.

## 23. Métriques
Packages Case avec provenance, exclusions justifiées, workspaces réutilisés, exports bloqués avant fuite et délai sélection→Case.

## 24. Classification de livraison
`defined` / `planned`, cible native. Promotion conditionnée par persistance versionnée, package Case sourcé et export permission-aware.

## 25. Critères d’acceptation
**Given** dix résultats et trois sélectionnés **When** ils sont transmis **Then** seuls trois liens sont créés et aucune Evidence n’apparaît.

**Given** un duplicate **When** il est exclu avec raison **Then** il reste traçable dans le replay.

**Given** des champs sensibles **When** un export est demandé **Then** redaction ou refus empêche toute fuite.

## 26. Questions ouvertes
Nécessité d’un objet Collection, rétention du workspace et permissions bulk/export restent ouvertes.

## 27. Consommateurs documentaires
Event Search, Case Workspace, Case Replay, Evidence candidates, Shared Export et Reporting.
