---
id: CAP-INV-324
title: Runtime Artifact Capture and Management
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
  - REQ-OBJ-003
  - REQ-PROD-020
open_decisions:
  - OPEN-013
  - OPEN-014
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-324 — Runtime Artifact Capture and Management

## 1. Définition
Recevoir, capturer, enregistrer, inspecter et relier les Runtime Artifacts produits ou modifiés pendant un Sandbox Run, notamment fichiers, contenus extraits, captures, logs et screenshots, sans remplacer l’Artifact source ni qualifier automatiquement Evidence.

## 2. Problème utilisateur
Les résultats matériels d’un Run peuvent être perdus, attribués au mauvais processus ou confondus avec le source. Sans lineage, ils ne sont ni auditables ni reproductibles.

## 3. Objectifs
- Capturer Runtime Artifacts avec Run, événement/processus source, Tool, timestamp et restrictions.
- Préserver l’Artifact source et les relations parent/enfant.
- Permettre inspection, descendants, export contrôlé et réanalyse statique.
- Préparer Evidence candidate sans qualification automatique.

## 4. Non-objectifs
Ne pas modifier le source, exécuter automatiquement un Runtime Artifact, définir stockage technique, moteur, format interne, reverse/debugger ou forensics.

## 5. Propriétaire
Investigate possède Runtime Artifact et son interprétation fonctionnelle; Studio possède Tool Calls; Settings stockage/rétention administratifs; Shared export/linking/versioning.

## 6. Utilisateurs
Malware Analyst; Dynamic Analysis Operator; Case Analyst; Evidence Reviewer.

## 7. Conditions d’entrée
Sandbox Run et source event/process identifiables; sortie autorisée; restrictions, permissions et rétention visibles.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Artifact source | CAP-INV-105 | parent original | oui | version référencée | lineage incomplete |
| Sandbox Run | CAP-INV-317 | source runtime | oui | état/phase visibles | capture orphaned |
| Événement/processus source | CAP-INV-318/319/320 | origine fonctionnelle | oui | même Run | source-missing |
| Output/capture candidate | Tool result | fichier, log, screenshot ou contenu | oui | encore résolvable | capture unavailable |
| Restrictions/rétention | Settings/policy | accès, export, conservation | oui | réévaluées | restricted |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Artifact source | Investigate | version, restrictions et lineage | consulter |
| Sandbox Run | Investigate concept | environnement, profil, timestamps | consulter |
| Behavioral/Process/System Observation | Investigate concepts | source de capture | consulter/lier |
| Tool Call | Studio | Tool/version, output et statut | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Runtime Artifact | créer/versionner/supersede/withdraw | Investigate concept | source Run/event/process obligatoire |
| Runtime lineage | créer | Investigate | parent et descendants navigables |
| Capture event | émettre | Investigate semantics/Shared Trace | méthode déclarée et timestamp |
| Evidence candidate context | préparer | Investigate | qualification 107/108 requise |

## 11. Fonctionnalités
Recevoir/capturer fichiers, contenus, captures, logs et screenshots; afficher source, Run, processus/événement, méthode, Tool/version, timestamp, transformation, restrictions, statut, Case et descendants.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Capturer résultat | Analyst/Tool | Runtime Artifact | 1 | source et permission | capture sourcée | non |
| Inspecter | Analyst | Runtime Artifact | 0 | disponible | preview/métadonnées | non |
| Annoter/reclassifier | Analyst | Runtime Artifact | 2 | justification | disposition versionnée | OPEN-013 |
| Retirer de l’usage actif | Reviewer | Runtime Artifact | 2 | raison | trace conservée | OPEN-013 |
| Exporter | Analyst | Export request | 1 | restriction/policy | package contrôlé | non |
| Réanalyser statiquement | Analyst | Runtime Artifact | 0 | compatible | CAP-INV-301 | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Détecter candidats à capturer | oui | règles de type/scope | oui | suggestion | liste manuelle |
| Construire lineage | non | relations déterministes | oui | non nécessaire | graph/metadata |
| Résumer le contenu | oui | extraction/preview | oui | résumé attribué | viewer |
| Qualifier Evidence/Finding | humain | contrôles seulement | workflow de revue | jamais autonome | CAP-INV-107/108/109 |

## 14. États fonctionnels
`proposed`, `capturing`, `available`, `partial`, `invalid`, `restricted`, `superseded`, `withdrawn-from-use`.

## 15. États d’interface
Loading conserve le Run; Empty explique; Partial montre capture incomplète; Error garde les éléments valides; Offline est read-only; Permission denied masque le contenu; restricted expose la raison sans fuite.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Runtime Artifact | concept Artifact | CAP-INV-301/105 | source Run/process/event et restrictions |
| Runtime lineage | relations | Case/Evidence Review | parent et descendants navigables |
| Capture event | business event | Trace/Activity | méthode, Tool/version et timestamp |
| Evidence candidate context | draft package | CAP-INV-328 | aucune qualification automatique |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Sandbox Run | output matériel | CAP-INV-324 | Run, event/process, output, restrictions | Run |
| CAP-INV-324 | inspection statique | CAP-INV-301 | Runtime Artifact, lineage, restrictions, Case | runtime view |
| CAP-INV-324 | préparer Evidence | CAP-INV-328 | Runtime Artifact, Run, observations, provenance | runtime view |
| Runtime Artifact | future reverse | future 4B.2B.2B | Artifact, dynamic context, provenance | runtime view |

## 18. Dépendances
CAP-INV-105/301/317/318/319/320/328; Studio Tool Calls; Settings retention/storage; Shared Linking/Versioning/Export; OPEN-013/014/015.

## 19. Source de vérité
Investigate possède le Runtime Artifact et son lineage; le source reste inchangé; Runtime Artifact ≠ Evidence; Tool output reste sourcé Studio.

## 20. Provenance et audit
Case, source Artifact/version, Run, environment/profile, event/process, capture method, Tool/version, timestamp, transformation, restrictions, status, descendants, annotations et disposition.

## 21. Permissions fonctionnelles
Runtime Artifact create/read/export; screenshot read; transcript/log read; lineage update; annotation; withdraw; Evidence candidate prepare; cross-tenant denied.

## 22. Limites et erreurs
Source event disparu; capture partielle/corrompue; contenu restreint; Tool/version manquant; stockage indisponible; export bloqué; permission révoquée.

## 23. Métriques
Captures proposées/réussies/partielles; lineage complet; exports bloqués; Runtime Artifacts réanalysés; handoffs; withdrawals.

## 24. Classification de livraison
`defined` / `planned`; aucun stockage, format, moteur, hyperviseur, API, protocole ou implémentation.

## 25. Critères d’acceptation
**Given** un fichier produit par un Run
**When** il est capturé comme Runtime Artifact
**Then** source Artifact, Run, processus/événement, Tool/version, timestamp, restrictions et lineage sont visibles

**Given** une capture partielle
**When** elle est ouverte
**Then** son statut partial est visible et elle ne devient pas Evidence

**Given** aucun modèle IA
**When** le Runtime Artifact est géré
**Then** règles, previews, relations et revue humaine couvrent le workflow

## 26. Questions ouvertes
Runtime Artifact reste un concept; OPEN-013, OPEN-014 et OPEN-015 restent ouvertes.

## 27. Consommateurs documentaires
Dynamic Sandbox, Analysis Workbench, Artifact Management, Evidence/Finding handoff, future Reverse, Objets, Trust et Permissions.
