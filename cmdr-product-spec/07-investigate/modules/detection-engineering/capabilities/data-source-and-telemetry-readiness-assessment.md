---
id: CAP-INV-404
title: Data Source and Telemetry Readiness Assessment
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
  - REQ-PROD-055
open_decisions:
  - OPEN-008
  - OPEN-013
source-of-truth: canonical
---
# CAP-INV-404 — Data Source and Telemetry Readiness Assessment
## 1. Définition
Évaluer fonctionnellement si les sources nécessaires sont disponibles, fraîches, suffisamment retenues et couvertes pour l’objectif de détection, sans administrer les sources ni confondre disponibilité et suffisance.

## 2. Problème utilisateur
Une source configurée ou active peut être stale, partielle, insuffisamment retenue ou incapable de couvrir les plateformes et comportements visés.

## 3. Objectifs
- identifier sources nécessaires, disponibles et manquantes
- voir tenants, environnements, producers, fraîcheur, rétention, gaps, interruptions et limitations
- comparer la readiness entre environnements et plateformes
- classer la readiness avec incertitude
- créer un Detection Gap ou préparer une demande Settings/Collection

## 4. Non-objectifs
- aucun moteur, langage, syntaxe vendor, API, protocole, parser, compilateur, AST, modèle ML, commande ou code
- aucune promotion, deployment, activation, deactivation, rollback, exception active ou mutation Signal/Alert
- aucune capability CAP-INV-5xx, Intelligence, Cloud/Mobile ou réécriture détaillée d’écran

## 5. Propriétaire
Investigate possède Telemetry Readiness Assessment et sa disposition humaine. Command conserve runtime Detection/Signal/Alert/Incident ; Settings les sources/parsers/schemas/health/retention ; Endpoint Agent ses capacités et résultats locaux ; Studio Tool/Tool Call/Workflow/Automation Run ; Govern l’autorité future ; Shared les mécanismes génériques.

## 6. Utilisateurs
Principal : **Detection Engineer**. Secondaires : Data Source Owner; Platform Administrator; Detection Reviewer.

## 7. Conditions d’entrée
Tenant, environnement, objective, scope, versions, sources, permissions, restrictions et return origin sont explicites. Une dépendance absente produit un état incomplet, partiel ou bloqué ; aucune donnée ou autorité n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | ---: | --- | --- |
| Detection Hypothesis and source requirements | CAP-INV-403 | comportement, plateformes, période et sources attendues | oui | version courante | assessment `unknown` |
| Data Source projections | Platform Settings | configuration, scope, fraîcheur, health et restrictions | oui | snapshot horodaté | `unavailable` |
| Endpoint telemetry capability projections | Endpoint Agent / Settings Fleet | capacités déclarées et limitations | non | état courant | `unsupported` ou `coverage-limited` |
| Retention and quality projections | Settings / Shared Data Quality | rétention, gaps, parser failures et complétude | oui | période évaluée | `retention-insufficient` ou `degraded` |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Data Source / Integration / Parser | Platform Settings | scope, état administratif, version et health | lecture uniquement |
| Telemetry Event / quality projection | Shared | fraîcheur, complétude et source | lecture/agrégation |
| Endpoint Agent/Fleet capability | Endpoint Agent / Settings | plateformes et télémétrie déclarée | lecture |
| Detection Hypothesis / Project | Investigate | objectif, période et consommateurs | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Telemetry Readiness Assessment | créer, annoter, comparer, contester, superseder | Investigate concept | ne modifie aucune source |
| Detection Gap | créer ou lier | Investigate concept | impact et owner probable visibles |
| Settings or Collection request context | préparer | owner destination | aucune configuration locale |

## 11. Fonctionnalités
- cartographier sources requises et disponibles
- examiner fraîcheur, rétention, gaps, volumes conceptuels et limitations
- comparer tenants, environnements et plateformes
- classifier readiness
- préparer gap ou demande propriétaire
- conserver versions, erreurs, partialité, restrictions et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | ---: | --- | --- | --- |
| Consulter/Comparer | Detection Engineer | Telemetry Readiness Assessment | 0 | lecture autorisée | projection sourcée | non |
| Exécuter traitement borné | Detection Engineer | Tool Call / Result | 1 | déclenchement explicite et permission | résultat attribué | policy |
| Créer/Modifier/Contester | Detection Engineer | Telemetry Readiness Assessment | 2 | mutation réversible | nouvelle version | OPEN-013 |
| Préparer handoff | Detection Engineer | candidate package | 2 | sources et limites visibles | package non effectif | destination |

