package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"testing/fstest"
	"time"
)

const (
	strictJavaScriptBehaviorDiagnostic = "queue-kanban: strict JavaScript behavior lane executed zero probes"
	strictJavaScriptBehaviorMarker     = "QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR"
	javaScriptBehaviorProbeMode        = "QUEUE_KANBAN_JAVASCRIPT_PROBES"
)

var javaScriptBehaviorProbeCount atomic.Int64
var liveIndexOnce sync.Once
var liveIndexHTML string

func TestEmbeddedAuthoredJavaScriptInventory(t *testing.T) {
	webEntries, readError := embeddedWebAssets.ReadDir("web")
	if readError != nil {
		t.Fatalf("read embedded web assets: %v", readError)
	}
	var authoredJavaScriptPaths []string
	for _, webEntry := range webEntries {
		if !webEntry.IsDir() && strings.HasSuffix(webEntry.Name(), ".js") {
			authoredJavaScriptPaths = append(authoredJavaScriptPaths, "web/"+webEntry.Name())
		}
	}
	wantAuthoredJavaScriptPaths := []string{
		"web/board-calendar.js",
		"web/board-cards.js",
		"web/board-clipboard.js",
		"web/board-controls.js",
		"web/board-core.js",
		"web/board-detail.js",
		"web/board-durations.js",
		"web/board-filters.js",
		"web/board-testing.js",
		"web/board-timeline.js",
		"web/board.js",
	}
	if strings.Join(authoredJavaScriptPaths, "\n") != strings.Join(wantAuthoredJavaScriptPaths, "\n") {
		t.Fatalf("embedded authored JavaScript paths = %q, want exact shell-plus-fragment inventory %q",
			authoredJavaScriptPaths, wantAuthoredJavaScriptPaths)
	}
}

func TestBoardJavaScriptAssemblyStructure(t *testing.T) {
	wantFragmentPaths := []string{
		"web/board-core.js",
		"web/board-filters.js",
		"web/board-cards.js",
		"web/board-calendar.js",
		"web/board-durations.js",
		"web/board-timeline.js",
		"web/board-testing.js",
		"web/board-detail.js",
		"web/board-controls.js",
		"web/board-clipboard.js",
	}
	if strings.Join(boardJavaScriptFragmentPaths[:], "\n") != strings.Join(wantFragmentPaths, "\n") {
		t.Fatalf("JavaScript fragment manifest = %q, want literal execution order %q",
			boardJavaScriptFragmentPaths, wantFragmentPaths)
	}

	shellText, shellError := embeddedWebAssets.ReadFile(boardJavaScriptShellPath)
	if shellError != nil {
		t.Fatalf("read JavaScript shell: %v", shellError)
	}
	placeholderBytes := []byte(boardJavaScriptPlaceholder)
	if placeholderCount := bytes.Count(shellText, placeholderBytes); placeholderCount != 1 {
		t.Fatalf("shell placeholder count = %d, want exactly 1", placeholderCount)
	}
	placeholderLineBytes := []byte(boardJavaScriptPlaceholderLine)
	if placeholderLineCount := bytes.Count(shellText, placeholderLineBytes); placeholderLineCount != 1 {
		t.Fatalf("canonical shell placeholder line count = %d, want exactly 1", placeholderLineCount)
	}

	assembledClient, assembleError := assembleBoardJavaScript(embeddedWebAssets)
	if assembleError != nil {
		t.Fatalf("assembleBoardJavaScript: %v", assembleError)
	}
	assembledAgain, secondAssembleError := assembleBoardJavaScript(embeddedWebAssets)
	if secondAssembleError != nil {
		t.Fatalf("second assembleBoardJavaScript: %v", secondAssembleError)
	}
	if !bytes.Equal(assembledClient, assembledAgain) {
		t.Fatal("JavaScript assembly is not deterministic")
	}
	if bytes.Contains(assembledClient, placeholderBytes) {
		t.Fatal("assembled JavaScript retained the private fragment placeholder")
	}

	fragmentTexts := make([][]byte, 0, len(wantFragmentPaths))
	previousFragmentIndex := -1
	for _, fragmentPath := range wantFragmentPaths {
		fragmentText, readError := embeddedWebAssets.ReadFile(fragmentPath)
		if readError != nil {
			t.Fatalf("read %s: %v", fragmentPath, readError)
		}
		if !bytes.HasSuffix(fragmentText, []byte("\n")) ||
			bytes.HasSuffix(fragmentText, []byte("\n\n")) ||
			bytes.HasSuffix(fragmentText, []byte("\r\n")) {
			t.Fatalf("%s must end with exactly one LF", fragmentPath)
		}
		if occurrenceCount := bytes.Count(assembledClient, fragmentText); occurrenceCount != 1 {
			t.Fatalf("%s occurs %d times in assembled client, want exactly 1", fragmentPath, occurrenceCount)
		}
		fragmentIndex := bytes.Index(assembledClient, fragmentText)
		if fragmentIndex <= previousFragmentIndex {
			t.Fatalf("%s starts at byte %d after previous fragment byte %d; manifest order was not preserved",
				fragmentPath, fragmentIndex, previousFragmentIndex)
		}
		previousFragmentIndex = fragmentIndex
		fragmentTexts = append(fragmentTexts, fragmentText)
	}

	wantFragmentBlock := bytes.Join(fragmentTexts, []byte("\n"))
	wantFragmentBlock = append(wantFragmentBlock, '\n')
	wantAssembledClient := bytes.Replace(shellText, placeholderLineBytes, wantFragmentBlock, 1)
	if !bytes.Equal(assembledClient, wantAssembledClient) {
		t.Fatal("assembled client does not use the deterministic one-blank-line fragment boundaries")
	}
}

func TestBoardJavaScriptAssemblerRejectsInvalidStructure(t *testing.T) {
	newFixtureAssets := func(shellText string) fstest.MapFS {
		fixtureAssets := fstest.MapFS{
			boardJavaScriptShellPath: &fstest.MapFile{Data: []byte(shellText)},
		}
		for _, fragmentPath := range boardJavaScriptFragmentPaths {
			fixtureAssets[fragmentPath] = &fstest.MapFile{Data: []byte("// " + fragmentPath + "\n")}
		}
		return fixtureAssets
	}

	t.Run("missing placeholder", func(t *testing.T) {
		fixtureAssets := newFixtureAssets("(function () {\n})();\n")
		if _, assembleError := assembleBoardJavaScript(fixtureAssets); assembleError == nil {
			t.Fatal("assembler accepted a shell with no fragment placeholder")
		}
	})
	t.Run("duplicate placeholder", func(t *testing.T) {
		fixtureAssets := newFixtureAssets(boardJavaScriptPlaceholderLine + boardJavaScriptPlaceholderLine)
		if _, assembleError := assembleBoardJavaScript(fixtureAssets); assembleError == nil {
			t.Fatal("assembler accepted a shell with duplicate fragment placeholders")
		}
	})
	t.Run("duplicate noncanonical placeholder", func(t *testing.T) {
		fixtureAssets := newFixtureAssets(boardJavaScriptPlaceholderLine + boardJavaScriptPlaceholder)
		if _, assembleError := assembleBoardJavaScript(fixtureAssets); assembleError == nil {
			t.Fatal("assembler accepted a second fragment placeholder without a trailing LF")
		}
	})
	t.Run("lone noncanonical placeholder", func(t *testing.T) {
		fixtureAssets := newFixtureAssets(boardJavaScriptPlaceholder)
		if _, assembleError := assembleBoardJavaScript(fixtureAssets); assembleError == nil {
			t.Fatal("assembler accepted a fragment placeholder without its canonical trailing LF")
		}
	})
	t.Run("missing fragment", func(t *testing.T) {
		fixtureAssets := newFixtureAssets(boardJavaScriptPlaceholderLine)
		delete(fixtureAssets, boardJavaScriptFragmentPaths[0])
		if _, assembleError := assembleBoardJavaScript(fixtureAssets); assembleError == nil {
			t.Fatal("assembler accepted an omitted manifest fragment")
		}
	})
	t.Run("noncanonical fragment ending", func(t *testing.T) {
		fixtureAssets := newFixtureAssets(boardJavaScriptPlaceholderLine)
		fixtureAssets[boardJavaScriptFragmentPaths[0]] = &fstest.MapFile{Data: []byte("no trailing LF")}
		if _, assembleError := assembleBoardJavaScript(fixtureAssets); assembleError == nil {
			t.Fatal("assembler accepted a fragment without exactly one trailing LF")
		}
	})
	t.Run("double fragment ending", func(t *testing.T) {
		fixtureAssets := newFixtureAssets(boardJavaScriptPlaceholderLine)
		fixtureAssets[boardJavaScriptFragmentPaths[0]] = &fstest.MapFile{Data: []byte("extra blank line\n\n")}
		if _, assembleError := assembleBoardJavaScript(fixtureAssets); assembleError == nil {
			t.Fatal("assembler accepted a fragment with two trailing LFs")
		}
	})
}

func TestMain(testMain *testing.M) {
	exitCode := testMain.Run()
	if strictProbeGuardFailed(exitCode, os.Getenv(strictJavaScriptBehaviorMarker), javaScriptBehaviorProbeCount.Load()) {
		fmt.Fprintln(os.Stderr, strictJavaScriptBehaviorDiagnostic)
		exitCode = 1
	}
	// Same guard for the browser lane (browser_probe_test.go): a strict run whose
	// probes all skipped must not report green, which is what makes the ordinary
	// skip safe. Both lanes gate here because TestMain is per-package, not per-file.
	if strictProbeGuardFailed(exitCode, os.Getenv(strictBrowserBehaviorMarker), browserBehaviorProbeCount.Load()) {
		fmt.Fprintln(os.Stderr, strictBrowserBehaviorDiagnostic)
		exitCode = 1
	}
	os.Exit(exitCode)
}

func strictProbeGuardFailed(exitCode int, strictMarker string, executedProbeCount int64) bool {
	return exitCode == 0 && strictMarker == "1" && executedProbeCount == 0
}

func testEnvironmentWithOverrides(baseEnvironment []string, overrides ...string) []string {
	overriddenKeys := make(map[string]bool, len(overrides))
	for _, override := range overrides {
		overrideKey, _, hasValue := strings.Cut(override, "=")
		if hasValue {
			overriddenKeys[overrideKey] = true
		}
	}

	cleanEnvironment := make([]string, 0, len(baseEnvironment)+len(overrides))
	for _, environmentEntry := range baseEnvironment {
		environmentKey, _, hasValue := strings.Cut(environmentEntry, "=")
		if hasValue && overriddenKeys[environmentKey] {
			continue
		}
		cleanEnvironment = append(cleanEnvironment, environmentEntry)
	}
	return append(cleanEnvironment, overrides...)
}

func TestStrictProbeGuard(t *testing.T) {
	testCases := []struct {
		name       string
		exitCode   int
		marker     string
		probeCount int64
		want       bool
	}{
		{name: "strict empty lane fails", exitCode: 0, marker: "1", probeCount: 0, want: true},
		{name: "ordinary empty lane passes", exitCode: 0, marker: "", probeCount: 0, want: false},
		{name: "executed strict lane passes", exitCode: 0, marker: "1", probeCount: 1, want: false},
		{name: "existing failure stays a failure", exitCode: 1, marker: "1", probeCount: 0, want: false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := strictProbeGuardFailed(testCase.exitCode, testCase.marker, testCase.probeCount); got != testCase.want {
				t.Fatalf("strictProbeGuardFailed(%d, %q, %d) = %t, want %t",
					testCase.exitCode, testCase.marker, testCase.probeCount, got, testCase.want)
			}
		})
	}
}

func lookupNodeForJavaScriptProbe(t *testing.T) string {
	t.Helper()
	if os.Getenv(javaScriptBehaviorProbeMode) != "on" {
		t.Skip("JavaScript behavior probes are heavy-only")
	}
	nodePath, lookupError := exec.LookPath("node")
	if lookupError != nil {
		t.Skip("node is unavailable; skipping JavaScript behavior probe")
	}
	return nodePath
}

func runJavaScriptBehaviorProbe(t *testing.T, probeName string, javascriptProbe string) []byte {
	t.Helper()
	nodePath := lookupNodeForJavaScriptProbe(t)

	// The probe arrives on stdin rather than as an "-e" argument: a probe that
	// embeds the assembled client exceeds Linux's 128 KiB per-argument limit and
	// fails the exec with "argument list too long" — a limit macOS does not have,
	// so an "-e" invocation passes for the maintainer and fails in CI-like Linux
	// environments on probe size alone.
	probeCommand := exec.Command(nodePath, "-")
	probeCommand.Stdin = strings.NewReader(javascriptProbe)
	javaScriptBehaviorProbeCount.Add(1)
	probeOutput, probeError := probeCommand.CombinedOutput()
	if probeError != nil {
		t.Fatalf("execute %s JavaScript behavior: %v\n%s", probeName, probeError, probeOutput)
	}
	return probeOutput
}

// Syntax participates in the maintainer-strict test selection but deliberately
// does not increment javaScriptBehaviorProbeCount: a parse-only check cannot
// satisfy the lane's requirement for executable behavior coverage.
func publicationTestBoard(title string, bodyMarkdown string, projectName string, generatedAt time.Time) *Board {
	ticket := &RequestTicket{
		RequestId:      "REQ-1",
		Title:          title,
		Status:         "pending",
		OriginalStatus: "pending",
		BodyMarkdown:   bodyMarkdown,
	}
	return &Board{
		GeneratedAt: generatedAt,
		ProjectName: projectName,
		AllRequests: []*RequestTicket{ticket},
		Columns: BoardColumns{
			Pending:      []*RequestTicket{ticket},
			PendingReady: []*RequestTicket{ticket},
		},
	}
}

func readStaticSiteTargets(t *testing.T, outputDirectory string) map[string]string {
	t.Helper()
	targetContents := make(map[string]string, 3)
	for _, targetName := range []string{boardDataJsFilename, boardMarkdownJsFilename, "index.html"} {
		targetBytes, readError := os.ReadFile(filepath.Join(outputDirectory, targetName))
		if readError != nil {
			t.Fatalf("read %s: %v", targetName, readError)
		}
		targetContents[targetName] = string(targetBytes)
	}
	return targetContents
}

// generateLiveSiteInDir builds the board against the REAL do-work tree and writes
// the static site into a temp dir, returning the output directory path. The git
// lookup is stubbed (mirrors board_live_test.go) so the whole-tree build is fast
// and deterministic. Tests that need board-data.js or other sibling files should
// call this helper directly.
func generateLiveSiteInDir(t *testing.T) string {
	t.Helper()
	workingDirectory, getwdError := os.Getwd()
	if getwdError != nil {
		t.Fatalf("getwd: %v", getwdError)
	}
	repoRoot, resolveError := resolveRepoRoot(workingDirectory)
	if resolveError != nil {
		t.Fatalf("resolveRepoRoot: %v", resolveError)
	}
	stubGitLookup := func(string, string) (time.Time, bool) { return time.Time{}, false }
	board, buildError := buildBoard(repoRoot, time.Now(), 7*24*time.Hour, stubGitLookup)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}

	outputDirectory := t.TempDir()
	if generateError := generateStaticSite(outputDirectory, board); generateError != nil {
		t.Fatalf("generateStaticSite: %v", generateError)
	}
	return outputDirectory
}

// generateLiveSite builds the board and returns the index.html contents. It is a
// convenience wrapper over generateLiveSiteInDir for tests that only need the
// main HTML page.
func generateLiveSite(t *testing.T) string {
	t.Helper()
	liveIndexOnce.Do(func() {
		outputDirectory := generateLiveSiteInDir(t)
		indexPath := filepath.Join(outputDirectory, "index.html")
		indexBytes, readError := os.ReadFile(indexPath)
		if readError != nil {
			t.Fatalf("reading generated index.html: %v", readError)
		}
		liveIndexHTML = string(indexBytes)
	})
	return liveIndexHTML
}

func TestGenerateWritesSelfContainedIndex(t *testing.T) {
	indexHtml := generateLiveSite(t)

	// The page must be self-contained: CSS + JS inlined, no CDN / external asset.
	if !strings.Contains(indexHtml, "<style>") {
		t.Fatalf("generated page has no inlined <style> block")
	}
	for _, externalMarker := range []string{
		`src="http`,
		`src='http`,
		`href="http`,
		`<link rel="stylesheet"`,
		"cdn.",
	} {
		if strings.Contains(indexHtml, externalMarker) {
			t.Fatalf("generated page is not self-contained: found external reference %q", externalMarker)
		}
	}
	// The inlined assembled client must be present (a known function name).
	if !strings.Contains(indexHtml, "renderColumns") {
		t.Fatalf("inlined assembled-client behaviour is missing from the page")
	}
	// The display placeholder must have been resolved.
	if strings.Contains(indexHtml, "GENERATED_AT_DISPLAY") {
		t.Fatalf("GENERATED_AT_DISPLAY placeholder was not substituted")
	}
}

func TestGenerateFirstPublicationAndSuccessfulReplacement(t *testing.T) {
	outputDirectory := filepath.Join(t.TempDir(), "static-board")
	oldBoard := publicationTestBoard(
		"Old title",
		"## What\n\nOld body.\n",
		"old-project",
		time.Date(2026, 8, 14, 9, 0, 0, 0, time.UTC),
	)
	if generateError := generateStaticSite(outputDirectory, oldBoard); generateError != nil {
		t.Fatalf("first generation: %v", generateError)
	}
	oldTargets := readStaticSiteTargets(t, outputDirectory)

	unrelatedPath := filepath.Join(outputDirectory, "keep-me.txt")
	const unrelatedContents = "unrelated output stays untouched\n"
	if writeError := os.WriteFile(unrelatedPath, []byte(unrelatedContents), 0o644); writeError != nil {
		t.Fatalf("write unrelated output fixture: %v", writeError)
	}

	newBoard := publicationTestBoard(
		"New title",
		"## What\n\nNew body.\n",
		"new-project",
		time.Date(2026, 8, 15, 10, 0, 0, 0, time.UTC),
	)
	if generateError := generateStaticSite(outputDirectory, newBoard); generateError != nil {
		t.Fatalf("replacement generation: %v", generateError)
	}
	newTargets := readStaticSiteTargets(t, outputDirectory)
	for targetName, oldContents := range oldTargets {
		if newTargets[targetName] == oldContents {
			t.Errorf("%s was not replaced", targetName)
		}
	}
	if !strings.Contains(newTargets[boardDataJsFilename], "New title") {
		t.Errorf("replacement %s is missing the new title", boardDataJsFilename)
	}
	if !strings.Contains(newTargets[boardMarkdownJsFilename], "New body.") {
		t.Errorf("replacement %s is missing the new body", boardMarkdownJsFilename)
	}
	if !strings.Contains(newTargets["index.html"], "new-project") {
		t.Errorf("replacement index.html is missing the new project name")
	}
	unrelatedBytes, unrelatedReadError := os.ReadFile(unrelatedPath)
	if unrelatedReadError != nil {
		t.Errorf("read unrelated output: %v", unrelatedReadError)
	} else if string(unrelatedBytes) != unrelatedContents {
		t.Errorf("unrelated output changed: got %q, want %q", unrelatedBytes, unrelatedContents)
	}
	outputEntries, readDirectoryError := os.ReadDir(outputDirectory)
	if readDirectoryError != nil {
		t.Fatalf("read output directory: %v", readDirectoryError)
	}
	if len(outputEntries) != 4 {
		entryNames := make([]string, 0, len(outputEntries))
		for _, entry := range outputEntries {
			entryNames = append(entryNames, entry.Name())
		}
		t.Errorf("private publication residue remains: entries = %v", entryNames)
	}
}

