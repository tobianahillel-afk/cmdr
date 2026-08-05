---
id: CAP-INV-406
title: Detection Content Authoring
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
# CAP-INV-406 — Detection Content Authoring
## 1. Définition
Créer et versionner un Detection Content Draft décrivant fonctionnellement sources, champs, logique, exclusions, enrichissements, paramètres, consommateurs, sévérité proposée, limites et dépendances sans imposer de langage ou moteur.

## 2. Problème utilisateur
Une Search Query, un Saved Search ou une note de Hunt ne constitue pas un contenu de détection versionné, documenté et testable.

## 3. Objectifs
- créer un draft et sélectionner une famille fonctionnelle
- définir problème, sources, champs, conditions, exclusions, enrichissements et paramètres
- documenter consommateurs, sévérité proposée, limites et dépendances
- lier Detection Hypothesis, tests et couverture
- versionner, comparer, cloner avec lineage, retirer et préparer validation

## 4. Non-objectifs
- aucun moteur, langage, syntaxe vendor, API, protocole, parser, compilateur, AST, modèle ML, commande ou code
- aucune promotion, deployment, activation, deactivation, rollback, exception active ou mutation Signal/Alert
- aucune capability CAP-INV-5xx, Intelligence, Cloud/Mobile ou réécriture détaillée d’écran

## 5. Propriétaire
Investigate possède Detection Content Draft et sa disposition humaine. Command conserve runtime Detection/Signal/Alert/Incident ; Settings les sources/parsers/schemas/health/retention ; Endpoint Agent ses capacités et résultats locaux ; Studio Tool/Tool Call/Workflow/Automation Run ; Govern l’autorité future ; Shared les mécanismes génériques.

## 6. Utilisateurs
Principal : **Detection Engineer**. Secondaires : Detection Reviewer; Threat Hunter; Content Owner.

## 7. Conditions d’entrée
Tenant, environnement, objective, scope, versions, sources, permissions, restrictions et return origin sont explicites. Une dépendance absente produit un état incomplet, partiel ou bloqué ; aucune donnée ou autorité n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | ---: | --- | --- |
| Detection Hypothesis | CAP-INV-403 | comportement, objectifs, exclusions et limites | oui | version acceptée pour authoring | draft `incomplete` |
| Readiness and field review | CAP-INV-404/405 | sources, fields, mappings, gaps et versions | oui | assessments liés | `blocked` ou warnings |
| Project context | CAP-INV-402 | owner, contributeurs, consommateurs et scope | oui | version courante | création refusée |
| Existing draft/version | Investigate concept | contenu, lineage et dispositions | non | version sélectionnée | nouveau draft |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Detection Engineering Project / Hypothesis | Investigate concepts | objectif, scope et relations | lecture/lien |
| Data Source / Field projections | Settings / Shared | sources, versions et mappings | lecture |
| Query / Saved Search / Hunt | Shared / Investigate | inspiration ou preuve exploratoire | lecture/référence seulement |
| Tool / Automation Run | Studio | assistance et validation éventuelles | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Detection Content Draft | créer, modifier, cloner, versionner, retirer, superseder | Investigate concept | distinct du runtime Detection |
| Detection Content Version | créer comme snapshot fonctionnel | Investigate concept | aucune version active ou déployée |
| Trace / Version event | émettre | Shared | lineage et auteur conservés |

## 11. Fonctionnalités
- créer un draft par famille événement/threshold/sequence/correlation/aggregation/behavior/anomaly candidate
- définir sources, fields, conditions, exclusions, enrichissements et paramètres
- documenter consommateurs, sévérité proposée et limites
- lier hypothesis, tests et coverage
- versionner, diff, clone avec lineage et retirer
- conserver versions, erreurs, partialité, restrictions et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | ---: | --- | --- | --- |
| Consulter/Comparer | Detection Engineer | Detection Content Draft | 0 | lecture autorisée | projection sourcée | non |
| Exécuter traitement borné | Detection Engineer | Tool Call / Result | 1 | déclenchement explicite et permission | résultat attribué | policy |
| Créer/Modifier/Contester | Detection Engineer | Detection Content Draft | 2 | mutation réversible | nouvelle version | OPEN-013 |
| Préparer handoff | Detection Engineer | candidate package | 2 | sources et limites visibles | package non effectif | destination |

