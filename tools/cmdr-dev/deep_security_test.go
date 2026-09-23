package main

import "testing"

func deepTestRegistry(readiness string) SecurityGateRegistry {
	key := ""
	if readiness == "active" {
		key = "test-implementation"
	}
	return SecurityGateRegistry{Gates: []SecurityGate{
		{ID: "SEC-FUZZ-001", Readiness: readiness, ImplementationKey: key, Stages: []string{"pr", "nightly", "on-demand"}, EvidenceKinds: []string{"fuzz-report"}},
		{ID: "SEC-DAST-001", Readiness: readiness, ImplementationKey: key, Stages: []string{"pr", "nightly", "release", "on-demand"}, EvidenceKinds: []string{"dast-report"}},
		{ID: "SEC-RACE-001", Readiness: readiness, ImplementationKey: key, Stages: []string{"pr", "nightly", "release", "on-demand"}, EvidenceKinds: []string{"test-result"}},
	}}
}

func baseDeepPolicy(targets ...DeepSecurityTarget) DeepSecurityPolicy {
	return DeepSecurityPolicy{
		SchemaVersion: 1, DefaultPolicy: "deny-unregistered-deep-security-target",
		EvidenceDirectory: "engineering/testing/deep-security-evidence",
		ControlPlanePaths: []string{"engineering/security/**", "engineering/architecture/**", ".github/workflows/**"},
		Gates: []DeepSecurityGatePolicy{
			{GateID: "SEC-FUZZ-001", TargetKinds: []string{"fuzz"}, PRMode: "security-sensitive-only"},
			{GateID: "SEC-DAST-001", TargetKinds: []string{"dast"}, PRMode: "security-sensitive-only"},
			{GateID: "SEC-RACE-001", TargetKinds: []string{"race", "sanitizer"}, PRMode: "security-sensitive-only"},
		},
		Targets: targets,
	}
}

func deepRuntimeArchitecture() ArchitectureRegistry {
	return ArchitectureRegistry{
		SchemaVersion: 1, DefaultPolicy: "deny-unregistered-path",
		Boundaries: []ArchitectureBoundary{
			{ID: "product-spec", Kind: "product-documentation", Roots: []string{"cmdr-product-spec/**"}},
			{ID: "api", Kind: "product-runtime", Roots: []string{"services/api/**"}},
		},
	}
}

func deepRuntimePolicy() RuntimeSecurityPolicy {
	return baseRuntimePolicy(RuntimeSecurityScope{
		BoundaryID: "api", Owner: "runtime-security",
		SecurityCriticalPaths: []string{"services/api/security/**"},
		AuthorizationRationale: "synthetic test", TenantIsolationRationale: "synthetic test",
	})
}

func TestDeepSecurityPRWithoutSensitiveChangeDoesNotSchedule(t *testing.T) {
	summary, err := evaluateDeepSecurityPolicy(
		"", "pr", "", []string{"docs/readme.md"},
		"0123456789abcdef0123456789abcdef01234567",
		securityTestArchitecture(), deepTestRegistry("deferred-runtime"),
		baseRuntimePolicy(), baseDeepPolicy(),
	)
	if err != nil {
		t.Fatal(err)
	}
	if summary.SecuritySensitiveChange || summary.SelectedTargets != 0 || summary.DeferredGates != 0 {
		t.Fatalf("unexpected PR plan: %#v", summary)
	}
}

func TestDeepSecuritySensitivePRDefersUnavailableRuntimeGates(t *testing.T) {
	summary, err := evaluateDeepSecurityPolicy(
		"", "pr", "", []string{"engineering/security/security-gates.json"},
		"0123456789abcdef0123456789abcdef01234567",
		securityTestArchitecture(), deepTestRegistry("deferred-runtime"),
		baseRuntimePolicy(), baseDeepPolicy(),
	)
	if err != nil {
		t.Fatal(err)
	}
	if !summary.SecuritySensitiveChange || summary.SelectedTargets != 0 || summary.DeferredGates != 3 {
		t.Fatalf("expected three deferred gates, got %#v", summary)
	}
	for _, decision := range summary.Decisions {
		if decision.Decision != "deferred" {
			t.Fatalf("expected deferred decision, got %#v", decision)
		}
	}
}

