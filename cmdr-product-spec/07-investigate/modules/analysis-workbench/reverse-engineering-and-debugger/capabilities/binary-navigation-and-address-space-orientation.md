---
id: CAP-INV-331
title: Binary Navigation and Address Space Orientation
product: investigate
module: analysis-workbench
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-05
requirement_ids:
  - REQ-INV-003
  - REQ-PROD-014
  - REQ-UX-002
  - REQ-UX-006
  - REQ-SEC-001
open_decisions:
  - OPEN-013
source-of-truth: canonical
---

# CAP-INV-331 — Binary Navigation and Address Space Orientation

## 1. Définition
Binary Navigation and Address Space Orientation définit le comportement produit permettant de s’orienter dans les régions, sections, segments et systèmes de repérage d’un Artifact sans confondre adresse, offset, emplacement et fonction.

## 2. Problème utilisateur
de s’orienter dans les régions, sections, segments et systèmes de repérage d’un Artifact sans confondre adresse, offset, emplacement et fonction.

## 3. Objectifs
- naviguer par adresse, offset, symbole ou référence
- distinguer les systèmes de repérage et zones non mappées
- conserver sélection, historique, retour et bookmarks
- relier un emplacement aux fonctions et données candidates sans certitude implicite

## 4. Non-objectifs
- définir un format d’adresse interne définitif
- détecter automatiquement une fonction certaine
- modifier l’Artifact source

## 5. Propriétaire
Investigate possède le contexte analytique, les annotations et les relations au Case. Studio conserve les Tools, Tool Calls et Automation Runs; Platform Settings administre les environnements; Govern conserve l’autorité sur toute cible réelle; Shared conserve les mécanismes transversaux.

## 6. Utilisateurs
Reverse Engineer principal; Malware Analyst et Reviewer secondaires.

## 7. Conditions d’entrée
- Reverse Analysis Session active ou consultable
- Artifact et structure disponibles au moins partiellement
- système de repérage déclaré ou ambiguïté visible

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Artifact layout | CAP-INV-305/307 | sections, segments, régions et métadonnées | Oui | version session | vue partielle si absent |
| Location query | utilisateur ou lien | adresse, offset, symbole ou référence | Oui | courante | afficher non résolu |
| Symbols/xrefs | CAP-INV-334 | repères et relations | Non | session active | navigation brute reste possible |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | contexte, Hypothesis et return origin | lecture et liaison uniquement |
| Artifact / Derived Artifact / Runtime Artifact | Investigate | source, versions, restrictions et provenance | lecture; aucun original modifié |
| Tool / Tool Call / Automation Run | CMDR Studio | outil, version, exécution et attribution | lecture et sélection; lifecycle Studio |
| Reverse Analysis Session | Investigate | Artifact actif, vues et sélection | lecture |
| Section/segment/location candidates | Investigate concepts | espaces, régions, adresses et offsets | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Trace / Activity event | Shared | émission de l’action et de son résultat | append-only; aucune suppression locale |
| Bookmark / location annotation | Investigate | create/update/version | n’établit pas une fonction confirmée |
| Navigation history | Shared/Investigate context | append/restore | préserve le focus et return path |

## 11. Fonctionnalités
- afficher espaces, régions, sections, segments et zones ambiguës
- naviguer et convertir visiblement entre systèmes de repérage sans les confondre
- conserver historique avant/arrière et retour de référence
- créer bookmark et annotation d’emplacement
- préserver la sélection entre vues synchronisées

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Naviguer vers un emplacement | Reverse Engineer | Code Location | 0 | repère interprétable | emplacement actif visible | Non |
| Créer un bookmark | Analyst | Bookmark | 2 | session modifiable | bookmark versionné | OPEN-013 |
| Annoter un emplacement | Analyst | Location annotation | 2 | permission d’annotation | annotation attribuée | OPEN-013 |
| Copier un identifiant autorisé | Reviewer | Location identifier | 0 | policy autorise copie | valeur avec système de repérage | Non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Inspecter l’espace d’adresses et les emplacements | Oui | Oui | Oui | Oui | viewer déterministe et navigation manuelle de l’espace d’adresses et les emplacements |
| Proposer un repère ou une région candidate | Oui | Oui | Oui | Oui | règles, heuristiques visibles, comparaison et revue humaine |
| Préparer un handoff | Oui | Oui | Oui | Oui | sélection manuelle, checklist et validation humaine |

