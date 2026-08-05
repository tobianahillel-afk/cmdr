---
id: CAP-INV-412
title: Detection Test Scenario and Dataset Management
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
  - REQ-AI-010
open_decisions:
  - OPEN-013
  - OPEN-014
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-412 — Detection Test Scenario and Dataset Management
## 1. Définition
Créer, versionner et gérer scénarios et datasets fonctionnels distinguant données synthetic, historical, curated et imported, avec restrictions, sensibilité, variantes, cas positifs, négatifs, ambigus, incomplets, manquants et hors ordre.

## 2. Problème utilisateur
Un dataset non attribué ou présenté comme ground truth peut conduire à des résultats trompeurs, exposer des données sensibles et masquer les cas négatifs ou incomplets.

## 3. Objectifs
- créer scenario lié à un draft et définir behavior/preconditions
- identifier dataset source, type, restrictions and sensitive elements
- définir positive, negative, ambiguous and incomplete variants
- inclure missing/out-of-order data cases
- versionner, comparer, annoter, retirer and superseder

## 4. Non-objectifs
- aucun moteur, langage, syntaxe vendor, API, protocole, parser, compilateur, AST, modèle ML, commande ou code
- aucune promotion, deployment, activation, deactivation, rollback, exception active ou mutation Signal/Alert
- aucune capability CAP-INV-5xx, Intelligence, Cloud/Mobile ou réécriture détaillée d’écran

## 5. Propriétaire
Investigate possède Detection Test Scenario and Dataset et sa disposition humaine. Command conserve runtime Detection/Signal/Alert/Incident ; Settings les sources/parsers/schemas/health/retention ; Endpoint Agent ses capacités et résultats locaux ; Studio Tool/Tool Call/Workflow/Automation Run ; Govern l’autorité future ; Shared les mécanismes génériques.

## 6. Utilisateurs
Principal : **Detection Test Designer**. Secondaires : Detection Engineer; Privacy Reviewer; Dataset Curator.

## 7. Conditions d’entrée
Tenant, environnement, objective, scope, versions, sources, permissions, restrictions et return origin sont explicites. Une dépendance absente produit un état incomplet, partiel ou bloqué ; aucune donnée ou autorité n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | ---: | --- | --- |
| Detection Content Draft | CAP-INV-406..411 | behavior, logic, validation and limits | oui | selected version | scenario `incomplete` |
| Dataset or Artifact projection | Studio Dataset / Investigate Artifact / Shared storage | samples and provenance | oui for execution | version/integrity visible | scenario without runnable dataset |
| Privacy and permission context | Security / source owner | classification, masking and permitted use | oui | current | `restricted` |
| Coverage needs | CAP-INV-403/416 | positive, negative and gap variants | non | current | limited scenario set |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Detection Content Draft / Validation | Investigate concepts | target version and diagnostics | lecture |
| Artifact / Dataset projection | Investigate / Studio / Shared | version, source, integrity and restrictions | lecture/use according to permission |
| Telemetry Event samples | Shared | controlled examples | lecture |
| Tool / Automation Run | Studio | generation or execution provenance | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Detection Test Scenario | créer, modifier, exécuter-record, retirer, superseder | Investigate concept | scenario ≠ production telemetry |
| Detection Test Dataset relation | lier/versionner usage | Investigate concept / source owner | dataset ownership unchanged |
| Trace / Version event | émettre | Shared | source and restrictions retained |

## 11. Fonctionnalités
- define behavior, preconditions and variants
- classify dataset source type
- manage sensitive restrictions and masking
- define positive/negative/ambiguous/incomplete/missing/out-of-order cases
- version, compare, withdraw and supersede
- conserver versions, erreurs, partialité, restrictions et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | ---: | --- | --- | --- |
| Consulter/Comparer | Detection Test Designer | Detection Test Scenario and Dataset | 0 | lecture autorisée | projection sourcée | non |
| Exécuter traitement borné | Detection Test Designer | Tool Call / Result | 1 | déclenchement explicite et permission | résultat attribué | policy |
| Créer/Modifier/Contester | Detection Test Designer | Detection Test Scenario and Dataset | 2 | mutation réversible | nouvelle version | OPEN-013 |
| Préparer handoff | Detection Test Designer | candidate package | 2 | sources et limites visibles | package non effectif | destination |

