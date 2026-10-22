---
id: ADR-0003-canonical-object-chain
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-OBJ-001
  - REQ-OBJ-002
  - REQ-OBJ-003
  - REQ-OBJ-004
  - REQ-OBJ-005
  - REQ-OBJ-007
---
# ADR-0003 — Chaîne nominale des objets opérationnels

## 1. Identifiant

`ADR-0003`

## 2. Titre

Chaîne nominale des objets opérationnels

## 3. Statut

Draft. Cette ADR n'est pas approuvée et ne remplace aucune décision source au-delà de ce qu'elle applique explicitement.

## 4. Date

2026-08-03

## 5. Propriétaire

Product Architecture.

## 6. Décideurs attendus

Head of Product et propriétaires des domaines affectés ; Security, UX ou Engineering selon les effets décrits.

## 7. Requirement IDs

`REQ-OBJ-001`, `REQ-OBJ-002`, `REQ-OBJ-003`, `REQ-OBJ-004`, `REQ-OBJ-005`, `REQ-OBJ-007`

## 8. Contexte

CMDR relie détection, investigation, décision et réponse. Une progression nominale commune est nécessaire sans imposer la création de tous les objets.

## 9. Problème

Sans chaîne commune, les produits emploient des termes interchangeables et les handoffs perdent provenance et responsabilité.

## 10. Forces en présence

Simplicité, variété des scénarios, cardinalités futures, audit et distinction entre faits, conclusions, décisions et résultats.

## 11. Options étudiées

1. Aucune chaîne canonique.
2. Une chaîne obligatoire et linéaire.
3. Une chaîne nominale avec embranchements et étapes optionnelles.

## 12. Décision

Utiliser Telemetry Event → Detection → Signal → Alert → Incident → Case → Evidence → Finding → Action Request → Decision → Response Run → Result comme progression nominale, non obligatoire.

## 13. Justification

La chaîne fournit une grammaire commune tout en laissant les scénarios directs et les relations multiples.

## 14. Conséquences positives

- Transitions compréhensibles.
- Provenance de bout en bout.
- Meilleure séparation Finding/Decision/Result.

## 15. Conséquences négatives

- Les cardinalités et erreurs restent à formaliser en Phase 7.

## 16. Risques

- Interprétation comme pipeline obligatoire.
- Confusion entre Artifact et Evidence.

## 17. Effets sur la navigation

La navigation permet des pivots et retours sans forcer une route unique.

## 18. Effets sur les objets

Chaque objet conserve son propriétaire et ses relations amont/aval.

## 19. Effets sur les permissions

Une transition ne transfère pas les permissions ; chaque étape réévalue l'autorisation.

## 20. Effets sur les parcours

Les huit parcours obligatoires utilisent la chaîne ou expliquent leurs étapes omises.

## 21. Effets sur les autres documents

`05-domain-model/`, `13-user-journeys/`, product boundaries et operating model.

## 22. Migration

Conserver les objets existants ; corriger les documents qui confondent les termes.

## 23. Critères de réévaluation

Réévaluer après formalisation des cardinalités et premiers parcours pilotes.

## 24. Questions encore ouvertes

- OPEN-014 — Artifact versus Attachment.
- OPEN-015 — Automation Run versus Response Run.
