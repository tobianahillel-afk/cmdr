---
id: investigate-automation-and-ai-model
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-AI-001
  - REQ-AI-002
  - REQ-AI-003
  - REQ-AI-006
  - REQ-AI-010
  - REQ-AI-011
open_decisions:
  - OPEN-007
  - OPEN-013
  - OPEN-015
---

# Automation and AI model — Investigate

## Authority

| Producteur | Observe | Propose | Modifie Investigate | Confirme Finding | Crée Decision |
|---|---:|---:|---:|---:|---:|
| Humain autorisé | oui | oui | selon permission/classe | selon rôle/revue | non, Govern |
| Règle | oui | oui | seulement contrat classe 2 explicite | non | non |
| Moteur déterministe | oui | oui | validation, liens ou transformation tracée | non | non |
| Workflow | oui | oui | étapes déployées et permissionnées | non sans humain autorisé | non |
| Automation Agent | oui | oui | proposal par défaut | non | non |
| Govern | consomme | retourne statut | lifecycle Action Request | ne confirme pas Finding | oui selon autorité |

## Usages autorisés

Aide à la Query, Hypothesis proposée, résumé et regroupement de résultats, Evidence candidate, Finding draft, Action Request draft et explication de relations.

## Interdictions

- cacher les sources ;
- exécuter automatiquement une Query sensible sans contrat ou déclencheur ;
- convertir automatiquement Event ou Artifact en Evidence ;
- confirmer un Finding ;
- créer ou approuver une Decision ;
- modifier une Evidence sans version ;
- supprimer une trace ;
- remplacer Case Workspace par un chatbot.

## Alternative sans IA

Chaque capability expose requêtes manuelles, builders, règles, moteurs déterministes, workflows non agentiques et actions humaines suffisants. L’absence de modèle ne change que la disponibilité des suggestions.

## Provenance requise

Initiateur, producer type, ID/version, run, sources, Tool Calls, facteurs, incertitude, objet proposé, owner humain, disposition accept/modify/reject et correlation ID.
