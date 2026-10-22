---
id: CAP-INV-420
title: Detection Change Request and Approval Handoff
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
# CAP-INV-420 — Detection Change Request and Approval Handoff

## 1. Définition
Préparer un Detection Change Request Draft complet, le transmettre explicitement vers l’Action Request Govern et suivre les retours sans auto-approbation ni exécution locale.

## 2. Problème utilisateur
Un plan ou une recommandation peut être confondu avec une Action Request acceptée, une Decision ou une Approval.

## 3. Objectifs
- documenter candidate, targets, type, classe, risques, impacts et dépendances
- documenter rollout, rollback, critères d’arrêt, validations, reviewers et urgence
- couvrir promotion, shadow, canary, activation, désactivation, tuning, suppression, exception, rollback et retirement
- transmettre vers Action Request avec return origin
- répondre aux demandes d’information, annuler avant soumission ou superseder

## 4. Non-objectifs
- ne pas définir API, protocole, compilateur, parser, AST, format de package, stockage, pipeline, streaming, commande ou code produit
- ne sélectionner aucun moteur, langage, représentation cible, syntaxe vendor, produit tiers ou modèle ML
- ne réaliser aucun déploiement, activation, désactivation, rollback, suppression ou exception réelle pendant la phase documentaire
- ne modifier ni supprimer silencieusement runtime Detection, Signal, Alert ou Incident
- ne créer aucune capability CAP-INV-5xx, aucun objet Threat Intelligence, Cloud Analysis, Mobile Forensics ou écran détaillé

## 5. Propriétaire
Investigate / Detection Engineering possède **Detection Change Request Draft** et ses dispositions fonctionnelles. Command conserve runtime Detection, Signal, Alert, Incident et le feedback opérationnel. Platform Settings conserve runtimes configurés, targets, environments, tenants, sources, parsers, schemas, health, providers, secrets et canaux administratifs. Endpoint Agent conserve ses capacités, versions, health et exécutions locales autorisées. Govern conserve Action Request, Decision, Approval, Response Run, Result et toute autorité de classe 3 ou 4. Studio conserve Tool, Tool Call, Workflow, Automation Agent, Human Gate et Automation Run. Shared conserve Jobs, Notifications, Trace, Activity, Linking, Versioning, Comparison, Reporting, Export, Collaboration, Audit Hooks et Recovery.

## 6. Utilisateurs
Principal : **Detection Owner**. Secondaires : Detection Engineer; Detection Reviewer; Detection Owner; Incident Commander; Platform Operator; Approver selon le contexte.

## 7. Conditions d’entrée
Tenant, environnement, version, owner, cible, période, permission, restrictions et return origin sont explicites. Les preuves de CAP-INV-401..417 sont référencées sans duplication. Toute dépendance absente conduit à un état incomplet, partiel, bloqué ou inconnu ; aucune version, santé, efficacité, autorité ou ground truth n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | --- | --- | --- |
| Release Candidate and Readiness | CAP-INV-418/419 | candidate, targets, blockers and conditions | oui | current versions | draft incomplete |
| Change-specific plan | CAP-INV-421..423/429/430/433/434 | scope, steps, stop and recovery | oui selon type | version linked | submission blocked |
| Risk and authority context | Govern / Security | class, risk, approval and separation | oui | current policy | clarification required |
| Evidence and provenance | CAP-INV-417 / Shared Trace | tests, review, sources and lineage | oui | resolvable chain | not ready |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Release Candidate / Readiness / Plans | Investigate | change evidence and scope | lecture |
| Action Request / Decision / Approval | Govern | submission and authority projections | lecture/lien |
| Environment / Deployment Target | Settings | target identity and restrictions | lecture |
| Trace / Comments | Shared | return and information requests | consommation |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Detection Change Request Draft | créer, modifier, annuler, superseder | Investigate concept | Draft ≠ accepted Action Request |
| Action Request context | soumettre explicitement | Govern | producer cannot self-approve |
| Handoff disposition | enregistrer | Investigate / Shared | submitted, returned, rejected or cancelled |

