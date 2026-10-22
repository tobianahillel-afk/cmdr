# AI start here — CMDR development

This is the universal entry point for any autonomous development agent.

## Read order

1. `AGENTS.md`
2. `engineering/state/current.json`
3. `work/graph.json`
4. the selected work-unit manifest
5. only the product/engineering sources referenced by that manifest

Do not crawl the entire product specification by default.

## Truth layers

1. **Product truth** — `cmdr-product-spec/**`
2. **Engineering truth** — `engineering/**`
3. **Execution truth** — `work/**`
4. **Repository reality** — Git commits, branches, PRs, checks and produced evidence

Recorded state is a cache. Repository reality is used to reconcile whether recorded work actually happened.

## Bootstrap state

This branch starts the autonomous engineering foundation from the validated Phase-6 documentary product-spec baseline at commit `99c3826cdcbc79bad51631cf6462450f92e012eb`.

The product specification is documentary-complete but runtime implementation is not claimed. The first objective is therefore to build the development system that will transform canonical product obligations into bounded, traceable, verifiable implementation work.

## Autonomous continuation

A fresh agent receiving only "continue development" must:

- recover state;
- select the first safe READY unit;
- load bounded context;
- execute or decompose it;
- verify it;
- persist evidence and handoff;
- leave the next state derivable.

If the recorded preferred next unit is stale, recompute it from the work graph rather than trusting the cache.
