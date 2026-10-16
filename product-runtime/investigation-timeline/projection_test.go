package investigationtimeline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type contractFixtures struct {
	Fixtures []struct {
		ID       string `json:"id"`
		Input    Input  `json:"input"`
		Expected struct {
			Status       string   `json:"status"`
			VisibleIDs   []string `json:"visible_ids"`
			VisibleKinds []string `json:"visible_kinds"`
		} `json:"expected"`
	} `json:"fixtures"`
}

func TestProjectMatchesContractFixtures(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "engineering", "implementation", "investigation-timeline", "fixtures.json"))
	if err != nil {
		t.Fatal(err)
	}
	var fixtures contractFixtures
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	if len(fixtures.Fixtures) != 18 {
		t.Fatalf("expected 18 contract cases, got %d", len(fixtures.Fixtures))
	}
	for _, f := range fixtures.Fixtures {
		t.Run(f.ID, func(t *testing.T) {
			got := Project(f.Input)
			if got.Status != f.Expected.Status {
				t.Fatalf("status got %s want %s", got.Status, f.Expected.Status)
			}
			ids, kinds := []string{}, []string{}
			for _, e := range got.Entries {
				ids = append(ids, e.ID)
				kinds = append(kinds, e.Kind)
			}
			if !reflect.DeepEqual(ids, f.Expected.VisibleIDs) || !reflect.DeepEqual(kinds, f.Expected.VisibleKinds) {
				t.Fatalf("projection got ids=%v kinds=%v, expected ids=%v kinds=%v", ids, kinds, f.Expected.VisibleIDs, f.Expected.VisibleKinds)
			}
		})
	}
}

func allowedFixture() Input {
	return Input{Scope: Scope{Tenant: "tenant-a", Environment: "prod", CasePermission: "allow", TimelinePermission: "allow"},
		Entries: []Entry{{ID: "x", Kind: "observed", Tenant: "tenant-a", SourcePermission: "allow",
			Timezone: "+00:00", SourceTime: "2026-10-01T10:00:00Z", SourceRef: "owner:x",
			Provenance: "telemetry:x", Version: 1}}}
}

func TestProjectRejectsInvalidFullScope(t *testing.T) {
	for _, bad := range []string{"", "*", "  "} {
		in := allowedFixture()
		in.Scope.Tenant = bad
		if got := Project(in); got.Status != "deny" || len(got.Entries) != 0 {
			t.Fatalf("scope %q leaked", bad)
		}
		in = allowedFixture()
		in.Scope.Environment = bad
		if got := Project(in); got.Status != "deny" || len(got.Entries) != 0 {
			t.Fatalf("environment %q leaked", bad)
		}
	}
	for _, p := range []string{"deny", "unknown", ""} {
		in := allowedFixture()
		in.Scope.CasePermission = p
		if Project(in).Status != "deny" {
			t.Fatal("case permission must fail closed")
		}
		in = allowedFixture()
		in.Scope.TimelinePermission = p
		if Project(in).Status != "deny" {
			t.Fatal("Timeline Entry permission must fail closed")
		}
	}
}

func TestProjectSuppressesSourceIdentifiers(t *testing.T) {
	for _, p := range []string{"deny", "unknown", "missing"} {
		in := allowedFixture()
		in.Entries[0].ID = "TOP-SECRET"
		in.Entries[0].SourceRef = "secret-ref"
		in.Entries[0].SourcePermission = p
		got := Project(in)
		data, err := json.Marshal(got)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(data), "TOP-SECRET") || strings.Contains(string(data), "secret-ref") {
			t.Fatalf("denied source disclosure: %s", data)
		}
	}
	in := allowedFixture()
	in.Entries[0].Tenant = "another-tenant"
	if len(Project(in).Entries) != 0 {
		t.Fatal("cross-tenant source visible")
	}
}

