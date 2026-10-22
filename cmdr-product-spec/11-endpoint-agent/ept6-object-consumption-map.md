---
id: endpoint-ept6-object-consumption-map
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-12
source-of-truth: canonical
---
# EPT-6 Object Consumption Map

| Concept | Owner | Endpoint use | Local operations | Boundary | Provenance |
|---|---|---|---|---|---|
| Update Assignment | Settings | desired target projection | receive/compare | Endpoint cannot administer assignment | assignment/Fleet/Policy refs |
| Update Channel/Wave | Settings | rollout context | read/project | no wave scheduling | source/version refs |
| Update Package Reference | release/admin source | metadata/reference | assess availability/compatibility | ref != contents/trust | release/version refs |
| Update Operation/Attempt | Endpoint | local lifecycle | create/transition | not Fleet job/Response Run | attempt/target/package |
| Update Verification | Endpoint | technical post-update checks | derive/version | not Govern verification | expected/observed/health |
| Update Reversion | Endpoint | local version recovery | execute/report if authorized | not Studio/Govern rollback | original/reversion refs |
| Buffered Item | Endpoint | local retained item | queue/state | no storage schema | source/item/time |
| Local Queue | Endpoint | local queue projection | derive/refresh | not Shared Job queue ownership | queue/pressure refs |
| Replay State | Endpoint | reconnect/replay facts | transition | replay != consumption/exactly-once | item/ack refs |
| Recovery State | Endpoint | restart/crash recovery | derive/transition | not Shared Recovery ownership | pre/post state refs |
| Resource Pressure | Endpoint | local pressure fact | observe/derive | no resource scheduler | monitor/policy refs |
| Dependency State | Endpoint | local dependency fact | reassess | restored != capability healthy | component/capability refs |
| Self-Protection State | Endpoint | protection observation | derive/refresh | no global security policy | component/config refs |
| Tamper Observation | Endpoint | technical candidate | create/version | not Finding/Incident | reason/source/time |
| Runtime Privilege State | Endpoint | local execution context | derive | privilege != permission | required/observed refs |
| Secret Reference | Settings/Security | opaque reference | consume status only | never raw secret | reference/status |
| Local Audit Event | Endpoint | local event/provenance | create/append conceptually | not Shared Trace/immutable ledger | source/actor/outcome/time |
| Endpoint Security State | Endpoint | consumer projection | derive/handoff | not Finding/Incident/Result | source-state refs |
| Endpoint Agent | Endpoint | primary local subject | observe/update local facts | Fleet admin external | Agent/version/tenant |
| Endpoint Agent Fleet | Settings | admin context | read projection | no Fleet mutation | Fleet refs |
| Endpoint Policy | Settings | effective constraints | read projection | no policy administration | Policy/version refs |
| Studio Deployment | Studio | non-overlap reference | none | Endpoint update != Studio deployment | deployment refs |
| Govern Response Run/Result | Govern | authority/consumer refs | none | no Result creation | Run/Decision/Result refs |
| Shared Job/Trace/Activity/Recovery | Shared | generic mechanism refs | consume/link | no ownership transfer | correlation refs |

No physical schema is introduced.