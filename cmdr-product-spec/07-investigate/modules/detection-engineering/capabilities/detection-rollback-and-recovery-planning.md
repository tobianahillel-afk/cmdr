---
id: CAP-INV-433
title: Detection Rollback and Recovery Planning
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
# CAP-INV-433 — Detection Rollback and Recovery Planning

## 1. Définition
Identifier une version de retour et ses targets, vérifier disponibilité/compatibilité, définir critères, autorité, validations, fenêtre, risques et données à préserver, puis suivre le Response Run sans confondre rollback et récupération complète.

## 2. Problème utilisateur
Un rollback demandé ou partiellement terminé peut être présenté comme retour complet alors qu’un target conserve la mauvaise version ou que la santé n’est pas restaurée.

## 3. Objectifs
- identifier return version and verify availability/compatibility
- définir targets, rollback and emergency criteria, authority and window
- définir validations, impacts, risks and data to preserve
- préparer Action Request and observe Response Run/Result
- reconcile observed version/health and open investigation on failure

## 4. Non-objectifs
Aucune API, protocole, syntaxe vendor, moteur, langage, compilateur, parser, AST, format technique, commande, code, modèle ML, mutation runtime directe, capability CAP-INV-5xx, Threat Intelligence, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède **Rollback Plan and Recovery Assessment** comme concept fonctionnel. Command garde Detection/Signal/Alert/Incident ; Settings les runtimes, targets, environments et health ; Endpoint ses capacités et exécutions locales ; Govern Action Request/Decision/Approval/Response Run/Result ; Studio Tools et Automation Runs ; Shared les mécanismes transverses.

## 6. Utilisateurs
Principal : **Detection Owner**. Secondaires : Detection Engineer, Detection Reviewer, Detection Owner, Incident Commander, Platform Operator et Approver autorisés.

## 7. Conditions d’entrée
Tenant, environnement, target, version, période, owner, permissions, restrictions et return origin sont explicites. Les preuves 4B.3A.1 sont référencées sans duplication. Toute absence produit un état incomplete, partial, blocked ou unknown.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | --- | --- | --- |
| Current/previous versions | CAP-INV-425 / Settings | expected, observed and available return version | oui | current snapshots | plan incomplete |
| Failure/stop evidence | CAP-INV-423/424/426/432 | failed targets, health, performance and stop criteria | oui | same incident/change | rollback need unclear |
| Authority and run context | Govern | Decision, Approval, Response Run and rollback policy | oui for execution | current and valid | request only |
| Recovery validation context | CAP-INV-419/411..417 | compatibility, tests and verification criteria | oui | selected versions | recovery unverifiable |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Runtime Detection / Version | Command | current/previous runtime state | lecture |
| Deployment Target / version availability | Settings / Endpoint | return capability and health | lecture |
| Decision / Approval / Response Run / Result | Govern | authority, progress and effects | lecture/lien |
| Validation/Readiness evidence | Investigate | compatibility and checks | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Rollback Plan | créer, modifier, comparer, superseder | Investigate concept | plan ≠ rollback execution |
| Recovery Assessment | enregistrer, contester, superseder | Investigate concept | rolled-back ≠ full recovery |
| Action Request/investigation context | préparer | Govern / CAP-INV-401 | no direct action |

## 11. Fonctionnalités
- identifier return version and verify availability/compatibility
- définir targets, rollback and emergency criteria, authority and window
- définir validations, impacts, risks and data to preserve
- préparer Action Request and observe Response Run/Result
- reconcile observed version/health and open investigation on failure
- conserver sources, versions, erreurs, partialité, attribution et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | --- | --- | --- | --- |
| consulter/comparer | Detection Owner | Rollback Plan and Recovery Assessment | 0 | lecture autorisée | projection sourcée | non |
| assessment bornée | Detection Owner | Tool Call / assessment | 1 | lancement explicite | résultat attribué | policy |
| créer/modifier/contester | Detection Owner | Rollback Plan and Recovery Assessment | 2 | mutation réversible | nouvelle version | OPEN-013 |
| préparer demande | Detection Owner | Action Request context | 2 | risque, cible et rollback visibles | package non effectif | Govern |
| changement réel | Govern/runtime owner | runtime target | 3 | Decision/Approval | Result projeté | owner |
| détruire historique | aucun rôle local | provenance | 4 | interdit | refus | strict |

Investigate n’exécute jamais les classes 3/4.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | --- | --- | --- | --- | --- |
| construire/revoir Rollback Plan and Recovery Assessment | oui | formulaires et contrôles | oui | proposition | matrices et revue humaine |
| comparer/valider | oui | diff et règles explicables | oui | explication | diagnostics et checklist |
| résumer risques/erreurs | oui | agrégations | oui | résumé sourcé | tables et timeline |
| approuver/exécuter | non localement | non | non | interdit | Govern/runtime owner |

Initiateur, agent/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner et disposition humaine sont visibles. Aucun choix silencieux.

