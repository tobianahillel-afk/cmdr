---
id: CAP-INV-314
title: Dynamic Analysis Intake and Preconditions
product: investigate
module: dynamic-sandbox
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-INV-005
  - REQ-PROD-020
  - REQ-SEC-002
open_decisions:
  - OPEN-005
  - OPEN-013
  - OPEN-014
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-314 — Dynamic Analysis Intake and Preconditions

## 1. Définition
Qualifier l’entrée d’un Artifact dans l’analyse dynamique, vérifier résultats statiques, permissions, policies, risques, environnements compatibles, limites de durée/réseau/interaction et créer ou reprendre une Dynamic Analysis Session sans lancer l’exécution.

## 2. Problème utilisateur
Un Artifact ne doit jamais être exécuté simplement parce qu’il est ouvert. Sans intake, les restrictions, dépendances, risques réseau et incompatibilités peuvent être ignorés.

## 3. Objectifs
- Afficher source, provenance, restrictions, type et résultats statiques existants.
- Conserver Case, Artifact, environnement, versions, limites et return origin.
- Produire des observations et résultats inspectables avec provenance.
- Préparer un handoff explicite sans qualification automatique.

## 4. Non-objectifs
- Ne pas définir virtualisation, instrumentation, hooks, APIs, protocoles, commandes ou schémas techniques.
- Ne pas choisir moteur, hyperviseur, fournisseur ou outil tiers final.
- Ne pas commencer Reverse Engineering, Debugger ou forensics.
- Ne pas exécuter sur un Endpoint réel ni confondre l’analyse avec un Response Run.

## 5. Propriétaire
Investigate possède le contexte et l’interprétation; Studio, Settings, Govern et Shared conservent leurs objets.

## 6. Utilisateurs
- Malware Analyst
- Case Analyst
- Dynamic Analysis Operator
- Evidence Reviewer

## 7. Conditions d’entrée
- Case et Artifact accessibles.
- Résultats statiques et restrictions visibles.
- Sandbox Environment autorisé et health connu.
- Permissions et policies réévaluées.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Artifact source | Investigate / CAP-INV-105 | entrée à exécuter en sandbox | oui | version immuable ou référencée | bloquer le Run |
| Case et objectif | Investigate | contexte analytique | oui | état courant | rester draft |
| Résultats statiques et restrictions | CAP-INV-301..313 | préconditions et limites | oui | versions visibles | marquer incomplete |
| Permissions et policies | Security/Settings | autorité fonctionnelle | oui | réévaluées avant Run | policy-blocked |
| Besoins durée/réseau/interaction | analyste | scope fonctionnel | oui | version d’intake | rester draft |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | scope, Hypothesis et return origin | consulter |
| Artifact | Investigate | source, version, restrictions et résultats statiques | consulter |
| Sandbox Environment | Platform Settings | compatibilité, health, capacités et policies | consulter seulement |
| Tool / Workflow | CMDR Studio | disponibilité et version | consulter seulement |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Dynamic Analysis Session | concept Investigate | créer ou reprendre | aucun schéma final |
| Dynamic analysis intake | Investigate | créer/modifier/supersede | aucun Run démarré |
| Sandbox Run | aucune création active pendant intake | Investigate concept | préparation seulement |

## 11. Fonctionnalités
- Afficher source, provenance, restrictions, type et résultats statiques existants.
- Vérifier environnements compatibles, permissions, policies et risques.
- Définir objectif, durée, limites, besoins réseau et interactions.
- Créer ou reprendre une Dynamic Analysis Session en conservant le return origin.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Ouvrir intake | Analyst | Artifact | 0 | Artifact accessible | préconditions affichées | non |
| Définir limites | Analyst | Intake draft | 2 | permission/policy visibles | draft versionné | OPEN-013 |
| Choisir objectif/profil demandé | Analyst | Intake draft | 2 | options visibles | sélection attribuée | OPEN-013 |
| Créer/reprendre session | Analyst | Dynamic Analysis Session | 2 | intake ready | session liée | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Vérifier les préconditions | oui | règles et policies | oui | explication | checklist et validations |
| Proposer un environnement ou profil | oui | compatibilité déclarée | oui | suggestion modifiable | catalogue manuel |
| Résumer les risques | oui | agrégation sourcée | oui | résumé attribué | champs bruts et filtres |
| Lancer le Run | humain explicite | contrôles seulement | workflow avec gate | jamais autonome | action humaine |

