# Documentation Completeness Checklist

A document is structurally complete when it is non-empty and contains a clear objective, scope, functional owner, affected objects, features, UX/interactions, permissions, states, dependencies, acceptance criteria and open questions.

A domain is semantically complete when:

- its pages cover the full user workflow;
- shared definitions are referenced rather than copied;
- every mutation has permission and audit implications;
- failure, empty, loading, partial and permission-denied behaviour is defined;
- cross-console entry and exit points are explicit;
- acceptance criteria can be converted into tests;
- remaining uncertainty is consolidated in `OPEN_GAPS.md`.

This foundation satisfies structural completeness. Semantic completeness is strongest for the already-defined three-console architecture and reference ransomware flow; technology and policy selections remain open where no decision has yet been made.
