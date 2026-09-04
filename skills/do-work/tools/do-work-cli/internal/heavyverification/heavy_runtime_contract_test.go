//go:build unix

package heavyverification

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func shippedRuntimeFiles(t *testing.T) (laneManifest, []byte) {
	t.Helper()
	repositoryRoot, err := filepath.Abs("../../../../../..")
	if err != nil {
		t.Fatal(err)
	}
	manifestBytes, err := os.ReadFile(filepath.Join(repositoryRoot, "_dev/tests/heavy-lanes.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest laneManifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		t.Fatal(err)
	}
	helper, err := os.ReadFile(filepath.Join(repositoryRoot, "_dev/tests/heavy-runtime-fingerprint.py"))
	if err != nil {
		t.Fatal(err)
	}
	return manifest, helper
}

func TestShippedLanesCoverRuntimeHelpersAndRefuseOpaqueBrowser(t *testing.T) {
	manifest, _ := shippedRuntimeFiles(t)
	reusableCount := 0
	for _, lane := range manifest.Lanes {
		for _, path := range []string{"_dev/tests/run-go-tests-with-budget.sh", "_dev/tests/fixture-repo.sh", "_dev/tests/test-duration-log.sh", "_dev/tests/heavy-runtime-fingerprint.py"} {
			if !laneCoversPath(lane.Coverage, path) {
				t.Errorf("%s misses executed helper %s", lane.ID, path)
			}
		}
		if lane.ID == "queue-kanban-browser" && lane.Fingerprint != nil {
			t.Fatal("browser runtime closure is not determinable; lane must execute")
		}
		if lane.Fingerprint != nil {
			reusableCount++
		}
	}
	if reusableCount == 0 {
		t.Fatal("complete lanes must retain evidence reuse")
	}
}

func newRuntimeEvidenceRepository(t *testing.T) string {
	t.Helper()
	root := newLaneEvidenceRepository(t)
	manifest, helper := shippedRuntimeFiles(t)
	var fixtureManifest laneManifest
	if err := json.Unmarshal([]byte(laneEvidenceTestManifest), &fixtureManifest); err != nil {
		t.Fatal(err)
	}
	// Use the shipped runtime discovery at the production runner seam, while
	// the cheap lane marker proves whether an execution actually happened.
	fixtureManifest.Lanes[0].Fingerprint = manifest.Lanes[2].Fingerprint
	for index, argument := range fixtureManifest.Lanes[0].Fingerprint.ToolchainProbes[0] {
		if argument == "skills/do-work/tools/do-work-cli" {
			fixtureManifest.Lanes[0].Fingerprint.ToolchainProbes[0][index] = "."
		}
	}
	encoded, err := json.Marshal(fixtureManifest)
	if err != nil {
		t.Fatal(err)
	}
	writeHeavyTestFile(t, root, "heavy-lanes.json", string(encoded))
	writeHeavyTestFile(t, root, "_dev/tests/heavy-runtime-fingerprint.py", string(helper))
	writeHeavyTestFile(t, root, "go.mod", "module example.test/heavyfixture\n\ngo 1.26.1\n")
	commitHeavyTestChanges(t, root, "shipped runtime fingerprint probe")
	return root
}

func TestShippedRuntimeEvidenceTracksEffectiveGoSettingsAndBinaryBytes(t *testing.T) {
	for _, setting := range []string{"GOOS", "GOARCH", "CGO_ENABLED", "GOENV", "DO_WORK_TEST_DO_WORK_CLI_BINARY"} {
		t.Run(setting, func(t *testing.T) {
			root := newRuntimeEvidenceRepository(t)
			now := time.Now()
			if setting == "GOENV" {
				envPath := filepath.Join(t.TempDir(), "go-env")
				if err := os.WriteFile(envPath, []byte("CGO_ENABLED=0\n"), 0600); err != nil {
					t.Fatal(err)
				}
				t.Setenv("GOENV", envPath)
			}
			if setting == "DO_WORK_TEST_DO_WORK_CLI_BINARY" {
				binaryBytes, err := os.ReadFile("/bin/sh")
				if err != nil {
					t.Fatal(err)
				}
				binaryPath := filepath.Join(t.TempDir(), "supplied-cli")
				if err := os.WriteFile(binaryPath, binaryBytes, 0700); err != nil {
					t.Fatal(err)
				}
				t.Setenv(setting, binaryPath)
			}
			first := verifyLanes(t, root, now, "alpha-lane")
			if laneExecutionRecord(t, first, "alpha-lane").FingerprintSHA256 == "" {
				t.Fatal("default runtime must have a determinable fingerprint")
			}
			second := verifyLanes(t, root, now.Add(time.Minute), "alpha-lane")
			assertLaneDisposition(t, laneExecutionRecord(t, second, "alpha-lane"), LaneDispositionReused, laneReasonFingerprintMatch)
			switch setting {
			case "GOOS":
				t.Setenv(setting, "freebsd")
			case "GOARCH":
				t.Setenv(setting, "386")
			case "CGO_ENABLED":
				t.Setenv(setting, "0")
			case "GOENV":
				if err := os.WriteFile(os.Getenv(setting), []byte("CGO_ENABLED=1\n"), 0600); err != nil {
					t.Fatal(err)
				}
			case "DO_WORK_TEST_DO_WORK_CLI_BINARY":
				file, err := os.OpenFile(os.Getenv(setting), os.O_APPEND|os.O_WRONLY, 0)
				if err != nil {
					t.Fatal(err)
				}
				if _, err := file.WriteString("changed binary bytes"); err != nil {
					t.Fatal(err)
				}
				if err := file.Close(); err != nil {
					t.Fatal(err)
				}
			}
			third := verifyLanes(t, root, now.Add(2*time.Minute), "alpha-lane")
			if laneExecutionRecord(t, third, "alpha-lane").Disposition != LaneDispositionExecuted {
				t.Fatal("changed effective tool/runtime input reused old success")
			}
		})
	}
}

func TestShippedRuntimeRefusesOpaqueOverrides(t *testing.T) {
	root := newRuntimeEvidenceRepository(t)
	now := time.Now()
	verifyLanes(t, root, now, "alpha-lane")
	t.Setenv("GOFLAGS", "-tags=alternate")
	result := verifyLanes(t, root, now.Add(time.Minute), "alpha-lane")
	assertLaneDisposition(t, laneExecutionRecord(t, result, "alpha-lane"), LaneDispositionExecuted, laneReasonFingerprintUncertain)
	t.Setenv("GOFLAGS", "")
	suppliedPath := filepath.Join(t.TempDir(), "cli-wrapper")
	if err := os.WriteFile(suppliedPath, []byte("#!/bin/sh\nexec /bin/true\n"), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DO_WORK_TEST_DO_WORK_CLI_BINARY", suppliedPath)
	result = verifyLanes(t, root, now.Add(2*time.Minute), "alpha-lane")
	assertLaneDisposition(t, laneExecutionRecord(t, result, "alpha-lane"), LaneDispositionExecuted, laneReasonFingerprintUncertain)
}

func TestShippedGitIsolationPreservesGenericLaneInheritance(t *testing.T) {
	root := newRuntimeEvidenceRepository(t)
	configuration := filepath.Join(t.TempDir(), "git-config")
	if err := os.WriteFile(configuration, []byte("[commit]\n gpgsign = true\n"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GIT_CONFIG_GLOBAL", configuration)
	argv := []string{"git", "config", "--get", "commit.gpgsign"}
	// An arbitrary manifest continues to receive the caller's configuration.
	if output, err := runFingerprintProbe(root, argv, time.Second); err != nil || string(output) != "true\n" {
		t.Fatalf("generic probe inheritance: %q %v", output, err)
	}
	manifest, _ := shippedRuntimeFiles(t)
	for _, lane := range manifest.Lanes {
		isolatedArgv := append(append([]string(nil), lane.Argv[:3]...), argv...)
		if output, err := runFingerprintProbe(root, isolatedArgv, time.Second); err == nil || len(output) != 0 {
			t.Fatalf("%s inherited global signing: %q %v", lane.ID, output, err)
		}
	}
	// The matching shipped probe stays complete with that host configuration.
	result := verifyLanes(t, root, time.Now(), "alpha-lane")
	if laneExecutionRecord(t, result, "alpha-lane").FingerprintSHA256 == "" {
		t.Fatal("shipped runtime probe did not isolate host Git configuration")
	}
}
