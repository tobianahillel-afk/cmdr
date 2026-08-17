---
id: CAP-INV-301
title: Artifact Analysis Intake and Routing
product: investigate
module: analysis-workbench
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-INV-002
  - REQ-UX-006
open_decisions:
  - OPEN-005
  - OPEN-014
source-of-truth: canonical
---
# CAP-INV-301 — Artifact Analysis Intake and Routing

## 1. Définition
Qualifier l’entrée d’un Artifact dans l’Analysis Workbench, exposer son identité, sa provenance, ses restrictions et les familles d’analyse compatibles, puis créer ou reprendre une Analysis Session sans exécuter le contenu.

## 2. Problème utilisateur
Un analyste peut ouvrir un fichier dont le type déclaré est faux, ambigu, chiffré, incomplet ou incompatible. Sans intake explicite, le produit risque d’exécuter un contenu, de masquer les restrictions ou d’envoyer l’Artifact vers un mauvais espace d’analyse.

## 3. Objectifs
- Préserver le Case, l’Artifact, l’Hypothesis éventuelle et le return origin.
- Comparer le type déclaré au type détecté et rendre les divergences visibles.
- Présenter analyses disponibles, indisponibles et futures avec leurs limites.
- Créer ou reprendre une Analysis Session uniquement après choix explicite.

## 4. Non-objectifs
- Ne pas exécuter l’Artifact pendant le routage.
- Ne pas choisir de moteur ou de format technique.
- Ne pas créer de capability dynamique, reverse, debugger ou forensic.

## 5. Propriétaire
Investigate possède le contexte et l’interprétation. Studio, Settings, Govern et Shared conservent leurs objets et mécanismes.

## 6. Utilisateurs
- Case Analyst
- Malware Analyst
- DFIR Analyst
- Evidence Reviewer

## 7. Conditions d’entrée
- Case et Artifact accessibles dans le même tenant.
- Identité, version, source et restrictions de l’Artifact disponibles ou lacunes explicites.
- Permission de lecture Artifact et d’ouverture d’analyse.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Artifact actif | CAP-INV-105 Artifact Management | objet d’analyse | oui | version courante | bloquer le lancement et conserver le Case |
| Case et Hypothesis | Investigate | contexte d’investigation | oui pour Case | état courant | ouvrir en analyse non liée interdit |
| Type déclaré et détecté | métadonnées/inspection sûre | classification fonctionnelle | oui | au moment de l’intake | marquer unknown ou ambiguous |
| Provenance et restrictions | CAP-INV-213/214 et policies | confiance et limites | oui | dernière version visible | routage restricted/incomplete |
| Analyses et Tools compatibles | module/Studio projection | options autorisées | non | health/version visibles | afficher indisponibilités |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | identité, scope, Hypothesis et liens | consulter |
| Artifact | Investigate | version, source, type, restrictions, relations | consulter |
| Tool | CMDR Studio | compatibilité, version et statut | consulter seulement |
| Sandbox Environment | Platform Settings | destination future reconnue | consulter seulement |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Analysis Session | concept Investigate | créer ou reprendre | pas de schéma final ; Case et Artifact obligatoires |
| Routing decision | Investigate | enregistrer/supersede | justification, type et limites visibles |
| Artifact | Investigate | aucune mutation du contenu | relations seulement |

## 11. Fonctionnalités
- Afficher identité, source, acquisition, versions, relations et restrictions.
- Détecter type unknown, ambiguous, unsupported, corrupted, encrypted ou incomplete.
- Présenter familles statique, dynamique future, reverse future et forensic future sans créer leurs capabilities.
- Conserver la navigation retour et les analyses déjà réalisées.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Ouvrir l’intake | Analyst | Artifact | 0 | Artifact readable | routage affiché | non |
| Choisir une famille | Analyst | Routing decision | 2 | compatibilité et permission visibles | sélection attribuée | OPEN-013 |
| Créer/reprendre session | Analyst | Analysis Session | 2 | Case et Artifact valides | session liée | OPEN-013 |
| Copier une valeur autorisée | Analyst | métadonnée | 0 | champ visible | valeur copiée | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| Suggérer une famille | oui | règles de type/compatibilité | oui | suggestion expliquée | catalogue et filtres |
| Expliquer une divergence de type | oui | comparaison explicite | oui | résumé | valeurs déclarée/détectée |
| Créer une session | oui | validation de préconditions | workflow explicite | non autonome | formulaire manuel |
| Choisir un Tool | oui | compatibilité permission-aware | oui | proposition modifiable | sélection humaine |

