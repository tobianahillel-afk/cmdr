---
id: CAP-INV-202
title: Collection Scope and Request Preparation
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
  - REQ-SEC-001
open_decisions:
  - OPEN-008
  - OPEN-013
source-of-truth: canonical
---
# CAP-INV-202 — Collection Scope and Request Preparation

## 1. Définition
Préparer une Collection Request bornée, classée et permission-aware depuis un Case.

## 2. Problème utilisateur
Sans préparation structurée, une collecte peut viser la mauvaise cible, être non bornée, ignorer la policy ou être lancée malgré un Endpoint offline ou unsupported.

## 3. Objectifs
- définir objectif, cible, Case, catégories, période, limites et impact conceptuel ;
- afficher capacité, policy, permission, classe et besoin de Govern ;
- enregistrer, valider, soumettre ou lancer selon la classe.

## 4. Non-objectifs
Ne pas administrer Fleet/Policies, définir paramètres techniques finaux, API, protocole, commande, moteur, format ou plateforme supportée, ni exécuter containment localement.

## 5. Propriétaire
Investigate possède la préparation et Collection Request. Platform Settings possède Fleet/Policy ; Endpoint Agent déclare/exécute ses capacités ; Govern possède l’autorité risquée.

## 6. Utilisateurs
Principal : Case Analyst ou Evidence Reviewer. Secondaires : Investigation Lead, DFIR Analyst, Response Operator et reviewer Govern.

## 7. Conditions d’entrée
Case actif, Endpoint résolu, objectif explicite, capacités et policy projetées, limites fonctionnelles, permission de préparation et classification d’action.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Case et objectif | Investigate | contexte et finalité | oui | version courante | rester `incomplete` |
| Endpoint et capacités | Settings / Endpoint Agent | cible et faisabilité | oui | dernière communication | offline/unsupported explicite |
| Catégories, période et limites | analyste / profile | scope borné | oui | validé avant soumission | soumission bloquée |
| Policy, permission et classe | Settings / Security / Govern | autorité | oui | snapshot à l’action | denied ou awaiting-approval |
| Impact et contraintes | Endpoint Agent / policy | coût conceptuel | selon collecte | estimation courante | avertissement ou blocage selon policy |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case / Finding / Evidence | Investigate | objectif, lacunes et justification | consulter et référencer |
| Endpoint / Endpoint Agent | partagé / Endpoint Agent | cible, disponibilité et capacité | consulter |
| Endpoint Policy / Fleet | Platform Settings | restrictions et affectation | consulter uniquement |
| Collection Request | Investigate | draft, version, scope et statut | créer/mettre à jour selon permission |
| Action Request / Decision | Govern | gate éventuel | consulter en projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Collection Request | créer, modifier, supersede, annuler ou soumettre | Investigate | cible, scope, limites, classe et Case obligatoires |
| Request–Case/Endpoint relations | créer/supersede | Investigate / Linking | relations sourcées et tenant-scoped |
| Approval context | préparer | Govern | aucune Decision créée localement |
| Fleet/Policy/Agent | aucune mutation | owner externe | projection seulement |

## 11. Fonctionnalités
Définir collecte d’observation, active, sensible, perturbatrice ou destructive ; sélectionner catégories ; afficher coût/impact, permissions, policy et classe ; contrôler complétude ; sauvegarder draft ; soumettre ou lancer uniquement selon autorité.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Lire capacités/policy | analyste | projections | 0 | read | contexte de faisabilité | non |
| Créer/modifier draft | analyste | Collection Request | 2 | Case/Endpoint | draft versionné | OPEN-013 |
| Valider scope | analyste/reviewer | request | 0 | champs accessibles | lacunes explicites | non |
| Soumettre/lancer collecte bornée | analyste autorisé | request | 1 ou 2 | ready, permission, target | request submitted/accepted | selon impact |
| Préparer classe 3/4 | Investigation Lead | Action Request | 3/4 | Finding/Evidence/impact/rollback | passage CAP-INV-215/113 | obligatoire |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Proposer un scope | oui | profiles/policy | oui | suggestion modifiable | formulaire et profiles |
| Détecter scope trop large | oui | validateurs | oui | explication | limites déterministes |
| Vérifier complétude | oui | oui | oui | explication | checklist |
| Classer l’action | oui | règles | oui | proposition | catalogue de classes |
| Soumettre | humain explicite | permission/gate | workflow possible | jamais autonome | action utilisateur |

