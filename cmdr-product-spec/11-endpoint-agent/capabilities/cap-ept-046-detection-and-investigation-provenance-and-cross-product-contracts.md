---
id: CAP-EPT-046
title: Detection and Investigation Provenance and Cross-Product Contracts
product: endpoint-agent
module: investigation
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-018, REQ-PROD-019, REQ-OBJ-008, REQ-INV-006, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008, OPEN-015, OPEN-017]
source-of-truth: canonical
---
# CAP-EPT-046 — Detection and Investigation Provenance and Cross-Product Contracts

## 1. Définition
Définir la chaîne de provenance EPT-3 et ses contrats de handoff : EPT-2 Observation → Detection Content/version → eligibility → local evaluation → match/local candidate → context → related observations → local timeline/investigation → summary → consumer, sans transfert d’ownership.

## 2. Problème utilisateur
Une investigation locale n’est fiable que si chaque résultat peut être retracé jusqu’aux observations et versions qui l’ont produit et si les consumers savent quelles parties restent Endpoint, Shared, Investigate, Command, Settings, Studio ou Govern.

## 3. Objectifs
Conserver IDs/versions/times/reasons/transformations ; distinguer observed/derived/AI-proposed ; préserver tenant/permissions/masking ; documenter cross-product owners et return origin ; garder local audit distinct de Shared Trace.

## 4. Non-objectifs
Aucun immutable-ledger implementation, cryptographic protocol final, physical schema, cross-product API, ownership transfer, Evidence custody claim, Response Run/Tool Call bridge final.

## 5. Propriétaire
Endpoint Agent possède la provenance de ses opérations techniques locales. Chaque source/destination conserve ses objets ; Shared conserve Trace/Linking génériques et Govern/Studio conservent leurs runs.

## 6. Utilisateurs
Endpoint Operator ; SOC/Investigate Analyst ; Detection Engineer ; Auditor ; Security/Privacy Reviewer ; Command/Govern/Settings/Studio consumers autorisés.

## 7. Conditions d’entrée
Stable refs des étapes EPT-2/EPT-3 ; versions/owners ; tenant/environment ; timestamps ; transformation/masking state ; destination permissions ; provenance manquante explicitement marquée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| observation/source provenance | CAP-EPT-015..030 | source chain | oui | source-owned | chain partial |
| content/evaluation/candidate provenance | CAP-EPT-031..036 | detection chain | oui si detection | operation time | chain partial |
| investigation/timeline/summary provenance | CAP-EPT-037..045 | investigation chain | oui si handoff | snapshot/version time | chain partial |
| cross-product owner/permission refs | governance/security sources | ownership context | oui | canonical/current | handoff restricted |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Local Observation/local-audit-event | Endpoint Agent | source/audit refs | read |
| telemetry-event/Trace/Timeline refs | Shared | generic refs only | read according source |
| Detection Content/Case/Evidence/Finding | Investigate | version/destination refs | no ownership transfer |
| Detection/Signal/Incident, Decision/Response Run/Result, Tool/Automation Run | Command/Govern/Studio | consumer/provenance refs | no local mutation |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| EPT-3 Provenance Chain | créer/étendre | Endpoint Agent | each hop typed/sourced |
| Cross-Product Handoff Reference | créer | Endpoint Agent | destination owner preserved |
| Provenance Gap Marker | dériver | Endpoint Agent | absence explicit, never invented |

## 11. Fonctionnalités
Retracer chaque hop, version, timestamp et transformation ; enregistrer rationale/context/AI attribution ; conserver destination/return-origin ; relier local audit ; exposer provenance gaps et permission restrictions ; ne jamais convertir Trace en owner des facts.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect provenance | Analyst/Auditor | Provenance Chain | 0 | read | chain visible | non |
| reconstruct chain | deterministic service | stable refs | 1 | refs resolvables | typed hops | non |
| forward provenance reference | authorized consumer | handoff | 2 | destination permission | reference only | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| reconstruct refs | oui | oui | oui | non nécessaire | stable graph traversal |
| detect provenance gaps | oui | oui | oui | explanation possible | completeness rules |
| summarize chain | oui | oui | oui | oui, attribuée | ordered hop list |
| invent missing observation/context | non | non | non | interdit | provenance-gap marker |

