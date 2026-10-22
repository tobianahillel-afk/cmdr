---
id: cross-product-transitions
domain: 04-experience-architecture
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-008
  - REQ-UX-006
  - REQ-UX-007
  - REQ-AI-002
  - REQ-SEC-001
---

# Transitions interproduits

| Transition | Contexte transmis | Contrôle | Retour |
|---|---|---|---|
| Command → Investigate | tenant, env, Incident, période, entités, origine | permission Case ; création/réouverture idempotente | Incident Detail exact |
| Investigate → Govern | Case, Finding, Evidence référencées, impact, Action Request | complétude, autorité, séparation des tâches | Case/Finding exact |
| Govern → Command | Decision, Response Run, Result, risque résiduel | projection autorisée | Incident et vue source |
| Settings → Investigate | Endpoint/Fleet projection, tenant, env | aucune administration transférée | Fleet avec sélection |
| Investigate → Endpoint operation | Endpoint, collecte/action demandée | collecte ou Govern selon classe | Case et Artifact/Run |
| Product → Studio | objet source, workflow/run/version | Studio n'acquiert pas l'objet métier | workspace source |
| Studio → Product | résultat proposé, Tool Calls, trace | acceptation par produit ; action risquée vers Govern | run Studio et objet source |

## Erreurs

Destination interdite, objet supprimé, contexte expiré et tenant incompatible sont affichés avant mutation. Aucune redirection silencieuse. Une transition échouée conserve le workspace source et une correlation ID.

## Trace

Chaque transition enregistre produit source/destination, acteur, objet source, objet cible, tenant, résultat et return origin. Les références sont permission-aware et bidirectionnelles lorsque autorisées.

## Critère d’acceptation

**Given** un Finding confirmé et une Action Request incomplète,  
**When** l'utilisateur tente Investigate → Govern,  
**Then** Govern n'est pas ouvert comme décision prête, les champs manquants sont expliqués, le Case reste intact et aucune autorité n'est transférée.
