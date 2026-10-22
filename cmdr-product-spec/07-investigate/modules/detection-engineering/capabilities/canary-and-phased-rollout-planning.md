---
id: CAP-INV-423
title: Canary and Phased Rollout Planning
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
# CAP-INV-423 — Canary and Phased Rollout Planning

## 1. Définition
Définir une population canary et un rollout progressif avec critères d’entrée, succès, arrêt, validations intermédiaires et rollback, puis analyser les résultats sans conclure qu’un succès partiel vaut succès global.

## 2. Problème utilisateur
Un canary réussi sur une population limitée peut être généralisé à tort, tandis qu’un target en erreur peut être masqué par un statut agrégé.

## 3. Objectifs
- définir population, environments, tenants, durée, étapes et pourcentages conceptuels
- définir critères d’entrée, succès, arrêt et limites
- définir validations intermédiaires, rollback, owner et observers
- comparer plusieurs stratégies et transmettre pour approbation
- recevoir les résultats et préparer l’étape suivante ou l’arrêt

## 4. Non-objectifs
- ne pas définir API, protocole, compilateur, parser, AST, format de package, stockage, pipeline, streaming, commande ou code produit
- ne sélectionner aucun moteur, langage, représentation cible, syntaxe vendor, produit tiers ou modèle ML
- ne réaliser aucun déploiement, activation, désactivation, rollback, suppression ou exception réelle pendant la phase documentaire
- ne modifier ni supprimer silencieusement runtime Detection, Signal, Alert ou Incident
- ne créer aucune capability CAP-INV-5xx, aucun objet Threat Intelligence, Cloud Analysis, Mobile Forensics ou écran détaillé

## 5. Propriétaire
Investigate / Detection Engineering possède **Canary Plan and Canary Assessment** et ses dispositions fonctionnelles. Command conserve runtime Detection, Signal, Alert, Incident et le feedback opérationnel. Platform Settings conserve runtimes configurés, targets, environments, tenants, sources, parsers, schemas, health, providers, secrets et canaux administratifs. Endpoint Agent conserve ses capacités, versions, health et exécutions locales autorisées. Govern conserve Action Request, Decision, Approval, Response Run, Result et toute autorité de classe 3 ou 4. Studio conserve Tool, Tool Call, Workflow, Automation Agent, Human Gate et Automation Run. Shared conserve Jobs, Notifications, Trace, Activity, Linking, Versioning, Comparison, Reporting, Export, Collaboration, Audit Hooks et Recovery.

## 6. Utilisateurs
Principal : **Detection Owner**. Secondaires : Detection Engineer; Detection Reviewer; Detection Owner; Incident Commander; Platform Operator; Approver selon le contexte.

## 7. Conditions d’entrée
Tenant, environnement, version, owner, cible, période, permission, restrictions et return origin sont explicites. Les preuves de CAP-INV-401..417 sont référencées sans duplication. Toute dépendance absente conduit à un état incomplet, partiel, bloqué ou inconnu ; aucune version, santé, efficacité, autorité ou ground truth n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | --- | --- | --- |
| Release Candidate / Readiness | CAP-INV-418/419 | candidate and target conditions | oui | current versions | plan incomplete |
| Promotion Plan | CAP-INV-421 | target order and environment differences | oui | selected plan | population undefined |
| Shadow Assessment | CAP-INV-422 | pre-canary observations and limits | non | selected assessment | risk remains explicit |
| Authority/target projections | Govern / Settings | approval, targets, health and rollout support | oui | current snapshots | request blocked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Release Candidate / Promotion Plan | Investigate | candidate and proposed sequence | lecture |
| Environment / Target / Fleet projection | Settings / Endpoint | population and health | lecture |
| Decision / Approval / Response Run / Result | Govern | authority and execution results | lecture/lien |
| Runtime Detection / health projection | Command / Settings | canary result context | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Canary Plan | créer, modifier, comparer, retirer, superseder | Investigate concept | plan ≠ rollout execution |
| Canary Assessment | enregistrer, comparer, contester | Investigate concept | partial target results preserved |
| Action Request context | préparer | Govern | no automatic stage advance |