func TestProjectRejectsDuplicateOrInvalidEntries(t *testing.T) {
	mutations := []func(*Input){
		func(i *Input) { i.Entries = append(i.Entries, i.Entries[0]) },
		func(i *Input) { i.Entries[0].ID = "" },
		func(i *Input) { i.Entries[0].Kind = "invented" },
		func(i *Input) { i.Entries[0].SourcePermission = "garbage" },
		func(i *Input) { i.Entries[0].Provenance = "" },
		func(i *Input) { i.Entries[0].SourceRef = "" },
		func(i *Input) { i.Entries[0].Version = 0 },
		func(i *Input) { i.Entries[0].Supersedes = "x" },
		func(i *Input) { i.Entries[0].ClockUncertaintyMS = -1 },
	}
	for _, change := range mutations {
		in := allowedFixture()
		change(&in)
		if r := Project(in); r.Status != "deny" || len(r.Entries) != 0 {
			t.Fatalf("malformed entry leaked: %+v", r)
		}
	}
}

func TestProjectBoundsAndCallerImmutability(t *testing.T) {
	in := allowedFixture()
	in.Entries = make([]Entry, MaxEntries+1)
	if got := Project(in); got.Status != "deny" {
		t.Fatal("entry bound not enforced")
	}
	in = allowedFixture()
	in.Entries[0].Provenance = strings.Repeat("x", MaxInputBytes+1)
	if got := Project(in); got.Status != "deny" {
		t.Fatal("input byte bound not enforced")
	}
	in = allowedFixture()
	out := Project(in)
	if out.Status != "allow" || len(out.Entries) != 1 {
		t.Fatal("expected authorized result")
	}
	out.Entries[0].ID = "changed"
	if in.Entries[0].ID != "x" {
		t.Fatal("mutation of caller state")
	}
}

func TestProjectTimeOrderingAndUncertainty(t *testing.T) {
	in := allowedFixture()
	second := in.Entries[0]
	second.ID = "a"
	second.SourceTime = "2026-10-01T11:00:00+01:00"
	second.Timezone = "+01:00"
	third := in.Entries[0]
	third.ID = "z"
	third.SourceTime = "2026-10-01T09:00:00Z"
	unknown := in.Entries[0]
	unknown.ID = "unknown"
	unknown.Timezone = "unknown"
	in.Entries = []Entry{second, unknown, third, in.Entries[0]}
	out := Project(in)
	want := []string{"z", "a", "x", "unknown"}
	var got []string
	for _, e := range out.Entries {
		got = append(got, e.ID)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("order=%v want=%v", got, want)
	}
	if out.Status != "uncertain" || !out.Entries[len(out.Entries)-1].TimeUncertain {
		t.Fatal("unknown timezone promoted to known UTC")
	}
	in = allowedFixture()
	in.Entries[0].SourceTime = "not-a-date"
	if got := Project(in); got.Status != "uncertain" || !got.Entries[0].TimeUncertain {
		t.Fatal("malformed timestamp promoted")
	}
	in = allowedFixture()
	in.Entries[0].ClockUncertaintyMS = 500
	if got := Project(in); got.Status != "uncertain" {
		t.Fatal("clock skew not marked")
	}
}

func TestProjectMissingAndSuperseded(t *testing.T) {
	in := allowedFixture()
	old := in.Entries[0]
	old.ID = "old"
	newer := in.Entries[0]
	newer.ID = "new"
	newer.Version = 2
	newer.Supersedes = "old"
	in.Entries = []Entry{old, newer}
	got := Project(in)
	if len(got.Entries) != 2 || got.Entries[0].ID != "new" || got.Entries[1].ID != "old" || got.Entries[0].Supersedes != "old" {
		t.Fatal("supersession or stable order not preserved")
	}
}

