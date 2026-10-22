---
id: CAP-INV-421
title: Environment Promotion Planning and Target Selection
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
# CAP-INV-421 — Environment Promotion Planning and Target Selection

## 1. Définition
Comparer les environnements, tenants, runtimes et targets administrés par Settings, sélectionner des cibles candidates et préparer un ordre de promotion sans administrer ni exécuter.

## 2. Problème utilisateur
Une target visible peut être sélectionnée sans être compatible, saine ou autorisée, et une sélection peut être prise pour une promotion réalisée.

## 3. Objectifs
- voir environments, tenants, runtimes, versions, capacities, health, restrictions and policies
- comparer data sources, schemas and environment differences
- sélectionner des targets candidates sans administrer Settings
- définir order, conditions, expected validations and return path
- annoter, contester et préparer le handoff

## 4. Non-objectifs
- ne pas définir API, protocole, compilateur, parser, AST, format de package, stockage, pipeline, streaming, commande ou code produit
- ne sélectionner aucun moteur, langage, représentation cible, syntaxe vendor, produit tiers ou modèle ML
- ne réaliser aucun déploiement, activation, désactivation, rollback, suppression ou exception réelle pendant la phase documentaire
- ne modifier ni supprimer silencieusement runtime Detection, Signal, Alert ou Incident
- ne créer aucune capability CAP-INV-5xx, aucun objet Threat Intelligence, Cloud Analysis, Mobile Forensics ou écran détaillé

## 5. Propriétaire
Investigate / Detection Engineering possède **Promotion Plan** et ses dispositions fonctionnelles. Command conserve runtime Detection, Signal, Alert, Incident et le feedback opérationnel. Platform Settings conserve runtimes configurés, targets, environments, tenants, sources, parsers, schemas, health, providers, secrets et canaux administratifs. Endpoint Agent conserve ses capacités, versions, health et exécutions locales autorisées. Govern conserve Action Request, Decision, Approval, Response Run, Result et toute autorité de classe 3 ou 4. Studio conserve Tool, Tool Call, Workflow, Automation Agent, Human Gate et Automation Run. Shared conserve Jobs, Notifications, Trace, Activity, Linking, Versioning, Comparison, Reporting, Export, Collaboration, Audit Hooks et Recovery.

## 6. Utilisateurs
Principal : **Detection Engineer**. Secondaires : Detection Engineer; Detection Reviewer; Detection Owner; Incident Commander; Platform Operator; Approver selon le contexte.

## 7. Conditions d’entrée
Tenant, environnement, version, owner, cible, période, permission, restrictions et return origin sont explicites. Les preuves de CAP-INV-401..417 sont référencées sans duplication. Toute dépendance absente conduit à un état incomplet, partiel, bloqué ou inconnu ; aucune version, santé, efficacité, autorité ou ground truth n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | --- | --- | --- |
| Readiness Assessment | CAP-INV-419 | candidate and target conditions | oui | current assessment | plan incomplete |
| Environment/tenant/target projections | Platform Settings | identity, runtime, version, capability, health and policy | oui | current snapshot | target unavailable |
| Data/schema differences | CAP-INV-404/405 / Settings | sources, fields and compatibility limits | oui | linked versions | conditions missing |
| Return and rollback context | CAP-INV-433 / Experience | return path and recovery constraints | oui | plan version | handoff blocked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Environment / Tenant / Deployment Target | Platform Settings | scope, labels, capabilities and health | lecture/sélection projection |
| Runtime capability projection | Settings / Endpoint Agent | declared support and version | lecture |
| Release Candidate / Readiness | Investigate | candidate and constraints | lecture |
| Policies / authority projection | Govern / Security | promotion constraints | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Promotion Plan | créer, modifier, comparer, contester, superseder | Investigate concept | plan ≠ execution |
| Target selection set | proposer/versionner | Investigate using Settings refs | selected ≠ compatible |
| Change Request context | préparer | CAP-INV-420 | no target administration |

