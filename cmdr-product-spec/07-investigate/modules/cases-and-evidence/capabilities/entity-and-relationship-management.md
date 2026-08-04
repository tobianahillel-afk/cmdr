---
id: CAP-INV-104
title: Entity and Relationship Management
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
open_decisions:
  - OPEN-013
---
# CAP-INV-104 — Entity and Relationship Management

## 1. Définition
Relier et contextualiser des Entity partagées dans un Case, avec conflits et propositions de fusion explicites, sans résolution silencieuse.

## 2. Problème utilisateur
Un même hôte, compte ou IP apparaît sous plusieurs représentations. Une fusion automatique peut perdre les distinctions temporelles et de source.

## 3. Objectifs
Afficher types, sources, relations, temporalité et confiance ; annoter et pivoter ; gérer conflits et propositions de fusion sans modifier silencieusement Entity.

## 4. Non-objectifs
Ne pas posséder Entity, définir l’algorithme de résolution, fusionner automatiquement ou définir le stockage du graph.

## 5. Propriétaire
Shared Capabilities possède Entity ; Investigate possède les relations et annotations Case-scoped.

## 6. Utilisateurs
Principal : Case Analyst. Secondaires : Threat Hunter, Entity Reviewer et consommateurs Evidence/Finding.

## 7. Conditions d’entrée
Entity IDs ou candidates, sources et temps connus, Case context et permissions de lecture/relation.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Entity | Shared Capabilities | identité/projection | oui | updated-at visible | candidate unresolved |
| Observations | Events/Artifacts/Evidence | représentations sources | oui pour relation | timestamps/source | aucune fusion |
| Relationship candidates | règles/analyst | liens proposés | non | facteurs/date | lien manuel possible |
| Case context | Investigate | scope | oui | version actuelle | inspection standalone seulement |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Entity | Shared Capabilities | type, IDs, versions | consulter/naviguer |
| Telemetry Event | Shared Capabilities | observations/temps | inspecter |
| Artifact/Evidence | Investigate | sources/relations | consulter/relier |
| Incident/Case | Command/Investigate | contextes | consulter |
| Relationship | Shared/Object Linking | source, type, confiance | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Case-scoped relation | créer, annoter, supersede | Investigate | source et temporalité requises |
| Entity merge proposal | créer | Shared/Investigate proposal | aucune fusion avant revue owner |
| Annotation | créer/versionner | Investigate | n’altère pas Entity |
| Entity | aucune fusion directe | Shared Capabilities | lifecycle Shared |

## 11. Fonctionnalités
Afficher Entity/représentations, relier objets, temporalité/sources, proposer relation/fusion, signaler conflit, annoter, pivoter, voir Cases/Incidents et comparer.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter/pivoter | Analyst | Entity | 0 | read | Search context | non |
| Créer relation | Case Analyst | relationship | 2 | source/context | relation auditée | OPEN-013 |
| Annoter | Analyst | annotation | 2 | permission | version | OPEN-013 |
| Proposer fusion | Reviewer | proposal | 2 | comparaison | proposal | OPEN-013 |
| Confirmer fusion | Shared owner | Entity | 2 | future policy | hors Investigate | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Extraction Entity | oui | oui | oui | oui | parsers/règles |
| Relation candidate | oui | oui | oui | oui | graph/règles |
| Merge proposal | oui | oui | oui | oui | comparaison manuelle |
| Fusion effective | hors Investigate | owner policy | future workflow | non autonome | revue owner |

## 14. États fonctionnels
`resolved`, `candidate`, `conflicted`, `merge-proposed`, `merged`, `superseded`, `stale`. États Draft.

## 15. États d’interface
Loading garde le Case ; Partial nomme les sources ; Error conserve proposals ; Offline lecture ; Permission denied masque les identifiants ; Stale indique relations dépassées.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Case relationships | relations | Case/Graph/Timeline | sourcées et temporelles |
| Merge proposal | proposal | Shared Entity owner | facteurs, sources, auteur |
| Pivot Query | Query draft | Event Search | IDs/aliases/scope |
| Conflict record | quality event | Reviewer | non destructif |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Event/Artifact | extraire Entity | Entity inspection | source, temps, valeurs | source |
| Entity | pivot | Event Search | IDs, aliases, temps | Entity |
| Entity conflict | proposer fusion | Shared review future | représentations, sources, raisons | Case |
| Entity | lier | Case Workspace | Entity, rôle, temps, source | Entity |

## 18. Dépendances
Shared Entity/Object Linking/Graph/Resolution, CAP-INV-002, Case/Artifact/Evidence, future identity resolution et OPEN-013.

## 19. Source de vérité
Entity lifecycle reste Shared ; annotations/relations Case restent Investigate ; merge proposal n’est pas une fusion.

## 20. Provenance et audit
Représentations/sources, auteur/règle de relation, validité temporelle, facteurs de merge proposal, disposition owner et version automation.

## 21. Permissions fonctionnelles
Entity read, relationship create, merge propose/approve distincts, sensitive identifiers et cross-tenant interdit.

## 22. Limites et erreurs
Entity inaccessible, sources conflictuelles, relation invalide, race de fusion, aliases stale, permission partielle ou tenant mismatch.

## 23. Métriques
Relations sourcées/temporalisées, proposals acceptées/rejetées, silent merge cible zéro, pivots réussis et âge des conflits.

## 24. Classification de livraison
`defined` / `planned`, cible native. Promotion conditionnée par workflow owner Shared, provenance et gestion des conflits.

## 25. Critères d’acceptation
**Given** deux Entity similaires **When** la similarité est détectée **Then** une proposal seulement apparaît et les IDs restent distincts.

**Given** une Entity Case **When** l’analyste pivote **Then** Query contient IDs/aliases/temps et le contexte Case est préservé.

**Given** des attributs contradictoires **When** conflit marqué **Then** aucun Finding n’est changé silencieusement.

## 26. Questions ouvertes
Algorithme de résolution, autorité de fusion et OPEN-013 restent ouverts.

## 27. Consommateurs documentaires
Case Workspace/Entity Graph, Search/Hunt, Evidence/Finding, phases Objets et Permissions.