// Source version references are sensitive: a permitted entry must not reveal
// identifiers of denied, foreign or missing Timeline Entry predecessors.
func TestProjectMasksInaccessibleSupersessionReferences(t *testing.T) {
	for _, permission := range []string{"deny", "unknown", "missing"} {
		t.Run(permission, func(t *testing.T) {
			in := allowedFixture()
			old := in.Entries[0]
			old.ID = "hidden-version"
			old.SourceRef = "hidden-source"
			old.SourcePermission = permission
			newer := in.Entries[0]
			newer.ID = "public-version"
			newer.Version = 2
			newer.Supersedes = old.ID
			in.Entries = []Entry{newer, old}
			got := Project(in)
			if got.Status != "partial" || len(got.Entries) != 1 || got.Entries[0].Supersedes != "" {
				t.Fatalf("inaccessible predecessor not masked: %+v", got)
			}
			blob, err := json.Marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(string(blob), "hidden-version") || strings.Contains(string(blob), "hidden-source") {
				t.Fatalf("protected predecessor identity disclosed: %s", blob)
			}
		})
	}
	in := allowedFixture()
	in.Entries[0].Supersedes = "unresolved-predecessor"
	got := Project(in)
	if got.Status != "partial" || len(got.Entries) != 1 || got.Entries[0].Supersedes != "" {
		t.Fatalf("unresolved predecessor leaked: %+v", got)
	}
}

func TestProjectTimestampOffsetMustMatchDeclaredTimezone(t *testing.T) {
	for _, zone := range []string{"+00:00", "UTC", "unknown", "+25:00"} {
		in := allowedFixture()
		in.Entries[0].SourceTime = "2026-10-01T11:00:00+01:00"
		in.Entries[0].Timezone = zone
		got := Project(in)
		if got.Status != "uncertain" || len(got.Entries) != 1 || !got.Entries[0].TimeUncertain {
			t.Fatalf("mismatched timezone %q treated as known: %+v", zone, got)
		}
	}
	in := allowedFixture()
	in.Entries[0].SourceTime = "2026-10-01T11:00:00+01:00"
	in.Entries[0].Timezone = "+01:00"
	got := Project(in)
	if got.Status != "allow" || len(got.Entries) != 1 || got.Entries[0].TimeUncertain {
		t.Fatalf("matching offset must be known: %+v", got)
	}
}

func TestProjectRejectsWhitespaceWrappedWildcards(t *testing.T) {
	for _, value := range []string{" * ", "\t*\t", "\n*\n"} {
		in := allowedFixture()
		in.Scope.Tenant = value
		if result := Project(in); result.Status != "deny" || len(result.Entries) != 0 {
			t.Fatalf("wildcard tenant %q not rejected: %+v", value, result)
		}
		in = allowedFixture()
		in.Scope.Environment = value
		if result := Project(in); result.Status != "deny" || len(result.Entries) != 0 {
			t.Fatalf("wildcard environment %q not rejected: %+v", value, result)
		}
	}
}

func TestHiddenEntryIdentifierCollisionCannotDenyVisibleProjection(t *testing.T) {
	base := allowedFixture()
	visible := base.Entries[0]
	hidden := visible
	hidden.SourcePermission = "deny"
	hidden.SourceRef = "protected-source"
	hidden.Provenance = "protected-provenance"
	base.Entries = []Entry{visible, hidden}
	result := Project(base)
	if result.Status != "filtered" || len(result.Entries) != 1 || result.Entries[0].ID != visible.ID {
		t.Fatalf("hidden collision changed authorized projection: %+v", result)
	}
	withoutCollision := allowedFixture()
	hidden.ID = "different-protected-id"
	withoutCollision.Entries = []Entry{visible, hidden}
	other := Project(withoutCollision)
	if !reflect.DeepEqual(result, other) {
		t.Fatalf("hidden identifier influenced visible output: %+v vs %+v", result, other)
	}
	for _, tenant := range []string{"foreign-tenant"} {
		otherHidden := visible
		otherHidden.Tenant = tenant
		withoutCollision.Entries = []Entry{visible, otherHidden}
		got := Project(withoutCollision)
		if got.Status != "filtered" || len(got.Entries) != 1 {
			t.Fatalf("foreign duplicate denied authorized projection: %+v", got)
		}
	}
	duplicateAuthorized := allowedFixture()
	duplicateAuthorized.Entries = []Entry{visible, visible}
	if result := Project(duplicateAuthorized); result.Status != "deny" {
		t.Fatalf("authorized duplicate must be denied: %+v", result)
	}
}