## 11. Fonctionnalités
- voir environments, tenants, runtimes, versions, capacities, health, restrictions and policies
- comparer data sources, schemas and environment differences
- sélectionner des targets candidates sans administrer Settings
- définir order, conditions, expected validations and return path
- annoter, contester et préparer le handoff
- conserver versions, sources, erreurs, résultats partiels, restrictions, attribution et return origin
- fonctionner entièrement sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | --- | --- | --- | --- |
| Consulter, filtrer ou comparer | Detection Engineer | Promotion Plan | 0 | lecture autorisée et scope explicite | projection sourcée | non |
| Exécuter une assessment ou comparaison bornée | Detection Engineer | Tool Call / Assessment Result | 1 | déclenchement explicite, sources, paramètres et permission | résultat attribué, partial ou complete | policy applicable |
| Créer, modifier, contester ou retirer | Detection Engineer | Promotion Plan | 2 | mutation réversible et versionnée | nouvelle version et disposition humaine | OPEN-013 |
| Préparer une demande gouvernée | Detection Engineer | Change/Action Request context | 2 | cible, effet, risque, rollback et autorité visibles | package non effectif | destination Govern |
| Exécuter le changement réel | Govern / Platform Operator | runtime target / Response Run | 3 | Decision/Approval et policy applicables | projection du Result seulement dans Investigate | Govern owner |
| Détruire historique ou provenance | aucun rôle local | historical content | 4 | interdit par défaut | action refusée | strict governance |

Les classes 3 et 4 ne sont jamais exécutées par Investigate. Une demande ou coordination locale de classe 2 ne vaut ni Decision, ni Approval, ni exécution.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | --- | --- | --- | --- | --- |
| Construire ou examiner Promotion Plan | oui | checklists, catalogues et contrôles explicables | oui | proposition attribuée | formulaire, tables, matrices et revue humaine |
| Comparer ou valider Promotion Plan | oui | comparateur et règles explicables | oui | explication facultative | diff, diagnostics et checklist |
| Résumer risques, erreurs ou contradictions | oui | catalogue et agrégations déterministes | oui | résumé sourcé | sources brutes, timeline et revue manuelle |
| Approuver, promouvoir, activer, désactiver ou rollback | non localement | non | non | interdit | Govern et owner runtime selon autorité |

Toute proposition automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet. Aucun target, état, approbation, suppression, exception ou rollback n’est sélectionné silencieusement.

## 14. États fonctionnels
`draft`, `incomplete`, `targets-proposed`, `targets-reviewed`, `ready-for-request`, `blocked`, `incompatible-target`, `disputed`, `superseded`, `withdrawn`. Ces états sont des projections fonctionnelles ; ils ne redéfinissent ni les machines d’état Govern, Command, Settings ou runtime, ni une machine d’état objet définitive.

## 15. États d’interface
Loading conserve version, target, phase et return origin. Empty distingue absence de résultat et absence d’accès. Partial expose chaque target et résultat utilisable. Error conserve les données valides et le correlation ID. Offline interdit les mutations non garanties. Permission denied masque les données protégées. Stale sépare dernière observation connue et état actuel. Les conflits de version fournissent diff et reprise sûre. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Promotion Plan | functional plan | CAP-INV-420/423/424 | targets, order, conditions and return visible |
| Target comparison | comparison result | Readiness / Reviewer | differences explicit |
| Settings information request | request context | Platform Settings | no configuration mutation |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-419 | targets ready/conditional | CAP-INV-421 | candidate, target assessments and restrictions | Readiness |
| CAP-INV-421 | plan ready | CAP-INV-420/423 | targets, order, conditions, validations and rollback | Promotion Plan |
| CAP-INV-421 | target incompatible/stale | CAP-INV-419 / Settings request | target, reason, version and impact | Promotion Plan |

Chaque transition conserve owner, tenant, environnement, target, version, permissions, restrictions, erreurs, autorité, provenance et return origin. Elle n’accorde aucun accès ni pouvoir supplémentaire.

## 18. Dépendances
CAP-INV-404,405,418,419,420,423,433; Settings Environment/Tenant/Health; Endpoint. OPEN-017 gouverne le futur choix runtime/langage/portabilité sans option sélectionnée. Shared Background Jobs, Notifications, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Collaboration, Comments, Assignments, Comparison, Inspector, Context Bar, Audit Hooks et Recovery sont consommés sans redéfinition.

