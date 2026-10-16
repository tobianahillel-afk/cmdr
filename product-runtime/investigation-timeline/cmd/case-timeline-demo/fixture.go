package main

import timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"

// demoData is deliberately synthetic and immutable after the owner adapter
// copies it. None of the source references point to a real external system.
func demoData() timeline.MemoryTimelineData {
	return timeline.MemoryTimelineData{
		Principals: []timeline.MemoryTimelinePrincipal{
			{ID: "analyst-a", Tenant: "demo-tenant-a", Environment: "lab"},
			{ID: "analyst-b", Tenant: "demo-tenant-b", Environment: "lab"},
		},
		Cases: []timeline.MemoryTimelineCase{
			{ID: "case-a", Tenant: "demo-tenant-a", Environment: "lab",
				Readers: []string{"analyst-a"}, EntryIDs: []string{"event-new", "event-old", "restricted"}},
			{ID: "case-b", Tenant: "demo-tenant-b", Environment: "lab",
				Readers: []string{"analyst-b"}, EntryIDs: []string{"event-b"}},
		},
		Entries: []timeline.MemoryTimelineEntry{
			{CaseID: "case-a", Value: timeline.Entry{
				ID: "event-new", Tenant: "demo-tenant-a", Kind: "observed",
				SourceTime: "2026-10-01T11:00:00Z", Timezone: "+00:00", Version: 2,
				Supersedes: "event-old", SourceRef: "demo-telemetry-a",
				Provenance: "synthetic-telemetry",
			}, Readers: []string{"analyst-a"}},
			{CaseID: "case-a", Value: timeline.Entry{
				ID: "event-old", Tenant: "demo-tenant-a", Kind: "inferred",
				SourceTime: "2026-10-01T10:00:00Z", Timezone: "+00:00", Version: 1,
				SourceRef: "demo-evidence-a", Provenance: "synthetic-evidence",
			}, Readers: []string{"analyst-a"}},
			{CaseID: "case-a", Value: timeline.Entry{
				ID: "restricted", Tenant: "demo-tenant-a", Kind: "observed",
				SourceTime: "2026-10-01T09:00:00Z", Timezone: "+00:00", Version: 1,
				SourceRef: "demo-hidden-source", Provenance: "restricted-synthetic",
			}, Readers: []string{"analyst-a"}},
			{CaseID: "case-b", Value: timeline.Entry{
				ID: "event-b", Tenant: "demo-tenant-b", Kind: "observed",
				SourceTime: "2026-10-01T08:00:00Z", Timezone: "+00:00", Version: 1,
				SourceRef: "demo-incident-b", Provenance: "synthetic-incident",
			}, Readers: []string{"analyst-b"}},
		},
		Sources: []timeline.MemoryTimelineSource{
			{Identity: timeline.SourceOwnerIdentity{
				Reference: "demo-telemetry-a", Kind: "telemetry-event", ObjectID: "demo-event-a",
				Tenant: "demo-tenant-a", Environment: "lab", CaseID: "case-a",
			}, Readers: []string{"analyst-a"}},
			{Identity: timeline.SourceOwnerIdentity{
				Reference: "demo-evidence-a", Kind: "evidence", ObjectID: "demo-evidence-a",
				Tenant: "demo-tenant-a", Environment: "lab", CaseID: "case-a",
			}, Readers: []string{"analyst-a"}},
			{Identity: timeline.SourceOwnerIdentity{
				Reference: "demo-hidden-source", Kind: "evidence", ObjectID: "demo-secret",
				Tenant: "demo-tenant-a", Environment: "lab", CaseID: "case-a",
			}, Readers: []string{}},
			{Identity: timeline.SourceOwnerIdentity{
				Reference: "demo-incident-b", Kind: "incident", ObjectID: "demo-incident",
				Tenant: "demo-tenant-b", Environment: "lab", CaseID: "case-b",
			}, Readers: []string{"analyst-b"}},
		},
	}
}
