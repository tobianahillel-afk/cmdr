---
id: experience-time-freshness
domain: 04-experience-architecture
status: draft
owner: UX Architecture Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-004
  - REQ-UX-005
---
# Temps et fraîcheur


Event time, ingestion time, processing time et last refresh sont nommés lorsqu'ils divergent. La timezone utilisateur est visible et l'UTC reste accessible.

`live`, `delayed`, `stale`, `partial`, `offline` et `unknown` sont textuels. Le seuil de stale appartient à la capability et n'est pas inventé par le composant. Les mises à jour ne déplacent pas la sélection.

**Given** une source retardée, **When** une Decision est ouverte, **Then** fraîcheur, dernière donnée et conséquence sur la décision sont visibles sans présenter le contenu comme live.
