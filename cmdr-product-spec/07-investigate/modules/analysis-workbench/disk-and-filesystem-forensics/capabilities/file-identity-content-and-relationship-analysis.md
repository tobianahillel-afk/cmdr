---
id: CAP-INV-368
title: File Identity, Content and Relationship Analysis
product: investigate
module: analysis-workbench
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-05
requirement_ids:
  - REQ-INV-001
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-OBJ-003
  - REQ-AI-002
  - REQ-SEC-001
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-368 — File Identity, Content and Relationship Analysis

## 1. Définition
Inspecter en lecture seule l’identité fonctionnelle, le contenu prévisualisable et les relations d’un fichier sans exécution, sans faire du path ou d’un indicateur d’identité une certitude historique.

## 2. Problème utilisateur
Copies, versions, conteneurs, liens et metadata peuvent suggérer une identité commune sans prouver qu’il s’agit du même objet historique; la prévisualisation peut aussi exposer du contenu sensible.

## 3. Objectifs
Voir name/path/context, declared/detected type, size, identity indicators, versions/copies candidates, links, parents/containers et related files; preview sans exécuter; comparer; annoter; relier à Entity/Artifact; router Static/Reverse; préserver source.

## 4. Non-objectifs
Ne pas exécuter, confirmer identité absolue, définir hashing/dedup/format final, extraire secrets, modifier source, créer Evidence automatiquement ou remplacer Static/Reverse.

## 5. Propriétaire
Investigate possède File Observation et relations; Artifact lifecycle reste CAP-INV-105; Studio Tools; Shared Preview/Linking/Export; Security contrôle contenu sensible.

## 6. Utilisateurs
Principal : Artifact Analyst. Secondaires : Disk Analyst, Malware Analyst, Reverse Engineer, Privacy Reviewer et Evidence Reviewer.

## 7. Conditions d’entrée
Entry sélectionnée; source/version et metadata visibles; permission content/preview; contenu jamais exécuté; limitations héritées.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Filesystem Entry | CAP-INV-367 | path/metadata/state | oui | source version | blocked |
| Content/reference | Disk Image | bytes/reference autorisés | selon analyse | immutable source | metadata-only |
| Type/identity results | Tool Call | declared/detected type/indicators | non | Tool/version visible | unknown |
| Related entries/containers | filesystem/Artifact graph | relations candidates | non | même scope | unlinked |
| Content permissions | Security | preview/sensitive access | oui | current | masked/denied |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Disk Image/Entry | Investigate | source/metadata/content ref | lire |
| Artifact/Derived Artifact | Investigate | versions/lineage | lire/lier |
| Tool/Tool Call | Studio | producer/type/identity result | lire |
| Entity/Case/Hypothesis | Shared/Investigate | context/relations | lire/lier |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| File Observation | créer/annoter/dispute | Investigate concept | path ≠ identity |
| Candidate relationship | créer/supersede | Investigate | source/evidence/confidence obligatoires |
| Derived Artifact request | préparer/créer via CAP-INV-311 | Investigate | parent/source/Tool requis |
| Source file | aucune mutation/exécution | Disk Image | invariant |

## 11. Fonctionnalités
Inspect name/path/context/type/size/identity indicators/copies/versions/links/parents/containers/related files; permission-aware preview; compare content/metadata; expose contradictions; annotate/link; route to Static/Reverse; preserve lineage.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter metadata/relations | Analyst | File Observation | 0 | read | view sourcée | non |
| Prévisualiser contenu | Analyst | content | 0 | preview permission | no execution | non |
| Comparer | Analyst | comparison | 1 | sources lisibles | result attribué | non |
| Annoter/relier/router | Analyst | observation/relation | 2 | permission | version/handoff | OPEN-013 |
| Extraire Derived Artifact | Analyst | Artifact | 1 | policy | lineage | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| détecter type/indicateurs | oui | oui | oui | suggestion | extracteurs/règles |
| proposer relations/copies | oui | oui | oui | oui | comparateur/metadata |
| résumer contradictions | oui | oui | oui | oui | table sourcée |
| confirmer identité | humain | non | non | assistance | revue multi-source |

Aucun fichier n’est exécuté; aucune certitude automatique.

## 14. États fonctionnels
`metadata-only`, `previewable`, `restricted`, `partial`, `inconsistent`, `related-candidate`, `derived`, `failed`, `superseded`.

## 15. États d’interface
Loading conserve entry; Empty distingue contenu absent/inaccessible; Partial expose missing ranges; Error garde metadata; Offline lecture; Permission denied masque contenu; Stale montre version.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| File Observation | observation | Workbench/Case | source/type/uncertainty visibles |
| Relationship/comparison | relation/result | Hypothesis/Timeline | candidate, non identité absolue |
| Derived Artifact/handoff | Artifact/package | Static/Reverse/CAP-INV-379 | lineage et restrictions conservées |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Entry | inspecter | CAP-INV-368 | source, path, metadata, permissions | Navigation |
| File/Derived Artifact | analyser | Static/Reverse | parent, content ref, type, objective, restrictions | File view |
| Observation | handoff | CAP-INV-379 | selected metadata/content relation, uncertainty | File view |

## 18. Dépendances
CAP-INV-105/301/311/329/367/369/375/376/379, Shared Preview/Linking/Export, Security, OPEN-005/008/013/014/015.

## 19. Source de vérité
Disk Image/Artifact remain sources; Investigate owns observation/interpretation; Tool results remain Studio; no identity indicator is absolute alone.

## 20. Provenance et audit
Image/session/filesystem/entry, path as displayed, metadata, content access, Tool/version, indicators, comparisons, related sources, annotations, exports/handoffs et actor.

## 21. Permissions fonctionnelles
Metadata read, file preview, sensitive content read, relationship create, comparison, Derived Artifact create/read/export, Static/Reverse handoff. Reveal/export gates remain future.

## 22. Limites et erreurs
Content missing/restricted/corrupt, type conflict, duplicate candidate, container unavailable, preview unsupported, permission denied. Copied file ≠ same provenance; path ≠ stable identity.

## 23. Métriques
Previews denied/restricted, identity conflicts, candidate relationships reviewed, Derived Artifacts with lineage et handoff success.

## 24. Classification de livraison
`defined` / `planned`; no execution, implementation, final identity algorithm or tool choice.

## 25. Critères d’acceptation
**Given** two files share an identity indicator but different provenance **When** compared **Then** candidate relation and differences are visible, without declaring historical identity.

**Given** content permission denied **When** preview opened **Then** metadata remains visible if authorized and content is not leaked or executed.

**Given** no AI **When** analyzed **Then** deterministic type detection, metadata, preview, comparison and human review work.

## 26. Questions ouvertes
OPEN-005/008/013/014/015 remain open; final identity, Artifact and permission models are future.

## 27. Consommateurs documentaires
INV-DSK-001, Static/Reverse, CAP-INV-369/375/376/379, Artifact/Evidence and future Objects/Permissions/Screens.
