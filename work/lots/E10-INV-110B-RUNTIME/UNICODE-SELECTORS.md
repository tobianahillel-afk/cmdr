# E10-INV-110B — Unicode authorization selector hardening

## Reuse-first decision

Use Go's existing standard-library `unicode/utf8` and a short explicit predicate in the existing `invalidScope` function. This avoids a new Unicode normalization library, a custom identity format, or a separate permission framework. Do not normalize and silently merge distinct IDs: fail closed for malformed selector bytes and a narrow set of invisible variation/selectors.

## Problem and correction

The pre-existing `unicode.IsControl` and `unicode.Cf` checks rejected ordinary control/format characters, but not the combining grapheme joiner U+034F, the Mongolian variation selectors U+180B–U+180D, variation selectors U+FE00–U+FE0F and U+E0100–U+E01EF (non-Cf), nor invalid UTF-8 bytes that `encoding/json` can replace with U+FFFD. These can make authorization object identifiers visually ambiguous or change identity across encodings.

The common selector guard now rejects these sequences. It is exercised through the read-only projection, both Case/Timeline/source access boundary adapters, and a negative assertion that visible combining accents remain permitted.

## Boundaries and verification

This is **not** a full Unicode confusables defense or a canonical-ID migration. Authenticated principal, Case membership, tenant/Timeline/source grants still require a real trusted provider; none is bundled. `OPEN-013` and the existing `E10-INV-110B-RUNTIME` READY status remain unchanged. Verify the *new SHA* with Go tests, race, vet, formatting, engineering-kernel checks and a separate L2 reading of this exact diff before claiming acceptance. Historical performance numbers or previous reviews must not be transplanted onto this commit.
