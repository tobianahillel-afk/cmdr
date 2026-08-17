---
id: CAP-CMD-305
title: Capability Readiness
product: command
module: readiness-and-operations
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids: [REQ-PROD-012, REQ-PROD-013, REQ-PROD-019, REQ-PROD-057]
open_decisions: [OPEN-010, OPEN-013]
source-of-truth: canonical
---

# CAP-CMD-305 — Capability Readiness

## 1. Définition
Évalue si une capability est available, partial, degraded, not-configured, not-tested, maintenance, planned ou unknown sans confondre readiness, document status et delivery mode.

## 2. Problème utilisateur
Une cible native peut être non configurée ou non livrée dans un tenant. Sans assessment séparée, le catalogue documentaire est pris pour une capacité disponible.

## 3. Objectifs
Séparer les statuts, afficher scope/source/date, lier blockers/evidence et créer une Task de validation/remediation ou ouvrir l’owner.

## 4. Non-objectifs
Ne promeut pas delivery mode, ne remplace pas Studio Assurance/Settings Health et n’invente pas de preuve.

## 5. Propriétaire
Command possède l’assessment tenant-scoped ; Capability Register et produits owners gardent delivery/config/health/evidence.

## 6. Utilisateurs
Principal : Readiness Coordinator. Secondaires : Incident Commander, Product Owner, Platform Administrator.

## 7. Conditions d’entrée
Capability ID, tenant/env et sources delivery/config/health/test.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Capability entry | Capability Register | owner, status, delivery mode et target | oui | dernière revue | `unknown`, aucune capability inventée |
| Configuration and health | produit owner ou Platform Settings | disponibilité tenant et santé | non | timestamp source | `not-configured`, `degraded` ou `unknown` selon preuve |
| Test and assurance evidence | Exercise, Studio Assurance ou validation | résultat et version | non | date/validité de preuve | `not-tested` |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Capability entry | Product Architecture / owner produit | ID, owner, delivery et cible | consulter et naviguer |
| Health / Assurance / Exercise evidence | propriétaires sources | statut, version et date | consulter en projection |
| Task | Command | validation ou remediation | consulter et modifier si autorisé |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Readiness assessment | créer ou mettre à jour statut, rationale et evidence refs | Command | classe 2, tenant-scoped et audité |
| Task | créer validation ou remediation | Command | classe 2, Capability ID et résultat attendu liés |
| Capability entry / delivery mode | aucune mutation | Capability Register / owner produit | readiness ne modifie jamais delivery classification |

## 11. Fonctionnalités
Résoudre ID/owner, afficher delivery séparément, produire assessment tenant, lier evidence/blockers et créer Task validation/remediation.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter | lecteur | Capability | 0 | registry entry | statut/source | non |
| Mettre à jour assessment | coordinateur | assessment | 2 | evidence/rationale | nouvelle assessment | OPEN-013 |
| Créer validation Task | coordinateur | Task | 2 | not-tested/unknown | Task | non |
| Ouvrir owner source | lecteur | product capability | 0 | permission | transition | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Résoudre le Capability ID | oui | oui | oui | non nécessaire | Capability Register |
| Agréger health/evidence | oui | oui | oui | résumé attribué | projections sources |
| Calculer un statut proposé | oui | mapping versionné | oui | proposition seulement | règles et validation humaine |
| Créer une Task | oui | validation/déduplication | oui | brouillon Task | création manuelle |

## 14. États fonctionnels
`available`, `partial`, `degraded`, `not-configured`, `not-tested`, `maintenance`, `planned`, `unknown`.

## 15. États d’interface
Chaque état expose scope/source/date ; Partial/Unknown ne sont jamais masqués ; Permission denied ne révèle pas evidence protégée ; Stale dégrade l’assessment.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Capability readiness | assessment Command | Readiness, Plans et Mission Control | tenant-scoped, datée et distincte du delivery mode |
| Validation/remediation Task | Task Command | owner source et Work Queue | Capability ID, evidence gap et résultat attendu liés |
| Readiness change event | événement d’audit | Reporting et historique | before/after, assessor et sources visibles |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Capability Readiness | ouvrir owner produit | produit owner | Capability ID, scope et return origin | même assessment restaurée |
| Capability Readiness | inspecter health | Platform Settings | Capability ID, tenant/env et evidence refs | assessment restaurée |
| Capability Readiness | inspecter assurance | Studio Assurance | Capability ID, test/evidence et return origin | assessment restaurée |
| Capability Readiness | gap | Improvement Actions | Capability ID, gap et résultat attendu | assessment restaurée |

## 18. Dépendances
Capability Register, CAP-CMD-301, CAP-CMD-303, Platform Health, Studio Assurance et Metrics Engine.

## 19. Source de vérité
Register possède delivery/owner ; sources possèdent health/evidence ; Command possède assessment tenant-scoped.

## 20. Provenance et audit
Capability ID, scope, evidence refs, assessor, mapping/version, rationale et date.

## 21. Permissions fonctionnelles
Command read/coordinate et source reads ; `OPEN-010/013` restent ouvertes.

## 22. Limites et erreurs
Entry absente, evidence stale, tenant mismatch, source down ou refus donnent unknown/partial, jamais available inventé.

## 23. Métriques
Capabilities avec assessment tenant, durée not-tested et assessments avec evidence ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; readiness ne modifie jamais delivery mode.

## 25. Critères d’acceptation
**Given** une capability `delivery_mode: planned`, **When** l’assessment est consultée, **Then** delivery reste planned et readiness est séparée.

**Given** une preuve expirée, **When** l’update est calculé, **Then** not-tested/partial est affiché avec source/date.

**Given** aucun modèle, **When** l’assessment est produit, **Then** registry, health et règles déterministes suffisent.

## 26. Questions ouvertes
Quelle source authoritative porte tenant availability et quels statuts seront normalisés en phase Objets ? — Requirement IDs ci-dessus ; `OPEN-010/013` restent ouvertes.

## 27. Consommateurs documentaires
Readiness, Mission Control, Plans, reports, parcours Phase 5, écrans Phase 6, phase Objets et permissions ultérieures.