import assert from "node:assert/strict";
import test from "node:test";

import {
  normalizeDeepLinkState,
  parseDeepLink,
  serializeDeepLinkState,
  StateValidationError,
} from "./state.mjs";

function validState() {
  return {
    tenantId: "tenant-a",
    environmentId: "prod-eu",
    from: "2026-09-24T10:00:00+02:00",
    to: "2026-09-24T11:00:00+02:00",
    query: "opaque query / not a selected dialect",
    sources: ["edr", "firewall", "edr"],
    filters: {
      severity: ["high", "critical", "high"],
      "event.kind": "alert",
    },
    activeRef: "telemetry-event:evt-42",
    returnTo: "/investigate/case/42?tab=evidence#event",
  };
}

test("round-trip-is-deterministic", () => {
  const first = serializeDeepLinkState(validState());
  const parsed = parseDeepLink(first, { tenantId: "tenant-a", environmentId: "prod-eu" });
  const second = serializeDeepLinkState(parsed);
  assert.equal(second, first);
  assert.deepEqual(parsed.sources, ["edr", "firewall"]);
  assert.deepEqual(parsed.filters.severity, ["critical", "high"]);
});

test("refresh-reconstruction-preserves-deep-link-state", () => {
  const url = serializeDeepLinkState(validState());
  const reconstructed = parseDeepLink(url);
  assert.equal(reconstructed.tenantId, "tenant-a");
  assert.equal(reconstructed.environmentId, "prod-eu");
  assert.equal(reconstructed.from, "2026-09-24T08:00:00.000Z");
  assert.equal(reconstructed.to, "2026-09-24T09:00:00.000Z");
  assert.equal(reconstructed.query, validState().query);
  assert.equal(reconstructed.returnTo, validState().returnTo);
  assert.ok(Object.isFrozen(reconstructed));
  assert.ok(Object.isFrozen(reconstructed.filters));
});

test("tenant-mismatch-is-rejected", () => {
  const url = serializeDeepLinkState(validState());
  assert.throws(
    () => parseDeepLink(url, { tenantId: "tenant-b" }),
    (error) => error instanceof StateValidationError && error.code === "tenant-mismatch",
  );
});

test("tenant-wildcard-is-rejected", () => {
  const state = validState();
  state.tenantId = "*";
  assert.throws(
    () => normalizeDeepLinkState(state),
    (error) => error instanceof StateValidationError && error.code === "wildcard-scope",
  );
});

test("environment-mismatch-is-rejected", () => {
  const url = serializeDeepLinkState(validState());
  assert.throws(
    () => parseDeepLink(url, { environmentId: "prod-us" }),
    (error) => error instanceof StateValidationError && error.code === "environment-mismatch",
  );
});

test("authorization-state-cannot-be-deeplinked", () => {
  const safe = serializeDeepLinkState(validState());
  assert.throws(
    () => parseDeepLink(`${safe}&permission=perm.investigate.search.execute`),
    (error) => error instanceof StateValidationError && error.code === "unknown-url-field",
  );
  assert.throws(
    () => serializeDeepLinkState({ ...validState(), permission: "admin" }),
    (error) => error instanceof StateValidationError && error.code === "unknown-state-field",
  );
});

test("transient-state-is-excluded-from-url", () => {
  const safe = serializeDeepLinkState(validState());
  for (const value of ["request-123", "focus-row-9", "window-20-40"]) {
    assert.equal(safe.includes(value), false);
  }
  assert.throws(
    () => serializeDeepLinkState({ ...validState(), inFlightRequestId: "request-123" }),
    (error) => error instanceof StateValidationError && error.code === "unknown-state-field",
  );
});

test("source-wildcard-and-empty-sources-are-rejected", () => {
  assert.throws(() => normalizeDeepLinkState({ ...validState(), sources: [] }));
  assert.throws(() => normalizeDeepLinkState({ ...validState(), sources: ["*"] }));
});

test("time-range-must-be-ordered", () => {
  assert.throws(() =>
    normalizeDeepLinkState({
      ...validState(),
      from: "2026-09-24T11:00:00Z",
      to: "2026-09-24T10:00:00Z",
    }),
  );
});

test("duplicate-singleton-url-state-is-rejected", () => {
  const url = serializeDeepLinkState(validState());
  assert.throws(() => parseDeepLink(`${url}&tenant=tenant-a`));
});

test("external-return-url-is-rejected", () => {
  assert.throws(() =>
    normalizeDeepLinkState({ ...validState(), returnTo: "//evil.example/path" }),
  );
  assert.throws(() =>
    normalizeDeepLinkState({ ...validState(), returnTo: "https://evil.example/path" }),
  );
});

test("unsafe-filter-key-is-rejected", () => {
  const url = serializeDeepLinkState(validState());
  assert.throws(() => parseDeepLink(`${url}&filter.__proto__=polluted`));
});

test("query-and-filter-values-remain-opaque", () => {
  const state = normalizeDeepLinkState({
    ...validState(),
    query: "field:* AND raw=\"<value>\"",
    filters: { syntax: ["*", "a:b", "<literal>"] },
  });
  const roundTrip = parseDeepLink(serializeDeepLinkState(state));
  assert.equal(roundTrip.query, state.query);
  assert.deepEqual(roundTrip.filters.syntax, ["*", "<literal>", "a:b"]);
});
