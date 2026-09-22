package main

import "testing"

func testSecurityRegistry() SecurityGateRegistry {
	return SecurityGateRegistry{
		SchemaVersion:   1,
		DefaultPolicy:   "deny-unregistered-security-control",
		ReadinessStates: []string{"specified", "active", "deferred-runtime"},
		ExecutionStages: []string{"pr", "nightly", "release", "on-demand"},
		DataPolicies:    []string{"source-local-only", "artifact-local-only", "synthetic-runtime-only"},
		EvidenceKinds:   []string{"process-exit", "finding-report", "sbom", "test-result", "fuzz-report", "dast-report", "provenance"},
		Gates: []SecurityGate{
			{ID: "SEC-SECRETS-001", Title: "Secrets", ControlFamily: "secrets", Readiness: "specified", Stages: []string{"pr"}, BlockingPolicy: "blocking", RiskDomains: []string{"security"}, TriggerPaths: []string{"**"}, DataPolicy: "source-local-only", EvidenceKinds: []string{"finding-report"}, ActivationCondition: "scanner required"},
			{ID: "SEC-SAST-001", Title: "SAST", ControlFamily: "sast", Readiness: "specified", Stages: []string{"pr"}, BlockingPolicy: "conditional-blocking", RiskDomains: []string{"security"}, TriggerPaths: []string{"**"}, DataPolicy: "source-local-only", EvidenceKinds: []string{"finding-report"}, ActivationCondition: "analyzer required"},
			{ID: "SEC-SCA-001", Title: "SCA", ControlFamily: "sca", Readiness: "specified", Stages: []string{"pr"}, BlockingPolicy: "blocking", RiskDomains: []string{"security"}, TriggerPaths: []string{"**/go.mod"}, DataPolicy: "source-local-only", EvidenceKinds: []string{"finding-report"}, ActivationCondition: "advisories required"},
			{ID: "SEC-SBOM-001", Title: "SBOM", ControlFamily: "sbom", Readiness: "specified", Stages: []string{"release"}, BlockingPolicy: "blocking", RiskDomains: []string{"security"}, TriggerPaths: []string{"**/go.mod"}, DataPolicy: "artifact-local-only", EvidenceKinds: []string{"sbom"}, ActivationCondition: "artifact required"},
			{ID: "SEC-IAC-001", Title: "IaC", ControlFamily: "iac", Readiness: "specified", Stages: []string{"pr"}, BlockingPolicy: "conditional-blocking", RiskDomains: []string{"security"}, TriggerPaths: []string{"infra/**"}, DataPolicy: "source-local-only", EvidenceKinds: []string{"finding-report"}, ActivationCondition: "IaC required"},
			{ID: "SEC-CONTAINER-001", Title: "Container", ControlFamily: "container", Readiness: "deferred-runtime", Stages: []string{"release"}, BlockingPolicy: "blocking", RiskDomains: []string{"security"}, TriggerPaths: []string{"**/Dockerfile"}, DataPolicy: "artifact-local-only", EvidenceKinds: []string{"finding-report"}, ActivationCondition: "container required"},
			{ID: "SEC-AUTH-NEG-001", Title: "Auth", ControlFamily: "authorization-negative", Readiness: "deferred-runtime", Stages: []string{"pr"}, BlockingPolicy: "blocking", RiskDomains: []string{"security"}, TriggerPaths: []string{"**"}, DataPolicy: "synthetic-runtime-only", EvidenceKinds: []string{"test-result"}, ActivationCondition: "auth runtime required"},
			{ID: "SEC-TENANT-ISO-001", Title: "Tenant", ControlFamily: "tenant-isolation", Readiness: "deferred-runtime", Stages: []string{"pr"}, BlockingPolicy: "blocking", RiskDomains: []string{"security"}, TriggerPaths: []string{"**"}, DataPolicy: "synthetic-runtime-only", EvidenceKinds: []string{"test-result"}, ActivationCondition: "tenant runtime required"},
			{ID: "SEC-FUZZ-001", Title: "Fuzz", ControlFamily: "fuzz", Readiness: "deferred-runtime", Stages: []string{"nightly"}, BlockingPolicy: "conditional-blocking", RiskDomains: []string{"security"}, TriggerPaths: []string{"**"}, DataPolicy: "synthetic-runtime-only", EvidenceKinds: []string{"fuzz-report"}, ActivationCondition: "target required"},
			{ID: "SEC-DAST-001", Title: "DAST", ControlFamily: "dast", Readiness: "deferred-runtime", Stages: []string{"nightly"}, BlockingPolicy: "conditional-blocking", RiskDomains: []string{"security"}, TriggerPaths: []string{"**"}, DataPolicy: "synthetic-runtime-only", EvidenceKinds: []string{"dast-report"}, ActivationCondition: "runtime required"},
			{ID: "SEC-RACE-001", Title: "Race", ControlFamily: "race-concurrency", Readiness: "deferred-runtime", Stages: []string{"nightly"}, BlockingPolicy: "conditional-blocking", RiskDomains: []string{"security"}, TriggerPaths: []string{"**"}, DataPolicy: "synthetic-runtime-only", EvidenceKinds: []string{"test-result"}, ActivationCondition: "concurrency required"},
			{ID: "SEC-SUPPLY-001", Title: "Supply", ControlFamily: "supply-chain", Readiness: "specified", Stages: []string{"release"}, BlockingPolicy: "blocking", RiskDomains: []string{"security"}, TriggerPaths: []string{"**"}, DataPolicy: "artifact-local-only", EvidenceKinds: []string{"provenance"}, ActivationCondition: "artifact required"},
		},
	}
}

func TestValidateSecurityGateRegistry(t *testing.T) {
	if err := validateSecurityGateRegistry(testSecurityRegistry()); err != nil {
		t.Fatal(err)
	}
}

func TestSecurityRegistryRejectsMissingFamily(t *testing.T) {
	registry := testSecurityRegistry()
	registry.Gates = registry.Gates[:len(registry.Gates)-1]
	if err := validateSecurityGateRegistry(registry); err == nil {
		t.Fatal("expected missing family rejection")
	}
}

func TestSecurityRegistryRejectsActiveGateWithoutImplementation(t *testing.T) {
	registry := testSecurityRegistry()
	registry.Gates[0].Readiness = "active"
	if err := validateSecurityGateRegistry(registry); err == nil {
		t.Fatal("expected active gate implementation rejection")
	}
}

func TestSecurityRegistryRejectsNonSecurityDomainGate(t *testing.T) {
	registry := testSecurityRegistry()
	registry.Gates[0].RiskDomains = []string{"code-quality"}
	if err := validateSecurityGateRegistry(registry); err == nil {
		t.Fatal("expected missing security risk-domain rejection")
	}
}

func TestSecurityRegistryRejectsUnknownEvidenceKind(t *testing.T) {
	registry := testSecurityRegistry()
	registry.Gates[0].EvidenceKinds = []string{"magic-report"}
	if err := validateSecurityGateRegistry(registry); err == nil {
		t.Fatal("expected unknown evidence rejection")
	}
}
