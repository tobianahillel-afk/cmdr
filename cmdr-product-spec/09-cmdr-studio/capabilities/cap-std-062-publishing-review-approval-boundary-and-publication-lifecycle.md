---
id: CAP-STD-062
title: Publishing Review, Approval Boundary and Publication Lifecycle
product: cmdr-studio
module: versions-and-deployment
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-10
requirement_ids: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-016, REQ-PROD-019, REQ-OBJ-009, REQ-AI-002, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-STD-062 — Publishing Review, Approval Boundary and Publication Lifecycle

## 1. Définition
Define Studio publishing review and publication lifecycle for candidate versions while preserving the boundary that publishing review is not automatically Govern Approval/Decision and publication does not deploy or enable.

## 2. Problème utilisateur
Sans **Publishing Review Record**, Studio ne peut pas distinguer de façon reproductible l’asset/version, le contexte, les preuves, les limites et le lifecycle; un score, une projection technique ou un état local pourrait alors être pris à tort pour une autorité, un succès production ou un objet d’un autre produit.

## 3. Objectifs
Version-pinner le sujet; conserver tenant/environment et sources; rendre partial/failure/warning/inconclusive explicites; préserver ownership et permissions; maintenir provenance et supersession; fournir un chemin sans IA.

## 4. Non-objectifs
Aucun moteur/API/protocole/package registry/scheduler/code/schéma physique/JSON Schema/RBAC final/provider imposé, aucun nouveau Screen ID, aucune capability Endpoint et aucune autorité parallèle à Govern.

## 5. Propriétaire
CMDR Studio / versions-and-deployment. Govern conserve Approval/Decision/Response Run/Result et response rollback; Settings conserve provider/integration/secret/environment admin; Shared conserve les moteurs génériques; Endpoint conserve ses primitives techniques.

## 6. Utilisateurs
Studio Asset Owner, Assurance/Release Operator et Reviewer; Platform Admin, Govern Reviewer, Auditor et consommateurs produit seulement selon leurs responsabilités et permissions source.

## 7. Conditions d’entrée
Asset/version, tenant/environment, références source et restrictions identifiables; permission de lecture puis de mutation locale si requise; Secret References seulement; aucune projection externe ne transfère ownership ou autorité.

## 8. Entrées fonctionnelles
| Entrée | Source | Requise | Fraîcheur | Si absente |
|---|---|---:|---|---|
| Publishing Candidate Package | Studio | oui | exact/current | blocked/incomplete if required |
| reviewer | Studio | oui | exact/current | blocked/incomplete if required |
| quality/readiness evidence | Studio | oui | exact/current | blocked/incomplete if required |
| publication permission | Studio | selon scope | exact/current | blocked/incomplete if required |
| Govern dependency projection | Studio | selon scope | exact/current | blocked/incomplete if required |
| changes requested | Studio | selon scope | exact/current | blocked/incomplete if required |
| version lineage | Studio | selon scope | exact/current | blocked/incomplete if required |
| withdrawal/supersession context | Studio | selon scope | exact/current | blocked/incomplete if required |
| tenant/environment restrictions | canonical owner | selon scope | exact/current | blocked/incomplete if required |
| provenance | Studio | selon scope | exact/current | blocked/incomplete if required |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| version | CMDR Studio | exact/versioned authorized projection | read/link only |
| evaluation | CMDR Studio | exact/versioned authorized projection | read/link only |
| simulation | CMDR Studio | exact/versioned authorized projection | read/link only |
| approval | Govern | exact/versioned authorized projection | read/link only |
| decision | Govern | exact/versioned authorized projection | read/link only |
| workflow | CMDR Studio | exact/versioned authorized projection | read/link only |
| skill | CMDR Studio | exact/versioned authorized projection | read/link only |
| automation-agent | CMDR Studio | exact/versioned authorized projection | read/link only |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Publishing Review Record | create/review/supersede as applicable | CMDR Studio | no implicit source-owner mutation |
| Published Studio Asset Version | create/review/supersede as applicable | CMDR Studio | no implicit source-owner mutation |