## 11. Fonctionnalités
- définir population, environments, tenants, durée, étapes et pourcentages conceptuels
- définir critères d’entrée, succès, arrêt et limites
- définir validations intermédiaires, rollback, owner et observers
- comparer plusieurs stratégies et transmettre pour approbation
- recevoir les résultats et préparer l’étape suivante ou l’arrêt
- conserver versions, sources, erreurs, résultats partiels, restrictions, attribution et return origin
- fonctionner entièrement sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | --- | --- | --- | --- |
| Consulter, filtrer ou comparer | Detection Owner | Canary Plan and Canary Assessment | 0 | lecture autorisée et scope explicite | projection sourcée | non |
| Exécuter une assessment ou comparaison bornée | Detection Owner | Tool Call / Assessment Result | 1 | déclenchement explicite, sources, paramètres et permission | résultat attribué, partial ou complete | policy applicable |
| Créer, modifier, contester ou retirer | Detection Owner | Canary Plan and Canary Assessment | 2 | mutation réversible et versionnée | nouvelle version et disposition humaine | OPEN-013 |
| Préparer une demande gouvernée | Detection Owner | Change/Action Request context | 2 | cible, effet, risque, rollback et autorité visibles | package non effectif | destination Govern |
| Exécuter le changement réel | Govern / Platform Operator | runtime target / Response Run | 3 | Decision/Approval et policy applicables | projection du Result seulement dans Investigate | Govern owner |
| Détruire historique ou provenance | aucun rôle local | historical content | 4 | interdit par défaut | action refusée | strict governance |

Les classes 3 et 4 ne sont jamais exécutées par Investigate. Une demande ou coordination locale de classe 2 ne vaut ni Decision, ni Approval, ni exécution.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | --- | --- | --- | --- | --- |
| Construire ou examiner Canary Plan and Canary Assessment | oui | checklists, catalogues et contrôles explicables | oui | proposition attribuée | formulaire, tables, matrices et revue humaine |
| Comparer ou valider Canary Plan and Canary Assessment | oui | comparateur et règles explicables | oui | explication facultative | diff, diagnostics et checklist |
| Résumer risques, erreurs ou contradictions | oui | catalogue et agrégations déterministes | oui | résumé sourcé | sources brutes, timeline et revue manuelle |
| Approuver, promouvoir, activer, désactiver ou rollback | non localement | non | non | interdit | Govern et owner runtime selon autorité |

Toute proposition automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet. Aucun target, état, approbation, suppression, exception ou rollback n’est sélectionné silencieusement.

## 14. États fonctionnels
`draft`, `incomplete`, `review-requested`, `approved-for-request`, `scheduled`, `running`, `partial`, `success-candidate`, `stop-triggered`, `failed`, `cancelled`, `superseded`. Ces états sont des projections fonctionnelles ; ils ne redéfinissent ni les machines d’état Govern, Command, Settings ou runtime, ni une machine d’état objet définitive.

## 15. États d’interface
Loading conserve version, target, phase et return origin. Empty distingue absence de résultat et absence d’accès. Partial expose chaque target et résultat utilisable. Error conserve les données valides et le correlation ID. Offline interdit les mutations non garanties. Permission denied masque les données protégées. Stale sépare dernière observation connue et état actuel. Les conflits de version fournissent diff et reprise sûre. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Canary Plan | phased plan | CAP-INV-420/Govern | population, stages, criteria and rollback |
| Canary Assessment | runtime assessment | Review / Promotion Decision | per-target results and limits |
| Next-stage recommendation | human proposal | CAP-INV-421/424/433 | not automatic promotion |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-421/422 | prepare canary | CAP-INV-423 | candidate, targets, shadow evidence, stages and criteria | Promotion/Shadow |
| CAP-INV-423 | submit for authority | CAP-INV-420/Govern | population, duration, stop, validation and rollback | Canary Plan |
| Govern Result | canary stage result | CAP-INV-423 | target results, errors, health, versions and stop triggers | Govern |
| CAP-INV-423 | review result | CAP-INV-424/433/421 | next-stage, stop or rollback recommendation | Canary Assessment |

Chaque transition conserve owner, tenant, environnement, target, version, permissions, restrictions, erreurs, autorité, provenance et return origin. Elle n’accorde aucun accès ni pouvoir supplémentaire.

## 18. Dépendances
CAP-INV-418..424,433; Settings targets/Fleet; Endpoint; Govern; Shared Jobs/Metrics/Trace. OPEN-017 gouverne le futur choix runtime/langage/portabilité sans option sélectionnée. Shared Background Jobs, Notifications, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Collaboration, Comments, Assignments, Comparison, Inspector, Context Bar, Audit Hooks et Recovery sont consommés sans redéfinition.

