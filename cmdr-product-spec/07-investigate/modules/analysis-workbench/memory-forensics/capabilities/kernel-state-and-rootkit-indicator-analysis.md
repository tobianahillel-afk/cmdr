---
id: CAP-INV-358
title: Kernel State and Rootkit Indicator Analysis
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
  - REQ-PROD-052
  - REQ-PROD-055
  - REQ-AI-002
  - REQ-SEC-001
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-358 — Kernel State and Rootkit Indicator Analysis

## 1. Définition
Inspecter les projections kernel et indicateurs d’altération disponibles selon la plateforme, sans décrire de méthode rootkit, dissimulation, bypass, moteur ou plugin.

## 2. Problème utilisateur
Sans prise en compte du support de plateforme et des projections manquantes, une incohérence kernel peut être présentée à tort comme rootkit confirmé.

## 3. Objectifs
- voir composants, drivers, mécanismes, listes/tables et modifications candidates disponibles.
- comparer plusieurs vues et exposer incohérences/éléments manquants.
- relier processus, modules et régions puis préparer Hypothesis/Evidence candidate.

## 4. Non-objectifs
Aucune technique rootkit, dissimulation, bypass, commande, signature, offset, moteur, API ou action Endpoint.

## 5. Propriétaire
Investigate possède observations/interprétations; Settings/Endpoint Agent exposent support; Studio Tools/Runs; Shared trace; Govern cibles réelles.

## 6. Utilisateurs
Principal : **Memory Forensics Analyst**; secondaires : Kernel Specialist, Evidence Reviewer et Investigation Lead autorisés.

## 7. Conditions d’entrée
Image/session/profil et support de plateforme visibles; Tool/version, permissions, partialité et projections absentes explicites.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Memory Image/Session | Investigate | source/scope | oui | versions liées | blocked |
| Platform/Profile/support | CAP-INV-350; Settings | interprétation/support | oui | courant | unsupported |
| Kernel projections | Tool Calls | observations candidates | non selon support | Tool/version visibles | partial/incompatible |
| Module/process/region observations | CAP-INV-351..353 | contexte | non | même image | unlinked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Memory Image | Investigate | source/limites | lecture |
| Driver/Module/Process Observation | Investigate concepts | relations | lecture/lien |
| Tool/Tool Call/Run | Studio | version/paramètres/statut | lecture |
| Platform support | Settings/Endpoint Agent | support déclaré | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Kernel State Observation | créer/annoter/contester | Investigate concept | support/source obligatoires |
| Rootkit indicator candidate | créer/lier/versionner | Investigate concept | ≠ rootkit confirmé |
| Hypothesis/Evidence relation | préparer | Investigate | qualification future |
| Trace event | émettre | Shared | append-only |

## 11. Fonctionnalités
- inspecter projections kernel disponibles et support associé.
- voir drivers, composants, callbacks/mécanismes fonctionnels, tables/listes exposées.
- comparer vues, éléments manquants et modifications candidates.
- annoter et relier processus/modules/régions.
- préparer Hypothesis/Evidence candidate sans confirmation automatique.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspecter/filtrer | Analyst | Kernel Observation | 0 | support/lecture | vue sourcée | non |
| Comparer vues | Analyst | Observation set | 0/1 | sources compatibles | divergences visibles | non |
| Annoter/classer | Analyst | Candidate | 2 | permission | version conservée | OPEN-013 |
| Préparer handoff | Analyst | Package | 2 | provenance complète | candidate | non |

Classes 3/4 indisponibles.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| produire projections supportées | oui | oui | oui | résumé | Tools déterministes |
| comparer vues | oui | oui | oui | explication | tables/comparateurs |
| signaler incohérences | oui | oui | oui | proposition | règles explicables |
| confirmer rootkit | oui | non | non | non automatique | revue humaine spécialisée |

Attribution complète et disposition humaine obligatoires.

## 14. États fonctionnels
`queued`, `processing`, `partial`, `available`, `failed`, `incompatible`, `disputed`, `superseded`. Machine objet finale reportée.

## 15. États d’interface
Loading/Empty/Partial/Error/Offline/Permission denied/Stale; support manquant et partialité ne sont jamais masqués.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Kernel State Observation | result | Workbench/Hypothesis | support/source visibles |
| Rootkit indicator candidate | candidate | CAP-INV-362 | aucune confirmation automatique |
| Unsupported/partial status | event | Case/Session | aucun résultat fictif |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Module/Session | inspecter kernel | CAP-INV-358 | image, profil, support, Tool | source |
| Observation | corréler | Processes/Modules/Regions | relations, incohérences, provenance | Kernel |
| Candidate | handoff | CAP-INV-362 | éléments pour/contre, limites | Kernel |

Tenant, Case, image, permissions et return origin préservés.

## 18. Dépendances
CAP-INV-350..353/359/362, Settings/Endpoint support, Studio, Shared, OPEN-005/008/013/015.

## 19. Source de vérité
Image/observations : Investigate; support : Settings/Endpoint Agent; Tool/Run : Studio; trace : Shared.

## 20. Provenance et audit
Image/session/profil, support déclaré, Tool/version, paramètres, projections disponibles/manquantes, relations, acteur, erreurs, partialité et disposition.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up | Séparation | Owner | Phase |
|---|---|---|---|---|---|---|
| kernel state read | système sensible | 0 | possible | specialist/reviewer | Investigate/Security | Permissions |
| compare/analyze | traitement | 0/1 | policy | initiateur/reviewer | Investigate | Permissions |
| classify/handoff | mutation | 2 | OPEN-013 | auteur/reviewer | Investigate | Permissions |

Modèle d’accès final reporté.

## 22. Limites et erreurs
- support de plateforme obligatoire et OPEN-008 ouverte.
- driver/incohérence ≠ rootkit confirmé.
- aucune méthode de dissimulation, bypass ou technique rootkit.
- données manquantes visibles.

## 23. Métriques
- résultats supported/partial/incompatible.
- projections manquantes par support.
- candidates confirmées/rejetées humainement.
- provenance complète.

## 24. Classification de livraison
`defined` / `planned`; aucune plateforme, implémentation ou moteur revendiqué.

## 25. Critères d’acceptation
### 1. Support absent
**Given** plateforme non supportée **When** vue kernel ouverte **Then** `unsupported` est visible et aucune projection fictive n’apparaît.
### 2. Incohérence
**Given** vues kernel divergent **When** examinées **Then** l’écart reste candidate et aucun rootkit n’est confirmé.
### 3. Sans IA
**Given** aucun modèle **When** analyse réalisée **Then** projections déterministes, tables et revue spécialisée fonctionnent.

## 26. Questions ouvertes
OPEN-005/008/013/015 restent ouvertes; plateformes/schémas/permissions reportés.

## 27. Consommateurs documentaires
INV-MEM-001, CAP-INV-351..353/359/362, Case/Evidence et futures phases Objects/Permissions/Screens.
