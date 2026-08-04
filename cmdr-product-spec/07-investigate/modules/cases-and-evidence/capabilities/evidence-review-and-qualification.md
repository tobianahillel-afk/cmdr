---
id: CAP-INV-108
title: Evidence Review and Qualification
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
# CAP-INV-108 — Evidence Review and Qualification

## 1. Définition
Évaluer séparément intégrité, pertinence, fiabilité, confiance, qualification et contestation d’une Evidence.

## 2. Problème utilisateur
Une étiquette unique mélange plusieurs dimensions et masque les désaccords. Les reviewers doivent justifier pourquoi une Evidence soutient ou non une Hypothesis.

## 3. Objectifs
Inspecter source/provenance/transformations, qualifier les dimensions séparées, marquer insuffisante/contestée, demander collecte supplémentaire et comparer alternatives.

## 4. Non-objectifs
Ne pas décider admissibilité juridique, confirmer Finding automatiquement, modifier source ou exécuter collecte 4B.2.

## 5. Propriétaire
Investigate / Cases and Evidence / Investigate Product Lead. La qualification d’Evidence est Investigate.

## 6. Utilisateurs
Principal : Evidence Reviewer. Secondaires : Case Analyst, Investigation Lead et consommateurs Finding/Govern.

## 7. Conditions d’entrée
Evidence accessible, source/provenance inspectable ou lacunes explicites, reviewer autorisé et claim/Hypothesis connu lorsque la pertinence est évaluée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Evidence | Investigate | objet revu | oui | version/status courants | qualification bloquée |
| Source/transformations | owners/producers | provenance | oui | version/time | integrity incomplete |
| Hypothesis/Finding context | Investigate | cible de pertinence | non | version courante | revue intrinsèque seulement |
| Alternatives/contradictions | Investigate | comparaison | non | statut courant | none found explicite |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Evidence | Investigate | versions/statut | lire/revoir |
| Artifact/source | owner respectif | origin/content ref | lire |
| Hypothesis/Finding | Investigate | claims/links | lire |
| Automation Run/Tool Call | Studio | transformations | consulter |
| Collection request | future/Govern | données demandées | projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Evidence qualification | créer/versionner | Investigate | dimensions séparées et raisons requises |
| Contest record | créer/résoudre/supersede | Investigate | ne supprime pas Evidence |
| Hypothesis–Evidence relation | qualifier | Investigate | support/contradiction/insufficient |
| Collection request draft | préparer | future 4B.2/Govern | aucune exécution 4B.1 |

## 11. Fonctionnalités
Inspecter Evidence/source, évaluer dimensions séparées, marquer insufficient/contested, relier Hypothesis, comparer alternatives, demander collecte, enregistrer reviewer/raisons et notifier Findings affectés.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Revoir | Reviewer | Evidence | 0 | read | review context | non |
| Qualifier dimensions | Reviewer | qualification | 2 | raison | assessment versionnée | OPEN-013 |
| Contester | Analyst | contest | 2 | grounds | contest record | OPEN-013 |
| Demander collecte | Reviewer | request draft | 1/3 | besoin/cible | draft | selon risque |
| Résoudre contest | Lead | contest | 2 | revue | resolution | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Completeness check | oui | oui | oui | explication | checklist |
| Reliability indicators | oui | oui | oui | proposition | source metadata/règles |
| Relevance suggestion | oui | oui | oui | oui | revue relation manuelle |
| Qualification decision | reviewer | non autonome | workflow revue | proposal | décision humaine |

## 14. États fonctionnels
`unreviewed`, `reviewing`, `qualified`, `insufficient`, `contested`, `needs-collection`, `review-complete`, `superseded-assessment`. États Draft.

## 15. États d’interface
Loading garde Evidence ; Partial nomme source/transformations ; Error conserve assessment draft ; Offline lecture ; Permission denied sans fuite ; Stale montre version dépassée.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Qualification record | assessment | Hypothesis/Finding | dimensionné, raisonné, versionné |
| Contest record | dispute | Case/Review | acteur, grounds, statut |
| Collection request draft | request context | 4B.2/Govern | target, need, risk visibles |
| Dependency notification | event | Finding owners | liens affectés nommés |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Evidence | open review | Evidence Review | source, versions, Case, links | Evidence |
| Review | qualifier relation | Hypothesis/Finding | dimensions, reason, reviewer | review |
| Evidence insufficient | request collection | future Collection/Action Request | target, missing data, risk | Case |
| Qualification changed | notify | Finding Management | old/new, affected claims | Evidence |

## 18. Dépendances
Evidence Creation, Hypothesis/Finding, future Collection 4B.2, Shared Trace/Notifications, review permissions et OPEN-013.

## 19. Source de vérité
Evidence/qualification restent Investigate ; source garde son owner ; collection execution future ; indicators automatisés sont dérivés et attribués.

## 20. Provenance et audit
Reviewer/time/version, dimensions/raisons, sources inspectées, indicators, contest/resolution et objets dépendants notifiés.

## 21. Permissions fonctionnelles
Evidence read/qualify, contest/resolve, sensitive source, request collection et séparation author/reviewer possible.

## 22. Limites et erreurs
Source unavailable, reviewer sans accès, dimensions conflictuelles, concurrent review, Finding déjà soumis, collection unavailable ou stale version.

## 23. Métriques
Assessments avec raisons par dimension, délai contest, Findings affectés, demandes de collecte et disposition des suggestions.

## 24. Classification de livraison
`defined` / `planned`, cible native. Promotion conditionnée par modèle de qualification, rôles reviewer et notifications de dépendance.

## 25. Critères d’acceptation
**Given** Evidence intacte mais peu pertinente **When** qualifiée **Then** intégrité et pertinence restent distinctes.

**Given** une Evidence contestée **When** contest enregistrée **Then** Evidence n’est pas supprimée et revue requise.

**Given** une Evidence insuffisante **When** collecte demandée **Then** un draft est créé sans exécution 4B.1.

## 26. Questions ouvertes
Vocabulaire final, admissibilité interne, séparation des tâches et OPEN-013 restent ouverts.

## 27. Consommateurs documentaires
Hypothesis, Finding, Action Request, Case Replay, future Collection et phases Objets/Permissions.
