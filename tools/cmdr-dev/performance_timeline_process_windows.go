//go:build windows

package main

import (
	"context"
	"os/exec"
	"strconv"
	"time"
)

// Best-effort Windows descendant cancellation via built-in taskkill.
// A Windows CI run is required before making OS-wide verification claims.
func configureTimelineProcessTree(cmd *exec.Cmd) {
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return nil
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		kill := exec.CommandContext(ctx, "taskkill", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid))
		if err := kill.Run(); err != nil {
			return cmd.Process.Kill()
		}
		return nil
	}
}
