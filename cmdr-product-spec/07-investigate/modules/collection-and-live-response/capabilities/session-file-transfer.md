---
id: CAP-INV-211
title: Session File Transfer
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
  - OPEN-014
source-of-truth: canonical
---
# CAP-INV-211 — Session File Transfer

## 1. Définition
Envoyer un fichier autorisé vers l’Endpoint ou récupérer un fichier depuis l’Endpoint dans une Live Session, avec classification, vérification, progression et lien Case.

## 2. Problème utilisateur
Le transfert de session peut être confondu avec acquisition, Attachment ou déploiement logiciel. L’opérateur doit savoir si le fichier est temporaire, collecté comme Artifact ou simplement envoyé comme outil autorisé.

## 3. Objectifs
- définir direction, source, destination fonctionnelle, taille, type, provenance et finalité
- vérifier session, Endpoint, policy, permission, classe et collision
- suivre progression, annulation, vérification et échecs partiels
- enregistrer correctement Artifact, temporary transfer ou Operation Result sans déploiement implicite

## 4. Non-objectifs
Ne pas définir protocole, chemin système exact, mécanisme d’upload/download, déploiement logiciel, package manager ou outil ; ne pas résoudre OPEN-014 ; ne pas qualifier automatiquement Evidence.

## 5. Propriétaire
Investigate / Collection and Live Response / Investigate Product Lead possède le contexte métier, les drafts et les relations au Case. Platform Settings reste propriétaire de Fleet et Endpoint Policies ; Endpoint Agent exécute et rapporte localement ; Govern conserve l’autorité, Decision, Response Run et Result.

## 6. Utilisateurs
Principal : Response Operator ou DFIR Analyst. Secondaires : Case Analyst, Evidence Reviewer et Investigation Lead.

## 7. Conditions d’entrée
Live Session active, Case/Endpoint accessibles, fichier/source autorisé, destination fonctionnelle bornée, policy/permission/classe validées, taille/type/provenance disponibles.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Live Session, Case et Endpoint | Investigate | contexte et cible | oui | états courants | transfert interdit |
| Direction et finalité | opérateur | upload outil temporaire ou récupération Artifact | oui | confirmées avant transfert | incomplete |
| Fichier/source, type, taille et provenance | Artifact/Attachment/authorised source | contenu référencé | oui | version courante | validation impossible |
| Destination fonctionnelle et collision policy | opérateur / policy | scope cible | oui | revalidée au lancement | restricted/collision |
| Permission, policy et classe | Security / Settings / Govern | gate | oui | snapshot à l’action | denied/awaiting-approval |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Live Session | Investigate, modèle futur | participants, statut et Endpoint | consulter/opérer |
| Artifact | Investigate | source à envoyer ou fichier récupéré | consulter/lier |
| Attachment | concept ouvert | fichier documentaire éventuel | consulter sans assimilation |
| Endpoint / Endpoint Agent / Policy | partagé / Agent / Settings | capability et restrictions | consulter |
| Operation Result | concept futur | progression, vérification et erreurs | consulter/lier |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Transfer record conceptuel | créer, actualiser, annuler ou supersede | Investigate, modèle futur | direction, finalité, source, destination et class obligatoires |
| Artifact | créer à récupération explicite | Investigate | source Endpoint, acquisition et Case conservés |
| Temporary transfer relation | créer puis clôturer | Investigate | ne devient ni Artifact ni deployment automatiquement |
| Attachment/Deployment/Fleet | aucune mutation implicite | owners respectifs | promotion ou administration séparée |

## 11. Fonctionnalités
- distinguer upload d’outil autorisé, récupération d’Artifact, Attachment et transfert temporaire
- afficher source, destination fonctionnelle, taille, type, provenance, policy et classe
- confirmer collisions et finalité avant démarrage
- suivre preparing/transferring/verifying/partial/failed/cancelled
- lier le résultat au Case et conserver l’Artifact uniquement lorsque requis

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Préparer récupération | DFIR Analyst | transfer draft | 1 | session active, source autorisée | ready | non normalement |
| Préparer upload temporaire | Response Operator | transfer draft | 2 | outil autorisé, destination/cleanup | ready | OPEN-013 selon policy |
| Confirmer collision | Response Operator | transfer | 2 | collision détectée et options permises | choix audité | OPEN-013 |
| Annuler | Response Operator | transfer | 2 | cancellable | cancelled/partial | selon policy |
| Enregistrer comme Artifact | Case Analyst | Artifact | 1/2 | récupération, provenance complète | Artifact lié | non automatique |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| Détecter collision/type/taille | oui | validateurs | oui | explication | contrôles déterministes |
| Proposer finalité | oui | règles | oui | suggestion | sélection manuelle |
| Suivre progression | oui | états déterministes | oui | résumé | progression brute |
| Proposer cleanup | oui | policy/catalogue | oui | suggestion | checklist |
| Transformer en Artifact/Evidence | humain explicite | règles de promotion | workflow possible | jamais automatique | CAP-INV-105/107 |

