---
id: CAP-INV-318
title: Behavioral Observation and Timeline
product: investigate
module: dynamic-sandbox
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-INV-005
  - REQ-PROD-020
  - REQ-AI-002
  - REQ-UX-007
open_decisions:
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-318 — Behavioral Observation and Timeline

## 1. Définition
Organiser les événements d’un Sandbox Run dans une timeline permission-aware en distinguant événement observé, inférence et annotation analyste, avec timestamps, sources, objets associés, erreurs et lacunes.

## 2. Problème utilisateur
Une timeline qui mélange événements capturés, inférences et commentaires transforme une hypothèse en fait et empêche la revue ou la reproduction.

## 3. Objectifs
- Organiser les événements observés, inférences et annotations analyste dans une timeline sourcée.
- Conserver Case, Artifact, session, Run, versions, limites et return origin.
- Exposer erreurs, résultats partiels, provenance et disposition humaine.
- Préparer les transitions autorisées sans qualification automatique.

## 4. Non-objectifs
Ne pas définir virtualisation, instrumentation, hooks, API, protocole, commande, code, moteur ou hyperviseur; ne pas commencer Reverse Engineering, Debugger ou forensics; ne pas agir sur un Endpoint réel.

## 5. Propriétaire
Investigate possède le contexte analytique et l’interprétation. Platform Settings possède l’administration des Sandbox Environments; Studio possède Tool, Tool Call et Automation Run; Govern possède l’autorité sur les cibles réelles; Shared possède les mécanismes transversaux.

## 6. Utilisateurs
Malware Analyst; Dynamic Analysis Operator; Case Analyst; Evidence Reviewer.

## 7. Conditions d’entrée
Case et Artifact accessibles; Dynamic Analysis Session ou Run identifiable; environnement autorisé; permissions, policies, versions et restrictions réévaluées.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Artifact source | Investigate / CAP-INV-105 | entrée exécutée en sandbox | oui | version immuable ou référencée | bloquer ou marquer partial |
| Case et objectif | Investigate | contexte analytique | oui | état courant | rester draft |
| Dynamic Analysis Session | concept Investigate | scope et paramètres fonctionnels | oui | version de session | créer ou reprendre |
| Sandbox Environment projection | Platform Settings | environnement autorisé et health | oui | réévaluée avant Run | environment-unavailable |
| Flux d’événements | Sandbox Run producers | événements observés, source et timestamp | oui pour timeline | fraîcheur visible | timeline partial |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | scope et return origin | consulter |
| Artifact / Sandbox Run | Investigate concepts | source, versions et contexte Run | consulter |
| Sandbox Environment / Tool | Settings / Studio | health, limites et version producteur | consulter seulement |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Dynamic Analysis Session | lier résultat/disposition | Investigate concept | aucun schéma final |
| Behavioral Observation | créer/annoter/supersede | Investigate concept | source Run obligatoire |
| Trace/relations | émettre/lier | Shared mechanism / owner semantics | aucun transfert d’ownership |

## 11. Fonctionnalités
- Organiser événements observés, inférences et annotations dans une timeline sourcée.
- Afficher sources, versions, paramètres, états et limites.
- Conserver les relations au Run, au Case et aux résultats associés.
- Distinguer observation, inférence, annotation, candidat, Evidence et Finding.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Lire la timeline | Analyst | Behavioral Timeline | 0 | permission et Run accessibles | événements affichés | non |
| Filtrer/grouper | Analyst | Timeline view | 0 | sources disponibles | vue conservée | non |
| Annoter | Analyst | Observation | 2 | observation accessible | annotation attribuée | OPEN-013 |
| Qualifier comme inférence | Reviewer | Observation | 2 | justification | type et justification visibles | OPEN-013 |
| Sélectionner pour handoff | Analyst | Selection set | 2 | sources visibles | sélection sourcée | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Organiser la timeline | oui | règles/timestamps | oui | regroupement expliqué | timeline et filtres |
| Résumer les résultats | oui | agrégation sourcée | oui | résumé attribué | champs bruts |
| Proposer une prochaine observation | oui | checklists/profils | oui | suggestion modifiable | catalogue manuel |
| Qualifier Evidence/Finding | humain | contrôles seulement | workflow de revue | jamais autonome | CAP-INV-107/108/109 |

