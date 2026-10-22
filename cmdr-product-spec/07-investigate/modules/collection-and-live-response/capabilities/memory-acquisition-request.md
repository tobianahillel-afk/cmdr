---
id: CAP-INV-207
title: Memory Acquisition Request
product: investigate
module: collection-and-live-response
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-PROD-052
  - REQ-PROD-055
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
source-of-truth: canonical
---
# CAP-INV-207 — Memory Acquisition Request

## 1. Définition
Préparer, autoriser et suivre une demande d’acquisition mémoire au niveau fonctionnel, sans sélectionner de moteur, format ou plateforme supportée.

## 2. Problème utilisateur
L’acquisition mémoire peut être indisponible, coûteuse ou perturbatrice. L’utilisateur doit voir les préconditions, l’impact conceptuel et le support déclaré avant toute demande, puis distinguer échec, partial et Artifact reçu.

## 3. Objectifs
- définir Case, cible, objectif et périmètre d’acquisition mémoire
- vérifier capacité déclarée, plateforme, état Agent, policy et permissions
- afficher contraintes conceptuelles de durée, taille et impact sans paramètre technique final
- suivre la demande et recevoir un Memory Image/Artifact avec provenance

## 4. Non-objectifs
Ne pas choisir moteur, format, méthode, commande, chiffrement ou stockage ; ne pas promettre le support d’une plateforme ; ne pas commencer Memory Forensics Workbench.

## 5. Propriétaire
Investigate / Collection and Live Response / Investigate Product Lead possède le contexte métier, les drafts et les relations au Case. Platform Settings reste propriétaire de Fleet et Endpoint Policies ; Endpoint Agent exécute et rapporte localement ; Govern conserve l’autorité, Decision, Response Run et Result.

## 6. Utilisateurs
Principal : DFIR Analyst. Secondaires : Investigation Lead, Evidence Reviewer et reviewer Govern.

## 7. Conditions d’entrée
Case et Endpoint accessibles, capability mémoire déclarée, support actuel explicitement available/unsupported, policy et permission vérifiées, impact conceptuel accepté ou gouverné.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Case et objectif | Investigate | raison et contexte | oui | version courante | rester draft |
| Endpoint, plateforme et Agent | Settings / Endpoint Agent | cible et support déclaré | oui | dernière communication/version | unsupported/offline |
| Scope mémoire fonctionnel | DFIR Analyst | type de demande et limites | oui | validé au lancement | incomplete |
| Impact, taille et durée conceptuels | Endpoint Agent / policy | prévision de coût | oui si disponible | estimation courante | warning ou blocage |
| Policy, permission et autorité | Settings / Security / Govern | gate | oui | snapshot à la soumission | denied/awaiting-approval |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case / Finding / Evidence | Investigate | besoin, contexte et justification | consulter/référencer |
| Endpoint / Endpoint Agent | partagé / Endpoint Agent | plateforme, état et capability mémoire | consulter |
| Endpoint Policy | Platform Settings | restrictions et impact permis | consulter uniquement |
| Collection Request / Job | Investigate / concept futur | scope, statut et progression | préparer/suivre |
| Memory Image / Artifact | Endpoint Agent / Investigate | sortie et métadonnées | consulter/lier |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Memory acquisition request context | créer, modifier ou supersede | Investigate | support, scope, impact et Case obligatoires |
| Collection Job relation | créer/actualiser | Investigate, modèle futur | statut métier distinct du Background Job |
| Memory Image/Artifact relation | créer à réception | Investigate | source, acquisition, taille déclarée et erreurs conservées |
| Platform support record | aucune mutation | Endpoint Agent / Settings | projection sous OPEN-008 |

## 11. Fonctionnalités
- afficher available, unsupported, degraded ou policy-blocked avant la demande
- présenter les contraintes conceptuelles de durée, taille et impact
- soumettre, annuler et suivre progression/erreurs lorsque permis
- conserver résultats partiels et absence de support sans créer un faux Artifact
- ouvrir le futur Memory Forensics Workbench uniquement comme handoff

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Consulter support/impact | DFIR Analyst | capability projection | 0 | read permission | faisabilité visible | non |
| Modifier scope/limites | DFIR Analyst | request draft | 2 | draft modifiable | nouvelle version | OPEN-013 |
| Soumettre l’acquisition | analyste autorisé | Collection Request | 1/2 | support available, policy/gate | Collection Job lié | selon impact |
| Annuler | analyste autorisé | Collection Job | 2 | cancellable | cancel demandé | OPEN-013 selon policy |
| Ouvrir l’Artifact | DFIR Analyst | Memory Image/Artifact | 0 | read permission | handoff futur | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| Proposer un scope | oui | profiles/policy | oui | suggestion modifiable | formulaire manuel |
| Vérifier support | oui | capability inventory | oui | explication | matrice déclarée |
| Estimer impact conceptuel | oui | règles déclarées | oui | résumé | valeurs/limites brutes |
| Expliquer une erreur | oui | catalogue d’erreurs | oui | oui | message déterministe |
| Lancer malgré unsupported | non | interdit | non | interdit | aucune exécution |

