package harness

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	rootOnce sync.Once
	repoRoot string
	rootErr  error
)

// RepoRoot returns the absolute path to the cmdr repository root.
func RepoRoot() (string, error) {
	rootOnce.Do(func() {
		dir, err := os.Getwd()
		if err != nil {
			rootErr = fmt.Errorf("failed to get current working directory: %w", err)
			return
		}
		for {
			if _, err := os.Stat(filepath.Join(dir, "cmdr-product-spec")); err == nil {
				repoRoot = filepath.Clean(dir)
				return
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				rootErr = fmt.Errorf("could not find cmdr repository root from working directory %s", dir)
				return
			}
			dir = parent
		}
	})
	return repoRoot, rootErr
}

// CmdResult captures the outcome of a subprocess invocation.
type CmdResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Duration time.Duration
	Err      error
}

// RunCmd executes a command with default timeout and retries on transient Windows file locks.
func RunCmd(ctx context.Context, dir string, env []string, name string, args ...string) CmdResult {
	start := time.Now()
	var stdout, stderr bytes.Buffer

	cmd := exec.CommandContext(ctx, name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	duration := time.Since(start)

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}

	return CmdResult{
		ExitCode: exitCode,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Duration: duration,
		Err:      err,
	}
}

// RunCmdrDev invokes a cmdr-dev CLI subcommand from repo root.
// It automatically sets --root <repoRoot> and passes extra arguments.
// If a precompiled cmdr-dev executable exists (e.g. in .cmdr-tmp, tools/cmdr-dev, or repo root),
// it is invoked directly to avoid Windows Defender ASR "Access is denied" on go run.
func RunCmdrDev(ctx context.Context, subcommand string, extraArgs ...string) CmdResult {
	root, err := RepoRoot()
	if err != nil {
		return CmdResult{
			ExitCode: -1,
			Err:      err,
		}
	}

	tmpDir := filepath.Join(root, ".cmdr-tmp")
	_ = os.MkdirAll(tmpDir, 0755)

	toolsDir := filepath.Join(root, "tools", "cmdr-dev")

	// Check for precompiled binary candidates
	candidates := []string{
		filepath.Join(tmpDir, "cmdr-dev.exe"),
		filepath.Join(tmpDir, "cmdr-dev"),
		filepath.Join(toolsDir, "cmdr-dev.exe"),
		filepath.Join(toolsDir, "cmdr-dev"),
		filepath.Join(root, "cmdr-dev.exe"),
		filepath.Join(root, "cmdr-dev"),
	}

	var precompiled string
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && !fi.IsDir() {
			precompiled = c
			break
		}
	}

	var cmdName string
	var args []string

	if precompiled != "" {
		cmdName = precompiled
		args = []string{subcommand, "--root", root}
		args = append(args, extraArgs...)
	} else {
		cmdName = "go"
		args = []string{"run", "-C", toolsDir, ".", subcommand, "--root", root}
		args = append(args, extraArgs...)
	}

	customEnv := []string{
		"GOTMPDIR=" + tmpDir,
	}

	// Retry on transient Windows Defender / process locking
	maxRetries := 3
	var res CmdResult
	for attempt := 1; attempt <= maxRetries; attempt++ {
		res = RunCmd(ctx, root, customEnv, cmdName, args...)
		if res.ExitCode == 0 || res.ExitCode == 1 || res.ExitCode == 2 {
			// Normal program exit (0 = success, 1 = failure, 2 = usage / unverified)
			break
		}
		// If transient fork/exec access denied, back off and retry
		if strings.Contains(res.Stderr, "Access is denied") || strings.Contains(res.Stderr, "Accès refusé") {
			time.Sleep(time.Duration(attempt*500) * time.Millisecond)
			continue
		}
		break
	}

	return res
}

// ParseJSON parses a JSON string or bytes into a target struct.
func ParseJSON(data string, v interface{}) error {
	trimmed := strings.TrimSpace(data)
	if trimmed == "" {
		return fmt.Errorf("empty JSON payload")
	}
	decoder := json.NewDecoder(strings.NewReader(trimmed))
	return decoder.Decode(v)
}

// TempDir creates an isolated temporary directory under .cmdr-tmp/ for a test.
func TempDir(t *testing.T, prefix string) string {
	root, err := RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}
	base := filepath.Join(root, ".cmdr-tmp")
	if err := os.MkdirAll(base, 0755); err != nil {
		t.Fatalf("failed to create base temp dir: %v", err)
	}
	dir, err := os.MkdirTemp(base, prefix+"-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})
	return dir
}
