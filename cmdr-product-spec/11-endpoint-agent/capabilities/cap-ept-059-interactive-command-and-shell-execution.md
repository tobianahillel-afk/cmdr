---
id: CAP-EPT-059
title: Interactive Command and Shell Execution
product: endpoint-agent
module: live-response
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-014, REQ-PROD-018, REQ-PROD-019, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-015, OPEN-017]
source-of-truth: canonical
---
# CAP-EPT-059 — Interactive Command and Shell Execution

## 1. Définition
Définir l’exécution interactive d’une Technical Command Request acceptée dans un contexte terminal/shell déclaré : execution start/state, user/privilege/working context, timeout, transcript/redaction, output/error/cancel/completion et provenance, sans choisir shell ni fournir commandes.

## 2. Problème utilisateur
Une session interactive peut être prise pour un shell illimité, et une commande acceptée pour une commande exécutée/successful. Le contexte de privilège et la redaction doivent rester auditables.

## 3. Objectifs
Conserver exact accepted request ; vérifier declared shell/runtime context ; enregistrer execution start confirmation ; préserver user/privilege/working context ; redacter transcript/output ; gérer timeout/cancel/error/completion.

## 4. Non-objectifs
Aucune commande offensive ou catalogue, aucun PowerShell/Bash/Python imposé, élévation de privilège, persistence, containment, bypass, remote-shell protocol ou unrestricted authority.

## 5. Propriétaire
Endpoint owns target-side interactive execution facts. Govern owns effectful authority/Response Run; Investigate owns business operation context; Studio owns Tool Call/Automation Run.

## 6. Utilisateurs
Response Operator, Endpoint Operator, Govern Reviewer, Security Reviewer, Auditor.

## 7. Conditions d’entrée
CAP-EPT-058 accepted request, active authorized session, declared interactive/runtime capability, user/privilege context known or explicit unknown, policy/authority current.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| accepted command request | CAP-EPT-058 | exact intent | oui | exact version | no execution |
| session state | CAP-EPT-057 | runtime context | oui | current | blocked |
| shell/runtime capability | Endpoint | declared support | oui | current | unsupported |
| user/privilege/working context | Endpoint/policy | execution context | oui when available | start time | unknown/restricted |
| authority | Govern/Security | effect authority | according action | current | blocked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Technical Command Request | Endpoint | mode/refs/authority | read |
| Technical Session | Endpoint | target/state/operator | read |
| Decision/Response Run | Govern | authority/correlation | read only |
| Endpoint Policy | Settings | runtime/restrictions | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Technical Execution Attempt | create/update | Endpoint | one accepted request/version |
| Interactive Transcript Ref | append/reference | Endpoint | redaction/classification applied |
| Runtime Context Record | derive | Endpoint | user/privilege/working context |

## 11. Fonctionnalités
Confirm execution start, bind runtime/user/privilege/working context, capture output/error/transcript refs, apply redaction, track running/partial/timeout/cancel/completion, preserve target-side unknown termination states.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect runtime context | Operator | attempt | 0 | read | context visible | non |
| execute non-mutating diagnostic command | authorized operator | attempt | 2 | accepted + authority as required | running | selon policy |
| execute effectful command | Govern-authorized path | attempt | 3 | Decision/authority | running | obligatoire |
| request cancel/stop | Operator | attempt | 2 | interruptible | request state | according class |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| runtime/precondition check | oui | oui | oui | explain | rules |
| redact transcript | oui | oui | oui | assist only | classification/masking |
| summarize output/errors | oui | oui | oui | oui | raw output |
| invent/execute command | non | explicit request only | non autonome | interdit | operator-approved input |

## 14. États fonctionnels
`accepted`, `start-requested`, `running`, `partial`, `completed`, `failed`, `timed-out`, `cancel-requested`, `cancelled`, `stop-requested`, `termination-unknown`, `runtime-unavailable`.

## 15. États d’interface
No terminal UI design. `start-requested != running`, `command success != Response Run success`, raw secrets/sensitive output masked.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Technical Execution Attempt state | Endpoint | CAP-EPT-062/063 | target-side facts only |
| transcript/output refs | Endpoint technical output | CAP-EPT-062/Investigate | not Govern Result/Evidence automatically |
| runtime context | metadata | Audit/Govern | user/privilege/source explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-058 | accepted interactive request | CAP-EPT-059 | request/session/authority | request |
| CAP-EPT-059 | output/state | CAP-EPT-062 | attempt/status/output/errors | same session |
| CAP-EPT-059 | operator control | CAP-EPT-063 | attempt/control/audit | same session |

## 18. Dépendances
Endpoint terminal/command-execution sources, CAP-EPT-057/058/062/063/064, Investigate CAP-INV-210, Govern boundary, OPEN-008/013/015/017.

## 19. Source de vérité
Endpoint SOT for local execution attempt/state/output; Govern SOT canonical Result; caller owner remains source of business intent.

## 20. Provenance et audit
Session/request/attempt, target, operator, runtime capability/version reference, user/privilege/working context, authority, start/end, output/error/transcript/redaction, controls and correlation.

## 21. Permissions fonctionnelles
Interactive execute, non-mutating/effectful differentiation, transcript/output read, sensitive output, cancel/stop, cross-tenant deny.

## 22. Limites et erreurs
Live Response Session ≠ unrestricted shell; session open ≠ command authorized; accepted ≠ started; started ≠ success; output ≠ Result/Evidence; stop request ≠ stopped.

## 23. Métriques
Start/failure/timeout/cancel, termination-unknown, redaction events, effectful commands with valid authority, provenance completeness.

## 24. Classification de livraison
`draft / defined / planned`; no shell/runtime/protocol/command catalog selected.

## 25. Critères d’acceptation
**Given** command accepted but runtime fails to start, **When** status is reconciled, **Then** attempt is failed/start-not-confirmed rather than running.

**Given** command times out, **When** no termination confirmation exists, **Then** timeout and target-side termination unknown are both visible.

**Given** a shell runtime is unavailable on platform, **When** execution is requested, **Then** no alternate universal shell is assumed.

## 26. Questions ouvertes
OPEN-008/013/015/017 remain open; no universal terminal/runtime.

## 27. Consommateurs documentaires
CAP-EPT-062/063/064, Investigate Live Response, Govern, Security, Quality.
