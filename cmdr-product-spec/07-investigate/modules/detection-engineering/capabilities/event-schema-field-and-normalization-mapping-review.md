---
id: CAP-INV-405
title: Event Schema, Field and Normalization Mapping Review
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
# CAP-INV-405 — Event Schema, Field and Normalization Mapping Review
## 1. Définition
Examiner schémas, champs source et normalisés, types, transformations, exemples, valeurs manquantes, ambiguïtés, mappings, contradictions et versions pour sélectionner les dépendances fonctionnelles d’un Detection Content Draft.

## 2. Problème utilisateur
La présence d’un champ ou d’un mapping ne garantit ni fiabilité de la valeur, ni équivalence sémantique, ni conservation complète de la source.

## 3. Objectifs
- voir schémas, champs source/normalisés, types, exemples autorisés et versions
- examiner valeurs manquantes, ambiguïtés, transformations, mappings et contradictions
- comparer plusieurs sources et versions
- proposer ou contester un mapping fonctionnel
- sélectionner les champs requis et préparer une demande Settings

## 4. Non-objectifs
- aucun moteur, langage, syntaxe vendor, API, protocole, parser, compilateur, AST, modèle ML, commande ou code
- aucune promotion, deployment, activation, deactivation, rollback, exception active ou mutation Signal/Alert
- aucune capability CAP-INV-5xx, Intelligence, Cloud/Mobile ou réécriture détaillée d’écran

## 5. Propriétaire
Investigate possède Field Mapping Review et sa disposition humaine. Command conserve runtime Detection/Signal/Alert/Incident ; Settings les sources/parsers/schemas/health/retention ; Endpoint Agent ses capacités et résultats locaux ; Studio Tool/Tool Call/Workflow/Automation Run ; Govern l’autorité future ; Shared les mécanismes génériques.

## 6. Utilisateurs
Principal : **Detection Engineer**. Secondaires : Data Engineer; Parser Owner; Detection Reviewer.

## 7. Conditions d’entrée
Tenant, environnement, objective, scope, versions, sources, permissions, restrictions et return origin sont explicites. Une dépendance absente produit un état incomplet, partiel ou bloqué ; aucune donnée ou autorité n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | ---: | --- | --- |
| Ready or partial data sources | CAP-INV-404 | sources, versions, gaps et restrictions | oui | assessment courant | review `blocked` ou `partial` |
| Schema and field catalog projections | Settings / Shared Normalization | source fields, normalized fields, types et versions | oui | version identifiée | `schema-unavailable` |
| Authorized examples | Telemetry Event projections | valeurs masquées ou échantillons | non | période/source indiquées | revue sans exemple |
| Detection Hypothesis | CAP-INV-403 | sémantique recherchée | oui | version liée | mapping sans objectif interdit |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Parser / Data Source | Platform Settings | versions et transformations déclarées | lecture |
| Schema / Field metadata | Settings / Shared | types, mappings, qualité et raw reference | lecture |
| Telemetry Event samples | Shared | exemples autorisés et provenance | lecture selon permission |
| Detection Hypothesis / Readiness Assessment | Investigate concepts | besoin et limites | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Field Mapping Review | créer, annoter, contester, superseder | Investigate concept | mapping reste fonctionnel et réversible |
| Field Reference / Mapping candidate | proposer ou sélectionner | Investigate concept / Settings destination | aucun parser modifié |
| Settings request context | préparer | Platform Settings | différences et exemples sourcés |

## 11. Fonctionnalités
- inspecter schémas et champs
- comparer source et normalisé
- examiner transformations, valeurs manquantes et contradictions
- sélectionner ou contester un mapping fonctionnel
- documenter dépendances et demande Settings
- conserver versions, erreurs, partialité, restrictions et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | ---: | --- | --- | --- |
| Consulter/Comparer | Detection Engineer | Field Mapping Review | 0 | lecture autorisée | projection sourcée | non |
| Exécuter traitement borné | Detection Engineer | Tool Call / Result | 1 | déclenchement explicite et permission | résultat attribué | policy |
| Créer/Modifier/Contester | Detection Engineer | Field Mapping Review | 2 | mutation réversible | nouvelle version | OPEN-013 |
| Préparer handoff | Detection Engineer | candidate package | 2 | sources et limites visibles | package non effectif | destination |