Toute automation expose initiateur, version, Automation Run/Tool Calls, sources, scope proposé, disposition humaine et trace. Aucun scope n’est étendu silencieusement.

## 14. États fonctionnels
`draft`, `incomplete`, `ready`, `submitted`, `awaiting-approval`, `accepted`, `rejected`, `superseded`, `cancelled`. Machine finale reportée.

## 15. États d’interface
Loading conserve draft ; Empty exige objectif/target ; Partial nomme capacité/policy absente ; Error préserve le draft ; Offline indique attente possible sans démarrage ; Permission denied ne divulgue rien ; Stale exige revalidation.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Collection Request draft | objet Investigate | analyste et CAP-INV-203 | versionné, borné et Case-scoped |
| Completeness/classification result | validation event | auteur/reviewer | règles et lacunes visibles |
| Approval package éventuel | contexte vers Govern | Review Queue | impact, classe, alternatives et autorité visibles |
| Request activity | événement | Timeline/Trace | acteur, version et disposition |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-201 | préparer collecte | CAP-INV-202 | Case, Endpoint, capacités, policy, permission, objectif | Endpoint Context |
| Collection Request ready | soumettre/lancer | CAP-INV-203 ou Govern | request/version, cible, scope, limites, classe, autorité | Case/Request |
| Evidence Review | demander complément | CAP-INV-202 | Case, Evidence, donnée manquante, target et raison | Evidence Review |
| Classe 3/4 | préparer réponse | CAP-INV-215/113 | Finding, Evidence, Endpoint, impact, rollback | Case |

## 18. Dépendances
CAP-INV-201/203/204..208/108/113/215, Collection Request object, Fleet/Policy/Agent capability projection, Shared validation/Trace/Linking et OPEN-008/013/007.

## 19. Source de vérité
Request et scope métier restent Investigate ; capacité/résultat Agent restent Endpoint Agent ; Policy/Fleet restent Settings ; autorité et Decision restent Govern.

## 20. Provenance et audit
Auteur, Case, Endpoint, objectif, catégories, bornes, profile/version, policy/version, capacité, permission, classe, impact, modifications, validation, soumission et correlation ID.

## 21. Permissions fonctionnelles
Collection prepare/submit, Endpoint/capability read, policy projection read, sensitive scope, cross-tenant/environment, class-2 step-up et containment request. Matrice atomique reportée.

## 22. Limites et erreurs
Scope non borné, target stale/offline/unsupported, catégories incompatibles, policy conflictuelle, impact inconnu, permission refusée, version concurrente ou Govern indisponible gardent un draft explicite.

## 23. Métriques
Drafts atteignant ready, scopes bloqués, demandes par classe, retours Govern, demandes offline, annulations et modifications après suggestion automatique.

## 24. Classification de livraison
`defined` / `planned`. Aucun moteur, protocole, paramètre bas niveau, plateforme ou release n’est déclaré. Promotion conditionnée par objets, permissions et OPEN-008/013.

## 25. Critères d’acceptation
**Given** un Case, Endpoint, policy et limite **When** l’analyste valide le scope **Then** catégorie, classe, impact, permission et besoin Govern sont explicites avant lancement.

**Given** un Endpoint offline **When** la request est soumise **Then** aucun démarrage n’est présenté, l’attente éventuelle et l’annulation sont visibles et le Case reste actif.

**Given** aucun modèle IA **When** une collecte est préparée **Then** formulaire, profiles, policy et validateurs permettent la préparation complète.

## 26. Questions ouvertes
OPEN-008, OPEN-013 et OPEN-007 restent ouvertes. Les champs finaux, états et cardinalités appartiennent à la phase Objets.

## 27. Consommateurs documentaires
Collection Job Management, profiles Triage/File/Memory/Network, Case Workspace, Evidence Review, Govern, Endpoint Agent, Permissions et parcours acquisition.