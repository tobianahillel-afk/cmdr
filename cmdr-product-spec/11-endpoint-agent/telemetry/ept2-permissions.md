# Endpoint EPT-2 — Functional Permissions

EPT-2 identifies functional access needs without choosing final RBAC/ABAC or renaming existing namespaces.

- telemetry source read;
- observation read;
- restricted telemetry read;
- sensitive process/user-session/network metadata read;
- sensor health and normalization metadata read;
- raw/source projection read when policy allows;
- capability declaration/availability/provenance read;
- cross-tenant telemetry access explicitly denied unless a future canonical authorization model permits it;
- telemetry export preparation only through Shared/authorized export mechanisms.

Existing Endpoint Agent and Settings Fleet/Policy permission families remain canonical. No new Endpoint Screen ID is created.
