---
id: CAP-INV-205
title: File and Directory Acquisition
product: investigate
module: collection-and-live-response
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-PROD-018
  - REQ-OBJ-004
open_decisions:
  - OPEN-008
  - OPEN-013
source-of-truth: canonical
---
# CAP-INV-205 — File and Directory Acquisition

## 1. Définition
Acquérir des fichiers ou répertoires ciblés avec bornes, permissions, métadonnées et résultats partiels.

## 2. Problème utilisateur
Sans File and Directory Acquisition, l’utilisateur perd le lien entre le Case, l’Endpoint, l’autorité applicable, l’exécution locale et les résultats. Les états partiels ou offline peuvent alors être pris pour un succès et les objets peuvent être confondus.

## 3. Objectifs
- fournir target paths/patterns, profondeur, taille, Case, policy;
- exposer cible, scope, fraîcheur, policy, permission et classe d’action;
- conserver erreurs, résultats partiels, provenance et retour au Case;
- produire Artifacts et métadonnées, erreurs missing/locked sans transférer l’ownership.

## 4. Non-objectifs
- ne pas administrer la Fleet ni les Endpoint Policies;
- ne pas définir protocole, API, commande, moteur, format, PKI, stockage ou plateforme supportée;
- ne pas créer automatiquement Evidence, Finding, Decision, Response Run ou Govern Result;
- ne pas commencer Analysis Workbench.

## 5. Propriétaire
Investigate / Collection and Live Response / Investigate Product Lead possède le contexte métier et les relations au Case. Platform Settings administre Fleet/Policies; Endpoint Agent exécute localement; Govern possède l’autorité risquée.

## 6. Utilisateurs
Principal : Case Analyst / DFIR Analyst. Secondaires : Investigation Lead, Evidence Reviewer, Incident Commander, approbateur Govern ou Platform Administrator en consultation selon la capability.

## 7. Conditions d’entrée
Tenant et environnement conservés, Case accessible, Endpoint résolu, fraîcheur et capacités visibles, policy projetée, permissions vérifiées et objectif explicite. Une dépendance absente produit un état partial/offline/unsupported, jamais un résultat inventé.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Case et objectif | Investigate | contexte métier | oui | version courante | rester draft ou refuser la mutation |
| Endpoint et état Agent | Platform Settings / Endpoint Agent | cible et disponibilité | oui | dernière communication visible | offline/unknown, aucune exécution présentée |
| Scope, limites et classe | utilisateur / policy | contrat d’action | oui | validés au déclenchement | incomplete ou policy-blocked |
| Autorité et permission | Security / Govern | droit et gate | selon classe | snapshot à l’action | denied ou awaiting-approval |
| Données spécifiques | target paths/patterns, profondeur, taille, Case, policy | données locales | selon opération | source/version visibles | résultat partiel explicite |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | identité, état, relations, fraîcheur ou autorité nécessaires | consulter et référencer; aucune administration implicite |
| Endpoint | owner canonique ou concept à formaliser | identité, état, relations, fraîcheur ou autorité nécessaires | consulter et référencer; aucune administration implicite |
| Endpoint Policy | Platform Settings | identité, état, relations, fraîcheur ou autorité nécessaires | consulter et référencer; aucune administration implicite |
| Collection Request | Investigate | identité, état, relations, fraîcheur ou autorité nécessaires | consulter et référencer; aucune administration implicite |
| Collection Job | owner canonique ou concept à formaliser | identité, état, relations, fraîcheur ou autorité nécessaires | consulter et référencer; aucune administration implicite |
| Artifact | Investigate | identité, état, relations, fraîcheur ou autorité nécessaires | consulter et référencer; aucune administration implicite |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Record métier File and Directory Acquisition | créer, mettre à jour ou supersede conceptuellement | Investigate, modèle final à Phase Objets | versionné, Case-scoped, sans machine finale |
| Artifact ou relation Artifact | créer/lier seulement lorsqu’un résultat matériel existe | Investigate | source et acquisition requises; Artifact ≠ Evidence |
| Événement métier | émettre vers Trace/Activity/Timeline/Audit Hooks | Shared mechanism, sémantique Investigate | acteur, cible, statut, erreur et correlation ID |
| Objet externe | aucune mutation administrative ou d’autorité | owner externe | projection uniquement, sauf commande locale autorisée par contrat |

