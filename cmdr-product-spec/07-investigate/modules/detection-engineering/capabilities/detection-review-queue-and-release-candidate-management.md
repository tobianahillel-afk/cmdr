---
id: CAP-INV-418
title: Detection Review Queue and Release Candidate Management
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
# CAP-INV-418 — Detection Review Queue and Release Candidate Management

## 1. Définition
Recevoir le Review Package de CAP-INV-417, créer explicitement un Release Candidate fondé sur une version immuable, organiser la revue humaine et enregistrer sa disposition sans créer d’Approval.

## 2. Problème utilisateur
Un Review Package peut être pris pour une version approuvée alors que ses preuves, reviewers, limites ou plans de retour sont incomplets.

## 3. Objectifs
- contrôler le Review Package sans dupliquer ses preuves
- créer un Release Candidate lié à une version immuable
- exposer validations, tests, replay, FP/FN candidates, couverture, gaps et limites
- assigner, commenter, comparer et demander des modifications
- accepter pour readiness, rejeter, retirer ou superseder

## 4. Non-objectifs
- ne pas définir API, protocole, compilateur, parser, AST, format de package, stockage, pipeline, streaming, commande ou code produit
- ne sélectionner aucun moteur, langage, représentation cible, syntaxe vendor, produit tiers ou modèle ML
- ne réaliser aucun déploiement, activation, désactivation, rollback, suppression ou exception réelle pendant la phase documentaire
- ne modifier ni supprimer silencieusement runtime Detection, Signal, Alert ou Incident
- ne créer aucune capability CAP-INV-5xx, aucun objet Threat Intelligence, Cloud Analysis, Mobile Forensics ou écran détaillé

## 5. Propriétaire
Investigate / Detection Engineering possède **Detection Release Candidate and Detection Review** et ses dispositions fonctionnelles. Command conserve runtime Detection, Signal, Alert, Incident et le feedback opérationnel. Platform Settings conserve runtimes configurés, targets, environments, tenants, sources, parsers, schemas, health, providers, secrets et canaux administratifs. Endpoint Agent conserve ses capacités, versions, health et exécutions locales autorisées. Govern conserve Action Request, Decision, Approval, Response Run, Result et toute autorité de classe 3 ou 4. Studio conserve Tool, Tool Call, Workflow, Automation Agent, Human Gate et Automation Run. Shared conserve Jobs, Notifications, Trace, Activity, Linking, Versioning, Comparison, Reporting, Export, Collaboration, Audit Hooks et Recovery.

## 6. Utilisateurs
Principal : **Detection Reviewer**. Secondaires : Detection Engineer; Detection Reviewer; Detection Owner; Incident Commander; Platform Operator; Approver selon le contexte.

## 7. Conditions d’entrée
Tenant, environnement, version, owner, cible, période, permission, restrictions et return origin sont explicites. Les preuves de CAP-INV-401..417 sont référencées sans duplication. Toute dépendance absente conduit à un état incomplet, partiel, bloqué ou inconnu ; aucune version, santé, efficacité, autorité ou ground truth n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | --- | --- | --- |
| Review Package | CAP-INV-417 | candidate version, preuves, limites, risques et unresolved items | oui | version soumise | candidate `incomplete` |
| Detection Content Version | CAP-INV-406/410 | snapshot fonctionnel immuable et metadata | oui | version résoluble | aucun Release Candidate |
| Validation/test/replay/coverage set | CAP-INV-411..416 | preuves et dispositions candidates | oui | versions liées | revue bloquée |
| Reviewer and permission context | Security / Shared Collaboration | reviewers, séparation et commentaires | oui | courant | assignment bloquée |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Review Package / Detection Content Version | Investigate | preuves, version, risques et lineage | lecture/comparaison |
| Validation/Test/Replay/Coverage concepts | Investigate | evidence set et limitations | lecture |
| Comments / Assignments / Versioning / Trace | Shared | collaboration et historique | consommation |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Detection Release Candidate | créer, retirer, superseder | Investigate concept | version source immuable; candidate ≠ active version |
| Detection Review | assigner, commenter, demander changements, disposer | Investigate concept | review-passed ≠ Approval |
| Review activity | émettre | Shared | auteur, version et disposition conservés |

