---
id: CAP-EPT-049
title: File and Filesystem Collection
product: endpoint-agent
module: collection
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-014, REQ-PROD-018, REQ-PROD-020, REQ-SEC-001, REQ-SEC-004]
open_decisions: [OPEN-008, OPEN-013, OPEN-014]
source-of-truth: canonical
---
# CAP-EPT-049 — File and Filesystem Collection

## 1. Définition
Acquérir techniquement un fichier, une liste ou un scope de chemin/répertoire explicitement borné et autorisé, en conservant metadata, hash/integrity metadata disponible, changements pendant collecte, résultats par item et provenance.

## 2. Problème utilisateur
Le fichier peut disparaître, changer, être locked/restricted ou trop volumineux pendant l’acquisition ; un succès global ne doit jamais masquer ces états.

## 3. Objectifs
Respecter exact targets/bounds ; distinguer metadata observation et content acquisition ; gérer missing/locked/restricted/changed/partial ; conserver before/after metadata/hashes si disponibles ; produire `Collection Item` neutre.

## 4. Non-objectifs
Aucune procédure d’exfiltration, commande, wildcard illimité, suppression/quarantine/restore, Artifact/Evidence automatique, format archive ou transport.

## 5. Propriétaire
Endpoint possède l’acquisition locale et le Collection Item technique. Investigate qualifie ensuite tout Artifact/Evidence sous OPEN-014.

## 6. Utilisateurs
DFIR Analyst, Endpoint Operator, Evidence Reviewer, Case Analyst, Security/Privacy Reviewer.

## 7. Conditions d’entrée
Planned item CAP-EPT-048, target/path scope borné, capability file collection available, policy/permission/authority valides, limits applicables.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| planned file item | CAP-EPT-048 | exact target/bounds | oui | plan version | no collection |
| path/file refs | Investigate/EPT-3 | target refs | oui | revalidated at start | missing |
| capability/platform | Endpoint | support | oui | current | unsupported |
| policy/authority | Settings/Govern | gate | oui selon class | current | blocked |
| metadata/hash baseline | Endpoint | pre-acquisition facts | non | start time | explicit unavailable |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Technical Plan/Item | Endpoint | target/bounds | read |
| Collection Request | Investigate | purpose/scope | read |
| Endpoint Policy | Settings | path/size restrictions | read |
| EPT-3 file context | Endpoint | source ref only | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| File Collection Attempt | create/update | Endpoint | one bounded attempt |
| Collection Item | create/mark partial/failure | Endpoint | neutral technical output |
| File Change Marker | derive | Endpoint | before/after mismatch explicit |

## 11. Fonctionnalités
Revalidate target, acquire authorized content, track byte/item progress conceptually, record pre/post metadata/integrity metadata, preserve sparse/large/locked states, detect disappearance/change, return per-item outcome.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect target readiness | Operator | file item | 0 | read | state visible | non |
| acquire bounded file | authorized path | attempt | 2 | plan/gates | running/item output | according impact |
| cancel attempt | Operator | attempt | 2 | cancellable | cancel-requested | OPEN-013 |
| delete/quarantine file | aucun en EPT-4 | target | 3+ | EPT-5/Govern | not executed | obligatoire |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| validate path/bounds | oui | oui | oui | explain | validators |
| compare metadata/hash refs | oui | oui | oui | summarize | direct comparison |
| group per-item errors | oui | oui | oui | summarize | raw errors |
| extend scope | non | interdit | non | interdit | new request/plan |

## 14. États fonctionnels
`preparing`, `acquiring`, `completed`, `partial`, `missing`, `locked`, `restricted`, `changed-during-collection`, `failed`, `timed-out`, `cancel-requested`, `cancelled`, `unsupported`.

## 15. États d’interface
No Screen ID. Partial preserves successful bytes/items; Permission denied masks path/content; changed-during-collection is not presented as stable snapshot.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Collection Item | Endpoint technical output | CAP-EPT-054/055/Investigate | not Artifact/Evidence automatically |
| per-item status/error | Endpoint state | CAP-EPT-053 | failure not hidden |
| integrity/change metadata | Endpoint metadata | CAP-EPT-054 | not cryptographic proof automatically |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-048 | file item ready | CAP-EPT-049 | target/bounds/authority | plan |
| CAP-EPT-049 | output available | CAP-EPT-054 | item/status/metadata | attempt |
| CAP-EPT-049 | package/transfer ready | CAP-EPT-055 | neutral item ref | attempt |

## 18. Dépendances
Endpoint file-collection source, CAP-EPT-048/053..055, Investigate CAP-INV-205, Settings Policy, OPEN-008/013/014.

## 19. Source de vérité
Endpoint is SOT for target-side acquisition facts/Collection Item. Investigate remains SOT for any canonical Artifact/Evidence qualification.

## 20. Provenance et audit
Request/plan/item refs, Agent/platform, target path reference, bounds, metadata/hash refs, start/end, bytes/items, change marker, errors/cancel, actor/authority and correlation.

## 21. Permissions fonctionnelles
File collection execute, sensitive path/content, cancel, output read, provenance, cross-tenant deny; no final RBAC.

## 22. Limites et erreurs
File metadata observation ≠ acquisition ; file acquisition ≠ execution ; Collection Item ≠ Artifact/Evidence/Finding ; cancel ≠ rollback ; timeout ≠ confirmed target-side termination.

## 23. Métriques
Per-item completed/partial/missing/locked/changed/failure, bytes conceptually processed, cancellations/timeouts, provenance completeness.

## 24. Classification de livraison
`draft / defined / planned`; no command, transport, archive format or tool selected.

## 25. Critères d’acceptation
**Given** un fichier disparaît pendant la collecte, **When** l’attempt se termine, **Then** l’état indique missing/partial selon bytes acquis et aucun succès complet n’est affiché.

**Given** un fichier change pendant la collecte, **When** pre/post metadata diffèrent, **Then** `changed-during-collection` est conservé avec provenance.

**Given** l’IA est indisponible, **When** un fichier est acquis, **Then** validation, progression et erreurs restent déterministes.

## 26. Questions ouvertes
OPEN-008/013/014 restent ouvertes ; le statut Artifact/Attachment/Evidence est laissé à Investigate.

## 27. Consommateurs documentaires
EPT-4 packaging/transfer, Investigate Collection/Artifact/Evidence, Security, Quality.
