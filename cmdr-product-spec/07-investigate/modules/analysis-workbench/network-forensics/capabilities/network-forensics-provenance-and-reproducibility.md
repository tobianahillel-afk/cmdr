---
id: CAP-INV-396
title: Network Forensics Provenance and Reproducibility
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
  - OPEN-014
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-396 — Network Forensics Provenance and Reproducibility

## 1. Définition
Retracer et évaluer la reproductibilité d’une analyse Network Forensics depuis le Case, la source et l’acquisition jusqu’aux Tools, paramètres, observations, Derived Artifacts, décisions humaines et dispositions finales, sans redéfinir Trace, Activity ou Audit.

## 2. Problème utilisateur
Une conclusion réseau non reliée à la capture, au capteur, à la timebase, aux paramètres et aux décisions humaines ne peut pas être reproduite, contestée ou qualifiée de manière fiable.

## 3. Objectifs
- relier Case, source, Collection Request, Collection Job et Capture Artifact ;
- conserver acquisition context, custody, capteurs, interfaces, scope, coverage et timebase ;
- conserver Tools, versions, Tool Calls, Automation Runs, paramètres, filtres et requêtes ;
- relier flows, sessions, conversations, transactions, observations et Derived Artifacts ;
- enregistrer erreurs, interruptions, accès sensibles, décisions humaines et dispositions ;
- évaluer la reproductibilité sans masquer les éléments manquants.

## 4. Non-objectifs
Ne définit aucun moteur, protocole, format, algorithme d’intégrité, schéma objet final, API, commande, stockage, mécanisme d’acquisition, règle Detection, objet Intelligence canonique, écran détaillé ou implémentation. Ne remplace pas Shared Trace, Activity Stream, Audit Hooks, Versioning ou Reporting.

## 5. Propriétaire
Investigate / Analysis Workbench / Investigate Product Lead possède l’évaluation locale de reproductibilité, les relations analytiques et la disposition humaine. Studio possède Tool, Tool Call, Workflow et Automation Run. Settings possède capteurs administrés, Fleet, Policies, providers, storage, retention, health, time synchronization et secrets. Shared possède Trace, Activity, Audit Hooks, Versioning, Timeline, Export, Reporting et Recovery. Govern possède l’autorité sur les cibles réelles.

## 6. Utilisateurs
Principal : **Audit Analyst**. Secondaires : Network Forensics Analyst, Evidence Reviewer, Investigation Lead, Quality Reviewer.

## 7. Conditions d’entrée
Case et Network Forensics Session accessibles ; capture et version identifiables ; acquisition, custody et limites disponibles ou explicitement absentes ; Tools et paramètres attribués ; permissions réévaluées ; aucun accès à une cible réelle.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Case, Collection et Capture lineage | Investigate / Collection projections | chaîne source et contexte | oui | versions liées à la session | `missing-capture` ou `partially-reproducible` |
| Session, scope, coverage et timebase | CAP-INV-381..394 | contexte analytique versionné | oui | snapshot de session | état spécifique de lacune |
| Tool, version, Tool Calls et Automation Run | CMDR Studio | producteur, paramètres et statut | oui si traitement exécuté | version du traitement | `missing-tool` ou `missing-version` |
| Observations, Derived Artifacts et dispositions | Investigate | résultats et décisions | oui pour évaluation complète | versions conservées | `partially-reproducible` |
| Permissions et restrictions | Security / Settings | accès et minimisation | oui | décision courante et historique disponible | `policy-blocked` ou `restricted-payload` |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case / Hypothesis / Artifact | Investigate | objectif, source, restrictions, versions et return origin | consulter et lier |
| Collection Request / Collection Job / Capture Artifact | Investigate / producteur | acquisition, résultat, erreurs, custody et lineage | consulter uniquement |
| Network Forensics Session et observations | Investigate | scope, coverage, timebase, sélections, résultats et dispositions | consulter, comparer et lier |
| Tool / Tool Call / Workflow / Automation Run | CMDR Studio | version, paramètres, initiateur, statut et résultat | consulter uniquement |
| Trace / Activity / Audit / Version / Timeline | Shared Capabilities | événements, corrélations et versions | consommer sans redéfinir |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Reproducibility Assessment | créer, réévaluer, contester ou superseder | Investigate | sources, lacunes, auteur, version et conclusion obligatoires |
| Provenance relation locale | créer ou corriger réversiblement | Investigate | référence les sources propriétaires sans les modifier |
| Disposition de reproductibilité | accepter, contester ou retirer | Investigate | historique conservé ; aucune suppression de trace |
| Trace / Activity event | émettre | Shared Capabilities | append-only avec correlation ID |

