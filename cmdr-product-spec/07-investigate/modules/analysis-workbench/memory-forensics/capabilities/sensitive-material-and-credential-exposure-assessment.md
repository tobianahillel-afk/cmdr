---
id: CAP-INV-356
title: Sensitive Material and Credential Exposure Assessment
product: investigate
module: analysis-workbench
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-05
requirement_ids:
  - REQ-INV-001
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-AI-002
  - REQ-SEC-001
  - REQ-SEC-002
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-356 — Sensitive Material and Credential Exposure Assessment

## 1. Définition
Détecter et évaluer défensivement une exposition candidate de données sensibles avec masquage par défaut, contrôle d’accès, audit et remédiation, sans méthode d’extraction exploitable.

## 2. Problème utilisateur
Sans contrôle spécifique, une analyse forensics peut révéler, copier ou exporter une valeur candidate au-delà du besoin et sans séparation des tâches.

## 3. Objectifs
- détecter une présence candidate et classer son type fonctionnel.
- masquer par défaut et séparer présence, révélation, copie et export.
- préserver source, contexte, confiance, restrictions, audit et provenance.

## 4. Non-objectifs
Aucune procédure de credential dumping, utilisation/validation de secret, commande, contournement, moteur, plugin, API, action Endpoint ou qualification automatique.

## 5. Propriétaire
Investigate possède candidate/interprétation; Security contrôle reveal/copy/export; Studio Tools/Runs; Settings policies; Shared audit/export; Govern cibles réelles.

## 6. Utilisateurs
Principal : **Sensitive Data Reviewer**; secondaires : Memory Analyst, Evidence Reviewer et Audit Analyst autorisés.

## 7. Conditions d’entrée
Image/session/profil/Tool lisibles; permissions de détection et de lecture masquée distinctes de reveal/copy/export.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Memory Image/Session | Investigate | source/scope | oui | versions liées | blocked |
| Sensitive candidate result | Tool Call | type/source/confiance | oui | Tool/version visibles | partial/failed |
| Process/context relation | CAP-INV-351..354 | attribution contextuelle | non | même image | unlinked |
| Policies/permissions | Security/Settings | masking/reveal/copy/export | oui | courant | masked/restricted |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Memory Image | Investigate | source/limites | lecture restreinte |
| Sensitive Material Candidate | Investigate concept | type, présence masquée, source | selon permission |
| Tool/Tool Call/Run | Studio | producteur/version/paramètres | lecture |
| Case/Hypothesis | Investigate | contexte | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Sensitive Material Candidate | créer/annoter/contester | Investigate concept | valeur masquée par défaut |
| Access event | émettre | Security/Shared | reveal/copy/export tracés |
| Remediation recommendation | préparer | Investigate | aucune utilisation du secret |
| Evidence candidate relation | préparer | Investigate | qualification future requise |

## 11. Fonctionnalités
- détecter une présence candidate, classer type, source, contexte et confiance.
- afficher valeur masquée par défaut et état complet/partiel/invalide.
- limiter, annoter, contester et préparer Evidence/remédiation.
- enregistrer chaque accès sensible.
- exporter uniquement sous policy stricte.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Voir présence masquée | Reviewer | Candidate | 0 | permission masked-read | valeur non révélée | non |
| Révéler | Reviewer | Candidate value | 0 | permission/step-up | révélation tracée | policy |
| Copier/exporter | Reviewer | Sensitive result | 1/2 | policy/step-up/SoD | action tracée | policy/Govern si cible |
| Annoter/contester | Reviewer | Candidate | 2 | permission | version conservée | OPEN-013 |

Classes 3/4 et usage du secret indisponibles.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| détecter/classer présence | oui | oui | oui | proposition | détection déterministe masquée |
| masquer/contrôler accès | oui | oui | oui | non nécessaire | viewer policy-aware |
| expliquer contexte | oui | oui | oui | résumé | table sourcée |
| confirmer exposition | oui | non | non | assistance | reviewer autorisé |

Attribution complète; aucune valeur envoyée au modèle sans permission explicite.

## 14. États fonctionnels
`detected`, `masked`, `restricted`, `under-review`, `invalid`, `confirmed-exposure`, `disputed`, `superseded`. `confirmed-exposure` ne signifie ni validité ni autorisation d’usage.

## 15. États d’interface
Loading/Empty/Partial/Error/Offline/Permission denied/Stale; la valeur reste masquée dans tous les états non autorisés.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Sensitive Material Candidate | result | Workbench/Case | masked by default |
| Access event | audit event | Security/Audit | reveal/copy/export attribués |
| Remediation/Evidence package | candidate | CAP-INV-362 | aucune valeur libre ni qualification automatique |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Session/process | examiner sensitive material | CAP-INV-356 | image, process, Tool, policy | source |
| Candidate | demander reveal | Security future step-up | type, justification, scope | Sensitive view |
| Candidate | handoff | CAP-INV-362 | masked reference, source, confidence | Sensitive view |

Tenant, Case, image et restrictions sont préservés.

## 18. Dépendances
CAP-INV-349/351/362, Security permission model, Shared Audit/Export, Studio, OPEN-005/008/013/014/015.

## 19. Source de vérité
Candidate/context : Investigate; access policy/audit : Security/Shared; Tool/Run : Studio.

## 20. Provenance et audit
Image/session/profil, Tool/version, type candidate, source/process, confidence, masked status, chaque reveal/copy/export, acteur, justification, timestamp, résultat et disposition.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up | Séparation | Owner | Phase |
|---|---|---|---|---|---|---|
| détecter | traitement sensible | 1 | policy | producer/reviewer | Investigate/Studio | Permissions |
| présence/masqué | information sensible | 0 | possible | least privilege | Security | Permissions |
| révéler | exposition secret | 0 | probable | demandeur/reviewer | Security | Permissions |
| copier | exfiltration | 1 | fort | double contrôle | Security | Permissions |
| exporter | diffusion externe | 1/2 | strict | approbateur séparé | Security/Shared | Permissions |

Matrice atomique et règles finales reportées.

## 22. Limites et erreurs
- aucune valeur affichée par défaut, copie libre, utilisation ou validation.
- aucune commande/procédure d’extraction ou contournement.
- candidate ≠ credential valide ≠ Evidence/Finding.
- export refusé conserve contexte sans révéler la valeur.

## 23. Métriques
- candidates masked/restricted/reviewed.
- reveal/copy/export autorisés/refusés.
- accès tracés et packages sans valeur libre.
- candidates invalidées/confirmées par humain.

## 24. Classification de livraison
`defined` / `planned`; aucune implémentation ou méthode d’extraction revendiquée.

## 25. Critères d’acceptation
### 1. Sans reveal
**Given** utilisateur sans permission **When** résultat ouvert **Then** présence/source visibles selon droits, valeur masquée, copie/export indisponibles et accès tracé.
### 2. Export refusé
**Given** policy bloque export **When** demandé **Then** aucune valeur n’est diffusée et le refus est audité.
### 3. Sans IA
**Given** aucun modèle **When** analyse réalisée **Then** détection déterministe, masking, revue et remédiation fonctionnent.

## 26. Questions ouvertes
OPEN-005/008/013/014/015 restent ouvertes; permissions/step-up/SoD finaux reportés.

## 27. Consommateurs documentaires
INV-MEM-001, Case/Evidence, Security/Audit/Export, CAP-INV-362 et futures phases Permissions/Screens.
