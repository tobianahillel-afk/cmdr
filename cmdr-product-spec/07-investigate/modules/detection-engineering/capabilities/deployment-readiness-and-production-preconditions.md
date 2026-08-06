---
id: CAP-INV-419
title: Deployment Readiness and Production Preconditions
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
  - OPEN-017
source-of-truth: canonical
---
# CAP-INV-419 — Deployment Readiness and Production Preconditions

## 1. Définition
Évaluer si un Release Candidate satisfait les préconditions fonctionnelles de promotion et de production pour chaque target, sans déployer ni confondre readiness et succès runtime.

## 2. Problème utilisateur
Une revue passée peut masquer un target incompatible, une dépendance absente, une séparation des tâches insuffisante ou l’absence de rollback.

## 3. Objectifs
- vérifier candidate, review, tests, replay, FP/FN limits, coverage, gaps et owner
- évaluer séparément chaque tenant, environnement, runtime déclaré et target
- contrôler sources, schemas, fields, enrichments, permissions, support et health
- vérifier séparation des tâches, rollback, observation plan et maintenance window
- classer readiness sans exécuter

## 4. Non-objectifs
- ne pas définir API, protocole, compilateur, parser, AST, format de package, stockage, pipeline, streaming, commande ou code produit
- ne sélectionner aucun moteur, langage, représentation cible, syntaxe vendor, produit tiers ou modèle ML
- ne réaliser aucun déploiement, activation, désactivation, rollback, suppression ou exception réelle pendant la phase documentaire
- ne modifier ni supprimer silencieusement runtime Detection, Signal, Alert ou Incident
- ne créer aucune capability CAP-INV-5xx, aucun objet Threat Intelligence, Cloud Analysis, Mobile Forensics ou écran détaillé

## 5. Propriétaire
Investigate / Detection Engineering possède **Deployment Readiness Assessment** et ses dispositions fonctionnelles. Command conserve runtime Detection, Signal, Alert, Incident et le feedback opérationnel. Platform Settings conserve runtimes configurés, targets, environments, tenants, sources, parsers, schemas, health, providers, secrets et canaux administratifs. Endpoint Agent conserve ses capacités, versions, health et exécutions locales autorisées. Govern conserve Action Request, Decision, Approval, Response Run, Result et toute autorité de classe 3 ou 4. Studio conserve Tool, Tool Call, Workflow, Automation Agent, Human Gate et Automation Run. Shared conserve Jobs, Notifications, Trace, Activity, Linking, Versioning, Comparison, Reporting, Export, Collaboration, Audit Hooks et Recovery.

## 6. Utilisateurs
Principal : **Detection Owner**. Secondaires : Detection Engineer; Detection Reviewer; Detection Owner; Incident Commander; Platform Operator; Approver selon le contexte.

## 7. Conditions d’entrée
Tenant, environnement, version, owner, cible, période, permission, restrictions et return origin sont explicites. Les preuves de CAP-INV-401..417 sont référencées sans duplication. Toute dépendance absente conduit à un état incomplet, partiel, bloqué ou inconnu ; aucune version, santé, efficacité, autorité ou ground truth n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | --- | --- | --- |
| Release Candidate and review | CAP-INV-418 | immutable version and review disposition | oui | current candidate | not-assessed |
| Authoring evidence | CAP-INV-411..417 | tests, replay, FP/FN, coverage, gaps et provenance | oui | selected snapshots | incomplete |
| Target/runtime projections | Platform Settings / Endpoint Agent | environment, tenant, capability, version and health | oui | current snapshot | target-unavailable |
| Governance preconditions | Govern / Security | authority, separation, rollback and observation requirements | oui | current policy | approval-required |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Detection Release Candidate | Investigate | candidate and evidence | lecture |
| Environment / Tenant / Deployment Target | Platform Settings | scope, capability, health and restrictions | lecture |
| Endpoint Agent capability/version/health | Endpoint Agent / Settings | execution capability projection | lecture |
| Decision/Approval policy projections | Govern / Security | authority and separation | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Deployment Readiness Assessment | créer, comparer, contester, superseder | Investigate concept | ready ≠ deployed |
| Target readiness disposition | enregistrer par target | Investigate concept | partial targets do not imply global success |
| Settings/Govern information request | préparer | destination owner | no admin mutation |

