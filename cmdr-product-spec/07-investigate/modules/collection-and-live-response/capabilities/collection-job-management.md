---
id: CAP-INV-203
title: Collection Job Management
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
  - REQ-PROD-019
open_decisions:
  - OPEN-008
source-of-truth: canonical
---
# CAP-INV-203 — Collection Job Management

## 1. Définition
Suivre l’exécution métier d’une collecte, ses résultats partiels, retries, annulations et Artifacts sans dupliquer Background Jobs.

## 2. Problème utilisateur
Sans ce suivi, un Background Job technique peut être confondu avec le résultat métier de collecte, et un état partiel ou offline peut être présenté comme un succès complet.

## 3. Objectifs
- relier Collection Request, cible, initiateur, policy, classe, progression et exécution Agent ;
- distinguer résultats réussis, échoués, expirés et en attente de reconnexion ;
- lier les Artifacts produits et préserver le retour au Case.

## 4. Non-objectifs
Ne pas posséder le moteur Background Jobs, définir transport/retry technique/API, créer automatiquement Evidence, assimiler Job à Automation Run/Response Run ou commencer Analysis Workbench.

## 5. Propriétaire
Investigate possède le contexte métier Collection Job. Shared possède le mécanisme Background Jobs ; Endpoint Agent exécute ; Govern possède Response Run ; Studio possède Automation Run.

## 6. Utilisateurs
Principal : Case Analyst ou Response Operator. Secondaires : DFIR Analyst, Investigation Lead, Evidence Reviewer et Incident Commander.

## 7. Conditions d’entrée
Collection Request acceptée, cible et Case résolus, capacité/policy snapshot, autorisation, projection Background Job et correlation ID.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Collection Request | Investigate | contrat métier | oui | version soumise | job non interprétable |
| Endpoint/Agent | Settings / Endpoint Agent | cible et disponibilité | oui | heartbeat visible | endpoint-offline/unknown |
| Background Job | Shared | queue/progression/cancel | oui après dispatch | événement courant | état `partial` ou `unknown` |
| Résultats et erreurs | Endpoint Agent | exécution locale | selon phase | horodatés | awaiting-result |
| Artifacts produits | Investigate | sorties matérielles | non | version/acquisition | aucun Artifact inventé |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Collection Request | Investigate | cible, scope, limites, classe | consulter |
| Collection Job concept | Investigate, modèle futur | statut métier, catégories, résultats | consulter/gérer selon permission |
| Background Job | Shared Capabilities | queue, progression, cancel et partial | consulter/invoquer mécanisme |
| Endpoint / Endpoint Agent | partagé / Endpoint Agent | disponibilité et exécution | consulter |
| Artifact / Case | Investigate | liens de sortie et contexte | consulter/lier |
| Automation Run / Response Run | Studio / Govern | distinction et provenance éventuelle | consulter uniquement |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Collection Job record conceptuel | créer, actualiser, cancel, expire ou supersede | Investigate, phase Objets future | distinct du Background Job |
| Artifact relation | créer pour chaque sortie réussie | Investigate | acquisition et résultat source obligatoires |
| Retry request | créer | Investigate | ciblé, idempotence visible, jamais succès global automatique |
| Background Job | aucune redéfinition | Shared | consommation seulement |

## 11. Fonctionnalités
Afficher queue, cible, état, progression, initiateur, policy, classe, phase, résultats partiels, erreurs, retry, cancel, timeout, expiration, reprise, succès/échecs par élément, Artifacts et retour au Case.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter progression | utilisateur autorisé | Job projection | 0 | read | état détaillé | non |
| Annuler avant effet final | analyste autorisé | Collection Job | 2 | cancellable | cancel demandé et tracé | OPEN-013 selon impact |
| Retry ciblé | analyste | éléments échoués | 1/2 | cause connue, idempotence | nouvelle tentative liée | selon classe |
| Ouvrir Artifact | analyste | Artifact | 0 | Artifact read | navigation CAP-INV-105 | non |
| Retourner au Case | utilisateur | Case | 0 | relation valide | contexte restauré | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Agréger progression | oui | oui | oui | non nécessaire | états Background Job |
| Détecter partial/timeout | oui | oui | oui | explication | règles d’état |
| Proposer retry ciblé | oui | règles | oui | suggestion | sélection manuelle |
| Résumer erreurs | oui | groupement | oui | résumé attribué | erreurs brutes et filtres |
| Annuler/reprendre | humain explicite | contrat | workflow possible | jamais autonome | action utilisateur |

