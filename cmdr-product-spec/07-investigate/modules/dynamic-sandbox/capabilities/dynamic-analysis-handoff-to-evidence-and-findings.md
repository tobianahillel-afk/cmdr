---
id: CAP-INV-328
title: Dynamic Analysis Handoff to Evidence and Findings
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
  - REQ-OBJ-004
  - REQ-AI-002
open_decisions:
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-328 — Dynamic Analysis Handoff to Evidence and Findings

## 1. Définition
Sélectionner observations, Runtime Artifacts, résultats réseau, changements, contradictions et sources d’une analyse dynamique pour préparer une Evidence candidate ou un Finding Draft, sans qualification ni confirmation automatique.

## 2. Problème utilisateur
Un comportement observé, un screenshot, une destination réseau ou un score de sandbox n’est pas une Evidence ou un Finding. Le handoff doit conserver Run, environnement, profil, sources et incertitude.

## 3. Objectifs
- Sélectionner résultats dynamiques et éléments contradictoires.
- Expliquer la relation avec l’Hypothesis.
- Préparer Evidence candidate et Finding Draft avec provenance/reproductibilité.
- Transmettre aux lifecycles 107/108/109 et revenir au Workbench.

## 4. Non-objectifs
Ne pas qualifier Evidence, confirmer Finding, convertir un score en conclusion, créer Report canonique, choisir moteur ou modifier les résultats sources.

## 5. Propriétaire
Investigate possède le handoff analytique; Evidence/Finding restent dans leurs capabilities propriétaires; Studio/Settings/Shared conservent leurs objets.

## 6. Utilisateurs
Malware Analyst; Case Analyst; Evidence Reviewer; Investigation Lead.

## 7. Conditions d’entrée
Case et session actifs; résultats, observations ou Runtime Artifacts sélectionnables; sources, Run, environnement et permissions accessibles.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Dynamic Analysis Session/Runs | CAP-INV-316/317 | contexte d’analyse | oui | versions visibles | handoff incomplete |
| Observations sélectionnées | CAP-INV-318..322 | comportements et contradictions | oui | sources résolubles | draft blocked |
| Runtime Artifacts | CAP-INV-324 | sorties matérielles | non | versions visibles | aucune pièce matérielle |
| Comparison/reproducibility | CAP-INV-325/327 | écarts et confiance | non | assessment courant | incertitude explicite |
| Hypothesis/Case | CAP-INV-103/102 | raisonnement/contexte | oui pour Finding Draft | version courante | relation manquante |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case / Hypothesis | Investigate | contexte et relation analytique | consulter/lier |
| Artifact / Runtime Artifact | Investigate | source, lineage et restrictions | consulter |
| Sandbox Run / Observations | Investigate concepts | résultats, erreurs et sources | consulter |
| Evidence / Finding | Investigate | lifecycle et état de revue | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Evidence candidate | préparer/modifier/retirer | Investigate | qualification CAP-INV-107/108 obligatoire |
| Finding Draft | préparer/modifier | Investigate | confirmation CAP-INV-109 obligatoire |
| Handoff disposition | enregistrer | Investigate | accept/modify/reject |
| Source objects | aucune mutation silencieuse | owners respectifs | références seulement |

## 11. Fonctionnalités
Sélectionner observations, Runtime Artifacts, réseau, changements, contradictions et sources; relier Hypothesis; préparer candidate/draft; transmettre à revue; conserver provenance, Run et return origin.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Sélectionner résultats | Analyst | Dynamic results | 0 | sources visibles | sélection conservée | non |
| Préparer Evidence candidate | Analyst | Evidence candidate | 2 | provenance | draft sourcé | OPEN-013 |
| Préparer Finding Draft | Analyst | Finding Draft | 2 | Evidence/contradictions | draft attribué | OPEN-013 |
| Transmettre à revue | Analyst/Lead | Handoff | 2 | permission/complétude | CAP-INV-107/108/109 | OPEN-013 |
| Revenir au Workbench | Reviewer | Return origin | 0 | contexte conservé | session/Run restaurés | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Assembler les sources | oui | relations/IDs | oui | résumé | sélection manuelle |
| Proposer candidate/draft | oui | templates/checks | oui | suggestion modifiable | formulaires |
| Signaler contradictions | oui | comparateurs | oui | explication | listes sourcées |
| Qualifier/confirmer | humain autorisé | contrôles seulement | workflow de revue | jamais autonome | CAP-INV-107/108/109 |

