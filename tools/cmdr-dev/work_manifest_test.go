package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func validManifestV2() WorkManifestV2 {
	return WorkManifestV2{
		SchemaVersion: 2,
		ID: "E1-WORK-001A",
		Title: "Validate manifests",
		Type: "sublot",
		Parent: "E1-WORK-001",
		Goal: "Validate future work manifests.",
		ProductRefs: ProductRefsV2{},
		DependsOn: []string{"E0-COVERAGE-001C"},
		AllowedPaths: []string{"tools/cmdr-dev/**"},
		ForbiddenPaths: []string{"cmdr-product-spec/**"},
		Inputs: []string{"work manifests"},
		Outputs: []string{"validation result"},
		Invariants: []string{"product spec remains read only"},
		SecurityProperties: []string{"unknown fields rejected"},
		AcceptanceTests: []string{"valid v2 passes"},
		PerformanceBudget: PerformanceBudgetV2{Applicable:false,Rationale:"control-plane validation only",Metrics:[]string{}},
		Migration: MigrationPlanV2{Required:false,Rationale:"no persistent runtime data",Steps:[]string{}},
		Rollback: RollbackPlanV2{Strategy:"revert the work-engine commit",Verification:[]string{"manifest validation returns to prior behavior"}},
		RequiredChecks: []string{"go test ./..."},
		DefinitionOfReady: []string{"dependency verified"},
		DefinitionOfDone: []string{"CI validates v2 manifests"},
		Complexity: ComplexityEstimateV2{PrimaryComponents:1,EstimatedFilesChanged:4,EstimatedNetLOC:350,ProductionFiles:0,NetProductionLOC:0,BoundedContexts:1,IndependentMigrations:0,DistinctSecurityModels:0,SeparatelyTestableBehaviors:1},
	}
}

func TestDecodeWorkManifestV2RejectsUnknownFields(t *testing.T) {
	_, err := decodeWorkManifestV2Bytes([]byte(`{"schema_version":2,"id":"E1-WORK-001A","unknown":true}`))
	if err == nil || !strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("expected unknown-field error, got %v", err)
	}
}

func TestValidateWorkManifestV2(t *testing.T) {
	root:=t.TempDir()
	manifest:=validManifestV2()
	known:=map[string]string{"E1-WORK-001":"parent","E0-COVERAGE-001C":"dependency","E1-WORK-001A":"self"}
	path:=filepath.Join(root,"work","lots","E1-WORK-001A","manifest.json")
	if err:=os.MkdirAll(filepath.Dir(path),0o755); err!=nil { t.Fatal(err) }
	if err:=validateWorkManifestV2(root,path,manifest,known); err!=nil { t.Fatal(err) }
}

func TestValidateWorkManifestV2RejectsPathConflict(t *testing.T) {
	root:=t.TempDir()
	manifest:=validManifestV2()
	manifest.AllowedPaths=[]string{"cmdr-product-spec/06-command/**"}
	known:=map[string]string{"E1-WORK-001":"parent","E0-COVERAGE-001C":"dependency","E1-WORK-001A":"self"}
	path:=filepath.Join(root,"work","lots","E1-WORK-001A","manifest.json")
	if err:=validateWorkManifestV2(root,path,manifest,known); err==nil {
		t.Fatal("expected allowed/forbidden conflict")
	}
}