Classes 3/4 exclues ; production et runtime appartiennent à 4B.3A.2/owners.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | ---: | ---: | ---: | ---: | --- |
| Construire Detection Test Scenario and Dataset | oui | éditeur/contrôles explicables | oui | proposition | formulaire/table/revue |
| Valider ou comparer | oui | validateur/comparateur | oui | explication | diagnostics/diff |
| Expliquer erreurs | oui | catalogue | oui | résumé sourcé | erreurs brutes/checklist |
| Promouvoir/déployer/qualifier runtime | non | non | non | interdit | future phase/owner |

Toute sortie expose initiateur, producteur/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun choix silencieux.

## 14. États fonctionnels
`draft`, `ready`, `restricted`, `incomplete`, `executed`, `partial`, `failed`, `superseded`, `withdrawn`. États fonctionnels, pas machine objet finale.

## 15. États d’interface
Loading conserve context/version ; Empty distingue absence et interdiction ; Partial expose gaps ; Error conserve le valide ; Offline bloque les nouveaux runs ; Permission denied masque ; Stale distingue ancien/courant ; conflits fournissent diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Detection Test Scenario | scenario | CAP-INV-413,411,415 | target version and variants visible |
| Dataset usage relation | versioned relation | Test execution / Audit | source ownership and restrictions preserved |
| Sensitive-data disposition | audit event | Security / Trace | read/use/copy boundaries visible |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-406/411 | préparer test | CAP-INV-412 | draft, behavior, logic, validation and limitations | authoring |
| CAP-INV-412 | définir attendu | CAP-INV-413 | scenario, variants, dataset source and confidence | tests |
| CAP-INV-412 | exécuter via validation/replay | CAP-INV-411/414 | scenario version, dataset ref, restrictions and expected outcome | tests |

Transitions conservent ownership, tenant/env, permissions, restrictions, versions, erreurs, provenance et return origin.

## 18. Dépendances
CAP-INV-105,406,411,416; Studio Datasets/Evaluation; Shared storage/versioning; OPEN-013/014/015. Shared Jobs/Trace/Versioning/Linking/Search/Export/Reporting/Collaboration/Comparison/Recovery consommés sans redéfinition.

## 19. Source de vérité
Investigate est source de Detection Test Scenario and Dataset ; tous les objets consommés restent chez leurs owners. Draft ≠ runtime Detection.

## 20. Provenance et audit
Source du besoin, Project, versions, sources/schemas/fields/mappings/logic, Tool/Calls/Runs, paramètres, datasets, résultats, erreurs, interruptions, auteurs, reviewers, dispositions, exports et correlation ID applicables à Detection Test Scenario and Dataset.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | ---: | --- | --- | --- | --- | --- |
| Test scenario read/create/update | test logic | 0/2 | sensitive values masked | OPEN-013 | author/reviewer | Investigate | Permissions |
| Test dataset read/use | sensitive or licensed data | 0/1 | minimization/redaction | step-up possible | requester/data owner | source owner | Permissions |
| Automated test generation request | invented assumptions | 1/2 | proposal and sources visible | possible | human reviewer | Studio/Investigate | Permissions |

Matrice atomique, namespaces, RBAC/ABAC, step-up et SoD finaux reportés ; permissions production exclues.

## 22. Limites et erreurs
- Test event ≠ production telemetry; synthetic fixture ≠ real incident.
- Dataset ≠ certain ground truth; labeled example ≠ certain ground truth.
- No unrestricted copy of sensitive data and no production deployment.
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
### 1. Restricted dataset
**Given** un dataset historique contenant des données sensibles  
**When** un scenario is opened  
**Then** metadata remains visible, content follows permission, and no unauthorized copy occurs

### 2. Ambiguous case
**Given** a variant with incomplete labels  
**When** the scenario is prepared  
**Then** confidence and ambiguity remain visible and expected outcome may be partial

### 3. Sans IA
**Given** aucun modèle  
**When** scenarios are created  
**Then** manual fixtures, controlled samples, forms and deterministic tooling remain available

## 26. Questions ouvertes
- OPEN-013 reste ouverte.
- OPEN-014 reste ouverte.
- OPEN-015 reste ouverte.
- Le moteur/langage Detection est une lacune future non couverte ; OPEN-005 reste forensic-only.
- Schémas, formats, permissions et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering module, Event Search/Hunt/Case/Evidence/technical handoffs, Command boundaries, Settings/Endpoint, Studio, Govern future review, Shared, Objects/Permissions/Screens/Journeys/Technique/4B.3A.2/Validation.
