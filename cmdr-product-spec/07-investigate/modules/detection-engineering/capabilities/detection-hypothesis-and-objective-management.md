---
id: CAP-INV-403
title: Detection Hypothesis and Objective Management
product: investigate
module: detection-engineering
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-06
requirement_ids:
  - REQ-INV-006
  - REQ-PROD-014
  - REQ-PROD-019
  - REQ-PROD-020
  - REQ-AI-002
  - REQ-SEC-001
open_decisions:
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-403 — Detection Hypothesis and Objective Management
## 1. Définition
Formuler et versionner la question de détection, le comportement recherché, les critères de succès et de non-succès, les exclusions, limites, alternatives et éléments favorables ou contradictoires.

## 2. Problème utilisateur
Une Case Hypothesis ou un Finding ne définit pas automatiquement ce qu’une détection doit observer, exclure ou considérer comme résultat utile.

## 3. Objectifs
- définir comportement, question, succès, non-succès et consommateurs
- documenter sources, contextes, exclusions et limites
- maintenir hypothèses alternatives et contradictions
- relier sans confondre Case Hypothesis, Finding et Evidence
- accepter pour authoring, contester, retirer ou superseder

## 4. Non-objectifs
- aucun moteur, langage, syntaxe vendor, API, protocole, parser, compilateur, AST, modèle ML, commande ou code
- aucune promotion, deployment, activation, deactivation, rollback, exception active ou mutation Signal/Alert
- aucune capability CAP-INV-5xx, Intelligence, Cloud/Mobile ou réécriture détaillée d’écran

## 5. Propriétaire
Investigate possède Detection Hypothesis et sa disposition humaine. Command conserve runtime Detection/Signal/Alert/Incident ; Settings les sources/parsers/schemas/health/retention ; Endpoint Agent ses capacités et résultats locaux ; Studio Tool/Tool Call/Workflow/Automation Run ; Govern l’autorité future ; Shared les mécanismes génériques.

## 6. Utilisateurs
Principal : **Detection Engineer**. Secondaires : Threat Hunter; Investigation Lead; Detection Reviewer.

## 7. Conditions d’entrée
Tenant, environnement, objective, scope, versions, sources, permissions, restrictions et return origin sont explicites. Une dépendance absente produit un état incomplet, partiel ou bloqué ; aucune donnée ou autorité n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | ---: | --- | --- |
| Project objective | CAP-INV-402 | scope, consommateurs et contraintes | oui | version project courante | hypothesis `draft` |
| Case Hypothesis / Finding / Evidence | Investigate | raisonnement et preuves liées | non | versions référencées | origine non qualifiée |
| Hunt and analysis observations | CAP-INV-005 / 313,328,346,362,379,397 | comportements et contradictions | non | sources résolubles | limitations visibles |
| Reviewer context | Investigate / Security | auteur, reviewer et permissions | oui | courant | modification refusée |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Case Hypothesis | Investigate | question d’enquête et statut | lecture/lien |
| Finding / Evidence | Investigate | conclusions et sources qualifiées | lecture/lien |
| Hunt / technical observations | Investigate | comportements candidats | lecture/lien |
| Detection Engineering Project | Investigate concept | objectif et scope | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Detection Hypothesis | créer, versionner, annoter, contester, accepter, retirer, superseder | Investigate concept | distincte de Case Hypothesis |
| Hypothesis relation | lier à sources et alternatives | Investigate | relations sourcées et réversibles |
| Trace event | émettre | Shared | auteur et disposition visibles |

## 11. Fonctionnalités
- définir comportement et question de détection
- définir succès, non-succès, exclusions et limites
- maintenir alternatives, observations favorables et contradictions
- lier Case Hypothesis, Finding et Evidence
- versionner et soumettre à authoring
- conserver versions, erreurs, partialité, restrictions et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | ---: | --- | --- | --- |
| Consulter/Comparer | Detection Engineer | Detection Hypothesis | 0 | lecture autorisée | projection sourcée | non |
| Exécuter traitement borné | Detection Engineer | Tool Call / Result | 1 | déclenchement explicite et permission | résultat attribué | policy |
| Créer/Modifier/Contester | Detection Engineer | Detection Hypothesis | 2 | mutation réversible | nouvelle version | OPEN-013 |
| Préparer handoff | Detection Engineer | candidate package | 2 | sources et limites visibles | package non effectif | destination |

Classes 3/4 exclues ; production et runtime appartiennent à 4B.3A.2/owners.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | ---: | ---: | ---: | ---: | --- |
| Construire Detection Hypothesis | oui | éditeur/contrôles explicables | oui | proposition | formulaire/table/revue |
| Valider ou comparer | oui | validateur/comparateur | oui | explication | diagnostics/diff |
| Expliquer erreurs | oui | catalogue | oui | résumé sourcé | erreurs brutes/checklist |
| Promouvoir/déployer/qualifier runtime | non | non | non | interdit | future phase/owner |

