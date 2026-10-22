---
id: CAP-INV-434
title: Detection Deactivation, Retirement and Replacement Management
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
  - OPEN-008
  - OPEN-013
  - OPEN-015
  - OPEN-017
source-of-truth: canonical
---
# CAP-INV-434 — Detection Deactivation, Retirement and Replacement Management

## 1. Définition
Préparer et suivre une désactivation, un retirement ou un remplacement gouverné tout en préservant versions, consommateurs, couverture, provenance et historique, sans confondre inactive, retired et deleted.

## 2. Problème utilisateur
Une règle peut être retirée sans analyser les consommateurs, les gaps de couverture, la période de coexistence ou le retour arrière, ou être supprimée alors que seul un retirement était requis.

## 3. Objectifs
- préparer deactivation and retirement with reason and scope
- identifier consumers, historical Signals/Incidents and dependencies
- identifier affected coverage and gaps
- identify and compare replacement candidates and coexistence period
- define transition, rollback and governed change
- confirm runtime state and archive without deleting history

## 4. Non-objectifs
Aucune API, protocole, syntaxe vendor, moteur, langage, compilateur, parser, AST, format technique, commande, code, modèle ML, mutation runtime directe, capability CAP-INV-5xx, Threat Intelligence, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède **Retirement Proposal**, **Replacement Relation** et l’assessment fonctionnelle de transition. Command garde Detection/Signal/Alert/Incident ; Settings les runtimes, targets, environments et health ; Endpoint ses capacités et exécutions locales ; Govern Action Request/Decision/Approval/Response Run/Result ; Studio Tools et Automation Runs ; Shared les mécanismes transverses.

## 6. Utilisateurs
Principal : **Detection Owner**. Secondaires : Detection Engineer, Detection Reviewer, Incident Commander, Platform Operator et Approver autorisés.

## 7. Conditions d’entrée
Tenant, environnement, target, version, owner, consumers, permissions, restrictions et return origin sont explicites. Les preuves 4B.3A.1 et assessments runtime sont référencées sans duplication. Toute absence produit un état incomplete, blocked, gap-identified ou replacement-pending.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | --- | --- | --- |
| Runtime/version/consumer context | CAP-INV-425..428 / Command | active state, version, Signals, Alerts, Incidents and consumers | oui | current and historical | proposal incomplete |
| Coverage and gap evidence | CAP-INV-416 / CAP-INV-431 | affected coverage, gaps and dependencies | oui | current assessment | retirement blocked |
| Replacement candidate | CAP-INV-406..417 / CAP-INV-429 | candidate version/content and comparative evidence | non | latest candidate | replacement-pending |
| Govern/runtime context | Govern / Settings / Endpoint | authority, target capability, Result and rollback constraints | oui for execution | current and valid | proposal only |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Runtime Detection / Version | Command | active/inactive state and history | lecture |
| Signal / Alert / Incident | Command | historical consumers and outcomes | lecture/lien |
| Coverage Assessment / Detection Gap | Investigate | affected coverage and unresolved gaps | lecture |
| Decision / Approval / Response Run / Result | Govern | authority and execution result | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Retirement Proposal | créer, modifier, comparer, retirer, superseder | Investigate concept | proposal ≠ execution |
| Replacement Relation | proposer, contester, versionner | Investigate concept | replacement ≠ equivalent coverage |
| Deactivation/retirement Action Request context | préparer | Govern | no direct runtime change |
| Retirement assessment | enregistrer | Investigate concept | retired ≠ deleted |

## 11. Fonctionnalités
- préparer deactivation and retirement with reason and scope
- identifier consumers, historical Signals/Incidents and dependencies
- identifier affected coverage and gaps
- identify and compare replacement candidates and coexistence period
- define transition, rollback and governed change
- confirm runtime state and archive without deleting history
- conserver sources, versions, erreurs, partialité, attribution et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | --- | --- | --- | --- |
| consulter/comparer | Detection Owner | Retirement Proposal / Replacement Relation | 0 | lecture autorisée | projection sourcée | non |
| assessment bornée | Detection Owner | comparison/coverage assessment | 1 | lancement explicite | résultat attribué | policy |
| créer/modifier/contester | Detection Owner | Retirement Proposal | 2 | mutation réversible | nouvelle version | OPEN-013 |
| préparer demande | Detection Owner | Action Request context | 2 | risque, target, gap and rollback visible | package non effectif | Govern |
| désactiver/retirer réellement | Govern/runtime owner | runtime target | 3 | Decision/Approval | Result projeté | obligatoire |
| supprimer contenu/historique | aucun rôle local | provenance | 4 | interdit par défaut | refus | strict |

Investigate n’exécute jamais les classes 3/4.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | --- | --- | --- | --- | --- |
| construire/revoir retirement and replacement | oui | formulaires et contrôles | oui | proposition | matrices et revue humaine |
| comparer coverage/consumers | oui | diff et règles explicables | oui | résumé | tables, reports and checklist |
| préparer transition/rollback | oui | templates et validations | oui | draft | workflow non agentique |
| approuver/exécuter | non localement | non | non | interdit | Govern/runtime owner |

Initiateur, agent/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner et disposition humaine sont visibles. Aucun choix silencieux.