## 11. Fonctionnalités
- prévisualiser; borner; demander; annuler;
- afficher capacités, limitations, policy, permission, classe, impact et autorité;
- gérer progression, partial, retry ciblé, cancel, timeout, déconnexion et reprise autorisée;
- lier les sorties au Case et aux Artifacts;
- préserver return origin et contexte.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter/inspecter | utilisateur autorisé | contexte et projections | 0 | read permission | vue sourcée et fraîcheur visible | non |
| Préparer ou lancer collecte bornée | analyste autorisé | request/job concept | 1 | scope, policy et capacité | demande ou exécution non destructive | selon impact |
| Modifier ou interrompre réversiblement | opérateur autorisé | session/opération/record | 2 | rollback/permission | transition auditée | OPEN-013 selon policy |
| Préparer containment | Investigation Lead | Action Request | 3 | Finding/Evidence/impact/rollback | demande vers CAP-INV-113/Govern | obligatoire |
| Préparer irréversible | Investigation Lead | Action Request | 4 | justification et alternatives | contexte seulement | obligatoire |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Construire scope/checklist | oui | profiles et règles | oui | suggestion modifiable | formulaire et profiles déterministes |
| Valider policy/capacité | oui | oui | oui | explication facultative | validateur et inventaire de capacités |
| Suivre progression/erreurs | oui | oui | oui | résumé | états et résultats bruts inspectables |
| Proposer prochaine action | oui | règles/workflow | oui | proposition attribuée | expertise humaine et procédures |
| Exécuter action sensible | humain explicite | contrat autorisé | workflow possible | jamais autonome | action humaine/Govern |

Toute sortie automatisée expose initiateur, moteur ou agent, version, Automation Run, Tool Calls, sources, paramètres fonctionnels, timestamp, statut, incertitude, owner humain, accept/modify/reject et trace.

## 14. États fonctionnels
`preparing`, `bounded`, `queued`, `acquiring`, `partial`, `completed`, `missing`, `locked`, `restricted`, `cancelled`. Ces dimensions sont Draft et ne finalisent aucune machine d’état objet.

## 15. États d’interface
Loading conserve Case et cible; Empty distingue absence de capacité et absence de résultat; Partial détaille les éléments réussis/échoués; Error préserve les données valides; Offline interdit toute présentation d’exécution démarrée; Permission denied ne révèle rien; Stale expose la dernière synchronisation.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Artifacts et métadonnées, erreurs missing/locked | record, relation ou événement métier | Case Workspace et capabilities dépendantes | cible, scope, acteur, statut, erreurs et version visibles |
| Artifact éventuel | Artifact Investigate | CAP-INV-105 puis CAP-INV-107 | source/acquisition conservées; aucune Evidence automatique |
| Progression et notification | Background Job/Notification projection | utilisateur et Case | succès partiels et échecs non masqués |
| Trace/provenance | événements métier | CAP-INV-110/112/214 et audit | correlation IDs, producteurs et corrections conservés |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Case Workspace | ouvrir activité endpoint | Endpoint Context / capability courante | tenant, environnement, Case, Incident, Endpoint, objectif, return origin | même Case et position |
| Endpoint Context | préparer/lancer | Collection Request, Job, Live Session ou opération | cible, capacités, policy, permission, classe, limites | Endpoint Context |
| Exécution locale | résultat/erreur | Operation Result / Artifact Management | opération, output, erreurs, fichiers, timestamps, provenance | Case ou session |
| Artifact | qualification humaine | CAP-INV-107 Evidence Creation | source, acquisition, Case, raison, transformations | Artifact |
| Finding/action risquée | préparer demande | CAP-INV-113 puis Govern | Finding, Evidence, Endpoint, impact, alternatives, rollback | Case avec projection Govern |