Toute sortie automatisée expose initiateur, producteur et version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres, timestamp, statut, incertitude, owner humain et disposition acceptée, modifiée ou rejetée.

## 14. États fonctionnels
- ready
- resolving
- mapped
- partially-mapped
- unmapped
- ambiguous
- stale
- restricted

Ces états sont fonctionnels et ne constituent pas une machine d’état objet définitive.

## 15. États d’interface
- **Loading** conserve le Workbench, l’Artifact, la sélection et le return origin.
- **Empty** explique l’absence de résultat sans simuler une analyse.
- **Partial** identifie les sources, vues ou événements manquants et leurs conséquences.
- **Error** conserve les résultats valides, l’erreur et une reprise sûre.
- **Offline** limite les mutations et affiche la dernière synchronisation.
- **Permission denied** ne révèle aucune donnée protégée.
- **Stale** distingue la dernière observation connue de l’état courant.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Active location context | Code Location concept | Disassembly/Decompilation/Debugger | système de repérage, valeur et ambiguïté explicites |
| Navigation history | History entries | Technical Workbench | retour stable et focus restaurable |
| Bookmark/annotation | Knowledge record | CAP-INV-337 | auteur, source et confiance conservés |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Reverse Session | ouvrir Artifact | Binary Navigation | Artifact, layout, sélection | retour session |
| Location | suivre référence | Target location | source, target, relation, historique | retour source |
| Location | ouvrir représentation | Disassembly/Decompilation | Artifact, emplacement, sélection | retour navigation |

Chaque transition conserve tenant, environnement, Case, Artifact, sélection, permissions et return origin. Une erreur ne détruit pas la source ni les résultats déjà valides.

## 18. Dépendances
- CAP-INV-305 file structure
- CAP-INV-307 static binary analysis
- CAP-INV-334 functions/symbols/xrefs
- Technical Workbench history/back
- Shared Inspector/Object Linking

## 19. Source de vérité
L’Artifact et son layout viennent des analyses propriétaires; les Code Locations, bookmarks et annotations restent dans le contexte Investigate; le rendu appartient au Design System.

## 20. Provenance et audit
Tracer système de repérage, source de résolution, conversions affichées, bookmark/annotation, utilisateur, timestamp, version et navigation inter-vues.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up éventuel | Séparation des tâches | Owner | Phase propriétaire |
|---|---|---:|---|---|---|---|
| binary navigation | exposition de contenu sensible | 0 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| bookmark create/update | mutation réversible | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| location annotation | connaissance partagée | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| raw Artifact read | accès aux octets autorisés | 0 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |

La matrice atomique, les namespaces et le modèle RBAC/ABAC restent reportés à la phase Permissions.

## 22. Limites et erreurs
- adresse et offset toujours étiquetés
- zone non mappée visible
- conversion impossible sans valeur inventée
- historique indisponible signalé
- sélection stale marquée

## 23. Métriques
- résolutions réussies/non résolues
- retours de navigation réussis
- bookmarks actifs/superseded
- zones ambiguës consultées

## 24. Classification de livraison
Delivery status `defined`; delivery mode `planned`. La promotion exige objets et permissions approuvés, Tools et environnements évalués, contrats techniques, tests de sécurité, écrans et release evidence. Aucun moteur, debugger, produit tiers ou implémentation n’est choisi.

## 25. Critères d’acceptation
### Scénario 1
**Given** un offset non mappé
**When** l’analyste navigue
**Then** la zone est indiquée non mappée; aucune fonction n’est inventée

### Scénario 2
**Given** une référence est suivie puis retour arrière
**When** l’utilisateur revient
**Then** emplacement, vue et focus source sont restaurés

### Scénario 3
**Given** aucun service IA
**When** la navigation est utilisée
**Then** adresse, offset, symboles et historique déterministes restent disponibles

## 26. Questions ouvertes
- OPEN-013 reste ouverte; aucune décision n’est fermée par cette spécification.

Les schémas, cardinalités, machines d’état, formats d’adresse et permissions atomiques sont reportés aux phases propriétaires.

## 27. Consommateurs documentaires
- INV-REV-001
- Technical Workbench
- CAP-INV-332/333/334/339
- future Code Location object specification
