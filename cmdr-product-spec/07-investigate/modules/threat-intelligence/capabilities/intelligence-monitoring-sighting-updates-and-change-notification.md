---
id: CAP-INV-532
title: Intelligence Monitoring, Sighting Updates and Change Notification
product: investigate
module: threat-intelligence
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-06
requirement_ids: [REQ-PROD-014, REQ-PROD-019, REQ-PROD-020, REQ-PROD-055, REQ-INV-006, REQ-AI-002, REQ-SEC-001, REQ-SEC-002, REQ-UX-010]
open_decisions: [OPEN-013, OPEN-018, OPEN-019]
source-of-truth: canonical
---
# CAP-INV-532 — Intelligence Monitoring, Sighting Updates and Change Notification

## 1. Définition
Définir et suivre un Intelligence Monitoring Plan autorisé, enregistrer de nouveaux Sightings et changements significatifs, produire une Change Assessment et notifier des consommateurs autorisés sans créer d’Alert opérationnelle ni surveiller activement une cible.

## 2. Problème utilisateur
Des changements peuvent rester invisibles, tandis qu’une notification Intelligence peut être confondue avec une alerte ou une collecte active.

## 3. Objectifs
Sélectionner connaissances/sources/fréquence ; suivre Sightings, expirations, contradictions et confidence ; produire Change Assessment ; notifier, interrompre et superseder.

## 4. Non-objectifs
Aucune API, protocole, format, provider, code, collecte active sur cible, Signal, Alert, blocage, réponse, partage externe, Cloud/Mobile ou écran détaillé.

## 5. Propriétaire
Investigate possède Monitoring Plan et Change Assessment. CAP-INV-513 conserve Sighting. Shared conserve Notifications, Jobs, Timeline, Trace et Versioning. Settings conserve sources/health/schedules administrés. Studio conserve Runs. Command conserve Signal/Alert/Incident.

## 6. Utilisateurs
Principal : **Threat Intelligence Analyst**. Secondaires : Senior Analyst, Intelligence Manager, SOC Analyst, Detection Engineer, Investigation Lead, Consumer Reviewer et Auditor autorisés.

## 7. Conditions d’entrée
Tenant, environnement, connaissances, sources autorisées, fréquence, consommateurs, permissions, restrictions, owner et return origin sont explicites ; toute lacune reste incomplete/partial/blocked/restricted.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Knowledge and product versions | CAP-INV-507..531 | selected monitored versions | oui | current versions | incomplete |
| Source and health context | Settings / source owners | availability and access | oui | current health | blocked/partial |
| Change criteria and consumers | analyst / CAP-INV-528/529 | bounded intent | oui | plan version | no notification |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Knowledge/Product/Watchlist projection | Investigate | version/lifecycle/restrictions | lecture/lien |
| Source health/schedule | Settings | availability/cadence | lecture limitée |
| Sighting/Notification/Automation Run | Investigate/Shared/Studio | change and provenance | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Intelligence Monitoring Plan | créer, suspendre, reprendre, versionner, superseder | Investigate concept | plan ≠ active target surveillance |
| Change Assessment | créer, contester, réviser, superseder | Investigate concept | change ≠ Alert |
| Sighting update request | préparer/lier | CAP-INV-513 owner | aucune observation inventée |

## 11. Fonctionnalités
Définir éléments/sources/cadence/critères ; voir observations, erreurs et sources indisponibles ; préparer Sighting et Change Assessment ; notifier ; interrompre et superseder.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter/filtrer/comparer | Analyst | monitoring/sources/changes | 0 | lecture autorisée | vue sourcée | non |
| exécuter/demander monitoring borné | Analyst | Plan/Automation Run | 1 | sources/cadence autorisées | résultats/erreurs visibles | selon politique |
| créer plan, assessment, notification | Analyst | concepts locaux | 2 | owner/consumers explicites | version réversible | OPEN-013/019 |
| surveiller activement cible ou runtime | aucun rôle local | cible/runtime | 3 | hors périmètre | aucune exécution | obligatoire |
| supprimer provenance | aucun rôle local | trace | 4 | interdit | refus audité | strict |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| définir Monitoring Plan | oui | formulaires/règles | oui | proposition sourcée | checklist |
| détecter différence/Sighting | oui | diff/filtres/jobs | oui | résumé incertain | timeline |
| préparer assessment/notification | oui | template | oui | draft attribué | matrice/revue |
| conclure malveillance ou créer Alert | owner humain | contrôles | jamais autonome | jamais décisionnaire | Command/Govern |