func TestDeepSecurityRegisteredTargetRequiresActiveGate(t *testing.T) {
	target := DeepSecurityTarget{
		ID: "api-fuzz", GateID: "SEC-FUZZ-001", Kind: "fuzz", BoundaryID: "api",
		Owner: "runtime-security", Target: "FuzzAPI", Seed: "corpus-v1",
		Configuration: "fuzz-default-v1", MinimumDurationSec: 60, SyntheticOnly: true,
	}
	policy := baseDeepPolicy(target)
	if err := validateDeepSecurityPolicy(policy, deepTestRegistry("deferred-runtime"), deepRuntimeArchitecture()); err != nil {
		t.Fatal(err)
	}
	if _, err := evaluateDeepSecurityPolicy(
		"", "nightly", "", nil, "0123456789abcdef0123456789abcdef01234567",
		deepRuntimeArchitecture(), deepTestRegistry("deferred-runtime"), deepRuntimePolicy(), policy,
	); err == nil {
		t.Fatal("expected registered target on inactive gate to fail closed")
	}
}

func TestDeepSecurityEvidenceEnforcesDurationCommitAndReleaseArtifact(t *testing.T) {
	target := DeepSecurityTarget{
		ID: "api-dast", GateID: "SEC-DAST-001", Kind: "dast", BoundaryID: "api",
		Owner: "runtime-security", Target: "https://synthetic.invalid", Seed: "dataset-v1",
		Configuration: "dast-v1", MinimumDurationSec: 120, SyntheticOnly: true,
	}
	gate := SecurityGate{ID: "SEC-DAST-001", EvidenceKinds: []string{"dast-report"}}
	commit := "0123456789abcdef0123456789abcdef01234567"
	valid := DeepSecurityEvidence{
		SchemaVersion: 1, GateID: target.GateID, TargetID: target.ID, Stage: "release",
		SourceCommit: commit, Target: target.Target, Seed: target.Seed, Configuration: target.Configuration,
		DurationSec: 120, EvidenceKind: "dast-report", SyntheticData: true, Result: "pass",
		ArtifactSHA256: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}
	if err := validateDeepSecurityEvidence(valid, target, gate, "release", commit); err != nil {
		t.Fatal(err)
	}
	short := valid
	short.DurationSec = 119
	if err := validateDeepSecurityEvidence(short, target, gate, "release", commit); err == nil {
		t.Fatal("expected short duration rejection")
	}
	stale := valid
	stale.SourceCommit = "abcdef0123456789abcdef0123456789abcdef01"
	if err := validateDeepSecurityEvidence(stale, target, gate, "release", commit); err == nil {
		t.Fatal("expected stale commit rejection")
	}
	noArtifact := valid
	noArtifact.ArtifactSHA256 = ""
	if err := validateDeepSecurityEvidence(noArtifact, target, gate, "release", commit); err == nil {
		t.Fatal("expected release artifact binding rejection")
	}
}

func TestDeepSecuritySanitizerUsesRaceGate(t *testing.T) {
	target := DeepSecurityTarget{
		ID: "native-asan", GateID: "SEC-RACE-001", Kind: "sanitizer", BoundaryID: "api",
		Owner: "runtime-security", Target: "native-suite", Seed: "synthetic-v1",
		Configuration: "asan-ubsan-v1", MinimumDurationSec: 90, SyntheticOnly: true,
	}
	if err := validateDeepSecurityPolicy(baseDeepPolicy(target), deepTestRegistry("active"), deepRuntimeArchitecture()); err != nil {
		t.Fatal(err)
	}
}

func TestDeepSecurityRejectsTargetOutsideRuntimeBoundary(t *testing.T) {
	target := DeepSecurityTarget{
		ID: "bad-target", GateID: "SEC-FUZZ-001", Kind: "fuzz", BoundaryID: "missing",
		Owner: "runtime-security", Target: "FuzzBad", Seed: "seed", Configuration: "cfg",
		MinimumDurationSec: 60, SyntheticOnly: true,
	}
	if err := validateDeepSecurityPolicy(baseDeepPolicy(target), deepTestRegistry("active"), deepRuntimeArchitecture()); err == nil {
		t.Fatal("expected unknown runtime boundary rejection")
	}
}
