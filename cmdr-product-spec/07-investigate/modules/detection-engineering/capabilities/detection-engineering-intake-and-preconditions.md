---
id: CAP-INV-401
title: Detection Engineering Intake and Preconditions
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
open_decisions:
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-401 — Detection Engineering Intake and Preconditions

## 1. Définition
Ouvrir, qualifier et borner un besoin de détection provenant d’un Finding, d’une Hypothesis, d’un Case, d’un Incident, d’une Hunt ou d’une analyse technique, sans créer silencieusement de Detection Content ni de runtime Detection.

## 2. Problème utilisateur
Un besoin transmis sans comportement, sources, limites ou consommateurs explicites peut devenir un draft ambigu et donner l’impression qu’une règle existe déjà.

## 3. Objectifs
- ouvrir le besoin depuis Finding, Hypothesis, Case, Incident, Hunt ou handoff technique
- afficher observations, Evidence, contradictions, limitations et provenance
- identifier données et sources de télémétrie candidates
- définir objectif, scope, consommateurs et owner initial
- créer ou reprendre explicitement un Detection Engineering Project en préservant le return origin

## 4. Non-objectifs
- ne pas choisir de moteur, langage de règle, syntaxe SIEM propriétaire, modèle ML ou produit tiers
- ne pas définir API, protocole, parser, compilateur, AST, format de stockage, pipeline ou architecture de streaming
- ne pas promouvoir, déployer, activer, désactiver, rollback, créer une exception active ou supprimer Signal/Alert
- ne pas créer de capability Threat Intelligence, CAP-INV-5xx, objet Intelligence canonique ou écran détaillé
- ne pas présenter un draft, test, replay, score, mapping ou package comme état runtime effectif

## 5. Propriétaire
Investigate / Detection Engineering / Investigate Product Lead possède Detection Engineering Intake, son contexte d’authoring et les dispositions humaines. Command conserve runtime Detection, Signal, Alert et Incident. Platform Settings conserve l’administration des sources, parsers, schemas, health, retention et environments. Endpoint Agent conserve ses capacités et résultats locaux. Studio conserve Tool, Tool Call, Workflow, Automation Agent, Human Gate et Automation Run. Govern conserve Decision, Approval et l’autorité future de production. Shared conserve Jobs, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Collaboration, Inspector et Recovery.

## 6. Utilisateurs
Principal : **Detection Engineer**. Secondaires : Investigation Lead; Case Analyst; Threat Hunter.

## 7. Conditions d’entrée
Tenant, environnement, objectif, scope, sources et permissions sont explicites. Les versions, restrictions, limites et return origin sont visibles. Toute dépendance absente conduit à un état `incomplete`, `partial`, `blocked` ou équivalent propre à la capability ; aucune donnée, ground truth ou autorité n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
| --- | --- | --- | ---: | --- | --- |
| Need source package | Finding / Hypothesis / Case / Incident / Hunt / CAP-INV-313,328,346,362,379,397 | besoin, observations, conditions, limites et provenance | oui | version source liée | intake `incomplete` |
| Existing Evidence and contradictions | CAP-INV-107..109 / analyses | contexte qualifié et éléments contraires | non | état courant | incertitude visible |
| Telemetry candidate context | Event Search / Settings projections | sources, environnements et disponibilité candidate | non | fraîcheur affichée | `data-source-required` |
| Permission and return origin | Security / Experience Architecture | accès, tenant, environnement et navigation | oui | réévalués à l’ouverture | `permission-blocked` |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
| --- | --- | --- | --- |
| Finding / Evidence / Hypothesis / Case | Investigate | besoin, raisonnement, sources, restrictions et provenance | lire et lier |
| Incident | Command | contexte opérationnel et consommateurs | lecture uniquement |
| Hunt / Query / Search Job | Investigate / Shared | question, recherches, résultats et fenêtres | lecture et référence |
| Technical analysis handoff packages | Investigate | comportements, conditions, Artifacts, limites | lecture et sélection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
| --- | --- | --- | --- |
| Detection Engineering Intake | créer, compléter, bloquer, retirer ou superseder | Investigate concept | aucun Detection Content silencieux |
| Detection Engineering Project relation | préparer ou lier explicitement | Investigate concept | création confirmée par l’utilisateur |
| Trace / Activity event | émettre | Shared | source, auteur et return origin conservés |

