---
id: CAP-STD-058
title: Permission, Governance and Boundary Assurance
product: cmdr-studio
module: assurance
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-10
requirement_ids: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-016, REQ-PROD-019, REQ-OBJ-009, REQ-AI-002, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-STD-058 — Permission, Governance and Boundary Assurance

## 1. Définition
Assess whether a Studio asset’s permissions, restricted inputs, Secret References, tenant/environment scope, Tool risk, Human Gates and Govern/Settings/Shared/Endpoint dependencies are coherent without creating authority.

## 2. Problème utilisateur
Sans **Boundary Assurance Assessment**, Studio ne peut pas distinguer de façon reproductible l’asset/version, le contexte, les preuves, les limites et le lifecycle; un score, une projection technique ou un état local pourrait alors être pris à tort pour une autorité, un succès production ou un objet d’un autre produit.

## 3. Objectifs
Version-pinner le sujet; conserver tenant/environment et sources; rendre partial/failure/warning/inconclusive explicites; préserver ownership et permissions; maintenir provenance et supersession; fournir un chemin sans IA.

## 4. Non-objectifs
Aucun moteur/API/protocole/package registry/scheduler/code/schéma physique/JSON Schema/RBAC final/provider imposé, aucun nouveau Screen ID, aucune capability Endpoint et aucune autorité parallèle à Govern.

## 5. Propriétaire
CMDR Studio / assurance. Govern conserve Approval/Decision/Response Run/Result et response rollback; Settings conserve provider/integration/secret/environment admin; Shared conserve les moteurs génériques; Endpoint conserve ses primitives techniques.

## 6. Utilisateurs
Studio Asset Owner, Assurance/Release Operator et Reviewer; Platform Admin, Govern Reviewer, Auditor et consommateurs produit seulement selon leurs responsabilités et permissions source.

## 7. Conditions d’entrée
Asset/version, tenant/environment, références source et restrictions identifiables; permission de lecture puis de mutation locale si requise; Secret References seulement; aucune projection externe ne transfère ownership ou autorité.

## 8. Entrées fonctionnelles
| Entrée | Source | Requise | Fraîcheur | Si absente |
|---|---|---:|---|---|
| asset/version | Studio | oui | exact/current | blocked/incomplete if required |
| functional permission needs | Studio | oui | exact/current | blocked/incomplete if required |
| restricted inputs | Studio | oui | exact/current | blocked/incomplete if required |
| Secret References | canonical owner | selon scope | exact/current | blocked/incomplete if required |
| tenant/environment | canonical owner | selon scope | exact/current | blocked/incomplete if required |
| Tool risk | Studio | selon scope | exact/current | blocked/incomplete if required |
| Human Gate requirements | Studio | selon scope | exact/current | blocked/incomplete if required |
| Govern dependency | Studio | selon scope | exact/current | blocked/incomplete if required |
| Settings/Shared/Endpoint boundaries | Studio | selon scope | exact/current | blocked/incomplete if required |
| permission namespace anomaly | Studio | selon scope | exact/current | blocked/incomplete if required |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| version | CMDR Studio | exact/versioned authorized projection | read/link only |
| evaluation | CMDR Studio | exact/versioned authorized projection | read/link only |
| human-gate | CMDR Studio | exact/versioned authorized projection | read/link only |
| secret-reference | Platform Settings | exact/versioned authorized projection | read/link only |
| environment | Platform Settings | exact/versioned authorized projection | read/link only |
| approval | Govern | exact/versioned authorized projection | read/link only |
| decision | Govern | exact/versioned authorized projection | read/link only |
| response-run | Govern | exact/versioned authorized projection | read/link only |
| endpoint-agent | Endpoint Agent | exact/versioned authorized projection | read/link only |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Boundary Assurance Assessment | create/review/supersede as applicable | CMDR Studio | no implicit source-owner mutation |
| Permission Assurance Finding | create/review/supersede as applicable | CMDR Studio | no implicit source-owner mutation |

