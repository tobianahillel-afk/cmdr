---
id: CAP-GOV-001
title: Govern Intake and Preconditions
product: govern
module: response-inbox
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-008, REQ-PROD-015, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-GOV-001 — Govern Intake and Preconditions

## 1. Définition
Recevoir une Action Request soumise à Govern, préserver sa version et son origine, vérifier les préconditions minimales de traitement et qualifier son état d’intake sans créer de Decision, Approval ou effet sur la cible.

## 2. Problème utilisateur
Une demande peut atteindre Govern avec cible ambiguë, contexte inaccessible, source stale, restrictions manquantes ou retour impossible. Sans intake explicite, une demande reçue peut être prise pour une demande complète ou autorisée.

## 3. Objectifs
- établir que la request reçue est identifiable et versionnée ;
- préserver source product, requester, return origin et provenance ;
- vérifier tenant/environnement, action proposée, target/scope et accès minimal ;
- distinguer `received`, `incomplete`, `blocked`, `duplicate-candidate` et `out-of-scope` ;
- router vers la Response Inbox sans décision implicite.

## 4. Non-objectifs
Ne pas qualifier Evidence/Finding, calculer la Decision finale, approuver, exécuter, créer un Response Run, résoudre un Policy conflict ou inventer une autorité.

## 5. Propriétaire
Govern / Response Inbox / Govern Product Lead. Le produit source reste propriétaire de ses objets et de son contenu métier ; Govern devient propriétaire du processing lifecycle de l’Action Request reçue.

## 6. Utilisateurs
Principal : Govern Reviewer. Secondaires : Govern Coordinator, Incident Commander/Investigation Lead comme requester ou source contributor, Auditor en lecture.

## 7. Conditions d’entrée
Action Request ou package de soumission explicite, requester attribué, tenant, return origin, action proposée et au moins une target/scope candidate. La soumission ne vaut ni complétude ni autorisation.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Action Request + version | Govern object / source submission | demande gouvernée | oui | version reçue immuable pour la review | intake refusé comme non identifiable |
| Source context | Command/Investigate/Detection/TI/analysis | Incident, Case, Finding, Evidence refs, package | non pour urgence justifiée, sinon selon demande | version/timestamp visibles | `incomplete` ou `information-required` |
| Requester + return origin | source product / identity projection | attribution et retour | oui | courant à la soumission | `blocked` |
| Tenant/environment | Settings/Security context | isolation | tenant oui; environment si pertinent | courant | permission denied / blocked |
| Proposed action + target/scope | requester/package | intention et cible | oui | version de request | `incomplete` |
| Restrictions/handling | source owners/Security | contraintes | selon contexte | version source | restricted/blocked plutôt qu’ignoré |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Action Request | Govern | id, version, requester, action, target/scope, status | lire et vérifier intake |
| Incident | Command | impact/urgence/return origin | projection lecture |
| Case/Finding/Evidence | Investigate | références et restrictions autorisées | projection lecture, aucune requalification |
| Principal/Tenant/Environment | Settings | identité/scope administratif | projection minimale |
| Workflow/Automation Run | Studio | provenance de préparation éventuelle | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Action Request processing state | enregistrer intake state et reason | Govern | versionné, aucune Decision |
| Govern Review Context | créer | Govern local concept | lie request/version/reviewer/return origin |
| duplicate candidate relation | proposer/lier | Govern/Shared Linking | candidat seulement, pas déduplication destructive |
| source objects | aucune mutation | source owners | liens et projections uniquement |

## 11. Fonctionnalités
Valider identifiants et version ; identifier source/return origin ; contrôler scope tenant/env ; vérifier présence minimale action/target ; signaler restrictions ; détecter candidat duplicate sans fusion automatique ; classer out-of-scope ; rendre les lacunes explicites ; préparer le queue item Govern.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter intake | Govern Reviewer | Action Request | 0 | read autorisé | contexte sourcé | non |
| exécuter precondition check | Govern Reviewer | Review Context | 1 | request/version accessible | résultat explicable | non |
| marquer information-required | Govern Reviewer | processing state | 2 | lacune nommée | source peut compléter | OPEN-013 |
| marquer out-of-scope | Govern Reviewer | processing state | 2 | raison documentée | retour source sans suppression | OPEN-013 |
| lier duplicate candidate | Govern Reviewer | relation | 2 | candidats comparables | relation réversible | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| vérifier champs minimaux | oui | oui, checklist | oui | expliquer lacunes | checklist structurée |
| détecter duplicate candidate | oui | matching explicable | oui | suggestion | recherche/compare manuels |
| résumer source context | oui | agrégation | oui | résumé sourcé | liens + tableau source |
| router vers queue | oui | règles explicites | oui | non nécessaire | affectation/routage manuel |
| décider/autoriser | oui par capacités ultérieures | validation seulement | jamais autonome | interdit | review humaine Govern |

