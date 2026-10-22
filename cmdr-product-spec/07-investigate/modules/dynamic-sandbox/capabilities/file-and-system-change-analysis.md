---
id: CAP-INV-320
title: File and System Change Analysis
product: investigate
module: dynamic-sandbox
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-INV-005
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-OBJ-003
open_decisions:
  - OPEN-013
  - OPEN-014
source-of-truth: canonical
---
# CAP-INV-320 — File and System Change Analysis

## 1. Définition
Analyser les fichiers créés, modifiés, supprimés ou renommés et les changements de configuration ou services observés pendant un Sandbox Run, avec processus source, timestamps et comparaison avant/après.

## 2. Problème utilisateur
Une modification observée ne prouve ni persistance ni malveillance. Sans relation au Run et au processus source, elle ne peut pas être interprétée ni reproduite.

## 3. Objectifs
- Afficher changements de fichiers et système avant/après.
- Relier chaque changement au processus, au Run et à sa source.
- Capturer un Runtime Artifact sélectionné sans modifier le source.
- Préparer Evidence candidate sans qualification automatique.

## 4. Non-objectifs
Ne pas définir hooks, instrumentation, commandes, formats internes, disk forensics, persistance confirmée, moteur ou implémentation.

## 5. Propriétaire
Investigate possède l’interprétation des changements et les relations analytiques; Settings, Studio et Shared conservent leurs responsabilités.

## 6. Utilisateurs
Malware Analyst; Dynamic Analysis Operator; Case Analyst; Evidence Reviewer.

## 7. Conditions d’entrée
Run accessible; baseline ou état avant disponible ou absence explicite; observations fichiers/système et permissions accessibles.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| État avant/après | Sandbox Run | comparaison fonctionnelle | oui | même Run/version | baseline-missing |
| Événements fichiers | Tool result | créations/modifications/suppressions/renommages | oui | timestamps visibles | partial |
| Événements système | Tool result | configuration/services équivalents | non | source/version visibles | unknown |
| Processus source | CAP-INV-319 | attribution temporelle | non | même Run | source-process-missing |
| Restrictions de capture | policies/Artifact | sécurité et rétention | oui pour capture | réévaluées | capture bloquée |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Sandbox Run | Investigate concept | environnement, profil et phases | consulter |
| Process Observation | Investigate concept | processus source | consulter/lier |
| Artifact | Investigate | original et restrictions | consulter |
| Tool Call | Studio | producteur, version et output | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| System Change observation | créer/annoter/supersede | Investigate concept | source Run/process obligatoire |
| Runtime Artifact candidate | proposer/capturer | Investigate concept | CAP-INV-324 et provenance obligatoires |
| Evidence candidate relation | préparer | Investigate | qualification 107/108 requise |

## 11. Fonctionnalités
Voir fichiers créés, modifiés, supprimés, renommés, emplacements, timestamps, processus sources, changements de configuration et services, comparaison avant/après, filtres, annotations et capture sélectionnée.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Lire les changements | Analyst | Change set | 0 | Run accessible | vue sourcée | non |
| Filtrer/comparer | Analyst | Change set | 0 | baseline disponible | diff visible | non |
| Capturer Runtime Artifact | Analyst | Runtime Artifact | 1 | permission/restrictions | capture sourcée | non |
| Annoter | Analyst | System Change | 2 | observation accessible | annotation attribuée | OPEN-013 |
| Préparer Evidence candidate | Reviewer | Candidate | 2 | sources/provenance | draft uniquement | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Construire le diff | oui | comparaison déterministe | oui | explication | vue avant/après |
| Regrouper les changements | oui | type/process/path | oui | regroupement suggéré | filtres |
| Proposer une capture | oui | règles de scope | oui | suggestion modifiable | sélection humaine |
| Confirmer persistance/Evidence | humain | contrôles seulement | workflow de revue | jamais autonome | revue humaine |

## 14. États fonctionnels
`collecting`, `available`, `partial`, `baseline-missing`, `created`, `modified`, `deleted`, `renamed`, `disputed`, `superseded`.

## 15. États d’interface
Loading conserve le Run; Empty distingue aucun changement et source absente; Partial nomme les événements manquants; Error garde les résultats valides; Offline est read-only; Permission denied masque les chemins sensibles.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| File/System Change set | Analysis Result | analyste/Case | before/after, process et source |
| Runtime Artifact candidates | candidate set | CAP-INV-324 | capture explicite seulement |
| Change observations | concepts analytiques | CAP-INV-322/328 | aucun verdict de persistance |
| Comparison input | result snapshot | CAP-INV-325 | préconditions visibles |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-318/319 | ouvrir changement | CAP-INV-320 | Run, event, process, path/config ref | timeline/tree |
| CAP-INV-320 | capturer sélection | CAP-INV-324 | source event/process, méthode, restrictions | changes |
| CAP-INV-320 | grouper mécanisme candidat | CAP-INV-322 | changes, process, temporalité, contradictions | changes |
| CAP-INV-320 | préparer handoff | CAP-INV-328 | sélection, Runtime Artifacts, provenance | changes |

## 18. Dépendances
CAP-INV-317/318/319/322/324/325/328; Artifact Management; Studio Tool Calls; Shared Trace/Export; OPEN-013/014.

## 19. Source de vérité
Investigate possède la disposition analytique; Tool Call reste Studio; l’Artifact source reste inchangé; changement observé ≠ persistance confirmée.

## 20. Provenance et audit
Run, environment/version, baseline, path/config ref, process source, timestamps, before/after, Tool/version, capture, annotation et disposition.

## 21. Permissions fonctionnelles
File/system behavior read; sensitive path read; Runtime Artifact create/read/export; change annotate; Evidence candidate prepare; cross-tenant denied.

## 22. Limites et erreurs
Baseline absente; événement tronqué; chemin redacted; processus source manquant; capture invalide; fichier disparu avant capture; permission révoquée.

## 23. Métriques
Changes par type; baselines absentes; captures réussies/partielles; observations contestées; handoffs; provenance complète.

## 24. Classification de livraison
`defined` / `planned`; aucune instrumentation, commande, moteur, API, protocole, disk forensics ou implémentation.

## 25. Critères d’acceptation
**Given** un Run avec modifications fichiers et système
**When** l’analyste ouvre File and System Change Analysis
**Then** before/after, timestamps, processus sources et relations au Run sont visibles

**Given** une baseline absente
**When** la comparaison est demandée
**Then** le résultat reste partial et aucune absence de changement n’est conclue

**Given** aucun modèle IA
**When** l’analyse est réalisée
**Then** viewers, comparateurs, filtres et revue humaine couvrent le workflow

## 26. Questions ouvertes
System Change et Runtime Artifact restent des concepts; OPEN-013 et OPEN-014 restent ouvertes.

## 27. Consommateurs documentaires
Dynamic Sandbox, Persistence Analysis, Runtime Artifact Management, Multi-Run Comparison, Evidence/Finding handoff, Objets et Permissions.
