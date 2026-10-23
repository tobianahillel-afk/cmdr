package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Approved one-time toolchain transitions for Go 1.26.10 stdlib security fix.
// Identity is derived from immutable old/new Git blob hashes, not the branch
// name or a runtime assertion from the caller. Future toolchain changes must
// be reviewed separately.
var approvedGo12610ModuleBlobs = map[string]struct {
	before, after string
}{
	"product-runtime/event-inspection": {
		"a2b90dca906b971f167a9108b6436bf9580a82e7",
		"9f45d10c3e854581bb835f326a8d3a5572dcea1f",
	},
	"product-runtime/event-search": {
		"aa7adaad1e3d7b57fcc698c842fb965a23c30068",
		"c954ca9b250acb749cda2f3bd6d98983777bd136",
	},
	"product-runtime/hunt-management": {
		"e05aba665edb6051273695bae4f5464f40f93559",
		"29d4f0aa7b40cd0c0e859f70e4d04a2d777a4e31",
	},
	"product-runtime/saved-query-assets": {
		"6c1ebc18e40d413a998916d42194cc0230d33a1a",
		"ac9ab9259a841ce23b1e9fa230341f6bd7c8f104",
	},
	"product-runtime/case-queue": {
		"fc6e5cd9f9be05481cb3774e87b2fb111db168ff",
		"ac66e292a7863ff0b93aeea3b979a5f608320966",
	},
	"product-runtime/context-envelope": {
		"dee4dd95ad9fb62dff955c9254176f180b0bfca2",
		"034cf1e675188c653a1449cbaa01006886a1ec8a",
	},
}

// checkPinnedGoToolchainDelta returns true ONLY for byte-identical historical
// sources. A false,nil result means exactly one reviewed go.mod file changed
// and requires a fresh execution proof; all other changes fail closed.
func checkPinnedGoToolchainDelta(root, evidenceHead, currentHead, runtimeDir string) (bool, error) {
	blobs, approved := approvedGo12610ModuleBlobs[runtimeDir]
	if !approved {
		return false, fmt.Errorf("unregistered Go toolchain evidence scope")
	}
	evidenceHead, err := validateFullCommitID(evidenceHead)
	if err != nil {
		return false, fmt.Errorf("invalid historical toolchain evidence: %w", err)
	}
	currentHead, err = validateFullCommitID(currentHead)
	if err != nil {
		return false, fmt.Errorf("invalid current toolchain revision: %w", err)
	}
	if err := ensureGitAncestor(root, evidenceHead, currentHead); err != nil {
		return false, fmt.Errorf("toolchain evidence ancestry: %w", err)
	}
	diff := exec.Command("git", "-C", root, "diff", "--name-only", evidenceHead, currentHead, "--", runtimeDir) // #nosec G204,G702 -- full validated SHA and fixed registered path.
	out, err := diff.Output()
	if err != nil {
		return false, fmt.Errorf("failed to inspect verified runtime source delta: %w", err)
	}
	deltas := strings.Fields(strings.TrimSpace(string(out)))
	if len(deltas) == 0 {
		return true, nil
	}
	expected := runtimeDir + "/go.mod"
	if len(deltas) != 1 || deltas[0] != expected {
		return false, fmt.Errorf("verified runtime contains unreviewed executable or module changes; full runtime review required")
	}
	for _, item := range []struct{ revision, blob string }{
		{evidenceHead, blobs.before},
		{currentHead, blobs.after},
	} {
		cmd := exec.Command("git", "-C", root, "rev-parse", "--verify", item.revision+":"+expected) // #nosec G204,G702 -- validated SHA and fixed registered module path.
		found, err := cmd.Output()
		if err != nil || strings.TrimSpace(string(found)) != item.blob {
			return false, fmt.Errorf("runtime compiler declaration does not match byte-reviewed Go 1.26.10 security transition")
		}
	}
	return false, nil
}

func reverifyPinnedGoToolchainTransition(root, evidenceHead, currentHead, runtimeDir string) (bool, error) {
	unchanged, err := checkPinnedGoToolchainDelta(root, evidenceHead, currentHead, runtimeDir)
	if err != nil || unchanged {
		return unchanged, err
	}
	// Mandatory global CHK-SCA-GO is separately evaluated. This fresh module
	// proof supplements it; neither historical E2E nor SCA is waived.
	dir, err := resolveRepoPath(root, runtimeDir, false)
	if err != nil {
		return false, err
	}
	for _, args := range [][]string{
		{"test", "-count=1", "./..."},
		{"vet", "./..."},
	} {
		cmd := exec.Command("go", args...) // #nosec G204 -- fixed executable and literal argument alternatives.
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local")
		if _, err := cmd.CombinedOutput(); err != nil {
			return false, fmt.Errorf("patched Go runtime failed fresh full module validation: %w", err)
		}
	}
	return false, nil // Do not assert that original VERIFIED source is unchanged.
}
