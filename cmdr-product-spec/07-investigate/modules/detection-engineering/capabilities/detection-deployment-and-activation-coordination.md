---
id: CAP-INV-424
title: Detection Deployment and Activation Coordination
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
# CAP-INV-424 — Detection Deployment and Activation Coordination

## 1. Définition
Coordonner dans Investigate la lecture des Decision, Approval, Response Run et Result d’un changement de détection, suivre progression et résultats par target et demander un arrêt selon l’autorité, sans administrer le runtime.

## 2. Problème utilisateur
Un statut agrégé peut masquer une activation non confirmée, des targets offline, une incompatibilité ou un résultat partiel, et faire croire qu’Investigate a exécuté le changement.

## 3. Objectifs
- recevoir Decision/Approval et voir Response Run, targets, version, phases and progress
- voir results, errors, incompatibilities, offline targets and partiality
- distinguer promotion, deployment, activation, deactivation and rollback projections
- demander information, interruption ou arrêt selon l’autorité
- relier Result et retourner au Project

## 4. Non-objectifs
- ne pas définir API, protocole, compilateur, parser, AST, format de package, stockage, pipeline, streaming, commande ou code produit
- ne sélectionner aucun moteur, langage, représentation cible, syntaxe vendor, produit tiers ou modèle ML
- ne réaliser aucun déploiement, activation, désactivation, rollback, suppression ou exception réelle pendant la phase documentaire
- ne modifier ni supprimer silencieusement runtime Detection, Signal, Alert ou Incident
- ne créer aucune capability CAP-INV-5xx, aucun objet Threat Intelligence, Cloud Analysis, Mobile Forensics ou écran détaillé

## 5. Propriétaire
Investigate / Detection Engineering possède **Deployment and Activation Coordination** et ses dispositions fonctionnelles. Command conserve runtime Detection, Signal, Alert, Incident et le feedback opérationnel. Platform Settings conserve runtimes configurés, targets, environments, tenants, sources, parsers, schemas, health, providers, secrets et canaux administratifs. Endpoint Agent conserve ses capacités, versions, health et exécutions locales autorisées. Govern conserve Action Request, Decision, Approval, Response Run, Result et toute autorité de classe 3 ou 4. Studio conserve Tool, Tool Call, Workflow, Automation Agent, Human Gate et Automation Run. Shared conserve Jobs, Notifications, Trace, Activity, Linking, Versioning, Comparison, Reporting, Export, Collaboration, Audit Hooks et Recovery.

## 6. Utilisateurs
Principal : **Platform Operator**. Secondaires : Detection Engineer; Detection Reviewer; Detection Owner; Incident Commander; Platform Operator; Approver selon le contexte.

## 7. Conditions d’entrée
Tenant, environnement, version, owner, cible, période, permission, restrictions et return origin sont explicites. Les preuves de CAP-INV-401..417 sont référencées sans duplication. Toute dépendance absente conduit à un état incomplet, partiel, bloqué ou inconnu ; aucune version, santé, efficacité, autorité ou ground truth n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | --- | --- | --- |
| Detection Change Request / Action Request | CAP-INV-420 / Govern | candidate, targets, change type and risks | oui | submitted/current | coordination unavailable |
| Decision and Approval | Govern | authority, conditions, expiry and scope | oui for class 3 | current and unrevoked | awaiting-approval |
| Response Run / Result | Govern | progress, target results, errors and rollback | oui for observation | live/versioned | unknown progress |
| Target/runtime projections | Settings / Endpoint / Command | observed application, activation and health | oui for reconciliation | current snapshot | state unconfirmed |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Action Request / Decision / Approval / Response Run / Result | Govern | authority, execution and verified effect | lecture/lien |
| Environment / Deployment Target / runtime configuration | Settings | target and application projection | lecture |
| Runtime Detection | Command | runtime identity and state projection | lecture |
| Release Candidate / Project | Investigate | expected version and return context | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Deployment coordination record | annoter, link, dispute, close | Investigate concept | not execution owner |
| Stop/information request | préparer and submit | Govern / runtime owner | authority checked |
| Project return disposition | enregistrer | Investigate / Shared | Result and partiality preserved |

