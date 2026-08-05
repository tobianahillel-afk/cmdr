---
id: CAP-INV-402
title: Detection Engineering Project and Workspace Management
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
  - REQ-AI-002
  - REQ-SEC-001
  - REQ-UX-006
open_decisions:
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-402 — Detection Engineering Project and Workspace Management

## 1. Définition
Créer et gérer un espace de travail durable regroupant objectifs, sources, drafts, tests, validations, replays, coverage assessments, gaps, versions, erreurs et handoffs d’authoring.

## 2. Problème utilisateur
Sans projet propriétaire et versionné, les drafts, datasets et décisions sont dispersés entre Cases, Hunts et Automation Runs, rendant l’authoring non reproductible.

## 3. Objectifs
- créer, reprendre, suspendre, clôturer, rouvrir, archiver ou superseder un projet
- lier Cases, Findings, Hypotheses, Hunts et Incidents sans les absorber
- gérer owner, contributeurs, scope et objectifs
- conserver drafts, tests, validations, replays, couverture, gaps, versions et erreurs
- préparer la transmission vers la future revue sans déploiement

## 4. Non-objectifs
- ne pas choisir de moteur, langage de règle, syntaxe SIEM propriétaire, modèle ML ou produit tiers
- ne pas définir API, protocole, parser, compilateur, AST, format de stockage, pipeline ou architecture de streaming
- ne pas promouvoir, déployer, activer, désactiver, rollback, créer une exception active ou supprimer Signal/Alert
- ne pas créer de capability Threat Intelligence, CAP-INV-5xx, objet Intelligence canonique ou écran détaillé
- ne pas présenter un draft, test, replay, score, mapping ou package comme état runtime effectif

## 5. Propriétaire
Investigate / Detection Engineering / Investigate Product Lead possède Detection Engineering Project, son contexte d’authoring et les dispositions humaines. Command conserve runtime Detection, Signal, Alert et Incident. Platform Settings conserve l’administration des sources, parsers, schemas, health, retention et environments. Endpoint Agent conserve ses capacités et résultats locaux. Studio conserve Tool, Tool Call, Workflow, Automation Agent, Human Gate et Automation Run. Govern conserve Decision, Approval et l’autorité future de production. Shared conserve Jobs, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Collaboration, Inspector et Recovery.

## 6. Utilisateurs
Principal : **Detection Engineering Lead**. Secondaires : Detection Engineer; Reviewer; Contributor.

## 7. Conditions d’entrée
Tenant, environnement, objectif, scope, sources et permissions sont explicites. Les versions, restrictions, limites et return origin sont visibles. Toute dépendance absente conduit à un état `incomplete`, `partial`, `blocked` ou équivalent propre à la capability ; aucune donnée, ground truth ou autorité n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | ---: | --- | --- |
| Qualified intake | CAP-INV-401 | objectif, scope, sources et consommateurs | oui | version courante | projet `draft` |
| Linked investigation context | Cases / Findings / Hypotheses / Hunts / Incidents | sources et return origins | non | références résolubles | projet sans lien |
| Contributor and permission context | Settings / Security | owner, contributeurs et restrictions | oui | courant | `blocked` |
| Existing project snapshot | Investigate | versions, drafts et résultats | non | dernière version autorisée | nouveau projet |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Case / Finding / Hypothesis / Hunt | Investigate | sources et reasoning context | lecture/lien |
| Incident | Command | contexte opérationnel | lecture/lien |
| Tool / Tool Call / Automation Run | Studio | exécutions et attribution | lecture/lien |
| Detection Content Drafts and test assets | Investigate concepts | contenu et versions du projet | lecture/organisation |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Detection Engineering Project | créer, modifier, suspendre, clôturer, rouvrir, archiver, superseder | Investigate concept | distinct du Case et de l’Automation Run |
| Project membership and links | ajouter ou retirer réversiblement | Investigate | aucun transfert d’ownership |
| Trace / Activity event | émettre | Shared | chaque changement attribué |

