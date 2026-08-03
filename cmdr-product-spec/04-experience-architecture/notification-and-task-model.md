---
id: experience-notification-task
domain: 04-experience-architecture
status: draft
owner: UX Architecture Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-009
  - REQ-UX-007
---
# Notifications et tâches


Une notification informe ; elle ne change pas l'état métier. Une Task possède owner, échéance, statut et objet source. Dédoublonnage, urgence et canal sont contrôlés par la capacité partagée.

Toute notification actionnable fournit un deep link et une raison. Une Task fermée ne ferme jamais automatiquement l'Incident, Case ou Decision source.

**Given** trois événements identiques, **When** ils sont notifiés, **Then** une notification agrégée expose count, période et lien ; aucune Task supplémentaire n'est créée sans règle propriétaire.
