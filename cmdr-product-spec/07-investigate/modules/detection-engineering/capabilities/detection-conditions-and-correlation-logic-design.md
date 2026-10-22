---
id: CAP-INV-407
title: Detection Conditions and Correlation Logic Design
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
# CAP-INV-407 — Detection Conditions and Correlation Logic Design
## 1. Définition
Concevoir une représentation fonctionnelle explicable des conditions positives et négatives, groupes logiques, relations, jointures, clés candidates, exclusions et dépendances d’enrichissement.

## 2. Problème utilisateur
Une logique valide en apparence peut corréler des événements sans relation réelle, masquer des champs manquants ou présenter une association comme causalité.

## 3. Objectifs
- définir conditions positives/négatives et groupes logiques
- définir relations, jointures fonctionnelles et clés candidates
- voir champs requis, données manquantes, ambiguïtés et enrichissements
- fournir représentation structurée et alternative tabulaire
- comparer, annoter, versionner et préparer des tests

## 4. Non-objectifs
- aucun moteur, langage, syntaxe vendor, API, protocole, parser, compilateur, AST, modèle ML, commande ou code
- aucune promotion, deployment, activation, deactivation, rollback, exception active ou mutation Signal/Alert
- aucune capability CAP-INV-5xx, Intelligence, Cloud/Mobile ou réécriture détaillée d’écran

## 5. Propriétaire
Investigate possède Detection Condition and Correlation Definition et sa disposition humaine. Command conserve runtime Detection/Signal/Alert/Incident ; Settings les sources/parsers/schemas/health/retention ; Endpoint Agent ses capacités et résultats locaux ; Studio Tool/Tool Call/Workflow/Automation Run ; Govern l’autorité future ; Shared les mécanismes génériques.

## 6. Utilisateurs
Principal : **Detection Engineer**. Secondaires : Detection Reviewer; Data Engineer.

## 7. Conditions d’entrée
Tenant, environnement, objective, scope, versions, sources, permissions, restrictions et return origin sont explicites. Une dépendance absente produit un état incomplet, partiel ou bloqué ; aucune donnée ou autorité n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | ---: | --- | --- |
| Detection Content Draft | CAP-INV-406 | version, objective et source context | oui | version sélectionnée | logic `draft` |
| Field references and mappings | CAP-INV-405 | fields, types, ambiguïtés et versions | oui | review liée | `blocked` |
| Data readiness | CAP-INV-404 | availability, gaps et restrictions | oui | assessment courant | warnings/partial |
| Enrichment requirements | CAP-INV-409 | contextes et fallback | non | version liée | logique sans enrichissement |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Detection Content Draft | Investigate concept | version et objective | lecture |
| Field Reference / Mapping | Investigate/Settings projections | fields et semantic limits | lecture |
| Telemetry Event / Entity projections | Shared | relations candidates | lecture limitée |
| Correlation Engine capability | Shared | capacité générique et limites | référence sans ownership |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Detection Condition | créer, modifier, versionner, retirer | Investigate concept | expression fonctionnelle, pas syntaxe de moteur |
| Correlation Definition | créer, comparer, contester, superseder | Investigate concept | correlation ≠ causation |
| Trace / Version event | émettre | Shared | diff et auteur visibles |

## 11. Fonctionnalités
- définir conditions et groupes logiques
- définir relations, joins et candidate keys
- documenter exclusions, missing data et enrichment dependencies
- offrir vue structurée et table accessible
- comparer variantes et préparer tests
- conserver versions, erreurs, partialité, restrictions et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | ---: | --- | --- | --- |
| Consulter/Comparer | Detection Engineer | Detection Condition and Correlation Definition | 0 | lecture autorisée | projection sourcée | non |
| Exécuter traitement borné | Detection Engineer | Tool Call / Result | 1 | déclenchement explicite et permission | résultat attribué | policy |
| Créer/Modifier/Contester | Detection Engineer | Detection Condition and Correlation Definition | 2 | mutation réversible | nouvelle version | OPEN-013 |
| Préparer handoff | Detection Engineer | candidate package | 2 | sources et limites visibles | package non effectif | destination |

