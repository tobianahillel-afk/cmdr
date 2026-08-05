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
  - REQ-AI-002
  - REQ-SEC-001
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-015
  - OPEN-014
source-of-truth: canonical
---
# CAP-INV-390 — File, Object and Content Transfer Reconstruction

## 1. Définition
Identifier un transfert candidat et reconstruire lorsque possible un fichier, objet ou contenu à partir d’une conversation, avec direction, endpoints, timestamps, tailles, parties manquantes, erreurs, metadata et partialité visibles, sans exécution.

## 2. Problème utilisateur
Une reconstruction incomplète peut perdre nom, contexte ou fragments et être présentée à tort comme fichier original complet ou Evidence.

## 3. Objectifs
- identifier transfert candidat, source, direction, endpoints et timestams
- voir taille déclarée ou reconstruite, parties manquantes, erreurs et conflits
- reconstruire un objet lorsque possible sans exécuter son contenu
- produire un Derived Artifact marqué partial si nécessaire
- comparer avec d’autres Artifacts, router vers Static/Reverse et retirer de l’usage actif

## 4. Non-objectifs
Aucune acquisition active, administration de capteur/Fleet, packet generation/injection/crafting, scanning/interception active, replay, interaction cible, déchiffrement non autorisé, extraction ou usage de secret, exploit, évasion, règle Detection, objet Intelligence canonique, API, protocole interne, moteur, commande, code ou écran détaillé.

## 5. Propriétaire
Investigate possède Transfer Observation and Reconstructed Object Candidate, l’interprétation et les packages candidats. Command conserve Detection/Signal/Alert/Incident. Collection/Endpoint Agent produit la capture et ses limites. Settings administre capteurs/Fleet/Policies/storage/retention/health/timebase/secrets. Studio possède Tool/Tool Call/Workflow/Automation Run. Govern possède l’autorité réelle. Shared possède Entity/Graph/Timeline et les mécanismes génériques.

## 6. Utilisateurs
Principal : **Network Artifact Analyst**. Secondaires : Static Analyst, Reverse Engineer, Evidence Reviewer.

## 7. Conditions d’entrée
Case, source et session accessibles ; provenance, coverage, timebase, restrictions et permissions visibles ; scope borné ; aucune interaction active avec une cible.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Conversation/transaction et fragments de transfert candidats | Conversation or Application Transaction | source et contexte spécifiques | oui | versions de session | `partial` ou bloqué |
| Case, Hypothesis et objectif | Investigate | contexte analytique | oui | état courant | rester draft |
| Tool/version et paramètres | Studio | traitement déterministe | oui pour classe 1 | réévalués au lancement | `tool-unavailable` |
| Permissions, payload policy et restrictions | Security/Settings | accès et minimisation | oui | décision courante | données masquées ou action refusée |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case / Hypothesis | Investigate | objectif, restrictions, return origin | consulter/lier |
| Conversation or Application Transaction | owner source / Investigate | contenu autorisé, version, limitations | consulter uniquement |
| Tool / Tool Call / Automation Run | Studio | version, paramètres, statut, résultats | sélectionner/invoquer/lire |
| Entity / Graph / Timeline / Trace | Shared | relations, ordre, provenance | consommer sans redéfinir |

## 10. Objets crép�s ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Transfer Observation and Reconstructed Object Candidate | créer, annoter, contester, superseder | Investigate | source, incertitude, auteur et version requis |
| Relation/sélection analytique | créer ou modifier réversiblement | Investigate | aucune fusion ou qualification silencieuse |
| Trace/Activity event | émettre | Shared | append-only et correlation ID |

## 11. Fonctionnalités
- identifier transfert candidat, source, direction, endpoints et timestamps
- voir taille déclarée ou reconstruite, parties manquantes, erreurs et conflits
- reconstruire un objet lorsque possible sans exécuter son contenu
- produire un Derived Artifact marqué partial si nécessaire
- comparer avec d’autres Artifacts, router vers Static/Reverse et retirer de l’usage actif
- préserver source brute autorisée, partialité, restrictions, erreurs et return origin ;
- fonctionner sans fournisseur de modèle.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter, filtrer ou comparer | Network Artifact Analyst | Transfer Observation and Reconstructed Object Candidate | 0 | lecture autorisée | vue sourcée | non |
| Reconstruire un transfert candidat | Network Artifact Analyst | Tool Call / résultat | 1 | scope, Tool/version, permission | résultat attribvé, éventuellement partial | selon policy |
| Annoter, contester ou relier | Network Artifact Analyst | Transfer Observation and Reconstructed Object Candidate | 2 | permission réversible | nouvelle disposition versionnée | OPEN-013 |
| Préparer un handoff | Network Artifact Analyst | package candidat | 2 | sources et limites présentes | package non qualifié | owner destination |

Classes 3/4 indisponibles ; toute action réelle est bloquée ou redirigée vers Collection/Live Response et Govern.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Identifier/grouper Transfer Observation and Reconstructed Object Candidate | oui | parsers/règles/Tools | oui | suggestion attribvée | viewer, tables et filtres |
| Comparer/reconstruire | oui | oui si préconditions | oui | proposition incertaine | comparateur/reconstructeur |
| Expliquer erreur/contradiction | oui | catalogue/contrôles | oui | résumé sourcé | erreurs brutes/checklist |
| Confirmer IOC/Evidence/Finding/règle/Intelligence | owner humain | contrôles seulement | revue | jamais autonome | capabilities proprétaires |