Toute sortie automatisée expose initiateur, producteur/version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres fonctionnels, timestamp, statut, incertitude, owner humain, acceptation/modification/rejet et trace.

## 14. États fonctionnels
`preparing`, `transferring`, `verifying`, `completed`, `partial`, `failed`, `cancelled`, `collision`, `restricted`. Machine finale reportée.

## 15. États d’interface
Loading conserve source/destination ; Empty distingue aucun fichier et permission ; Partial affiche octets/éléments reçus sans succès global ; Error garde l’Artifact reçu ; Offline n’indique pas completed ; Permission denied masque contenus ; Stale exige revalidation.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Transfer status/result | record conceptuel | Live Session/CAP-INV-212 | direction, finalité, progression, verification et errors |
| Recovered Artifact | Artifact | CAP-INV-105/107 | source Endpoint et acquisition conservées |
| Temporary upload record | relation/event | Session/Trace | cleanup/finalité visibles, pas de deployment |
| Collision/cancel event | business event | Timeline/Notifications | choix, acteur et résultat |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Live Session | préparer transfert | CAP-INV-211 | session, Endpoint, direction, source, destination, policy | Live Session |
| CAP-INV-211 | démarrer | Endpoint Agent | transfer context, class, authority, limits | CAP-INV-211/Session |
| Récupération terminée | enregistrer Artifact | CAP-INV-105 | file, source, acquisition, verification, errors | CAP-INV-211 |
| Transfer result | revoir | CAP-INV-212 | status, output, errors, Artifact/temp relation | Live Session/Case |

## 18. Dépendances
CAP-INV-105/107/209/212/213/214, Endpoint Agent transfer, Settings Policy, Shared Jobs/Preview/Trace, OPEN-008/013/014.

## 19. Source de vérité
Investigate possède le record métier et l’Artifact récupéré. Endpoint Agent reste source de l’exécution locale. Attachment reste ouvert sous OPEN-014 ; deployment et Fleet restent hors Investigate.

## 20. Provenance et audit
Case, session, Endpoint/Agent, direction, finalité, source/version, destination fonctionnelle, type/taille, policy, initiateur, progress, collision, verification, Artifact/temp relation et correlation IDs.

## 21. Permissions fonctionnelles
File download/upload, Artifact read/create/export, session operation, sensitive content, collision override, cross-tenant/environment et step-up pour upload.

## 22. Limites et erreurs
Fichier absent/verrouillé, source stale, type/taille interdits, collision, destination restricted, partial transfer, disconnect, cleanup non confirmé ou permission révoquée. Aucun software deployment implicite.

## 23. Métriques
Transfers par direction/finalité, collision/partial/failure, cancellations, verified transfers, Artifacts créés et temporary uploads clôturés.

## 24. Classification de livraison
`defined` / `planned`. Aucun protocole, path exact, outil, package ou deployment mechanism n’est déclaré.

## 25. Critères d’acceptation
**Given** un fichier autorisé à récupérer **When** le transfert se termine **Then** progression, vérification, erreurs et Artifact source sont visibles sans qualification Evidence automatique.

**Given** un upload temporaire **When** l’opérateur le prépare **Then** finalité, destination, classe et cleanup attendu sont visibles et il n’est pas présenté comme deployment.

**Given** aucun modèle IA **When** un transfert est préparé **Then** formulaire, validateurs, confirmation de collision et suivi déterministe suffisent.

## 26. Questions ouvertes
OPEN-008, OPEN-013 et OPEN-014 restent ouvertes. Les chemins, protocoles, cleanup technique et identité Attachment/Artifact restent hors phase.

## 27. Consommateurs documentaires
Live Session, Result Handling, Artifact/Evidence, Custody, Provenance, Endpoint Agent transfer, Shared Preview/Export et phases Objets/Permissions.
