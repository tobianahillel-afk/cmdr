# Intelligence

## Objective

Relate threat-intelligence knowledge to current operations and explain why an indicator or actor context matters.

## Scope

This specification owns the page-local behaviour of **Intelligence** in **Command Center**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

SOC Operations Lead

## Affected objects

- Indicator
- Threat actor/profile
- Campaign
- Alert
- Incident
- Entity
- Finding

## Features

- Indicator search and enrichment
- Campaign and actor summaries
- Matches to current incidents
- Source reliability and confidence
- Expiry and revocation handling
- Promote relevant context into an investigation

## UX and interactions

- Intel is clearly separated from observed evidence
- Users can inspect source and retrieval time
- Expired or disputed intel is visibly degraded
- Pivoting preserves filters and tenant context

## Permissions

`command.read`; external enrichment may require `integration.manage`; promoting context requires `case.manage` or `evidence.collect` depending on representation.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Intelligence lifecycle is source-specific; any derived Finding uses canonical Finding states.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Threat-intelligence integrations, entity model, case service.

## Acceptance criteria

- Every enrichment shows provenance and licence restrictions
- Indicator matches are reproducible
- Revoked intel cannot silently remain authoritative
- Cross-tenant matches expose only authorised summaries

## Open questions

- Which intel standards and providers are first-class?
- How is confidence normalised across providers?