## 11. Fonctionnalités
- gérer le cycle fonctionnel du projet
- maintenir owner et contributeurs
- regrouper drafts, scenarios, datasets, validations, replays, coverage et gaps
- conserver erreurs, résultats partiels et versions
- comparer ou superseder plusieurs projets
- conserver erreurs, résultats partiels, restrictions, attribution et return origin
- fonctionner intégralement sans fournisseur de modèle

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | ---: | --- | --- | --- |
| Consulter Detection Engineering Project | Detection Engineering Lead | Detection Engineering Project | 0 | lecture et scope autorisés | projection sourcée | non |
| Exécuter un traitement borné pour Detection Engineering Project | Detection Engineering Lead | Tool Call / Analysis Result | 1 | déclenchement explicite, Tool/version et permission | résultat attribué, partiel ou complet | non; policy applicable |
| Créer ou modifier Detection Engineering Project | Detection Engineering Lead | Detection Engineering Project | 2 | permission de mutation réversible | nouvelle version et disposition | OPEN-013 |
| Préparer un handoff | Detection Engineering Lead | Candidate package / relation | 2 | sources, limites et destination visibles | package non effectif | destination owner |

Les classes 3 et 4 sont exclues. Toute promotion, deployment, activation, deactivation, rollback, production exception ou modification runtime appartient à la future Phase 4B.3A.2 et aux owners applicables.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | ---: | ---: | ---: | ---: | --- |
| Construire ou examiner Detection Engineering Project | oui | règles, catalogues et contrôles explicables | oui | suggestion attribuée | éditeur structuré, tables et revue manuelle |
| Comparer ou valider Detection Engineering Project | oui | oui lorsque les préconditions existent | oui | explication facultative | comparateur, validateur et diagnostics déterministes |
| Expliquer erreurs ou contradictions de Detection Engineering Project | oui | catalogue et contrôles | oui | résumé sourcé | erreurs brutes, diff et checklist |
| Promouvoir, déployer, activer, désactiver ou qualifier runtime | non dans cette phase | non | non | interdit | future 4B.3A.2 et owners applicables |

Toute sortie automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres fonctionnels, timestamp, statut, erreurs, incertitude, owner humain et disposition acceptée, modifiée ou rejetée. Aucun Tool, field, mapping, oracle, label, permission ou état effectif n’est choisi silencieusement.

## 14. États fonctionnels
`draft`, `ready`, `active`, `paused`, `blocked`, `partial`, `completed-authoring`, `failed`, `archived`, `superseded`. Ces états sont fonctionnels et ne constituent pas une machine d’état objet définitive.

## 15. États d’interface
Loading conserve Project, version, sélection et return origin ; Empty distingue absence de contenu et absence d’accès ; Partial expose gaps et résultats utilisables ; Error conserve les résultats valides ; Offline interdit les nouveaux traitements non garantis ; Permission denied masque les données protégées ; Stale distingue la dernière version connue de l’état courant ; conflit de version fournit un diff et une reprise sûre. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Detection Engineering Project | project context | CAP-INV-403..417 | scope, owner, versions et liens visibles |
| Project snapshot | versioned projection | Comparison / Recovery | contenu et erreurs conservés |
| Review handoff context | candidate package context | CAP-INV-417 | aucune Approval ou promotion |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| CAP-INV-401 | création explicite | CAP-INV-402 | objectif, scope, owner, sources et return origin | intake |
| CAP-INV-402 | définir le besoin | CAP-INV-403 | project, sources, observations et limitations | project |
| CAP-INV-402 | évaluer les données | CAP-INV-404 | project, sources requises, environnements et période | project |
| CAP-INV-402 | authoring terminé | CAP-INV-417 | versions candidates, résultats, gaps et provenance | project |

Chaque transition conserve tenant, environnement, ownership, restrictions, permissions, erreurs, versions, provenance et return origin. Une transition ne crée jamais d’autorité ou d’accès supplémentaire.

## 18. Dépendances
CAP-INV-401; Case/Hunt/Incident projections; Shared Collaboration/Versioning/Recovery; Studio Runs; OPEN-013/015. Shared Background Jobs, Notifications, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Collaboration, Comparison, Inspector, Context Bar, Audit Hooks et Recovery sont consommés sans redéfinition. Aucun moteur ou contrat technique bas niveau n’est sélectionné.