Classes 3/4 exclues ; production et runtime appartiennent à 4B.3A.2/owners.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | ---: | ---: | ---: | ---: | --- |
| Construire Field Mapping Review | oui | éditeur/contrôles explicables | oui | proposition | formulaire/table/revue |
| Valider ou comparer | oui | validateur/comparateur | oui | explication | diagnostics/diff |
| Expliquer erreurs | oui | catalogue | oui | résumé sourcé | erreurs brutes/checklist |
| Promouvoir/déployer/qualifier runtime | non | non | non | interdit | future phase/owner |

Toute sortie expose initiateur, producteur/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun choix silencieux.

## 14. États fonctionnels
`not-reviewed`, `reviewing`, `selected`, `ambiguous`, `conflicting`, `partial`, `schema-unavailable`, `field-missing`, `disputed`, `superseded`. États fonctionnels, pas machine objet finale.

## 15. États d’interface
Loading conserve context/version ; Empty distingue absence et interdiction ; Partial expose gaps ; Error conserve le valide ; Offline bloque les nouveaux runs ; Permission denied masque ; Stale distingue ancien/courant ; conflits fournissent diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Selected Field References | functional references | CAP-INV-406..411 | versions et limites visibles |
| Field Mapping Review | review result | Project / Validation | équivalence non affirmée |
| Settings mapping request | request context | Platform Settings | aucun parser ou schema modifié |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-404 | sources identifiées | CAP-INV-405 | sources, versions, gaps, freshness et restrictions | assessment |
| CAP-INV-405 | champs sélectionnés | CAP-INV-406 | field refs, mappings, ambiguïtés et versions | project |
| CAP-INV-405 | mapping incomplet ou contradictoire | Settings request preparation | source fields, candidates, exemples et impact | review |

Transitions conservent ownership, tenant/env, permissions, restrictions, versions, erreurs, provenance et return origin.

## 18. Dépendances
CAP-INV-003/004/404; Settings Data Source/Parser; Shared Telemetry Normalization/Data Quality; OPEN-008/013. Shared Jobs/Trace/Versioning/Linking/Search/Export/Reporting/Collaboration/Comparison/Recovery consommés sans redéfinition.

## 19. Source de vérité
Investigate est source de Field Mapping Review ; tous les objets consommés restent chez leurs owners. Draft ≠ runtime Detection.

## 20. Provenance et audit
Source du besoin, Project, versions, sources/schemas/fields/mappings/logic, Tool/Calls/Runs, paramètres, datasets, résultats, erreurs, interruptions, auteurs, reviewers, dispositions, exports et correlation ID applicables à Field Mapping Review.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | ---: | --- | --- | --- | --- | --- |
| Schema and field catalog read | structure et exemples sensibles | 0 | valeurs masquées | possible | viewer/parser owner | Settings/Shared | Permissions |
| Field mapping propose/review | impact sémantique | 2 | raw restreint | OPEN-013 | author/parser owner | Investigate/Settings | Permissions |
| Sensitive example read | données de production | 0/1 | redaction et minimisation | step-up probable | requester/reviewer | Security | Permissions |

Matrice atomique, namespaces, RBAC/ABAC, step-up et SoD finaux reportés ; permissions production exclues.

## 22. Limites et erreurs
- Field present ≠ field reliable.
- Mapping ≠ certaine semantic equivalence; normalized field ≠ source field.
- Aucun parser, schéma physique, AST ou format de champ final n’est défini.
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
### 1. Mapping ambigu
**Given** un champ source et deux mappings aux sémantiques différentes  
**When** l’analyste sélectionne un mapping  
**Then** les candidats et différences restent visibles, la sélection est attribuée et contestable

### 2. Exemple restreint
**Given** un exemple contient des données non autorisées  
**When** la revue est ouverte  
**Then** les metadata restent disponibles, la valeur est masquée et aucun accès n’est élargi

### 3. Sans IA
**Given** aucun modèle  
**When** le mapping est revu  
**Then** catalogues, diff, exemples autorisés et revue humaine restent disponibles

## 26. Questions ouvertes
- OPEN-008 reste ouverte.
- OPEN-013 reste ouverte.
- Le moteur/langage Detection est une lacune future non couverte ; OPEN-005 reste forensic-only.
- Schémas, formats, permissions et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering module, Event Search/Hunt/Case/Evidence/technical handoffs, Command boundaries, Settings/Endpoint, Studio, Govern future review, Shared, Objects/Permissions/Screens/Journeys/Technique/4B.3A.2/Validation.
