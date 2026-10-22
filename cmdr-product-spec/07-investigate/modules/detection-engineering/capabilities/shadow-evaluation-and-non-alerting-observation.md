---
id: CAP-INV-422
title: Shadow Evaluation and Non-Alerting Observation
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
# CAP-INV-422 — Shadow Evaluation and Non-Alerting Observation

## 1. Définition
Préparer, faire autoriser et observer une évaluation shadow non alertante d’une version candidate sur des targets et une période définis, sans produire de Signal opérationnel.

## 2. Problème utilisateur
Une observation shadow peut être confondue avec un replay historique, une activation ou un runtime alertant, et ses matches peuvent être pris pour des Signals.

## 3. Objectifs
- sélectionner version, targets, période, sources, métriques et limites
- définir critères d’arrêt et autorisation requise
- observer matches, non-matches pertinents, erreurs, latency et volume
- comparer au runtime existant sans alerting
- interrompre, clôturer et préparer une recommandation

## 4. Non-objectifs
- ne pas définir API, protocole, compilateur, parser, AST, format de package, stockage, pipeline, streaming, commande ou code produit
- ne sélectionner aucun moteur, langage, représentation cible, syntaxe vendor, produit tiers ou modèle ML
- ne réaliser aucun déploiement, activation, désactivation, rollback, suppression ou exception réelle pendant la phase documentaire
- ne modifier ni supprimer silencieusement runtime Detection, Signal, Alert ou Incident
- ne créer aucune capability CAP-INV-5xx, aucun objet Threat Intelligence, Cloud Analysis, Mobile Forensics ou écran détaillé

## 5. Propriétaire
Investigate / Detection Engineering possède **Shadow Evaluation Plan and Shadow Assessment** et ses dispositions fonctionnelles. Command conserve runtime Detection, Signal, Alert, Incident et le feedback opérationnel. Platform Settings conserve runtimes configurés, targets, environments, tenants, sources, parsers, schemas, health, providers, secrets et canaux administratifs. Endpoint Agent conserve ses capacités, versions, health et exécutions locales autorisées. Govern conserve Action Request, Decision, Approval, Response Run, Result et toute autorité de classe 3 ou 4. Studio conserve Tool, Tool Call, Workflow, Automation Agent, Human Gate et Automation Run. Shared conserve Jobs, Notifications, Trace, Activity, Linking, Versioning, Comparison, Reporting, Export, Collaboration, Audit Hooks et Recovery.

## 6. Utilisateurs
Principal : **Detection Engineer**. Secondaires : Detection Engineer; Detection Reviewer; Detection Owner; Incident Commander; Platform Operator; Approver selon le contexte.

## 7. Conditions d’entrée
Tenant, environnement, version, owner, cible, période, permission, restrictions et return origin sont explicites. Les preuves de CAP-INV-401..417 sont référencées sans duplication. Toute dépendance absente conduit à un état incomplet, partiel, bloqué ou inconnu ; aucune version, santé, efficacité, autorité ou ground truth n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | --- | --- | --- |
| Release Candidate / Readiness | CAP-INV-418/419 | candidate version and authorized scope | oui | selected immutable version | plan blocked |
| Target and source projections | Settings / Endpoint / Shared | targets, runtime capabilities, source health | oui | current snapshot | target/source unavailable |
| Shadow authority and policy | Govern / Security | non-alerting permission, duration and stop criteria | oui | current Decision/Approval if required | authorization-required |
| Comparison baseline | Command runtime Detection projection | existing runtime version and metrics | non | time-bounded | assessment without baseline |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Detection Release Candidate | Investigate | candidate logic/version | lecture |
| Environment / Target / runtime capability | Settings / Endpoint | authorized observation scope | lecture |
| Runtime Detection baseline | Command | version and aggregate behavior | lecture projection |
| Decision / Approval / Response Run / Result | Govern | authority and execution projection | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Shadow Evaluation Plan | créer, modifier, retirer, superseder | Investigate concept | plan ≠ execution |
| Shadow Assessment | enregistrer, comparer, contester | Investigate concept | shadow match ≠ Signal |
| Action Request context | préparer | Govern | no silent execution |

