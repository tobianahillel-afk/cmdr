package contextenvelope

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

type fixtureSet struct {
	SchemaVersion int           `json:"schema_version"`
	ContractID    string        `json:"contract_id"`
	Cases         []fixtureCase `json:"cases"`
}

type fixtureCase struct {
	ID       string `json:"id"`
	Input    Input  `json:"input"`
	Expected Result `json:"expected"`
}

func TestProjectMatchesPredeclaredContractFixtures(t *testing.T) {
	data, err := os.ReadFile("../../engineering/pilot/context-envelope-fixtures.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures fixtureSet
	decErr := json.Unmarshal(data, &fixtures)
	if decErr != nil {
		t.Fatal(decErr)
	}
	if fixtures.SchemaVersion != 1 || fixtures.ContractID != "PILOT-CONTEXT-ENVELOPE-V1" {
		t.Fatalf("unexpected fixture identity: version=%d contract=%s", fixtures.SchemaVersion, fixtures.ContractID)
	}
	if len(fixtures.Cases) != 11 {
		t.Fatalf("expected 11 predeclared fixtures, got %d", len(fixtures.Cases))
	}

	for _, fixture := range fixtures.Cases {
		fixture := fixture
		t.Run(fixture.ID, func(t *testing.T) {
			got := Project(fixture.Input)
			if !reflect.DeepEqual(got, fixture.Expected) {
				t.Fatalf("projection mismatch\n got: %#v\nwant: %#v", got, fixture.Expected)
			}
		})
	}
}

func TestProjectRejectsInvalidTenantReferences(t *testing.T) {
	base := Input{
		SourceTenantRef:          "tenant-a",
		DestinationTenantRef:     "tenant-a",
		DestinationAuthorization: AuthorizationAllow,
		ContextFresh:             true,
		ReturnOrigin:             "/safe",
	}
	bad := []string{
		"*",
		"tenant *",
		"tenant\tA",
		"tenant\nA",
		strings.Repeat("x", 129),
	}
	for _, value := range bad {
		in := base
		in.SourceTenantRef = value
		got := Project(in)
		if got.Error != ErrorInvalidTenantRef || !got.TraceMinimalOnly || got.ProtectedRefsVisible {
			t.Fatalf("reference %q did not fail closed: %#v", value, got)
		}
	}

	in := base
	in.DestinationTenantRef = ""
	got := Project(in)
	if got.Error != ErrorInvalidTenantRef {
		t.Fatalf("empty destination tenant must be invalid: %#v", got)
	}
}

func TestProjectRejectsUnsafeReturnOrigins(t *testing.T) {
	base := Input{
		SourceTenantRef:          "tenant-a",
		DestinationTenantRef:     "tenant-a",
		DestinationAuthorization: AuthorizationAllow,
		ContextFresh:             true,
	}
	bad := []string{
		"https://example.invalid/path",
		"//example.invalid/path",
		"/safe?token=synthetic",
		"/safe#fragment",
		"/safe\\windows",
		"/safe\x00control",
		strings.Repeat("/x", 600),
	}
	for _, origin := range bad {
		in := base
		in.ReturnOrigin = origin
		got := Project(in)
		if got.Error != ErrorUnsafeReturnOrigin || !got.TraceMinimalOnly || got.ProtectedRefsVisible {
			t.Fatalf("origin %q did not fail closed: %#v", origin, got)
		}
	}
}

func TestProjectUnknownAuthorizationFailsClosed(t *testing.T) {
	got := Project(Input{
		SourceTenantRef:          "tenant-a",
		DestinationTenantRef:     "tenant-a",
		DestinationAuthorization: AuthorizationOutcome("unknown"),
		ContextFresh:             true,
		ReturnOrigin:             "/safe",
	})
	if got.State != StatePermissionDenied ||
		got.TenantRef != "" ||
		got.EnvironmentRef != "" ||
		got.ProtectedRefsVisible ||
		!got.TraceMinimalOnly {
		t.Fatalf("unknown authorization did not fail closed: %#v", got)
	}
}