## 11. Fonctionnalités
- afficher la chaîne Case → Collection → Capture → Session → analyse → handoff ;
- afficher capteurs, interfaces, scope, coverage, timebase et restrictions ;
- afficher Tools, versions, Tool Calls, Automation Runs, paramètres, filtres et requêtes ;
- afficher les relations entre packets, flows, sessions, conversations, transactions, observations et Derived Artifacts ;
- afficher accès sensibles, erreurs, interruptions, résultats partiels et décisions humaines ;
- identifier les éléments manquants ou non rejouables ;
- produire une évaluation attribuée et contestable ;
- fonctionner intégralement sans fournisseur de modèle.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter la provenance | Audit Analyst | Provenance projections | 0 | lecture autorisée | chaîne sourcée et lacunes visibles | non |
| Reproduire une analyse autorisée | Audit Analyst | Tool Call / Analysis Result | 1 | capture, Tool/version, paramètres et permission disponibles | nouveau résultat attribué et comparable | non, sauf policy |
| Créer ou réévaluer l’assessment | Audit Analyst | Reproducibility Assessment | 2 | sources et limites présentes | nouvelle version contestable | OPEN-013 |
| Contester ou superseder | Quality Reviewer | Reproducibility Assessment | 2 | justification et permission | disposition versionnée | OPEN-013 |
| Exporter une provenance autorisée | Audit Analyst | provenance package | 1 | classification et export permission | export sourcé et audité | selon policy |

Les classes 3 et 4 ne sont jamais exécutées. Toute interaction réelle est bloquée ou routée vers les capabilities propriétaires et Govern.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Collecter les références de provenance | oui | oui | oui | résumé facultatif | tables, Trace et relations explicites |
| Vérifier les préconditions de reproduction | oui | validateurs explicables | oui | explication facultative | checklist et diagnostics déterministes |
| Comparer résultat initial et reproduction | oui | comparateur | oui | synthèse attribuée | diff, tables et inspection manuelle |
| Proposer un état de reproductibilité | oui | règles explicables | oui | suggestion avec incertitude | règles et revue humaine |
| Confirmer Evidence, Finding, règle ou Intelligence | owner humain | contrôles seulement | workflow de revue | jamais autonome | capabilities propriétaires |

Toute sortie automatisée expose initiateur, agent ou moteur, version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet.

## 14. États fonctionnels
`reproducible`, `partially-reproducible`, `not-reproducible`, `missing-capture`, `incomplete-capture`, `missing-sensor-context`, `missing-timebase`, `missing-tool`, `missing-version`, `unsupported-link-context`, `restricted-payload`, `policy-blocked`, `disputed`. Ces états sont fonctionnels et ne constituent pas une machine d’état objet définitive.

## 15. États d’interface
Loading conserve le Case, la session et la sélection ; Empty distingue absence de provenance et analyse non exécutée ; Partial nomme chaque lien manquant ; Error conserve les références valides ; Offline interdit une nouvelle reproduction ; Permission denied masque les valeurs sensibles sans masquer l’existence de la restriction ; Stale expose Tool, policy ou source obsolète. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Reproducibility Assessment | assessment Investigate | analyste, Evidence Reviewer et Quality | état, facteurs, lacunes, auteur et version visibles |
| Provenance package | projection sourcée | CAP-INV-397 / Export / Reporting | sources propriétaires et restrictions conservées |
| Reproduction comparison | Analysis Result / comparison | Network Workbench | paramètres, Tool/version et différences visibles |
| Missing-context event | quality event | Case / owner source / Settings ou Studio | lacune nommée sans invention de contexte |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-380..395 | demander provenance | CAP-INV-396 | Case, capture, session, observations, Tools, paramètres et restrictions | Workbench avec sélection préservée |
| CAP-INV-396 | reproduire un traitement | CMDR Studio Tool Call | source, version, paramètres, permissions et correlation ID | assessment et vue d’origine |
| CAP-INV-396 | comparer plusieurs résultats | CAP-INV-394 | résultats, versions, timebases, sources et limites | provenance avec return origin |
| CAP-INV-396 | préparer handoff | CAP-INV-397 | assessment, sources, erreurs, décisions et restrictions | provenance avec return origin |

Chaque transition conserve tenant, Case, ownership, permissions, restrictions, erreurs et return origin.

