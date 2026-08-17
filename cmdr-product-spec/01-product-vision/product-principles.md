---
id: product-principles
domain: 01-product-vision
status: draft
owner: Head of Product
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-002
  - REQ-PROD-003
  - REQ-PROD-004
  - REQ-PROD-005
  - REQ-PROD-006
  - REQ-PROD-007
  - REQ-PROD-008
  - REQ-PROD-009
  - REQ-PROD-010
  - REQ-PROD-011
  - REQ-PROD-012
  - REQ-AI-007
  - REQ-AI-010
---
# Principes produit

Ces principes arbitrent les décisions produit, UX et fonctionnelles. Ils ne sont pas des slogans : chaque principe décrit une conséquence et un anti-pattern.

| Requirement | Principe | Définition | Raison | Conséquence UX | Conséquence fonctionnelle | Anti-pattern | Exemple |
|---|---|---|---|---|---|---|---|
| REQ-PROD-002 | Evidence first | Une affirmation importante est soutenue par des sources et une provenance. | Éviter les décisions fondées sur une vue isolée. | Montrer preuve, fraîcheur et intégrité avant la recommandation. | Les workflows lient Artifact, Evidence et Finding. | Afficher un score sans sources. | Un Finding cite les Evidence favorables et contradictoires. |
| REQ-PROD-003 | Human accountable | Une personne ou une autorité identifiable demeure responsable. | L'automatisation ne doit pas dissoudre la responsabilité. | Afficher initiateur, owner, approbateur et dernier acteur. | Les décisions et actions enregistrent l'autorité. | « L'IA a décidé » sans responsable. | Un approbateur humain assume une Decision de classe élevée. |
| REQ-PROD-004 | Actions safely governed | Le risque détermine contrôles, approbations et rollback. | Réduire les dommages et actions hors périmètre. | Présenter classe, cible, impact et retour arrière. | Créer une Action Request et appliquer les policies. | Bouton destructif direct depuis une analyse. | Une isolation de classe 3 passe par Govern. |
| REQ-PROD-005 | Every conclusion traceable | Une conclusion expose données, méthode, auteur et incertitude. | Permettre revue, contestation et audit. | Ouvrir la trace depuis Finding, résumé ou rapport. | Conserver versions, Tool Calls et sources. | Résumé sans lien vers les éléments sources. | Une proposition IA montre run, outils et Evidence. |
| REQ-PROD-006 | No duplicated object | Un objet possède une seule définition et un propriétaire. | Éviter divergence et cycles concurrents. | Afficher des projections reliées à la source. | Les changements passent par le propriétaire. | Redéfinir Evidence dans Command. | Command affiche une projection Investigate. |
| REQ-PROD-007 | Progressive disclosure | La densité augmente avec l'intention et l'expertise. | Rendre CMDR utilisable sans cacher l'essentiel. | Priorité d'abord, détail et technique à la demande. | Les workspaces révèlent outils et traces progressivement. | Tout afficher dans un HUD dense. | Mission Control résume ; l'Inspector révèle la trace. |
| REQ-PROD-008 | Context preserved | Tenant, objets, filtres et retour survivent aux transitions. | Éviter ressaisie et erreurs de périmètre. | Context Bar, liens profonds et retour réel. | Propager identifiants et sélection autorisés. | Ouvrir Govern sans Incident ni Case. | Le lien vers Govern conserve l'Action Request. |
| REQ-PROD-009 | One workflow, one owner | Chaque workflow métier possède un responsable. | Éviter les files et lifecycles parallèles. | Montrer owner et prochaine action. | Le produit propriétaire définit états et handoffs. | Deux Work Queues pour le même travail. | Command possède le lifecycle de l'Incident. |
| REQ-PROD-009 | One capability, one canonical owner | Une capability partagée possède un owner unique. | Maintenir contrat et roadmap cohérents. | Référencer l'owner dans l'UI d'administration. | Registre et source unique. | Chaque produit crée son export engine. | Reporting Engine reste Shared. |
| REQ-PROD-006 | One definition, one source of truth | Une règle partagée n'est définie qu'une fois. | Réduire les contradictions. | Liens contextuels vers la source. | Registres et matrice assurent la résolution. | Copier une palette dans un écran. | L'écran cite les tokens de marque. |
| REQ-PROD-010 | Essential workflows work without AI | L'usage normal ne dépend pas d'un modèle. | Résilience, contrôle et accessibilité économique. | Présenter les fonctions IA comme options. | GUI, règles, moteurs déterministes et API restent disponibles. | Bloquer Case Workspace sans modèle. | La recherche déterministe fonctionne hors IA. |
| REQ-PROD-011 | Deterministic engines are first-class | Règles et moteurs reproductibles ont le même statut produit. | Les agents ne remplacent pas la détection vérifiable. | Distinguer règle, modèle et suggestion. | Tests, versioning et déploiement propres. | Masquer les règles derrière un assistant. | Detection Engineering conserve tests et rollback. |
| REQ-PROD-011 | Rules and human workflows are not replaced by agents | Un agent assiste un workflow propriétaire. | Préserver expertise et contrôle. | L'agent apparaît comme participant attribué. | Les produits restent propriétaires des objets. | Déplacer le triage entier dans Studio. | Un agent propose, l'analyste dispose. |
| REQ-SEC-002 | Dangerous actions require explicit governance | Les classes élevées exigent autorité applicable. | Éviter les contournements. | Afficher le passage dans Govern. | Policies, Approval et Decision précèdent le run. | Exécuter une suppression depuis un Tool Call. | Une action classe 4 exige Decision. |
| REQ-PROD-003 | Every important action attributable | Toute mutation significative possède initiateur et contexte. | Responsabilité et investigation d'incident. | Activity et Trace exposent les acteurs. | Audit immuable et correlation ID. | Action système anonyme. | Une commande indique humain, workflow et agent. |
| REQ-AI-010 | Generated conclusions expose provenance | Les contenus générés indiquent données, modèle, run et incertitude. | Ne pas confondre synthèse et fait. | Badge de provenance et accès à la trace. | Stocker inputs référencés et version. | Texte généré indiscernable d'un Finding confirmé. | Une hypothèse IA reste proposée. |
| REQ-AI-007 | Automation interruptible and auditable | Une Automation Run peut être observée, suspendue et expliquée. | Limiter les échecs en cascade. | Control Room expose stop, pause et trace. | État, Tool Calls et outputs sont persistés. | Automation opaque sans contrôle. | Un humain interrompt avant une action risquée. |
| REQ-PROD-019 | Integration is not automatically native | Native implique expérience, contrat et moteur principal CMDR. | Éviter les promesses trompeuses. | Masquer la marque dans l'activité mais montrer la provenance du moteur. | Classifier chaque capability. | Appeler natif un simple wrapper. | Un moteur forensic temporaire est étiqueté. |
| REQ-PROD-012 | Planned is not implemented | La cible future est séparée de la livraison prouvée. | Décisions et attentes réalistes. | Badges target/planned/current. | Preuve de release requise pour implemented. | Présenter Endpoint Agent complet comme livré. | Le README indique cible native planifiée. |
| REQ-PROD-003 | No hidden autonomous decision | Aucune décision métier ne résulte d'une autonomie invisible. | Maintenir légitimité et recours. | Différencier proposition, approval et Decision. | Govern enregistre l'autorité. | Agent modifie l'état critique sans trace. | La Decision cite approbateur et policy. |

## Arbitrage

En cas de tension :

1. intégrité, sécurité et responsabilité précèdent la vitesse ;
2. la source canonique précède la convenance locale ;
3. l'usage essentiel sans IA précède l'optimisation agentique ;
4. une décision ouverte reste visible plutôt qu'inventée ;
5. un produit propriétaire conserve son workflow, même si Studio l'assiste.

## Critère d'acceptation

**Given** une proposition d'automatiser entièrement une décision de containment,  
**When** elle est évaluée contre les principes,  
**Then** le workflow conserve un propriétaire humain, expose les Evidence et la classe d'action, passe par Govern, permet interruption et audit, et reste exécutable sans agent lorsque l'automatisation est indisponible.