## 18. Dépendances
CAP-INV-102, 105, 107, 108, 110, 112, 113; Platform Settings Fleet/Policies/Health; Endpoint Agent capabilities; Shared Background Jobs, Notifications, Trace, Timeline, Inspector, Context Bar, Linking, Export et recovery; Govern; Studio optional; décisions ouvertes listées au front matter.

## 19. Source de vérité
Investigate est source du contexte métier et des relations Case. Platform Settings reste source de Fleet/Policy; Endpoint Agent de son état, commandes et résultats locaux; Govern de Decision/Response Run/Result; Studio d’Automation Run/Tool Calls; Shared des mécanismes génériques.

## 20. Provenance et audit
Case, Endpoint, Agent, policy/version, initiateur, permission, classe, scope, paramètres fonctionnels, autorité, timestamps, transitions, erreurs, résultats partiels, fichiers, Artifacts, Automation Run/Tool Calls, Action Request/Decision/Run/Result et disposition humaine.

## 21. Permissions fonctionnelles
Endpoint read, capability read, collection prepare/submit/cancel/retry, raw result read, Artifact receive/export, Live Session request/open/join/extend/close, operation execute/interrupt, file transfer, inspection, memory/network request, containment request, sensitive output, cross-tenant/environment, transcript read, result verify et custody review selon la capability. Step-up, séparation des tâches et matrice atomique sont reportés.

## 22. Limites et erreurs
Endpoint offline/stale/unsupported, Agent absent ou degraded, policy blocked, scope trop large, permission révoquée, timeout, déconnexion, conflit de session, résultat partiel, fichier manquant/verrouillé, cible changée, Govern indisponible ou tenant mismatch. Aucun retry ne duplique silencieusement l’effet.

## 23. Métriques
Temps de préparation et d’exécution, demandes bloquées par capacité/policy, résultats partiels, retries ciblés, annulations, déconnexions, Artifacts avec origine complète, opérations avec provenance complète, erreurs par catégorie et retours Case réussis.

## 24. Classification de livraison
`defined` / `planned`. Cible native via Endpoint Agent, mais aucune plateforme, moteur, protocole, commande, API ou release n’est prouvée. Promotion conditionnée par OPEN-008, objets, permissions, contrats d’autorité, preuve d’implémentation et validation.

## 25. Critères d’acceptation
**Given** un Case, un Endpoint disponible et un utilisateur autorisé **When** il utilise File and Directory Acquisition **Then** cible, scope, policy, classe, progression, résultat et retour au Case sont visibles sans transfert d’ownership.

**Given** un Endpoint offline, unsupported ou une permission refusée **When** l’action est demandée **Then** aucune exécution n’est présentée comme démarrée, l’état et les options sûres sont explicites et le Case reste accessible.

**Given** aucun fournisseur de modèle **When** le workflow est exécuté **Then** formulaires, profiles, règles, validateurs, Jobs, revue et actions humaines permettent le résultat essentiel.

## 26. Questions ouvertes
OPEN-008 conserve les plateformes; OPEN-013 la gouvernance classe 2; OPEN-007 Human Gate/Govern; OPEN-015 Automation Run/Response Run; OPEN-005 les moteurs forensics futurs lorsque référencé. Les objets et permissions détaillés restent à leurs phases.

## 27. Consommateurs documentaires
Module Collection and Live Response, Case Workspace, Evidence Board, Platform Settings Fleet/Policies/Health, Govern Action Center et Runs, parcours Endpoint investigation/containment/offline recovery, phases Objets/Permissions/Technique et future Phase 4B.2B uniquement comme handoff Artifact.