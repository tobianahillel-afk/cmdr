---
id: CAP-INV-380
title: Network Forensics Intake and Preconditions
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
# CAP-INV-380 — Network Forensics Intake and Preconditions

## 1. Définition
Ouvrir une capture ou un Artifact réseau depuis un Case, Collection Job ou Artifact Management, exposer sa provenance, son scope, sa couverture, ses restrictions et ses préconditions, puis créer ou reprendre explicitement une Network Forensics Session sans démarrage automatique.

## 2. Problème utilisateur
Une capture partielle, tronquée, restreinte ou non supportée peut être analysée comme complète si l’intake masque ses limites, ce qui fausse toutes les reconstructions aval.

## 3. Objectifs
- ouvrir depuis Case, Collection Job ou Artifact Management
- voir source, période, taille, intégrité, scope et couverture
- voir pertes, troncatures, interfaces, directions et restrictions
- voir résultats antérieurs, Tools disponibles et permissions
- définir objectif/scope et créer ou reprendre une session

## 4. Non-objectifs
Ne définit ni acquisition active, capteur, produit tiers, moteur, format, commande, API, protocole interne, packet crafting, injection, replay, interception, exploit, déchiffrement non autorisé, extraction/utilisation de secret, règle de détection, objet Intelligence canonique, écran détaillé ou implémentation.

## 5. Propriétaire
Investigate / Analysis Workbench / Investigate Product Lead possède le contexte, l’interprétation, les observations et packages candidats. Command conserve Detection, Signal, Alert et Incident. Endpoint Agent/Collection produit les résultats bruts. Settings administre capteurs, Fleet, Policies, stockage, rétention, health, time synchronization et secrets. Studio possède Tool, Tool Call, Workflow et Automation Run. Govern possède l’autorité sur les cibles réelles. Shared possède Entity, Graph, Timeline et mécanismes transversaux.

## 6. Utilisateurs
Principal : **DFIR Analyst**. Secondaires : Case Analyst, Evidence Reviewer, Investigation Lead.

## 7. Conditions d’entrée
Case et source accessibles, tenant/environnement cohérents, provenance et limites visibles, permission fonctionnelle réévaluée, objectif et scope bornés, aucune interaction active avec une cible.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Network Capture Artifact | Investigate / owner projection | source analytique et contexte | oui | version liée à la session | bloquer ou marquer partial |
| Case, objectif et Hypothesis | Investigate | contexte et question analytique | oui | état courant | rester draft |
| Contexte Network Intake Assessment | analyste / résultat précédent | scope, sélection et interprétation | oui | snapshot versionné | demander complétude |
| Tool, version, permission et restrictions | CMDR Studio / Security / Settings | capacité et autorisation | oui avant traitement | réévaluées au lancement | tool-unavailable ou permission-denied |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case / Hypothesis | Investigate | objectif, scope, restrictions et return origin | consulter et lier |
| Network Capture Artifact | Investigate ou source propriétaire | source, version, acquisition, limitations et lineage | consulter uniquement |
| Tool / Tool Call / Automation Run | CMDR Studio | compatibilité, version, paramètres, résultat et attribution | sélectionner/invoquer/lire selon permission |
| Entity / Timeline / Trace projections | Shared Capabilities | relations, ordering, provenance et audit | consommer sans redéfinir |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Network Intake Assessment | créer, annoter, contester ou superseder selon capability | Investigate | source, incertitude, auteur et version obligatoires |
| Derived relation / selection | créer ou modifier réversiblement | Investigate | ne transfère ni ownership ni permission |
| Trace / Activity event | émettre | Shared Capabilities | append-only et correlation ID |

