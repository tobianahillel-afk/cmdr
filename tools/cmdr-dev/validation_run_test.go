package main

import "testing"

func TestValidationExecutorModesCoverKnownExecutorKeys(t *testing.T) {
	for key := range knownExecutorKeys {
		if _, err := validationExecutorMode(key); err != nil {
			t.Fatalf("known executor key %s is not covered: %v", key, err)
		}
	}
}

func TestValidationExecutorRejectsUnknownKey(t *testing.T) {
	if _, err := validationExecutorMode("shell:arbitrary"); err == nil {
		t.Fatal("expected arbitrary executor rejection")
	}
}

func TestValidationExecutorPreflightSetIsExplicit(t *testing.T) {
	for _, key := range []string{"git-changes", "impact", "validation-plan", "validation-run"} {
		mode, err := validationExecutorMode(key)
		if err != nil {
			t.Fatal(err)
		}
		if mode != "preflight" {
			t.Fatalf("%s should be preflight, got %s", key, mode)
		}
	}
	for _, key := range []string{"gofmt", "go-vet", "go-unit", "architecture-audit"} {
		mode, err := validationExecutorMode(key)
		if err != nil {
			t.Fatal(err)
		}
		if mode != "execute" {
			t.Fatalf("%s should execute, got %s", key, mode)
		}
	}
}
