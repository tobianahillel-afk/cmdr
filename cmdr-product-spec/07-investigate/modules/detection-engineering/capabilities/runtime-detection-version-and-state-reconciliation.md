---
id: CAP-INV-425
title: Runtime Detection Version and State Reconciliation
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
# CAP-INV-425 — Runtime Detection Version and State Reconciliation

## 1. Définition
Comparer la version et l’état attendus avec les observations runtime par target, détecter divergences, stale or partial application et préparer une remediation request sans écrasement silencieux.

## 2. Problème utilisateur
Une configuration acceptée ou copiée peut ne pas être appliquée, et une vue agrégée peut masquer des versions différentes entre targets.

## 3. Objectifs
- voir expected, observed, active, partial and stale versions by target
- voir configuration acceptance/application status, timestamps, sources and errors
- comparer et contester observations
- ouvrir une divergence investigation
- préparer remediation request et conserver l’historique

## 4. Non-objectifs
- ne pas définir API, protocole, compilateur, parser, AST, format de package, stockage, pipeline, streaming, commande ou code produit
- ne sélectionner aucun moteur, langage, représentation cible, syntaxe vendor, produit tiers ou modèle ML
- ne réaliser aucun déploiement, activation, désactivation, rollback, suppression ou exception réelle pendant la phase documentaire
- ne modifier ni supprimer silencieusement runtime Detection, Signal, Alert ou Incident
- ne créer aucune capability CAP-INV-5xx, aucun objet Threat Intelligence, Cloud Analysis, Mobile Forensics ou écran détaillé

## 5. Propriétaire
Investigate / Detection Engineering possède **Runtime Version Observation and Reconciliation Assessment** et ses dispositions fonctionnelles. Command conserve runtime Detection, Signal, Alert, Incident et le feedback opérationnel. Platform Settings conserve runtimes configurés, targets, environments, tenants, sources, parsers, schemas, health, providers, secrets et canaux administratifs. Endpoint Agent conserve ses capacités, versions, health et exécutions locales autorisées. Govern conserve Action Request, Decision, Approval, Response Run, Result et toute autorité de classe 3 ou 4. Studio conserve Tool, Tool Call, Workflow, Automation Agent, Human Gate et Automation Run. Shared conserve Jobs, Notifications, Trace, Activity, Linking, Versioning, Comparison, Reporting, Export, Collaboration, Audit Hooks et Recovery.

## 6. Utilisateurs
Principal : **Detection Owner**. Secondaires : Detection Engineer; Detection Reviewer; Detection Owner; Incident Commander; Platform Operator; Approver selon le contexte.

## 7. Conditions d’entrée
Tenant, environnement, version, owner, cible, période, permission, restrictions et return origin sont explicites. Les preuves de CAP-INV-401..417 sont référencées sans duplication. Toute dépendance absente conduit à un état incomplet, partiel, bloqué ou inconnu ; aucune version, santé, efficacité, autorité ou ground truth n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | --- | --- | --- |
| Expected candidate/version state | CAP-INV-418/421/424 | candidate version, targets and intended state | oui | selected change | expected unknown |
| Observed runtime version/state | Command / Settings / Endpoint | version, application, activation and timestamps | oui | current observation | unknown/unavailable |
| Response Run / Result | Govern | reported effects and target outcomes | non | linked run | reconciliation without canonical Result |
| Health/source context | CAP-INV-426 / Settings | freshness, errors and capability | oui | current | observation confidence limited |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Release Candidate / Promotion Plan | Investigate | expected version and scope | lecture |
| Runtime Detection | Command | runtime identity, version and state | lecture |
| Target configuration/version observation | Settings / Endpoint | applied version and timestamp | lecture |
| Response Run / Result | Govern | change result and verification | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Runtime Version Reconciliation Assessment | créer, comparer, contester, superseder | Investigate concept | does not change runtime |
| Divergence investigation | ouvrir/link | Investigate | expected and observed preserved |
| Remediation request context | préparer | Govern/Settings | no overwrite |

