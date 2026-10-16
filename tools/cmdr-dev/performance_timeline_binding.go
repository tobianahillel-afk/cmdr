package main

import (
	"fmt"
	"strings"
)

const timelineBoundaryID = "investigation-timeline-runtime"

type timelineWorkloadBinding struct {
	performanceID string
	deepID        string
	workloadKey   string
	handlerKey    string
}

// Bind workload identity to an actual compiled operation. A mere handler allowlist
// does not prevent a Case Queue probe being substituted for a Timeline probe.
var timelineWorkloadBindings = []timelineWorkloadBinding{
	{"PERF-TGT-INV-TIMELINE-PROJECTION", "PERF-DEEP-INV-TIMELINE-PROJECTION", "investigation-timeline-projection-v1", "builtin-investigation-timeline-projection-preimplementation-v1"},
	{"PERF-TGT-INV-TIMELINE-NAVIGATION", "PERF-DEEP-INV-TIMELINE-NAVIGATION", "investigation-timeline-navigation-v1", "builtin-investigation-timeline-navigation-preimplementation-v1"},
}

func validateTimelinePerformanceTargetBinding(target PerformanceTarget) error {
	for _, binding := range timelineWorkloadBindings {
		if target.ID == binding.performanceID {
			if target.BoundaryID != timelineBoundaryID || target.WorkloadKey != binding.workloadKey ||
				target.OwnerPath != "product-runtime/investigation-timeline/**" {
				return fmt.Errorf("Timeline performance target %s has mismatched boundary/workload/owner", target.ID)
			}
			return nil
		}
	}
	if target.BoundaryID == timelineBoundaryID || strings.HasPrefix(target.ID, "PERF-TGT-INV-TIMELINE-") ||
		strings.HasPrefix(target.WorkloadKey, "investigation-timeline-") {
		return fmt.Errorf("unrecognized or substituted Timeline performance target %s", target.ID)
	}
	return nil
}

func validateTimelineBenchmarkWorkloadBinding(workload BenchmarkWorkloadDefinition) error {
	for _, binding := range timelineWorkloadBindings {
		if workload.Key == binding.workloadKey {
			if workload.HandlerKey != binding.handlerKey {
				return fmt.Errorf("Timeline workload %s must use its exact compiled handler", workload.Key)
			}
			return nil
		}
	}
	if strings.HasPrefix(workload.HandlerKey, "builtin-investigation-timeline-") {
		return fmt.Errorf("Timeline handler cannot be bound to unrelated workload %s", workload.Key)
	}
	return nil
}

func validateTimelineDeepPerformanceBinding(target DeepPerformanceTarget) error {
	for _, binding := range timelineWorkloadBindings {
		if target.ID == binding.deepID || target.PerformanceTargetID == binding.performanceID {
			if target.ID != binding.deepID || target.PerformanceTargetID != binding.performanceID ||
				target.HandlerKey != binding.handlerKey {
				return fmt.Errorf("Timeline deep target %s must match exact target and handler", target.ID)
			}
			return nil
		}
	}
	if strings.HasPrefix(target.HandlerKey, "builtin-investigation-timeline-") ||
		strings.HasPrefix(target.ID, "PERF-DEEP-INV-TIMELINE-") {
		return fmt.Errorf("unrecognized or substituted Timeline deep target %s", target.ID)
	}
	return nil
}