## 11. Fonctionnalités
- ouvrir depuis Case, Collection Job ou Artifact Management
- voir source, période, taille, intégrité, scope et couverture
- voir pertes, troncatures, interfaces, directions et restrictions
- voir résultats antérieurs, Tools disponibles et permissions
- définir objectif/scope et créer ou reprendre une session
- conserver les sources, restrictions, erreurs, incertitudes et return origin ;
- fonctionner sans fournisseur de modèle.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter Network Intake Assessment | DFIR Analyst | Network Intake Assessment | 0 | source et lecture autorisées | vue sourcée et limites visibles | non |
| Exécuter traitement borné pour Network Intake Assessment | DFIR Analyst | Tool Call / Analysis Result | 1 | scope, Tool/version et permission explicites | résultat partiel ou complet attribué | non, sauf policy |
| Annoter ou contester Network Intake Assessment | DFIR Analyst | Network Intake Assessment | 2 | permission de mutation réversible | nouvelle disposition versionnée | OPEN-013 |
| Préparer un handoff | DFIR Analyst | Candidate package | 2 | sources et limitations présentes | package non qualifié automatiquement | destination owner |

Les classes 3 et 4 ne sont jamais exécutées dans Network Forensics. Toute action sur une cible réelle est bloquée ou redirigée vers les capabilities propriétaires et Govern.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Identifier ou grouper pour Network Intake Assessment | oui | règles/Tools explicables | oui | suggestion attribuée | viewer, filtres et sélection manuelle |
| Comparer ou reconstruire Network Intake Assessment | oui | oui lorsque les préconditions existent | oui | proposition avec incertitude | comparateur ou reconstructeur déterministe |
| Expliquer erreur ou contradiction de Network Intake Assessment | oui | catalogue et contrôles | oui | résumé sourcé | erreurs brutes et checklist |
| Qualifier Evidence, Finding, IOC, règle ou objet Intelligence | humain owner | contrôles seulement | workflow de revue | jamais autonome | CAP-INV-107/108/109 et phases futures |

Aucun Tool, protocole, Entity, payload, anomalie ou qualification n’est choisi silencieusement. Toute sortie automatisée expose initiateur, agent ou moteur, version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet.

## 14. États fonctionnels
`draft`, `incomplete`, `ready`, `partial-capture`, `corrupted`, `truncated`, `restricted`, `unsupported`, `coverage-unknown`, `tool-unavailable`, `policy-blocked`. Ces états sont fonctionnels et ne constituent pas une machine d’état objet définitive.

## 15. États d’interface
Loading conserve le Case, la source et la sélection ; Empty distingue absence observée, données absentes et scope vide ; Partial expose pertes, gaps, troncatures et résultats utilisables ; Error conserve les résultats valides ; Offline interdit tout nouveau traitement non garanti ; Permission denied masque brut et payload ; Stale affiche source, timebase, policy ou Tool obsolète. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Network Intake Assessment | Network Intake Assessment | Network Workbench / analyste | source, scope, partialité et auteur visibles |
| Session creation context | relation, événement ou package | Case / Timeline / capability suivante | ownership destination conservé |
| Visible limitation set | projection ou candidate | CAP-INV-397 / Evidence / future phase | aucune qualification automatique |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Case / Collection Job / Artifact | ouvrir la capture | CAP-INV-380 | Case, capture, acquisition, custody, restrictions, return origin | source ou Workbench avec sélection préservée |
| CAP-INV-380 | intégrité et scope suffisants | CAP-INV-382 | capture, provenance, limites et objectif | source ou Workbench avec sélection préservée |
| CAP-INV-380 | créer ou reprendre | CAP-INV-381 | capture(s), Case, objectif, scope et sélection | source ou Workbench avec sélection préservée |

Chaque transition conserve tenant, Case, source, permissions, restrictions, erreurs, return origin et ownership de la destination.

## 18. Dépendances
CAP-INV-102/105/107/108/109/208/212/213/214/301..313/321/329..346/355/359/362/371/374/378/379 selon le handoff ; CMDR Studio Tool/Tool Call/Automation Run ; Platform Settings capteurs/Fleet/Policies/storage/retention/health/timebase ; Shared Entity/Graph/Timeline/Jobs/Trace/Activity/Search/Linking/Export/Reporting/Recovery ; Govern pour toute autorité réelle ; OPEN-005/008/013/014/015 selon le front matter.

