---
id: CAP-INV-702
title: Mobile Investigation Session and Workspace Management
product: investigate
module: mobile-forensics
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-07
requirement_ids: [REQ-INV-001, REQ-PROD-014, REQ-PROD-020, REQ-OBJ-009, REQ-AI-002, REQ-SEC-001]
open_decisions: [OPEN-011, OPEN-013, OPEN-014, OPEN-015]
source-of-truth: canonical
---
# CAP-INV-702 — Mobile Investigation Session and Workspace Management

## 1. Définition
Gérer une Mobile Forensics Session durable qui conserve Case/origine, Mobile Evidence Packages, device context, plateforme candidate, objectif, scope, owner/contributeurs, Tools/Runs, paramètres fonctionnels, vues/filtres/requêtes, observations, Derived Artifacts, Hypotheses, erreurs, interruptions, handoffs, versions et supersession.

## 2. Problème utilisateur
Sans Session bornée, les analyses de plusieurs packages/backups, les résultats Tool, les hypothèses et les données privées peuvent perdre leur origine, leur version, leur scope ou leur retour au Case.

## 3. Objectifs
- créer/reprendre/pauser/compléter/archiver/supersede une Session;
- préserver scope, source versions, permissions, sélection, vues et return origin;
- relier observations, Derived Artifacts, Hypotheses et handoffs sans les confondre avec Evidence;
- rendre erreurs, interruptions et contributions automatisées/humaines traçables.

## 4. Non-objectifs
Ne pas créer Evidence, Finding, Collection Job, Tool/Automation Run lifecycle, final object schema, Mobile screen ID, acquisition ou real-device action.

## 5. Propriétaire
Investigate possède la Session et ses relations analytiques. Shared possède generic versioning/recovery/linking; Studio possède Tools/Runs; destination owners conservent leurs objets.

## 6. Utilisateurs
Principal : Investigation Lead. Secondaires : Mobile Forensics Analyst, DFIR Analyst, Evidence Reviewer, SOC Analyst et Sensitive Data Reviewer.

## 7. Conditions d’entrée
Intake existant, origine et source refs, scope/permissions connus ou limitations acceptées, owner humain et return origin. Une Session peut être créée `blocked` ou `partial` si les limites sont explicites.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Mobile Intake | CAP-INV-701 | objectif/scope/origine | oui | version sélectionnée | Session reste `draft` |
| Mobile Evidence Package(s) | source/Investigate | matériaux analysés | au moins une pour analyse | versionnées | `blocked`/`source-required` |
| Device/platform context | CAP-INV-703 | contexte candidate | non au create, requis pour certaines analyses | version courante | partial explicite |
| Owner/contributeurs/permissions | IAM/Security | responsabilité et accès | oui | revalidés à mutation | read-only/blocked |
| Tools/Runs/results | Studio | exécution attribuée | non | versionnés | aucune exécution inventée |
| Views/filters/queries | utilisateur/Shared | état workspace | non | dernière sauvegarde | defaults sûrs |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Intake / Case / Artifact | Investigate | objective/source/return | lire/lier |
| Package/backup/extraction concepts | source/Investigate | source material | lire selon permission |
| Tool / Tool Call / Automation Run | Studio | provenance/results | lire/lier |
| Hypothesis / Evidence / Finding | Investigate owner workflows | context/handoff status | lire/lier; pas qualifier localement |
| Shared version/link/trace | Shared | generic mechanisms | consommer |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Mobile Forensics Session | créer, modifier, pause, complete, archive, reopen, supersede | Investigate | scope/source/owner/version/return origin conservés |
| Session membership/view state | ajouter/retirer/versionner | Investigate/Shared projection | n’accorde aucune permission source |
| Observation/Hypothesis links | lier/délier/supersede | Investigate | relation sourcée, pas d’ownership transfer |
| Tool/Run | aucune redéfinition | Studio | référence seulement |