## 11. Fonctionnalités
- ouvrir depuis chaque source autorisée
- inspecter sources, observations, Evidence et contradictions
- qualifier les informations manquantes
- définir objectif, scope et consommateurs
- créer ou reprendre explicitement le projet
- conserver erreurs, résultats partiels, restrictions, attribution et return origin
- fonctionner intégralement sans fournisseur de modèle

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
| --- | --- | --- | ---: | --- | --- | --- |
| Consulter Detection Engineering Intake | Detection Engineer | Detection Engineering Intake | 0 | lecture et scope autorisés | projection sourcée | non |
| Exécuter un traitement borné pour Detection Engineering Intake | Detection Engineer | Tool Call / Analysis Result | 1 | déclenchement explicite, Tool/version et permission | résultat attribué, partiel ou complet | non; policy applicable |
| Créer ou modifier Detection Engineering Intake | Detection Engineer | Detection Engineering Intake | 2 | permission de mutation réversible | nouvelle version et disposition | OPEN-013 |
| Préparer un handoff | Detection Engineer | Candidate package / relation | 2 | sources, limites et destination visibles | package non effectif | destination owner |

Les classes 3 et 4 sont exclues. Toute promotion, deployment, activation, deactivation, rollback, production exception ou modification runtime appartient à la future Phase 4B.3A.2 et aux owners applicables.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
| --- | ---: | ---: | ---: | ---: | --- |
| Construire ou examiner Detection Engineering Intake | oui | règles, catalogues et contrôles explicables | oui | suggestion attribuée | éditeur structuré, tables et revue manuelle |
| Comparer ou valider Detection Engineering Intake | oui | oui lorsque les préconditions existent | oui | explication facultative | comparateur, validateur et diagnostics déterministes |
| Expliquer erreurs ou contradictions de Detection Engineering Intake | oui | catalogue et contrôles | oui | résumé sourcé | erreurs brutes, diff et checklist |
| Promouvoir, déployer, activer, désactiver ou qualifier runtime | non dans cette phase | non | non | interdit | future 4B.3A.2 et owners applicables |

Toute sortie automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres fonctionnels, timestamp, statut, erreurs, incertitude, owner humain et disposition acceptée, modifiée ou rejetée. Aucun Tool, field, mapping, oracle, label, permission ou état effectif n’est choisi silencieusement.

## 14. États fonctionnels
`draft`, `incomplete`, `ready`, `blocked`, `insufficient-evidence`, `data-source-required`, `permission-blocked`, `out-of-scope`, `superseded`. Ces états sont fonctionnels et ne constituent pas une machine d’état objet définitive.

## 15. États d’interface
Loading conserve Project, version, sélection et return origin ; Empty distingue absence de contenu et absence d’accès ; Partial expose gaps et résultats utilisables ; Error conserve les résultats valides ; Offline interdit les nouveaux traitements non garantis ; Permission denied masque les données protégées ; Stale distingue la dernière version connue de l’état courant ; conflit de version fournit un diff et une reprise sûre. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
| --- | --- | --- | --- |
| Qualified intake | Detection Engineering Intake | CAP-INV-402 | source, lacunes, objectif et auteur visibles |
| Data-readiness question | source requirement context | CAP-INV-404 | aucune disponibilité supposée |
| Return context | navigation context | source origin | tenant, sélection et position conservés |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
| --- | --- | --- | --- | --- |
| Finding / Hypothesis / Case / Incident / Hunt / technical analysis | ouvrir un besoin | CAP-INV-401 | source, observations, Evidence, contradictions, limites et provenance | source origin |
| CAP-INV-401 | intake `ready` | CAP-INV-402 | objectif, scope, owner, sources candidates, consommateurs | intake |
| CAP-INV-401 | source inconnue ou insuffisante | CAP-INV-404 | besoins de données, plateformes, environnements et période | intake |

Chaque transition conserve tenant, environnement, ownership, restrictions, permissions, erreurs, versions, provenance et return origin. Une transition ne crée jamais d’autorité ou d’accès supplémentaire.

## 18. Dépendances
CAP-INV-005,103,107..109,313,328,346,362,379,397; Command Incident; Shared Linking/Trace; OPEN-013/015. Shared Background Jobs, Notifications, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Collaboration, Comparison, Inspector, Context Bar, Audit Hooks et Recovery sont consommés sans redéfinition. Aucun moteur ou contrat technique bas niveau n’est sélectionné.

