---
id: CAP-INV-537
title: Intelligence Lifecycle Provenance, Closure and Continuous Improvement
product: investigate
module: threat-intelligence
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-06
requirement_ids: [REQ-PROD-014, REQ-PROD-019, REQ-PROD-020, REQ-PROD-055, REQ-INV-006, REQ-AI-002, REQ-SEC-001, REQ-SEC-002, REQ-UX-010]
open_decisions: [OPEN-013, OPEN-014, OPEN-015, OPEN-018, OPEN-019]
source-of-truth: canonical
---
# CAP-INV-537 — Intelligence Lifecycle Provenance, Closure and Continuous Improvement

## 1. Définition
Retracer le lifecycle Intelligence complet, vérifier sa clôture déclarée et préparer un Continuous Improvement Package vers Requirement, Knowledge Project, Case, Hunt, Detection Engineering, Settings, Studio ou reporting sans appliquer de modification active.

## 2. Problème utilisateur
La fin d’un cycle peut perdre les sources, décisions, corrections, consommateurs ou erreurs et transformer un package d’amélioration en action automatique.

## 3. Objectifs
- reconstruire la lineage de bout en bout.
- exposer gaps, restrictions, erreurs et décisions humaines.
- préparer des handoffs d’amélioration sans transfert d’ownership ni exécution.

## 4. Non-objectifs
Aucune API, protocole, format, provider, schéma final, action runtime, collecte, règle, blocage, réponse, partage externe, Cloud/Mobile ou écran détaillé.

## 5. Propriétaire
Investigate possède la provenance analytique et le Continuous Improvement Package. Chaque produit conserve ses objets. Shared possède Trace, Activity, Linking, Versioning, Notifications et Recovery. Govern conserve Decision/Approval/Response Run/Result. Studio conserve Tools/Runs.

## 6. Utilisateurs
Principal : **Intelligence Manager**. Secondaires : Threat Intelligence Analyst, Investigation Lead, Detection Engineer, SOC Analyst, Settings Admin, Studio Owner, Reviewer et Auditor autorisés.

## 7. Conditions d’entrée
Scope de clôture, origins, Requirements, Projects, versions, sources, assessments, products, publications, consumers, corrections, permissions, restrictions, owner et return origin sont explicites ; sinon partial/blocked.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Foundation and analysis lineage | CAP-INV-501..536 | complete selected lifecycle | oui | immutable/current refs | closure incomplete |
| Cross-product origins and outcomes | Case/Hunt/Command/Detection/Workbench | source and consumer context | selon scope | resolvable versions | gap visible |
| Tool, authority and trace context | Studio / Govern / Shared | runs, decisions and audit | selon scope | current links | provenance partial |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Threat Intelligence concepts | Investigate | complete lifecycle projection | lecture/lien |
| Case/Hunt/Incident/Detection/Finding/Evidence | respective owners | origin/outcome context | lecture limitée |
| Tool Call/Automation Run/Decision/Approval/Trace | Studio/Govern/Shared | contributor and authority lineage | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Intelligence Lifecycle Provenance Assessment | créer, compléter, contester, superseder | Investigate concept | missing links explicit |
| Continuous Improvement Package | créer, versionner, retirer, superseder | Investigate concept | package ≠ active change |
| Closure disposition | créer/revoir | Investigate concept | closure declared for scope only |

## 11. Fonctionnalités
Retracer origins, Requirements, Projects, sources, materials, candidates, Sightings, relations, assessments, Sessions, hypotheses, fusion, attribution, products, reviews, publications, watchlists, handoffs, monitoring, feedback, satisfaction, corrections, Runs, Decisions, errors et human dispositions ; préparer amélioration.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter/naviguer/comparer | Intelligence Manager | lifecycle/trace | 0 | lecture autorisée | vue sourcée | non |
| reconstruire lineage/gaps | Intelligence Manager | assessment/Tool Call | 1 | scope explicite | résultat attribué | selon politique |
| créer closure/package/handoff | Intelligence Manager | concepts locaux | 2 | owner/provenance explicites | version réversible | OPEN-013 |
| appliquer amélioration/action | aucun rôle local | destination object/runtime | 3 | destination owner/Govern | aucune exécution | obligatoire |
| supprimer trace/historique | aucun rôle local | provenance | 4 | interdit | refus audité | strict |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| reconstruire lineage | oui | linking/trace/versioning | oui | résumé sourcé | tables/timeline |
| détecter gaps/incohérences | oui | règles/diff | oui | suggestion incertaine | checklist/audit |
| préparer improvement package | oui | template versionné | oui | draft attribué | formulaire/workflow |
| appliquer ou approuver changement | destination/Govern humain | contrôles | jamais autonome | jamais décisionnaire | handoff humain |

