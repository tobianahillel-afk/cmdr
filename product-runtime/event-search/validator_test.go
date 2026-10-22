package eventsearch

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

type contractFixtureSet struct {
	Cases []struct {
		ID    string `json:"id"`
		Input struct {
			TenantRef            string   `json:"tenant_ref"`
			AuthorizedTenants    []string `json:"authorized_tenants"`
			EnvironmentRef       string   `json:"environment_ref"`
			TimeStart            string   `json:"time_start"`
			TimeEnd              string   `json:"time_end"`
			Sources              []string `json:"sources"`
			AuthorizedSources    []string `json:"authorized_sources"`
			Permissions          []string `json:"permissions"`
			QueryPresent         bool     `json:"query_present"`
			QueryValid           bool     `json:"query_valid"`
			CaseLinkRequested    bool     `json:"case_link_requested"`
		} `json:"input"`
		Expected struct {
			Allowed              bool   `json:"allowed"`
			ErrorCode            string `json:"error_code"`
			ProtectedDataVisible bool   `json:"protected_data_visible"`
		} `json:"expected"`
	} `json:"cases"`
}

func TestValidateContractFixtures(t *testing.T) {
	data, err := os.ReadFile("../../engineering/implementation/event-search/fixtures.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures contractFixtureSet
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}

	lifecycleOwnedByNextTask := map[string]bool{
		"illegal-terminal-transition-is-rejected": true,
		"running-job-can-be-cancelled":            true,
	}
	consumed := 0
	for _, fixture := range fixtures.Cases {
		fixture := fixture
		if lifecycleOwnedByNextTask[fixture.ID] {
			continue
		}
		consumed++
		t.Run(fixture.ID, func(t *testing.T) {
			in := Input{
				TenantRef:         fixture.Input.TenantRef,
				AuthorizedTenants: fixture.Input.AuthorizedTenants,
				EnvironmentRef:    fixture.Input.EnvironmentRef,
				TimeStart:         fixture.Input.TimeStart,
				TimeEnd:           fixture.Input.TimeEnd,
				Sources:           fixture.Input.Sources,
				AuthorizedSources: fixture.Input.AuthorizedSources,
				Permissions:       fixture.Input.Permissions,
				QueryPresent:      fixture.Input.QueryPresent,
				QueryValid:        fixture.Input.QueryValid,
				CorrelationID:     "corr-fixture-001",
				QueryVersion:      "query-v1",
				CaseLinkRequested: fixture.Input.CaseLinkRequested,
			}
			got := Validate(in)
			if got.Allowed != fixture.Expected.Allowed {
				t.Fatalf("allowed: got %t want %t", got.Allowed, fixture.Expected.Allowed)
			}
			if string(got.ErrorCode) != fixture.Expected.ErrorCode {
				t.Fatalf("error code: got %q want %q", got.ErrorCode, fixture.Expected.ErrorCode)
			}
			if got.ProtectedDataVisible != fixture.Expected.ProtectedDataVisible {
				t.Fatalf("protected visibility: got %t want %t", got.ProtectedDataVisible, fixture.Expected.ProtectedDataVisible)
			}
			if !got.AuditRequired {
				t.Fatal("every validation attempt must remain auditable")
			}
			if got.Allowed {
				wantStart, _ := time.Parse(time.RFC3339, fixture.Input.TimeStart)
				wantEnd, _ := time.Parse(time.RFC3339, fixture.Input.TimeEnd)
				if got.Envelope.TenantRef != fixture.Input.TenantRef ||
					got.Envelope.EnvironmentRef != fixture.Input.EnvironmentRef ||
					!got.Envelope.TimeStart.Equal(wantStart) ||
					!got.Envelope.TimeEnd.Equal(wantEnd) ||
					!reflect.DeepEqual(got.Envelope.Sources, fixture.Input.Sources) ||
					got.Envelope.CorrelationID != "corr-fixture-001" ||
					got.Envelope.QueryVersion != "query-v1" {
					t.Fatalf("validated envelope lost scope/provenance: %#v", got.Envelope)
				}
			} else if !isZeroEnvelope(got.Envelope) {
				t.Fatalf("denied validation leaked a populated envelope: %#v", got.Envelope)
			}
		})
	}
	if consumed != 12 {
		t.Fatalf("consumed %d B-owned fixtures, want 12", consumed)
	}
}