Classes 3/4 exclues ; production et runtime appartiennent à 4B.3A.2/owners.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | ---: | ---: | ---: | ---: | --- |
| Construire Telemetry Readiness Assessment | oui | éditeur/contrôles explicables | oui | proposition | formulaire/table/revue |
| Valider ou comparer | oui | validateur/comparateur | oui | explication | diagnostics/diff |
| Expliquer erreurs | oui | catalogue | oui | résumé sourcé | erreurs brutes/checklist |
| Promouvoir/déployer/qualifier runtime | non | non | non | interdit | future phase/owner |

Toute sortie expose initiateur, producteur/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun choix silencieux.

## 14. États fonctionnels
`unknown`, `unavailable`, `partially-available`, `available`, `stale`, `degraded`, `coverage-limited`, `retention-insufficient`, `unsupported`, `policy-blocked`, `disputed`. États fonctionnels, pas machine objet finale.

## 15. États d’interface
Loading conserve context/version ; Empty distingue absence et interdiction ; Partial expose gaps ; Error conserve le valide ; Offline bloque les nouveaux runs ; Permission denied masque ; Stale distingue ancien/courant ; conflits fournissent diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Telemetry Readiness Assessment | assessment | CAP-INV-405,406,416 | availability et sufficiency séparées |
| Data Source Requirement | functional requirement | Platform Settings request preparation | aucune configuration exécutée |
| Detection Gap candidate | gap | CAP-INV-416 | impact, source et environnement visibles |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-402/403 | évaluer les données | CAP-INV-404 | sources requises, plateformes, environnements, période | project |
| CAP-INV-404 | sources exploitables | CAP-INV-405 | sources, versions, quality et restrictions | assessment |
| CAP-INV-404 | source ou rétention insuffisante | CAP-INV-416 / Settings request | gap, impact et contexte administratif | assessment |

Transitions conservent ownership, tenant/env, permissions, restrictions, versions, erreurs, provenance et return origin.

## 18. Dépendances
Platform Settings Data Source/Parser/Health/Retention; Endpoint Agent telemetry projections; Shared Data Quality; OPEN-008/013. Shared Jobs/Trace/Versioning/Linking/Search/Export/Reporting/Collaboration/Comparison/Recovery consommés sans redéfinition.

## 19. Source de vérité
Investigate est source de Telemetry Readiness Assessment ; tous les objets consommés restent chez leurs owners. Draft ≠ runtime Detection.

## 20. Provenance et audit
Source du besoin, Project, versions, sources/schemas/fields/mappings/logic, Tool/Calls/Runs, paramètres, datasets, résultats, erreurs, interruptions, auteurs, reviewers, dispositions, exports et correlation ID applicables à Telemetry Readiness Assessment.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | ---: | --- | --- | --- | --- | --- |
| Data Source readiness read | métadonnées d’infrastructure | 0 | secrets et valeurs sensibles masqués | possible | viewer/admin | Settings projection | Permissions |
| Readiness assessment create/update | impact de couverture | 2 | tenant/env bornés | OPEN-013 | author/source owner review | Investigate | Permissions |
| Cross-tenant readiness compare | isolement tenant | 1/2 | agrégation et masquage | step-up probable | reviewer indépendant | Security/Settings | Permissions |

Matrice atomique, namespaces, RBAC/ABAC, step-up et SoD finaux reportés ; permissions production exclues.

## 22. Limites et erreurs
- Data source configured ≠ healthy.
- Healthy ≠ telemetry complete; telemetry present ≠ sufficient.
- Aucune source, parser, retention ou Fleet policy n’est modifiée.
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
### 1. Source insuffisante
**Given** une source active avec gaps et rétention trop courte  
**When** la readiness est évaluée  
**Then** elle n’est pas pleinement ready, les gaps/rétention sont visibles et un Detection Gap peut être créé

### 2. Source indisponible
**Given** une projection Settings inaccessible  
**When** l’assessment est ouvert  
**Then** l’état reste `unknown` ou `policy-blocked`, aucune santé n’est inventée

### 3. Sans IA
**Given** aucun modèle  
**When** la readiness est examinée  
**Then** tables, health projections, règles explicables et revue humaine fonctionnent

## 26. Questions ouvertes
- OPEN-008 reste ouverte.
- OPEN-013 reste ouverte.
- Le moteur/langage Detection est une lacune future non couverte ; OPEN-005 reste forensic-only.
- Schémas, formats, permissions et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering module, Event Search/Hunt/Case/Evidence/technical handoffs, Command boundaries, Settings/Endpoint, Studio, Govern future review, Shared, Objects/Permissions/Screens/Journeys/Technique/4B.3A.2/Validation.
