---
id: CAP-INV-321
title: Network Behavior Analysis
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
  - REQ-SEC-002
open_decisions:
  - OPEN-005
  - OPEN-013
source-of-truth: canonical
---
# CAP-INV-321 — Network Behavior Analysis

## 1. Définition
Analyser les connexions, résolutions, destinations, protocoles fonctionnels déclarés, volumes, timestamps, processus associés, erreurs, tentatives bloquées et réponses simulées d’un Sandbox Run, en distinguant simulation et réseau réel.

## 2. Problème utilisateur
Une destination observée n’est pas automatiquement un IOC. Une sortie réseau réelle non visible ou silencieusement autorisée peut exposer des tiers ou fausser l’interprétation.

## 3. Objectifs
- Afficher connexions, résolutions, destinations, volumes, processus et erreurs.
- Distinguer réseau simulé, bloqué et réel explicitement autorisé.
- Préparer des indicateurs candidats sans confirmation automatique.
- Déclencher un signal de sécurité si une connectivité réelle inattendue est observée.

## 4. Non-objectifs
Ne pas définir protocoles internes, instrumentation, packet capture technique, network forensics avancé, commandes, moteur, infrastructure ou méthode de contournement.

## 5. Propriétaire
Investigate possède l’interprétation réseau du Run; Settings possède les policies d’environnement; Studio possède les Tools; Shared possède trace/export.

## 6. Utilisateurs
Malware Analyst; Dynamic Analysis Operator; Case Analyst; Evidence Reviewer.

## 7. Conditions d’entrée
Run identifiable; politique réseau visible; observations réseau disponibles ou absence explicite; permission de lecture des résultats.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Politique réseau | Sandbox Environment | simulated/blocked/authorized scope | oui | version Run | policy-unknown |
| Connexions et résolutions | Tool result | observations réseau | oui pour analyse | timestamps visibles | partial |
| Processus associés | CAP-INV-319 | attribution | non | même Run | process-unknown |
| Réponses simulées/bloquées | environnement/Tool | contexte de simulation | non | même Run | unknown |
| Restrictions d’export | policy/classification | protection | oui pour export | réévaluées | export bloqué |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Sandbox Run | Investigate concept | profil, environnement, phases | consulter |
| Sandbox Environment | Platform Settings | politique réseau et isolation déclarée | consulter seulement |
| Process Observation | Investigate concept | processus associé | consulter/lier |
| Tool Call | Studio | producteur/version/output | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Network Observation | créer/annoter/supersede | Investigate concept | source Run/Tool obligatoire |
| Indicator Candidate | proposer/exclure | concept Investigate/future Intelligence | aucune confirmation automatique |
| Network safety event | émettre | Investigate → Settings | connectivité réelle inattendue visible |

## 11. Fonctionnalités
Voir connexions, résolutions, destinations, protocoles déclarés, volumes, timestamps, processus, erreurs, tentatives bloquées et réponses simulées; filtrer, grouper, comparer et préparer un candidat.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Lire résultats réseau | Analyst | Network Observation | 0 | Run/permission | vue sourcée | non |
| Filtrer/grouper/comparer | Analyst | Network set | 0 | données accessibles | vue/diff | non |
| Proposer indicateur candidat | Analyst | Candidate | 2 | source/contexte | candidat non confirmé | OPEN-013 |
| Exclure faux candidat | Analyst | Candidate | 2 | justification | disposition attribuée | OPEN-013 |
| Arrêter si réseau inattendu | Operator | Sandbox Run | 2 | observation/policy mismatch | stop demandé | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Classer simulated/blocked/real | oui | policy + source | oui | explication | champs bruts |
| Grouper destinations | oui | règles explicables | oui | suggestion | filtres |
| Proposer candidat | oui | patterns | oui | suggestion modifiable | sélection humaine |
| Confirmer IOC/Finding | humain | contrôles seulement | workflow de revue | jamais autonome | future Intelligence/Evidence review |

## 14. États fonctionnels
`collecting`, `available`, `partial`, `simulated`, `blocked`, `allowed`, `unexpected-real-network`, `resolution-failed`, `disputed`, `superseded`.

## 15. États d’interface
Loading conserve le Run; Empty distingue aucune observation et source absente; Partial nomme pertes/erreurs; Error garde les données valides; Offline interdit nouveau Run; Permission denied redacted.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Network observation set | Analysis Result | analyste/Case | policy, source et process visibles |
| Indicator candidates | candidate set | future Intelligence/CAP-INV-328 | non confirmés |
| Network safety event | event | CAP-INV-326/Settings | mismatch explicite |
| Comparison input | snapshot | CAP-INV-325 | préconditions visibles |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Sandbox Run | observation réseau | CAP-INV-321 | Run, policy, destination, process, timestamp | Run |
| CAP-INV-321 | candidat/handoff | future Intelligence ou CAP-INV-328 | valeur, contexte, simulation, source, incertitude | network view |
| CAP-INV-321 | réseau réel inattendu | CAP-INV-326/Settings | environment, policy, observation, stop status | network view |

## 18. Dépendances
CAP-INV-317/318/319/325/326/328; Settings Network Policy/Health; Studio Tools; Shared Trace/Export; OPEN-005/013.

## 19. Source de vérité
Investigate possède l’interprétation; Settings reste source de la policy réseau; une destination observée n’est pas un IOC confirmé.

## 20. Provenance et audit
Run, environment/version, policy, source Tool/version, process, destination/résolution, timestamps, volumes, simulated/blocked/real disposition, annotations et handoffs.

## 21. Permissions fonctionnelles
Network behavior read; sensitive destination read; network profile select; indicator candidate prepare; Run stop; export; cross-tenant denied.

## 22. Limites et erreurs
Résolution échouée; données tronquées; process inconnu; policy stale; simulation ambiguë; connectivité réelle inattendue; redaction; permission révoquée.

## 23. Métriques
Observations simulated/blocked/allowed; réseau inattendu; résolutions échouées; candidats proposés/rejetés; Runs arrêtés; provenance complète.

## 24. Classification de livraison
`defined` / `planned`; aucun moteur, protocole, instrumentation, commande, network forensics avancé ou implémentation.

## 25. Critères d’acceptation
**Given** un Run avec réseau simulé et connexions bloquées
**When** l’analyste ouvre Network Behavior Analysis
**Then** simulation, blocage, destinations, processus et policy sont distincts et visibles

**Given** une connectivité réelle non prévue
**When** elle est observée
**Then** elle est signalée, le Run peut être arrêté et Settings reçoit le contexte sans silence

**Given** aucun modèle IA
**When** l’analyse est réalisée
**Then** policies, règles, filtres et revue humaine couvrent le workflow

## 26. Questions ouvertes
Network Observation reste un concept; OPEN-005 et OPEN-013 restent ouvertes; Intelligence appartient à une phase future.

## 27. Consommateurs documentaires
Dynamic Sandbox, Safety, Multi-Run Comparison, future Intelligence, Evidence/Finding handoff, Settings, Objets et Permissions.
