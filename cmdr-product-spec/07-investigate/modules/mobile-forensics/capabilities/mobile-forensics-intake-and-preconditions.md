---
id: CAP-INV-701
title: Mobile Forensics Intake and Preconditions
product: investigate
module: mobile-forensics
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-07
requirement_ids: [REQ-INV-001, REQ-PROD-014, REQ-PROD-020, REQ-PROD-055, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-005, OPEN-008, OPEN-011, OPEN-013, OPEN-014, OPEN-015]
source-of-truth: canonical
---
# CAP-INV-701 — Mobile Forensics Intake and Preconditions

## 1. Définition
Qualifier une demande d’analyse Mobile avant toute interprétation, en préservant origine, device/source, représentation disponible, période, autorité, restrictions, intégrité/completude connues, permissions, Tools compatibles déclarés et return origin.

## 2. Problème utilisateur
Un analyste peut sinon confondre un device avec une extraction, une extraction partielle avec un état complet, une autorisation de collecte avec une permission d’analyse, ou lancer une analyse sur une source verrouillée, chiffrée, restreinte ou hors scope.

## 3. Objectifs
- ouvrir depuis Case, Incident, Finding, Hunt, Signal, Collection Job, Artifact ou Mobile Evidence Package;
- afficher device source/candidats, plateforme déclarée/candidate, acquisition et représentation disponibles;
- rendre période, restrictions, permissions, custody, intégrité, complétude, lock/encryption/partial/unsupported et résultats précédents explicites;
- définir objectif/scope et préserver le contexte de retour.

## 4. Non-objectifs
Aucune analyse automatique à l’ouverture; aucune acquisition, plateforme supportée, outil imposé, unlock, bypass, rooting/jailbreak, commande, API, protocole, format propriétaire, exploitation, MDM ou action sur appareil réel.

## 5. Propriétaire
Investigate possède le Mobile Forensics Intake. Collection possède acquisition/request/job/result; Settings possède sources/Fleet/policies/secrets; Endpoint possède ses capacités déclarées; Govern possède l’autorité réelle; Studio possède Tools/Runs; Shared possède les mécanismes génériques.

## 6. Utilisateurs
Principal : Mobile Forensics Analyst. Secondaires : DFIR Analyst, Investigation Lead, Evidence Reviewer, SOC Analyst et Sensitive Data Reviewer.

## 7. Conditions d’entrée
Origine autorisée, tenant/environnement, objectif ou lacune, device/package/source candidate, période, restrictions, permission de lecture et return origin. Les éléments absents restent visibles comme préconditions manquantes.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Case/Incident/Finding/Hunt/Signal | owner source | objectif et contexte | non, au moins une origine | version courante | origine explicite requise |
| Collection Job / Artifact / Mobile Evidence Package | Collection / Investigate | source matérielle et limites | oui pour analyse de données | version disponible | `source-required` |
| Device et plateforme déclarée/candidate | source/Settings/analyste | contexte technique | oui comme candidate | source horodatée | `incomplete`, aucune supposition |
| Acquisition/représentation/custody | Collection/source | provenance et couverture | oui si connue | version d’acquisition | `incomplete` ou `disputed` |
| Scope, période, restrictions et permissions | analyste/Security | autorité analytique | oui | revalidé à l’ouverture | `permission-blocked`/`out-of-scope` |
| Résultats/Tools antérieurs | Studio/Investigate | contexte réutilisable | non | versionnés | aucun résultat inventé |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case / Finding / Artifact | Investigate | objectif, relations, return origin | lire/lier |
| Incident / Signal | Command | origine opérationnelle | lire/lier |
| Collection Request / Job | Collection/Investigate | acquisition, status, erreurs | lire uniquement |
| Mobile Evidence Package / backup/extraction concept | source/Investigate concept | représentation, scope, limites | lire selon permission |
| Device/Fleet/Policy/source projection | Settings / Endpoint | device context et support déclaré | lire uniquement |
| Tool / Tool Call / Automation Run | Studio | résultats et compatibilité déclarée | lire/lier |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Mobile Forensics Intake | créer, compléter, classer, supersede | Investigate | origine, objectif, source, scope et return origin obligatoires avant ready |
| Scope inclusion/exclusion | créer/modifier | Investigate | aucune permission source n’est créée |
| Collection gap proposal | préparer | Collection | aucune acquisition exécutée |
| Platform/Device candidate relation | proposer/disputer | Investigate | candidate sourcée; aucune attribution à une personne |

