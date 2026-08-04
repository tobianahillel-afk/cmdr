---
id: CAP-INV-316
title: Dynamic Analysis Session Management
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
  - REQ-AI-002
open_decisions:
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-316 — Dynamic Analysis Session Management

## 1. Définition
Gérer le contexte durable d’une analyse dynamique, ses objectifs, participants, Artifact, environnements, profils, Tools, paramètres, Sandbox Runs, erreurs, interruptions et dispositions, distinct de chaque Run et Automation Run.

## 2. Problème utilisateur
Une investigation dynamique peut nécessiter plusieurs Runs dans des environnements ou profils différents. Sans session, les comparaisons, décisions et reprises perdent leur contexte.

## 3. Objectifs
- Créer, reprendre, suspendre, clôturer et comparer une Dynamic Analysis Session sans la confondre avec un Sandbox Run.
- Conserver Case, Artifact, session, Runs, versions, limites et return origin.
- Exposer erreurs, résultats partiels, provenance et disposition humaine.
- Transmettre les résultats sans qualification automatique.

## 4. Non-objectifs
Ne pas définir virtualisation, instrumentation, API, protocole, commande, code, moteur ou hyperviseur; ne pas commencer Reverse Engineering, Debugger ou forensics; ne pas agir sur un Endpoint réel.

## 5. Propriétaire
Investigate possède le contexte analytique. Settings possède les Sandbox Environments; Studio possède Tool, Tool Call et Automation Run; Govern possède l’autorité sur les cibles réelles; Shared possède les mécanismes transversaux.

## 6. Utilisateurs
Malware Analyst; Dynamic Analysis Operator; Case Analyst; Evidence Reviewer.

## 7. Conditions d’entrée
Case et Artifact accessibles; objectif et owner définis; environnement autorisé ou besoin explicite; permissions, policies, versions et restrictions réévaluées.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Case et objectif | Investigate | contexte durable | oui | état courant | rester draft |
| Artifact source | CAP-INV-105 | entrée analysée | oui | version référencée | session incomplete |
| Owner et contributeurs | utilisateur/identité | responsabilité | oui | réévalués | blocked |
| Environnements et profils | Settings/Studio projections | préconditions des Runs | selon Run | versions visibles | Run non préparé |
| Runs, Tools et résultats | Investigate/Studio | historique de session | non au départ | états courants | session sans Run |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case / Artifact | Investigate | scope, source et return origin | consulter/lier |
| Sandbox Environment | Platform Settings | version, health et limites | consulter |
| Tool / Tool Call / Automation Run | Studio | version, statut, output et provenance | consulter |
| Sandbox Run | concept Investigate | états, profils et résultats | consulter/lier |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Dynamic Analysis Session | créer/modifier/clôturer/reouvrir/supersede | Investigate concept | schéma final reporté |
| Session relation | lier Case, Artifact, Runs, Tools et résultats | Investigate | relation sourcée |
| Session disposition | enregistrer | Investigate | completed ne confirme aucun Finding |

## 11. Fonctionnalités
Créer/reprendre; définir objectif, owner et contributeurs; lier plusieurs Runs; enregistrer environnements, Tools, profils, paramètres, erreurs et interruptions; suspendre, clôturer, rouvrir, comparer et transmettre.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Créer session | Analyst | Dynamic Analysis Session | 2 | Case/Artifact/objective | draft créé | OPEN-013 |
| Ajouter contributeur | Owner | Session | 2 | permission | participants versionnés | OPEN-013 |
| Suspendre/reprendre | Contributor | Session | 2 | état compatible | contexte conservé | OPEN-013 |
| Clôturer/réouvrir | Owner/Lead | Session | 2 | justification | historique conservé | OPEN-013 |
| Comparer Runs | Reviewer | Sandbox Runs | 0 | Runs accessibles | CAP-INV-325 | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Résumer la session | oui | agrégation | oui | résumé attribué | timeline et filtres |
| Proposer un prochain Run | oui | profils/checklists | oui | suggestion modifiable | sélection manuelle |
| Capturer Runs/Tool Calls | non | relations déterministes | oui | non nécessaire | Trace |
| Clôturer ou qualifier | humain | contrôles seulement | workflow de revue | jamais autonome | checklist humaine |

Toute automatisation expose initiateur, producteur/version, Runs, Tool Calls, sources, paramètres, timestamp, statut, incertitude, owner humain et disposition.

## 14. États fonctionnels
`draft`, `ready`, `active`, `paused`, `blocked`, `partial`, `completed`, `failed`, `archived`, `superseded`. Machines finales reportées.

## 15. États d’interface
Loading conserve le contexte; Empty explique; Partial nomme les lacunes; Error garde les résultats valides; Offline est stale/read-only; Permission denied ne fuit rien.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Dynamic Analysis Session | contexte | Case/Workbench | scope, Runs et dispositions |
| Session summary | Analysis Result | Reviewer | résultats et lacunes visibles |
| Run comparison entry | comparison context | CAP-INV-325 | préconditions et versions |
| Handoff context | draft package | CAP-INV-328 | sources et contradictions |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-314 | créer/reprendre | CAP-INV-316 | Case, Artifact, objectif, limites, owner | intake |
| CAP-INV-316 | préparer Run | CAP-INV-317 | session, Artifact, environnement, profil, Tools | session |
| CAP-INV-316 | transmettre résultats | CAP-INV-328 | Runs, observations, Runtime Artifacts, contradictions | session |

## 18. Dépendances
CAP-INV-314/315/317/325/328; CAP-INV-105; Settings Environments; Studio Tools/Runs; Shared Trace/Collaboration/Versioning; OPEN-013/015.

## 19. Source de vérité
Investigate est source de la session et de sa disposition; chaque Sandbox Run reste distinct; Studio reste source des Tool Calls et Automation Runs.

## 20. Provenance et audit
Création, versions, owner, contributeurs, objectif, Artifacts, environnements, profils, Runs, Tools, paramètres, erreurs, interruptions, dispositions et return origin.

## 21. Permissions fonctionnelles
Dynamic Analysis read; session create/update/close/reopen; contributor manage; Sandbox Run link; Tool output read; multi-Run comparison; cross-tenant denied.

## 22. Limites et erreurs
Artifact inaccessible; contributeur retiré; environnement ou Tool retiré; Run orphelin; résultat expiré; réouverture avec versions manquantes; permission révoquée.

## 23. Métriques
Sessions actives/paused/partial; durée et reprises; Runs par session; sessions avec provenance complète; handoffs réalisés.

## 24. Classification de livraison
`defined` / `planned`; aucune implémentation, moteur, hyperviseur, API, protocole ou commande.

## 25. Critères d’acceptation
**Given** une session comportant plusieurs Runs
**When** l’analyste la suspend puis la reprend
**Then** objectif, Artifact, environnements, profils, Runs, paramètres et résultats restent liés

**Given** un utilisateur sans droit de réouverture
**When** il tente de rouvrir une session
**Then** l’action est refusée et l’historique reste inchangé

**Given** aucun modèle IA
**When** la session est gérée de bout en bout
**Then** formulaires, états, timelines, comparateurs et revue humaine restent disponibles

## 26. Questions ouvertes
Dynamic Analysis Session reste un concept; OPEN-013 et OPEN-015 restent ouvertes.

## 27. Consommateurs documentaires
Dynamic Sandbox, Case Workspace, Analysis Workbench, Multi-Run Comparison, Evidence/Finding handoff, Objets, Permissions et Journeys.