## 19. Source de vérité
Investigate est source de vérité de Promotion Plan comme concept fonctionnel. Command reste source de runtime Detection, Signal, Alert et Incident. Settings reste source des environments, targets, runtimes configurés et health administratif. Govern reste source des Action Requests, Decisions, Approvals, Response Runs et Results. Studio reste source des Tools et Automation Runs. Une projection locale ne remplace jamais son objet canonique propriétaire.

## 20. Provenance et audit
Enregistrer le besoin initial, Project, Hypothesis, Drafts, versions, Review Packages, Release Candidates, reviewers, readiness, targets, plans, Action Requests, Decisions, Approvals, Response Runs, Results, runtime observations, health, Signals/Alerts/Incidents liés, assessments, propositions, erreurs, interruptions, auteurs, timestamps, Tool Calls, Automation Runs, paramètres, dispositions humaines, exports et correlation IDs applicables à Promotion Plan. Aucune trace n’est supprimée ni réécrite silencieusement.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Target environment read/select | infrastructure and tenant scope | 0/2 | secrets/hidden tenants masked | step-up possible | planner/Settings admin | Settings/Investigate | Permissions |
| Promotion Plan create/update | future production sequence | 2 | target details scoped | OPEN-013 | author/reviewer | Investigate | Permissions |
| Cross-tenant promotion prepare | tenant isolation | 2 | explicit tenant boundaries | step-up required | independent approval future | Security/Govern | Permissions |

Les permissions atomiques, namespaces, RBAC/ABAC, step-up définitif et séparation des tâches finale restent reportés. Toute exécution réelle de classe 3 ou 4 reste chez Govern et l’owner runtime.

## 22. Limites et erreurs
- Investigate does not create or administer targets.
- Target selected ≠ target compatible.
- Promotion Plan ≠ promotion execution.
- Runtime/language choice remains OPEN-017.
- Les données peuvent être stale, partielles, restreintes, incohérentes entre tenants ou indisponibles.
- Un Result, score, match, non-match, health state, disposition Command ou sortie IA ne vaut pas conclusion universelle.
- Timeout, cancellation, revocation, target offline, version superseded et partial result restent attribués et visibles.

## 23. Métriques conceptuelles
- plans by state and environment
- target incompatibility/staleness
- environment differences documented
- silent target selection — target zero
- complétude de provenance, versions, targets, autorité et dispositions humaines
- nombre d’exécutions silencieuses, d’auto-approbations ou de suppressions de traces — cible conceptuelle zéro

Aucun seuil universel de latency, coût, précision, recall, health ou business value n’est imposé.

## 24. Classification de livraison
`defined` / `planned`. Preuve documentaire uniquement. Aucun document n’est `validated`, `implemented`, `deployed`, `active`, `native` ou `integrated`. Aucun moteur, langage, syntaxe, produit tiers, modèle ML, API, protocole, commande ou code n’est choisi.

## 25. Critères d’acceptation
### 1. Mixed targets
**Given** un candidate avec targets compatible et incompatible  
**When** le plan est préparé  
**Then** les targets restent distincts et aucun succès global n’est déclaré

### 2. Stale health
**Given** un target avec health stale  
**When** il est sélectionné  
**Then** le stale state reste visible et readiness doit être réévaluée

### 3. Sans IA
**Given** aucun modèle  
**When** les targets sont comparés  
**Then** tables, matrices, diffs and human selection suffice

## 26. Questions ouvertes
- OPEN-008 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- OPEN-017 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- Les schémas, cardinalités, machines d’état, formats de déploiement, permissions atomiques, contrats techniques et composition détaillée des écrans restent futurs.
- OPEN-005 demeure forensic-only et n’est pas utilisée pour Detection.

## 27. Consommateurs documentaires
Detection Engineering lifecycle and Capability Maps ; Command runtime Detection/Signal/Alert/Incident ; Platform Settings environments/targets/health ; Endpoint Agent capabilities/version/health ; Govern Action Request/Decision/Approval/Response Run/Result ; Studio Tools/Runs ; Shared mechanisms ; phases Objects, Permissions, Screens, Journeys, Technique, validation et future 4B.3B uniquement comme handoff non canonique.
