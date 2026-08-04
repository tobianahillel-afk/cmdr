---
id: CAP-INV-205
title: File and Directory Acquisition
product: investigate
module: collection-and-live-response
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-PROD-018
  - REQ-OBJ-004
open_decisions:
  - OPEN-008
  - OPEN-013
source-of-truth: canonical
---
# CAP-INV-205 — File and Directory Acquisition

## 1. Définition
Préparer et suivre l’acquisition bornée de fichiers, ensembles ou répertoires autorisés, en conservant métadonnées, erreurs par élément et provenance.

## 2. Problème utilisateur
Une acquisition non bornée peut saturer l’Endpoint, récupérer des données excessives ou masquer les fichiers manquants et verrouillés. L’analyste doit connaître le périmètre exact avant l’exécution et recevoir un résultat par élément.

## 3. Objectifs
- sélectionner des cibles exactes, motifs autorisés ou répertoires avec profondeur et limites fonctionnelles
- prévisualiser l’estimation disponible sans révéler un contenu non autorisé
- bloquer les scopes non bornés et afficher taille, policy, permission et impact conceptuel
- créer un Artifact distinct pour chaque sortie réussie et conserver les échecs

## 4. Non-objectifs
Ne pas fournir de commande de collecte, définir le moteur, les chemins système finaux, le format, le transport ou la gestion technique des fichiers verrouillés ; ne pas devenir un outil de déploiement.

## 5. Propriétaire
Investigate / Collection and Live Response / Investigate Product Lead possède le contexte métier, les drafts et les relations au Case. Platform Settings reste propriétaire de Fleet et Endpoint Policies ; Endpoint Agent exécute et rapporte localement ; Govern conserve l’autorité, Decision, Response Run et Result.

## 6. Utilisateurs
Principal : DFIR Analyst. Secondaires : Case Analyst, Evidence Reviewer et Investigation Lead.

## 7. Conditions d’entrée
Case et Endpoint accessibles, capacité de file collection déclarée, cible fonctionnelle saisie, profondeur/taille/volume bornés, policy et permissions validées.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Case et justification | Investigate | contexte de collecte | oui | version courante | rester draft |
| Cible fichier/répertoire/motif | analyste | scope fonctionnel | oui | revalidé au lancement | incomplete |
| Profondeur, taille et nombre maximum | analyste / policy | bornes | oui si répertoire/motif | snapshot au lancement | scope refusé |
| Endpoint/Agent capability | Endpoint Agent | faisabilité et limites | oui | dernière communication | offline/unsupported |
| Policy et permission | Settings / Security / Govern | restrictions et autorité | oui | version/snapshot visibles | denied ou policy-blocked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | objectif et restrictions | consulter |
| Endpoint / Endpoint Agent | partagé / Endpoint Agent | disponibilité et file-collection capability | consulter |
| Endpoint Policy | Platform Settings | paths/patterns autorisés et bornes | consulter uniquement |
| Collection Request / Job | Investigate / concept futur | scope, statut et résultats | préparer/suivre |
| Artifact | Investigate | doublons, versions et relations existantes | consulter/lier |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Collection Request scope | créer/modifier/supersede | Investigate | cible et bornes explicites ; aucun wildcard non borné |
| Per-item acquisition result | enregistrer réussite, missing, locked, restricted ou failed | Investigate, modèle futur | un résultat par élément |
| Artifact | créer pour chaque fichier reçu | Investigate | métadonnées, source et acquisition requises |
| Endpoint/Fleet/Policy | aucune mutation | owners externes | lecture seule |

## 11. Fonctionnalités
- sélectionner un fichier, une liste, un motif autorisé ou un répertoire
- définir profondeur, taille, nombre d’éléments et limite globale
- prévisualiser les métadonnées autorisées et l’estimation de volume
- suivre progression par élément et conserver missing, locked, restricted et partial
- annuler, retry uniquement les éléments éligibles et ouvrir les Artifacts liés

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Prévisualiser la cible | DFIR Analyst | scope | 0 | metadata read et policy | estimation bornée | non |
| Modifier les bornes | DFIR Analyst | Collection Request draft | 2 | draft modifiable | nouvelle version | OPEN-013 |
| Soumettre l’acquisition | analyste autorisé | Collection Request | 1 | scope ready et capacité disponible | Collection Job lié | selon impact |
| Annuler | analyste autorisé | Collection Job | 2 | cancellable | cancel demandé | OPEN-013 selon policy |
| Retry un élément | DFIR Analyst | per-item result | 1 | erreur retryable | tentative ciblée | selon impact |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| Suggérer des bornes | oui | defaults de policy | oui | suggestion modifiable | saisie manuelle |
| Détecter un scope non borné | oui | validateur | oui | explication | règles explicites |
| Estimer le volume | oui | métadonnées déclarées | oui | résumé | estimation déterministe |
| Regrouper les erreurs | oui | codes/catégories | oui | résumé | filtres |
| Étendre automatiquement le scope | non | interdit | non | interdit | action humaine explicite |

