# Security, Privacy and Resilience

## Objective

Define cross-cutting product requirements that protect CMDR, its tenants and sensitive evidence.

## Scope

This specification owns the page-local behaviour of **Security, Privacy and Resilience** in **Shared Platform**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

Platform Product Lead

## Affected objects

- All domain objects
- Secret
- Encryption key
- Backup
- Service
- Security event

## Features

- Encryption in transit and at rest
- Secret isolation
- Secure sample handling
- Least privilege
- Backups and disaster recovery
- Service degradation modes
- Privacy minimisation
- Security monitoring and vulnerability management

## UX and interactions

- Security controls fail closed where authority or isolation is at risk
- Degraded mode explains unavailable guarantees
- Sensitive content is masked by default
- Recovery does not create duplicate actions

## Permissions

Security controls apply independently of ordinary product permissions.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Service states healthy/degraded/unavailable/recovering; security incidents follow Incident state model.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

All platform components, key management, observability, incident response.

## Acceptance criteria

- RTO/RPO are defined by data class
- Backups are encrypted and restoration-tested
- Malware samples are isolated
- Secrets never enter logs
- Failover preserves idempotency and audit continuity

## Open questions

- What deployment models are supported?
- Which compliance baselines are launch targets?