## 14. États fonctionnels
`draft`, `incomplete`, `ready-for-request`, `submitted`, `approved-projection`, `rollback-requested`, `rolling-back`, `partially-rolled-back`, `rolled-back`, `recovery-verifying`, `recovery-incomplete`, `failed`, `cancelled`, `superseded`. Projections fonctionnelles, pas machine d’état objet définitive.

## 15. États d’interface
Loading conserve le contexte ; Empty distingue absence et interdiction ; Partial détaille les targets ; Error conserve le valide ; Offline bloque les mutations ; Permission denied masque ; Stale distingue ancien/courant ; les conflits offrent diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Rollback Plan | governed plan | CAP-INV-420/Govern | return version, targets, criteria, risks and validation |
| Recovery Assessment | post-run assessment | CAP-INV-425/426/435 | per-target version, health and residual issues |
| Failure investigation | investigation context | CAP-INV-401/402 | failed target, versions, result and provenance |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| Deployment/canary/health/performance failure | rollback criteria met | CAP-INV-433 | current version, targets, errors, impact and stop evidence | source assessment |
| CAP-INV-433 | submit rollback request | CAP-INV-420/Govern | return version, targets, authority, window, validation and risks | Rollback Plan |
| Govern Response Run/Result | rollback progresses/completes | CAP-INV-433 | per-target results, partiality, errors and timestamps | Govern |
| CAP-INV-433 | verify recovery | CAP-INV-425/426/401 | observed versions, health, failed targets and residual risk | Rollback |

Chaque transition conserve ownership, tenant/env, target, version, permissions, restrictions, erreurs, autorité, provenance et return origin.

## 18. Dépendances
CAP-INV-419,420,423..426,432,434,435; Govern rollback model; Settings/Endpoint; Command. OPEN-017 couvre la future stratégie runtime/langage/portabilité sans option sélectionnée. Shared Jobs, Notifications, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Collaboration, Comparison, Inspector, Audit Hooks et Recovery sont consommés sans redéfinition.

## 19. Source de vérité
Investigate possède l’assessment/proposition locale. Les objets Command, Settings, Endpoint, Govern, Studio et Shared restent canoniques chez leurs owners. Une projection ne remplace jamais sa source.

## 20. Provenance et audit
Conserver besoin, Project, Hypothesis, Drafts/Versions, Review Package, Release Candidate, reviews, readiness, targets, plans, Action Requests, Decisions, Approvals, Runs, Results, runtime observations, health, Signals/Alerts/Incidents, assessments, propositions, erreurs, auteurs, timestamps, Tools/Runs, paramètres et dispositions. Aucune trace supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Rollback Plan create/update | production recovery risk | 2 | targets and versions scoped | OPEN-013 | planner/reviewer | Investigate | Permissions |
| Rollback request submit/read Run | authority boundary | 0/2 | no hidden execution | step-up required | requester/approver/operator | Govern | Permissions |
| Rollback execution | production change | 3 | not available locally | Govern mandatory | operator/approver | Govern/Settings/runtime | Permissions |

Permissions atomiques, namespaces, RBAC/ABAC, step-up et séparation finale restent reportés.

## 22. Limites et erreurs
- Rollback requested ≠ rollback completed.
- Rollback completed ≠ full recovery.
- Failed deployment ≠ rollback completed.
- Per-target partiality and preserved data remain visible.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, target offline et version superseded restent visibles.
- Result, match, health, disposition Command ou sortie IA ne vaut pas conclusion universelle.

## 23. Métriques conceptuelles
- rollback plans/runs by outcome
- targets partially rolled back
- time to recovery verification
- direct rollbacks by Investigate — target zero
- provenance et dispositions humaines complètes
- exécutions silencieuses, auto-approbations et suppressions de trace : cible zéro

Aucun seuil technique universel n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; documentation only. Aucun `validated`, `implemented`, `deployed`, `active`, `native` ou `integrated`.

## 25. Critères d’acceptation
### 1. Rollback partiel
**Given** un approved rollback on multiple targets and one target fails to return  
**When** Result is received  
**Then** rollback remains partial, target remains visible, full recovery is not declared and investigation is opened

### 2. Rollback completed, health unknown
**Given** all versions appear reverted but health is unavailable  
**When** recovery is assessed  
**Then** rolled-back may be observed while recovery remains incomplete/unknown

### 3. Sans IA
**Given** aucun modèle  
**When** rollback is planned and verified  
**Then** checklists, version comparison, Run/Result views and human verification suffice

## 26. Questions ouvertes
- OPEN-008 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-015 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-017 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-005 reste forensic-only.
- Schémas, formats, cardinalités, permissions atomiques, contrats techniques et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering lifecycle, Command, Settings, Endpoint, Govern, Studio, Shared, Objects, Permissions, Screens, Journeys, Technique, validation et futur handoff 4B.3B non canonique.