Toute sortie automatisée expose initiateur, moteur/version, Automation Run/Tool Calls éventuels, sources, statut, incertitude, owner humain et trace.

## 14. États fonctionnels
`queued`, `dispatched`, `running`, `partial`, `completed`, `failed`, `cancelled`, `expired`, `endpoint-offline`, `awaiting-reconnect`. Machine finale reportée.

## 15. États d’interface
Loading conserve request/cible ; Empty signifie aucun dispatch ; Partial détaille chaque catégorie ; Error conserve résultats valides ; Offline distingue attente autorisée et échec ; Permission denied masque sorties sensibles ; Stale expose dernière mise à jour.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Collection Job status | record métier conceptuel | Case Workspace | Request, target, phase et erreurs visibles |
| Résultats par élément | result entries | CAP-INV-204..208/212 | succès et échec séparés |
| Artifact links | relations | CAP-INV-105/107 | aucune Evidence automatique |
| Progress/notification events | Shared events | utilisateur/Timeline/Trace | deep link et correlation ID |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-202 | request acceptée | Collection Job Management | request/version, cible, scope, limites, autorité | Request/Case |
| Collection Job | sortie réussie | CAP-INV-105 | source, acquisition, timestamps, transformations, erreurs | Job/Case |
| Job partial | retry ciblé | nouvelle tentative liée | éléments échoués, cause, idempotency context | même Job group |
| Job terminal | ouvrir résultat | CAP-INV-212 ou Case | job, outputs, Artifacts, erreurs, trace | Job/Case |

## 18. Dépendances
CAP-INV-202/204..208/105/107/212/214, Shared Background Jobs/Notifications/Trace, Endpoint Agent queue/result, Settings health/policy et OPEN-008.

## 19. Source de vérité
Investigate possède le statut métier et ses relations ; Shared possède la mécanique de Job ; Endpoint Agent possède l’exécution locale ; Govern/Studio restent owners de leurs Runs.

## 20. Provenance et audit
Request/version, Case, Endpoint, Agent, policy, initiateur, autorité, dispatch, phases, catégories, retries, cancel, timeout, résultats, erreurs, Artifacts et correlation IDs.

## 21. Permissions fonctionnelles
Collection read/cancel/retry, raw result read, Artifact receive, sensitive output, cross-tenant/environment et audit read. Permission atomique et step-up reportés.

## 22. Limites et erreurs
Endpoint offline, queue indisponible, expiry, résultat tardif, duplication potentielle, cancel trop tardif, résultat partiel, Artifact inaccessible ou permissions révoquées restent explicites ; aucun retry silencieux.

## 23. Métriques
Jobs par état, temps en queue, partial rate, retries, cancellations, expiry, awaiting-reconnect, Artifacts produits et retour Case réussi.

## 24. Classification de livraison
`defined` / `planned`. Aucun moteur Background Job, protocole, file technique ou plateforme n’est déclaré livré.

## 25. Critères d’acceptation
**Given** une collecte multi-catégories **When** certaines réussissent et d’autres échouent **Then** l’état est partial, chaque résultat est visible et un retry ciblé est possible.

**Given** un Endpoint offline **When** la demande est dispatchée **Then** l’état indique attente/reconnexion ou échec selon policy, sans faux démarrage.

**Given** aucun modèle IA **When** le Job est suivi **Then** progression, filtres, erreurs, cancel et retry restent disponibles.

## 26. Questions ouvertes
L’objet Collection Job, sa relation à Collection Request/Background Job, les états finaux et permissions appartiennent aux phases Objets/Technique ; OPEN-008 reste ouverte.

## 27. Consommateurs documentaires
Module Collection and Live Response, profiles de collecte, Case Workspace, Timeline/Replay, Artifact/Evidence, Endpoint Agent, Shared Background Jobs et futurs contrats de retry.