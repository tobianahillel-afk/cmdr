---
id: CAP-INV-394
title: Multi-Capture, Sensor and Time-Window Comparison
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
  - REQ-UX-006
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-394 — Multi-Capture, Sensor and Time-Window Comparison

## 1. Définition
Comparer plusieurs captures, capteurs et fenêtres temporelles avec scopes, interfaces, couverture, timebases, flows, conversations, transactions, Entities, objets transférés, pertes, duplications et incompatibilités explicitement conservés.

## 2. Problème utilisateur
Comparer des sources de couverture différente sans préconditions peut transformer `non observé` en `absent` et supprimer silencieusement des duplications utiles.

## 3. Objectifs
- sélectionner captures, capteurs, fenêtres et relations
- comparer scopes, interfaces, coverage et timebases
- comparer flows, conversations, transactions, Entities et objets transférés
- distinguer common, added, removed, modified, not-observed et incomparable
- voir pertes, duplications, incompatibilités, enregistrer la comparaison et préparer Finding Draft

## 4. Non-objectifs
Aucune acquisition active, administration de capteur/Fleet, packet generation/injection/crafting, scanning/interception active, replay, interaction cible, déchiffrement non autorisé, extraction ou usage de secret, exploit, évasion, règle Detection, objet Intelligence canonique, API, protocole interne, moteur, commande, code ou écran détaillé.

## 5. Propriétaire
Investigate possède Network Comparison Result, l’interprétation et les packages candidats. Command conserve Detection/Signal/Alert/Incident. Collection/Endpoint Agent produit la capture et ses limites. Settings administre capteurs/Fleet/Policies/storage/retention/health/timebase/secrets. Studio possède Tool/Tool Call/Workflow/Automation Run. Govern possède l’autorité réelle. Shared possède Entity/Graph/Timeline et les mécanismes génériques.

## 6. Utilisateurs
Principal : **Comparison Analyst**. Secondaires : Network Forensics Analyst, Timeline Analyst, Investigation Lead.

## 7. Conditions d’entrée
Case, source et session accessibles ; provenance, coverage, timebase, restrictions et permissions visibles ; scope borné ; aucune interaction active avec une cible.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Captures, capteurs, fenêtres, coverage et timebases comparables | Capture Artifacts and Network Sessions | source et contexte spécifiques | oui | versions de session | `partial` ou bloqué |
| Case, Hypothesis et objectif | Investigate | contexte analytique | oui | état courant | rester draft |
| Tool/version et paramètres | Studio | traitement déterministe | oui pour classe 1 | réévalués au lancement | `tool-unavailable` |
| Permissions, payload policy et restrictions | Security/Settings | accès et minimisation | oui | décision courante | données masquées ou action refusée |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case / Hypothesis | Investigate | objectif, restrictions, return origin | consulter/lier |
| Capture Artifacts and Network Sessions | owner source / Investigate | contenu autorisé, version, limitations | consulter uniquement |
| Tool / Tool Call / Automation Run | Studio | version, paramètres, statut, résultats | sélectionner/invoquer/lire |
| Entity / Graph / Timeline / Trace | Shared | relations, ordre, provenance | consommer sans redéfinir |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Network Comparison Result | créer, annoter, contester, superseder | Investigate | source, incertitude, auteur et version requis |
| Relation/sélection analytique | créer ou modifier réversiblement | Investigate | aucune fusion ou qualification silencieuse |
| Trace/Activity event | émettre | Shared | append-only et correlation ID |

## 11. Fonctionnalités
- sélectionner captures, capteurs, fenêtres et relations
- comparer scopes, interfaces, coverage et timebases
- comparer flows, conversations, transactions, Entities et objets transférés
- distinguer common, added, removed, modified, not-observed et incomparable
- voir pertes, duplications, incompatibilités, enregistrer la comparaison et préparer Finding Draft
- préserver source brute autorisée, partialité, restrictions, erreurs et return origin ;
- fonctionner sans fournisseur de modèle.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter, filtrer ou comparer | Comparison Analyst | Network Comparison Result | 0 | lecture autorisée | vue sourcée | non |
| Comparer captures, capteurs et fenêtres | Comparison Analyst | Tool Call / résultat | 1 | scope, Tool/version, permission | résultat attribué, éventuellement partial | selon policy |
| Annoter, contester ou relier | Comparison Analyst | Network Comparison Result | 2 | permission réversible | nouvelle disposition versionnée | OPEN-013 |
| Préparer un handoff | Comparison Analyst | package candidat | 2 | sources et limites présentes | package non qualifié | owner destination |

Classes 3/4 indisponibles ; toute action réelle est bloquée ou redirigée vers Collection/Live Response et Govern.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Identifier/grouper Network Comparison Result | oui | parsers/règles/Tools | oui | suggestion attribuée | viewer, tables et filtres |
| Comparer/reconstruire | oui | oui si préconditions | oui | proposition incertaine | comparateur/reconstructeur |
| Expliquer erreur/contradiction | oui | catalogue/contrôles | oui | résumé sourcé | erreurs brutes/checklist |
| Confirmer IOC/Evidence/Finding/règle/Intelligence | owner humain | contrôles seulement | revue | jamais autonome | capabilities propriétaires |

