---
id: CAP-INV-409
title: Detection Enrichment and Context Requirement Design
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
# CAP-INV-409 — Detection Enrichment and Context Requirement Design
## 1. Définition
Définir les contextes, Entities, relations, informations Asset/Endpoint/Identity, historiques et projections Case/Incident nécessaires ou optionnels, leurs fraîcheurs, permissions, fallbacks, owners et limitations.

## 2. Problème utilisateur
Un enrichissement implicite peut remplacer silencieusement l’événement source, élargir les permissions ou rendre la logique inutilisable lorsque le contexte manque.

## 3. Objectifs
- définir enrichissements requis et optionnels
- définir owners, fraîcheur, permissions et restrictions
- définir comportement et fallback en absence
- voir dépendances et limitations
- préparer une demande à la capability propriétaire

## 4. Non-objectifs
- aucun moteur, langage, syntaxe vendor, API, protocole, parser, compilateur, AST, modèle ML, commande ou code
- aucune promotion, deployment, activation, deactivation, rollback, exception active ou mutation Signal/Alert
- aucune capability CAP-INV-5xx, Intelligence, Cloud/Mobile ou réécriture détaillée d’écran

## 5. Propriétaire
Investigate possède Enrichment Requirement et sa disposition humaine. Command conserve runtime Detection/Signal/Alert/Incident ; Settings les sources/parsers/schemas/health/retention ; Endpoint Agent ses capacités et résultats locaux ; Studio Tool/Tool Call/Workflow/Automation Run ; Govern l’autorité future ; Shared les mécanismes génériques.

## 6. Utilisateurs
Principal : **Detection Engineer**. Secondaires : Entity Analyst; Asset Owner; Platform Administrator.

## 7. Conditions d’entrée
Tenant, environnement, objective, scope, versions, sources, permissions, restrictions et return origin sont explicites. Une dépendance absente produit un état incomplet, partiel ou bloqué ; aucune donnée ou autorité n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | ---: | --- | --- |
| Detection Hypothesis and logic | CAP-INV-403/407/408 | behavior, entities, keys et context needs | oui | versions liées | requirements `draft` |
| Entity / Asset / Endpoint / Identity projections | Shared / Settings / Endpoint Agent | contextes disponibles et owners | non | freshness visible | fallback requis |
| Case / Incident context policy | Investigate / Command | contexte autorisé | non | permission courante | context omitted |
| Historical and source metadata | Shared / Settings | history, quality et restrictions | non | version/time bound | limitation visible |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Entity / relationship | Shared | candidate identity and context | lecture/lien sans merge |
| Endpoint/Fleet/Asset/Identity projections | Settings / Endpoint Agent | authorized context | lecture |
| Case / Incident | Investigate / Command | allowed contextual projection | lecture |
| Telemetry Event | Shared | source event preserved | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Enrichment Requirement | créer, modifier, versionner, superseder | Investigate concept | ne possède aucune source d’enrichissement |
| Fallback specification | définir | Investigate concept | absence explicitement gérée |
| Owner request context | préparer | destination owner | aucune donnée créée localement |

## 11. Fonctionnalités
- définir required/optional context
- documenter owners, freshness and permission
- définir absence behavior and fallback
- préserver source event alongside enrichment
- préparer requests to owners
- conserver versions, erreurs, partialité, restrictions et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | ---: | --- | --- | --- |
| Consulter/Comparer | Detection Engineer | Enrichment Requirement | 0 | lecture autorisée | projection sourcée | non |
| Exécuter traitement borné | Detection Engineer | Tool Call / Result | 1 | déclenchement explicite et permission | résultat attribué | policy |
| Créer/Modifier/Contester | Detection Engineer | Enrichment Requirement | 2 | mutation réversible | nouvelle version | OPEN-013 |
| Préparer handoff | Detection Engineer | candidate package | 2 | sources et limites visibles | package non effectif | destination |