## 11. Fonctionnalités
- documenter candidate, targets, type, classe, risques, impacts et dépendances
- documenter rollout, rollback, critères d’arrêt, validations, reviewers et urgence
- couvrir promotion, shadow, canary, activation, désactivation, tuning, suppression, exception, rollback et retirement
- transmettre vers Action Request avec return origin
- répondre aux demandes d’information, annuler avant soumission ou superseder
- conserver versions, sources, erreurs, résultats partiels, restrictions, attribution et return origin
- fonctionner entièrement sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | --- | --- | --- | --- |
| Consulter, filtrer ou comparer | Detection Owner | Detection Change Request Draft | 0 | lecture autorisée et scope explicite | projection sourcée | non |
| Exécuter une assessment ou comparaison bornée | Detection Owner | Tool Call / Assessment Result | 1 | déclenchement explicite, sources, paramètres et permission | résultat attribué, partial ou complete | policy applicable |
| Créer, modifier, contester ou retirer | Detection Owner | Detection Change Request Draft | 2 | mutation réversible et versionnée | nouvelle version et disposition humaine | OPEN-013 |
| Préparer une demande gouvernée | Detection Owner | Change/Action Request context | 2 | cible, effet, risque, rollback et autorité visibles | package non effectif | destination Govern |
| Exécuter le changement réel | Govern / Platform Operator | runtime target / Response Run | 3 | Decision/Approval et policy applicables | projection du Result seulement dans Investigate | Govern owner |
| Détruire historique ou provenance | aucun rôle local | historical content | 4 | interdit par défaut | action refusée | strict governance |

Les classes 3 et 4 ne sont jamais exécutées par Investigate. Une demande ou coordination locale de classe 2 ne vaut ni Decision, ni Approval, ni exécution.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | --- | --- | --- | --- | --- |
| Construire ou examiner Detection Change Request Draft | oui | checklists, catalogues et contrôles explicables | oui | proposition attribuée | formulaire, tables, matrices et revue humaine |
| Comparer ou valider Detection Change Request Draft | oui | comparateur et règles explicables | oui | explication facultative | diff, diagnostics et checklist |
| Résumer risques, erreurs ou contradictions | oui | catalogue et agrégations déterministes | oui | résumé sourcé | sources brutes, timeline et revue manuelle |
| Approuver, promouvoir, activer, désactiver ou rollback | non localement | non | non | interdit | Govern et owner runtime selon autorité |

Toute proposition automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet. Aucun target, état, approbation, suppression, exception ou rollback n’est sélectionné silencieusement.

## 14. États fonctionnels
`draft`, `incomplete`, `ready`, `submitted`, `awaiting-decision`, `approved`, `rejected`, `returned-for-information`, `cancelled`, `superseded`. Ces états sont des projections fonctionnelles ; ils ne redéfinissent ni les machines d’état Govern, Command, Settings ou runtime, ni une machine d’état objet définitive.

## 15. États d’interface
Loading conserve version, target, phase et return origin. Empty distingue absence de résultat et absence d’accès. Partial expose chaque target et résultat utilisable. Error conserve les données valides et le correlation ID. Offline interdit les mutations non garanties. Permission denied masque les données protégées. Stale sépare dernière observation connue et état actuel. Les conflits de version fournissent diff et reprise sûre. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Detection Change Request Draft | draft package | Govern Action Request | complete context, no authority |
| Action Request link | canonical Govern reference | CAP-INV-424/433/434 | ownership remains Govern |
| Return disposition | event | Project / source plan | conditions and return origin retained |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| Lifecycle plan | prepare governed change | CAP-INV-420 | candidate, targets, type, risk, rollout, rollback and evidence | source capability |
| CAP-INV-420 | explicit submit | Govern Action Request | requester, targets, action type, risks and rollback | Detection Project |
| Govern | Decision/info request/reject | CAP-INV-420 | canonical refs, conditions, rationale and requested info | Govern |

Chaque transition conserve owner, tenant, environnement, target, version, permissions, restrictions, erreurs, autorité, provenance et return origin. Elle n’accorde aucun accès ni pouvoir supplémentaire.

## 18. Dépendances
CAP-INV-417..419,421..423,429,430,433,434; Govern Action Request/Decision/Approval; Security SoD. OPEN-017 gouverne le futur choix runtime/langage/portabilité sans option sélectionnée. Shared Background Jobs, Notifications, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Collaboration, Comments, Assignments, Comparison, Inspector, Context Bar, Audit Hooks et Recovery sont consommés sans redéfinition.