func TestValidateRejectsInvalidProvenance(t *testing.T) {
	base := validInput()
	tests := []struct {
		name string
		edit func(*Input)
		code ErrorCode
	}{
		{name: "missing correlation", edit: func(in *Input) { in.CorrelationID = "" }, code: ErrorInvalidCorrelationID},
		{name: "correlation whitespace", edit: func(in *Input) { in.CorrelationID = "corr 1" }, code: ErrorInvalidCorrelationID},
		{name: "correlation control", edit: func(in *Input) { in.CorrelationID = "corr\n1" }, code: ErrorInvalidCorrelationID},
		{name: "missing query version", edit: func(in *Input) { in.QueryVersion = "" }, code: ErrorInvalidQueryVersion},
		{name: "query version unsafe", edit: func(in *Input) { in.QueryVersion = "v1@example" }, code: ErrorInvalidQueryVersion},
		{name: "oversized correlation", edit: func(in *Input) { in.CorrelationID = strings.Repeat("a", 257) }, code: ErrorInvalidCorrelationID},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			in := base
			tt.edit(&in)
			got := Validate(in)
			if got.Allowed || got.ErrorCode != tt.code || !isZeroEnvelope(got.Envelope) {
				t.Fatalf("unexpected provenance rejection: %#v", got)
			}
		})
	}
}

func TestValidateRejectsAdditionalFailClosedInputs(t *testing.T) {
	base := validInput()
	tests := []struct {
		name string
		edit func(*Input)
		code ErrorCode
	}{
		{name: "unparseable start", edit: func(in *Input) { in.TimeStart = "not-time" }, code: ErrorInvalidTimeRange},
		{name: "equal range", edit: func(in *Input) { in.TimeEnd = in.TimeStart }, code: ErrorInvalidTimeRange},
		{name: "blank source", edit: func(in *Input) { in.Sources = []string{""}; in.AuthorizedSources = []string{""} }, code: ErrorUnauthorizedSource},
		{name: "query absent", edit: func(in *Input) { in.QueryPresent = false }, code: ErrorInvalidQuery},
		{name: "query invalid", edit: func(in *Input) { in.QueryValid = false }, code: ErrorInvalidQuery},
		{name: "permission missing", edit: func(in *Input) { in.Permissions = in.Permissions[:3] }, code: ErrorPermissionDenied},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			in := base
			tt.edit(&in)
			got := Validate(in)
			if got.Allowed || got.ErrorCode != tt.code || got.ProtectedDataVisible {
				t.Fatalf("unexpected fail-closed result: %#v", got)
			}
		})
	}
}

func TestValidateCopiesSources(t *testing.T) {
	in := validInput()
	got := Validate(in)
	if !got.Allowed {
		t.Fatalf("valid input rejected: %#v", got)
	}
	in.Sources[0] = "mutated-after-validation"
	if got.Envelope.Sources[0] != "source-a" {
		t.Fatal("validated envelope aliases caller-owned source slice")
	}
}

func TestAllowedReferenceRune(t *testing.T) {
	for _, r := range []rune("Az09-_.:/") {
		if !allowedReferenceRune(r) {
			t.Fatalf("expected %q to be allowed", r)
		}
	}
	for _, r := range []rune{' ', '@', '\n', 'é'} {
		if allowedReferenceRune(r) {
			t.Fatalf("expected %q to be rejected", r)
		}
	}
}

func validInput() Input {
	return Input{
		TenantRef:         "tenant-a",
		AuthorizedTenants: []string{"tenant-a"},
		EnvironmentRef:    "env-prod",
		TimeStart:         "2026-09-24T10:00:00Z",
		TimeEnd:           "2026-09-24T11:00:00Z",
		Sources:           []string{"source-a"},
		AuthorizedSources: []string{"source-a"},
		Permissions: []string{
			PermissionSearchExecute,
			PermissionQueryRead,
			PermissionSearchManage,
			PermissionTelemetryRead,
		},
		QueryPresent:  true,
		QueryValid:    true,
		CorrelationID: "corr-valid-001",
		QueryVersion:  "query-v1",
	}
}

func isZeroEnvelope(got Envelope) bool {
	return got.TenantRef == "" && got.EnvironmentRef == "" &&
		got.TimeStart.IsZero() && got.TimeEnd.IsZero() &&
		len(got.Sources) == 0 && got.CorrelationID == "" && got.QueryVersion == ""
}
