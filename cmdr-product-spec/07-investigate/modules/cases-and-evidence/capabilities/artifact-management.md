---
id: CAP-INV-105
title: Artifact Management
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
  - OPEN-014
---
# CAP-INV-105 — Artifact Management

## 1. Définition
Recevoir, importer, versionner, relier et rendre inspectable un Artifact technique sans définir stockage, hashing final ou moteur d’analyse.

## 2. Problème utilisateur
Les Artifacts perdent origine, versions, dérivés ou restrictions et peuvent être confondus avec des Evidence.

## 3. Objectifs
Conserver origine, acquisition, type, versions, dérivés, relations et restrictions ; lier au Case ; ouvrir un futur workbench ; distinguer Artifact/Evidence/Attachment.

## 4. Non-objectifs
Ne pas définir stockage, déduplication, hashing ou chain of custody technique ; ne pas exécuter 4B.2 ; ne pas résoudre OPEN-014.

## 5. Propriétaire
Investigate / Cases and Evidence / Investigate Product Lead. Artifact est Investigate ; acquisition/moteurs sont reportés.

## 6. Utilisateurs
Principal : Case Analyst. Secondaires : Forensic Analyst futur, Evidence Reviewer et consommateurs Reporting/Action Request.

## 7. Conditions d’entrée
Source ou import autorisé, tenant/Case, type identifiable, permission import/read et origine/acquisition connue ou explicitement incomplete.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Content/reference | collection/import/user | objet technique | oui | disponibilité/rétention | reject ou candidate incomplete |
| Origin/acquisition | Endpoint/user/source | provenance | oui | timestamp/version | bloquer promotion Evidence |
| Case | Investigate | contexte | oui pour usage actif | permission courante | unattached seulement sous policy |
| Type/metadata | producer/analyst | classification | oui | versionnée | unknown/restricted workbench |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Artifact | Investigate | metadata, versions, derivatives | lire, inspecter, lier, exporter |
| Case | Investigate | scope/restrictions | lire/lier |
| Endpoint | Endpoint Agent/Settings projection | source context | consulter |
| Automation Run/Tool Call | Studio | provenance producer | consulter |
| Attachment | concept ouvert | source collaborative | lire relation sans identité supposée |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Artifact | créer/importer/versionner/supersede | Investigate | source et lineage requis |
| Derivative relation | créer | Investigate | parent, transformation, producer |
| Case relation | créer/supersede | Investigate | ne rend pas Evidence |
| Evidence candidate | proposer | Investigate | qualification CAP-INV-107/108 |

## 11. Fonctionnalités
Importer/recevoir, afficher origine/acquisition/type/statut, gérer versions/dérivés, relier Case/Entity/Hypothesis, ouvrir workbench futur, restreindre accès/export et proposer candidate Evidence.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Importer/recevoir | Case Analyst | Artifact | 1 | source/permission | Artifact | selon sensibilité |
| Lier Case | Analyst | relation | 2 | Case accessible | lien | OPEN-013 |
| Enregistrer dérivé | futur analyst | Artifact | 1/2 | output transformation | lineage | selon action |
| Exporter | Analyst | Export Job | 1 | permission/redaction | job | selon données |
| Proposer Evidence | Reviewer | candidate | 2 | provenance complète | candidate | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Extraction metadata | oui | oui | oui | résumé | extracteurs |
| Suggestion type | oui | oui | oui | oui | signatures/règles |
| Lien dérivé | oui | oui | oui | oui | producer IDs |
| Candidate Evidence | oui | oui | oui | oui | qualification humaine |

## 14. États fonctionnels
`received`, `metadata-pending`, `available`, `restricted`, `derivative`, `superseded`, `unavailable`, `candidate-evidence`, `retired`. États Draft.

## 15. États d’interface
Loading expose producer ; Partial nomme metadata/provenance ; Error conserve référence ; Offline lecture ; Permission denied masque content ; Stale indique disponibilité.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Artifact | objet | Case/Workbench | source, type, version, owner |
| Derivative lineage | relation | Review/Evidence | transformation/producer |
| Evidence candidate | candidate relation | Evidence Creation | non qualifiée automatiquement |
| Export request | Export Job | Analyst/Reporting | permission-aware |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Collection/import | recevoir | Artifact Management | source, acquisition, target, producer | Case/source |
| Artifact | ouvrir | Workbench 4B.2 futur | type, version, restrictions, Case | Artifact |
| Artifact | qualifier candidate | Evidence Creation | origin, Case, transformation, reason | Artifact |
| Attachment | promotion explicite | Artifact Management | parent, metadata, permission | collaboration context |

## 18. Dépendances
Case, future Collection/Workbench 4B.2, Shared Versioning/Preview/Export/Trace, Evidence et OPEN-014.

## 19. Source de vérité
Artifact metadata/lineage restent Investigate ; content storage futur ; producer records restent owners ; Evidence distincte ; Attachment ouvert.

## 20. Provenance et audit
Origin/acquisition, importer/producer/version, transformations/dérivés, liens Case/Entity, access/export et disposition candidate.

## 21. Permissions fonctionnelles
Artifact import/read/export, sensitive content, derivative create, Case link, Evidence candidate, redaction/retention.

## 22. Limites et erreurs
Content missing/corrupt, origin absent, type unknown, duplicate candidate, restricted source, parent unavailable, export denied ou tenant mismatch.

## 23. Métriques
Artifacts avec origine complète, lineage des dérivés, candidates acceptées/rejetées, accès interdits bloqués et références cassées.

## 24. Classification de livraison
`defined` / `planned`, cible native. Promotion conditionnée par IDs stables, provenance, lineage, access et futurs contrats stockage/intégrité.

## 25. Critères d’acceptation
**Given** un Artifact lié au Case **When** ouvert **Then** aucune Evidence n’est implicite et une action de qualification existe.

**Given** un dérivé **When** enregistré **Then** parent, tool/version et timestamp sont visibles et source inchangée.

**Given** un fichier joint à une Note **When** analyse demandée **Then** promotion Artifact explicite et origine Note conservée.

## 26. Questions ouvertes
Stockage/hashing/dedup sont reportés ; OPEN-014 demeure ouverte ; autorité de collecte sera traitée en 4B.2/Govern.

## 27. Consommateurs documentaires
Case Workspace, futur Workbench, Evidence, Reporting, Action Request et phases Objets/Permissions.
