---
id: CAP-INV-704
title: Mobile Acquisition, Backup and Extraction Context Review
product: investigate
module: mobile-forensics
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-07
requirement_ids: [REQ-INV-001, REQ-PROD-014, REQ-PROD-020, REQ-OBJ-004, REQ-SEC-001]
open_decisions: [OPEN-005, OPEN-008, OPEN-011, OPEN-013, OPEN-014]
source-of-truth: canonical
---
# CAP-INV-704 — Mobile Acquisition, Backup and Extraction Context Review

## 1. Définition
Revoir fonctionnellement l’origine et le contexte d’une représentation Mobile sans prescrire sa méthode technique : live acquisition request, logical extraction, filesystem extraction, local backup, synchronized backup, selected-file collection, application-data export ou autre evidence representation.

## 2. Problème utilisateur
Une représentation peut être partielle, transformée, copiée ou issue d’un backup sans représenter l’état courant du device. Sans contexte explicite, l’analyste peut confondre acquisition, extraction, appareil original et couverture réelle.

## 3. Objectifs
- distinguer toutes les représentations disponibles;
- exposer initiateur, autorité, méthode déclarée, période, device state, lock/encryption state déclarés, disponibilité, erreurs et éléments manquants;
- conserver transformations, transferts, copies, custody, restrictions et limitations;
- préparer l’Integrity/Completeness/Accessibility Assessment.

## 4. Non-objectifs
Aucune procédure de déverrouillage/extraction, acquisition réelle, bypass, rooting/jailbreak, password/code test, exploit, outil imposé, API, protocole, format propriétaire, commande, driver, profile installation ou device mutation.

## 5. Propriétaire
Investigate possède l’Acquisition Context Review. Collection possède requests/jobs/exécution/résultats/custody de collecte. Govern possède l’autorité risquée. Settings/Endpoint conservent configuration et capacités déclarées.

## 6. Utilisateurs
Principal : Evidence Reviewer. Secondaires : Mobile Forensics Analyst, DFIR Analyst, Investigation Lead et Sensitive Data Reviewer.

## 7. Conditions d’entrée
Session/Intake, Mobile Evidence Package ou source representation, acquisition/custody references disponibles ou explicitement manquantes, permission de lecture et purpose/scope.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Collection Request/Job/result | Collection | acquisition origin/status | si acquisition CMDR | version terminale/courante | origin gap explicite |
| Mobile Evidence Package | source/Investigate concept | representation envelope | oui | version sélectionnée | review impossible |
| Declared method/representation | collector/source | acquisition descriptor | oui si connu | acquisition version | `incomplete` |
| Device/lock/encryption/power availability state | source | declared context | non | acquisition time | unknown, jamais inféré |
| Transformations/transfers/copies | custody/source | lineage | selon history | event timestamps | custody gap |
| Authority/restrictions | Govern/Collection/Security | legal/operational boundary projection | oui si applicable | acquisition version | `disputed`/`restricted` |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Collection Request/Job | Collection/Investigate | request, target, status, errors | read |
| Mobile Evidence Package | source/Investigate concept | package/version/manifest context | read |
| Device Backup / Logical / Filesystem Extraction concepts | source | declared representation | read |
| Custody/provenance events | Collection/Shared/Trust | transfers, copies, transformations | read/review |
| Action Request/Decision | Govern | acquisition authority projection | read only |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Mobile Acquisition Context Assessment | create/review/dispute/supersede | Investigate concept | representation and limitations explicit |
| Representation relation | create/dispute | Investigate | Backup ≠ Filesystem ≠ Logical; extraction ≠ device |
| Custody gap annotation | create/supersede | Investigate/Trust projection | source record not rewritten |
| Complementary Collection proposal | prepare | Collection | no acquisition execution |