## 11. Fonctionnalités
- recevoir Decision/Approval et voir Response Run, targets, version, phases and progress
- voir results, errors, incompatibilities, offline targets and partiality
- distinguer promotion, deployment, activation, deactivation and rollback projections
- demander information, interruption ou arrêt selon l’autorité
- relier Result et retourner au Project
- conserver versions, sources, erreurs, résultats partiels, restrictions, attribution et return origin
- fonctionner entièrement sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | --- | --- | --- | --- |
| Consulter, filtrer ou comparer | Platform Operator | Deployment and Activation Coordination | 0 | lecture autorisée et scope explicite | projection sourcée | non |
| Exécuter une assessment ou comparaison bornée | Platform Operator | Tool Call / Assessment Result | 1 | déclenchement explicite, sources, paramètres et permission | résultat attribué, partial ou complete | policy applicable |
| Créer, modifier, contester ou retirer | Platform Operator | Deployment and Activation Coordination | 2 | mutation réversible et versionnée | nouvelle version et disposition humaine | OPEN-013 |
| Préparer une demande gouvernée | Platform Operator | Change/Action Request context | 2 | cible, effet, risque, rollback et autorité visibles | package non effectif | destination Govern |
| Exécuter le changement réel | Govern / Platform Operator | runtime target / Response Run | 3 | Decision/Approval et policy applicables | projection du Result seulement dans Investigate | Govern owner |
| Détruire historique ou provenance | aucun rôle local | historical content | 4 | interdit par défaut | action refusée | strict governance |

Les classes 3 et 4 ne sont jamais exécutées par Investigate. Une demande ou coordination locale de classe 2 ne vaut ni Decision, ni Approval, ni exécution.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | --- | --- | --- | --- | --- |
| Construire ou examiner Deployment and Activation Coordination | oui | checklists, catalogues et contrôles explicables | oui | proposition attribuée | formulaire, tables, matrices et revue humaine |
| Comparer ou valider Deployment and Activation Coordination | oui | comparateur et règles explicables | oui | explication facultative | diff, diagnostics et checklist |
| Résumer risques, erreurs ou contradictions | oui | catalogue et agrégations déterministes | oui | résumé sourcé | sources brutes, timeline et revue manuelle |
| Approuver, promouvoir, activer, désactiver ou rollback | non localement | non | non | interdit | Govern et owner runtime selon autorité |

Toute proposition automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet. Aucun target, état, approbation, suppression, exception ou rollback n’est sélectionné silencieusement.

## 14. États fonctionnels
`planned`, `scheduled`, `awaiting-approval`, `approved`, `deploying`, `partially-deployed`, `deployed`, `activating`, `active`, `activation-failed`, `deactivating`, `inactive`, `rollback-requested`, `rolling-back`, `rolled-back`, `failed`, `cancelled`, `superseded`. Ces états sont des projections fonctionnelles ; ils ne redéfinissent ni les machines d’état Govern, Command, Settings ou runtime, ni une machine d’état objet définitive.

## 15. États d’interface
Loading conserve version, target, phase et return origin. Empty distingue absence de résultat et absence d’accès. Partial expose chaque target et résultat utilisable. Error conserve les données valides et le correlation ID. Offline interdit les mutations non garanties. Permission denied masque les données protégées. Stale sépare dernière observation connue et état actuel. Les conflits de version fournissent diff et reprise sûre. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Deployment/activation projection | coordination view | CAP-INV-425/426 | per-target expected and reported states |
| Stop or information request | governed request | Govern | no direct runtime command |
| Result-linked project update | lifecycle event | Detection Project | canonical Govern refs and return origin |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| Govern Decision/Approval | authorized execution | Response Run / Settings runtime target | canonical Govern-owned context; Investigate observes | Govern |
| Response Run / Result | progress or completion | CAP-INV-424 | phases, targets, errors, partiality and rollback state | Govern |
| CAP-INV-424 | result or observed state available | CAP-INV-425/426 | expected version, observed hints, targets and timestamps | Project |
| CAP-INV-424 | stop requested | Govern | reason, authority, criteria and current phase | Coordination |

Chaque transition conserve owner, tenant, environnement, target, version, permissions, restrictions, erreurs, autorité, provenance et return origin. Elle n’accorde aucun accès ni pouvoir supplémentaire.

## 18. Dépendances
CAP-INV-420..423,425,426,433,434; Govern; Settings; Endpoint; Command; Shared Run Shell/Trace. OPEN-017 gouverne le futur choix runtime/langage/portabilité sans option sélectionnée. Shared Background Jobs, Notifications, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Collaboration, Comments, Assignments, Comparison, Inspector, Context Bar, Audit Hooks et Recovery sont consommés sans redéfinition.

