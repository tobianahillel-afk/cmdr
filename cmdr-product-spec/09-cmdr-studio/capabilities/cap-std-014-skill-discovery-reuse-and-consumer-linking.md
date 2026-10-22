---
id: CAP-STD-014
title: Skill Discovery, Reuse and Consumer Linking
product: cmdr-studio
module: skills
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-OBJ-009]
open_decisions: []
source-of-truth: canonical
---
# CAP-STD-014 — Skill Discovery, Reuse and Consumer Linking
## 1. Définition
Discovery/reuse permission-aware de Skills via Studio Library et Shared Search, avec explicit consumer linking et aucun auto-binding silencieux.
## 2. Problème utilisateur
Trouver une Skill ne signifie pas que ses dépendances, Tool permissions ou consumer compatibility sont satisfaites.
## 3. Objectifs
Search/filter, owner/status/version/risk/availability, dependency preview, compatibility and explicit reuse links.
## 4. Non-objectifs
Pas de generic Search engine, automatic Workflow/Agent binding, Tool permission inheritance ou deployment.
## 5. Propriétaire
Studio / Skills+Library possède discovery semantics; Shared possède generic Search; consumer possède son object/link usage.
## 6. Utilisateurs
Automation Designer, SOC Analyst L2, Investigate Analyst et consuming product users.
## 7. Conditions d’entrée
Authorized Library context, Skill metadata, consumer context and tenant/env resolved.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Library query | user | search/filter | oui | request | permitted catalog |
| Skill metadata | Studio | owner/version/status/risk | oui | source freshness | Partial |
| consumer context | Command/Investigate/Studio | authorized use | conditionnel | current | reuse blocked |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Skill | CMDR Studio | catalog/compatibility | read |
| Version | CMDR Studio | exact version | read |
| Search query mechanism | Shared | generic search | consume |
## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| consumer Skill reference | create explicit link | consumer owner / Studio ref semantics | no auto binding |
## 11. Fonctionnalités
Search/filter, dependency preview, compatible-consumer assessment, explicit reuse/linking and provenance.
## 12. Actions utilisateur
Class 0 search/filter/inspect; Class 1 compatibility check; Class 2 create explicit reuse/reference link.
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| discovery | oui | oui | oui | oui | filters + Shared Search |
| compatibility suggestion | oui | oui | oui | oui | deterministic fields |
## 14. États fonctionnels
visible, compatible, incompatible, partial, restricted, deprecated; visibility != permission to execute dependencies.
## 15. États d’interface
Permission denied masks restricted Skill/dependency fields; Partial shows missing source/dependency.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Skill discovery results | projection | user/consumer | authorized metadata only |
| reuse link | exact reference | Command/Investigate/future Studio | version/provenance pinned |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| consumer | discover Skill | Library | tenant/env/consumer/return | consumer |
| Library | link Skill | consumer object | Skill/version/provenance | consumer confirms |
## 18. Dépendances
CAP-STD-001/010/011/013, Shared Search/Linking, Security and consumer owner semantics.
## 19. Source de vérité
Skill metadata/lifecycle remains Studio; generic search remains Shared; consumer link remains consumer-owned.
## 20. Provenance et audit
Search/link actor, Skill/version, consumer object, compatibility snapshot and return-origin retained.
## 21. Permissions fonctionnelles
Skill read/restricted read/reuse; dependent Tool permissions remain independent.
## 22. Limites et erreurs
Missing dependency, incompatible consumer, restricted Skill, deprecated version or Shared Search degradation explicit.
## 23. Métriques
Discovery results, compatibility failures, explicit reuse links, deprecated reuse and silent-binding violations.
## 24. Classification de livraison
`defined / planned`; no auto-binding/deployment implementation.
## 25. Critères d’acceptation
**Given** Skill found but Tool dependency permission absent, **When** reuse inspected, **Then** Skill visibility remains distinct from Tool authorization.  
**Given** missing dependency, **When** search result opened, **Then** Partial/blocked state and dependency visible.  
**Given** IA off, **When** Skill searched/reused, **Then** filters/Search and explicit link remain functional.
## 26. Questions ouvertes
No new OPEN; Workflow/Agent consumers remain future references only.
## 27. Consommateurs documentaires
Library, Skills, Command/Investigate consumers, future STD-2/3, Shared Search/Linking, Quality.
