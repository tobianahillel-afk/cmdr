---
id: CAP-INV-323
title: Interaction and Scenario Profiles
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
  - REQ-AI-002
open_decisions:
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-323 — Interaction and Scenario Profiles

## 1. Définition
Sélectionner un profil fonctionnel d’interaction ou de scénario autorisé — aucune interaction, activité utilisateur contrôlée, navigation, attente, réseau simulé ou paramètres régionaux et temporels — sans définir les scripts ou commandes internes.

## 2. Problème utilisateur
Certains Artifacts n’exposent un comportement qu’après une interaction. Si le profil est invisible ou lancé automatiquement, l’analyste ne peut pas distinguer comportement observé et simulation.

## 3. Objectifs
- Présenter objectif, étapes fonctionnelles, limites, durée et environnement requis.
- Distinguer interaction simulée et comportement réel.
- Conserver profil/version, paramètres et disposition humaine.
- Exiger un lancement de Run explicite après sélection.

## 4. Non-objectifs
Ne pas définir scripts, commandes, automatisation interne, instrumentation, moteur, hyperviseur, reverse, debugger ou action sur cible réelle.

## 5. Propriétaire
Studio/Settings possèdent la définition et l’administration de leurs profils; Investigate possède la sélection dans le contexte analytique et son interprétation.

## 6. Utilisateurs
Malware Analyst; Dynamic Analysis Operator; Case Analyst; Reviewer.

## 7. Conditions d’entrée
Dynamic Analysis Session et Artifact accessibles; environnement compatible; catalogue permission-aware; policies et restrictions visibles.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Objectif d’analyse | CAP-INV-314/316 | besoin fonctionnel | oui | version de session | rester draft |
| Profils autorisés | Studio/Settings projection | catalogue versionné | oui | version/status visibles | aucun profil disponible |
| Environnement requis | Platform Settings | compatibilité | oui | health/version courants | incompatible |
| Limites et durée | profil/policy | scope d’interaction | oui | snapshot du Run | policy-blocked |
| Suggestion automatisée | Studio | proposition attribuée | non | sources/Tool Calls visibles | sélection manuelle |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Dynamic Analysis Session | Investigate concept | objectif, Artifact et Runs | consulter/modifier localement |
| Interaction Profile | Studio/Settings | version, étapes fonctionnelles, limites | consulter/sélectionner |
| Sandbox Environment | Platform Settings | compatibilité et health | consulter |
| Automation Run / Tool Calls | Studio | provenance de suggestion | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Profile selection | créer/supersede | Investigate | profil/version et justification |
| Dynamic Analysis Session | lier sélection | Investigate concept | aucune mutation du profil source |
| Simulated interaction event | émettre | Investigate semantics | clairement marqué simulated |

## 11. Fonctionnalités
Consulter et sélectionner profils; voir objectif, étapes, limites, durée, environnement, actions simulées et provenance; modifier seulement des paramètres fonctionnels autorisés; comparer versions.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter profils | Analyst | Profiles | 0 | catalogue autorisé | options visibles | non |
| Sélectionner profil | Analyst | Selection record | 2 | compatibilité/policy | choix attribué | OPEN-013 |
| Modifier paramètre fonctionnel | Analyst | Selection draft | 2 | champ autorisé | nouvelle version | OPEN-013 |
| Accepter/modifier/rejeter suggestion | Analyst | Suggestion | 2 | provenance visible | disposition tracée | OPEN-013 |
| Lancer avec profil | Operator | Sandbox Run | 1 | confirmation explicite | CAP-INV-317 | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Filtrer profils compatibles | oui | matrice de compatibilité | oui | explication | filtres/catalogue |
| Proposer un profil | oui | règles/checklists | oui | suggestion modifiable | sélection manuelle |
| Résumer les étapes | oui | métadonnées versionnées | oui | résumé attribué | description brute |
| Lancer le Run | humain explicite | jamais silencieux | workflow avec gate | jamais autonome | action humaine |

## 14. États fonctionnels
`available`, `incompatible`, `restricted`, `selected`, `modified`, `superseded`, `simulation-partial`, `failed`.

## 15. États d’interface
Loading conserve la session; Empty explique; Partial montre les étapes indisponibles; Error garde les profils valides; Offline est read-only; Permission denied masque les profils interdits.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Profile selection | selection record | CAP-INV-316/317 | profil/version, objectif et limites |
| Simulated interaction events | behavior events | CAP-INV-318 | distincts du comportement réel |
| Profile incompatibility | availability event | analyste/Settings/Studio | cause et owner visibles |
| Suggestion disposition | business event | Trace | accept/modify/reject |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-314/316 | profil requis | CAP-INV-323 | Artifact, objectif, environnement, restrictions | session |
| CAP-INV-323 | profil sélectionné | CAP-INV-317 | profil/version, paramètres, limites, initiateur | profile view |
| Studio suggestion | ouvrir | CAP-INV-323 | proposition, rationale, sources, Tool Calls | Studio/Session |

## 18. Dépendances
CAP-INV-314/315/316/317/318; Settings Environments; Studio Profile/Tool/Automation Run; Shared Trace/Versioning; OPEN-013/015.

## 19. Source de vérité
La définition/version du profil reste chez son owner; Investigate conserve seulement le choix, les paramètres autorisés et la disposition analytique.

## 20. Provenance et audit
Session, Artifact, environment/version, profile/version, paramètres, initiateur, suggestion, Automation Run/Tool Calls, accept/modify/reject et Run lié.

## 21. Permissions fonctionnelles
Interaction profile read/select; functional parameter update; automated analysis request; simulated event read; Run start séparé; cross-tenant denied.

## 22. Limites et erreurs
Profil retiré/incompatible; environnement indisponible; paramètre interdit; version manquante; suggestion sans provenance; policy change avant Run.

## 23. Métriques
Profils sélectionnés/modifiés; incompatibilités; suggestions acceptées/modifiées/rejetées; Runs lancés explicitement; simulation partial.

## 24. Classification de livraison
`defined` / `planned`; aucun script, commande, moteur, hyperviseur, API ou implémentation.

## 25. Critères d’acceptation
**Given** un profil compatible et versionné
**When** l’analyste le sélectionne
**Then** objectif, étapes, limites, durée, environnement et simulation sont visibles avant tout Run

**Given** une suggestion IA
**When** l’analyste l’ouvre
**Then** sources, agent/version et Tool Calls sont visibles et aucune exécution ne démarre

**Given** aucun modèle IA
**When** un profil est sélectionné
**Then** catalogue, filtres et revue humaine couvrent le workflow

## 26. Questions ouvertes
Interaction Profile reste à formaliser; OPEN-013 et OPEN-015 restent ouvertes.

## 27. Consommateurs documentaires
Dynamic Sandbox, Run Management, Behavioral Timeline, Settings, Studio, Objets, Permissions et Journeys.
