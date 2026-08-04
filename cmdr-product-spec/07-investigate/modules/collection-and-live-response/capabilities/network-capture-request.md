---
id: CAP-INV-208
title: Network Capture Request
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
  - REQ-PROD-055
open_decisions:
  - OPEN-008
  - OPEN-013
source-of-truth: canonical
---
# CAP-INV-208 — Network Capture Request

## 1. Définition
Préparer, démarrer, suivre et arrêter une capture réseau bornée par cible, durée, volume et filtre conceptuel, sans choisir de moteur ou format.

## 2. Problème utilisateur
Une capture réseau non bornée peut consommer des ressources, dépasser le Case ou collecter des données sensibles. L’analyste doit voir le scope, l’impact et les pertes avant et après l’exécution.

## 3. Objectifs
- sélectionner l’Endpoint et une interface ou un scope conceptuel autorisé
- définir durée, limite de volume et filtre conceptuel
- afficher policy, permission, impact et capacité déclarée
- arrêter la capture et recevoir un Artifact avec pertes/erreurs visibles

## 4. Non-objectifs
Ne pas définir interface bas niveau, moteur, syntaxe de filtre, commande, protocole ou format final ; ne pas promettre le support de toutes les plateformes.

## 5. Propriétaire
Investigate / Collection and Live Response / Investigate Product Lead possède le contexte métier, les drafts et les relations au Case. Platform Settings reste propriétaire de Fleet et Endpoint Policies ; Endpoint Agent exécute et rapporte localement ; Govern conserve l’autorité, Decision, Response Run et Result.

## 6. Utilisateurs
Principal : DFIR Analyst. Secondaires : Case Analyst, Response Operator et Evidence Reviewer.

## 7. Conditions d’entrée
Case/Endpoint accessibles, network capture capability déclarée, scope/durée/limite définis, policy/permission validées et impact conceptuel affiché.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Case et objectif | Investigate | contexte | oui | version courante | rester draft |
| Endpoint/Agent et capability | Endpoint Agent | cible et support | oui | dernière communication | offline/unsupported |
| Interface ou scope conceptuel | analyste | périmètre réseau | oui | revalidé au lancement | incomplete |
| Durée, volume et filtre conceptuel | analyste / policy | bornes | oui | snapshot au lancement | scope refusé |
| Policy, permission et impact | Settings / Security / Govern | gate et coût | oui | version/snapshot visibles | denied/policy-blocked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | objectif et restrictions | consulter |
| Endpoint / Endpoint Agent | partagé / Endpoint Agent | interface/scope disponible et état | consulter |
| Endpoint Policy | Platform Settings | capture permise, limites et sensibilité | consulter uniquement |
| Collection Request / Job | Investigate / concept futur | scope et progression | préparer/suivre |
| Artifact | Investigate | captures reçues et versions | consulter/lier |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Network capture request context | créer/modifier/supersede | Investigate | durée, volume, filtre et cible obligatoires |
| Capture status/result | enregistrer start/stop/partial/loss/error | Investigate, modèle futur | pertes et durée réelle visibles |
| Capture Artifact | créer à réception | Investigate | source, période, scope et erreurs conservés |
| Endpoint/Policy | aucune mutation | owners externes | projection uniquement |

## 11. Fonctionnalités
- sélectionner cible et scope conceptuel sans exposer de paramètres bas niveau
- définir durée, volume et filtre conceptuel bornés
- afficher impact et autorité avant démarrage
- arrêter ou annuler la capture lorsque permis
- recevoir l’Artifact avec durée réelle, pertes, erreurs et fraîcheur

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Prévisualiser scope/impact | DFIR Analyst | request draft | 0 | read permission | bornes visibles | non |
| Modifier durée/limite/filtre | DFIR Analyst | request draft | 2 | draft modifiable | nouvelle version | OPEN-013 |
| Démarrer capture | analyste autorisé | Collection Request | 1/2 | ready, capability/policy | capture Job lié | selon impact |
| Arrêter capture | Response Operator | capture Job | 2 | running et permission | stop demandé | OPEN-013 selon policy |
| Ouvrir Artifact | analyste | capture Artifact | 0 | read | Artifact Detail | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| Proposer bornes | oui | defaults/policy | oui | suggestion modifiable | saisie manuelle |
| Détecter filtre trop large | oui | validateur conceptuel | oui | explication | règles |
| Suivre volume/durée | oui | compteurs déclarés | oui | résumé | progression brute |
| Expliquer pertes/erreurs | oui | catalogue déterministe | oui | oui | détails bruts |
| Prolonger automatiquement | non | interdit | non | interdit | action humaine |

