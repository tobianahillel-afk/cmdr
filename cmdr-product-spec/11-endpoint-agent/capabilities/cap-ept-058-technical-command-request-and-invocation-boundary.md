---
id: CAP-EPT-058
title: Technical Command Request and Invocation Boundary
product: endpoint-agent
module: live-response
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-006, REQ-PROD-014, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-015, OPEN-017]
source-of-truth: canonical
---
# CAP-EPT-058 — Technical Command Request and Invocation Boundary

## 1. Définition
Recevoir une demande d’exécution technique bornée dans une session autorisée, conserver origin/session/operator/target/purpose/mode/parameter references, masquer sensitive values, vérifier capability/policy/authority et produire accepted/rejected avant toute invocation.

## 2. Problème utilisateur
Une chaîne de commande ou un clic UI peut être pris pour une exécution réelle ou un Tool Call/Govern Decision. Les gates et l’autorité doivent être visibles avant invocation.

## 3. Objectifs
Distinguer request/accepted/invoked ; conserver exact origin/authority refs ; vérifier target/session/capability/mode ; traiter parameters sans raw secrets ; classifier impact ; rejeter unsupported/effectful operations sans authority.

## 4. Non-objectifs
Aucun catalogue final de commandes, syntaxe, shell, protocol, secret reveal, Tool Call, Decision/Approval/Response Run creation, containment primitive or EPT-5 effect.

## 5. Propriétaire
Endpoint owns technical request intake/invocation facts. Investigate owns business operation context, Govern authority, Studio Tool Call, Settings Secret References/config.

## 6. Utilisateurs
Response Operator, Endpoint Operator, Investigation Lead, Govern Reviewer, Security Reviewer, Auditor.

## 7. Conditions d’entrée
Active/eligible technical session or explicitly permitted non-session context, operation purpose and mode, target, parameter refs, capability/policy and authority snapshot.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| session/origin | CAP-EPT-056/057 or authorized handoff | execution context | oui | current | reject |
| operator/target/purpose | caller/session | attribution | oui | request time | incomplete |
| execution mode/parameter refs | caller/catalogue concept | technical intent | oui | exact request | incomplete |
| capability/policy | Endpoint/Settings | eligibility | oui | current | unsupported/blocked |
| authority/Secret References | Govern/Settings | authority/input refs | according class | current | reject/awaiting-authority |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Technical Session | Endpoint | target/state/operator | read |
| Endpoint Operation business context | Investigate | purpose/class | read ref |
| Decision/Response Run | Govern | authority ref | read only |
| Tool Call/Automation Run | Studio | origin ref | read only |
| Secret Reference | Settings | reference metadata only | restricted read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Technical Command Request | create/version | Endpoint | request != execution |
| Invocation Eligibility | derive | Endpoint | capability/policy/authority reasons |
| Invocation Attempt Ref | create only after acceptance | Endpoint | exact request/version |

## 11. Fonctionnalités
Validate session/target, parse only conceptual requested parameters, classify sensitivity/effect, resolve only approved Secret References at executor boundary without exposing values, check capability/policy/authority, accept/reject, correlate invocation attempt.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect request | Operator/Reviewer | command request | 0 | read | intent/gates visible | non |
| validate eligibility | service/reviewer | request | 1 | facts | pass/block | non |
| invoke non-mutating diagnostic request | authorized operator | request | 2 | session + authority as required | invocation-attempt | selon policy |
| invoke effectful request | authorized Govern path | request | 3 | Decision/authority | invocation-attempt | obligatoire |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| validate request/gates | oui | oui | oui | explain | validators |
| mask sensitive params | oui | oui | oui | non needed | classification rules |
| suggest diagnostic action | oui | catalogue/rules | oui | suggestion | manual catalogue |
| authorize/invoke autonomously | non | authority contract | non autonome | interdit | operator/Govern path |

## 14. États fonctionnels
`draft`, `validation-required`, `eligible`, `ineligible`, `awaiting-authority`, `accepted`, `rejected`, `invocation-requested`, `superseded`.

## 15. États d’interface
No Screen ID. Accepted != started; sensitive values remain masked; authority/eligibility shown separately.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Technical Command Request | Endpoint concept | CAP-EPT-059/062 | exact intent/version |
| acceptance/rejection | state | caller/Govern | reasons explicit |
| invocation attempt ref | Endpoint ref | CAP-EPT-062/064 | no Tool Call/Response Run identity merge |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| session/caller | operation requested | CAP-EPT-058 | origin/target/operator/authority | source |
| CAP-EPT-058 accepted | interactive command | CAP-EPT-059 | request/mode/refs | same session |
| CAP-EPT-058 accepted | technical execution | CAP-EPT-062 | attempt correlation | source |

## 18. Dépendances
CAP-EPT-056/057/059/062/064, Investigate CAP-INV-210, Govern CAP-GOV-025, Studio Tool Call, Settings Secret References, OPEN-008/013/015/017.

## 19. Source de vérité
Endpoint SOT of technical request/acceptance/invocation attempt; Govern SOT authority; Studio SOT Tool Call; Settings SOT Secret Reference.

## 20. Provenance et audit
Origin/session, operator, target, purpose, execution mode, masked parameter refs, capability/policy versions, authority refs, acceptance/rejection, timestamps and correlation.

## 21. Permissions fonctionnelles
Command request/read/invoke non-mutating/effectful, sensitive parameter metadata, Secret Reference use at executor boundary, cross-tenant deny; no final RBAC.

## 22. Limites et erreurs
Command request ≠ execution/Tool Call/Decision; accepted ≠ started; capability ≠ authority; raw secret never stored/displayed; no final command syntax/runtime.

## 23. Métriques
Accepted/rejected/authority-missing, masked sensitive inputs, unsupported modes, requests incorrectly merged with Tool Call/Response Run target zero.

## 24. Classification de livraison
`draft / defined / planned`; no command catalog/executor protocol/shell implementation.

## 25. Critères d’acceptation
**Given** a command request lacks required authority, **When** validated, **Then** it remains awaiting-authority/rejected and no invocation occurs.

**Given** a sensitive parameter uses a Secret Reference, **When** request is inspected, **Then** raw secret is never displayed or persisted in request/audit.

**Given** a Studio Tool Call originated the request, **When** invocation begins, **Then** Tool Call remains Studio-owned and technical execution remains Endpoint-owned.

## 26. Questions ouvertes
OPEN-008/013/015/017 remain open; no universal shell/runtime or bridge selected.

## 27. Consommateurs documentaires
CAP-EPT-059/062/064, Investigate, Govern, Studio, Settings, Security, Quality.