## 11. Fonctionnalités
asset/version; functional permission needs; restricted inputs; Secret References; tenant/environment; Tool risk; Human Gate requirements; Govern dependency; Settings/Shared/Endpoint boundaries; permission namespace anomaly; validation des dépendances; version pinning; partial/inconclusive; warnings; supersession; historique.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspecter Boundary Assurance Assessment | authorized user | Boundary Assurance Assessment | 0 | read | contexte exact | non |
| évaluer préconditions/compatibilité | reviewer | Boundary Assurance Assessment | 1 | sources lisibles | pass/warning/blocked explicable | non |
| créer/mettre à jour Boundary Assurance Assessment | owner/reviewer | Boundary Assurance Assessment | 2 | manage permission | record versionné | OPEN-013 |
| effectuer une transition pouvant produire un effet | authorized operator | Boundary Assurance Assessment | 3 | permissions + dépendances Govern/Settings applicables | transition explicite seulement | governed |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| vérifier entrées/dépendances | oui | oui | oui | explication sourcée | checklist/règles |
| comparer/résumer observations | oui | oui si structuré | oui | oui, attribuée | tables/diff/review humain |
| publier/déployer/autoriser ou inventer un résultat | accountable path only | validation only | non autonome | interdit | contrôle explicite |

L’IA ne peut jamais déclarer seule production-safe, accorder permission, contourner Govern, publier/déployer seule, inventer un résultat, masquer une régression ou révéler un secret.

## 14. États fonctionnels
`pass-with-evidence`, `warning`, `failed`, `inconclusive`, `review-required`. Lifecycle et outcome restent distincts quand les deux existent.

## 15. États d’interface
Loading conserve le contexte; Empty nomme la pièce manquante; Partial liste les sources manquantes; Error garde le dernier état valide; Offline interdit les mutations non garanties; Permission denied masque; Stale exige revalidation. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Boundary Assurance Assessment | Studio record | downstream Studio assurance/lifecycle | version/scope/limits/provenance explicit |
| Permission Assurance Finding | Studio record | downstream Studio assurance/lifecycle | version/scope/limits/provenance explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| sources canoniques | préconditions satisfaites | Boundary Assurance Assessment | versions, tenant/env, permissions, provenance | source |
| Boundary Assurance Assessment | gap/contradiction | Assurance/source owner | missing refs + limites | même lineage |
| Boundary Assurance Assessment | résultat accepté | étape Studio suivante | outcome + warnings + provenance | return origin |

## 18. Dépendances
Studio STD-1/2/3 selon les assets consommés, capabilities STD-4 adjacentes, Security permission model, Govern/Settings/Shared/Endpoint boundaries et OPEN-007, OPEN-008, OPEN-013, OPEN-015. Aucune dépendance n’est déclarée implémentée par simple référence.

## 19. Source de vérité
**Boundary Assurance Assessment** est source de sa sémantique Studio seulement. Les objets/faits externes restent chez leurs owners. Evaluation/Simulation/health/output ne deviennent pas Govern Approval/Decision/Result/Evidence automatiquement.

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
**Given** permission assurance fails, **When** readiness is assessed, **Then** the affected path is blocked or warned according to the declared gate and no permission is invented.

**Given** a Govern dependency is required but absent, **When** boundary assurance runs, **Then** the missing Approval/Decision dependency is surfaced and no Studio substitute is created.

**Given** a cross-tenant deployment is proposed without permission, **When** assurance evaluates scope, **Then** the candidate is blocked and tenant isolation remains intact.

**Given** both perm.studio and perm.cmdr-studio references exist, **When** assurance reviews namespace consistency, **Then** the anomaly is reported without bulk rename or arbitrary normalization.

## 26. Questions ouvertes
OPEN-007, OPEN-008, OPEN-013, OPEN-015 restent ouvertes. Le lot ne ferme aucune des 18 OPEN par rédaction et n’en crée aucune nouvelle.

## 27. Consommateurs documentaires
Studio Assurance/Evaluations/Simulations/Versions & Deployment, Library/Builder/Skills/Workflows/Agents/Control Room selon besoin, registers, Requirements/baseline, screen maps, Quality, Roadmap et external-owner boundaries.
