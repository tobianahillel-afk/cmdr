---
id: CAP-INV-114
title: Case Reporting Preparation
product: investigate
module: cases-and-evidence
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-PROD-018
open_decisions:
  - OPEN-014
---
# CAP-INV-114 — Case Reporting Preparation

## 1. Définition

Permet de préparer un rapport d’investigation à partir de Findings, Evidence, Timeline et objets sources, en préservant citations, redactions, versions et provenance, tout en consommant le Reporting Engine Shared sans définir l’implémentation d’export.

## 2. Problème utilisateur

Un rapport construit manuellement peut perdre les références vers les Evidence, exposer des données sensibles ou devenir obsolète après un nouveau Finding. L’analyste doit produire un brouillon révisable sans dupliquer le moteur de reporting.

## 3. Objectifs

- sélectionner Findings, Evidence et éléments temporels ;
- définir audience, portée et structure ;
- conserver citations et liens vers les objets sources ;
- gérer redactions, contenu sensible, versions, revue et commentaires ;
- préparer un brouillon exploitable par le Reporting Engine et actualisable.

## 4. Non-objectifs

- ne pas posséder le Reporting Engine ;
- ne pas définir un format d’export ou moteur de rendu final ;
- ne pas transformer un Report en Finding ;
- ne pas inclure automatiquement une Attachment ou Evidence non qualifiée.

## 5. Propriétaire

Investigate / Cases and Evidence / Investigate Product Lead possède la sélection et l’interprétation du contenu métier. Shared Capabilities possède le Reporting Engine, la mécanique de version/export et l’objet Report selon le registre.

## 6. Utilisateurs

Principal : Case Analyst ou Investigation Lead. Secondaires : Reviewer, Incident Commander, Business Owner selon audience et destinataire autorisé.

## 7. Conditions d’entrée

Case accessible, audience et portée définies, Findings/Evidence sélectionnables, permissions de données sensibles et de reporting, politique de redaction et versions sources disponibles.

## 8. Entrées fonctionnelles

| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Case | Investigate | contexte du rapport | oui | version/snapshot explicite | préparation impossible |
| Findings | Investigate | conclusions à communiquer | oui pour rapport conclusif | statut et version visibles | rester draft ou rapport de situation |
| Evidence et citations | Investigate | support et références | oui pour chaque assertion probatoire | qualification courante | marquer l’assertion non soutenue |
| Timeline selection | CAP-INV-110 / Shared | séquence et événements | non | période et timezone visibles | rapport sans chronologie détaillée |
| Audience, portée et redactions | utilisateur / politique | contrat éditorial | oui | validés avant revue | empêcher `ready-for-review` |
| Comments et Attachments | Shared Capabilities | revue et éléments documentaires | non | versions et provenance visibles | revue textuelle disponible |

## 9. Objets lus

| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case / Finding / Evidence / Artifact / Hypothesis | Investigate | contenu, statut, versions et citations | consulter et sélectionner |
| Timeline Entry | Shared Capabilities | événements, sources et timestamps | consulter et citer |
| Report | Shared Capabilities | draft, versions, audience et statut éditorial | consulter et contribuer |
| Note / Comment / Attachment | Shared Capabilities | revue, commentaires et pièces jointes | consulter selon visibilité |
| Incident / Result | Command / Govern | contexte opérationnel et résultat vérifié | consulter et citer sans modifier |

## 10. Objets créés ou modifiés

| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Report draft | créer, modifier, versionner et soumettre à revue | Shared Capabilities, contenu Investigate | classe 2 ; chaque assertion cite une source ou indique sa nature |
| Citation relation | créer ou mettre à jour | Object Linking / Reporting Engine | référence versionnée vers l’objet source |
| Redaction instruction | ajouter, modifier ou retirer avant export | Shared Capabilities / policy | permission et raison obligatoires |
| Review comment | créer ou résoudre | Shared Capabilities | version, auteur et cible visibles |

## 11. Fonctionnalités

- sélectionner et ordonner les sections du rapport ;
- choisir Findings, Evidence et Timeline ;
- générer des citations navigables ;
- définir audience, portée et niveau de détail ;
- appliquer et revoir les redactions ;
- créer, comparer et actualiser les versions ;
- gérer commentaires et validation éditoriale ;
- préparer un export futur via Reporting Engine sans en définir le format.

## 12. Actions utilisateur

| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Créer ou modifier draft | auteur autorisé | Report draft | 2 | Case et audience | version éditoriale | OPEN-013 selon politique |
| Ajouter une citation | auteur | relation source | 2 | objet accessible | citation versionnée | non |
| Appliquer une redaction | auteur/reviewer autorisé | contenu sensible | 2 | politique et raison | instruction auditée | selon sensibilité |
| Soumettre à revue | auteur | Report draft | 2 | contrôles de complétude | `in-review` | non par défaut |
| Approuver éditorialement | reviewer distinct si requis | Report version | 2 | revue terminée | `ready-for-export` conceptuel | politique locale |
| Préparer export | utilisateur autorisé | export request | 1 | version approuvée | demande au Reporting Engine | selon données |

