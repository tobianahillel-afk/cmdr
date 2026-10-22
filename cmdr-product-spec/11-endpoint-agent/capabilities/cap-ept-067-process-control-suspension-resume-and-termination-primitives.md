---
id: CAP-EPT-067
title: Process Control, Suspension, Resume and Termination Primitives
product: endpoint-agent
module: containment
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-006, REQ-PROD-016, REQ-PROD-018, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-EPT-067 — Process Control, Suspension, Resume and Termination Primitives

## 1. Définition
Définir les primitives techniques Endpoint de suspension, reprise et terminaison d’un processus ou scope d’arbre explicitement autorisé, avec identité de cible, état d’exécution, résultat technique et vérification locale.

## 2. Problème utilisateur
Un PID ou process reference peut disparaître ou être réutilisé entre demande et exécution. Une réponse `success` ne prouve pas que le processus visé a effectivement atteint l’état souhaité.

## 3. Objectifs
Résoudre une identité de processus fraîche ; distinguer suspend/resume/terminate ; borner tree scope ; exiger precheck/authority ; conserver requested/started/confirmed states ; produire target-state observation et provenance.

## 4. Non-objectifs
Fournir une commande OS, déterminer qu’un process est malveillant, conclure incident resolved, tuer automatiquement un arbre non autorisé, créer Response Run/Result ou définir support plateforme.

## 5. Propriétaire
Endpoint possède la primitive et les faits process target-side. Govern possède authority/Decision/Run/verification métier. Investigate possède Finding/Evidence ; EPT-3 possède inspection context.

## 6. Utilisateurs
Response Operator, Endpoint Operator, Govern Reviewer, Verification Reviewer, DFIR Analyst, Auditor.

## 7. Conditions d’entrée
CAP-EPT-065/066 ready, exact process reference plus freshness/identity evidence, requested operation, authority current, target Agent reachable enough for the requested primitive.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| process ref + identity facts | EPT-3/EPT-4 | target | oui | fresh precheck | stale/unknown |
| operation suspend/resume/terminate | Govern handoff | requested effect | oui | exact Step | reject |
| tree-scope bound | Govern/technical plan | scope | si tree action | pinned | no expansion |
| primitive readiness | CAP-EPT-066 | precondition | oui | current | blocked |
| authority ref | Govern | authorization | oui | current/effective | authority-missing |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Response Run/Step/Decision | Govern | operation/scope/authority | read/ref |
| process observation/context | Endpoint | PID/identity/start-time/tree/current state | read |
| Endpoint Policy | Settings | process-control restriction | read only |
| Finding/Case | Investigate | provenance/purpose only | linked read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Process Control Execution | create/transition | Endpoint | target identity pinned |
| Process Control State | observe/update | Endpoint | requested ≠ confirmed |
| Technical Outcome | create | Endpoint | raw technical fact only |
| process | suspend/resume/terminate effect | target system | only after Govern-authorized handoff |

## 11. Fonctionnalités
Revalidate process identity; reject stale/reused identity; request bounded suspend/resume/terminate; track accepted/start/partial/fail/unknown; preserve already-exited/inaccessible; verify target state after effect; retain tree-scope and per-process outcome.

## 12. Actions utilisateur
Inspect = Class 0. Revalidate identity = Class 1. Process suspend/resume/terminate = Class 3 by default and always Govern-dependent. No arbitrary Class 2 promotion under OPEN-013.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| validate identity/scope | oui | oui | oui | explain | identity diff |
| observe state after effect | oui | oui | oui | summarize | local observation |
| propose candidate primitive | oui | catalog | oui | oui | manual selection |
| authorize/execute autonomously | non | contract-controlled | non autonome | interdit | Govern + operator |

## 14. États fonctionnels
`requested`, `identity-check`, `stale-identity`, `already-exited`, `ready`, `executing`, `suspend-requested`, `suspended-observed`, `resume-requested`, `running-observed`, `terminate-requested`, `terminated-observed`, `partial`, `failed`, `inaccessible`, `unknown`, `cancel-requested`.

## 15. États d’interface
Aucun Screen ID. PID/identity freshness, requested vs observed state, partial tree results and unknown termination state remain visible conceptually.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Process Control Execution | technical record | CAP-EPT-073/079 | exact target/op/state |
| process target-state observation | technical fact | CAP-EPT-074/075 | timestamp/source |
| per-target/tree outcome | technical outcome | Govern reconciliation | partiality preserved |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-066 | process control ready | CAP-EPT-067 | process identity/op/scope/authority | Govern Run |
| CAP-EPT-067 | execution state | CAP-EPT-073 | technical execution/outcome | same primitive |
| CAP-EPT-067 | effect confirmation needed | CAP-EPT-074/075 | expected process state + observations | same Run |

## 18. Dépendances
CAP-EPT-037/050/058/062/065/066/073..080, CAP-GOV-025..032, Settings Policy, Security authority, OPEN-008/013/015.

## 19. Source de vérité
Endpoint is SOT of process primitive execution and target observations. Govern is SOT of authorization and response outcome.

## 20. Provenance et audit
Process identity/PID/start context/tree, operation, precheck, authority/Run/Step, operator, Agent/version, requested/observed states, timestamps, errors, partials, verification refs and correlation id.

## 21. Permissions fonctionnelles
Process control request/read, suspend, resume, terminate, sensitive process context read, verification read, cross-tenant deny; step-up/SoD/Govern dependency explicit, no final RBAC.

## 22. Limites et erreurs
Terminate request ≠ terminated; process absent ≠ threat resolved; suspended ≠ terminated; resume ≠ rollback; PID reused/stale ≠ same process; technical success ≠ Response Run success.

## 23. Métriques
Stale identities, already-exited, suspend/resume/terminate technical outcomes, partial tree controls, verification mismatch, unauthorized attempts target zero.

## 24. Classification de livraison
`draft / defined / planned`; no OS command, process-kill syntax, API, driver or implementation selected.

## 25. Critères d’acceptation
**Given** le processus disparaît avant terminate, **When** l’exécution commence, **Then** `already-exited` est enregistré et aucun autre PID n’est ciblé.

**Given** l’identité process a changé malgré le même PID, **When** precheck revalide, **Then** la primitive est bloquée comme stale/mismatch.

**Given** terminate succeeds technically, **When** l’état absent est observé, **Then** le technical outcome peut être succeeded mais l’incident n’est pas déclaré résolu.

## 26. Questions ouvertes
OPEN-008/013/015 remain open; no supported OS/process API or default low-risk governance selected.

## 27. Consommateurs documentaires
EPT-5 verification/reversal/provenance, Govern Runs & Rollback, Investigate/Command consumers, Security, Quality.