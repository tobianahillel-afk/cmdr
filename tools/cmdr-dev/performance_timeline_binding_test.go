package main

import "testing"

func TestTimelinePerformanceTargetBindings(t *testing.T) {
	for _, binding := range timelineWorkloadBindings {
		valid := PerformanceTarget{ID: binding.performanceID, BoundaryID: timelineBoundaryID,
			OwnerPath: "product-runtime/investigation-timeline/**", WorkloadKey: binding.workloadKey}
		if err := validateTimelinePerformanceTargetBinding(valid); err != nil {
			t.Fatal(err)
		}
		mutations := []PerformanceTarget{
			{ID: valid.ID, BoundaryID: "case-queue-runtime", OwnerPath: valid.OwnerPath, WorkloadKey: valid.WorkloadKey},
			{ID: valid.ID, BoundaryID: valid.BoundaryID, OwnerPath: valid.OwnerPath, WorkloadKey: "case-queue-projection-v1"},
			{ID: valid.ID, BoundaryID: valid.BoundaryID, OwnerPath: "product-runtime/case-queue/**", WorkloadKey: valid.WorkloadKey},
		}
		for _, bad := range mutations {
			if err := validateTimelinePerformanceTargetBinding(bad); err == nil {
				t.Fatalf("accepted substituted target %#v", bad)
			}
		}
	}
}

func TestTimelineBenchmarkWorkloadBindings(t *testing.T) {
	for _, binding := range timelineWorkloadBindings {
		valid := BenchmarkWorkloadDefinition{Key: binding.workloadKey, HandlerKey: binding.handlerKey}
		if err := validateTimelineBenchmarkWorkloadBinding(valid); err != nil {
			t.Fatal(err)
		}
		valid.HandlerKey = "builtin-case-queue-projection-v1"
		if err := validateTimelineBenchmarkWorkloadBinding(valid); err == nil {
			t.Fatalf("accepted substituted benchmark handler for %s", binding.workloadKey)
		}
		valid.Key = "case-queue-projection-v1"
		valid.HandlerKey = binding.handlerKey
		if err := validateTimelineBenchmarkWorkloadBinding(valid); err == nil {
			t.Fatalf("accepted Timeline handler bound to Case Queue workload")
		}
	}
}

func TestTimelineDeepPerformanceBindings(t *testing.T) {
	for _, binding := range timelineWorkloadBindings {
		valid := DeepPerformanceTarget{ID: binding.deepID, PerformanceTargetID: binding.performanceID, HandlerKey: binding.handlerKey}
		if err := validateTimelineDeepPerformanceBinding(valid); err != nil {
			t.Fatal(err)
		}
		for _, bad := range []DeepPerformanceTarget{
			{ID: binding.deepID, PerformanceTargetID: binding.performanceID, HandlerKey: "builtin-case-queue-projection-v1"},
			{ID: "PERF-DEEP-CASE-QUEUE-PROJECTION", PerformanceTargetID: binding.performanceID, HandlerKey: binding.handlerKey},
			{ID: binding.deepID, PerformanceTargetID: "PERF-TGT-CASE-QUEUE-PROJECTION", HandlerKey: binding.handlerKey},
		} {
			if err := validateTimelineDeepPerformanceBinding(bad); err == nil {
				t.Fatalf("accepted substituted deep target %#v", bad)
			}
		}
	}
}
