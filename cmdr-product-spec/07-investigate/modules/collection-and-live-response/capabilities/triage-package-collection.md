---
id: CAP-INV-204
title: Triage Package Collection
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
# CAP-INV-204 — Triage Package Collection

## 1. Définition
Demander un package de triage profilé, borné et Case-scoped, puis traiter séparément les résultats et erreurs de chaque catégorie.

## 2. Problème utilisateur
Une collecte de triage monolithique peut masquer les catégories non exécutées, produire un faux succès global ou récupérer plus de données que nécessaire. L’analyste doit choisir un profil explicite et comprendre exactement ce qui a réussi.

## 3. Objectifs
- sélectionner un profil de triage versionné et les catégories réellement demandées
- afficher cible, Case, limites, impact conceptuel, policy, permission et classe
- suivre chaque catégorie indépendamment et permettre un retry ciblé
- enregistrer les Artifacts réussis sans convertir automatiquement leur contenu en Evidence

## 4. Non-objectifs
Ne pas définir le format du package, les chemins, commandes, outils, compression, transport ou manifeste technique ; ne pas déclarer une plateforme supportée ; ne pas exécuter de containment.

## 5. Propriétaire
Investigate / Collection and Live Response / Investigate Product Lead possède le contexte métier, les drafts et les relations au Case. Platform Settings reste propriétaire de Fleet et Endpoint Policies ; Endpoint Agent exécute et rapporte localement ; Govern conserve l’autorité, Decision, Response Run et Result.

## 6. Utilisateurs
Principal : Case Analyst ou DFIR Analyst. Secondaires : Investigation Lead, Evidence Reviewer et Response Operator.

## 7. Conditions d’entrée
Case actif, Endpoint résolu, profil disponible pour la capacité déclarée, catégories et limites sélectionnées, policy/permission vérifiées et classe d’action déterminée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Case et besoin de triage | Investigate | objectif et contexte | oui | version courante | rester draft |
| Endpoint, Agent et capacités | Settings / Endpoint Agent | cible et faisabilité | oui | dernière communication | offline ou unsupported |
| Profil et catégories | catalogue/policy | scope de triage | oui | version du profil | aucun dispatch |
| Bornes et impact conceptuel | analyste / policy | limites de collecte | oui | revalidés avant soumission | request incomplete |
| Autorité et permission | Security / Govern | gate d’exécution | selon classe | snapshot à l’action | denied ou awaiting-approval |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | objectif, scope et restrictions | consulter et référencer |
| Endpoint / Endpoint Agent | partagé / Endpoint Agent | disponibilité et catégories déclarées | consulter |
| Endpoint Policy | Platform Settings | catégories permises, limites et restrictions | consulter uniquement |
| Collection Request / Collection Job | Investigate / concept futur | profil, catégories, statut et résultats | préparer et suivre |
| Artifact | Investigate | sorties déjà reçues et métadonnées d’acquisition | consulter et lier |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Triage request profile selection | créer, modifier ou supersede | Investigate | profil/version, catégories et limites obligatoires |
| Collection Job category result | enregistrer résultat ou erreur par catégorie | Investigate, modèle futur | aucune agrégation ne masque un échec |
| Artifact relation | créer pour une catégorie réussie | Investigate | origine, catégorie, acquisition et Case conservés |
| Endpoint Policy / Fleet | aucune mutation | Platform Settings | projection en lecture seule |

## 11. Fonctionnalités
- choisir un profil complet ou un sous-ensemble de catégories autorisées
- prévisualiser catégories, limites et impact conceptuel avant soumission
- afficher progression et résultat pour système, processus, connexions, comptes/sessions, persistence indicators, journaux et fichiers ciblés
- conserver les catégories non disponibles comme unsupported ou policy-blocked
- annuler la collecte ou relancer uniquement les catégories échouées lorsque permis

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Choisir un profil | DFIR Analyst | profil/version | 0 | catalogue accessible | scope explicite | non |
| Retirer ou ajouter une catégorie autorisée | Case Analyst | Collection Request draft | 2 | policy et limites | nouvelle version du draft | OPEN-013 |
| Soumettre la collecte de triage | analyste autorisé | Collection Request | 1 | request ready et Endpoint compatible | Collection Job créé/lié | selon impact |
| Annuler avant fin | analyste autorisé | Collection Job | 2 | job cancellable | annulation demandée et tracée | OPEN-013 selon policy |
| Relancer une catégorie | DFIR Analyst | résultat échoué | 1 | cause connue et retry permis | tentative ciblée liée | selon impact |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| Proposer un profil | oui | règles par objectif/policy | oui | suggestion modifiable | catalogue et choix manuel |
| Détecter un scope excessif | oui | limites déterministes | oui | explication facultative | validateur de bornes |
| Suivre les catégories | oui | états par catégorie | oui | résumé | table de résultats brute |
| Proposer un retry ciblé | oui | règles sur erreurs | oui | suggestion | sélection manuelle |
| Lancer ou étendre la collecte | humain explicite | permission/gate | workflow possible | jamais autonome | action utilisateur |