## 11. Fonctionnalités
- sélectionner version, targets, période, sources, métriques et limites
- définir critères d’arrêt et autorisation requise
- observer matches, non-matches pertinents, erreurs, latency et volume
- comparer au runtime existant sans alerting
- interrompre, clôturer et préparer une recommandation
- conserver versions, sources, erreurs, résultats partiels, restrictions, attribution et return origin
- fonctionner entièrement sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | --- | --- | --- | --- |
| Consulter, filtrer ou comparer | Detection Engineer | Shadow Evaluation Plan and Shadow Assessment | 0 | lecture autorisée et scope explicite | projection sourcée | non |
| Exécuter une assessment ou comparaison bornée | Detection Engineer | Tool Call / Assessment Result | 1 | déclenchement explicite, sources, paramètres et permission | résultat attribué, partial ou complete | policy applicable |
| Créer, modifier, contester ou retirer | Detection Engineer | Shadow Evaluation Plan and Shadow Assessment | 2 | mutation réversible et versionnée | nouvelle version et disposition humaine | OPEN-013 |
| Préparer une demande gouvernée | Detection Engineer | Change/Action Request context | 2 | cible, effet, risque, rollback et autorité visibles | package non effectif | destination Govern |
| Exécuter le changement réel | Govern / Platform Operator | runtime target / Response Run | 3 | Decision/Approval et policy applicables | projection du Result seulement dans Investigate | Govern owner |
| Détruire historique ou provenance | aucun rôle local | historical content | 4 | interdit par défaut | action refusée | strict governance |

Les classes 3 et 4 ne sont jamais exécutées par Investigate. Une demande ou coordination locale de classe 2 ne vaut ni Decision, ni Approval, ni exécution.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | --- | --- | --- | --- | --- |
| Construire ou examiner Shadow Evaluation Plan and Shadow Assessment | oui | checklists, catalogues et contrôles explicables | oui | proposition attribuée | formulaire, tables, matrices et revue humaine |
| Comparer ou valider Shadow Evaluation Plan and Shadow Assessment | oui | comparateur et règles explicables | oui | explication facultative | diff, diagnostics et checklist |
| Résumer risques, erreurs ou contradictions | oui | catalogue et agrégations déterministes | oui | résumé sourcé | sources brutes, timeline et revue manuelle |
| Approuver, promouvoir, activer, désactiver ou rollback | non localement | non | non | interdit | Govern et owner runtime selon autorité |

Toute proposition automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet. Aucun target, état, approbation, suppression, exception ou rollback n’est sélectionné silencieusement.

## 14. États fonctionnels
`draft`, `authorization-required`, `scheduled`, `observing`, `partial`, `completed`, `failed`, `cancelled`, `data-gap`, `permission-blocked`, `superseded`. Ces états sont des projections fonctionnelles ; ils ne redéfinissent ni les machines d’état Govern, Command, Settings ou runtime, ni une machine d’état objet définitive.

## 15. États d’interface
Loading conserve version, target, phase et return origin. Empty distingue absence de résultat et absence d’accès. Partial expose chaque target et résultat utilisable. Error conserve les données valides et le correlation ID. Offline interdit les mutations non garanties. Permission denied masque les données protégées. Stale sépare dernière observation connue et état actuel. Les conflits de version fournissent diff et reprise sûre. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Shadow Evaluation Plan | observation plan | CAP-INV-420/Govern | version, targets, period, metrics and stop criteria |
| Shadow Assessment | runtime non-alerting assessment | CAP-INV-418/419/427 | matches, errors, volume, latency and limitations |
| Recommendation | human disposition | Review / Project | no automatic activation |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-418/419 | candidate suitable for shadow | CAP-INV-422 | version, targets, sources, metrics, limitations and authority needs | Review/Readiness |
| CAP-INV-422 | authorization required | CAP-INV-420/Govern | non-alerting scope, duration, stop and rollback context | Shadow Plan |
| Govern/Settings runtime | shadow result available | CAP-INV-422 | matches, errors, latency, volume, target results and provenance | Govern/Run |
| CAP-INV-422 | assessment complete | CAP-INV-418/419/429 | recommendation, limitations and changes proposed | Shadow Assessment |

Chaque transition conserve owner, tenant, environnement, target, version, permissions, restrictions, erreurs, autorité, provenance et return origin. Elle n’accorde aucun accès ni pouvoir supplémentaire.

## 18. Dépendances
CAP-INV-418..421,424..429; Command runtime projection; Settings/Endpoint; Govern; Shared Metrics/Trace. OPEN-017 gouverne le futur choix runtime/langage/portabilité sans option sélectionnée. Shared Background Jobs, Notifications, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Collaboration, Comments, Assignments, Comparison, Inspector, Context Bar, Audit Hooks et Recovery sont consommés sans redéfinition.

