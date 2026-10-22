---
id: experience-command-palette
domain: 04-experience-architecture
status: draft
owner: UX Architecture Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-001
  - REQ-UX-007
  - REQ-AI-001
  - REQ-SEC-002
---

# Command Palette — modèle d’expérience

La Command Palette permet navigation, commandes autorisées, raccourcis et actions contextuelles. Elle n'est pas un chatbot.

Les résultats sont groupés : destinations, objets récents, commandes locales, commandes globales, automatisations optionnelles. Chaque action indique produit, portée, raccourci, permission et, si applicable, nécessité de Govern.

Une action classe 2+ n'est jamais exécutée directement : la palette ouvre la confirmation, l'Action Request ou le workspace propriétaire. Les suggestions IA sont libellées et attribuées ; la recherche de commandes fonctionne sans modèle.

**Given** une commande d'isolation classe 3,  
**When** elle est trouvée dans la palette,  
**Then** l'utilisateur voit cible, classe et permission, puis est dirigé vers Govern ; aucune exécution immédiate n'a lieu.
