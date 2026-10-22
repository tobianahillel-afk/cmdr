---
id: CAP-EPT-056
title: Live Response Session Definition, Eligibility and Target Binding
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
# CAP-EPT-056 — Live Response Session Definition, Eligibility and Target Binding

## 1. Définition
Définir la sémantique technique locale d’une Live Response Session : target binding, capability/eligibility, requester/initiator/operator refs, contraintes temporelles et authority context, sans devenir le record métier Investigate ni un Govern Response Run.

## 2. Problème utilisateur
Une session peut être demandée alors que l’Endpoint est offline, unsupported, policy-blocked ou non autorisé. Une simple session ouverte peut aussi être prise à tort pour une autorité de commande illimitée.

## 3. Objectifs
Lier exactement l’Agent/Endpoint/tenant/environment ; vérifier capability/platform/policy ; conserver requester/initiator/operator et purpose ; consommer authority refs ; exposer contraintes et eligibility ; séparer session technique, Live Session métier, Automation Run et Response Run.

## 4. Non-objectifs
Aucun protocole remote-shell, PKI, transport, terminal UX, shell universel, commande, script, containment, Decision, Approval, Response Run ou Result.

## 5. Propriétaire
Endpoint possède la session technique côté Agent. Investigate conserve le contexte/record métier Live Session et Case ; Govern l’autorité/Response Run ; Studio Automation Run/Tool Call ; Settings policy/configuration.

## 6. Utilisateurs
Response Operator, Endpoint Operator, Case Analyst, Investigation Lead, Govern Reviewer, Security Reviewer et Auditor.

## 7. Conditions d’entrée
Request origin traçable, target Endpoint/Agent résolu, purpose et session constraints explicites, capability Live Response déclarée, policy/permission et authority context disponibles selon action.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Live Session/request ref | Investigate/Govern/authorized caller | contexte métier | oui | version courante | no technical session |
| Endpoint/Agent binding | EPT-1 | cible | oui | current binding | reject |
| capability/support/health | EPT-1/EPT-2 | faisabilité | oui | last seen/current | unavailable/unknown |
| policy/session limits | Settings | restriction | oui | current | block/unknown |
| authority reference | Govern/Security | authority | selon operation class | valid at open | awaiting-authority |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Live Session business ref | Investigate | case/purpose/participants | référence seulement |
| endpoint-agent | Endpoint | identity/health/capabilities | local read |
| Endpoint Policy | Settings | session duration/restrictions | read projection |
| Response Run/Decision | Govern | authority/correlation | read only |
| Automation Run/Tool Call | Studio | caller provenance éventuelle | read only |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Endpoint Technical Session | créer/actualiser | Endpoint | distinct du Live Session métier |
| Session Eligibility Assessment | dériver | Endpoint | support/policy/reasons |
| Target Binding Record | dériver/versionner | Endpoint | exact Agent/Endpoint/tenant/env |

## 11. Fonctionnalités
Valider target/capability/support, policy et authority reference ; créer un Session ID conceptuel ; attacher requester/initiator/operator/purpose/tenant/environment ; appliquer expiration/inactivity constraints conceptuelles ; rejeter mismatch/cross-tenant/unsupported.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect session eligibility | Operator | Eligibility Assessment | 0 | read | reasons visible | non |
| validate target binding | Endpoint service | binding | 1 | stable refs | pass/block | non |
| request/open authorized diagnostic session | Response/Endpoint Operator | Technical Session | 2 | capability + permission + authority as required | requested/opening | OPEN-013/015 selon contexte |
| open cross-tenant/unbounded session | aucun | session | — | interdit | denied | obligatoire |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| resolve target/capability | oui | oui | oui | explanation only | deterministic refs |
| validate policy/constraints | oui | oui | oui | oui | checklist/rules |
| summarize purpose/context | oui | oui | oui | oui, attributed | structured fields |
| authorize/open autonomously | non | contract only | non autonome | interdit | explicit authority path |

## 14. États fonctionnels
`requested`, `validating`, `eligible`, `ineligible`, `awaiting-authority`, `target-mismatch`, `unsupported`, `unavailable`, `opening`, `open-rejected`.

## 15. États d’interface
Aucun Screen ID. Loading/Offline/Stale/Permission denied/Unavailable montrent last-known target, eligibility et authority status sans représenter une session comme active avant confirmation.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Endpoint Technical Session ref | concept Endpoint | CAP-EPT-057..063 | exact target/constraints |
| eligibility/binding result | projection | Investigate/Govern | source/reason/time |
| rejected/blocked reason | diagnostic | caller | no hidden authority |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Investigate/Govern request | handoff | CAP-EPT-056 | origin/target/purpose/authority | source origin |
| CAP-EPT-056 | eligible/open requested | CAP-EPT-057 | session/binding/constraints | request retained |
| blocked | eligibility terminal | Investigate/Govern | reasons/limitations | no session active |

## 18. Dépendances
CAP-EPT-001..030/045/046 ; Investigate CAP-INV-209 ; Govern CAP-GOV-025/execution boundaries ; Settings Policy/Secrets ; Studio run provenance ; OPEN-008/013/015/017.

## 19. Source de vérité
Endpoint est SOT de la session technique/target-side connectivity state. Investigate reste SOT du record métier de session ; Govern de l’autorité ; Studio de ses Runs/Calls.

## 20. Provenance et audit
Session ref, origin, target Agent/Endpoint, tenant/environment, requester/initiator/operator, purpose, capability/platform state, policy/version, authority refs, constraints, assessment/open timestamps et correlation IDs.

## 21. Permissions fonctionnelles
Technical session read/request/open/close needs, target/capability read, restricted session metadata, authority refs, cross-tenant deny. Aucun final RBAC/ABAC.

## 22. Limites et erreurs
Session ≠ Response Run/Automation Run/unrestricted shell ; eligible ≠ authorized ; session available ≠ command allowed ; session open ≠ command authority ; unsupported/offline/denied restent explicites.

## 23. Métriques
Eligibility distribution, open rejections, target mismatch, authority-missing, cross-tenant denials, unsupported/unavailable rates, provenance completeness.

## 24. Classification de livraison
`draft / defined / planned`; aucun protocole, terminal, shell, runtime universel ou transport implémenté.

## 25. Critères d’acceptation
**Given** une session techniquement eligible mais sans authority requise, **When** l’ouverture est demandée, **Then** elle reste `awaiting-authority` et aucune session active n’est déclarée.

**Given** le target appartient à un autre tenant, **When** le binding est validé, **Then** la session est refusée sans fuite de contexte.

**Given** l’IA est indisponible, **When** l’éligibilité/session binding est évalué, **Then** le chemin déterministe fournit le résultat complet.

## 26. Questions ouvertes
OPEN-008/013/015/017 restent ouvertes. Aucun bridge Run, support universel ou runtime final n’est choisi.

## 27. Consommateurs documentaires
EPT-4 Live Response, Investigate Collection/Live Response, Govern, Studio, Settings, Shared Trace, Security, Quality et future EPT-5 comme frontière seulement.
