---
id: CAP-INV-431
title: Data, Schema and Dependency Drift Assessment
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
# CAP-INV-431 — Data, Schema and Dependency Drift Assessment

## 1. Définition
Détecter ou examiner les changements de source, schéma, field, type, mapping, normalisation, enrichissement, runtime, policy, rétention, qualité, volume ou dépendance et évaluer leur compatibilité sans déclarer automatiquement la règle cassée.

## 2. Problème utilisateur
Un changement de schéma ou de dépendance peut être assimilé à une panne certaine, ou rester invisible jusqu’à une dégradation runtime.

## 3. Objectifs
- identifier source/new-source disappearance and version changes
- examiner schema, field, type, mapping and normalization changes
- examiner enrichment, runtime, policy, retention, quality and volume changes
- identifier affected contents/targets and tests to rerun
- classer compatibility and prepare revalidation, new Draft or remediation

## 4. Non-objectifs
Aucune API, protocole, syntaxe vendor, moteur, langage, compilateur, parser, AST, format technique, commande, code, modèle ML, mutation runtime directe, capability CAP-INV-5xx, Threat Intelligence, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède **Drift Assessment** comme concept fonctionnel. Command garde Detection/Signal/Alert/Incident ; Settings les runtimes, targets, environments et health ; Endpoint ses capacités et exécutions locales ; Govern Action Request/Decision/Approval/Response Run/Result ; Studio Tools et Automation Runs ; Shared les mécanismes transverses.

## 6. Utilisateurs
Principal : **Detection Engineer**. Secondaires : Detection Engineer, Detection Reviewer, Detection Owner, Incident Commander, Platform Operator et Approver autorisés.

## 7. Conditions d’entrée
Tenant, environnement, target, version, période, owner, permissions, restrictions et return origin sont explicites. Les preuves 4B.3A.1 sont référencées sans duplication. Toute absence produit un état incomplete, partial, blocked ou unknown.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | --- | --- | --- |
| Current dependency snapshots | Settings / Shared / Endpoint | source, schema, field, mapping, runtime, policy and health versions | oui | current and historical | compatibility unknown |
| Detection dependencies | CAP-INV-404..410 / runtime version | required sources, fields, enrichments and target versions | oui | linked versions | affected content unknown |
| Health/performance observations | CAP-INV-426/432 | errors, latency, gaps and affected targets | non | same period | impact uncertain |
| Change/audit events | Settings / Shared Activity / Govern Result | declared changes and timestamps | non | resolvable | drift candidate only |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Data Source / Parser / Schema / Environment | Settings | versions, changes and health | lecture |
| Field/Normalization/Data Quality projections | Shared / Settings | semantic and quality changes | lecture |
| Runtime Detection / Version | Command | affected runtime projection | lecture |
| Detection Content dependencies | Investigate | expected refs and limits | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Drift Assessment | créer, confirmer, comparer, contester, superseder | Investigate concept | drift ≠ failure |
| Affected-content relation | lier/versionner | Investigate | scope and confidence visible |
| Revalidation/new Draft/remediation context | préparer | CAP-INV-411/406/Settings | no active mutation |

## 11. Fonctionnalités
- identifier source/new-source disappearance and version changes
- examiner schema, field, type, mapping and normalization changes
- examiner enrichment, runtime, policy, retention, quality and volume changes
- identifier affected contents/targets and tests to rerun
- classer compatibility and prepare revalidation, new Draft or remediation
- conserver sources, versions, erreurs, partialité, attribution et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | --- | --- | --- | --- |
| consulter/comparer | Detection Engineer | Drift Assessment | 0 | lecture autorisée | projection sourcée | non |
| assessment bornée | Detection Engineer | Tool Call / assessment | 1 | lancement explicite | résultat attribué | policy |
| créer/modifier/contester | Detection Engineer | Drift Assessment | 2 | mutation réversible | nouvelle version | OPEN-013 |
| préparer demande | Detection Engineer | Action Request context | 2 | risque, cible et rollback visibles | package non effectif | Govern |
| changement réel | Govern/runtime owner | runtime target | 3 | Decision/Approval | Result projeté | owner |
| détruire historique | aucun rôle local | provenance | 4 | interdit | refus | strict |

Investigate n’exécute jamais les classes 3/4.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | --- | --- | --- | --- | --- |
| construire/revoir Drift Assessment | oui | formulaires et contrôles | oui | proposition | matrices et revue humaine |
| comparer/valider | oui | diff et règles explicables | oui | explication | diagnostics et checklist |
| résumer risques/erreurs | oui | agrégations | oui | résumé sourcé | tables et timeline |
| approuver/exécuter | non localement | non | non | interdit | Govern/runtime owner |

