---
id: CAP-INV-430
title: Suppression and Exception Proposal Management
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
# CAP-INV-430 — Suppression and Exception Proposal Management

## 1. Définition
Préparer des suppressions et exceptions bornées par scope, targets, durée, expiration, risques et monitoring compensatoire, puis les transmettre à Govern sans créer d’état actif local.

## 2. Problème utilisateur
Une suppression temporaire peut devenir permanente, être confondue avec la suppression d’un Signal historique ou contourner silencieusement la gouvernance.

## 3. Objectifs
- préparer suppression or exception with justification, scope, targets, duration and expiry
- documenter requester, required approver, risks, alternatives and compensating monitoring
- documenter removal criteria and affected objects
- préparer Action Request and receive Decision/runtime state
- préparer renewal or revocation with full audit

## 4. Non-objectifs
Aucune API, protocole, syntaxe vendor, moteur, langage, compilateur, parser, AST, format technique, commande, code, modèle ML, mutation runtime directe, capability CAP-INV-5xx, Threat Intelligence, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède **Suppression Proposal and Exception Proposal** comme concept fonctionnel. Command garde Detection/Signal/Alert/Incident ; Settings les runtimes, targets, environments et health ; Endpoint ses capacités et exécutions locales ; Govern Action Request/Decision/Approval/Response Run/Result ; Studio Tools et Automation Runs ; Shared les mécanismes transverses.

## 6. Utilisateurs
Principal : **Detection Owner**. Secondaires : Detection Engineer, Detection Reviewer, Detection Owner, Incident Commander, Platform Operator et Approver autorisés.

## 7. Conditions d’entrée
Tenant, environnement, target, version, période, owner, permissions, restrictions et return origin sont explicites. Les preuves 4B.3A.1 sont référencées sans duplication. Toute absence produit un état incomplete, partial, blocked ou unknown.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | --- | --- | --- |
| Runtime quality / match evidence | CAP-INV-427/428 | noise, benign context, version and limits | oui | selected evidence | proposal incomplete |
| Scope/target/runtime projections | Settings / Command | targets, affected runtime version and population | oui | current snapshot | scope blocked |
| Policy/authority context | Govern / Security | approval, duration, emergency and revocation rules | oui | current policy | not ready |
| Alternatives/monitoring/rollback | CAP-INV-429/426/433 | safer changes, compensating controls and recovery | oui | linked versions | risk incomplete |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Runtime Detection / Signal / Alert history | Command | affected runtime and historical objects | lecture only |
| Environment / Target / policy | Settings | scope and restrictions | lecture |
| Action Request / Decision / Approval / Result | Govern | authority and applied-state projection | lecture/lien |
| Quality/Review/Health assessments | Investigate | evidence and limitations | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Suppression Proposal | créer, modifier, renouveler-proposer, retirer, superseder | Investigate concept | proposal ≠ active suppression |
| Exception Proposal | créer, modifier, renouveler-proposer, révocation-proposer | Investigate concept | bounded and expiring |
| Action Request context | préparer | Govern | no auto-approval |

## 11. Fonctionnalités
- préparer suppression or exception with justification, scope, targets, duration and expiry
- documenter requester, required approver, risks, alternatives and compensating monitoring
- documenter removal criteria and affected objects
- préparer Action Request and receive Decision/runtime state
- préparer renewal or revocation with full audit
- conserver sources, versions, erreurs, partialité, attribution et return origin
- fonctionner sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | --- | --- | --- | --- |
| consulter/comparer | Detection Owner | Suppression Proposal and Exception Proposal | 0 | lecture autorisée | projection sourcée | non |
| assessment bornée | Detection Owner | Tool Call / assessment | 1 | lancement explicite | résultat attribué | policy |
| créer/modifier/contester | Detection Owner | Suppression Proposal and Exception Proposal | 2 | mutation réversible | nouvelle version | OPEN-013 |
| préparer demande | Detection Owner | Action Request context | 2 | risque, cible et rollback visibles | package non effectif | Govern |
| changement réel | Govern/runtime owner | runtime target | 3 | Decision/Approval | Result projeté | owner |
| détruire historique | aucun rôle local | provenance | 4 | interdit | refus | strict |

Investigate n’exécute jamais les classes 3/4.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | --- | --- | --- | --- | --- |
| construire/revoir Suppression Proposal and Exception Proposal | oui | formulaires et contrôles | oui | proposition | matrices et revue humaine |
| comparer/valider | oui | diff et règles explicables | oui | explication | diagnostics et checklist |
| résumer risques/erreurs | oui | agrégations | oui | résumé sourcé | tables et timeline |
| approuver/exécuter | non localement | non | non | interdit | Govern/runtime owner |

Initiateur, agent/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner et disposition humaine sont visibles. Aucun choix silencieux.