## 19. Source de vérité
Investigate est source de Detection Engineering Project et de sa disposition humaine. Les objets canoniques ou projections consommés restent chez leurs owners. Un Detection Content Draft n’est jamais la source de vérité d’un runtime Detection Command. Tool, Tool Call, Automation Run, Query, Telemetry Event, Data Source, Parser, Decision et Approval conservent leurs sources canoniques.

## 20. Provenance et audit
Enregistrer Project, source du besoin, Case/Incident/Finding/Hypothesis/Hunt ou analyse technique, versions, Data Sources, schemas, fields, mappings, logic, enrichments, Tool/version, Tool Calls, Automation Runs, paramètres, datasets, restrictions, résultats, erreurs, interruptions, auteurs, reviewers, timestamps, dispositions, exports et correlation ID applicables à Detection Engineering Project. Aucune trace n’est supprimée ou remplacée silencieusement.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | ---: | --- | --- | --- | --- | --- |
| Detection Engineering Project read | contexte et datasets liés | 0 | selon classification | possible | viewer/contributor | Investigate | Permissions |
| Project create/update/close/reopen | cycle analytique | 2 | liens restreints masqués | OPEN-013 | owner/reviewer | Investigate | Permissions |
| Project contributor management | accès et responsabilité | 2 | aucune auto-attribution | possible | owner/admin distinct | Investigate/Settings | Permissions |

La matrice atomique, les namespaces, le RBAC/ABAC, le step-up définitif et la séparation des tâches finale restent reportés à la phase Permissions. Les droits de deployment, activation, deactivation, rollback et production exception sont hors périmètre.

## 22. Limites et erreurs
- Le projet n’est ni un Case, ni un Workflow, ni un Automation Run.
- `completed-authoring` ne signifie ni approved, ni promoted, ni deployed.
- Aucun état `production` ou `active-rule` n’est créé.
- source inaccessible, stale, superseded, partielle ou policy-blocked ; tenant/environnement incohérent
- Tool/version indisponible, timeout, cancellation, permission révoquée ou résultat incomplet
- aucun résultat Tool, score, match, non-match, sortie IA ou mapping ne vaut conclusion ou autorité à lui seul

## 23. Métriques conceptuelles
- nombre de Detection Engineering Project par état et version
- proportion de résultats partial, blocked, disputed ou failed
- complétude de provenance, sources, paramètres et dispositions humaines
- temps conceptuel entre entrée, revue et handoff
- nombre d’actions silencieuses, promotions ou déploiements — cible conceptuelle zéro

Aucune cible de performance de production, précision, recall, coût runtime, drift ou santé de règle déployée n’est fixée dans 4B.3A.1.

## 24. Classification de livraison
`defined` / `planned`. Preuve documentaire uniquement. Aucun document n’est `validated`, `implemented`, `deployed`, `active`, `native` ou `integrated`. Aucun moteur, langage, syntaxe, produit tiers, modèle ML, API, protocole, commande ou code n’est choisi.

## 25. Critères d’acceptation
### 1. Projet lié
**Given** un intake `ready` avec plusieurs sources d’investigation  
**When** le projet est créé  
**Then** les sources restent chez leurs owners et le return origin est conservé

### 2. Projet partiel
**Given** un replay ou dataset devient inaccessible  
**When** le projet est repris  
**Then** les résultats valides restent visibles et l’état peut devenir `partial`

### 3. Sans IA
**Given** aucun modèle disponible  
**When** le projet est géré  
**Then** navigation, versioning, collaboration, diff et recovery restent disponibles

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune décision n’est fermée par cette capability.
- OPEN-015 reste ouverte ; aucune décision n’est fermée par cette capability.
- Le choix futur d’un moteur ou langage de détection n’est couvert par aucune décision ouverte existante ; cette lacune est enregistrée sans bloquer la définition fonctionnelle et sans détourner OPEN-005.
- Les schémas, cardinalités, machines d’état, formats de règle, permissions atomiques, contrats techniques et composition détaillée des écrans restent futurs.

## 27. Consommateurs documentaires
Detection Engineering module and Capability Map ; Event Search, Hunt, Case, Evidence, Finding and technical Analysis handoffs ; Command runtime boundaries ; Platform Settings and Endpoint projections ; Studio Tools/Evaluations ; Govern future review ; Shared mechanisms ; phases Objects, Permissions, Screens, Journeys, Technique, 4B.3A.2 and validation.
