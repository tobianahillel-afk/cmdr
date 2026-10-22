---
id: terminology-rules
domain: 00-governance
status: draft
owner: Documentation Governance Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-006
  - REQ-UX-001
  - REQ-AI-002
  - REQ-OBJ-009
---
# Règles de terminologie structurelle

## Objet

Ces définitions sont normatives pour l'architecture documentaire. Le glossaire éditorial détaillé appartient à `15-content-and-language/`.

| Terme | Définition normative | Ne signifie pas |
|---|---|---|
| Product | Ensemble cohérent de responsabilités, workspaces et objets sous un propriétaire produit | Une marque tierce, un simple dossier ou un moteur |
| Domain | Zone documentaire correspondant à une responsabilité canonique | Un nouveau produit implicite |
| Module | Sous-ensemble fonctionnel d'un produit, orienté activité | Une capacité partagée redéfinie |
| Capability | Résultat réutilisable possédé par un produit ou Shared Capabilities | Une liste d'écrans |
| Feature | Comportement livrable au sein d'une capability | Un produit autonome |
| Workspace | Espace durable pour une activité longue, complexe ou multi-objet | Une vue filtrée |
| Page | Destination routable ayant un objectif utilisateur distinct | Chaque filtre ou représentation |
| View | Sous-ensemble ou organisation enregistrée des mêmes objets et du même objectif | Une nouvelle source de vérité |
| Mode | Représentation ou configuration alternative du même espace | Un nouvel objet ou workflow |
| Filter | Restriction temporaire et visible du jeu de données | Une page |
| Panel | Région persistante ou redimensionnable d'un workspace | Un produit |
| Inspector | Panneau canonique de contexte sur l'objet sélectionné | Un panneau droit recréé par module |
| Drawer | Surface temporaire pour une action courte ou un formulaire secondaire | Un workspace d'analyse |
| Modal | Confirmation ou décision courte bloquant temporairement l'arrière-plan | Une activité longue |
| Object | Entité canonique avec identité, propriétaire et cycle de vie | Une simple projection UI |
| Projection | Représentation locale d'un objet canonique, sans changement de propriété | Une nouvelle définition |
| Source of truth | Fichier actif de plus haute autorité pour un concept | Toute mention du concept |
| Owner | Domaine responsable de la sémantique, des changements et de la validation | Tout producteur ou consommateur |
| Consumer | Produit ou document utilisant une source canonique | Copropriétaire implicite |
| Native | CMDR possède l'expérience, le contrat fonctionnel et le moteur principal | Une interface CMDR autour d'un moteur externe |
| Integrated | Capacité durablement fournie par un système externe et adaptée au modèle CMDR | Une capability native |
| Temporary integration | Moteur externe temporaire avec remplacement ou décision attendue | Une dépendance permanente cachée |
| Planned | Capacité cible non encore suffisamment spécifiée ou livrée | Une fonctionnalité disponible |
| Out of scope | Capacité explicitement exclue du périmètre considéré | Une lacune non documentée |
| AI | Modèle ou capacité probabiliste optionnelle augmentant un workflow | L'interface ou l'autorité principale |
| Deterministic engine | Moteur dont les entrées, règles et résultats sont reproductibles selon un contrat | Un agent |
| Automation | Exécution gouvernée d'étapes prédéfinies, avec ou sans IA | Autonomie illimitée |
| Automation Agent | Entité Studio configurée avec rôle, objectifs, Skills, Tools, limites et permissions | Propriétaire d'une décision métier |
| Workflow | Orchestration explicite d'étapes, conditions, outils, agents et Human Gates | Une intention libre |
| Tool | Interface exécutable définissant entrées, sorties, permissions, effets et risque | Une Skill ou une capability |
| Skill | Connaissance, méthode ou procédure versionnée pouvant utiliser des Tools | Une exécution |
| Human Gate | Suspension d'une Automation Run exigeant une validation humaine | Une Decision Govern par défaut |
| Automation Run | Exécution Studio d'un Workflow ou d'un Automation Agent | Un Response Run |
| Response Run | Exécution gouvernée d'une action de réponse autorisée par Govern | Toute automatisation |
| Result | Résultat vérifié d'un Response Run, consommable par Command et Investigate | Un Finding analytique |

## Règles d'usage

- Employer les noms canoniques `Command`, `Investigate`, `Govern`, `CMDR Studio`, `Platform Settings`, `Endpoint Agent`.
- Écrire les termes d'objet avec une majuscule lorsqu'ils désignent l'entité canonique.
- Ne pas utiliser `agent` seul lorsqu'il peut désigner un Endpoint Agent ou un Automation Agent.
- Ne pas employer `approval` comme synonyme de `Decision`.
- Ne pas employer `native` sans classification et preuve.
- Ne pas employer `screen`, `page`, `workspace` et `view` de manière interchangeable.

## Critère d'acceptation

**Given** un document utilisant « workspace », « native », « agent » ou « run »,  
**When** le terme est interprété,  
**Then** son sens correspond à une seule définition de ce document, son propriétaire est identifiable et aucune ambiguïté n'existe entre Automation Run et Response Run, Endpoint Agent et Automation Agent, ou page et vue.
