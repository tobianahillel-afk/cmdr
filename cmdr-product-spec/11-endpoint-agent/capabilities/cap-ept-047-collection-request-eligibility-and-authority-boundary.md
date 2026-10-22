---
id: CAP-EPT-047
title: Collection Request, Eligibility and Authority Boundary
product: endpoint-agent
module: collection
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-006, REQ-PROD-014, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-014, OPEN-015]
source-of-truth: canonical
---
# CAP-EPT-047 — Collection Request, Eligibility and Authority Boundary

## 1. Définition
Définir comment Endpoint reçoit une `Collection Request` Investigate-owned, vérifie l’éligibilité technique locale, les capabilities, la plateforme, les limites et une référence d’autorité, puis accepte ou rejette un handoff technique sans devenir propriétaire de la request ni de la Decision.

## 2. Problème utilisateur
Une demande métier peut être correctement formulée tout en étant techniquement impossible, non supportée ou non autorisée. Sans frontière explicite, `eligible`, `authorized`, `accepted` et `started` peuvent être confondus.

## 3. Objectifs
Préserver l’owner Investigate de la Collection Request ; calculer l’éligibilité locale ; exposer support/degradation et raisons ; consommer un authority reference sans l’inventer ; produire un technical intake corrélable ; arrêter avant tout démarrage si une précondition manque.

## 4. Non-objectifs
Aucune création de Decision/Approval/Response Run, aucune qualification Artifact/Evidence, aucun protocole, API, commande de collecte, transport, plateforme officiellement supportée, containment ou EPT-5.

## 5. Propriétaire
Investigate possède la Collection Request et son contexte Case. Endpoint possède seulement le `Collection Technical Intake` et l’`Eligibility Assessment` locaux. Govern possède l’autorité ; Settings Fleet/Policy ; Shared les Jobs génériques.

## 6. Utilisateurs
DFIR Analyst, Case Analyst, Response Operator, Endpoint Operator, Govern Reviewer, Security Reviewer et Auditor.

## 7. Conditions d’entrée
Collection Request/version résoluble, Endpoint/Agent tenant-bound, requested collection type/scope lisibles, capability state disponible ou explicitement unknown, policy/permission et authority context fournis selon classe.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Collection Request + version | Investigate | demande métier | oui | version soumise | reject `request-unresolved` |
| EPT-3 collection-required ref | CAP-EPT-043/044/045 | besoin contextuel | non | source freshness | aucune inférence |
| platform/capability/health | CAP-EPT-004/005/008/027/028 | faits techniques | oui | dernière observation | `eligibility-unknown` |
| policy/limits | Platform Settings | restriction | oui si applicable | snapshot courant | `policy-unknown`/block |
| authority reference | Govern/Security | autorité externe | selon classe | valide à l’action | `authority-missing` |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Collection Request | Investigate | id/version/target/scope/purpose/class | lecture référencée |
| endpoint-agent | Endpoint | identity/platform/health/capabilities | lecture locale |
| Endpoint Policy | Settings | restrictions/limits | lecture projection |
| Action Request/Decision/Response Run | Govern | authority/correlation refs | lecture seulement |
| Endpoint Investigation Summary | Endpoint | missing-context reason | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Collection Technical Intake | créer/actualiser | Endpoint | request ref seulement, aucune copie d’ownership |
| Collection Eligibility Assessment | dériver | Endpoint | raisons et dépendances explicites |
| Authority Reference Projection | attacher/rafraîchir | Endpoint | aucune Decision créée |

## 11. Fonctionnalités
Résoudre request/target, vérifier type/scope demandé, capability/support/degradation, contraintes platform/version, policy et authority ref, puis produire `eligible`, `ineligible`, `unknown`, `blocked` ou `accepted-for-technical-planning`. `Accepted` ne signifie jamais `started`.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect request/eligibility | Analyst/Operator | Technical Intake | 0 | read | raisons visibles | non |
| recompute eligibility | Endpoint service | Eligibility Assessment | 1 | facts courants | assessment versionné | non |
| accept bounded technical intake | Endpoint Operator/service | Technical Intake | 2 | eligible + authority ref si requise | planning allowed | selon classe source |
| bypass missing authority | aucun | request | — | interdit | no effect | obligatoire si requis |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| vérifier capability/support | oui | oui | oui | explication seulement | règles déclarées |
| vérifier completeness | oui | oui | oui | oui | checklist |
| résumer raisons | oui | oui | oui | oui, attribuée | reason codes |
| autoriser/démarrer collecte | non autonome | authority contract | non autonome | interdit | chemin Govern/opérateur |

