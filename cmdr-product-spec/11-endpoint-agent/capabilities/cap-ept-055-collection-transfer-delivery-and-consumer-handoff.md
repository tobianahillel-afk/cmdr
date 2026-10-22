---
id: CAP-EPT-055
title: Collection Transfer, Delivery and Consumer Handoff
product: endpoint-agent
module: collection
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-014, REQ-PROD-018, REQ-PROD-020, REQ-SEC-001, REQ-SEC-004]
open_decisions: [OPEN-008, OPEN-014, OPEN-015]
source-of-truth: canonical
---
# CAP-EPT-055 — Collection Transfer, Delivery and Consumer Handoff

## 1. Définition
Définir le transfert/delivery fonctionnel d’un Collection Item/Package depuis Endpoint vers une destination autorisée, avec destination ref, transfer state, acknowledgement, partial/failure, retention reference et provenance, sans choisir protocole ni promouvoir l’output en Evidence/Artifact.

## 2. Problème utilisateur
Acquisition réussie peut être suivie d’un transfert échoué ; `transfer complete` peut être confondu avec contenu validé ou Evidence reçue.

## 3. Objectifs
Séparer acquisition et transfer ; préserver package/item refs ; exposer destination/consumer/acknowledgement ; gérer partial/failure/retry eligibility ; conserver restrictions/retention ref ; handoff permission-aware.

## 4. Non-objectifs
Aucun transfer protocol, storage engine, Export ownership, Artifact/Evidence auto-creation, content validation, external sharing or raw secret handling.

## 5. Propriétaire
Endpoint owns local transfer state and technical handoff refs. Shared may own generic transfer/export infrastructure; Investigate owns destination qualification/Artifact/Evidence.

## 6. Utilisateurs
DFIR Analyst, Endpoint Operator, Evidence Reviewer, Auditor, Platform/Storage consumer as authorized.

## 7. Conditions d’entrée
Collected output/package, destination/consumer ref, permission and tenant scope, transfer capability/availability, retention/restriction context when applicable.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Collection Item/Package | CAP-EPT-049..054 | technical output | oui | current version | no transfer |
| destination/consumer | Investigate/Shared/authorized owner | target ref | oui | current | blocked |
| transfer capability/state | Endpoint/Shared projection | availability | oui | current | unavailable |
| permission/tenant/retention refs | Security/Settings | restrictions | oui selon output | current | denied/restricted |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Collection Package/Item | Endpoint | output refs/status | read |
| Case/Collection Request | Investigate | destination context | read |
| Export/Job | Shared | optional mechanism refs | read/link only |
| retention/policy | Settings/Trust | restrictions | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Collection Transfer | create/update | Endpoint | technical state only |
| Delivery/Acknowledgement Ref | attach | Endpoint | source/destination attributed |
| Consumer Handoff Ref | create | Endpoint | no ownership/permission transfer |

## 11. Fonctionnalités
Prepare bounded transfer, validate destination/tenant/permissions, track transferring/partial/failed/completed, record acknowledgement separately, preserve acquisition completeness and errors, expose retry eligibility without silent retry.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect transfer state | Analyst | transfer | 0 | read | state visible | non |
| validate destination/restrictions | service/reviewer | handoff | 1 | refs | pass/block | non |
| request/retry bounded transfer | authorized operator | transfer | 2 | permission/capability | transfer attempt | according sensitivity |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| validate destination | oui | oui | oui | explain | rules |
| track transfer | oui | oui | oui | summarize | raw progress |
| explain failure | oui | catalogue | oui | oui | error details |
| qualify Evidence/Artifact | non Endpoint | external workflow | non Endpoint | interdit | Investigate qualification |

## 14. États fonctionnels
`preparing`, `transferring`, `partial`, `failed`, `timed-out`, `cancelled`, `completed`, `acknowledgement-pending`, `acknowledged`, `destination-denied`, `retention-restricted`.

## 15. États d’interface
No Screen ID. Transfer complete, acknowledgement, acquisition completeness and content validation are shown as separate facts.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Collection Transfer state | Endpoint | Investigate/Shared | technical only |
| delivered output refs | reference set | Investigate | Collected Output != Evidence/Artifact automatically |
| acknowledgement/restriction | metadata | Audit/Quality | destination/time/source explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-054 | package ready | CAP-EPT-055 | package/gaps/restrictions | package |
| CAP-EPT-055 | delivered | Investigate | neutral output/provenance/limits | transfer ref |
| CAP-EPT-055 | generic mechanism | Shared Job/Export | technical transfer request ref | mechanism remains Shared |

## 18. Dépendances
CAP-EPT-049..054/064, Investigate CAP-INV-203/211/213/214, Shared Jobs/Export, Settings/Trust restrictions, OPEN-008/014/015.

## 19. Source de vérité
Endpoint SOT of local transfer facts; destination owner SOT of receipt/qualification; Shared SOT of generic mechanism if used.

## 20. Provenance et audit
Output/package version, source Agent, destination/consumer, permissions, transfer attempts/progress/errors, acknowledgement, retention refs, actor and correlation.

## 21. Permissions fonctionnelles
Transfer request/read/cancel/retry, sensitive output/destination read, provenance, cross-tenant deny; no final RBAC.

## 22. Limites et erreurs
Transfer complete ≠ content validated ; transfer failed after acquisition does not erase acquired output ; output ≠ Evidence/Finding/Artifact automatically ; retry ≠ duplicate-free guarantee.

## 23. Métriques
Transfer completed/partial/failure/timeout, acknowledgement latency, retry candidates, destination denials, qualification-boundary violations target zero.

## 24. Classification de livraison
`draft / defined / planned`; no transfer protocol/storage/export implementation.

## 25. Critères d’acceptation
**Given** acquisition succeeded but transfer fails, **When** status is reviewed, **Then** acquisition remains successful locally while delivery is failed/partial.

**Given** transfer completes, **When** destination acknowledges, **Then** acknowledgement is recorded separately from content validation/Evidence qualification.

**Given** output reaches Investigate, **When** handoff completes, **Then** Investigate performs any Artifact/Evidence qualification and OPEN-014 remains open.

## 26. Questions ouvertes
OPEN-008/014/015 remain open; transfer implementation and Artifact identity remain undecided.

## 27. Consommateurs documentaires
Investigate Collection/Artifact/Evidence, Shared Jobs/Export, CAP-EPT-064, Security/Trust, Quality.
