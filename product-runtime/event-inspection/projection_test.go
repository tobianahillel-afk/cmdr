package eventinspection

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

type contractFixtureSet struct {
	Cases []struct {
		ID string `json:"id"`
		Input struct {
			TenantRef string `json:"tenant_ref"`
			AuthorizedTenants []string `json:"authorized_tenants"`
			EnvironmentRef string `json:"environment_ref"`
			EventRef string `json:"event_ref"`
			EventTenantRef string `json:"event_tenant_ref"`
			Permissions []string `json:"permissions"`
			RawAccessDecision string `json:"raw_access_decision"`
			RenderedAccessDecision string `json:"rendered_access_decision"`
			RawPayloadPresent bool `json:"raw_payload_present"`
			RawDerivedSensitivePresent bool `json:"raw_derived_sensitive_present"`
			RenderedPayloadPresent bool `json:"rendered_payload_present"`
			EnrichmentsPresent bool `json:"enrichments_present"`
			EnrichmentProducer string `json:"enrichment_producer"`
			EnrichmentVersion string `json:"enrichment_version"`
			EnrichmentFreshness string `json:"enrichment_freshness"`
			EnrichmentPresentedAsSource bool `json:"enrichment_presented_as_source"`
			PivotRequested bool `json:"pivot_requested"`
			MutationRequested string `json:"mutation_requested"`
		} `json:"input"`
		Expected struct {
			Allowed bool `json:"allowed"`
			ErrorCode string `json:"error_code"`
			RawVisible bool `json:"raw_visible"`
			RenderedVisible bool `json:"rendered_visible"`
			RawDerivedSensitiveVisible bool `json:"raw_derived_sensitive_visible"`
			AuditRequired bool `json:"audit_required"`
		} `json:"expected"`
	} `json:"cases"`
}

func TestProjectMatchesPredeclaredContractFixtures(t *testing.T) {
	data, err := os.ReadFile("../../engineering/implementation/event-inspection/fixtures.json")
	if err != nil { t.Fatal(err) }
	var fixtures contractFixtureSet
	if err := json.Unmarshal(data, &fixtures); err != nil { t.Fatal(err) }

	consumed := 0
	for _, fixture := range fixtures.Cases {
		fixture := fixture
		if fixture.Input.PivotRequested { continue }
		consumed++
		t.Run(fixture.ID, func(t *testing.T) {
			in := fixtureInput(fixture.Input.TenantRef, fixture.Input.AuthorizedTenants,
				fixture.Input.EnvironmentRef, fixture.Input.EventRef, fixture.Input.EventTenantRef,
				fixture.Input.Permissions, AccessDecision(fixture.Input.RawAccessDecision),
				AccessDecision(fixture.Input.RenderedAccessDecision))
			in.MutationRequested = fixture.Input.MutationRequested
			if fixture.Input.RawPayloadPresent { in.Event.RawPayload = []byte("raw-payload") }
			if fixture.Input.RawDerivedSensitivePresent { in.Event.RawDerivedSensitive = map[string]string{"credential":"sensitive"} }
			if fixture.Input.RenderedPayloadPresent { in.Event.RenderedPayload = []byte("rendered-payload") }
			if fixture.Input.EnrichmentsPresent {
				in.Event.Enrichments = []EventEnrichment{{
					Producer: fixture.Input.EnrichmentProducer, Version: fixture.Input.EnrichmentVersion,
					Freshness: fixture.Input.EnrichmentFreshness,
					PresentedAsSource: fixture.Input.EnrichmentPresentedAsSource,
					Values: map[string]string{"risk":"high"},
				}}
			}
			got := Project(in)
			if got.Allowed != fixture.Expected.Allowed { t.Fatalf("allowed: got %t want %t", got.Allowed, fixture.Expected.Allowed) }
			if string(got.ErrorCode) != fixture.Expected.ErrorCode { t.Fatalf("error: got %q want %q", got.ErrorCode, fixture.Expected.ErrorCode) }
			if got.AuditRequired != fixture.Expected.AuditRequired { t.Fatalf("audit mismatch") }
			if !got.Allowed {
				if got.Projection != nil { t.Fatalf("denied request leaked projection: %#v", got.Projection) }
				return
			}
			if got.Projection == nil { t.Fatal("allowed request has no projection") }
			if (len(got.Projection.RawPayload)>0) != fixture.Expected.RawVisible { t.Fatal("raw visibility mismatch") }
			if (len(got.Projection.RenderedPayload)>0) != fixture.Expected.RenderedVisible { t.Fatal("rendered visibility mismatch") }
			if (len(got.Projection.RawDerivedSensitive)>0) != fixture.Expected.RawDerivedSensitiveVisible { t.Fatal("raw-derived visibility mismatch") }
			assertProjectionStructure(t, got.Projection, in)
		})
	}
	if consumed != 15 { t.Fatalf("consumed %d B-owned fixtures, want 15", consumed) }
}