func TestSupersessionRequiresStrictVersionProgression(t *testing.T) {
	base := allowedFixture()
	before := base.Entries[0]
	before.ID = "before"
	after := base.Entries[0]
	after.ID = "after"
	after.Version = 2
	after.Supersedes = before.ID
	base.Entries = []Entry{after, before}
	if result := Project(base); result.Status != "allow" || len(result.Entries) != 2 {
		t.Fatalf("valid version progression rejected: %+v", result)
	}
	bad := base
	bad.Entries = append([]Entry(nil), base.Entries...)
	bad.Entries[0].Version = 1
	if result := Project(bad); result.Status != "deny" || len(result.Entries) != 0 {
		t.Fatalf("same version predecessor accepted: %+v", result)
	}
	bad.Entries[0].Version = 2
	bad.Entries[1].Supersedes = after.ID
	if result := Project(bad); result.Status != "deny" || len(result.Entries) != 0 {
		t.Fatalf("cyclic supersession accepted: %+v", result)
	}
}

func TestOpaqueSourcesCannotControlVisibleValidation(t *testing.T) {
	// A denied, missing or foreign source is already inaccessible: even
	// invalid IDs, duplicate IDs and oversized protected fields cannot
	// affect validation of the one authorized Timeline Entry.
	cases := []struct {
		name       string
		permission string
		tenant     string
		status     string
	}{
		{name: "denied", permission: "deny", tenant: "tenant-a", status: "filtered"},
		{name: "unknown", permission: "unknown", tenant: "tenant-a", status: "filtered"},
		{name: "missing", permission: "missing", tenant: "tenant-a", status: "partial"},
		{name: "foreign", permission: "allow", tenant: "tenant-b", status: "filtered"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			base := allowedFixture()
			hidden := base.Entries[0]
			hidden.Tenant = tc.tenant
			hidden.SourcePermission = tc.permission
			hidden.ID = "protected"
			hidden.SourceRef = "hidden-source"
			base.Entries = append(base.Entries, hidden)
			baseline := Project(base)
			if baseline.Status != tc.status || len(baseline.Entries) != 1 {
				t.Fatalf("invalid baseline: %+v", baseline)
			}
			for _, mutate := range []func(*Entry){
				func(e *Entry) { e.ID = "" },
				func(e *Entry) { e.ID = "x" },
				func(e *Entry) { e.Provenance = strings.Repeat("protected", MaxInputBytes/9+1) },
				func(e *Entry) { e.SourceRef = strings.Repeat("secret", MaxInputBytes/6+1) },
				func(e *Entry) { e.Version = -999; e.Kind = "invalid"; e.SourceTime = "malformed" },
			} {
				in := allowedFixture()
				opaque := hidden
				mutate(&opaque)
				in.Entries = append(in.Entries, opaque)
				got := Project(in)
				if !reflect.DeepEqual(got, baseline) {
					t.Fatalf("protected fields changed visible result: got %+v want %+v", got, baseline)
				}
			}
		})
	}
}

func TestAuthorizedSourceStillEnforcesIdentityAndByteLimits(t *testing.T) {
	for _, mutate := range []func(*Input){
		func(i *Input) { i.Entries[0].ID = "" },
		func(i *Input) { i.Entries = append(i.Entries, i.Entries[0]) },
		func(i *Input) { i.Entries[0].SourceRef = strings.Repeat("x", MaxInputBytes+1) },
	} {
		in := allowedFixture()
		mutate(&in)
		if result := Project(in); result.Status != "deny" || len(result.Entries) != 0 {
			t.Fatalf("authorized invalid input not denied: %+v", result)
		}
	}
	in := allowedFixture()
	in.Entries = make([]Entry, MaxEntries+1)
	if got := Project(in); got.Status != "deny" {
		t.Fatalf("global entry count limit not enforced: %+v", got)
	}
}

func TestProjectRejectsInvisibleVisibleIDsWithoutChargingHiddenEntries(t *testing.T) {
	for _, id := range []string{"x\x00bad", "x\nspoof", "x\u200bhidden", "x\u202ebad"} {
		in := allowedFixture()
		in.Entries[0].ID = id
		if got := Project(in); got.Status != "deny" || len(got.Entries) != 0 {
			t.Fatalf("visible malformed ID accepted: %+v", got)
		}
		in = allowedFixture()
		in.Entries = append(in.Entries, Entry{ID: id, Tenant: "other-tenant", SourcePermission: "allow"})
		if got := Project(in); got.Status != "filtered" || len(got.Entries) != 1 {
			t.Fatalf("opaque foreign ID interfered with visible result: %+v", got)
		}
	}
}