Toute sortie automatisée expose initiateur, producteur/version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres, timestamp, statut, incertitude, owner humain et disposition. La voie sans IA reste complète.

## 14. États fonctionnels
`recognized`, `ambiguous`, `unsupported`, `restricted`, `corrupted`, `encrypted`, `incomplete`, `ready-for-analysis`. Machines finales reportées.

## 15. États d’interface
Loading conserve le contexte; Empty explique; Partial nomme les lacunes; Error garde les données valides; Offline est stale/read-only; Permission denied ne fuit rien; Stale montre les versions.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Intake result | routing decision | Analysis Workbench | types, restrictions et limites explicites |
| Analysis Session link | relation fonctionnelle | CAP-INV-302 | Case/Artifact/owner conservés |
| Unavailable capability notice | événement de disponibilité | analyste/Settings | raison et owner visibles |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Case Workspace | ouvrir Artifact | CAP-INV-301 | tenant, Case, Artifact, Hypothesis, objectif, return origin | Case Workspace |
| CAP-INV-301 | choisir analyse statique | CAP-INV-302 | Artifact, type, provenance, restrictions, owner | intake |
| CAP-INV-301 | reconnaître future destination | future dynamic/reverse/forensic workspace | Artifact, objectifs, limites, provenance | intake |

Chaque transition conserve tenant, environnement, Case, sélection, permissions et return origin; l’erreur conserve les drafts et résultats partiels.

## 18. Dépendances
- CAP-INV-105 Artifact Management
- CAP-INV-213 Collection Integrity and Custody
- CAP-INV-214 Collection and Live Response Provenance
- CMDR Studio Tool projection
- Platform Settings environment projection
- OPEN-005/014

## 19. Source de vérité
Artifact et Case restent propriétaires Investigate ; Tool reste Studio ; environnement reste Settings. Le routage est une décision analytique Investigate, pas une mutation de l’objet source.

## 20. Provenance et audit
Enregistrer acteur, tenant, Case, Artifact/version, types déclaré et détecté, restrictions, options présentées, choix, justification, timestamp et return origin.

## 21. Permissions fonctionnelles
- Artifact read/preview
- sensitive Artifact read
- Analysis Session create
- Tool compatibility read
- cross-tenant analysis interdit par défaut

La matrice atomique est reportée.

## 22. Limites et erreurs
- Type détecté absent ou contradictoire.
- Artifact inaccessible, corrompu, chiffré ou superseded.
- Tool ou environnement incompatible ou indisponible.
- Permission révoquée entre ouverture et création de session.

## 23. Métriques
- Intakes par état de routage.
- Divergences type déclaré/détecté.
- Sessions créées ou reprises.
- Routages bloqués par restriction ou permission.

## 24. Classification de livraison
`defined` / `planned`; aucune preuve d’implémentation, moteur, hyperviseur, API, protocole, commande ou support final.

## 25. Critères d’acceptation
**Given** un Artifact lié à un Case avec types déclaré et détecté différents
**When** l’analyste ouvre l’intake
**Then** les deux types, la divergence, la provenance et les analyses compatibles sont visibles ; aucun contenu n’est exécuté

**Given** un Artifact restreint et un utilisateur sans permission
**When** il tente de créer une session
**Then** la création est refusée sans fuite et le Case reste ouvert

**Given** aucun fournisseur IA
**When** l’analyste route un Artifact reconnu
**Then** les règles déterministes et le catalogue permettent la sélection complète

## 26. Questions ouvertes
- Le statut canonique d’Analysis Session reste à décider en phase Objets.
- OPEN-005 conserve le choix futur des moteurs.
- OPEN-014 conserve Artifact versus Attachment.

## 27. Consommateurs documentaires
- Analysis Workbench README et capability map
- Case Workspace et Artifact Detail
- écrans Static Analysis et futurs workbenches
- phases Objets, Permissions et Journeys
