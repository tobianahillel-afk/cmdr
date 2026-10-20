package eventinspection

import "testing"

// A missing Go map entry yields the zero string and must not be treated as a present field.
func TestPreparePivotRejectsMissingEmptyValuedField(t *testing.T) {
	cases := []struct {
		name string
		prepare func(*Projection)
	}{
		{"nil maps", func(p *Projection) {
			p.SourceFields = nil
			p.NormalizedFields = nil
			p.RawDerivedSensitive = nil
			p.Enrichments = nil
		}},
		{"unrelated fields", func(p *Projection) {
			p.SourceFields = map[string]string{"source.present": ""}
			p.NormalizedFields = map[string]string{"normalized.present": ""}
			p.RawDerivedSensitive = map[string]string{"raw.present": ""}
			p.Enrichments = []ProjectionEnrichment{{
				Producer: "engine", Version: "v1",
				Values: map[string]string{"enrichment.present": ""},
			}}
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := mustProject(t)
			tc.prepare(p)
			got := PreparePivot(p, PivotSelection{
				Field: "absent.field", Value: "", ValuePresent: true, ValueAuthorized: true,
				TimeStart: "2026-09-25T10:00:00Z", TimeEnd: "2026-09-25T11:00:00Z",
				ReturnOrigin: "event-search:run-001:row-17",
			})
			if got.Allowed || got.Draft != nil || got.ErrorCode != ErrorInvalidPivotContext || !got.AuditRequired {
				t.Fatalf("missing field accepted as pivot: %#v", got)
			}
		})
	}
}

func TestPreparePivotAllowsPresentEmptyValuedField(t *testing.T) {
	cases := []struct {
		name, field, origin string
		prepare func(*Projection)
	}{
		{"source", "source.empty", "source", func(p *Projection) {
			p.SourceFields = map[string]string{"source.empty": ""}
		}},
		{"normalized", "normalized.empty", "normalized", func(p *Projection) {
			p.NormalizedFields = map[string]string{"normalized.empty": ""}
		}},
		{"raw-derived", "raw.empty", "raw-derived-sensitive", func(p *Projection) {
			p.RawDerivedSensitive = map[string]string{"raw.empty": ""}
		}},
		{"enrichment", "enrichment.empty", "enrichment:engine@v1", func(p *Projection) {
			p.Enrichments = []ProjectionEnrichment{{
				Producer: "engine", Version: "v1",
				Values: map[string]string{"enrichment.empty": ""},
			}}
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := mustProject(t)
			tc.prepare(p)
			got := PreparePivot(p, PivotSelection{
				Field: tc.field, Value: "", ValuePresent: true, ValueAuthorized: true,
				TimeStart: "2026-09-25T10:00:00Z", TimeEnd: "2026-09-25T11:00:00Z",
				ReturnOrigin: "event-search:run-001:row-17",
			})
			if !got.Allowed || got.Draft == nil || got.ErrorCode != ErrorNone || got.Draft.FieldOrigin != tc.origin {
				t.Fatalf("present empty field rejected or misattributed: %#v", got)
			}
		})
	}
}