Toute sortie automatisée expose initiateur, producteur/version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres fonctionnels, timestamp, statut, incertitude, owner humain, acceptation/modification/rejet et trace.

## 14. États fonctionnels
`draft`, `validating`, `queued`, `collecting`, `partial`, `completed`, `failed`, `cancelled`, `unsupported`, `policy-blocked`. La machine d’état finale du Job est reportée.

## 15. États d’interface
Loading conserve profil et catégories ; Empty distingue profil absent et catégorie non sélectionnée ; Partial affiche chaque catégorie ; Error conserve les Artifacts valides ; Offline indique attente ou échec sans faux démarrage ; Permission denied masque les données ; Stale exige revalidation du profil.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Résultat de catégorie | Collection result entry | CAP-INV-203/212 | statut, erreur, timestamps et scope visibles |
| Artifact de triage | Artifact | CAP-INV-105/107 | catégorie, source et acquisition conservées |
| Résumé de package | projection métier | Case Workspace | succès partiels jamais présentés comme succès complet |
| Événements de progression | Background Job/Trace events | Notifications/Timeline | correlation ID et deep link |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-202 | profil validé | CAP-INV-203/204 | request, profil/version, catégories, bornes, autorité | Collection Request ou Case |
| Catégorie réussie | sortie reçue | CAP-INV-105 | source, catégorie, acquisition, timestamps, erreurs | Triage Job |
| Package partial | retry ciblé | CAP-INV-203 | catégorie échouée, cause, limites, autorité | même groupe de Job |
| Artifact | qualification humaine | CAP-INV-107 | Artifact, Case, provenance, raison | Artifact Detail |

## 18. Dépendances
CAP-INV-202/203/105/107/212/213/214, Endpoint Agent collection, Settings Policy/Health, Shared Background Jobs/Notifications/Trace et OPEN-008/013.

## 19. Source de vérité
Investigate possède le profil choisi, la request, le statut métier par catégorie et les relations au Case. Endpoint Agent reste source de l’exécution locale ; Settings reste source de la Policy ; Shared reste source du mécanisme Background Job.

## 20. Provenance et audit
Case, Endpoint/Agent, profil/version, catégories, bornes, policy/version, initiateur, autorité, dispatch, résultat/erreur par catégorie, retries, Artifacts et correlation IDs.

## 21. Permissions fonctionnelles
Endpoint/capability read, collection prepare/submit/cancel/retry, raw category result read, Artifact receive et sensitive-output read. Step-up et séparation des tâches restent à la phase Permissions.

## 22. Limites et erreurs
Profil absent ou incompatible, catégorie unsupported, Endpoint offline, policy blocked, scope trop large, timeout, catégorie partielle, résultat tardif ou permission révoquée. Aucun package complet n’est déclaré si une catégorie requise échoue.

## 23. Métriques
Packages par profil, catégories réussies/échouées, taux partial, retries ciblés, annulations, Artifacts par catégorie et scopes réduits après validation.

## 24. Classification de livraison
`defined` / `planned`. Aucun format, moteur, commande, protocole ou support plateforme n’est déclaré livré.

## 25. Critères d’acceptation
**Given** une collecte comportant plusieurs catégories **When** certaines réussissent et d’autres échouent **Then** l’état est `partial`, chaque catégorie expose son résultat et seuls les Artifacts réussis sont liés.

**Given** une catégorie bloquée par policy **When** l’analyste valide le profil **Then** la catégorie et la raison sont visibles et aucun dispatch silencieux n’a lieu.

**Given** aucun modèle IA **When** un package est préparé **Then** le catalogue, les limites, la checklist et le suivi déterministe permettent le workflow complet.

## 26. Questions ouvertes
OPEN-008 conserve le support plateforme ; OPEN-013 la gouvernance de certaines mutations/retries. Le format du package et les catégories finales restent aux phases Technique/Objets.

## 27. Consommateurs documentaires
Collection Job Management, Artifact Management, Evidence Review, Case Workspace, Endpoint Agent collection, parcours Endpoint investigation et phases Objets/Permissions/Technique.
