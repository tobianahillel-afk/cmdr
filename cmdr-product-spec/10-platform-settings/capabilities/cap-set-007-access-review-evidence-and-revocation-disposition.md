---
id: CAP-SET-007
title: Access Review, Evidence and Revocation Disposition
product: platform-settings
module: users-and-roles
owner: Platform Settings Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-12
requirement_ids: [REQ-PROD-006, REQ-OBJ-001, REQ-SEC-001, REQ-SEC-002, REQ-SEC-006]
open_decisions: [OPEN-007, OPEN-013, OPEN-014, OPEN-015, OPEN-019]
source-of-truth: canonical
---
# CAP-SET-007 — Access Review, Evidence and Revocation Disposition
## 1. Définition
Capability Settings qui définit l’administration périodique d’Access Review, son owner, son evidence/provenance et une disposition keep/revoke, sans inventer d’objet AccessReview ni de mécanique générique de retrait d’assignation.
## 2. Problème utilisateur
Un reviewer doit examiner l’accès visible à partir des informations Principal/Role sourcées, conserver les éléments ayant fondé son examen et produire une disposition auditable sans devenir moteur Security ou workflow Govern.
## 3. Objectifs
Définir review scope, owner, evidence/provenance, résultat administratif et handoff de révocation. En l’absence de mécanique canonique explicite de retrait Principal/Role, la révocation reste une disposition/handoff, pas une mutation inventée.
## 4. Non-objectifs
Ne pas définir AccessReview/AccessAssignment/RoleAssignment/PermissionAssignment object, Group, Permission revocation générique, effective access, Compliance engine, Approval, Decision, suppression Principal ou historique destructif.
## 5. Propriétaire
Platform Settings Product Lead est capability owner. Security reste owner de l’autorisation et des permissions; Govern reste owner d’Approval/Decision/Decision Authority; Investigate reste owner de son Evidence canonique.
## 6. Utilisateurs
Review owner ou administrateur Platform Settings autorisé; Security Administrator peut fournir/consommer le contexte Security. Aucun reviewer ne reçoit d’autorité implicite au-delà de ses permissions.
## 7. Conditions d’entrée
Review owner identifiable, Tenant résolu, Principals/Roles/références disponibles selon permissions, contexte Security courant, source `access-reviews.md` disponible et historique administratif consultable.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| review scope | administrateur/source | ensemble de références sourcées | oui | explicite | review bloquée/partielle |
| review owner | Settings | Principal/ref responsable | oui | courant | aucune clôture |
| Principal/Role information | objets canoniques | références/projections | oui selon scope | versionnée | signaler gap |
| evidence/provenance context | audit/sources | références sourcées | oui pour disposition | traçable | disposition non justifiée |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Principal | Platform Settings | id, tenant-id, state, version, relations sourcées | read |
| Role | Platform Settings | id, tenant-id, state, version, contraintes/relations sourcées | read |
| Administrative audit context | Platform Settings/Security constraints | événements/références utiles | consume |
| Permission context | Security | tenant/RBAC/ABAC/SoD/step-up | consume |
## 10. Objets créés ou modifiés
| Objet/artefact | Opération | Propriétaire | Règle |
|---|---|---|---|
| aucun nouvel objet canonique AccessReview/Assignment | none | n/a | interdit de l’inventer pour remplir le contrat |
| administrative review outcome/provenance | enregistrer comme résultat documentaire/administratif sourcé | Platform Settings | conserve scope, owner, evidence refs, disposition, justification et corrélation |
| Principal/Role | référence uniquement par défaut | Platform Settings | aucune mutation de relation sans mécanique canonique explicite |
## 11. Fonctionnalités
Lancer/inspecter une review périodique, fixer owner/scope, rassembler références autorisées, qualifier gaps, enregistrer justification, produire keep/revoke-type disposition et transmettre une revocation-required handoff quand l’exécution mécanique n’est pas sourcée.
## 12. Actions utilisateur
Class 0: inspecter review/scope/evidence. Class 1: valider complétude, tenant scope et prévisualiser disposition. Class 2: démarrer/enregistrer review et disposition lorsque sourcé. Aucun Class 3/4 créé par Settings.
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| résumer evidence/provenance | oui | oui pour agrégation | oui | oui | liste/règles déterministes |
| signaler SoD potentiel | oui | oui selon règle Security | oui | oui | règle Security + revue humaine |
| proposer disposition | oui | oui selon critères explicites | non autonome | oui | reviewer + règles |
| exécuter révocation | non autonome | seulement si source future l’autorise | non | non | handoff explicite |
## 14. États fonctionnels
La source ne définit pas de state machine canonique AccessReview; cette capability n’en invente aucune. Les états Principal/Role restent ceux de leurs objets. La review expose seulement son progrès/complétude documentaire et sa disposition sourcée.
## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied et Stale. Partial/Stale identifient précisément les références manquantes ou obsolètes et leur impact sur la disposition.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| review disposition | résultat administratif | review owner/Settings | keep/revoke-type + justification, pas mutation implicite |
| revocation-required handoff | paquet/référence documentaire | owner d’exécution canonique futur | cible/scope/evidence refs sans inventer assignment removal |
| audit/provenance outcome | événement/référence | SET-AUD-001/Security audit | acteur/tenant/review/result/correlation |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| SET-IAM-001 | lancer/ouvrir review | Access Review capability | tenant/scope/owner/refs | même contexte |
| review | keep disposition | SET-IAM-001/SET-AUD-001 | scope/result/justification | audit |
| review | revoke disposition sans mécanique sourcée | handoff vers owner canonique approprié | cible/scope/evidence/provenance | statut de handoff seulement |
## 18. Dépendances
Users & Roles, `OBJ-PRINCIPAL`, `OBJ-ROLE`, Tenant, Administrative Audit, Security Permission Model/ABAC/SoD/step-up/audit, et Govern uniquement lorsqu’une autorité Govern est réellement requise par une règle canonique séparée.
## 19. Source de vérité
`users-and-roles/access-reviews.md` pour periodic review/owner/evidence/revocation; `SET-IAM-001` pour “Lancer une revue”; Principal/Role pour références; Security pour autorisation; aucun source ne définit actuellement un assignment-removal générique.
## 20. Provenance et audit
Chaque review conserve owner, Tenant, scope, versions/références consultées, evidence/provenance refs, gaps, justification, disposition, timestamps et correlation id. Review evidence n’est pas automatiquement Investigate Evidence.
## 21. Permissions fonctionnelles
Consomme `perm.settings.identity.read/manage` et les lectures objet Principal/Role applicables. Aucune `review.execute`, `role.assign` ou Permission ID nouvelle n’est créée.
## 22. Limites et erreurs
Access Review ≠ Govern Approval/Decision/Compliance engine. Revocation disposition ≠ direct Permission revocation/Group removal/Principal deletion. Une source insuffisante produit un handoff/gap, jamais une mutation inventée.
## 23. Métriques
Reviews lancées/complètes/partielles, gaps de provenance, keep/revoke dispositions, handoffs requis et refus d’autorisation. Aucune métrique ne prouve automatiquement conformité ou effectiveness.
## 24. Classification de livraison
`defined / planned`; aucune engine de recertification, directory sync, assignment removal, conformité ou implémentation runtime n’est revendiquée.
## 25. Critères d’acceptation
**Given** un scope Principal/Role autorisé et sourcé, **When** une review est lancée, **Then** owner, Tenant, versions, evidence refs et gaps sont conservés sans créer d’objet AccessReview canonique.  
**Given** une review concluant à `keep`, **When** la disposition est enregistrée, **Then** elle reste un résultat administratif audité et ne crée aucune permission nouvelle.  
**Given** une review concluant à `revoke` alors qu’aucune mécanique canonique de retrait d’assignation n’existe, **When** la review est clôturée, **Then** une revocation-required disposition/handoff est produite et aucune relation n’est supprimée implicitement.  
**Given** l’IA indisponible, **When** une review est conduite, **Then** scope, evidence, décision humaine et handoff restent réalisables par chemin déterministe/manual.
## 26. Questions ouvertes
OPEN-007/013/014/015/019 restent ouverts selon leur domaine. Cette capability n’invente ni relation Human Gate/Approval, ni politique C2, ni rétention globale, ni provenance cross-run, ni partage externe.
## 27. Consommateurs documentaires
Users & Roles, Administrative Audit, Security, Govern quand explicitement requis, Command, Investigate, Studio, Endpoint, Shared, Quality et Roadmap; aucune consommation ne transfère ownership.
