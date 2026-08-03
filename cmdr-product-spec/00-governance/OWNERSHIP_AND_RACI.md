# Ownership and RACI

## Objective

Assign functional accountability without coupling ownership to a named individual.

| Domain | Accountable owner | Responsible contributors | Consulted | Informed |
|---|---|---|---|---|
| Product vision and IA | Head of Product | Product Architecture | SOC, DFIR, Governance leads | Engineering, Design |
| Command Center | SOC Operations Lead | Command product squad | Threat Intel, IT Ops | SOC users |
| Investigation Lab | DFIR Lead | Investigation product squad | Detection Engineering, Malware Research | SOC and Legal |
| Response & Governance | Response Governance Lead | Response product squad | CISO, IT owners, Legal | SOC and auditors |
| Shared platform | Platform Product Lead | Platform engineering | Security, SRE, Data | All product squads |
| UX system | Design Lead | Product Design | Accessibility, domain leads | Engineering |
| Permissions and audit | Security Architecture Lead | IAM and Platform teams | Compliance, Legal | Product owners |
| Documentation quality | Product Architecture | All document owners | QA, Engineering | Stakeholders |

## Rules

Accountability stays with the functional role even when implementation is delegated. A page may have only one accountable functional owner. Cross-console workflows require approval from every affected accountable owner.

## Escalation

Conflicting requirements are recorded in `DECISION_LOG.md`. Security, privacy and legal constraints override convenience, but the trade-off and user impact must be documented.