Toute automatisation expose initiateur, moteur/agent et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude et acceptation/modification/rejet.

## 14. États fonctionnels
`draft`, `ready`, `processing`, `available`, `partial`, `incompatible`, `timebase-conflict`, `coverage-limited`, `disputed`, `superseded`. États fonctionnels, pas machine objet définitive.

## 15. États d’interface
Loading conserve contexte/sélection ; Empty distingue absence observée et donnée absente ; Partial expose gaps/troncatures ; Error conserve les résultats valides ; Offline bloque les traitements ; Permission denied masque brut/payload ; Stale expose source/timebase/Tool obsolète.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Network Comparison Result | Network Comparison Result | Network Workbench, Case ou owner destination | source, partialité et provenance visibles |
| Difference relation set | Difference relation set | Network Workbench, Case ou owner destination | source, partialité et provenance visibles |
| Finding Draft context | Finding Draft context | Network Workbench, Case ou owner destination | source, partialité et provenance visibles |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-381 / Capture Artifacts | préparer comparaison | CAP-INV-394 | sources, scopes, sensors, time windows et permissions | Workbench avec return origin |
| CAP-INV-394 | corréler les différences | CAP-INV-393 | differences, sources, timebases et uncertainty | Workbench avec return origin |
| CAP-INV-394 | préparer un draft | CAP-INV-397 | selected differences, contradictions et provenance | Workbench avec return origin |

Les transitions conservent tenant, Case, ownership, restrictions, erreurs, permissions et return origin.

## 18. Dépendances
Case/CAP-INV-105/107/108/109, CAP-INV-208/212/213/214, related CAP-INV-380..397, Static CAP-INV-301..313, Reverse CAP-INV-329..346, Memory CAP-INV-347..362, Disk CAP-INV-363..379, Studio Tool/Tool Call/Automation Run, Settings sensor/Fleet/Policy/health, Shared Entity/Graph/Timeline/Trace/Jobs/Export/Recovery, Govern and OPEN-005/008/013/014/015.

## 19. Source de vérité
Investigate est source de Network Comparison Result et de la disposition humaine. Les captures/résultats bruts restent chez leurs producteurs ; Studio, Settings, Shared, Command et Govern conservent leurs sources propriétaires.

## 20. Provenance et audit
Case, capture/version, Collection Request/Job, custody, session, capteur/interface, scope/coverage/timebase, Tool/version/Calls/Run, paramètres/filtres, entrées/sorties, erreurs, partialité, accès sensibles, annotations, décisions et disposition.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Network Comparison Result metadata/read | corrélation réseau | 0 | classification | possible | viewer/reviewer | Investigate | Permissions |
| Network Comparison Result processing/read | volume, brut ou payload | 1 | payload masqué par défaut | possible | initiateur/reviewer | Investigate/Studio | Permissions/Technique |
| Network Comparison Result annotate/link/handoff | mutation analytique | 2 | valeurs protégées | OPEN-013 | auteur/reviewer | Investigate | Permissions |

Matrice atomique, namespaces, RBAC/ABAC et step-up final reportés.

## 22. Limites et erreurs
- source partial/corrupted/truncated/restricted/stale, gaps, timebase incertaine, Tool indisponible, timeout, permission ou tenant mismatch ;
- `Non observé` ne signifie pas absent lorsque la couverture diffère.
- Aucune source n’est fusionnée ou modifiée.
- Les duplications candidates et timebases restent attribuées.
- aucun résultat Tool, score ou sortie IA ne vaut conclusion.

## 23. Métriques
Usages, états partial/failed/disputed, gaps et limites héritées, provenance complète, accès sensibles refusés, dispositions humaines et reproductibilité ; aucune cible chiffrée définitive.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire seulement. Aucun moteur, produit, format, protocole final, API, commande, modèle ML, écran détaillé ou code.

## 25. Critères d’acceptation
### 1. Scénario
**Given** deux capteurs avec packets candidats dupliqués et timestamps légèrement différents
**When** l’analyste compare
**Then** les sources restent distinctes et aucune suppression silencieuse n’a lieu

### 2. Scénario
**Given** deux captures incomparables par scope
**When** la comparaison est lancée
**Then** l’état `incompatible` ou `partial` est visible et aucun diff absolu n’est affirmé

### 3. Scénario
**Given** aucun modèle IA
**When** la comparaison est réalisée
**Then** comparateurs déterministes et revue humaine couvrent le workflow

## 26. Questions ouvertes
OPEN-005 moteurs ; OPEN-008 support capteurs/plateformes ; OPEN-013 classe 2 ; OPEN-014 Artifact/Attachment lorsque pertinent ; OPEN-015 bridge des Runs. OPEN-011/012 restent ouvertes et hors périmètre. Objets, permissions atomiques, écrans et contrats techniques sont futurs.

## 27. Consommateurs documentaires
Network Workbench, Case/Evidence/Finding, Event Search, Dynamic Sandbox, Memory/Disk, Static/Reverse, Entity Graph, future Detection Engineering/Intelligence et phases Objets/Permissions/Écrans/Journeys/Technique/Validation.
