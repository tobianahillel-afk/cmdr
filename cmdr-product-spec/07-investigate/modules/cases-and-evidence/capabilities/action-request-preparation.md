---
id: CAP-INV-113
title: Action Request Preparation
product: investigate
module: cases-and-evidence
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-PROD-016
open_decisions:
  - OPEN-007
  - OPEN-013
  - OPEN-015
---
# CAP-INV-113 — Action Request Preparation

## 1. Définition

Permet à Investigate de préparer et soumettre à Govern une Action Request complète à partir de Findings et Evidence, sans posséder son lifecycle, créer une Decision ni s’auto-approuver.

## 2. Problème utilisateur

Une conclusion technique utile à la réponse peut perdre ses sources, ses incertitudes ou son rollback lorsqu’elle est transmise sous forme de message libre. À l’inverse, Investigate ne doit pas transformer sa recommandation en autorité de décision.

## 3. Objectifs

- sélectionner Findings et Evidence pertinents ;
- décrire action, cible, impacts, urgence, alternatives, conditions et rollback attendu ;
- vérifier la complétude avant soumission ;
- gérer retour pour informations, annulation avant soumission et supersession ;
- conserver Case, provenance et lien vers la future Decision.

## 4. Non-objectifs

- ne pas créer ou approuver une Decision ;
- ne pas posséder le lifecycle de l’Action Request ;
- ne pas exécuter une action de classe 3 ou 4 ;
- ne pas transférer l’ownership des Evidence ou Findings.

## 5. Propriétaire

Investigate / Cases and Evidence / Investigate Product Lead possède la préparation analytique. Govern possède l’Action Request lifecycle, Decision, Approval, Response Run et Result.

## 6. Utilisateurs

Principal : Investigation Lead ou Senior Analyst. Secondaires : Case Analyst contributeur, Incident Commander fournissant l’impact opérationnel et reviewer Govern consommateur.

## 7. Conditions d’entrée

Case accessible, au moins un Finding ou une justification équivalente, Evidence référencées, cible résolvable, action proposée, permissions de préparation et soumission, tenant/environnement conservés.

## 8. Entrées fonctionnelles

| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Case | Investigate | contexte analytique | oui | version courante | préparation impossible |
| Findings | Investigate | conclusions justifiées | oui sauf demande d’information explicitement permise | statut et version visibles | rester `incomplete` |
| Evidence | Investigate | références et provenance | oui pour assertion soutenue | qualification courante | bloquer `ready` ou déclarer l’incertitude |
| Action, cible et conditions | utilisateur / workflow | demande structurée | oui | validées à la soumission | rester draft |
| Impact et urgence | Command / utilisateur | contexte opérationnel et technique | oui selon classe | horodatés et sourcés | demander information |
| Alternatives et rollback attendu | analyste / procédures | maîtrise du risque | requis pour action sensible | version visible | soumission bloquée selon politique |

## 9. Objets lus

| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case / Finding / Evidence | Investigate | contexte, assertions, sources et incertitudes | consulter et sélectionner |
| Incident / Service / Task | Command | impact, urgence, owner et coordination | consulter et relier |
| Action Request | Govern | draft reçu, statut, retour d’information et supersession | consulter en projection après soumission |
| Decision / Approval / Response Run / Result | Govern | traitement et résultat ultérieurs | consulter et naviguer |
| Workflow / Automation Run | CMDR Studio | préparation automatisée et provenance | inspecter sans transférer l’autorité |

## 10. Objets créés ou modifiés

| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Action Request draft | préparer, compléter, annuler avant soumission ou soumettre | Govern lifecycle, Investigate producer | classe 2 avant soumission ; aucune Decision créée |
| Case relation | lier demande, retour, Decision future et Result | Investigate / Object Linking | historique conservé, Evidence non transférées |
| Finding relation | sélectionner ou retirer avant soumission | Investigate | version et rôle visibles |
| Information response | compléter une demande retournée | Govern lifecycle, contribution Investigate | conserve demande, liens et historique |

## 11. Fonctionnalités

- sélectionner Findings et Evidence ;
- définir action, cible, impacts, urgence, alternatives, conditions, rollback et incertitudes ;
- afficher un contrôle de complétude ;
- préparer manuellement ou via workflow ;
- soumettre vers Govern ;
- traiter un retour pour informations ;
- annuler avant soumission ou supersede sans effacer l’historique ;
- suivre la relation vers Decision, Run et Result.

## 12. Actions utilisateur

| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Créer ou modifier draft | analyste autorisé | Action Request draft | 2 | Case et cible | draft versionné | OPEN-013 |
| Valider la complétude | analyste/reviewer | draft | 0 | champs accessibles | liste déterministe des lacunes | non |
| Soumettre | Investigation Lead | Action Request | 2, transition d’autorité | permissions et état `ready` | réception Govern | obligatoire |
| Annuler avant soumission | auteur autorisé | draft | 2 | non soumis | annulation auditée | non |
| Répondre à une demande d’information | analyste | request contribution | 2 | request retournée | nouvelle version soumise | Govern reçoit |
| Préparer containment | analyste | action classe 3 | 3, demande uniquement | impact et rollback | Action Request, aucune exécution | obligatoire |

