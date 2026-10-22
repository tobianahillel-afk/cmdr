---
id: CAP-STD-053
title: Evaluation Suite, Cases, Inputs and Expected Outcomes
product: cmdr-studio
module: evaluations
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-10
requirement_ids: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-016, REQ-PROD-019, REQ-OBJ-009, REQ-AI-002, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-013, OPEN-014]
source-of-truth: canonical
---
# CAP-STD-053 — Evaluation Suite, Cases, Inputs and Expected Outcomes

## 1. Définition
Define reusable Evaluation Suites and Cases with version applicability, input references, conceptual fixtures/reference context, expected outcomes/errors, dependency requirements and provenance without selecting a physical test format.

## 2. Problème utilisateur
Sans **Evaluation Suite**, Studio ne peut pas distinguer de façon reproductible l’asset/version, le contexte, les preuves, les limites et le lifecycle; un score, une projection technique ou un état local pourrait alors être pris à tort pour une autorité, un succès production ou un objet d’un autre produit.

## 3. Objectifs
Version-pinner le sujet; conserver tenant/environment et sources; rendre partial/failure/warning/inconclusive explicites; préserver ownership et permissions; maintenir provenance et supersession; fournir un chemin sans IA.

## 4. Non-objectifs
Aucun moteur/API/protocole/package registry/scheduler/code/schéma physique/JSON Schema/RBAC final/provider imposé, aucun nouveau Screen ID, aucune capability Endpoint et aucune autorité parallèle à Govern.

## 5. Propriétaire
CMDR Studio / evaluations. Govern conserve Approval/Decision/Response Run/Result et response rollback; Settings conserve provider/integration/secret/environment admin; Shared conserve les moteurs génériques; Endpoint conserve ses primitives techniques.

## 6. Utilisateurs
Studio Asset Owner, Assurance/Release Operator et Reviewer; Platform Admin, Govern Reviewer, Auditor et consommateurs produit seulement selon leurs responsabilités et permissions source.

## 7. Conditions d’entrée
Asset/version, tenant/environment, références source et restrictions identifiables; permission de lecture puis de mutation locale si requise; Secret References seulement; aucune projection externe ne transfère ownership ou autorité.

## 8. Entrées fonctionnelles
| Entrée | Source | Requise | Fraîcheur | Si absente |
|---|---|---:|---|---|
| Evaluation definition | Studio | oui | exact/current | blocked/incomplete if required |
| suite scope | Studio | oui | exact/current | blocked/incomplete if required |
| case input references | Studio | oui | exact/current | blocked/incomplete if required |
| fixture/reference context | Studio | selon scope | exact/current | blocked/incomplete if required |
| expected outcomes/errors | Studio | selon scope | exact/current | blocked/incomplete if required |
| environment constraints | Studio | selon scope | exact/current | blocked/incomplete if required |
| version applicability | Studio | selon scope | exact/current | blocked/incomplete if required |
| Tool/Skill dependencies | Studio | selon scope | exact/current | blocked/incomplete if required |
| sensitive-input classification | Studio | selon scope | exact/current | blocked/incomplete if required |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| evaluation | CMDR Studio | exact/versioned authorized projection | read/link only |
| version | CMDR Studio | exact/versioned authorized projection | read/link only |
| skill | CMDR Studio | exact/versioned authorized projection | read/link only |
| workflow | CMDR Studio | exact/versioned authorized projection | read/link only |
| automation-agent | CMDR Studio | exact/versioned authorized projection | read/link only |
| secret-reference | Platform Settings | exact/versioned authorized projection | read/link only |
| environment | Platform Settings | exact/versioned authorized projection | read/link only |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Evaluation Suite | create/review/supersede as applicable | CMDR Studio | no implicit source-owner mutation |
| Evaluation Case | create/review/supersede as applicable | CMDR Studio | no implicit source-owner mutation |

## 11. Fonctionnalités
Evaluation definition; suite scope; case input references; fixture/reference context; expected outcomes/errors; environment constraints; version applicability; Tool/Skill dependencies; sensitive-input classification; validation des dépendances; version pinning; partial/inconclusive; warnings; supersession; historique.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspecter Evaluation Suite | authorized user | Evaluation Suite | 0 | read | contexte exact | non |
| évaluer préconditions/compatibilité | reviewer | Evaluation Suite | 1 | sources lisibles | pass/warning/blocked explicable | non |
| créer/mettre à jour Evaluation Suite | owner/reviewer | Evaluation Suite | 2 | manage permission | record versionné | OPEN-013 |
| effectuer une transition pouvant produire un effet | authorized operator | Evaluation Suite | 3 | permissions + dépendances Govern/Settings applicables | transition explicite seulement | governed |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| vérifier entrées/dépendances | oui | oui | oui | explication sourcée | checklist/règles |
| comparer/résumer observations | oui | oui si structuré | oui | oui, attribuée | tables/diff/review humain |
| publier/déployer/autoriser ou inventer un résultat | accountable path only | validation only | non autonome | interdit | contrôle explicite |

L’IA ne peut jamais déclarer seule production-safe, accorder permission, contourner Govern, publier/déployer seule, inventer un résultat, masquer une régression ou révéler un secret.

## 14. États fonctionnels
`draft`, `ready`, `incomplete`, `version-inapplicable`, `superseded`. Lifecycle et outcome restent distincts quand les deux existent.

## 15. États d’interface
Loading conserve le contexte; Empty nomme la pièce manquante; Partial liste les sources manquantes; Error garde le dernier état valide; Offline interdit les mutations non garanties; Permission denied masque; Stale exige revalidation. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Evaluation Suite | Studio record | downstream Studio assurance/lifecycle | version/scope/limits/provenance explicit |
| Evaluation Case | Studio record | downstream Studio assurance/lifecycle | version/scope/limits/provenance explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| sources canoniques | préconditions satisfaites | Evaluation Suite | versions, tenant/env, permissions, provenance | source |
| Evaluation Suite | gap/contradiction | Assurance/source owner | missing refs + limites | même lineage |
| Evaluation Suite | résultat accepté | étape Studio suivante | outcome + warnings + provenance | return origin |

## 18. Dépendances
Studio STD-1/2/3 selon les assets consommés, capabilities STD-4 adjacentes, Security permission model, Govern/Settings/Shared/Endpoint boundaries et OPEN-013, OPEN-014. Aucune dépendance n’est déclarée implémentée par simple référence.

## 19. Source de vérité
**Evaluation Suite** est source de sa sémantique Studio seulement. Les objets/faits externes restent chez leurs owners. Evaluation/Simulation/health/output ne deviennent pas Govern Approval/Decision/Result/Evidence automatiquement.

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
**Given** an Evaluation Case has no expected outcome, **When** the suite is reviewed, **Then** the case is marked incomplete and cannot contribute a passing result.

**Given** a Tool version changes after the suite baseline, **When** suite applicability is checked, **Then** affected cases are flagged for re-evaluation rather than silently reused.

**Given** an input is sensitive, **When** a user lacks the necessary source permission, **Then** only an authorized reference or masked projection is exposed.

**Given** no AI is available, **When** cases are authored, **Then** human-authored cases and deterministic validation remain sufficient.

## 26. Questions ouvertes
OPEN-013, OPEN-014 restent ouvertes. Le lot ne ferme aucune des 18 OPEN par rédaction et n’en crée aucune nouvelle.

## 27. Consommateurs documentaires
Studio Assurance/Evaluations/Simulations/Versions & Deployment, Library/Builder/Skills/Workflows/Agents/Control Room selon besoin, registers, Requirements/baseline, screen maps, Quality, Roadmap et external-owner boundaries.