func TestGenerateRefusesNonRegularOutputTargetsWithoutMutation(t *testing.T) {
	testCases := []struct {
		name       string
		targetName string
		obstruct   func(*testing.T, string, string) func(*testing.T)
	}{
		{
			name:       "directory",
			targetName: "index.html",
			obstruct: func(t *testing.T, _ string, targetPath string) func(*testing.T) {
				if mkdirError := os.Mkdir(targetPath, 0o755); mkdirError != nil {
					t.Fatalf("create target directory: %v", mkdirError)
				}
				keptPath := filepath.Join(targetPath, "kept.txt")
				const keptContents = "directory contents stay untouched\n"
				if writeError := os.WriteFile(keptPath, []byte(keptContents), 0o644); writeError != nil {
					t.Fatalf("write nested target fixture: %v", writeError)
				}
				return func(t *testing.T) {
					keptBytes, readError := os.ReadFile(keptPath)
					if readError != nil {
						t.Fatalf("read preserved nested target: %v", readError)
					}
					if string(keptBytes) != keptContents {
						t.Errorf("nested target contents = %q, want %q", keptBytes, keptContents)
					}
				}
			},
		},
		{
			name:       "symlink",
			targetName: boardMarkdownJsFilename,
			obstruct: func(t *testing.T, caseRoot string, targetPath string) func(*testing.T) {
				symlinkTargetPath := filepath.Join(caseRoot, "symlink-target.txt")
				const targetContents = "symlink target stays untouched\n"
				if writeError := os.WriteFile(symlinkTargetPath, []byte(targetContents), 0o644); writeError != nil {
					t.Fatalf("write symlink target fixture: %v", writeError)
				}
				if symlinkError := os.Symlink(symlinkTargetPath, targetPath); symlinkError != nil {
					t.Skipf("symlinks unavailable: %v", symlinkError)
				}
				return func(t *testing.T) {
					linkDestination, readlinkError := os.Readlink(targetPath)
					if readlinkError != nil {
						t.Fatalf("read preserved target symlink: %v", readlinkError)
					}
					if linkDestination != symlinkTargetPath {
						t.Errorf("target symlink destination = %q, want %q", linkDestination, symlinkTargetPath)
					}
					targetBytes, readError := os.ReadFile(symlinkTargetPath)
					if readError != nil {
						t.Fatalf("read symlink destination: %v", readError)
					}
					if string(targetBytes) != targetContents {
						t.Errorf("symlink destination contents = %q, want %q", targetBytes, targetContents)
					}
				}
			},
		},
		{
			name:       "special file",
			targetName: boardDataJsFilename,
			obstruct: func(t *testing.T, _ string, targetPath string) func(*testing.T) {
				mkfifoCommand := exec.Command("mkfifo", targetPath)
				if mkfifoOutput, mkfifoError := mkfifoCommand.CombinedOutput(); mkfifoError != nil {
					t.Skipf("named pipes unavailable: %v\n%s", mkfifoError, mkfifoOutput)
				}
				return func(t *testing.T) {
					targetInfo, lstatError := os.Lstat(targetPath)
					if lstatError != nil {
						t.Fatalf("inspect preserved named-pipe target: %v", lstatError)
					}
					if targetInfo.Mode()&os.ModeNamedPipe == 0 {
						t.Errorf("target mode = %v, want named pipe", targetInfo.Mode())
					}
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			caseRoot := t.TempDir()
			outputDirectory := filepath.Join(caseRoot, "static-board")
			if mkdirError := os.Mkdir(outputDirectory, 0o755); mkdirError != nil {
				t.Fatalf("create output directory: %v", mkdirError)
			}
			originalTargets := map[string]string{
				boardDataJsFilename:     "old board data\n",
				boardMarkdownJsFilename: "old board markdown\n",
				"index.html":            "old index\n",
			}
			for targetName, targetContents := range originalTargets {
				targetPath := filepath.Join(outputDirectory, targetName)
				if writeError := os.WriteFile(targetPath, []byte(targetContents), 0o644); writeError != nil {
					t.Fatalf("write %s fixture: %v", targetName, writeError)
				}
			}

			obstructedPath := filepath.Join(outputDirectory, testCase.targetName)
			if removeError := os.Remove(obstructedPath); removeError != nil {
				t.Fatalf("remove regular %s fixture: %v", testCase.targetName, removeError)
			}
			assertObstructionPreserved := testCase.obstruct(t, caseRoot, obstructedPath)

			board := publicationTestBoard(
				"Replacement title",
				"## What\n\nReplacement body.\n",
				"replacement-project",
				time.Date(2026, 8, 16, 11, 0, 0, 0, time.UTC),
			)
			publishCalls := 0
			countedPublisher := func(stagedPath string, targetPath string) error {
				publishCalls++
				return os.Rename(stagedPath, targetPath)
			}
			if generationError := generateStaticSiteWithPublisher(outputDirectory, board, countedPublisher); generationError == nil {
				t.Error("generateStaticSiteWithPublisher succeeded with a non-regular output target")
			}
			if publishCalls != 0 {
				t.Errorf("publication calls = %d, want 0 before refusing a non-regular target", publishCalls)
			}

			assertObstructionPreserved(t)
			for targetName, targetContents := range originalTargets {
				if targetName == testCase.targetName {
					continue
				}
				targetBytes, readError := os.ReadFile(filepath.Join(outputDirectory, targetName))
				if readError != nil {
					t.Errorf("read preserved %s: %v", targetName, readError)
					continue
				}
				if string(targetBytes) != targetContents {
					t.Errorf("%s contents = %q, want %q", targetName, targetBytes, targetContents)
				}
			}

			outputEntries, readDirectoryError := os.ReadDir(outputDirectory)
			if readDirectoryError != nil {
				t.Fatalf("read output directory: %v", readDirectoryError)
			}
			if len(outputEntries) != len(originalTargets) {
				entryNames := make([]string, 0, len(outputEntries))
				for _, outputEntry := range outputEntries {
					entryNames = append(entryNames, outputEntry.Name())
				}
				t.Errorf("refused publication left private residue: entries = %v", entryNames)
			}
		})
	}
}

func TestGeneratePublicationFailureRestoresThePreviousBundle(t *testing.T) {
	outputDirectory := t.TempDir()
	oldBoard := publicationTestBoard(
		"Old title",
		"## What\n\nOld body.\n",
		"old-project",
		time.Date(2026, 8, 14, 9, 0, 0, 0, time.UTC),
	)
	if generateError := generateStaticSite(outputDirectory, oldBoard); generateError != nil {
		t.Fatalf("generate initial static site: %v", generateError)
	}

	unrelatedPath := filepath.Join(outputDirectory, "keep-me.txt")
	const unrelatedContents = "unrelated output stays untouched\n"
	if writeError := os.WriteFile(unrelatedPath, []byte(unrelatedContents), 0o644); writeError != nil {
		t.Fatalf("write unrelated output fixture: %v", writeError)
	}
	oldTargets := readStaticSiteTargets(t, outputDirectory)

	newBoard := publicationTestBoard(
		"New title",
		"## What\n\nNew body.\n",
		"new-project",
		time.Date(2026, 8, 15, 10, 0, 0, 0, time.UTC),
	)
	seededPublicationFailure := errors.New("seeded static publication failure")
	publishCalls := 0
	failSecondPublication := func(stagedPath string, targetPath string) error {
		publishCalls++
		if publishCalls == 2 {
			return seededPublicationFailure
		}
		return os.Rename(stagedPath, targetPath)
	}
	generationError := generateStaticSiteWithPublisher(outputDirectory, newBoard, failSecondPublication)
	if !errors.Is(generationError, seededPublicationFailure) {
		t.Errorf("generation error = %v, want seeded publication failure", generationError)
	}
	if publishCalls != 2 {
		t.Errorf("publication calls = %d, want 2", publishCalls)
	}

	for targetName, oldContents := range oldTargets {
		currentBytes, readError := os.ReadFile(filepath.Join(outputDirectory, targetName))
		if readError != nil {
			t.Errorf("read restored %s: %v", targetName, readError)
			continue
		}
		if string(currentBytes) != oldContents {
			t.Errorf("%s differs after failed publication", targetName)
		}
	}
	unrelatedBytes, unrelatedReadError := os.ReadFile(unrelatedPath)
	if unrelatedReadError != nil {
		t.Errorf("read unrelated output: %v", unrelatedReadError)
	} else if string(unrelatedBytes) != unrelatedContents {
		t.Errorf("unrelated output changed: got %q, want %q", unrelatedBytes, unrelatedContents)
	}
	outputEntries, readDirectoryError := os.ReadDir(outputDirectory)
	if readDirectoryError != nil {
		t.Fatalf("read output directory: %v", readDirectoryError)
	}
	if len(outputEntries) != 4 {
		entryNames := make([]string, 0, len(outputEntries))
		for _, entry := range outputEntries {
			entryNames = append(entryNames, entry.Name())
		}
		t.Errorf("private publication residue remains: entries = %v", entryNames)
	}
}

func TestGenerateRendersColumnHeaders(t *testing.T) {
	indexHtml := generateLiveSite(t)
	for _, columnHeader := range []string{
		"Pending",
		"Claimed",
		"Needs input",
		"Recently done",
	} {
		if !strings.Contains(indexHtml, columnHeader) {
			t.Fatalf("column header %q not found in the generated page", columnHeader)
		}
	}
}

func TestGenerateEmbedsLivePendingCards(t *testing.T) {
	// After REQ-1213 the card data (including REQ IDs) lives in board-data.js. The
	// expected ids are derived from the live board, not hard-coded — the old test
	// pinned REQ-1207..1210 from the source monorepo, which don't exist in this
	// extraction. Exact seeded-card coverage lives in the synthetic board tests.
	board := liveBoard(t)
	if len(board.Columns.Pending) == 0 {
		t.Skip("no pending REQs in the live tree; nothing to assert")
	}

	outputDirectory := generateLiveSiteInDir(t)
	boardDataPath := filepath.Join(outputDirectory, "board-data.js")
	boardDataBytes, readError := os.ReadFile(boardDataPath)
	if readError != nil {
		t.Fatalf("reading board-data.js: %v", readError)
	}
	boardDataJs := string(boardDataBytes)

	checks := 0
	for _, ticket := range board.Columns.Pending {
		if !strings.Contains(boardDataJs, ticket.RequestId) {
			t.Fatalf("live pending id %q not found in board-data.js", ticket.RequestId)
		}
		checks++
		if checks >= 25 {
			break // a representative sample is enough
		}
	}
}

func TestGenerateEmitsBlockedFields(t *testing.T) {
	// The synthetic tree seeds REQ-9006 as status: blocked with a free-text
	// blocked_by, a blocked_at, and a blocked_check. Those must survive into the
	// generated payload so the frontend can render the "blocked by" badge/drawer.
	board := syntheticBoard(t)
	generatedData, buildError := buildGeneratedBoardData(board)
	if buildError != nil {
		t.Fatalf("buildGeneratedBoardData: %v", buildError)
	}
	blockedRequest, present := generatedData.Requests["REQ-9006"]
	if !present {
		t.Fatalf("REQ-9006 (blocked) missing from generated requests")
	}
	if blockedRequest.Status != "blocked" {
		t.Fatalf("REQ-9006 status = %q, want blocked", blockedRequest.Status)
	}
	if len(blockedRequest.BlockedBy) != 1 || blockedRequest.BlockedBy[0] != "LM Studio running locally" {
		t.Fatalf("REQ-9006 blockedBy = %+v, want [\"LM Studio running locally\"]", blockedRequest.BlockedBy)
	}
	if blockedRequest.BlockedCheck == "" || blockedRequest.BlockedAt == "" {
		t.Fatalf("REQ-9006 blockedCheck/blockedAt not populated: check=%q at=%q", blockedRequest.BlockedCheck, blockedRequest.BlockedAt)
	}
	// The fields must also survive JSON marshaling under their camelCase keys.
	marshaledBytes, marshalError := json.Marshal(blockedRequest)
	if marshalError != nil {
		t.Fatalf("marshal generated request: %v", marshalError)
	}
	marshaledJson := string(marshaledBytes)
	for _, expectedKey := range []string{`"blockedBy"`, `"blockedAt"`, `"blockedCheck"`} {
		if !strings.Contains(marshaledJson, expectedKey) {
			t.Fatalf("generated JSON missing %s: %s", expectedKey, marshaledJson)
		}
	}
}

func TestGenerateIndexHtmlUnderSizeBudget(t *testing.T) {
	// The JSON data island (all pre-rendered REQ bodies) must be externalized to
	// board-data.js so index.html stays well under 1 MB. Before REQ-1213 the
	// monolithic file weighed ~14 MB.
	const maxIndexHtmlBytes = 1 << 20 // 1 MiB
	indexHtml := generateLiveSite(t)
	actualBytes := len(indexHtml)
	if actualBytes >= maxIndexHtmlBytes {
		t.Fatalf("index.html is %d bytes (%.1f MB) — exceeds the 1 MB budget; externalize the JSON data island to board-data.js",
			actualBytes, float64(actualBytes)/(1<<20))
	}
}

func TestGenerateHasCalendarAndUserRequestLensHooks(t *testing.T) {
	outputDirectory := generateLiveSiteInDir(t)

	indexPath := filepath.Join(outputDirectory, "index.html")
	indexBytes, readError := os.ReadFile(indexPath)
	if readError != nil {
		t.Fatalf("reading generated index.html: %v", readError)
	}
	indexHtml := string(indexBytes)

	if !strings.Contains(indexHtml, `data-view-target="calendar"`) {
		t.Fatalf("calendar view hook not found")
	}
	if !strings.Contains(indexHtml, `data-lens-target="user-request"`) {
		t.Fatalf("by-UR lens toggle hook not found")
	}

	// Calendar day-keyed completion entries live in the externalized board-data.js.
	boardDataPath := filepath.Join(outputDirectory, "board-data.js")
	boardDataBytes, bdReadError := os.ReadFile(boardDataPath)
	if bdReadError != nil {
		t.Fatalf("reading board-data.js: %v", bdReadError)
	}
	if !strings.Contains(string(boardDataBytes), `"dayKey"`) {
		t.Fatalf("calendar entries (dayKey) not found in board-data.js")
	}
}

func TestGenerateEmbedsGoldmarkRenderedBody(t *testing.T) {
	// After REQ-1213 the JSON data island (including pre-rendered bodies) lives in
	// board-data.js, not in index.html. Read the sibling file for assertions.
	outputDirectory := generateLiveSiteInDir(t)
	boardDataPath := filepath.Join(outputDirectory, "board-data.js")
	boardDataBytes, readError := os.ReadFile(boardDataPath)
	if readError != nil {
		t.Fatalf("reading board-data.js: %v", readError)
	}
	boardDataJs := string(boardDataBytes)

	// Every REQ body in this repo has `## ` headings; goldmark (with auto heading
	// IDs) renders them to `<h2 id="...">`. Asserting the id form proves the
	// marker came from a rendered REQ body — not from the page chrome.
	if !strings.Contains(boardDataJs, `<h2 id=`) {
		t.Fatalf("no goldmark-rendered `<h2 id=` body heading found in board-data.js")
	}
	// The data island must carry pre-rendered bodies under the bodyHtml key.
	if !strings.Contains(boardDataJs, `"bodyHtml"`) {
		t.Fatalf("board-data.js has no bodyHtml field")
	}
}

func TestGenerateSeparatesRawMarkdownForLazyCopy(t *testing.T) {
	// Copy still needs exact source, but shipping it beside bodyHtml nearly
	// doubles the initial payload. Raw bodies belong in a lazy sibling script.
	outputDirectory := generateLiveSiteInDir(t)

	boardDataBytes, readError := os.ReadFile(filepath.Join(outputDirectory, "board-data.js"))
	if readError != nil {
		t.Fatalf("reading board-data.js: %v", readError)
	}
	if strings.Contains(string(boardDataBytes), `"bodyMarkdown"`) {
		t.Fatalf("board-data.js still carries bodyMarkdown — raw source must stay out of the initial payload")
	}

	boardMarkdownBytes, markdownReadError := os.ReadFile(filepath.Join(outputDirectory, boardMarkdownJsFilename))
	if markdownReadError != nil {
		t.Fatalf("reading %s: %v", boardMarkdownJsFilename, markdownReadError)
	}
	if !strings.HasPrefix(string(boardMarkdownBytes), "window.queueKanbanBoardMarkdownData = ") {
		t.Fatalf("%s does not assign the lazy Markdown global", boardMarkdownJsFilename)
	}

	indexBytes, indexReadError := os.ReadFile(filepath.Join(outputDirectory, "index.html"))
	if indexReadError != nil {
		t.Fatalf("reading generated index.html: %v", indexReadError)
	}
	if !strings.Contains(string(indexBytes), `id="detail-copy"`) {
		t.Fatalf("detail drawer Copy button (id=\"detail-copy\") not found in index.html")
	}
	if strings.Contains(string(indexBytes), `<script src="board-markdown.js"></script>`) {
		t.Fatalf("index.html eagerly loads board-markdown.js; raw source must load only after Copy")
	}
	if !strings.Contains(string(indexBytes), `markdownScript.src = "board-markdown.js"`) {
		t.Fatalf("inlined board-clipboard.js has no lazy board-markdown.js loader")
	}
	// Since REQ-089 the lazy payload holds whole FILES (frontmatter fence + body),
	// so the primary Copy path keeps the file's own bytes — no synthesized heading, or the
	// paste stops round-tripping back into a valid REQ file. The identifying heading
	// belongs to the rendered-text fallback alone, which has no frontmatter to carry.
	if !strings.Contains(string(indexBytes), "copyTextWithHeading(requestedKind, requestedId, renderedTextFallback)") {
		t.Fatalf("inlined board-clipboard.js does not prepend the id/title heading on the rendered-text fallback path")
	}
	if strings.Contains(string(indexBytes), "copyTextWithHeading(requestedKind, requestedId, bodyText)") {
		t.Fatalf("inlined board-clipboard.js still routes the lazy payload through the heading builder — the primary path must start from the file's own bytes, annotating only the body")
	}
}

func TestBuildGeneratedBoardMarkdownDataKeepsExactSources(t *testing.T) {
	board := &Board{
		AllRequests: []*RequestTicket{
			{RequestId: "REQ-1", BodyMarkdown: "## What\n\n- [ ] keep formatting\n"},
		},
		UserRequests: []*UserRequestTicket{
			{UserRequestId: "UR-1", InputFilePresent: true, BodyMarkdown: "# Original request\n\nExact text.\n"},
		},
	}

	markdownData := buildGeneratedBoardMarkdownData(board)
	if got := markdownData.Requests["REQ-1"]; got != board.AllRequests[0].BodyMarkdown {
		t.Fatalf("REQ raw Markdown changed: got %q, want %q", got, board.AllRequests[0].BodyMarkdown)
	}
	if got := markdownData.UserRequests["UR-1"]; got != board.UserRequests[0].BodyMarkdown {
		t.Fatalf("UR raw Markdown changed: got %q, want %q", got, board.UserRequests[0].BodyMarkdown)
	}
}

func TestStaticAndLiveCitationDataShareOneAnalysisAndRefreshTogether(t *testing.T) {
	requestText := "---\nid: REQ-378\ntitle: Find referenced work\nstatus: pending\nrelated: REQ-1679\n---\n\nEmoji 😀 cites REQ-1679.\n"
	userRequestText := "---\nid: UR-076\ntitle: Reference search\nrequests: [REQ-378]\n---\n\n[REQ-1679](https://example.test/)\n"
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"do-work/queue/REQ-378-search.md", requestText},
		{"do-work/queue/REQ-1679-target.md", "---\nid: REQ-1679\ntitle: Target work\nstatus: pending\n---\n\nTarget body.\n"},
		{"do-work/user-requests/UR-076/input.md", userRequestText},
	})
	originalParser := markdownToHtmlRenderer.Parser()
	countingParser := &citationCountingParser{Parser: originalParser, BodyParses: map[string]int{}}
	markdownToHtmlRenderer.SetParser(countingParser)
	t.Cleanup(func() { markdownToHtmlRenderer.SetParser(originalParser) })
	board, err := buildBoard(repoRoot, time.Now(), defaultRecentWindow, nil)
	if err != nil {
		t.Fatal(err)
	}
	outputDirectory := t.TempDir()
	if err := generateStaticSite(outputDirectory, board); err != nil {
		t.Fatal(err)
	}
	assertSingleAnalysis := func() {
		t.Helper()
		for _, raw := range []string{requestText, userRequestText} {
			_, body, _, _ := splitFrontmatter(raw)
			if got := countingParser.BodyParses[body]; got != 2 {
				t.Errorf("body parsed %d times, want exactly 2 (HTML + shared raw analysis): %q", got, body)
			}
		}
	}
	assertSingleAnalysis()
	payload, err := os.ReadFile(filepath.Join(outputDirectory, boardDataJsFilename))
	if err != nil {
		t.Fatal(err)
	}
	var staticData generatedBoardData
	jsonText := strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(string(payload)), "window.queueKanbanBoardData = "), ";")
	if err := json.Unmarshal([]byte(jsonText), &staticData); err != nil {
		t.Fatal(err)
	}
	for id, record := range staticData.Requests {
		want := []string{}
		if id == "REQ-378" {
			want = []string{"REQ-1679"}
		}
		if !reflect.DeepEqual(record.CitedTicketIds, want) {
			t.Errorf("static %s citations = %v, want %v", id, record.CitedTicketIds, want)
		}
	}
	if !reflect.DeepEqual(staticData.UserRequests["UR-076"].CitedTicketIds, []string{"REQ-1679"}) {
		t.Fatal("static UR lost authored-link citation")
	}
	countingParser.BodyParses = map[string]int{}
	server := httptest.NewServer(newLiveBoardServer(repoRoot, defaultRecentWindow))
	defer server.Close()
	liveData := fetchServedBoardData(t, server.URL)
	markdownData := fetchServedBoardMarkdownData(t, server.URL)
	assertSingleAnalysis() // lazy fetch reused the same fresh model and analysis
	if !reflect.DeepEqual(liveData.Requests, staticData.Requests) || !reflect.DeepEqual(liveData.UserRequests, staticData.UserRequests) {
		t.Fatal("static/live citation records disagree")
	}
	if markdownData.Requests["REQ-378"] != requestText || markdownData.UserRequests["UR-076"] != userRequestText {
		t.Fatal("search generation changed raw source bytes")
	}
	if got := describeTicketMentions(requestText, markdownData.RequestMentions["REQ-378"]); !reflect.DeepEqual(got, []string{`req REQ-1679 EXPAND "REQ-1679"`}) {
		t.Fatalf("clipboard offsets/annotation shape changed: %v", got)
	}
	// The tree changes after initial paint. Lazy Copy and the refreshed eager
	// index must come from the new text, not a long-lived analysis cache.
	requestText = strings.ReplaceAll(requestText, "REQ-1679", "UR-076")
	requestPath := filepath.Join(repoRoot, "do-work/queue/REQ-378-search.md")
	if err := os.WriteFile(requestPath, []byte(requestText), 0o644); err != nil {
		t.Fatal(err)
	}
	modifiedAt := time.Now().Add(2 * time.Second)
	if err := os.Chtimes(requestPath, modifiedAt, modifiedAt); err != nil {
		t.Fatal(err)
	}
	countingParser.BodyParses = map[string]int{}
	markdownData = fetchServedBoardMarkdownData(t, server.URL)
	liveData = fetchServedBoardData(t, server.URL)
	assertSingleAnalysis()
	if !reflect.DeepEqual(liveData.Requests["REQ-378"].CitedTicketIds, []string{"UR-076"}) || markdownData.Requests["REQ-378"] != requestText {
		t.Fatal("live citation index and raw source did not refresh together")
	}
}

