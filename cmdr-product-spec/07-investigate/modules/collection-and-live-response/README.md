---
id: investigate-collection-and-live-response
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-014
  - REQ-PROD-018
  - REQ-PROD-020
open_decisions:
  - OPEN-007
  - OPEN-008
  - OPEN-013
  - OPEN-015
---
# Collection and Live Response

## Mission
Permettre à Investigate de préparer, suivre et qualifier des collectes et opérations endpoint liées à un Case, sans administrer la Fleet ni transférer l’autorité de Govern.

## Capabilities
- CAP-INV-201..208 : contexte Endpoint et collecte.
- CAP-INV-209..212 : Live Session, opérations, transfert et résultats.
- CAP-INV-213..215 : integrity/custody, provenance et préparation containment.

## Ownership
Investigate possède le contexte métier, Collection Request, relations Case et qualification analytique. Platform Settings possède Fleet/Policies. Endpoint Agent exécute localement et rapporte. Govern possède Action Request lifecycle, Decision, Response Run et Result. Studio possède Workflow/Automation Run/Tool Calls. Shared possède Jobs, Trace, Timeline, Notifications, Linking, Export et audit mechanisms.

## Invariants
Collection Request ≠ Collection Job ≠ Background Job. Live Session ≠ Terminal ≠ Automation Run ≠ Response Run. Operation Result ≠ Govern Result. Artifact ≠ Evidence. Containment request ≠ containment execution.

## Limites
Aucune API, protocole, commande, moteur, format, plateforme supportée, schéma objet final, permission atomique, écran détaillé ou capability CAP-INV-3xx.

## Shared Capabilities
| Shared Capability | Usage Collection/Live Response | Données locales | Source canonique |
|---|---|---|---|
| Background Jobs | queue, progression, cancel, partial | request/job/cible/Case | `12-shared-capabilities/background-jobs.md` |
| Notifications | états, erreurs et deep links | Case, Endpoint, job/session | `12-shared-capabilities/notification-center.md` |
| Object Linking | liens Case/Endpoint/Artifact/Run | typed relations | `12-shared-capabilities/object-linking-service.md` |
| Timeline / Trace | événements métier et provenance | source, acteur, statut, erreurs | Timeline Engine et mécanismes Shared |
| Inspector / Context Bar | inspection et contexte | projections permission-aware | Design System |
| Export | sortie contrôlée | Artifacts, transcript, custody | `12-shared-capabilities/export-engine.md` |

## Critère d’acceptation
Un utilisateur peut préparer et suivre une collecte ou une Live Session sans IA, voir les résultats partiels et revenir au Case, tandis que Fleet, Policy, exécution locale et autorité restent chez leurs owners.