## 11. Fonctionnalités
- contrôler le Review Package sans dupliquer ses preuves
- créer un Release Candidate lié à une version immuable
- exposer validations, tests, replay, FP/FN candidates, couverture, gaps et limites
- assigner, commenter, comparer et demander des modifications
- accepter pour readiness, rejeter, retirer ou superseder
- conserver versions, sources, erreurs, résultats partiels, restrictions, attribution et return origin
- fonctionner entièrement sans modèle IA

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | --- | --- | --- | --- |
| Consulter, filtrer ou comparer | Detection Reviewer | Detection Release Candidate and Detection Review | 0 | lecture autorisée et scope explicite | projection sourcée | non |
| Exécuter une assessment ou comparaison bornée | Detection Reviewer | Tool Call / Assessment Result | 1 | déclenchement explicite, sources, paramètres et permission | résultat attribué, partial ou complete | policy applicable |
| Créer, modifier, contester ou retirer | Detection Reviewer | Detection Release Candidate and Detection Review | 2 | mutation réversible et versionnée | nouvelle version et disposition humaine | OPEN-013 |
| Préparer une demande gouvernée | Detection Reviewer | Change/Action Request context | 2 | cible, effet, risque, rollback et autorité visibles | package non effectif | destination Govern |
| Exécuter le changement réel | Govern / Platform Operator | runtime target / Response Run | 3 | Decision/Approval et policy applicables | projection du Result seulement dans Investigate | Govern owner |
| Détruire historique ou provenance | aucun rôle local | historical content | 4 | interdit par défaut | action refusée | strict governance |

Les classes 3 et 4 ne sont jamais exécutées par Investigate. Une demande ou coordination locale de classe 2 ne vaut ni Decision, ni Approval, ni exécution.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | --- | --- | --- | --- | --- |
| Construire ou examiner Detection Release Candidate and Detection Review | oui | checklists, catalogues et contrôles explicables | oui | proposition attribuée | formulaire, tables, matrices et revue humaine |
| Comparer ou valider Detection Release Candidate and Detection Review | oui | comparateur et règles explicables | oui | explication facultative | diff, diagnostics et checklist |
| Résumer risques, erreurs ou contradictions | oui | catalogue et agrégations déterministes | oui | résumé sourcé | sources brutes, timeline et revue manuelle |
| Approuver, promouvoir, activer, désactiver ou rollback | non localement | non | non | interdit | Govern et owner runtime selon autorité |

Toute proposition automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet. Aucun target, état, approbation, suppression, exception ou rollback n’est sélectionné silencieusement.

## 14. États fonctionnels
`draft`, `incomplete`, `review-requested`, `under-review`, `changes-requested`, `review-passed`, `review-failed`, `readiness-pending`, `blocked`, `superseded`, `withdrawn`. Ces états sont des projections fonctionnelles ; ils ne redéfinissent ni les machines d’état Govern, Command, Settings ou runtime, ni une machine d’état objet définitive.

## 15. États d’interface
Loading conserve version, target, phase et return origin. Empty distingue absence de résultat et absence d’accès. Partial expose chaque target et résultat utilisable. Error conserve les données valides et le correlation ID. Offline interdit les mutations non garanties. Permission denied masque les données protégées. Stale sépare dernière observation connue et état actuel. Les conflits de version fournissent diff et reprise sûre. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Detection Release Candidate | functional candidate | CAP-INV-419 | version, preuves et unresolved items visibles |
| Detection Review disposition | review event | Project / CAP-INV-419 | recommendation, pas Approval |
| Returned-for-changes package | change context | CAP-INV-406/417 | Draft source et demandes conservés |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-417 | soumettre Review Package | CAP-INV-418 | version, validation, tests, replay, FP/FN, coverage, gaps, risques et provenance | CAP-INV-417/Project |
| CAP-INV-418 | review-passed | CAP-INV-419 | candidate, review disposition, conditions et unresolved items | Review Queue |
| CAP-INV-418 | changes-requested or rejected | CAP-INV-406/417 | source version, comments et required changes | Review Queue |

Chaque transition conserve owner, tenant, environnement, target, version, permissions, restrictions, erreurs, autorité, provenance et return origin. Elle n’accorde aucun accès ni pouvoir supplémentaire.

## 18. Dépendances
CAP-INV-406,410..417; Shared Collaboration/Comments/Assignments/Comparison/Versioning/Trace. OPEN-017 gouverne le futur choix runtime/langage/portabilité sans option sélectionnée. Shared Background Jobs, Notifications, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Collaboration, Comments, Assignments, Comparison, Inspector, Context Bar, Audit Hooks et Recovery sont consommés sans redéfinition.