Toute sortie automatisée expose initiateur, producteur/version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres fonctionnels, timestamp, statut, incertitude, owner humain, acceptation/modification/rejet et trace.

## 14. États fonctionnels
`draft`, `validating`, `awaiting-approval`, `queued`, `capturing`, `stopping`, `partial`, `completed`, `failed`, `cancelled`, `unsupported`. Machine finale reportée.

## 15. États d’interface
Loading conserve bornes ; Empty distingue scope absent et capture vide ; Partial affiche pertes/segments ; Error conserve l’Artifact valide ; Offline n’indique pas capturing ; Permission denied masque le contenu ; Stale expose dernière mise à jour.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Capture request/status | Collection context | Case Workspace | cible, durée, limites et statut visibles |
| Capture Artifact | Artifact | CAP-INV-105/107 et future 4B.2B | scope, période, pertes et provenance conservés |
| Loss/error summary | Operation/collection result | CAP-INV-212 | aucune perte masquée |
| Progress/stop events | Background Job/Trace events | Notifications/Timeline | acteur et correlation ID |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Endpoint Context | préparer capture | CAP-INV-208 | Case, Endpoint, capability, policy, reason | Endpoint Context |
| CAP-INV-208 | démarrer | CAP-INV-203 | request, durée, volume, filtre, autorité | CAP-INV-208/Case |
| Capture running | arrêter | CAP-INV-203/208 | job, initiateur, raison, timestamp | CAP-INV-208 |
| Capture reçue | enregistrer | CAP-INV-105 | Artifact, source, période, pertes, errors | CAP-INV-208 |

## 18. Dépendances
CAP-INV-201/202/203/105/107/212/213/214, Endpoint Agent network capture capability, Settings Policy, Shared Jobs/Trace et OPEN-008/013.

## 19. Source de vérité
Investigate possède la request, le statut métier et l’Artifact. Endpoint Agent reste source de l’exécution et des pertes/erreurs locales ; Settings reste source de la Policy.

## 20. Provenance et audit
Case, Endpoint/Agent, scope/interface conceptuels, durée/volume/filtre, policy, initiateur, start/stop, pertes, erreurs, Artifact et correlation IDs.

## 21. Permissions fonctionnelles
Endpoint/capability read, network capture prepare/start/stop/cancel, raw result read, capture Artifact receive/export et sensitive network output read.

## 22. Limites et erreurs
Scope/interface unsupported, Endpoint offline, filtre refusé, limite atteinte, pertes, arrêt tardif, partial, permission révoquée ou policy conflict. Aucun moteur ou paramètre bas niveau n’est inféré.

## 23. Métriques
Captures par statut, durée/volume demandés et réalisés, pertes, partial/failure rate, arrêts, Artifacts reçus et scopes bloqués.

## 24. Classification de livraison
`defined` / `planned`. Aucun moteur, format, commande, protocole ou support plateforme n’est prouvé.

## 25. Critères d’acceptation
**Given** une durée, une limite et un filtre conceptuel **When** la capture démarre **Then** le scope reste borné et l’arrêt, la progression et l’impact sont visibles.

**Given** des pertes ou erreurs **When** la capture se termine **Then** elles sont affichées et le résultat n’est pas présenté comme complet.

**Given** aucun modèle IA **When** la capture est préparée **Then** formulaire, bornes, policy et validateurs permettent le workflow.

## 26. Questions ouvertes
OPEN-008 conserve le support plateforme ; OPEN-013 la gouvernance du start/stop classe 2. Moteur, filtre bas niveau et format restent hors phase.

## 27. Consommateurs documentaires
Collection Job, Artifact Management, Custody, future Network Artifact Analysis, Endpoint Agent network capability et phases Objets/Permissions/Technique.