Toute automatisation expose initiateur, moteur/agent et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude et acceptation/modification/rejet.

## 14. États fonctionnels
`candidate`, `processing`, `available`, `partial`, `missing-fragments`, `invalid`, `restricted`, `disputed`, `superseded`, `withdrawn-from-use`. États fonctionnels, pas machine objet définitive.

## 15. États d’interface
Loading conserve contexte/sélection ; Empty distingue absence observée et donnée absente ; Partial expose gaps/troncatures ; Error conserve les résultats valides ; Offline bloque les traitements ; Permission denied masque brut/payload ; Stale expose source/timebase/Tool obsolète.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Transfer Observation | Transfer Observation | Network Workbench, Case ou owner destination | source, partialité et provenance visibles |
| Reconstructed Object Candidate | Reconstructed Object Candidate | Network Workbench, Case ou owner destination | source, partialité et provenance visibles |
| Derived Artifact | Derived Artifact | Network Workbench, Case ou owner destination | source, partialité et provenance visibles |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-388 / CAP-INV-386 | transfert candidat | CAP-INV-390 | conversation/transaction, direction, fragments, metadata et restrictions | Workbench avec return origin |
| CAP-INV-390 | produire un Artifact | CAP-INV-395 | source transfer, reconstruction, limitations et Tool context | Workbench avec return origin |
| Derived Artifact | analyser sans exécution | CAP-INV-301 / CAP-INV-329 | Artifact, parent, Case, restrictions et return origin | Workbench avec return origin |

Les transitions conservent tenant, Case, ownership, restrictions, erreurs, permissions et return origin.

## 18. Dépendances
Case/CAP-INV-105/107/108/109, CAP-INV-208/212/213/214, related CAP-INV-380..397, Static CAP-INV-301..313, Reverse CAP-INV-329..346, Memory CAP-INV-347..362, Disk CAP-INV-363..379, Studio Tool/Tool Call/Automation Run, Settings sensor/Fleet/Policy/health, Shared Entity/Graph/Timeline/Trace/Jobs/Export/Recovery, Govern and OPEN-005/008/013/014/015.

## 19. Source de vérité
Investigate est source de Transfer Observation and Reconstructed Object Candidate et de la disposition humaine. Les captures/résultats bruts restent chez leurs producteurs ; Studio, Settings, Shared, Command et Govern conservent leurs sources propriétaires.

## 20. Provenance et audit
Case, capture/version, Collection Request/Job, custody, session, capteur/interface, scope/coverage/timebase, Tool/version/Calls/Run, paramètres/filtres, entrées/sorties, erreurs, partialité, accès sensibles, annotations, décisions et disposition.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|--:|---|---|---|---|---|
| Transfer Observation and Reconstructed Object Candidate metadata/read | corrélation réseau | 0 | classification | possible | viewer/reviewer | Investigate | Permissions |
| Transfer Observation and Reconstructed Object Candidate processing/read | volume, brut ou payload | 1 | payload masqué par défaut | possible | initiateur/reviewer | Investigate/Studio | Permissions/Technique |
| Transfer Observation and Reconstructed Object Candidate annotate/link/handoff | mutation analytique | 2 | valeurs protégées | OPEN-013 | auteur/reviewer | Investigate | Permissions |

Matrice atomique, namespaces, RBAC/ABAC et step-up final reportés.

## 22. Limites et erreurs
- source partial/corrupted/truncated/restricted/stale, gaps, timebase incertaine, Tool indisponible, timeout, permission ou tenant mismatch ;
- Incomplete reassembly ≠ transfert complet.
- Reconstructed file ≠ fichier original certain.
- Le content n’est jamais exécuté et ne devient pas automatiquement Evidence.
- aucun résultat Tool, score ou sortie IA ne vaut conclusion.

## 23. Métriques
Usages, états partial/failed/disputed, gaps et limites héritées, provenance complète, accès sensibles refusés, dispositions humaines et reproductibilité ; aucune cible chiffrée définitive.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire seulement. Aucun moteur, produit, format, protocole final, API, commande, modèle ML, écran détaillé ou code.

## 25. Critères d’acceptation
### 1. Scénario
**Given** un transfert candidat dont des fragments manquent
**When** une extraction autorisée produit un Derived Artifact
**Then** le résultat est `partial`, les parties manquantes et la source sont visibles

### 2. Scénario
**Given** un contenu reconstruit restreint
**When** l’analyste tente de le prévisualiser
**Then** les permissions sont réévaluées et aucun contenu n’est révélé sans droit

### 3. Scénario
**Given** aucun modèle IA
**When** la reconstruction est lancée
**Then** reassembly déterministe, contrôles et revue humaine fonctionnent

## 26. Questions ouvertes
OPEN-005 moteurs ; OPEN-008 support capteurs/plateformes ; OPEN-013 classe 2 ; OPEN-014 Artifact/Attachment lorsque pertinent ; OPEN-015 bridge des Runs. OPEN-011/012 restent ouvertes et hors périmètre. Objets, permissions atomiques, écrans et contrats techniques sont futurs.

## 27. Consommateurs documentaires
Network Workbench, Case/Evidence/Finding, Event Search, Dynamic Sandbox, Memory/Disk, Static/Reverse, Entity Graph, future Detection Engineering/Intelligence et phases Objets/Permissions/Écrans/Journeys/Technique/Validation.
