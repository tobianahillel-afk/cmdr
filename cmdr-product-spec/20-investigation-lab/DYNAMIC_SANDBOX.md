# Dynamic Sandbox

## Objective

Execute suspicious artefacts in controlled environments and capture behavioural evidence.

## Scope

This specification owns the page-local behaviour of **Dynamic Sandbox** in **Investigation Lab**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

DFIR Lead

## Affected objects

- Evidence
- Sandbox profile
- Analysis job
- Process entity
- Network entity
- Derived artefact

## Features

- Profile selection and detonation controls
- Process, file, registry and network telemetry
- Screenshots and memory captures
- Behaviour timeline
- Anti-analysis indicators
- Extracted artefacts
- Replayable result package

## UX and interactions

- Impact and isolation policy are shown before run
- Internet simulation/egress mode is explicit
- Live run view cannot be confused with host response
- Analysts can mark notable behaviours and promote evidence
- Timeout and partial collection are clear

## Permissions

`sandbox.execute`; sensitive profiles or external egress may require elevated ABAC conditions; export requires `evidence.export`.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Sandbox jobs use analysis job states; agent/worker health is separate.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Sandbox infrastructure, isolated networks, evidence storage, threat intel.

## Acceptance criteria

- Samples cannot reach production networks
- Profile and image versions are recorded
- All captured artefacts have provenance
- Timeout produces explicit partial status
- Run is reproducible from stored parameters

## Open questions

- Which operating-system images are required?
- When is controlled real internet access permitted?
