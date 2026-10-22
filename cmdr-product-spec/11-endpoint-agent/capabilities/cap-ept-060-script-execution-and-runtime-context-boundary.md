---
id: CAP-EPT-060
title: Script Execution and Runtime Context Boundary
product: endpoint-agent
module: live-response
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-014, REQ-PROD-018, REQ-PROD-019, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008, OPEN-013, OPEN-014, OPEN-015, OPEN-017]
source-of-truth: canonical
---
# CAP-EPT-060 — Script Execution and Runtime Context Boundary

## 1. Définition
Définir l’exécution d’un `Script Reference` autorisé : exact source/version/hash-like reference lorsqu’elle existe, runtime/interpreter dependency, parameters, eligibility/authority, execution context/state/output/failure/cleanup-provenance, sans fournir script ni choisir runtime universel.

## 2. Problème utilisateur
Un script signé/référencé peut être incompatible, unsafe pour le contexte ou manquer d’autorité. `script uploaded/accepted` ne doit jamais signifier executed/safe/successful.

## 3. Objectifs
Conserver source/version/reference ; valider runtime/platform/capability ; appliquer authority/policy ; préserver sensitive params via refs ; enregistrer exact execution attempt/output/errors ; distinguer temporary transfer et execution.

## 4. Non-objectifs
Aucun script fourni, PowerShell/Bash/Python imposé, package format, signing implementation, deployment, persistence, bypass, containment ou automatic safety verdict.

## 5. Propriétaire
Endpoint owns target-side script execution facts. Source script/Artifact identity remains its owner under OPEN-014; Govern owns authority; Settings runtime/admin/Secret References.

## 6. Utilisateurs
Response Operator, Endpoint Operator, Automation/Studio caller as source, Govern Reviewer, Security Reviewer.

## 7. Conditions d’entrée
Authorized Script Reference/version, declared runtime/capability, target/session, parameters/Secret References, policy/authority, expected cleanup semantics if temporary material is used.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Script Reference/version | Investigate/Studio/authorized source | content ref | oui | exact version | reject |
| runtime requirement | script metadata | dependency | oui | versioned | unsupported |
| target/session | CAP-EPT-056/057 | execution context | oui | current | blocked |
| parameters/Secret References | caller/Settings | inputs | conditionnel | exact request | incomplete |
| policy/authority | Settings/Govern | gate | oui selon effect | current | blocked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Script Reference / source file | source owner | id/version/hash-like ref | read/use according permission |
| Technical Session | Endpoint | target/state | read |
| Secret Reference | Settings | reference only | restricted use |
| Decision/Response Run | Govern | authority refs | read only |
| Tool Call/Automation Run | Studio | origin refs | read only |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Script Execution Request | create/version | Endpoint | exact script/runtime/params refs |
| Technical Execution Attempt | create/update | Endpoint | distinct from Tool/Response Run |
| Cleanup Status Ref | record | Endpoint | does not imply rollback/state restoration |

## 11. Fonctionnalités
Validate script source/version/runtime/platform, sensitive parameters and authority; accept/reject; start execution only after confirmation; capture output/errors/timeout/cancel; record cleanup request/result where applicable without declaring host state restored.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect script/runtime eligibility | Operator | script request | 0 | read | reasons visible | non |
| execute diagnostic script | authorized operator | attempt | 2 | runtime + permission + authority as required | execution | according effect |
| execute effectful script | Govern-authorized path | attempt | 3 | Decision/authority | execution | obligatoire |
| cancel | Operator | attempt | 2 | interruptible | cancel-requested | according class |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| validate version/runtime | oui | oui | oui | explain | metadata checks |
| validate parameter schema concept | oui | oui | oui | assist | deterministic validation |
| summarize output/errors | oui | oui | oui | oui | raw output |
| generate/alter/execute script autonomously | non | explicit governed path | non autonome | interdit | approved source reference |

## 14. États fonctionnels
`draft`, `validating`, `runtime-unavailable`, `awaiting-authority`, `accepted`, `start-requested`, `running`, `partial`, `completed`, `failed`, `timed-out`, `cancel-requested`, `cancelled`, `cleanup-pending`, `cleanup-unknown`.

## 15. États d’interface
No Screen ID. Script reference/version/runtime and authority remain visible; sensitive values masked; accepted != executed.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Script Execution Attempt | Endpoint | CAP-EPT-062/064 | exact source/runtime refs |
| technical output/error | Endpoint output | Investigate/Govern reconciliation | not Result/Evidence automatically |
| cleanup status | Endpoint fact | Audit | cleanup != rollback/state restored |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| script source/caller | request | CAP-EPT-060 | script/version/runtime/params/authority | source |
| CAP-EPT-060 | execution state/output | CAP-EPT-062 | attempt/output/errors | session |
| CAP-EPT-060 | audit/provenance | CAP-EPT-063/064 | source/runtime/controls | source |

## 18. Dépendances
Endpoint script-execution source, CAP-EPT-056/057/061/062/064, Investigate CAP-INV-210/211, Settings Secrets, Govern, Studio, OPEN-008/013/014/015/017.

## 19. Source de vérité
Endpoint SOT for execution facts; source owner SOT for script identity/version; Settings SOT Secret References/runtime admin; Govern SOT authority.

## 20. Provenance et audit
Script ref/version/hash-like metadata, source owner, session/target/operator, runtime dependency, masked params/secret refs, policy/authority, start/end/errors/output/cleanup and correlation.

## 21. Permissions fonctionnelles
Script reference read/use, script execute diagnostic/effectful, sensitive parameter/output, runtime capability, cancel, provenance, cross-tenant deny.

## 22. Limites et erreurs
Script uploaded ≠ executed; accepted ≠ safe; script success ≠ authorized result; runtime unavailable remains unavailable; cleanup ≠ rollback; no universal interpreter.

## 23. Métriques
Runtime-unavailable, validation/authority blocks, execution partial/failure/timeout/cancel, cleanup-unknown, provenance completeness.

## 24. Classification de livraison
`draft / defined / planned`; no script, interpreter, signing, transport or executor implementation selected.

## 25. Critères d’acceptation
**Given** runtime unavailable, **When** script execution is requested, **Then** request is blocked without substituting another universal runtime.

**Given** script requires stronger authority, **When** diagnostic operator submits it, **Then** it remains awaiting-authority and does not execute.

**Given** a script succeeds technically, **When** result is handed to Govern, **Then** it remains technical output and Govern separately reconciles canonical Result.

## 26. Questions ouvertes
OPEN-008/013/014/015/017 remain open.

## 27. Consommateurs documentaires
CAP-EPT-062..064, Investigate/Studio/Govern/Settings, Security/Quality.
