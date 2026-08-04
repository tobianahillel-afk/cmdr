---
id: CAP-INV-008
title: Search and Hunt Provenance
product: investigate
module: signals-and-hunt
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-AI-002
open_decisions:
  - OPEN-015
---
# CAP-INV-008 — Search and Hunt Provenance

## 1. Définition
Retracer Queries, runs, paramètres, sources, transformations, erreurs, résultats, assistance et liens Case d’une recherche ou d’un Hunt.

## 2. Problème utilisateur
Une recherche non reproductible ou non attribuée ne peut pas soutenir une investigation fiable. Trace et Activity Stream ont besoin de la sémantique métier Investigate.

## 3. Objectifs
Définir les événements métier de provenance, conserver Query/version/auteur/paramètres/période/sources/run/erreurs/transformations/exports et exposer toute contribution automatisée.

## 4. Non-objectifs
Ne pas recréer Trace ou Activity Stream ; ne pas définir stockage/rétention ; ne pas transformer une trace en Evidence.

## 5. Propriétaire
Investigate définit la sémantique ; Shared Capabilities possède Trace, Activity Stream, Audit Hooks et Timeline Engine.

## 6. Utilisateurs
Principal : Threat Hunter. Secondaires : Case Reviewer, Auditor, Report Author et Quality reviewer.

## 7. Conditions d’entrée
Query ou Hunt identifiable, correlation IDs, sources/périodes connues ou lacunes explicites et permissions d’audit.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Query/version | Shared Query | origine logique | oui | version immuable | provenance incomplete |
| Search Job | Shared Search | exécution | oui pour run | timestamps/statut | conserver Query seule |
| Sources/paramètres | Data Sources/user | scope | oui | snapshot au lancement | marquer inconnus |
| Automation metadata | Studio/workflow | proposition/transformation | non | run/version/Tool Calls | indiquer aucune automation |
| Case/Hunt links | Investigate | contexte métier | non | état courant | standalone |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Query | Shared Capabilities | version et paramètres | consulter/référencer |
| Search Job | Shared Capabilities | étapes, erreurs, résultats | consulter/corréler |
| Automation Run | CMDR Studio | agent, Tool Calls, outputs | consulter sans posséder |
| Case/Hunt workspace | Investigate | relations et dispositions | consulter |
| Timeline Entry | Shared Capabilities | événements ordonnés | consulter/filtrer |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Provenance record | émettre événements métier | Shared Capabilities | immuable ou superseded |
| Trace link | lier Query/run/result/Case | Shared Capabilities | sourcé et permission-aware |
| Hunt activity | ajouter événements | Investigate via mécanisme Shared | sémantique locale |
| Evidence | aucune création | Investigate | trace ≠ Evidence |

## 11. Fonctionnalités
Capturer Query/version/auteur/paramètres/période/sources, Search Job/erreurs/limites, transformations, résultats, annotations, exclusions, exports, liens Hunt/Case/Incident et agent/workflow/Tool Calls.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter provenance | Analyst/Auditor | Trace | 0 | audit read | chaîne affichée | non |
| Comparer runs | Threat Hunter | Search Jobs | 0 | runs accessibles | diff | non |
| Ajouter justification | Analyst | activity event | 2 | objet accessible | annotation auditée | OPEN-013 |
| Exporter trace | Auditor | Export Job | 1 | permission/redaction | package | selon sensibilité |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Capture événements | non | oui | oui | non | Audit Hooks |
| Corrélation | oui | oui | oui | explication | correlation IDs/Object Linking |
| Résumé trace | oui | oui | oui | oui | timeline filtrée |
| Détection lacunes | oui | oui | oui | oui | règles de complétude |

## 14. États fonctionnels
`complete`, `partial`, `broken-link`, `redacted`, `superseded`, `expired-source`, `disputed`. États Draft.

## 15. États d’interface
Loading indique la source ; Partial nomme les éléments absents ; Error garde les IDs ; Offline affiche le cache stale ; Permission denied conserve les relations sans fuite ; Redacted garde la raison.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Search provenance chain | Trace relations | Analyst/Auditor | corrélée, versionnée, permission-aware |
| Hunt activity timeline | Timeline entries | Hunt/Case Replay | ordonnée et sourcée |
| Completeness event | quality event | owner | lacunes nommées sans conclusion |
| Trace export | Export package | Audit/Reporting | redacted et référencé |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Query Authoring | exécution | Trace | Query/version, auteur, paramètres, scope | Search Job |
| Search Job | fin/erreur | Trace | statut, sources, limites, résultats | Event Search |
| Hunt | promotion Case | Case Trace | Queries, results, Hypotheses, disposition | Hunt |
| Automation Run | proposition utilisée | Investigate Trace | agent/version/Tool Calls/sources/disposition | objet proposé |

## 18. Dépendances
Shared Trace, Activity Stream, Timeline, Audit Hooks, Object Linking, CAP-INV-002/005/112, Studio Automation Run, OPEN-015 et policies de rétention/redaction.

## 19. Source de vérité
Query/Search Job/Timeline restent Shared ; Hunt/Case dispositions restent Investigate ; Automation Run reste Studio ; Trace est un mécanisme Shared.

## 20. Provenance et audit
Ce contrat impose source, acteur/producteur, version, run, paramètres, transformations, disposition, correction par supersession et redaction avec raison.

## 21. Permissions fonctionnelles
Trace read, audit export, sensitive source visibility, Automation Run projection read et cross-tenant audit interdit par défaut.

## 22. Limites et erreurs
Correlation ID absent, run expiré, source inaccessible, trace partielle, redaction excessive, clock skew, permission mismatch ou duplicate events.

## 23. Métriques
Runs avec provenance complète, broken links, propositions agentiques avec disposition, exports redacted et corrections superseded.

## 24. Classification de livraison
`defined` / `planned`, cible native. Promotion conditionnée par hooks déterministes, IDs stables, redaction et tests de liens.

## 25. Critères d’acceptation
**Given** une Query exécutée et un Case lié **When** la provenance est ouverte **Then** version, auteur, période, sources, erreurs et transformations sont visibles.

**Given** une Query proposée puis modifiée **When** le run est inspecté **Then** agent/version/Tool Calls, diff et disposition humaine sont visibles.

**Given** un résultat expiré **When** la trace est ouverte **Then** le lien cassé est explicite et aucune donnée n’est inventée.

## 26. Questions ouvertes
Rétention et stockage d’audit sont reportés ; OPEN-015 reste ouverte ; la redaction des Tool Calls doit être précisée.

## 27. Consommateurs documentaires
Signals and Hunt, Case Replay, Evidence/Finding review, Reporting Preparation, Quality et Audit.
