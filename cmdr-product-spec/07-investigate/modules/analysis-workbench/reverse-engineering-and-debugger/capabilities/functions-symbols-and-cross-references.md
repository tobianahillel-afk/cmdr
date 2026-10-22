---
id: CAP-INV-334
title: Functions, Symbols and Cross-References
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
  - REQ-PROD-020
  - REQ-AI-002
  - REQ-UX-006
open_decisions:
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---

# CAP-INV-334 — Functions, Symbols and Cross-References

## 1. Définition
Functions, Symbols and Cross-References définit le comportement produit permettant de distinguer fonctions détectées ou proposées, symboles importés ou attribués et références entrantes ou sortantes tout en conservant l’incertitude et l’historique.

## 2. Problème utilisateur
de distinguer fonctions détectées ou proposées, symboles importés ou attribués et références entrantes ou sortantes tout en conservant l’incertitude et l’historique.

## 3. Objectifs
- afficher fonctions, callers, callees, symboles et xrefs
- relier données, chaînes et ressources référencées
- naviguer, filtrer, comparer, annoter et renommer
- marquer une fonction ambiguë et conserver l’évolution de l’interprétation

## 4. Non-objectifs
- présenter un symbole importé ou proposé comme vérité absolue
- transformer une xref en dépendance métier
- modifier l’Artifact source ou définir les algorithmes de détection

## 5. Propriétaire
Investigate possède le contexte analytique, les annotations et les relations au Case. Studio conserve les Tools, Tool Calls et Automation Runs; Platform Settings administre les environnements; Govern conserve l’autorité sur toute cible réelle; Shared conserve les mécanismes transversaux.

## 6. Utilisateurs
Reverse Engineer principal; Malware Analyst, Reviewer et Detection Engineer viewer secondaires.

## 7. Conditions d’entrée
- Reverse Analysis Session active
- au moins une représentation statique disponible
- sources de symbole ou de détection attribuées

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Function candidates | Disassembly/Decompiler Tools | boundaries et confiance | Oui | Tool/version active | liste partielle avec source manquante |
| Symbols | Artifact/Tool/debug data | imports, exports, noms proposés | Non | session active | emplacements restent accessibles |
| Cross-references | Static analysis result | relations code/données | Non | version active | relation unresolved visible |
| Strings/resources | CAP-INV-306 | contenu référencé | Non | source disponible | lien absent sans donnée inventée |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | contexte, Hypothesis et return origin | lecture et liaison uniquement |
| Artifact / Derived Artifact / Runtime Artifact | Investigate | source, versions, restrictions et provenance | lecture; aucun original modifié |
| Tool / Tool Call / Automation Run | CMDR Studio | outil, version, exécution et attribution | lecture et sélection; lifecycle Studio |
| Function/Symbol/Xref candidates | Investigate concepts / Tool results | détections, propositions et relations | lecture |
| Strings/resources/data references | CAP-INV-306/307 | contenus liés | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Trace / Activity event | Shared | émission de l’action et de son résultat | append-only; aucune suppression locale |
| Function interpretation | Investigate | confirm/ambiguous/conflicting/supersede | confirmation humaine explicite |
| Symbol/rename record | Investigate | create/update/version | ne modifie pas le binaire |

## 11. Fonctionnalités
- cataloguer fonctions détectées, proposées et confirmées par analyste
- afficher symboles disponibles, importés et attribués
- explorer xrefs, callers, callees, données, chaînes et ressources
- filtrer, comparer et naviguer sans perte de contexte
- renommer et marquer ambigu/conflicting avec historique

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter une fonction | Reverse Engineer | Function candidate | 0 | résultat disponible | sources et confiance visibles | Non |
| Confirmer ou marquer ambiguë | Reviewer | Function interpretation | 2 | éléments examinés | statut humain attribué | OPEN-013 |
| Renommer un symbole | Analyst | Symbol knowledge | 2 | session modifiable | nouveau nom versionné | OPEN-013 |
| Suivre une xref | Analyst | Cross Reference | 0 | target résolvable ou unresolved | navigation et retour conservés | Non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Inspecter les fonctions, symboles et xrefs | Oui | Oui | Oui | Oui | viewer déterministe et navigation manuelle de les fonctions, symboles et xrefs |
| Proposer une fonction, un nom ou une relation pertinente | Oui | Oui | Oui | Oui | règles, heuristiques visibles, comparaison et revue humaine |
| Préparer un handoff | Oui | Oui | Oui | Oui | sélection manuelle, checklist et validation humaine |

