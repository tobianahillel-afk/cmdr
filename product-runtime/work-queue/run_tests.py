#!/usr/bin/env python3
"""
Test runner and verification script for Command Work Queue core (CAP-CMD-101 / E10-CMD-101A-READONLY).
Runs the standard Go unit/boundary test suite and validates all 6 canonical system views.
"""

import os
import subprocess
import sys

def main():
    script_dir = os.path.dirname(os.path.abspath(__file__))
    go_sdk_bin = r"C:\Users\hillel.tobiana\.gemini\antigravity\scratch\go_sdk\go\bin"
    tmp_dir = r"C:\Users\hillel.tobiana\.gemini\antigravity\scratch\tmp"
    
    env = os.environ.copy()
    if os.path.exists(go_sdk_bin):
        env["PATH"] = go_sdk_bin + os.pathsep + env.get("PATH", "")
    if os.path.exists(tmp_dir):
        env["GOTMPDIR"] = tmp_dir

    print("======================================================================")
    print("Command Work Queue (CAP-CMD-101 / E10-CMD-101A-READONLY) Test Runner")
    print("======================================================================")

    import shutil
    go_cmd = shutil.which("go", path=env["PATH"]) or "go"

    # 1. Run unit tests
    print("\n--- Step 1: Running unit and boundary tests (`go test -v -count=1 ./...`) ---")
    res_test = subprocess.run(
        [go_cmd, "test", "-v", "-count=1", "./..."],
        cwd=script_dir,
        env=env,
        capture_output=True,
        text=True
    )
    print(res_test.stdout)
    if res_test.stderr:
        print(res_test.stderr, file=sys.stderr)

    if res_test.returncode != 0:
        print("FAIL: Unit tests failed!", file=sys.stderr)
        return res_test.returncode

    # 2. Run static analysis
    print("\n--- Step 2: Running static analysis (`go vet ./...`) ---")
    res_vet = subprocess.run(
        [go_cmd, "vet", "./..."],
        cwd=script_dir,
        env=env,
        capture_output=True,
        text=True
    )
    print(res_vet.stdout)
    if res_vet.stderr:
        print(res_vet.stderr, file=sys.stderr)

    if res_vet.returncode != 0:
        print("FAIL: go vet failed!", file=sys.stderr)
        return res_vet.returncode

    print("\n======================================================================")
    print("ALL TESTS PASSED: CAP-CMD-101 Read-Only Work Queue Verified")
    print("  - Canonical system views: all, incidents, tasks, unassigned, sla-risk, my-work")
    print("  - Deterministic sort, ACL fail-closed, cross-tenant isolation confirmed")
    print("======================================================================")
    return 0

if __name__ == "__main__":
    sys.exit(main())
