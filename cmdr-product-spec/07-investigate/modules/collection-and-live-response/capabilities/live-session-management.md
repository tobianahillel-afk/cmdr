---
id: CAP-INV-209
title: Live Session Management
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
  - REQ-SEC-002
open_decisions:
  - OPEN-007
  - OPEN-008
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-209 — Live Session Management

## 1. Définition
Demander, ouvrir, rejoindre, superviser, suspendre, reconnecter et fermer une Live Session visible, Case-scoped et limitée dans le temps.

## 2. Problème utilisateur
Une session invisible ou assimilée à un terminal, une Automation Run ou un Response Run rend l’autorité et les actions impossibles à auditer. Chaque participant doit connaître la cible, la policy, l’expiration et l’état réel de connexion.

## 3. Objectifs
- ouvrir une session liée à un Case et un Endpoint après vérification de disponibilité, policy et permission
- afficher initiateur, participants autorisés, début, expiration et inactivité
- gérer prolongation, suspension, déconnexion, reconnexion, conflit et révocation
- conserver transcript/trace et retour exact au Case

## 4. Non-objectifs
Ne pas définir shell, terminal, protocole, transport, PKI ou commande ; ne pas confondre la session avec Automation Run ou Response Run ; ne pas rendre une session invisible.

## 5. Propriétaire
Investigate / Collection and Live Response / Investigate Product Lead possède le contexte métier, les drafts et les relations au Case. Platform Settings reste propriétaire de Fleet et Endpoint Policies ; Endpoint Agent exécute et rapporte localement ; Govern conserve l’autorité, Decision, Response Run et Result.

## 6. Utilisateurs
Principal : Response Operator. Secondaires : Investigation Lead, Case Analyst, reviewer Govern et participant autorisé.

## 7. Conditions d’entrée
Case/Endpoint accessibles, Agent available ou état explicite, session reason, policy, classe, permission, participants et expiration définis ; gate Govern/Human Gate appliqué lorsque requis.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Case, Endpoint et raison | Investigate | contexte de session | oui | versions courantes | request incomplete |
| Agent availability/capabilities | Endpoint Agent | faisabilité | oui | dernière communication | offline/unsupported |
| Policy, classe et permission | Settings / Security / Govern | autorité | oui | snapshot à l’ouverture | denied/awaiting-approval |
| Participants et rôles | Identity / initiateur | collaboration autorisée | oui | revalidés à l’entrée | join interdit |
| Expiration et inactivity limit | policy / initiateur | limites temporelles | oui | valeurs à l’ouverture | session non ouverte |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | objectif, owner et restrictions | consulter/lier |
| Endpoint / Endpoint Agent | partagé / Endpoint Agent | état et capability live response | consulter |
| Endpoint Policy | Platform Settings | session permise, durée et restrictions | consulter uniquement |
| Live Session concept | Investigate, modèle futur | participants, statut, expiration et transcript refs | gérer selon permission |
| Decision / Response Run | Govern | autorité éventuelle, distincte de la session | consulter uniquement |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Live Session record conceptuel | créer, actualiser, suspendre, prolonger, fermer ou révoquer | Investigate, modèle futur | visible, Case-scoped et expiration obligatoire |
| Participant relation | ajouter/retirer selon autorisation | Investigate / Identity | aucun accès implicite |
| Session transcript/trace relation | créer/supersede | Investigate / Shared mechanisms | horodatage, acteur et redaction |
| Decision/Response Run | aucune création | Govern | projection seulement |

## 11. Fonctionnalités
- demander et ouvrir une session avec raison et expiration
- afficher initiateur, participants, statut, début, inactivité et capacité
- rejoindre, prolonger, suspendre, reprendre, fermer ou révoquer selon permission
- gérer Endpoint offline, reconnecting et conflit de sessions
- conserver transcript/trace, opérations liées et return origin

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Demander une session | Response Operator | Live Session request | 2 | Case/Endpoint/reason/policy | requested ou awaiting-approval | OPEN-007/013 selon policy |
| Ouvrir/rejoindre | participant autorisé | Live Session | 2 | approval/permission et Endpoint available | active avec participant visible | selon policy |
| Prolonger | Session owner | Live Session | 2 | active et durée permise | nouvelle expiration auditée | OPEN-013 |
| Suspendre/reprendre | Session owner | Live Session | 2 | état compatible | transition visible | OPEN-013 selon policy |
| Fermer/révoquer | owner ou autorité | Live Session | 2 | permission et raison si revoke | closed/revoked, transcript conservé | selon policy |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| Préremplir la raison | oui | template Case | oui | brouillon | saisie manuelle |
| Vérifier participants/expiration | oui | policy/règles | oui | explication | checklist |
| Résumer le transcript | oui | agrégation | oui | résumé attribué | transcript filtré |
| Signaler inactivité/conflit | oui | timers/règles | oui | explication | indicateurs déterministes |
| Ouvrir/prolonger automatiquement | non | interdit sans contrat/autorité | non par défaut | jamais autonome | action humaine explicite |

