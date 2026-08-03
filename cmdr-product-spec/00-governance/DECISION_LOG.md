# Decision Log

This log records durable product decisions. New entries must describe context, decision, consequences and affected canonical documents.

## D-001 — Three-console operating model

**Decision:** CMDR keeps three visible top-level consoles: Command Center, Investigation Lab and Response & Governance.

**Rationale:** Operational prioritisation, deep evidence work and controlled authority are distinct cognitive modes. Combining them would overload pages and weaken role separation.

**Consequences:** Shared objects remain continuous across consoles; transitions preserve context and are first-class product behaviour.

## D-002 — Evidence before authority

**Decision:** Investigation produces findings and evidence packages; Response & Governance owns approval and execution authority.

**Consequences:** Investigation may recommend an action but cannot bypass policy gates or approval requirements.

## D-003 — Complete forensic chain

**Decision:** Investigation Lab includes separate capabilities for static analysis, dynamic sandboxing, reverse engineering, debugging, memory forensics and disk/artifact forensics.

**Consequences:** These capabilities may share evidence and navigation but must not be collapsed into a generic “malware analysis” page.

## D-004 — Single source of truth

**Decision:** Shared entities, states, permissions, design patterns and transition contracts are defined once and referenced by page specifications.

**Consequences:** Local convenience never justifies duplicated canonical definitions.

## D-005 — Decision re-enters operations

**Decision:** Approved, rejected, conditioned and executed response decisions feed back into the Command Center operational picture and incident timeline.

**Consequences:** The operating picture must show governance status, execution progress, rollback availability and business impact.