## 18. Dépendances
CAP-INV-105/107/108/109/203/208/212/213/214/301..313/329..346/347..362/363..379/380..397 ; CMDR Studio Tool/Tool Call/Workflow/Automation Run ; Platform Settings sensors/Fleet/Policies/storage/retention/health/timebase ; Shared Trace/Activity/Audit Hooks/Versioning/Timeline/Export/Reporting/Recovery ; OPEN-005/008/013/014/015.

## 19. Source de vérité
Investigate est source du Reproducibility Assessment et des relations analytiques locales. Les captures et résultats bruts restent chez leurs producteurs. Studio est source des Tools, Tool Calls et Automation Runs. Settings est source des projections administratives. Shared est source de Trace, Activity, Audit et Versioning. Aucun propriétaire n’est remplacé.

## 20. Provenance et audit
Conserver Case, source ou Endpoint, Collection Request, Collection Job, Capture Artifact, acquisition context, custody, Network Forensics Session, capteurs, interfaces, scope, coverage, timebase, Tools, versions, Tool Calls, Automation Runs, paramètres, filtres, requêtes, observations, Derived Artifacts, accès aux payloads sensibles, erreurs, interruptions, décisions humaines, dispositions et correlation IDs.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Provenance et assessment read | données réseau, utilisateurs et outils | 0 | valeurs sensibles masquées | possible | viewer/reviewer selon policy | Investigate / Shared | Permissions |
| Reproduction run | coût, volume et accès au brut | 1 | payload masqué selon permission | possible | initiateur distinct du reviewer si requis | Investigate / Studio | Permissions/Technique |
| Assessment create/dispute/supersede | mutation analytique | 2 | justification protégée si sensible | OPEN-013 | auteur/reviewer | Investigate | Permissions |
| Provenance export | fuite ou contexte sensible | 1 | minimisation et redaction | possible | auteur/export reviewer | Shared Export / Security | Permissions |

La matrice atomique, les namespaces, le RBAC/ABAC et les règles finales de step-up restent reportés.

## 22. Limites et erreurs
Capture absente, incomplète, corrompue ou restreinte ; contexte capteur ou timebase manquant ; Tool/version indisponible ; paramètres non conservés ; lien cross-source non supporté ; payload restreint ; policy block ; permission révoquée ; tenant mismatch ; résultat partiel ou contesté. Une analyse non reproductible reste consultable avec ses limites ; aucune donnée manquante n’est inventée.

## 23. Métriques
Part des sessions avec chaîne complète, assessments par état, liens manquants par type, reproductions réussies ou partielles, divergences expliquées, accès sensibles tracés, dispositions contestées et exports autorisés. Aucune cible chiffrée définitive n’est fixée.

## 24. Classification de livraison
`defined` / `planned`. La preuve est documentaire uniquement. Aucun moteur, produit tiers, protocole, API, commande, format final, écran détaillé ou code n’est revendiqué. Promotion conditionnée par objets, permissions, contrats Shared/Studio et implémentation future.

## 25. Critères d’acceptation
### 1. Capture incomplète
**Given** une capture partielle et des pertes connues
**When** l’analyste évalue la reproductibilité
**Then** l’état est `partially-reproducible` ou `incomplete-capture`, les limitations sont héritées et aucune complétude n’est affirmée

### 2. Tool ou version manquante
**Given** une analyse avec Tool Call mais version indisponible
**When** une reproduction est demandée
**Then** l’exécution n’est pas présentée comme équivalente, `missing-version` est visible et le résultat initial reste consultable

### 3. Sans IA
**Given** aucun fournisseur de modèle
**When** la provenance est inspectée et une reproduction autorisée est comparée
**Then** tables, Trace, validateurs, comparateur et revue humaine couvrent le workflow

### 4. Payload restreint
**Given** une analyse dont le payload source est restreint
**When** un reviewer sans permission ouvre l’assessment
**Then** l’existence et l’impact de la restriction sont visibles, le contenu reste masqué et aucun export n’est possible

## 26. Questions ouvertes
OPEN-005 conserve le choix des moteurs ; OPEN-008 le support capteurs et plateformes ; OPEN-013 la gouvernance classe 2 ; OPEN-014 la frontière Artifact/Attachment ; OPEN-015 les relations Automation Run/Response Run. OPEN-011 et OPEN-012 restent ouvertes et hors périmètre. Les schémas, permissions atomiques, écrans et contrats techniques restent futurs.

## 27. Consommateurs documentaires
Network Forensics Workbench, Case Workspace, Evidence Review, Finding Management, Static/Reverse/Memory/Disk correlation, CMDR Studio, Platform Settings, Shared Trace/Reporting/Export, futurs packages Detection Engineering et Intelligence, phases Objets, Permissions, Écrans, Journeys, Technique et Validation.
