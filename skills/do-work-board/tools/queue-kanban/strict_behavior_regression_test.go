package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestLegacyStrictEntryPointsRejectZeroProbes(t *testing.T) {
	for _, lane := range []struct {
		name       string
		entryPoint string
		diagnostic string
	}{
		{"JavaScript", "TestMaintainerStrictJavaScriptBehaviorLane", strictJavaScriptBehaviorDiagnostic},
		{"browser", "TestMaintainerStrictBrowserBehaviorLane", strictBrowserBehaviorDiagnostic},
	} {
		t.Run(lane.name, func(t *testing.T) {
			command := exec.Command(os.Args[0], "-test.run=^"+lane.entryPoint+"$", "-test.count=1", "-test.timeout=30s")
			// The wrapper must enable its own strict marker. No installed runtime
			// can be discovered, so all selected probes must leave its count zero.
			command.Env = testEnvironmentWithOverrides(os.Environ(),
				"PATH="+t.TempDir(), browserProbeBinaryOverride+"=",
				javaScriptBehaviorProbeMode+"=off", browserBehaviorProbeMode+"=off",
				strictJavaScriptBehaviorMarker+"=", strictBrowserBehaviorMarker+"=",
			)
			output, err := command.CombinedOutput()
			if err == nil || !strings.Contains(string(output), lane.diagnostic) {
				t.Fatalf("legacy %s entry point did not reject zero probes: error=%v\n%s", lane.name, err, output)
			}
		})
	}
}

func TestLegacyStrictEntryPointsExecuteProbes(t *testing.T) {
	for _, lane := range []struct {
		name       string
		entryPoint string
		probeMode  string
		probeName  string
	}{
		{"JavaScript", "TestMaintainerStrictJavaScriptBehaviorLane", javaScriptBehaviorProbeMode, "TestJavaScriptBehavior"},
		{"browser", "TestMaintainerStrictBrowserBehaviorLane", browserBehaviorProbeMode, "TestBrowserBehavior"},
	} {
		t.Run(lane.name, func(t *testing.T) {
			if os.Getenv(lane.probeMode) != "on" {
				t.Skip("real probe execution is heavy-only")
			}
			command := exec.Command(os.Args[0], "-test.run=^"+lane.entryPoint+"$", "-test.count=1", "-test.timeout=180s", "-test.v")
			// Start disabled even in a heavy parent: the legacy entry point must
			// enable its child rather than relying on the caller's environment.
			command.Env = testEnvironmentWithOverrides(os.Environ(),
				javaScriptBehaviorProbeMode+"=off", browserBehaviorProbeMode+"=off",
				strictJavaScriptBehaviorMarker+"=", strictBrowserBehaviorMarker+"=",
			)
			output, err := command.CombinedOutput()
			if err != nil || !strings.Contains(string(output), "--- PASS: "+lane.probeName) {
				t.Fatalf("legacy %s entry point did not execute successful probes: error=%v\n%s", lane.name, err, output)
			}
		})
	}
}