## 14. États fonctionnels
`complete`, `partial`, `gap-present`, `restricted`, `stale-reference`, `destination-denied`, `superseded`, `unknown`.

## 15. États d’interface
Aucun Screen ID. Future UI distingue source facts, derived states et AI proposal, expose gaps/restrictions et ne révèle pas des refs non autorisées.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| EPT-3 Provenance Chain | concept Endpoint | Auditor/Investigate/Quality | typed hops + owner/version/time |
| cross-product handoff refs | references | Command/Investigate/Govern/Studio/Settings | no ownership/permission transfer |
| provenance gap report | diagnostic | Analyst/Quality | missing facts not invented |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| EPT-2 observation | EPT-3 evaluation | detection provenance | source/content/eval refs | original source retained |
| detection context | local investigation | investigation provenance | candidate/context/timeline refs | detection origin retained |
| Endpoint Summary | consumer handoff | cross-product consumer | typed refs/limits/provenance | destination owner retained |

## 18. Dépendances
CAP-EPT-001..045 ; Endpoint local audit ; Shared Trace/Timeline/Linking ; Investigate Detection Engineering/Case/Evidence/Finding ; Command Detection/Signal/Incident ; Govern Decision/Response Run/Result ; Studio Tool/Automation Run ; Settings administration ; `OPEN-008`, `OPEN-015`, `OPEN-017`.

## 19. Source de vérité
Endpoint est SOT de la provenance de ses opérations EPT-3. Shared Trace reste un mécanisme générique, non propriétaire des facts. Chaque produit reste SOT de ses objets canoniques.

## 20. Provenance et audit
Ce contrat est lui-même normatif pour : Agent/version, tenant/environment, content/version/owner, observations/timestamps, evaluation/match/candidate, context transformations, masking, AI attribution, timeline, summary, destination et correlation ids.

## 21. Permissions fonctionnelles
Provenance read, restricted rationale/content/context read, cross-product ref read selon source, cross-tenant deny ; navigation n’accorde aucune source permission ; pas de final RBAC/ABAC.

## 22. Limites et erreurs
Provenance ≠ Evidence custody ; local audit ≠ Shared Trace ; Tool/Automation Run ≠ Endpoint evaluation ; Response Run/Result restent Govern ; missing hop reste gap ; aucune cryptographic protocol finalisée.

## 23. Métriques
Provenance completeness, gap/stale/denied refs, owner/version coverage, AI attribution completeness, handoff reconstruction success.

## 24. Classification de livraison
`draft / defined / planned`; documentation seulement, aucun ledger/protocol/API implémenté.

## 25. Critères d’acceptation
**Given** un summary est consommé par Investigate, **When** sa provenance est reconstruite, **Then** les hops observation→content→evaluation→candidate→context→timeline→summary sont résolubles ou marqués comme gaps.

**Given** une observation source est inaccessible à un consumer, **When** le handoff est inspecté, **Then** la référence reste permission-aware et n’accorde pas l’accès source.

**Given** l’IA a proposé une hypothèse de contexte, **When** la chaîne est auditée, **Then** cette étape est attribuée comme proposition et non comme fact sourcé.

## 26. Questions ouvertes
`OPEN-015` conserve les bridges Run cross-product ; `OPEN-017` conserve runtime/language ; `OPEN-008` conserve support plateforme/source. Aucun choix n’est fermé.

## 27. Consommateurs documentaires
Endpoint EPT-3 maps/quality ; Investigate ; Command ; Govern ; Studio ; Settings ; Shared ; Security/Audit ; Requirements ; Roadmap ; future EPT-4..6 comme consumers seulement.