Toute sortie automatisée expose initiateur, producteur/version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres fonctionnels, timestamp, statut, incertitude, owner humain, acceptation/modification/rejet et trace.

## 14. États fonctionnels
`draft`, `bounded`, `queued`, `acquiring`, `partial`, `completed`, `missing`, `locked`, `restricted`, `failed`, `cancelled`. États fonctionnels, pas machine finale.

## 15. États d’interface
Loading conserve cible et bornes ; Empty distingue cible absente et résultat vide ; Partial liste chaque élément ; Error garde les Artifacts reçus ; Offline ne prétend pas démarrer ; Permission denied masque chemins/contenus protégés ; Stale exige revalidation.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Résultat par élément | acquisition result | CAP-INV-203/212 | path/pattern fonctionnel, statut et erreur visibles |
| Fichier collecté | Artifact | CAP-INV-105/107 | source, métadonnées, acquisition et Case conservés |
| Résumé de scope | projection | Case Workspace | bornes et volume réellement traités |
| Événements de progression | Background Job/Trace events | Notifications/Timeline | aucun échec masqué |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-202 | scope prêt | CAP-INV-205/203 | Case, Endpoint, cibles, bornes, policy, autorité | Collection Request |
| Élément reçu | vérification fonctionnelle | CAP-INV-105 | fichier, métadonnées, source, timestamps, erreurs | Acquisition Job |
| Élément missing/locked | retry ciblé | CAP-INV-203 | élément, cause, nouvelle limite éventuelle | même acquisition |
| Artifact | qualifier | CAP-INV-107 | Artifact, Case, provenance et raison | Artifact Detail |

## 18. Dépendances
CAP-INV-202/203/105/107/212/213/214, Endpoint Agent file collection, Settings Policy, Shared Jobs/Preview/Trace et OPEN-008/013.

## 19. Source de vérité
Investigate possède le scope métier, les résultats par élément et les Artifacts. Endpoint Agent reste source de l’exécution et des erreurs locales ; Settings reste source des restrictions.

## 20. Provenance et audit
Case, Endpoint/Agent, cibles, motif, bornes, estimation, policy/version, initiateur, résultats par élément, retries, Artifacts, accès et correlation IDs.

## 21. Permissions fonctionnelles
Endpoint read, file/directory acquisition prepare/submit/cancel/retry, raw result read, Artifact receive/export, sensitive path/content read et cross-tenant restrictions.

## 22. Limites et erreurs
Cible inexistante, fichier verrouillé, permissions OS insuffisantes, motif interdit, volume dépassé, path changed, partial transfer, Endpoint offline ou policy conflict. La cible n’est jamais élargie automatiquement.

## 23. Métriques
Scopes rejetés comme non bornés, volume estimé/réel, éléments missing/locked/restricted, partial rate, retries et Artifacts avec provenance complète.

## 24. Classification de livraison
`defined` / `planned`. Aucun moteur, commande, path technique, format ou protocole n’est choisi.

## 25. Critères d’acceptation
**Given** un fichier ciblé, une limite et une policy **When** l’acquisition s’exécute **Then** la cible reste bornée, la progression et les erreurs sont visibles et l’Artifact conserve sa source.

**Given** un répertoire contenant des éléments verrouillés **When** le Job se termine **Then** les éléments reçus et échoués sont séparés et l’état global peut être `partial`.

**Given** aucun modèle IA **When** le scope est préparé **Then** champs, bornes, preview et validateurs déterministes suffisent.

## 26. Questions ouvertes
OPEN-008 conserve les différences de plateforme ; OPEN-013 la gouvernance de cancel/retry classe 2. Le moteur et le traitement technique des fichiers restent hors phase.

## 27. Consommateurs documentaires
Collection Job, Artifact Management, Evidence Creation, Custody, Endpoint Agent file collection, Case Workspace et phases Objets/Permissions/Technique.
