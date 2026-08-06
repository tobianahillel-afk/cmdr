---
id: unresolved-decisions
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-06
source-of-truth: canonical
---
# Unresolved decisions

The programme retains **17 open decisions**. OPEN-009 remains the only historically resolved item. Phase 4B.3B.1 closes none.

| Relevant decision | Application in 4B.3B.1 | Disposition |
|---|---|---|
| OPEN-005 | forensic engine selection only | remains open; not reused for Intelligence |
| OPEN-008 | source/platform availability and support | remains open |
| OPEN-011 / OPEN-012 | Mobile and Cloud scope | remain open and out of scope |
| OPEN-013 | reversible Class-2 knowledge mutations | remains open |
| OPEN-014 | Artifact/Attachment/Intelligence Material relations and retention | remains open |
| OPEN-015 | Tool Call/Automation Run and cross-product provenance | remains open |
| OPEN-017 | Detection runtime, target language and portability | remains open and Detection-only |
| **OPEN-018** | Threat Intelligence ontology, interoperability and exchange strategy | **created open; no option selected** |

## OPEN-018 — Threat intelligence ontology, interoperability and exchange strategy
**Status:** open  
**Owner:** Product Architecture  
**Consumers:** Threat Intelligence, Cases/Hunts, Analysis Workbench, Detection Engineering, Command, Platform Settings, Studio, Govern, Shared, Objects, Permissions, Technique and future 4B.3B.2.

### Decision question
How will CMDR represent Observables, Indicators, Threat Entity candidates, Malware/Tool/Infrastructure/Campaign knowledge, TTPs, Sightings, relationships, confidence, markings and lifecycle while preserving portability and future interoperability?

### Options retained without selection
1. CMDR canonical ontology with adapters.
2. Primary alignment on an external standard.
3. Hybrid canonical model and external representations.
4. Native models per source with a projection layer.
5. Progressive combination by object and consumer.

### Constraints
- no standard, protocol, message format, schema, provider or proprietary extension model is selected;
- no STIX-like schema, TAXII-like protocol, API, graph database, identifier format or implementation is created;
- Observable, Indicator, Entity, relationship and confidence distinctions remain explicit;
- source markings, restrictions, versions, contradictions and provenance survive any future conversion;
- external exchange, dissemination and publication remain 4B.3B.2.

All other open decisions retain their prior wording and status.
