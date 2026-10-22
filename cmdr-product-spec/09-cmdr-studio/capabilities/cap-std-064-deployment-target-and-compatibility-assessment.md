---
id: CAP-STD-064
title: Deployment Target and Compatibility Assessment
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
# CAP-STD-064 — Deployment Target and Compatibility Assessment

## 1. Définition
Assess Studio asset/version compatibility with a deployment target context including tenant, environment, runtime/provider/integration, dependency versions, Tools/Skills/Workflow/Agent requirements and Secret References without authorizing deployment.

## 2. Problème utilisateur
Sans **Deployment Compatibility Assessment**, Studio ne peut pas distinguer de façon reproductible l’asset/version, le contexte, les preuves, les limites et le lifecycle; un score, une projection technique ou un état local pourrait alors être pris à tort pour une autorité, un succès production ou un objet d’un autre produit.

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
| asset/version | Studio | oui | exact/current | blocked/incomplete if required |
| deployment target type | Studio | oui | exact/current | blocked/incomplete if required |
| tenant | canonical owner | oui | exact/current | blocked/incomplete if required |
| environment | Studio | selon scope | exact/current | blocked/incomplete if required |
| runtime/provider references | canonical owner | selon scope | exact/current | blocked/incomplete if required |
| integration references | Studio | selon scope | exact/current | blocked/incomplete if required |
| dependency versions | Studio | selon scope | exact/current | blocked/incomplete if required |
| Tool/Skill availability | Studio | selon scope | exact/current | blocked/incomplete if required |
| Workflow/Agent requirements | Studio | selon scope | exact/current | blocked/incomplete if required |
| Secret/config references | canonical owner | selon scope | exact/current | blocked/incomplete if required |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| version | CMDR Studio | exact/versioned authorized projection | read/link only |
| deployment | CMDR Studio | exact/versioned authorized projection | read/link only |
| environment | Platform Settings | exact/versioned authorized projection | read/link only |
| integration | Platform Settings | exact/versioned authorized projection | read/link only |
| model-provider | Platform Settings | exact/versioned authorized projection | read/link only |
| secret-reference | Platform Settings | exact/versioned authorized projection | read/link only |
| workflow | CMDR Studio | exact/versioned authorized projection | read/link only |
| automation-agent | CMDR Studio | exact/versioned authorized projection | read/link only |
| endpoint-agent | Endpoint Agent | exact/versioned authorized projection | read/link only |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Deployment Compatibility Assessment | create/review/supersede as applicable | CMDR Studio | no implicit source-owner mutation |

## 11. Fonctionnalités
asset/version; deployment target type; tenant; environment; runtime/provider references; integration references; dependency versions; Tool/Skill availability; Workflow/Agent requirements; Secret/config references; validation des dépendances; version pinning; partial/inconclusive; warnings; supersession; historique.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspecter Deployment Compatibility Assessment | authorized user | Deployment Compatibility Assessment | 0 | read | contexte exact | non |
| évaluer préconditions/compatibilité | reviewer | Deployment Compatibility Assessment | 1 | sources lisibles | pass/warning/blocked explicable | non |
| créer/mettre à jour Deployment Compatibility Assessment | owner/reviewer | Deployment Compatibility Assessment | 2 | manage permission | record versionné | OPEN-013 |
| effectuer une transition pouvant produire un effet | authorized operator | Deployment Compatibility Assessment | 3 | permissions + dépendances Govern/Settings applicables | transition explicite seulement | governed |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| vérifier entrées/dépendances | oui | oui | oui | explication sourcée | checklist/règles |
| comparer/résumer observations | oui | oui si structuré | oui | oui, attribuée | tables/diff/review humain |
| publier/déployer/autoriser ou inventer un résultat | accountable path only | validation only | non autonome | interdit | contrôle explicite |

L’IA ne peut jamais déclarer seule production-safe, accorder permission, contourner Govern, publier/déployer seule, inventer un résultat, masquer une régression ou révéler un secret.

## 14. États fonctionnels
`compatible`, `compatible-with-warnings`, `incompatible`, `incomplete`, `unknown`, `blocked`. Lifecycle et outcome restent distincts quand les deux existent.

## 15. États d’interface
Loading conserve le contexte; Empty nomme la pièce manquante; Partial liste les sources manquantes; Error garde le dernier état valide; Offline interdit les mutations non garanties; Permission denied masque; Stale exige revalidation. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Deployment Compatibility Assessment | Studio record | downstream Studio assurance/lifecycle | version/scope/limits/provenance explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| sources canoniques | préconditions satisfaites | Deployment Compatibility Assessment | versions, tenant/env, permissions, provenance | source |
| Deployment Compatibility Assessment | gap/contradiction | Assurance/source owner | missing refs + limites | même lineage |
| Deployment Compatibility Assessment | résultat accepté | étape Studio suivante | outcome + warnings + provenance | return origin |

## 18. Dépendances
Studio STD-1/2/3 selon les assets consommés, capabilities STD-4 adjacentes, Security permission model, Govern/Settings/Shared/Endpoint boundaries et OPEN-008, OPEN-013, OPEN-015. Aucune dépendance n’est déclarée implémentée par simple référence.

## 19. Source de vérité
**Deployment Compatibility Assessment** est source de sa sémantique Studio seulement. Les objets/faits externes restent chez leurs owners. Evaluation/Simulation/health/output ne deviennent pas Govern Approval/Decision/Result/Evidence automatiquement.

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
**Given** a deployment target is incompatible, **When** compatibility is assessed, **Then** the candidate is blocked and no deployment is started.

**Given** a required Secret Reference is missing, **When** compatibility is checked, **Then** the missing reference is reported without revealing or inventing a secret.

**Given** a runtime/provider is unavailable, **When** compatibility is assessed, **Then** the state is incompatible or unknown rather than silently rerouted.

**Given** the target is an Endpoint Agent, **When** Studio evaluates only its asset dependency boundary, **Then** Endpoint technical deployment remains Endpoint-owned and no CAP-EPT behavior is created.

## 26. Questions ouvertes
OPEN-008, OPEN-013, OPEN-015 restent ouvertes. Le lot ne ferme aucune des 18 OPEN par rédaction et n’en crée aucune nouvelle.

## 27. Consommateurs documentaires
Studio Assurance/Evaluations/Simulations/Versions & Deployment, Library/Builder/Skills/Workflows/Agents/Control Room selon besoin, registers, Requirements/baseline, screen maps, Quality, Roadmap et external-owner boundaries.