func TestBrowserBehaviorCitationSearchShowsReasonsAcrossViews(t *testing.T) {
	lookupBrowserForBehaviorProbe(t)
	board := citationSearchFixtureBoard()
	board.GeneratedAt = time.Now().UTC()
	board.ProjectName = "Citation search acceptance"
	board.Columns.PendingReady = board.AllRequests
	board.Columns.Pending = board.AllRequests
	for _, parent := range board.UserRequests {
		if parent.UserRequestId == "UR-075" {
			parent.RequestIds = []string{"REQ-378"}
		}
	}
	// Testing uses the same request card constructor and must retain its reason.
	completed := ticketMentionFixtureTicket("REQ-379", "Completed citation", "completed", "REQ-1679\n")
	completed.CompletedAt = board.GeneratedAt.Add(-time.Hour).Format(time.RFC3339)
	completed.CompletionTime = board.GeneratedAt.Add(-time.Hour)
	board.AllRequests = append(board.AllRequests, completed)
	siteDirectory := t.TempDir()
	if err := generateStaticSite(siteDirectory, board); err != nil {
		t.Fatal(err)
	}
	indexBytes, err := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	// Capture mount errors as well as interaction errors, before any client code.
	page := strings.Replace(string(indexBytes), "<head>", `<head><script>
window.citationProbeErrors = [];
addEventListener("error", function(event) { window.citationProbeErrors.push(event.message); });
addEventListener("unhandledrejection", function(event) { window.citationProbeErrors.push(String(event.reason)); });
console.warn = console.error = function() { window.citationProbeErrors.push(Array.from(arguments).join(" ")); };
</script>`, 1)
	session := startTrustedInputBrowserSession(t, "citation search", siteDirectory, page)
	session.waitForPageCondition(t, "initial request cards", `document.querySelector('.req-card[data-detail-id="REQ-378"]') !== null`)
	session.evaluateInPage(t, `(function(){document.getElementById("filter-search").focus();return true;})()`)
	session.callDevToolsMethod(t, "Input.insertText", map[string]any{"text": "REQ-1679"}, true)
	session.waitForPageCondition(t, "citation-only match", `document.querySelector('[data-cards="pending"] .req-card[data-detail-id="REQ-378"] .citation-match') !== null`)
	for _, scheme := range []string{"light", "dark"} {
		for _, width := range []int{320, 768, 1280} {
			t.Run(fmt.Sprintf("%s-%d", scheme, width), func(t *testing.T) {
				session.callDevToolsMethod(t, "Emulation.setEmulatedMedia", map[string]any{"features": []map[string]string{{"name": "prefers-color-scheme", "value": scheme}}}, true)
				session.callDevToolsMethod(t, "Emulation.setDeviceMetricsOverride", map[string]any{"width": width, "height": 900, "deviceScaleFactor": 1, "mobile": false}, true)
				for _, view := range []string{"columns", "user-request", "testing"} {
					var result struct {
						Href, Browser, Marker, AccessibleName string
						Ids, Errors                           []string
						MarkerWidth, MarkerHeight             float64
						Contained, Focusable, LazyUnloaded    bool
					}
					probe := `(function(){
var view = ` + mustMarshalJSONString(t, view) + `;
document.querySelector('[data-view-target="' + (view === "testing" ? "testing" : "board") + '"]').click();
if (view !== "testing") document.querySelector('[data-lens-target="' + (view === "columns" ? "flat" : "user-request") + '"]:not([data-ur-cards])').click();
var root = document.querySelector(view === "columns" ? '[data-cards="pending"]' : view === "testing" ? '#view-testing' : '#user-request-lens');
var button = root.querySelector(view === "user-request" ? '.ur-group-head[data-detail-id="UR-075"]' : '.req-card[data-detail-id="' + (view === "testing" ? 'REQ-379' : 'REQ-378') + '"]');
var marker = button && button.querySelector('.citation-match');
var bounds = marker && marker.getBoundingClientRect(), container = button && button.getBoundingClientRect();
if (button) button.focus();
return {href:location.href, browser:navigator.userAgent, marker:marker && marker.textContent,
accessibleName:button && (button.getAttribute('aria-label') || button.textContent + ' ' + (marker && marker.getAttribute('aria-label'))),
ids:Array.from(root.querySelectorAll('.req-card'), function(card){return card.dataset.detailId;}),
markerWidth:bounds && bounds.width, markerHeight:bounds && bounds.height,
contained:!!bounds && bounds.left >= container.left && bounds.right <= container.right && bounds.top >= container.top && bounds.bottom <= container.bottom,
focusable:!!button && document.activeElement === button,
lazyUnloaded:!window.queueKanbanBoardMarkdownData, errors:window.citationProbeErrors};
})()`
					session.decodeResult(t, "rendered citation reason", session.evaluateInPage(t, probe), &result)
					t.Logf("%s/%s/%d: %+v", view, scheme, width, result)
					if result.Marker == "" || result.MarkerWidth <= 0 || result.MarkerHeight <= 0 || !result.Contained || !result.Focusable || !result.LazyUnloaded || len(result.Errors) != 0 || !strings.Contains(strings.ToLower(result.AccessibleName), "cites req-1679") {
						t.Fatalf("citation reason is missing, clipped, inaccessible, or dependent on lazy Copy: %+v", result)
					}
					if view == "columns" && !reflect.DeepEqual(result.Ids, []string{"REQ-1679", "REQ-378"}) {
						t.Fatalf("citation search returned %v", result.Ids)
					}
					if view == "testing" && !reflect.DeepEqual(result.Ids, []string{"REQ-379"}) {
						t.Fatalf("Testing citation search returned %v", result.Ids)
					}
				}
			})
		}
	}
	// A partial title remains a plain hit, and clearing removes stale reasons.
	session.evaluateInPage(t, `(function(){
document.querySelector('[data-view-target="board"]').click();
document.querySelector('[data-lens-target="flat"]').click();
var search=document.getElementById('filter-search');search.value='referenced';search.dispatchEvent(new Event('input',{bubbles:true}));return true;
})()`)
	session.waitForPageCondition(t, "ordinary partial title match without a citation reason", `document.querySelector('[data-cards="pending"] .req-card[data-detail-id="REQ-378"]') !== null && !document.querySelector('[data-cards="pending"] .citation-match')`)
	session.evaluateInPage(t, `(function(){document.getElementById('filter-clear').click();return true;})()`)
	session.waitForPageCondition(t, "cleared citation markers", `document.getElementById('filter-search').value === '' && !document.querySelector('[data-cards="pending"] .citation-match')`)
}

// The first pass must decorate author-owned anchors without replacing their
// destinations; the second pass must neither nest links nor spend expansion again.
func TestBrowserBehaviorAuthoredTicketLinksPreserveDestinationsAndTwoPassDOM(t *testing.T) {
	lookupBrowserForBehaviorProbe(t)
	longTitle := "Make REQ-7777 and every referenced request identifier in a drawer body carry its own title"
	body := "See [REQ-1108](https://example.com/spec?x=1&y=%2F) for the shape. " +
		"Again [REQ-1108](https://example.com/later), then REQ-1108.\n\n" +
		"[`REQ-1679`](https://example.com/code) then [REQ-1679](https://example.com/prose).\n\n" +
		"REQ-2000 before [REQ-2000](https://example.com/after).\n\n" +
		"[REQ-3000](REQ-3000), [REQ-3000 and UR-075 and REQ-9999](../notes.md#part).\n\n" +
		"[`REQ-4000`](https://example.com/only-code).\n\n" +
		"https://example.com/REQ-5000 and www.example.com/REQ-5000 and REQ-5000@example.com and <mailto:REQ-5000@example.com>.\n\n" +
		"<custom:REQ-5000/é> <mailto:REQ-5000@éxample.com> <custom:REQ-5000/%2F/é> " +
		"<mailto:REQ-5000+%2F@éxample.com> <custom:REQ-5000/%zz/é> <mailto:REQ-5000+%@éxample.com> <custom:REQ-5000/%C3%A9>.\n"
	board := &Board{GeneratedAt: time.Now().UTC(), ProjectName: "Authored ticket links"}
	for _, entry := range []struct{ id, title, body string }{
		{"REQ-382", "Authored references", body}, {"REQ-1108", longTitle, ""},
		{"REQ-1679", "Code does not spend expansion", ""}, {"REQ-2000", "Plain before authored", ""},
		{"REQ-3000", "Relative destination", ""}, {"REQ-4000", "Code only reference", ""},
		{"REQ-5000", "Autolink is not a citation", ""}, {"REQ-7777", "Inserted title is not rescanned", ""},
	} {
		board.AllRequests = append(board.AllRequests, ticketMentionFixtureTicket(entry.id, entry.title, "pending", entry.body))
	}
	board.Columns.PendingReady = board.AllRequests
	board.Columns.Pending = board.AllRequests
	board.UserRequests = []*UserRequestTicket{{UserRequestId: "UR-075", Title: "Parent reference"}}
	siteDirectory := t.TempDir()
	if err := generateStaticSite(siteDirectory, board); err != nil {
		t.Fatal(err)
	}
	indexBytes, err := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	page := strings.Replace(string(indexBytes), "<head>", `<head><script>
window.authoredLinkErrors=[];
addEventListener('error',function(event){window.authoredLinkErrors.push(event.message);});
addEventListener('unhandledrejection',function(event){window.authoredLinkErrors.push(String(event.reason));});
console.warn=console.error=function(){window.authoredLinkErrors.push(Array.from(arguments).join(' '));};
</script>`, 1)
	page = strings.Replace(page, "function linkifyDetailBody(", "window.authoredLinkProbe=linkifyDetailBody;\nfunction linkifyDetailBody(", 1)
	session := startTrustedInputBrowserSession(t, "authored ticket links", siteDirectory, page)
	session.waitForPageCondition(t, "request card", `document.querySelector('.req-card[data-detail-id="REQ-382"]') !== null`)
	session.evaluateInPage(t, `(function(){document.querySelector('.req-card[data-detail-id="REQ-382"]').click();return true;})()`)
	type anchorResult struct {
		Href, Text, Title               string
		Nested, Navigation, Decorations int
	}
	var result struct {
		Href, Browser            string
		Anchors                  []anchorResult
		Glossary, SecondGlossary []string
		GlossaryText             string
		NavigationPreserved      bool
		Identical                bool
		Missing, Nested          int
		Errors                   []string
	}
	session.decodeResult(t, "authored links and second pass", session.evaluateInPage(t, `(function(){
var root=document.getElementById('detail-body');var first=root.innerHTML;
var anchors=Array.from(root.querySelectorAll('a:not(.ticket-link)')).map(function(a){return {
href:a.getAttribute('href'),text:a.textContent,title:(a.querySelector('[title]')||a).title,
nested:a.querySelectorAll('a').length,decorations:a.querySelectorAll('.ticket-link').length,navigation:a.querySelectorAll('[data-detail-kind]').length+(a.hasAttribute('data-detail-kind')?1:0)};});
var glossary=Array.from(document.querySelectorAll('#detail-glossary dt')).map(function(dt){return dt.textContent;});
var navigationPreserved=true;
function observeClick(event){navigationPreserved=navigationPreserved&&!event.defaultPrevented;event.preventDefault();}
document.addEventListener('click',observeClick);
root.querySelector('a').click();root.querySelector('a[href="REQ-3000"]').click();
document.removeEventListener('click',observeClick);
navigationPreserved=navigationPreserved&&document.getElementById('detail-id').textContent==='REQ-382';
var second=window.authoredLinkProbe(root,'Authored references');
return {href:location.href,browser:navigator.userAgent,anchors:anchors,glossary:glossary,secondGlossary:second.map(function(e){return e.id;}),
glossaryText:document.getElementById('detail-glossary').textContent,navigationPreserved:navigationPreserved,identical:root.innerHTML===first,missing:root.querySelectorAll('a .ticket-missing').length,nested:root.querySelectorAll('a a').length,errors:window.authoredLinkErrors};
})()`), &result)
	t.Logf("rendered authored links: %+v", result)
	if !result.NavigationPreserved || !result.Identical || result.Nested != 0 || result.Missing != 0 || len(result.Errors) != 0 {
		t.Fatalf("second pass changed DOM, nested links, flagged author text, or failed: %+v", result)
	}
	wantTexts := []string{
		"REQ-1108 Make REQ-7777 and every referenced request identifier in a…", "REQ-1108", "REQ-1679", "REQ-1679 Code does not spend expansion",
		"REQ-2000", "REQ-3000 Relative destination", "REQ-3000 and UR-075 Parent reference and REQ-9999", "REQ-4000",
		"https://example.com/REQ-5000", "www.example.com/REQ-5000", "REQ-5000@example.com", "mailto:REQ-5000@example.com",
		"custom:REQ-5000/é", "mailto:REQ-5000@éxample.com", "custom:REQ-5000/%2F/é",
		"mailto:REQ-5000+%2F@éxample.com", "custom:REQ-5000/%zz/é", "mailto:REQ-5000+%@éxample.com", "custom:REQ-5000/%C3%A9",
	}
	wantHrefs := []string{"https://example.com/spec?x=1&y=%2F", "https://example.com/later", "https://example.com/code", "https://example.com/prose", "https://example.com/after", "REQ-3000", "../notes.md#part", "https://example.com/only-code", "https://example.com/REQ-5000", "http://www.example.com/REQ-5000", "mailto:REQ-5000@example.com", "mailto:REQ-5000@example.com",
		"custom:REQ-5000/%C3%A9", "mailto:REQ-5000@%C3%A9xample.com", "custom:REQ-5000/%2F/%C3%A9",
		"mailto:REQ-5000+%2F@%C3%A9xample.com", "custom:REQ-5000/%25zz/%C3%A9", "mailto:REQ-5000+%25@%C3%A9xample.com", "custom:REQ-5000/%C3%A9",
	}
	if len(result.Anchors) != len(wantTexts) {
		t.Fatalf("authored/autolink count=%d want %d", len(result.Anchors), len(wantTexts))
	}
	for i, anchor := range result.Anchors {
		if anchor.Text != wantTexts[i] || anchor.Href != wantHrefs[i] || anchor.Nested != 0 || anchor.Navigation != 0 || (i >= 8 && anchor.Decorations != 0) {
			t.Errorf("anchor %d=%+v want text=%q href=%q", i, anchor, wantTexts[i], wantHrefs[i])
		}
	}
	if !strings.Contains(result.GlossaryText, longTitle) || !strings.Contains(result.GlossaryText, "pending") || !strings.Contains(result.GlossaryText, "user request") {
		t.Errorf("glossary lost full title or record status: %q", result.GlossaryText)
	}
	if result.Anchors[0].Title != longTitle {
		t.Errorf("first anchor tooltip=%q want full title=%q", result.Anchors[0].Title, longTitle)
	}
	wantGlossary := []string{"REQ-1108", "REQ-1679", "REQ-2000", "REQ-3000", "UR-075", "REQ-4000"}
	if !reflect.DeepEqual(result.Glossary, wantGlossary) || !reflect.DeepEqual(result.SecondGlossary, wantGlossary) {
		t.Errorf("glossary first=%v second=%v want %v", result.Glossary, result.SecondGlossary, wantGlossary)
	}
	// Replacing the same drawer root must not inherit a previous body's state.
	renderedBody, err := renderMarkdownBodyToHtml(body)
	if err != nil {
		t.Fatal(err)
	}
	var resetEqual bool
	session.decodeResult(t, "drawer root reuse", session.evaluateInPage(t, `(function(){var root=document.getElementById('detail-body'),before=root.innerHTML;root.innerHTML=`+mustMarshalJSONString(t, renderedBody)+`;window.authoredLinkProbe(root,'Authored references');return root.innerHTML===before;})()`), &resetEqual)
	if !resetEqual {
		t.Error("replacing drawer body reused stale expansion state")
	}
	for _, scheme := range []string{"light", "dark"} {
		for _, width := range []int{320, 768, 1280} {
			session.callDevToolsMethod(t, "Emulation.setEmulatedMedia", map[string]any{"features": []map[string]string{{"name": "prefers-color-scheme", "value": scheme}}}, true)
			session.callDevToolsMethod(t, "Emulation.setDeviceMetricsOverride", map[string]any{"width": width, "height": 900, "deviceScaleFactor": 1, "mobile": false}, true)
			var visible struct {
				Href, Browser string
				Width, Height float64
				Focusable     bool
			}
			session.decodeResult(t, "visible authored link", session.evaluateInPage(t, `(function(){var a=document.querySelector('#detail-body a'),r=a.getBoundingClientRect();a.focus();return {href:location.href,browser:navigator.userAgent,width:r.width,height:r.height,focusable:document.activeElement===a};})()`), &visible)
			if visible.Width <= 0 || visible.Height <= 0 || !visible.Focusable {
				t.Errorf("%s/%d link not visible/focusable: %+v", scheme, width, visible)
			}
		}
	}
}

func TestBuildGeneratedBoardDataCarriesDomainAndRouteProvenance(t *testing.T) {
	board := &Board{
		AllRequests: []*RequestTicket{
			{
				RequestId:          "REQ-1",
				Domain:             "backend",
				OriginalDomain:     "back-end",
				DomainUnrecognized: false,
				Route:              "A",
				OriginalRoute:      "a",
				RouteUnrecognized:  false,
			},
			{
				RequestId:          "REQ-2",
				Domain:             "general",
				OriginalDomain:     "quantum",
				DomainUnrecognized: true,
				Route:              "Z",
				OriginalRoute:      "z",
				RouteUnrecognized:  true,
			},
		},
	}

	generatedData, buildError := buildGeneratedBoardData(board)
	if buildError != nil {
		t.Fatalf("buildGeneratedBoardData: %v", buildError)
	}
	request := generatedData.Requests["REQ-1"]
	if request.OriginalDomain != "back-end" || request.DomainUnrecognized {
		t.Fatalf("domain provenance = (%q, %v), want (%q, false)",
			request.OriginalDomain, request.DomainUnrecognized, "back-end")
	}
	if request.OriginalRoute != "a" || request.RouteUnrecognized {
		t.Fatalf("route provenance = (%q, %v), want (%q, false)",
			request.OriginalRoute, request.RouteUnrecognized, "a")
	}
	invalidRequest := generatedData.Requests["REQ-2"]
	if invalidRequest.OriginalDomain != "quantum" || !invalidRequest.DomainUnrecognized {
		t.Fatalf("invalid domain provenance = (%q, %v), want (%q, true)",
			invalidRequest.OriginalDomain, invalidRequest.DomainUnrecognized, "quantum")
	}
	if invalidRequest.OriginalRoute != "z" || !invalidRequest.RouteUnrecognized {
		t.Fatalf("invalid route provenance = (%q, %v), want (%q, true)",
			invalidRequest.OriginalRoute, invalidRequest.RouteUnrecognized, "z")
	}
}

func TestDomainAndRouteProvenanceRenderAtFieldLevel(t *testing.T) {
	indexHtml := generateLiveSite(t)
	for _, requiredToken := range []string{
		"request.originalDomain || request.domain",
		"request.domainUnrecognized",
		"request.originalRoute || request.route",
		"request.routeUnrecognized",
		"schemaFieldDetailValue(request.originalDomain, request.domain",
		"schemaFieldDetailValue(request.originalRoute, request.route",
	} {
		if !strings.Contains(indexHtml, requiredToken) {
			t.Fatalf("domain/route field provenance is not rendered: %q missing", requiredToken)
		}
	}
}

func TestBuildGeneratedBoardDataCarriesFailureDetails(t *testing.T) {
	board := &Board{
		AllRequests: []*RequestTicket{
			{
				RequestId:             "REQ-1",
				Error:                 "compiler exploded",
				ErrorType:             "code",
				OriginalErrorType:     "cosmic-ray",
				ErrorTypeUnrecognized: true,
			},
			{
				RequestId: "REQ-2",
				Error:     "unclassified failure",
			},
		},
	}

	generatedData, buildError := buildGeneratedBoardData(board)
	if buildError != nil {
		t.Fatalf("buildGeneratedBoardData: %v", buildError)
	}
	invalidRequest := generatedData.Requests["REQ-1"]
	if invalidRequest.Error != "compiler exploded" || invalidRequest.ErrorType != "code" ||
		invalidRequest.OriginalErrorType != "cosmic-ray" || !invalidRequest.ErrorTypeUnrecognized {
		t.Fatalf("failure projection = (%q, %q, %q, %v), want (%q, %q, %q, true)",
			invalidRequest.Error, invalidRequest.ErrorType, invalidRequest.OriginalErrorType,
			invalidRequest.ErrorTypeUnrecognized, "compiler exploded", "code", "cosmic-ray")
	}
	unclassifiedRequest := generatedData.Requests["REQ-2"]
	if unclassifiedRequest.Error != "unclassified failure" {
		t.Fatalf("unclassified failure Error = %q, want recorded failure text", unclassifiedRequest.Error)
	}
	if unclassifiedRequest.ErrorType != "" || unclassifiedRequest.OriginalErrorType != "" || unclassifiedRequest.ErrorTypeUnrecognized {
		t.Fatalf("absent error_type projection = (%q, %q, %v), want empty/empty/false",
			unclassifiedRequest.ErrorType, unclassifiedRequest.OriginalErrorType,
			unclassifiedRequest.ErrorTypeUnrecognized)
	}
	encodedRequest, encodeError := json.Marshal(unclassifiedRequest)
	if encodeError != nil {
		t.Fatalf("marshal unclassified request: %v", encodeError)
	}
	if strings.Contains(string(encodedRequest), `"errorType"`) ||
		strings.Contains(string(encodedRequest), `"originalErrorType"`) ||
		strings.Contains(string(encodedRequest), `"errorTypeUnrecognized"`) {
		t.Fatalf("absent error_type leaked into JSON: %s", encodedRequest)
	}
}

func TestFailureDetailsRenderInTheDrawer(t *testing.T) {
	indexHtml := generateLiveSite(t)
	for _, requiredToken := range []string{
		`appendMetaRow("Error", request.error)`,
		"request.originalErrorType || request.errorType",
		"request.errorTypeUnrecognized",
		"schemaFieldDetailValue(request.originalErrorType, request.errorType",
	} {
		if !strings.Contains(indexHtml, requiredToken) {
			t.Fatalf("failure details are not rendered in the drawer: %q missing", requiredToken)
		}
	}
}

func TestDrawerDropsOnlyAMatchingLeadingHeading(t *testing.T) {
	indexHtml := generateLiveSite(t)
	for _, requiredToken := range []string{
		"function linkifyDetailBody(bodyRootElement, recordTitle)",
		"bodyRootElement.firstElementChild",
		`firstBodyElement.tagName === "H1"`,
		"normalizeHeadingText(firstBodyElement.textContent) === normalizeHeadingText(recordTitle)",
		"linkifyDetailBody(drawerBody, request.title)",
		"linkifyDetailBody(drawerBody, userRequest.title)",
	} {
		if !strings.Contains(indexHtml, requiredToken) {
			t.Fatalf("matching leading drawer-heading de-dup is not wired for REQ and UR bodies: %q missing", requiredToken)
		}
	}
}

// Execute the drawer post-processor under Node so the title de-duplication is
// pinned to behavior: matching leading H1s disappear for both record kinds,
// while a meaningful nonmatching H1 survives.
// sliceDeclarationAfter returns the source of the declaration that begins at
// anchorToken and ends at the first semicolon outside a string literal.
//
// sliceBalancedBlockAfter cannot take bodyMentionPattern: its regex source
// carries a "{0,7}" quantifier inside a string literal, and the brace counter
// reads that closing brace as the end of the block. A probe that re-declared the
// pattern beside the slice would stop testing the shipped one (REQ-322).
func sliceDeclarationAfter(t *testing.T, sourceText string, anchorToken string) string {
	t.Helper()
	anchorIndex := strings.Index(sourceText, anchorToken)
	if anchorIndex == -1 {
		t.Fatalf("anchor %q not found in the generated page", anchorToken)
	}
	var openQuoteByte byte
	for scanOffset := anchorIndex; scanOffset < len(sourceText); scanOffset++ {
		currentByte := sourceText[scanOffset]
		if openQuoteByte != 0 {
			if currentByte == '\\' {
				scanOffset++
			} else if currentByte == openQuoteByte {
				openQuoteByte = 0
			}
			continue
		}
		switch currentByte {
		case '"', '\'':
			openQuoteByte = currentByte
		case ';':
			return sourceText[anchorIndex : scanOffset+1]
		}
	}
	t.Fatalf("no terminating semicolon found after anchor %q", anchorToken)
	return ""
}

// One resolver, defined once. REQ-375 consumes the same helpers from the
// clipboard fragment, and board-core.js is the fragment that runs first — a
// second copy left behind in board-detail.js would drift the moment either side
// learns about a new id shape.
func TestTicketMentionResolverLivesOnlyInBoardCore(t *testing.T) {
	coreBytes, coreReadError := embeddedWebAssets.ReadFile("web/board-core.js")
	if coreReadError != nil {
		t.Fatalf("read web/board-core.js: %v", coreReadError)
	}
	detailBytes, detailReadError := embeddedWebAssets.ReadFile("web/board-detail.js")
	if detailReadError != nil {
		t.Fatalf("read web/board-detail.js: %v", detailReadError)
	}
	for _, sharedDefinition := range []string{
		"function buildRequestIdByReqSegment(",
		"function resolveTicketMention(",
		"function isAmbiguousTicketMention(",
		"function ticketTitleFor(",
		"function shortTicketTitle(",
	} {
		if !strings.Contains(string(coreBytes), sharedDefinition) {
			t.Errorf("web/board-core.js does not define %q", sharedDefinition)
		}
		if strings.Contains(string(detailBytes), sharedDefinition) {
			t.Errorf("web/board-detail.js still defines %q — one definition, or the two drift", sharedDefinition)
		}
	}
}

