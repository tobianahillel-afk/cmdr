---
id: CAP-STD-066
title: Deployment Health, Reversion and Previous-Version Recovery
product: cmdr-studio
module: versions-and-deployment
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-10
requirement_ids: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-016, REQ-PROD-019, REQ-OBJ-009, REQ-AI-002, REQ-SEC-001, REQ-SEC-002, REQ-PROD-017]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-STD-066 — Deployment Health, Reversion and Previous-Version Recovery

## 1. Définition
Assess Studio Deployment health and define previous-version reversion/recovery semantics, including partial rollout/failure/manual intervention, while explicitly separating Studio asset reversion from Govern Response Rollback and from restoration of business state.

## 2. Problème utilisateur
Sans **Deployment Health Assessment**, Studio ne peut pas distinguer de façon reproductible l’asset/version, le contexte, les preuves, les limites et le lifecycle; un score, une projection technique ou un état local pourrait alors être pris à tort pour une autorité, un succès production ou un objet d’un autre produit.

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
| Studio Deployment | Studio | oui | exact/current | blocked/incomplete if required |
| deployment health observations | Studio | oui | exact/current | blocked/incomplete if required |
| runtime/asset availability | Studio | oui | exact/current | blocked/incomplete if required |
| compatibility issues | Studio | selon scope | exact/current | blocked/incomplete if required |
| rollout state | Studio | selon scope | exact/current | blocked/incomplete if required |
| previous version | Studio | selon scope | exact/current | blocked/incomplete if required |
| reversion eligibility | Studio | selon scope | exact/current | blocked/incomplete if required |
| reversion authority/context | Studio | selon scope | exact/current | blocked/incomplete if required |
| manual intervention notes | Studio | selon scope | exact/current | blocked/incomplete if required |
| post-reversion observations | Studio | selon scope | exact/current | blocked/incomplete if required |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| deployment | CMDR Studio | exact/versioned authorized projection | read/link only |
| version | CMDR Studio | exact/versioned authorized projection | read/link only |
| environment | Platform Settings | exact/versioned authorized projection | read/link only |
| response-rollback | Govern | exact/versioned authorized projection | read/link only |
| response-run | Govern | exact/versioned authorized projection | read/link only |
| result | Govern | exact/versioned authorized projection | read/link only |
| workflow | CMDR Studio | exact/versioned authorized projection | read/link only |
| automation-agent | CMDR Studio | exact/versioned authorized projection | read/link only |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Deployment Health Assessment | create/review/supersede as applicable | CMDR Studio | no implicit source-owner mutation |
| Deployment Reversion Record | create/review/supersede as applicable | CMDR Studio | no implicit source-owner mutation |
| Post-Reversion Assessment | create/review/supersede as applicable | CMDR Studio | no implicit source-owner mutation |

## 11. Fonctionnalités
Studio Deployment; deployment health observations; runtime/asset availability; compatibility issues; rollout state; previous version; reversion eligibility; reversion authority/context; manual intervention notes; post-reversion observations; validation des dépendances; version pinning; partial/inconclusive; warnings; supersession; historique.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspecter Deployment Health Assessment | authorized user | Deployment Health Assessment | 0 | read | contexte exact | non |
| évaluer préconditions/compatibilité | reviewer | Deployment Health Assessment | 1 | sources lisibles | pass/warning/blocked explicable | non |
| créer/mettre à jour Deployment Health Assessment | owner/reviewer | Deployment Health Assessment | 2 | manage permission | record versionné | OPEN-013 |
| effectuer une transition pouvant produire un effet | authorized operator | Deployment Health Assessment | 3 | permissions + dépendances Govern/Settings applicables | transition explicite seulement | governed |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| vérifier entrées/dépendances | oui | oui | oui | explication sourcée | checklist/règles |
| comparer/résumer observations | oui | oui si structuré | oui | oui, attribuée | tables/diff/review humain |
| publier/déployer/autoriser ou inventer un résultat | accountable path only | validation only | non autonome | interdit | contrôle explicite |

L’IA ne peut jamais déclarer seule production-safe, accorder permission, contourner Govern, publier/déployer seule, inventer un résultat, masquer une régression ou révéler un secret.

## 14. États fonctionnels
`healthy`, `degraded`, `failed`, `reversion-candidate`, `reverting`, `reverted`, `partial-reversion`, `recovery-required`, `inconclusive`. Lifecycle et outcome restent distincts quand les deux existent.

## 15. États d’interface
Loading conserve le contexte; Empty nomme la pièce manquante; Partial liste les sources manquantes; Error garde le dernier état valide; Offline interdit les mutations non garanties; Permission denied masque; Stale exige revalidation. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Deployment Health Assessment | Studio record | downstream Studio assurance/lifecycle | version/scope/limits/provenance explicit |
| Deployment Reversion Record | Studio record | downstream Studio assurance/lifecycle | version/scope/limits/provenance explicit |
| Post-Reversion Assessment | Studio record | downstream Studio assurance/lifecycle | version/scope/limits/provenance explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| sources canoniques | préconditions satisfaites | Deployment Health Assessment | versions, tenant/env, permissions, provenance | source |
| Deployment Health Assessment | gap/contradiction | Assurance/source owner | missing refs + limites | même lineage |
| Deployment Health Assessment | résultat accepté | étape Studio suivante | outcome + warnings + provenance | return origin |

## 18. Dépendances
Studio STD-1/2/3 selon les assets consommés, capabilities STD-4 adjacentes, Security permission model, Govern/Settings/Shared/Endpoint boundaries et OPEN-008, OPEN-013, OPEN-015. Aucune dépendance n’est déclarée implémentée par simple référence.

## 19. Source de vérité
**Deployment Health Assessment** est source de sa sémantique Studio seulement. Les objets/faits externes restent chez leurs owners. Evaluation/Simulation/health/output ne deviennent pas Govern Approval/Decision/Result/Evidence automatiquement.

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
**Given** a deployment reports healthy but a later Automation Run fails, **When** health is reviewed, **Then** deployment health remains distinct from runtime business success and the Run failure is linked separately.

**Given** reversion to a previous version completes, **When** post-reversion assessment runs, **Then** the previous asset version is restored without claiming production or business state restoration.

**Given** reversion partially succeeds, **When** health is assessed, **Then** partial-reversion and remaining affected scope stay explicit.

**Given** the previous version is unavailable, **When** reversion is requested, **Then** the request is blocked and manual recovery is required rather than fabricated rollback.

## 26. Questions ouvertes
OPEN-008, OPEN-013, OPEN-015 restent ouvertes. Le lot ne ferme aucune des 18 OPEN par rédaction et n’en crée aucune nouvelle.

## 27. Consommateurs documentaires
Studio Assurance/Evaluations/Simulations/Versions & Deployment, Library/Builder/Skills/Workflows/Agents/Control Room selon besoin, registers, Requirements/baseline, screen maps, Quality, Roadmap et external-owner boundaries.
