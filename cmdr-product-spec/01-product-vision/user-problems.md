---
id: user-problems
domain: 01-product-vision
status: draft
owner: Head of Product
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-001
  - REQ-PROD-008
  - REQ-PROD-009
---
# Problèmes utilisateurs

| Problème | Utilisateurs | Situation | Conséquence | Réponse produit attendue | Propriétaire | Indicateur possible |
|---|---|---|---|---|---|---|
| Fragmentation des outils | SOC, Incident Commander | Les informations et actions sont réparties entre consoles. | Ressaisie, latence, erreurs. | Contexte et objets communs. | Plateforme / Shared | temps et transitions sans contexte |
| Perte de contexte | Tous rôles opérationnels | Une transition ouvre un nouvel outil sans tenant, objet ou sélection. | Reconstruction manuelle et erreur de cible. | Context preservation et liens profonds. | Experience Architecture | taux de transitions restaurées |
| Duplication des objets | Product owners, Engineering | Incident, Case ou Evidence sont redéfinis localement. | États divergents et audit incohérent. | Source unique et projections. | Product Architecture | doublons actifs |
| Investigation dispersée | L2, Hunter, DFIR | Recherche, collecte et analyse ne partagent pas un Case. | Faits non reliés et collaboration difficile. | Case Workspace et workbench relié. | Investigate | temps Case → Finding |
| Décisions peu justifiées | Approvers, Business Owner | La demande ne présente pas preuves, impact et alternatives. | Risque de refus arbitraire ou action dangereuse. | Decision Workspace. | Govern | décisions avec package complet |
| Provenance insuffisante | Analystes, Auditor | Un résumé ou Finding ne montre pas sa construction. | Contestabilité et confiance faibles. | Trace et Evidence links. | Investigate / Shared | conclusions avec provenance |
| Actions mal gouvernées | Operator, Approver | Une action risquée peut être lancée hors policy. | Dommage et non-conformité. | Action classes et Govern. | Govern | actions élevées avec Decision |
| Handovers incomplets | Incident Commander, SOC | Le changement d'équipe perd ownership et prochaine action. | Travail dupliqué et délai. | Handover structuré. | Command | handovers complets |
| Vue du travail incohérente | SOC, Managers | Plusieurs files montrent des états différents. | Priorité et charge illisibles. | Work Queue unique avec vues. | Command | objets sans owner / SLA |
| Sécurité et impact métier séparés | Incident Commander, Business Owner | Les faits techniques ne sont pas reliés aux services. | Mauvaises priorités. | Service impact dans Command et Decision. | Command / Govern | incidents avec impact qualifié |
| Automatisations opaques | Automation Designer, Auditor | Les étapes, outils et reprises ne sont pas visibles. | Échec non explicable. | Automation Run, Tool Calls, Control Room. | Studio | runs entièrement traçables |
| IA non contrôlée | Tous rôles | Une suggestion peut être confondue avec un fait ou une décision. | Biais, contournement et perte de confiance. | Provenance, Human Gates et optionnalité. | Studio / Govern | outputs attribués |
| Collecte et preuve déconnectées | DFIR, Auditor | Une collecte n'est pas reliée à intégrité et Case. | Preuve inutilisable. | Collection → Artifact → Evidence. | Investigate | collectes avec chain of custody |
| Signal difficile à transformer en résultat | SOC, Command | Chaque étape utilise une logique et un owner différents. | Temps et pertes de handoff. | Chaîne nominale et transitions. | Tous produits | temps signal → Result |
| Apprentissage peu intégré | Detection Engineer, Hunter | Les Findings et Results n'alimentent pas les règles. | Répétition et faible couverture. | Boucle d'amélioration Detection Engineering. | Investigate / Command | Findings reliés à règles |

## Critère d'acceptation

**Given** un problème décrit,  
**When** une feature est proposée,  
**Then** elle identifie le problème, l'utilisateur, le propriétaire produit et un résultat observable ; une feature sans problème explicite n'entre pas automatiquement dans le scope.
