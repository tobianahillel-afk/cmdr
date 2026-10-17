//go:build linux || darwin

package main

import (
	"os/exec"
	"syscall"
)

// A process group ensures cancellation reaches children spawned by go build.
func configureTimelineProcessTree(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return syscall.ESRCH
		}
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
}