## 13. Automatisation et IA

| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Assembler les sections | oui | template et règles | oui | brouillon | structure manuelle |
| Générer les citations | oui | oui, relations et versions | oui | non nécessaire | sélection et lien manuels |
| Résumer Findings/Timeline | oui | agrégation factuelle | oui | brouillon attribué | rédaction humaine |
| Détecter une source manquante | oui | règles de complétude | oui | explication facultative | checklist déterministe |
| Proposer une redaction | oui | politique/règles | oui | suggestion | revue manuelle des champs sensibles |

Toute sortie automatisée expose initiateur, producteur, version, sources, Automation Run ou Tool Calls, timestamp, incertitude, owner humain, rejet et trace. Aucun texte généré n’est présenté comme Finding ou citation vérifiée sans source.

## 14. États fonctionnels

`draft`, `incomplete`, `ready-for-review`, `in-review`, `changes-requested`, `editorially-approved`, `ready-for-export`, `superseded`, `withdrawn-with-trace`. Le statut de livraison/export reste distinct.

## 15. États d’interface

Loading conserve structure et sélection ; Empty propose un plan sans inventer de contenu ; Partial liste citations ou permissions manquantes ; Error préserve le draft ; Offline autorise la rédaction locale contrôlée mais pas la soumission ; Permission denied masque les sources ; Stale signale les Findings/Evidence modifiés.

## 16. Sorties

| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Report draft | Report version | reviewers et Reporting Engine | attribué, versionné et relié au Case |
| Citation set | relations versionnées | lecteurs autorisés | source, version et navigation conservées |
| Redaction plan | instructions de traitement | Reporting Engine / reviewer | permission-aware, justifié et auditable |
| Export request | événement ou job futur | Reporting Engine | non destructif, version précise et audience explicite |

## 17. Transitions

| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Case Workspace | préparer rapport | Reporting Preparation | Case, audience, scope et snapshot | Case restauré |
| Finding/Evidence/Timeline | ajouter au rapport | Report draft | objet, version, citation, rôle et return origin | objet source restauré |
| Report draft | soumettre à revue | reviewer / collaboration | version, audience, citations, redactions et commentaires | même draft après revue |
| Version approuvée | demander export | Reporting Engine | Report version, format demandé futur, audience et redactions | report avec job projeté |

## 18. Dépendances

Reporting Engine, Export, Comments and Notes, Versioning, Object Linking, CAP-INV-109, CAP-INV-110, CAP-INV-111, OPEN-014, permissions sensibles et futures décisions de format/export.

## 19. Source de vérité

Case, Findings et Evidence restent Investigate. Report, versioning et export mechanics restent Shared. Investigate est source de vérité pour la sélection et l’interprétation du contenu lié au Case.

## 20. Provenance et audit

Case ID, snapshot, auteur, audience, portée, sources et versions citées, texte automatisé, redactions, reviewers, commentaires, validations éditoriales, export request et correlation IDs.

## 21. Permissions fonctionnelles

Report prepare/review, Finding/Evidence read, sensitive-data read, redaction apply/review, citation link, export request, share et cross-tenant restrictions. Matrice atomique et step-up reportés.

## 22. Limites et erreurs

Source inaccessible, Finding contesté, Evidence superseded, citation cassée, redaction incomplète, audience non autorisée, conflit de version, Reporting Engine indisponible ou export refusé gardent le draft et exposent la cause.

## 23. Métriques

- assertions avec citations valides ;
- versions actualisées après nouveau Finding ;
- redactions revues avant export ;
- retours de revue par type ;
- navigation rapport → source réussie.

## 24. Classification de livraison

`defined` / `planned`, cible native pour la préparation Investigate et intégration Shared. Aucun moteur ou format d’export n’est déclaré livré.

## 25. Critères d’acceptation

**Given** un Case avec Findings, Evidence et Timeline  
**When** l’auteur prépare un rapport  
**Then** chaque élément sélectionné conserve sa version et une citation navigable, et le draft reste modifiable.

**Given** une Evidence superseded après la rédaction  
**When** le rapport est rouvert  
**Then** la citation est marquée stale, la nouvelle version est proposée et aucune mise à jour silencieuse n’a lieu.

**Given** aucun fournisseur de modèle  
**When** le rapport est préparé  
**Then** templates, citations, redactions et rédaction manuelle restent entièrement disponibles.

## 26. Questions ouvertes

La politique finale de Report ownership, formats d’export, redaction, validation éditoriale et relation Attachment/Artifact restent à trancher dans leurs phases propriétaires ; OPEN-014 reste ouverte.

## 27. Consommateurs documentaires

Investigation Report screen, Case Workspace, Reporting Engine, futurs parcours de revue/export, phase Permissions, objets Report et rollout des écrans.