func TestDrawerTicketMentionsCarryTitlesAndAGlossary(t *testing.T) {
	indexHtml := generateLiveSite(t)
	for _, requiredToken := range []string{
		`<section class="detail-glossary" id="detail-glossary"`,
		"function makeTicketLink(detailKind, detailId, linkText, expandTitle, insideAuthoredAnchor)",
		`createElement("span", "ticket-link-id"`,
		`"ticket-link-title is-fallback" : "ticket-link-title"`,
		`createElement("span", "ticket-missing"`,
		`"Not found in this queue"`,
		"function describeTicketTitle(detailKind, detailId)",
		`"no input.md — synthesized from REQ pointers"`,
		".ticket-link-title.is-fallback,",
		"renderDetailGlossary(linkifyDetailBody(drawerBody, request.title))",
		"renderDetailGlossary(linkifyDetailBody(drawerBody, userRequest.title))",
		".ticket-link-title {",
		".ticket-missing {",
		".detail-glossary {",
	} {
		if !strings.Contains(indexHtml, requiredToken) {
			t.Errorf("title-bearing ticket mentions are not wired into the generated page: %q missing", requiredToken)
		}
	}
}

type ticketMentionNodeProbe struct {
	Tag             string   `json:"tag"`
	ClassName       string   `json:"className"`
	Text            string   `json:"text"`
	Title           string   `json:"title"`
	DetailKind      string   `json:"detailKind"`
	DetailId        string   `json:"detailId"`
	ChildClassNames []string `json:"childClassNames"`
}

type ticketGlossaryRowProbe struct {
	TermTag    string `json:"termTag"`
	Identifier string `json:"identifier"`
	DetailKind string `json:"detailKind"`
	Title      string `json:"title"`
	Status     string `json:"status"`
}

type ticketFallbackTitleProbe struct {
	Text       string `json:"text"`
	IsFallback bool   `json:"isFallback"`
}

type ticketGlossaryFallbackProbe struct {
	Id    string                   `json:"id"`
	Title ticketFallbackTitleProbe `json:"title"`
}

type ticketMentionProbeResult struct {
	ShortTitles             []string                      `json:"shortTitles"`
	CodeSpanFragment        []ticketMentionNodeProbe      `json:"codeSpanFragment"`
	InlineCodeMissing       []ticketMentionNodeProbe      `json:"inlineCodeMissingFragment"`
	FencedMissing           []ticketMentionNodeProbe      `json:"fencedMissingFragment"`
	ProseMissing            []ticketMentionNodeProbe      `json:"proseMissingFragment"`
	ProseFragment           []ticketMentionNodeProbe      `json:"proseFragment"`
	SynthesizedFragment     []ticketMentionNodeProbe      `json:"synthesizedFragment"`
	SynthesizedGlossary     []ticketGlossaryFallbackProbe `json:"synthesizedGlossaryTitles"`
	RepeatFragment          []ticketMentionNodeProbe      `json:"repeatFragment"`
	AmbiguousOnlyLinked     bool                          `json:"ambiguousOnlyLinked"`
	MetaRowLink             ticketMentionNodeProbe        `json:"metaRowLink"`
	Glossary                []ticketGlossaryRowProbe      `json:"glossary"`
	GlossaryHidden          bool                          `json:"glossaryHidden"`
	EmptyGlossaryHidden     bool                          `json:"emptyGlossaryHidden"`
	EmptyGlossaryChildCount int                           `json:"emptyGlossaryChildCount"`
}

// Execute the shipped mention pipeline under Node. The five behaviors this pins
// are the ones a title-bearing link can silently get wrong: the first prose
// mention expands and a later one does not, a code-span mention never gains
// prose (but still earns its glossary line), a dead id is FLAGGED while an
// ambiguous segment is LEFT ALONE, and the glossary lists each resolved id once
// with its untruncated title.
// referencedRequestsGlossaryHeading is the appendix heading REQ-379 specifies
// for the clipboard payload. Pinned as a literal in the tests rather than read
// back off the client: it is the sentence a paste's reader sees, so a reworded
// heading is a behavior change and must fail, not follow the assertions along.
const referencedRequestsGlossaryHeading = "## Referenced requests (added by the board — not part of the file)"

type clipboardAnnotationProbeResult struct {
	AnnotatedHostDocument string   `json:"annotatedHostDocument"`
	HostReferencedIds     []string `json:"hostReferencedIds"`
	JoinedPayload         string   `json:"joinedPayload"`
	GlossaryHeadingCount  int      `json:"glossaryHeadingCount"`
	ExcludedPayload       string   `json:"excludedPayload"`
	UnclosedFencePayload  string   `json:"unclosedFencePayload"`
	CarriageReturnPayload string   `json:"carriageReturnPayload"`
	FencelessPayload      string   `json:"fencelessPayload"`
	LoneFencePayload      string   `json:"loneFencePayload"`
	AmbiguousPayload      string   `json:"ambiguousPayload"`
	NoReferencePayload    string   `json:"noReferencePayload"`
	BlockquotedFence      string   `json:"blockquotedFencePayload"`
	InvalidInfoString     string   `json:"invalidInfoStringPayload"`
	ListItemFence         string   `json:"listItemFencePayload"`
	MultiLineCodeSpan     string   `json:"multiLineCodeSpanPayload"`
}

// Execute the shipped clipboard annotator under Node. Each case below names a
// way a payload that passes a looser assertion still pastes as a broken file:
// an annotated frontmatter fence stops parsing as YAML, a second document's
// fence annotated after a join is the same failure one ticket later, an
// unclosed fence read as "all fence" silently skips a whole document, a fenced
// block flagged as a dead reference asserts something false about an
// illustration, and a repeat mention expanded twice turns prose into noise.
// The three Copy handlers each hand their own payload and their own exclusion
// set to one annotator. A handler wired to the wrong exclusion set glosses the
// tickets it already contains, which no pure-function probe can see.
func TestClipboardAnnotationWiresEveryCopyHandler(t *testing.T) {
	indexHtml := generateLiveSite(t)
	if !strings.Contains(indexHtml, "function annotateClipboardPayload(") {
		t.Fatal("the generated page defines no annotateClipboardPayload — every assertion below is vacuous")
	}
	if callSiteCount := strings.Count(indexHtml, "return annotateClipboardPayload("); callSiteCount != 3 {
		t.Errorf("annotateClipboardPayload call sites = %d, want one per Copy handler", callSiteCount)
	}
	for _, requiredCallSite := range []string{
		// The drawer Copy annotates its RAW branch only. Its fallback input is
		// drawerBody.innerText, which the drawer already expanded, so annotating
		// that again duplicated every title ("REQ-1679 (Short one) Short one").
		"return annotateClipboardPayload([clipboardDocument], [requestedId]);",
		"[requestedUserRequestId].concat(requestedRequestIds)",
		"return annotateClipboardPayload(clipboardDocumentsForRequests(markdownData, requestIds), requestIds);",
	} {
		if !strings.Contains(indexHtml, requiredCallSite) {
			t.Errorf("Copy handler wiring missing: %q", requiredCallSite)
		}
	}
	if !strings.Contains(indexHtml, referencedRequestsGlossaryHeading) {
		t.Errorf("the generated page does not carry the glossary heading %q", referencedRequestsGlossaryHeading)
	}
}

// The client carries no Markdown scanner any more. Each name below was a piece
// of one, and each got a CommonMark rule wrong: fence recognition, the
// closing-fence length rule, inline code-span matching, container prefixes, and
// the frontmatter split. Asserting their absence is what makes the deletion a
// fact rather than a claim — a rewrite that left one behind for "just this case"
// is how a scanner grows back.
//
// The positive half is not decoration: without it every absence below would also
// pass if the whole Copy feature were deleted.
func TestClipboardCarriesNoMarkdownScanner(t *testing.T) {
	indexHtml := generateLiveSite(t)

	for _, scannerPiece := range []string{
		"codeFenceRunFor",
		"codeFenceRunCloses",
		"findMatchingBacktickRun",
		"stripContainerPrefix",
		"frontmatterFenceEndOffset",
		"annotateLineOutsideFence",
		"annotateMarkdownBody",
		"annotateMentionRun",
	} {
		if strings.Contains(indexHtml, scannerPiece) {
			t.Errorf("the generated page still carries %s — Markdown structure is decided in citations.go now", scannerPiece)
		}
	}

	for _, splicerPiece := range []string{
		"function annotateTicketMentions(documentText, ticketMentions)",
		"function ticketMentionsForDetail(",
		"function clipboardDocumentFor(",
		"markdownData.userRequestMentions : markdownData.requestMentions",
	} {
		if !strings.Contains(indexHtml, splicerPiece) {
			t.Errorf("the generated page is missing %q — the absences above would pass on an empty page", splicerPiece)
		}
	}
}

// mustMarshalJSONString renders a Go string as a JavaScript string literal for
// splicing into a probe, so a fixture title carrying quotes or a non-ASCII
// character cannot break the probe's syntax.
func mustMarshalJSONString(t *testing.T, plainText string) string {
	t.Helper()
	encoded, encodeError := json.Marshal(plainText)
	if encodeError != nil {
		t.Fatalf("encode probe string %q: %v", plainText, encodeError)
	}
	return string(encoded)
}

// The Copy payload must be the ticket file exactly as it exists on disk —
// frontmatter fence included — so a paste can be saved straight back as a valid
// REQ or UR file. Parsed from real files rather than hand-built structs, because
// the whole point is that the ORIGINAL bytes survive: a reconstructed fence would
// pass a struct-level assertion while losing key order, comments, and line
// endings.
func TestBuildGeneratedBoardMarkdownDataRoundTripsTheWholeFile(t *testing.T) {
	fixtureDirectory := t.TempDir()

	requestPath := filepath.Join(fixtureDirectory, "REQ-4242-round-trip.md")
	requestFileText := "---\nid: REQ-4242\nstatus: pending\ntitle: round trip\n" +
		"# a comment the fence must keep\ndomain:   general\n---\n\n## What\n\n- [ ] keep formatting\n"
	if writeError := os.WriteFile(requestPath, []byte(requestFileText), 0o644); writeError != nil {
		t.Fatalf("write REQ fixture: %v", writeError)
	}

	userRequestPath := filepath.Join(fixtureDirectory, "UR-4242-input.md")
	userRequestFileText := "---\nid: UR-4242\ntitle: the original ask\nrequests: [REQ-4242]\n---\n\n# Original request\n\nExact text.\n"
	if writeError := os.WriteFile(userRequestPath, []byte(userRequestFileText), 0o644); writeError != nil {
		t.Fatalf("write UR fixture: %v", writeError)
	}

	parsedRequest, requestParseError := parseRequestTicket(requestPath, "queue")
	if requestParseError != nil {
		t.Fatalf("parseRequestTicket: %v", requestParseError)
	}
	parsedUserRequest, userRequestParseError := parseUserRequestTicket(userRequestPath)
	if userRequestParseError != nil {
		t.Fatalf("parseUserRequestTicket: %v", userRequestParseError)
	}

	board := &Board{
		AllRequests:  []*RequestTicket{parsedRequest},
		UserRequests: []*UserRequestTicket{parsedUserRequest},
	}
	markdownData := buildGeneratedBoardMarkdownData(board)

	if got := markdownData.Requests["REQ-4242"]; got != requestFileText {
		t.Errorf("REQ Copy payload is not the file on disk:\n got: %q\nwant: %q", got, requestFileText)
	}
	if got := markdownData.UserRequests["UR-4242"]; got != userRequestFileText {
		t.Errorf("UR Copy payload is not the file on disk:\n got: %q\nwant: %q", got, userRequestFileText)
	}
}

// A file with no frontmatter at all must still copy as itself — the fence field
// is empty, not a fabricated one.
func TestBuildGeneratedBoardMarkdownDataHandlesAFenceLessFile(t *testing.T) {
	fixtureDirectory := t.TempDir()
	requestPath := filepath.Join(fixtureDirectory, "REQ-4243-no-fence.md")
	requestFileText := "# REQ-4243\n\nA legacy file with no frontmatter.\n"
	if writeError := os.WriteFile(requestPath, []byte(requestFileText), 0o644); writeError != nil {
		t.Fatalf("write fixture: %v", writeError)
	}

	parsedRequest, parseError := parseRequestTicket(requestPath, "queue")
	if parseError != nil {
		t.Fatalf("parseRequestTicket: %v", parseError)
	}
	if parsedRequest.FrontmatterMarkdown != "" {
		t.Errorf("a file with no frontmatter must yield an empty fence, got %q", parsedRequest.FrontmatterMarkdown)
	}

	board := &Board{AllRequests: []*RequestTicket{parsedRequest}}
	markdownData := buildGeneratedBoardMarkdownData(board)
	if got := markdownData.Requests[parsedRequest.RequestId]; got != requestFileText {
		t.Errorf("fence-less Copy payload changed:\n got: %q\nwant: %q", got, requestFileText)
	}
}

// A synthesized UR node (a REQ points at a UR whose input.md was never found)
// must NOT get a markdown-map entry: the frontend reads key presence as "the
// real file is available" and copies the value verbatim, so an empty entry
// makes the drawer's Copy button write an empty string instead of falling back
// to the rendered text with its identifying heading.
func TestBuildGeneratedBoardMarkdownDataOmitsSynthesizedUserRequests(t *testing.T) {
	board := &Board{
		UserRequests: []*UserRequestTicket{
			{UserRequestId: "UR-7", InputFilePresent: true, BodyMarkdown: "# Real request\n"},
			{UserRequestId: "UR-8", InputFilePresent: false},
		},
	}

	markdownData := buildGeneratedBoardMarkdownData(board)
	if _, exists := markdownData.UserRequests["UR-7"]; !exists {
		t.Errorf("a UR with a real input.md must keep its markdown-map entry")
	}
	if _, exists := markdownData.UserRequests["UR-8"]; exists {
		t.Errorf("a synthesized UR must have NO markdown-map entry — key presence sends the frontend down the verbatim-copy path with an empty payload")
	}
}

func TestRenderMarkdownBodyToHtmlHeadingsAndTaskLists(t *testing.T) {
	body := "## What\n\nA paragraph.\n\n- [ ] unchecked item\n- [x] checked item\n"
	rendered, renderError := renderMarkdownBodyToHtml(body)
	if renderError != nil {
		t.Fatalf("renderMarkdownBodyToHtml: %v", renderError)
	}
	if !strings.Contains(rendered, "<h2") {
		t.Fatalf("expected an <h2> from a ## heading, got: %s", rendered)
	}
	if !strings.Contains(rendered, `type="checkbox"`) {
		t.Fatalf("expected GFM task-list checkboxes, got: %s", rendered)
	}
}

func TestRenderMarkdownQuestionOptionsKeepTheirOwnLines(t *testing.T) {
	// The Open Questions format (actions/capture.md) indents Recommended:/Also:
	// continuation lines under the checkbox item; plain Markdown would lazily
	// merge them into the question paragraph. The renderer must emit a <br>
	// before each so they stay separate visual lines in the drawer.
	body := "## Open Questions\n\n" +
		"- [ ] Should I process this as a new task?\n" +
		"  Recommended: Yes, add to queue.\n" +
		"  Also: No, discard it.\n"
	rendered, renderError := renderMarkdownBodyToHtml(body)
	if renderError != nil {
		t.Fatalf("renderMarkdownBodyToHtml: %v", renderError)
	}
	if strings.Count(rendered, "<br") != 2 {
		t.Fatalf("expected 2 hard breaks (before Recommended: and Also:), got: %s", rendered)
	}
	if !strings.Contains(rendered, `type="checkbox"`) {
		t.Fatalf("checkbox item must survive the option-line preprocessing, got: %s", rendered)
	}
}

func TestRenderMarkdownInvalidBacktickInfoRemainsQuestionProse(t *testing.T) {
	body := "## Open Questions\n\n" +
		"- [ ] Should I process this as a new task?\n" +
		"  ```lang`invalid\n" +
		"  Recommended: Yes, add to queue.\n" +
		"  Also: No, discard it.\n"
	rendered, renderError := renderMarkdownBodyToHtml(body)
	if renderError != nil {
		t.Fatalf("renderMarkdownBodyToHtml: %v", renderError)
	}
	if strings.Count(rendered, "<br") != 2 {
		t.Fatalf("invalid backtick-info prose must keep 2 hard breaks (before Recommended: and Also:), got: %s", rendered)
	}
}

func TestRenderMarkdownLeavesValidCodeFencesVerbatim(t *testing.T) {
	// A fenced block whose content happens to start with an option keyword must
	// not have hard-break backslashes injected into its verbatim content.
	testCases := map[string]string{
		"backtick fence": "```go\nsome output\nRecommended: not a question option\n```\n",
		"tilde fence":    "~~~lang`valid\nsome output\nRecommended: not a question option\n~~~\n",
	}
	for testName, body := range testCases {
		t.Run(testName, func(t *testing.T) {
			if preprocessedBody := insertQuestionOptionHardBreaks(body); preprocessedBody != body {
				t.Fatalf("valid fence changed during preprocessing:\n got: %q\nwant: %q", preprocessedBody, body)
			}
			rendered, renderError := renderMarkdownBodyToHtml(body)
			if renderError != nil {
				t.Fatalf("renderMarkdownBodyToHtml: %v", renderError)
			}
			if strings.Contains(rendered, "\\") || strings.Contains(rendered, "<br") {
				t.Fatalf("code fence content must stay verbatim, got: %s", rendered)
			}
		})
	}
}

func TestRenderMarkdownEscapesRawHtml(t *testing.T) {
	rendered, renderError := renderMarkdownBodyToHtml("a <script>alert(1)</script> b")
	if renderError != nil {
		t.Fatalf("renderMarkdownBodyToHtml: %v", renderError)
	}
	if strings.Contains(rendered, "<script>") {
		t.Fatalf("raw <script> should be escaped, got: %s", rendered)
	}
}

// TestEncodeBoardDataJsAssignmentPreservesRawHtml covers the one encoder both
// generate and serve actually ship (board-data.js is a plain .js file, never
// HTML-parsed, so no </script> neutralization is involved): the assignment
// wrapper must be exact and pre-rendered body HTML must survive unescaped
// (SetEscapeHTML off — the goldmark proof the GREEN test greps for).
func TestEncodeBoardDataJsAssignmentPreservesRawHtml(t *testing.T) {
	data := generatedBoardData{
		Requests: map[string]generatedRequest{
			"REQ-1": {RequestId: "REQ-1", BodyHtml: "<h2>Lessons & Notes</h2>"},
		},
	}
	encoded, encodeError := encodeBoardDataForJsAssignment(data)
	if encodeError != nil {
		t.Fatalf("encodeBoardDataForJsAssignment: %v", encodeError)
	}
	if !strings.HasPrefix(encoded, "window.queueKanbanBoardData = ") {
		t.Fatalf("expected the window.queueKanbanBoardData assignment prefix: %s", encoded)
	}
	if !strings.HasSuffix(encoded, ";\n") {
		t.Fatalf("expected the assignment to end with a semicolon + newline: %s", encoded)
	}
	if !strings.Contains(encoded, "<h2>Lessons & Notes</h2>") {
		t.Fatalf("expected pre-rendered HTML to survive verbatim (HTML escaping off): %s", encoded)
	}
	escapedLessThan := "\\u003c"
	escapedAmpersand := "\\u0026"
	if strings.Contains(encoded, escapedLessThan) || strings.Contains(encoded, escapedAmpersand) {
		t.Fatalf("body HTML was unicode-escaped by the JSON encoder: %s", encoded)
	}
}

func TestEncodeBoardMarkdownJsAssignmentRoundTripsRawSource(t *testing.T) {
	want := generatedBoardMarkdownData{
		Requests:     map[string]string{"REQ-1": "## What\n\nA <literal> & text.\n"},
		UserRequests: map[string]string{"UR-1": "# Ask\n\nCopy me.\n"},
	}
	encoded, encodeError := encodeBoardMarkdownForJsAssignment(want)
	if encodeError != nil {
		t.Fatalf("encodeBoardMarkdownForJsAssignment: %v", encodeError)
	}

	const prefix = "window.queueKanbanBoardMarkdownData = "
	if !strings.HasPrefix(encoded, prefix) || !strings.HasSuffix(encoded, ";\n") {
		t.Fatalf("unexpected lazy Markdown assignment envelope: %q", encoded)
	}
	jsonText := strings.TrimSuffix(strings.TrimPrefix(encoded, prefix), ";\n")
	var got generatedBoardMarkdownData
	if decodeError := json.Unmarshal([]byte(jsonText), &got); decodeError != nil {
		t.Fatalf("decode lazy Markdown assignment: %v", decodeError)
	}
	if got.Requests["REQ-1"] != want.Requests["REQ-1"] || got.UserRequests["UR-1"] != want.UserRequests["UR-1"] {
		t.Fatalf("raw Markdown did not round-trip: got %#v, want %#v", got, want)
	}
}

// TestRecentlyDoneWindowDefaultsTo24h asserts that a fresh board load defaults
// the RECENTLY DONE column to the 24h window: the 24h toggle button must carry
// aria-pressed="true" and the 7d (168h) button must NOT be the default-active one.
// The assertion also verifies that the inlined board.js shell initialises windowHours to
// 24, not 168, so the JS runtime agrees with the HTML button state on load.
func TestRecentlyDoneWindowDefaultsTo24h(t *testing.T) {
	indexHtml := generateLiveSite(t)

	// The 24h button must be the active one on load.
	activeMarker24h := `data-window-hours="24" aria-pressed="true"`
	if !strings.Contains(indexHtml, activeMarker24h) {
		t.Fatalf("24h window button is not the default-active toggle: expected %q in the generated page", activeMarker24h)
	}

	// The 7d button must NOT carry aria-pressed="true" (it is the old default).
	staleActive7d := `data-window-hours="168" aria-pressed="true"`
	if strings.Contains(indexHtml, staleActive7d) {
		t.Fatalf("7d window button is still marked as the default-active toggle: %q must not appear in the generated page", staleActive7d)
	}

	// The inlined board.js shell default must match the HTML button state.
	jsDefaultWindow24h := "windowHours: 24"
	if !strings.Contains(indexHtml, jsDefaultWindow24h) {
		t.Fatalf("board.js windowHours default is not 24: expected %q in the inlined script", jsDefaultWindow24h)
	}
	jsDefaultWindow168 := "windowHours: 168"
	if strings.Contains(indexHtml, jsDefaultWindow168) {
		t.Fatalf("board.js still initialises windowHours to 168: %q must not appear in the inlined script", jsDefaultWindow168)
	}
}

