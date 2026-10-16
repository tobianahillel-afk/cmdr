package investigationtimeline

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// FuzzProjectionOpaqueFieldsDoNotInterfere uses independent input variation
// rather than repeating the contract's expected fixture IDs/statuses.
// Fuzz seed cases also run during ordinary go test on GitHub Actions.
func FuzzProjectionOpaqueFieldsDoNotInterfere(f *testing.F) {
	for _, seed := range []struct {
		secret   string
		category uint8
	}{
		{"confidential", 0}, {"跨租户", 1}, {"\u200b", 2}, {"private/path:123", 3},
	} {
		f.Add(seed.secret, seed.category)
	}
	f.Fuzz(func(t *testing.T, secret string, category uint8) {
		// The fuzz engine supplies arbitrarily long or non-UTF8 strings;
		// keep individual cases bounded without losing adversarial bytes.
		if len(secret) > 256 {
			return
		}
		in := allowedFixture()
		visible := in.Entries[0]
		hidden := visible
		hidden.ID = "protected-id:" + secret
		hidden.SourceRef = "protected-ref:" + secret
		hidden.Provenance = "protected-provenance:" + secret
		hidden.Supersedes = "protected-prior:" + secret
		switch category % 4 {
		case 0:
			hidden.SourcePermission = "deny"
		case 1:
			hidden.SourcePermission = "unknown"
		case 2:
			hidden.SourcePermission = "missing"
		case 3:
			hidden.Tenant = "foreign-tenant"
		}
		in.Entries = append(in.Entries, hidden)
		before := append([]Entry(nil), in.Entries...)
		result := Project(in)
		wantStatus := "filtered"
		if category%4 == 2 {
			wantStatus = "partial"
		}
		if result.Status != wantStatus || len(result.Entries) != 1 || result.Entries[0].ID != visible.ID {
			t.Fatalf("hidden entry changed visible projection status/content: %+v", result)
		}
		if !reflect.DeepEqual(in.Entries, before) {
			t.Fatal("projection mutated caller-owned input")
		}
		encoded, err := json.Marshal(result)
		if err != nil {
			t.Fatal(err)
		}
		for _, forbidden := range []string{"protected-id:", "protected-ref:", "protected-provenance:", "protected-prior:"} {
			if strings.Contains(string(encoded), forbidden) {
				t.Fatalf("non-authorized field appeared in projection: %s", forbidden)
			}
		}
		// Mutate only opaque field contents, preserving the source's
		// non-accessible classification. The visible result must be identical.
		in.Entries[1].ID = ""
		in.Entries[1].SourceRef = strings.Repeat("hidden-", len(secret)+1)
		in.Entries[1].Provenance = ""
		in.Entries[1].Version = -1
		if got := Project(in); !reflect.DeepEqual(got, result) {
			t.Fatalf("masked source data influenced visible result: got %+v want %+v", got, result)
		}
	})
}