## 14. États fonctionnels
`draft`, `incomplete`, `under-review`, `ready-for-request`, `submitted`, `approved-projection`, `rejected`, `active-projection`, `expiring`, `expired`, `revocation-proposed`, `revoked-projection`, `superseded`. Projections fonctionnelles, pas machine d’état objet définitive.

## 15. États d’interface
Loading conserve le contexte ; Empty distingue absence et interdiction ; Partial détaille les targets ; Error conserve le valide ; Offline bloque les mutations ; Permission denied masque ; Stale distingue ancien/courant ; les conflits offrent diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Suppression/Exception Proposal | bounded proposal | CAP-INV-420/Govern | scope, expiry, risks, alternatives and monitoring |
| Runtime state projection | status reference | Lifecycle / CAP-INV-425 | canonical owner result, no local effect |
| Renewal/revocation context | governed request | Govern | previous Decision and expiry retained |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-427/428 | noise or exception need identified | CAP-INV-430 | evidence, affected version, targets and uncertainty | Quality/Review |
| CAP-INV-430 | submit proposal | CAP-INV-420/Govern | scope, duration, expiry, risk, alternatives, monitoring and rollback | Proposal |
| Govern Result | applied/expired/revoked state | CAP-INV-430/425 | canonical refs, target state and timestamps | Govern |
| CAP-INV-430 | prefer content change | CAP-INV-429/406 | evidence and safer tuning context | Proposal |

Chaque transition conserve ownership, tenant/env, target, version, permissions, restrictions, erreurs, autorité, provenance et return origin.

## 18. Dépendances
CAP-INV-420,425..429,433; Command runtime objects; Settings targets/policies; Govern; Security. OPEN-017 couvre la future stratégie runtime/langage/portabilité sans option sélectionnée. Shared Jobs, Notifications, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Collaboration, Comparison, Inspector, Audit Hooks et Recovery sont consommés sans redéfinition.

## 19. Source de vérité
Investigate possède l’assessment/proposition locale. Les objets Command, Settings, Endpoint, Govern, Studio et Shared restent canoniques chez leurs owners. Une projection ne remplace jamais sa source.

## 20. Provenance et audit
Conserver besoin, Project, Hypothesis, Drafts/Versions, Review Package, Release Candidate, reviews, readiness, targets, plans, Action Requests, Decisions, Approvals, Runs, Results, runtime observations, health, Signals/Alerts/Incidents, assessments, propositions, erreurs, auteurs, timestamps, Tools/Runs, paramètres et dispositions. Aucune trace supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Suppression/Exception Proposal create/update | coverage bypass risk | 2 | scope and affected objects visible | OPEN-013 | requester/reviewer | Investigate | Permissions |
| Proposal submit/read Decision | authority boundary | 0/2 | no hidden approval | step-up likely | requester/approver | Govern | Permissions |
| Apply/renew/revoke active suppression/exception | production change | 3 | not available locally | Govern mandatory | requester/approver/operator | Govern/runtime | Permissions |

Permissions atomiques, namespaces, RBAC/ABAC, step-up et séparation finale restent reportés.

## 22. Limites et erreurs
- Suppression ≠ deletion of historical Signal.
- Exception ≠ permanent bypass.
- Proposal ≠ active suppression/exception.
- Expiry, scope and compensating monitoring remain visible.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, target offline et version superseded restent visibles.
- Result, match, health, disposition Command ou sortie IA ne vaut pas conclusion universelle.

## 23. Métriques conceptuelles
- proposals by type/state/scope
- time-bounded vs missing-expiry count
- renewal/revocation outcomes
- historical Signals deleted — target zero
- provenance et dispositions humaines complètes
- exécutions silencieuses, auto-approbations et suppressions de trace : cible zéro

Aucun seuil technique universel n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; documentation only. Aucun `validated`, `implemented`, `deployed`, `active`, `native` ou `integrated`.

## 25. Critères d’acceptation
### 1. Suppression temporaire
**Given** une proposal limitée avec approver requis  
**When** elle est soumise  
**Then** elle ne devient active without Decision, expiry and scope remain mandatory and no historical Signal is deleted

### 2. Expired exception
**Given** an exception projection reaches expiry  
**When** lifecycle is reviewed  
**Then** expired remains visible and renewal requires a new governed request

### 3. Sans IA
**Given** aucun modèle  
**When** proposal is prepared  
**Then** forms, risk matrices, expiry checks and human approval suffice

## 26. Questions ouvertes
- OPEN-008 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-015 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-017 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-005 reste forensic-only.
- Schémas, formats, cardinalités, permissions atomiques, contrats techniques et écrans détaillés restent futurs.

## 27. Consommateurs documentaires
Detection Engineering lifecycle, Command, Settings, Endpoint, Govern, Studio, Shared, Objects, Permissions, Screens, Journeys, Technique, validation et futur handoff 4B.3B non canonique.