## 19. Source de vérité
Investigate est source de vérité de Detection Change Request Draft comme concept fonctionnel. Command reste source de runtime Detection, Signal, Alert et Incident. Settings reste source des environments, targets, runtimes configurés et health administratif. Govern reste source des Action Requests, Decisions, Approvals, Response Runs et Results. Studio reste source des Tools et Automation Runs. Une projection locale ne remplace jamais son objet canonique propriétaire.

## 20. Provenance et audit
Enregistrer le besoin initial, Project, Hypothesis, Drafts, versions, Review Packages, Release Candidates, reviewers, readiness, targets, plans, Action Requests, Decisions, Approvals, Response Runs, Results, runtime observations, health, Signals/Alerts/Incidents liés, assessments, propositions, erreurs, interruptions, auteurs, timestamps, Tool Calls, Automation Runs, paramètres, dispositions humaines, exports et correlation IDs applicables à Detection Change Request Draft. Aucune trace n’est supprimée ni réécrite silencieusement.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Change Request prepare/update/cancel | future production request | 2 | target/risk scoped | OPEN-013 | requester/reviewer | Investigate | Permissions |
| Change Request submit | authority boundary | 2 | no hidden approval | step-up possible | requester cannot approve | Investigate/Govern | Permissions |
| Govern Decision/Approval read | authority data | 0 | policy-controlled | possible | read-only consumer | Govern | Permissions |

Les permissions atomiques, namespaces, RBAC/ABAC, step-up définitif et séparation des tâches finale restent reportés. Toute exécution réelle de classe 3 ou 4 reste chez Govern et l’owner runtime.

## 22. Limites et erreurs
- Change Request Draft ≠ accepted Action Request.
- Action Request ≠ Decision.
- Decision ≠ Approval ≠ Response Run.
- No self-approval or implicit emergency execution.
- Les données peuvent être stale, partielles, restreintes, incohérentes entre tenants ou indisponibles.
- Un Result, score, match, non-match, health state, disposition Command ou sortie IA ne vaut pas conclusion universelle.
- Timeout, cancellation, revocation, target offline, version superseded et partial result restent attribués et visibles.

## 23. Métriques conceptuelles
- drafts/submissions/returns by type
- missing rollback or stop criteria
- returned-for-information duration
- automatic approval — target zero
- complétude de provenance, versions, targets, autorité et dispositions humaines
- nombre d’exécutions silencieuses, d’auto-approbations ou de suppressions de traces — cible conceptuelle zéro

Aucun seuil universel de latency, coût, précision, recall, health ou business value n’est imposé.

## 24. Classification de livraison
`defined` / `planned`. Preuve documentaire uniquement. Aucun document n’est `validated`, `implemented`, `deployed`, `active`, `native` ou `integrated`. Aucun moteur, langage, syntaxe, produit tiers, modèle ML, API, protocole, commande ou code n’est choisi.

## 25. Critères d’acceptation
### 1. Review passed sans Govern
**Given** une review-passed et aucune Decision  
**When** le request est préparé  
**Then** seul un draft/Action Request existe; aucune promotion ou active version n’est inventée

### 2. Information requested
**Given** une Action Request returned-for-information  
**When** le owner répond  
**Then** la réponse et versions sont liées et l’autorité reste Govern

### 3. Sans IA
**Given** aucun modèle  
**When** un Change Request est préparé  
**Then** forms, checklists and risk matrices suffice

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- OPEN-015 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- OPEN-017 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- Les schémas, cardinalités, machines d’état, formats de déploiement, permissions atomiques, contrats techniques et composition détaillée des écrans restent futurs.
- OPEN-005 demeure forensic-only et n’est pas utilisée pour Detection.

## 27. Consommateurs documentaires
Detection Engineering lifecycle and Capability Maps ; Command runtime Detection/Signal/Alert/Incident ; Platform Settings environments/targets/health ; Endpoint Agent capabilities/version/health ; Govern Action Request/Decision/Approval/Response Run/Result ; Studio Tools/Runs ; Shared mechanisms ; phases Objects, Permissions, Screens, Journeys, Technique, validation et future 4B.3B uniquement comme handoff non canonique.
