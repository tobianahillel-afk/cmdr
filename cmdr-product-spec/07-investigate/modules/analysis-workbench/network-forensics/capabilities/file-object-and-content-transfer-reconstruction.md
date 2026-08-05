---
id: CAP-INV-390
title: File, Object and Content Transfer Reconstruction
product: investigate
module: analysis-workbench
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-05
requirement_ids:
  - REQ-INV-001
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-OBJ-003
  - REQ-AI-002
  - REQ-SEC-001
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-390 — File, Object and Content Transfer Reconstruction

## 1. Définition
Identifier un transfert candidat et reconstruire, lorsque les fragments disponibles le permettent, un fichier, objet ou contenu avec source, direction candidate, endpoints, timestamps, taille, parties manquantes, erreurs et limites visibles, sans jamais exécuter le contenu.

## 2. Problème utilisateur
Une reconstruction incomplète peut perdre son nom, son contexte ou des fragments et être présentée à tort comme le fichier original complet ou comme une Evidence qualifiée.

## 3. Objectifs
- identifier le transfert candidat et sa conversation source ;
- afficher direction candidate, endpoints, timestamps et tailles déclarées ou reconstruites ;
- afficher fragments manquants, conflits, erreurs et métadonnées disponibles ;
- reconstruire un objet lorsque possible et le marquer partiel si nécessaire ;
- produire un Derived Artifact sourcé et router explicitement vers Static ou Reverse ;
- permettre le retrait de l’usage actif sans supprimer la trace.

## 4. Non-objectifs
Ne pas exécuter le contenu, appeler une cible, rejouer une transaction, contourner un chiffrement, extraire ou utiliser un secret, choisir un protocole ou moteur, créer une règle Detection, qualifier une Evidence, créer un objet Intelligence, définir une API, une commande, un format final ou un écran détaillé.

## 5. Propriétaire
Investigate / Analysis Workbench / Investigate Product Lead possède Transfer Observation, Reconstructed Object Candidate, Derived Artifact et disposition humaine. La capture et les résultats bruts restent liés à leurs producteurs. Studio possède Tool, Tool Call et Automation Run. Shared possède Trace, Activity, Linking, Export et Versioning. Govern possède toute autorité réelle.

## 6. Utilisateurs
Principal : **Network Artifact Analyst**. Secondaires : Static Analyst, Reverse Engineer, Evidence Reviewer, Investigation Lead.

## 7. Conditions d’entrée
Case, source, conversation ou transaction accessible ; capture, coverage, timebase et restrictions visibles ; permission d’accès et d’extraction réévaluée ; Tool/version explicites ; aucun appel à une cible réelle.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Conversation ou transaction source | CAP-INV-386/388 | messages, direction candidate et fragments | oui | version de session | rester `candidate` ou bloquer |
| Capture, coverage et timebase | CAP-INV-382/383 | source et limitations | oui | snapshot lié | marquer `partial` |
| Permission et restrictions payload | Security / Settings | accès, minimisation et export | oui | réévaluées | masquer ou refuser |
| Tool, version et paramètres | CMDR Studio | reconstruction déterministe | oui pour traitement | version du lancement | `tool-unavailable` |
| Case et Hypothesis | Investigate | objectif et return origin | oui | état courant | rester draft |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case / Hypothesis | Investigate | objectif, restrictions et return origin | consulter/lier |
| Capture Artifact | Investigate / producteur | source, version, acquisition et limites | consulter uniquement |
| Conversation Candidate / Application Transaction | Investigate | fragments, direction, endpoints et timestamps | consulter/sélectionner |
| Tool / Tool Call / Automation Run | CMDR Studio | version, paramètres, statut et résultats | sélectionner/invoquer/lire |
| Artifact / Derived Artifact | Investigate | parent, restrictions et lineage | consulter/comparer |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Transfer Observation | créer, annoter, contester ou superseder | Investigate | conversation source et incertitude obligatoires |
| Reconstructed Object Candidate | créer, comparer ou retirer de l’usage actif | Investigate | partialité et fragments manquants visibles |
| Derived Artifact | créer via CAP-INV-395 | Investigate | parent/enfant, Tool/version, paramètres et restrictions conservés |
| Trace / Activity event | émettre | Shared | append-only avec correlation ID |

