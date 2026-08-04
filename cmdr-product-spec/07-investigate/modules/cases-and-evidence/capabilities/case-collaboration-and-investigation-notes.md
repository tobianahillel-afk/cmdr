---
id: CAP-INV-111
title: Case Collaboration and Investigation Notes
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
  - OPEN-013
  - OPEN-014
---
# CAP-INV-111 — Case Collaboration and Investigation Notes

## 1. Définition

Applique Notes, Comments, mentions, présence et liens d’objets au contexte d’un Case afin de conserver les observations, questions et décisions d’analyse sans créer une seconde Task générale ni redéfinir les Shared Capabilities.

## 2. Problème utilisateur

Une investigation distribuée perd rapidement les raisons, les questions ouvertes et les responsabilités lorsqu’elles restent dans des messages externes ou dans des notes non reliées aux objets du Case.

## 3. Objectifs

- conserver notes, commentaires, observations, questions et décisions d’analyse ;
- attribuer les contributions et gérer visibilité, version, résolution et rétention ;
- relier les contributions aux objets du Case ;
- promouvoir une demande persistante vers Task plutôt que créer un mécanisme concurrent.

## 4. Non-objectifs

- ne pas posséder Comment, Note, Presence ou Collaboration ;
- ne pas remplacer Task opérationnelle Command ;
- ne pas décider la fusion Attachment/Artifact ;
- ne pas définir un outil de messagerie général.

## 5. Propriétaire

Investigate / Cases and Evidence / Investigate Product Lead possède l’usage métier dans le Case. Shared Capabilities reste propriétaire des mécanismes Note, Comment, Collaboration, Presence, Versioning et Attachment selon OPEN-014.

## 6. Utilisateurs

Principal : Case Analyst. Secondaires : Investigation Lead, Threat Hunter, Reviewer, Incident Commander en projection et spécialistes invités.

## 7. Conditions d’entrée

Case accessible, politique de visibilité connue, auteur identifié, permissions de contribution et de lecture des objets liés, tenant et environnement conservés.

## 8. Entrées fonctionnelles

| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Case | Investigate | contexte de collaboration | oui | version courante | contribution impossible |
| Note ou Comment | utilisateur / Shared Capabilities | contenu collaboratif | oui pour création | versionnée | rester en brouillon |
| Objets liés | Investigate / Command / Shared / Govern | références analytiques | non | résolution à l’ouverture | contribution conservée sans lien non prouvé |
| Participants et présence | Collaboration / Presence | contexte humain | non | présence temps réel ou dernière activité | collaboration asynchrone disponible |
| Attachment éventuelle | Shared Capabilities | pièce jointe | non | provenance et politique visibles | contenu textuel reste utilisable |

## 9. Objets lus

| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Note / Comment | Shared Capabilities | contenu, auteur, version, visibilité et résolution | consulter et répondre |
| Case / Hypothesis / Evidence / Finding / Artifact | Investigate | contexte et liens | consulter, citer et naviguer |
| Incident / Task | Command | contexte ou travail persistant | consulter et préparer une liaison |
| Attachment | Shared Capabilities, relation ouverte | fichier, provenance et rétention | consulter sans promouvoir automatiquement |
| Presence / Activity | Shared Capabilities | participants et activité | consulter selon permission |

## 10. Objets créés ou modifiés

| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Note / Comment | créer, modifier, répondre, résoudre ou supersede | Shared Capabilities | classe 2 ; auteur, version et visibilité obligatoires |
| Object link | ajouter ou retirer une relation réversible | Object Linking / propriétaire de l’objet | classe 2 ; aucun transfert d’ownership |
| Attachment relation | joindre ou retirer selon politique | Shared Capabilities | OPEN-014 ; aucune conversion Artifact/Evidence automatique |
| Task handoff | préparer une demande de travail persistant | Command | création via contrat Command, pas de Task locale concurrente |

## 11. Fonctionnalités

- créer des Notes et Comments liés au Case ;
- mentionner des participants et gérer visibilité ;
- relier observations, questions et décisions d’analyse aux objets ;
- versionner, résoudre, corriger et supersede sans effacer l’historique ;
- joindre un fichier sous OPEN-014 ;
- convertir une demande persistante en Task Command par transition explicite.

## 12. Actions utilisateur

| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Lire ou rechercher | utilisateur autorisé | Note/Comment | 0 | visibilité accordée | contenu attribué | non |
| Créer ou modifier | contributeur | Note/Comment | 2 | Case accessible | version auditée | OPEN-013 |
| Mentionner ou résoudre | contributeur | collaboration state | 2 | participant résolvable | activité et notification | OPEN-013 |
| Joindre un fichier | contributeur autorisé | Attachment relation | 2 | politique et provenance | pièce jointe distincte d’Artifact/Evidence | OPEN-014 |
| Préparer une Task | analyste | Task handoff | 2 | demande persistante | brouillon Command avec return origin | OPEN-013 |