## 11. Fonctionnalités
Compare representations and versions; display initiator, authority, declared method, acquisition period, device state, lock/encryption, availability, errors, missing areas, transformations, transfers, copies, custody, restrictions and limitations; link source events and identify unsupported/inconsistent claims.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter acquisition/custody | reviewer | source projections | 0 | read | context visible | non |
| Annoter/disputer contexte | reviewer | assessment | 2 | reason/source | new version | OPEN-013 |
| Comparer deux représentations | analyste | representations | 1 | read both | bounded comparison | non |
| Demander complément | lead | collection proposal | 2 | gap + purpose | handoff Collection | future authority rechecked |
| Ouvrir intégrité | reviewer | CAP-INV-705 | 0 | context available | next review | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Classer representation déclarée | oui | metadata rules | oui | suggestion | source fields/table |
| Assembler custody events | oui | IDs/timestamps | oui | non nécessaire | deterministic timeline |
| Signaler gaps/incohérences | oui | rules | oui | explanation | checklist/diff |
| Résumer limitations | oui | structured fields | oui | summary | table |
| Proposer acquisition technique | non | non | non | interdit | Collection future decision |

## 14. États fonctionnels
`declared`, `under-review`, `accepted-with-limitations`, `incomplete`, `inconsistent`, `disputed`, `unsupported`, `superseded`. These do not state whether underlying data is malicious or relevant.

## 15. États d’interface
Loading preserves package/version; Empty identifies missing acquisition record; Partial lists missing context; Error keeps valid custody; Offline permits read-only cached context; Permission denied hides restricted details; Stale identifies superseded versions.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Acquisition Context Assessment | concept | CAP-INV-705..719 | method declared, authority, transformations, gaps and limitations |
| Representation comparison | result | analyst/QA | representations remain distinct |
| Custody gap annotation | event | CAP-INV-705/719 | source and reason retained |
| Collection gap package | package | Collection | bounded need, no execution |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-701/702 | acquisition review | CAP-INV-704 | Session, package, source, scope | Session/Intake |
| CAP-INV-704 | trust review | CAP-INV-705 | representation, custody, transformations, gaps | CAP-INV-704 |
| CAP-INV-704 | data analysis | CAP-INV-706..715 | accepted representation + limitations | Session |
| Gap | complementary source needed | CAP-INV-202 | need, source, device, scope, authority | CAP-INV-704 |
| Closure | provenance | CAP-INV-719 | assessment versions, sources, decisions | Session |

## 18. Dépendances
CAP-INV-202/203/213/214, CAP-INV-701..703/705..719, Evidence Trust, Shared Trace/Timeline, Govern authority, OPEN-005/008/011/013/014.

## 19. Source de vérité
Collection remains source for acquisition execution/result/custody events; source package owner remains source for representation; Mobile assessment is Investigate; authority is Govern. Review never rewrites collection history.

## 20. Provenance et audit
Record Session, package/version, Collection refs, initiator, authority, declared representation/method, period, device/lock/encryption state, transformations, transfers, copies, custody events, errors, gaps, reviewer, disputes and supersession.

## 21. Permissions fonctionnelles
Acquisition context read/review, custody read, raw manifest read, representation comparison, dispute, restricted source read, complementary collection prepare and provenance export. No acquisition execution permission is granted here.

## 22. Limites et erreurs
Missing manifest, unknown collector/method, conflicting custody, encrypted/locked source, corrupted package, unsupported representation, permission denial or cross-tenant restriction keeps assessment incomplete/disputed; no technical bypass is suggested.

## 23. Métriques
Assessments by state/representation, custody gaps, unknown methods, disputed contexts, partial packages, complementary collection proposals and time to accepted-with-limitations.

## 24. Classification de livraison
`defined` / `planned`. No acquisition method, tool, format, platform, engine, protocol or implementation selected.

## 25. Critères d’acceptation
**Given** a synchronized backup **When** acquisition context is reviewed **Then** it remains a synchronized backup and is not presented as current local device state.

**Given** a locked device with a partial extraction **When** the review completes **Then** lock and partiality remain explicit and no unlock procedure is produced.

**Given** no AI **When** context is reviewed **Then** source metadata, custody timeline, rules and checklist provide full functionality.

## 26. Questions ouvertes
OPEN-011 keeps platform/acquisition/tool strategy unresolved; OPEN-005 forensic engines; OPEN-008 collection/Endpoint support; OPEN-013 mutations; OPEN-014 package/material relations.

## 27. Consommateurs documentaires
CAP-INV-705..719, Collection, Evidence, Trust, Objects, Permissions, Quality, Roadmap and Mobile source migration.
