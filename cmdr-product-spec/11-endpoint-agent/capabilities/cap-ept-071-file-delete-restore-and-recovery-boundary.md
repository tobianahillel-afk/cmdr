---
id: CAP-EPT-071
title: File Delete, Restore and Recovery Boundary
product: endpoint-agent
module: containment
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-006, REQ-PROD-016, REQ-PROD-018, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008, OPEN-013, OPEN-014, OPEN-015]
source-of-truth: canonical
---
# CAP-EPT-071 — File Delete, Restore and Recovery Boundary

## 1. Définition
Définir les frontières techniques pour delete/restore/recovery d’un fichier lorsque les sources l’autorisent, avec authority explicite, risque d’irréversibilité, état de suppression/restauration observé et provenance conservée.

## 2. Problème utilisateur
Delete peut être irréversible ou seulement logique ; restore peut échouer ou restaurer un contenu non sûr. L’historique ne doit jamais disparaître avec le fichier.

## 3. Objectifs
Pinner identity ; classify destructive risk ; require stronger authority ; execute bounded delete or restore only when supported ; preserve provenance ; observe deleted/restored state ; expose recovery availability/limitations.

## 4. Non-objectifs
Garantir unrecoverability, fournir deletion/restore commands, promettre recovery, qualifier restored file safe, supprimer provenance/Evidence, ou définir filesystem implementation.

## 5. Propriétaire
Endpoint owns technical file delete/restore primitive facts. Govern owns authority and rollback/recovery governance; Investigate owns Artifact/Evidence/Finding.

## 6. Utilisateurs
Response Operator, Endpoint Operator, Govern Reviewer, Verification Reviewer, DFIR Analyst, Auditor.

## 7. Conditions d’entrée
Exact file identity, CAP-EPT-066 readiness, explicit delete/restore request, irreversible-risk context, authority current, recovery/restore capability state known where applicable.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| file identity/current state | Endpoint | target | oui | fresh | blocked/unknown |
| delete or restore request | Govern | effect intent | oui | pinned | no action |
| authority + risk context | Govern/Security | authorization | oui | current | blocked |
| recovery/restore availability | Endpoint | capability | restore/recovery | current | unavailable/unknown |
| provenance/related refs | Endpoint/Investigate | history | oui | stable refs | execution blocked if audit impossible |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Decision/Response Run/Rollback Plan | Govern | effect/authority/scope | read |
| file/quarantine state | Endpoint | identity/current state | read |
| Artifact/Evidence refs | Investigate | protected provenance | link/read |
| Endpoint Policy | Settings | destructive action restriction | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| File State Execution | create/transition | Endpoint | exact operation/identity |
| Delete/Restore State | observe/update | Endpoint | technical state only |
| Technical Outcome | create | Endpoint | not Result |
| file | delete/restore/recover effect | target | history/provenance retained separately |

## 11. Fonctionnalités
Revalidate identity; expose destructive/irreversible risk; require authority; execute supported bounded delete/restore; track requested/start/partial/fail/unknown; observe existence/restored state; retain immutable provenance; expose recovery availability without guarantee.

## 12. Actions utilisateur
Inspect/availability Class 0/1. Delete is Class 3/4 depending irreversibility and always governed. Restore/recovery is Class 3 unless source policy explicitly classifies otherwise; OPEN-013 remains unresolved.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| validate identity/authority | oui | oui | oui | explain | checklist |
| observe existence/restore state | oui | oui | oui | summarize | local facts |
| suggest restore/recovery | oui | capability rules | oui | recommendation | operator selection |
| authorize destructive action | non | no autonomous | non | interdit | Govern |

## 14. États fonctionnels
`delete-requested`, `deleting`, `deleted-observed`, `delete-partial`, `delete-failed`, `delete-unknown`, `recovery-unavailable`, `recovery-available`, `restore-requested`, `restoring`, `restored-observed`, `restore-partial`, `restore-failed`, `restore-unknown`, `identity-changed`, `superseded`.

## 15. États d’interface
No Screen ID. Delete requested/observed, recovery availability, restore state and file safety remain separate concepts.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Delete/Restore State | technical fact | CAP-EPT-073/075 | source/time |
| destructive/recovery outcome | technical outcome | Govern | irreversibility/limits retained |
| provenance continuity | audit refs | Audit/Investigate | no historical deletion |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-066/070 | file action ready | CAP-EPT-071 | identity/op/authority/risk | Run |
| CAP-EPT-071 | outcome | CAP-EPT-073/076 | technical state/limits | same Run |
| restore/reversal request | Govern | CAP-EPT-077/071 | prior state/recovery capability | Rollback relation retained |

## 18. Dépendances
CAP-EPT-038/049/054/070/073..080, Govern rollback/result, Investigate evidence, Settings Policy, Security, OPEN-008/013/014/015.

## 19. Source de vérité
Endpoint SOT of local delete/restore technical state; Govern SOT of response rollback/result; Investigate SOT of evidence/artifact semantics.

## 20. Provenance et audit
Identity/path/hash, before/after state, requested operation, authority/risk/Run refs, recovery source, operator, Agent/version, timestamps, partial/error, verification, retained provenance links.

## 21. Permissions fonctionnelles
File delete, restore/recovery request, state read, sensitive context, verification, cross-tenant deny; step-up/SoD and Govern dependency mandatory; no final RBAC.

## 22. Limites et erreurs
Deleted ≠ guaranteed unrecoverable; restore ≠ safe; recovery available ≠ successful; technical reversal ≠ Govern rollback; provenance is never deleted.

## 23. Métriques
Delete/restore attempts/outcomes, recovery unavailable, partials, verification mismatch, provenance gaps target zero.

## 24. Classification de livraison
`draft / defined / planned`; no deletion/restore command, filesystem mechanism, API/protocol or implementation.

## 25. Critères d’acceptation
**Given** deletion lacks sufficient authority, **When** requested, **Then** no effect occurs and the denial/authority gap is audited.

**Given** restore technically succeeds, **When** file existence is observed, **Then** restored state is recorded but file safety remains unknown.

**Given** deletion completes, **When** audit is reconstructed, **Then** prior identity/authority/outcome provenance remains available regardless of file availability.

## 26. Questions ouvertes
OPEN-008/013/014/015 remain open; final recovery guarantees and Artifact/Attachment relations are not selected.

## 27. Consommateurs documentaires
EPT-5 verification/reversal, Govern, Investigate, Security, Quality, Audit.