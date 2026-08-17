---
id: govern-action-center-decision-package
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-004, REQ-PROD-015, REQ-PROD-020]
open_decisions: [OPEN-007, OPEN-013, OPEN-015]
---
# Decision Preparation Package — GOV-1 Supporting Contract

## Purpose

Support `CAP-GOV-014` with a stable list of decision-preparation inputs. This package is a review composition, not a canonical Decision, Approval, Response Run or execution contract.

## Required composition for decision readiness

Where applicable to the request, the package exposes by reference and exact version:
- Action Request and requester;
- source product and return origin;
- Incident, Case, Finding and Evidence references with restrictions;
- exact target and included/excluded scope;
- impact, risk, uncertainty and reversibility assessments;
- expected outcome and alternatives;
- applicable Policy Evaluations and outcomes;
- Policy conflicts and Exception Candidates/authorized exception context;
- Authority Requirement and Authority Context;
- approver eligibility/SoD result;
- Approval Requests and Approvals;
- delegation/substitution/escalation/emergency context;
- unresolved questions and information requests;
- recommendations, including AI recommendations, with producer/provenance;
- proposed Decision dispositions, conditions, time bounds, rollback and verification requirements;
- dissent/challenges and reviewer annotations.

## Readiness

A package may be `draft`, `incomplete`, `under-review`, `changes-requested`, `review-ready` or `decision-ready`. `decision-ready` means that the review prerequisites required by the current Policy/authority context are present and current; it is **not** an Approval or Decision.

## Recommendations

Human, deterministic, workflow or AI recommendations are inputs only. They must expose producer, source factors and uncertainty where applicable. No recommendation may mutate the effective disposition.

## Changes

A material change to Action Request, target/scope, Policy, authority, Approval or risk/reversibility input invalidates affected readiness and requires explicit diff/review. Prior package versions remain historical evidence.

## Boundary with Execution Handoff

Decision Preparation Package feeds `CAP-GOV-015`. Only after a valid `approve` or `approve-with-conditions` Decision may `CAP-GOV-016` create an **Execution Handoff Package**. The two packages are distinct and neither is a Response Run.
