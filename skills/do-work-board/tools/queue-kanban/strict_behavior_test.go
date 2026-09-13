package main

import (
	"flag"
	"os"
	"os/exec"
	"testing"
	"time"
)

// These exact names are compatibility entry points used by consumer gate scripts.
// Ordinary package runs keep heavy probes opt-in, as they were before the wrappers.
func TestMaintainerStrictJavaScriptBehaviorLane(t *testing.T) {
	runLegacyStrictBehaviorLane(t, "^TestJavaScriptBehavior", javaScriptBehaviorProbeMode, strictJavaScriptBehaviorMarker)
}

func TestMaintainerStrictBrowserBehaviorLane(t *testing.T) {
	runLegacyStrictBehaviorLane(t, "^TestBrowserBehavior", browserBehaviorProbeMode, strictBrowserBehaviorMarker)
}

func runLegacyStrictBehaviorLane(t *testing.T, probePattern, probeMode, strictMarker string) {
	t.Helper()
	selection := flag.Lookup("test.run")
	if selection == nil || selection.Value.String() != "^"+t.Name()+"$" {
		t.Skip("legacy strict behavior lane runs only when selected directly")
	}
	probeTimeout := 180 * time.Second
	if deadline, ok := t.Deadline(); ok {
		probeTimeout = time.Until(deadline) - time.Second
		if probeTimeout <= 0 {
			t.Fatal("no time remains to run strict behavior probes")
		}
	}
	// Select only real probes, excluding both wrappers and their regression tests.
	// TestMain in this child remains the authority that rejects zero probes.
	command := exec.Command(os.Args[0], "-test.run="+probePattern, "-test.count=1", "-test.timeout="+probeTimeout.String(), "-test.v")
	command.Env = testEnvironmentWithOverrides(os.Environ(),
		javaScriptBehaviorProbeMode+"=off", browserBehaviorProbeMode+"=off",
		strictJavaScriptBehaviorMarker+"=", strictBrowserBehaviorMarker+"=",
	)
	command.Env = testEnvironmentWithOverrides(command.Env, probeMode+"=on", strictMarker+"=1")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("strict behavior lane failed: %v\n%s", err, output)
	}
	t.Logf("%s", output)
}
