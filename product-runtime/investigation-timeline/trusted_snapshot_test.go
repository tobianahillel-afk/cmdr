package investigationtimeline

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type failingTrustedReader struct{}
func (failingTrustedReader) Read([]byte) (int, error) { return 0, errors.New("unavailable") }

func trustedSnapshotFixture() MemoryTimelineData {
	return MemoryTimelineData{
		Principals: []MemoryTimelinePrincipal{{ID:"reader-a", Tenant:"tenant-a",Environment:"lab"}},
		Cases: []MemoryTimelineCase{{ID:"case-a",Tenant:"tenant-a",Environment:"lab",
			Readers:[]string{"reader-a"},EntryIDs:[]string{"event-a"}}},
		Entries: []MemoryTimelineEntry{{CaseID:"case-a",Readers:[]string{"reader-a"},
			Value:Entry{ID:"event-a",Tenant:"tenant-a",Kind:"observed",SourceTime:"2026-10-01T10:00:00Z",
				Timezone:"+00:00",Version:1,SourceRef:"ref-a",Provenance:"trusted-local"}}},
		Sources: []MemoryTimelineSource{{Identity:SourceOwnerIdentity{
			Reference:"ref-a",Kind:"telemetry-event",ObjectID:"object-a",
			Tenant:"tenant-a",Environment:"lab",CaseID:"case-a"}, Readers:[]string{"reader-a"}}},
	}
}

func TestTrustedSnapshotDecodeCanAssembleRealOwnerChain(t *testing.T) {
	original := trustedSnapshotFixture()
	raw, err := json.Marshal(original)
	if err != nil { t.Fatal(err) }
	decoded, err := DecodeTrustedMemoryTimelineJSON(bytes.NewReader(raw))
	if err != nil { t.Fatal(err) }
	if len(decoded.Principals)!=1 || decoded.Principals[0].ID!="reader-a" ||
		len(decoded.Cases)!=1 || len(decoded.Entries)!=1 || len(decoded.Sources)!=1 {
		t.Fatalf("trusted dataset not preserved: %+v", decoded)
	}
	handler, err := NewMemoryCaseTimelineHTTPHandler(
		func(r *http.Request)(string,error){ return "reader-a",nil }, decoded)
	if err != nil || handler == nil { t.Fatalf("provider rejected imported records: %v",err) }
}

func TestTrustedSnapshotRejectsUnsafeEnvelope(t *testing.T) {
	good, err := json.Marshal(trustedSnapshotFixture())
	if err != nil { t.Fatal(err) }
	cases := []struct{name string;input io.Reader}{
		{"nil-reader",nil},
		{"read-error",failingTrustedReader{}},
		{"empty",strings.NewReader("")},
		{"invalid",strings.NewReader("{bad")},
		{"unknown-field",strings.NewReader(`{"Principals":[],"Cases":[],"Entries":[],"Sources":[],"Surprise":1}`)},
		{"second-object",strings.NewReader(string(good)+" {}")},
		{"second-invalid-value",strings.NewReader(string(good)+" unparsed")},
		{"oversized",strings.NewReader(strings.Repeat("x",MaxTrustedSnapshotBytes+1))},
		{"missing-principals",strings.NewReader(`{"Cases":[{"ID":"case-a"}]}`)},
		{"missing-cases",strings.NewReader(`{"Principals":[{"ID":"reader-a"}]}`)},
	}
	for _,tc:=range cases {
		t.Run(tc.name,func(t *testing.T){
			data,err:=DecodeTrustedMemoryTimelineJSON(tc.input)
			if err==nil || len(data.Principals)>0 || len(data.Cases)>0 {
				t.Fatalf("untrusted envelope accepted: %+v err=%v",data,err)
			}
		})
	}
}