Toute sortie automatisée expose initiateur, producteur/version, Automation Run/Tool Calls, sources, paramètres, timestamp, statut, incertitude, owner humain et disposition. La voie sans IA est complète.

## 14. États fonctionnels
`draft`, `incomplete`, `ready`, `blocked`, `environment-unavailable`, `policy-blocked`, `unsupported`. Machines finales reportées.

## 15. États d’interface
Loading conserve le contexte; Empty explique; Partial nomme les lacunes; Error garde les résultats valides; Offline est stale/read-only; Permission denied ne fuit rien.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Dynamic intake assessment | intake record | CAP-INV-315/316 | préconditions et blocages explicites |
| Session link | Dynamic Analysis Session | Case/Workbench | Case, Artifact et objectif conservés |
| Missing dependency notice | availability event | analyste/Settings/Studio | owner et cause visibles |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Case Workspace | ouvrir analyse dynamique | CAP-INV-314 | tenant, Case, Artifact, Hypothesis, objectif, return origin | Case |
| Static Analysis Session | demander exécution contrôlée | CAP-INV-314 | Artifact, dérivés, résultats statiques, restrictions | Static Analysis |
| CAP-INV-314 | préconditions satisfaites | CAP-INV-316 | Artifact, objectif, limites, owner, policies | intake |

Le return origin, les permissions, versions et résultats partiels sont conservés.

## 18. Dépendances
CAP-INV-105, CAP-INV-301..313, CAP-INV-213/214, Shared Jobs/Trace/Linking, Studio Tools/Automation Runs et Settings Sandbox Environments/Health.

## 19. Source de vérité
Investigate possède le contexte et la disposition analytique. Platform Settings reste source de l’administration des environnements; Studio reste source des Tools et Automation Runs; Govern reste source de Decision/Response Run/Result.

## 20. Provenance et audit
Enregistrer Case, session, Artifact/version, résultats statiques, environnement requis, limites, permissions, policies, initiateur, propositions automatisées, dispositions et correlation IDs.

## 21. Permissions fonctionnelles
Dynamic Analysis read; Dynamic Analysis Session create/update; Sandbox Environment view; sensitive Artifact execute-in-sandbox; network profile select; interaction profile select; cross-tenant denied.

## 22. Limites et erreurs
Artifact inaccessible, superseded, chiffré ou restreint; résultats statiques absents; environnement incompatible; policy ou permission refusée; dépendance manquante; contexte tenant incohérent.

## 23. Métriques
Intakes ready/blocked; dépendances manquantes; divergences statiques; temps jusqu’à session; propositions acceptées/modifiées/rejetées.

## 24. Classification de livraison
`defined` / `planned`; aucune preuve d’implémentation, moteur, hyperviseur, API, protocole ou commande.

## 25. Critères d’acceptation
**Given** un Artifact compatible et un Case actif
**When** l’analyste ouvre Dynamic Analysis Intake and Preconditions
**Then** source, provenance, restrictions, résultats statiques, permissions, policies, risques, limites et return origin sont visibles et aucun Run ne démarre

**Given** un environnement indisponible ou une permission refusée
**When** l’utilisateur demande une session
**Then** l’action est bloquée, les données valides sont conservées et aucun succès n’est inventé

**Given** aucun modèle IA
**When** l’analyste complète l’intake
**Then** règles, profils, formulaires et checklists déterministes couvrent le workflow essentiel

## 26. Questions ouvertes
Dynamic Analysis Session et intake record restent des concepts; OPEN-005/013/014/015 restent ouvertes.

## 27. Consommateurs documentaires
Dynamic Sandbox, Analysis Workbench, Case Workspace, Settings Sandbox Environments, Studio Control Room et phases Objets/Permissions/Journeys/Technique.
