---
id: CAP-INV-106
title: Attachment Handling
product: investigate
module: cases-and-evidence
owner: Investigate Product Lead
status: draft
delivery_status: proposed
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-PROD-020
open_decisions:
  - OPEN-014
---
# CAP-INV-106 — Attachment Handling

## 1. Définition
Encadrer les fichiers joints à Notes, Comments, Reports ou communications sans décider leur fusion ou distinction finale avec Artifact.

## 2. Problème utilisateur
Un Attachment collaboratif peut être utilisé dans une conclusion sans provenance suffisante. Attachment, Artifact et Evidence ont des permissions, rétentions et usages distincts.

## 3. Objectifs
Décrire upload, lecture, retrait logique, rétention, provenance, promotion explicite vers Artifact et blocage d’une conversion automatique vers Evidence.

## 4. Non-objectifs
Ne pas décider l’objet final, fusionner avec Artifact, définir stockage, qualifier Evidence ou fermer OPEN-014.

## 5. Propriétaire
Capability proposée dans Investigate ; l’ownership objet reste ouvert et peut relever de Shared Capabilities selon OPEN-014.

## 6. Utilisateurs
Principal : Case Collaborator. Secondaires : Report Author, Evidence Reviewer et futurs owners Product Architecture/Shared/Investigate.

## 7. Conditions d’entrée
Contexte Note/Comment/Report/Case, permission upload/read, source identifiée et sensibilité/rétention évaluées.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| File/reference | utilisateur ou collaboration | attachment candidate | oui | disponibilité/rétention | reject ou broken-reference |
| Parent context | Note/Comment/Report/Case | relation | oui | parent version current | pas d’orphelin silencieux |
| Uploader/source | identity/service | provenance | oui | identité au dépôt | reject anonyme |
| Sensitivity | user/policy | accès | oui | réévaluée | restrict/classify |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Note/Comment | Shared/Investigate | parent/visibility | consulter/lier |
| Report draft | Shared Reporting | section/citation | consulter/lier |
| Artifact | Investigate | éventuel objet promu | consulter relation |
| Evidence | Investigate | objet qualifié | consulter sans équivalence |
| Attachment | concept ouvert | metadata/content ref | consulter sous contrat proposé |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Attachment record | créer/versionner/retirer logiquement | ownership ouvert | parent/source obligatoires, non canonique avant OPEN-014 |
| Artifact | promotion explicite | Investigate | lineage depuis Attachment |
| Evidence | aucune création directe | Investigate | qualification séparée |
| Parent relation | créer/supersede | owner parent | versionnée et permission-aware |

## 11. Fonctionnalités
Joindre, montrer uploader/parent/temps/sensibilité/rétention, lire, retirer logiquement, demander promotion Artifact, voir promotion, bloquer usage Evidence non qualifié et exposer les options OPEN-014.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Upload | Collaborator | Attachment proposal | 1/2 | parent/permission | relation | OPEN-013 |
| Lire/télécharger | Collaborator | Attachment | 0/1 | access | contenu audité | selon données |
| Retirer logiquement | Owner | relation | 2 | policy | inactive avec trace | OPEN-013 |
| Promouvoir Artifact | Case Analyst | Artifact | 2 | provenance | Artifact explicite | OPEN-013 |
| Utiliser en Evidence | Reviewer | candidate | 2 | qualification | candidate seulement | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Extraction metadata | oui | oui | oui | résumé | extracteurs |
| Classification sensibilité | oui | oui | oui | proposition | policy/manuelle |
| Suggestion promotion | oui | oui | oui | oui | action humaine |
| Usage Evidence | oui | règles bloquantes | workflow | proposition | qualification humaine |

## 14. États fonctionnels
`proposed`, `uploaded`, `restricted`, `broken-reference`, `removed`, `promoted-to-artifact`, `candidate-evidence`, `retention-pending`. États Draft.

## 15. États d’interface
Loading conserve parent ; Partial nomme metadata ; Error évite orphelin ; Offline lecture ; Permission denied sans fuite ; Stale indique rétention/disponibilité.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Attachment relation | proposed record | parent | parent, uploader, accès visibles |
| Promotion request | Artifact context | Artifact Management | explicite, traçable, non destructive |
| Retention/access event | audit event | Policy/Audit | raison/acteur |
| Evidence candidate | candidate relation | Evidence workflow | non qualifiée |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Note/Comment | attach | Attachment Handling | parent, uploader, file, visibility | parent |
| Attachment | promote | Artifact Management | source, parent, metadata, reason | parent/attachment |
| Attachment/Artifact | qualify | Evidence Creation | origin, Case, transformations, reviewer | source |
| Report draft | include/cite | Reporting Engine | ref, redaction, permission | report |

## 18. Dépendances
OPEN-014, Shared Comments/Notes/Attachments possible, Artifact/Evidence, Reporting Engine, Retention/Versioning/Export/Audit et OPEN-013.

## 19. Source de vérité
Aucun owner final Attachment décidé ; parent reste owner ; Artifact après promotion est Investigate ; Evidence distincte ; mécanisme Shared reste option.

## 20. Provenance et audit
Uploader/source, parent/version, upload/access/remove, promotion lineage, rétention et qualification éventuelle.

## 21. Permissions fonctionnelles
Upload/read/delete logical, sensitive download, promote Artifact, use in report, Evidence candidate et retention admin future.

## 22. Limites et erreurs
Parent inaccessible, source anonyme, contenu bloqué, conflit rétention, broken ref, duplicate promotion, permission révoquée ou cross-tenant.

## 23. Métriques
Attachments avec parent/source, promotions Artifact, usages Evidence non qualifiés bloqués, broken refs et conflits de rétention.

## 24. Classification de livraison
`proposed` / `planned`. Aucune promotion avant décision OPEN-014 avec preuves objet, permissions, rétention et migration.

## 25. Critères d’acceptation
**Given** un fichier joint à une Note **When** utilisé dans une conclusion **Then** il n’est pas Evidence et une qualification est requise.

**Given** un Attachment sourcé **When** promu **Then** un Artifact distinct est créé et le lineage est conservé.

**Given** un Attachment cité **When** retiré **Then** trace et citation impactée restent visibles.

## 26. Questions ouvertes
OPEN-014 conserve quatre options ; owner, schéma, rétention et export restent indécis.

## 27. Consommateurs documentaires
Case Collaboration, Artifact, Evidence, Reporting, phases Objets/Permissions et document OPEN-014.
