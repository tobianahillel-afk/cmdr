---
id: investigate-cross-product-links
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-004
  - REQ-PROD-005
  - REQ-PROD-008
open_decisions:
  - OPEN-007
  - OPEN-013
  - OPEN-015
---

# Cross-product links — Investigate Phase 4B.1

| Transition | Source / déclencheur | Destination | Contexte transmis | Ownership | Retour et erreurs |
|---|---|---|---|---|---|
| Signal → investigation | analyste depuis Signal Triage ; qualification ou pivot | Event Search, Hunt ou Case Lifecycle | tenant, environnement, Signal, Alert éventuel, Detection, Events, Incident, période, raison et return origin | Signal/Alert/Incident restent Command ; Case/Hypothesis futurs restent Investigate | retour exact au Signal ; déduplication conceptuelle avant nouveau Case ; erreur conserve le triage |
| Incident → Case | utilisateur Command choisit Open in Investigate | Case existant ou nouveau | tenant, environnement, service, Incident, Signal/Alert, objectif initial, période, requester et owner proposé | Incident reste Command ; Case appartient à Investigate | retour vers Incident Detail ; Case inaccessible produit une proposition ou une erreur sans divulgation |
| Search → Case | sélection de résultats ou Add to Case | Case Workspace / Case Lifecycle | Query et version, paramètres, période, Search Job, résultats sélectionnés, sources, annotations, provenance et Evidence candidates | Event/Query/Search Job restent Shared ; Case reste Investigate | retour à la recherche avec sélection et filtres ; aucun résultat ne devient Evidence automatiquement |
| Hunt → Case | promotion humaine d’un Hunt | Case existant ou nouveau | question, scope, période, Queries, runs, résultats, Hypotheses, auteurs, statut, conclusion et limites | Hunt workspace et Case restent Investigate ; objets Shared conservés | retour au Hunt ; package sélectif et versionné ; Case inaccessible conserve le Hunt |
| Artifact → Evidence | analyste qualifie un Artifact ou une autre source | Evidence Creation | Artifact/source, origine, acquisition/import, Case, auteur, raison, transformations, version, provenance et rôle | Artifact et Evidence restent deux objets Investigate distincts | retour à Artifact Detail ; aucune conversion automatique ; échec conserve l’Artifact et un candidat non qualifié |
| Evidence → Finding | analyste/reviewer crée ou révise une assertion | Finding Management | Evidence favorables et contradictoires, versions, auteur, revue, incertitude, Hypotheses et Case | Evidence et Finding restent Investigate | retour à Evidence Review ; Finding reste draft/proposed jusqu’à confirmation humaine autorisée |
| Finding → Action Request | Investigation Lead prépare une action | Action Request Preparation | Finding, Evidence, cible, action, impacts technique/opérationnel, urgence, alternatives, conditions, rollback, incertitudes et requester | Finding/Evidence restent Investigate ; Action Request lifecycle appartient à Govern | retour au Finding/Case ; draft conservé si incomplet ; aucune Decision créée |
| Action Request → Govern | soumission humaine ou workflow autorisé | Govern Review Queue | Action Request version, Case, Findings/Evidence refs, cible, impacts, alternatives, rollback, provenance, permission et return origin | Govern possède request lifecycle, Decision, Approval, Run et Result | retour pour informations vers le même draft ; future Decision liée au Case ; refus ne supprime pas l’historique |
| Result → Case | Result Govern vérifié ou projeté | Case Workspace / Finding Review | Action Request, Decision, Response Run, Result, statut de vérification, impact, risque résiduel, Findings concernés et prochaine action | Result reste Govern ; relations et réactions analytiques restent Investigate | retour vers Result ou Incident ; Hypotheses peuvent être réévaluées sans modification automatique |
| Investigate → Studio | utilisateur demande une assistance déployée | Workflow ou Automation Run inspectable | Case ou Query context, version de workflow, permissions, limites, sources et initiateur | Workflow/Agent/Automation Run restent Studio | retour au workspace source avec proposition attribuée ; Tool Calls et run inspectables |

## Invariants de transition

- tenant, environnement, sélection, filtres et return origin sont conservés ;
- la permission est réévaluée à l’entrée et au retour ;
- une transition ne transfère jamais l’ownership des objets source ;
- une erreur conserve le workspace et les drafts récupérables avec correlation ID ;
- les sorties automatisées restent des propositions tant qu’une action humaine ou un contrat déterministe autorisé ne les rend pas effectives ;
- les classes 3 et 4 passent par Govern et ne sont jamais exécutées directement dans Investigate.
