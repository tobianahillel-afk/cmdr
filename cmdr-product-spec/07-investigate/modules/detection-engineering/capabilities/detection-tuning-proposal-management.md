---
id: CAP-INV-429
title: Detection Tuning Proposal Management
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
  - REQ-PROD-055
  - REQ-AI-002
  - REQ-SEC-001
  - REQ-SEC-002
  - REQ-UX-010
open_decisions:
  - OPEN-013
  - OPEN-015
  - OPEN-017
source-of-truth: canonical
---
# CAP-INV-429 — Detection Tuning Proposal Management

## 1. Définition
Préparer et versionner une proposition de tuning fondée sur les observations runtime, sans modifier silencieusement la version active et en retournant vers un nouveau Detection Content Draft.

## 2. Problème utilisateur
Une recommandation de seuil, exclusion ou mapping peut être appliquée directement en production sans tests, rollback ou revue.

## 3. Objectifs
- proposer changes to conditions, exclusions, thresholds, windows, sequences, enrichments, mappings and proposed severity
- proposer target population, version, data source and dependency changes
- documenter observed problem, sources, expected impact, risks and candidate FP/FN
- documenter environments, tests required, rollback, reviewer and justification
- créer ou alimenter explicitement un nouveau Draft with lineage

## 4. Non-objectifs
Aucune API, protocole, syntaxe vendor, moteur, langage, compilateur, parser, AST, format technique, commande, code, modèle ML, mutation runtime directe, capability CAP-INV-5xx, Threat Intelligence, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède **Tuning Proposal** comme concept fonctionnel. Command garde Detection/Signal/Alert/Incident ; Settings les runtimes, targets, environments et health ; Endpoint ses capacités et exécutions locales ; Govern Action Request/Decision/Approval/Response Run/Result ; Studio Tools et Automation Runs ; Shared les mécanismes transverses.

## 6. Utilisateurs
Principal : **Detection Engineer**. Secondaires : Detection Engineer, Detection Reviewer, Detection Owner, Incident Commander, Platform Operator et Approver autorisés.

## 7. Conditions d’entrée
Tenant, environnement, target, version, période, owner, permissions, restrictions et return origin sont explicites. Les preuves 4B.3A.1 sont référencées sans duplication. Toute absence produit un état incomplete, partial, blocked ou unknown.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | --- | --- | --- |
| Runtime Quality / Production Review | CAP-INV-427/428 | problem, evidence, candidate FP/FN and version | oui | selected assessment | proposal incomplete |
| Health/drift/performance context | CAP-INV-426/431/432 | technical limits and dependencies | non | same period/version | risk limited |
| Active and source Draft versions | Command / CAP-INV-406 | runtime version and authoring lineage | oui | resolvable versions | no safe lineage |
| Test/rollback requirements | CAP-INV-412..417/433 | tests to rerun and recovery context | oui | current | not reviewable |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Runtime Detection / Version | Command | active source version projection | lecture |
| Runtime Quality / Production Match Review | Investigate | evidence and dispositions | lecture |
| Detection Content Draft / logic components | Investigate | lineage and editable new version | lecture/reference |
| Data/schema/health projections | Settings / Shared | dependencies and limits | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Tuning Proposal | créer, modifier, revoir, retirer, superseder | Investigate concept | proposal ≠ active tuning |
| New Detection Content Draft | créer/alimenter explicitement | Investigate via CAP-INV-406 | new version with lineage |
| Change/test context | préparer | CAP-INV-411..420 | no runtime effect |

## 11. Fonctionnalités
- proposer changes to conditions, exclusions, thresholds, windows, sequences, enrichments, mappings and proposed severity
- proposer target population, version, data source and dependency changes
- documenter observed problem, sources, expected impact, risks and candidate FP/FN
- documenter environments, tests required, rollback, reviewer and justification
- créer ou alimenter explicitement un nouveau Draft with lineage
- conserver sources, versions, erreurs, partialité, attribution et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | --- | --- | --- | --- |
| consulter/comparer | Detection Engineer | Tuning Proposal | 0 | lecture autorisée | projection sourcée | non |
| assessment bornée | Detection Engineer | Tool Call / assessment | 1 | lancement explicite | résultat attribué | policy |
| créer/modifier/contester | Detection Engineer | Tuning Proposal | 2 | mutation réversible | nouvelle version | OPEN-013 |
| préparer demande | Detection Engineer | Action Request context | 2 | risque, cible et rollback visibles | package non effectif | Govern |
| changement réel | Govern/runtime owner | runtime target | 3 | Decision/Approval | Result projeté | owner |
| détruire historique | aucun rôle local | provenance | 4 | interdit | refus | strict |