## 11. Fonctionnalités
Workspace durable avec breadcrumb origine, packages/sources, scope, device/platform context, collaborateurs, onglets analytiques conceptuels, views/filters/queries, bookmarks, activity, Tool results, errors, Derived Artifacts, hypotheses, handoffs, versions et restore/reopen contrôlé.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Créer Session | analyste autorisé | Session | 2 | Intake | draft/ready | OPEN-013 |
| Modifier scope/owner/view | owner/contributeur | Session | 2 | manage | version nouvelle | non |
| Pauser/reprendre | analyste | Session | 2 | état compatible | état tracé | non |
| Clore/archiver/reopen | owner/reviewer | Session | 2 | reason | historique préservé | non |
| Ouvrir source/Tool result | utilisateur | projection | 0 | source permission | navigation | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Restaurer workspace | oui | oui | oui | non | saved state |
| Résumer activité | oui | oui par événements | oui | résumé | activity log |
| Proposer regroupement d’observations | oui | règles | oui | suggestion | tags/filters |
| Détecter source stale/permission change | oui | oui | oui | explication | validation déterministe |
| Clore Session | humain | règles de complétude | workflow | jamais autonome | checklist/review |

## 14. États fonctionnels
`draft`, `ready`, `active`, `paused`, `blocked`, `partial`, `completed`, `failed`, `archived`, `superseded`. `completed` ne signifie pas que l’extraction était complète.

## 15. États d’interface
Loading restaure le dernier contexte; Empty montre la source manquante; Partial liste packages/zones manquants; Error garde les observations valides; Offline limite aux données disponibles; Permission denied masque source; Stale identifie la version à revalider.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Mobile Forensics Session | concept | CAP-INV-703..719 | scope, owners, sources, version et return origin |
| Session activity | événement | Trace/Timeline/QA | acteur/action/result/reason |
| Workspace state | projection | utilisateur | filtres/vues sans élargir permissions |
| Handoff references | relations | destination owners | package/draft uniquement |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-701 | ouvrir workspace | CAP-INV-702 | Intake, source refs, scope, return origin | Intake/origine |
| Session active | contexte device | CAP-INV-703 | Session/version, source, scope | Session |
| Session + source | analyser | CAP-INV-704..717 | selected refs, permissions, filters | Session |
| Session result | préparer handoff | CAP-INV-718 | observations, artifacts, restrictions | Session |
| Session close/review | audit | CAP-INV-719 | versions, activity, Tools/Runs, errors | Session/origine |

## 18. Dépendances
CAP-INV-701/703..719, Shared Linking/Versioning/Recovery/Trace/Collaboration, Studio Tools/Runs, Security permissions, OPEN-011/013/014/015.

## 19. Source de vérité
Session state and local analytical relations: Investigate. Generic mechanisms: Shared. Tools/Runs: Studio. Source/package authority and destination objects remain their owners.

## 20. Provenance et audit
Session ID/version, Intake/origine, owner/contributeurs, package versions, device/platform context, scope, permissions, views/filters/queries, Tools/Calls/Runs/parameters, observations, Derived Artifacts, hypotheses, errors, interruptions, handoffs, state changes and human reasons.

## 21. Permissions fonctionnelles
Session create/read/update/close/reopen/archive, membership update, source link, view/filter/query save, Tool result read, sensitive-source read separately, cross-tenant relation and provenance read/export. No Session permission grants raw source access.

## 22. Limites et erreurs
Missing source, revoked access, concurrent version, stale package, Tool failure, cross-tenant denial, invalid state transition or recovery failure preserves valid history and marks partial/blocked rather than dropping work.

## 23. Métriques
Sessions by state, reopen/supersession rates, stale-source blocks, permission changes, lost-context errors, return-origin success, Tool failures and time from ready to completed.

## 24. Classification de livraison
`defined` / `planned`. No final object schema, UI screen, storage, collaboration engine or runtime is delivered.

## 25. Critères d’acceptation
**Given** deux versions d’un backup **When** l’analyste change de source **Then** la Session conserve la version sélectionnée et les observations antérieures restent attribuées à leur source.

**Given** une permission source révoquée **When** la Session est reprise **Then** le contenu protégé est masqué/bloqué sans perdre l’historique d’analyse autorisé.

**Given** aucun modèle IA **When** la Session est gérée **Then** views, tables, filters, activity, checklists et human review couvrent le workflow.

## 26. Questions ouvertes
Final Session schema/state machine, atomic permissions, Tool/Run relations and Mobile delivery remain Objects/Permissions/Technique; OPEN-011/013/014/015 remain open.

## 27. Consommateurs documentaires
All Mobile capabilities, Case Workspace, Evidence, Analysis Workbench, Studio, Shared, Security, Objects, Screens, Quality and Roadmap.