Toute automatisation expose initiateur, producteur/version, Automation Run et Tool Calls, sources, paramètres, timestamp, statut, incertitude, owner humain et acceptation, modification ou rejet. Aucun Run n’est lancé silencieusement.

## 14. États fonctionnels
`collecting`, `available`, `partial`, `clock-skew`, `source-missing`, `inferred`, `annotated`, `disputed`, `superseded`. Machines finales reportées à la phase Objets.

## 15. États d’interface
Loading conserve le contexte; Empty explique; Partial nomme les lacunes; Error garde les résultats valides; Offline est stale/read-only; Permission denied ne fuit rien; unsafe ou stale restent visibles.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Behavioral Timeline | Timeline projection | analyste/Case | observed, inferred et annotation distincts |
| Behavioral Observations | concepts d’observation | CAP-INV-319..328 | sources et timestamps |
| Timeline selection | selection set | CAP-INV-328 | relations et incertitude |
| Completeness event | quality event | Reviewer | lacunes explicites |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Sandbox Run | événements reçus | CAP-INV-318 | Run, source, timestamps, object refs, errors | Run |
| CAP-INV-318 | ouvrir relation | CAP-INV-319/320/321/322/324 | observation, objet, timestamp, return origin | timeline |
| CAP-INV-318 | préparer handoff | CAP-INV-328 | observations sélectionnées, sources, contradictions | timeline |

Chaque transition conserve tenant, Case, Artifact, session, Run, versions, permissions, erreurs et return origin.

## 18. Dépendances
CAP-INV-105, CAP-INV-301..313, CAP-INV-314..328 selon le flux, CAP-INV-107/108/109, CAP-INV-213/214, Settings Sandbox Environments, Studio Tools/Automation Runs, Shared Jobs/Trace/Timeline/Linking/Export et OPEN applicables.

## 19. Source de vérité
Investigate est source de l’interprétation et des relations analytiques; Settings, Studio, Govern et Shared restent sources de leurs objets. Une projection ne transfère jamais l’ownership.

## 20. Provenance et audit
Enregistrer Case, session, Artifact/version, Sandbox Environment/version, profil, Run, Tools/Tool Calls, initiateur, paramètres, timestamps, observations, erreurs, interruption, nettoyage, reset, décisions humaines et disposition finale.

## 21. Permissions fonctionnelles
Dynamic Analysis read; Sandbox Run read; behavioral results read; timeline annotate; observation classify; selection for handoff. La matrice atomique, le step-up et la séparation des tâches sont reportés; cross-tenant est refusé par défaut.

## 22. Limites et erreurs
Artifact indisponible ou restreint; environnement unavailable, unhealthy, resetting ou policy-blocked; version/profile/Tool manquant; clock skew; événements manquants; timeout, crash, résultat partiel, permission révoquée ou tenant mismatch.

## 23. Métriques
Ouvertures, états partial/failed, événements sans source, inférences contestées, provenance complète, propositions acceptées/modifiées/rejetées et handoffs revus.

## 24. Classification de livraison
`defined` / `planned`; aucune preuve d’implémentation, moteur, hyperviseur, API, protocole, commande ou outil tiers final.

## 25. Critères d’acceptation
**Given** un Sandbox Run contenant événements et inférences
**When** l’analyste ouvre Behavioral Observation and Timeline
**Then** observed, inferred et annotations sont distincts, sourcés et navigables

**Given** des événements manquants ou une permission refusée
**When** la timeline est ouverte
**Then** les lacunes sont explicites, les données valides restent visibles et aucune observation n’est inventée

**Given** aucun modèle IA
**When** le workflow est réalisé
**Then** timeline, filtres, règles et revue humaine couvrent toutes les fonctions essentielles

## 26. Questions ouvertes
Behavioral Observation reste un concept; OPEN-013 et OPEN-015 restent ouvertes; aucun moteur n’est choisi.

## 27. Consommateurs documentaires
Dynamic Sandbox, Analysis Workbench, Case Workspace, Evidence Board, Settings Sandbox Environments, Studio Control Room et phases Objets, Permissions, Journeys, Trust et Technique.
