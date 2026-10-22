---
id: CAP-CMD-301
title: Readiness Overview
product: command
module: readiness-and-operations
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids: [REQ-PROD-005, REQ-PROD-013, REQ-PROD-021, REQ-PROD-057]
open_decisions: [OPEN-010, OPEN-013]
source-of-truth: canonical
---

# CAP-CMD-301 — Readiness Overview

## 1. Définition
Présente la préparation opérationnelle par capability, scénario, équipe, plan, exercice et action d’amélioration, avec owner, échéance, validation et lacunes.

## 2. Problème utilisateur
La présence d’une capability documentée ne prouve ni configuration, test ou disponibilité tenant. Sans vue Readiness, les lacunes sont découvertes pendant l’Incident.

## 3. Objectifs
Distinguer readiness et delivery, montrer sources/date/scope, relier les gaps à des Tasks et ouvrir l’owner technique sans dupliquer Assurance/Health.

## 4. Non-objectifs
Ne certifie pas un agent, n’administre pas la plateforme, ne promeut pas le delivery mode et ne confond pas readiness/document status.

## 5. Propriétaire
Command possède l’assessment opérationnel ; Studio, Settings et produits sources possèdent assurance, configuration et health.

## 6. Utilisateurs
Principal : Readiness Coordinator. Secondaires : Incident Commander, Team Lead, Business Owner.

## 7. Conditions d’entrée
Scope/scénario, projections de capabilities/plans/exercices/health/tasks et permission.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Capability inventory | Capability Register | owner, delivery et cible | oui | dernière revue | état `unknown` |
| Plans, exercises and Tasks | Command | preuves Readiness | non | version courante | gap visible, jamais `ready` inventé |
| Health and assurance | Settings, Studio et produits | références de preuve | non | timestamp/version source | `not-tested` ou `unavailable` |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Capability entry | Product Architecture / owner produit | owner, delivery status/mode et cible | consulter et relier |
| Operational Plan / Exercise | Command Readiness | état, owner et résultats | consulter et naviguer |
| Task | Command | amélioration et échéance | consulter et modifier si autorisé |
| Health / Assurance evidence | Settings / Studio / produit source | statut et preuve datée | consulter en projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Readiness assessment | créer ou mettre à jour statut, scope et rationale | Command | classe 2, source et date obligatoires |
| Task | créer une action d’amélioration | Command | classe 2, gap et validation attendue liés |
| Capability delivery classification | aucune mutation | Capability Register / owner produit | readiness ne promeut jamais delivery mode |

## 11. Fonctionnalités
Afficher readiness par scope, séparer delivery/document status, identifier gaps/dépendances, lier action/owner et ouvrir Studio Assurance ou Settings Health.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter | lecteur | assessment | 0 | scope | statut et sources | non |
| Créer Task | coordinateur | Task | 2 | gap et résultat attendu | Task | OPEN-013 |
| Mettre à jour assessment | coordinateur | assessment | 2 | source et rationale | statut mis à jour | OPEN-013 |
| Ouvrir owner source | lecteur | capability source | 0 | permission | transition | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Agréger les preuves | oui | oui | oui | résumé attribué | projections et références sources |
| Calculer un statut proposé | oui | mapping versionné | oui | proposition seulement | règles Readiness et validation humaine |
| Détecter les gaps | oui | oui | oui | explication possible | comparaison déterministe |
| Créer une action | oui | validation/déduplication | oui | brouillon Task | création manuelle de Task |

## 14. États fonctionnels
`ready`, `partial`, `degraded`, `not-tested`, `unavailable`, `planned`, `unknown`.

## 15. États d’interface
Partial/Unknown montrent sources manquantes ; Offline conserve la dernière assessment datée ; Permission denied masque la preuve protégée ; Stale bloque la validation.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Readiness assessment | record Command | Mission Control, Plans et Reporting | scope, source, date et assessor visibles |
| Readiness gap | événement fonctionnel | Readiness Coordinator et owner source | attribué et distinct du delivery mode |
| Improvement Task | Task Command | Work Queue | owner, due et validation attendue liés |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Readiness Overview | détail assurance | Studio Assurance | Capability ID, scope, evidence refs et return origin | même assessment restaurée |
| Readiness Overview | détail health/configuration | Platform Settings | Capability ID, tenant/env et source | même assessment restaurée |
| Readiness Overview | gap à traiter | Improvement Actions | gap, source, owner et résultat attendu | Readiness Overview restauré |

## 18. Dépendances
CAP-CMD-203, CAP-CMD-303, CAP-CMD-304, CAP-CMD-305, Studio Assurance, Platform Health et Metrics Engine.

## 19. Source de vérité
Assessment : Command. Evidence : propriétaire source. Delivery mode/status : Capability Register.

## 20. Provenance et audit
Scope, source/evidence, date, assessor, rationale, before/after, tenant et correlation ID.

## 21. Permissions fonctionnelles
Command read/coordinate/task manage et source-specific reads ; `OPEN-010` et `OPEN-013` restent ouvertes.

## 22. Limites et erreurs
Evidence stale/absente, capability inconnue, tenant mismatch, health down ou refus donnent partial/unknown, jamais ready inventé.

## 23. Métriques
Capabilities avec source/date, gaps sans owner/due et délai gap→closure validée ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; l’assessment n’est pas une preuve de livraison.

## 25. Critères d’acceptation
**Given** une capability planned non testée, **When** Readiness est consultée, **Then** delivery mode et readiness restent distincts.

**Given** un gap, **When** une Task est créée, **Then** owner, validation attendue et source sont liés.

**Given** aucun modèle, **When** la vue est utilisée, **Then** sources, règles et Tasks suffisent.

## 26. Questions ouvertes
Quel record porte assessment et quels scénarios sont obligatoires par tenant ? — Requirement IDs ci-dessus ; `OPEN-010` et `OPEN-013` restent ouvertes.

## 27. Consommateurs documentaires
Readiness screen, Mission Control, reports, parcours Phase 5, écrans Phase 6, phase Objets et permissions ultérieures.