func TestProjectRejectsInvisibleAuthorizedSourceReferences(t *testing.T) {
	for _, ref := range []string{"", " * ", "owner:x\x00suffix", "owner:x\nspoof", "owner:\u200bsecret", "owner:\u202ereversed"} {
		in := allowedFixture()
		in.Entries[0].SourceRef = ref
		got := Project(in)
		if got.Status != "deny" || len(got.Entries) != 0 {
			t.Fatalf("visible malformed source reference %q escaped projection: %+v", ref, got)
		}
	}
}

func TestProjectIgnoresMalformedOpaqueSourceReferences(t *testing.T) {
	for _, permission := range []string{"deny", "unknown", "missing"} {
		in := allowedFixture()
		hidden := in.Entries[0]
		hidden.ID = "unavailable"
		hidden.SourcePermission = permission
		hidden.SourceRef = "owner:\u202ehidden\x00"
		in.Entries = append(in.Entries, hidden)
		got := Project(in)
		wantStatus := "filtered"
		if permission == "missing" {
			wantStatus = "partial"
		}
		if got.Status != wantStatus || len(got.Entries) != 1 || got.Entries[0].SourceRef != "owner:x" {
			t.Fatalf("opaque malformed SourceRef influenced visible result: %+v", got)
		}
	}
}

func TestProjectNegativeZeroLocalOffsetIsUncertain(t *testing.T) {
	cases := []struct {
		name, sourceTime, timezone string
		uncertain bool
	}{
		{"negative-zero-time-positive-label", "2026-10-01T10:00:00-00:00", "+00:00", true},
		{"negative-zero-time-negative-label", "2026-10-01T10:00:00-00:00", "-00:00", true},
		{"positive-zero-time-negative-label", "2026-10-01T10:00:00+00:00", "-00:00", true},
		{"z-time-negative-label", "2026-10-01T10:00:00Z", "-00:00", true},
		{"positive-zero-known", "2026-10-01T10:00:00+00:00", "+00:00", false},
		{"z-time-explicit-offset", "2026-10-01T10:00:00Z", "+00:00", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			in := allowedFixture()
			in.Entries[0].SourceTime = tc.sourceTime
			in.Entries[0].Timezone = tc.timezone
			got := Project(in)
			if len(got.Entries) != 1 || got.Entries[0].TimeUncertain != tc.uncertain {
				t.Fatalf("offset uncertainty %q %q: %+v", tc.sourceTime, tc.timezone, got)
			}
			wantStatus := "allow"
			if tc.uncertain {
				wantStatus = "uncertain"
			}
			if got.Status != wantStatus {
				t.Fatalf("offset %q %q status %s want %s", tc.sourceTime, tc.timezone, got.Status, wantStatus)
			}
			if got.Entries[0].SourceTime != tc.sourceTime || got.Entries[0].Timezone != tc.timezone {
				t.Fatal("original forensic timestamp/offset was overwritten")
			}
		})
	}
}

func TestProjectUnknownNegativeZeroSortsAfterKnownTime(t *testing.T) {
	in := allowedFixture()
	unknown := in.Entries[0]
	unknown.ID = "a-unknown"
	unknown.SourceTime = "2026-10-01T09:00:00-00:00"
	unknown.Timezone = "+00:00"
	known := in.Entries[0]
	known.ID = "z-known"
	known.SourceTime = "2026-10-01T11:00:00+00:00"
	known.Timezone = "+00:00"
	in.Entries = []Entry{unknown, known}
	got := Project(in)
	if got.Status != "uncertain" || len(got.Entries) != 2 ||
		got.Entries[0].ID != known.ID || got.Entries[1].ID != unknown.ID ||
		got.Entries[0].TimeUncertain || !got.Entries[1].TimeUncertain {
		t.Fatalf("unknown local offset promoted to known chronology: %+v", got)
	}
}
