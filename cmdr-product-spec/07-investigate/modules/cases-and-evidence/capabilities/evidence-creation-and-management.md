---
id: CAP-INV-107
title: Evidence Creation and Management
product: investigate
module: cases-and-evidence
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-PROD-020
open_decisions:
  - OPEN-013
  - OPEN-014
---
# CAP-INV-107 — Evidence Creation and Management

## 1. Définition
Créer et gérer une Evidence qualifiée depuis une source autorisée tout en conservant source, transformations, versions et contestations.

## 2. Problème utilisateur
Des données ou Artifacts sont parfois présentés comme preuves sans contexte, auteur ou transformation. Une mutation silencieuse détruit la confiance.

## 3. Objectifs
Créer Evidence avec origine/acquisition, Case, auteur/moteur et transformations ; lier Hypotheses/Findings ; versionner, contester, retirer d’usage et exporter.

## 4. Non-objectifs
Ne pas définir cryptographie/storage ; ne pas créer Evidence automatiquement ; ne pas confondre interprétation/source ; ne pas décider l’admissibilité juridique.

## 5. Propriétaire
Investigate / Cases and Evidence / Investigate Product Lead. Evidence est un objet canonique Investigate.

## 6. Utilisateurs
Principal : Evidence Reviewer. Secondaires : Case Analyst, Investigation Lead et consommateurs Finding/Action Request/Reporting.

## 7. Conditions d’entrée
Case accessible, source autorisée/résoluble, provenance suffisante, auteur/moteur identifié et permission Evidence create.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Source element | Artifact/Event/result/document | basis | oui | disponibilité/source time | candidate seulement |
| Origin/acquisition | producer/import/collection | provenance | oui | timestamp/version | block qualification |
| Case context | Investigate | scope | oui | permission courante | reject cross-tenant |
| Qualification reason | reviewer | contexte interprétatif | oui | versionnée | rester candidate |
| Transformations | tools/workflow | lineage | non | producer/version | none/unknown explicite |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Artifact | Investigate | source, versions, derivatives | lire/lier |
| Telemetry Event/Search result | Shared Capabilities | faits sources | lire/référencer |
| Case/Hypothesis | Investigate | scope/question | lire/lier |
| Finding | Investigate | claims | lire/lier |
| Automation Run/Tool Call | Studio | provenance producer | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Evidence | créer, versionner, supersede, retirer usage | Investigate | source conservée, aucune mutation silencieuse |
| Evidence relations | lier Case/Hypothesis/Finding | Investigate | typées et sourcées |
| Qualification event | créer | Timeline/Trace | reviewer/raison/version |
| Source object | aucune mutation | owner source | référence uniquement |

## 11. Fonctionnalités
Créer depuis source, enregistrer origine/acquisition/auteur/transformations, lier Case/Hypothesis/Finding, afficher intégrité conceptuelle/versions/accès/contestations, retirer sans supprimer et exporter.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Créer Evidence | Reviewer | Evidence | 2 | provenance complète | draft/qualified | OPEN-013 |
| Lier Hypothesis/Finding | Analyst | relation | 2 | objets accessibles | relation | OPEN-013 |
| Versionner/supersede | Reviewer | Evidence | 2 | raison/nouvelle source | version | OPEN-013 |
| Retirer usage | Lead | Evidence | 2 | raison | inactive avec trace | OPEN-013 |
| Exporter | Reviewer | Export Job | 1 | permission/redaction | package | selon données |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Capture provenance | oui | oui | oui | non | Audit Hooks |
| Suggestion candidate | oui | oui | oui | oui | sélection manuelle/règles |
| Lineage transformation | oui | oui | oui | explication | records producer |
| Qualification Evidence | reviewer | completeness rules | workflow revue | proposition seulement | décision humaine |

## 14. États fonctionnels
`candidate`, `draft`, `qualified`, `restricted`, `contested`, `superseded`, `withdrawn`, `unavailable`. États Draft.

## 15. États d’interface
Loading conserve source ; Partial nomme provenance ; Error garde candidate ; Offline lecture ; Permission denied sans fuite ; Stale expose versions/source.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Evidence | objet | Hypothesis/Finding/Action Request | source, auteur, version, statut visibles |
| Evidence relation | typed link | Case | auditée sans transfert |
| Qualification event | Timeline/Trace | Review/Audit | raison/reviewer |
| Evidence export | Export package | Govern/Reporting | permission-aware, cité, redacted |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Artifact/Event/result | qualifier | Evidence Creation | source, Case, origin, transformations, reason | source |
| Evidence | lier | Hypothesis/Finding | role, relevance, reviewer | Evidence |
| Evidence | contester/supersede | Evidence Review | reason, alternative, version | prior/new version |
| Evidence set | soumettre | Action Request Preparation | refs, status, restrictions | Case |

## 18. Dépendances
Artifact, Case, Hypothesis, Finding, Shared Trace/Versioning/Export/Object Linking, future collection 4B.2, permissions et OPEN-013/014.

## 19. Source de vérité
Evidence/qualification restent Investigate ; la source conserve son owner ; transformations conservent leur producer ; storage/intégrité futurs ; Export Shared.

## 20. Provenance et audit
Source IDs/versions, origine/acquisition, auteur/engine, transformations/Tool Calls, qualification/contest/supersede et access/export.

## 21. Permissions fonctionnelles
Evidence create/read/qualify, sensitive source, export/redaction, contest/withdraw, cross-tenant interdit et step-up pour Evidence restreinte.

## 22. Limites et erreurs
Source inaccessible, provenance incomplete, transformation unknown, duplicate candidate, Case mismatch, conflit de version, permission denied ou export restricted.

## 23. Métriques
Evidence avec provenance complète, candidates qualifiées/rejetées, contestations, silent mutation cible zéro et broken links.

## 24. Classification de livraison
`defined` / `planned`, cible native. Promotion conditionnée par modèle objet/version, provenance, access et futurs contrats d’intégrité/storage.

## 25. Critères d’acceptation
**Given** un Artifact sourcé **When** qualifié **Then** Artifact reste, Evidence distincte, auteur/raison/transformations visibles.

**Given** une candidate automatisée **When** affichée **Then** statut candidate et accept/reject, sans qualification automatique.

**Given** une interprétation corrigée **When** superseded **Then** ancienne version reste et les Findings sont notifiés.

## 26. Questions ouvertes
Algorithmes d’intégrité/storage, qualification levels et OPEN-014 pour Attachment restent ouverts.

## 27. Consommateurs documentaires
Hypothesis, Finding, Action Request, Reporting, Case Replay et phases Objets/Permissions.
