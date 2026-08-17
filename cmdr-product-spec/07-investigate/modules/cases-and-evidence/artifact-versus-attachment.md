---
id: investigate-artifact-versus-attachment
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: proposal
requirements:
  - REQ-PROD-061
open_decisions:
  - OPEN-014
---

# Artifact versus Attachment — OPEN-014

## 1. Titre et statut

**Artifact versus Attachment.** Document de proposition `draft`; aucune option n’est approuvée et aucune migration n’est autorisée par ce document.

## 2. Décision OPEN associée

`OPEN-014 — Artifact versus Attachment` reste ouverte. `OPEN-009` reste la seule décision historiquement résolue du registre.

## 3. Owner

Investigate Product Lead prépare l’analyse fonctionnelle. Product Architecture et les owners Shared Capabilities, Objets, Trust, Reporting et UX doivent participer à la décision finale.

## 4. Phase cible

Décision cible : future phase Objets et validation transversale, après inventaire des usages réels et des exigences de rétention, permissions, provenance et migration.

## 5. Caractère bloquant

OPEN-014 bloque le contrat objet définitif, les permissions atomiques, certaines règles de rétention et la conversion UX. Elle ne bloque pas l’usage prudent d’Attachments ni l’enregistrement d’Artifacts distincts.

## 6. Contexte

Le corpus utilise `Artifact` pour les éléments techniques analysables et `Attachment` pour les fichiers joints à des Notes, Comments, Reports ou communications. Les mêmes octets peuvent apparaître dans les deux usages sans prouver une identité de lifecycle.

## 7. Problème

Une fusion prématurée ferait perdre la distinction documentaire/analytique. Une séparation absolue pourrait dupliquer contenu, stockage, permissions et références. La décision doit porter sur l’identité fonctionnelle, pas seulement sur le format fichier.

## 8. Objets concernés

Artifact, Attachment, Case, Note, Comment, Report, Evidence, Finding, Timeline Entry, provenance record, export package et les futurs objets de collection 4B.2.

## 9. Usages actuels d’Artifact

- fichier, archive, binaire, capture réseau, image mémoire/disque, log exporté ou document analysable ;
- enregistrement lié à un Case ;
- version, dérivé et relation vers un workbench ;
- source possible d’Evidence après qualification explicite ;
- provenance technique et restrictions d’accès.

## 10. Usages actuels d’Attachment

- fichier joint à une Note ou un Comment ;
- pièce de communication ou de revue ;
- document ajouté à un Report draft ;
- support collaboratif qui n’est pas nécessairement analysable ni probatoire ;
- objet ou mécanisme encore absent du registre canonique.

## 11. Chevauchements

Les deux peuvent représenter un fichier, nécessiter upload, preview, versioning, scan de sécurité, rétention et export. Un Attachment peut être promu vers Artifact ; un Artifact peut être cité ou joint sans devenir un nouvel objet.

## 12. Différences fonctionnelles

Artifact porte un rôle analytique, un contexte d’acquisition et des dérivés. Attachment porte d’abord une relation à un conteneur collaboratif ou documentaire. Le même contenu peut donc avoir des rôles différents et des lifecycles distincts.

## 13. Provenance

Artifact exige origine, acquisition/import, auteur ou moteur, Case, transformations et dérivés. Attachment exige au minimum auteur, conteneur, timestamp, version et source d’upload. Une promotion doit conserver les deux chaînes et leur relation.

## 14. Permissions

Artifact peut nécessiter accès aux données techniques sensibles, workbench, dérivés et export. Attachment dépend de la visibilité de Note/Comment/Report. Une permission de lecture Attachment ne doit pas accorder automatiquement l’accès analytique ou Evidence.

## 15. Rétention

La rétention peut dépendre du Case, du legal hold, du conteneur, de la classification et du statut Evidence. La suppression logique d’une relation Attachment ne doit pas supprimer silencieusement un Artifact ou une Evidence référencée.

## 16. Export

L’export d’un Artifact peut inclure provenance et dérivés ; l’export d’une Attachment peut suivre le package documentaire. Chaque export doit citer la version, appliquer les redactions et conserver la trace de la relation.

## 17. Relation avec Note

Une Note peut référencer ou joindre un fichier. Le fichier reste Attachment tant qu’une action explicite ne l’enregistre pas comme Artifact ou Evidence candidate.

## 18. Relation avec Comment

Un Comment peut porter une Attachment pour illustrer une discussion. La résolution du Comment ne modifie pas automatiquement le lifecycle du fichier ni ses éventuelles promotions.

## 19. Relation avec Report

Un Report peut citer un Artifact, référencer une Evidence ou inclure une Attachment documentaire. Les citations versionnées sont préférables à la duplication lorsque la politique d’accès le permet.

## 20. Relation avec Evidence

Ni Attachment ni Artifact ne devient automatiquement Evidence. La qualification Evidence exige source, Case, contexte, raison, auteur ou moteur, transformations, version et trace de revue.

## 21. Relation avec le workbench

Un workbench consomme un Artifact ou une référence autorisée. Une Attachment doit être promue ou résolue vers un Artifact avant analyse technique, sans préjuger du moteur de 4B.2.

## 22. Options ouvertes

1. **Objets distincts** avec relations explicites.
2. **Attachment comme rôle ou relation d’Artifact**.
3. **Mécanisme Shared Attachment avec promotion vers Artifact**.
4. **Convergence partielle** : contenu commun, métadonnées et lifecycles séparés.

## 23. Avantages par option

| Option | Avantages |
|---|---|
| Objets distincts | ownership et permissions lisibles ; distinction documentaire/analytique forte |
| Rôle/relation d’Artifact | réduit les duplications et conserve un seul contenu adressable |
| Shared avec promotion | collaboration générique, promotion explicite et ownership Investigate préservé |
| Convergence partielle | mutualise stockage/versioning tout en séparant règles métier |

## 24. Risques par option

| Option | Risques |
|---|---|
| Objets distincts | doublons, migrations et références divergentes |
| Rôle/relation d’Artifact | Artifact trop générique et permissions difficiles à isoler |
| Shared avec promotion | double lifecycle et UX de promotion complexe |
| Convergence partielle | modèle abstrait, responsabilité technique et contractuelle ambiguë |

## 25. Migration

Toute migration doit inventorier les Attachments existantes, détecter les références, préserver IDs ou aliases, conserver l’historique, éviter la qualification Evidence automatique et fournir un rollback documentaire. Aucun plan n’est exécuté en 4B.1.

## 26. Preuves nécessaires

- inventaire des usages réels par conteneur ;
- exigences de legal hold, rétention et suppression ;
- volumes, déduplication et versioning ;
- besoins de preview, workbench et export ;
- permissions et séparation des tâches ;
- tests UX de promotion et de citation ;
- compatibilité avec Evidence et Reporting.

## 27. Critères de décision et questions restantes

La décision devra minimiser la duplication sans perdre les distinctions de provenance, permission et lifecycle ; permettre une promotion explicite et réversible dans l’usage ; préserver les liens et les audits ; et rester compréhensible dans Case Workspace, Notes, Evidence Board, Artifact Detail et Reporting. Questions restantes : owner d’Attachment, identité du contenu, politique de suppression, rétention, legal hold, promotion, export et migration.

Aucune option n’est recommandée comme décision finale. CAP-INV-106 reste `proposed`, `delivery_mode: planned` et liée à OPEN-014.