## 19. Source de vérité
Investigate est source de vérité de Shadow Evaluation Plan and Shadow Assessment comme concept fonctionnel. Command reste source de runtime Detection, Signal, Alert et Incident. Settings reste source des environments, targets, runtimes configurés et health administratif. Govern reste source des Action Requests, Decisions, Approvals, Response Runs et Results. Studio reste source des Tools et Automation Runs. Une projection locale ne remplace jamais son objet canonique propriétaire.

## 20. Provenance et audit
Enregistrer le besoin initial, Project, Hypothesis, Drafts, versions, Review Packages, Release Candidates, reviewers, readiness, targets, plans, Action Requests, Decisions, Approvals, Response Runs, Results, runtime observations, health, Signals/Alerts/Incidents liés, assessments, propositions, erreurs, interruptions, auteurs, timestamps, Tool Calls, Automation Runs, paramètres, dispositions humaines, exports et correlation IDs applicables à Shadow Evaluation Plan and Shadow Assessment. Aucune trace n’est supprimée ni réécrite silencieusement.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Shadow evaluation prepare/request/read | runtime observation | 1/2 | event details permission-bound | step-up possible | requester/approver | Investigate/Govern | Permissions |
| Shadow results read/compare | sensitive production behavior | 0/1 | aggregate or mask fields | possible | observer/data owner | Command/Settings | Permissions |
| Shadow cancel request | operational impact | 2 request / 3 execution | no hidden stop | authority required | requester/operator | Govern/Settings | Permissions |

Les permissions atomiques, namespaces, RBAC/ABAC, step-up définitif et séparation des tâches finale restent reportés. Toute exécution réelle de classe 3 ou 4 reste chez Govern et l’owner runtime.

## 22. Limites et erreurs
- Shadow evaluation ≠ historical replay and ≠ activation.
- Shadow match ≠ Signal.
- No alerting authority is implied.
- Absence of errors ≠ production readiness certainty.
- Les données peuvent être stale, partielles, restreintes, incohérentes entre tenants ou indisponibles.
- Un Result, score, match, non-match, health state, disposition Command ou sortie IA ne vaut pas conclusion universelle.
- Timeout, cancellation, revocation, target offline, version superseded et partial result restent attribués et visibles.

## 23. Métriques conceptuelles
- shadow matches/non-matches by target/version
- error, data-gap, latency and volume observations
- comparison with runtime baseline
- Signals created by shadow — target zero
- complétude de provenance, versions, targets, autorité et dispositions humaines
- nombre d’exécutions silencieuses, d’auto-approbations ou de suppressions de traces — cible conceptuelle zéro

Aucun seuil universel de latency, coût, précision, recall, health ou business value n’est imposé.

## 24. Classification de livraison
`defined` / `planned`. Preuve documentaire uniquement. Aucun document n’est `validated`, `implemented`, `deployed`, `active`, `native` ou `integrated`. Aucun moteur, langage, syntaxe, produit tiers, modèle ML, API, protocole, commande ou code n’est choisi.

## 25. Critères d’acceptation
### 1. Shadow evaluation
**Given** une version autorisée en shadow sans alerting  
**When** des matches sont observés  
**Then** aucun Signal opérationnel n’est créé, les matches restent shadow et aucune activation automatique n’a lieu

### 2. Shadow data gap
**Given** une source becomes unavailable during observation  
**When** l’assessment est clôturé  
**Then** le résultat reste partial/data-gap et la recommandation expose la limite

### 3. Sans IA
**Given** aucun modèle  
**When** la shadow evaluation est préparée et revue  
**Then** forms, metrics tables, comparison and human review suffice

## 26. Questions ouvertes
- OPEN-008 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- OPEN-015 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- OPEN-017 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- Les schémas, cardinalités, machines d’état, formats de déploiement, permissions atomiques, contrats techniques et composition détaillée des écrans restent futurs.
- OPEN-005 demeure forensic-only et n’est pas utilisée pour Detection.

## 27. Consommateurs documentaires
Detection Engineering lifecycle and Capability Maps ; Command runtime Detection/Signal/Alert/Incident ; Platform Settings environments/targets/health ; Endpoint Agent capabilities/version/health ; Govern Action Request/Decision/Approval/Response Run/Result ; Studio Tools/Runs ; Shared mechanisms ; phases Objects, Permissions, Screens, Journeys, Technique, validation et future 4B.3B uniquement comme handoff non canonique.
