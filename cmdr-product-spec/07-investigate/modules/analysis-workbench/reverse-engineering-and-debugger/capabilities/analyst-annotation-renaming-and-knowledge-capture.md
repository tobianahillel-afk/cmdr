---
id: CAP-INV-337
title: Analyst Annotation, Renaming and Knowledge Capture
product: investigate
module: analysis-workbench
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-05
requirement_ids:
  - REQ-INV-003
  - REQ-INV-004
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-AI-002
  - REQ-PROD-061
open_decisions:
  - OPEN-013
  - OPEN-014
  - OPEN-015
source-of-truth: canonical
---

# CAP-INV-337 — Analyst Annotation, Renaming and Knowledge Capture

## 1. Définition
Analyst Annotation, Renaming and Knowledge Capture définit le comportement produit permettant de capturer annotations, commentaires, bookmarks et renommages attribués, versionnés et réversibles sans transformer une interprétation en vérité.

## 2. Problème utilisateur
de capturer annotations, commentaires, bookmarks et renommages attribués, versionnés et réversibles sans transformer une interprétation en vérité.

## 3. Objectifs
- annoter emplacements, fonctions, variables, symboles et structures
- attribuer auteur, source et niveau de confiance
- lier Hypothesis, observation ou Evidence existante
- comparer, restaurer et partager les versions selon permission
- accepter, modifier ou rejeter une suggestion automatisée

## 4. Non-objectifs
- modifier l’Artifact ou son identité canonique
- qualifier automatiquement une Evidence ou un Finding
- faire d’un nom analytique une identité certaine

## 5. Propriétaire
Investigate possède le contexte analytique, les annotations et les relations au Case. Studio conserve les Tools, Tool Calls et Automation Runs; Platform Settings administre les environnements; Govern conserve l’autorité sur toute cible réelle; Shared conserve les mécanismes transversaux.

## 6. Utilisateurs
Reverse Engineer et Malware Analyst principaux; Reviewer et Evidence Reviewer secondaires.

## 7. Conditions d’entrée
- session accessible
- cible d’annotation identifiée
- droit de contribution ou lecture évalué

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Analytical target | Reverse/Debugger capability | location, function, symbol, type, snapshot or event | Oui | version courante | annotation non créée |
| Author/source/confidence | user or automation provenance | attribution | Oui | au moment de l’action | action bloquée si attribution impossible |
| Evidence/Hypothesis reference | Investigate | relation optionnelle | Non | Case courant | annotation autonome |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | contexte, Hypothesis et return origin | lecture et liaison uniquement |
| Artifact / Derived Artifact / Runtime Artifact | Investigate | source, versions, restrictions et provenance | lecture; aucun original modifié |
| Tool / Tool Call / Automation Run | CMDR Studio | outil, version, exécution et attribution | lecture et sélection; lifecycle Studio |
| Reverse/Debugger observations | CAP-INV-331..343 | emplacements, fonctions, valeurs et événements | lecture |
| Existing Evidence | Investigate | references qualifiées | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Trace / Activity event | Shared | émission de l’action et de son résultat | append-only; aucune suppression locale |
| Annotation / Comment / Bookmark / Rename | Investigate | create/update/version/restore | connaissance analytique, jamais identité certaine |
| Suggestion disposition | Investigate | accept/modify/reject | producteur et sources conservés |

## 11. Fonctionnalités
- ajouter annotation, commentaire et bookmark
- renommer fonction, variable, symbole ou structure
- attribuer confiance et source
- lier Hypothesis, observation ou Evidence existante
- versionner, comparer, restaurer et exporter selon policy
- gérer accept/modify/reject des suggestions

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Ajouter une annotation | Analyst | Annotation | 2 | cible accessible | version attribuée | OPEN-013 |
| Renommer analytiquement | Reverse Engineer | Rename record | 2 | session modifiable | nom local, historique conservé | OPEN-013 |
| Restaurer une version | Reviewer | Knowledge version | 2 | version accessible | ancienne version réactivée | OPEN-013 |
| Lire/comparer | Contributor | Knowledge history | 0 | permission lecture | différences et auteurs visibles | Non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Inspecter les connaissances analytiques | Oui | Oui | Oui | Oui | viewer déterministe et navigation manuelle de les connaissances analytiques |
| Proposer un nom, un commentaire ou une synthèse | Oui | Oui | Oui | Oui | règles, heuristiques visibles, comparaison et revue humaine |
| Préparer un handoff | Oui | Oui | Oui | Oui | sélection manuelle, checklist et validation humaine |

