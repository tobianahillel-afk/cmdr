# E10-INV-110B — legacy Timeline selector preflight deduplication

The owner-resolved adapter `ProjectFromOwnerSnapshots` already rejects duplicate caller-selected Timeline Entry IDs before calling any authority. The older caller-snapshot adapter `ProjectWithAccessChecker` bounded ID bytes but did not pre-reject duplicates. Repeated identical visible IDs therefore triggered redundant Case/Shared/source authority calls before the pure `Project` eventually denied the result.

This change reuses the owner-boundary's existing Go map-based selector preflight, instead of building a new authorization framework, third-party dependency or cache. It rejects repeated **visible, same-tenant, syntactically valid** IDs before any provider I/O. Opaque cross-tenant or malformed IDs do not enter the deduplication set: they cannot veto an otherwise visible Case projection. The legacy snapshot path is still explicitly untrusted for source metadata; use `ProjectFromOwnerSnapshots` with trusted owner resolver for serious integrations.

Tests in `access_test.go` assert zero provider calls for a duplicate visible selection, preserve legitimate output with a foreign same-ID entry, foreign oversized ID or malformed opaque same-tenant ID, and check that only the expected Case/Timeline/source obligations ran.

No Product Spec / Shared ownership / OPEN-013 changes and no claim of real provider connectivity or production SLO. CI and independent L2 must be verified for the exact new commit before this evidence can be promoted; B remains READY.