## 11. Fonctionnalités
- vérifier candidate, review, tests, replay, FP/FN limits, coverage, gaps et owner
- évaluer séparément chaque tenant, environnement, runtime déclaré et target
- contrôler sources, schemas, fields, enrichments, permissions, support et health
- vérifier séparation des tâches, rollback, observation plan et maintenance window
- classer readiness sans exécuter
- conserver versions, sources, erreurs, résultats partiels, restrictions, attribution et return origin
- fonctionner entièrement sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | --- | --- | --- | --- |
| Consulter, filtrer ou comparer | Detection Owner | Deployment Readiness Assessment | 0 | lecture autorisée et scope explicite | projection sourcée | non |
| Exécuter une assessment ou comparaison bornée | Detection Owner | Tool Call / Assessment Result | 1 | déclenchement explicite, sources, paramètres et permission | résultat attribué, partial ou complete | policy applicable |
| Créer, modifier, contester ou retirer | Detection Owner | Deployment Readiness Assessment | 2 | mutation réversible et versionnée | nouvelle version et disposition humaine | OPEN-013 |
| Préparer une demande gouvernée | Detection Owner | Change/Action Request context | 2 | cible, effet, risque, rollback et autorité visibles | package non effectif | destination Govern |
| Exécuter le changement réel | Govern / Platform Operator | runtime target / Response Run | 3 | Decision/Approval et policy applicables | projection du Result seulement dans Investigate | Govern owner |
| Détruire historique ou provenance | aucun rôle local | historical content | 4 | interdit par défaut | action refusée | strict governance |

Les classes 3 et 4 ne sont jamais exécutées par Investigate. Une demande ou coordination locale de classe 2 ne vaut ni Decision, ni Approval, ni exécution.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | --- | --- | --- | --- | --- |
| Construire ou examiner Deployment Readiness Assessment | oui | checklists, catalogues et contrôles explicables | oui | proposition attribuée | formulaire, tables, matrices et revue humaine |
| Comparer ou valider Deployment Readiness Assessment | oui | comparateur et règles explicables | oui | explication facultative | diff, diagnostics et checklist |
| Résumer risques, erreurs ou contradictions | oui | catalogue et agrégations déterministes | oui | résumé sourcé | sources brutes, timeline et revue manuelle |
| Approuver, promouvoir, activer, désactiver ou rollback | non localement | non | non | interdit | Govern et owner runtime selon autorité |

Toute proposition automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet. Aucun target, état, approbation, suppression, exception ou rollback n’est sélectionné silencieusement.

## 14. États fonctionnels
`not-assessed`, `assessing`, `ready`, `ready-with-conditions`, `incomplete`, `incompatible`, `dependency-missing`, `target-unavailable`, `rollback-required`, `approval-required`, `blocked`, `disputed`, `superseded`. Ces états sont des projections fonctionnelles ; ils ne redéfinissent ni les machines d’état Govern, Command, Settings ou runtime, ni une machine d’état objet définitive.

## 15. États d’interface
Loading conserve version, target, phase et return origin. Empty distingue absence de résultat et absence d’accès. Partial expose chaque target et résultat utilisable. Error conserve les données valides et le correlation ID. Offline interdit les mutations non garanties. Permission denied masque les données protégées. Stale sépare dernière observation connue et état actuel. Les conflits de version fournissent diff et reprise sûre. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Deployment Readiness Assessment | assessment | CAP-INV-421/420 | readiness separated by target and condition |
| Target blocker set | functional blockers | Project / Settings request | reasons and owners visible |
| Approval requirement context | handoff context | CAP-INV-420 | no authority invented |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-418 | review-passed | CAP-INV-419 | candidate, evidence, conditions and unresolved items | Review Queue |
| CAP-INV-419 | ready or ready-with-conditions | CAP-INV-421 | targets, conditions, restrictions and rollback requirements | Readiness |
| CAP-INV-419 | approval-required or blocked | CAP-INV-420 / source owner | risk, blocker and missing dependency | Readiness |

Chaque transition conserve owner, tenant, environnement, target, version, permissions, restrictions, erreurs, autorité, provenance et return origin. Elle n’accorde aucun accès ni pouvoir supplémentaire.

## 18. Dépendances
CAP-INV-404,405,409,411..418; Settings Environment/Target/Health; Endpoint; Govern/Security. OPEN-017 gouverne le futur choix runtime/langage/portabilité sans option sélectionnée. Shared Background Jobs, Notifications, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Collaboration, Comments, Assignments, Comparison, Inspector, Context Bar, Audit Hooks et Recovery sont consommés sans redéfinition.

