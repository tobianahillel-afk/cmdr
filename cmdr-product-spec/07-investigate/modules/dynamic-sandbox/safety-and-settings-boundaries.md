---
id: dynamic-sandbox-safety-settings-boundaries
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-SEC-001
  - REQ-SEC-002
open_decisions:
  - OPEN-013
---
# Safety and Platform Settings boundaries

Platform Settings crée, valide, met hors service et supervise les Sandbox Environments, providers, secrets, capacité, health, rétention et politiques. Investigate sélectionne seulement une projection autorisée, vérifie les préconditions et peut demander une alternative.

La sécurité fonctionnelle exige limites de réseau, durée, ressources et interactions ; interdiction de propagation et d’accès non autorisé ; stop d’urgence ; nettoyage ; reset ; blocage d’un environnement suspect ; signalement à Settings ; conservation de la trace. Aucun mécanisme d’isolation ou hyperviseur n’est choisi et aucune sandbox n’est déclarée invulnérable.