## 11. Fonctionnalités
- voir expected, observed, active, partial and stale versions by target
- voir configuration acceptance/application status, timestamps, sources and errors
- comparer et contester observations
- ouvrir une divergence investigation
- préparer remediation request et conserver l’historique
- conserver versions, sources, erreurs, résultats partiels, restrictions, attribution et return origin
- fonctionner entièrement sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | --- | --- | --- | --- |
| Consulter, filtrer ou comparer | Detection Owner | Runtime Version Observation and Reconciliation Assessment | 0 | lecture autorisée et scope explicite | projection sourcée | non |
| Exécuter une assessment ou comparaison bornée | Detection Owner | Tool Call / Assessment Result | 1 | déclenchement explicite, sources, paramètres et permission | résultat attribué, partial ou complete | policy applicable |
| Créer, modifier, contester ou retirer | Detection Owner | Runtime Version Observation and Reconciliation Assessment | 2 | mutation réversible et versionnée | nouvelle version et disposition humaine | OPEN-013 |
| Préparer une demande gouvernée | Detection Owner | Change/Action Request context | 2 | cible, effet, risque, rollback et autorité visibles | package non effectif | destination Govern |
| Exécuter le changement réel | Govern / Platform Operator | runtime target / Response Run | 3 | Decision/Approval et policy applicables | projection du Result seulement dans Investigate | Govern owner |
| Détruire historique ou provenance | aucun rôle local | historical content | 4 | interdit par défaut | action refusée | strict governance |

Les classes 3 et 4 ne sont jamais exécutées par Investigate. Une demande ou coordination locale de classe 2 ne vaut ni Decision, ni Approval, ni exécution.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | --- | --- | --- | --- | --- |
| Construire ou examiner Runtime Version Observation and Reconciliation Assessment | oui | checklists, catalogues et contrôles explicables | oui | proposition attribuée | formulaire, tables, matrices et revue humaine |
| Comparer ou valider Runtime Version Observation and Reconciliation Assessment | oui | comparateur et règles explicables | oui | explication facultative | diff, diagnostics et checklist |
| Résumer risques, erreurs ou contradictions | oui | catalogue et agrégations déterministes | oui | résumé sourcé | sources brutes, timeline et revue manuelle |
| Approuver, promouvoir, activer, désactiver ou rollback | non localement | non | non | interdit | Govern et owner runtime selon autorité |

Toute proposition automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet. Aucun target, état, approbation, suppression, exception ou rollback n’est sélectionné silencieusement.

## 14. États fonctionnels
`expected`, `observed`, `active`, `inactive`, `partially-active`, `stale`, `divergent`, `unknown`, `unavailable`, `rollback-observed`, `superseded`. Ces états sont des projections fonctionnelles ; ils ne redéfinissent ni les machines d’état Govern, Command, Settings ou runtime, ni une machine d’état objet définitive.

## 15. États d’interface
Loading conserve version, target, phase et return origin. Empty distingue absence de résultat et absence d’accès. Partial expose chaque target et résultat utilisable. Error conserve les données valides et le correlation ID. Offline interdit les mutations non garanties. Permission denied masque les données protégées. Stale sépare dernière observation connue et état actuel. Les conflits de version fournissent diff et reprise sûre. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Version/state reconciliation | assessment | CAP-INV-426/433/434 | per-target expected/observed differences |
| Divergence case | investigation context | Project / Settings/Govern request | sources, timestamps and errors visible |
| Reconciliation history | versioned record | Audit / Lifecycle provenance | no observation overwritten |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-424 | deployment result or runtime observation | CAP-INV-425 | expected version, targets, reported result and timestamps | Coordination |
| Runtime owner | state/version update | CAP-INV-425 | observed version, source, timestamp, freshness and error | Runtime owner |
| CAP-INV-425 | divergent or stale | CAP-INV-420/433 / investigation | expected/observed, target, risk and remediation context | Reconciliation |
| CAP-INV-425 | reconciled | CAP-INV-426 | runtime version, state, freshness and target scope | Reconciliation |

Chaque transition conserve owner, tenant, environnement, target, version, permissions, restrictions, erreurs, autorité, provenance et return origin. Elle n’accorde aucun accès ni pouvoir supplémentaire.

## 18. Dépendances
CAP-INV-418,421,424,426,433,434; Command Detection; Settings/Endpoint version projections; Govern Result. OPEN-017 gouverne le futur choix runtime/langage/portabilité sans option sélectionnée. Shared Background Jobs, Notifications, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Collaboration, Comments, Assignments, Comparison, Inspector, Context Bar, Audit Hooks et Recovery sont consommés sans redéfinition.