## 11. Fonctionnalités
Publishing Candidate Package; reviewer; quality/readiness evidence; publication permission; Govern dependency projection; changes requested; version lineage; withdrawal/supersession context; tenant/environment restrictions; provenance; validation des dépendances; version pinning; partial/inconclusive; warnings; supersession; historique.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspecter Publishing Review Record | authorized user | Publishing Review Record | 0 | read | contexte exact | non |
| évaluer préconditions/compatibilité | reviewer | Publishing Review Record | 1 | sources lisibles | pass/warning/blocked explicable | non |
| créer/mettre à jour Publishing Review Record | owner/reviewer | Publishing Review Record | 2 | manage permission | record versionné | OPEN-013 |
| effectuer une transition pouvant produire un effet | authorized operator | Publishing Review Record | 3 | permissions + dépendances Govern/Settings applicables | transition explicite seulement | governed |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| vérifier entrées/dépendances | oui | oui | oui | explication sourcée | checklist/règles |
| comparer/résumer observations | oui | oui si structuré | oui | oui, attribuée | tables/diff/review humain |
| publier/déployer/autoriser ou inventer un résultat | accountable path only | validation only | non autonome | interdit | contrôle explicite |

L’IA ne peut jamais déclarer seule production-safe, accorder permission, contourner Govern, publier/déployer seule, inventer un résultat, masquer une régression ou révéler un secret.

## 14. États fonctionnels
`review-pending`, `changes-requested`, `publishable`, `published`, `rejected`, `withdrawn`, `superseded`. Lifecycle et outcome restent distincts quand les deux existent.

## 15. États d’interface
Loading conserve le contexte; Empty nomme la pièce manquante; Partial liste les sources manquantes; Error garde le dernier état valide; Offline interdit les mutations non garanties; Permission denied masque; Stale exige revalidation. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Publishing Review Record | Studio record | downstream Studio assurance/lifecycle | version/scope/limits/provenance explicit |
| Published Studio Asset Version | Studio record | downstream Studio assurance/lifecycle | version/scope/limits/provenance explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| sources canoniques | préconditions satisfaites | Publishing Review Record | versions, tenant/env, permissions, provenance | source |
| Publishing Review Record | gap/contradiction | Assurance/source owner | missing refs + limites | même lineage |
| Publishing Review Record | résultat accepté | étape Studio suivante | outcome + warnings + provenance | return origin |

## 18. Dépendances
Studio STD-1/2/3 selon les assets consommés, capabilities STD-4 adjacentes, Security permission model, Govern/Settings/Shared/Endpoint boundaries et OPEN-007, OPEN-013, OPEN-015. Aucune dépendance n’est déclarée implémentée par simple référence.

## 19. Source de vérité
**Publishing Review Record** est source de sa sémantique Studio seulement. Les objets/faits externes restent chez leurs owners. Evaluation/Simulation/health/output ne deviennent pas Govern Approval/Decision/Result/Evidence automatiquement.

## 20. Provenance et audit
Asset/version, tenant/env, sources, critères/dépendances, restrictions, permission/authority refs, auteur/reviewer, règle ou IA, observations, failures/warnings, timestamps, supersession et correlation ids sont conservés; aucun historique n’est effacé.

## 21. Permissions fonctionnelles
Read; create/update; no-effect evaluate/run quand applicable; sensitive read masqué; review; lifecycle transition; provenance/export preparation. `perm.studio.*` et `perm.cmdr-studio.*` restent non normalisés. Toute action effectful exige le contrôle d’autorité applicable.

## 22. Limites et erreurs
Stale version, tenant mismatch, permission denied, Secret Reference manquante, provider/runtime indisponible, dépendance incompatible, conflit d’observations, incomplete evidence ou concurrent supersession donnent blocked/partial/inconclusive/error sans fallback élargissant le scope.

## 23. Métriques
Records complets/incomplets, blocked dependencies/permissions, partial/inconclusive, warnings, supersessions, reviews humains, assistance IA attribuée et boundary violations; aucune cible/SLO numérique non approuvé.

## 24. Classification de livraison
`defined` / `planned`; couverture documentaire fonctionnelle uniquement, sans preuve d’implémentation, déploiement réel ou disponibilité.

## 25. Critères d’acceptation
**Given** a publishing reviewer accepts a candidate, **When** publication lifecycle proceeds, **Then** the act remains a Studio publishing review and is not silently converted into Govern Approval.

**Given** a Govern Approval is required, **When** publication is requested without it, **Then** the publication is blocked or routed to Govern rather than self-authorized.

**Given** a version is published, **When** deployment is inspected, **Then** no Deployment exists automatically.

**Given** AI proposes publishing, **When** the proposal is reviewed, **Then** AI cannot publish autonomously and the accountable path remains explicit.

## 26. Questions ouvertes
OPEN-007, OPEN-013, OPEN-015 restent ouvertes. Le lot ne ferme aucune des 18 OPEN par rédaction et n’en crée aucune nouvelle.

## 27. Consommateurs documentaires
Studio Assurance/Evaluations/Simulations/Versions & Deployment, Library/Builder/Skills/Workflows/Agents/Control Room selon besoin, registers, Requirements/baseline, screen maps, Quality, Roadmap et external-owner boundaries.