// TestGenerateInlinesWriteSetOverlapBadgeRenderPath guards the frontend half of
// the write_set overlap annotation. The Go tests cover annotateWriteSetOverlap
// (model_test.go), but nothing proved the derived list still gets *rendered*:
// a refactor that dropped the badge renderer from web/board-cards.js would ship a
// silent regression, since the badge only appears when the live tree happens to
// have overlapping REQs. These are code tokens from the assembled client/board.css,
// so the assertion holds regardless of what the queue currently contains.
func TestGenerateInlinesWriteSetOverlapBadgeRenderPath(t *testing.T) {
	indexHtml := generateLiveSite(t)

	for _, renderToken := range []string{
		// The makeBadge() call that emits the card badge. The quoted form only
		// occurs in board-cards.js — the bare class name would also match the CSS rule.
		`"badge-write-overlap"`,
		// The generated payload field (generate.go's writeSetOverlaps key) that
		// the badge gates on.
		"request.writeSetOverlaps",
		// The drawer row that makes the contending REQ ids clickable.
		"Overlapping write sets",
		// Without the stylesheet rule the badge renders unstyled and invisible.
		".badge-write-overlap",
	} {
		if !strings.Contains(indexHtml, renderToken) {
			t.Fatalf("write_set overlap badge render path is missing from the generated page: %q not found in the assembled client/board.css", renderToken)
		}
	}
}

// TestGenerateInlinesSemanticStatusCardStyles pins the visual status contract
// at the generated-site seam shared by static and live boards. The By-UR lens
// renders the same request cards as the column lens, so card-level status
// variables must win without a lens-specific neutral override.
func TestGenerateInlinesSemanticStatusCardStyles(t *testing.T) {
	indexHtml := generateLiveSite(t)

	if strings.Contains(indexHtml, `.ur-group-cards .req-card`) {
		t.Fatalf("By-UR still overrides request-card status styling with a neutral accent")
	}
	if strings.Contains(indexHtml, `border-left: 2px solid var(--card-accent)`) {
		t.Fatalf("request-card status rail is still 2px instead of the required 3px")
	}
	if !strings.Contains(indexHtml, `border-left: 3px solid var(--card-accent)`) {
		t.Fatalf("request-card 3px status rail is missing from the inlined stylesheet")
	}

	statusStyleCases := []struct {
		name              string
		selectorAnchor    string
		requiredAccent    string
		requiredTint      string
		requiredSelectors []string
	}{
		{
			name:           "pending",
			selectorAnchor: `.req-card[data-status="pending"]`,
			requiredAccent: `--card-accent: var(--accent-pending)`,
			requiredTint:   `--card-tint: var(--tint-pending)`,
		},
		{
			name:           "claimed",
			selectorAnchor: `.req-card[data-status="claimed"]`,
			requiredAccent: `--card-accent: var(--accent-claimed)`,
			requiredTint:   `--card-tint: var(--tint-claimed)`,
		},
		{
			name:           "blocked and failed",
			selectorAnchor: `.req-card[data-status="pending-answers"]`,
			requiredAccent: `--card-accent: var(--accent-blocked)`,
			requiredTint:   `--card-tint: var(--tint-blocked)`,
			requiredSelectors: []string{
				`.req-card[data-status="pending-heavy-testing"]`,
				`.req-card[data-status="blocked"]`,
				`.req-card[data-status="blocked-archive-collision"]`,
				`.req-card[data-status="blocked-dependency-cycle"]`,
				`.req-card[data-status="failed"]`,
			},
		},
		{
			name:           "completed",
			selectorAnchor: `.req-card[data-status="completed"]`,
			requiredAccent: `--card-accent: var(--accent-done)`,
			requiredTint:   `--card-tint: var(--tint-done)`,
			requiredSelectors: []string{
				`.req-card[data-status="completed-with-issues"]`,
			},
		},
		{
			name:           "cancelled",
			selectorAnchor: `.req-card[data-status="cancelled"]`,
			requiredAccent: `--card-accent: var(--ink-faint)`,
			requiredTint:   `--card-tint: var(--surface-3)`,
		},
		{
			name:           "unrecognized",
			selectorAnchor: `.req-card.is-status-unrecognized`,
			requiredAccent: `--card-accent: var(--accent-blocked)`,
			requiredTint:   `--card-tint: var(--tint-blocked)`,
		},
	}
	for _, styleCase := range statusStyleCases {
		t.Run(styleCase.name, func(t *testing.T) {
			styleRule := sliceBalancedBlockAfter(t, indexHtml, styleCase.selectorAnchor)
			for _, requiredToken := range append(
				[]string{styleCase.requiredAccent, styleCase.requiredTint},
				styleCase.requiredSelectors...,
			) {
				if !strings.Contains(styleRule, requiredToken) {
					t.Fatalf("status style rule %q is missing %q:\n%s", styleCase.selectorAnchor, requiredToken, styleRule)
				}
			}
		})
	}

	statusPillRule := sliceBalancedBlockAfter(t, indexHtml, `.req-card-status {`)
	for _, requiredToken := range []string{
		`background-color: var(--card-tint)`,
		`border-radius: var(--radius-pill)`,
		`color: var(--ink-soft)`,
	} {
		if !strings.Contains(statusPillRule, requiredToken) {
			t.Fatalf("status pill is missing %q:\n%s", requiredToken, statusPillRule)
		}
	}

	if !strings.Contains(indexHtml, `card.className += " is-status-unrecognized"`) {
		t.Fatalf("unrecognized status does not mark the request card for red rail and pill styling")
	}
}

// sliceBalancedBlockAfter returns the source text of the first brace-balanced
// block that starts at or after anchorToken, including the anchor itself. It
// brace-matches rather than scanning to a blank line so the slice stays exact
// under reformatting; the blocks it is pointed at contain no braces inside
// string literals, which is what makes the naive counter safe here.
func sliceBalancedBlockAfter(t *testing.T, sourceText string, anchorToken string) string {
	t.Helper()
	anchorIndex := strings.Index(sourceText, anchorToken)
	if anchorIndex == -1 {
		t.Fatalf("anchor %q not found in the generated page", anchorToken)
	}
	braceDepth := 0
	sawOpeningBrace := false
	for scanOffset := anchorIndex; scanOffset < len(sourceText); scanOffset++ {
		switch sourceText[scanOffset] {
		case '{':
			braceDepth++
			sawOpeningBrace = true
		case '}':
			braceDepth--
			if sawOpeningBrace && braceDepth == 0 {
				return sourceText[anchorIndex : scanOffset+1]
			}
		}
	}
	t.Fatalf("no brace-balanced block found after anchor %q", anchorToken)
	return ""
}

// The Active by-UR lens and Recently done column share one wall-clock window.
// Execute both production predicates so terminal work crosses the scope boundary
// when it ages out, while non-terminal work remains active.
func TestGenerateOffersDurationsWindowControls(t *testing.T) {
	indexHTML := generateLiveSite(t)
	for _, requiredToken := range []string{
		`id="durations-window-group" hidden`,
		`data-durations-window="30" aria-pressed="true"`,
		`data-durations-window="90" aria-pressed="false"`,
		`data-durations-window="all" aria-pressed="false"`,
		`document.getElementById("durations-window-group").hidden = viewState.view !== "durations"`,
	} {
		if !strings.Contains(indexHTML, requiredToken) {
			t.Fatalf("generated board is missing Durations-window contract %q", requiredToken)
		}
	}
	if strings.Count(indexHTML, `data-durations-window=`) != 3 {
		t.Fatalf("generated board has %d Durations-window choices, want exactly 3", strings.Count(indexHTML, `data-durations-window=`))
	}
	if strings.Contains(indexHTML, `data-durations-window="30" data-window-hours=`) {
		t.Fatal("Durations window control is coupled to the board's recently-done window attribute")
	}
}

// Panel A's overflow lane keeps every mark, while the adjacent ranked list
// carries the complete text record for spans whose y position is capped. The
// list must be ordinary HTML rather than SVG text so density changes scrolling,
// not which samples remain reachable.
func TestGenerateCarriesAdjacentCompleteDurationsLongestSpansList(t *testing.T) {
	indexHTML := generateLiveSite(t)

	for _, requiredToken := range []string{
		`class="durations-chart-list"`,
		`id="durations-chart"`,
		`class="durations-longest-spans"`,
		`id="durations-longest-count"`,
		`id="durations-longest-list"`,
		`function renderDurationsLongestSpans(`,
	} {
		if !strings.Contains(indexHTML, requiredToken) {
			t.Errorf("generated board is missing the complete longest-spans contract %q", requiredToken)
		}
	}
	chartIndex := strings.Index(indexHTML, `id="durations-chart"`)
	listIndex := strings.Index(indexHTML, `id="durations-longest-list"`)
	if chartIndex < 0 || listIndex < 0 || listIndex < chartIndex {
		t.Errorf("longest-spans list must follow the chart inside their shared wrapper (chart=%d list=%d)", chartIndex, listIndex)
	}
}

// Execute the pure empty-state decision under Node so the regression is pinned
// to state transitions rather than the presence of reassuring source strings.
// Exercise the production lens caller, not only its pure predicate: a terminal
// REQ inside the selected window renders while an older terminal REQ stays
// hidden and is counted in the scope note.
// TestUserRequestActivityToggleDocumentsWidenedRule pins the template half: the
// Active chip must explain the widened rule on hover, because "Active" alone no
// longer means what a reader would assume it means.
func TestUserRequestActivityToggleDocumentsWidenedRule(t *testing.T) {
	indexHtml := generateLiveSite(t)

	if !strings.Contains(indexHtml, `data-ur-activity="active"`) {
		t.Fatalf("the Active user-request scope button is missing from the generated page")
	}
	activeButtonSource := indexHtml[strings.Index(indexHtml, `data-ur-activity="active"`):]
	if closingIndex := strings.Index(activeButtonSource, "</button>"); closingIndex != -1 {
		activeButtonSource = activeButtonSource[:closingIndex]
	}
	if !strings.Contains(activeButtonSource, "title=") {
		t.Fatalf("the Active scope button has no title explaining the widened rule: %s", activeButtonSource)
	}
	if !strings.Contains(activeButtonSource, "RECENTLY DONE window") {
		t.Fatalf("the Active scope button's title does not mention the RECENTLY DONE window: %s", activeButtonSource)
	}

	// Two comments restated the old rule in prose — one above the template's
	// control group, one on the filterState declaration in board.js. Renaming the
	// predicate did not touch either, so a grep for the identifier missed them
	// both; this substring is the phrasing they shared. Per the skill's
	// closed-enumerations rule the fix was to point at the predicate as the
	// canonical statement rather than to re-copy the widened rule a third time.
	if strings.Contains(indexHtml, "whose REQs are all resolved") {
		t.Fatalf("a stale prose restatement of the by-UR Active rule is still present in the generated page")
	}
}

// ---- panel B's slowest-day annotation, and the faces around it ------------
//
// Both faces below are the browser's answer, because a face is the browser's
// answer and no arithmetic over the constants under test can produce it. Panel
// B's title draws at .durations-axis-title (12px, weight 600) and the
// annotation at .durations-mark-label (11px). Every number is rounded UP, which
// GROWS the boxes the assertion holds apart, so a pass can never be an artefact
// of the rounding.
//
// Procedure, reproducible from any board directory `queue-kanban generate`
// wrote: load index.html, activate the Durations view, and read getBBox()
// against the node's own `y` on the two <text> nodes, at 1400x1200 — the SVG is
// a fixed viewBox at width:100%, so user units are zoom- and window-independent.
//
// A measured face is PER-BROWSER, so each constant's own doc comment names the
// Chromium build its number was taken on — TestMeasuredFaceConstantsNameTheirBuild
// (below) enforces that for every durationsMeasured constant in the package, and
// its comment records both the collision that earned the rule and why the
// enforcer lives here now. A re-measurement on another build may only
// RAISE a constant, never lower it: a box that reaches further makes every
// clearance test demand more room than the render needs, which is the safe
// direction for every caller.

// The 12px title face's ascent, rounded up from the 12.0372 REQ-242 measured on
// Chromium 146 (Playwright 1.59). This constant is declared here for the whole
// package rather than once per test file because REQ-241's clearance assertion
// measured the same face at 11.2300 on its own Chromium (recorded only as
// browser build chromium-1228) and declared its own constant; the two REQs
// merged into one package, disagreed by 0.86 units, and failed to compile. 12.1
// is the larger, so it is the one kept. Chromium 141.0.7390.37 (Playwright
// 1.56.1, REQ-252) measures 11.1112 — smaller, so the constant stands.
const durationsMeasuredAxisTitleAscentUnits = 12.1

// The same title box's descent, rounded up from the 2.7778 REQ-242 measured on
// Chromium 146 (Playwright 1.59); Chromium 141.0.7390.37 (REQ-252) measures the
// same 2.7778.
const durationsMeasuredAxisTitleDescentUnits = 2.8

// The 11px annotation face's ascent, rounded to 10.5 so it covers both the
// 10.1853 REQ-242 measured on Chromium 146 (Playwright 1.59) and the 10.4278
// the same face measured for REQ-241 on its chromium-1228 build; Chromium
// 141.0.7390.37 (REQ-252) measures 10.1853.
const durationsMeasuredMarkLabelAscentUnits = 10.5

// The annotation box's descent, rounded up from the 2.7778 REQ-242 measured on
// Chromium 146 (Playwright 1.59); Chromium 141.0.7390.37 (REQ-252) measures the
// same 2.7778.
const durationsMeasuredMarkLabelDescentUnits = 2.8

// durationsAnnotationCase is one (position, bar height) the annotation can be
// asked for — the two dimensions that decided whether the original defect was
// visible.
type durationsAnnotationCase struct {
	caseName      string
	dayCentreX    float64
	medianMinutes float64
}

// The named extremes. Each is a position or a height this repository's own data
// has never produced, and the first is the one the defect hid behind — the
// annotation is centred on its day, so its x follows the slowest day wherever
// that lands.
var durationsAnnotationNamedExtremes = []durationsAnnotationCase{
	{"leftmost day, over the ceiling", 54, 209},
	{"leftmost day, a hair left of the plot", 38.9, 209},
	{"leftmost day, at the ceiling", 54, 45},
	{"leftmost day, a flat bar", 54, 0},
	{"mid-plot day, over the ceiling", 618, 209},
	{"rightmost day, over the ceiling", 1182, 209},
}

// Sweep bounds. x runs WIDER than the plot on purpose: the day-anchored domain
// (REQ-248) keeps every real day centre inside it, but this sweep is a property
// of the drawing function at ANY x it could ever be handed, so it still covers
// the off-plot band the pre-REQ-248 renderer produced. The median range covers
// everything under the four-hour read-time ceiling that admits a sample at all.
const (
	durationsAnnotationSweepLeftX      = -400.0
	durationsAnnotationSweepRightX     = 1400.0
	durationsAnnotationSweepMaxMinutes = 240.0
	durationsAnnotationSweepXCount     = 200
	durationsAnnotationSweepMedians    = 50
)

// durationsSlowestDayAnnotationCaseList is the named extremes plus a
// deterministic pseudo-random sweep of the whole plane the annotation can be
// asked for, crossed so that a rule keyed on x AND median together has nowhere
// to hide either.
//
// Why a sweep and not six points: six points is a SAMPLE, and the property
// under test is a claim about every x. A mutant that returned the original
// defect's baseline for x in (700, 1100) — which is where the maintainer's own
// slowest day sits, at x=881.4 — passed the six-point version while
// reproducing the defect on the real board. So did one banded on the median
// between 1 and 44, where the real board's 42-minute median lives. Ten thousand
// pairs cost about a third of a second; the structural check below closes the
// gap the sweep still leaves for an arbitrarily narrow band.
func durationsSlowestDayAnnotationCaseList() []durationsAnnotationCase {
	annotationCases := append([]durationsAnnotationCase(nil), durationsAnnotationNamedExtremes...)
	// A fixed seed, so a failure names coordinates the next run reproduces.
	randomState := uint64(20260818)
	nextUnitInterval := func() float64 {
		randomState = randomState*6364136223846793005 + 1442695040888963407
		return float64(randomState>>11) / float64(uint64(1)<<53)
	}
	// One decimal place, because the renderer writes x through toFixed(1) and a
	// case the probe cannot echo back exactly would test the rounding instead.
	sweptPositions := make([]float64, durationsAnnotationSweepXCount)
	for positionIndex := range sweptPositions {
		spread := durationsAnnotationSweepRightX - durationsAnnotationSweepLeftX
		sweptPositions[positionIndex] = math.Round((durationsAnnotationSweepLeftX+nextUnitInterval()*spread)*10) / 10
	}
	sweptMedians := make([]float64, durationsAnnotationSweepMedians)
	for medianIndex := range sweptMedians {
		sweptMedians[medianIndex] = math.Round(nextUnitInterval()*durationsAnnotationSweepMaxMinutes*100) / 100
	}
	for _, sweptX := range sweptPositions {
		for _, sweptMedian := range sweptMedians {
			annotationCases = append(annotationCases, durationsAnnotationCase{
				caseName:      fmt.Sprintf("swept day at x=%.1f with a %.2f-minute median", sweptX, sweptMedian),
				dayCentreX:    sweptX,
				medianMinutes: sweptMedian,
			})
		}
	}
	return annotationCases
}

// The defect this pins (REQ-242): the slowest-day annotation was drawn 7 units
// above its bar's top, which for an over-ceiling day is a baseline of 355
// against panel B's own title baseline of 350 — two boxes sharing eight units of
// vertical space, colliding at every x where their text ranges met. It never
// showed here because this repository's slowest day falls to the right of where
// the title's text ends, which is luck and not design.
//
// The annotation now sits below panel B's baseline, so the clearance is the same
// for every x and every bar height — which is what the assertion states by
// driving the renderer's own drawing function across ten thousand positions and
// heights and demanding one unchanging baseline that clears every neighbour that
// strip has: panel B's title, panel B's "0" axis tick, panel C's title, and the
// plot the bars occupy. The strip's fourth occupant, the month rule, cannot be
// cleared by any baseline and is handled as an accepted crossing at the end.
// The sweep above is dense, but it is still a sample: a band narrower than its
// spacing would slip through it. This closes that gap exactly, by reading the
// shipped function and requiring its baseline expression to mention neither
// input. That is what makes one measurement at one x a statement about every x,
// which is the whole claim this REQ rests on.
func assertDurationsAnnotationBaselineIgnoresItsInputs(t *testing.T, annotationSource string) {
	t.Helper()
	baselineExpression := ""
	for _, sourceLine := range strings.Split(annotationSource, "\n") {
		if strings.HasPrefix(strings.TrimSpace(sourceLine), "y:") {
			baselineExpression = strings.TrimSpace(sourceLine)
			break
		}
	}
	if baselineExpression == "" {
		t.Fatal("drawDurationsSlowestDayAnnotation declares no single-line `y:` baseline — if the expression became multi-line, restate this check against its new shape rather than deleting it")
	}
	for _, inputName := range []string{"dayCentreX", "medianMinutes"} {
		if strings.Contains(baselineExpression, inputName) {
			t.Fatalf("the annotation's baseline expression %q reads %s — a baseline that depends on where the slowest day sits or how tall its bar is puts the clearance back at the mercy of the data, which is the defect REQ-242 fixed",
				baselineExpression, inputName)
		}
	}
}

// durationsStyleDeclaration reads one property out of one rule in the board's
// own stylesheet, so a test can hold a claim about how something is DRAWN
// against the sheet that draws it rather than against a copied value.
func durationsStyleDeclaration(t *testing.T, ruleSelector string, propertyName string) string {
	t.Helper()
	styleSheet, readError := embeddedWebAssets.ReadFile("web/board.css")
	if readError != nil {
		t.Fatalf("read web/board.css: %v", readError)
	}
	rulePattern := regexp.MustCompile(`(?s)` + regexp.QuoteMeta(ruleSelector) + `\s*\{(.*?)\}`)
	ruleMatch := rulePattern.FindSubmatch(styleSheet)
	if ruleMatch == nil {
		t.Fatalf("web/board.css declares no rule for %s", ruleSelector)
	}
	declarationPattern := regexp.MustCompile(`(?m)^\s*` + regexp.QuoteMeta(propertyName) + `\s*:\s*([^;]+);`)
	declarationMatch := declarationPattern.FindSubmatch(ruleMatch[1])
	if declarationMatch == nil {
		t.Fatalf("%s declares no %s", ruleSelector, propertyName)
	}
	return strings.TrimSpace(string(declarationMatch[1]))
}

// durationsSlowestDayAnnotationProbeCases hands the probe only the two fields it
// drives the renderer with; the case name stays on the Go side for the failure
// messages.
func durationsSlowestDayAnnotationProbeCases(annotationCases []durationsAnnotationCase) []map[string]float64 {
	probeCases := make([]map[string]float64, 0, len(annotationCases))
	for _, probeCase := range annotationCases {
		probeCases = append(probeCases, map[string]float64{
			"dayCentreX":    probeCase.dayCentreX,
			"medianMinutes": probeCase.medianMinutes,
		})
	}
	return probeCases
}

// ---- panel B/C day buckets and the shared axis domain ----------------------

// durationsRenderDomStubPreamble is the smallest DOM renderDurationsView
// touches, so the probe below can run the REAL renderer rather than sliced
// fragments of it. Every SVG node records its tag and attributes; nothing here
// re-implements layout.
const durationsRenderDomStubPreamble = `
function makeStubNode(nodeName) {
  return {
    stubName: nodeName,
    attributes: {},
    children: [],
    textContent: "",
    setAttribute: function (attributeName, attributeValue) { this.attributes[attributeName] = String(attributeValue); },
    appendChild: function (childNode) { this.children.push(childNode); return childNode; },
    removeChild: function (childNode) { this.children = this.children.filter(function (candidateNode) { return candidateNode !== childNode; }); },
    getComputedTextLength: function () { return 20; },
    addEventListener: function () {}
  };
}
var durationsStubHosts = {
  "durations-chart": makeStubNode("div"),
  "durations-summary": makeStubNode("p"),
  "durations-stat-median": makeStubNode("dd"),
  "durations-stat-p90": makeStubNode("dd"),
  "durations-stat-active-days": makeStubNode("dd"),
  "durations-stat-reqs-per-day": makeStubNode("dd"),
  "durations-readout": makeStubNode("p"),
  "durations-table-body": makeStubNode("tbody")
};
var document = {
  getElementById: function (nodeId) { return durationsStubHosts[nodeId] || null; },
  createElementNS: function (namespaceUri, nodeName) { return makeStubNode(nodeName); },
  createElement: function (nodeName) { return makeStubNode(nodeName); },
  createTextNode: function (nodeText) { return { textContent: nodeText }; }
};
`