## 19. Source de vérité
Investigate est source de Network Intake Assessment, de l’interprétation et de la disposition humaine. La capture, l’acquisition et les résultats bruts restent liés à leurs producteurs. Studio reste source de Tool/Tool Call/Automation Run ; Settings des projections administratives ; Shared des mécanismes transversaux ; Command de Detection/Signal/Alert/Incident ; Govern de l’autorité réelle.

## 20. Provenance et audit
Enregistrer Case, source/capture et version, Collection Request/Job lorsque présents, acquisition/custody, session, capteur/interface, scope/coverage/timebase, Tool/version/Tool Calls/Automation Run, paramètres/filtres/requêtes, entrées, sorties, erreurs, partialité, accès aux données sensibles, annotations, décisions humaines, disposition et correlation ID.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Network Intake Assessment read | données réseau et corrélations | 0 | masquage selon classification | possible | viewer/reviewer selon policy | Investigate | Permissions |
| Network Intake Assessment analyse/reconstruction | traitement et volume | 1 | payload ou brut masqué si nécessaire | possible | initiateur distinct du reviewer si requis | Investigate/Studio | Permissions/Technique |
| Network Intake Assessment annotate/link/handoff | mutation analytique | 2 | valeurs sensibles protégées | OPEN-013 | auteur/reviewer | Investigate | Permissions |

La matrice atomique, les namespaces, le RBAC/ABAC et les règles finales de step-up restent reportés.

## 22. Limites et erreurs
- source inaccessible, partial, corrupted, truncated, stale, restricted ou superseded ;
- coverage ou timebase insuffisante, packets manquants, données malformées ou décodage incompatible ;
- Tool/version indisponible, timeout, interruption, permission révoquée ou tenant mismatch ;
- aucune absence, identité, protocole, intention, IOC, anomalie, Evidence ou Finding n’est inventé.
- Aucune analyse ne démarre automatiquement.
- Une capture partielle, corrompue, tronquée, restreinte ou unsupported reste explicitement signalée.
- L’intake ne définit ni acquisition, ni capteur, ni format.

## 23. Métriques
Mesurer usages de Network Intake Assessment, taux partial/failed/disputed, pertes et limitations héritées, temps jusqu’à annotation ou handoff, résultats avec provenance complète, accès sensibles refusés, candidates acceptées/modifiées/rejetées par l’owner et reproductions réussies ou bloquées. Aucune cible numérique n’est fixée.

## 24. Classification de livraison
`defined` / `planned`. La preuve est documentaire uniquement. Aucun moteur Network Forensics, produit tiers, support protocolaire final, API, commande, format, modèle ML, écran détaillé ou code n’est revendiqué.

## 25. Critères d’acceptation
### 1. Scénario
**Given** une capture liée à un Case avec pertes et troncatures
**When** l’analyste ouvre Network Intake
**Then** l’état reste `partial-capture`, les limites sont visibles et les analyses héritent de la partialité.

### 2. Scénario
**Given** une capture restreinte et un utilisateur sans droit payload
**When** l’intake est affiché
**Then** les métadonnées autorisées restent visibles, le contenu est masqué et aucune permission n’est accordée.

### 3. Scénario
**Given** aucun fournisseur de modèle et une capture exploitable
**When** l’analyste prépare l’intake
**Then** formulaires, viewers, règles et Tools déterministes suffisent.

## 26. Questions ouvertes
OPEN-005 conserve le choix des moteurs ; OPEN-008 le support Endpoint/capteur ; OPEN-013 la gouvernance de classe 2 ; OPEN-014 l’ambiguïté Artifact/Attachment lorsque pertinente ; OPEN-015 le bridge Automation Run/Response Run. Les schémas, cardinalités, permissions atomiques, écrans et contrats techniques sont reportés. OPEN-011/012 restent ouvertes et hors périmètre.

## 27. Consommateurs documentaires
Network Forensics Workbench, Case Workspace, Evidence Board, Event Search, Dynamic Sandbox, Memory/Disk Forensics, Static/Reverse, Entity Graph, future Detection Engineering et Intelligence, ainsi que les phases Objets, Permissions, Écrans, Journeys, Technique et Validation.