Toute sortie expose initiateur, producteur/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun choix silencieux.

## 14. États fonctionnels
`draft`, `incomplete`, `under-review`, `accepted-for-authoring`, `disputed`, `withdrawn`, `superseded`. États fonctionnels, pas machine objet finale.

## 15. États d’interface
Loading conserve context/version ; Empty distingue absence et interdiction ; Partial expose gaps ; Error conserve le valide ; Offline bloque les nouveaux runs ; Permission denied masque ; Stale distingue ancien/courant ; conflits fournissent diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Detection Hypothesis | functional hypothesis | CAP-INV-404,406,416 | source, limites et auteur visibles |
| Objective criteria | criteria set | CAP-INV-410,413 | aucune ground truth supposée |
| Contradiction set | linked observations | Validation / Review | éléments favorables et contraires conservés |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| Case Hypothesis / Finding / Hunt | formaliser un besoin | CAP-INV-403 | observations, contradictions, objectif et provenance | source |
| CAP-INV-403 | acceptée pour authoring | CAP-INV-406 | comportement, conditions, exclusions, limites | project |
| CAP-INV-403 | identifier les données | CAP-INV-404 | sources et contextes requis | project |

Transitions conservent ownership, tenant/env, permissions, restrictions, versions, erreurs, provenance et return origin.

## 18. Dépendances
CAP-INV-103,109,005,401,402; analysis handoffs; Shared Linking/Versioning; OPEN-013/015. Shared Jobs/Trace/Versioning/Linking/Search/Export/Reporting/Collaboration/Comparison/Recovery consommés sans redéfinition.

## 19. Source de vérité
Investigate est source de Detection Hypothesis ; tous les objets consommés restent chez leurs owners. Draft ≠ runtime Detection.

## 20. Provenance et audit
Source du besoin, Project, versions, sources/schemas/fields/mappings/logic, Tool/Calls/Runs, paramètres, datasets, résultats, erreurs, interruptions, auteurs, reviewers, dispositions, exports et correlation ID applicables à Detection Hypothesis.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | ---: | --- | --- | --- | --- | --- |
| Detection Hypothesis read | raisonnement sensible | 0 | selon Case et Project | possible | viewer/reviewer | Investigate | Permissions |
| Detection Hypothesis create/update | orientation de détection | 2 | sources protégées masquées | OPEN-013 | author/reviewer | Investigate | Permissions |
| Hypothesis review/dispute | qualité et responsabilité | 2 | aucune clôture silencieuse | possible | reviewer distinct si requis | Investigate | Permissions |

Matrice atomique, namespaces, RBAC/ABAC, step-up et SoD finaux reportés ; permissions production exclues.

## 22. Limites et erreurs
- Detection Hypothesis ≠ Case Hypothesis.
- Elle ne confirme ni attaque, ni comportement malveillant, ni couverture.
- Une suggestion IA reste une proposition attribuée.
- source stale/restricted/partial, tenant mismatch, permission revoked, Tool/version unavailable, timeout or cancellation
- Tool result, score, match, non-match, AI output or mapping is not a conclusion by itself

## 23. Métriques conceptuelles
- volume par état/version
- partial/blocked/disputed/failed
- provenance et dispositions humaines complètes
- silent promotion/deployment count — cible zéro

Aucune cible runtime, precision/recall garantie, drift, health ou coût production.

## 24. Classification de livraison
`defined` / `planned` ; documentation only. Aucun `validated`, `implemented`, `deployed`, `active`, `native` ou `integrated`.

## 25. Critères d’acceptation
### 1. Hypothèse dérivée
**Given** une Case Hypothesis avec observations et contradictions  
**When** une Detection Hypothesis est créée  
**Then** les deux objets restent distincts et la relation source est visible

### 2. Contradiction
**Given** un Finding contredit une condition envisagée  
**When** la hypothesis est revue  
**Then** la contradiction reste visible et le statut peut être `disputed`

### 3. Sans IA
**Given** aucun modèle  
**When** l’objectif est défini  
**Then** formulaire structuré, catalogues et revue humaine restent suffisants

## 26. Questions ouvertes
- OPEN-013 reste ouverte.
- OPEN-015 reste ouverte.
- Le moteur/langage Detection est une lacune future non couverte ; OPEN-005 reste forensic-only.
- Schémas, formats, permissions et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering module, Event Search/Hunt/Case/Evidence/technical handoffs, Command boundaries, Settings/Endpoint, Studio, Govern future review, Shared, Objects/Permissions/Screens/Journeys/Technique/4B.3A.2/Validation.