func TestRawDeniedNeverLeaksOrSynthesizesRawContent(t *testing.T) {
	in := validInput()
	in.RawAccessDecision = AccessDeny
	in.Event.RawPayload = []byte("raw-secret")
	in.Event.RawDerivedSensitive = map[string]string{"token":"raw-derived-secret"}
	in.Event.RenderedPayload = []byte("rendered-safe")
	got := Project(in)
	if !got.Allowed || got.Projection == nil { t.Fatalf("valid rendered-only projection rejected: %#v", got) }
	if got.Projection.RawPayload != nil || got.Projection.RawDerivedSensitive != nil { t.Fatalf("raw denial leaked raw content: %#v", got.Projection) }
	if string(got.Projection.RenderedPayload) != "rendered-safe" { t.Fatal("rendered payload hidden or synthesized") }
}

func TestProjectionDeepCopiesCallerOwnedData(t *testing.T) {
	in := validInput()
	in.Event.RawPayload=[]byte("raw"); in.Event.RenderedPayload=[]byte("rendered")
	in.Event.RawDerivedSensitive=map[string]string{"secret":"one"}
	in.Event.Enrichments=[]EventEnrichment{{Producer:"engine",Version:"v1",Freshness:"fresh",Values:map[string]string{"risk":"medium"}}}
	got:=Project(in)
	if !got.Allowed || got.Projection==nil { t.Fatalf("valid input rejected: %#v",got) }
	in.AuthorizedTenants[0]="mutated"; in.Event.RawPayload[0]='X'; in.Event.RenderedPayload[0]='X'
	in.Event.SourceFields["source.ip"]="mutated"; in.Event.NormalizedFields["source.ip"]="mutated"
	in.Event.MissingFields[0]="mutated"; in.Event.RawDerivedSensitive["secret"]="mutated"
	in.Event.Enrichments[0].Values["risk"]="mutated"
	p:=got.Projection
	if string(p.RawPayload)!="raw" || string(p.RenderedPayload)!="rendered" ||
		p.SourceFields["source.ip"]!="192.0.2.10" || p.NormalizedFields["source.ip"]!="192.0.2.10" ||
		p.MissingFields[0]!="user.name" || p.RawDerivedSensitive["secret"]!="one" ||
		p.Enrichments[0].Values["risk"]!="medium" { t.Fatalf("projection aliases caller data: %#v",p) }
}

func TestProjectionPreservesBoundedEventStatesWithoutSynthesis(t *testing.T) {
	for _, state := range []EventState{StateSourceUnavailable,StateTombstone,StatePartial,StateRestricted} {
		state:=state
		t.Run(string(state),func(t *testing.T){
			in:=validInput(); in.Event.State=state; in.Event.RawPayload=nil; in.Event.RenderedPayload=nil
			got:=Project(in)
			if !got.Allowed || got.Projection==nil { t.Fatalf("state rejected: %#v",got) }
			if got.Projection.State!=state { t.Fatalf("state changed: %s",got.Projection.State) }
			if got.Projection.RawPayload!=nil || got.Projection.RenderedPayload!=nil { t.Fatal("payload synthesized") }
		})
	}
}

func TestProjectionMarksStaleEnrichmentAndPreservesProvenance(t *testing.T) {
	in:=validInput()
	in.Event.Enrichments=[]EventEnrichment{{Producer:"intel-engine",Version:"2026.09",Freshness:"stale",Values:map[string]string{"reputation":"known"}}}
	got:=Project(in)
	if !got.Allowed || got.Projection==nil || len(got.Projection.Enrichments)!=1 { t.Fatalf("enrichment failed: %#v",got) }
	e:=got.Projection.Enrichments[0]
	if e.Producer!="intel-engine" || e.Version!="2026.09" || e.Freshness!="stale" || !e.Stale { t.Fatalf("provenance lost: %#v",e) }
}

