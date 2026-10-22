---
id: CAP-EPT-061
title: Live Response File Read, Download and Transfer Operations
product: endpoint-agent
module: live-response
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-014, REQ-PROD-018, REQ-PROD-020, REQ-SEC-001, REQ-SEC-004]
open_decisions: [OPEN-008, OPEN-013, OPEN-014, OPEN-015]
source-of-truth: canonical
---
# CAP-EPT-061 — Live Response File Read, Download and Transfer Operations

## 1. Définition
Définir pendant une Live Response Session des opérations bornées de file read/download/transfer et, uniquement parce que le corpus Investigate le source explicitement, d’upload temporaire authority-bound, avec source/destination/finality/size/type/provenance/progress/collision/cleanup refs.

## 2. Problème utilisateur
File transfer peut être confondu avec Collection, Attachment, deployment ou remediation. Un upload peut modifier la cible et ne doit jamais être implicite ou Class 1.

## 3. Objectifs
Distinguer read/download/upload temporary/finality ; valider exact source/destination/collision/policy ; suivre progress/partial/failure/cancel ; préserver classification/provenance ; classer upload comme effectful avec authority stricte ; empêcher delete/quarantine/restore EPT-5 actions.

## 4. Non-objectifs
Aucun protocole/path technique universel, deployment/package manager, delete/move/quarantine/restore remediation, Artifact/Evidence auto qualification, software installation or containment.

## 5. Propriétaire
Endpoint owns target-side file operation facts. Investigate owns business transfer context/Artifact qualification; Settings policy/secrets; Govern authority for effectful upload/collision overrides.

## 6. Utilisateurs
Response Operator, DFIR Analyst, Endpoint Operator, Evidence Reviewer, Govern/Security Reviewer.

## 7. Conditions d’entrée
Active technical session, exact file/source reference, direction/finality, bounded destination, policy/permission, size/type/provenance metadata and authority according direction/effect.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| session/target | CAP-EPT-056/057 | execution context | oui | current | no operation |
| direction/finality | Investigate/caller | read/download/upload-temp | oui | request version | incomplete |
| file/source ref + metadata | source owner | content reference | oui | exact version | invalid |
| destination/collision policy | Settings/caller | bounds | oui when transfer | start time | blocked |
| authority | Govern/Security | gate | upload/effectful especially | current | awaiting-authority |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Technical Session | Endpoint | target/state | read |
| business transfer/Case | Investigate | direction/finality | read ref |
| Artifact/Attachment | Investigate/open | source/destination concepts | no identity inference |
| Endpoint Policy | Settings | path/size/collision restrictions | read |
| Decision/Response Run | Govern | authority ref | read only |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Technical File Transfer Attempt | create/update | Endpoint | bounded direction/finality |
| File Read/Download Output Ref | create | Endpoint | neutral technical output |
| Temporary Upload Ref | create/close | Endpoint | effectful, no deployment identity |
| Collision/Cleanup Status | record | Endpoint | confirmation separate from rollback |

## 11. Fonctionnalités
Validate direction/source/destination/size/type/collision; read/download authorized content; permit temporary upload only with explicit authority and finality; track transfer/verifying/partial/failure/cancel; record cleanup intent/status without declaring state restoration.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| bounded file read/download | DFIR/Response Operator | transfer | 2 | session + permission | transfer attempt | generally low-impact; policy applies |
| temporary upload | authorized Response Operator | transfer | 3 | authority + source/destination + cleanup | upload attempt | obligatoire/strict |
| collision override | authorized operator | transfer | 3 | explicit decision/policy | audited attempt | obligatoire |
| delete/quarantine/restore | aucun EPT-4 | target file | 3+ | EPT-5 | not executed | obligatoire |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| validate source/type/size/collision | oui | oui | oui | explain | validators |
| track progress | oui | oui | oui | summarize | raw progress |
| propose cleanup | oui | policy | oui | suggestion | checklist |
| upload/override autonomously | non | authority contract | non autonome | interdit | Govern/operator path |

## 14. États fonctionnels
`preparing`, `ready`, `awaiting-authority`, `transferring`, `verifying`, `completed`, `partial`, `failed`, `timed-out`, `cancel-requested`, `cancelled`, `collision`, `restricted`, `cleanup-pending`, `cleanup-unknown`.

## 15. États d’interface
No Screen ID. Upload is visibly effectful; transfer completion, verification and cleanup are separate facts.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Technical File Transfer result | Endpoint | CAP-EPT-062/064/Investigate | direction/finality/progress/errors |
| downloaded/read output ref | neutral output | Investigate | not Artifact/Evidence automatically |
| temporary upload/cleanup refs | technical facts | Govern/Audit | no deployment/rollback inference |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Live Session | transfer requested | CAP-EPT-061 | session/file/direction/authority | session |
| CAP-EPT-061 | output/error | CAP-EPT-062 | attempt/status/output | session |
| CAP-EPT-061 | handoff | Investigate | neutral output/provenance | qualification external |

## 18. Dépendances
Endpoint file-actions boundary, Investigate CAP-INV-211, CAP-EPT-056/057/062..064, Govern, Settings, OPEN-008/013/014/015.

## 19. Source de vérité
Endpoint SOT of local file operation facts; Investigate SOT business transfer and Artifact qualification; Govern SOT effectful authority.

## 20. Provenance et audit
Session, target, operator, direction/finality, source/version, destination functional ref, type/size, policy/authority, collision choice, progress/errors, output/upload/cleanup refs, timestamps/correlation.

## 21. Permissions fonctionnelles
File read/download, temporary upload, sensitive content, collision override, cancel, output read, provenance, cross-tenant deny; upload always stronger authority than passive read.

## 22. Limites et erreurs
File transfer ≠ Collection qualification automatically; file upload ≠ script executed; upload ≠ safe/deployed; cancellation ≠ rollback; cleanup ≠ target state restored; destructive file actions are EPT-5 boundary.

## 23. Métriques
Transfers by direction/finality, upload authority blocks, collision/partial/failure/cancel, cleanup-unknown, qualification-boundary violations target zero.

## 24. Classification de livraison
`draft / defined / planned`; no transfer protocol/path/deployment implementation.

## 25. Critères d’acceptation
**Given** a file download succeeds, **When** output is handed to Investigate, **Then** it remains a neutral technical output until Investigate qualifies it.

**Given** an upload is requested, **When** authority is absent, **Then** it is blocked and never treated as Class 1 implicit transfer.

**Given** a destructive file action is requested, **When** EPT-4 evaluates it, **Then** it stops at EPT-5/Govern boundary without executing it.

## 26. Questions ouvertes
OPEN-008/013/014/015 remain open; file identity, protocol and cleanup implementation unresolved.

## 27. Consommateurs documentaires
CAP-EPT-062..064, Investigate Session File Transfer/Artifact/Evidence, Govern, Settings, Security, Quality.
