# Authentication and SSO

## Objective

Authenticate human and service principals securely and provide session assurance appropriate to action risk.

## Scope

This specification owns the page-local behaviour of **Authentication and SSO** in **Shared Platform**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

Platform Product Lead

## Affected objects

- Principal
- Session
- Identity provider
- Credential
- Tenant

## Features

- OIDC/SAML SSO
- Local break-glass only where approved
- MFA and step-up authentication
- Service identities
- Session management and revocation
- Login audit and risk signals

## UX and interactions

- Authentication failures do not disclose account existence
- Step-up explains why it is required
- Tenant selection occurs only after identity validation
- Session expiry protects unsaved work where possible

## Permissions

Authentication is prerequisite; administration requires `platform.admin` or delegated identity permissions.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Sessions are pending/active/step_up_required/expired/revoked.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

External IdP, secrets/keys, audit, tenancy.

## Acceptance criteria

- SSO metadata and certificates rotate safely
- Revocation takes effect within defined SLA
- Critical approvals can require fresh authentication
- Service identities cannot use human UI sessions
- All outcomes are audited

## Open questions

- Which IdPs are launch requirements?
- What is the break-glass governance process?
