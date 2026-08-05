---
id: CAP-INV-376
title: Encrypted, Compressed and Restricted Content Handling
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
  - REQ-SEC-001
  - REQ-SEC-002
  - REQ-AI-002
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-376 — Encrypted, Compressed and Restricted Content Handling

## 1. Définition
Identifier et gérer fonctionnellement un contenu candidat chiffré, compressé ou restreint, conserver restrictions et tentatives autorisées, et produire un Derived Artifact uniquement après accès permis, sans cassage, bypass ou attaque de mot de passe.

## 2. Problème utilisateur
Un contenu inaccessible peut être pris pour absent ou malveillant, tandis qu’un accès non gouverné peut exposer des données privées ou secrets.

## 3. Objectifs
Show candidate type/container/available information/restrictions/permissions/declared authorized methods; mark inaccessible; request information/authorization; use only authorized means and bounded attempts; record errors/attempts; create Derived Artifact only after authorized access; preserve original; annotate/dispute/handoff.

## 4. Non-objectifs
No cracking, password suggestion, automated attack, protection bypass, secret reuse without authorization, decryption algorithm, command, tool mandate or malicious inference.

## 5. Propriétaire
Investigate owns Restricted Content Record and analysis context; Security/Settings own policy/authorization/secret references; Studio owns Tool/Call/Run; Shared owns audit/export.

## 6. Utilisateurs
Principal : Authorized Forensics Analyst. Secondaires : Security Reviewer, Privacy Reviewer, Artifact Analyst and Audit Analyst.

## 7. Conditions d’entrée
Source/container identified; restrictions and permissions evaluated; any authorized information represented by reference, not exposed value; original immutable.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Content/container candidate | CAP-INV-368/309 | source/type | yes | source version | blocked |
| Restriction/protection result | Tool Call | encrypted/compressed/restricted status | yes | Tool/version | unknown |
| Policy/permission/authorization | Security/Settings/Govern | allowed actions | yes | current | restricted |
| Authorized access information | approved source/reference | permitted input | no | validity visible | inaccessible |
| Attempt history | Investigate/Trace | prior authorized attempts/errors | no | append-only | none |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Artifact/File Observation | Investigate | source/container/relations | read |
| Restricted Content Record | Investigate concept | status/restrictions/attempts | read |
| Policy/Secret Reference | Settings/Security | authorization projection | reference only |
| Tool/Tool Call/Run | Studio | method/version/result | read/select |
| Decision/Approval | Govern | sensitive authority if required | read/link |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Restricted Content Record | create/update/dispute | Investigate concept | inaccessible ≠ absent |
| Access request/authorization relation | prepare/link | owner capability | no self-approval |
| Attempt event | append | Shared/Investigate | actor/reason/result required |
| Derived Artifact | create after authorized access | Investigate | parent/method/restrictions retained |

## 11. Fonctionnalités
Detect candidate protection, show container/info/restrictions/permissions/authorized methods, indicate inaccessible, request authorization/information, limit attempts, preserve attempt/errors, produce Derived Artifact only after authorization, preserve original, annotate/dispute/handoff.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspect status | Analyst | Restricted Content | 0 | read | protection/restrictions visible | no |
| Request access/information | Analyst | request | 2 | reason/scope | owner handoff | policy/Govern if required |
| Attempt authorized access | Authorized Analyst | Tool Call | 1 | valid authorization/bounds | attributed result | according policy |
| Annotate/dispute | Reviewer | record | 2 | permission | versioned disposition | OPEN-013 |
| Export derived content | Analyst | Derived Artifact | 1/2 | export policy | controlled output | policy |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| detect protection/container | yes | Tool/rules | yes | suggestion | deterministic inspection |
| enforce limits/policy | yes | yes | yes | no | policy engine/checklist |
| explain error | yes | diagnostics | yes | summary | error catalog |
| choose/guess password or bypass | no | prohibited | no | prohibited | request authorized information |

## 14. États fonctionnels
`detected`, `compressed`, `encrypted`, `restricted`, `access-requested`, `authorized`, `processing`, `accessible`, `inaccessible`, `partial`, `failed`, `disputed`, `superseded`.

## 15. États d’interface
Loading preserves source; Empty is not absence; Partial shows accessible members; Error keeps valid outputs; Offline blocks new attempt; Permission denied hides protected information; Stale shows authorization/version expiry.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Restricted Content Record | concept | Workbench/Security | status/restrictions/attempt history |
| Access request/result | relation/event | owner capability/Audit | no self-permission or hidden attempt |
| Authorized Derived Artifact | Artifact | Static/Reverse/CAP-INV-379 | parent/method/policy lineage |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| File/container | detect restriction | CAP-INV-376 | source, container, protection, permissions | File view |
| Restricted record | request authority | Security/Settings/Govern | reason, scope, source, policy | Restricted view |
| Accessible output | analyze | Static/Reverse/CAP-INV-379 | Derived Artifact, parent, authorization, restrictions | Restricted view |

## 18. Dépendances
CAP-INV-309/311/368/375/378/379, Settings/Security/Govern, Studio Tools, Shared Audit/Export, OPEN-005/008/013/014/015.

## 19. Source de vérité
Investigate owns record/interpretation; policy/authorization/secret references remain owner sources; no value copied into documentation or model without permission.

## 20. Provenance et audit
Source/container, protection status, policy/version, request, authorization, reference used, Tool/version/Call/Run, each bounded attempt, errors, access result, Derived Artifact, actor and export.

## 21. Permissions fonctionnelles
Restricted content presence/read, access request, authorized method use, sensitive content read, Derived Artifact create/export and audit review. Reveal/copy/export/step-up/SoD remain separate and future.

## 22. Limites et erreurs
No cracking/bypass/password attack or secret reuse. Inaccessible ≠ absent; encrypted ≠ malicious. Expired authorization, unsupported protection, malformed container, attempt limit, permission denial and Tool failure are explicit.

## 23. Métriques
Restricted/encrypted/compressed records, access requests approved/denied, bounded attempts, policy blocks, accessible/partial results and exports with lineage.

## 24. Classification de livraison
`defined` / `planned`; no algorithm, product, secret method, engine or implementation chosen.

## 25. Critères d’acceptation
**Given** encrypted content and no authorized access method **When** opened **Then** it remains inaccessible, no bypass/password attack is suggested or launched, a request may be prepared and attempts are traced.

**Given** compressed container with restricted member **When** inspected **Then** accessible structure and restricted member states remain distinct and no content executes.

**Given** no AI **When** handled **Then** deterministic detection, policy checks, request workflow and human review work.

## 26. Questions ouvertes
OPEN-005/008/013/014/015 remain open; support, authorization and permission rules final are future.

## 27. Consommateurs documentaires
INV-DSK-001, CAP-INV-368/375/378/379, Static/Reverse, Security/Settings/Govern and Objects/Permissions/Technique.