## 13. Automatisation et IA

| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Résumer une discussion | oui | agrégation par thread | oui | brouillon attribué | lecture des versions et décisions |
| Suggérer des liens d’objets | oui | matching d’identifiants | oui | proposition expliquée | recherche et liaison manuelles |
| Détecter une question non résolue | oui | statut et mentions | oui | classement facultatif | filtres déterministes |
| Préparer une Task | oui | modèle de handoff | oui | brouillon | saisie manuelle structurée |

Toute sortie automatisée expose initiateur, producteur, version, sources, éventuels Tool Calls, Automation Run, incertitude, owner humain, rejet et trace. Elle ne publie ni ne résout silencieusement une contribution.

## 14. États fonctionnels

Notes et Comments peuvent être `draft`, `published`, `edited`, `resolved`, `reopened`, `superseded`, `restricted` ou `withdrawn-with-trace`. Les pièces jointes ont des états séparés de qualification ; une Attachment n’acquiert aucun statut Artifact/Evidence implicitement.

## 15. États d’interface

Loading conserve le brouillon ; Empty invite à contribuer sans créer de contenu automatique ; Partial nomme les objets ou participants indisponibles ; Error préserve le texte local ; Offline autorise le brouillon ; Permission denied ne révèle rien ; Stale signale une version concurrente.

## 16. Sorties

| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Contribution collaborative | Note ou Comment | Case Workspace et Activity Stream | attribuée, versionnée et permission-aware |
| Décision d’analyse documentée | Note liée | Hypothesis/Finding review | raison, auteur et objets cités |
| Attachment relation | relation Shared | collaboration et audit | provenance visible, non qualifiée automatiquement |
| Task handoff | brouillon ou relation | Command | contexte, demande et return origin conservés |

## 17. Transitions

| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Case Workspace | ouvrir collaboration | Notes and Comments | Case, participants, visibilité et sélection | Case restauré |
| Note/Comment | ouvrir objet lié | objet propriétaire | objet, citation, auteur et return origin | contribution restaurée |
| Demande interne | promouvoir en Task | Command Work Queue | Case, demande, owner proposé, échéance et source | retour au thread |
| Attachment | demander qualification | CAP-INV-105 ou CAP-INV-107 | fichier, origine, auteur et relation Case | Attachment inchangée |

## 18. Dépendances

Comments and Notes, Collaboration, Presence, Notifications, Versioning, Audit Hooks, Object Linking, CAP-INV-102, CAP-INV-105, CAP-INV-107, Command Task, OPEN-013 et OPEN-014.

## 19. Source de vérité

Case et objets analytiques restent Investigate ; Note, Comment, Presence et Attachment mechanics restent Shared ; Task reste Command. Investigate est source de vérité pour la signification locale des contributions dans le Case.

## 20. Provenance et audit

Case ID, auteur, mentions, version, visibilité, objet lié, Attachment source, édition avant/après, résolution, supersession, promotion Task et correlation ID.

## 21. Permissions fonctionnelles

Case collaboration read/write, mention, resolve/reopen, restricted-note access, Attachment upload/read/remove/promote, Task handoff et cross-tenant restrictions. Les permissions atomiques et la rétention finale restent reportées.

## 22. Limites et erreurs

Conflit de version, participant inaccessible, objet supprimé ou interdit, Attachment bloquée, changement de tenant, rétention expirée ou création Task refusée produisent un état visible sans perte du brouillon ni invention d’ownership.

## 23. Métriques

- contributions liées à un objet ;
- questions résolues avec raison ;
- conflits de version récupérés ;
- demandes persistantes correctement promues vers Task ;
- Attachments utilisées sans qualification implicite.

## 24. Classification de livraison

`defined` / `planned`, cible native pour l’usage Investigate et consommation des Shared Capabilities. La relation Attachment/Artifact reste conditionnée par OPEN-014.

## 25. Critères d’acceptation

**Given** une Note liée à une Evidence et visible par l’équipe  
**When** un reviewer la corrige  
**Then** une nouvelle version est créée, l’auteur et la raison sont visibles et l’ancienne version reste auditée.

**Given** un fichier joint à une Note  
**When** l’utilisateur l’utilise dans une conclusion  
**Then** il reste Attachment jusqu’à qualification explicite et n’apparaît pas automatiquement comme Artifact ou Evidence.

**Given** une demande interne persistante  
**When** l’analyste la promeut  
**Then** une Task Command est préparée avec le Case comme contexte et aucun objet Task concurrent n’est créé dans Investigate.

## 26. Questions ouvertes

OPEN-014 gouverne Attachment/Artifact ; OPEN-013 gouverne les actions réversibles. La rétention, la visibilité fine et les règles de promotion Task restent à préciser.

## 27. Consommateurs documentaires

Case Workspace, Evidence Review, Finding Review, futurs écrans Collaboration, phase Permissions, Objets, OPEN-014 et parcours Case→Task.