Investigate n’exécute jamais les classes 3/4.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | --- | --- | --- | --- | --- |
| construire/revoir Tuning Proposal | oui | formulaires et contrôles | oui | proposition | matrices et revue humaine |
| comparer/valider | oui | diff et règles explicables | oui | explication | diagnostics et checklist |
| résumer risques/erreurs | oui | agrégations | oui | résumé sourcé | tables et timeline |
| approuver/exécuter | non localement | non | non | interdit | Govern/runtime owner |

Initiateur, agent/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner et disposition humaine sont visibles. Aucun choix silencieux.

## 14. États fonctionnels
`draft`, `incomplete`, `under-review`, `accepted-for-draft`, `rejected`, `returned-for-evidence`, `superseded`, `withdrawn`. Projections fonctionnelles, pas machine d’état objet définitive.

## 15. États d’interface
Loading conserve le contexte ; Empty distingue absence et interdiction ; Partial détaille les targets ; Error conserve le valide ; Offline bloque les mutations ; Permission denied masque ; Stale distingue ancien/courant ; les conflits offrent diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Tuning Proposal | proposal | CAP-INV-406/418 | problem, changes, risks and evidence |
| New Draft context | authoring handoff | CAP-INV-406 | source active version and lineage |
| Required-test package | test request context | CAP-INV-411..417 | no promotion without evidence |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-427/428/432 | tuning need identified | CAP-INV-429 | runtime problem, version, targets, evidence and limits | source assessment |
| CAP-INV-429 | accepted for draft | CAP-INV-406 | proposed logic/metadata changes, tests and lineage | Tuning Proposal |
| CAP-INV-429 | future change ready | CAP-INV-418/420 | new candidate evidence, risks, rollback and review | Tuning Proposal |

Chaque transition conserve ownership, tenant/env, target, version, permissions, restrictions, erreurs, autorité, provenance et return origin.

## 18. Dépendances
CAP-INV-406..418,426..428,431..433; Command runtime version; Settings/Shared projections. OPEN-017 couvre la future stratégie runtime/langage/portabilité sans option sélectionnée. Shared Jobs, Notifications, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Collaboration, Comparison, Inspector, Audit Hooks et Recovery sont consommés sans redéfinition.

## 19. Source de vérité
Investigate possède l’assessment/proposition locale. Les objets Command, Settings, Endpoint, Govern, Studio et Shared restent canoniques chez leurs owners. Une projection ne remplace jamais sa source.

## 20. Provenance et audit
Conserver besoin, Project, Hypothesis, Drafts/Versions, Review Package, Release Candidate, reviews, readiness, targets, plans, Action Requests, Decisions, Approvals, Runs, Results, runtime observations, health, Signals/Alerts/Incidents, assessments, propositions, erreurs, auteurs, timestamps, Tools/Runs, paramètres et dispositions. Aucune trace supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Tuning Proposal create/update | future active behavior | 2 | sensitive examples masked | OPEN-013 | author/reviewer | Investigate | Permissions |
| Automated tuning proposal | overfitting/unsafe exclusion | 1/2 | sources, uncertainty and diff visible | possible | human acceptance | Studio/Investigate | Permissions |
| Active tuning application | production change | 3 | not available locally | Govern required | requester/approver/operator | Govern/Settings/runtime | Permissions |

Permissions atomiques, namespaces, RBAC/ABAC, step-up et séparation finale restent reportés.

## 22. Limites et erreurs
- Tuning Proposal ≠ active tuning.
- No active version is modified silently.
- Observed noise does not prove all matches are false positives.
- Tests, review and rollback remain required.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, target offline et version superseded restent visibles.
- Result, match, health, disposition Command ou sortie IA ne vaut pas conclusion universelle.

## 23. Métriques conceptuelles
- proposals by source/problem/state
- accepted/rejected/returned proportions
- required tests and risk categories
- silent active tuning — target zero
- provenance et dispositions humaines complètes
- exécutions silencieuses, auto-approbations et suppressions de trace : cible zéro

Aucun seuil technique universel n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; documentation only. Aucun `validated`, `implemented`, `deployed`, `active`, `native` ou `integrated`.

## 25. Critères d’acceptation
### 1. Healthy noisy runtime
**Given** a healthy runtime and noise feedback  
**When** tuning is proposed  
**Then** the active version remains unchanged and a new Draft with tests can be created

### 2. Insufficient evidence
**Given** a proposed exclusion with weak ground truth  
**When** review is performed  
**Then** proposal may be returned-for-evidence and no active suppression is created

### 3. Sans IA
**Given** aucun modèle  
**When** tuning is prepared  
**Then** forms, diffs, source evidence and human review suffice

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-015 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-017 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-005 reste forensic-only.
- Schémas, formats, cardinalités, permissions atomiques, contrats techniques et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering lifecycle, Command, Settings, Endpoint, Govern, Studio, Shared, Objects, Permissions, Screens, Journeys, Technique, validation et futur handoff 4B.3B non canonique.
