---
id: CAP-EPT-051
title: Memory Acquisition and Collection Boundary
product: endpoint-agent
module: collection
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-014, REQ-PROD-018, REQ-PROD-020, REQ-PROD-052, REQ-PROD-055, REQ-SEC-001, REQ-SEC-004]
open_decisions: [OPEN-005, OPEN-008, OPEN-013, OPEN-014]
source-of-truth: canonical
---
# CAP-EPT-051 — Memory Acquisition and Collection Boundary

## 1. Définition
Définir l’acquisition mémoire technique full/process lorsque la capability est déclarée disponible, avec scope, platform/privilege dependencies, impact preview, progress/partial/failure et output reference, sans choisir d’outil ni de format.

## 2. Problème utilisateur
Memory acquisition peut être unsupported, privileged, sensible ou perturbatrice. L’utilisateur doit distinguer demande, éligibilité, acquisition, output et analyse forensic.

## 3. Objectifs
Vérifier support/capability/privilege ; borner full/process scope ; exposer impact/size/duration conceptuels ; suivre progress/partial/failure/cancel ; produire neutral collected output/provenance.

## 4. Non-objectifs
Aucun memory tool, command, image format, crypto/storage protocol, Memory Forensics analysis, credential extraction, bypass or platform promise.

## 5. Propriétaire
Endpoint owns acquisition facts/output; Investigate owns request/business context and future memory analysis/Artifact qualification.

## 6. Utilisateurs
DFIR Analyst, Endpoint Operator, Evidence Reviewer, Govern Reviewer, Security/Privacy Reviewer.

## 7. Conditions d’entrée
CAP-EPT-048 memory item, declared memory capability, target/platform/privilege context, policy/authority, adequate resource/impact information or explicit unknown.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| memory plan item | CAP-EPT-048 | full/process scope | oui | current plan | no acquire |
| platform/capability | Endpoint | support | oui | current | unsupported |
| privilege/resource/impact | Endpoint/Policy | prerequisite | oui si known | start time | blocked/warning |
| authority | Govern/Security | gate | according impact | current | awaiting-authority |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Collection Request/Case | Investigate | purpose | read |
| Technical Plan | Endpoint | scope/limits | read |
| endpoint-agent | Endpoint | platform/capability/resources | read |
| Endpoint Policy | Settings | restrictions | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Memory Acquisition Attempt | create/update | Endpoint | exact scope/authority refs |
| Memory Collection Item | create/partial/fail | Endpoint | neutral output reference |
| Impact/Support Assessment | derive | Endpoint | no support assumption |

## 11. Fonctionnalités
Validate full/process scope and declared support, record estimated impact, begin only with gates satisfied, expose progress/partial/errors/cancel, preserve sensitive-output restrictions and output reference.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect support/impact | Analyst | assessment | 0 | read | reasons visible | non |
| acquire memory | authorized operator/path | attempt | 2/3 by real impact | support + authority | acquiring/output | according classification |
| cancel | Operator | attempt | 2 | cancellable | request recorded | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| support/privilege check | oui | oui | oui | explain | declared capability matrix |
| impact estimate | oui | source rules | oui | summarize | raw estimate |
| progress/error summary | oui | oui | oui | oui | raw status |
| bypass unsupported/authority | non | interdit | non | interdit | no operation |

## 14. États fonctionnels
`preparing`, `unsupported`, `policy-blocked`, `awaiting-authority`, `acquiring`, `partial`, `completed`, `failed`, `timed-out`, `cancel-requested`, `cancelled`.

## 15. États d’interface
No Screen ID. Sensitive memory output is masked/restricted; unsupported remains explicit.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Memory Collection Item | Endpoint output | CAP-EPT-054/055/Investigate | acquired != analyzed |
| progress/error | Endpoint state | CAP-EPT-053 | no false completion |
| impact/provenance | metadata | Audit/Govern | source-backed |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-048 | memory item | CAP-EPT-051 | scope/gates | plan |
| CAP-EPT-051 | output | CAP-EPT-054/055 | item/status/provenance | attempt |
| Investigate | analyze later | Memory Forensics | qualified reference | collection source retained |

## 18. Dépendances
Endpoint memory-collection source, Investigate CAP-INV-207, CAP-EPT-048/053..055, OPEN-005/008/013/014.

## 19. Source de vérité
Endpoint SOT of acquisition facts; Investigate owns request, qualification and forensic analysis.

## 20. Provenance et audit
Request/plan, Agent/platform, full/process scope, capability/privilege, impact estimate, authority, start/end/progress/errors, output ref and correlation.

## 21. Permissions fonctionnelles
Memory acquisition request/execute/cancel, sensitive output read, provenance, cross-tenant deny, step-up according impact.

## 22. Limites et erreurs
Memory acquired ≠ analyzed ; collected output ≠ Artifact/Evidence automatically ; unsupported platform cannot be overridden by documentation ; cancel ≠ rollback.

## 23. Métriques
Unsupported/policy-blocked, partial/failure/timeout/cancel, impact-estimate coverage, output/provenance completeness.

## 24. Classification de livraison
`draft / defined / planned`; no acquisition tool, format, command or platform delivered.

## 25. Critères d’acceptation
**Given** memory acquisition is unsupported, **When** requested, **Then** no attempt starts and unsupported reason is visible.

**Given** acquisition returns partial output, **When** completed/failed state is reconciled, **Then** partial remains explicit and output is not presented complete.

**Given** memory bytes are received, **When** handed off, **Then** they are not considered analyzed or Evidence automatically.

## 26. Questions ouvertes
OPEN-005/008/013/014 remain open.

## 27. Consommateurs documentaires
EPT-4 packaging/transfer, Investigate Memory/Artifact/Evidence, Govern, Security, Quality.
