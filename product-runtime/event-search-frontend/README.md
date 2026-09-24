# Event Search frontend runtime

This directory is reserved for the first-party browser runtime of `CAP-INV-002 / INV-EVS-001`.

Current state: **preimplementation marker only**.

The accepted engineering architecture is `ENG-DEC-0001`: first-party Web Platform ES modules, native semantic HTML by default, bounded incremental DOM updates, URL/History deep-link state and no third-party UI framework/runtime by default.

No executable frontend code is allowed in the boundary-registration lot.

Explicitly unresolved/out of scope here:
- Case-link mutation while `OPEN-013` remains open;
- final query dialect;
- search index or storage engine;
- search provider;
- Saved Search implementation;
- persisted result model;
- production Search SLO.