## 11. Fonctionnalités
- détecter ou sélectionner un transfert candidat ;
- inspecter fragments, métadonnées, direction, tailles et erreurs ;
- reconstruire de manière déterministe lorsque possible ;
- comparer le résultat à d’autres Artifacts ;
- créer un Derived Artifact sans exécution ;
- afficher les limites de reconstruction et conserver le lien à la conversation ;
- router vers Static, Reverse ou CAP-INV-395 ;
- fonctionner sans fournisseur de modèle.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter le transfert | Network Artifact Analyst | Transfer Observation | 0 | lecture autorisée | fragments et limites visibles | non |
| Reconstruire un objet | Network Artifact Analyst | Tool Call / Reconstructed Object Candidate | 1 | scope, Tool/version et permission explicites | résultat attribué et éventuellement partiel | selon policy |
| Annoter ou contester | Network Artifact Analyst | observation ou candidat | 2 | permission réversible | nouvelle disposition versionnée | OPEN-013 |
| Retirer de l’usage actif | Network Artifact Analyst | Reconstructed Object Candidate | 2 | justification | retrait sans suppression de trace | OPEN-013 |
| Préparer une extraction ou un handoff | Network Artifact Analyst | candidate package | 2 | source et limites présentes | package non qualifié | owner destination |

Les classes 3 et 4 sont exclues. Aucun contenu n’est exécuté et aucune interaction cible n’est permise.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Identifier des transferts candidats | oui | parsers et règles | oui | suggestion attribuée | filtres et inspection manuelle |
| Reconstruire les fragments | oui | reconstructeur déterministe | oui | explication des limites | reconstructeur et journal brut |
| Comparer des objets | oui | comparateur | oui | résumé sourcé | diff et tables |
| Expliquer un échec | oui | codes et diagnostics | oui | résumé facultatif | erreurs brutes et checklist |
| Qualifier Evidence, Finding ou IOC | owner humain | contrôles seulement | revue | jamais autonome | CAP-INV-107/108/109 |

Toute automatisation expose initiateur, producteur/version, Tool Calls, Automation Run, sources, paramètres, timestamp, statut, erreurs, incertitude et disposition humaine.

## 14. États fonctionnels
`candidate`, `queued`, `processing`, `available`, `partial`, `missing-fragments`, `conflicting`, `invalid`, `restricted`, `failed`, `disputed`, `superseded`, `withdrawn-from-use`. États fonctionnels, pas machine objet définitive.

## 15. États d’interface
Loading conserve conversation et sélection ; Empty distingue absence de transfert et fragments absents ; Partial expose les parties manquantes ; Error conserve les fragments valides ; Offline bloque un nouveau traitement ; Permission denied masque payload et contenu ; Stale expose source ou Tool obsolète.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Transfer Observation | observation Investigate | Workbench / Case | source, direction candidate et limites visibles |
| Reconstructed Object Candidate | candidat | CAP-INV-395 / comparaison | partialité et erreurs conservées |
| Derived Artifact context | package d’extraction | CAP-INV-311/395 | parent, Tool/version et fragments manquants visibles |
| Static/Reverse handoff | transition | CAP-INV-301 ou CAP-INV-329 | aucun contenu exécuté et return origin conservé |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-386/388 | transfert candidat | CAP-INV-390 | conversation, transaction, direction, fragments et restrictions | Workbench avec sélection préservée |
| CAP-INV-390 | créer un Derived Artifact | CAP-INV-395 | source, fragments, Tool/version, limites et Case | CAP-INV-390 |
| Derived Artifact | analyser sans exécution | CAP-INV-301 / CAP-INV-329 | Artifact, parent, restrictions, provenance et return origin | CAP-INV-390 |
| CAP-INV-390 | préparer Evidence/Finding | CAP-INV-397 | candidat, parties manquantes, contradictions et provenance | CAP-INV-390 |