func TestGenerateDurationsHeadlineStatsUseASemanticDefinitionList(t *testing.T) {
	indexHTML := generateLiveSite(t)

	requiredFragments := []string{
		`<dl class="durations-stats" id="durations-stats">`,
		`<dt>Median · all plotted spans</dt>`,
		`<dd id="durations-stat-median">`,
		`<dt>P90 · all plotted spans</dt>`,
		`<dd id="durations-stat-p90">`,
		`<dt>Active completion days</dt>`,
		`<dd id="durations-stat-active-days">`,
		`<dt>Projected REQs per active day</dt>`,
		`<dd id="durations-stat-reqs-per-day">`,
		`grid-template-columns: repeat(auto-fit, minmax(`,
	}
	lastOffset := -1
	for _, fragment := range requiredFragments {
		offset := strings.Index(indexHTML, fragment)
		if offset < 0 {
			t.Errorf("generated Durations view is missing %q", fragment)
			continue
		}
		if strings.HasPrefix(fragment, "<") && offset < lastOffset {
			t.Errorf("generated Durations stat fragment %q is out of semantic reading order", fragment)
		}
		if strings.HasPrefix(fragment, "<") {
			lastOffset = offset
		}
	}
	if strings.Count(indexHTML, `<dl class="durations-stats"`) != 1 ||
		strings.Count(indexHTML, `<dt>`) < 4 {
		t.Errorf("generated page does not carry one four-item Durations definition list")
	}
}

func durationHeadlineFixtureData(t *testing.T) generatedDurations {
	t.Helper()
	fixtureSpecs := []struct {
		requestID string
		completed time.Time
		minutes   float64
	}{
		{requestID: "REQ-801", completed: time.Date(2026, 4, 26, 10, 0, 0, 0, time.UTC), minutes: 60},
		{requestID: "REQ-802", completed: time.Date(2026, 5, 20, 10, 0, 0, 0, time.UTC), minutes: 70},
		{requestID: "REQ-803", completed: time.Date(2026, 6, 10, 10, 0, 0, 0, time.UTC), minutes: 40},
		{requestID: "REQ-804", completed: time.Date(2026, 7, 1, 10, 0, 0, 0, time.UTC), minutes: 50},
		{requestID: "REQ-805", completed: time.Date(2026, 8, 20, 9, 0, 0, 0, time.UTC), minutes: -5},
		{requestID: "REQ-806", completed: time.Date(2026, 8, 20, 10, 0, 0, 0, time.UTC), minutes: 10},
		{requestID: "REQ-807", completed: time.Date(2026, 8, 20, 11, 0, 0, 0, time.UTC), minutes: 300},
		{requestID: "REQ-808", completed: time.Date(2026, 8, 22, 10, 0, 0, 0, time.UTC), minutes: 20},
		{requestID: "REQ-809", completed: time.Date(2026, 8, 22, 11, 0, 0, 0, time.UTC), minutes: 400},
		{requestID: "REQ-810", completed: time.Date(2026, 8, 24, 10, 0, 0, 0, time.UTC), minutes: 30},
	}
	tickets := make([]*RequestTicket, 0, len(fixtureSpecs))
	for _, fixture := range fixtureSpecs {
		tickets = append(tickets, durationTicket(
			fixture.requestID,
			"B",
			fixture.completed.Add(-time.Duration(fixture.minutes*float64(time.Minute))).Format(time.RFC3339),
			fixture.completed.Format(time.RFC3339),
		))
	}
	generatedData, buildError := buildGeneratedBoardData(&Board{AllRequests: tickets})
	if buildError != nil {
		t.Fatalf("build headline fixture data: %v", buildError)
	}
	return generatedData.Durations
}

func durationRollingFixtureData(t *testing.T, eligibleDayCount int) generatedDurations {
	t.Helper()
	eligibleDays := []struct {
		completed time.Time
		minutes   time.Duration
	}{
		{completed: time.Date(2026, 7, 1, 10, 0, 0, 0, time.UTC), minutes: 10 * time.Minute},
		{completed: time.Date(2026, 7, 2, 10, 0, 0, 0, time.UTC), minutes: 70 * time.Minute},
		{completed: time.Date(2026, 7, 5, 10, 0, 0, 0, time.UTC), minutes: 20 * time.Minute},
		{completed: time.Date(2026, 7, 10, 10, 0, 0, 0, time.UTC), minutes: 60 * time.Minute},
		{completed: time.Date(2026, 7, 11, 10, 0, 0, 0, time.UTC), minutes: 30 * time.Minute},
		{completed: time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC), minutes: 50 * time.Minute},
		{completed: time.Date(2026, 8, 10, 10, 0, 0, 0, time.UTC), minutes: 40 * time.Minute},
		{completed: time.Date(2026, 9, 20, 10, 0, 0, 0, time.UTC), minutes: 80 * time.Minute},
	}
	if eligibleDayCount < 0 || eligibleDayCount > len(eligibleDays) {
		t.Fatalf("eligible day count %d is outside fixture", eligibleDayCount)
	}
	tickets := make([]*RequestTicket, 0, eligibleDayCount+5)
	for dayIndex, eligibleDay := range eligibleDays[:eligibleDayCount] {
		tickets = append(tickets, durationTicket(
			fmt.Sprintf("REQ-%03d", 820+dayIndex),
			"B",
			eligibleDay.completed.Add(-eligibleDay.minutes).Format(time.RFC3339),
			eligibleDay.completed.Format(time.RFC3339),
		))
	}
	// Five paused spans make 7 July an excluded-only active day and Panel C's
	// odd peak. It must not become a zero in the rolling median.
	for pausedIndex := 0; pausedIndex < 5; pausedIndex++ {
		completed := time.Date(2026, 7, 7, 10+pausedIndex, 0, 0, 0, time.UTC)
		tickets = append(tickets, durationTicket(
			fmt.Sprintf("REQ-%03d", 850+pausedIndex),
			"C",
			completed.Add(-8*time.Hour).Format(time.RFC3339),
			completed.Format(time.RFC3339),
		))
	}
	generatedData, buildError := buildGeneratedBoardData(&Board{AllRequests: tickets})
	if buildError != nil {
		t.Fatalf("build rolling fixture data: %v", buildError)
	}
	return generatedData.Durations
}

// durationsRenderProbeDriver renders the fixture board twice and reports every
// drawn bar's x-interval, the slowest-day annotation's anchor x (the only
// mid-anchored text at the annotation baseline), and every Panel A mark centre
// in draw order. Two renders make deterministic jitter a directly observed
// property instead of an inference from one set of coordinates.
const durationsRenderProbeDriver = `
function captureDurationsGeometry() {
  var drawnBars = [], drawnAnnotationXs = [], drawnMarkCxs = [];
  function walkDrawnNodes(parentNode) {
  (parentNode.children || []).forEach(function (childNode) {
    var attributes = childNode.attributes || {};
    if (childNode.stubName === "rect" && String(attributes["class"] || "").indexOf("durations-bar") !== -1) {
      drawnBars.push({ class: attributes["class"], x: Number(attributes.x), width: Number(attributes.width) });
    }
    if (childNode.stubName === "circle" && String(attributes["class"] || "").indexOf("durations-mark") !== -1) {
      drawnMarkCxs.push(Number(attributes.cx));
    }
    if (childNode.stubName === "text" && attributes["text-anchor"] === "middle" &&
        Number(attributes.y) === DURATIONS_MEDIAN_ANNOTATION_BASELINE_Y) {
      drawnAnnotationXs.push(Number(attributes.x));
    }
    walkDrawnNodes(childNode);
  });
  }
  walkDrawnNodes(durationsStubHosts["durations-chart"]);
  return { bars: drawnBars, annotationXs: drawnAnnotationXs, markCxs: drawnMarkCxs };
}
renderDurationsView();
var firstGeometry = captureDurationsGeometry();
durationsStubHosts["durations-chart"].children = [];
durationsStubHosts["durations-table-body"].children = [];
renderDurationsView();
var secondGeometry = captureDurationsGeometry();
process.stdout.write(JSON.stringify({
  bars: secondGeometry.bars,
  annotationXs: secondGeometry.annotationXs,
  markCxs: secondGeometry.markCxs,
  firstMarkCxs: firstGeometry.markCxs
}));
`

// durationsDayCountFixtureTickets builds dayCount active days of archived REQs.
// The first completion of each day lands 13:54 into it — the offset that put
// the real board's leftmost bar in the axis gutter (REQ-248).
func durationsDayCountFixtureTickets(dayCount int) []*RequestTicket {
	firstDay := time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC)
	tickets := make([]*RequestTicket, 0, dayCount*2)
	for dayIndex := 0; dayIndex < dayCount; dayIndex++ {
		dayStart := firstDay.AddDate(0, 0, dayIndex)
		for sampleIndex, completionOffset := range []time.Duration{
			13*time.Hour + 54*time.Minute,
			16*time.Hour + 54*time.Minute,
		} {
			completedAt := dayStart.Add(completionOffset)
			claimedAt := completedAt.Add(-time.Duration(24+10*sampleIndex) * time.Minute)
			tickets = append(tickets, durationTicket(
				fmt.Sprintf("REQ-%03d%d", dayIndex, sampleIndex),
				"B",
				claimedAt.Format(time.RFC3339),
				completedAt.Format(time.RFC3339),
			))
		}
	}
	return tickets
}

// REQ-248: Panel B placed its day bars at xOfEpoch(day midnight) while the axis
// domain began at the FIRST COMPLETION INSTANT, so the leftmost bar always sat
// left of the plot by however far into its day the first sample fell — and on a
// one- or two-day board the disagreement dominated the whole span and Panels B
// and C rendered off canvas entirely (live-DOM baselines: bar x=-5184.4 at one
// active day, x=-564.5 at two, x=-8.3 at fourteen against a left margin of 54;
// Chromium 141.0.7390.37 via Playwright 1.56.1).
//
// This drives the REAL renderDurationsView over a stub DOM at one, two,
// fourteen and four hundred active days and asserts, from the drawn attributes:
// every Panel B and C bar inside the plot area, exactly one slowest-day
// annotation anchored inside it, and every Panel A mark inside its own UTC-day
// slot. REQ-349 intentionally replaces exact completion-instant x with stable
// within-day jitter, so the old exact-x assertion would pin the defect. The
// bounds still hold the renderer's axis domain against the Go-side definition.
// REQ-349: Panel A needs enough vertical resolution for ordinary work, enough
// horizontal resolution for a busy day, and a quiet daily distribution behind
// the individual REQs. This drives the complete renderer so the assertions pin
// SVG output and draw order rather than helpers that could be left unwired.
// durationsUserRequestLaneFixtureTickets builds a board whose UR lane must
// overflow: userRequestCount URs whose brackets are far wider than the gap
// between their start instants, so they compete for the same rows, plus samples
// carrying no `user_request` at all. Both conditions are asserted in the test
// rather than assumed — a fixture that stopped overflowing would make the
// remainder assertion pass while proving nothing.
func durationsUserRequestLaneFixtureTickets(userRequestCount int, noUserRequestCount int) []*RequestTicket {
	weekStart := time.Date(2026, 8, 10, 9, 0, 0, 0, time.UTC)
	tickets := make([]*RequestTicket, 0, userRequestCount*2+noUserRequestCount)
	for userRequestIndex := 0; userRequestIndex < userRequestCount; userRequestIndex++ {
		// Each UR opens half an hour after the last and runs for twenty, so every
		// bracket overlaps almost every other one.
		firstCompletion := weekStart.Add(time.Duration(userRequestIndex) * 30 * time.Minute)
		for sampleIndex, completedAt := range []time.Time{firstCompletion, firstCompletion.Add(20 * time.Hour)} {
			ticket := durationTicket(
				fmt.Sprintf("REQ-%03d", userRequestIndex*2+sampleIndex+1),
				"B",
				completedAt.Add(-12*time.Minute).Format(time.RFC3339),
				completedAt.Format(time.RFC3339),
			)
			ticket.UserRequestId = fmt.Sprintf("UR-%03d", userRequestIndex+1)
			tickets = append(tickets, ticket)
		}
	}
	for noUserRequestIndex := 0; noUserRequestIndex < noUserRequestCount; noUserRequestIndex++ {
		// The pre-UR era: parseable stamps, no `user_request`. Placed well left of
		// the crowd so the bucket's bracket cannot be mistaken for a UR's.
		completedAt := weekStart.Add(time.Duration(-72-noUserRequestIndex*6) * time.Hour)
		tickets = append(tickets, durationTicket(
			fmt.Sprintf("REQ-9%02d", noUserRequestIndex+1),
			"A",
			completedAt.Add(-9*time.Minute).Format(time.RFC3339),
			completedAt.Format(time.RFC3339),
		))
	}
	return tickets
}

// REQ-346's two rules that were explicitly not the builder's call, pinned
// against the real renderer rather than against the packer it calls.
//
// RULE ONE: every UR that finds no lane row is COUNTED. The lane packs into a
// fixed number of rows, so on any board with more overlapping URs than rows some
// bracket cannot be drawn — and a reader who takes the drawn brackets for all of
// them reads the lane wrong. The invariant is a relationship, not a number:
// brackets drawn plus the count the lane states must equal the URs the samples
// carry, whatever the row count is later set to.
//
// RULE TWO: a sample whose REQ carries no `user_request` is NAMED, on every
// surface, and named apart from rule one's remainder. Twelve of this
// repository's own samples are in that state — REQ-001 through REQ-011 and
// REQ-060 pre-date the UR system — and buildDurationAggregate measures every one
// of them, so "blank" and "some default UR" are both wrong answers that a table
// full of real URs would hide.
//
// It drives renderDurationsView itself over the DOM stub, because the join, the
// packing, the remainder sentence and the table cell are four call sites of one
// rule and a probe calling the packer directly would hold none of them (REQ-305).
// REQ-347 changes Panel A's fill only. This drives the complete renderer across
// each selectable channel, because a helper-only test could keep the selector
// green while the circles still read sample.route. The fixture deliberately has
// more URs than the categorical palette, absent UR/domain values, and a reversed
// stamp: the overflow bucket and unknown state must be named while the critical
// stamp treatment remains stronger than any chosen colour channel.
func floatPointer(value float64) *float64 {
	return &value
}

func formatOptionalFloat(value *float64) string {
	if value == nil {
		return "null"
	}
	return fmt.Sprintf("%v", *value)
}

// timelineProbePreamble declares the renderer's own constants, read from the
// fragment rather than copied, so a probe cannot pass against numbers the
// shipped view does not use.
func timelineProbePreamble(t *testing.T, constantNames ...string) string {
	t.Helper()
	preamble := ""
	for _, constantName := range constantNames {
		preamble += fmt.Sprintf("var %s = %v;\n",
			constantName, rendererNumericConstant(t, "web/board-timeline.js", constantName))
	}
	return preamble
}

// Zoom is the one piece of this view that can be wrong in a way no screenshot
// shows: the instant under the pointer has to stay under the pointer, or every
// zoom drifts the chart sideways a little and the reader loses their place. It
// is a pure transform precisely so this can be asserted.
// Virtualization is what makes 560 rows cost the same as 40, and it is invisible
// in a screenshot: a wrong slice shows blank strips only while scrolling fast.
// Window-scoped rows: a REQ is listed when the bar it would draw overlaps the
// visible window, and is absent otherwise. Before REQ-319 every row was listed
// at every zoom level and the ones outside the window drew nothing — 305 labels
// above empty space on a 309-REQ board zoomed to one day.
//
// OVERLAP is the property under test, not containment. The four rows below are
// the four ways a row can relate to a window, and a containment rule would
// silently drop two of them: the REQ that started before the window and is still
// running is exactly the one a reader asking "what was happening this week"
// needs to see.
// Typed dates are the fourth way to move the window, and the one a reader can
// aim: the other three (pointer, keyboard, period chips) are all relative. What
// has to hold is that typing cannot reach a window the other three cannot —
// same floor, same ceiling, same edge clamp — because a control with its own
// clamp is how a view starts disagreeing with itself.
func TestTimelineRendererGroupsTheWindowedRowsByUserRequest(t *testing.T) {
	rendererBytes, readError := embeddedWebAssets.ReadFile("web/board-timeline.js")
	if readError != nil {
		t.Fatalf("read web/board-timeline.js: %v", readError)
	}
	if !strings.Contains(string(rendererBytes), "function timelineGroupWindowRows(") {
		t.Error("the Timeline has no window-scoped user-request grouping function")
	}
}

// The label's truncation arithmetic. The MEASUREMENT it depends on is a browser
// question and lives in the browser lane; this pins what the module does with
// the number once it has it.
//
// The id is the row's identity, so the rule under test is that it is never the
// thing that gets cut: a budget too small for id-plus-a-useful-title drops the
// title entirely rather than shipping a half-drawn id.
// timelineForecastDomStub is the smallest DOM renderTimelineForecast touches. It
// is a stub rather than a headless browser because what is being pinned is the
// SENTENCE — which figures reach the reader — not the layout.
const timelineForecastDomStub = `
function makeStubNode() {
  var node = {
    storedText: "",
    className: "",
    children: [],
    classList: { add: function () {}, remove: function () {} },
    setAttribute: function () {},
    appendChild: function (child) { this.children.push(child); return child; }
  };
  // Assigning textContent replaces ALL children, exactly as the DOM does. A
  // plain data property would leave the old children in place and make "was it
  // cleared?" unanswerable — which is how the first version of this probe
  // passed against code that cleared nothing.
  Object.defineProperty(node, "textContent", {
    get: function () { return this.storedText; },
    set: function (value) { this.storedText = value; this.children = []; }
  });
  return node;
}
var stubNodes = { "timeline-forecast": makeStubNode(), "timeline-excluded": makeStubNode() };
var document = {
  getElementById: function (id) { return stubNodes[id] || null; },
  createElement: function () { return makeStubNode(); },
  createTextNode: function (text) { return { textContent: text }; }
};
function collectText(node) {
  var text = node.textContent || "";
  (node.children || []).forEach(function (child) { text += " " + collectText(child); });
  return text;
}
`

// The REQ calls the honesty requirements load-bearing rather than decoration: a
// forecast that states a date without stating what it assumed is exactly the
// artifact people screenshot and quote. These pin that the sentence carries its
// own assumptions, and that thin history declines instead of guessing.
// At Fit all over three months a completed REQ was two adjacent 1.5px marks of
// different hues: a wait/work split drawn in pixels that cannot carry one. The
// collapse is the honest alternative — one marker for the row — and what has to
// hold is that it fires only where the split is genuinely unreadable, and never
// on the rows whose visible breakage is the reason they are drawn at all.
// The threshold decision, without a DOM. The browser probe next door proves the
// behaviour end to end; this pins the boundary itself, which is the part a
// pointer-event probe measures only at whatever offsets it happens to dispatch.
// The axis draws ticks and the plot draws gridlines at the same instants. Two
// loops doing the same arithmetic is one edit away from a gridline that means a
// slightly different time than the tick above it, so there is one source and
// this is what keeps it one.
// Rows advertise role="button" and take focus, but a <g> is not a native button:
// Enter and Space never synthesize the click the drawer listens for. The role is
// a promise, and this is the code that keeps it.
// Zoom and pan now have two drivers — a pointer and a keyboard — and two drivers
// that each compute a window are two definitions of where the window goes. The
// keyboard path is written as pure transforms over the SAME timelineZoomedWindow
// the wheel and the zoom buttons call; this drives both to their edges and
// requires them to arrive at the same ones.
// The interaction has to be discoverable without seeing the hint line beside the
// chart, and focusable on the element that actually takes the keys. Both are one
// attribute each, and both are silently droppable in a template edit.
func TestTimelinePanelStatesItsKeyboardInteraction(t *testing.T) {
	indexHtml := generateLiveSite(t)

	panelStart := strings.Index(indexHtml, `<section id="view-timeline"`)
	if panelStart == -1 {
		t.Fatal("the generated page carries no timeline panel")
	}
	panelOpeningTag := indexHtml[panelStart : panelStart+strings.Index(indexHtml[panelStart:], ">")+1]
	for _, wantPhrase := range []struct {
		requirement string
		phrase      string
	}{
		{"the panel is still named", "Timeline"},
		{"which keys pan", "arrow keys"},
		{"which keys zoom", "plus and minus"},
	} {
		if !strings.Contains(panelOpeningTag, wantPhrase.phrase) {
			t.Fatalf("the timeline panel's accessible name does not state %s (wanted %q in %q)",
				wantPhrase.requirement, wantPhrase.phrase, panelOpeningTag)
		}
	}

	scrollStart := strings.Index(indexHtml, `id="timeline-scroll"`)
	if scrollStart == -1 {
		t.Fatal("the generated page carries no timeline scroll container")
	}
	scrollOpeningTag := indexHtml[strings.LastIndex(indexHtml[:scrollStart], "<") : scrollStart+strings.Index(indexHtml[scrollStart:], ">")+1]
	if !strings.Contains(scrollOpeningTag, `tabindex="0"`) {
		t.Fatalf("the chart cannot be focused, so its keyboard path is unreachable: %q", scrollOpeningTag)
	}
}

// The Lens group must offer all three readings. URs only selects the by-UR lens
// and adds the fold, so it carries the same data-lens-target as By UR plus the
// fold attribute — that is what keeps every other "is the UR lens on screen"
// test (filters, testing transitions, the recently-done chip) correct for it.
func TestGenerateOffersThreeLensButtons(t *testing.T) {
	indexHtml := generateLiveSite(t)
	for _, requiredToken := range []string{
		`data-lens-target="flat"`,
		`data-lens-target="user-request"`,
		`data-ur-cards="folded"`,
		"URs&nbsp;only",
		`applyLensSelection(button.getAttribute("data-lens-target"), button.getAttribute("data-ur-cards"))`,
	} {
		if !strings.Contains(indexHtml, requiredToken) {
			t.Fatalf("lens group is missing %q in the generated page", requiredToken)
		}
	}
}

// One UR group as the probe reports it back: which UR it names, whether its
// header announces an expanded fold, which REQ cards are actually in the DOM,
// and which nodes inside it still open the UR detail drawer.
type renderedUserRequestRow struct {
	UserRequestId  string   `json:"userRequestId"`
	Expanded       string   `json:"expanded"`
	CardIds        []string `json:"cardIds"`
	DrawerTriggers []string `json:"drawerTriggers"`
}

// URs only is the by-UR lens with its REQ cards folded away, rendered by the
// same function, so this drives the production renderer under Node: folded it
// must emit one row per UR and no cards until a row is opened, opening a row
// must reveal exactly that UR's filtered cards, and the Active scope plus a
// status filter must hide exactly the URs they hide in the by-UR lens.
func userRequestIdsOf(rows []renderedUserRequestRow) []string {
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.UserRequestId)
	}
	return ids
}