## 19. Source de vérité
Investigate est source de Detection Engineering Intake et de sa disposition humaine. Les objets canoniques ou projections consommés restent chez leurs owners. Un Detection Content Draft n’est jamais la source de vérité d’un runtime Detection Command. Tool, Tool Call, Automation Run, Query, Telemetry Event, Data Source, Parser, Decision et Approval conservent leurs sources canoniques.

## 20. Provenance et audit
Enregistrer Project, source du besoin, Case/Incident/Finding/Hypothesis/Hunt ou analyse technique, versions, Data Sources, schemas, fields, mappings, logic, enrichments, Tool/version, Tool Calls, Automation Runs, paramètres, datasets, restrictions, résultats, erreurs, interruptions, auteurs, reviewers, timestamps, dispositions, exports et correlation ID applicables à Detection Engineering Intake. Aucune trace n’est supprimée ou remplacée silencieusement.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
| --- | --- | ---: | --- | --- | --- | --- | --- |
| Detection Engineering Intake read | contexte sensible | 0 | relations selon droits | possible | viewer/reviewer | Investigate | Permissions |
| Detection Engineering Intake create/update | mutation analytique | 2 | sources restreintes masquées | OPEN-013 | auteur/reviewer | Investigate | Permissions |
| Cross-product source link | élargissement d’accès | 2 | aucune projection non autorisée | possible | source owner distinct | Investigate/owners | Permissions |

La matrice atomique, les namespaces, le RBAC/ABAC, le step-up définitif et la séparation des tâches finale restent reportés à la phase Permissions. Les droits de deployment, activation, deactivation, rollback et production exception sont hors périmètre.

## 22. Limites et erreurs
- Un besoin incomplet reste `incomplete`.
- Finding, Evidence, Incident ou résultat technique ne devient pas Detection Content.
- Aucun runtime Detection, Signal, Alert, déploiement ou activation n’est créé.
- source inaccessible, stale, superseded, partielle ou policy-blocked ; tenant/environnement incohérent
- Tool/version indisponible, timeout, cancellation, permission révoquée ou résultat incomplet
- aucun résultat Tool, score, match, non-match, sortie IA ou mapping ne vaut conclusion ou autorité à lui seul

## 23. Métriques conceptuelles
- nombre de Detection Engineering Intake par état et version
- proportion de résultats partial, blocked, disputed ou failed
- complétude de provenance, sources, paramètres et dispositions humaines
- temps conceptuel entre entrée, revue et handoff
- nombre d’actions silencieuses, promotions ou déploiements — cible conceptuelle zéro

Aucune cible de performance de production, précision, recall, coût runtime, drift ou santé de règle déployée n’est fixée dans 4B.3A.1.

## 24. Classification de livraison
`defined` / `planned`. Preuve documentaire uniquement. Aucun document n’est `validated`, `implemented`, `deployed`, `active`, `native` ou `integrated`. Aucun moteur, langage, syntaxe, produit tiers, modèle ML, API, protocole, commande ou code n’est choisi.

## 25. Critères d’acceptation
### 1. Besoin incomplet
**Given** un Finding sans source de données et un comportement insuffisamment défini  
**When** le Detection Engineer ouvre l’intake  
**Then** l’état reste `incomplete`, les informations manquantes sont visibles et aucun Detection Content n’est créé

### 2. Permission refusée
**Given** une source liée mais inaccessible  
**When** l’utilisateur ouvre l’intake  
**Then** la référence peut rester visible selon policy, le contenu reste masqué et aucune permission n’est élargie

### 3. Sans IA
**Given** aucun fournisseur de modèle  
**When** le besoin est qualifié  
**Then** formulaires, checklists, liens et revue humaine couvrent le workflow

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune décision n’est fermée par cette capability.
- OPEN-015 reste ouverte ; aucune décision n’est fermée par cette capability.
- Le choix futur d’un moteur ou langage de détection n’est couvert par aucune décision ouverte existante ; cette lacune est enregistrée sans bloquer la définition fonctionnelle et sans détourner OPEN-005.
- Les schémas, cardinalités, machines d’état, formats de règle, permissions atomiques, contrats techniques et composition détaillée des écrans restent futurs.

## 27. Consommateurs documentaires
Detection Engineering module and Capability Map ; Event Search, Hunt, Case, Evidence, Finding and technical Analysis handoffs ; Command runtime boundaries ; Platform Settings and Endpoint projections ; Studio Tools/Evaluations ; Govern future review ; Shared mechanisms ; phases Objects, Permissions, Screens, Journeys, Technique, 4B.3A.2 and validation.
