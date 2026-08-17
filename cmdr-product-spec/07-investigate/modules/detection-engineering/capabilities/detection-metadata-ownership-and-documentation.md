---
id: CAP-INV-410
title: Detection Metadata, Ownership and Documentation
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
# CAP-INV-410 — Detection Metadata, Ownership and Documentation
## 1. Définition
Documenter nom, objectif, owner, contributeurs, statut, version, sources, hypothesis, fields, logique, exclusions, enrichissements, limites, sévérité et priorité proposées, consommateurs, mappings, tests, risques, dépendances, historique, provenance et politique de revue future.

## 2. Problème utilisateur
Un draft techniquement lisible mais sans owner, objectif, limites, risques FP/FN ou provenance ne peut pas être revu de façon responsable.

## 3. Objectifs
- documenter identity, objective, owner, contributors, status and version
- relier need source, hypothesis, data sources, fields and logic
- documenter exclusions, enrichments, limitations, severity/priority proposals and consumers
- documenter mappings, tests, expected outcomes, FP/FN risks and dependencies
- conserver history, provenance, future review policy and handoff

## 4. Non-objectifs
- aucun moteur, langage, syntaxe vendor, API, protocole, parser, compilateur, AST, modèle ML, commande ou code
- aucune promotion, deployment, activation, deactivation, rollback, exception active ou mutation Signal/Alert
- aucune capability CAP-INV-5xx, Intelligence, Cloud/Mobile ou réécriture détaillée d’écran

## 5. Propriétaire
Investigate possède Detection Metadata Record et sa disposition humaine. Command conserve runtime Detection/Signal/Alert/Incident ; Settings les sources/parsers/schemas/health/retention ; Endpoint Agent ses capacités et résultats locaux ; Studio Tool/Tool Call/Workflow/Automation Run ; Govern l’autorité future ; Shared les mécanismes génériques.

## 6. Utilisateurs
Principal : **Detection Content Owner**. Secondaires : Detection Engineer; Reviewer; Investigation Lead.

## 7. Conditions d’entrée
Tenant, environnement, objective, scope, versions, sources, permissions, restrictions et return origin sont explicites. Une dépendance absente produit un état incomplet, partiel ou bloqué ; aucune donnée ou autorité n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | ---: | --- | --- |
| Detection Content Draft | CAP-INV-406..409 | version, sources, logic and enrichments | oui | version candidate | documentation `incomplete` |
| Project and Hypothesis | CAP-INV-402/403 | owner, objective, need source and limits | oui | current linked versions | missing ownership |
| Test/validation/coverage references | CAP-INV-411..416 | scenarios, outcomes, results and gaps | non | latest selected | pre-test documentation |
| Taxonomy and mapping references | Shared / product catalogues | candidate mappings | non | version identified | no mapping claim |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Detection Content Draft / Version | Investigate concepts | authoring details | lecture |
| Project / Detection Hypothesis | Investigate concepts | objective and ownership | lecture |
| Test / Validation / Replay / Coverage results | Investigate concepts | evidence for documentation | lecture |
| Taxonomy projections | Shared / source owner | mapping and version | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Detection Metadata Record | créer, modifier, versionner, superseder | Investigate concept | metadata proposed, not runtime parameters |
| Documentation history | append/update through versions | Investigate / Shared Versioning | author and diff visible |
| Review policy reference | documenter | future owner | aucune Approval créée |

## 11. Fonctionnalités
- documenter identity, objective, owner and version
- documenter sources, fields, logic, exclusions and enrichments
- documenter proposed severity/priority and consumers
- documenter tests, risks, gaps and dependencies
- maintenir history, provenance and future review policy
- conserver versions, erreurs, partialité, restrictions et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | ---: | --- | --- | --- |
| Consulter/Comparer | Detection Content Owner | Detection Metadata Record | 0 | lecture autorisée | projection sourcée | non |
| Exécuter traitement borné | Detection Content Owner | Tool Call / Result | 1 | déclenchement explicite et permission | résultat attribué | policy |
| Créer/Modifier/Contester | Detection Content Owner | Detection Metadata Record | 2 | mutation réversible | nouvelle version | OPEN-013 |
| Préparer handoff | Detection Content Owner | candidate package | 2 | sources et limites visibles | package non effectif | destination |

