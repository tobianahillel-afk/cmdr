---
id: CAP-EPT-054
title: Collection Packaging, Completeness, Integrity and Metadata
product: endpoint-agent
module: collection
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-014, REQ-PROD-020, REQ-OBJ-004, REQ-SEC-004]
open_decisions: [OPEN-014]
source-of-truth: canonical
---
# CAP-EPT-054 — Collection Packaging, Completeness, Integrity and Metadata

## 1. Définition
Assembler conceptuellement des `Collection Item` dans un `Collection Package` technique avec manifest/reference, item status, completeness, missing/partial markers, size/time/source, integrity metadata et acquisition context, sans définir format d’archive ni Evidence package.

## 2. Problème utilisateur
Un package peut contenir des éléments partiels/manquants et des integrity metadata qui ne constituent pas une preuve cryptographique. Ces limites doivent survivre au handoff.

## 3. Objectifs
Lister exactement items attendus/reçus/manquants ; calculer completeness ; conserver metadata/provenance/integrity indicators ; préserver category failures ; produire package neutral sous OPEN-014.

## 4. Non-objectifs
Aucun archive/container format, hash algorithm imposé, PKI/signature protocol, Evidence package, Artifact/Attachment identity, storage engine or legal conclusion.

## 5. Propriétaire
Endpoint owns technical packaging/completeness metadata. Investigate/Trust retain custody/Artifact/Evidence qualification.

## 6. Utilisateurs
DFIR Analyst, Evidence Reviewer, Endpoint Operator, Auditor, Security/Trust Reviewer.

## 7. Conditions d’entrée
One or more Collection Items or explicit expected/missing items, source plan/request refs, acquisition contexts and available integrity metadata.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| planned item manifest refs | CAP-EPT-048 | expected set | oui | plan version | completeness unknown |
| collected items/status | CAP-EPT-049..053 | outputs | oui at least status | terminal/current | partial |
| source/acquisition context | Endpoint | provenance | oui | item time | gap marker |
| integrity metadata | source-specific | checksum/hash/signature-like refs | non | item time | not-available |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Technical Plan/Items | Endpoint | expected/actual | read |
| Collection Request | Investigate | business context | read |
| local-audit-event | Endpoint | source audit refs | read |
| Artifact/Evidence | Investigate | destination concepts only | no creation |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Collection Package | create/version | Endpoint | conceptual neutral package |
| Completeness Assessment | derive | Endpoint | expected vs actual |
| Integrity Metadata Set | attach | Endpoint | source-attributed, not proof automatically |

## 11. Fonctionnalités
Assemble item refs, manifest expected/actual statuses, include size/time/source/acquisition context, mark partial/missing/corrupt-or-unverified states, preserve item-level errors and version package when late items arrive.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect package/items | Analyst/Reviewer | package | 0 | read | contents/gaps visible | non |
| compute completeness | service | package | 1 | expected/actual refs | assessment | non |
| verify available integrity metadata | reviewer/service | metadata | 1/2 | source metadata | verified-at-metadata-level/disputed | no automatic Evidence authority |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| build manifest | oui | oui | oui | non needed | deterministic refs |
| compute completeness | oui | oui | oui | explain | expected/actual diff |
| summarize gaps | oui | oui | oui | oui | raw gap list |
| fabricate missing item/proof | non | interdit | non | interdit | gap marker |

## 14. États fonctionnels
`assembling`, `complete`, `partial`, `missing-items`, `integrity-metadata-present`, `integrity-unverified`, `disputed`, `superseded`.

## 15. États d’interface
No Screen ID. Package complete status must be separate from downstream validation/Evidence qualification.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Collection Package | Endpoint concept | CAP-EPT-055/Investigate | not Evidence package |
| Completeness Assessment | metadata | Investigate/Quality | expected/actual transparent |
| Integrity Metadata | source-attributed metadata | Trust/Investigate | not crypto proof automatically |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-049..053 | items/status available | CAP-EPT-054 | expected/actual/provenance | attempts |
| CAP-EPT-054 | ready for delivery | CAP-EPT-055 | package/gaps/integrity metadata | package |
| Investigate | qualify | Artifact/Evidence workflows | neutral refs + custody context | source package retained |

## 18. Dépendances
Collection triage-package/chain-of-custody sources, CAP-EPT-048..055, Investigate CAP-INV-204/213, OPEN-014.

## 19. Source de vérité
Endpoint SOT of package technical content/completeness metadata; Investigate SOT of Artifact/Evidence/custody business qualification.

## 20. Provenance et audit
Package/version, request/plan, expected/actual items, sizes/times/sources, acquisition contexts, integrity metadata source/method label, gaps/disputes, actor and correlation.

## 21. Permissions fonctionnelles
Package/item metadata read, sensitive content read, completeness/integrity review, provenance, cross-tenant deny.

## 22. Limites et erreurs
Collection Package ≠ Evidence package; integrity metadata ≠ cryptographic proof automatically; item ≠ Artifact/Evidence; completeness does not prove content validity.

## 23. Métriques
Complete/partial packages, missing items, integrity-metadata coverage, disputed/gap counts, provenance completeness.

## 24. Classification de livraison
`draft / defined / planned`; no archive format, hash algorithm, PKI, storage or physical schema.

## 25. Critères d’acceptation
**Given** one expected item fails, **When** package is assembled, **Then** it remains partial/missing-item and failure is in manifest.

**Given** integrity metadata exists, **When** reviewed, **Then** it is source-attributed and not presented as cryptographic proof automatically.

**Given** OPEN-014 remains open, **When** package is handed to Investigate, **Then** no Artifact/Attachment/Evidence identity is assigned by Endpoint.

## 26. Questions ouvertes
OPEN-014 remains open; legal hold/retention/Artifact identity finalization stays outside Endpoint EPT-4.

## 27. Consommateurs documentaires
CAP-EPT-055/064, Investigate Artifact/Evidence/Custody, Trust, Quality.
