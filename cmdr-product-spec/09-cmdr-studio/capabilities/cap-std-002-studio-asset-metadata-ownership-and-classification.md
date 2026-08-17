---
id: CAP-STD-002
title: Studio Asset Metadata, Ownership and Classification
product: cmdr-studio
module: library
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-OBJ-009]
open_decisions: []
source-of-truth: canonical
---
# CAP-STD-002 — Studio Asset Metadata, Ownership and Classification
## 1. Définition
Contrat fonctionnel de metadata Studio: identity, type, owner, contributors, version, lifecycle, delivery, risk, execution class, AI/non-AI, dependencies, consumers, supersession et provenance.
## 2. Problème utilisateur
Sans metadata commune, Library et consommateurs peuvent confondre asset visible, exécutable, owner et runtime.
## 3. Objectifs
Rendre ownership, classification et limitations observables sans créer un schéma physique final.
## 4. Non-objectifs
Aucun JSON Schema final, storage schema, permission atomique ou transfert d'ownership.
## 5. Propriétaire
CMDR Studio possède la sémantique metadata des assets Studio; Product Architecture conserve les règles d'ownership globales.
## 6. Utilisateurs
Automation Designer, Studio Reviewer, Auditor et consommateurs autorisés.
## 7. Conditions d’entrée
Asset résolu, source canonique identifiable, tenant/env et actor connus.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| asset identity | source Studio | stable ref | oui | courant | erreur qualité |
| owner/contributors | Ownership Register/Studio | identities | oui | courant | owner unknown explicite |
| classifications | Studio/Security | lifecycle/risk/delivery | oui | version courante | Partial |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Skill | CMDR Studio | identity/version/lifecycle | read |
| Version | CMDR Studio | compatibility/supersession | read |
| Principal | Platform Settings | owner identity projection | read |
## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Studio asset metadata | create/update draft metadata | CMDR Studio | fonctionnel seulement; pas de schéma physique |
## 11. Fonctionnalités
Metadata validation, ownership projection, classification risk/execution/AI, dependency/consumer links et supersession.
## 12. Actions utilisateur
Class 0 inspect metadata; Class 2 edit draft metadata, assign owner reference, record supersession.
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| validate required fields | oui | oui | oui | non requis | checklist |
| suggest tags/classification | oui | oui | oui | oui | règles/manuelle |
## 14. États fonctionnels
draft, complete, partial, deprecated, superseded; ces états ne remplacent aucun lifecycle objet externe.
## 15. États d’interface
Les six états obligatoires + Stale indiquent source/freshness et masquent les champs non autorisés.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| metadata projection | metadata | Library/consumers | owner/version/lifecycle explicites |
| supersession reference | link | consumers | histoire préservée |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| owning surface | update metadata | Library | asset/version/provenance | owning surface |
| Library | inspect source | canonical owner | asset ref/return-origin | Library |
## 18. Dépendances
Ownership Register, canonical objects, Security classification et Shared Linking.
## 19. Source de vérité
L'asset source conserve sa vérité; CAP-STD-002 définit uniquement la projection metadata Studio.
## 20. Provenance et audit
Chaque changement garde actor, reason, old/new metadata, version, tenant/env et supersession.
## 21. Permissions fonctionnelles
Metadata read/restricted-read, draft update, ownership reference change et classification update; aucun namespace atomique final.
## 22. Limites et erreurs
Owner absent, source inaccessible, classification contradictoire, cross-tenant ref ou stale metadata restent explicites.
## 23. Métriques
Completeness, unknown owner, stale metadata, superseded-reference usage, classification gaps.
## 24. Classification de livraison
`defined / planned`; aucune base de données ou catalogue runtime n'est livré.
## 25. Critères d’acceptation
**Given** owner inconnu, **When** asset est affiché, **Then** owner unknown est visible et jamais inventé.  
**Given** asset superseded, **When** historique est consulté, **Then** old/new refs restent navigables.  
**Given** IA absente, **When** metadata est validée, **Then** règles/checklists restent suffisantes.
## 26. Questions ouvertes
Aucune nouvelle OPEN; final schema et atomic permissions restent futurs.
## 27. Consommateurs documentaires
Library, Skills, Builder references, future STD lots, Registers, Quality et Roadmap.