Toute sortie automatisée expose initiateur, producteur/version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres fonctionnels, timestamp, statut, incertitude, owner humain, acceptation/modification/rejet et trace.

## 14. États fonctionnels
`requested`, `awaiting-approval`, `opening`, `active`, `idle`, `suspended`, `reconnecting`, `closing`, `closed`, `expired`, `failed`, `revoked`. Machine finale reportée.

## 15. États d’interface
Loading conserve Case/Endpoint ; Empty signifie aucune session ; Partial nomme participant/capability manquants ; Error garde transcript/opérations valides ; Offline montre reconnecting/closed sans succès implicite ; Permission denied ne divulgue pas le transcript ; Stale indique dernière activité.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Live Session status | record conceptuel | Case Workspace/CAP-INV-210 | participants, timing, policy et état visibles |
| Transcript/trace | relation/événements | CAP-INV-214/Timeline/Audit | aucune session invisible |
| Session notification | Notification event | participants | deep link, expiration et changement d’état |
| Closure context | business event | Case Workspace | raison, dernière opération et return origin |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Case/Endpoint Context | request session | CAP-INV-209 | Case, Endpoint, reason, class, policy, participants, expiration | Case/Endpoint Context |
| Live Session active | sélectionner opération | CAP-INV-210 | session, participant, Endpoint, permissions, trace | Live Session |
| Live Session | transférer fichier | CAP-INV-211 | session, source/destination, policy, class | Live Session |
| Live Session closed | retour | Case Workspace | session, transcript, results, next action, return origin | même Case |

## 18. Dépendances
CAP-INV-201/210/211/212/214, Endpoint Agent session capability, Settings Policy/Health, Identity, Shared Notifications/Trace/Recovery, Govern/Studio boundaries et OPEN-007/008/013/015.

## 19. Source de vérité
Investigate possède le contexte et le record métier de session. Endpoint Agent reste source de la connexion/exécution locale ; Settings de la Policy ; Govern de l’autorité ; Studio d’Automation Run. Une session n’est aucun de ces Runs.

## 20. Provenance et audit
Case, Endpoint/Agent, initiateur, participants, permissions, policy/version, approval éventuelle, start/end, expiration, inactivity, reconnects, operations, transcript, closure et correlation IDs.

## 21. Permissions fonctionnelles
Live Session request/open/join/extend/suspend/resume/close, transcript read, participant management, sensitive output et cross-tenant restrictions. Step-up et séparation des tâches restent ouvertes.

## 22. Limites et erreurs
Endpoint offline, capability unsupported, session concurrente, participant non autorisé, expiration, inactivity, reconnect failure, revocation, policy change ou permission retirée. Déconnexion ≠ succès d’opération.

## 23. Métriques
Requests/open failures, session duration, reconnects, conflicts, expirations, revoked sessions, operations/session et retours Case réussis.

## 24. Classification de livraison
`defined` / `planned`. Aucun protocole, terminal, transport, plateforme ou implementation n’est déclaré livré.

## 25. Critères d’acceptation
**Given** aucun fournisseur de modèle, un Endpoint available et une policy permise **When** un opérateur ouvre une session **Then** la session fonctionne manuellement, les classes sont visibles, le transcript est conservé et elle peut être fermée.

**Given** une déconnexion pendant une opération **When** la session passe en reconnecting **Then** l’opération n’est pas déclarée réussie et son statut reste indépendant.

**Given** un participant non autorisé **When** il tente de rejoindre **Then** l’accès est refusé sans révéler le transcript et la session reste visible aux participants autorisés.

## 26. Questions ouvertes
OPEN-007 traite Human Gate/Govern ; OPEN-008 le support ; OPEN-013 la gouvernance classe 2 ; OPEN-015 le bridge des Runs. Le protocole et la machine finale restent hors phase.

## 27. Consommateurs documentaires
Case Workspace, Interactive Operations, Session File Transfer, Result Handling, Provenance, Govern projections, Endpoint Agent Live Response et phases Objets/Permissions.