Classes 3/4 exclues ; production et runtime appartiennent à 4B.3A.2/owners.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | ---: | ---: | ---: | ---: | --- |
| Construire Detection Content Draft | oui | éditeur/contrôles explicables | oui | proposition | formulaire/table/revue |
| Valider ou comparer | oui | validateur/comparateur | oui | explication | diagnostics/diff |
| Expliquer erreurs | oui | catalogue | oui | résumé sourcé | erreurs brutes/checklist |
| Promouvoir/déployer/qualifier runtime | non | non | non | interdit | future phase/owner |

Toute sortie expose initiateur, producteur/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun choix silencieux.

## 14. États fonctionnels
`draft`, `incomplete`, `structurally-valid`, `valid-with-warnings`, `semantically-reviewed`, `test-ready`, `replay-ready`, `review-ready`, `blocked`, `invalid`, `superseded`, `withdrawn`. États fonctionnels, pas machine objet finale.

## 15. États d’interface
Loading conserve context/version ; Empty distingue absence et interdiction ; Partial expose gaps ; Error conserve le valide ; Offline bloque les nouveaux runs ; Permission denied masque ; Stale distingue ancien/courant ; conflits fournissent diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Detection Content Draft | draft/version | CAP-INV-407..417 | non déployé et auteur visible |
| Authoring diff | version comparison | Reviewer / Project | lineage et changements conservés |
| Validation candidate | draft reference | CAP-INV-411 | sources et dépendances explicites |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-403/405 | démarrer authoring | CAP-INV-406 | hypothesis, fields, mappings, sources et limites | project |
| CAP-INV-406 | concevoir la logique | CAP-INV-407/408/409 | draft, dependencies et objective | authoring |
| CAP-INV-406 | préparer validation | CAP-INV-411 | version, sources, fields, logic et metadata | authoring |
| CAP-INV-406 | préparer tests | CAP-INV-412 | version, behavior, variants et limitations | authoring |

Transitions conservent ownership, tenant/env, permissions, restrictions, versions, erreurs, provenance et return origin.

## 18. Dépendances
CAP-INV-003,403..405; Shared Versioning/Comparison; Studio optional Tools; OPEN-013/015. Shared Jobs/Trace/Versioning/Linking/Search/Export/Reporting/Collaboration/Comparison/Recovery consommés sans redéfinition.

## 19. Source de vérité
Investigate est source de Detection Content Draft ; tous les objets consommés restent chez leurs owners. Draft ≠ runtime Detection.

## 20. Provenance et audit
Source du besoin, Project, versions, sources/schemas/fields/mappings/logic, Tool/Calls/Runs, paramètres, datasets, résultats, erreurs, interruptions, auteurs, reviewers, dispositions, exports et correlation ID applicables à Detection Content Draft.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | ---: | --- | --- | --- | --- | --- |
| Detection Content read | contenu sensible et techniques internes | 0 | selon project | possible | viewer/reviewer | Investigate | Permissions |
| Detection Content create/update/clone/archive | mutation du draft | 2 | secrets et exemples exclus | OPEN-013 | author/reviewer | Investigate | Permissions |
| Automated authoring request | suggestion potentiellement erronée | 1/2 | sources et output visibles | possible | requester/reviewer | Studio/Investigate | Permissions |

Matrice atomique, namespaces, RBAC/ABAC, step-up et SoD finaux reportés ; permissions production exclues.

## 22. Limites et erreurs
- Detection Content Draft ≠ runtime Detection; Rule Draft ≠ deployed rule.
- Detection Rule ≠ Search Query; Saved Search ≠ Detection Content.
- Aucun langage, syntaxe propriétaire, moteur, code ou état deployed/active/production.
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
### 1. Sans IA
**Given** des sources accessibles et aucun fournisseur de modèle  
**When** un draft est créé  
**Then** l’éditeur structuré, les conditions, séquences, seuils, validations, tests et versions restent disponibles

### 2. Query source
**Given** une Query exploratoire issue d’une Hunt  
**When** elle est référencée dans un draft  
**Then** la Query reste distincte, aucune conversion silencieuse n’a lieu

### 3. Draft incomplet
**Given** un field requis est absent  
**When** la validation d’authoring est demandée  
**Then** le draft reste `incomplete` ou `blocked` et aucune promotion n’est possible

## 26. Questions ouvertes
- OPEN-013 reste ouverte.
- OPEN-015 reste ouverte.
- Le moteur/langage Detection est une lacune future non couverte ; OPEN-005 reste forensic-only.
- Schémas, formats, permissions et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering module, Event Search/Hunt/Case/Evidence/technical handoffs, Command boundaries, Settings/Endpoint, Studio, Govern future review, Shared, Objects/Permissions/Screens/Journeys/Technique/4B.3A.2/Validation.
