---
id: context-preservation
domain: 04-experience-architecture
status: draft
owner: UX Architecture Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-008
  - REQ-UX-006
  - REQ-UX-007
  - REQ-UX-009
  - REQ-SEC-001
---

# Conservation du contexte

## Enveloppe de contexte

| Élément | Transmission | Condition |
|---|---|---|
| tenant | automatique | toujours visible et autorisé |
| environnement | automatique | compatible avec destination |
| service | automatique | si pertinent |
| Incident/Case/Finding/Decision/Run | automatique par référence | permission réévaluée |
| Endpoint/Artifact/Evidence | sur action explicite ou lien | sensibilité et périmètre |
| filtre, tri, vue, mode | au retour et deep link sûr | aucun secret |
| sélection, onglet, scroll, panneaux | session locale | objet encore accessible |
| travail non enregistré | confirmation/récupération | jamais transmis à un autre tenant |
| Action Request | explicite vers Govern | autorité et complétude |
| payload brut ou secret | jamais | utiliser référence protégée |

## Transitions

Une transition crée une `return origin` contenant route et état non sensible. La destination affiche le contexte hérité et son produit source. Elle ne transforme pas une référence en propriété locale.

## Cas limites

- tenant changé : effacer tout contexte incompatible après confirmation ;
- environnement indisponible : demander un choix, sans fallback silencieux ;
- objet supprimé : conserver l'origine, afficher tombstone autorisé et action de retour ;
- permission retirée : masquer contenu, préserver une trace minimale non sensible ;
- contexte expiré : expliquer ce qui n'a pas été restauré ;
- multi-tenant agrégé : chaque objet garde son tenant, mutation inter-tenant interdite.

## Exemple end-to-end

Command Incident → Investigate Case → Govern Decision → Response Run → Command. Le retour conserve les références Incident, Case, Finding, Decision, Run et Result sans recopier leurs données.

## Critère d’acceptation

**Given** Incident et Case liés, viewport wide et permissions valides,  
**When** l'utilisateur ouvre le Case depuis Command,  
**Then** tenant/environnement/Incident restent visibles, le Case devient actif, le bouton Retour restaure Incident Detail, et le Context Bar ne duplique pas le breadcrumb.
