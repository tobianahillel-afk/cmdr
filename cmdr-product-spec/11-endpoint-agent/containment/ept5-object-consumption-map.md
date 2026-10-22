---
id: endpoint-ept5-object-consumption-map
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# Endpoint EPT-5 — Object Consumption Map

| Concept | Canonical owner | Endpoint usage | Local operations | Boundary | Provenance |
|---|---|---|---|---|---|
| Response Primitive | Endpoint | technical primitive ref | resolve/read | not Action Request | type/version/target |
| Primitive Eligibility | Endpoint | local assessment | derive/version | not authority | sources/reasons |
| Requested Target State | Govern intent / Endpoint projection | technical comparator | read/project | requested != achieved | Run/Step ref |
| Technical Precheck | Endpoint | readiness | run/version | PASS != universally safe | facts/time |
| Technical Execution | Endpoint | target-side effect | execute under authority | != Response Run | Run/Step/Agent |
| Technical Outcome | Endpoint | raw technical result | record/update | != Result | executor refs |
| Technical Verification Request | Endpoint | technical check contract | create/read | != Govern Verification Plan | execution/criteria |
| Technical Verification Observation | Endpoint | target-side fact | derive/read | != Govern Verification | observations/time |
| Target State Observation | Endpoint | current target fact | observe/version | observed != business objective | source/freshness |
| Response Drift | Endpoint | technical regression | derive | != response failure automatically | before/after refs |
| Technical Reversal | Endpoint | inverse primitive | execute under Govern rollback authority | != Response Rollback | original effect/ref |
| Containment State | Endpoint | isolation/quarantine/control state | observe/update | technical only | target/time |
| Process Control State | Endpoint | process target state | observe | absent != threat resolved | process identity |
| Isolation State | Endpoint | host isolation state | observe | isolation != offline | target/connectivity |
| Network Control State | Endpoint | bounded block state | observe | != host isolation | scope/expiry |
| File Quarantine State | Endpoint | quarantine/release state | observe | != Evidence/verdict | file identity |
| Service Control State | Endpoint | component state | observe | stopped != remediation | component ref |
| Local Session Control State | Endpoint | local session lock/termination | observe | directory state external | session/principal refs |
| Action Request | Govern | authority lineage ref | read only | Endpoint never owns | request/version |
| Approval | Govern | authority evidence ref | read only | Human Gate != Approval | approval/version |
| Decision | Govern | authorized bounds ref | read only | Decision != command | decision/version |
| Response Run | Govern | parent execution context | read/correlate | != Endpoint execution | run/step refs |
| Result | Govern | downstream outcome ref | read only | Endpoint never creates | result/version |
| Endpoint Agent | Endpoint | executor identity | read | no Fleet admin | Agent/version |
| Technical Capability | Endpoint | availability/readiness | read | available != allowed | declaration/time |
| Endpoint Policy | Settings | local restriction projection | read | Policy != authority | version/assignment |
| Human Gate | Studio | workflow provenance | read/link | != Approval/Decision | gate/response |
| Tool Call | Studio | execution provenance | read/link | != Endpoint action | call/version |
| Automation Run | Studio | workflow runtime provenance | read/link | != Response Run | run/version |
| Finding/Evidence | Investigate | source justification refs | read/link | no requalification | source/version |
| Trace/Activity | Shared | generic history mechanisms | link | not fact owner | correlation refs |