## 14. États fonctionnels
`received`, `request-unresolved`, `validating`, `eligible`, `ineligible`, `eligibility-unknown`, `policy-blocked`, `authority-missing`, `accepted-for-technical-planning`, `rejected`, `stale`.

## 15. États d’interface
Aucun Screen ID. Loading/Partial/Stale/Offline/Permission denied/Unsupported doivent préserver la request source, la dernière assessment valide et la distinction `eligible != authorized`.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Eligibility Assessment | projection Endpoint | Investigate/CAP-EPT-048 | reason/source/version |
| Technical Intake | concept Endpoint | CAP-EPT-048/053 | request ref + authority context |
| rejection/block reason | diagnostic | Investigate/Govern | aucun faux démarrage |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| EPT-3 collection-required | request préparée | Investigate Collection Request | need/context refs | EPT-3 origin |
| Collection Request | handoff reçu | CAP-EPT-047 | request/version/target/scope | Investigate owner retained |
| CAP-EPT-047 | eligible + bounds | CAP-EPT-048 | intake/authority/capabilities | request origin retained |
| blocked/ineligible | assessment terminal | Investigate/Govern | reasons/limitations | no operation |

## 18. Dépendances
CAP-EPT-001..045 ; Investigate CAP-INV-202/203 ; Settings Endpoint Policy/Fleet ; Govern authority/execution handoff ; Shared Jobs/Trace ; OPEN-008/013/014/015.

## 19. Source de vérité
Investigate reste SOT de Collection Request ; Endpoint SOT de l’éligibilité technique locale ; Govern SOT de l’autorité ; Settings SOT de policy/admin configuration.

## 20. Provenance et audit
Request id/version, Case ref, Endpoint/Agent, requester, purpose, requested type/scope, platform/capability facts, policy/version, authority ref, assessment reasons, actor/service, timestamps et correlation ID.

## 21. Permissions fonctionnelles
Request read, local eligibility assess, capability/policy read, authority-reference read, sensitive scope metadata, cross-tenant deny. Aucun RBAC/ABAC final.

## 22. Limites et erreurs
Eligible ≠ authorized ; authorized ≠ started ; Request ≠ operation/Decision/Response Run ; capability available ≠ authority ; offline/stale/unsupported restent explicites ; aucun fallback ne lance une collecte.

## 23. Métriques
Assessments par état/reason, authority-missing, unsupported/degraded, stale assessments, cross-tenant denials, requests acceptées sans start ; aucun SLO numérique.

## 24. Classification de livraison
`draft / defined / planned`. Contrat documentaire seulement ; aucun collector, protocole, transport ou support plateforme livré.

## 25. Critères d’acceptation
**Given** une Collection Request sans référence d’autorité requise, **When** l’éligibilité est évaluée, **Then** elle peut être techniquement eligible mais reste `authority-missing` et aucune opération ne démarre.

**Given** une capability de collecte indisponible, **When** la request arrive, **Then** l’état est ineligible/unsupported avec raison et aucune collecte n’est simulée.

**Given** l’IA est indisponible, **When** la request est évaluée, **Then** capability, policy, authority et reasons sont calculés par règles déterministes.

## 26. Questions ouvertes
OPEN-008/013/014/015 restent ouvertes. Ce contrat ne décide ni Artifact/Attachment/Evidence ni bridge Run final.

## 27. Consommateurs documentaires
EPT-4 Collection, Investigate Collection and Live Response, Govern, Settings, Shared, Security, Quality, Roadmap Phase 5 et EPT-5 futur comme frontière uniquement.