Classes 3/4 exclues ; production et runtime appartiennent à 4B.3A.2/owners.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | ---: | ---: | ---: | ---: | --- |
| Construire Detection Metadata Record | oui | éditeur/contrôles explicables | oui | proposition | formulaire/table/revue |
| Valider ou comparer | oui | validateur/comparateur | oui | explication | diagnostics/diff |
| Expliquer erreurs | oui | catalogue | oui | résumé sourcé | erreurs brutes/checklist |
| Promouvoir/déployer/qualifier runtime | non | non | non | interdit | future phase/owner |

Toute sortie expose initiateur, producteur/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun choix silencieux.

## 14. États fonctionnels
`draft`, `incomplete`, `reviewable`, `valid-with-warnings`, `documented`, `stale`, `disputed`, `superseded`. États fonctionnels, pas machine objet finale.

## 15. États d’interface
Loading conserve context/version ; Empty distingue absence et interdiction ; Partial expose gaps ; Error conserve le valide ; Offline bloque les nouveaux runs ; Permission denied masque ; Stale distingue ancien/courant ; conflits fournissent diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Detection Metadata Record | documentation | CAP-INV-411,417 | version and owner explicit |
| Documentation completeness result | assessment | Project / Reviewer | missing fields visible |
| Review context | metadata projection | future 4B.3A.2 | no runtime activation |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-406..409 | documenter le draft | CAP-INV-410 | version, objective, sources, fields, logic and limitations | authoring |
| CAP-INV-410 | valider documentation | CAP-INV-411 | metadata record and dependencies | documentation |
| CAP-INV-410 | préparer handoff | CAP-INV-417 | owner, risks, mappings, history and unresolved items | documentation |

Transitions conservent ownership, tenant/env, permissions, restrictions, versions, erreurs, provenance et return origin.

## 18. Dépendances
CAP-INV-402..409,411..416; Shared Versioning/Reporting/Taxonomy projections; OPEN-013/015. Shared Jobs/Trace/Versioning/Linking/Search/Export/Reporting/Collaboration/Comparison/Recovery consommés sans redéfinition.

## 19. Source de vérité
Investigate est source de Detection Metadata Record ; tous les objets consommés restent chez leurs owners. Draft ≠ runtime Detection.

## 20. Provenance et audit
Source du besoin, Project, versions, sources/schemas/fields/mappings/logic, Tool/Calls/Runs, paramètres, datasets, résultats, erreurs, interruptions, auteurs, reviewers, dispositions, exports et correlation ID applicables à Detection Metadata Record.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | ---: | --- | --- | --- | --- | --- |
| Detection metadata read | documentation sensible | 0 | restricted references masked | possible | viewer/reviewer | Investigate | Permissions |
| Detection metadata update | ownership and risk statements | 2 | no runtime secret | OPEN-013 | content owner/reviewer | Investigate | Permissions |
| Taxonomy mapping propose | overclaim coverage | 2 | candidate status visible | possible | author/taxonomy reviewer | Investigate/owner | Permissions |

Matrice atomique, namespaces, RBAC/ABAC, step-up et SoD finaux reportés ; permissions production exclues.

## 22. Limites et erreurs
- Proposed severity ≠ operational severity.
- Metadata proposed ≠ active runtime parameters.
- Technique mapping ≠ proof of attack coverage; local documentation ≠ canonical Report.
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
### 1. Owner missing
**Given** un draft sans owner  
**When** la documentation est validée  
**Then** le résultat reste incomplete et aucun handoff review-ready n’est attribué

### 2. Severity proposal
**Given** une sévérité proposée  
**When** Command consulte le contexte  
**Then** aucune sévérité opérationnelle active n’est créée

### 3. Sans IA
**Given** aucun modèle  
**When** la documentation est produite  
**Then** forms, templates, checklists and manual review suffice

## 26. Questions ouvertes
- OPEN-013 reste ouverte.
- OPEN-015 reste ouverte.
- Le moteur/langage Detection est une lacune future non couverte ; OPEN-005 reste forensic-only.
- Schémas, formats, permissions et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering module, Event Search/Hunt/Case/Evidence/technical handoffs, Command boundaries, Settings/Endpoint, Studio, Govern future review, Shared, Objects/Permissions/Screens/Journeys/Technique/4B.3A.2/Validation.
