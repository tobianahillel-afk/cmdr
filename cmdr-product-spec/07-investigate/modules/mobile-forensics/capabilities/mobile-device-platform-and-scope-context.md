---
id: CAP-INV-703
title: Mobile Device, Platform and Scope Context
product: investigate
module: mobile-forensics
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-07
requirement_ids: [REQ-INV-001, REQ-PROD-014, REQ-PROD-020, REQ-PROD-055, REQ-SEC-001]
open_decisions: [OPEN-008, OPEN-011, OPEN-013, OPEN-014]
source-of-truth: canonical
---
# CAP-INV-703 — Mobile Device, Platform and Scope Context

## 1. Définition
Construire un contexte provider-neutral de device, plateforme et scope à partir de sources autorisées, en séparant device candidates, classe, plateforme/version déclarée ou candidate, identifiers/aliases, manufacturer/model metadata, ownership déclaré, management/enrollment state, temps/région et inclusions/exclusions.

## 2. Problème utilisateur
Les identifiants mobiles, comptes et metadata peuvent être ambigus. Sans contexte explicite, un analyste peut attribuer un appareil à une personne, traiter une plateforme candidate comme certaine, ou mélanger données hors scope/cross-tenant.

## 3. Objectifs
- conserver plusieurs device/platform/version candidates avec source, confidence et contradiction;
- enregistrer device class, identifiers autorisés, aliases et manufacturer/model metadata;
- distinguer ownership déclaré, management/enrollment state et utilisateur effectif inconnu;
- définir included/excluded scope, time/regional context et restrictions.

## 4. Non-objectifs
Aucune résolution automatique de personne, aucune sélection de plateforme supportée, aucun enrollment, MDM/EMM change, agent assumption, device fingerprinting actif, commande, scan, connector ou schéma final.

## 5. Propriétaire
Investigate possède Mobile Device/Platform candidate context and scope semantics. Settings owns configured Fleet/enrollment/policies/source projections; Endpoint owns agent-reported capability when present; Shared owns Entity/Graph mechanisms.

## 6. Utilisateurs
Principal : Mobile Forensics Analyst. Secondaires : Investigation Lead, DFIR Analyst, Evidence Reviewer, Sensitive Data Reviewer et Platform Administrator en consultation.

## 7. Conditions d’entrée
Mobile Session/Intake, at least one device/package/source reference, permission to view permitted identifiers, tenant/environment and scope purpose.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Device identifiers/aliases | package/source/Settings | technical identity candidates | oui, au moins une ref | source version | unknown candidate, jamais personne inférée |
| Declared/candidate platform/version | source/metadata/analyste | platform context | non | source timestamp | unknown, no default platform |
| Manufacturer/model/device class | source | descriptive metadata | non | source version | unknown |
| Declared owner/management/enrollment | source/Settings | administrative context | non | last authoritative projection | unknown/stale explicit |
| Time/region context | source/Session | interpretation context | non | source scoped | unknown/time quality gap |
| Included/excluded scope/restrictions | Intake/Session/Security | analytical boundary | oui | current Session version | analysis blocked or incomplete |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Mobile Session / Intake | Investigate | purpose/scope/source refs | lire/lier |
| Mobile Evidence Package | source/Investigate concept | device/platform metadata | lire selon permission |
| Endpoint/Fleet/Policy | Endpoint/Settings | management/capability projection | lire seulement |
| Tenant / Environment / Principal | Settings | boundary/admin context | lire seulement |
| Entity / Graph | Shared | candidate relation mechanism | lire/lier sans merge automatique |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Mobile Device Candidate | créer/review/dispute/supersede | Investigate concept | source + scope + confidence; identifier ≠ person |
| Mobile Platform Candidate | créer/review/dispute/supersede | Investigate concept | declared vs inferred explicit |
| Scope inclusion/exclusion set | modifier/versionner | Investigate | permission/source restrictions preserved |
| Candidate relations | créer/disputer | Investigate using Shared | no identity merge or permission transfer |