## 14. États fonctionnels
`received`, `incomplete`, `ready-for-review`, `blocked`, `duplicate-candidate`, `policy-context-required`, `authority-context-required`, `information-required`, `out-of-scope`, `superseded`, `withdrawn`.

## 15. États d’interface
Loading conserve request/version ; Empty signifie aucune request sélectionnée ; Partial nomme les projections absentes ; Error conserve les données valides ; Offline interdit toute transition non garantie ; Permission denied ne révèle pas les références protégées ; Stale expose la version source obsolète.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| intake disposition | processing event | CAP-GOV-002/003 | reason, actor, request version et return origin |
| Review Context | Govern local context | Response Inbox / Action Center | aucune autorité ajoutée |
| information request need | transition context | source product | questions minimales et provenance |
| out-of-scope/withdrawn disposition | event | source product / history | historique conservé |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| source product | explicit submit | CAP-GOV-001 | request/version/source/return origin | source workspace |
| CAP-GOV-001 | preconditions sufficient | CAP-GOV-002 | Review Context + intake disposition | request intake |
| CAP-GOV-001 | missing source data | source product | question, missing fields, request/version | same Action Request lineage |
| CAP-GOV-001 | withdrawn/superseded | Decision/Request history | disposition and links | source product |

## 18. Dépendances
Action Request object, Settings tenant/environment/identity projections, Shared Linking/Trace, source-product links, CAP-GOV-002/003/004/006, OPEN-013 and OPEN-015 for future automation-run lineage details.

## 19. Source de vérité
Govern owns Action Request processing state and Review Context. Command/Investigate/Studio/Settings remain authoritative for their source objects and projections. Intake never copies them into a new owner.

## 20. Provenance et audit
Record submitter/requester, source product/object refs, request id/version, tenant/env, action/target summary, restrictions, precondition results, duplicate candidates, missing fields, human/automation producer, timestamps, reasons, return origin and correlation ids.

## 21. Permissions fonctionnelles
Govern Inbox/read, Action Request read/review, restricted context read, information request, Govern assignment preparation and cross-tenant review where explicitly authorized. Atomic permissions and step-up remain future.

## 22. Limites et erreurs
Missing/stale/inaccessible source, tenant mismatch, target ambiguity, withdrawn version, duplicate candidate, unsupported action, destination unavailable or permission denial must not become a fabricated complete request. Safe partial context is retained.

## 23. Métriques
Requests by intake disposition; missing-field categories; blocked/out-of-scope/duplicate candidates; source-return success; stale/restricted projections; zero silent Decision/Approval/execution from intake.

## 24. Classification de livraison
`defined` / `planned`. Documentary functional contract only; no ingestion API, queue engine or native runtime is claimed.

## 25. Critères d’acceptation
**Given** an Action Request without an exact target, **When** intake runs, **Then** it becomes `incomplete`/`information-required`, the missing target is named and no Decision or Approval is created.

**Given** a request produced by an Automation Run, **When** intake opens it, **Then** the Run reference is preserved as provenance and does not become authority.

**Given** no AI provider, **When** intake is performed, **Then** deterministic field checks, source links and manual review provide the complete path.

## 26. Questions ouvertes
OPEN-013 remains open for default class-2 governance. OPEN-015 remains open for final Automation Run/Response Run provenance bridge. No new OPEN is needed for intake.

## 27. Consommateurs documentaires
Response Inbox, Action Center, source-product handoffs, Govern capability/register maps, future Objects/Permissions/Screens/Journeys/Technique, quality report and GOV-2 handoff design.