Classes 3/4 exclues ; production et runtime appartiennent à 4B.3A.2/owners.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | ---: | ---: | ---: | ---: | --- |
| Construire Enrichment Requirement | oui | éditeur/contrôles explicables | oui | proposition | formulaire/table/revue |
| Valider ou comparer | oui | validateur/comparateur | oui | explication | diagnostics/diff |
| Expliquer erreurs | oui | catalogue | oui | résumé sourcé | erreurs brutes/checklist |
| Promouvoir/déployer/qualifier runtime | non | non | non | interdit | future phase/owner |

Toute sortie expose initiateur, producteur/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun choix silencieux.

## 14. États fonctionnels
`draft`, `ready`, `optional`, `required`, `unavailable`, `stale`, `permission-limited`, `fallback-defined`, `blocked`, `superseded`. États fonctionnels, pas machine objet finale.

## 15. États d’interface
Loading conserve context/version ; Empty distingue absence et interdiction ; Partial expose gaps ; Error conserve le valide ; Offline bloque les nouveaux runs ; Permission denied masque ; Stale distingue ancien/courant ; conflits fournissent diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Enrichment Requirements | functional requirement set | CAP-INV-406,411..414 | owner et fallback visibles |
| Owner request | request context | Shared/Settings/Endpoint/Command owner | aucune mutation source |
| Fallback behavior | functional specification | Validation / Tests | source event jamais remplacé |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-403/407 | identifier contexte | CAP-INV-409 | behavior, fields, entities and constraints | authoring |
| CAP-INV-409 | requirements disponibles | CAP-INV-406/411 | requirements, freshness, permissions and fallback | authoring |
| CAP-INV-409 | context missing | owner request preparation / CAP-INV-416 | gap, owner, impact and limitation | authoring |

Transitions conservent ownership, tenant/env, permissions, restrictions, versions, erreurs, provenance et return origin.

## 18. Dépendances
Shared Entity/Linking; Settings data/identity/fleet; Endpoint projections; Command/Case context; OPEN-013/015. Shared Jobs/Trace/Versioning/Linking/Search/Export/Reporting/Collaboration/Comparison/Recovery consommés sans redéfinition.

## 19. Source de vérité
Investigate est source de Enrichment Requirement ; tous les objets consommés restent chez leurs owners. Draft ≠ runtime Detection.

## 20. Provenance et audit
Source du besoin, Project, versions, sources/schemas/fields/mappings/logic, Tool/Calls/Runs, paramètres, datasets, résultats, erreurs, interruptions, auteurs, reviewers, dispositions, exports et correlation ID applicables à Enrichment Requirement.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | ---: | --- | --- | --- | --- | --- |
| Enrichment metadata read | cross-product context | 0 | sensitive values masked | possible | viewer/source owner | owners | Permissions |
| Enrichment requirement update | logic dependency | 2 | no source copy | OPEN-013 | author/reviewer | Investigate | Permissions |
| Context request prepare | access expansion | 2 | least data disclosure | step-up possible | requester/source owner | owners/Security | Permissions |

Matrice atomique, namespaces, RBAC/ABAC, step-up et SoD finaux reportés ; permissions production exclues.

## 22. Limites et erreurs
- Detection Engineering ne possède aucune source d’enrichissement.
- Enriched data ne remplace jamais silencieusement source event.
- Aucune permission, Entity merge ou source configuration n’est accordée localement.
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
### 1. Contexte absent
**Given** un enrichment requis indisponible  
**When** la logique est validée  
**Then** le fallback ou l’état blocked est visible et aucun contexte n’est inventé

### 2. Permission
**Given** un contexte Identity non autorisé  
**When** l’utilisateur inspecte les requirements  
**Then** l’existence peut être signalée, les valeurs restent masquées et une demande peut être préparée

### 3. Sans IA
**Given** aucun modèle  
**When** les enrichissements sont définis  
**Then** catalogues, owners, forms et checklists restent disponibles

## 26. Questions ouvertes
- OPEN-013 reste ouverte.
- OPEN-015 reste ouverte.
- Le moteur/langage Detection est une lacune future non couverte ; OPEN-005 reste forensic-only.
- Schémas, formats, permissions et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering module, Event Search/Hunt/Case/Evidence/technical handoffs, Command boundaries, Settings/Endpoint, Studio, Govern future review, Shared, Objects/Permissions/Screens/Journeys/Technique/4B.3A.2/Validation.