## 13. Automatisation et IA

| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Assembler les références | oui | oui, liens et versions | oui | non nécessaire | sélection manuelle et checklist |
| Vérifier la complétude | oui | règles explicites | oui | explication facultative | validation déterministe |
| Proposer action ou alternatives | oui | procédures/règles | oui | brouillon seulement | procédures et expertise humaine |
| Rédiger le package | oui | template structuré | oui | brouillon attribué | formulaire manuel |
| Soumettre | oui | contrôle de permission | workflow possible | jamais autonome | action humaine explicite |

Toute sortie automatisée expose initiateur, moteur/agent, version, Automation Run, Tool Calls, sources, timestamp, statut, incertitude, owner humain, rejet et trace. Aucune suggestion ne devient Decision ou Approval.

## 14. États fonctionnels

Préparation Investigate : `draft`, `incomplete`, `ready`, `submitted`, `returned-for-information`, `superseded`, `cancelled-before-submission`. Les états Govern post-soumission restent propriété de Govern.

## 15. États d’interface

Loading conserve le draft ; Empty exige Case et objectif ; Partial nomme les Evidence ou impacts manquants ; Error préserve les données ; Offline autorise le brouillon mais pas la soumission ; Permission denied masque les sources protégées ; Stale bloque la soumission tant que les versions ne sont pas revues.

## 16. Sorties

| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Action Request draft | objet Govern en préparation | analystes Investigate | versionné, attribué et non décisionnel |
| Submission package | Action Request submitted event | Govern | Findings/Evidence référencés, impacts, alternatives et rollback présents |
| Information response | nouvelle contribution/version | Govern reviewer | demande et liens d’origine conservés |
| Case status link | relation Action Request/Decision/Result | Case Workspace | navigation bidirectionnelle et ownership préservé |

## 17. Transitions

| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Finding Review | préparer une action | Action Request Preparation | Case, Findings, Evidence, auteur et incertitudes | Finding restauré |
| Action Request draft | soumission humaine | Govern Review Queue | request, cible, action, impacts, urgence, alternatives, rollback et sources | Case avec statut Govern projeté |
| Govern | retour pour information | Case Workspace / draft | demande, questions, version et liens | resoumission vers même request |
| Govern Result | résultat disponible | Case Workspace | Result, Decision, Run, request, Findings et prochaine action | retour au Result ou Incident |

## 18. Dépendances

CAP-INV-107, CAP-INV-108, CAP-INV-109, Case Lifecycle, Object Linking, Govern Action Request/Decision/Run/Result, Command impact context, Studio workflow provenance, OPEN-007, OPEN-013 et OPEN-015.

## 19. Source de vérité

Evidence, Findings et Case restent Investigate. Action Request, Decision, Approval, Response Run et Result restent Govern. Investigate est source de vérité pour le contenu analytique préparé avant soumission, pas pour l’autorité.

## 20. Provenance et audit

Case ID, auteur, Findings/Evidence et versions, action/cible, impacts, urgence, alternatives, rollback, incertitudes, validations de complétude, automation provenance, soumission, retours, supersession et correlation IDs.

## 21. Permissions fonctionnelles

Action Request prepare/submit/cancel-draft/respond, Finding/Evidence read, sensitive-target access, cross-tenant restrictions, step-up pour soumission sensible et séparation analyste/approbateur. Matrice atomique reportée.

## 22. Limites et erreurs

Finding contesté, Evidence inaccessible, cible stale, impact absent, conflit de version, permission refusée, Govern indisponible ou changement d’environnement maintiennent un draft explicite ; aucune Decision ou exécution n’est simulée.

## 23. Métriques

- drafts atteignant `ready` sans retour ;
- retours pour informations par type de lacune ;
- packages avec alternatives et rollback ;
- délai Finding confirmé → soumission ;
- soumissions automatisées : toujours zéro sans action humaine.

## 24. Classification de livraison

`defined` / `planned`, cible native pour la préparation Investigate et intégration Govern. Promotion conditionnée par contrat d’objet, permissions, OPEN-007/013/015 et workflow Govern.

## 25. Critères d’acceptation

**Given** un Finding confirmé, des Evidence liées et une action classe 3  
**When** l’analyste soumet la demande  
**Then** Govern reçoit l’Action Request, Investigate conserve Findings/Evidence, aucune Decision ni exécution n’est créée localement et le lien retour Case est conservé.

**Given** une demande incomplète  
**When** l’utilisateur demande la validation  
**Then** les lacunes sont nommées, l’état reste `incomplete` et aucune soumission silencieuse n’a lieu.

**Given** une request retournée par Govern  
**When** l’analyste ajoute les informations  
**Then** la même request et son historique sont conservés, les liens restent valides et une nouvelle version peut être soumise.

## 26. Questions ouvertes

OPEN-007 traite Human Gate/Govern, OPEN-013 la classe 2 et OPEN-015 le bridge Automation Run/Response Run. Les permissions de soumission et le contrat final Action Request restent à finaliser.

## 27. Consommateurs documentaires

Govern Review Queue, Case Workspace, Finding Review, Command Incident Detail, futures phases Objets/Permissions, workflows Studio et parcours Finding→Govern→Result.