Toute sortie automatisée expose initiateur, producteur et version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres, timestamp, statut, incertitude, owner humain et disposition acceptée, modifiée ou rejetée.

## 14. États fonctionnels
- detected
- proposed
- confirmed-by-analyst
- ambiguous
- conflicting
- superseded
- withdrawn

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
| Function catalog | Function concepts | CAP-INV-335/337/338/339 | chaque entrée conserve source et confiance |
| Symbol/xref index | Symbol/Cross Reference concepts | Navigation and graphs | relations techniques, jamais métier implicite |
| Interpretation history | Versioned knowledge | Reviewer/Case | ancien état consultable |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Disassembly/Decompilation | sélection fonction | Function detail | Artifact, emplacement, sources, confiance | retour vue source |
| Function | ouvrir xrefs | Xref view | callers, callees, data/string/resource refs | retour fonction |
| Function | préparer analyse runtime | Debugger Session | Artifact/copie, emplacement, objectif, restrictions | retour Reverse Session |

Chaque transition conserve tenant, environnement, Case, Artifact, sélection, permissions et return origin. Une erreur ne détruit pas la source ni les résultats déjà valides.

## 18. Dépendances
- CAP-INV-306/307
- CAP-INV-331/332/333
- CAP-INV-335/337/338
- Shared Graph/Linking/Versioning
- OPEN-013/015

## 19. Source de vérité
Les résultats de détection viennent des Tools; les confirmations, ambiguïtés, noms et annotations sont des connaissances Investigate; l’Artifact reste intact.

## 20. Provenance et audit
Tracer source de chaque fonction/symbole/xref, Tool/version, confiance, auteur, renommages, confirmation/ambiguïté, comparaisons et navigation.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up éventuel | Séparation des tâches | Owner | Phase propriétaire |
|---|---|---:|---|---|---|---|
| function/symbol/xref read | contenu sensible | 0 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| function confirmation | qualification analytique | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| symbol rename | mutation réversible | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| cross-reference review | relation incertaine | 0 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |

La matrice atomique, les namespaces et le modèle RBAC/ABAC restent reportés à la phase Permissions.

## 22. Limites et erreurs
- fonction candidate non certaine
- symbole importé peut être trompeur ou stale
- xref unresolved visible
- conflit de renommage non écrasé
- aucune dépendance métier inférée

## 23. Métriques
- fonctions par statut
- xrefs resolved/unresolved
- renommages et rollbacks
- conflits d’interprétation

## 24. Classification de livraison
Delivery status `defined`; delivery mode `planned`. La promotion exige objets et permissions approuvés, Tools et environnements évalués, contrats techniques, tests de sécurité, écrans et release evidence. Aucun moteur, debugger, produit tiers ou implémentation n’est choisi.

## 25. Critères d’acceptation
### Scénario 1
**Given** une fonction détectée avec confiance limitée
**When** l’analyste l’ouvre
**Then** elle reste detected/proposed et ses sources sont visibles

### Scénario 2
**Given** un renommage concurrent
**When** le reviewer compare les versions
**Then** les deux noms et auteurs restent visibles avant résolution

### Scénario 3
**Given** aucune IA
**When** l’analyste explore
**Then** détection Tool, symboles, xrefs, navigation et annotation manuelle restent disponibles

## 26. Questions ouvertes
- OPEN-013 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-015 reste ouverte; aucune décision n’est fermée par cette spécification.

Les schémas, cardinalités, machines d’état, formats d’adresse et permissions atomiques sont reportés aux phases propriétaires.

## 27. Consommateurs documentaires
- INV-REV-001
- Functions/Xrefs views
- CAP-INV-331/332/333/335/337/338/339
- future Function/Symbol/Xref object specifications