## 14. États fonctionnels
`proposed`, `under-review`, `approved-projection`, `scheduled`, `deactivating`, `inactive`, `retired`, `replacement-pending`, `completed`, `rollback-required`, `cancelled`, `blocked`, `superseded`. Projections fonctionnelles, pas machine d’état objet définitive.

## 15. États d’interface
Loading conserve le contexte ; Empty distingue absence et interdiction ; Partial détaille targets et consumers ; Error conserve le valide ; Offline bloque les mutations ; Permission denied masque ; Stale distingue ancien/courant ; les conflits offrent diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Retirement Proposal | governed proposal | CAP-INV-420 / Govern | scope, consumers, coverage, gaps and rollback visible |
| Replacement Relation | relation candidate | CAP-INV-416/435 | evidence and non-equivalence limits visible |
| Retirement Result assessment | assessment | CAP-INV-425/435 | runtime state and residual gaps reconciled |
| Archived lineage | provenance package | Reporting/Audit | history retained, not deleted |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-427/428/431/432 | retirement considered | CAP-INV-434 | quality, performance, drift, consumers and gaps | source assessment |
| CAP-INV-434 | replacement required | CAP-INV-406/416/429 | required behavior, gap, candidate relation and tests | Retirement |
| CAP-INV-434 | submit deactivation/retirement | CAP-INV-420/Govern | reason, scope, targets, consumers, rollback and authority | Retirement Proposal |
| Govern Result / runtime observation | state changed | CAP-INV-425/434 | target state, version, partiality, errors and timestamps | Govern/runtime |
| CAP-INV-434 | completed or blocked | CAP-INV-435 | final disposition, coverage impact and unresolved items | Retirement |

Chaque transition conserve ownership, tenant/env, target, version, permissions, restrictions, erreurs, autorité, provenance et return origin.

## 18. Dépendances
CAP-INV-406,416,420,425..433,435; Command runtime and historical outcomes; Settings/Endpoint targets; Govern change authority. OPEN-017 couvre la future stratégie runtime/langage/portabilité sans option sélectionnée. Shared Jobs, Notifications, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Collaboration, Comparison, Inspector, Audit Hooks et Recovery sont consommés sans redéfinition.

## 19. Source de vérité
Investigate possède l’assessment/proposition locale. Les objets Command, Settings, Endpoint, Govern, Studio et Shared restent canoniques chez leurs owners. Une projection ne remplace jamais sa source.

## 20. Provenance et audit
Conserver besoin, Project, Hypothesis, Drafts/Versions, Review Package, Release Candidate, reviews, readiness, targets, plans, Action Requests, Decisions, Approvals, Runs, Results, runtime observations, health, Signals/Alerts/Incidents, assessments, propositions, erreurs, auteurs, timestamps, Tools/Runs, paramètres et dispositions. Aucune trace supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Retirement Proposal create/update | loss of coverage | 2 | target/consumer scope visible | OPEN-013 | proposer/reviewer | Investigate | Permissions |
| Consumer/history/coverage read | operational and historical sensitivity | 0 | scoped and minimized | possible | viewer/owner | Command/Investigate | Permissions |
| Deactivation/retirement request submit | production change | 2 | no direct execution | step-up required | requester/approver/operator | Govern | Permissions |
| Runtime deactivation/retirement | production action | 3 | not available locally | Govern mandatory | approver/operator | Govern/Settings/runtime | Permissions |

Permissions atomiques, namespaces, RBAC/ABAC, step-up et séparation finale restent reportés.

## 22. Limites et erreurs
- Deactivated ≠ retired; retired ≠ deleted.
- Replacement ≠ equivalent coverage certaine.
- Retirement with unresolved gap may remain blocked.
- Historical Signals, Alerts, Incidents, versions and provenance are not deleted.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, target offline et version superseded restent visibles.
- Result, match, health, disposition Command ou sortie IA ne vaut pas conclusion universelle.

## 23. Métriques conceptuelles
- retirement proposals by disposition
- affected consumers and coverage gaps
- replacements by evidence level
- direct deletions by Investigate — target zero
- provenance et dispositions humaines complètes
- exécutions silencieuses, auto-approbations et suppressions de trace : cible zéro

Aucun seuil technique universel n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; documentation only. Aucun `validated`, `implemented`, `deployed`, `active`, `native` ou `integrated`.

## 25. Critères d’acceptation
### 1. Retirement avec gap
**Given** un retirement proposal sans replacement et un coverage gap  
**When** la revue est conduite  
**Then** le gap reste visible, le retirement peut être blocked, aucun contenu n’est supprimé et un remplacement peut être demandé

### 2. Désactivation réussie
**Given** Govern Result confirms inactive on all targets  
**When** retirement is assessed  
**Then** inactive remains distinct from retired and historical lineage is preserved

### 3. Sans IA
**Given** aucun modèle  
**When** deactivation/retirement is prepared  
**Then** checklists, consumer maps, coverage comparison and human review suffice

## 26. Questions ouvertes
- OPEN-008 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-015 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-017 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-005 reste forensic-only.
- Schémas, formats, cardinalités, permissions atomiques, contrats techniques et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering lifecycle, Command, Settings, Endpoint, Govern, Studio, Shared, Objects, Permissions, Screens, Journeys, Technique, validation et futur handoff 4B.3B non canonique.