## 19. Source de vérité
Investigate est source de vérité de Runtime Version Observation and Reconciliation Assessment comme concept fonctionnel. Command reste source de runtime Detection, Signal, Alert et Incident. Settings reste source des environments, targets, runtimes configurés et health administratif. Govern reste source des Action Requests, Decisions, Approvals, Response Runs et Results. Studio reste source des Tools et Automation Runs. Une projection locale ne remplace jamais son objet canonique propriétaire.

## 20. Provenance et audit
Enregistrer le besoin initial, Project, Hypothesis, Drafts, versions, Review Packages, Release Candidates, reviewers, readiness, targets, plans, Action Requests, Decisions, Approvals, Response Runs, Results, runtime observations, health, Signals/Alerts/Incidents liés, assessments, propositions, erreurs, interruptions, auteurs, timestamps, Tool Calls, Automation Runs, paramètres, dispositions humaines, exports et correlation IDs applicables à Runtime Version Observation and Reconciliation Assessment. Aucune trace n’est supprimée ni réécrite silencieusement.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Runtime Detection/version read | production state | 0 | target/tenant scoped | possible | viewer/runtime admin | Command/Settings | Permissions |
| Reconciliation create/update/dispute | quality assessment | 2 | no hidden overwrite | OPEN-013 | assessor/reviewer | Investigate | Permissions |
| Remediation request prepare | production impact | 2 | expected/observed explicit | step-up possible | requester/operator | Govern/Settings | Permissions |

Les permissions atomiques, namespaces, RBAC/ABAC, step-up définitif et séparation des tâches finale restent reportés. Toute exécution réelle de classe 3 ou 4 reste chez Govern et l’owner runtime.

## 22. Limites et erreurs
- Configuration accepted ≠ configuration applied.
- Expected version ≠ observed version.
- Observed active ≠ healthy or effective.
- No automatic overwrite or remediation.
- Les données peuvent être stale, partielles, restreintes, incohérentes entre tenants ou indisponibles.
- Un Result, score, match, non-match, health state, disposition Command ou sortie IA ne vaut pas conclusion universelle.
- Timeout, cancellation, revocation, target offline, version superseded et partial result restent attribués et visibles.

## 23. Métriques conceptuelles
- targets reconciled/divergent/stale/unknown
- observation freshness
- time to divergence disposition
- silent overwrites — target zero
- complétude de provenance, versions, targets, autorité et dispositions humaines
- nombre d’exécutions silencieuses, d’auto-approbations ou de suppressions de traces — cible conceptuelle zéro

Aucun seuil universel de latency, coût, précision, recall, health ou business value n’est imposé.

## 24. Classification de livraison
`defined` / `planned`. Preuve documentaire uniquement. Aucun document n’est `validated`, `implemented`, `deployed`, `active`, `native` ou `integrated`. Aucun moteur, langage, syntaxe, produit tiers, modèle ML, API, protocole, commande ou code n’est choisi.

## 25. Critères d’acceptation
### 1. Version divergente
**Given** une expected version différente de l’observed target version  
**When** la reconciliation s’exécute  
**Then** le target est divergent, les deux versions restent visibles et une remediation request peut être préparée

### 2. Partial versions
**Given** une population affiche plusieurs versions  
**When** l’état global est consulté  
**Then** partially-active remains and no single active version is claimed

### 3. Sans IA
**Given** aucun modèle  
**When** la reconciliation est réalisée  
**Then** deterministic comparison, timestamps, diff and human disposition suffice

## 26. Questions ouvertes
- OPEN-008 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- OPEN-017 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- Les schémas, cardinalités, machines d’état, formats de déploiement, permissions atomiques, contrats techniques et composition détaillée des écrans restent futurs.
- OPEN-005 demeure forensic-only et n’est pas utilisée pour Detection.

## 27. Consommateurs documentaires
Detection Engineering lifecycle and Capability Maps ; Command runtime Detection/Signal/Alert/Incident ; Platform Settings environments/targets/health ; Endpoint Agent capabilities/version/health ; Govern Action Request/Decision/Approval/Response Run/Result ; Studio Tools/Runs ; Shared mechanisms ; phases Objects, Permissions, Screens, Journeys, Technique, validation et future 4B.3B uniquement comme handoff non canonique.
