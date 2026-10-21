# Event Search frontend runtime

This directory is reserved for the first-party browser runtime of `CAP-INV-002 / INV-EVS-001`.

Current state: **bounded state/deep-link core, accessible semantic shell, bounded incremental result renderer and backend-neutral transport integration implemented and verified; adversarial end-to-end frontend validation remains next**.

The accepted engineering architecture is `ENG-DEC-0001`: first-party Web Platform ES modules, native semantic HTML by default, bounded incremental DOM updates, URL/History deep-link state and no third-party UI framework/runtime by default.

The boundary-registration, state, shell, result-rendering and transport-integration lots are complete. The executable frontend surface includes first-party deep-link state, semantic shell states, keyed bounded incremental result rendering with position preservation, backend-neutral execution/retry transport, AbortController cancellation, permission-denied fail-closed handling, partial-result preservation and retry-as-new-run semantics. No final streaming transport, Saved Search behavior or Case-link mutation is implemented yet.

Explicitly unresolved/out of scope here:
- Case-link mutation while `OPEN-013` remains open;
- final query dialect;
- search index or storage engine;
- search provider;
- Saved Search implementation;
- persisted result model;
- production Search SLO.
