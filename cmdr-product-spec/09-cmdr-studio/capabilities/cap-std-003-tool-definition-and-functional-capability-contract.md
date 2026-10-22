---
id: CAP-STD-003
title: Tool Definition and Functional Capability Contract
product: cmdr-studio
module: studio-foundations
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-OBJ-009, REQ-AI-002]
open_decisions: []
source-of-truth: canonical
---
# CAP-STD-003 — Tool Definition and Functional Capability Contract
## 1. Définition
Un Tool est une capacité exécutable décrite et possédée fonctionnellement par Studio. Tool ≠ Tool Call ≠ Skill ≠ Workflow ≠ Agent ≠ Endpoint primitive ≠ Govern Playbook.
## 2. Problème utilisateur
Les produits consommateurs doivent connaître purpose, operations, dependencies, risk et limits avant toute invocation.
## 3. Objectifs
Définir purpose, owner, users, operations, execution nature, consumers, I/O refs, runtime/provider deps, risk, availability, compatibility, version et limitations.
## 4. Non-objectifs
Aucun moteur, API, protocol, command, final Tool object schema ou orchestration.
## 5. Propriétaire
CMDR Studio / Studio foundations; le schéma objet final Tool reste différé conformément à l'Ownership Register.
## 6. Utilisateurs
Automation Designer, Studio Operator, Investigate Analyst et consommateurs produits autorisés.
## 7. Conditions d’entrée
Purpose/owner connus, operations bornées, references de dépendance explicites et tenant/env résolus.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Tool intent | designer | purpose/operations | oui | draft | bloque definition |
| runtime/provider refs | Settings projection | dependencies | conditionnel | freshness visible | unavailable/partial |
| risk/ownership | Studio/Security | classification | oui | courant | pas de claim executable |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Capability | Product Architecture | classification | read |
| Integration/Model Provider | Platform Settings | reference/availability | read |
| Version | CMDR Studio | version ref | read |
## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Tool functional definition | create/update draft | CMDR Studio | semantics only; final object deferred |
## 11. Fonctionnalités
Déclarer supported operations, deterministic/AI-capable/manual boundary, I/O refs, dependency refs, action/risk class, side-effect context et limits.
## 12. Actions utilisateur
Class 0 inspect Tool; Class 2 create/update draft; Class 1 validate contract. Exécution réelle appartient à CAP-STD-007/008 et au runtime owner.
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| contract validation | oui | oui | oui | oui | checklist déterministe |
| description assistance | oui | oui | oui | oui | authoring manuel |
## 14. États fonctionnels
draft, review-ready, available-reference, deprecated, disabled, superseded; published/deployed ne sont pas inférés.
## 15. États d’interface
Loading/Empty/Partial/Error/Offline/Permission denied/Stale suivent les patterns communs sans changer les semantics Tool.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Tool functional contract | definition | Library/Skills/consumers | purpose/ops/deps/limits explicites |
| eligibility prerequisites | metadata | CAP-STD-007 | aucune autorisation implicite |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Library | open Tool | Tool definition | Tool/version/return-origin | Library |
| Tool | reference in Skill | Skill | exact Tool ref/version/permission needs | Tool |
## 18. Dépendances
Security permissions, Settings providers/runtime, Shared Versioning/Trace et Govern pour actions gouvernées.
## 19. Source de vérité
Ce fichier porte la capability; Studio porte le Tool functional definition; les dependencies restent chez leurs owners.
## 20. Provenance et audit
Conserver creator, owner, version, changes, dependency refs, risk class et source requirement.
## 21. Permissions fonctionnelles
Tool metadata read, definition create/update, version management et invoke need séparés; aucun RBAC atomique final.
## 22. Limites et erreurs
Unsupported operation, missing owner, unavailable runtime, incompatible dependency, permission/governance requirement et invalid contract.
## 23. Métriques
Defined Tools, invalid contracts, dependency gaps, deprecated references et provenance completeness.
## 24. Classification de livraison
`defined / planned`; definition documentaire ≠ runtime disponible.
## 25. Critères d’acceptation
**Given** Tool défini mais runtime absent, **When** Tool est inspecté, **Then** définition reste valide et availability est séparée.  
**Given** Tool AI-capable, **When** utilisateur inspecte nature d'exécution, **Then** Tool n'est pas transformé en Agent.  
**Given** IA absente, **When** Tool est créé/validé, **Then** authoring et validation déterministes restent disponibles.
## 26. Questions ouvertes
Final Tool object/API/runtime restent futurs; aucune nouvelle OPEN nécessaire pour STD-1.
## 27. Consommateurs documentaires
Library, Skills, future Workflow/Agent lots, Command, Investigate, Govern, Settings, Security, Quality.
