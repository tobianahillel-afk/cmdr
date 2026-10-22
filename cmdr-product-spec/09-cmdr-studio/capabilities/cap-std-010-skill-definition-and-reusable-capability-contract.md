---
id: CAP-STD-010
title: Skill Definition and Reusable Capability Contract
product: cmdr-studio
module: skills
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-OBJ-009, REQ-AI-002]
open_decisions: []
source-of-truth: canonical
---
# CAP-STD-010 — Skill Definition and Reusable Capability Contract
## 1. Définition
Skill est un contrat réutilisable de comportement/connaissance Studio avec I/O, dependencies, preconditions, constraints, version et provenance. Skill ≠ Tool ≠ Workflow ≠ Agent.
## 2. Problème utilisateur
Un consommateur doit comprendre ce qu'une Skill fournit sans lui attribuer automatiquement exécution, orchestration ou permissions Tool.
## 3. Objectifs
Définir purpose, owner, consumers, reusable contract, I/O, Tool/Skill deps, preconditions, AI dependency éventuelle, non-AI path, version/lifecycle/provenance.
## 4. Non-objectifs
Pas de Workflow orchestration, Agent runtime, publishing/deployment détaillé ou schéma technique final.
## 5. Propriétaire
CMDR Studio / Skills; l'objet Skill canonique reste la source du lifecycle métier.
## 6. Utilisateurs
Automation Designer principal; Studio Reviewer et consommateurs Command/Investigate autorisés secondaires.
## 7. Conditions d’entrée
Skill identity/owner, tenant/env, consumer intent et dependency refs disponibles.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Skill intent | designer | purpose/consumer | oui | draft | bloque définition |
| Tool/Skill refs | Studio | dependencies | conditionnel | version visible | Partial |
| constraints | Studio/Security | permission/AI limits | oui | courant | bloque availability claim |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Skill | CMDR Studio | identity/lifecycle | read/write owner path |
| Version | CMDR Studio | version | read |
| Tool | CMDR Studio | dependency contract | read/reference |
## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Skill | create/update draft | CMDR Studio | ne devient pas Tool/Workflow/Agent |
## 11. Fonctionnalités
Purpose, reusable behavior, I/O, Tool/Skill dependencies, preconditions, constraints, optional AI dependency, version and provenance.
## 12. Actions utilisateur
Class 0 inspect; Class 2 create/update draft; Class 1 validate reusable contract.
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| validate Skill | oui | oui | oui | oui | checklist |
| author description | oui | oui | oui | oui | manuel |
## 14. États fonctionnels
draft, review-ready, available, partial, deprecated, disabled, superseded; deployment détaillé deferred STD-4.
## 15. États d’interface
Les états communs exposent source/dependency/freshness sans masquer les limitations.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Skill contract | Skill | Library/future consumers | I/O/deps/constraints explicites |
| reusable ref | reference | consumer | Tool permissions non héritées |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Library | open Skill | Skill surface | Skill/version/return-origin | Library |
| consumer | reference Skill | Studio | consumer context/version | consumer |
## 18. Dépendances
Library, Tool contracts, Versioning, Security et futurs Workflow/Agent seulement comme consumers.
## 19. Source de vérité
Skill object/source Studio; dependencies référencées chez leurs owners.
## 20. Provenance et audit
Creator/owner, exact version, dependency versions, changes, consumer links et source refs.
## 21. Permissions fonctionnelles
Skill read, restricted read, create/update, version/deprecate/reuse; Tool invocation permissions restent distinctes.
## 22. Limites et erreurs
Missing dependency, cycle candidate, incompatible version, restricted input/dependency ou unavailable runtime sont explicites.
## 23. Métriques
Valid Skills, dependency gaps, reuse links, deprecated usage, no-AI coverage.
## 24. Classification de livraison
`defined / planned`; Skill definition ≠ deployment.
## 25. Critères d’acceptation
**Given** Skill référence Tool, **When** Skill est visible, **Then** Tool invoke permission n'est pas héritée.  
**Given** dependency manquante, **When** Skill est validée, **Then** state devient Partial/unavailable sans substitution.  
**Given** IA off, **When** Skill est découverte/validée, **Then** core path reste fonctionnel.
## 26. Questions ouvertes
Aucune nouvelle OPEN; publishing/deployment reste STD-4.
## 27. Consommateurs documentaires
Library, Builder, future STD-2/3/4, Command/Investigate autorisés, Quality et Roadmap.