## 19. Source de vérité
Investigate est source de vérité de Canary Plan and Canary Assessment comme concept fonctionnel. Command reste source de runtime Detection, Signal, Alert et Incident. Settings reste source des environments, targets, runtimes configurés et health administratif. Govern reste source des Action Requests, Decisions, Approvals, Response Runs et Results. Studio reste source des Tools et Automation Runs. Une projection locale ne remplace jamais son objet canonique propriétaire.

## 20. Provenance et audit
Enregistrer le besoin initial, Project, Hypothesis, Drafts, versions, Review Packages, Release Candidates, reviewers, readiness, targets, plans, Action Requests, Decisions, Approvals, Response Runs, Results, runtime observations, health, Signals/Alerts/Incidents liés, assessments, propositions, erreurs, interruptions, auteurs, timestamps, Tool Calls, Automation Runs, paramètres, dispositions humaines, exports et correlation IDs applicables à Canary Plan and Canary Assessment. Aucune trace n’est supprimée ni réécrite silencieusement.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Canary plan prepare/request/read | phased production risk | 2 | target populations scoped | step-up likely | planner/approver/operator | Investigate/Govern | Permissions |
| Canary assessment create/update | quality interpretation | 1/2 | target results permission-bound | possible | observer/reviewer | Investigate | Permissions |
| Advance/stop/rollback request | production impact | 2 request / 3 execution | criteria visible | authority required | requester/operator | Govern | Permissions |

Les permissions atomiques, namespaces, RBAC/ABAC, step-up définitif et séparation des tâches finale restent reportés. Toute exécution réelle de classe 3 ou 4 reste chez Govern et l’owner runtime.

## 22. Limites et erreurs
- Canary target ≠ full population.
- Canary success ≠ global success.
- No canary error ≠ no production risk.
- Percentages are functional planning, not engine configuration.
- Les données peuvent être stale, partielles, restreintes, incohérentes entre tenants ou indisponibles.
- Un Result, score, match, non-match, health state, disposition Command ou sortie IA ne vaut pas conclusion universelle.
- Timeout, cancellation, revocation, target offline, version superseded et partial result restent attribués et visibles.

## 23. Métriques conceptuelles
- canary stages and target outcomes
- stop triggers and partial results
- time to human disposition
- automatic stage advances — target zero
- complétude de provenance, versions, targets, autorité et dispositions humaines
- nombre d’exécutions silencieuses, d’auto-approbations ou de suppressions de traces — cible conceptuelle zéro

Aucun seuil universel de latency, coût, précision, recall, health ou business value n’est imposé.

## 24. Classification de livraison
`defined` / `planned`. Preuve documentaire uniquement. Aucun document n’est `validated`, `implemented`, `deployed`, `active`, `native` ou `integrated`. Aucun moteur, langage, syntaxe, produit tiers, modèle ML, API, protocole, commande ou code n’est choisi.

## 25. Critères d’acceptation
### 1. Canary partiel
**Given** un canary multi-target avec un target en erreur  
**When** le résultat est examiné  
**Then** le résultat reste partial, le target est visible, aucune promotion globale automatique et un rollback peut être préparé

### 2. Canary success candidate
**Given** tous les targets du canary réussissent  
**When** la prochaine étape est proposée  
**Then** la proposition reste soumise aux conditions/authority et n’est pas global success

### 3. Sans IA
**Given** aucun modèle  
**When** le canary est planifié  
**Then** matrices, target tables, explicit criteria and human review suffice

## 26. Questions ouvertes
- OPEN-008 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- OPEN-015 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- OPEN-017 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- Les schémas, cardinalités, machines d’état, formats de déploiement, permissions atomiques, contrats techniques et composition détaillée des écrans restent futurs.
- OPEN-005 demeure forensic-only et n’est pas utilisée pour Detection.

## 27. Consommateurs documentaires
Detection Engineering lifecycle and Capability Maps ; Command runtime Detection/Signal/Alert/Incident ; Platform Settings environments/targets/health ; Endpoint Agent capabilities/version/health ; Govern Action Request/Decision/Approval/Response Run/Result ; Studio Tools/Runs ; Shared mechanisms ; phases Objects, Permissions, Screens, Journeys, Technique, validation et future 4B.3B uniquement comme handoff non canonique.