## 11. Fonctionnalités
Afficher origines, source package, device, plateforme/version candidates, représentation, lock/encryption, intégrité/completude connues, permissions/restrictions, période, résultats antérieurs, Tools compatibles déclarés, objectif, included/excluded scope, readiness et deep link retour.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter contexte | analyste | projections | 0 | read | readiness visible | non |
| Définir objectif/scope | analyste | Intake | 2 | permission manage | draft versionné | OPEN-013 |
| Sélectionner candidate plateforme/device | analyste | relation candidate | 2 | source visible | candidate attribuée à l’analyste | non |
| Demander source complémentaire | Investigation Lead | gap package | 2 | lacune justifiée | handoff Collection | selon acquisition future |
| Ouvrir Session | analyste | CAP-INV-702 | 0/2 | Intake ready ou limitations acceptées | Session liée | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Vérifier complétude des préconditions | oui | oui | oui | explication | checklist |
| Proposer scope | oui | règles/profils | oui | suggestion | formulaire/profils |
| Proposer plateforme candidate | oui | metadata déclarée | oui | suggestion | comparaison metadata |
| Détecter restriction/permission manquante | oui | oui | oui | non nécessaire | policy/permission checks |
| Ouvrir une Session | humain | oui | workflow | jamais autonome | action explicite |

## 14. États fonctionnels
`draft`, `incomplete`, `ready`, `partial-extraction`, `encrypted`, `locked`, `restricted`, `unsupported`, `source-required`, `permission-blocked`, `out-of-scope`, `superseded`. Aucun état implique que des données absentes n’existent pas.

## 15. États d’interface
Loading conserve origine/scope; Empty demande une source; Partial nomme les zones manquantes; Error préserve le draft; Offline utilise le dernier contexte sûr sans mutation; Permission denied ne révèle aucune donnée; Stale exige revalidation des permissions/source.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Mobile Forensics Intake | concept Investigate | CAP-INV-702/703/704 | origine, objectif, source, scope, restrictions et return origin |
| Readiness assessment | résultat fonctionnel | analyste/reviewer | aucune complétude supposée |
| Collection gap proposal | package | CAP-INV-202 / Collection | lacune, source, scope et autorité requise; aucune exécution |
| Intake activity | événement | Shared Trace/Timeline | acteur, source, permission, version et disposition |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Case/Incident/Finding/Hunt/Signal | besoin Mobile | CAP-INV-701 | objectif, refs, période, restrictions | origine |
| Collection Job/Artifact/package | résultat disponible | CAP-INV-701 | source, custody, erreurs, partialité | job/Artifact |
| Intake ready | ouvrir workspace | CAP-INV-702 | Intake/version, scope, source, return origin | Intake/origine |
| Acquisition inconnue/à revoir | continuer | CAP-INV-704 | source, méthode déclarée, restrictions | Intake |
| Donnée manquante | préparer complément | CAP-INV-202 | besoin, source, scope, raison, autorité | Intake |

## 18. Dépendances
CAP-INV-101..114, CAP-INV-202/203/213/214, Platform Settings, Endpoint Agent projections, Govern authority, Studio Tools/Runs, Shared Linking/Trace, OPEN-005/008/011/013/014/015.

## 19. Source de vérité
L’Intake et son scope analytique restent Investigate. Acquisition/job/custody restent Collection; source/Fleet/policy restent Settings/Endpoint; Tool/Run restent Studio; authority reste Govern.

## 20. Provenance et audit
Enregistrer origine, tenant/environnement, auteur, device/package refs, plateforme candidate et source, représentation, acquisition/custody, période, restrictions, permissions, lock/encryption/partial state, résultats antérieurs, Tools, objectif, scope, modifications et return origin.

## 21. Permissions fonctionnelles
Mobile package read, raw/restricted source read, Intake create/update, platform candidate select, scope update, prior Tool result read, cross-tenant relation read, sensitive presence read et collection-gap prepare. Matrice atomique reportée.

## 22. Limites et erreurs
Source absente, permission refusée, representation unsupported, package partial/corrupted, lock/encryption, stale device metadata, custody lacunaire ou tenant mismatch bloquent ou limitent l’analyse sans effacer le contexte valide.

## 23. Métriques
Intakes par origine/état, time-to-ready, sources manquantes, permission blocks, partial/locked/encrypted rates, scope revisions, collection-gap handoffs et retours origine réussis.

## 24. Classification de livraison
`defined` / `planned`. Aucun moteur, provider, plateforme, acquisition method, format, API, connector, commande ou code livré.

## 25. Critères d’acceptation
**Given** un Collection Job partiel **When** l’analyste ouvre Mobile Intake **Then** la partialité, les éléments manquants, la custody et le return origin sont visibles et aucun succès complet n’est affiché.

**Given** un device déclaré verrouillé **When** l’Intake est évalué **Then** l’état `locked` n’est pas interprété comme absence de données et aucune procédure de déverrouillage n’est proposée.

**Given** aucun modèle IA **When** un Intake est préparé **Then** formulaires, metadata, permission checks et checklist permettent le workflow complet.

## 26. Questions ouvertes
OPEN-011 conserve plateforme/méthode/tool delivery; OPEN-005 moteur forensics; OPEN-008 support source/Endpoint; OPEN-013 classes 2; OPEN-014 material relations; OPEN-015 run provenance.

## 27. Consommateurs documentaires
Mobile module, Case Workspace, Collection, Evidence, Analysis Workbench, Settings, Endpoint, Studio, Govern, Shared, Permissions, Objects, Quality, Roadmap et futurs écrans Mobile conceptuels.