Toute automatisation expose initiateur, moteur/version, Tool Calls, Run, sources, paramètres, erreurs, incertitude et disposition humaine.

## 14. États fonctionnels
`draft`, `ready`, `active-bounded`, `paused`, `partial`, `source-unavailable`, `change-detected`, `assessment-required`, `completed`, `stopped`, `superseded`.

## 15. États d’interface
Loading conserve contexte ; Empty distingue absence/non-collecte ; Partial nomme sources manquantes ; Error conserve valide ; Offline stale/read-only ; Permission denied ne révèle rien ; Stale expose dates ; Conflict offre diff/recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Change Assessment | assessment versionné | CAP-INV-529/534/536/537 | changements/sources/limites visibles |
| Sighting update context | handoff | CAP-INV-513 | observation sourcée |
| Consumer notification context | notification event | authorized consumers | notification ≠ Alert |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Monitoring Plan | observation/version change | CAP-INV-513/532 | source, time, version, errors | Plan |
| CAP-INV-513 | new/revised Sighting | CAP-INV-532 | occurrence, context, confidence | Sighting |
| CAP-INV-532 | significant change reviewed | Shared Notifications / CAP-INV-534/536 | assessment, consumers, markings | Change Assessment |

Chaque transition conserve owner, tenant, versions, permissions, restrictions, erreurs, autorité, provenance et return origin.

## 18. Dépendances
CAP-INV-513/517/528..531/534/536/537 ; Shared ; Settings ; Studio ; OPEN-013/018/019.

## 19. Source de vérité
Investigate est source du Plan/Assessment ; sources, Sightings, notifications et runtimes restent chez leurs owners.

## 20. Provenance et audit
Conserver plan, critères, cadence, sources, versions, Sightings, changes, Runs, Tool Calls, erreurs, interruptions, consommateurs, notifications, décisions humaines et timestamps.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Monitoring Plan create/run/review | unauthorized collection/disclosure | 1/2 | sources/consumers scoped | OPEN-013/019 | author/reviewer | Investigate/Settings | Permissions |
| Change notification | context leakage | 2 | minimized payload | possible | assessor/notifier | Shared/Investigate | Permissions |

## 22. Limites et erreurs
Monitoring ≠ surveillance active ; Change Notification ≠ Alert ; absence de Sighting ≠ absence de menace ; source unavailable, stale, timeout et superseded restent visibles.

## 23. Métriques conceptuelles
Plans par état/source ; changes avec provenance ; Alert/collecte active créée localement — cible zéro ; notification non autorisée — cible zéro.

## 24. Classification de livraison
`defined` / `planned` ; aucune surveillance active n’est revendiquée.

## 25. Critères d’acceptation
### 1. Source indisponible
**Given** une source autorisée indisponible  
**When** le monitoring s’exécute  
**Then** résultat partial, erreur visible, aucune absence de menace conclue.

### 2. Changement significatif
**Given** nouveau Sighting et confidence révisée  
**When** Change Assessment est créée  
**Then** sources/versions/contradictions restent visibles et aucune Alert n’est créée.

### 3. Sans IA
**Given** aucun modèle  
**When** le workflow est exécuté  
**Then** diff, jobs, timelines, formulaires et revue humaine suffisent.

## 26. Questions ouvertes
OPEN-013/018/019 restent ouvertes ; aucun choix technique/politique final.

## 27. Consommateurs documentaires
Threat Intelligence, Investigate, Detection Engineering, Command, Settings, Studio, Govern, Shared, Quality et Roadmap. Aucun Cloud/Mobile n’est lancé.
