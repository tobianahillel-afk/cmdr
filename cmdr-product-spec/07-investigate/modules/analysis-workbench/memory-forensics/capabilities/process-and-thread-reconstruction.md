---
id: CAP-INV-351
title: Process and Thread Reconstruction
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
  - REQ-AI-002
  - REQ-SEC-001
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-351 — Process and Thread Reconstruction

## 1. Définition
Permet de reconstruire les processus, threads et relations candidates observables dans une image mémoire en conservant incohérences et limites, sans moteur, plugin ni implémentation imposés.

## 2. Problème utilisateur
Sans reconstruction multi-vue, les processus ou threads absents d’une projection peuvent être déclarés cachés à tort, alors que l’image ou la méthode peut être partielle.

## 3. Objectifs
- reconstruire processus, threads, relations parent/enfant et contextes candidats.
- Préserver Case, image, limites, incertitude, provenance et return origin.
- Séparer reconstruction, état runtime courant, Evidence et Finding.

## 4. Non-objectifs
Aucune acquisition, commande, méthode offensive, moteur, plugin, offset, algorithme, API, protocole, action Endpoint ou Disk/Filesystem/full Network Forensics.

## 5. Propriétaire
Investigate possède contexte et interprétation; Endpoint Agent l’acquisition; Settings l’administration; Studio Tools/Runs; Govern les cibles réelles; Shared les mécanismes transversaux.

## 6. Utilisateurs
Principal : **Memory Forensics Analyst**. Secondaires : Investigation Lead, Evidence Reviewer et Audit Analyst.

## 7. Conditions d’entrée
Memory Image, session, profil, permissions et Tool/version visibles. Toute absence devient `partial`, `unsupported` ou `blocked`.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Memory Image et Session | Investigate | source, scope et restrictions | oui | versions liées | blocked ou partial |
| Selected platform/profile | CAP-INV-350 | interprétation et limites | oui | sélection courante | unsupported |
| Reconstruction results | Tool Calls / analyst | projections sourcées | oui | Tool/version visibles | partial ou failed |
| Case, permissions, policies | Investigate/Security/Settings | contexte et autorisation | oui | courant | denied |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Memory Image / Artifact | Investigate | source, limites, lineage | lecture |
| Process / Thread Observation | Investigate concept | projection, source, confiance | lecture |
| Tool / Tool Call / Automation Run | Studio | version, paramètres, statut | lecture |
| Case / Hypothesis | Investigate | contexte et relations | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Process Observation | créer/annoter/contester/superseder | Investigate concept | aucun schéma final |
| Thread Observation | créer/lier/versionner | Investigate concept | source et incertitude obligatoires |
| Trace / Activity event | émettre | Shared | append-only |
| Source objects | aucune mutation | owners respectifs | projections seulement |

## 11. Fonctionnalités
- reconstruire processus, threads et relations candidates.
- voir identités, timestamps, états, modules, connexions et handles associés.
- voir processus terminés, orphelins et divergences entre méthodes.
- filtrer, comparer, annoter et relier à une Hypothesis.
- préparer un handoff sans qualification automatique.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspecter/filtrer | Analyst | Process Observation | 0 | lecture autorisée | vue sourcée | non |
| Comparer vues/images | Analyst | Comparison | 0 | projections lisibles | divergences visibles | non |
| Annoter/contester/relier | Analyst | Observation | 2 | permission/justification | version conservée | OPEN-013 |
| Lancer traitement borné | Analyst | Tool Call | 1 | scope/Tool/policy explicites | job tracé | non |

Classes 3/4 bloquées et routées vers Collection/Live Response et Govern.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| reconstruire processus/threads | oui | oui | oui | proposition | reconstructeurs déterministes et tables |
| afficher relations | oui | oui | oui | résumé | arbre avec alternative tabulaire |
| comparer méthodes | oui | oui | oui | explication | comparateur multi-vue |
| qualifier conclusion | oui | non | non | assistance | revue humaine |

Attribution obligatoire : initiateur, producteur/version, Tool Calls, Run, sources, paramètres, statut, erreurs, incertitude et disposition humaine.

## 14. États fonctionnels
`queued`, `processing`, `partial`, `available`, `failed`, `incompatible`, `disputed`, `superseded`. États objet finaux reportés.

## 15. États d’interface
Loading conserve le contexte; Empty n’invente rien; Partial expose les manques; Error conserve le valide; Offline limite les mutations; Permission denied masque; Stale distingue ancien/courant.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Process Observation | Reconstruction result | Memory Workbench | source, partialité et incertitude visibles |
| Thread Observation | candidate relation | Case/Hypothesis | aucune conclusion automatique |
| Handoff selection | candidate context | CAP-INV-362 | provenance conservée |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Session/Profile | reconstruire | CAP-INV-351 | image, profil, Tool/version, limites | Session |
| Process | inspecter relations | CAP-INV-352..355 | identifiants, relations, timestamps, provenance | Processes |
| Observation | préparer handoff | CAP-INV-362 | candidates, contradictions, incertitude | Processes |

Tenant, Case, image, permissions, sélection et return origin sont préservés.

## 18. Dépendances
CAP-INV-350, CAP-INV-352..355, CAP-INV-359, CAP-INV-362, Studio, Shared, OPEN-005/008/013/015. Aucune dépendance bas niveau.

## 19. Source de vérité
Image/contexte : Investigate; acquisition : Collection/Endpoint Agent; administration : Settings; Tools/Runs : Studio; Trace/Timeline : Shared.

## 20. Provenance et audit
Case, Endpoint, image, session, profil, Tool/version, Calls/Run, paramètres, méthodes, acteur, timestamps, partialité, erreurs, annotations et dispositions; aucune suppression silencieuse.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up | Séparation | Owner | Phase |
|---|---|---|---|---|---|---|
| process/thread read | données système sensibles | 0 | possible | reviewer si requis | Investigate/Security | Permissions |
| comparison | corrélation sensible | 0/1 | selon policy | initiateur/reviewer | Investigate/Security | Permissions |
| annotation/dispute | mutation réversible | 2 | OPEN-013 | auteur/reviewer | Investigate | Permissions |

Matrice atomique, RBAC/ABAC, step-up et SoD finaux reportés.

## 22. Limites et erreurs
- processus reconstruit ≠ processus live.
- absence d’une vue ≠ dissimulation.
- image partielle et méthodes divergentes visibles.
- aucune action Endpoint ou qualification malveillante automatique.

## 23. Métriques
- résultats par état et taux partial/failed.
- divergences multi-vue.
- observations avec provenance complète.
- candidates acceptées/modifiées/rejetées.

## 24. Classification de livraison
`defined` / `planned`; aucune preuve d’implémentation, moteur, plugin ou plateforme.

## 25. Critères d’acceptation
### 1. Incohérence processus
**Given** un processus apparaît dans une méthode et pas une autre **When** l’analyste inspecte **Then** l’écart et l’impact de l’image partielle sont visibles sans déclarer un processus caché.
### 2. Résultat partiel
**Given** des structures manquent **When** le traitement termine **Then** `partial` est visible et aucune donnée n’est inventée.
### 3. Sans IA
**Given** aucun modèle **When** l’analyste travaille **Then** reconstructeurs, tables, filtres et annotations fonctionnent.

## 26. Questions ouvertes
OPEN-005, OPEN-008, OPEN-013 et OPEN-015 restent ouvertes. Schémas, plateformes et permissions finales sont reportés.

## 27. Consommateurs documentaires
INV-MEM-001, Case/Evidence, CAP-INV-352..362, futures phases Objects/Permissions/Journeys/Screens/Contracts.
