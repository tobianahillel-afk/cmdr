# CMDR Product Specification

## Purpose

This directory is the canonical product repository for CMDR, a unified cyber-defence platform organised around three complementary consoles:

1. **Command Center** — shared operational picture, prioritisation and coordination.
2. **Investigation Lab** — evidence-led investigation, forensics and analytical reasoning.
3. **Response & Governance** — controlled decisions, authorisation, execution, rollback and audit.

The specification describes product behaviour, information architecture, functional ownership, object lifecycles, permissions, UX expectations and acceptance criteria. It is intentionally independent from a particular frontend framework or backend implementation.

## How to use this repository

Start with:

- [Source of truth rules](00-governance/SOURCE_OF_TRUTH.md)
- [Product vision](01-product/PRODUCT_VISION.md)
- [Information architecture](01-product/INFORMATION_ARCHITECTURE.md)
- [Domain model](02-domain-model/DOMAIN_MODEL.md)
- [State models](02-domain-model/STATE_MODELS.md)
- [Permission model](02-domain-model/PERMISSION_MODEL.md)
- [Global navigation](03-experience/GLOBAL_NAVIGATION.md)
- [Cross-console transitions](04-cross-console/TRANSITIONS.md)

Console specifications:

- [Command Center](10-command-center/README.md)
- [Investigation Lab](20-investigation-lab/README.md)
- [Response & Governance](30-response-governance/README.md)
- [Shared platform capabilities](40-platform/README.md)
- [Quality and acceptance](50-quality/README.md)

## Documentation contract

A page specification owns the behaviour of that page only. Shared objects, states, permissions, visual tokens and cross-console transitions are defined once in their canonical documents and referenced elsewhere. When two documents conflict, the precedence rules in `00-governance/SOURCE_OF_TRUTH.md` apply.

## Current maturity

This foundation captures the complete known CMDR architecture and substantially specifies the workflows already defined. Areas that still require product decisions are explicitly recorded in each document's **Open questions** section and consolidated in `50-quality/OPEN_GAPS.md`.
