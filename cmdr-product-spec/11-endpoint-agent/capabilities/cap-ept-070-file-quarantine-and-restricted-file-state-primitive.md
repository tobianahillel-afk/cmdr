---
id: CAP-EPT-070
title: File Quarantine and Restricted File-State Primitive
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
# CAP-EPT-070 — File Quarantine and Restricted File-State Primitive

## 1. Définition
Définir la primitive technique de mise en quarantaine/restriction d’un fichier identifié, avec path/hash/identity precheck, requested/applied state, conservation de l’origine, release/restore reference et verification.

## 2. Problème utilisateur
Le fichier peut disparaître ou changer entre analyse et action. `quarantine requested`, `quarantined`, `deleted` et `malicious` ne sont pas équivalents.

## 3. Objectifs
Pinner file identity ; detect missing/changed target ; quarantine/restrict under authority ; preserve original path/hash metadata ; observe applied state ; expose inaccessible/partial/failure ; prepare release/restore without claiming verdict.

## 4. Non-objectifs
Définir secure-store implementation, antivirus verdict, Evidence qualification, delete semantics, OS commands, restore safety or supported platform.

## 5. Propriétaire
Endpoint owns quarantine primitive and File Quarantine State. Investigate owns Artifact/Evidence/Finding qualification. Govern owns authority/Run/Result.

## 6. Utilisateurs
Response Operator, Endpoint Operator, DFIR Analyst, Govern Reviewer, Verification Reviewer, Auditor.

## 7. Conditions d’entrée
CAP-EPT-065/066 ready; exact file ref/path/hash or equivalent identity; latest file-state observation; authority; policy; quarantine capability available.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| file ref/path/hash/identity | EPT-3/EPT-4 | target | oui | fresh precheck | target-unknown |
| requested quarantine state | Govern | intent | oui | pinned | reject |
| current file state | Endpoint | precheck | oui | timestamped | unknown |
| quarantine capability/policy | Endpoint/Settings | readiness | oui | current | unsupported/blocked |
| authority ref | Govern | authorization | oui | effective | no effect |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Response Run/Decision | Govern | scope/authority | read |
| file observations/Collection Item refs | Endpoint | identity/current state | read |
| Artifact/Evidence/Finding refs | Investigate | provenance only | read/link |
| Endpoint Policy | Settings | quarantine restriction | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Quarantine Execution | create/transition | Endpoint | identity pinned |
| File Quarantine State | observe/update | Endpoint | quarantined ≠ deleted |
| Technical Outcome | create | Endpoint | no malware verdict |
| target file access/state | quarantine/restrict effect | target | implementation unspecified |

## 11. Fonctionnalités
Revalidate file identity; reject changed/missing target safely; request quarantine; preserve original identity metadata; observe applied/restricted state; track inaccessible/partial/failure/unknown; expose release/restore eligibility and provenance.

## 12. Actions utilisateur
Inspect Class 0; identity/readiness/verification Class 1; quarantine/release effect Class 3 by default, Govern-dependent.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| compare file identity | oui | oui | oui | explain | path/hash/metadata diff |
| observe quarantine state | oui | oui | oui | summary | local observation |
| propose target | oui | source-backed | oui | suggestion | analyst selection |
| authorize/quarantine autonomously | non | interdit | non | interdit | Govern/operator |

## 14. États fonctionnels
`not-quarantined`, `target-missing`, `identity-changed`, `quarantine-requested`, `quarantining`, `quarantined-observed`, `restricted-observed`, `partial`, `failed`, `inaccessible`, `unknown`, `release-eligible`, `release-requested`, `released-observed`, `drifted`.

## 15. États d’interface
No Screen ID. Missing, inaccessible, changed identity, requested and observed states are distinct; maliciousness is never inferred from quarantine status.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| File Quarantine State | technical state | CAP-EPT-073/075/078 | source/time/identity |
| quarantine outcome | technical fact | Govern | partial/failure explicit |
| release/restore ref | technical eligibility | CAP-EPT-071/077/078 | no automatic reversal |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-066 | quarantine ready | CAP-EPT-070 | file identity/state/authority | Run |
| CAP-EPT-070 | outcome | CAP-EPT-073/076 | requested/observed state | same Run |
| CAP-EPT-070 | verify/release | CAP-EPT-074/075/078 | file/quarantine observations | Run retained |

## 18. Dépendances
CAP-EPT-018/038/049/061/062/065/066/071/073..080, Investigate Evidence/Artifact, Govern, Settings, OPEN-008/013/014/015.

## 19. Source de vérité
Endpoint SOT of local quarantine execution/state. Investigate remains SOT for evidence/artifact/finding semantics; Govern for response outcome.

## 20. Provenance et audit
File ref/path/hash/metadata, before/after observations, original location metadata, authority/Run/Step, policy, operator, Agent/version, timestamps, errors, verification/release refs.

## 21. Permissions fonctionnelles
Quarantine read/request/release, sensitive file context, verification, cross-tenant deny, step-up/SoD/Govern dependency. Evidence permission remains separate.

## 22. Limites et erreurs
Quarantine requested ≠ applied; quarantined ≠ deleted; quarantined ≠ malicious verdict; changed file identity invalidates stale request; release ≠ safe file.

## 23. Métriques
Missing/changed target, applied/partial/fail/unknown, release, verification mismatch, unauthorized attempts target zero.

## 24. Classification de livraison
`draft / defined / planned`; no secure-store format, quarantine command, API, physical schema or implementation.

## 25. Critères d’acceptation
**Given** target file no longer exists, **When** quarantine begins, **Then** target-missing is recorded and no substitute file is affected.

**Given** the file changed before quarantine, **When** identity is revalidated, **Then** execution blocks or requires fresh governed scope rather than acting on stale identity.

**Given** quarantine is observed, **When** outcome is handed off, **Then** threat status remains unknown unless Investigate separately determines it.

## 26. Questions ouvertes
OPEN-008/013/014/015 remain open; quarantine representation and Artifact/Attachment relation are not finalized.

## 27. Consommateurs documentaires
EPT-5 verification/release/reversal, Govern, Investigate, Settings, Security, Quality.