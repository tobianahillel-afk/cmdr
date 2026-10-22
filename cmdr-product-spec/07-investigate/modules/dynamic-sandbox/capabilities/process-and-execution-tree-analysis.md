---
id: CAP-INV-319
title: Process and Execution Tree Analysis
product: investigate
module: dynamic-sandbox
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-INV-005
  - REQ-PROD-014
  - REQ-PROD-020
open_decisions:
  - OPEN-005
  - OPEN-013
source-of-truth: canonical
---
# CAP-INV-319 — Process and Execution Tree Analysis

## 1. Définition
Analyser le processus initial, les relations parent/enfant, créations, terminaisons, identités fonctionnelles, paramètres déclarés, relations temporelles et liens vers fichiers ou réseau, sans définir l’instrumentation ni fournir de commande.

## 2. Problème utilisateur
Un arbre incomplet ou non sourcé peut conduire à attribuer une intention malveillante à une relation d’exécution ordinaire.

## 3. Objectifs
- Afficher le processus initial et les relations parent/enfant.
- Relier chaque nœud aux timestamps, fichiers, connexions et observations.
- Comparer plusieurs arbres sans conclure automatiquement.
- Préparer un futur handoff Reverse ou Debugger sans commencer ces capabilities.

## 4. Non-objectifs
Ne pas définir instrumentation, collecte technique, commande, désassemblage, décompilation, debugger, patching, moteur ou hyperviseur.

## 5. Propriétaire
Investigate possède l’interprétation du process tree; Studio possède les Tools; Settings l’environnement; Shared les mécanismes de trace.

## 6. Utilisateurs
Malware Analyst; Dynamic Analysis Operator; Case Analyst; Reviewer.

## 7. Conditions d’entrée
Sandbox Run identifiable; événements processus disponibles ou lacunes explicites; Artifact, environnement, Tools et permissions accessibles.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Processus initial | Sandbox Run | racine observée | oui | timestamp Run | arbre orphaned |
| Événements création/fin | Behavioral Timeline | relations d’exécution | oui | ordonnés | partial |
| Identités et paramètres déclarés | Tool result | contexte fonctionnel | non | source/version visibles | unknown |
| Relations fichiers/réseau | CAP-INV-320/321 | pivots analytiques | non | même Run | relation absente explicite |
| Run de comparaison | CAP-INV-325 | second arbre | non | préconditions visibles | comparaison indisponible |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Sandbox Run | Investigate concept | source, profil, environnement et timestamps | consulter |
| Behavioral Observation | Investigate concept | créations, terminaisons et relations | consulter |
| Artifact / Runtime Artifact | Investigate | fichiers liés | consulter/lier |
| Network Observation | Investigate concept | connexions associées | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Process Observation | créer/annoter/supersede | Investigate concept | source event et Run obligatoires |
| Execution Tree result | créer/versionner | Investigate concept | nœuds orphelins visibles |
| Future handoff context | préparer | Investigate | aucune capability Reverse/Debugger créée |

## 11. Fonctionnalités
Voir racine, parents/enfants, créations, terminaisons, identités, paramètres autorisés, relations temporelles, fichiers et connexions; filtrer, sélectionner, annoter et comparer.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Ouvrir l’arbre | Analyst | Execution Tree | 0 | Run accessible | arbre sourcé | non |
| Filtrer/sélectionner | Analyst | Process Observation | 0 | nœuds visibles | sélection conservée | non |
| Annoter une anomalie | Analyst | Observation | 2 | source visible | annotation attribuée | OPEN-013 |
| Comparer deux arbres | Reviewer | Execution Trees | 0 | préconditions comparables | diff sourcé | non |
| Préparer handoff futur | Analyst | Handoff context | 2 | Artifact/process sélectionné | contexte seulement | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Construire les relations | oui | IDs/timestamps | oui | non nécessaire | arbre déterministe |
| Signaler nœuds orphelins | oui | règles | oui | explication | filtres et statuts |
| Résumer l’arbre | oui | agrégation | oui | résumé attribué | navigation manuelle |
| Qualifier un comportement | humain | contrôles seulement | workflow de revue | jamais autonome | annotation/revue |

## 14. États fonctionnels
`collecting`, `available`, `partial`, `orphaned-node`, `identity-unknown`, `terminated`, `disputed`, `superseded`.

## 15. États d’interface
Loading conserve le Run; Empty distingue absence d’événement et source indisponible; Partial montre les nœuds manquants; Error garde les branches valides; Offline est read-only; Permission denied masque les paramètres sensibles.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Execution Tree | Analysis Result | analyste/Case | relations et sources visibles |
| Process Observations | concepts analytiques | CAP-INV-320/321/322 | Run et timestamps conservés |
| Comparison input | tree snapshot | CAP-INV-325 | version et profil visibles |
| Future handoff context | context package | 4B.2B.2B future | aucun outil ou commande défini |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-318 | sélectionner événement processus | CAP-INV-319 | Run, event, process refs, timestamp | timeline |
| CAP-INV-319 | ouvrir fichier/réseau | CAP-INV-320/321 | process, observation, Run, return origin | tree |
| CAP-INV-319 | préparer future analyse | future Reverse/Debugger | Artifact, process, static/dynamic results, provenance | tree |

## 18. Dépendances
CAP-INV-317/318/320/321/322/325; Studio Tool Calls; Shared Graph/Trace/Timeline; OPEN-005/013.

## 19. Source de vérité
Investigate possède l’arbre analytique; les événements producteurs et Tool Calls restent sourcés; aucune observation ne devient une intention attribuée.

## 20. Provenance et audit
Run, environment/version, process identities, parent-child edges, timestamps, Tool/version, missing events, annotations, comparisons et dispositions.

## 21. Permissions fonctionnelles
Process tree read; sensitive parameters read; process observation annotate; multi-Run comparison; future handoff prepare; cross-tenant denied.

## 22. Limites et erreurs
Processus racine absent; nœud orphelin; événements hors ordre; identité ambiguë; paramètres redacted; clock skew; Tool output partiel; permission révoquée.

## 23. Métriques
Arbres complets/partiels; nœuds orphelins; processus sans identité; pivots fichiers/réseau; comparaisons; annotations contestées.

## 24. Classification de livraison
`defined` / `planned`; aucune instrumentation, commande, moteur, API, protocole ou implémentation.

## 25. Critères d’acceptation
**Given** un Run avec créations et terminaisons de processus
**When** l’analyste ouvre Process and Execution Tree Analysis
**Then** racine, relations, timestamps, sources et objets associés sont navigables

**Given** un nœud sans parent observable
**When** l’arbre est construit
**Then** il reste orphaned et n’est pas rattaché artificiellement

**Given** aucun modèle IA
**When** l’arbre est analysé
**Then** relations déterministes, filtres et revue humaine couvrent le workflow

## 26. Questions ouvertes
Process Observation et Execution Tree restent des concepts; OPEN-005 et OPEN-013 restent ouvertes.

## 27. Consommateurs documentaires
Dynamic Sandbox, Behavioral Timeline, File/System and Network Analysis, Multi-Run Comparison, future Reverse/Debugger, Objets et Permissions.