## 19. Source de vérité
Investigate est source de vérité de Deployment and Activation Coordination comme concept fonctionnel. Command reste source de runtime Detection, Signal, Alert et Incident. Settings reste source des environments, targets, runtimes configurés et health administratif. Govern reste source des Action Requests, Decisions, Approvals, Response Runs et Results. Studio reste source des Tools et Automation Runs. Une projection locale ne remplace jamais son objet canonique propriétaire.

## 20. Provenance et audit
Enregistrer le besoin initial, Project, Hypothesis, Drafts, versions, Review Packages, Release Candidates, reviewers, readiness, targets, plans, Action Requests, Decisions, Approvals, Response Runs, Results, runtime observations, health, Signals/Alerts/Incidents liés, assessments, propositions, erreurs, interruptions, auteurs, timestamps, Tool Calls, Automation Runs, paramètres, dispositions humaines, exports et correlation IDs applicables à Deployment and Activation Coordination. Aucune trace n’est supprimée ni réécrite silencieusement.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Deployment/Result read | sensitive production state | 0 | target details scoped | possible | observer/operator | Govern/Settings | Permissions |
| Stop/information request | operational interruption | 2 request / 3 execution | authority visible | step-up likely | requester/operator | Govern | Permissions |
| Deployment/activation execution | production change | 3 | not available locally | mandatory Govern | requester/approver/operator | Govern/Settings/runtime | Permissions |

Les permissions atomiques, namespaces, RBAC/ABAC, step-up définitif et séparation des tâches finale restent reportés. Toute exécution réelle de classe 3 ou 4 reste chez Govern et l’owner runtime.

## 22. Limites et erreurs
- Investigate does not administer or execute the runtime.
- Approval ≠ Response Run; Response Run ≠ successful deployment.
- Deployment started ≠ activation completed; deployed ≠ active.
- Partial deployment ≠ success.
- Les données peuvent être stale, partielles, restreintes, incohérentes entre tenants ou indisponibles.
- Un Result, score, match, non-match, health state, disposition Command ou sortie IA ne vaut pas conclusion universelle.
- Timeout, cancellation, revocation, target offline, version superseded et partial result restent attribués et visibles.

## 23. Métriques conceptuelles
- runs by phase and target outcome
- offline/incompatible/partial targets
- time to verified Result
- direct runtime mutations by Investigate — target zero
- complétude de provenance, versions, targets, autorité et dispositions humaines
- nombre d’exécutions silencieuses, d’auto-approbations ou de suppressions de traces — cible conceptuelle zéro

Aucun seuil universel de latency, coût, précision, recall, health ou business value n’est imposé.

## 24. Classification de livraison
`defined` / `planned`. Preuve documentaire uniquement. Aucun document n’est `validated`, `implemented`, `deployed`, `active`, `native` ou `integrated`. Aucun moteur, langage, syntaxe, produit tiers, modèle ML, API, protocole, commande ou code n’est choisi.

## 25. Critères d’acceptation
### 1. Deployment sans activation confirmée
**Given** un Result indique configuration copied mais aucune confirmation runtime  
**When** l’état est consulté  
**Then** le statut ne devient pas active, expected/observed remain distinct and health stays unknown

### 2. Partial deployment
**Given** des targets réussissent et un target échoue  
**When** le run est observé  
**Then** partially-deployed remains visible and no success is claimed

### 3. Sans IA
**Given** aucun modèle  
**When** le changement est coordonné  
**Then** Run Shell, per-target tables, traces and human requests suffice

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- OPEN-015 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- OPEN-017 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- Les schémas, cardinalités, machines d’état, formats de déploiement, permissions atomiques, contrats techniques et composition détaillée des écrans restent futurs.
- OPEN-005 demeure forensic-only et n’est pas utilisée pour Detection.

## 27. Consommateurs documentaires
Detection Engineering lifecycle and Capability Maps ; Command runtime Detection/Signal/Alert/Incident ; Platform Settings environments/targets/health ; Endpoint Agent capabilities/version/health ; Govern Action Request/Decision/Approval/Response Run/Result ; Studio Tools/Runs ; Shared mechanisms ; phases Objects, Permissions, Screens, Journeys, Technique, validation et future 4B.3B uniquement comme handoff non canonique.