Initiateur, agent/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner et disposition humaine sont visibles. Aucun choix silencieux.

## 14. États fonctionnels
`no-drift-observed`, `drift-candidate`, `confirmed-change`, `compatibility-unknown`, `compatible`, `partially-compatible`, `incompatible`, `blocked`, `remediation-required`, `disputed`, `superseded`. Projections fonctionnelles, pas machine d’état objet définitive.

## 15. États d’interface
Loading conserve le contexte ; Empty distingue absence et interdiction ; Partial détaille les targets ; Error conserve le valide ; Offline bloque les mutations ; Permission denied masque ; Stale distingue ancien/courant ; les conflits offrent diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Drift Assessment | assessment | CAP-INV-411/429/432/435 | change, compatibility and affected scope visible |
| Revalidation request | test context | CAP-INV-411..417 | versions and changed dependency linked |
| New Draft or Settings request | change context | CAP-INV-406 / Settings | no silent correction |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| Settings/Shared/Endpoint | change detected or declared | CAP-INV-431 | before/after versions, scope, source and timestamp | source owner |
| CAP-INV-431 | compatibility unknown/change confirmed | CAP-INV-411..417 | affected candidate/runtime refs, tests to rerun and limitations | Drift |
| CAP-INV-431 | content change required | CAP-INV-406/429 | change, impact, target versions and provenance | Drift |
| CAP-INV-431 | admin remediation needed | Settings/Govern request | dependency, target and risk context | Drift |

Chaque transition conserve ownership, tenant/env, target, version, permissions, restrictions, erreurs, autorité, provenance et return origin.

## 18. Dépendances
CAP-INV-404..417,425,426,429,432,435; Settings/Shared Data Quality; Endpoint; Command runtime. OPEN-017 couvre la future stratégie runtime/langage/portabilité sans option sélectionnée. Shared Jobs, Notifications, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Collaboration, Comparison, Inspector, Audit Hooks et Recovery sont consommés sans redéfinition.

## 19. Source de vérité
Investigate possède l’assessment/proposition locale. Les objets Command, Settings, Endpoint, Govern, Studio et Shared restent canoniques chez leurs owners. Une projection ne remplace jamais sa source.

## 20. Provenance et audit
Conserver besoin, Project, Hypothesis, Drafts/Versions, Review Package, Release Candidate, reviews, readiness, targets, plans, Action Requests, Decisions, Approvals, Runs, Results, runtime observations, health, Signals/Alerts/Incidents, assessments, propositions, erreurs, auteurs, timestamps, Tools/Runs, paramètres et dispositions. Aucune trace supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Drift assessment read/create/update | production compatibility risk | 0/2 | infrastructure details scoped | OPEN-013 | assessor/source owner | Investigate | Permissions |
| Change/audit projections read | administrative data | 0 | secrets masked | possible | viewer/admin | Settings/Shared | Permissions |
| Revalidation/remediation request | future impact | 2 | affected targets explicit | step-up possible | requester/owner | Investigate/Settings/Govern | Permissions |

Permissions atomiques, namespaces, RBAC/ABAC, step-up et séparation finale restent reportés.

## 22. Limites et erreurs
- Schema drift ≠ certain rule failure.
- Dependency change ≠ certain incompatibility.
- No active content is changed automatically.
- Compatibility remains unknown until evidence supports a disposition.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, target offline et version superseded restent visibles.
- Result, match, health, disposition Command ou sortie IA ne vaut pas conclusion universelle.

## 23. Métriques conceptuelles
- drift candidates/confirmed changes by type
- affected contents/targets
- time to compatibility disposition
- automatic active fixes — target zero
- provenance et dispositions humaines complètes
- exécutions silencieuses, auto-approbations et suppressions de trace : cible zéro

Aucun seuil technique universel n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; documentation only. Aucun `validated`, `implemented`, `deployed`, `active`, `native` ou `integrated`.

## 25. Critères d’acceptation
### 1. Schema drift
**Given** un field changes and compatibility is unknown  
**When** drift is detected  
**Then** the rule is not automatically broken, versions remain visible and revalidation/new Draft may be requested

### 2. Compatible change
**Given** a parser version changes without semantic impact evidence  
**When** assessment completes  
**Then** compatible may be recorded with sources and no active mutation

### 3. Sans IA
**Given** aucun modèle  
**When** drift is assessed  
**Then** version diffs, catalogs, tests and human review suffice

## 26. Questions ouvertes
- OPEN-008 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-015 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-017 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-005 reste forensic-only.
- Schémas, formats, cardinalités, permissions atomiques, contrats techniques et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering lifecycle, Command, Settings, Endpoint, Govern, Studio, Shared, Objects, Permissions, Screens, Journeys, Technique, validation et futur handoff 4B.3B non canonique.