## 14. États fonctionnels
`draft`, `incomplete`, `ready-for-review`, `submitted`, `returned`, `accepted-as-candidate`, `rejected`, `superseded`.

## 15. États d’interface
Loading conserve session/Run; Partial nomme sources manquantes; Error garde le draft; Offline interdit soumission non garantie; Permission denied ne fuit rien; stale montre versions superseded.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Dynamic Evidence candidate package | draft context | CAP-INV-107/108 | Run, environment, sources et qualification nécessaire |
| Finding Draft package | draft context | CAP-INV-109 | observations, Evidence et contradictions |
| Handoff disposition | business event | session/Trace | accepté, modifié ou rejeté |
| Return context | navigation context | Dynamic Workbench | session, Run et sélection restaurés |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Dynamic Analysis Session | préparer Evidence | CAP-INV-328 | Runs, observations, Runtime Artifacts, Hypothesis, provenance | session |
| CAP-INV-328 | soumettre candidate | CAP-INV-107/108 | candidate, Artifact/Runtime Artifact, Case, qualification requise | Dynamic Workbench |
| CAP-INV-328 | soumettre Finding Draft | CAP-INV-109 | draft, Evidence, contradictions, incertitude, reproductibilité | Dynamic Workbench |

## 18. Dépendances
CAP-INV-103/105/107/108/109/316..327; Shared Linking/Trace; Studio provenance; OPEN-013/015.

## 19. Source de vérité
Les résultats restent distincts des Evidence/Findings; CAP-INV-107/108/109 conservent qualification et confirmation; Tool output ≠ Report canonique.

## 20. Provenance et audit
Case, Hypothesis, session, Runs, environment/profile/Tool versions, observations, Runtime Artifacts, comparisons, reproducibility, sources, contradictions, auteur et disposition.

## 21. Permissions fonctionnelles
Dynamic results read; Runtime Artifact read; Evidence candidate prepare; Finding Draft prepare; result annotate; handoff submit; sensitive sources read; cross-tenant denied.

## 22. Limites et erreurs
Source inaccessible; provenance incomplète; Run partial/timed-out; environment/version manquant; contradiction non incluse; permission révoquée; draft returned.

## 23. Métriques
Candidates/drafts préparés; retours/rejets; sources manquantes; contradictions incluses; temps jusqu’à qualification; confirmations automatiques cible zéro.

## 24. Classification de livraison
`defined` / `planned`; aucune qualification automatique, moteur, API, protocole, commande ou implémentation.

## 25. Critères d’acceptation
**Given** des résultats dynamiques reliés à un Case
**When** l’analyste prépare une Evidence candidate
**Then** observation, Runtime Artifact, Run, environnement, sources et qualification nécessaire restent distincts et visibles

**Given** un Finding Draft basé sur un Run timed-out
**When** il est soumis
**Then** timeout, résultats partiels et contradictions sont conservés et aucune confirmation n’est automatique

**Given** aucun modèle IA
**When** le handoff est réalisé
**Then** formulaires, templates, relations et revue humaine couvrent le workflow

## 26. Questions ouvertes
Les packages restent des concepts; OPEN-013 et OPEN-015 restent ouvertes.

## 27. Consommateurs documentaires
Dynamic Sandbox, Evidence Creation/Review, Finding Management, Case Replay, Reporting Preparation, Objets, Trust et Permissions.