Toute sortie automatisée expose initiateur, producteur et version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres, timestamp, statut, incertitude, owner humain et disposition acceptée, modifiée ou rejetée.

## 14. États fonctionnels
- draft
- active
- conflicting
- accepted
- modified
- rejected
- superseded
- restored
- withdrawn

Ces états sont fonctionnels et ne constituent pas une machine d’état objet définitive.

## 15. États d’interface
- **Loading** conserve le Workbench, l’Artifact, la sélection et le return origin.
- **Empty** explique l’absence de résultat sans simuler une analyse.
- **Partial** identifie les sources, vues ou événements manquants et leurs conséquences.
- **Error** conserve les résultats valides, l’erreur et une reprise sûre.
- **Offline** limite les mutations et affiche la dernière synchronisation.
- **Permission denied** ne révèle aucune donnée protégée.
- **Stale** distingue la dernière observation connue de l’état courant.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Versioned knowledge record | Annotation/Rename concept | Reverse and Debug views | auteur, source, confiance et historique |
| Suggestion disposition | Automation review event | Studio provenance / audit | acceptation, modification ou rejet explicite |
| Share/export package | Export result | authorized collaborator | restrictions et redaction appliquées |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Any Reverse/Debugger view | annoter/renommer | Knowledge record | target, selection, source, confidence | retour vue et focus |
| Automated suggestion | review humaine | Knowledge record | producer/version, sources, Tool Calls, uncertainty | retour suggestion |
| Knowledge record | lier Evidence/Hypothesis | Case workspace | relation et justification | retour Workbench |

Chaque transition conserve tenant, environnement, Case, Artifact, sélection, permissions et return origin. Une erreur ne détruit pas la source ni les résultats déjà valides.

## 18. Dépendances
- CAP-INV-103 Hypothesis
- CAP-INV-107 Evidence
- CAP-INV-331..343
- Shared Collaboration/Versioning/Export
- OPEN-013/014/015

## 19. Source de vérité
Les annotations et renommages appartiennent à Investigate; les références Evidence restent qualifiées dans leurs capabilities; les suggestions conservent provenance Studio.

## 20. Provenance et audit
Tracer auteur/producteur, target, ancienne/nouvelle valeur, source, confiance, justification, version, accept/modify/reject, partage, export et restauration.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up éventuel | Séparation des tâches | Owner | Phase propriétaire |
|---|---|---:|---|---|---|---|
| annotation/comment create | connaissance partagée | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| function/symbol/type rename | interprétation réversible | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| knowledge share/export | diffusion sensible | 1 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| Evidence link | relation à contenu qualifié | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |

La matrice atomique, les namespaces et le modèle RBAC/ABAC restent reportés à la phase Permissions.

## 22. Limites et erreurs
- conflit conservé, pas d’écrasement silencieux
- suggestion IA jamais appliquée automatiquement
- ancien nom toujours historique
- Evidence liée non modifiée
- export soumis à policy

## 23. Métriques
- annotations actives/superseded
- renommages et restores
- conflits
- suggestions acceptées/modifiées/rejetées

## 24. Classification de livraison
Delivery status `defined`; delivery mode `planned`. La promotion exige objets et permissions approuvés, Tools et environnements évalués, contrats techniques, tests de sécurité, écrans et release evidence. Aucun moteur, debugger, produit tiers ou implémentation n’est choisi.

## 25. Critères d’acceptation
### Scénario 1
**Given** une suggestion IA de renommage
**When** l’analyste l’examine
**Then** agent/version/sources/Tool Calls visibles; aucun apply automatique; ancien nom conservé

### Scénario 2
**Given** deux noms concurrents
**When** le reviewer compare
**Then** auteurs, sources et confiance restent visibles avant disposition

### Scénario 3
**Given** aucun fournisseur IA
**When** l’analyste annote
**Then** toutes les fonctions manuelles de connaissance et versioning restent disponibles

## 26. Questions ouvertes
- OPEN-013 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-014 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-015 reste ouverte; aucune décision n’est fermée par cette spécification.

Les schémas, cardinalités, machines d’état, formats d’adresse et permissions atomiques sont reportés aux phases propriétaires.

## 27. Consommateurs documentaires
- INV-REV-001
- INV-DBG-001
- Case Workspace and Evidence Board
- future Annotation object specification
