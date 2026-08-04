---
id: CAP-INV-315
title: Sandbox Environment Selection
product: investigate
module: dynamic-sandbox
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-INV-005
  - REQ-PROD-014
  - REQ-PROD-017
  - REQ-SEC-001
open_decisions:
  - OPEN-005
source-of-truth: canonical
---
# CAP-INV-315 — Sandbox Environment Selection

## 1. Définition
Présenter les Sandbox Environments autorisés, leurs types fonctionnels, profils déclarés, disponibilité, health, version, capacités, limitations, politiques réseau, reset et isolation déclarée, puis permettre une sélection explicite sans administration locale.

## 2. Problème utilisateur
Un environnement peut être incompatible, unhealthy, en reset ou soumis à une politique réseau différente. Le masquer peut produire un Run trompeur ou dangereux.

## 3. Objectifs
- Afficher les environnements autorisés et permission-aware.
- Exposer compatibilité, health, version, capacités, limites et politique réseau.
- Sélectionner explicitement ou demander une alternative.
- Conserver la décision et sa provenance.

## 4. Non-objectifs
Ne pas créer, modifier, réparer ou administrer l’environnement, ses secrets, providers, infrastructure ou policies; ne choisir aucun moteur ou hyperviseur.

## 5. Propriétaire
Platform Settings possède l’administration du Sandbox Environment. Investigate possède seulement la sélection et son contexte analytique.

## 6. Utilisateurs
Malware Analyst; Dynamic Analysis Operator; Case Analyst; Platform Administrator en consultation de handoff.

## 7. Conditions d’entrée
Intake existant; exigences fonctionnelles définies; environnement projeté par Settings; permissions et policy réévaluées.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Exigences de l’Artifact | CAP-INV-314 | compatibilité demandée | oui | version d’intake | aucune sélection |
| Catalogue autorisé | Platform Settings | environnements permission-aware | oui | health/version visibles | environment-unavailable |
| Capacités et limitations | Platform Settings | profil fonctionnel | oui | état courant ou stale explicite | incompatible |
| Politique réseau/isolation/reset | Platform Settings | contraintes de sécurité | oui | version courante | policy-blocked |
| Tools disponibles | Studio/Settings projection | disponibilité fonctionnelle | non | version/status visibles | option indisponible |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Sandbox Environment | Platform Settings | type, version, health, capacités, limitations, réseau et reset | consulter |
| Artifact | Investigate | exigences et restrictions | consulter |
| Tool | Studio | disponibilité compatible | consulter |
| Dynamic Analysis Session | concept Investigate | environnement requis/sélectionné | consulter/modifier localement |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Environment selection record | Investigate | créer/supersede | environment/version et justification visibles |
| Dynamic Analysis Session | Investigate concept | lier sélection | aucune administration Settings |
| Environment request | Platform Settings handoff | préparer | besoin fonctionnel seulement |

## 11. Fonctionnalités
Afficher type, système/profil déclaré, disponibilité, health, version, capacités, limitations, politique réseau, reset, isolation déclarée et Tools disponibles; signaler incompatibilité; sélectionner; demander une alternative.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter environnements | Analyst | Sandbox Environments | 0 | projection autorisée | liste permission-aware | non |
| Sélectionner environnement | Analyst | Selection record | 2 | compatibilité/policy visibles | choix attribué | OPEN-013 |
| Changer sélection | Analyst | Selection record | 2 | justification | ancienne version conservée | OPEN-013 |
| Demander alternative | Analyst | Environment request | 2 | aucun environnement compatible | handoff Settings | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Évaluer compatibilité | oui | règles de capacité/policy | oui | explication | matrice et filtres |
| Proposer un environnement | oui | ranking explicable | oui | suggestion modifiable | sélection manuelle |
| Expliquer indisponibilité | oui | statuts/health | oui | résumé | champs bruts |
| Lancer un Run | humain explicite | jamais par sélection seule | non | jamais autonome | CAP-INV-317 |

Aucune fonction essentielle ne dépend de l’IA; la proposition n’est jamais une sélection cachée.

## 14. États fonctionnels
`available`, `degraded`, `unavailable`, `resetting`, `restricted`, `incompatible`, `retired`, `unknown-health`. Machines finales reportées.

## 15. États d’interface
Loading conserve l’intake; Empty explique l’absence; Partial distingue health stale; Error garde les options valides; Offline est read-only; Permission denied ne révèle aucun environnement interdit.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Environment selection | selection record | CAP-INV-316/317 | version, health et limitations |
| Compatibility assessment | Analysis Result | analyste | raisons et restrictions explicites |
| Environment request | request context | Platform Settings | aucune administration Investigate |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-314 | sélection requise | CAP-INV-315 | Artifact, exigences, limites, policies | intake |
| CAP-INV-315 | environnement sélectionné | CAP-INV-316 | environment/version, health, limitations, Tools | selection |
| CAP-INV-315 | aucun compatible | Platform Settings | exigences, raison, policy, return origin | selection |

## 18. Dépendances
CAP-INV-314/316/317; Settings Sandbox Environments, Health, Secrets et Retention; Studio Tool catalogue; Shared Linking/Trace.

## 19. Source de vérité
Platform Settings reste l’unique source de l’environnement et de son health; Investigate conserve seulement le choix analytique et sa justification.

## 20. Provenance et audit
Enregistrer exigences, options présentées, health/version, incompatibilités, choix, justification, initiateur, proposition automatisée et disposition.

## 21. Permissions fonctionnelles
Sandbox Environment view; restricted environment use; environment version select; environment request; Dynamic Analysis Session update; cross-tenant denied.

## 22. Limites et erreurs
Health stale/inconnu; environnement resetting/retired; capacités insuffisantes; policy réseau incompatible; version retirée; permission révoquée.

## 23. Métriques
Environnements disponibles/incompatibles; sélections modifiées; demandes alternatives; Run bloqués par health/policy.

## 24. Classification de livraison
`defined` / `planned`; aucune implémentation, moteur, hyperviseur, API, protocole ou commande.

## 25. Critères d’acceptation
**Given** plusieurs environnements autorisés
**When** l’analyste ouvre Sandbox Environment Selection
**Then** leurs versions, health, capacités, limitations, politiques réseau et reset sont visibles et la sélection reste explicite

**Given** aucun environnement compatible
**When** l’analyste demande une alternative
**Then** Settings reçoit le besoin fonctionnel sans modification locale d’infrastructure

**Given** aucun modèle IA
**When** l’environnement est sélectionné
**Then** matrice de compatibilité, filtres et revue humaine couvrent le workflow

## 26. Questions ouvertes
Le choix futur des moteurs reste OPEN-005; le modèle final de Sandbox Environment appartient à Settings/Objets.

## 27. Consommateurs documentaires
Dynamic Sandbox, Settings Sandbox Environments, Run Management, Studio Control Room, Permissions et Journeys.