Toute automatisation expose initiateur, moteur/version, Tool Calls, Run, sources, paramètres, erreurs, incertitude et disposition humaine.

## 14. États fonctionnels
`draft`, `collecting-lineage`, `partial`, `complete-for-declared-scope`, `reviewed`, `closure-ready`, `closed-for-scope`, `improvement-ready`, `returned`, `withdrawn`, `superseded`.

## 15. États d’interface
Loading conserve contexte ; Empty distingue aucune lineage/interdiction ; Partial nomme liens manquants ; Error conserve valide ; Offline stale/read-only ; Permission denied ne révèle rien ; Stale expose versions ; Conflict offre diff/recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Lifecycle Provenance Assessment | assessment versionné | reviewers/auditors | lineage, gaps et restrictions visibles |
| Closure disposition | business event | Threat Intelligence owner / Quality | scope et réserves explicites |
| Continuous Improvement Package | package versionné | Intake/Requirement/Case/Hunt/Detection/Settings/Studio | aucune mutation active |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-501..536 | lifecycle event or closure request | CAP-INV-537 | IDs, versions, sources, decisions, errors | source concept |
| CAP-INV-537 | improvement need approved for handoff | Intake/Requirement/Case/Hunt/Detection/Settings/Studio | gap, evidence, restrictions, return origin | Package |
| destination owner | disposition returned | CAP-INV-537 | accepted/rejected scope, errors, status | Handoff |

Chaque transition conserve owner, tenant, versions, permissions, restrictions, erreurs, autorité, provenance et return origin.

## 18. Dépendances
CAP-INV-501..536 ; Command ; Cases/Hunts/Workbench ; Detection Engineering ; Settings ; Studio ; Govern ; Shared ; OPEN-013/014/015/018/019.

## 19. Source de vérité
Investigate est source de la provenance/closure Intelligence ; chaque objet canonique et destination reste chez son owner. Le package ne remplace jamais la source et ne crée aucun droit.

## 20. Provenance et audit
Conserver toutes origins, versions, sources, restrictions, permissions, Tools, Runs, assessments, products, reviews, publications, consumers, feedback, corrections, Decisions, errors, timestamps, owners et dispositions humaines. Aucune trace supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| lifecycle provenance read/export | cross-product disclosure | 0/1 | restricted lineage masked | possible | viewer/auditor | owners/Shared | Permissions |
| closure/package create-review | premature closure or unintended influence | 2 | scope minimized | OPEN-013 | author/reviewer | Investigate | Permissions |

## 22. Limites et erreurs
Closure ≠ truth achieved ; package ≠ action ; destination acceptance ≠ implementation ; missing links, restricted data, stale versions, partial outcomes et unresolved OPEN restent visibles.

## 23. Métriques conceptuelles
Cycles par closure state ; lineage completeness ; packages/dispositions ; active changes executed locally — cible zéro ; trace supprimée — cible zéro.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire uniquement. Aucun statut implemented/operational n’est revendiqué.

## 25. Critères d’acceptation
### 1. Provenance partielle
**Given** un Tool Call ou consumer link manquant  
**When** la clôture est évaluée  
**Then** état partial, gap explicite et aucune clôture complète inventée.

### 2. Improvement Package
**Given** un gap validé  
**When** le package est transmis  
**Then** destination, scope, restrictions et return origin sont conservés sans modification active.

### 3. Sans IA
**Given** aucun modèle  
**When** le lifecycle est retracé  
**Then** linking, trace, versioning, tables, timelines et revue humaine suffisent.

## 26. Questions ouvertes
OPEN-013/014/015/018/019 restent ouvertes ; objets finaux, permissions atomiques, contrats techniques et implémentation restent futurs.

## 27. Consommateurs documentaires
Threat Intelligence, Investigate, Command, Cases/Hunts, Analysis Workbench, Detection Engineering, Settings, Studio, Govern, Shared, Quality, Technique et Roadmap. Cloud/Mobile restent différés.