Chaque transition conserve tenant, Case, ownership, restrictions, erreurs, permissions et return origin.

## 18. Dépendances
CAP-INV-105/107/108/109/208/212/213/214/301..313/311/329..346/386/388/395/397 ; Studio Tool/Tool Call/Automation Run ; Shared Trace/Activity/Linking/Export/Versioning ; OPEN-005/008/013/014/015.

## 19. Source de vérité
Investigate est source de Transfer Observation, Reconstructed Object Candidate et disposition humaine. La capture et les fragments bruts restent liés à leur producteur. Studio reste source des Tools et Runs ; Shared des traces et versions.

## 20. Provenance et audit
Conserver Case, capture/version, Collection Request/Job, custody, session, capteur/interface, coverage/timebase, conversation/transaction, fragments, direction candidate, Tool/version/Calls/Run, paramètres, accès payload, erreurs, résultats partiels, Derived Artifact, annotations et dispositions.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Transfer metadata read | données réseau et contenu | 0 | metadata sensible masquée | possible | viewer/reviewer | Investigate | Permissions |
| Reconstruction run | volume et accès aux fragments | 1 | payload masqué selon droit | possible | initiateur/reviewer | Investigate / Studio | Permissions/Technique |
| Content preview/read | données privées ou secrets | 0 | masqué par défaut | possible | lecteur distinct de l’exporteur | Security / Investigate | Permissions |
| Derived Artifact create/export | extraction et fuite | 1/2 | minimisation et redaction | possible | auteur/reviewer | Investigate / Shared Export | Permissions |
| Annotate/dispute/withdraw/handoff | mutation analytique | 2 | justification protégée | OPEN-013 | auteur/reviewer | Investigate | Permissions |

La matrice atomique, les namespaces, RBAC/ABAC et règles finales de step-up restent futurs.

## 22. Limites et erreurs
Fragments manquants, capture tronquée, direction ambiguë, conflit de taille, nom absent, contenu restreint, Tool indisponible, timeout, interruption, permission révoquée ou tenant mismatch. Reconstructed object ≠ original certain ; incomplete reassembly ≠ transfert complet ; Derived Artifact ≠ Evidence.

## 23. Métriques
Transferts candidats, reconstructions disponibles/partielles/échouées, fragments manquants, objets retirés, provenance complète, accès sensibles refusés et handoffs qualifiés ultérieurement. Aucune cible chiffrée définitive.

## 24. Classification de livraison
`defined` / `planned`. Preuve documentaire seulement ; aucun moteur, produit tiers, format, protocole final, API, commande, code ou écran détaillé.

## 25. Critères d’acceptation
### 1. Objet transféré partiel
**Given** un transfert candidat avec des fragments manquants
**When** la reconstruction produit un objet
**Then** le résultat est `partial`, les fragments manquants sont visibles et la source reste inchangée

### 2. Contenu restreint
**Given** un candidat contenant un payload restreint
**When** un utilisateur sans permission l’ouvre
**Then** les métadonnées autorisées restent visibles et le contenu, la copie et l’extraction sont bloqués

### 3. Sans IA
**Given** aucun fournisseur de modèle
**When** la reconstruction est réalisée
**Then** parsers, reconstructeur, diff, journal et revue humaine couvrent le workflow

## 26. Questions ouvertes
OPEN-005 moteurs ; OPEN-008 support ; OPEN-013 classe 2 ; OPEN-014 Artifact/Attachment ; OPEN-015 provenance des Runs. Les objets, permissions atomiques, écrans et contrats techniques restent futurs.

## 27. Consommateurs documentaires
Network Workbench, Static Analysis, Reverse Engineering, Derived Artifact Management, Evidence/Finding handoff, Case Workspace, phases Objets, Permissions, Écrans, Technique et Validation.