Classes 3/4 exclues ; production et runtime appartiennent à 4B.3A.2/owners.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | ---: | ---: | ---: | ---: | --- |
| Construire Detection Condition and Correlation Definition | oui | éditeur/contrôles explicables | oui | proposition | formulaire/table/revue |
| Valider ou comparer | oui | validateur/comparateur | oui | explication | diagnostics/diff |
| Expliquer erreurs | oui | catalogue | oui | résumé sourcé | erreurs brutes/checklist |
| Promouvoir/déployer/qualifier runtime | non | non | non | interdit | future phase/owner |

Toute sortie expose initiateur, producteur/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun choix silencieux.

## 14. États fonctionnels
`draft`, `incomplete`, `reviewable`, `ambiguous`, `conflicting`, `valid-with-warnings`, `blocked`, `superseded`, `withdrawn`. États fonctionnels, pas machine objet finale.

## 15. États d’interface
Loading conserve context/version ; Empty distingue absence et interdiction ; Partial expose gaps ; Error conserve le valide ; Offline bloque les nouveaux runs ; Permission denied masque ; Stale distingue ancien/courant ; conflits fournissent diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Condition set | Detection Condition concepts | CAP-INV-411..413 | fields et versions explicites |
| Correlation Definition | functional logic | CAP-INV-408,411,414 | aucune causalité affirmée |
| Logic comparison | comparison result | Reviewer / Project | variantes et limites visibles |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-406 | concevoir conditions | CAP-INV-407 | draft, hypothesis, fields et sources | authoring |
| CAP-INV-407 | ajouter temporalité | CAP-INV-408 | conditions, keys, order constraints et gaps | logic |
| CAP-INV-407 | valider/tester | CAP-INV-411/412 | version, fields, relations et exclusions | logic |

Transitions conservent ownership, tenant/env, permissions, restrictions, versions, erreurs, provenance et return origin.

## 18. Dépendances
CAP-INV-404..406,409; Shared Correlation/Entity/Comparison; OPEN-013/015. Shared Jobs/Trace/Versioning/Linking/Search/Export/Reporting/Collaboration/Comparison/Recovery consommés sans redéfinition.

## 19. Source de vérité
Investigate est source de Detection Condition and Correlation Definition ; tous les objets consommés restent chez leurs owners. Draft ≠ runtime Detection.

## 20. Provenance et audit
Source du besoin, Project, versions, sources/schemas/fields/mappings/logic, Tool/Calls/Runs, paramètres, datasets, résultats, erreurs, interruptions, auteurs, reviewers, dispositions, exports et correlation ID applicables à Detection Condition and Correlation Definition.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | ---: | --- | --- | --- | --- | --- |
| Detection condition read | logique sensible | 0 | selon project | possible | viewer/reviewer | Investigate | Permissions |
| Condition/correlation update | impact sur comportement | 2 | field values masquées | OPEN-013 | author/reviewer | Investigate | Permissions |
| Automated logic proposal | erreur ou sur-corrélation | 1/2 | proposition attribuée | possible | human acceptance | Studio/Investigate | Permissions |

Matrice atomique, namespaces, RBAC/ABAC, step-up et SoD finaux reportés ; permissions production exclues.

## 22. Limites et erreurs
- Correlation ≠ causation.
- Syntax valid ≠ logic correct; logic valid ≠ useful detection.
- Aucune syntaxe, AST, compilateur, moteur ou implémentation.
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
### 1. Corrélation ambiguë
**Given** deux relations candidates et des fields incomplets  
**When** l’analyste sélectionne une variante  
**Then** les alternatives et données manquantes restent visibles et aucune causalité n’est affirmée

### 2. Field inaccessible
**Given** une condition dépend d’un field restreint  
**When** la logique est ouverte  
**Then** la dépendance reste visible, la valeur est masquée et l’état peut être `blocked`

### 3. Sans IA
**Given** aucun modèle  
**When** la corrélation est conçue  
**Then** builder structuré, table, catalogues et revue humaine restent disponibles

## 26. Questions ouvertes
- OPEN-013 reste ouverte.
- OPEN-015 reste ouverte.
- Le moteur/langage Detection est une lacune future non couverte ; OPEN-005 reste forensic-only.
- Schémas, formats, permissions et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering module, Event Search/Hunt/Case/Evidence/technical handoffs, Command boundaries, Settings/Endpoint, Studio, Govern future review, Shared, Objects/Permissions/Screens/Journeys/Technique/4B.3A.2/Validation.
