# Acceptance Strategy

## Levels

1. **Document acceptance** — required sections, ownership, canonical references and internal links are valid.
2. **Design acceptance** — user flows, states, permissions and accessibility are represented.
3. **Contract acceptance** — API, event and permission contracts trace to canonical product behaviour.
4. **Implementation acceptance** — automated and manual evidence proves criteria.
5. **Operational acceptance** — realistic tenant, workload, failure and recovery scenarios pass.

## Required evidence

Each implemented page must link test cases to its acceptance criteria. High-risk response, tenancy, audit and evidence-integrity workflows require negative tests. Canvas or graph views require accessible alternatives. Cross-console transitions require refresh, back-navigation and insufficient-permission tests.

## Release gates

No page is “complete” because it renders. It is complete only when empty/loading/error/partial/offline/permission-denied states, audit events and data-scope controls pass.