// rendererDeclarationLine returns one `var NAME = ...;` line verbatim from a
// shipped renderer fragment. It is rendererNumericConstant's sibling for the
// declarations that are not a single number — a probe that needs the view's list
// of period levels must drive THAT list, not a copy beside it that goes stale.
func rendererDeclarationLine(t *testing.T, assetPath string, constantName string) string {
	t.Helper()
	rendererText, readError := embeddedWebAssets.ReadFile(assetPath)
	if readError != nil {
		t.Fatalf("read %s: %v", assetPath, readError)
	}
	pattern := regexp.MustCompile(`(?m)^\s*(var ` + regexp.QuoteMeta(constantName) + ` = .+;)$`)
	match := pattern.FindSubmatch(rendererText)
	if match == nil {
		t.Fatalf("%s declares no single-line constant %s", assetPath, constantName)
	}
	return string(match[1])
}

// rendererBracketDeclaration slices a multi-line `var NAME = [ ... ];` out of a
// shipped asset, balancing SQUARE brackets.
//
// sliceBalancedBlockAfter balances braces, so pointing it at an array of object
// literals returns the first element and nothing else — the probe then fails to
// parse, which is how this helper came to exist. rendererDeclarationLine is the
// single-line case; this is the same contract for a block.
func rendererBracketDeclaration(t *testing.T, assetPath string, constantName string) string {
	t.Helper()
	rendererBytes, readError := embeddedWebAssets.ReadFile(assetPath)
	if readError != nil {
		t.Fatalf("read %s: %v", assetPath, readError)
	}
	rendererText := string(rendererBytes)
	anchor := "var " + constantName + " = ["
	anchorIndex := strings.Index(rendererText, anchor)
	if anchorIndex == -1 {
		t.Fatalf("%s declares no bracketed constant %s", assetPath, constantName)
	}
	bracketDepth := 0
	for scanOffset := anchorIndex; scanOffset < len(rendererText); scanOffset++ {
		switch rendererText[scanOffset] {
		case '[':
			bracketDepth++
		case ']':
			bracketDepth--
			if bracketDepth == 0 {
				return rendererText[anchorIndex:scanOffset+1] + ";"
			}
		}
	}
	t.Fatalf("no bracket-balanced declaration found for %s in %s", constantName, assetPath)
	return ""
}

// The period toolbar's windows TRAIL NOW instead of naming a calendar period, and
// what can be wrong about a trailing window is invisible in a screenshot: the bars
// still look like bars while the window quietly no longer ends where it claims to.
//
// Four properties, each a defect this arithmetic can have:
//
//	a chip's window ends at NOW and spans exactly what it asked for;
//	a chip wider than the archive is CUT SHORT at the range start and STILL ends
//	  at now — the settle preserves a WIDTH and slides, so an unclamped candidate
//	  gets pinned to the range start and dragged forward off now;
//	All days is the recorded range exactly, so nothing on the chart is unreachable
//	  from the button that claims to show all of it;
//	timelinePannedWindow moves one screenful and is its own inverse. This is the
//	  CONTINUOUS pan the keyboard and the drag rest on, mid-range where its clamp
//	  cannot fire. What the toolbar's arrows do with it near a bound — where the
//	  clamp does fire — is TestJavaScriptBehaviorTimelineWindowStepArrowsAreInversesAtTheBounds.
//
// It drives the shipped functions rather than reimplementing them (REQ-305), and
// reads the renderer's own constants rather than restating them (REQ-322).
// The five trailing-window chips have to keep working after the queue DRAINS,
// which is the state a queue-draining tool spends its life heading toward.
//
// timelineRange in timeline.go stretches the payload's range end to now only
// while some row still has WaitOpen or WorkOpen. Finish everything and the bounds
// end at the last completion instant, now falls OUTSIDE them, and a window
// computed as [now - N days, now] lands entirely past the bound end, clamps to
// zero width and settles onto the one-hour zoom floor. Measured on the merged
// tree at 59105df: three days idle collapsed "Last day"; ten days idle put "Last
// day" and "Last 7 days" on the SAME one-hour window, so pressing the second lit
// the first; a hundred days idle collapsed four of the five.
//
// Four properties, and the sweep over board ages is what makes them properties
// rather than four fixed cases — a board is idle for however long it is idle, and
// a fix keyed on the ages that happened to get measured is not a fix:
//
//	no chip settles onto the zoom floor while the board has a range to show;
//	the chips stay pairwise DISTINCT, which is what makes the chip the reader
//	  pressed the one that lights — renderTrailingWindowControls lights the first
//	  chip whose candidate window matches the window on screen;
//	a chip's window is the last N days OF THE RECORDED RANGE: it ends at now while
//	  now is inside the bounds and at the range end once now has left them;
//	the window carries whether it got the whole span the chip asked for, so the
//	  "part of" readout cannot become a second opinion about the clamp.
//
// The control set is read out of the SHIPPED PAGE rather than restated here, so a
// chip added to template.html is swept with the rest (REQ-322: read the value the
// decision turns on, never restate it beside the test).
// A pure timelineTrailingWindow case supplied with unpadded bounds cannot catch
// the production ordering bug: drawTimeline pads boundEndMs before its nested
// chip caller invokes the helper. Drive those shipped production statements so
// the test sees both the semantic end and the cosmetic display range.
// The toolbar's ‹ and › are the reader's undo for each other, and near the right
// bound they stopped being it: timelinePannedWindow CLAMPS, so a forward step
// with less than a screenful of room ahead moved a partial screenful while the
// back step that followed moved a full one. Measured on the merged tree at
// 59105df with a 7-day window on a 90-day range: -120.00h of drift a partial
// screenful from the right bound, -168.00h flush against it, 0.00h mid-range.
//
// This is a REGRESSION with a history: the calendar-period test REQ-390 deleted
// pinned the inverse property, and its replacement checks only the mid-range case
// — the one case where the clamp cannot fire. So the sweep here walks the window
// across the whole range, and the property is stated as the pair of rules that
// makes a discrete step coherent rather than as one arithmetic identity:
//
//	a step MOVES A WHOLE SCREENFUL or it does not move at all, never a partial one;
//	wherever it moves, the opposite step returns to exactly where it started;
//	and it does move wherever there is a screenful of room, which is the half that
//	  keeps "refuse the partial step" from being satisfied by refusing everything.
//
// It follows the ARROWS' OWN CALL SITE rather than naming a window function: the
// probe reads which function steppedWindowFor hands the press to and drives that
// one, so a step function swapped in without this property comes back red instead
// of leaving the test measuring something the arrows no longer call (REQ-305).
// The Now button is TWO movements, and the second used to be missing: recentring
// the time window never moved the ROW list, so "jump to the remaining work" landed
// the reader on whichever archived rows happened to be scrolled into view.
//
// It also pins the window half's DEGENERATE cases, which are the ordinary ones on
// a queue that is nearly drained: the forecast's queue-empty instant minutes from
// the now-line, or no forecast at all, so the span the margin is a fraction OF is
// effectively zero and only the floor decides the window. Flooring on half the
// zoom floor put Now on a one-hour window — the floor itself — and the obvious
// next move was dead: the + button, ctrl+wheel and the + key were all silent
// no-ops with no disabled state and no message.
//
// Split out of the calendar-period navigation test when REQ-390 replaced the
// Day/Week/Month chips with trailing windows. This half never touched the calendar
// arithmetic that went with them, and deleting the test wholesale would have taken
// the only Node-lane coverage of timelineNowJump and timelineFirstOpenRowIndex —
// including the NaN forecast and the decide-the-scroll-after-the-refresh ordering
// — with it.
// The axis is the one part of this view whose defect is pure text. The ticks are
// drawn where the code says, so nothing about the drawing looks wrong while the
// labels describe instants the ticks are not at. REQ-227 wrote the minute as the
// literal ":00", and REQ-235's Now button made a sub-hour window the landing
// state of the view's most-used control: seven ticks, two distinct labels.
//
// REQ-327 ADDED THE ASSERTION THIS TEST WAS MISSING, and it is worth saying why
// the old one passed over the defect. It required that every number in a label be
// one the tick's own instant carries — and a bare "9 Jul" on a tick at
// 9 Jul 12:00 satisfies that: the day IS 9 July. What it never required was that
// a label OMITTING the time be at midnight. So the week of 6 July rendered
// "6, 7, 8, 9, 10, 11, 13 Jul", every intermediate tick four hours off the date it
// named, 12 July absent, and the whole suite green.
//
// It also DRIVES THE REAL TICK SOURCE now. It used to reimplement the equal-parts
// spacing inline, so it could not have noticed the spacing changing under it
// (REQ-305: a probe that reimplements the function under test cannot hold its call
// site).
//
// Four properties, per window: no two ticks share a label; every number in a label
// belongs to its tick; a label with no time is at UTC midnight; and the ticks
// ascend and stay inside the window.
// The Timeline's rows are keyboard-focusable SVG <g> nodes whose outline board.css
// used to switch off outright, leaving a focused row with only a one-step background
// tint while every other focusable thing on the board drew a 2px ring. This pins the
// row's ring to the SAME width and token as the board's reference ring, so a later
// change to one cannot quietly weaken the other, and pins the offset as inset: at a
// positive offset the ring is clipped on three sides — the rows SVG's own viewport
// takes the left and right edges and the scroll container takes the top — which
// paints a divider under the next row instead of a ring around this one.
func TestGenerateGivesTimelineRowsTheBoardsFocusRing(t *testing.T) {
	indexHtml := generateLiveSite(t)

	const rowFocusSelector = ".timeline-row:focus-visible {"
	if !strings.Contains(indexHtml, rowFocusSelector) {
		t.Fatalf("Timeline rows carry no %q rule: a keyboard-focused row has no ring, only the tint", rowFocusSelector)
	}
	rowFocusRule := sliceBalancedBlockAfter(t, indexHtml, rowFocusSelector)
	referenceRingRule := sliceBalancedBlockAfter(t, indexHtml, ".control-button:focus-visible {")

	ringPattern := regexp.MustCompile(`outline:\s*(\d+)px\s+solid\s+var\((--[a-z0-9-]+)\)`)
	referenceRing := ringPattern.FindStringSubmatch(referenceRingRule)
	if referenceRing == nil {
		t.Fatalf("the board's reference ring no longer declares a token-coloured outline: %q", referenceRingRule)
	}
	rowRing := ringPattern.FindStringSubmatch(rowFocusRule)
	if rowRing == nil {
		t.Fatalf("a keyboard-focused Timeline row draws no token-coloured outline: %q", rowFocusRule)
	}
	if rowRing[1] != referenceRing[1] || rowRing[2] != referenceRing[2] {
		t.Fatalf("the Timeline row ring is %spx var(%s), but the board's rings are %spx var(%s)",
			rowRing[1], rowRing[2], referenceRing[1], referenceRing[2])
	}

	offsetPattern := regexp.MustCompile(`outline-offset:\s*(-?\d+)px`)
	rowOffset := offsetPattern.FindStringSubmatch(rowFocusRule)
	if rowOffset == nil || !strings.HasPrefix(rowOffset[1], "-") {
		t.Fatalf("the Timeline row ring must be drawn inward — an outward ring is clipped by the rows SVG and the scroll container; rule is %q", rowFocusRule)
	}
}

// spanRoundingCase is one duration and the exact text each renderer formatter
// must draw for it. Both formatters are pinned in one probe because they had the
// same defect for the same reason, and a fix applied to one is a regression in
// the other's agreement with it.
type spanRoundingCase struct {
	Minutes           float64 `json:"minutes"`
	WantDurationsText string  `json:"wantDurationsText"`
	WantTimelineText  string  `json:"wantTimelineText"`
	Requirement       string  `json:"-"`
}

// The renderers split a magnitude into units and round each unit on its own, so
// a remainder that rounded up printed a value its own field cannot hold: "1h 60m"
// for 119.5 minutes, "1d 24h" for 2879, "60 min" for 59.96. These are ordinary
// REQ durations (1h59m30s is not a corner case), and the labels are what the
// charts, hover text, tables and forecast sentence all read from.
// TestJavaScriptBehaviorCalendarDayBreakdownGroupsStatuses executes the real
// calendarDayBreakdown from the assembled client. The count line it feeds ("12
// done · 2 cancelled") is the calendar's main at-a-glance signal, so a status
// landing in the wrong group misreports the day — counting abandoned or
// still-open work as "done" is the failure that matters. It also pins that only
// non-zero groups render, that the three blocked variants collapse into one
// group while an unrecognized status does NOT join them, and the fixed group
// order the colours in board.css are written against.
// REQ-284's second captured RED: the emitted payload must carry what verify sees,
// minus the three categories the board already renders by other means. Forwarding
// those would print the same prose a second or third time — anomalies reach the page
// through their own column, a per-card badge, AND board.Warnings.
func TestGeneratedBoardDataCarriesVerifyFindingsWithoutTheOnesTheBoardAlreadyShows(t *testing.T) {
	claimedAt := time.Date(2026, 8, 19, 9, 0, 0, 0, time.UTC)
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		// A stale claim: must appear in verifyFindings.
		{"do-work/working/REQ-820-stale-claim.md",
			"---\nid: REQ-820\ntitle: fixture\nstatus: claimed\n" +
				"claimed_at: 2026-08-19T09:00:00Z\n---\n"},
		// A completion anomaly: must NOT appear — the board renders it three ways.
		{"do-work/archive/REQ-821-anomaly.md",
			"---\nid: REQ-821\ntitle: fixture\nstatus: completed\n" +
				"claimed_at: 2026-08-19T10:00:00Z\ncompleted_at: 2026-08-19T09:00:00Z\n---\n"},
		// An unrecognized status reaches the verifier but already appears in the
		// board's Needs input / Blocked strip, so its prose is suppressed.
		{"do-work/queue/REQ-822-unrecognized-status.md",
			"---\nid: REQ-822\ntitle: fixture\nstatus: pnding\nuser_request: UR-900\n---\n"},
		// Structural damage has no separate board strip and must keep reaching the
		// payload even as unrecognized status is suppressed.
		{"do-work/queue/REQ-823-missing-user-request.md",
			"---\nid: REQ-823\ntitle: fixture\nstatus: pending\n---\n"},
	})
	moment := claimedAt.Add(4 * time.Hour)

	board, buildError := buildBoard(repoRoot, moment, defaultRecentWindow, lookupGitCommitDate)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}
	boardData, projectError := buildGeneratedBoardData(board)
	if projectError != nil {
		t.Fatalf("buildGeneratedBoardData: %v", projectError)
	}
	attachVerifyFindings(&boardData, board, moment)

	sawStaleClaim := false
	sawStructuralDamage := false
	for _, finding := range boardData.VerifyFindings {
		if boardRenderedVerifyCategories[finding.Category] {
			t.Errorf("suppressed category %q reached verifyFindings: %s", finding.Category, finding.Detail)
		}
		if finding.Category == verifyCategoryClaimNeedsAttention {
			sawStaleClaim = true
			if finding.Remedy == "" {
				t.Error("stale-claim finding reached the page without its remedy")
			}
		}
		if finding.Category == verifyCategoryStructurallyDamagedRequest {
			sawStructuralDamage = true
		}
	}
	if !sawStaleClaim {
		t.Errorf("verifyFindings carries no stale-claim finding: %+v", boardData.VerifyFindings)
	}
	if !sawStructuralDamage {
		t.Errorf("verifyFindings lost structural damage while suppressing board-rendered findings: %+v", boardData.VerifyFindings)
	}
	// The anomaly must still exist as a finding — it is suppressed from the page,
	// not from verify. Otherwise this test would pass on a probe that stopped working.
	report := collectVerifyFindings(repoRoot, board, moment)
	if anomalies := findingsMentioning(report, verifyCategoryCompletionAnomaly); len(anomalies) == 0 {
		t.Error("the fixture produced no completion anomaly, so the suppression assertion proves nothing")
	}
	if statuses := findingsMentioning(report, verifyCategoryUnrecognizedRequestStatus); len(statuses) == 0 {
		t.Error("the fixture produced no unrecognized status, so its suppression assertion proves nothing")
	}
}

// The no-absolute-paths capture decision, kept honest by assertion. The worktree
// probe's detail is absolute at its source (`git worktree list --porcelain`), so
// this fails the moment someone forwards a finding unreduced. A shared snapshot
// must not describe the filesystem of the machine that produced it.
func TestGeneratedVerifyPayloadCarriesNoAbsolutePaths(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/working/REQ-830-stale-claim.md",
			"---\nid: REQ-830\ntitle: fixture\nstatus: claimed\n" +
				"claimed_at: 2026-08-19T09:00:00Z\n---\n"},
	})
	moment := time.Date(2026, 8, 19, 13, 0, 0, 0, time.UTC)

	board, buildError := buildBoard(repoRoot, moment, defaultRecentWindow, lookupGitCommitDate)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}
	boardData, projectError := buildGeneratedBoardData(board)
	if projectError != nil {
		t.Fatalf("buildGeneratedBoardData: %v", projectError)
	}
	// Seed a finding whose detail is absolute at the source, the way the worktree
	// probe's is, so the reduction is under test rather than merely unexercised.
	boardData.VerifyFindings = nil
	syntheticReport := VerifyReport{Findings: []VerifyFinding{{
		Category: verifyCategoryUnmergedWorktreeLeftover,
		Detail:   "worktree " + filepath.Join(repoRoot, "..", "repo-worktrees", "worktree-agent-REQ-999") + " is unmerged",
		Remedy:   "inspect " + filepath.Join(repoRoot, "do-work", "queue") + " first",
	}}}
	for _, finding := range syntheticReport.Findings {
		boardData.VerifyFindings = append(boardData.VerifyFindings, generatedVerifyFinding{
			Category: finding.Category,
			Detail:   reduceAbsolutePaths(finding.Detail, board.RepoRoot),
			Remedy:   reduceAbsolutePaths(finding.Remedy, board.RepoRoot),
		})
	}
	attachVerifyFindings(&boardData, board, moment)

	encoded, encodeError := encodeBoardDataForJsAssignment(boardData)
	if encodeError != nil {
		t.Fatalf("encodeBoardDataForJsAssignment: %v", encodeError)
	}
	if strings.Contains(encoded, repoRoot) {
		t.Errorf("emitted payload contains the repo root %q", repoRoot)
	}
	for _, finding := range boardData.VerifyFindings {
		for _, text := range []string{finding.Detail, finding.Remedy} {
			if remainingAbsolutePath.MatchString(text) {
				t.Errorf("verify payload still carries an absolute path: %q", text)
			}
		}
	}
}

// The reduction must strip absolute paths WITHOUT touching relative ones. RE2 has
// no lookbehind, so the boundary in remainingAbsolutePath is captured and restored;
// the first version asserted nothing and matched the `/` inside an already-relative
// path, turning "do-work/calibration-log.tsv" into "do-work<path…>" — mangling
// precisely the paths the repo-root reduction had just produced.
func TestReduceAbsolutePathsLeavesRelativePathsIntact(t *testing.T) {
	repoRoot := filepath.Join("/tmp", "fixture-repo")
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			"repo-root path becomes relative and stays whole",
			"calibration-log probe: " + filepath.Join(repoRoot, "do-work", "calibration-log.tsv") + " is absent",
			"calibration-log probe: do-work/calibration-log.tsv is absent",
		},
		{
			"an already-relative path is untouched",
			"open do-work/queue/REQ-042-thing.md first",
			"open do-work/queue/REQ-042-thing.md first",
		},
		{
			"a path outside the repo is replaced wholesale",
			"worktree /elsewhere/repo-worktrees/worktree-agent-REQ-999 is unmerged",
			"worktree <path outside this repository> is unmerged",
		},
		{
			"a Windows drive path with backslashes is replaced",
			`worktree C:\Users\alice\proj\worktree-agent-REQ-999 is unmerged`,
			"worktree <path outside this repository> is unmerged",
		},
		{
			// Git on Windows emits this form too. A drive pattern that accepts only
			// the backslash form let it through whole, and `C:` is not a boundary
			// character so the POSIX branch could not catch the `/` after it either.
			"a Windows drive path with forward slashes is replaced",
			"worktree C:/Users/alice/proj/worktree-agent-REQ-999 is unmerged",
			"worktree <path outside this repository> is unmerged",
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := reduceAbsolutePaths(testCase.input, repoRoot); got != testCase.want {
				t.Errorf("reduceAbsolutePaths(%q)\n got %q\nwant %q", testCase.input, got, testCase.want)
			}
		})
	}
}

// Every measured-face constant names the build its number came from.
//
// A face is per-browser: the same .durations-mark-label ascent measured 10.1853
// on one Chromium and 10.4278 on another, and REQ-241 and REQ-242 collided
// because a number measured on one build was read as fact on another. The rule
// that came out of that is enforced, not remembered.
//
// WHY THIS LIVES HERE NOW. The rule used to be enforced by
// durations_test.go's TestDurationsMeasuredConstantsNameTheirChromiumBuild.
// REQ-292 deleted that test along with the width model, on the reasoning that no
// measured constant survived to name a build — which was true of durations.go and
// board-durations.js, the files REQ-292 was clearing, and NOT true of this file,
// where four such constants live and are read. REQ-277's sweep found the rule's
// stated enforcer pointing at a deleted test, which is the exact defect class
// that REQ asks about: a comment claiming an arrangement that is no longer real.
//
// The vacuity guard is the same one the deleted test carried and for the same
// reason: a scan that finds zero constants must fail rather than pass, or this
// becomes a green check over an empty set the day someone renames the prefix.
func TestMeasuredFaceConstantsNameTheirBuild(t *testing.T) {
	measuredConstantDeclaration := regexp.MustCompile(`(?m)^const (durationsMeasured[A-Za-z]+)\s*=`)
	buildAnchor := regexp.MustCompile(`(?i)chromium|playwright|firefox|webkit|safari`)

	sourceText := readPackageSourceForTest(t, "generate_test.go")
	declarations := measuredConstantDeclaration.FindAllStringSubmatchIndex(sourceText, -1)
	if len(declarations) == 0 {
		t.Fatal("found no durationsMeasured* constants to check — this scan must never pass on an empty set")
	}

	sourceLines := strings.Split(sourceText, "\n")
	for _, declaration := range declarations {
		constantName := sourceText[declaration[2]:declaration[3]]
		declarationLine := strings.Count(sourceText[:declaration[0]], "\n")
		// The doc comment is the contiguous run of // lines immediately above.
		commentLines := []string{}
		for lineIndex := declarationLine - 1; lineIndex >= 0; lineIndex-- {
			if !strings.HasPrefix(strings.TrimSpace(sourceLines[lineIndex]), "//") {
				break
			}
			commentLines = append(commentLines, sourceLines[lineIndex])
		}
		if !buildAnchor.MatchString(strings.Join(commentLines, "\n")) {
			t.Errorf("%s has no browser build named in its doc comment — a measured face is per-browser, "+
				"so an unnamed build makes the number read as timeless fact (REQ-241/REQ-242 collided on exactly that)",
				constantName)
		}
	}
}

