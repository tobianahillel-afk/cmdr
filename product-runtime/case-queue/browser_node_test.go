package casequeue

import (
 "os/exec"
 "testing"
)

// Native UI uses no new JS dependency. This optional behavioral test runs
// with the host Node runtime when available; Go HTTP tests always run.
func TestNativeQueueBrowserClientBehavior(t *testing.T) {
 node, err := exec.LookPath("node")
 if err != nil { t.Skip("Node is unavailable on this host") }
 cmd := exec.Command(node, "--test", "browser.test.cjs")
 output, err := cmd.CombinedOutput()
 if err != nil { t.Fatalf("native Case Queue browser behavior: %v\n%s",err,output) }
}
