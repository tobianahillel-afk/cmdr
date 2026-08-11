---
id: CAP-EPT-064
title: Collection and Live Response Provenance and Cross-Product Execution Contracts
product: endpoint-agent
module: collection-live-response
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-006, REQ-PROD-014, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-OBJ-004, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008, OPEN-013, OPEN-014, OPEN-015, OPEN-017]
source-of-truth: canonical
---
# CAP-EPT-064 — Collection and Live Response Provenance and Cross-Product Execution Contracts

## 1. Définition
Retracer les chaînes EPT-4 Collection et Live Response depuis source request/authority jusqu’aux opérations, outputs, handoffs et consumers, en conservant owners et sans fusionner Collection Request, Artifact/Evidence, Tool Call/Automation Run, Response Run ou Result avec les objets techniques Endpoint.

## 2. Problème utilisateur
Cross-product execution peut perdre l’origine, l’autorité, les limites ou convertir un technical success/output en canonical business object. Une provenance reconstructible est nécessaire pour audit et reconciliation.

## 3. Objectifs
Chaîne Collection complète ; chaîne Live Response complète ; exact owner/version/time/correlation ; authority/Secret refs sans raw values ; handoff permission-aware ; gaps explicites ; return origin conservé ; OPEN-014/015 non résolues.

## 4. Non-objectifs
Aucun cross-product API/protocol, immutable ledger, cryptographic protocol, physical schema, object identity merger, Evidence custody finalization, Response Run bridge final or containment.

## 5. Propriétaire
Endpoint owns provenance of its technical collection/execution facts. Each external owner retains its objects: Investigate, Govern, Studio, Settings, Shared.

## 6. Utilisateurs
Auditor, DFIR Analyst, Response Operator, Endpoint Operator, Govern Reviewer, Studio Operator, Security/Trust Reviewer.

## 7. Conditions d’entrée
Stable refs from EPT-3/EPT-4 stages, owner/version/timestamps, tenant/environment, authority/provenance refs, masking state and destination permissions.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| EPT-3 collection-required/summary refs | CAP-EPT-043..046 | upstream context | conditionnel | source time | chain starts at request |
| collection chain refs | CAP-EPT-047..055 | technical acquisition chain | conditionnel | operation time | partial chain |
| live response chain refs | CAP-EPT-056..063 | technical execution chain | conditionnel | operation time | partial chain |
| external owner refs | Investigate/Govern/Studio/Settings/Shared | business/authority/mechanism refs | according workflow | source-owned | gap/restricted |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Collection Request/Case/Artifact/Evidence/Finding | Investigate | source/destination refs | no ownership transfer |
| Action Request/Decision/Response Run/Result | Govern | authority/reconciliation refs | no mutation |
| Tool/Tool Call/Automation Run | Studio | origin/provenance refs | no mutation |
| Secret Reference/Fleet/Policy | Settings | refs/version only | restricted read |
| Job/Trace/Activity/Export | Shared | generic mechanism refs | read/link only |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| EPT-4 Provenance Chain | create/extend | Endpoint | typed hops/owners |
| Cross-Product Handoff Reference | create | Endpoint | destination owner retained |
| Provenance Gap/Restriction Marker | derive | Endpoint | missing data never invented |

## 11. Fonctionnalités
Collection chain: EPT-3 need → Investigate request → local eligibility/authority → plan → operation/item → package → transfer → consumer. Live chain: request origin → technical session → authority context → command/script/file request → execution → technical output → operator/session closure → consumer. Preserve refs and limitations at every hop.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect provenance | Auditor/Reviewer | chain | 0 | read | typed hops/gaps | non |
| reconstruct chain | service | stable refs | 1 | refs resolvable | chain | non |
| forward refs to consumer | authorized operator/service | handoff | 2 | destination permission | reference only | according source action |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| correlate IDs | oui | oui | oui | no need | stable refs |
| detect gaps | oui | oui | oui | explain | completeness rules |
| summarize chain | oui | oui | oui | oui, attributed | ordered hops |
| invent authority/output/Evidence/Result | non | interdit | non | interdit | gap/owner workflow |

## 14. États fonctionnels
`complete`, `partial`, `gap-present`, `restricted`, `stale-reference`, `destination-denied`, `superseded`, `unknown`.

## 15. États d’interface
No Screen ID. Source facts, derived states, AI summaries and external canonical objects remain distinguishable; inaccessible refs are not leaked.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| EPT-4 Provenance Chain | Endpoint concept | Audit/Investigate/Govern | owner/version/time typed |
| handoff refs | references | external owners | no permission transfer |
| provenance gaps | diagnostic | Quality/Audit | no invention |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Collection Request | technical handoff | EPT-4 Collection | request/authority refs | Investigate owner retained |
| Govern Response Run | bounded handoff | Endpoint execution | Decision/Run/step/target/scope | Govern owner retained |
| Studio Tool Call/Automation Run | caller handoff | Endpoint execution | source run/call refs | Studio owner retained |
| Endpoint output | reconciliation | Investigate/Govern | technical output/status/limits | destination qualifies |

## 18. Dépendances
CAP-EPT-043..063, Investigate CAP-INV-202..214, Govern execution boundaries/CAP-GOV-025, Studio Tool Call/Automation Run, Settings Secrets/Policy, Shared Jobs/Trace/Activity/Export, OPEN-008/013/014/015/017.

## 19. Source de vérité
Endpoint SOT of its technical provenance chain; each external product remains SOT of its canonical objects and generic mechanisms.

## 20. Provenance et audit
All request/session/operation/item/output/transfer/control IDs, target/tenant, versions, actor/operator, policy, authority, Secret References only, timestamps, masking, errors, external run/call refs, destination and correlation IDs.

## 21. Permissions fonctionnelles
Provenance read, restricted output/session refs, external object refs under source permission, export preparation, cross-tenant deny; navigation never grants permission.

## 22. Limites et erreurs
Collection Item/Output ≠ Artifact/Evidence/Finding automatically; Endpoint Technical Execution ≠ Tool Call/Automation Run/Response Run; technical output ≠ Result; local cancel/cleanup ≠ Govern rollback.

## 23. Métriques
Chain completeness/gaps, orphaned external refs, owner/version coverage, cross-tenant denials, technical objects incorrectly promoted/merged target zero.

## 24. Classification de livraison
`draft / defined / planned`; no API/protocol/ledger/schema/bridge implementation.

## 25. Critères d’acceptation
**Given** a Studio Tool Call invokes a future Endpoint bridge, **When** provenance is reconstructed, **Then** Tool Call remains Studio-owned and Endpoint technical execution is a distinct hop.

**Given** a Govern Response Run hands off an effectful request, **When** technical output returns, **Then** Response Run remains Govern-owned and Govern separately reconciles Result.

**Given** collected bytes reach Investigate, **When** chain is inspected, **Then** they remain neutral Endpoint output until Investigate explicitly qualifies Artifact/Evidence under OPEN-014.

## 26. Questions ouvertes
OPEN-008/013/014/015/017 remain open. No bridge/object identity is finalized.

## 27. Consommateurs documentaires
Investigate, Govern, Studio, Settings, Shared, Security/Trust, Quality, Roadmap and future EPT-5/EPT-6 as consumers/boundaries only.