// readPackageSourceForTest reads one of this package's own source files off disk.
// The measured constants live in a _test.go file, which go:embed cannot reach, so
// this is a plain read rather than an embedded asset.
func readPackageSourceForTest(t *testing.T, fileName string) string {
	t.Helper()
	sourceBytes, readError := os.ReadFile(fileName)
	if readError != nil {
		t.Fatalf("read %s: %v", fileName, readError)
	}
	return string(sourceBytes)
}

// TestGenerateInlinesImpactAndEffortChipRenderPath guards the frontend half of
// the impact/effort split (REQ-289), which had NO test coverage of any kind until
// REQ-293 — neither `badge-impact` nor `badge-effort-estimate` appeared in any
// test file. The Go side pins the vocabulary and the parser; nothing proved the
// chips still get rendered, so a refactor that dropped either renderer from
// web/board-cards.js would ship a silent regression. The chips only appear when
// the live tree happens to carry a non-default value, so a queue-dependent
// assertion would be no assertion at all.
//
// These are code tokens from the assembled client and board.css, so the check
// holds regardless of what the queue currently contains — the same shape
// TestGenerateInlinesWriteSetOverlapBadgeRenderPath established for the overlap
// badge. REQ-289's own Discovered Task framed this as needing a JavaScript
// behavior probe; the cheaper precedent was in this file already.
func TestGenerateInlinesImpactAndEffortChipRenderPath(t *testing.T) {
	indexHtml := generateLiveSite(t)

	for _, renderToken := range []string{
		// The makeBadge() calls that emit the two chips. The quoted forms occur
		// only in board-cards.js — the bare class names would also match the CSS.
		`"badge-impact"`,
		`"badge-effort-estimate"`,
		// The payload field the impact chip reads, together with the default it
		// falls back to. That default is the property REQ-290's
		// --skip-impact-negligible depends on: absence must read as
		// impact-user-visible, never as the user's stop signal.
		`request.impact || "impact-user-visible"`,
		// Without the stylesheet rules the chips render unstyled and invisible.
		".badge-impact",
		".badge-effort-estimate",
	} {
		if !strings.Contains(indexHtml, renderToken) {
			t.Fatalf("impact/effort chip render path is missing from the generated page: %q not found in the assembled client/board.css", renderToken)
		}
	}
}

// ---- timeline: a reversed span is drawn as a break, never as a bar ----------

// timelineRenderDomStubPreamble is the smallest DOM renderTimelineView touches,
// so the probe below runs the REAL renderer rather than a sliced fragment of it.
// Every SVG node records its tag and attributes; nothing here re-implements
// layout. Same shape as durationsRenderDomStubPreamble, with the extra surface
// the timeline uses: a scroll host with measurable geometry (row virtualization
// asks for it), classList on the forecast nodes, and the two globals the module
// reads from its siblings in the assembled client.
const timelineRenderDomStubPreamble = `
function makeStubNode(nodeName) {
  return {
    stubName: nodeName,
    attributes: {},
    children: [],
    textContent: "",
    clientWidth: 900,
    clientHeight: 400,
    scrollTop: 0,
    style: {},
    classList: { add: function () {}, remove: function () {}, toggle: function () {} },
    setAttribute: function (attributeName, attributeValue) { this.attributes[attributeName] = String(attributeValue); },
    appendChild: function (childNode) { this.children.push(childNode); return childNode; },
    addEventListener: function () {},
    removeEventListener: function () {},
    getBoundingClientRect: function () { return { width: 900, height: 400, left: 0, top: 0 }; }
  };
}
var timelineStubHosts = {};
[
  "timeline-summary", "timeline-axis", "timeline-scroll", "timeline-readout",
  "timeline-table-body", "timeline-forecast", "timeline-excluded", "timeline-period-state"
].forEach(function (hostId) { timelineStubHosts[hostId] = makeStubNode("div"); });
var document = {
  getElementById: function (nodeId) { return timelineStubHosts[nodeId] || null; },
  createElementNS: function (namespaceUri, nodeName) { return makeStubNode(nodeName); },
  createElement: function (nodeName) { return makeStubNode(nodeName); },
  createTextNode: function (nodeText) { return { textContent: nodeText }; },
  querySelectorAll: function () { return []; },
  querySelector: function () { return null; }
};
var window = { addEventListener: function () {}, removeEventListener: function () {} };
var requestsById = {};
var generatedAtMs = Date.parse("2026-08-18T12:00:00Z");
// Overwrite timelineStubVisibleIds to simulate an active filter; null means no
// filter, which is what the real requestMatchesFilters reports with none set.
var timelineStubVisibleIds = null;
function requestMatchesFilters(requestId) {
  return timelineStubVisibleIds === null || timelineStubVisibleIds.indexOf(requestId) !== -1;
}
function setActiveButton() {}
`

// timelineRenderProbeDriver reports every drawn rect's class and width, in draw
// order, grouped by the row group that carries them — the row id is the only
// thing that ties a segment back to the fixture that produced it.
const timelineRenderProbeDriver = `
renderTimelineView();
var drawnRows = [];
function collectRowRects(node, sink) {
  (node.children || []).forEach(function (childNode) {
    var attributes = childNode.attributes || {};
    if (childNode.stubName === "rect" && attributes["class"]) {
      sink.push({ class: attributes["class"], width: Number(attributes.width) });
    }
    collectRowRects(childNode, sink);
  });
}
function walkRowGroups(node) {
  (node.children || []).forEach(function (childNode) {
    var attributes = childNode.attributes || {};
    if (childNode.stubName === "g" && attributes["data-detail-id"]) {
      var rowRects = [];
      collectRowRects(childNode, rowRects);
      drawnRows.push({ id: attributes["data-detail-id"], rects: rowRects });
      return;
    }
    walkRowGroups(childNode);
  });
}
walkRowGroups(timelineStubHosts["timeline-scroll"]);
process.stdout.write(JSON.stringify({ rows: drawnRows }));
`

// TestJavaScriptBehaviorTimelineRefusesToRenderAgainstAnUnmeasurableHost pins the
// other half of the plot-width fix, the half no click can reach.
//
// `addTimelineListener(window, "resize", renderAll)` fires while #view-timeline is
// `hidden`, where the scroll host's clientWidth is 0. plotWidth's floor turned
// that into Math.max(120, 0 - 184 - 12) = 120 and MEMOISED it, and
// timelineVisibleRowRange turned clientHeight 0 into eight rows — so a browser
// resize taken on another view left the Timeline showing three months of archive
// crushed into a 120-pixel strip with eight rows in it, and only a window move
// repaired it.
//
// The property is that an unmeasurable host is "not yet", never 120: the render is
// skipped rather than performed against numbers that describe nothing. The
// recovery half is the ResizeObserver, which a DOM stub has none of — that is what
// TestBrowserBehaviorTimelineBarsSurviveTheDetailDrawerOpening covers in a real
// engine. Here the question is only whether the wrong numbers get written.
// TestJavaScriptBehaviorReversedWaitDrawsAsABreak pins the wait segment to the
// rule the work segment already followed: a span whose end precedes its start
// has no width to draw honestly, so it is a break marker rather than a bar.
//
// drawSegment sorts its endpoints with Math.min/Math.max — correctly, for every
// caller that should reach it — so a reversed wait handed to it painted as an
// ordinary positive-width waiting bar while the table beside it printed the
// signed value. A row whose numbers say "−60 min" and whose bar says "60 min of
// waiting" is broken bookkeeping rendering as healthy.
//
// All three cases render in ONE pass over one payload, deliberately: the bug was
// a missing branch, and a fix that turned every wait into a break would satisfy
// a reversed-only test.
// TestJavaScriptBehaviorTimelineFallbackBoundsSpanTheWholeMatchedSet pins the
// render's bounds fallback, the branch the Go producer cannot currently reach.
//
// It read filterMatchedRows[0].createdTime plus one hour. Rows are newest-first
// (REQ-318), so [0] is the NEWEST capture — the bounds collapsed to a one-hour
// window around it, and bounds are what every control clamps against, so no
// control could leave. On this repo's board that is 287 of 317 REQs permanently
// out of reach.
//
// The branch is UNREACHABLE from the producer today (timelineRange always returns
// real instants for a non-empty row set), which is exactly why it is worth a test:
// a fallback nobody exercises is a fallback nobody notices rotting, and this one
// rotted into the worst possible failure mode.
// TestJavaScriptBehaviorTimelineNoMatchStateRetiresTheToolbar pins the one render
// path that draws nothing.
//
// renderTimelineView returns early when the filters match no REQ, after
// releaseTimelineListeners but BEFORE the toolbar is wired. The toolbar is bound
// with `button.onclick =`, which is outside the teardown registry, so every
// handler from the previous render survived — holding that render's rows, its
// detached rows SVG and its renderAll. One press of Fit all then refilled the
// summary, the forecast and the details table with the REQs the filter had
// excluded, over a chart that stayed empty.
//
// The property: after a no-match render, no control this view owns has a handler
// or can be pressed; after a render that matches again, they all work.
// TestJavaScriptBehaviorTimelineSummaryCountsRowsDrawnAsBreaks drives the whole
// renderer because the contract spans two seams: the filtered row population
// chosen by renderTimelineView and every reason its drawing pass represents as
// broken. Counting causes would double-count REQ-914; counting the unfiltered
// payload would make the filtered case report four instead of one.
//
// REQ-328 CHANGED WHAT THIS COUNTS, and the change is deliberate and stated here
// because a quietly-edited expectation looks identical in a diff.
//
// The old rule was `row.anomaly || waitMinutes < 0 || workMinutes < 0`. row.anomaly
// is the board's BROADER bookkeeping verdict and includes rows whose spans are
// perfectly drawable, so the sentence announced breaks the chart never drew: nine
// such rows on this repo's own board produced "9 with broken stamps, drawn as
// breaks" over a chart with zero break markers on it. The clause now counts
// exactly what the drawing pass turns into a break, through the one predicate both
// read (timelineRowDrawsABreak).
//
// So REQ-911 — flagged anomalous, every span drawn — is now expected to count
// ZERO. It stays in the fixture as the guard against `row.anomaly ||` returning.
// REQ-916 is the shape that replaces it as a real cause: a REQ that STOPPED with
// no resolvable end instant, which carries a zero completedTime and no open flag
// and is drawn as a break. The unfiltered total is still 4, for a different set of
// four rows: 912, 913, 914 and 916.
// TestJavaScriptBehaviorTimelineForecastLabelsAFilteredView drives the WHOLE
// renderTimelineView, not renderTimelineForecast alone, because the defect lived
// in the wiring: rows were filtered, projection never was, and the call site
// handed the forecast a filtered row list it then ignored. A probe that calls
// the forecast function directly cannot tell a correct call site from one that
// always says "unfiltered" — this one can.
// ---- the done card's implementation span -----------------------------------

// implementationSpanPausedFixtureSpan is the over-ceiling case from the REQ's
// Red-Green Proof (an overnight pause). The helper below asserts it still
// exceeds the read-time ceiling, so moving the ceiling cannot quietly turn this
// fixture into an ordinary span that witnesses nothing.
const implementationSpanPausedFixtureSpan = 18 * time.Hour

// implementationSpanFixtureCommitHash is a plausible hash for the git-dated
// completion case. The stub lookup below is what dates it — no git runs.
const implementationSpanFixtureCommitHash = "0f1e2d3c4b5a69788796a5b4c3d2e1f009182736"

// spanFixtureFrontmatter writes one archived REQ. Stamps are omitted rather than
// written empty when blank, because an empty frontmatter value and an absent key
// are different inputs to the parser.
func spanFixtureFrontmatter(requestId string, title string, status string, claimedAt string, completedAt string, extraLines ...string) string {
	frontmatter := "---\nid: " + requestId + "\ntitle: " + title + "\nstatus: " + status + "\nuser_request: UR-900\n"
	if claimedAt != "" {
		frontmatter += "claimed_at: " + claimedAt + "\n"
	}
	if completedAt != "" {
		frontmatter += "completed_at: " + completedAt + "\n"
	}
	for _, extraLine := range extraLines {
		frontmatter += extraLine + "\n"
	}
	return frontmatter + "---\n\n# " + title + "\n"
}

// buildImplementationSpanFixturePayload projects the done-card span cases the
// REQ's Red-Green Proof names through the PRODUCTION pipeline — literal
// frontmatter, buildBoard, buildGeneratedBoardData — and returns the complete
// board payload. The payload assertions and the rendered-card probe both read
// this, so neither holds a hand-written copy of its field names or badge text.
func buildImplementationSpanFixturePayload(t *testing.T) generatedBoardData {
	t.Helper()
	if implementationSpanPausedFixtureSpan <= analysisOutlierCeiling {
		t.Fatalf("the paused fixture spans %v, which the read-time ceiling (%v) no longer excludes — that case would witness nothing",
			implementationSpanPausedFixtureSpan, analysisOutlierCeiling)
	}
	claimInstant := time.Date(2026, 8, 24, 10, 5, 0, 0, time.UTC)
	ordinaryCompletion := time.Date(2026, 8, 24, 12, 45, 0, 0, time.UTC)
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"do-work/archive/REQ-901-ordinary-span.md", spanFixtureFrontmatter(
			"REQ-901", "ordinary span", "completed",
			claimInstant.Format(time.RFC3339), ordinaryCompletion.Format(time.RFC3339),
			"planning_at: 2026-08-24T10:10:00Z",
			"dispatch_at: 2026-08-24T10:15:00Z",
			"builder_handback_at: 2026-08-24T12:00:00Z",
			"integration_at: 2026-08-24T12:10:00Z",
			"review_at: 2026-08-24T12:40:00Z",
			"release_at: 2026-08-24T13:00:00Z")},
		{"do-work/archive/REQ-902-paused-span.md", spanFixtureFrontmatter(
			"REQ-902", "overnight span", "completed",
			claimInstant.Format(time.RFC3339), claimInstant.Add(implementationSpanPausedFixtureSpan).Format(time.RFC3339))},
		{"do-work/archive/REQ-903-reversed-span.md", spanFixtureFrontmatter(
			"REQ-903", "reversed span", "completed",
			ordinaryCompletion.Format(time.RFC3339), claimInstant.Format(time.RFC3339))},
		{"do-work/archive/REQ-904-no-claim-stamp.md", spanFixtureFrontmatter(
			"REQ-904", "no claim stamp", "completed",
			"", ordinaryCompletion.Format(time.RFC3339))},
		{"do-work/archive/REQ-905-cancelled.md", spanFixtureFrontmatter(
			"REQ-905", "cancelled with both stamps", "cancelled",
			claimInstant.Format(time.RFC3339), ordinaryCompletion.Format(time.RFC3339))},
		// D-01's lock-in: this card's rendered completion instant comes from the
		// commit's git date, so the card HAS a done line but must state no span —
		// a commit delta is not an implementation span.
		{"do-work/archive/REQ-906-git-dated-completion.md", spanFixtureFrontmatter(
			"REQ-906", "git-dated completion", "completed",
			claimInstant.Format(time.RFC3339), "",
			"commit: "+implementationSpanFixtureCommitHash)},
		// The sub-hour case. It is the one that separates the card's stopwatch
		// vocabulary from the Durations table's chart vocabulary: 34 minutes reads
		// "34m 00s" here and "34.0 min" there.
		{"do-work/archive/REQ-907-sub-hour-span.md", spanFixtureFrontmatter(
			"REQ-907", "sub-hour span", "completed-with-issues",
			claimInstant.Format(time.RFC3339), claimInstant.Add(34*time.Minute).Format(time.RFC3339))},
		// The zero boundary. A REQ claimed and completed at the same instant has a
		// real span of zero, which is NOT the same state as an unmeasured one — the
		// distinction hasImplementationSpan exists to carry (D-06). It is the case an
		// `omitempty` float silently drops, leaving the client to multiply undefined
		// and draw "took NaNs".
		{"do-work/archive/REQ-908-zero-span.md", spanFixtureFrontmatter(
			"REQ-908", "claimed and completed at the same instant", "completed",
			claimInstant.Format(time.RFC3339), claimInstant.Format(time.RFC3339))},
	})
	gitDateLookupStub := func(_ string, commitHash string) (time.Time, bool) {
		if commitHash == implementationSpanFixtureCommitHash {
			return ordinaryCompletion, true
		}
		return time.Time{}, false
	}
	moment := time.Date(2026, 8, 26, 9, 0, 0, 0, time.UTC)
	board, buildError := buildBoard(repoRoot, moment, defaultRecentWindow, gitDateLookupStub)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}
	boardData, projectError := buildGeneratedBoardData(board)
	if projectError != nil {
		t.Fatalf("buildGeneratedBoardData: %v", projectError)
	}
	return boardData
}

// The span and its verdict are decided in Go and shipped in the per-request
// payload (D-02), so the client never restates the read-time ceiling. This pins
// what reaches the client for each case the REQ names.
func TestGeneratedRequestCarriesTheDoneCardImplementationSpan(t *testing.T) {
	boardData := buildImplementationSpanFixturePayload(t)
	requestsById := boardData.Requests
	if boardData.ImplementationSpanPausedBadgeText != "over 4h · assumed pause" {
		t.Fatalf("implementationSpanPausedBadgeText = %q, want the label derived from the current ceiling", boardData.ImplementationSpanPausedBadgeText)
	}

	spanExpectations := []struct {
		requestId   string
		hasSpan     bool
		wantMinutes float64
		wantReason  string
		requirement string
	}{
		{"REQ-901", true, 160, "", "both stamps parse and the span is inside the ceiling"},
		{"REQ-902", true, implementationSpanPausedFixtureSpan.Minutes(), "paused", "an over-ceiling span ships the verdict, never the ceiling"},
		{"REQ-903", true, -160, "reversed", "a reversed span ships raw and signed with the reversed verdict"},
		{"REQ-904", false, 0, "", "no claimed_at means no measurable span"},
		{"REQ-905", false, 0, "", "cancelled is terminal but not completed — the request was scoped to completed work"},
		{"REQ-906", false, 0, "", "a git-dated completion instant is not an implementation span (D-01)"},
		{"REQ-907", true, 34, "", "completed-with-issues is terminal success and states its span"},
		{"REQ-908", true, 0, "", "a zero-minute span is measured, not unmeasured — the flag says present, the value says zero"},
	}
	if len(requestsById) != len(spanExpectations) {
		t.Fatalf("payload holds %d requests, want %d — a fixture that never parsed asserts nothing", len(requestsById), len(spanExpectations))
	}

	for _, expectation := range spanExpectations {
		payload, present := requestsById[expectation.requestId]
		if !present {
			t.Fatalf("%s never reached the payload", expectation.requestId)
		}
		if payload.HasImplementationSpan != expectation.hasSpan {
			t.Errorf("%s hasImplementationSpan = %v, want %v (%s)",
				expectation.requestId, payload.HasImplementationSpan, expectation.hasSpan, expectation.requirement)
		}
		if payload.ImplementationSpanMinutes != expectation.wantMinutes {
			t.Errorf("%s implementationSpanMinutes = %v, want %v (%s)",
				expectation.requestId, payload.ImplementationSpanMinutes, expectation.wantMinutes, expectation.requirement)
		}
		if payload.ImplementationSpanReason != expectation.wantReason {
			t.Errorf("%s implementationSpanReason = %q, want %q (%s)",
				expectation.requestId, payload.ImplementationSpanReason, expectation.wantReason, expectation.requirement)
		}
	}

	// The D-01 case is only about the SPAN: the card still renders a completion
	// instant, and it came from git. Without this the case would pass for a
	// fixture that simply failed to produce a done line at all.
	gitDatedPayload := requestsById["REQ-906"]
	if gitDatedPayload.CompletionTime == "" {
		t.Fatalf("REQ-906 resolved no completion instant, so it cannot witness the git-dated case")
	}
	if gitDatedPayload.CompletionTimeSource != string(CompletionFromGitLog) {
		t.Fatalf("REQ-906 completionTimeSource = %q, want %q", gitDatedPayload.CompletionTimeSource, CompletionFromGitLog)
	}

	// The zero span has to survive MARSHALLING, not just the struct: an omitempty
	// float drops a genuine 0 from the wire while hasImplementationSpan still ships
	// true, and the client then multiplies undefined into NaN. Asserting the struct
	// field alone would pass with that tag restored, so this reads the JSON.
	zeroSpanJson, marshalError := json.Marshal(requestsById["REQ-908"])
	if marshalError != nil {
		t.Fatalf("marshal REQ-908: %v", marshalError)
	}
	if !strings.Contains(string(zeroSpanJson), `"implementationSpanMinutes":0`) {
		t.Errorf("REQ-908 marshalled without its zero span: %s\n"+
			"a present-but-zero span must reach the client as 0, or the card renders \"took NaNs\"", zeroSpanJson)
	}
}

func TestGeneratedRequestCarriesObservedPhaseBreakdown(t *testing.T) {
	boardData := buildImplementationSpanFixturePayload(t)
	request := boardData.Requests["REQ-901"]
	if request.PlanningAt != "2026-08-24T10:10:00Z" || request.ReleaseAt != "2026-08-24T13:00:00Z" {
		t.Fatalf("raw phase observations did not reach generated request: %#v", request)
	}
	wantFields := []string{
		"planning_at", "dispatch_at", "builder_handback_at", "integration_at", "review_at", "completed_at", "release_at",
	}
	if len(request.PhaseBreakdown) != len(wantFields) {
		t.Fatalf("phase payload has %d entries, want %d: %#v", len(request.PhaseBreakdown), len(wantFields), request.PhaseBreakdown)
	}
	for index, wantField := range wantFields {
		if request.PhaseBreakdown[index].FieldName != wantField {
			t.Errorf("phase payload entry %d field = %q, want %q", index, request.PhaseBreakdown[index].FieldName, wantField)
		}
	}
	if historical := boardData.Requests["REQ-902"].PhaseBreakdown; len(historical) != 0 {
		t.Fatalf("historical request without phase observations gained a breakdown: %#v", historical)
	}
}

// The rendered done line, built by the REAL makeRequestCard out of the generated
// page and fed the REAL payload. The point of the whole REQ is what the card
// SAYS, so a stubbed card would assert nothing: the instant node is the one
// stub, because its text has its own coverage and this probe is about the span
// fragment and where it lands.
// clipboardProbeDocument renders one fixture as the { text, ticketMentions } pair
// a Copy handler hands the client, with the index computed by the SAME Go walk
// the build ships. That is what makes the probe an end-to-end check of the two
// halves together rather than of the client alone: a walk that mislocates a
// mention shows up here as a wrong payload, not as a passing client test.
func clipboardProbeDocument(t *testing.T, documentText string, resolver *ticketMentionResolver) string {
	t.Helper()
	ticketMentions, encodeError := json.Marshal(collectDocumentTicketMentions(documentText, resolver))
	if encodeError != nil {
		t.Fatalf("encode probe ticket mentions: %v", encodeError)
	}
	return "{ text: " + mustMarshalJSONString(t, documentText) + ", ticketMentions: " + string(ticketMentions) + " }"
}