Toute sortie automatisée expose initiateur, producteur/version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres fonctionnels, timestamp, statut, incertitude, owner humain, acceptation/modification/rejet et trace.

## 14. États fonctionnels
`draft`, `unsupported`, `policy-blocked`, `awaiting-approval`, `queued`, `acquiring`, `partial`, `completed`, `failed`, `cancelled`. Machine finale reportée.

## 15. États d’interface
Loading conserve cible/scope ; Empty distingue support inconnu et aucun Artifact ; Partial expose parties reçues/erreurs ; Error conserve la request ; Offline n’indique pas acquiring ; Permission denied masque les sorties ; Stale revalide support et impact.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Memory acquisition request | Collection Request context | CAP-INV-203/207 | support, scope, impact et autorité visibles |
| Progress/status | Collection Job projection | Case Workspace | partial/offline/unsupported explicites |
| Memory Image/Artifact | Artifact | CAP-INV-105 puis future 4B.2B | source/acquisition et limites conservées |
| Error/unsupported event | business event | Timeline/Trace | aucun support ou résultat inventé |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Endpoint Context | préparer mémoire | CAP-INV-207 | Case, Endpoint, plateforme, capability, policy, return origin | Endpoint Context |
| CAP-INV-207 | request autorisée | CAP-INV-203 | request/version, impact, scope, autorité | CAP-INV-207/Case |
| Acquisition result | Artifact reçu | CAP-INV-105 | Memory Image, source, acquisition, errors, Case | CAP-INV-207 |
| Memory Artifact | ouvrir analyse future | Phase 4B.2B future | Artifact/version, restrictions, Case | Artifact Detail |

## 18. Dépendances
CAP-INV-201/202/203/105/107/213/214, Endpoint Agent memory capability, Settings Policy/Health, Shared Jobs/Trace, OPEN-005/008/013.

## 19. Source de vérité
Investigate possède la request et les relations au Case/Artifact. Endpoint Agent reste source du support déclaré et de l’exécution locale. Aucun moteur forensic n’est choisi ; OPEN-005 et OPEN-008 restent ouverts.

## 20. Provenance et audit
Case, Endpoint/Agent/platform, support/version, scope, impact, policy, permission, initiateur, autorité, progression, erreurs, Memory Image/Artifact et correlation IDs.

## 21. Permissions fonctionnelles
Endpoint/capability read, memory acquisition prepare/submit/cancel, raw result read, Memory Artifact receive/export et sensitive output read. Step-up selon impact reste à finaliser.

## 22. Limites et erreurs
Plateforme unsupported, Agent degraded/offline, impact non calculable, espace/temps conceptuellement insuffisant, permission refusée, timeout, partial ou Artifact indisponible. Aucune promesse universelle.

## 23. Métriques
Demandes unsupported/policy-blocked, durée conceptuelle/réelle disponible, partial/failure rate, annulations et Memory Artifacts avec provenance complète.

## 24. Classification de livraison
`defined` / `planned`. Aucun moteur, format, plateforme ou protocole n’est choisi. OPEN-005 et OPEN-008 bloquent toute déclaration plus forte.

## 25. Critères d’acceptation
**Given** une plateforme non supportée **When** l’utilisateur prépare l’acquisition **Then** l’indisponibilité est visible, aucune demande invalide ne démarre et le Case reste accessible.

**Given** une acquisition partielle **When** le résultat arrive **Then** les parties disponibles, erreurs et Artifact éventuel sont séparés.

**Given** aucun modèle IA **When** la demande est préparée **Then** support déclaré, formulaire, policy et validateurs suffisent.

## 26. Questions ouvertes
OPEN-008 conserve les plateformes initiales ; OPEN-005 les moteurs forensics futurs ; OPEN-013 la gouvernance classe 2. Le format et l’exécution technique sont hors phase.

## 27. Consommateurs documentaires
Collection Job, Artifact Management, Custody, Evidence, future Memory Forensics Workbench, Endpoint Agent memory collection et phases Objets/Permissions/Technique.
