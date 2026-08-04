---
id: CAP-INV-327
title: Dynamic Analysis Provenance and Reproducibility
product: investigate
module: dynamic-sandbox
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-020
  - REQ-INV-005
  - REQ-AI-002
  - REQ-OBJ-009
open_decisions:
  - OPEN-005
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-327 — Dynamic Analysis Provenance and Reproducibility

## 1. Définition
Retracer Case, Dynamic Analysis Session, Artifact et dérivés, Runtime Artifacts, environnement/version, profils, Tools, Tool Calls, Automation Runs, paramètres, observations, erreurs, interruptions, nettoyage et dispositions, puis évaluer la reproductibilité.

## 2. Problème utilisateur
Un comportement peut dépendre d’un environnement, d’un profil, d’un timing ou d’une interaction. Sans ces préconditions, un Run reproduit partiellement ne doit pas être présenté comme équivalent.

## 3. Objectifs
- Reconstruire la chaîne complète d’une analyse dynamique.
- Distinguer reproduced, partial, not-reproduced et behavior-not-reproduced.
- Préparer un nouveau Run avec préconditions visibles.
- Conserver les écarts et contestations sans réécrire l’historique.

## 4. Non-objectifs
Ne pas définir infrastructure de replay, stockage, API, protocole, commande, moteur, hyperviseur, reverse/debugger ou schéma final.

## 5. Propriétaire
Investigate possède la sémantique analytique; Studio possède Automation Run/Tool Calls; Settings possède Environment/version; Shared possède Trace/Activity/Audit mechanisms.

## 6. Utilisateurs
Malware Analyst; Reviewer; Auditor; Dynamic Analysis Operator.

## 7. Conditions d’entrée
Session ou Run identifiable; inputs, versions, environnement, profils, Tools et dispositions accessibles ou lacunes explicites.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Case/Session/Artifact | Investigate | contexte et inputs | oui | versions référencées | provenance broken |
| Runs et observations | CAP-INV-317..325 | exécutions/résultats | oui | états et timestamps visibles | assessment partial |
| Environment/profile versions | Settings/Studio | préconditions | oui | snapshot Run | missing-environment/version |
| Tool Calls/Automation Runs | Studio | producteur et paramètres | si utilisés | références résolubles | producer incomplete |
| Dispositions humaines | analyste/reviewer | accept/modify/reject/dispute | oui selon proposition | horodatées | unresolved |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Dynamic Analysis Session / Sandbox Run | Investigate concepts | scope, phases et outputs | consulter |
| Artifact / Runtime Artifact | Investigate | inputs, versions et lineage | consulter |
| Sandbox Environment | Platform Settings | version, health et restrictions | consulter |
| Tool Call / Automation Run | Studio | versions, paramètres et résultats | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Provenance chain relation/event | émettre/supersede | Investigate semantics/Shared mechanisms | source et version obligatoires |
| Reproducibility assessment | créer/versionner/contester | Investigate concept | préconditions et écarts visibles |
| New Run request context | préparer | Investigate | aucun lancement automatique |

## 11. Fonctionnalités
Retracer inputs, versions, environnements, profils, Tools, Runs, observations, erreurs, cleanup/reset et dispositions; comparer reproduction et source; signaler toute précondition manquante.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Ouvrir provenance | Reviewer | Provenance chain | 0 | read | chaîne affichée | non |
| Demander reproduction | Analyst | New Run request | 1 | préconditions/permission | intake préparé | non |
| Comparer reproduction | Reviewer | Sandbox Runs | 0 | Runs accessibles | CAP-INV-325 | non |
| Contester assessment | Reviewer | Assessment | 2 | justification | disputed | OPEN-013 |
| Exporter trace | Auditor | Export request | 1 | redaction/policy | package contrôlé | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Corréler sources/versions | oui | IDs et relations | oui | explication | graph/timeline |
| Détecter précondition manquante | oui | règles | oui | résumé | checklist |
| Résumer les écarts | oui | comparaison déterministe | oui | résumé attribué | diff brut |
| Déclarer reproduced | reviewer | contrôles explicites | workflow de revue | jamais autonome | revue humaine |

## 14. États fonctionnels
`reproduced`, `partially-reproduced`, `not-reproduced`, `missing-environment`, `missing-version`, `missing-input`, `policy-blocked`, `behavior-not-reproduced`, `disputed`.

## 15. États d’interface
Loading conserve les références; Partial nomme chaque lien manquant; Error garde les événements valides; Offline est read-only; Permission denied redacted; stale montre les versions retirées.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Dynamic provenance chain | relations/events | Case Replay/Audit | sources, versions et dispositions |
| Reproducibility assessment | Analysis Result | Reviewer | préconditions et écarts |
| New Run context | request | CAP-INV-314/317 | scope, environment, profile, versions |
| Comparison result | Multi-Run comparison | Case | non reproduction distincte d’absence de menace |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Dynamic Analysis Session | ouvrir trace | CAP-INV-327 | session, Runs, inputs, versions, outputs | session |
| CAP-INV-327 | reproduire | CAP-INV-314/317 | Artifact, environment/version, profile, Tools, parameters, limits | provenance |
| New Sandbox Run | comparer | CAP-INV-325/327 | results, observations, errors, differences | provenance |
| CAP-INV-327 | préparer handoff | CAP-INV-328 | résultats, assessment, sources, contradictions | provenance |

## 18. Dépendances
CAP-INV-312/314/317/325/328; Settings versions/health; Studio Tools/Runs; Shared Trace/Activity/Export; OPEN-005/015.

## 19. Source de vérité
Investigate émet la sémantique; Settings et Studio restent sources de leurs versions; reproduced exige des préconditions démontrées.

## 20. Provenance et audit
Case, session, source/runtime Artifacts, Runs, environment/profile/Tool versions, paramètres, observations, erreurs, interruptions, cleanup/reset, human disposition et correlation IDs.

## 21. Permissions fonctionnelles
Reproducibility review; new Run request; Tool/Automation Run projection read; dynamic provenance export; sensitive result read; cross-tenant denied.

## 22. Limites et erreurs
Environment/version/input manquant; policy-blocked; behavior not reproduced; timing/interactions différents; Run partial; permission révoquée; source redacted.

## 23. Métriques
Assessments reproduced/partial/not; préconditions manquantes; behavior-not-reproduced; disputes; replays; provenance complète.

## 24. Classification de livraison
`defined` / `planned`; aucune infrastructure de reproduction, moteur, API, protocole, commande ou implémentation.

## 25. Critères d’acceptation
**Given** un ancien Run et une version d’environnement indisponible
**When** l’analyste tente de le reproduire
**Then** missing-version est visible et le résultat n’est pas présenté comme reproduced

**Given** un nouveau Run qui ne reproduit pas le comportement
**When** l’assessment est créé
**Then** behavior-not-reproduced est distinct de not-malicious

**Given** aucun modèle IA
**When** la provenance est revue
**Then** relations, timelines, comparateurs et checklists couvrent le workflow

## 26. Questions ouvertes
Les objets finaux restent à la phase Objets; OPEN-005 et OPEN-015 restent ouvertes.

## 27. Consommateurs documentaires
Dynamic Sandbox, Case Replay, Multi-Run Comparison, Evidence/Finding handoff, Audit, Objets, Trust et Permissions.