## 19. Source de vérité
Investigate est source de vérité de Deployment Readiness Assessment comme concept fonctionnel. Command reste source de runtime Detection, Signal, Alert et Incident. Settings reste source des environments, targets, runtimes configurés et health administratif. Govern reste source des Action Requests, Decisions, Approvals, Response Runs et Results. Studio reste source des Tools et Automation Runs. Une projection locale ne remplace jamais son objet canonique propriétaire.

## 20. Provenance et audit
Enregistrer le besoin initial, Project, Hypothesis, Drafts, versions, Review Packages, Release Candidates, reviewers, readiness, targets, plans, Action Requests, Decisions, Approvals, Response Runs, Results, runtime observations, health, Signals/Alerts/Incidents liés, assessments, propositions, erreurs, interruptions, auteurs, timestamps, Tool Calls, Automation Runs, paramètres, dispositions humaines, exports et correlation IDs applicables à Deployment Readiness Assessment. Aucune trace n’est supprimée ni réécrite silencieusement.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Deployment Readiness read/create/update | future production gate | 0/2 | target details scoped | OPEN-013 | assessor/owner/reviewer | Investigate | Permissions |
| Target environment read | infrastructure exposure | 0 | secrets masked | possible | viewer/admin distinct | Settings | Permissions |
| Cross-tenant readiness assess | tenant isolation | 1/2 | aggregate/minimize | step-up probable | independent reviewer | Security/Settings | Permissions |

Les permissions atomiques, namespaces, RBAC/ABAC, step-up définitif et séparation des tâches finale restent reportés. Toute exécution réelle de classe 3 ou 4 reste chez Govern et l’owner runtime.

## 22. Limites et erreurs
- Deployment Readiness ≠ deployment.
- Target selected or configured ≠ compatible.
- Runtime capability declared ≠ available or healthy.
- Ready-with-conditions does not replace Approval.
- Les données peuvent être stale, partielles, restreintes, incohérentes entre tenants ou indisponibles.
- Un Result, score, match, non-match, health state, disposition Command ou sortie IA ne vaut pas conclusion universelle.
- Timeout, cancellation, revocation, target offline, version superseded et partial result restent attribués et visibles.

## 23. Métriques conceptuelles
- assessments by state/target
- incompatibility causes
- stale target projections
- readiness treated as deployment — target zero
- complétude de provenance, versions, targets, autorité et dispositions humaines
- nombre d’exécutions silencieuses, d’auto-approbations ou de suppressions de traces — cible conceptuelle zéro

Aucun seuil universel de latency, coût, précision, recall, health ou business value n’est imposé.

## 24. Classification de livraison
`defined` / `planned`. Preuve documentaire uniquement. Aucun document n’est `validated`, `implemented`, `deployed`, `active`, `native` ou `integrated`. Aucun moteur, langage, syntaxe, produit tiers, modèle ML, API, protocole, commande ou code n’est choisi.

## 25. Critères d’acceptation
### 1. Target incompatible
**Given** un candidate, un target compatible et un incompatible  
**When** la readiness est évaluée  
**Then** les targets restent distincts, l’incompatible reste blocked et aucun succès global n’est déclaré

### 2. Evidence missing
**Given** un candidate review-passed avec replay manquant  
**When** la readiness est évaluée  
**Then** l’assessment reste incomplete ou conditional et aucune promotion n’est exécutée

### 3. Sans IA
**Given** aucun modèle  
**When** la readiness est construite  
**Then** checklists, target matrices et contrôles déterministes suffisent

## 26. Questions ouvertes
- OPEN-008 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- OPEN-017 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- Les schémas, cardinalités, machines d’état, formats de déploiement, permissions atomiques, contrats techniques et composition détaillée des écrans restent futurs.
- OPEN-005 demeure forensic-only et n’est pas utilisée pour Detection.

## 27. Consommateurs documentaires
Detection Engineering lifecycle and Capability Maps ; Command runtime Detection/Signal/Alert/Incident ; Platform Settings environments/targets/health ; Endpoint Agent capabilities/version/health ; Govern Action Request/Decision/Approval/Response Run/Result ; Studio Tools/Runs ; Shared mechanisms ; phases Objects, Permissions, Screens, Journeys, Technique, validation et future 4B.3B uniquement comme handoff non canonique.
