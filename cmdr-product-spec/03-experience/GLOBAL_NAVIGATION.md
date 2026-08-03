# Global Navigation

## Objective

Preserve orientation and context across the three-console product.

## Persistent header

The header contains CMDR identity, console tabs, tenant/context selector, time range when applicable, global search, notifications, tasks and user menu.

## Sidebar

Each console owns a stable left sidebar. The active page, console colour and keyboard focus are visible. Sidebar order follows user workflow rather than alphabetical order.

## Context bar

When an Incident, Case, Finding, Response request or Run is active, a context bar displays the object identifier, status, severity/impact, owner and relevant linked objects. It persists through allowed transitions.

## Deep links

Every primary object and stable filtered view has a shareable URL. Links preserve tenant, selected object and safe filter state, but do not embed secrets or raw sample content.

## Back behaviour

Browser back returns to the previous meaningful product state without losing unsaved work. Cross-console transitions provide an explicit “return to source” link.

## Acceptance criteria

Navigation is keyboard operable, URL-addressable and resilient to refresh. Users never lose tenant context silently.
