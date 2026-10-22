---
id: CAP-INV-372
title: User Activity and Application Artifact Analysis
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
  - REQ-SEC-002
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-372 — User Activity and Application Artifact Analysis

## 1. Définition
Examiner des artefacts persistants de profils, applications, documents récents, navigation, téléchargements, caches, communications et sessions comme observations soumises à minimisation, sans déduire automatiquement intention, auteur ou action humaine.

## 2. Problème utilisateur
Un compte, profil ou historique associé peut provenir d’une automatisation, synchronisation, autre utilisateur ou donnée ancienne. Une attribution automatique viole rigueur et confidentialité.

## 3. Objectifs
Show user profiles, persistent application activity, recent documents/open-history candidates, browsing/use/download candidates, caches, application history, locally stored communication traces when authorized, stored sessions/connections, timestamps, identities, conflicts and missing data; annotate/correlate with privacy controls.

## 4. Non-objectifs
Ne pas affirmer intention/auteur/action humaine, révéler données privées sans permission, devenir browser/network/cloud/mobile forensics complete, exécuter contenu, définir parser/moteur/format ou auto-create Finding.

## 5. Propriétaire
Investigate owns observations/interpre­tation; Security/Privacy controls access/minimization; application/user identity sources remain contextual; Studio Tools; Shared Timeline/Linking.

## 6. Utilisateurs
Principal : User Activity Analyst. Secondaires : DFIR Analyst, Investigation Lead, Privacy Reviewer, Evidence Reviewer and Audit Analyst.

## 7. Conditions d’entrée
Case purpose and authorization; image/session/filesystem accessible; relevant artifacts and age/source visible; sensitive content/path permissions; minimization applied.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| User/application artifacts | filesystem/Tool Calls | persistent observations | oui | image version | unavailable |
| Profile/account context | artifacts/identity projection | association candidate | non | source time visible | unattributed |
| Timestamp/session/download/cache data | source parsers | activity context | non | per record | partial |
| Privacy/purpose restrictions | Security/Case | minimization/access | oui | current | restricted |
| Correlation context | CAP-INV-370/374/377 | timeline/comparison | non | versions visible | local-only |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Disk Image/File Observation | Investigate | source/artifact | lire |
| User Activity/Application Artifact | Investigate concepts | persistent record/context | lire |
| Principal/Entity | Settings/Shared | identity context only | projection/lier |
| Tool/Tool Call | Studio | producer/version | lire |
| Timeline/Case/Hypothesis | Shared/Investigate | correlation/reasoning | lire/lier |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| User Activity Observation | create/annotate/dispute | Investigate concept | artifact ≠ intent |
| Application Artifact relation | create/supersede | Investigate | source/time/limits required |
| Identity context link | create/dispute | Shared/Investigate | associated account ≠ author |
| Evidence candidate | prepare | Investigate | qualification future |

## 11. Fonctionnalités
Inspect profiles, recent/open candidates, browsing/use/download/cache/application/communication/session traces when permitted; show timestamps/identities/conflicts/missing data; filter/compare/annotate/correlate; mask/minimize private data; prepare Hypothesis/Evidence without attribution certainty.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspect/filter | Analyst | artifacts | 0 | authorized read | minimized view | non |
| Compare/correlate | Analyst | relation | 0/1 | sources permitted | conflicts/limits | non |
| Annotate/dispute/link | Analyst | observation | 2 | permission | versioned context | OPEN-013 |
| Prepare candidate/handoff | Reviewer | package | 2 | provenance | CAP-INV-379 | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| parse/group artifacts | oui | parsers/rules | oui | suggestion | tables/categories |
| correlate timestamps/accounts | oui | oui | oui | suggestion | explicit filters/links |
| summarize activity | oui | source aggregation | oui | oui | timeline/table |
| infer intent/author | humain corroboré | non | non | no automatic conclusion | evidence review |

No sensitive data to model without authorization; no automatic attribution.

## 14. États fonctionnels
`available`, `partial`, `unattributed`, `conflicting`, `restricted`, `masked`, `unsupported`, `disputed`, `superseded`.

## 15. États d’interface
Loading preserves purpose; Empty does not mean no activity; Partial exposes missing data; Error keeps valid records; Offline read-only; Permission denied masks private content; Stale shows artifact age.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| User/Application Observation | observation | Case/Hypothesis | artifact, profile context and uncertainty distinct |
| Correlation relation | relation | CAP-INV-374/377 | source/time/privacy visible |
| Evidence/handoff selection | candidate package | CAP-INV-379 | no intention or authorship certainty |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| User/application artifact | inspect | CAP-INV-372 | source, profile context, timestamp, restrictions | Navigation |
| Observation | correlate | CAP-INV-374 | source/time/identity/conflicts | User view |
| Observation | handoff | CAP-INV-379 | selected records, minimization, uncertainty | User view |

## 18. Dépendances
CAP-INV-367/368/370/374/377/379, Security/Privacy, Studio Tools, Shared Timeline/Entity/Linking, OPEN-005/008/013/014/015.

## 19. Source de vérité
Persistent artifacts remain source records; Investigate owns interpretation; principal/profile is context, not automatic author; private data access remains Security-owned.

## 20. Provenance et audit
Case purpose, image/session/filesystem/file, Tool/version, artifact category, profile/account context, timestamps, masking/access, filters, correlations, annotations, exports/handoffs and disposition.

## 21. Permissions fonctionnelles
User activity read, application artifact read, sensitive/private content read, communication trace read, compare/correlate, annotate/link and candidate prepare. Purpose limitation, step-up and SoD future.

## 22. Limites et erreurs
Artifact ≠ intention; profile/account ≠ author; persistent data ≠ current state. Missing/corrupt/unsupported/private data, timezone conflicts, sync/automation ambiguity and permission denial remain explicit.

## 23. Métriques
Artifacts by category/state, masked/restricted views, unattributed/conflicting relations, privacy blocks and handoffs with explicit uncertainty.

## 24. Classification de livraison
`defined` / `planned`; no browser/cloud/mobile/network complete scope or implementation.

## 25. Critères d’acceptation
**Given** application history and associated profile but no corroboration **When** observation prepared **Then** artifact and profile are visible as context, not intent or certain author.

**Given** private communication artifact and insufficient permission **When** opened **Then** content stays masked/denied while safe metadata may remain visible.

**Given** no AI **When** analyzed **Then** parsers, filters, timeline, comparisons and human review work.

## 26. Questions ouvertes
OPEN-005/008/013/014/015 remain open; privacy, objects and permissions final are future.

## 27. Consommateurs documentaires
INV-DSK-001, Case/Hypothesis/Evidence, CAP-INV-374/377/379 and Objects/Permissions/Trust/Technique phases.