## 19. Source de vérité
Investigate est source de vérité de Detection Release Candidate and Detection Review comme concept fonctionnel. Command reste source de runtime Detection, Signal, Alert et Incident. Settings reste source des environments, targets, runtimes configurés et health administratif. Govern reste source des Action Requests, Decisions, Approvals, Response Runs et Results. Studio reste source des Tools et Automation Runs. Une projection locale ne remplace jamais son objet canonique propriétaire.

## 20. Provenance et audit
Enregistrer le besoin initial, Project, Hypothesis, Drafts, versions, Review Packages, Release Candidates, reviewers, readiness, targets, plans, Action Requests, Decisions, Approvals, Response Runs, Results, runtime observations, health, Signals/Alerts/Incidents liés, assessments, propositions, erreurs, interruptions, auteurs, timestamps, Tool Calls, Automation Runs, paramètres, dispositions humaines, exports et correlation IDs applicables à Detection Release Candidate and Detection Review. Aucune trace n’est supprimée ni réécrite silencieusement.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Release Candidate read/create/update/withdraw | future production impact | 0/2 | restricted evidence masked | OPEN-013 | author/reviewer | Investigate | Permissions |
| Detection Review assign/comment/disposition | review authority | 2 | source permissions retained | step-up possible | reviewer distinct | Investigate | Permissions |
| Automated review proposal | biased recommendation | 1/2 | sources and uncertainty visible | possible | human reviewer mandatory | Studio/Investigate | Permissions |

Les permissions atomiques, namespaces, RBAC/ABAC, step-up définitif et séparation des tâches finale restent reportés. Toute exécution réelle de classe 3 ou 4 reste chez Govern et l’owner runtime.

## 22. Limites et erreurs
- Review Package ≠ Release Candidate.
- Release Candidate ≠ active version and ≠ runtime Detection.
- Detection Review ≠ Approval or Govern Decision.
- Blocking gaps prevent readiness.
- Les données peuvent être stale, partielles, restreintes, incohérentes entre tenants ou indisponibles.
- Un Result, score, match, non-match, health state, disposition Command ou sortie IA ne vaut pas conclusion universelle.
- Timeout, cancellation, revocation, target offline, version superseded et partial result restent attribués et visibles.

## 23. Métriques conceptuelles
- candidates by state/reviewer
- review duration and returned-for-changes rate
- missing-evidence categories
- automatic approvals — target zero
- complétude de provenance, versions, targets, autorité et dispositions humaines
- nombre d’exécutions silencieuses, d’auto-approbations ou de suppressions de traces — cible conceptuelle zéro

Aucun seuil universel de latency, coût, précision, recall, health ou business value n’est imposé.

## 24. Classification de livraison
`defined` / `planned`. Preuve documentaire uniquement. Aucun document n’est `validated`, `implemented`, `deployed`, `active`, `native` ou `integrated`. Aucun moteur, langage, syntaxe, produit tiers, modèle ML, API, protocole, commande ou code n’est choisi.

## 25. Critères d’acceptation
### 1. Review Package incomplet
**Given** un package sans rollback plan et avec gaps non résolus  
**When** le reviewer l’ouvre  
**Then** aucun candidate ready n’est déclaré, les manques sont visibles, changes can be requested et aucun déploiement n’est possible

### 2. Review passed sans Approval
**Given** une review-passed sans Govern Decision  
**When** le candidat est consulté  
**Then** il reste non approuvé et seule une Action Request peut être préparée

### 3. Sans IA
**Given** aucun modèle et un package complet  
**When** la revue est conduite  
**Then** checklists, matrices, diff, comments et décision humaine suffisent

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- OPEN-015 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- OPEN-017 reste ouverte ; aucune option n’est sélectionnée ou fermée par cette capability.
- Les schémas, cardinalités, machines d’état, formats de déploiement, permissions atomiques, contrats techniques et composition détaillée des écrans restent futurs.
- OPEN-005 demeure forensic-only et n’est pas utilisée pour Detection.

## 27. Consommateurs documentaires
Detection Engineering lifecycle and Capability Maps ; Command runtime Detection/Signal/Alert/Incident ; Platform Settings environments/targets/health ; Endpoint Agent capabilities/version/health ; Govern Action Request/Decision/Approval/Response Run/Result ; Studio Tools/Runs ; Shared mechanisms ; phases Objects, Permissions, Screens, Journeys, Technique, validation et future 4B.3B uniquement comme handoff non canonique.