## 11. Fonctionnalités
Compare device candidates, normalized display aliases, platform/version claims, device class, manufacturer/model, declared ownership, management/enrollment, source/package relations, time/region context, scope tree, exclusions, confidence, contradictions and supporting sources.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter metadata | analyste | projection | 0 | read | context visible | non |
| Ajouter/disputer candidate | analyste | candidate | 2 | source | versioned assessment | non |
| Sélectionner candidate pour Session | analyste | Session relation | 2 | scope permission | selected candidate, not truth | OPEN-013 |
| Modifier included/excluded scope | owner | scope | 2 | authority | versioned scope | OPEN-013 |
| Ouvrir Fleet/source | utilisateur autorisé | Settings projection | 0 | destination permission | navigation | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Normaliser labels/aliases | oui | oui | oui | non nécessaire | parser/table |
| Comparer metadata | oui | oui | oui | résumé | diff/table |
| Proposer plateforme candidate | oui | metadata rules | oui | suggestion | source matrix |
| Signaler contradictions | oui | oui | oui | explanation | comparison |
| Associer à une personne | jamais automatique | non | non | interdit | human sourced relation only |

## 14. États fonctionnels
Candidates: `proposed`, `under-review`, `supported`, `weakly-supported`, `contradicted`, `inconclusive`, `disputed`, `superseded`, `withdrawn`. Scope may be `draft`, `ready`, `restricted`, `out-of-scope`, `superseded`.

## 15. États d’interface
Loading preserves Session; Empty means no identity metadata; Partial lists missing/stale attributes; Error preserves sourced candidates; Offline uses last projection; Permission denied masks identifiers; Stale marks administrative metadata age.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Device Context | candidate assessment | CAP-INV-704..717 | source/confidence/contradictions visible |
| Platform Candidate | candidate relation | CAP-INV-704/706..708 | no platform implementation implied |
| Scoped source set | versioned scope | all Mobile analysis | includes/excludes and restrictions explicit |
| Context activity | event | Trace/QA | author/source/version/disposition |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-702 | orient analysis | CAP-INV-703 | Session, package refs, scope | Session |
| CAP-INV-703 | review acquisition | CAP-INV-704/705 | device/platform candidates, scope | Session |
| CAP-INV-703 | storage/app analysis | CAP-INV-706/707 | selected candidate, package refs, restrictions | Session |
| CAP-INV-703 | location/connectivity | CAP-INV-711/713 | device context, time/region, scope | Session |
| Contradiction | review | CAP-INV-717 | candidate/support/contradictions | CAP-INV-703 |

## 18. Dépendances
CAP-INV-701/702/704..719, Settings Tenant/Fleet/Policies/Sources, Endpoint declared capabilities, Shared Entity/Graph/Linking, Security permission model, OPEN-008/011/013/014.

## 19. Source de vérité
Candidate assessments and analytical scope: Investigate. Configured device/Fleet/enrollment/policy: Settings. Agent state: Endpoint. Generic Entity/Graph: Shared. No local candidate becomes a canonical person/identity automatically.

## 20. Provenance et audit
Session, source/package, raw permitted identifier references, aliases, declared/inferred flags, platform/version source, device class, manufacturer/model, declared owner, management/enrollment, time/region, scope changes, confidence, contradictions, actor and versions.

## 21. Permissions fonctionnelles
Device identifier read, restricted identifier read, platform candidate review/select, Fleet/source projection read, scope update, cross-tenant relation and candidate export. Raw identifiers remain masked where required.

## 22. Limites et erreurs
Identifier collision/rotation, stale enrollment, unsupported metadata, missing platform version, conflicting aliases/owners, redacted identifiers or cross-tenant relations produce candidate uncertainty; no resolution is silently forced.

## 23. Métriques
Candidates per Session, unresolved contradictions, unknown platform/version, scope revisions, stale admin projections, permission masks and manual changes to automated candidate suggestions.

## 24. Classification de livraison
`defined` / `planned`. No mobile platform list, MDM provider, identifier schema, entity resolution algorithm, connector or code selected.

## 25. Critères d’acceptation
**Given** two identifiers pointing to conflicting owner metadata **When** the analyst reviews context **Then** both sources and contradiction remain visible and no person is auto-selected.

**Given** no platform metadata **When** analysis starts **Then** the platform remains unknown/candidate and no iOS/Android default is assumed.

**Given** no AI **When** device context is reviewed **Then** source tables, filters, diff and human disposition provide complete functionality.

## 26. Questions ouvertes
OPEN-011 defines future platform priority/version/device support; OPEN-008 source/Endpoint support; OPEN-013 class-2 scope mutations; OPEN-014 final identity/material relations remain open.

## 27. Consommateurs documentaires
Mobile acquisition/integrity/filesystem/application/location/connectivity capabilities, Case/Entity views, Settings/Endpoint, Shared Graph, Security, Objects, Screens, Quality and Roadmap.