func TestProjectRejectsAdditionalFailClosedInputs(t *testing.T) {
	tests:=[]struct{name string; edit func(*Input); code ErrorCode}{
		{"missing environment",func(in *Input){in.EnvironmentRef=""},ErrorMissingEnvironment},
		{"invalid raw decision",func(in *Input){in.RawAccessDecision="unknown"},ErrorInvalidAccessDecision},
		{"invalid rendered decision",func(in *Input){in.RenderedAccessDecision=""},ErrorInvalidAccessDecision},
		{"invalid correlation",func(in *Input){in.CorrelationID="corr bad"},ErrorInvalidCorrelationID},
		{"invalid state",func(in *Input){in.Event.State="invented"},ErrorInvalidEventState},
		{"unsupported mutation",func(in *Input){in.MutationRequested="unknown"},ErrorUnsupportedMutation},
	}
	for _,tt:=range tests { tt:=tt; t.Run(tt.name,func(t *testing.T){
		in:=validInput(); tt.edit(&in); got:=Project(in)
		if got.Allowed || got.ErrorCode!=tt.code || got.Projection!=nil || !got.AuditRequired { t.Fatalf("unexpected rejection: %#v",got) }
	})}
}

func TestProjectRejectsInvalidEnrichmentProvenance(t *testing.T) {
	tests:=[]EventEnrichment{
		{Producer:"",Version:"v1",Freshness:"fresh"},
		{Producer:"engine",Version:"",Freshness:"fresh"},
		{Producer:"engine",Version:"v1",Freshness:""},
		{Producer:"engine",Version:"v1",Freshness:"fresh",PresentedAsSource:true},
	}
	for i,e:=range tests {
		in:=validInput(); in.Event.Enrichments=[]EventEnrichment{e}; got:=Project(in)
		if got.Allowed || got.Projection!=nil { t.Fatalf("case %d allowed: %#v",i,got) }
		want:=ErrorInvalidEnrichmentProvenance; if e.PresentedAsSource { want=ErrorDerivedAsSource }
		if got.ErrorCode!=want { t.Fatalf("case %d got %s want %s",i,got.ErrorCode,want) }
	}
}

func TestOpaqueReferenceValidation(t *testing.T) {
	for _,v:=range []string{"corr-1","a_b.c:d/e","A09"} { if !validOpaqueReference(v){t.Fatalf("expected valid %q",v)} }
	for _,v:=range []string{""," bad","bad ","a@b","é",strings.Repeat("a",257)} { if validOpaqueReference(v){t.Fatalf("expected invalid %q",v)} }
}

func fixtureInput(tenant string, authorized []string, environment,eventRef,eventTenant string, permissions []string, raw,rendered AccessDecision) Input {
	in:=validInput(); in.TenantRef=tenant; in.AuthorizedTenants=append([]string(nil),authorized...)
	in.EnvironmentRef=environment; in.Event.Ref=eventRef; in.Event.TenantRef=eventTenant
	in.Permissions=append([]string(nil),permissions...); in.RawAccessDecision=raw; in.RenderedAccessDecision=rendered
	in.Event.RawPayload=nil; in.Event.RenderedPayload=nil; in.Event.RawDerivedSensitive=nil; in.Event.Enrichments=nil
	return in
}

func validInput() Input {
	return Input{
		TenantRef:"tenant-a",AuthorizedTenants:[]string{"tenant-a"},EnvironmentRef:"env-prod",
		Permissions:[]string{PermissionTelemetryRead},RawAccessDecision:AccessAllow,RenderedAccessDecision:AccessAllow,
		CorrelationID:"corr-event-001",
		Event:Event{
			Ref:"evt-001",TenantRef:"tenant-a",State:StateAvailable,
			Source:SourceMetadata{SourceRef:"source-a",DataSourceRef:"data-source-a",Kind:"telemetry"},
			Parser:ParserMetadata{ParserRef:"parser-a",Version:"v1",Status:"parsed"},
			SourceFields:map[string]string{"source.ip":"192.0.2.10"},
			NormalizedFields:map[string]string{"source.ip":"192.0.2.10"},
			MissingFields:[]string{"user.name"},RawPayload:[]byte("raw"),RenderedPayload:[]byte("rendered"),
		},
	}
}

func assertProjectionStructure(t *testing.T,p *Projection,in Input){
	t.Helper()
	if p.TenantRef!=in.TenantRef || p.EnvironmentRef!=in.EnvironmentRef || p.EventRef!=in.Event.Ref ||
		p.CorrelationID!=in.CorrelationID || p.State!=in.Event.State || p.Source!=in.Event.Source || p.Parser!=in.Event.Parser {
		t.Fatalf("projection lost identity/provenance: %#v",p)
	}
	if !reflect.DeepEqual(p.SourceFields,in.Event.SourceFields) || !reflect.DeepEqual(p.NormalizedFields,in.Event.NormalizedFields) ||
		!reflect.DeepEqual(p.MissingFields,in.Event.MissingFields) { t.Fatalf("projection lost distinguished fields: %#v",p) }
}
