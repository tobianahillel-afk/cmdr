package main

import "testing"

func TestEvaluateComplexityWithinBudget(t *testing.T) {
	manifest := validManifestV2()
	manifest.Complexity = ComplexityEstimateV2{
		PrimaryComponents: 1, EstimatedFilesChanged: 6, EstimatedNetLOC: 300,
		ProductionFiles: 5, NetProductionLOC: 250, BoundedContexts: 1,
		SeparatelyTestableBehaviors: 1,
	}
	result := evaluateComplexity(manifest, "READY")
	if result.Decision != "within-budget" || result.ReadinessViolation {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestEvaluateComplexitySoftWarningDoesNotBlockReady(t *testing.T) {
	manifest := validManifestV2()
	manifest.Complexity = ComplexityEstimateV2{
		PrimaryComponents: 2, EstimatedFilesChanged: 9, EstimatedNetLOC: 450,
		ProductionFiles: 8, NetProductionLOC: 400, BoundedContexts: 1,
		SeparatelyTestableBehaviors: 1,
	}
	result := evaluateComplexity(manifest, "READY")
	if result.Decision != "warning" || len(result.Warnings) != 2 || result.ReadinessViolation {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestEvaluateComplexityHardSplitBlocksReady(t *testing.T) {
	manifest := validManifestV2()
	manifest.Complexity = ComplexityEstimateV2{
		PrimaryComponents: 2, EstimatedFilesChanged: 13, EstimatedNetLOC: 500,
		ProductionFiles: 10, NetProductionLOC: 500, BoundedContexts: 2,
		DistinctSecurityModels: 2, SeparatelyTestableBehaviors: 2,
	}
	result := evaluateComplexity(manifest, "READY")
	if result.Decision != "split-required" || !result.ReadinessViolation {
		t.Fatalf("unexpected result: %#v", result)
	}
	if len(result.SplitTriggers) != 4 {
		t.Fatalf("expected 4 split triggers, got %v", result.SplitTriggers)
	}
}

func TestEvaluateComplexityHardSplitAllowedWhileDecomposed(t *testing.T) {
	manifest := validManifestV2()
	manifest.Complexity = ComplexityEstimateV2{
		PrimaryComponents: 3, EstimatedFilesChanged: 20, EstimatedNetLOC: 900,
		ProductionFiles: 15, NetProductionLOC: 850, BoundedContexts: 2,
		IndependentMigrations: 2, DistinctSecurityModels: 2, SeparatelyTestableBehaviors: 2,
	}
	result := evaluateComplexity(manifest, "DECOMPOSED")
	if result.Decision != "split-required" || result.ReadinessViolation {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestValidateComplexityReadiness(t *testing.T) {
	summary := ComplexityAuditSummary{
		ReadinessViolations: 1,
		Results:             []ComplexityResult{{ID: "E2-X-001", ReadinessViolation: true}},
	}
	if err := validateComplexityReadiness(summary); err == nil {
		t.Fatal("expected readiness violation")
	}
}
