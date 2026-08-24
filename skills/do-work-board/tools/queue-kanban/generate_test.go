package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"testing/fstest"
	"time"
)

const (
	strictJavaScriptBehaviorDiagnostic = "queue-kanban: strict JavaScript behavior lane executed zero probes"
	strictJavaScriptBehaviorMarker     = "QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR"
	strictJavaScriptBehaviorRunPattern = "^TestMaintainerStrictJavaScriptBehaviorLane$"
)

var javaScriptBehaviorProbeCount atomic.Int64

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
	if exitCode == 0 && os.Getenv(strictJavaScriptBehaviorMarker) == "1" && javaScriptBehaviorProbeCount.Load() == 0 {
		fmt.Fprintln(os.Stderr, strictJavaScriptBehaviorDiagnostic)
		exitCode = 1
	}
	// Same guard for the browser lane (browser_probe_test.go): a strict run whose
	// probes all skipped must not report green, which is what makes the ordinary
	// skip safe. Both lanes gate here because TestMain is per-package, not per-file.
	if exitCode == 0 && os.Getenv(strictBrowserBehaviorMarker) == "1" && browserBehaviorProbeCount.Load() == 0 {
		fmt.Fprintln(os.Stderr, strictBrowserBehaviorDiagnostic)
		exitCode = 1
	}
	os.Exit(exitCode)
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

func TestMaintainerStrictJavaScriptBehaviorLaneRejectsZeroProbes(t *testing.T) {
	strictCommand := exec.Command(os.Args[0], "-test.run=^TestJavaScriptBehavior", "-test.count=1")
	strictCommand.Env = testEnvironmentWithOverrides(
		os.Environ(),
		"PATH="+t.TempDir(),
		strictJavaScriptBehaviorMarker+"=1",
	)
	strictOutput, strictError := strictCommand.CombinedOutput()
	if strictError == nil {
		t.Fatalf("strict JavaScript behavior lane exited zero without Node; output:\n%s", strictOutput)
	}
	if !strings.Contains(string(strictOutput), strictJavaScriptBehaviorDiagnostic) {
		t.Fatalf("strict JavaScript behavior lane output = %q, want %q", strictOutput, strictJavaScriptBehaviorDiagnostic)
	}
}

func TestMaintainerStrictJavaScriptBehaviorLane(t *testing.T) {
	testRunFlag := flag.Lookup("test.run")
	if testRunFlag == nil || testRunFlag.Value.String() != strictJavaScriptBehaviorRunPattern {
		t.Skip("maintainer strict JavaScript behavior lane runs only when selected directly")
	}

	strictCommand := exec.Command(os.Args[0], "-test.run=^TestJavaScriptBehavior", "-test.count=1")
	strictCommand.Env = testEnvironmentWithOverrides(
		os.Environ(),
		strictJavaScriptBehaviorMarker+"=1",
	)
	strictOutput, strictError := strictCommand.CombinedOutput()
	if strictError != nil {
		t.Fatalf("strict JavaScript behavior lane failed: %v\n%s", strictError, strictOutput)
	}
}

func lookupNodeForJavaScriptProbe(t *testing.T) string {
	t.Helper()
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
func TestJavaScriptBehaviorAssembledClientSyntax(t *testing.T) {
	behaviorProbeCountBefore := javaScriptBehaviorProbeCount.Load()
	defer func() {
		if behaviorProbeCountAfter := javaScriptBehaviorProbeCount.Load(); behaviorProbeCountAfter != behaviorProbeCountBefore {
			t.Errorf("assembled syntax changed behavior-probe count from %d to %d",
				behaviorProbeCountBefore, behaviorProbeCountAfter)
		}
	}()
	assembledClient, assembleError := assembleBoardJavaScript(embeddedWebAssets)
	if assembleError != nil {
		t.Fatalf("assembleBoardJavaScript: %v", assembleError)
	}
	nodePath := lookupNodeForJavaScriptProbe(t)
	syntaxCommand := exec.Command(nodePath, "--check", "-")
	syntaxCommand.Stdin = bytes.NewReader(assembledClient)
	syntaxOutput, syntaxError := syntaxCommand.CombinedOutput()
	if syntaxError != nil {
		t.Fatalf("node --check assembled client: %v\n%s", syntaxError, syntaxOutput)
	}
}

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
	outputDirectory := generateLiveSiteInDir(t)
	indexPath := filepath.Join(outputDirectory, "index.html")
	indexBytes, readError := os.ReadFile(indexPath)
	if readError != nil {
		t.Fatalf("reading generated index.html: %v", readError)
	}
	return string(indexBytes)
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
	// so the primary Copy path writes them verbatim — no synthesized heading, or the
	// paste stops round-tripping back into a valid REQ file. The identifying heading
	// belongs to the rendered-text fallback alone, which has no frontmatter to carry.
	if !strings.Contains(string(indexBytes), "copyTextWithHeading(requestedKind, requestedId, renderedTextFallback)") {
		t.Fatalf("inlined board-clipboard.js does not prepend the id/title heading on the rendered-text fallback path")
	}
	if strings.Contains(string(indexBytes), "copyTextWithHeading(requestedKind, requestedId, bodyText)") {
		t.Fatalf("inlined board-clipboard.js still routes the lazy payload through the heading builder — the primary path must copy the file verbatim")
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
func TestJavaScriptBehaviorDrawerHeadingDeduplication(t *testing.T) {
	indexHtml := generateLiveSite(t)
	functionBlocks := []string{
		sliceBalancedBlockAfter(t, indexHtml, "function normalizeHeadingText("),
		sliceBalancedBlockAfter(t, indexHtml, "function linkifyDetailBody("),
	}
	javascriptProbe := `
var NodeFilter = { SHOW_TEXT: 4 };
var document = {
  createTreeWalker: function () {
    return { nextNode: function () { return false; } };
  }
};
function makeDrawerBody(headingText) {
  var bodyRoot = {
    firstElementChild: null,
    querySelectorAll: function () { return []; }
  };
  var heading = {
    tagName: "H1",
    textContent: headingText,
    removed: false,
    remove: function () {
      this.removed = true;
      bodyRoot.firstElementChild = null;
    }
  };
  bodyRoot.firstElementChild = heading;
  return { root: bodyRoot, heading: heading };
}
` + strings.Join(functionBlocks, "\n") + `
var requestMatch = makeDrawerBody("  Compile   assets ");
var requestMismatch = makeDrawerBody("Implementation notes");
var userRequestMatch = makeDrawerBody("Launch plan");
var userRequestMismatch = makeDrawerBody("Background");
linkifyDetailBody(requestMatch.root, "Compile assets");
linkifyDetailBody(requestMismatch.root, "Compile assets");
linkifyDetailBody(userRequestMatch.root, "LAUNCH PLAN");
linkifyDetailBody(userRequestMismatch.root, "Launch plan");
process.stdout.write(JSON.stringify([
  requestMatch.heading.removed,
  requestMismatch.heading.removed,
  userRequestMatch.heading.removed,
  userRequestMismatch.heading.removed
]));`
	probeOutput := runJavaScriptBehaviorProbe(t, "drawer heading", javascriptProbe)
	var removedResults []bool
	if decodeError := json.Unmarshal(probeOutput, &removedResults); decodeError != nil {
		t.Fatalf("decode drawer heading behavior: %v (output %q)", decodeError, probeOutput)
	}
	wantedResults := []bool{true, false, true, false}
	if len(removedResults) != len(wantedResults) {
		t.Fatalf("drawer heading result count = %d, want %d: %#v", len(removedResults), len(wantedResults), removedResults)
	}
	for resultIndex := range wantedResults {
		if removedResults[resultIndex] != wantedResults[resultIndex] {
			t.Fatalf("drawer heading result[%d] = %v, want %v; all results=%#v",
				resultIndex, removedResults[resultIndex], wantedResults[resultIndex], removedResults)
		}
	}
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
func TestJavaScriptBehaviorByUserRequestLensCountsRecentlyDoneAsActive(t *testing.T) {
	indexHtml := generateLiveSite(t)
	for _, requiredToken := range []string{
		"userRequestHasOpenOrRecentWork(userRequest, recentlyDoneIdSet)",
		"recentlyDoneIds(viewState.windowHours)",
	} {
		if !strings.Contains(indexHtml, requiredToken) {
			t.Fatalf("by-UR recent-work behavior is not wired into the generated asset: %q missing", requiredToken)
		}
	}

	functionBlocks := []string{
		sliceBalancedBlockAfter(t, indexHtml, "function isTerminalResolvedStatus("),
		sliceBalancedBlockAfter(t, indexHtml, "function userRequestHasOpenOrRecentWork("),
		sliceBalancedBlockAfter(t, indexHtml, "function recentlyDoneIds("),
	}
	// The calendar carries EVERY REQ, so the stub holds a claimed-an-hour-ago and
	// a failed-an-hour-ago entry alongside the completions. Both are inside the
	// 24h window and neither may reach Recently done: recentlyDoneIds is gated on
	// terminal-resolved, not on mere presence in the array.
	javascriptProbe := `
Date.now = function () { return Date.parse("2026-08-15T12:00:00Z"); };
var boardData = { calendar: [
  { id: "REQ-claimed", status: "claimed", entryTime: "2026-08-15T11:00:00Z" },
  { id: "REQ-failed", status: "failed", entryTime: "2026-08-15T11:00:00Z" },
  { id: "REQ-queued", status: "pending", entryTime: "" },
  { id: "REQ-recent", status: "completed", entryTime: "2026-08-15T06:00:00Z" },
  { id: "REQ-old", status: "completed", entryTime: "2026-08-13T06:00:00Z" }
] };
var requestsById = {
  "REQ-recent": { status: "completed" },
  "REQ-old": { status: "completed" },
  "REQ-open": { status: "pending" }
};
` + strings.Join(functionBlocks, "\n") + `
var recentIds = recentlyDoneIds(24);
var recentlyDoneIdSet = {};
recentIds.forEach(function (requestId) { recentlyDoneIdSet[requestId] = true; });
process.stdout.write(JSON.stringify({
  recentIds: recentIds,
  recentActive: userRequestHasOpenOrRecentWork({ requestIds: ["REQ-recent"] }, recentlyDoneIdSet),
  oldActive: userRequestHasOpenOrRecentWork({ requestIds: ["REQ-old"] }, recentlyDoneIdSet),
  openActive: userRequestHasOpenOrRecentWork({ requestIds: ["REQ-open"] }, recentlyDoneIdSet)
}));`
	probeOutput := runJavaScriptBehaviorProbe(t, "by-UR recent-work predicate", javascriptProbe)
	var result struct {
		RecentIds    []string `json:"recentIds"`
		RecentActive bool     `json:"recentActive"`
		OldActive    bool     `json:"oldActive"`
		OpenActive   bool     `json:"openActive"`
	}
	if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
		t.Fatalf("decode by-UR recent-work result: %v (output %q)", decodeError, probeOutput)
	}
	if len(result.RecentIds) != 1 || result.RecentIds[0] != "REQ-recent" {
		t.Fatalf("recentlyDoneIds(24) = %#v, want only REQ-recent — a claimed, failed, or queued calendar "+
			"entry inside the window must never reach the Recently done column", result.RecentIds)
	}
	if !result.RecentActive || result.OldActive || !result.OpenActive {
		t.Fatalf("Active predicate result = recent:%v old:%v open:%v, want true, false, true",
			result.RecentActive, result.OldActive, result.OpenActive)
	}
}

func TestJavaScriptBehaviorRecentlyDoneWindowRefreshesVisibleLens(t *testing.T) {
	indexHtml := generateLiveSite(t)
	const wiringToken = `applyRecentWindowSelection(parseInt(button.getAttribute("data-window-hours"), 10))`
	if !strings.Contains(indexHtml, wiringToken) {
		t.Fatalf("recent-window click handler is not wired to the transition helper: %q missing", wiringToken)
	}

	recentWindowFunction := sliceBalancedBlockAfter(t, indexHtml, "function applyRecentWindowSelection(")
	javascriptProbe := `
var viewState = { windowHours: 24, view: "board", lens: "user-request" };
var renderedOnce = { userRequestLens: true };
var selectedWindow = "";
var columnRenderCount = 0;
var lensRenderCount = 0;
function setActiveButton(selector, attributeName, attributeValue) { selectedWindow = attributeValue; }
function renderColumns() { columnRenderCount += 1; }
function renderUserRequestLens() { lensRenderCount += 1; }
` + recentWindowFunction + `
applyRecentWindowSelection(168);
var visibleLensState = {
  windowHours: viewState.windowHours,
  selectedWindow: selectedWindow,
  columnRenderCount: columnRenderCount,
  lensRenderCount: lensRenderCount,
  lensFresh: renderedOnce.userRequestLens
};
viewState.lens = "columns";
applyRecentWindowSelection(48);
process.stdout.write(JSON.stringify({
  visibleLensState: visibleLensState,
  hiddenLensState: {
    windowHours: viewState.windowHours,
    columnRenderCount: columnRenderCount,
    lensRenderCount: lensRenderCount,
    lensFresh: renderedOnce.userRequestLens
  }
}));`
	probeOutput := runJavaScriptBehaviorProbe(t, "recent-window transition", javascriptProbe)
	var result struct {
		VisibleLensState struct {
			WindowHours       int    `json:"windowHours"`
			SelectedWindow    string `json:"selectedWindow"`
			ColumnRenderCount int    `json:"columnRenderCount"`
			LensRenderCount   int    `json:"lensRenderCount"`
			LensFresh         bool   `json:"lensFresh"`
		} `json:"visibleLensState"`
		HiddenLensState struct {
			WindowHours       int  `json:"windowHours"`
			ColumnRenderCount int  `json:"columnRenderCount"`
			LensRenderCount   int  `json:"lensRenderCount"`
			LensFresh         bool `json:"lensFresh"`
		} `json:"hiddenLensState"`
	}
	if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
		t.Fatalf("decode recent-window transition: %v (output %q)", decodeError, probeOutput)
	}
	visibleState := result.VisibleLensState
	if visibleState.WindowHours != 168 || visibleState.SelectedWindow != "168" || visibleState.ColumnRenderCount != 1 || visibleState.LensRenderCount != 1 || !visibleState.LensFresh {
		t.Fatalf("visible by-UR transition = %#v, want selected window 168 with both lenses refreshed", visibleState)
	}
	hiddenState := result.HiddenLensState
	if hiddenState.WindowHours != 48 || hiddenState.ColumnRenderCount != 2 || hiddenState.LensRenderCount != 1 || hiddenState.LensFresh {
		t.Fatalf("hidden by-UR transition = %#v, want columns refreshed and by-UR marked stale", hiddenState)
	}
}

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

func TestJavaScriptBehaviorDurationsWindowSelectionRefreshesOnlyDurations(t *testing.T) {
	indexHTML := generateLiveSite(t)
	transitionFunction := sliceBalancedBlockAfter(t, indexHTML, "function applyDurationsWindowSelection(")
	javascriptProbe := `
var viewState = { windowHours: 24, view: "durations" };
var renderedOnce = { durations: true };
var chosenWindows = [];
var activeWindows = [];
var renderCount = 0;
function setDurationsWindow(windowName) { chosenWindows.push(windowName); }
function setActiveButton(selector, attributeName, attributeValue) { activeWindows.push(attributeValue); }
function renderDurationsView() { renderCount += 1; }
` + transitionFunction + `
["30", "90", "all"].forEach(function (windowName) { applyDurationsWindowSelection(windowName); });
var visibleState = { renderCount: renderCount, rendered: renderedOnce.durations };
viewState.view = "board";
applyDurationsWindowSelection("30");
process.stdout.write(JSON.stringify({
  chosenWindows: chosenWindows,
  activeWindows: activeWindows,
  windowHours: viewState.windowHours,
  visibleState: visibleState,
  hiddenRenderCount: renderCount,
  hiddenRendered: renderedOnce.durations
}));`
	probeOutput := runJavaScriptBehaviorProbe(t, "Durations window transitions", javascriptProbe)
	var result struct {
		ChosenWindows []string `json:"chosenWindows"`
		ActiveWindows []string `json:"activeWindows"`
		WindowHours   int      `json:"windowHours"`
		VisibleState  struct {
			RenderCount int  `json:"renderCount"`
			Rendered    bool `json:"rendered"`
		} `json:"visibleState"`
		HiddenRenderCount int  `json:"hiddenRenderCount"`
		HiddenRendered    bool `json:"hiddenRendered"`
	}
	if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
		t.Fatalf("decode Durations-window transition: %v (output %q)", decodeError, probeOutput)
	}
	wantWindows := []string{"30", "90", "all", "30"}
	if strings.Join(result.ChosenWindows, ",") != strings.Join(wantWindows, ",") ||
		strings.Join(result.ActiveWindows, ",") != strings.Join(wantWindows, ",") {
		t.Fatalf("Durations selections = chosen %#v active %#v, want %#v", result.ChosenWindows, result.ActiveWindows, wantWindows)
	}
	if result.WindowHours != 24 {
		t.Fatalf("Durations selection changed viewState.windowHours to %d, want unchanged 24", result.WindowHours)
	}
	if result.VisibleState.RenderCount != 3 || !result.VisibleState.Rendered {
		t.Fatalf("three visible selections produced %#v, want one render per selection and fresh state", result.VisibleState)
	}
	if result.HiddenRenderCount != 3 || result.HiddenRendered {
		t.Fatalf("hidden selection produced renderCount=%d rendered=%v, want no render and stale cache", result.HiddenRenderCount, result.HiddenRendered)
	}
}

// Execute the pure empty-state decision under Node so the regression is pinned
// to state transitions rather than the presence of reassuring source strings.
func TestJavaScriptBehaviorByUserRequestLensEmptyState(t *testing.T) {
	indexHtml := generateLiveSite(t)
	functionBlocks := []string{
		sliceBalancedBlockAfter(t, indexHtml, "function recentWindowPhrase("),
		sliceBalancedBlockAfter(t, indexHtml, "function userRequestLensEmptyText("),
	}
	javascriptProbe := strings.Join(functionBlocks, "\n") + `
const results = [
  recentWindowPhrase(1),
  recentWindowPhrase(168),
  userRequestLensEmptyText(true, 4, 2, "the last 24 hours"),
  userRequestLensEmptyText(true, 4, 0, "the last 24 hours"),
  userRequestLensEmptyText(false, 4, 0, "the last 24 hours"),
  userRequestLensEmptyText(false, 0, 0, "the last 24 hours")
];
process.stdout.write(JSON.stringify(results));`
	probeOutput := runJavaScriptBehaviorProbe(t, "by-UR empty-state decision", javascriptProbe)
	var results []string
	if decodeError := json.Unmarshal(probeOutput, &results); decodeError != nil {
		t.Fatalf("decode assembled-client empty-state results: %v (output %q)", decodeError, probeOutput)
	}
	if len(results) != 6 {
		t.Fatalf("empty-state result count = %d, want 6", len(results))
	}
	if results[0] != "the last 1 hour" || results[1] != "the last 7 days" {
		t.Fatalf("recent-window phrases = %q, %q, want singular hour and seven-day copy", results[0], results[1])
	}
	if !strings.Contains(results[2], "switch URs to All") || !strings.Contains(results[2], "2 resolved matches") {
		t.Fatalf("scope-hidden search result = %q, want an All-scope escape with the match count", results[2])
	}
	if results[3] != "No user requests match the current filters." {
		t.Fatalf("genuine filter miss = %q, want the generic no-match message", results[3])
	}
	if !strings.Contains(results[4], "widen the RECENTLY DONE window") || !strings.Contains(results[4], "switch URs to All") {
		t.Fatalf("scope-only empty state = %q, want both scope escapes", results[4])
	}
	if results[5] != "No user requests in this tree yet." {
		t.Fatalf("empty tree state = %q, want the empty-tree message", results[5])
	}
}

// Exercise the production lens caller, not only its pure predicate: a terminal
// REQ inside the selected window renders while an older terminal REQ stays
// hidden and is counted in the scope note.
func TestJavaScriptBehaviorByUserRequestLensUsesRecentWindowAtCaller(t *testing.T) {
	indexHtml := generateLiveSite(t)
	functionBlocks := []string{
		sliceBalancedBlockAfter(t, indexHtml, "function createElement("),
		sliceBalancedBlockAfter(t, indexHtml, "function isTerminalResolvedStatus("),
		sliceBalancedBlockAfter(t, indexHtml, "function hasActiveFilters("),
		sliceBalancedBlockAfter(t, indexHtml, "function searchMatchesRequest("),
		sliceBalancedBlockAfter(t, indexHtml, "function searchMatchesUserRequest("),
		sliceBalancedBlockAfter(t, indexHtml, "function requestMatchesFilters("),
		sliceBalancedBlockAfter(t, indexHtml, "function userRequestHasOpenOrRecentWork("),
		sliceBalancedBlockAfter(t, indexHtml, "function recentWindowPhrase("),
		sliceBalancedBlockAfter(t, indexHtml, "function userRequestLensEmptyText("),
		sliceBalancedBlockAfter(t, indexHtml, "function recentlyDoneIds("),
		sliceBalancedBlockAfter(t, indexHtml, "function renderUserRequestLens("),
	}
	javascriptProbe := `
Date.now = function () { return Date.parse("2026-08-15T12:00:00Z"); };
var boardData = {
  requests: {
    "REQ-501": { status: "completed", title: "old work" },
    "REQ-502": { status: "completed", title: "recent work" }
  },
  userRequests: {
    "UR-301": { requestIds: ["REQ-501"], title: "old request", inputFilePresent: true },
    "UR-302": { requestIds: ["REQ-502"], title: "recent request", inputFilePresent: true }
  },
  userRequestOrder: ["UR-301", "UR-302"],
  calendar: [
    { id: "REQ-502", status: "completed", entryTime: "2026-08-15T06:00:00Z" },
    { id: "REQ-501", status: "completed", entryTime: "2026-08-13T06:00:00Z" }
  ]
};
var requestsById = boardData.requests;
var userRequestsById = boardData.userRequests;
var viewState = { windowHours: 24 };
var filterState = { searchText: "", domain: "", status: "", userRequestActivity: "active" };
// The always-open reading of this lens; the folded one has its own probe.
var userRequestCardsFolded = false;
function makeNode() {
  return {
    childNodes: [],
    dataset: {},
    appendChild: function (childNode) { this.childNodes.push(childNode); return childNode; }
  };
}
var userRequestLensNode = makeNode();
var document = {
  getElementById: function (nodeId) { return nodeId === "user-request-lens" ? userRequestLensNode : null; },
  createElement: function () { return makeNode(); }
};
function makeRequestCard(requestId) { return { requestId: requestId }; }
` + strings.Join(functionBlocks, "\n") + `
renderUserRequestLens();
var renderedUserRequestIds = userRequestLensNode.childNodes
  .filter(function (node) { return node.className === "ur-group"; })
  .map(function (groupNode) { return groupNode.childNodes[0].dataset.detailId; });
var scopeNotes = userRequestLensNode.childNodes
  .filter(function (node) { return node.className === "ur-lens-hidden-note"; })
  .map(function (node) { return node.textContent; });
process.stdout.write(JSON.stringify({ renderedUserRequestIds: renderedUserRequestIds, scopeNotes: scopeNotes }));
`
	probeOutput := runJavaScriptBehaviorProbe(t, "by-UR caller", javascriptProbe)
	var result struct {
		RenderedUserRequestIds []string `json:"renderedUserRequestIds"`
		ScopeNotes             []string `json:"scopeNotes"`
	}
	if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
		t.Fatalf("decode assembled-client by-UR caller output: %v (output %q)", decodeError, probeOutput)
	}
	if len(result.RenderedUserRequestIds) != 1 || result.RenderedUserRequestIds[0] != "UR-302" {
		t.Fatalf("Active by-UR caller rendered %#v, want only recent terminal UR-302", result.RenderedUserRequestIds)
	}
	if len(result.ScopeNotes) != 1 || !strings.Contains(result.ScopeNotes[0], "1 UR with no open work or activity in the last 24 hours") {
		t.Fatalf("Active by-UR caller scope notes = %#v, want one old hidden UR and the selected window", result.ScopeNotes)
	}
}

func TestJavaScriptBehaviorTestingDoneWindowIsViewSpecific(t *testing.T) {
	indexHtml := generateLiveSite(t)
	functionBlocks := []string{
		sliceBalancedBlockAfter(t, indexHtml, "function createElement("),
		sliceBalancedBlockAfter(t, indexHtml, "function hasActiveFilters("),
		sliceBalancedBlockAfter(t, indexHtml, "function hasActiveVisibleFilters("),
		sliceBalancedBlockAfter(t, indexHtml, "function formatFilteredCount("),
		sliceBalancedBlockAfter(t, indexHtml, "function columnEmptyText("),
		sliceBalancedBlockAfter(t, indexHtml, "function fillColumn("),
		sliceBalancedBlockAfter(t, indexHtml, "function fillTestingColumn("),
	}
	javascriptProbe := `
var filterState = { searchText: "", domain: "", status: "", doneWindow: "168" };
var viewState = { view: "board" };
var nodesBySelector = {};
function makeNode() {
  return {
    childNodes: [],
    textContent: "",
    appendChild: function (childNode) { this.childNodes.push(childNode); return childNode; }
  };
}
var document = {
  createElement: function () { return makeNode(); },
  querySelector: function (selector) {
    if (!nodesBySelector[selector]) {
      nodesBySelector[selector] = makeNode();
    }
    return nodesBySelector[selector];
  }
};
` + strings.Join(functionBlocks, "\n") + `
fillColumn("board", [], null, 1);
var boardCopy = nodesBySelector['[data-cards="board"]'].childNodes[0].textContent;
viewState.view = "testing";
fillColumn("hidden-board", [], null, 1);
var hiddenBoardCopy = nodesBySelector['[data-cards="hidden-board"]'].childNodes[0].textContent;
fillTestingColumn("testing-ready", [], 1);
var testingCopy = nodesBySelector['[data-cards="testing-ready"]'].childNodes[0].textContent;
process.stdout.write(JSON.stringify([boardCopy, hiddenBoardCopy, testingCopy]));`
	probeOutput := runJavaScriptBehaviorProbe(t, "testing empty-copy decision", javascriptProbe)
	var results []string
	if decodeError := json.Unmarshal(probeOutput, &results); decodeError != nil {
		t.Fatalf("decode assembled-client empty-copy results: %v (output %q)", decodeError, probeOutput)
	}
	if len(results) != 3 {
		t.Fatalf("empty-copy result count = %d, want 3: %#v", len(results), results)
	}
	if results[0] != "Nothing here" {
		t.Fatalf("Board empty copy with only doneWindow = %q, want Nothing here", results[0])
	}
	if results[1] != "Nothing here" {
		t.Fatalf("hidden Board empty copy during Testing view = %q, want Nothing here", results[1])
	}
	if results[2] != "No matches" {
		t.Fatalf("Testing empty copy with doneWindow = %q, want No matches", results[2])
	}
}

func TestJavaScriptBehaviorTestingStatusUpdateInvalidatesUserRequestLens(t *testing.T) {
	indexHtml := generateLiveSite(t)
	postTestingSource := sliceBalancedBlockAfter(t, indexHtml, "function postTestingStatus(")
	updateCallback := sliceBalancedBlockAfter(t, postTestingSource, ".then(function (payload) {")
	const wiringToken = "applyConfirmedTestingTransition(requestId, testingState, feedbackText, payload)"
	if !strings.Contains(updateCallback, wiringToken) {
		t.Fatalf("testing-status success callback is not wired to its confirmed transition: %q missing", wiringToken)
	}

	transitionFunction := sliceBalancedBlockAfter(t, indexHtml, "function applyConfirmedTestingTransition(")
	javascriptProbe := `
var requestsById = {
  "REQ-1": {
    testingStatus: "",
    testedBy: "",
    testingUpdatedAt: "",
    testingFeedback: "",
    testingStatusUnrecognized: true,
    originalTestingStatus: "invalid"
  }
};
var feedbackFormRequestId = "REQ-1";
var feedbackDraftText = "draft";
var renderedOnce = { userRequestLens: true };
var viewState = { view: "board", lens: "user-request" };
var testingRenderCount = 0;
var columnRenderCount = 0;
var lensRenderCount = 0;
function renderTestingView() { testingRenderCount += 1; }
function renderColumns() { columnRenderCount += 1; }
function renderUserRequestLens() { lensRenderCount += 1; }
` + transitionFunction + `
applyConfirmedTestingTransition("REQ-1", "returned", "needs revision", {
  testingStatus: "returned",
  testedBy: "Alex",
  testingUpdatedAt: "2026-08-15T12:00:00Z"
});
var visibleTransition = {
  request: Object.assign({}, requestsById["REQ-1"]),
  feedbackFormRequestId: feedbackFormRequestId,
  feedbackDraftText: feedbackDraftText,
  testingRenderCount: testingRenderCount,
  columnRenderCount: columnRenderCount,
  lensRenderCount: lensRenderCount,
  lensFresh: renderedOnce.userRequestLens
};
viewState.lens = "columns";
renderedOnce.userRequestLens = true;
applyConfirmedTestingTransition("REQ-1", "tested", "", {
  testingStatus: "tested",
  testedBy: "Alex",
  testingUpdatedAt: "2026-08-15T12:05:00Z"
});
process.stdout.write(JSON.stringify({
  visibleTransition: visibleTransition,
  hiddenLensFresh: renderedOnce.userRequestLens,
  hiddenLensRenderCount: lensRenderCount,
  hiddenTestingRenderCount: testingRenderCount,
  hiddenColumnRenderCount: columnRenderCount
}));`
	probeOutput := runJavaScriptBehaviorProbe(t, "confirmed testing transition", javascriptProbe)
	var result struct {
		VisibleTransition struct {
			Request struct {
				TestingStatus             string `json:"testingStatus"`
				TestedBy                  string `json:"testedBy"`
				TestingUpdatedAt          string `json:"testingUpdatedAt"`
				TestingFeedback           string `json:"testingFeedback"`
				TestingStatusUnrecognized bool   `json:"testingStatusUnrecognized"`
				OriginalTestingStatus     string `json:"originalTestingStatus"`
			} `json:"request"`
			FeedbackFormRequestId *string `json:"feedbackFormRequestId"`
			FeedbackDraftText     string  `json:"feedbackDraftText"`
			TestingRenderCount    int     `json:"testingRenderCount"`
			ColumnRenderCount     int     `json:"columnRenderCount"`
			LensRenderCount       int     `json:"lensRenderCount"`
			LensFresh             bool    `json:"lensFresh"`
		} `json:"visibleTransition"`
		HiddenLensFresh          bool `json:"hiddenLensFresh"`
		HiddenLensRenderCount    int  `json:"hiddenLensRenderCount"`
		HiddenTestingRenderCount int  `json:"hiddenTestingRenderCount"`
		HiddenColumnRenderCount  int  `json:"hiddenColumnRenderCount"`
	}
	if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
		t.Fatalf("decode confirmed testing transition: %v (output %q)", decodeError, probeOutput)
	}
	visibleTransition := result.VisibleTransition
	request := visibleTransition.Request
	if request.TestingStatus != "returned" || request.TestedBy != "Alex" || request.TestingUpdatedAt != "2026-08-15T12:00:00Z" || request.TestingFeedback != "needs revision" || request.TestingStatusUnrecognized || request.OriginalTestingStatus != "returned" {
		t.Fatalf("confirmed testing request state = %#v, want server-confirmed returned state", request)
	}
	if visibleTransition.FeedbackFormRequestId != nil || visibleTransition.FeedbackDraftText != "" {
		t.Fatalf("confirmed testing feedback form state = id:%v draft:%q, want cleared", visibleTransition.FeedbackFormRequestId, visibleTransition.FeedbackDraftText)
	}
	if visibleTransition.TestingRenderCount != 1 || visibleTransition.ColumnRenderCount != 1 || visibleTransition.LensRenderCount != 1 || !visibleTransition.LensFresh {
		t.Fatalf("confirmed testing render state = testing:%d columns:%d lens:%d fresh:%v, want one refresh for each visible surface",
			visibleTransition.TestingRenderCount, visibleTransition.ColumnRenderCount, visibleTransition.LensRenderCount, visibleTransition.LensFresh)
	}
	if result.HiddenLensFresh || result.HiddenLensRenderCount != 1 || result.HiddenTestingRenderCount != 2 || result.HiddenColumnRenderCount != 2 {
		t.Fatalf("hidden by-UR render state = testing:%d columns:%d lens:%d fresh:%v, want lens uncalled and marked stale",
			result.HiddenTestingRenderCount, result.HiddenColumnRenderCount, result.HiddenLensRenderCount, result.HiddenLensFresh)
	}
}

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

// Band-and-row geometry, and the remainder sentence's all-or-nothing rule.
//
// Before REQ-292 this probe also pinned "draw the payload's verdict, do not
// re-derive it" — the renderer is now the placer, so there is no payload verdict
// left to obey and that half of the property is gone by construction rather than
// by omission. What survives is still real and still worth pinning: a row index
// maps to a baseline on the sample's OWN band, an out-of-range row is no label at
// all rather than a label at a wrong y, and a band with nothing hidden prints no
// remainder while a nonzero one states the count. The original defect this test
// was written for — a pass that labelled every overflow sample from an index
// cycle and had no concept of a remainder — is caught by the second half.
func TestJavaScriptBehaviorDurationsLabelRowsAndRemainders(t *testing.T) {
	indexHtml := generateLiveSite(t)

	constantPreamble := ""
	for _, constantName := range []string{
		"DURATIONS_LABEL_ROW_COUNT",
		"DURATIONS_LABEL_ROW_HEIGHT",
		"DURATIONS_LANE_LABEL_ROW_Y",
		"DURATIONS_REVERSED_LABEL_ROW_Y",
		"DURATIONS_VIEW_WIDTH",
		"DURATIONS_MARGIN_RIGHT",
	} {
		constantPreamble += fmt.Sprintf("var %s = %v;\n", constantName, durationsRendererConstant(t, constantName))
	}

	javascriptProbe := constantPreamble +
		"var svg = null;\n" +
		"var drawnRemainders = [];\n" +
		"function makeDurationsSvgNode(svg, name, attributes, textContent) { drawnRemainders.push(textContent); }\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function durationsBandRowY(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function durationsLabelBaselineY(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function durationsRemainderBaselineY(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function composeDurationsRemainderText(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function drawDurationsRemainder(") + `
drawDurationsRemainder(0, durationsRemainderBaselineY(DURATIONS_LANE_LABEL_ROW_Y), "over 60 min");
drawDurationsRemainder(23, durationsRemainderBaselineY(DURATIONS_LANE_LABEL_ROW_Y), "over 60 min");
drawDurationsRemainder(2, durationsRemainderBaselineY(DURATIONS_REVERSED_LABEL_ROW_Y), "reversed");
process.stdout.write(JSON.stringify({
  remainderBaselines: [
    durationsRemainderBaselineY(DURATIONS_LANE_LABEL_ROW_Y),
    durationsRemainderBaselineY(DURATIONS_REVERSED_LABEL_ROW_Y)
  ],
  baselines: [
    durationsLabelBaselineY({ wallMinutes: 95 }, 0),
    durationsLabelBaselineY({ wallMinutes: 95 }, 1),
    durationsLabelBaselineY({ wallMinutes: 95 }, -1),
    durationsLabelBaselineY({ wallMinutes: -20 }, 0),
    durationsLabelBaselineY({ wallMinutes: -20 }, 1),
    durationsLabelBaselineY({ wallMinutes: -20 }, -1),
    durationsLabelBaselineY({ wallMinutes: 95 }, undefined)
  ],
  remainders: drawnRemainders
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "durations label verdict", javascriptProbe)
	var probeResult struct {
		Baselines          []*float64 `json:"baselines"`
		RemainderBaselines []float64  `json:"remainderBaselines"`
		Remainders         []string   `json:"remainders"`
	}
	if decodeError := json.Unmarshal(probeOutput, &probeResult); decodeError != nil {
		t.Fatalf("decode durations label behavior: %v (output %q)", decodeError, probeOutput)
	}

	laneRowY := durationsRendererConstant(t, "DURATIONS_LANE_LABEL_ROW_Y")
	reversedRowY := durationsRendererConstant(t, "DURATIONS_REVERSED_LABEL_ROW_Y")
	rowHeight := durationsRendererConstant(t, "DURATIONS_LABEL_ROW_HEIGHT")
	wantBaselines := []*float64{
		&laneRowY,
		floatPointer(laneRowY + rowHeight),
		nil,
		&reversedRowY,
		floatPointer(reversedRowY + rowHeight),
		nil,
		nil,
	}
	if len(probeResult.Baselines) != len(wantBaselines) {
		t.Fatalf("baseline count = %d, want %d", len(probeResult.Baselines), len(wantBaselines))
	}
	for baselineIndex := range wantBaselines {
		got := probeResult.Baselines[baselineIndex]
		want := wantBaselines[baselineIndex]
		if (got == nil) != (want == nil) || (got != nil && *got != *want) {
			t.Fatalf("baseline[%d] = %v, want %v (nil means the sample carries no direct label)",
				baselineIndex, formatOptionalFloat(got), formatOptionalFloat(want))
		}
	}

	// The remainder must land on the band's LAST row. On the first row it sits at
	// the marks' own height, and the dense render showed it overprinted by the
	// very blob it was describing — the defect reproduced inside its own fix.
	lastRowOffset := (durationsRendererConstant(t, "DURATIONS_LABEL_ROW_COUNT") - 1) *
		durationsRendererConstant(t, "DURATIONS_LABEL_ROW_HEIGHT")
	wantRemainderBaselines := []float64{laneRowY + lastRowOffset, reversedRowY + lastRowOffset}
	if len(probeResult.RemainderBaselines) != len(wantRemainderBaselines) {
		t.Fatalf("remainder baseline count = %d, want %d",
			len(probeResult.RemainderBaselines), len(wantRemainderBaselines))
	}
	for baselineIndex, wantBaseline := range wantRemainderBaselines {
		if probeResult.RemainderBaselines[baselineIndex] != wantBaseline {
			t.Fatalf("remainder baseline[%d] = %v, want %v — the sentence must clear the mark row",
				baselineIndex, probeResult.RemainderBaselines[baselineIndex], wantBaseline)
		}
	}

	wantRemainders := []string{"+23 more over 60 min", "+2 more reversed"}
	if strings.Join(probeResult.Remainders, "|") != strings.Join(wantRemainders, "|") {
		t.Fatalf("drawn remainders = %q, want %q — a zero remainder must draw nothing and a nonzero one must state its count",
			probeResult.Remainders, wantRemainders)
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
func TestJavaScriptBehaviorDurationsSlowestDayAnnotationClearsItsNeighbours(t *testing.T) {
	indexHtml := generateLiveSite(t)
	annotationCases := durationsSlowestDayAnnotationCaseList()

	probeCases, encodeError := json.Marshal(durationsSlowestDayAnnotationProbeCases(annotationCases))
	if encodeError != nil {
		t.Fatalf("encode annotation probe cases: %v", encodeError)
	}
	annotationSource := sliceBalancedBlockAfter(t, indexHtml, "function drawDurationsSlowestDayAnnotation(")
	assertDurationsAnnotationBaselineIgnoresItsInputs(t, annotationSource)

	javascriptProbe := fmt.Sprintf("var DURATIONS_MEDIAN_ANNOTATION_BASELINE_Y = %v;\n",
		durationsRendererConstant(t, "DURATIONS_MEDIAN_ANNOTATION_BASELINE_Y")) +
		"var drawnNodes = [];\n" +
		"function makeDurationsSvgNode(svg, name, attributes, textContent) { drawnNodes.push({ name: name, attributes: attributes, text: textContent }); }\n" +
		annotationSource + `
var probeCases = ` + string(probeCases) + `;
probeCases.forEach(function (probeCase) {
  drawDurationsSlowestDayAnnotation(null, { medianMinutes: probeCase.medianMinutes }, probeCase.dayCentreX);
});
process.stdout.write(JSON.stringify(drawnNodes.map(function (node) {
  return { y: Number(node.attributes.y), x: Number(node.attributes.x), anchor: node.attributes["text-anchor"], text: node.text };
})));`

	probeOutput := runJavaScriptBehaviorProbe(t, "durations slowest-day annotation", javascriptProbe)
	var drawnAnnotations []struct {
		Y      float64 `json:"y"`
		X      float64 `json:"x"`
		Anchor string  `json:"anchor"`
		Text   string  `json:"text"`
	}
	if decodeError := json.Unmarshal(probeOutput, &drawnAnnotations); decodeError != nil {
		t.Fatalf("decode slowest-day annotation behavior: %v (output starts %q)",
			decodeError, string(probeOutput[:min(len(probeOutput), 400)]))
	}
	if len(drawnAnnotations) != len(annotationCases) {
		t.Fatalf("the renderer drew %d annotations for %d cases — it must draw exactly one per slowest day",
			len(drawnAnnotations), len(annotationCases))
	}

	// Every box below is the taller of the two models of its own face: the one
	// the renderer declares and the one the browser draws.
	annotationAscent := math.Max(durationsRendererConstant(t, "DURATIONS_LABEL_TEXT_ASCENT"), durationsMeasuredMarkLabelAscentUnits)
	annotationDescent := math.Max(durationsLabelTextDescentUnits, durationsMeasuredMarkLabelDescentUnits)
	medianBaseline := durationsRendererConstant(t, "DURATIONS_MEDIAN_BOTTOM")

	// The annotation's neighbours in the strip it now occupies. Panel B's "0"
	// tick is the one a render caught and no arithmetic would have: it sits in
	// the y-axis gutter, so it is only ever in the annotation's way when the
	// slowest day is the leftmost — the same luck-of-x that hid the original
	// defect. The ticks for 15/30/45 need no case of their own: their baselines
	// are above DURATIONS_MEDIAN_BOTTOM, which the baseline check below covers.
	neighbourBoxes := []struct {
		neighbourName string
		baseline      float64
		ascent        float64
		descent       float64
	}{
		{
			"panel B's title",
			durationsRendererConstant(t, "DURATIONS_MEDIAN_TITLE_Y"),
			durationsMeasuredAxisTitleAscentUnits,
			durationsMeasuredAxisTitleDescentUnits,
		},
		{
			"panel B's \"0\" axis tick",
			durationsRendererConstant(t, "DURATIONS_MEDIAN_BOTTOM") + durationsRendererConstant(t, "DURATIONS_TICK_BASELINE_DROP"),
			annotationAscent,
			annotationDescent,
		},
		{
			"panel C's title",
			durationsRendererConstant(t, "DURATIONS_COUNT_TITLE_Y"),
			durationsMeasuredAxisTitleAscentUnits,
			durationsMeasuredAxisTitleDescentUnits,
		},
	}
	for _, neighbour := range neighbourBoxes {
		neighbourTop := neighbour.baseline - neighbour.ascent
		neighbourBottom := neighbour.baseline + neighbour.descent
		for caseIndex, probeCase := range annotationCases {
			drawn := drawnAnnotations[caseIndex]
			annotationTop := drawn.Y - annotationAscent
			annotationBottom := drawn.Y + annotationDescent
			if annotationBottom >= neighbourTop && annotationTop <= neighbourBottom {
				t.Fatalf("%s: the annotation's text box [%.2f, %.2f] intersects %s's box [%.2f, %.2f] — the two overprint wherever their x ranges meet, and x follows whichever day is slowest",
					probeCase.caseName, annotationTop, annotationBottom, neighbour.neighbourName, neighbourTop, neighbourBottom)
			}
		}
	}
	for caseIndex, probeCase := range annotationCases {
		drawn := drawnAnnotations[caseIndex]
		if drawn.Y-annotationAscent <= medianBaseline {
			t.Fatalf("%s: the annotation's text box starts at %.2f, above panel B's baseline at %.2f — inside the plot it overprints the bars, which are 4 units wide and shoulder to shoulder on a dense board",
				probeCase.caseName, drawn.Y-annotationAscent, medianBaseline)
		}
		if drawn.Y != drawnAnnotations[0].Y {
			t.Fatalf("%s: the annotation's baseline is %.2f but %s put it at %.2f — a baseline that moves with the day's position or its bar's height is a clearance that holds only for the days this repository happens to have",
				probeCase.caseName, drawn.Y, annotationCases[0].caseName, drawnAnnotations[0].Y)
		}
		if drawn.X != probeCase.dayCentreX || drawn.Anchor != "middle" {
			t.Fatalf("%s: the annotation was drawn at x=%.2f anchored %q, want x=%.2f anchored \"middle\" — it must stay centred on the day it describes",
				probeCase.caseName, drawn.X, drawn.Anchor, probeCase.dayCentreX)
		}
		if wantText := fmt.Sprintf("%.0f min", math.Round(probeCase.medianMinutes)); drawn.Text != wantText {
			t.Fatalf("%s: the annotation reads %q, want %q — moving it must not cost it the value it exists to state",
				probeCase.caseName, drawn.Text, wantText)
		}
	}

	// The strip's fourth occupant is the month rule, and unlike the other three
	// it cannot be cleared: .durations-month-line spans DURATIONS_MAIN_TOP to
	// DURATIONS_COUNT_BOTTOM, so it crosses EVERY baseline the annotation could
	// legally take, and it crosses panel A's reversed-band labels the same way.
	// The crossing is accepted, not overlooked — on a fixture whose slowest day
	// falls on a month boundary the rule passes between the "9" and the " min".
	// What makes it acceptable is that it is a one-unit soft rule, so that is
	// what gets asserted in place of a clearance: if the month rule ever grows
	// wide or firm, this fires and the acceptance has to be re-argued.
	annotationTop := drawnAnnotations[0].Y - annotationAscent
	annotationBottom := drawnAnnotations[0].Y + annotationDescent
	monthRuleTop := durationsRendererConstant(t, "DURATIONS_MAIN_TOP")
	monthRuleBottom := durationsRendererConstant(t, "DURATIONS_COUNT_BOTTOM")
	if annotationTop <= monthRuleTop || annotationBottom >= monthRuleBottom {
		t.Fatalf("the annotation's text box [%.2f, %.2f] is no longer inside the month rule's span [%.2f, %.2f] — the crossing this test ACCEPTS has become avoidable, so it belongs in the clearance list above instead",
			annotationTop, annotationBottom, monthRuleTop, monthRuleBottom)
	}
	if strokeWidth := durationsStyleDeclaration(t, ".durations-month-line", "stroke-width"); strokeWidth != "1" {
		t.Fatalf("the month rule's stroke-width is %q, not \"1\" — it is allowed to cross the slowest-day annotation only because it is a hairline",
			strokeWidth)
	}
	if strokeColour := durationsStyleDeclaration(t, ".durations-month-line", "stroke"); strokeColour != "var(--line-soft)" {
		t.Fatalf("the month rule is stroked %q, not the soft line token — it is allowed to cross the slowest-day annotation only because it is soft",
			strokeColour)
	}
}

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

// durationLabelWidthSampleMinutes spans the renderer's formatting branches: sub-hour
// with a decimal, negative, exactly on the hour, and multi-hour. The last two are
// the rounding-carry values: they are the only ones where a per-unit rounding
// regression changes the character count ("1h 60m" against "2h 0m"), so they keep
// this width lock-step sensitive to it.
var durationLabelWidthSampleMinutes = []float64{7.5, -25, 60, 95.4, 655.2, 1440, 119.5, 59.96}

// Placement sizes a label from the text the renderer will draw, so it carries its
// own width model of that text. The renderer stays the definition of the copy —
// this pins the model to it. Without this the two agree today and drift the first
// time the renderer's formatting gains a character, and the only symptom would be
// labels overlapping again at exactly the densities this REQ was about.
func TestJavaScriptBehaviorDurationsLabelWidthModelMatchesTheRendererFormatter(t *testing.T) {
	indexHtml := generateLiveSite(t)
	probeValues, encodeError := json.Marshal(durationLabelWidthSampleMinutes)
	if encodeError != nil {
		t.Fatalf("encode probe values: %v", encodeError)
	}
	javascriptProbe := sliceBalancedBlockAfter(t, indexHtml, "function formatDurationMinutes(") + `
process.stdout.write(JSON.stringify(` + string(probeValues) + `.map(formatDurationMinutes)));`

	probeOutput := runJavaScriptBehaviorProbe(t, "durations label width model", javascriptProbe)
	var rendererTexts []string
	if decodeError := json.Unmarshal(probeOutput, &rendererTexts); decodeError != nil {
		t.Fatalf("decode renderer formatting: %v (output %q)", decodeError, probeOutput)
	}
	if len(rendererTexts) != len(durationLabelWidthSampleMinutes) {
		t.Fatalf("renderer produced %d strings, want %d", len(rendererTexts), len(durationLabelWidthSampleMinutes))
	}
	for valueIndex, minutes := range durationLabelWidthSampleMinutes {
		// The renderer writes U+2212 for a negative sign; only the character
		// COUNT reaches the width model, so compare lengths in runes.
		rendererLength := len([]rune(rendererTexts[valueIndex]))
		modelLength := len([]rune(formatDurationLabelMinutes(minutes)))
		if rendererLength != modelLength {
			t.Fatalf("%.1f min: renderer draws %q (%d chars) but the width model assumes %q (%d chars)",
				minutes, rendererTexts[valueIndex], rendererLength,
				formatDurationLabelMinutes(minutes), modelLength)
		}
	}
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

func TestJavaScriptBehaviorDurationsWindowsProjectOneSharedRealTimeDomain(t *testing.T) {
	rendererFragment, readError := embeddedWebAssets.ReadFile("web/board-durations.js")
	if readError != nil {
		t.Fatalf("read web/board-durations.js: %v", readError)
	}

	completionInstants := []time.Time{
		time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC),
		time.Date(2026, 5, 27, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC),
		time.Date(2026, 7, 26, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC),
		time.Date(2026, 8, 11, 12, 0, 0, 0, time.UTC),
		time.Date(2026, 8, 12, 12, 0, 0, 0, time.UTC),
		time.Date(2026, 8, 24, 12, 0, 0, 0, time.UTC),
	}
	fixtureTickets := make([]*RequestTicket, 0, len(completionInstants))
	for completionIndex, completedAt := range completionInstants {
		fixtureTickets = append(fixtureTickets, durationTicket(
			fmt.Sprintf("REQ-%03d", completionIndex+1),
			"B",
			completedAt.Add(-10*time.Minute).Format(time.RFC3339),
			completedAt.Format(time.RFC3339),
		))
	}
	fixtureBoard := &Board{
		GeneratedAt: time.Date(2030, 1, 1, 12, 0, 0, 0, time.UTC),
		AllRequests: fixtureTickets,
	}
	generatedData, buildError := buildGeneratedBoardData(fixtureBoard)
	if buildError != nil {
		t.Fatalf("buildGeneratedBoardData: %v", buildError)
	}
	boardDataJSON, encodeError := json.Marshal(generatedData)
	if encodeError != nil {
		t.Fatalf("encode board payload: %v", encodeError)
	}

	probeDriver := `
function renderedWindow(windowName) {
  durationsStubHosts["durations-chart"] = makeStubNode("div");
  durationsStubHosts["durations-summary"] = makeStubNode("p");
  durationsStubHosts["durations-readout"] = makeStubNode("p");
  durationsStubHosts["durations-table-body"] = makeStubNode("tbody");
  setDurationsWindow(windowName);
  renderDurationsView();
  var markCentres = [], panelBBarCentres = [], panelCBars = 0;
  function walkDrawnNodes(parentNode) {
    (parentNode.children || []).forEach(function (childNode) {
      var attributes = childNode.attributes || {};
      var nodeClass = String(attributes["class"] || "");
      if (childNode.stubName === "circle" && nodeClass.indexOf("durations-mark") !== -1) {
        markCentres.push(Number(attributes.cx));
      }
      if (childNode.stubName === "rect" && nodeClass === "durations-bar") {
        panelBBarCentres.push(Number(attributes.x) + Number(attributes.width) / 2);
      }
      if (childNode.stubName === "rect" && nodeClass.indexOf("durations-bar-count") !== -1) {
        panelCBars += 1;
      }
      walkDrawnNodes(childNode);
    });
  }
  var svg = durationsStubHosts["durations-chart"].children[0];
  walkDrawnNodes(svg);
  return {
    summary: durationsStubHosts["durations-summary"].textContent,
    ariaLabel: svg.attributes["aria-label"],
    markCentres: markCentres,
    panelBBarCentres: panelBBarCentres,
    panelCBars: panelCBars,
    tableRows: durationsStubHosts["durations-table-body"].children.length
  };
}
process.stdout.write(JSON.stringify({
  last30: renderedWindow("30"),
  last90: renderedWindow("90"),
  all: renderedWindow("all")
}));
`
	javascriptProbe := durationsRenderDomStubPreamble +
		"var boardData = " + string(boardDataJSON) + ";\n" +
		string(rendererFragment) +
		probeDriver
	probeOutput := runJavaScriptBehaviorProbe(t, "Durations projected windows", javascriptProbe)

	type renderedWindow struct {
		Summary          string    `json:"summary"`
		AriaLabel        string    `json:"ariaLabel"`
		MarkCentres      []float64 `json:"markCentres"`
		PanelBBarCentres []float64 `json:"panelBBarCentres"`
		PanelCBars       int       `json:"panelCBars"`
		TableRows        int       `json:"tableRows"`
	}
	var result struct {
		Last30 renderedWindow `json:"last30"`
		Last90 renderedWindow `json:"last90"`
		All    renderedWindow `json:"all"`
	}
	if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
		t.Fatalf("decode projected Durations windows: %v (output starts %q)",
			decodeError, string(probeOutput[:min(len(probeOutput), 400)]))
	}

	for _, windowCase := range []struct {
		name          string
		window        renderedWindow
		wantCount     int
		wantStartDate string
		wantEndDate   string
	}{
		{name: "Last 30 days", window: result.Last30, wantCount: 5, wantStartDate: "26 Jul", wantEndDate: "25 Aug"},
		{name: "Last 90 days", window: result.Last90, wantCount: 7, wantStartDate: "27 May", wantEndDate: "25 Aug"},
		{name: "All history", window: result.All, wantCount: 8, wantStartDate: "1 Apr", wantEndDate: "25 Aug"},
	} {
		if len(windowCase.window.MarkCentres) != windowCase.wantCount ||
			windowCase.window.TableRows != windowCase.wantCount ||
			len(windowCase.window.PanelBBarCentres) != windowCase.wantCount ||
			windowCase.window.PanelCBars != windowCase.wantCount {
			t.Errorf("%s counts = marks %d, table %d, Panel B %d, Panel C %d; want %d projected samples/days on every surface",
				windowCase.name, len(windowCase.window.MarkCentres), windowCase.window.TableRows,
				len(windowCase.window.PanelBBarCentres), windowCase.window.PanelCBars, windowCase.wantCount)
		}
		for _, surfaceCopy := range []string{windowCase.window.Summary, windowCase.window.AriaLabel} {
			for _, requiredText := range []string{windowCase.name, windowCase.wantStartDate, windowCase.wantEndDate, "end exclusive", fmt.Sprintf("%d archived REQ", windowCase.wantCount)} {
				if !strings.Contains(surfaceCopy, requiredText) {
					t.Errorf("%s accessibility copy %q is missing %q", windowCase.name, surfaceCopy, requiredText)
				}
			}
		}
	}

	marginLeft := durationsRendererConstant(t, "DURATIONS_MARGIN_LEFT")
	plotWidth := durationsRendererConstant(t, "DURATIONS_VIEW_WIDTH") - marginLeft - durationsRendererConstant(t, "DURATIONS_MARGIN_RIGHT")
	firstDayRight := marginLeft + plotWidth/30
	if result.Last30.MarkCentres[0] < marginLeft || result.Last30.MarkCentres[0] > firstDayRight {
		t.Errorf("Last 30 days left-boundary sample x=%.2f, want inside first day slot [%.2f, %.2f]", result.Last30.MarkCentres[0], marginLeft, firstDayRight)
	}
	firstDayGap := result.Last30.PanelBBarCentres[2] - result.Last30.PanelBBarCentres[1]
	secondDayGap := result.Last30.PanelBBarCentres[3] - result.Last30.PanelBBarCentres[2]
	if math.Abs(firstDayGap-secondDayGap) > 0.11 || firstDayGap <= 0 {
		t.Errorf("equal UTC-day gaps draw %.2f and %.2f units apart, want equal positive affine spacing", firstDayGap, secondDayGap)
	}
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
    if (childNode.stubName === "circle") {
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
func TestJavaScriptBehaviorDurationsDayBucketsStayInsideThePlot(t *testing.T) {
	rendererFragment, readError := embeddedWebAssets.ReadFile("web/board-durations.js")
	if readError != nil {
		t.Fatalf("read web/board-durations.js: %v", readError)
	}

	marginLeft := durationsRendererConstant(t, "DURATIONS_MARGIN_LEFT")
	plotRight := durationsRendererConstant(t, "DURATIONS_VIEW_WIDTH") -
		durationsRendererConstant(t, "DURATIONS_MARGIN_RIGHT")

	for _, dayCount := range []int{1, 2, 14, 400} {
		dayCount := dayCount
		t.Run(fmt.Sprintf("%d-active-days", dayCount), func(t *testing.T) {
			fixtureTickets := durationsDayCountFixtureTickets(dayCount)
			aggregate := buildDurationAggregate(fixtureTickets)
			fixtureBoard := &Board{
				GeneratedAt: time.Date(2026, 8, 18, 12, 0, 0, 0, time.UTC),
				AllRequests: fixtureTickets,
			}
			generatedData, buildError := buildGeneratedBoardData(fixtureBoard)
			if buildError != nil {
				t.Fatalf("buildGeneratedBoardData: %v", buildError)
			}
			durationsJson, encodeError := json.Marshal(generatedData.Durations)
			if encodeError != nil {
				t.Fatalf("encode durations payload: %v", encodeError)
			}

			javascriptProbe := durationsRenderDomStubPreamble +
				"var boardData = { durations: " + string(durationsJson) + " };\n" +
				string(rendererFragment) +
				"setDurationsWindow(\"all\");\n" +
				durationsRenderProbeDriver
			probeOutput := runJavaScriptBehaviorProbe(t,
				fmt.Sprintf("durations day buckets (%d active days)", dayCount), javascriptProbe)

			var drawn struct {
				Bars []struct {
					Class string  `json:"class"`
					X     float64 `json:"x"`
					Width float64 `json:"width"`
				} `json:"bars"`
				AnnotationXs []float64 `json:"annotationXs"`
				MarkCxs      []float64 `json:"markCxs"`
				FirstMarkCxs []float64 `json:"firstMarkCxs"`
			}
			if decodeError := json.Unmarshal(probeOutput, &drawn); decodeError != nil {
				t.Fatalf("decode drawn durations geometry: %v (output starts %q)",
					decodeError, string(probeOutput[:min(len(probeOutput), 400)]))
			}

			// (1) Every Panel B and C bar inside the plot area. 0.05 covers the
			// renderer's toFixed(1) rounding and nothing else.
			panelBBarCount := 0
			for _, bar := range drawn.Bars {
				if !strings.Contains(bar.Class, "durations-bar-count") &&
					!strings.Contains(bar.Class, "durations-bar-over-ceiling") {
					panelBBarCount++
				}
				if bar.X < marginLeft-0.05 || bar.X+bar.Width > plotRight+0.05 {
					t.Errorf("%q bar spans x %.1f–%.1f, outside the plot area [%.0f, %.0f]",
						bar.Class, bar.X, bar.X+bar.Width, marginLeft, plotRight)
				}
			}

			// (2) One Panel B bar per day with a median — the guard against a
			// render that went off the rails and drew nothing to check.
			medianDayCount := 0
			for _, day := range aggregate.Days {
				if day.HasMedian {
					medianDayCount++
				}
			}
			if panelBBarCount != medianDayCount {
				t.Errorf("drew %d Panel B bars for %d days with a median", panelBBarCount, medianDayCount)
			}

			// (3) Exactly one slowest-day annotation, anchored inside the plot —
			// it exists to state a value a clipped bar cannot, and cannot do
			// that from off-screen.
			if len(drawn.AnnotationXs) != 1 {
				t.Fatalf("drew %d slowest-day annotations, want exactly 1", len(drawn.AnnotationXs))
			}
			if annotationX := drawn.AnnotationXs[0]; annotationX < marginLeft-0.05 || annotationX > plotRight+0.05 {
				t.Errorf("slowest-day annotation anchored at x=%.1f, outside the plot area [%.0f, %.0f]",
					annotationX, marginLeft, plotRight)
			}

			// (4) Every Panel A mark stays inside its own UTC day slot. This pins
			// the shared whole-day axis domain while allowing the intentional
			// deterministic jitter within that slot.
			rangeStart, rangeEnd, hasRange := durationLabelTimeRange(aggregate.Samples)
			if !hasRange {
				t.Fatal("fixture produced no label time range")
			}
			if len(drawn.MarkCxs) != len(aggregate.Samples) {
				t.Fatalf("drew %d marks for %d samples", len(drawn.MarkCxs), len(aggregate.Samples))
			}
			for sampleIndex, sample := range aggregate.Samples {
				dayStart := sample.CompletionTime.UTC().Truncate(24 * time.Hour)
				dayLeft := marginLeft + durationLabelPlotX(dayStart, rangeStart, rangeEnd)
				dayRight := marginLeft + durationLabelPlotX(dayStart.Add(24*time.Hour), rangeStart, rangeEnd)
				if drawn.MarkCxs[sampleIndex] < dayLeft-0.06 || drawn.MarkCxs[sampleIndex] > dayRight+0.06 {
					t.Errorf("%s mark drawn at x=%.2f outside its UTC-day slot [%.2f, %.2f]",
						sample.RequestId, drawn.MarkCxs[sampleIndex], dayLeft, dayRight)
				}
			}

			// (5) The same payload renders to the same jittered coordinates, and
			// a busy day uses more than one x instead of restacking every mark.
			if len(drawn.FirstMarkCxs) != len(drawn.MarkCxs) {
				t.Fatalf("first render drew %d marks and second drew %d", len(drawn.FirstMarkCxs), len(drawn.MarkCxs))
			}
			busyDaySpread := 0.0
			for sampleIndex := range drawn.MarkCxs {
				if math.Abs(drawn.FirstMarkCxs[sampleIndex]-drawn.MarkCxs[sampleIndex]) > 0.001 {
					t.Errorf("mark %d moved from x=%.2f to x=%.2f across identical renders",
						sampleIndex, drawn.FirstMarkCxs[sampleIndex], drawn.MarkCxs[sampleIndex])
				}
				if sampleIndex%2 == 1 {
					busyDaySpread = math.Max(busyDaySpread, math.Abs(drawn.MarkCxs[sampleIndex]-drawn.MarkCxs[sampleIndex-1]))
				}
			}
			if busyDaySpread < 0.1 {
				t.Errorf("busy-day marks have only %.2f units of x spread; jitter is degenerate", busyDaySpread)
			}
		})
	}
}

// REQ-349: Panel A needs enough vertical resolution for ordinary work, enough
// horizontal resolution for a busy day, and a quiet daily distribution behind
// the individual REQs. This drives the complete renderer so the assertions pin
// SVG output and draw order rather than helpers that could be left unwired.
func TestJavaScriptBehaviorDurationsPanelASpreadAndDailyDistribution(t *testing.T) {
	rendererFragment, readError := embeddedWebAssets.ReadFile("web/board-durations.js")
	if readError != nil {
		t.Fatalf("read web/board-durations.js: %v", readError)
	}

	dayStart := time.Date(2026, 8, 10, 9, 0, 0, 0, time.UTC)
	fixtureMinutes := []float64{0, 5, 15, 30, 45, 60, 300, -20}
	fixtureTickets := make([]*RequestTicket, 0, len(fixtureMinutes)+2)
	for sampleIndex, minutes := range fixtureMinutes {
		completedAt := dayStart.Add(time.Duration(sampleIndex) * time.Hour)
		claimedAt := completedAt.Add(-time.Duration(minutes * float64(time.Minute)))
		fixtureTickets = append(fixtureTickets, durationTicket(
			fmt.Sprintf("REQ-%03d", sampleIndex+1), "B",
			claimedAt.Format(time.RFC3339), completedAt.Format(time.RFC3339),
		))
	}
	// A missing-route sample proves that the lower ordinary opacity does not
	// erase the outlined unknown category.
	unknownCompletedAt := dayStart.Add(24*time.Hour + 2*time.Hour)
	fixtureTickets = append(fixtureTickets, durationTicket(
		"REQ-901", "", unknownCompletedAt.Add(-10*time.Minute).Format(time.RFC3339), unknownCompletedAt.Format(time.RFC3339),
	))
	thirdDayCompletedAt := dayStart.Add(48*time.Hour + 3*time.Hour)
	fixtureTickets = append(fixtureTickets, durationTicket(
		"REQ-902", "A", thirdDayCompletedAt.Add(-20*time.Minute).Format(time.RFC3339), thirdDayCompletedAt.Format(time.RFC3339),
	))

	aggregate := buildDurationAggregate(fixtureTickets)
	generatedData, buildError := buildGeneratedBoardData(&Board{AllRequests: fixtureTickets})
	if buildError != nil {
		t.Fatalf("build generated board data: %v", buildError)
	}
	durationsJSON, encodeError := json.Marshal(generatedData.Durations)
	if encodeError != nil {
		t.Fatalf("encode durations payload: %v", encodeError)
	}

	probeDriver := `
renderDurationsView();
var svg = durationsStubHosts["durations-chart"].children[0];
var ticks = [], circles = [], paths = [];
(svg.children || []).forEach(function (childNode, childIndex) {
  var attributes = childNode.attributes || {};
  if (childNode.stubName === "text" && String(attributes["class"] || "").indexOf("durations-tick") !== -1) {
    ticks.push({ text: ((childNode.children || [])[0] || {}).textContent || "", x: Number(attributes.x), y: Number(attributes.y) });
  }
  if (childNode.stubName === "circle") {
    circles.push({
      cx: Number(attributes.cx), cy: Number(attributes.cy), opacity: attributes.opacity || "",
      fill: attributes.fill || "", class: attributes["class"] || "", childIndex: childIndex
    });
  }
  if (childNode.stubName === "path") {
    paths.push({ class: attributes["class"] || "", d: attributes.d || "", childIndex: childIndex });
  }
});
process.stdout.write(JSON.stringify({ ticks: ticks, circles: circles, paths: paths }));
`
	probeOutput := runJavaScriptBehaviorProbe(t, "durations panel A spread and distribution",
		durationsRenderDomStubPreamble+
			"var boardData = { durations: "+string(durationsJSON)+" };\n"+
			string(rendererFragment)+probeDriver)

	var drawn struct {
		Ticks []struct {
			Text string  `json:"text"`
			X    float64 `json:"x"`
			Y    float64 `json:"y"`
		} `json:"ticks"`
		Circles []struct {
			CX         float64 `json:"cx"`
			CY         float64 `json:"cy"`
			Opacity    string  `json:"opacity"`
			Fill       string  `json:"fill"`
			Class      string  `json:"class"`
			ChildIndex int     `json:"childIndex"`
		} `json:"circles"`
		Paths []struct {
			Class      string `json:"class"`
			D          string `json:"d"`
			ChildIndex int    `json:"childIndex"`
		} `json:"paths"`
	}
	if decodeError := json.Unmarshal(probeOutput, &drawn); decodeError != nil {
		t.Fatalf("decode panel A geometry: %v (output starts %q)", decodeError, string(probeOutput[:min(len(probeOutput), 400)]))
	}
	if len(drawn.Circles) != len(aggregate.Samples) {
		t.Fatalf("drew %d marks for %d samples", len(drawn.Circles), len(aggregate.Samples))
	}

	// The tick set is exact, including the new 5-minute foothold. Its y positions
	// follow sqrt(minutes / 60), not the old linear scale.
	panelATicks := map[string]float64{}
	for _, tick := range drawn.Ticks {
		if math.Abs(tick.X-(durationsRendererConstant(t, "DURATIONS_MARGIN_LEFT")-8)) < 0.01 &&
			tick.Y <= durationsRendererConstant(t, "DURATIONS_MAIN_BOTTOM")+durationsRendererConstant(t, "DURATIONS_TICK_BASELINE_DROP")+0.01 {
			panelATicks[tick.Text] = tick.Y - durationsRendererConstant(t, "DURATIONS_TICK_BASELINE_DROP")
		}
	}
	wantTickMinutes := []float64{0, 5, 15, 30, 45, 60}
	if len(panelATicks) != len(wantTickMinutes)+1 { // plus the separate 60+ overflow-lane tick
		t.Fatalf("Panel A ticks = %v, want 60+ and exactly 0, 5, 15, 30, 45, 60", panelATicks)
	}
	mainTop := durationsRendererConstant(t, "DURATIONS_MAIN_TOP")
	mainBottom := durationsRendererConstant(t, "DURATIONS_MAIN_BOTTOM")
	for _, minutes := range wantTickMinutes {
		label := strconv.FormatFloat(minutes, 'f', -1, 64)
		gotY, exists := panelATicks[label]
		if !exists {
			t.Errorf("Panel A has no %s-minute tick", label)
			continue
		}
		wantY := mainBottom - math.Sqrt(minutes/60)*(mainBottom-mainTop)
		if math.Abs(gotY-wantY) > 0.01 {
			t.Errorf("%s-minute tick y=%.3f, want sqrt-scale y=%.3f", label, gotY, wantY)
		}
	}

	ordinaryMark := drawn.Circles[1]
	if ordinaryMark.Opacity == "" || ordinaryMark.Opacity == "1" {
		t.Errorf("ordinary mark opacity = %q, want a stated lower opacity", ordinaryMark.Opacity)
	}
	reversedMark := drawn.Circles[7]
	if reversedMark.Fill != "var(--durations-critical)" || reversedMark.Opacity != "" {
		t.Errorf("reversed mark fill/opacity = %q / %q, want undimmed critical red", reversedMark.Fill, reversedMark.Opacity)
	}
	unknownMark := drawn.Circles[8]
	if !strings.Contains(unknownMark.Class, "durations-mark-unknown") || unknownMark.Opacity != "" {
		t.Errorf("unknown mark class/opacity = %q / %q, want an undimmed outlined unknown", unknownMark.Class, unknownMark.Opacity)
	}

	if len(drawn.Paths) != 2 {
		t.Fatalf("drew %d Panel A distribution paths, want one p25-p75 ribbon and one median line: %+v", len(drawn.Paths), drawn.Paths)
	}
	pathByClass := map[string]struct {
		D          string
		ChildIndex int
	}{}
	for _, path := range drawn.Paths {
		pathByClass[path.Class] = struct {
			D          string
			ChildIndex int
		}{path.D, path.ChildIndex}
		if path.D == "" || strings.Contains(path.D, "NaN") || strings.Contains(path.D, "Infinity") {
			t.Errorf("%q path has invalid geometry %q", path.Class, path.D)
		}
	}
	ribbon, ribbonExists := pathByClass["durations-quantile-ribbon"]
	median, medianExists := pathByClass["durations-quantile-median"]
	if !ribbonExists || !medianExists {
		t.Fatalf("distribution path classes = %v, want ribbon and median", pathByClass)
	}
	if ribbon.ChildIndex >= drawn.Circles[0].ChildIndex || median.ChildIndex >= drawn.Circles[0].ChildIndex {
		t.Errorf("distribution paths at child indexes %d/%d are not behind first circle at %d",
			ribbon.ChildIndex, median.ChildIndex, drawn.Circles[0].ChildIndex)
	}
}

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
func TestJavaScriptBehaviorDurationsUserRequestLaneNamesEveryUserRequest(t *testing.T) {
	rendererFragment, readError := embeddedWebAssets.ReadFile("web/board-durations.js")
	if readError != nil {
		t.Fatalf("read web/board-durations.js: %v", readError)
	}

	const fixtureUserRequestCount = 40
	const fixtureNoUserRequestCount = 5
	fixtureTickets := durationsUserRequestLaneFixtureTickets(fixtureUserRequestCount, fixtureNoUserRequestCount)
	fixtureBoard := &Board{
		GeneratedAt: time.Date(2026, 8, 18, 12, 0, 0, 0, time.UTC),
		AllRequests: fixtureTickets,
	}
	generatedData, buildError := buildGeneratedBoardData(fixtureBoard)
	if buildError != nil {
		t.Fatalf("buildGeneratedBoardData: %v", buildError)
	}
	// The WHOLE payload, not just the durations slice: the join reads
	// boardData.requests, so a probe handed only the samples would take the
	// unknown-UR path for every row and never exercise the join at all.
	boardDataJson, encodeError := json.Marshal(generatedData)
	if encodeError != nil {
		t.Fatalf("encode board payload: %v", encodeError)
	}

	probeDriver := `
renderDurationsView();
var drawnBrackets = [], drawnTexts = [];
function walkDrawnNodes(parentNode) {
  (parentNode.children || []).forEach(function (childNode) {
    var attributes = childNode.attributes || {};
    var nodeClass = String(attributes["class"] || "");
    if (childNode.stubName === "rect" && nodeClass.indexOf("durations-ur-bracket") !== -1) {
      drawnBrackets.push({ class: nodeClass });
    }
    if (childNode.stubName === "text") {
      drawnTexts.push((childNode.children[0] || {}).textContent || "");
    }
    walkDrawnNodes(childNode);
  });
}
walkDrawnNodes(durationsStubHosts["durations-chart"]);
var userRequestCells = durationsStubHosts["durations-table-body"].children.map(function (tableRow) {
  return tableRow.children[1].textContent;
});
process.stdout.write(JSON.stringify({
  bracketClasses: drawnBrackets.map(function (bracket) { return bracket.class; }),
  drawnTexts: drawnTexts,
  unknownName: DURATIONS_UNKNOWN_USER_REQUEST_NAME,
  remainderSentenceForOne: composeDurationsUserRequestRemainderText(1),
  userRequestCells: userRequestCells
}));
`

	javascriptProbe := durationsRenderDomStubPreamble +
		"var boardData = " + string(boardDataJson) + ";\n" +
		string(rendererFragment) +
		probeDriver
	probeOutput := runJavaScriptBehaviorProbe(t, "durations UR lane", javascriptProbe)

	var drawn struct {
		BracketClasses          []string `json:"bracketClasses"`
		DrawnTexts              []string `json:"drawnTexts"`
		UnknownName             string   `json:"unknownName"`
		RemainderSentenceForOne string   `json:"remainderSentenceForOne"`
		UserRequestCells        []string `json:"userRequestCells"`
	}
	if decodeError := json.Unmarshal(probeOutput, &drawn); decodeError != nil {
		t.Fatalf("decode drawn UR lane: %v (output starts %q)",
			decodeError, string(probeOutput[:min(len(probeOutput), 400)]))
	}

	// ---- the fixture does what it claims, checked before anything is read from it.
	userRequestBrackets, unknownBrackets := 0, 0
	for _, bracketClass := range drawn.BracketClasses {
		if strings.Contains(bracketClass, "durations-ur-bracket-unknown") {
			unknownBrackets++
			continue
		}
		userRequestBrackets++
	}
	if userRequestBrackets == 0 {
		t.Fatal("the lane drew no UR brackets at all, so nothing below is a test of the lane")
	}
	if userRequestBrackets >= fixtureUserRequestCount {
		t.Fatalf("all %d fixture URs found a row, so the remainder path was never taken — this fixture no "+
			"longer overflows the lane and the rule it exists to pin is untested", fixtureUserRequestCount)
	}

	// ---- rule one: drawn brackets plus the stated remainder account for every UR.
	statedRemainder := -1
	remainderPattern := regexp.MustCompile(`^\+([0-9]+) URs? with no free row$`)
	for _, drawnText := range drawn.DrawnTexts {
		if remainderPattern.MatchString(drawnText) {
			if _, scanError := fmt.Sscanf(drawnText, "+%d", &statedRemainder); scanError != nil {
				t.Fatalf("read the remainder count out of %q: %v", drawnText, scanError)
			}
		}
	}
	if statedRemainder < 0 {
		t.Fatalf("%d of %d fixture URs found no row and the lane said nothing about it — a reader takes the "+
			"brackets they can see for all of them",
			fixtureUserRequestCount-userRequestBrackets, fixtureUserRequestCount)
	}
	if userRequestBrackets+statedRemainder != fixtureUserRequestCount {
		t.Errorf("the lane drew %d brackets and stated %d more, accounting for %d of the fixture's %d URs — "+
			"every UR is either on a row or in the remainder, and the two must add up at any row count",
			userRequestBrackets, statedRemainder, userRequestBrackets+statedRemainder, fixtureUserRequestCount)
	}

	// ---- rule two: the samples with no UR are named, on every surface, and named apart.
	if unknownBrackets != 1 {
		t.Errorf("the lane drew %d unknown-UR brackets for %d samples carrying no user_request, want exactly 1 — "+
			"the bucket holds one reserved row, it does not compete for one",
			unknownBrackets, fixtureNoUserRequestCount)
	}
	if drawn.UnknownName == "" {
		t.Fatal("the unknown-UR name is empty, so every surface below states nothing")
	}
	if strings.Contains(drawn.RemainderSentenceForOne, drawn.UnknownName) {
		t.Errorf("the remainder sentence %q contains the unknown-UR name %q — a UR that found no row and a "+
			"sample that has no UR at all are different facts and must not read as one",
			drawn.RemainderSentenceForOne, drawn.UnknownName)
	}
	if len(drawn.UserRequestCells) != len(fixtureTickets) {
		t.Fatalf("the table has %d UR cells for %d samples", len(drawn.UserRequestCells), len(fixtureTickets))
	}
	namedUnknownCells := 0
	for cellIndex, cellText := range drawn.UserRequestCells {
		if cellText == "" {
			t.Fatalf("table row %d has a blank UR cell — a sample with no UR must SAY so, and a blank cell "+
				"reads as a rendering fault rather than as a fact about the REQ", cellIndex)
		}
		if cellText == drawn.UnknownName {
			namedUnknownCells++
		}
	}
	if namedUnknownCells != fixtureNoUserRequestCount {
		t.Errorf("%d table cells carry the unknown-UR name for %d samples with no user_request — a sample "+
			"with no UR must never be given one, and one with a UR must never lose it",
			namedUnknownCells, fixtureNoUserRequestCount)
	}
}

// REQ-347 changes Panel A's fill only. This drives the complete renderer across
// each selectable channel, because a helper-only test could keep the selector
// green while the circles still read sample.route. The fixture deliberately has
// more URs than the categorical palette, absent UR/domain values, and a reversed
// stamp: the overflow bucket and unknown state must be named while the critical
// stamp treatment remains stronger than any chosen colour channel.
func TestJavaScriptBehaviorDurationsColourChannelsNameAndRecolourPanelA(t *testing.T) {
	rendererFragment, readError := embeddedWebAssets.ReadFile("web/board-durations.js")
	if readError != nil {
		t.Fatalf("read web/board-durations.js: %v", readError)
	}

	fixtureTickets := make([]*RequestTicket, 0, 16)
	fixtureStart := time.Date(2026, 8, 10, 9, 0, 0, 0, time.UTC)
	for ticketIndex := 0; ticketIndex < 14; ticketIndex++ {
		completedAt := fixtureStart.Add(time.Duration(ticketIndex) * time.Hour)
		ticket := durationTicket(
			fmt.Sprintf("REQ-%03d", ticketIndex+1),
			"B",
			completedAt.Add(-12*time.Minute).Format(time.RFC3339),
			completedAt.Format(time.RFC3339),
		)
		ticket.UserRequestId = fmt.Sprintf("UR-%03d", ticketIndex+1)
		ticket.Domain = []string{"frontend", "backend", "testing"}[ticketIndex%3]
		fixtureTickets = append(fixtureTickets, ticket)
	}
	missingTicket := durationTicket("REQ-901", "", fixtureStart.Add(15*time.Hour).Format(time.RFC3339), fixtureStart.Add(15*time.Hour+12*time.Minute).Format(time.RFC3339))
	fixtureTickets = append(fixtureTickets, missingTicket)
	reversedTicket := durationTicket("REQ-902", "B", fixtureStart.Add(17*time.Hour).Format(time.RFC3339), fixtureStart.Add(16*time.Hour).Format(time.RFC3339))
	reversedTicket.UserRequestId = "UR-015"
	reversedTicket.Domain = "frontend"
	fixtureTickets = append(fixtureTickets, reversedTicket)

	fixtureBoard := &Board{
		GeneratedAt: time.Date(2026, 8, 18, 12, 0, 0, 0, time.UTC),
		AllRequests: fixtureTickets,
	}
	generatedData, buildError := buildGeneratedBoardData(fixtureBoard)
	if buildError != nil {
		t.Fatalf("buildGeneratedBoardData: %v", buildError)
	}
	boardDataJSON, encodeError := json.Marshal(generatedData)
	if encodeError != nil {
		t.Fatalf("encode board payload: %v", encodeError)
	}

	probeDriver := `
function renderedChannel(colourChannel) {
  durationsStubHosts["durations-chart"] = makeStubNode("div");
  durationsStubHosts["durations-colour-legend"] = makeStubNode("p");
  setDurationsColourChannel(colourChannel);
  renderDurationsView();
  var marks = [];
  function walkDrawnNodes(parentNode) {
    (parentNode.children || []).forEach(function (childNode) {
      var attributes = childNode.attributes || {};
      if (childNode.stubName === "circle" && String(attributes["class"] || "").indexOf("durations-mark") !== -1) {
        marks.push({ fill: attributes.fill || "", class: attributes["class"] || "" });
      }
      walkDrawnNodes(childNode);
    });
  }
  walkDrawnNodes(durationsStubHosts["durations-chart"]);
  return {
    marks: marks,
    legend: durationsStubHosts["durations-colour-legend"].textContent,
    ariaLabel: (durationsStubHosts["durations-chart"].children[0] || { attributes: {} }).attributes["aria-label"] || ""
  };
}
process.stdout.write(JSON.stringify({
  route: renderedChannel("route"),
  userRequest: renderedChannel("user-request"),
  domain: renderedChannel("domain")
}));
`

	javascriptProbe := durationsRenderDomStubPreamble +
		"var boardData = " + string(boardDataJSON) + ";\n" +
		string(rendererFragment) +
		probeDriver
	probeOutput := runJavaScriptBehaviorProbe(t, "durations colour channels", javascriptProbe)

	type drawnMark struct {
		Fill  string `json:"fill"`
		Class string `json:"class"`
	}
	type renderedChannel struct {
		Marks     []drawnMark `json:"marks"`
		Legend    string      `json:"legend"`
		AriaLabel string      `json:"ariaLabel"`
	}
	var rendered struct {
		Route       renderedChannel `json:"route"`
		UserRequest renderedChannel `json:"userRequest"`
		Domain      renderedChannel `json:"domain"`
	}
	if decodeError := json.Unmarshal(probeOutput, &rendered); decodeError != nil {
		t.Fatalf("decode rendered colour channels: %v (output starts %q)", decodeError, string(probeOutput[:min(len(probeOutput), 400)]))
	}

	if len(rendered.Route.Marks) != len(fixtureTickets) {
		t.Fatalf("route render drew %d marks for %d fixture samples", len(rendered.Route.Marks), len(fixtureTickets))
	}
	if rendered.Route.Marks[0].Fill != "var(--route-b)" {
		t.Errorf("route is not the default channel: first ordinary mark fill = %q, want route B", rendered.Route.Marks[0].Fill)
	}
	if rendered.UserRequest.Marks[0].Fill == rendered.UserRequest.Marks[1].Fill {
		t.Errorf("two named URs share fill %q before the categorical palette is exhausted", rendered.UserRequest.Marks[0].Fill)
	}
	if rendered.UserRequest.Marks[12].Fill != rendered.UserRequest.Marks[13].Fill {
		t.Errorf("the thirteenth and fourteenth named UR fills %q / %q differ — the stated Other URs bucket must be shared", rendered.UserRequest.Marks[12].Fill, rendered.UserRequest.Marks[13].Fill)
	}
	if !strings.Contains(rendered.UserRequest.Legend, "UR") || !strings.Contains(rendered.UserRequest.Legend, "Other URs") || !strings.Contains(rendered.UserRequest.Legend, "No UR recorded") {
		t.Errorf("UR legend %q does not name its active channel, overflow, and missing-value rule", rendered.UserRequest.Legend)
	}
	if !strings.Contains(rendered.Domain.Legend, "Domain") || !strings.Contains(rendered.Domain.Legend, "No domain recorded") {
		t.Errorf("domain legend %q does not name its active channel and missing-value rule", rendered.Domain.Legend)
	}
	if rendered.Domain.Marks[0].Fill == rendered.Domain.Marks[1].Fill {
		t.Errorf("different named domains share fill %q before the palette is exhausted", rendered.Domain.Marks[0].Fill)
	}
	missingMarkIndex := len(fixtureTickets) - 2
	if !strings.Contains(rendered.UserRequest.Marks[missingMarkIndex].Class, "unknown") || !strings.Contains(rendered.Domain.Marks[missingMarkIndex].Class, "unknown") {
		t.Errorf("missing UR/domain sample classes = %q / %q, want the visually distinct unknown mark", rendered.UserRequest.Marks[missingMarkIndex].Class, rendered.Domain.Marks[missingMarkIndex].Class)
	}
	criticalMarkIndex := len(fixtureTickets) - 1
	for channelName, channel := range map[string]renderedChannel{"route": rendered.Route, "user-request": rendered.UserRequest, "domain": rendered.Domain} {
		if channel.Marks[criticalMarkIndex].Fill != "var(--durations-critical)" {
			t.Errorf("%s colour channel changed reversed-stamp fill to %q", channelName, channel.Marks[criticalMarkIndex].Fill)
		}
		if !strings.Contains(channel.AriaLabel, "coloured by ") {
			t.Errorf("%s chart aria label %q does not name the active colour channel", channelName, channel.AriaLabel)
		}
	}
}

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
func TestJavaScriptBehaviorTimelineZoomHoldsTheAnchorInstant(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := timelineProbePreamble(t, "TIMELINE_MIN_SPAN_MS") +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineZoomedWindow(") + `
var boundStart = 0;
var boundEnd = 30 * 24 * 3600 * 1000;   // a 30-day board
var startWindow = { windowStartMs: boundStart, windowEndMs: boundEnd };

function anchorInstant(window, fraction) {
  return window.windowStartMs + (window.windowEndMs - window.windowStartMs) * fraction;
}

// Zoom in three times at the same off-centre anchor; the instant under it must
// not move.
var anchorFraction = 0.25;
var wantAnchor = anchorInstant(startWindow, anchorFraction);
var zoomed = startWindow;
for (var step = 0; step < 3; step++) {
  zoomed = timelineZoomedWindow(zoomed.windowStartMs, zoomed.windowEndMs, 1.6, anchorFraction, boundStart, boundEnd);
}
var anchorDriftMs = Math.abs(anchorInstant(zoomed, anchorFraction) - wantAnchor);

// Zooming all the way back out clamps to the bounds rather than overshooting.
var wideOpen = zoomed;
for (var back = 0; back < 12; back++) {
  wideOpen = timelineZoomedWindow(wideOpen.windowStartMs, wideOpen.windowEndMs, 1 / 1.6, 0.5, boundStart, boundEnd);
}

// Zooming all the way in stops at the floor rather than collapsing to zero.
var deep = startWindow;
for (var deeper = 0; deeper < 40; deeper++) {
  deep = timelineZoomedWindow(deep.windowStartMs, deep.windowEndMs, 1.6, 0.5, boundStart, boundEnd);
}

process.stdout.write(JSON.stringify({
  anchorDriftMs: anchorDriftMs,
  widestSpanMs: wideOpen.windowEndMs - wideOpen.windowStartMs,
  boundSpanMs: boundEnd - boundStart,
  deepestSpanMs: deep.windowEndMs - deep.windowStartMs,
  minSpanMs: TIMELINE_MIN_SPAN_MS,
  withinBounds: wideOpen.windowStartMs >= boundStart && wideOpen.windowEndMs <= boundEnd
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline zoom", javascriptProbe)
	var zoomResult struct {
		AnchorDriftMs float64 `json:"anchorDriftMs"`
		WidestSpanMs  float64 `json:"widestSpanMs"`
		BoundSpanMs   float64 `json:"boundSpanMs"`
		DeepestSpanMs float64 `json:"deepestSpanMs"`
		MinSpanMs     float64 `json:"minSpanMs"`
		WithinBounds  bool    `json:"withinBounds"`
	}
	if decodeError := json.Unmarshal(probeOutput, &zoomResult); decodeError != nil {
		t.Fatalf("decode timeline zoom behavior: %v (output %q)", decodeError, probeOutput)
	}
	if zoomResult.AnchorDriftMs > 1 {
		t.Fatalf("the anchored instant drifted %.0f ms over three zoom steps; it must stay put",
			zoomResult.AnchorDriftMs)
	}
	if zoomResult.WidestSpanMs != zoomResult.BoundSpanMs {
		t.Fatalf("zooming out settled at %.0f ms, want the full bound span %.0f ms",
			zoomResult.WidestSpanMs, zoomResult.BoundSpanMs)
	}
	if !zoomResult.WithinBounds {
		t.Fatal("the zoomed-out window escaped its bounds")
	}
	if zoomResult.DeepestSpanMs != zoomResult.MinSpanMs {
		t.Fatalf("zooming in settled at %.0f ms, want the %.0f ms floor",
			zoomResult.DeepestSpanMs, zoomResult.MinSpanMs)
	}
}

// Virtualization is what makes 560 rows cost the same as 40, and it is invisible
// in a screenshot: a wrong slice shows blank strips only while scrolling fast.
func TestJavaScriptBehaviorTimelineVirtualizesRowsAtScale(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := timelineProbePreamble(t, "TIMELINE_ROW_HEIGHT", "TIMELINE_OVERSCAN_ROWS") +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineVisibleRowRange(") + `
var rowCount = 560;
var viewportHeight = 600;
var atTop = timelineVisibleRowRange(0, viewportHeight, rowCount);
var midway = timelineVisibleRowRange(rowCount * TIMELINE_ROW_HEIGHT / 2, viewportHeight, rowCount);
var atBottom = timelineVisibleRowRange(rowCount * TIMELINE_ROW_HEIGHT, viewportHeight, rowCount);
process.stdout.write(JSON.stringify({
  atTopCount: atTop.lastRow - atTop.firstRow,
  midwayCount: midway.lastRow - midway.firstRow,
  midwayCoversScrollPosition:
    midway.firstRow <= rowCount / 2 && midway.lastRow > rowCount / 2,
  atBottomLastRow: atBottom.lastRow,
  rowCount: rowCount
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline virtualization", javascriptProbe)
	var sliceResult struct {
		AtTopCount                 int  `json:"atTopCount"`
		MidwayCount                int  `json:"midwayCount"`
		MidwayCoversScrollPosition bool `json:"midwayCoversScrollPosition"`
		AtBottomLastRow            int  `json:"atBottomLastRow"`
		RowCount                   int  `json:"rowCount"`
	}
	if decodeError := json.Unmarshal(probeOutput, &sliceResult); decodeError != nil {
		t.Fatalf("decode timeline virtualization behavior: %v (output %q)", decodeError, probeOutput)
	}
	// A 600px viewport at the shipped row height holds well under 60 rows; the
	// point of the assertion is that the slice is bounded by the VIEWPORT and
	// never by the row count, which is what a non-virtualized render would do.
	if sliceResult.AtTopCount >= sliceResult.RowCount/4 {
		t.Fatalf("a 600px viewport rendered %d of %d rows; the slice must be viewport-bounded",
			sliceResult.AtTopCount, sliceResult.RowCount)
	}
	if sliceResult.MidwayCount != sliceResult.AtTopCount {
		t.Fatalf("slice size changed with scroll position (%d at top, %d midway); it must not",
			sliceResult.AtTopCount, sliceResult.MidwayCount)
	}
	if !sliceResult.MidwayCoversScrollPosition {
		t.Fatal("the midway slice does not contain the row at the scroll position")
	}
	if sliceResult.AtBottomLastRow != sliceResult.RowCount {
		t.Fatalf("scrolled past the end the slice reached row %d, want it clamped to %d",
			sliceResult.AtBottomLastRow, sliceResult.RowCount)
	}
}

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
func TestJavaScriptBehaviorTimelineRowsFollowTheWindow(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := sliceBalancedBlockAfter(t, indexHtml, "function timelineWaitEndMs(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineWorkEndMs(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineRowSegments(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineRowsInWindow(") + `
var nowMs = Date.UTC(2026, 5, 20, 12, 0);
var windowStartMs = Date.UTC(2026, 5, 10);
var windowEndMs = Date.UTC(2026, 5, 17);

var rows = [
  // Wholly before the window.
  { id: "REQ-001", createdTime: "2026-06-01T09:00:00Z", claimedTime: "2026-06-02T09:00:00Z",
    completedTime: "2026-06-03T09:00:00Z", hasWork: true },
  // Wholly inside it.
  { id: "REQ-002", createdTime: "2026-06-12T09:00:00Z", claimedTime: "2026-06-13T09:00:00Z",
    completedTime: "2026-06-14T09:00:00Z", hasWork: true },
  // Started before, finished after: STRADDLES the window and never sits inside it.
  { id: "REQ-003", createdTime: "2026-06-05T09:00:00Z", claimedTime: "2026-06-06T09:00:00Z",
    completedTime: "2026-06-25T09:00:00Z", hasWork: true },
  // Captured before the window and STILL RUNNING at the now-line beyond it.
  { id: "REQ-004", createdTime: "2026-06-02T09:00:00Z", claimedTime: "2026-06-03T09:00:00Z",
    hasWork: true, workOpen: true },
  // Wholly after the window.
  { id: "REQ-005", createdTime: "2026-06-18T09:00:00Z", claimedTime: "2026-06-19T09:00:00Z",
    completedTime: "2026-06-19T18:00:00Z", hasWork: true }
];
// A reversed stamp: claimed BEFORE created. It carries waitMinutes because that
// is the field BOTH the renderer and the segment list branch on — without it the
// fixture is silently a forward row and tests nothing about reversal.
//
// No segment may come out with its start after its end, or the overlap test
// inverts and the row vanishes from every window. A reversed wait satisfies that
// as a point at the created instant, which is also where the break marker goes.
var reversedRow = { id: "REQ-006", createdTime: "2026-06-14T09:00:00Z",
  claimedTime: "2026-06-12T09:00:00Z", completedTime: "2026-06-15T09:00:00Z",
  hasWork: true, anomaly: true, waitMinutes: -2880, workMinutes: 4320 };
// A pending REQ whose only presence in a FUTURE window is its forecast bar. Its
// measured extent is the open wait, which stops at the now-line; the window sits
// entirely after that, so the row reaches it only if the extent follows the
// projected segment the chart actually draws.
var projectedRow = { id: "REQ-007", createdTime: "2026-06-01T09:00:00Z", waitOpen: true };
var projection = { id: "REQ-007", startTime: "2026-06-21T00:00:00Z", endTime: "2026-06-21T06:00:00Z" };
var futureWindowStartMs = Date.UTC(2026, 5, 20, 18, 0);
var futureWindowEndMs = Date.UTC(2026, 5, 22);

function idsInSpan(rowList, projectedById, spanStartMs, spanEndMs) {
  var extents = rowList.map(function (row) {
    return timelineRowSegments(row, nowMs, (projectedById || {})[row.id]);
  });
  return timelineRowsInWindow(rowList, extents, spanStartMs, spanEndMs)
    .map(function (row) { return row.id; });
}
function idsInWindow(rowList, projectedById) {
  return idsInSpan(rowList, projectedById, windowStartMs, windowEndMs);
}

var reversedSegments = timelineRowSegments(reversedRow, nowMs, undefined);

// THE EXCLUSIVE END. Every window this view can build ends on the NEXT period's
// first instant — timelinePeriodWindow returns [start, nextStart) and the end
// date field renders windowEndMs - 1 to say so. A segment that begins exactly at
// that instant belongs to the next window, and drawSegment puts its floored
// rectangle at the right edge where it is clipped: listed, nothing drawn.
var edgeWindowStartMs = Date.UTC(2026, 5, 10);
var edgeWindowEndMs = Date.UTC(2026, 5, 17);
var startsAtWindowEndRow = {
  id: "REQ-008",
  createdTime: new Date(edgeWindowEndMs).toISOString(),
  claimedTime: new Date(edgeWindowEndMs + 3600000).toISOString(),
  completedTime: new Date(edgeWindowEndMs + 7200000).toISOString(),
  hasWork: true, waitMinutes: 60, workMinutes: 60
};
// The control, one millisecond earlier: genuinely inside, and must stay listed.
var startsJustInsideRow = {
  id: "REQ-009",
  createdTime: new Date(edgeWindowEndMs - 1).toISOString(),
  claimedTime: new Date(edgeWindowEndMs + 3600000).toISOString(),
  completedTime: new Date(edgeWindowEndMs + 7200000).toISOString(),
  hasWork: true, waitMinutes: 60, workMinutes: 60
};
// The symmetric case at the other end, which must NOT change: windowStartMs is
// inclusive, and a span ending exactly there draws a floored mark at x=0 that
// the reader can see.
var endsAtWindowStartRow = {
  id: "REQ-010",
  createdTime: new Date(edgeWindowStartMs - 7200000).toISOString(),
  claimedTime: new Date(edgeWindowStartMs).toISOString(),
  hasWork: false, waitMinutes: 120
};

// A REVERSED WAIT drawn as what the renderer actually draws. renderVisibleRows
// puts a 6px break marker at the CREATED instant for a reversed wait and at the
// CLAIMED instant for reversed work — it does not draw a bar across the min/max
// interval. Modelling the row as that interval is the forecast-gap defect again
// in another costume: created 14 Jun, claimed 12 Jun, completed 12 Jun 06:00
// puts the hull across 13 June while both drawn marks sit outside it.
var reversedHullRow = {
  id: "REQ-011",
  createdTime: "2026-06-14T12:00:00Z",
  claimedTime: "2026-06-12T00:00:00Z",
  completedTime: "2026-06-12T06:00:00Z",
  hasWork: true, waitMinutes: -3600, workMinutes: 360, anomaly: true
};
var reversedHullGapStartMs = Date.UTC(2026, 5, 13);
var reversedHullGapEndMs = Date.UTC(2026, 5, 13, 23, 59);
// The control: a window over the break marker itself must still list it.
var reversedHullMarkerStartMs = Date.UTC(2026, 5, 14, 6, 0);
var reversedHullMarkerEndMs = Date.UTC(2026, 5, 14, 18, 0);

// THE GAP. A pending REQ draws two disjoint marks: the open wait ending at the
// now-line, and the forecast bar starting after in-flight work finishes. A hull
// over both spans the gap between them, so a window sitting in that gap listed
// the row with nothing drawn on it. Segment-wise overlap is what makes the REQ's
// GREEN — every listed row has something drawn on it — actually true.
var gapWindowStartMs = Date.UTC(2026, 5, 20, 14, 0);   // after the now-line
var gapWindowEndMs = Date.UTC(2026, 5, 20, 18, 0);     // before the forecast bar
var gapProjection = { id: "REQ-007", startTime: "2026-06-21T00:00:00Z", endTime: "2026-06-21T06:00:00Z" };
process.stdout.write(JSON.stringify({
  inWindow: idsInWindow(rows, {}),
  reversedInWindow: idsInWindow([reversedRow], {}),
  reversedExtentOrdered: reversedSegments.every(function (s) { return s.startMs <= s.endMs; }),
  reversedExtentStartIso: new Date(reversedSegments[0].startMs).toISOString(),
  projectedOnlyInWindow:
    idsInSpan([projectedRow], { "REQ-007": projection }, futureWindowStartMs, futureWindowEndMs),
  projectedIgnoredWithoutForecast:
    idsInSpan([projectedRow], {}, futureWindowStartMs, futureWindowEndMs),
  everythingInAWideWindow: idsInSpan(rows, {}, Date.UTC(2026, 0, 1), Date.UTC(2027, 0, 1)),
  inTheForecastGap:
    idsInSpan([projectedRow], { "REQ-007": gapProjection }, gapWindowStartMs, gapWindowEndMs),
  spanningTheForecastGap:
    idsInSpan([projectedRow], { "REQ-007": gapProjection }, nowMs - 3600000, Date.UTC(2026, 5, 21, 3, 0)),

  startsAtWindowEnd: idsInSpan([startsAtWindowEndRow], {}, edgeWindowStartMs, edgeWindowEndMs),
  startsJustInside: idsInSpan([startsJustInsideRow], {}, edgeWindowStartMs, edgeWindowEndMs),
  endsAtWindowStart: idsInSpan([endsAtWindowStartRow], {}, edgeWindowStartMs, edgeWindowEndMs),

  reversedHullInTheGap:
    idsInSpan([reversedHullRow], {}, reversedHullGapStartMs, reversedHullGapEndMs),
  reversedHullAtItsMarker:
    idsInSpan([reversedHullRow], {}, reversedHullMarkerStartMs, reversedHullMarkerEndMs)
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline window rows", javascriptProbe)
	var windowResult struct {
		InWindow                        []string `json:"inWindow"`
		ReversedInWindow                []string `json:"reversedInWindow"`
		ReversedExtentOrdered           bool     `json:"reversedExtentOrdered"`
		ReversedExtentStartIso          string   `json:"reversedExtentStartIso"`
		ProjectedOnlyInWindow           []string `json:"projectedOnlyInWindow"`
		ProjectedIgnoredWithoutForecast []string `json:"projectedIgnoredWithoutForecast"`
		EverythingInAWideWindow         []string `json:"everythingInAWideWindow"`
		InTheForecastGap                []string `json:"inTheForecastGap"`
		SpanningTheForecastGap          []string `json:"spanningTheForecastGap"`

		StartsAtWindowEnd []string `json:"startsAtWindowEnd"`
		StartsJustInside  []string `json:"startsJustInside"`
		EndsAtWindowStart []string `json:"endsAtWindowStart"`

		ReversedHullInTheGap    []string `json:"reversedHullInTheGap"`
		ReversedHullAtItsMarker []string `json:"reversedHullAtItsMarker"`
	}
	if decodeError := json.Unmarshal(probeOutput, &windowResult); decodeError != nil {
		t.Fatalf("decode timeline window rows behavior: %v (output %q)", decodeError, probeOutput)
	}

	wantInWindow := "REQ-002,REQ-003,REQ-004"
	if gotInWindow := strings.Join(windowResult.InWindow, ","); gotInWindow != wantInWindow {
		t.Fatalf("rows in the window = %q, want %q — the straddling and still-running REQs "+
			"overlap it and must be listed; the two outside it must not be",
			gotInWindow, wantInWindow)
	}
	if !windowResult.ReversedExtentOrdered {
		t.Fatalf("a reversed-stamp row produced a segment whose start (%s) is after its end; "+
			"a reversed span is a point at the break marker's own instant, and an inverted "+
			"segment makes the overlap test drop broken rows from every window",
			windowResult.ReversedExtentStartIso)
	}
	if len(windowResult.ReversedInWindow) != 1 {
		t.Fatalf("the reversed-stamp row is inside the window and was not listed (got %v)",
			windowResult.ReversedInWindow)
	}
	if len(windowResult.ProjectedOnlyInWindow) != 1 {
		t.Fatal("a pending REQ whose forecast bar is the only thing it draws in the window " +
			"was not listed; the extent must reach the projected segment")
	}
	// The control for the assertion above: same row, same window, no forecast.
	// Without it the extent stops at the now-line and the row is genuinely absent,
	// which is what proves the projection — not the fixture — put it in the window.
	if len(windowResult.ProjectedIgnoredWithoutForecast) != 0 {
		t.Fatalf("the same pending REQ reached a window beyond the now-line with no forecast "+
			"attached (got %v); its measured extent ends at the now-line",
			windowResult.ProjectedIgnoredWithoutForecast)
	}
	if len(windowResult.EverythingInAWideWindow) != 5 {
		t.Fatalf("a window spanning the whole year listed %d of 5 rows; widening the window "+
			"must never drop a row", len(windowResult.EverythingInAWideWindow))
	}
	// The pair that forces segment-wise overlap rather than a hull. Both windows
	// sit between the same two marks; only the second one touches either.
	if len(windowResult.InTheForecastGap) != 0 {
		t.Fatalf("a window in the gap between a REQ's open wait and its forecast bar listed "+
			"%v; the row draws nothing there, and listing it is what window-scoping exists to "+
			"stop — a hull over the two marks spans the gap between them",
			windowResult.InTheForecastGap)
	}
	if len(windowResult.SpanningTheForecastGap) != 1 {
		t.Fatal("a window reaching across both the open wait and the forecast bar did not list " +
			"the row; the gap rule must not cost a row that genuinely draws inside the window")
	}

	// THE EXCLUSIVE END, and its two controls. windowEndMs is the next window's
	// first instant everywhere else in this module — timelinePeriodWindow builds
	// [start, nextStart) and the end field renders windowEndMs - 1 to say so — so
	// admitting a segment that begins exactly there lists a row whose only mark is
	// clipped at the right edge.
	if len(windowResult.StartsAtWindowEnd) != 0 {
		t.Errorf("a REQ whose span begins exactly at the window's exclusive end was listed "+
			"(got %v); its floored rectangle lands at the clipped right edge, so the row shows "+
			"nothing — and it belongs to the NEXT window", windowResult.StartsAtWindowEnd)
	}
	if len(windowResult.StartsJustInside) != 1 {
		t.Errorf("the same REQ one millisecond earlier was not listed (got %v); the end is "+
			"exclusive by one instant, not by a margin", windowResult.StartsJustInside)
	}
	// Deliberately asymmetric. The START instant IS in the window and a span
	// ending on it draws a visible floored mark at x=0, so this stays inclusive.
	if len(windowResult.EndsAtWindowStart) != 1 {
		t.Errorf("a REQ whose span ends exactly at the window's start was dropped (got %v); "+
			"the start is inclusive and that span draws a mark at the left edge",
			windowResult.EndsAtWindowStart)
	}

	// A REVERSED SPAN IS A POINT, because that is what the renderer draws. The
	// same defect as the forecast gap: a hull over an interval nothing is drawn
	// across lists rows with nothing on them.
	if len(windowResult.ReversedHullInTheGap) != 0 {
		t.Errorf("a window inside a reversed span's min/max interval listed the row (got %v), "+
			"but renderVisibleRows draws only a break marker at the created instant and the "+
			"row's other bar sits elsewhere — nothing is drawn in that window",
			windowResult.ReversedHullInTheGap)
	}
	if len(windowResult.ReversedHullAtItsMarker) != 1 {
		t.Errorf("a window over the reversed span's own break marker did not list the row "+
			"(got %v); the point has to sit where the renderer puts the marker, or a broken "+
			"row becomes unfindable — which is what the min/max existed to prevent",
			windowResult.ReversedHullAtItsMarker)
	}
}

// Typed dates are the fourth way to move the window, and the one a reader can
// aim: the other three (pointer, keyboard, period chips) are all relative. What
// has to hold is that typing cannot reach a window the other three cannot —
// same floor, same ceiling, same edge clamp — because a control with its own
// clamp is how a view starts disagreeing with itself.
func TestJavaScriptBehaviorTimelineTypedDatesMoveTheWindow(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := timelineProbePreamble(t, "TIMELINE_MIN_SPAN_MS", "TIMELINE_DAY_MS") +
		rendererDeclarationLine(t, "web/board-timeline.js", "TIMELINE_PERIOD_LEVEL_NAMES") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineZoomedWindow(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineDateFieldToEpoch(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineStartEpochToDateField(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineEndEpochToDateField(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelinePeriodStart(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineSteppedPeriodStart(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelinePeriodLevelOfWindow(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineTypedWindow(") + `
var boundStart = Date.UTC(2026, 3, 7);
var boundEnd = Date.UTC(2026, 8, 2);
var windowStart = Date.UTC(2026, 5, 1);
var windowEnd = Date.UTC(2026, 5, 8);

function typed(startText, endText) {
  return timelineTypedWindow(startText, endText, windowStart, windowEnd, boundStart, boundEnd);
}
function iso(epochMs) { return new Date(epochMs).toISOString(); }

var bothFields = typed("2026-06-01", "2026-06-15");
var startOnly = typed("2026-06-03", "");
var endOnly = typed("", "2026-06-20");
var sameDay = typed("2026-07-04", "2026-07-04");
var reversed = typed("2026-07-10", "2026-07-01");
var beforeRange = typed("2020-01-01", "2020-01-31");
// A date PAST THE END of the range, typed into From with the end field left
// alone. It has to land on the last day the board HAS, not collapse both
// endpoints onto the bound and leave an empty zoom-floor sliver behind the frame.
var pastRangeStartOnly = typed("2026-09-30", "");
// And the mirror, typed into the end field.
var pastRangeEndOnly = typed("", "2026-09-30");
// A start typed while the end field still holds the board's last day. The
// implied span overruns the ceiling, and a span-preserving settle would pin the
// end to the bound and drag this start backwards to keep the width.
var startAgainstCeiling = timelineTypedWindow(
  "2026-08-01", timelineEndEpochToDateField(boundEnd), windowStart, windowEnd, boundStart, boundEnd);
var neither = typed("", "");
var rubbish = typed("not-a-date", "2026-13-45");
var rolled = timelineDateFieldToEpoch("2026-02-31");

process.stdout.write(JSON.stringify({
  bothStartIso: iso(bothFields.windowStartMs),
  bothEndIso: iso(bothFields.windowEndMs),
  startOnlyStartIso: iso(startOnly.windowStartMs),
  startOnlyKeptEnd: startOnly.windowEndMs === windowEnd,
  endOnlyKeptStart: endOnly.windowStartMs === windowStart,
  endOnlyEndIso: iso(endOnly.windowEndMs),
  sameDaySpanMs: sameDay.windowEndMs - sameDay.windowStartMs,
  reversedOrdered: reversed.windowStartMs < reversed.windowEndMs,
  reversedStartIso: iso(reversed.windowStartMs),
  beforeRangeClampedToBound: beforeRange.windowStartMs >= boundStart,
  beforeRangeStartIso: iso(beforeRange.windowStartMs),
  startAgainstCeilingIso: iso(startAgainstCeiling.windowStartMs),
  pastRangeStartIso: iso(pastRangeStartOnly.windowStartMs),
  pastRangeEndIso: iso(pastRangeStartOnly.windowEndMs),
  pastRangeSpanMs: pastRangeStartOnly.windowEndMs - pastRangeStartOnly.windowStartMs,
  pastRangeEndOnlyEndIso: iso(pastRangeEndOnly.windowEndMs),
  lastDayStartMs: timelinePeriodStart(boundEnd - 1, "day"),
  pastRangeStartMs: pastRangeStartOnly.windowStartMs,
  minSpanMs: TIMELINE_MIN_SPAN_MS,
  // THE ROUND TRIP. Render a real calendar window into the two fields, parse them
  // straight back, and the same window has to come out — otherwise editing one
  // field re-applies a mangled version of the other, and no typed pair can ever
  // equal what a period chip produces.
  periodRoundTrips: (function () {
    return ["day", "week", "month"].map(function (levelName) {
      var periodStartMs = timelinePeriodStart(Date.UTC(2026, 6, 15, 9, 30), levelName);
      var periodEndMs = timelineSteppedPeriodStart(periodStartMs, levelName, 1);
      var reparsed = timelineTypedWindow(
        timelineStartEpochToDateField(periodStartMs),
        timelineEndEpochToDateField(periodEndMs),
        0, 0, boundStart, boundEnd);
      return {
        level: levelName,
        fields: timelineStartEpochToDateField(periodStartMs) + ".." + timelineEndEpochToDateField(periodEndMs),
        exact: reparsed.windowStartMs === periodStartMs && reparsed.windowEndMs === periodEndMs,
        reparsedLevel: timelinePeriodLevelOfWindow(reparsed.windowStartMs, reparsed.windowEndMs)
      };
    });
  })(),
  neitherIsNull: neither === null,
  rubbishIsNull: rubbish === null,
  rolledIsNaN: isNaN(rolled),
  roundTrip: timelineStartEpochToDateField(Date.UTC(2026, 5, 9, 13, 45))
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline typed dates", javascriptProbe)
	var typedResult struct {
		BothStartIso              string  `json:"bothStartIso"`
		BothEndIso                string  `json:"bothEndIso"`
		StartOnlyStartIso         string  `json:"startOnlyStartIso"`
		StartOnlyKeptEnd          bool    `json:"startOnlyKeptEnd"`
		EndOnlyKeptStart          bool    `json:"endOnlyKeptStart"`
		EndOnlyEndIso             string  `json:"endOnlyEndIso"`
		SameDaySpanMs             float64 `json:"sameDaySpanMs"`
		ReversedOrdered           bool    `json:"reversedOrdered"`
		ReversedStartIso          string  `json:"reversedStartIso"`
		BeforeRangeClampedToBound bool    `json:"beforeRangeClampedToBound"`
		BeforeRangeStartIso       string  `json:"beforeRangeStartIso"`
		StartAgainstCeilingIso    string  `json:"startAgainstCeilingIso"`
		PastRangeStartIso         string  `json:"pastRangeStartIso"`
		PastRangeEndIso           string  `json:"pastRangeEndIso"`
		PastRangeSpanMs           float64 `json:"pastRangeSpanMs"`
		PastRangeEndOnlyEndIso    string  `json:"pastRangeEndOnlyEndIso"`
		LastDayStartMs            float64 `json:"lastDayStartMs"`
		PastRangeStartMs          float64 `json:"pastRangeStartMs"`
		MinSpanMs                 float64 `json:"minSpanMs"`
		PeriodRoundTrips          []struct {
			Level         string `json:"level"`
			Fields        string `json:"fields"`
			Exact         bool   `json:"exact"`
			ReparsedLevel string `json:"reparsedLevel"`
		} `json:"periodRoundTrips"`
		NeitherIsNull bool   `json:"neitherIsNull"`
		RubbishIsNull bool   `json:"rubbishIsNull"`
		RolledIsNaN   bool   `json:"rolledIsNaN"`
		RoundTrip     string `json:"roundTrip"`
	}
	if decodeError := json.Unmarshal(probeOutput, &typedResult); decodeError != nil {
		t.Fatalf("decode timeline typed dates behavior: %v (output %q)", decodeError, probeOutput)
	}

	if typedResult.BothStartIso != "2026-06-01T00:00:00.000Z" {
		t.Fatalf("typed start = %s, want the UTC midnight opening that day", typedResult.BothStartIso)
	}
	// The end field names a day to INCLUDE and the window's end is EXCLUSIVE, so
	// the day typed resolves to the FOLLOWING midnight. That is not cosmetic: it
	// is what makes a typed pair produce byte-identical windows to the period
	// chips, which is what the round-trip block below turns into a hard rule.
	// (This assertion previously wanted 23:59:59.999 — the inclusive end that
	// made render and parse non-inverses. Changed deliberately, not bent to fit.)
	if typedResult.BothEndIso != "2026-06-16T00:00:00.000Z" {
		t.Fatalf("typed end = %s, want the midnight following the day typed", typedResult.BothEndIso)
	}
	if !typedResult.StartOnlyKeptEnd || typedResult.StartOnlyStartIso != "2026-06-03T00:00:00.000Z" {
		t.Fatalf("typing only a start moved the end too (start %s, kept end %v); each field must "+
			"resolve against the window already on screen",
			typedResult.StartOnlyStartIso, typedResult.StartOnlyKeptEnd)
	}
	if !typedResult.EndOnlyKeptStart || typedResult.EndOnlyEndIso != "2026-06-21T00:00:00.000Z" {
		t.Fatalf("typing only an end moved the start too (end %s, kept start %v)",
			typedResult.EndOnlyEndIso, typedResult.EndOnlyKeptStart)
	}
	// One date in both fields is the commonest thing a reader will do with a date
	// picker, and it must mean that day rather than an empty window.
	if typedResult.SameDaySpanMs != 86400000 {
		t.Fatalf("the same date in both fields spanned %.0f ms, want exactly one day — the same "+
			"window the Day chip produces, so the chip can light for it",
			typedResult.SameDaySpanMs)
	}
	if !typedResult.ReversedOrdered || typedResult.ReversedStartIso != "2026-07-10T00:00:00.000Z" {
		t.Fatalf("an end before the start produced %s and ordered=%v; it must clamp forward from "+
			"the start the reader typed, never silently swap the two",
			typedResult.ReversedStartIso, typedResult.ReversedOrdered)
	}
	// The clamp has to be the shared one. A control with its own bounds is how
	// the reader reaches a window no other control can, and then cannot get back.
	if !typedResult.BeforeRangeClampedToBound {
		t.Fatalf("a date before the board's range escaped the bounds (start %s)",
			typedResult.BeforeRangeStartIso)
	}
	// The defect a browser found and the unit fixture had missed: each endpoint is
	// clamped on its own, because a typed date is a position and the shared settle
	// preserves a span. Pinning the end to the ceiling must never move the start
	// the reader typed.
	if typedResult.StartAgainstCeilingIso != "2026-08-01T00:00:00.000Z" {
		t.Fatalf("a start typed against the range ceiling came back as %s, want the date typed; "+
			"the settle preserves a span and would drag the start back to keep the width",
			typedResult.StartAgainstCeilingIso)
	}
	// One assertion covering the whole class: the fields must be a lossless view
	// of a calendar window, and a typed pair must be able to name the same window
	// a chip does. Rendering an exclusive end instant's own date failed both —
	// every period window came back a day long, and no typed pair could ever
	// light a chip.
	if len(typedResult.PeriodRoundTrips) != 3 {
		t.Fatalf("want a round trip for each of day, week and month; got %d",
			len(typedResult.PeriodRoundTrips))
	}
	for _, roundTrip := range typedResult.PeriodRoundTrips {
		if !roundTrip.Exact {
			t.Errorf("a %s window rendered to fields %s and did not parse back to itself; "+
				"render and parse must be inverses or editing one field mangles the other",
				roundTrip.Level, roundTrip.Fields)
		}
		if roundTrip.ReparsedLevel != roundTrip.Level {
			t.Errorf("a %s window typed back as fields %s reads as level %q; a typed pair must be "+
				"able to name the same window a period chip produces",
				roundTrip.Level, roundTrip.Fields, roundTrip.ReparsedLevel)
		}
	}
	if !typedResult.NeitherIsNull {
		t.Fatal("two empty fields must return null — a cleared field is not a request to move")
	}
	if !typedResult.RubbishIsNull {
		t.Fatal("unparseable text in both fields must return null rather than moving the window")
	}
	if !typedResult.RolledIsNaN {
		t.Fatal("2026-02-31 must be rejected; Date.UTC rolls it into March and a rolled date is " +
			"not the one that was typed")
	}
	// A DATE PAST THE END OF THE RANGE lands on the last day the board has. Before
	// this, both endpoints collapsed onto the bound and the settle turned that into
	// an empty one-hour window tucked behind the right edge, while the field went on
	// showing the rejected date.
	if typedResult.PastRangeStartMs != typedResult.LastDayStartMs {
		t.Errorf("a From date past the end of the range put the window start at %s, want the last "+
			"day the board has", typedResult.PastRangeStartIso)
	}
	if typedResult.PastRangeSpanMs <= typedResult.MinSpanMs {
		t.Errorf("a From date past the end of the range produced a %.0f ms window (%s → %s); at the "+
			"zoom floor or below it is the empty sliver this clamp exists to prevent",
			typedResult.PastRangeSpanMs, typedResult.PastRangeStartIso, typedResult.PastRangeEndIso)
	}
	if typedResult.PastRangeEndOnlyEndIso != "2026-09-02T00:00:00.000Z" {
		t.Errorf("a `to` date past the end of the range ended the window at %s, want the range's "+
			"own end", typedResult.PastRangeEndOnlyEndIso)
	}

	if typedResult.RoundTrip != "2026-06-09" {
		t.Fatalf("an instant mid-day rendered into the date field as %q, want its UTC date",
			typedResult.RoundTrip)
	}
}

// The label's truncation arithmetic. The MEASUREMENT it depends on is a browser
// question and lives in the browser lane; this pins what the module does with
// the number once it has it.
//
// The id is the row's identity, so the rule under test is that it is never the
// thing that gets cut: a budget too small for id-plus-a-useful-title drops the
// title entirely rather than shipping a half-drawn id.
func TestJavaScriptBehaviorTimelineRowLabelTruncation(t *testing.T) {
	indexHtml := generateLiveSite(t)
	// The column width is READ FROM THE RENDERER, not restated here. Restating it
	// is how the first version of this test passed with the constant reverted to
	// its old value: the budget floor below was measuring a column the board had
	// stopped using (REQ-265 — grep the quantity, not the constant name).
	javascriptProbe := timelineProbePreamble(t, "TIMELINE_LABEL_WIDTH") +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineLabelCharacterBudget(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineLabelCellCount(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineLabelPrefixWithinCells(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineRowLabelText(") + `
var longTitle = "Colour timeline bars by REQ status";
// The shipped column, minus the label's own 6px x-offset and a 6px gap before
// the first bar — the same expression the renderer uses.
var shippedColumnCells = timelineLabelCharacterBudget(6.0219, TIMELINE_LABEL_WIDTH - 12);
process.stdout.write(JSON.stringify({
  shippedColumnCells: shippedColumnCells,
  shippedLabelSample: timelineRowLabelText("REQ-042", longTitle, shippedColumnCells),
  taggedLabelSample: timelineRowLabelText(
    "REQ-306", "[impact-rule-change] Judge effort_estimate on review-minted follow-ups too",
    shippedColumnCells),
  cjkCells: timelineLabelCellCount("修复部署管道"),
  asciiCells: timelineLabelCellCount("abcdef"),
  astralCells: timelineLabelCellCount("🚀"),
  cjkLabel: timelineRowLabelText("REQ-042", "修复部署管道的一个严重问题在这里", 28),
  astralBoundary: timelineRowLabelText("REQ-042", "Fix the deploy 🚀 pipeline right now", 26),
  budgetAt6px: timelineLabelCharacterBudget(6, 172),
  budgetWithNoAdvance: timelineLabelCharacterBudget(0, 172),
  budgetWithNegativeAdvance: timelineLabelCharacterBudget(-3, 172),

  roomy: timelineRowLabelText("REQ-042", longTitle, 60),
  cut: timelineRowLabelText("REQ-042", longTitle, 30),
  cutLength: timelineRowLabelText("REQ-042", longTitle, 30).length,
  tight: timelineRowLabelText("REQ-042", longTitle, 14),
  noBudget: timelineRowLabelText("REQ-042", longTitle, 0),
  noTitle: timelineRowLabelText("REQ-042", "", 60),
  exactFit: timelineRowLabelText("REQ-042", "abcdef", 15),
  oneOver: timelineRowLabelText("REQ-042", "abcdefghijkl", 20),
  longId: timelineRowLabelText("REQ-100042", longTitle, 14)
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline row label", javascriptProbe)
	var labelResult struct {
		ShippedColumnCells        int    `json:"shippedColumnCells"`
		ShippedLabelSample        string `json:"shippedLabelSample"`
		TaggedLabelSample         string `json:"taggedLabelSample"`
		CjkCells                  int    `json:"cjkCells"`
		AsciiCells                int    `json:"asciiCells"`
		AstralCells               int    `json:"astralCells"`
		CjkLabel                  string `json:"cjkLabel"`
		AstralBoundary            string `json:"astralBoundary"`
		BudgetAt6px               int    `json:"budgetAt6px"`
		BudgetWithNoAdvance       int    `json:"budgetWithNoAdvance"`
		BudgetWithNegativeAdvance int    `json:"budgetWithNegativeAdvance"`
		Roomy                     string `json:"roomy"`
		Cut                       string `json:"cut"`
		CutLength                 int    `json:"cutLength"`
		Tight                     string `json:"tight"`
		NoBudget                  string `json:"noBudget"`
		NoTitle                   string `json:"noTitle"`
		ExactFit                  string `json:"exactFit"`
		OneOver                   string `json:"oneOver"`
		LongId                    string `json:"longId"`
	}
	if decodeError := json.Unmarshal(probeOutput, &labelResult); decodeError != nil {
		t.Fatalf("decode timeline row label behavior: %v (output %q)", decodeError, probeOutput)
	}

	// THE COLUMN WIDTH, pinned to the shipped constant. Reverting
	// TIMELINE_LABEL_WIDTH to its pre-REQ 104 used to pass this whole file; now it
	// fails here, which is what makes the width a decision the tests hold rather
	// than a number in a comment.
	if labelResult.ShippedColumnCells < 20 {
		t.Errorf("the shipped label column fits %d cells; below about 20 the title is a stub and "+
			"the column stops earning the plot width it costs (label sample %q)",
			labelResult.ShippedColumnCells, labelResult.ShippedLabelSample)
	}
	// A classification tag is metadata for the board's search box, not title text.
	// Unstripped it consumed the entire budget and every review-minted REQ read
	// "[impact-user-visib…".
	if strings.Contains(labelResult.TaggedLabelSample, "[impact-") {
		t.Errorf("a tagged title rendered %q; the leading [impact-…] tag has to come off in the "+
			"label or it is the whole budget", labelResult.TaggedLabelSample)
	}
	if !strings.Contains(labelResult.TaggedLabelSample, "Judge") {
		t.Errorf("a tagged title rendered %q, want the actual title after the tag came off",
			labelResult.TaggedLabelSample)
	}

	// Cells, not characters. The measured advance describes the face's Latin cell,
	// and a face that draws 中 at 10px against a 6.02px cell overruns the column.
	if labelResult.AsciiCells != 6 {
		t.Errorf("six ASCII characters counted %d cells, want 6", labelResult.AsciiCells)
	}
	if labelResult.CjkCells != 12 {
		t.Errorf("six CJK characters counted %d cells, want 12 — one cell each is what let a CJK "+
			"title draw 36px into the plot", labelResult.CjkCells)
	}
	if labelResult.AstralCells != 2 {
		t.Errorf("one astral character counted %d cells, want 2", labelResult.AstralCells)
	}
	// The cut lands on a code point boundary, so an astral character is never
	// split into a lone surrogate that renders as a fallback box.
	for _, character := range labelResult.AstralBoundary {
		if character == '\uFFFD' {
			t.Errorf("a title cut near an astral character produced %q, which contains a "+
				"replacement character — the cut split a surrogate pair",
				labelResult.AstralBoundary)
		}
	}
	if labelResult.CjkLabel == "" || !strings.HasPrefix(labelResult.CjkLabel, "REQ-042") {
		t.Errorf("a CJK title rendered %q, want a label led by its id", labelResult.CjkLabel)
	}

	if labelResult.BudgetAt6px != 28 {
		t.Errorf("a 172px column at a 6px advance fits %d characters, want 28", labelResult.BudgetAt6px)
	}
	// An unmeasurable face must produce NO budget rather than a plausible one.
	// A guessed advance is the REQ-292 defect: a number that looks like a
	// measurement and does not move when the face does.
	if labelResult.BudgetWithNoAdvance != 0 || labelResult.BudgetWithNegativeAdvance != 0 {
		t.Errorf("an unmeasured face produced budgets %d and %d, want 0 for both",
			labelResult.BudgetWithNoAdvance, labelResult.BudgetWithNegativeAdvance)
	}

	if labelResult.Roomy != "REQ-042  Colour timeline bars by REQ status" {
		t.Errorf("a roomy budget rendered %q, want the id and the whole title", labelResult.Roomy)
	}
	if !strings.HasPrefix(labelResult.Cut, "REQ-042  ") || !strings.HasSuffix(labelResult.Cut, "…") {
		t.Errorf("a cut label rendered %q, want the id then a truncated title ending in an ellipsis",
			labelResult.Cut)
	}
	// The ellipsis is inside the budget, not on top of it — otherwise the label
	// the arithmetic says fits is one character wider than the column.
	if labelResult.CutLength > 30 {
		t.Errorf("a 30-character budget produced a %d-character label %q; the ellipsis has to fit "+
			"inside the budget", labelResult.CutLength, labelResult.Cut)
	}
	// Too tight for a useful title: the id survives whole, alone.
	if labelResult.Tight != "REQ-042" {
		t.Errorf("a tight budget rendered %q, want the id alone — a half-drawn id is worse than "+
			"no title", labelResult.Tight)
	}
	if labelResult.NoBudget != "REQ-042" || labelResult.NoTitle != "REQ-042" {
		t.Errorf("no budget rendered %q and no title rendered %q, want the id in both cases",
			labelResult.NoBudget, labelResult.NoTitle)
	}
	if labelResult.ExactFit != "REQ-042  abcdef" {
		t.Errorf("a title that exactly fills the budget rendered %q, want it whole and unmarked",
			labelResult.ExactFit)
	}
	if !strings.HasSuffix(labelResult.OneOver, "…") || len([]rune(labelResult.OneOver)) > 20 {
		t.Errorf("a title one character over the budget rendered %q (%d runes); it must be cut, "+
			"marked, and still inside the budget",
			labelResult.OneOver, len([]rune(labelResult.OneOver)))
	}
	// A longer id eats the same budget. The rule holds by id length, not by a
	// hard-coded seven characters.
	if labelResult.LongId != "REQ-100042" {
		t.Errorf("a ten-character id at a 14-character budget rendered %q, want the id alone",
			labelResult.LongId)
	}
}

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
func TestJavaScriptBehaviorTimelineNarrowRowsDrawOneMarker(t *testing.T) {
	indexHtml := generateLiveSite(t)
	// The threshold is READ FROM THE RENDERER. Restating 7 here would let the
	// shipped constant drift to 1 with this test still green — REQ-265's lesson,
	// and REQ-322 shipped exactly that mistake in this file.
	javascriptProbe := timelineProbePreamble(t, "TIMELINE_MIN_SPLIT_WIDTH", "TIMELINE_MIN_SEGMENT_WIDTH") +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineWaitEndMs(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineWorkEndMs(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineRowSegments(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineCollapsedRowMark(") + `
var nowMs = Date.UTC(2026, 7, 23, 12, 0);
// Fit all over three months at the shipped plot width: 90 days across 1200px.
var windowStartMs = Date.UTC(2026, 5, 1);
var windowEndMs = Date.UTC(2026, 7, 30);
var plotWidthPx = 1200;
var windowSpanMs = windowEndMs - windowStartMs;
var msPerPixel = windowSpanMs / plotWidthPx;

function markFor(row, projectedRow, spanStartMs, spanEndMs, widthPx) {
  return timelineCollapsedRowMark(
    timelineRowSegments(row, nowMs, projectedRow),
    spanStartMs === undefined ? windowStartMs : spanStartMs,
    spanEndMs === undefined ? windowEndMs : spanEndMs,
    widthPx === undefined ? plotWidthPx : widthPx);
}

// A REQ whose whole life is four pixels wide at this zoom: the wait and the work
// each floor to TIMELINE_MIN_SEGMENT_WIDTH and together claim more room than the
// REQ occupies.
var fourPixelStartMs = Date.UTC(2026, 6, 1);
var narrowRow = {
  id: "REQ-001",
  createdTime: new Date(fourPixelStartMs).toISOString(),
  claimedTime: new Date(fourPixelStartMs + msPerPixel * 2).toISOString(),
  completedTime: new Date(fourPixelStartMs + msPerPixel * 4).toISOString(),
  hasWork: true, waitMinutes: 60, workMinutes: 60
};
// The same REQ at a zoom where four pixels became four hundred.
var wideWindowStartMs = fourPixelStartMs - msPerPixel * 10;
var wideWindowEndMs = fourPixelStartMs + msPerPixel * 14;

// A row with reversed stamps. Its break markers are the point of drawing it, so
// the caller excludes it — but the mark function must not be the thing that
// silently rescues a caller that forgets.
var reversedRow = {
  id: "REQ-002",
  createdTime: new Date(fourPixelStartMs).toISOString(),
  claimedTime: new Date(fourPixelStartMs - msPerPixel * 2).toISOString(),
  completedTime: new Date(fourPixelStartMs + msPerPixel).toISOString(),
  hasWork: true, waitMinutes: -60, workMinutes: 30
};

// A REQ still waiting: ONE drawn segment, so there is no split to withdraw and
// no marker to collapse to, however narrow it is.
var singleSegmentRow = {
  id: "REQ-003",
  createdTime: new Date(nowMs - msPerPixel).toISOString(),
  waitOpen: true, waitMinutes: 5
};

// The unparseable row. timelineRowSegments hands it the -Infinity → Infinity
// sentinel so it is listed in every window; collapsing that to one marker would
// draw a bar across the entire chart.
var unparseableRow = { id: "REQ-004", createdTime: "not-a-date", waitMinutes: 0 };

// A pending REQ whose forecast bar sits right beside its open wait, both inside
// the same handful of pixels. Collapsing it would draw one SOLID marker over a
// span that is mostly projection — work claimed as measured.
var projectedNarrowRow = {
  id: "REQ-005",
  createdTime: new Date(nowMs - msPerPixel).toISOString(),
  waitOpen: true, waitMinutes: 5
};
var narrowProjection = {
  id: "REQ-005",
  startTime: new Date(nowMs).toISOString(),
  endTime: new Date(nowMs + msPerPixel * 2).toISOString()
};

var narrowMark = markFor(narrowRow);
var wideMark = markFor(narrowRow, undefined, wideWindowStartMs, wideWindowEndMs);

process.stdout.write(JSON.stringify({
  splitWidth: TIMELINE_MIN_SPLIT_WIDTH,
  segmentWidth: TIMELINE_MIN_SEGMENT_WIDTH,
  narrowCollapsed: narrowMark !== null,
  narrowMarkStartIso: narrowMark ? new Date(narrowMark.startMs).toISOString() : "",
  narrowMarkEndIso: narrowMark ? new Date(narrowMark.endMs).toISOString() : "",
  narrowRowCreatedIso: narrowRow.createdTime,
  narrowRowCompletedIso: narrowRow.completedTime,
  wideCollapsed: wideMark !== null,
  reversedCollapsed: markFor(reversedRow) !== null,
  singleSegmentCollapsed: markFor(singleSegmentRow) !== null,
  unparseableCollapsed: markFor(unparseableRow) !== null,
  // A row sitting entirely outside the window has no visible width at all, so it
  // is at the collapsing end of the scale rather than the splitting end.
  offscreenCollapsed: markFor(narrowRow, undefined, Date.UTC(2027, 0, 1), Date.UTC(2027, 1, 1)) !== null,
  // Two segments, four pixels: collapsible by width alone. The caller is what
  // has to spare it.
  projectedNarrowCollapsibleByWidth: markFor(projectedNarrowRow, narrowProjection) !== null
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline narrow rows", javascriptProbe)
	var markResult struct {
		SplitWidth             float64 `json:"splitWidth"`
		SegmentWidth           float64 `json:"segmentWidth"`
		NarrowCollapsed        bool    `json:"narrowCollapsed"`
		NarrowMarkStartIso     string  `json:"narrowMarkStartIso"`
		NarrowMarkEndIso       string  `json:"narrowMarkEndIso"`
		NarrowRowCreatedIso    string  `json:"narrowRowCreatedIso"`
		NarrowRowCompletedIso  string  `json:"narrowRowCompletedIso"`
		WideCollapsed          bool    `json:"wideCollapsed"`
		ReversedCollapsed      bool    `json:"reversedCollapsed"`
		SingleSegmentCollapsed bool    `json:"singleSegmentCollapsed"`
		UnparseableCollapsed   bool    `json:"unparseableCollapsed"`
		OffscreenCollapsed     bool    `json:"offscreenCollapsed"`

		ProjectedNarrowCollapsibleByWidth bool `json:"projectedNarrowCollapsibleByWidth"`
	}
	if decodeError := json.Unmarshal(probeOutput, &markResult); decodeError != nil {
		t.Fatalf("decode timeline narrow row behavior: %v (output %q)", decodeError, probeOutput)
	}

	// THE THRESHOLD, pinned to what a two-segment bar physically needs: two
	// floored segments plus a pixel of boundary between them. Dropping
	// TIMELINE_MIN_SPLIT_WIDTH below that would make the collapse fire only where
	// the split was already invisible, which is the defect this REQ removed.
	if wantFloor := 2*markResult.SegmentWidth + 1; markResult.SplitWidth < wantFloor {
		t.Errorf("the split threshold is %g and a readable two-segment bar needs %g "+
			"(2 x %g floored segments plus a boundary); below it the collapse never fires "+
			"where the split is unreadable", markResult.SplitWidth, wantFloor, markResult.SegmentWidth)
	}

	// The pair that makes the rule about WIDTH and not about the row. Same REQ,
	// same stamps, two zooms.
	if !markResult.NarrowCollapsed {
		t.Error("a REQ whose whole span is four pixels wide kept its wait/work split; " +
			"two floored segments in four pixels is a split the pixels cannot carry")
	}
	if markResult.WideCollapsed {
		t.Error("the same REQ collapsed to one marker at a zoom where its span is hundreds " +
			"of pixels wide; the collapse must cost nothing once there is room to split")
	}
	// The marker has to cover the row's real extent, not a floored stub anchored
	// at one end: the reader still reads its POSITION against the gridlines.
	if markResult.NarrowMarkStartIso != markResult.NarrowRowCreatedIso ||
		markResult.NarrowMarkEndIso != markResult.NarrowRowCompletedIso {
		t.Errorf("the collapsed marker covers %s → %s, want the row's own extent %s → %s",
			markResult.NarrowMarkStartIso, markResult.NarrowMarkEndIso,
			markResult.NarrowRowCreatedIso, markResult.NarrowRowCompletedIso)
	}

	if markResult.SingleSegmentCollapsed {
		t.Error("a REQ drawing one segment was reported as collapsible; there is no split " +
			"to withdraw, and collapsing it would replace its open wait with a closed marker")
	}
	// Same guard as the case above, reached by a different route: the sentinel
	// segment timelineRowSegments emits for an unreadable created_at arrives alone,
	// so "one segment" is what spares it. Collapsing it would draw one marker
	// across the whole chart.
	if markResult.UnparseableCollapsed {
		t.Error("a row with an unreadable created_at was reported as collapsible; its segment " +
			"is the -Infinity sentinel, and collapsing it draws one bar across the whole chart")
	}
	if !markResult.OffscreenCollapsed {
		t.Error("a row with no visible width in the window was reported as splittable; " +
			"zero pixels cannot carry two segments")
	}
	// Not an assertion about the mark function's own judgement — it has no way to
	// know — but a record of what the caller must keep doing. renderVisibleRows
	// excludes broken rows before asking, and this pins the reason: asked
	// directly, the function WOULD collapse one.
	if !markResult.ReversedCollapsed {
		t.Error("a reversed-stamp row stopped being collapsible by width alone; if that is " +
			"now handled here, renderVisibleRows's own broken-stamp guard is dead code")
	}
	// The forecast is the second distinction the collapse would erase, and the
	// same pair applies: collapsible by width alone, spared by the caller.
	if !markResult.ProjectedNarrowCollapsibleByWidth {
		t.Error("a narrow pending REQ with a forecast stopped being collapsible by width " +
			"alone; if that is now handled in the mark function, the caller's projection " +
			"guard is dead code")
	}
	// And the other half of that pair: the guard has to still be at the call site.
	// This is a source check rather than a behavioral one because the collapse
	// decision is the pure function above and the exclusion is the caller's — the
	// failure it names is a broken row quietly drawn as one clean marker, with its
	// break markers, the only reason it is on the chart, gone.
	renderVisibleRowsBody := sliceBalancedBlockAfter(t, indexHtml, "function renderVisibleRows(")
	if !strings.Contains(renderVisibleRowsBody, "rowHasBrokenStamps") {
		t.Error("renderVisibleRows no longer excludes broken-stamp rows before asking whether " +
			"to collapse; a narrow reversed row would draw one clean marker instead of the " +
			"break markers that are the reason it is drawn at all")
	}
	if !strings.Contains(renderVisibleRowsBody, "rowHasBrokenStamps || projectedRow") {
		t.Error("renderVisibleRows no longer excludes forecast-carrying rows before asking " +
			"whether to collapse; a narrow pending REQ would draw one solid marker over a span " +
			"that is mostly projection, claiming work that has not happened")
	}
}

// The threshold decision, without a DOM. The browser probe next door proves the
// behaviour end to end; this pins the boundary itself, which is the part a
// pointer-event probe measures only at whatever offsets it happens to dispatch.
func TestJavaScriptBehaviorTimelinePanThreshold(t *testing.T) {
	indexHtml := generateLiveSite(t)
	// The distance is READ FROM THE RENDERER, so a probe cannot keep passing
	// against a threshold the board stopped using.
	javascriptProbe := timelineProbePreamble(t, "TIMELINE_PAN_THRESHOLD_PX") +
		sliceBalancedBlockAfter(t, indexHtml, "function timelinePanEngaged(") + `
var pressX = 500;
process.stdout.write(JSON.stringify({
  threshold: TIMELINE_PAN_THRESHOLD_PX,
  atRest: timelinePanEngaged(false, pressX, pressX),
  justUnder: timelinePanEngaged(false, pressX, pressX + TIMELINE_PAN_THRESHOLD_PX - 0.01),
  exactlyAt: timelinePanEngaged(false, pressX, pressX + TIMELINE_PAN_THRESHOLD_PX),
  // Leftward drags are drags. An unsigned comparison would engage in one
  // direction only, and the bug would look like "panning left is sticky".
  justUnderLeftward: timelinePanEngaged(false, pressX, pressX - TIMELINE_PAN_THRESHOLD_PX + 0.01),
  exactlyAtLeftward: timelinePanEngaged(false, pressX, pressX - TIMELINE_PAN_THRESHOLD_PX),
  // Latched: once engaged, back at the press point is still engaged.
  latchedAtRest: timelinePanEngaged(true, pressX, pressX)
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline pan threshold", javascriptProbe)
	var thresholdResult struct {
		Threshold         float64 `json:"threshold"`
		AtRest            bool    `json:"atRest"`
		JustUnder         bool    `json:"justUnder"`
		ExactlyAt         bool    `json:"exactlyAt"`
		JustUnderLeftward bool    `json:"justUnderLeftward"`
		ExactlyAtLeftward bool    `json:"exactlyAtLeftward"`
		LatchedAtRest     bool    `json:"latchedAtRest"`
	}
	if decodeError := json.Unmarshal(probeOutput, &thresholdResult); decodeError != nil {
		t.Fatalf("decode timeline pan threshold behavior: %v (output %q)", decodeError, probeOutput)
	}

	// A range, not a value. Below about 2px a trackpad tremor still trips it and
	// the click is lost again; much above 8px a deliberate short drag feels stuck.
	if thresholdResult.Threshold < 2 || thresholdResult.Threshold > 8 {
		t.Errorf("the pan threshold is %g px; under 2 a hand tremor trips it and the click is "+
			"lost again, over 8 a short deliberate drag feels stuck", thresholdResult.Threshold)
	}
	if thresholdResult.AtRest {
		t.Error("a press that has not moved at all engaged the pan")
	}
	if thresholdResult.JustUnder || thresholdResult.JustUnderLeftward {
		t.Errorf("a press just under the %g px threshold engaged the pan (rightward %v, "+
			"leftward %v)", thresholdResult.Threshold,
			thresholdResult.JustUnder, thresholdResult.JustUnderLeftward)
	}
	if !thresholdResult.ExactlyAt || !thresholdResult.ExactlyAtLeftward {
		t.Errorf("a press exactly at the %g px threshold did not engage the pan (rightward %v, "+
			"leftward %v); the comparison has to be inclusive and unsigned in distance",
			thresholdResult.Threshold,
			thresholdResult.ExactlyAt, thresholdResult.ExactlyAtLeftward)
	}
	if !thresholdResult.LatchedAtRest {
		t.Error("an already-engaged drag disengaged on returning to the press point; " +
			"engagement latches, or a wandering drag flickers in and out of panning")
	}
}

// The axis draws ticks and the plot draws gridlines at the same instants. Two
// loops doing the same arithmetic is one edit away from a gridline that means a
// slightly different time than the tick above it, so there is one source and
// this is what keeps it one.
func TestJavaScriptBehaviorTimelineGridlinesShareTheAxisTicks(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := timelineProbePreamble(t, "TIMELINE_AXIS_TICK_COUNT", "TIMELINE_DAY_MS", "TIMELINE_AXIS_TICK_LIMIT") +
		rendererDeclarationLine(t, "web/board-timeline.js", "TIMELINE_WEEK_ALIGNMENT_MS") + "\n" +
		rendererBracketDeclaration(t, "web/board-timeline.js", "TIMELINE_AXIS_TICK_STEPS") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineTickStepSpanMs(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineAxisTickStep(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineTickAtOrBefore(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineSteppedTick(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineAxisTicks(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineAxisTickInstants(") + `
var windowStartMs = Date.UTC(2026, 5, 1);
var windowEndMs = Date.UTC(2026, 5, 8);
var instants = timelineAxisTickInstants(windowStartMs, windowEndMs);
process.stdout.write(JSON.stringify({
  tickCount: TIMELINE_AXIS_TICK_COUNT,
  instantCount: instants.length,
  firstIso: new Date(instants[0]).toISOString(),
  lastIso: new Date(instants[instants.length - 1]).toISOString(),
  ascending: instants.every(function (instant, index) {
    return index === 0 || instant > instants[index - 1];
  }),
  // A zero-width window is reachable: the fields accept one day in both boxes
  // before the settle widens it. It must produce a tick list, not NaNs.
  degenerateFinite: timelineAxisTickInstants(windowStartMs, windowStartMs)
    .every(function (instant) { return isFinite(instant); }),
  degenerateCount: timelineAxisTickInstants(windowStartMs, windowStartMs).length
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline axis ticks", javascriptProbe)
	var tickResult struct {
		TickCount        int    `json:"tickCount"`
		InstantCount     int    `json:"instantCount"`
		FirstIso         string `json:"firstIso"`
		LastIso          string `json:"lastIso"`
		Ascending        bool   `json:"ascending"`
		DegenerateFinite bool   `json:"degenerateFinite"`
		DegenerateCount  int    `json:"degenerateCount"`
	}
	if decodeError := json.Unmarshal(probeOutput, &tickResult); decodeError != nil {
		t.Fatalf("decode timeline axis tick behavior: %v (output %q)", decodeError, probeOutput)
	}

	// REQ-327 CHANGED WHAT THIS COUNTS. The axis no longer divides the window into
	// TIMELINE_AXIS_TICK_COUNT equal parts, so the old "one instant per tick plus
	// both window edges" arithmetic no longer describes it — that arithmetic is the
	// defect: on a seven-day window it put ticks 28 hours apart and the formatter
	// labelled each with a bare date. A June week now gets eight midnights, and the
	// count varies with the rung of the ladder the span picks.
	if tickResult.InstantCount != 8 {
		t.Errorf("a seven-day window produced %d ticks, want the 8 midnights that span it "+
			"(1 June through 8 June inclusive)", tickResult.InstantCount)
	}
	if tickResult.FirstIso != "2026-06-01T00:00:00.000Z" || tickResult.LastIso != "2026-06-08T00:00:00.000Z" {
		t.Errorf("the tick list spans %s → %s, want 2026-06-01 → 2026-06-08",
			tickResult.FirstIso, tickResult.LastIso)
	}
	if tickResult.DegenerateCount != 1 {
		t.Errorf("a zero-width window produced %d ticks, want exactly 1 — six copies of one "+
			"instant is what the old equal-parts loop emitted", tickResult.DegenerateCount)
	}
	if !tickResult.Ascending {
		t.Error("the tick instants are not strictly ascending")
	}
	if !tickResult.DegenerateFinite {
		t.Error("a zero-width window produced non-finite tick instants; the fields can reach " +
			"one before the settle widens it, and NaN ticks draw nothing and log nothing")
	}

	// ONE source, and this is the half that bites. The probe above would pass just
	// as well with renderAxis keeping a private copy of the walk, which is exactly
	// the drift the extraction removed — so count the callers instead of trusting
	// the function's existence.
	//
	// The BOUNDARY WALK specifically: the loop that steps from one tick to the next
	// is what places a tick, and timelineSteppedTick is the only thing entitled to
	// do it. Inlining it back into either caller puts a second copy in the page and
	// fails here. (This replaced a count of "tickIndex) / TIMELINE_AXIS_TICK_COUNT",
	// the equal-parts expression REQ-327 deleted.)
	// Counted as "calls from anywhere else", not as a raw string count: the walk
	// legitimately calls timelineSteppedTick twice inside timelineAxisTicks (once to
	// skip a boundary that precedes the window, once per step), so a bare count of 1
	// would be wrong about the healthy shape and could only be satisfied by making
	// the code worse.
	axisTicksBody := sliceBalancedBlockAfter(t, indexHtml, "function timelineAxisTicks(")
	steppedTickDefinition := sliceBalancedBlockAfter(t, indexHtml, "function timelineSteppedTick(")
	callsEverywhere := strings.Count(indexHtml, "timelineSteppedTick(")
	callsInTheWalk := strings.Count(axisTicksBody, "timelineSteppedTick(")
	callsInItsOwnDefinition := strings.Count(steppedTickDefinition, "timelineSteppedTick(")
	if callsElsewhere := callsEverywhere - callsInTheWalk - callsInItsOwnDefinition; callsElsewhere != 0 {
		t.Errorf("timelineSteppedTick is called from %d place(s) outside timelineAxisTicks; the "+
			"boundary walk lives in one function or the gridlines can start meaning a different "+
			"instant than the ticks above them", callsElsewhere)
	}
	if callsInTheWalk == 0 {
		t.Error("timelineAxisTicks does not call timelineSteppedTick, so the check above is " +
			"comparing two numbers that no longer describe the walk")
	}
	// renderAxis needs the GAP as well as the instants, so it calls timelineAxisTicks
	// directly; drawGridlines needs only the instants. Both bottom out in one walk,
	// and each is named here so neither can quietly grow its own.
	for caller, wantCall := range map[string]string{
		"function renderAxis(":    "timelineAxisTicks(",
		"function drawGridlines(": "timelineAxisTickInstants(",
	} {
		callerBody := sliceBalancedBlockAfter(t, indexHtml, caller)
		if !strings.Contains(callerBody, wantCall) {
			t.Errorf("%s does not read %s; the axis and the gridlines have to come from one "+
				"list or they can disagree", caller, wantCall)
		}
	}
	// And the gap the labels are formatted against comes from the same call that
	// positioned the ticks, rather than being derived a second time.
	renderAxisBody := sliceBalancedBlockAfter(t, indexHtml, "function renderAxis(")
	if !strings.Contains(renderAxisBody, "axisTicks.gapMs") {
		t.Error("renderAxis does not pass the chosen gap to timelineFormatAxisTick; a gap " +
			"derived a second time is how a date-only label came to sit on a 04:00 tick")
	}
}

func TestJavaScriptBehaviorTimelineForecastStatesItsAssumptions(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := timelineForecastDomStub +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineFormatSpanMinutes(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineFormatStamp(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function clearTimelineForecast(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function renderTimelineForecast(") + `
var confidentProjection = {
  confident: true,
  chainStart: "2026-06-20T12:00:00Z",
  queueEnd: "2026-06-20T14:30:00Z",
  windowSamples: 60, windowSize: 60, minimumSamples: 5,
  normalSamples: 55, normalMinutes: 40,
  trivialSamples: 5, trivialMinutes: 10,
  rows: [{ id: "REQ-401" }, { id: "REQ-402" }],
  excluded: [{ id: "REQ-404", reason: "waiting on an external condition" }]
};
renderTimelineForecast(confidentProjection, false);
var confidentText = collectText(stubNodes["timeline-forecast"]);
var confidentExcludedText = collectText(stubNodes["timeline-excluded"]);

stubNodes["timeline-forecast"] = makeStubNode();
stubNodes["timeline-excluded"] = makeStubNode();
renderTimelineForecast({
  confident: false,
  declinedReason: "only 2 completed REQs inside the read-time rule; 5 are needed before a median means anything",
  rows: [], excluded: [], windowSamples: 2, minimumSamples: 5
}, false);
var declinedText = collectText(stubNodes["timeline-forecast"]);

// The same projection with filters ON. The rows are a subset; the forecast is
// not, and the copy has to say so.
stubNodes["timeline-forecast"] = makeStubNode();
stubNodes["timeline-excluded"] = makeStubNode();
renderTimelineForecast(confidentProjection, true);
var filteredText = collectText(stubNodes["timeline-forecast"]);
var filteredExcludedText = collectText(stubNodes["timeline-excluded"]);

// Declining with filters on carries the same label: the history it declined on
// is the whole queue's, not the subset's.
stubNodes["timeline-forecast"] = makeStubNode();
stubNodes["timeline-excluded"] = makeStubNode();
renderTimelineForecast({
  confident: false,
  declinedReason: "only 2 completed REQs inside the read-time rule; 5 are needed before a median means anything",
  rows: [], excluded: [], windowSamples: 2, minimumSamples: 5
}, true);
var filteredDeclinedText = collectText(stubNodes["timeline-forecast"]);

// The no-rows path clears both nodes without rendering anything: a forecast left
// standing beside "no REQ matches" describes rows that are not on screen.
clearTimelineForecast();
var clearedText = collectText(stubNodes["timeline-forecast"]);
var clearedExcludedText = collectText(stubNodes["timeline-excluded"]);

process.stdout.write(JSON.stringify({
  confidentText: confidentText,
  confidentExcludedText: confidentExcludedText,
  declinedText: declinedText,
  filteredText: filteredText,
  filteredExcludedText: filteredExcludedText,
  filteredDeclinedText: filteredDeclinedText,
  clearedText: clearedText,
  clearedExcludedText: clearedExcludedText
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline forecast", javascriptProbe)
	var forecastResult struct {
		ConfidentText         string `json:"confidentText"`
		ConfidentExcludedText string `json:"confidentExcludedText"`
		DeclinedText          string `json:"declinedText"`
		FilteredText          string `json:"filteredText"`
		FilteredExcludedText  string `json:"filteredExcludedText"`
		FilteredDeclinedText  string `json:"filteredDeclinedText"`
		ClearedText           string `json:"clearedText"`
		ClearedExcludedText   string `json:"clearedExcludedText"`
	}
	if decodeError := json.Unmarshal(probeOutput, &forecastResult); decodeError != nil {
		t.Fatalf("decode timeline forecast behavior: %v (output %q)", decodeError, probeOutput)
	}

	// Each of these is a separate requirement, so each is asserted separately —
	// a single "contains everything" check would report one failure for any of
	// four different regressions.
	for _, wantFragment := range []struct {
		requirement string
		fragment    string
	}{
		{"the end instant itself", "2026-06-20 14:30 UTC"},
		{"the window's sample size", "last 60 completed REQs"},
		{"each bucket's sample count and median", "55 substantive at 40 min"},
		{"the serial assumption", "one REQ at a time"},
		{"the no-parallelism assumption", "no parallel builders"},
		{"the static-queue assumption", "queue that stops growing"},
		{"the read-time rule's exclusions", "Paused and reversed spans are excluded"},
	} {
		if !strings.Contains(forecastResult.ConfidentText, wantFragment.fragment) {
			t.Fatalf("the forecast sentence does not state %s (wanted %q in %q)",
				wantFragment.requirement, wantFragment.fragment, forecastResult.ConfidentText)
		}
	}
	if !strings.Contains(forecastResult.ConfidentExcludedText, "REQ-404") ||
		!strings.Contains(forecastResult.ConfidentExcludedText, "waiting on an external condition") {
		t.Fatalf("the excluded list must name every unschedulable REQ and its reason; got %q",
			forecastResult.ConfidentExcludedText)
	}

	if strings.TrimSpace(forecastResult.ClearedText) != "" ||
		strings.TrimSpace(forecastResult.ClearedExcludedText) != "" {
		t.Fatalf("clearing left forecast %q and excluded %q; a filter matching no rows must leave neither standing beside \"no REQ matches\"",
			forecastResult.ClearedText, forecastResult.ClearedExcludedText)
	}

	if strings.Contains(forecastResult.DeclinedText, "Queue empties") {
		t.Fatalf("thin history produced an end date: %q", forecastResult.DeclinedText)
	}
	if !strings.Contains(forecastResult.DeclinedText, "No end estimate") ||
		!strings.Contains(forecastResult.DeclinedText, "5 are needed") {
		t.Fatalf("declining must say so and carry the reason; got %q", forecastResult.DeclinedText)
	}

	// REQ-305: rows are filtered, the projection never is. With a subset on
	// screen the forecast schedules the whole queue and the excluded list names
	// IDs no visible row carries, so the copy has to name its own population.
	// The label has to read correctly alone, because this paragraph is the one
	// people screenshot.
	if !strings.Contains(forecastResult.FilteredText, "whole queue") {
		t.Errorf("with filters on, the forecast must say it covers the whole queue rather than the rows shown; got %q",
			forecastResult.FilteredText)
	}
	if !strings.Contains(forecastResult.FilteredExcludedText, "whole queue") {
		t.Errorf("with filters on, the excluded list must say it lists the whole queue's exclusions — it names IDs no visible row carries; got %q",
			forecastResult.FilteredExcludedText)
	}
	if !strings.Contains(forecastResult.FilteredDeclinedText, "whole queue") {
		t.Errorf("a declined forecast declined on the whole queue's history, and must say so under a filter too; got %q",
			forecastResult.FilteredDeclinedText)
	}
	// The label is added, never substituted: everything the unfiltered sentence
	// promised is still in the filtered one.
	if !strings.Contains(forecastResult.FilteredText, "Queue empties around") ||
		!strings.Contains(forecastResult.FilteredText, "no parallel builders") {
		t.Errorf("the filtered forecast must still carry the estimate and its assumptions; got %q",
			forecastResult.FilteredText)
	}
	// And the unfiltered copy is untouched — the settled case stays settled.
	if strings.Contains(forecastResult.ConfidentText, "whole queue") ||
		strings.Contains(forecastResult.ConfidentExcludedText, "whole queue") {
		t.Errorf("with no filter active there is nothing to disambiguate, so the label must not appear; got %q / %q",
			forecastResult.ConfidentText, forecastResult.ConfidentExcludedText)
	}
}

// Rows advertise role="button" and take focus, but a <g> is not a native button:
// Enter and Space never synthesize the click the drawer listens for. The role is
// a promise, and this is the code that keeps it.
func TestJavaScriptBehaviorTimelineRowsActivateFromTheKeyboard(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := sliceBalancedBlockAfter(t, indexHtml, "function timelineKeyboardActivationTarget(") + `
function rowEvent(key, detailId) {
  var trigger = detailId === null ? null : {
    getAttribute: function (name) {
      return name === "data-detail-kind" ? "request" : detailId;
    }
  };
  return { key: key, target: { closest: function () { return trigger; } } };
}
process.stdout.write(JSON.stringify({
  enter: timelineKeyboardActivationTarget(rowEvent("Enter", "REQ-401")),
  space: timelineKeyboardActivationTarget(rowEvent(" ", "REQ-402")),
  legacySpace: timelineKeyboardActivationTarget(rowEvent("Spacebar", "REQ-403")),
  tab: timelineKeyboardActivationTarget(rowEvent("Tab", "REQ-404")),
  arrow: timelineKeyboardActivationTarget(rowEvent("ArrowDown", "REQ-405")),
  offRow: timelineKeyboardActivationTarget(rowEvent("Enter", null))
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline keyboard activation", javascriptProbe)
	var activationResult struct {
		Enter       *struct{ DetailKind, DetailId string } `json:"enter"`
		Space       *struct{ DetailKind, DetailId string } `json:"space"`
		LegacySpace *struct{ DetailKind, DetailId string } `json:"legacySpace"`
		Tab         *struct{ DetailKind, DetailId string } `json:"tab"`
		Arrow       *struct{ DetailKind, DetailId string } `json:"arrow"`
		OffRow      *struct{ DetailKind, DetailId string } `json:"offRow"`
	}
	if decodeError := json.Unmarshal(probeOutput, &activationResult); decodeError != nil {
		t.Fatalf("decode timeline keyboard activation: %v (output %q)", decodeError, probeOutput)
	}
	for _, activated := range []struct {
		keyName string
		result  *struct{ DetailKind, DetailId string }
		wantId  string
	}{
		{"Enter", activationResult.Enter, "REQ-401"},
		{"Space", activationResult.Space, "REQ-402"},
		{"Spacebar (legacy)", activationResult.LegacySpace, "REQ-403"},
	} {
		if activated.result == nil {
			t.Fatalf("%s on a focused row activated nothing; the row advertises role=button", activated.keyName)
		}
		if activated.result.DetailId != activated.wantId || activated.result.DetailKind != "request" {
			t.Fatalf("%s activated %+v, want request/%s", activated.keyName, *activated.result, activated.wantId)
		}
	}
	// Navigation keys and keys pressed off a row must not open anything.
	for _, ignored := range []struct {
		keyName string
		result  *struct{ DetailKind, DetailId string }
	}{
		{"Tab", activationResult.Tab},
		{"ArrowDown", activationResult.Arrow},
		{"Enter off a row", activationResult.OffRow},
	} {
		if ignored.result != nil {
			t.Fatalf("%s activated %+v; it must open nothing", ignored.keyName, *ignored.result)
		}
	}
}

// Zoom and pan now have two drivers — a pointer and a keyboard — and two drivers
// that each compute a window are two definitions of where the window goes. The
// keyboard path is written as pure transforms over the SAME timelineZoomedWindow
// the wheel and the zoom buttons call; this drives both to their edges and
// requires them to arrive at the same ones.
func TestJavaScriptBehaviorTimelineKeyboardMovesTheSameWindowAsThePointer(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := timelineProbePreamble(t, "TIMELINE_MIN_SPAN_MS", "TIMELINE_ZOOM_STEP", "TIMELINE_PAN_FRACTION") +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineZoomedWindow(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelinePannedWindow(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineKeyboardWindow(") + `
var boundStart = 0;
var boundEnd = 30 * 24 * 3600 * 1000;   // a 30-day board

// Start half zoomed in, so a pan has room to run in both directions before it
// meets a bound.
var halfway = timelineZoomedWindow(boundStart, boundEnd, 2, 0.5, boundStart, boundEnd);
var halfSpanMs = halfway.windowEndMs - halfway.windowStartMs;

function pressKey(currentWindow, keyName) {
  var moved = timelineKeyboardWindow(
    keyName, currentWindow.windowStartMs, currentWindow.windowEndMs, boundStart, boundEnd);
  return moved || currentWindow;
}
function repeatKey(currentWindow, keyName, pressCount) {
  for (var press = 0; press < pressCount; press++) {
    currentWindow = pressKey(currentWindow, keyName);
  }
  return currentWindow;
}

var pannedRight = pressKey(halfway, "ArrowRight");
var pannedLeft = pressKey(halfway, "ArrowLeft");

// Held to the edges: the window must stop AT the bound, keeping its span.
var atRightEdge = repeatKey(halfway, "ArrowRight", 40);
var atLeftEdge = repeatKey(halfway, "ArrowLeft", 40);

// The pointer path's own floor and ceiling, reached through the wheel's
// off-centre anchor rather than the keyboard's centred one.
var pointerFloor = halfway;
var pointerCeiling = halfway;
for (var pointerStep = 0; pointerStep < 40; pointerStep++) {
  pointerFloor = timelineZoomedWindow(
    pointerFloor.windowStartMs, pointerFloor.windowEndMs, TIMELINE_ZOOM_STEP, 0.25, boundStart, boundEnd);
  pointerCeiling = timelineZoomedWindow(
    pointerCeiling.windowStartMs, pointerCeiling.windowEndMs, 1 / TIMELINE_ZOOM_STEP, 0.25, boundStart, boundEnd);
}
var keyboardFloor = repeatKey(halfway, "+", 40);
var keyboardCeiling = repeatKey(halfway, "-", 40);

process.stdout.write(JSON.stringify({
  panStepMs: pannedRight.windowStartMs - halfway.windowStartMs,
  panBackStepMs: halfway.windowStartMs - pannedLeft.windowStartMs,
  wantPanStepMs: halfSpanMs * TIMELINE_PAN_FRACTION,
  windowSpanMs: halfSpanMs,
  panKeepsSpan:
    pannedRight.windowEndMs - pannedRight.windowStartMs === halfSpanMs &&
    pannedLeft.windowEndMs - pannedLeft.windowStartMs === halfSpanMs,
  rightEdgeMs: atRightEdge.windowEndMs,
  leftEdgeMs: atLeftEdge.windowStartMs,
  boundStartMs: boundStart,
  boundEndMs: boundEnd,
  edgesKeepSpan:
    atRightEdge.windowEndMs - atRightEdge.windowStartMs === halfSpanMs &&
    atLeftEdge.windowEndMs - atLeftEdge.windowStartMs === halfSpanMs,
  keyboardFloorSpanMs: keyboardFloor.windowEndMs - keyboardFloor.windowStartMs,
  pointerFloorSpanMs: pointerFloor.windowEndMs - pointerFloor.windowStartMs,
  minSpanMs: TIMELINE_MIN_SPAN_MS,
  keyboardCeilingSpanMs: keyboardCeiling.windowEndMs - keyboardCeiling.windowStartMs,
  pointerCeilingSpanMs: pointerCeiling.windowEndMs - pointerCeiling.windowStartMs,
  boundSpanMs: boundEnd - boundStart,
  unownedKeys: ["Enter", " ", "Spacebar", "Tab", "ArrowUp", "ArrowDown", "a"].map(function (keyName) {
    return timelineKeyboardWindow(keyName, halfway.windowStartMs, halfway.windowEndMs, boundStart, boundEnd);
  })
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline keyboard pan and zoom", javascriptProbe)
	var keyboardResult struct {
		PanStepMs             float64            `json:"panStepMs"`
		PanBackStepMs         float64            `json:"panBackStepMs"`
		WantPanStepMs         float64            `json:"wantPanStepMs"`
		WindowSpanMs          float64            `json:"windowSpanMs"`
		PanKeepsSpan          bool               `json:"panKeepsSpan"`
		RightEdgeMs           float64            `json:"rightEdgeMs"`
		LeftEdgeMs            float64            `json:"leftEdgeMs"`
		BoundStartMs          float64            `json:"boundStartMs"`
		BoundEndMs            float64            `json:"boundEndMs"`
		EdgesKeepSpan         bool               `json:"edgesKeepSpan"`
		KeyboardFloorSpanMs   float64            `json:"keyboardFloorSpanMs"`
		PointerFloorSpanMs    float64            `json:"pointerFloorSpanMs"`
		MinSpanMs             float64            `json:"minSpanMs"`
		KeyboardCeilingSpanMs float64            `json:"keyboardCeilingSpanMs"`
		PointerCeilingSpanMs  float64            `json:"pointerCeilingSpanMs"`
		BoundSpanMs           float64            `json:"boundSpanMs"`
		UnownedKeys           []*json.RawMessage `json:"unownedKeys"`
	}
	if decodeError := json.Unmarshal(probeOutput, &keyboardResult); decodeError != nil {
		t.Fatalf("decode timeline keyboard behavior: %v (output %q)", decodeError, probeOutput)
	}

	// A pan step has to be a fraction of what is on screen: a fixed number of
	// milliseconds is either imperceptible zoomed out or a jump zoomed in.
	if math.Abs(keyboardResult.PanStepMs-keyboardResult.WantPanStepMs) > 1 {
		t.Fatalf("ArrowRight moved the window %.0f ms, want %.0f ms — one step of the visible span",
			keyboardResult.PanStepMs, keyboardResult.WantPanStepMs)
	}
	if math.Abs(keyboardResult.PanBackStepMs-keyboardResult.WantPanStepMs) > 1 {
		t.Fatalf("ArrowLeft moved the window %.0f ms, want %.0f ms back",
			keyboardResult.PanBackStepMs, keyboardResult.WantPanStepMs)
	}
	if keyboardResult.PanStepMs <= 0 || keyboardResult.PanStepMs >= keyboardResult.WindowSpanMs {
		t.Fatalf("a pan step of %.0f ms against a %.0f ms window is not a bounded step; a reader loses their place",
			keyboardResult.PanStepMs, keyboardResult.WindowSpanMs)
	}
	if !keyboardResult.PanKeepsSpan {
		t.Fatal("panning changed the window span; panning moves the window, zooming resizes it")
	}

	// Held down, a pan must stop at the range edge rather than walking the window
	// off the data — the same clamp the drag path applies.
	if math.Abs(keyboardResult.RightEdgeMs-keyboardResult.BoundEndMs) > 1 {
		t.Fatalf("panning right settled with the window ending at %.0f ms, want the range edge %.0f ms",
			keyboardResult.RightEdgeMs, keyboardResult.BoundEndMs)
	}
	if math.Abs(keyboardResult.LeftEdgeMs-keyboardResult.BoundStartMs) > 1 {
		t.Fatalf("panning left settled with the window starting at %.0f ms, want the range edge %.0f ms",
			keyboardResult.LeftEdgeMs, keyboardResult.BoundStartMs)
	}
	if !keyboardResult.EdgesKeepSpan {
		t.Fatal("clamping at a range edge changed the window span; it must only stop the window, not shrink it")
	}

	// The point of routing the keys through timelineZoomedWindow: one floor and
	// one ceiling, whichever driver arrives at them.
	if keyboardResult.KeyboardFloorSpanMs != keyboardResult.PointerFloorSpanMs {
		t.Fatalf("`+` bottomed out at %.0f ms but the pointer path bottoms out at %.0f ms; the two have diverged",
			keyboardResult.KeyboardFloorSpanMs, keyboardResult.PointerFloorSpanMs)
	}
	if keyboardResult.KeyboardFloorSpanMs != keyboardResult.MinSpanMs {
		t.Fatalf("`+` bottomed out at %.0f ms, want the renderer's %.0f ms floor",
			keyboardResult.KeyboardFloorSpanMs, keyboardResult.MinSpanMs)
	}
	if keyboardResult.KeyboardCeilingSpanMs != keyboardResult.PointerCeilingSpanMs {
		t.Fatalf("`-` topped out at %.0f ms but the pointer path tops out at %.0f ms; the two have diverged",
			keyboardResult.KeyboardCeilingSpanMs, keyboardResult.PointerCeilingSpanMs)
	}
	if keyboardResult.KeyboardCeilingSpanMs != keyboardResult.BoundSpanMs {
		t.Fatalf("`-` topped out at %.0f ms, want the full range span %.0f ms",
			keyboardResult.KeyboardCeilingSpanMs, keyboardResult.BoundSpanMs)
	}

	// Enter and Space belong to row activation, and Up/Down to scrolling the
	// queue. Claiming any of them would take a working interaction away.
	unownedKeyNames := []string{"Enter", "Space", "Spacebar", "Tab", "ArrowUp", "ArrowDown", "a"}
	if len(keyboardResult.UnownedKeys) != len(unownedKeyNames) {
		t.Fatalf("probe reported %d unowned-key results, want %d", len(keyboardResult.UnownedKeys), len(unownedKeyNames))
	}
	for keyIndex, unownedKeyName := range unownedKeyNames {
		if keyboardResult.UnownedKeys[keyIndex] != nil {
			t.Fatalf("%s moved the time window to %s; that key belongs to row activation or to scrolling",
				unownedKeyName, string(*keyboardResult.UnownedKeys[keyIndex]))
		}
	}
}

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
func TestJavaScriptBehaviorUserRequestsOnlyLensFoldsCardsUntilARowIsOpened(t *testing.T) {
	indexHtml := generateLiveSite(t)
	functionBlocks := []string{
		sliceBalancedBlockAfter(t, indexHtml, "function createElement("),
		sliceBalancedBlockAfter(t, indexHtml, "function isTerminalResolvedStatus("),
		sliceBalancedBlockAfter(t, indexHtml, "function hasActiveFilters("),
		sliceBalancedBlockAfter(t, indexHtml, "function searchMatchesRequest("),
		sliceBalancedBlockAfter(t, indexHtml, "function searchMatchesUserRequest("),
		sliceBalancedBlockAfter(t, indexHtml, "function requestMatchesFilters("),
		sliceBalancedBlockAfter(t, indexHtml, "function userRequestHasOpenOrRecentWork("),
		sliceBalancedBlockAfter(t, indexHtml, "function recentWindowPhrase("),
		sliceBalancedBlockAfter(t, indexHtml, "function userRequestLensEmptyText("),
		sliceBalancedBlockAfter(t, indexHtml, "function recentlyDoneIds("),
		sliceBalancedBlockAfter(t, indexHtml, "function renderUserRequestLens("),
	}
	javascriptProbe := `
Date.now = function () { return Date.parse("2026-08-15T12:00:00Z"); };
var boardData = {
  requests: {
    "REQ-601": { status: "pending", title: "alpha open", domain: "general" },
    "REQ-602": { status: "completed", title: "alpha shipped", domain: "general" },
    "REQ-603": { status: "claimed", title: "beta running", domain: "general" },
    "REQ-604": { status: "completed", title: "gamma archived", domain: "general" }
  },
  userRequests: {
    "UR-401": { requestIds: ["REQ-601", "REQ-602"], title: "alpha request", inputFilePresent: true },
    "UR-402": { requestIds: ["REQ-603"], title: "beta request", inputFilePresent: true },
    "UR-403": { requestIds: ["REQ-604"], title: "gamma request", inputFilePresent: false }
  },
  userRequestOrder: ["UR-401", "UR-402", "UR-403"],
  calendar: [
    { id: "REQ-602", completionTime: "2026-08-15T06:00:00Z" },
    { id: "REQ-604", completionTime: "2026-08-01T06:00:00Z" }
  ]
};
var requestsById = boardData.requests;
var userRequestsById = boardData.userRequests;
var viewState = { view: "board", lens: "user-request", windowHours: 24 };
var filterState = { searchText: "", domain: "", status: "", userRequestActivity: "all" };
var userRequestCardsFolded = true;

function makeNode() {
  return {
    childNodes: [],
    dataset: {},
    attributes: {},
    listeners: {},
    appendChild: function (childNode) { this.childNodes.push(childNode); return childNode; },
    removeChild: function (childNode) {
      var childIndex = this.childNodes.indexOf(childNode);
      if (childIndex !== -1) { this.childNodes.splice(childIndex, 1); }
      return childNode;
    },
    setAttribute: function (attributeName, attributeValue) { this.attributes[attributeName] = attributeValue; },
    getAttribute: function (attributeName) {
      return Object.prototype.hasOwnProperty.call(this.attributes, attributeName)
        ? this.attributes[attributeName]
        : null;
    },
    addEventListener: function (eventName, handler) {
      this.listeners[eventName] = (this.listeners[eventName] || []).concat([handler]);
    },
    dispatch: function (eventName) {
      (this.listeners[eventName] || []).forEach(function (handler) { handler(); });
    }
  };
}
var userRequestLensNode = makeNode();
var document = {
  getElementById: function (nodeId) { return nodeId === "user-request-lens" ? userRequestLensNode : null; },
  createElement: function () { return makeNode(); }
};
function makeRequestCard(requestId) { return { className: "req-card", requestId: requestId }; }
` + strings.Join(functionBlocks, "\n") + `
function collectByClassName(node, wantedClassName, found) {
  found = found || [];
  if (node.className === wantedClassName) { found.push(node); }
  (node.childNodes || []).forEach(function (child) { collectByClassName(child, wantedClassName, found); });
  return found;
}
function collectDrawerTriggers(node, found) {
  found = found || [];
  if (node.dataset && node.dataset.detailKind === "ur") { found.push(node.dataset.detailId); }
  (node.childNodes || []).forEach(function (child) { collectDrawerTriggers(child, found); });
  return found;
}
function renderGroups() {
  userRequestLensNode = makeNode();
  renderUserRequestLens();
  return userRequestLensNode.childNodes.filter(function (node) { return node.className === "ur-group"; });
}
function headOf(group) { return collectByClassName(group, "ur-group-head")[0]; }
function describeGroups(groups) {
  return groups.map(function (group) {
    var cardIds = [];
    collectByClassName(group, "ur-group-cards").forEach(function (cardsNode) {
      cardsNode.childNodes.forEach(function (card) { cardIds.push(card.requestId); });
    });
    return {
      userRequestId: collectByClassName(group, "ur-id")[0].textContent,
      expanded: headOf(group).getAttribute("aria-expanded") || "",
      cardIds: cardIds,
      drawerTriggers: collectDrawerTriggers(group)
    };
  });
}

var foldedGroups = renderGroups();
var foldedInitial = describeGroups(foldedGroups);
headOf(foldedGroups[0]).dispatch("click");
var afterOpen = describeGroups(foldedGroups);
headOf(foldedGroups[0]).dispatch("click");
var afterClose = describeGroups(foldedGroups);

filterState.userRequestActivity = "active";
filterState.status = "pending";
var scopedFoldedGroups = renderGroups();
headOf(scopedFoldedGroups[0]).dispatch("click");
var scopedFolded = describeGroups(scopedFoldedGroups);
userRequestCardsFolded = false;
var scopedByUserRequest = describeGroups(renderGroups());

process.stdout.write(JSON.stringify({
  foldedInitial: foldedInitial,
  afterOpen: afterOpen,
  afterClose: afterClose,
  scopedFolded: scopedFolded,
  scopedByUserRequest: scopedByUserRequest
}));
`
	probeOutput := runJavaScriptBehaviorProbe(t, "URs-only fold", javascriptProbe)

	var result struct {
		FoldedInitial       []renderedUserRequestRow `json:"foldedInitial"`
		AfterOpen           []renderedUserRequestRow `json:"afterOpen"`
		AfterClose          []renderedUserRequestRow `json:"afterClose"`
		ScopedFolded        []renderedUserRequestRow `json:"scopedFolded"`
		ScopedByUserRequest []renderedUserRequestRow `json:"scopedByUserRequest"`
	}
	if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
		t.Fatalf("decode URs-only fold output: %v (output %q)", decodeError, probeOutput)
	}

	// Folded: one row per UR, no cards, and the drawer still reachable from each row.
	wantUserRequestIds := []string{"UR-401", "UR-402", "UR-403"}
	if len(result.FoldedInitial) != len(wantUserRequestIds) {
		t.Fatalf("URs only rendered %d rows, want %d: %#v", len(result.FoldedInitial), len(wantUserRequestIds), result.FoldedInitial)
	}
	for rowIndex, row := range result.FoldedInitial {
		if row.UserRequestId != wantUserRequestIds[rowIndex] {
			t.Fatalf("URs only row %d = %q, want %q", rowIndex, row.UserRequestId, wantUserRequestIds[rowIndex])
		}
		if len(row.CardIds) != 0 {
			t.Fatalf("URs only row %s rendered cards %#v before it was opened, want none", row.UserRequestId, row.CardIds)
		}
		if row.Expanded != "false" {
			t.Fatalf("URs only row %s aria-expanded = %q, want \"false\"", row.UserRequestId, row.Expanded)
		}
		if len(row.DrawerTriggers) != 1 || row.DrawerTriggers[0] != row.UserRequestId {
			t.Fatalf("URs only row %s drawer triggers = %#v, want exactly one for itself", row.UserRequestId, row.DrawerTriggers)
		}
	}

	// Opening one row reveals exactly that UR's cards and leaves the others folded.
	openedRow := result.AfterOpen[0]
	if openedRow.Expanded != "true" {
		t.Fatalf("opened row %s aria-expanded = %q, want \"true\"", openedRow.UserRequestId, openedRow.Expanded)
	}
	if strings.Join(openedRow.CardIds, ",") != "REQ-601,REQ-602" {
		t.Fatalf("opened row %s cards = %#v, want REQ-601 and REQ-602", openedRow.UserRequestId, openedRow.CardIds)
	}
	for _, stillFolded := range result.AfterOpen[1:] {
		if len(stillFolded.CardIds) != 0 {
			t.Fatalf("row %s unfolded with a sibling: cards %#v", stillFolded.UserRequestId, stillFolded.CardIds)
		}
	}

	// Activating it again folds it back.
	closedRow := result.AfterClose[0]
	if closedRow.Expanded != "false" || len(closedRow.CardIds) != 0 {
		t.Fatalf("re-activated row = %#v, want aria-expanded false and no cards", closedRow)
	}

	// Active scope plus a status filter must decide identically in both lenses.
	foldedScopedIds := userRequestIdsOf(result.ScopedFolded)
	byUserRequestScopedIds := userRequestIdsOf(result.ScopedByUserRequest)
	if strings.Join(foldedScopedIds, ",") != strings.Join(byUserRequestScopedIds, ",") {
		t.Fatalf("URs only showed %#v under Active+status:pending, by-UR showed %#v; the two lenses must hide the same URs",
			foldedScopedIds, byUserRequestScopedIds)
	}
	if strings.Join(foldedScopedIds, ",") != "UR-401" {
		t.Fatalf("Active+status:pending showed %#v, want only UR-401", foldedScopedIds)
	}
	if strings.Join(result.ScopedFolded[0].CardIds, ",") != "REQ-601" {
		t.Fatalf("opened row under a status filter showed %#v, want only the matching REQ-601", result.ScopedFolded[0].CardIds)
	}
	// The fold is the folded lens's alone: a by-UR head announces no expanded state.
	if result.ScopedByUserRequest[0].Expanded != "" {
		t.Fatalf("by-UR head carries aria-expanded=%q; the fold must not leak into the always-open lens",
			result.ScopedByUserRequest[0].Expanded)
	}
}

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

// Period navigation is the THIRD way to move the timeline's window, after the
// pointer and the keyboard, and what it can get wrong is a calendar that is only
// nearly a calendar: a "week" starting at whatever instant sat in the middle of
// the view, a "next" that adds 7×24h to an unaligned start, or a step that walks
// off the end of the range. Those are invisible in a screenshot — the bars still
// look like bars. This drives the pure transforms over a five-month fixture and
// requires real Monday-to-Monday boundaries, a clamp at the end of the range, and
// a level that stops claiming to be exact once a free zoom has moved the window.
//
// It also pins the other half of the Now button: recentring the time window never
// moved the ROW list, so "jump to the remaining work" landed the reader on
// whichever archived rows happened to be scrolled into view.
func TestJavaScriptBehaviorTimelinePeriodStepsOnCalendarBoundariesAndJumpsToNow(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := timelineProbePreamble(t, "TIMELINE_MIN_SPAN_MS", "TIMELINE_DAY_MS", "TIMELINE_ROW_HEIGHT",
		"TIMELINE_NOW_JUMP_MARGIN_FRACTION", "TIMELINE_NOW_JUMP_MINIMUM_SPAN_MS") +
		rendererDeclarationLine(t, "web/board-timeline.js", "TIMELINE_PERIOD_LEVEL_NAMES") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineZoomedWindow(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelinePeriodStart(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineSteppedPeriodStart(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelinePeriodWindow(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelinePeriodAnchor(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelinePeriodLevelOfWindow(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineFirstOpenRowIndex(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineNowJump(") + `
// A board whose range spans five months — the shape the user reported, where the
// only sideways movement was a very long drag.
var boundStart = Date.UTC(2026, 3, 7);        // 7 Apr 2026
var boundEnd = Date.UTC(2026, 8, 2);          // 2 Sep 2026
var nowMs = Date.UTC(2026, 7, 18, 10, 30);
var queueEndMs = Date.UTC(2026, 7, 30, 6, 0); // the projection's queue-empty instant

var fitted = { windowStartMs: boundStart, windowEndMs: boundEnd };
// The PRODUCTION anchor, not a local copy of it. A probe that reimplements the
// decision under test cannot hold its call site (REQ-305), and the local midpoint
// copy this replaced is exactly why the chips could land months from the now-line
// with the whole suite green.
function anchorOf(movedWindow, stepCount) {
  return timelinePeriodAnchor(stepCount, movedWindow.windowStartMs, movedWindow.windowEndMs, nowMs);
}

var weekWindow = timelinePeriodWindow(anchorOf(fitted, 0), "week", 0, boundStart, boundEnd);
var nextWeek = timelinePeriodWindow(anchorOf(weekWindow, 1), "week", 1, boundStart, boundEnd);
var prevWeek = timelinePeriodWindow(anchorOf(nextWeek, -1), "week", -1, boundStart, boundEnd);

// Held down against the end of the range: stepping has to stop, keeping the
// window on the data instead of running off it.
var atRangeEnd = weekWindow;
for (var step = 0; step < 60; step++) {
  atRangeEnd = timelinePeriodWindow(anchorOf(atRangeEnd, 1), "week", 1, boundStart, boundEnd);
}
var pastRangeEnd = timelinePeriodWindow(anchorOf(atRangeEnd, 1), "week", 1, boundStart, boundEnd);

var dayWindow = timelinePeriodWindow(anchorOf(fitted, 0), "day", 0, boundStart, boundEnd);
var monthWindow = timelinePeriodWindow(anchorOf(fitted, 0), "month", 0, boundStart, boundEnd);
var nextMonth = timelinePeriodWindow(anchorOf(monthWindow, 1), "month", 1, boundStart, boundEnd);

// A free zoom through the pointer path's own transform: the level must stop
// reading as an exact week rather than keep claiming one.
var freelyZoomed = timelineZoomedWindow(
  weekWindow.windowStartMs, weekWindow.windowEndMs, 1.6, 0.5, boundStart, boundEnd);

// Closed rows above the still-open ones: the case the row-list jump exists for.
// Under newest-first order (REQ-318) the newest open REQ is usually row 0 and the
// jump is a no-op, so the fixture deliberately puts the open work lower — an old
// REQ still running under newer finished ones is what makes the second movement
// do anything.
var rows = [
  { waitOpen: false, workOpen: false },
  { waitOpen: false, workOpen: false },
  { waitOpen: false, workOpen: false },
  { waitOpen: true, workOpen: false },
  { waitOpen: false, workOpen: true }
];
var scrollHostStub = { scrollTop: 0 };
// The Now button's three steps, in the order the handler runs them: take the
// window, let the rows follow it, then scroll among THOSE rows.
//
// The two row sets below are what makes this an ORDERING assertion rather than a
// restatement. The rows array is what was on screen before the jump;
// rowsAfterJump is
// what the window the jump chose admits — a narrower set whose first open row
// sits at a different index. Deciding the scroll before the refresh, which is
// what timelineNowJump used to do internally, yields 3; deciding it after yields
// 1. Only the second is right, and only a fixture where the two disagree can
// tell them apart.
var rowsAfterJump = [
  { waitOpen: false, workOpen: false },
  { waitOpen: true, workOpen: false },
  { waitOpen: false, workOpen: true }
];
var scrollTopIfDecidedBeforeRefresh = timelineFirstOpenRowIndex(rows) * TIMELINE_ROW_HEIGHT;
var nowWindow = timelineNowJump(nowMs, queueEndMs, boundStart, boundEnd);
// THE DEGENERATE CASE, which is the ordinary one on a queue that is nearly
// drained: the forecast's queue-empty instant is minutes from the now-line, or
// there is no forecast at all, so the span the margin is a fraction OF is
// effectively zero and only the floor decides the window. Flooring on half the
// zoom floor put Now on a one-hour window — the floor itself — and the obvious
// next move was dead.
var nowWindowWithNothingScheduled = timelineNowJump(nowMs, nowMs, boundStart, boundEnd);
var nowWindowWithNoForecast = timelineNowJump(nowMs, NaN, boundStart, boundEnd);
var nowJump = { window: nowWindow };
var openRowIndex = timelineFirstOpenRowIndex(rowsAfterJump);
if (openRowIndex >= 0) {
  scrollHostStub.scrollTop = openRowIndex * TIMELINE_ROW_HEIGHT;
}

function utcParts(epochMs) {
  var instant = new Date(epochMs);
  return {
    iso: instant.toISOString(),
    weekday: instant.getUTCDay(),
    dayOfMonth: instant.getUTCDate(),
    hour: instant.getUTCHours(),
    minute: instant.getUTCMinutes()
  };
}

process.stdout.write(JSON.stringify({
  week: utcParts(weekWindow.windowStartMs),
  weekEnd: utcParts(weekWindow.windowEndMs),
  weekSpanMs: weekWindow.windowEndMs - weekWindow.windowStartMs,
  nextWeek: utcParts(nextWeek.windowStartMs),
  nextWeekStepMs: nextWeek.windowStartMs - weekWindow.windowStartMs,
  nextWeekSpanMs: nextWeek.windowEndMs - nextWeek.windowStartMs,
  prevWeekStartMs: prevWeek.windowStartMs,
  weekStartMs: weekWindow.windowStartMs,
  dayMs: TIMELINE_DAY_MS,

  atRangeEndStartMs: atRangeEnd.windowStartMs,
  atRangeEndEndMs: atRangeEnd.windowEndMs,
  atRangeEndIso: new Date(atRangeEnd.windowStartMs).toISOString() + " → " + new Date(atRangeEnd.windowEndMs).toISOString(),
  pastRangeEndStartMs: pastRangeEnd.windowStartMs,
  pastRangeEndEndMs: pastRangeEnd.windowEndMs,
  boundStartMs: boundStart,
  boundEndMs: boundEnd,

  day: utcParts(dayWindow.windowStartMs),
  daySpanMs: dayWindow.windowEndMs - dayWindow.windowStartMs,
  month: utcParts(monthWindow.windowStartMs),
  monthEnd: utcParts(monthWindow.windowEndMs),
  nextMonth: utcParts(nextMonth.windowStartMs),

  exactWeekLevel: timelinePeriodLevelOfWindow(weekWindow.windowStartMs, weekWindow.windowEndMs),
  exactDayLevel: timelinePeriodLevelOfWindow(dayWindow.windowStartMs, dayWindow.windowEndMs),
  exactMonthLevel: timelinePeriodLevelOfWindow(monthWindow.windowStartMs, monthWindow.windowEndMs),
  zoomedLevel: timelinePeriodLevelOfWindow(freelyZoomed.windowStartMs, freelyZoomed.windowEndMs),

  nowWindowIso: new Date(nowJump.window.windowStartMs).toISOString() + " → " + new Date(nowJump.window.windowEndMs).toISOString(),
  nothingScheduledSpanMs:
    nowWindowWithNothingScheduled.windowEndMs - nowWindowWithNothingScheduled.windowStartMs,
  nothingScheduledHoldsNow:
    nowMs >= nowWindowWithNothingScheduled.windowStartMs &&
    nowMs <= nowWindowWithNothingScheduled.windowEndMs,
  noForecastSpanMs: nowWindowWithNoForecast.windowEndMs - nowWindowWithNoForecast.windowStartMs,
  noForecastHoldsNow:
    nowMs >= nowWindowWithNoForecast.windowStartMs && nowMs <= nowWindowWithNoForecast.windowEndMs,
  minSpanMs: TIMELINE_MIN_SPAN_MS,
  nowInsideWindow: nowMs >= nowJump.window.windowStartMs && nowMs <= nowJump.window.windowEndMs,
  queueEndInsideWindow: queueEndMs >= nowJump.window.windowStartMs && queueEndMs <= nowJump.window.windowEndMs,
  scrollTop: scrollHostStub.scrollTop,
  wantScrollTop: 1 * TIMELINE_ROW_HEIGHT,
  scrollTopIfDecidedBeforeRefresh: scrollTopIfDecidedBeforeRefresh
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline period navigation", javascriptProbe)
	type utcParts struct {
		Iso        string `json:"iso"`
		Weekday    int    `json:"weekday"`
		DayOfMonth int    `json:"dayOfMonth"`
		Hour       int    `json:"hour"`
		Minute     int    `json:"minute"`
	}
	var periodResult struct {
		Week            utcParts `json:"week"`
		WeekEnd         utcParts `json:"weekEnd"`
		WeekSpanMs      float64  `json:"weekSpanMs"`
		NextWeek        utcParts `json:"nextWeek"`
		NextWeekStepMs  float64  `json:"nextWeekStepMs"`
		NextWeekSpanMs  float64  `json:"nextWeekSpanMs"`
		PrevWeekStartMs float64  `json:"prevWeekStartMs"`
		WeekStartMs     float64  `json:"weekStartMs"`
		DayMs           float64  `json:"dayMs"`

		AtRangeEndStartMs   float64 `json:"atRangeEndStartMs"`
		AtRangeEndEndMs     float64 `json:"atRangeEndEndMs"`
		AtRangeEndIso       string  `json:"atRangeEndIso"`
		PastRangeEndStartMs float64 `json:"pastRangeEndStartMs"`
		PastRangeEndEndMs   float64 `json:"pastRangeEndEndMs"`
		BoundStartMs        float64 `json:"boundStartMs"`
		BoundEndMs          float64 `json:"boundEndMs"`

		Day       utcParts `json:"day"`
		DaySpanMs float64  `json:"daySpanMs"`
		Month     utcParts `json:"month"`
		MonthEnd  utcParts `json:"monthEnd"`
		NextMonth utcParts `json:"nextMonth"`

		ExactWeekLevel  *string `json:"exactWeekLevel"`
		ExactDayLevel   *string `json:"exactDayLevel"`
		ExactMonthLevel *string `json:"exactMonthLevel"`
		ZoomedLevel     *string `json:"zoomedLevel"`

		NowWindowIso                    string  `json:"nowWindowIso"`
		NothingScheduledSpanMs          float64 `json:"nothingScheduledSpanMs"`
		NothingScheduledHoldsNow        bool    `json:"nothingScheduledHoldsNow"`
		NoForecastSpanMs                float64 `json:"noForecastSpanMs"`
		NoForecastHoldsNow              bool    `json:"noForecastHoldsNow"`
		MinSpanMs                       float64 `json:"minSpanMs"`
		NowInsideWindow                 bool    `json:"nowInsideWindow"`
		QueueEndInsideWindow            bool    `json:"queueEndInsideWindow"`
		ScrollTop                       float64 `json:"scrollTop"`
		WantScrollTop                   float64 `json:"wantScrollTop"`
		ScrollTopIfDecidedBeforeRefresh float64 `json:"scrollTopIfDecidedBeforeRefresh"`
	}
	if decodeError := json.Unmarshal(probeOutput, &periodResult); decodeError != nil {
		t.Fatalf("decode timeline period behavior: %v (output %q)", decodeError, probeOutput)
	}

	// A week is Monday→Sunday, starting at midnight UTC. "Seven days long" alone
	// would pass for a window starting at 14:37 on a Thursday.
	if periodResult.Week.Weekday != 1 || periodResult.Week.Hour != 0 || periodResult.Week.Minute != 0 {
		t.Fatalf("the week window starts at %s (weekday %d); want a Monday at 00:00 UTC",
			periodResult.Week.Iso, periodResult.Week.Weekday)
	}
	if periodResult.WeekSpanMs != 7*periodResult.DayMs {
		t.Fatalf("the week window spans %.0f ms (%s → %s), want exactly seven days",
			periodResult.WeekSpanMs, periodResult.Week.Iso, periodResult.WeekEnd.Iso)
	}
	if periodResult.NextWeekStepMs != 7*periodResult.DayMs {
		t.Fatalf("next period moved the window %.0f ms (%s → %s), want exactly seven days",
			periodResult.NextWeekStepMs, periodResult.Week.Iso, periodResult.NextWeek.Iso)
	}
	if periodResult.NextWeek.Weekday != 1 || periodResult.NextWeek.Hour != 0 {
		t.Fatalf("next period landed on %s (weekday %d); stepping must stay Monday-aligned",
			periodResult.NextWeek.Iso, periodResult.NextWeek.Weekday)
	}
	if periodResult.NextWeekSpanMs != periodResult.WeekSpanMs {
		t.Fatalf("next period changed the window span from %.0f ms to %.0f ms; a step moves the window, it does not resize it",
			periodResult.WeekSpanMs, periodResult.NextWeekSpanMs)
	}
	// Forward then back is the reader's undo; it has to return to the same week.
	if periodResult.PrevWeekStartMs != periodResult.WeekStartMs {
		t.Fatalf("next then previous landed on %s, want the week it started from %s",
			time.UnixMilli(int64(periodResult.PrevWeekStartMs)).UTC().Format(time.RFC3339),
			periodResult.Week.Iso)
	}

	// Held against the end of the range the window has to stop ON the data.
	// Without a clamp it walks into empty months forever, and every bar leaves
	// the screen.
	if periodResult.AtRangeEndEndMs > periodResult.BoundEndMs ||
		periodResult.AtRangeEndStartMs < periodResult.BoundStartMs {
		t.Fatalf("stepping to the end of the range left the window at %s, outside the range %.0f–%.0f",
			periodResult.AtRangeEndIso, periodResult.BoundStartMs, periodResult.BoundEndMs)
	}
	if periodResult.PastRangeEndStartMs != periodResult.AtRangeEndStartMs ||
		periodResult.PastRangeEndEndMs != periodResult.AtRangeEndEndMs {
		t.Fatalf("one more step past the last period moved the window from %s to %.0f–%.0f; it must clamp",
			periodResult.AtRangeEndIso, periodResult.PastRangeEndStartMs, periodResult.PastRangeEndEndMs)
	}

	// Day and month are calendar periods too: midnight→midnight, and the 1st to
	// the 1st — not 24h and not 30 days from wherever the view happened to be.
	if periodResult.Day.Hour != 0 || periodResult.Day.Minute != 0 || periodResult.DaySpanMs != periodResult.DayMs {
		t.Fatalf("the day window is %s spanning %.0f ms; want midnight UTC to midnight UTC",
			periodResult.Day.Iso, periodResult.DaySpanMs)
	}
	if periodResult.Month.DayOfMonth != 1 || periodResult.Month.Hour != 0 {
		t.Fatalf("the month window starts at %s; want the 1st at 00:00 UTC", periodResult.Month.Iso)
	}
	if periodResult.MonthEnd.DayOfMonth != 1 || periodResult.MonthEnd.Hour != 0 {
		t.Fatalf("the month window ends at %s; want the 1st of the next month", periodResult.MonthEnd.Iso)
	}
	if periodResult.NextMonth.DayOfMonth != 1 {
		t.Fatalf("next month landed on %s; a month step is calendar arithmetic, not 30 days", periodResult.NextMonth.Iso)
	}

	// The level is read back off the window, so it cannot claim a level the
	// window no longer has once a free zoom has resized it.
	for _, exactLevel := range []struct {
		levelName string
		reported  *string
	}{
		{"day", periodResult.ExactDayLevel},
		{"week", periodResult.ExactWeekLevel},
		{"month", periodResult.ExactMonthLevel},
	} {
		if exactLevel.reported == nil || *exactLevel.reported != exactLevel.levelName {
			t.Fatalf("a %s window reads back as %v, want %q", exactLevel.levelName, exactLevel.reported, exactLevel.levelName)
		}
	}
	if periodResult.ZoomedLevel != nil {
		t.Fatalf("after a free zoom the level still reads %q; it must report no exact level rather than claim one",
			*periodResult.ZoomedLevel)
	}

	// Now is two movements. The window has to carry both the now-line and the
	// forecast, and the ROW list has to land on the still-open work — the half
	// that recentring the time window never did.
	// Now on a drained queue: a window, not the zoom floor. This is the state the
	// button lands in whenever the forecast has nothing left to schedule, which on a
	// healthy queue is most of the time.
	for _, degenerate := range []struct {
		name     string
		spanMs   float64
		holdsNow bool
	}{
		{"the queue-empty instant equal to now", periodResult.NothingScheduledSpanMs, periodResult.NothingScheduledHoldsNow},
		{"no forecast at all", periodResult.NoForecastSpanMs, periodResult.NoForecastHoldsNow},
	} {
		if degenerate.spanMs <= periodResult.MinSpanMs {
			t.Errorf("with %s, Now lands on a %.0f ms window against a %.0f ms zoom floor; at or "+
				"below the floor there is nowhere left to zoom and no context around the now-line",
				degenerate.name, degenerate.spanMs, periodResult.MinSpanMs)
		}
		if !degenerate.holdsNow {
			t.Errorf("with %s, the window Now lands on does not contain the now-line", degenerate.name)
		}
	}

	if !periodResult.NowInsideWindow || !periodResult.QueueEndInsideWindow {
		t.Fatalf("the Now window %s does not cover both the now-line and the projected queue end", periodResult.NowWindowIso)
	}
	if periodResult.ScrollTop != periodResult.WantScrollTop {
		t.Fatalf("Now left the row list at scrollTop %.0f, want %.0f — the first still-open row",
			periodResult.ScrollTop, periodResult.WantScrollTop)
	}
	// The ordering half, and the reason the fixture carries two row sets. Deciding
	// the scroll from the PRE-jump rows — which is what timelineNowJump did before
	// REQ-319 split it — lands somewhere else entirely. If these two ever agree the
	// fixture has stopped being able to tell the orders apart, and the assertion
	// above has quietly become a restatement.
	if periodResult.ScrollTopIfDecidedBeforeRefresh == periodResult.WantScrollTop {
		t.Fatalf("the fixture's pre-jump and post-jump row sets both give scrollTop %.0f, so this "+
			"test cannot tell a scroll decided before the row refresh from one decided after; "+
			"give the two sets different first-open indices",
			periodResult.WantScrollTop)
	}
}

// The period controls are the FIRST thing a reader reaches for and, before this
// test, the first thing that took them somewhere useless. Three separate bugs
// shared one cause: timelinePeriodWindow handed an out-of-range calendar period
// to timelineZoomedWindow, which preserves a WIDTH and slides — so an edge step
// kept a week's length and lost a week's alignment — and applyPeriodWindow then
// anchored the next press on the slid window's MIDPOINT, which by then sat in the
// neighbouring period.
//
// This drives the production functions rather than reimplementing their anchor
// (REQ-305's lesson: a probe that reimplements the function under test cannot
// hold its call site), and it pins the four properties the controls owe:
//
//	a chip lands on the period containing NOW while now is on screen;
//	a chip lands near the READER when it is not;
//	a step preserves calendar alignment, at a range edge included;
//	next-then-previous is an identity, at a range edge included;
//	a step on a window that is NOT a period MOVES it rather than resizing it.
//
// The fixture's range deliberately ENDS MID-WEEK (25 Aug 2026 04:23, a Tuesday),
// because a range ending on a period boundary cannot tell a clamped period from
// an aligned one — which is exactly why the old suite passed.
func TestJavaScriptBehaviorTimelinePeriodChipsLandOnNowAndStepsStayAligned(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := timelineProbePreamble(t, "TIMELINE_MIN_SPAN_MS", "TIMELINE_DAY_MS", "TIMELINE_PAN_FRACTION") +
		rendererDeclarationLine(t, "web/board-timeline.js", "TIMELINE_PERIOD_LEVEL_NAMES") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineZoomedWindow(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelinePannedWindow(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelinePeriodStart(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineSteppedPeriodStart(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelinePeriodWindow(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelinePeriodAnchor(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelinePeriodLevelOfWindow(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelinePeriodGridOfWindow(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineSteppedWindow(") + `
// A board shaped like this repo's own: months of archive, and a range END that
// falls mid-week rather than on a boundary.
var boundStart = Date.UTC(2026, 3, 7);            // 7 Apr 2026
var boundEnd = Date.UTC(2026, 7, 25, 4, 23);      // 25 Aug 2026 04:23 UTC, a Tuesday
var nowMs = Date.UTC(2026, 7, 23, 11, 13);        // 23 Aug 2026 11:13 UTC, a Sunday
var fitted = { windowStartMs: boundStart, windowEndMs: boundEnd };

// What the view actually does when a chip is pressed, spelled with the same two
// production calls applyPeriodWindow makes.
function chipWindow(window, levelName) {
  return timelinePeriodWindow(
    timelinePeriodAnchor(0, window.windowStartMs, window.windowEndMs, nowMs),
    levelName, 0, boundStart, boundEnd);
}
// What the view actually does when an arrow is pressed.
function stepWindow(window, stepCount) {
  return timelineSteppedWindow(
    window.windowStartMs, window.windowEndMs, stepCount, boundStart, boundEnd, nowMs);
}
function spanOf(window) {
  return window.windowEndMs - window.windowStartMs;
}
function holds(window, instantMs) {
  return instantMs >= window.windowStartMs && instantMs <= window.windowEndMs;
}

// 1. A chip from the view's OPENING window. Fit all always covers the now-line,
//    so Week means the current week.
var weekFromFitAll = chipWindow(fitted, "week");
var dayFromFitAll = chipWindow(fitted, "day");
var monthFromFitAll = chipWindow(fitted, "month");

// 2. The same chip from a window the reader panned to, which does NOT hold now.
var pannedAway = { windowStartMs: Date.UTC(2026, 5, 1), windowEndMs: Date.UTC(2026, 5, 15) };
var weekFromPannedAway = chipWindow(pannedAway, "week");

// 3. Pressing a chip twice never moves the window.
var weekPressedTwice = chipWindow(weekFromFitAll, "week");

// 4. Stepping at the RIGHT EDGE. The week containing now is the last full week
//    the range reaches; one step forward lands on a week the range only partly
//    covers, and the step back has to come home.
var weekAtEdge = stepWindow(weekFromFitAll, 1);
var backFromEdge = stepWindow(weekAtEdge, -1);
var pastEdge = stepWindow(weekAtEdge, 1);

// 5. Mid-range next-then-previous, for the ordinary case.
var midWeek = chipWindow({ windowStartMs: Date.UTC(2026, 4, 11), windowEndMs: Date.UTC(2026, 4, 18) }, "week");
var midNext = stepWindow(midWeek, 1);
var midBack = stepWindow(midNext, -1);

// 6. A step on a window that is NOT a period. Nineteen days is nearer a month
//    (30d) than a week (7d), so the old nearest-level rule RESIZED it to a month;
//    an arrow has to move the window it was given.
var readersOwn = { windowStartMs: Date.UTC(2026, 5, 1), windowEndMs: Date.UTC(2026, 5, 20) };
var readersOwnStepped = stepWindow(readersOwn, 1);

// 7. Held down against the end of the range: it has to stop, on the data.
var heldAtEnd = weekFromFitAll;
for (var pressIndex = 0; pressIndex < 40; pressIndex++) {
  heldAtEnd = stepWindow(heldAtEnd, 1);
}
var heldPastEnd = stepWindow(heldAtEnd, 1);

function utcParts(epochMs) {
  var instant = new Date(epochMs);
  return {
    iso: instant.toISOString(),
    weekday: instant.getUTCDay(),
    dayOfMonth: instant.getUTCDate(),
    hour: instant.getUTCHours(),
    minute: instant.getUTCMinutes()
  };
}

process.stdout.write(JSON.stringify({
  nowIso: new Date(nowMs).toISOString(),
  boundEndIso: new Date(boundEnd).toISOString(),
  boundEndWeekday: new Date(boundEnd).getUTCDay(),

  weekFromFitAllStart: utcParts(weekFromFitAll.windowStartMs),
  weekFromFitAllEnd: utcParts(weekFromFitAll.windowEndMs),
  weekFromFitAllHoldsNow: holds(weekFromFitAll, nowMs),
  weekFromFitAllLevel: timelinePeriodLevelOfWindow(weekFromFitAll.windowStartMs, weekFromFitAll.windowEndMs),
  dayFromFitAllStart: utcParts(dayFromFitAll.windowStartMs),
  dayFromFitAllHoldsNow: holds(dayFromFitAll, nowMs),
  monthFromFitAllStart: utcParts(monthFromFitAll.windowStartMs),
  monthFromFitAllHoldsNow: holds(monthFromFitAll, nowMs),

  weekFromPannedAwayStart: utcParts(weekFromPannedAway.windowStartMs),
  weekFromPannedAwayHoldsNow: holds(weekFromPannedAway, nowMs),
  wantPannedAwayWeekStartMs: timelinePeriodStart(
    (pannedAway.windowStartMs + pannedAway.windowEndMs) / 2, "week"),
  pannedAwayWeekStartMs: weekFromPannedAway.windowStartMs,

  weekPressedTwiceStartMs: weekPressedTwice.windowStartMs,
  weekPressedTwiceEndMs: weekPressedTwice.windowEndMs,
  weekFromFitAllStartMs: weekFromFitAll.windowStartMs,
  weekFromFitAllEndMs: weekFromFitAll.windowEndMs,

  weekAtEdgeStart: utcParts(weekAtEdge.windowStartMs),
  weekAtEdgeEnd: utcParts(weekAtEdge.windowEndMs),
  weekAtEdgeStartIsAligned:
    weekAtEdge.windowStartMs === timelinePeriodStart(weekAtEdge.windowStartMs, "week"),
  weekAtEdgeStartMs: weekAtEdge.windowStartMs,
  weekAtEdgeEndMs: weekAtEdge.windowEndMs,
  wantWeekAtEdgeStartMs: timelineSteppedPeriodStart(weekFromFitAll.windowStartMs, "week", 1),
  backFromEdgeStartMs: backFromEdge.windowStartMs,
  backFromEdgeEndMs: backFromEdge.windowEndMs,
  pastEdgeStartMs: pastEdge.windowStartMs,
  pastEdgeEndMs: pastEdge.windowEndMs,

  midWeekStartMs: midWeek.windowStartMs,
  midNextStepMs: midNext.windowStartMs - midWeek.windowStartMs,
  midNextSpanMs: spanOf(midNext),
  midWeekSpanMs: spanOf(midWeek),
  midBackStartMs: midBack.windowStartMs,

  readersOwnSpanMs: spanOf(readersOwn),
  readersOwnSteppedSpanMs: spanOf(readersOwnStepped),
  readersOwnSteppedStepMs: readersOwnStepped.windowStartMs - readersOwn.windowStartMs,

  heldAtEndStartMs: heldAtEnd.windowStartMs,
  heldAtEndEndMs: heldAtEnd.windowEndMs,
  heldPastEndStartMs: heldPastEnd.windowStartMs,
  heldPastEndEndMs: heldPastEnd.windowEndMs,
  heldAtEndIso: new Date(heldAtEnd.windowStartMs).toISOString() + " → " + new Date(heldAtEnd.windowEndMs).toISOString(),
  boundStartMs: boundStart,
  boundEndMs: boundEnd,
  dayMs: TIMELINE_DAY_MS
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline period chips and steps", javascriptProbe)
	type utcParts struct {
		Iso        string `json:"iso"`
		Weekday    int    `json:"weekday"`
		DayOfMonth int    `json:"dayOfMonth"`
		Hour       int    `json:"hour"`
		Minute     int    `json:"minute"`
	}
	var chipResult struct {
		NowIso          string `json:"nowIso"`
		BoundEndIso     string `json:"boundEndIso"`
		BoundEndWeekday int    `json:"boundEndWeekday"`

		WeekFromFitAllStart     utcParts `json:"weekFromFitAllStart"`
		WeekFromFitAllEnd       utcParts `json:"weekFromFitAllEnd"`
		WeekFromFitAllHoldsNow  bool     `json:"weekFromFitAllHoldsNow"`
		WeekFromFitAllLevel     *string  `json:"weekFromFitAllLevel"`
		DayFromFitAllStart      utcParts `json:"dayFromFitAllStart"`
		DayFromFitAllHoldsNow   bool     `json:"dayFromFitAllHoldsNow"`
		MonthFromFitAllStart    utcParts `json:"monthFromFitAllStart"`
		MonthFromFitAllHoldsNow bool     `json:"monthFromFitAllHoldsNow"`

		WeekFromPannedAwayStart    utcParts `json:"weekFromPannedAwayStart"`
		WeekFromPannedAwayHoldsNow bool     `json:"weekFromPannedAwayHoldsNow"`
		WantPannedAwayWeekStartMs  float64  `json:"wantPannedAwayWeekStartMs"`
		PannedAwayWeekStartMs      float64  `json:"pannedAwayWeekStartMs"`

		WeekPressedTwiceStartMs float64 `json:"weekPressedTwiceStartMs"`
		WeekPressedTwiceEndMs   float64 `json:"weekPressedTwiceEndMs"`
		WeekFromFitAllStartMs   float64 `json:"weekFromFitAllStartMs"`
		WeekFromFitAllEndMs     float64 `json:"weekFromFitAllEndMs"`

		WeekAtEdgeStart          utcParts `json:"weekAtEdgeStart"`
		WeekAtEdgeEnd            utcParts `json:"weekAtEdgeEnd"`
		WeekAtEdgeStartIsAligned bool     `json:"weekAtEdgeStartIsAligned"`
		WeekAtEdgeStartMs        float64  `json:"weekAtEdgeStartMs"`
		WeekAtEdgeEndMs          float64  `json:"weekAtEdgeEndMs"`
		WantWeekAtEdgeStartMs    float64  `json:"wantWeekAtEdgeStartMs"`
		BackFromEdgeStartMs      float64  `json:"backFromEdgeStartMs"`
		BackFromEdgeEndMs        float64  `json:"backFromEdgeEndMs"`
		PastEdgeStartMs          float64  `json:"pastEdgeStartMs"`
		PastEdgeEndMs            float64  `json:"pastEdgeEndMs"`

		MidWeekStartMs float64 `json:"midWeekStartMs"`
		MidNextStepMs  float64 `json:"midNextStepMs"`
		MidNextSpanMs  float64 `json:"midNextSpanMs"`
		MidWeekSpanMs  float64 `json:"midWeekSpanMs"`
		MidBackStartMs float64 `json:"midBackStartMs"`

		ReadersOwnSpanMs        float64 `json:"readersOwnSpanMs"`
		ReadersOwnSteppedSpanMs float64 `json:"readersOwnSteppedSpanMs"`
		ReadersOwnSteppedStepMs float64 `json:"readersOwnSteppedStepMs"`

		HeldAtEndStartMs   float64 `json:"heldAtEndStartMs"`
		HeldAtEndEndMs     float64 `json:"heldAtEndEndMs"`
		HeldPastEndStartMs float64 `json:"heldPastEndStartMs"`
		HeldPastEndEndMs   float64 `json:"heldPastEndEndMs"`
		HeldAtEndIso       string  `json:"heldAtEndIso"`
		BoundStartMs       float64 `json:"boundStartMs"`
		BoundEndMs         float64 `json:"boundEndMs"`
		DayMs              float64 `json:"dayMs"`
	}
	if decodeError := json.Unmarshal(probeOutput, &chipResult); decodeError != nil {
		t.Fatalf("decode timeline period chip behavior: %v (output %q)", decodeError, probeOutput)
	}

	// The fixture's whole discriminating power is that the range ends mid-week. If
	// someone edits boundEnd onto a Monday this test silently stops being able to
	// tell a clamped period from a slid one.
	if chipResult.BoundEndWeekday == 1 {
		t.Fatalf("the fixture's range ends on a Monday (%s), so a slid edge window and a "+
			"cut-short one are indistinguishable; move boundEnd off the week boundary",
			chipResult.BoundEndIso)
	}

	// 1. A chip while the now-line is on screen means the period containing now.
	// This is the reported defect: from Fit all the midpoint anchor landed months
	// back, on a week with nothing drawn in it.
	if !chipResult.WeekFromFitAllHoldsNow {
		t.Fatalf("Week from the Fit-all window landed on %s → %s, which does not contain the "+
			"now-line at %s; a chip pressed while now is on screen must land on the period "+
			"holding it", chipResult.WeekFromFitAllStart.Iso, chipResult.WeekFromFitAllEnd.Iso,
			chipResult.NowIso)
	}
	if !chipResult.DayFromFitAllHoldsNow {
		t.Fatalf("Day from the Fit-all window landed on %s, which does not contain the now-line at %s",
			chipResult.DayFromFitAllStart.Iso, chipResult.NowIso)
	}
	if !chipResult.MonthFromFitAllHoldsNow {
		t.Fatalf("Month from the Fit-all window landed on %s, which does not contain the now-line at %s",
			chipResult.MonthFromFitAllStart.Iso, chipResult.NowIso)
	}
	// It must still be a real calendar week, not merely a seven-day span holding now.
	if chipResult.WeekFromFitAllStart.Weekday != 1 || chipResult.WeekFromFitAllStart.Hour != 0 ||
		chipResult.WeekFromFitAllStart.Minute != 0 {
		t.Fatalf("Week landed on %s (weekday %d); want a Monday at 00:00 UTC",
			chipResult.WeekFromFitAllStart.Iso, chipResult.WeekFromFitAllStart.Weekday)
	}
	if chipResult.WeekFromFitAllLevel == nil || *chipResult.WeekFromFitAllLevel != "week" {
		t.Fatalf("the window Week produced reads back as level %v, so the chip would not light",
			chipResult.WeekFromFitAllLevel)
	}
	if chipResult.MonthFromFitAllStart.DayOfMonth != 1 || chipResult.MonthFromFitAllStart.Hour != 0 {
		t.Fatalf("Month landed on %s; want the 1st at 00:00 UTC", chipResult.MonthFromFitAllStart.Iso)
	}

	// 2. The same chip from a window the reader panned to must stay THERE. Keying
	// the anchor on now unconditionally would drag them back to the present.
	if chipResult.WeekFromPannedAwayHoldsNow {
		t.Fatalf("Week from a window that does not contain the now-line jumped to %s, which does "+
			"contain it; a chip must not drag the reader out of the span they panned to",
			chipResult.WeekFromPannedAwayStart.Iso)
	}
	if chipResult.PannedAwayWeekStartMs != chipResult.WantPannedAwayWeekStartMs {
		t.Fatalf("Week from a panned-away window landed on %s, want the week holding that window's "+
			"own midpoint", chipResult.WeekFromPannedAwayStart.Iso)
	}

	// 3. Idempotence. A reader who presses Week twice has asked one question.
	if chipResult.WeekPressedTwiceStartMs != chipResult.WeekFromFitAllStartMs ||
		chipResult.WeekPressedTwiceEndMs != chipResult.WeekFromFitAllEndMs {
		t.Fatalf("pressing Week twice moved the window from %s to %s; a chip press is idempotent",
			chipResult.WeekFromFitAllStart.Iso,
			time.UnixMilli(int64(chipResult.WeekPressedTwiceStartMs)).UTC().Format(time.RFC3339))
	}

	// 4. A step at a range edge keeps its calendar alignment. The reported failure
	// was a window still seven days long that no longer started on a Monday.
	if !chipResult.WeekAtEdgeStartIsAligned {
		t.Fatalf("stepping into the last, partly-covered week started the window at %s, which is "+
			"not a week boundary; an edge period is cut short, never slid",
			chipResult.WeekAtEdgeStart.Iso)
	}
	if chipResult.WeekAtEdgeStartMs != chipResult.WantWeekAtEdgeStartMs {
		t.Fatalf("one step forward from %s started the window at %s, want the next week boundary %s",
			chipResult.WeekFromFitAllStart.Iso, chipResult.WeekAtEdgeStart.Iso,
			time.UnixMilli(int64(chipResult.WantWeekAtEdgeStartMs)).UTC().Format(time.RFC3339))
	}
	if chipResult.WeekAtEdgeEnd.Iso != chipResult.BoundEndIso {
		t.Fatalf("the last week's window ends at %s, want the range end %s — cut short, not slid",
			chipResult.WeekAtEdgeEnd.Iso, chipResult.BoundEndIso)
	}
	// The identity, at the edge, which is where it broke: forward then back
	// returned a whole week EARLIER than where the reader started.
	if chipResult.BackFromEdgeStartMs != chipResult.WeekFromFitAllStartMs ||
		chipResult.BackFromEdgeEndMs != chipResult.WeekFromFitAllEndMs {
		t.Fatalf("next then previous at the range edge landed on %s → %s, want the window it "+
			"started from %s → %s",
			time.UnixMilli(int64(chipResult.BackFromEdgeStartMs)).UTC().Format(time.RFC3339),
			time.UnixMilli(int64(chipResult.BackFromEdgeEndMs)).UTC().Format(time.RFC3339),
			chipResult.WeekFromFitAllStart.Iso, chipResult.WeekFromFitAllEnd.Iso)
	}
	if chipResult.PastEdgeStartMs != chipResult.WeekAtEdgeStartMs ||
		chipResult.PastEdgeEndMs != chipResult.WeekAtEdgeEndMs {
		t.Fatalf("one more step past the last period moved the window; it must clamp on the period")
	}

	// 5. The ordinary mid-range case still steps exactly one period and returns.
	if chipResult.MidNextStepMs != 7*chipResult.DayMs {
		t.Fatalf("a mid-range next moved the window %.0f ms, want exactly seven days", chipResult.MidNextStepMs)
	}
	if chipResult.MidNextSpanMs != chipResult.MidWeekSpanMs {
		t.Fatalf("a mid-range next resized the window from %.0f ms to %.0f ms; a step moves, it does not resize",
			chipResult.MidWeekSpanMs, chipResult.MidNextSpanMs)
	}
	if chipResult.MidBackStartMs != chipResult.MidWeekStartMs {
		t.Fatalf("mid-range next then previous did not return to the same week")
	}

	// 6. An arrow on a window of the reader's own MOVES it. Snapping to the
	// nearest period level resized nineteen days into a calendar month.
	if chipResult.ReadersOwnSteppedSpanMs != chipResult.ReadersOwnSpanMs {
		t.Fatalf("stepping a %.0f-day window resized it to %.0f days; an arrow moves the window, "+
			"a chip resizes it", chipResult.ReadersOwnSpanMs/chipResult.DayMs,
			chipResult.ReadersOwnSteppedSpanMs/chipResult.DayMs)
	}
	if chipResult.ReadersOwnSteppedStepMs != chipResult.ReadersOwnSpanMs {
		t.Fatalf("stepping a window of the reader's own moved it %.0f ms, want one screenful (%.0f ms)",
			chipResult.ReadersOwnSteppedStepMs, chipResult.ReadersOwnSpanMs)
	}

	// 7. Held down, it stops on the data rather than walking off it.
	if chipResult.HeldAtEndEndMs > chipResult.BoundEndMs || chipResult.HeldAtEndStartMs < chipResult.BoundStartMs {
		t.Fatalf("holding next left the window at %s, outside the range", chipResult.HeldAtEndIso)
	}
	if chipResult.HeldPastEndStartMs != chipResult.HeldAtEndStartMs ||
		chipResult.HeldPastEndEndMs != chipResult.HeldAtEndEndMs {
		t.Fatalf("one more press after holding next moved the window from %s; it must clamp",
			chipResult.HeldAtEndIso)
	}
}

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
func TestJavaScriptBehaviorTimelineAxisLabelsNameTheirOwnInstant(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := timelineProbePreamble(t, "TIMELINE_MIN_SPAN_MS", "TIMELINE_DAY_MS",
		"TIMELINE_AXIS_TICK_COUNT", "TIMELINE_AXIS_TICK_LIMIT") +
		rendererDeclarationLine(t, "web/board-timeline.js", "TIMELINE_WEEK_ALIGNMENT_MS") + "\n" +
		rendererDeclarationLine(t, "web/board-timeline.js", "TIMELINE_MONTHS") + "\n" +
		rendererBracketDeclaration(t, "web/board-timeline.js", "TIMELINE_AXIS_TICK_STEPS") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineTickStepSpanMs(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineAxisTickStep(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineTickAtOrBefore(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineSteppedTick(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineAxisTicks(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelinePeriodStart(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineFormatAxisTick(") + `
// The REAL tick source and the REAL formatter, called the way renderAxis calls
// them — including passing the gap that positioned the ticks rather than deriving
// one here.
function axisTicks(name, startMs, endMs) {
  var chosen = timelineAxisTicks(startMs, endMs);
  var ticks = chosen.instants.map(function (tickMs) {
    var instant = new Date(tickMs);
    return {
      epochMs: tickMs,
      label: timelineFormatAxisTick(tickMs, chosen.gapMs, startMs, endMs),
      dayOfMonth: instant.getUTCDate(),
      hour: instant.getUTCHours(),
      minute: instant.getUTCMinutes(),
      year: instant.getUTCFullYear()
    };
  });
  return {
    name: name, ticks: ticks, startMs: startMs, endMs: endMs, gapMs: chosen.gapMs,
    // Whether every tick sits where timelinePeriodStart puts a week — the SHIPPED
    // week rule, read rather than restated, so moving Monday moves both together
    // or fails here (REQ-322).
    everyTickOnAWeekBoundary: chosen.instants.every(function (tickMs) {
      return timelinePeriodStart(tickMs, "week") === tickMs;
    }),
    tickWeekdays: chosen.instants.map(function (tickMs) { return new Date(tickMs).getUTCDay(); })
  };
}

var mondayMs = Date.UTC(2026, 7, 17);      // 17 Aug 2026 is a Monday
process.stdout.write(JSON.stringify({
  windows: [
    // Where the Now button lands: a window covering the now-line and the
    // forecast's queue-empty instant, which on a healthy queue is well under an
    // hour, so the span settles near the view's floor and the start is wherever
    // "now" fell — 11:26, not the top of an hour.
    axisTicks("Now", Date.UTC(2026, 7, 18, 11, 26), Date.UTC(2026, 7, 18, 11, 26) + TIMELINE_MIN_SPAN_MS),
    axisTicks("Day", Date.UTC(2026, 7, 18), Date.UTC(2026, 7, 19)),
    // A free zoom between the period levels: not a whole number of anything.
    axisTicks("free zoom, four days", Date.UTC(2026, 7, 15), Date.UTC(2026, 7, 19)),
    axisTicks("Week", mondayMs, mondayMs + 7 * TIMELINE_DAY_MS),
    axisTicks("Month", Date.UTC(2026, 7, 1), Date.UTC(2026, 8, 1)),
    axisTicks("Fit all", Date.UTC(2026, 3, 7), Date.UTC(2026, 7, 18)),
    // Three months, which is what Fit all measures on this repo's own board. It
    // picks the FORTNIGHT rung, and is the only fixture here that does — without
    // it the week-boundary alignment of that rung is never checked.
    axisTicks("Fit all, three months", Date.UTC(2026, 4, 27), Date.UTC(2026, 7, 25)),
    // Fit all is the whole capture history, and it only grows. Once it crosses a
    // calendar year one day-and-month comes round twice.
    axisTicks("Fit all across two years", Date.UTC(2025, 7, 18), Date.UTC(2027, 7, 18)),
    // Nine days, crossing New Year. Shorter than a year and still ambiguous
    // without one — the case the old spanMs >= TIMELINE_YEAR_MS threshold missed.
    axisTicks("across New Year", Date.UTC(2026, 11, 28), Date.UTC(2027, 0, 6))
  ]
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline axis labels", javascriptProbe)
	var axisResult struct {
		Windows []struct {
			Name  string `json:"name"`
			Ticks []struct {
				EpochMs    float64 `json:"epochMs"`
				Label      string  `json:"label"`
				DayOfMonth int     `json:"dayOfMonth"`
				Hour       int     `json:"hour"`
				Minute     int     `json:"minute"`
				Year       int     `json:"year"`
			} `json:"ticks"`
			StartMs                  float64 `json:"startMs"`
			EndMs                    float64 `json:"endMs"`
			GapMs                    float64 `json:"gapMs"`
			EveryTickOnAWeekBoundary bool    `json:"everyTickOnAWeekBoundary"`
			TickWeekdays             []int   `json:"tickWeekdays"`
		} `json:"windows"`
	}
	if decodeError := json.Unmarshal(probeOutput, &axisResult); decodeError != nil {
		t.Fatalf("decode timeline axis behavior: %v (output %q)", decodeError, probeOutput)
	}

	// What each window's labels have to look like. The three period windows and
	// Fit all are here to hold their EXISTING labels: this is a formatting fix,
	// and the formats that were already right may not move.
	const (
		axisLabelWithTime = "date and time"
		axisLabelDateOnly = "date alone"
		axisLabelWithYear = "date and year"
	)
	wantAxisLabelShape := map[string]string{
		"Now":                   axisLabelWithTime,
		"Day":                   axisLabelWithTime,
		"free zoom, four days":  axisLabelDateOnly,
		"Week":                  axisLabelDateOnly,
		"Month":                 axisLabelDateOnly,
		"Fit all":               axisLabelDateOnly,
		"Fit all, three months": axisLabelDateOnly,
		// Both of these cross a calendar year, which is what earns the year — not
		// being longer than 365 days. "across New Year" is nine days long.
		"Fit all across two years": axisLabelWithYear,
		"across New Year":          axisLabelWithYear,
	}
	if len(axisResult.Windows) != len(wantAxisLabelShape) {
		t.Fatalf("the probe drove %d windows, want the %d named", len(axisResult.Windows), len(wantAxisLabelShape))
	}

	weekGapsSeen := 0
	for _, window := range axisResult.Windows {
		labelShape, isNamed := wantAxisLabelShape[window.Name]
		if !isNamed {
			t.Fatalf("the probe drove an unnamed window %q", window.Name)
		}
		renderedLabels := make([]string, 0, len(window.Ticks))
		distinctLabels := map[string]bool{}
		for _, tick := range window.Ticks {
			renderedLabels = append(renderedLabels, tick.Label)
			distinctLabels[tick.Label] = true
		}
		// Two ticks at different instants reading the same label is what makes
		// the axis unreadable rather than merely imprecise.
		if len(distinctLabels) != len(window.Ticks) {
			t.Fatalf("the %s window draws %d ticks with only %d distinct labels: %q",
				window.Name, len(window.Ticks), len(distinctLabels), renderedLabels)
		}
		// THE TICKS THEMSELVES. Ascending, inside the window, and there at all —
		// without this the label assertions below would pass over an empty axis.
		if len(window.Ticks) < 3 {
			t.Fatalf("the %s window drew %d ticks; an axis with fewer than three is not one",
				window.Name, len(window.Ticks))
		}
		for tickIndex, tick := range window.Ticks {
			if tick.EpochMs < window.StartMs || tick.EpochMs > window.EndMs {
				t.Fatalf("the %s window drew a tick at %s, outside the window it describes",
					window.Name, tick.Label)
			}
			if tickIndex > 0 && tick.EpochMs <= window.Ticks[tickIndex-1].EpochMs {
				t.Fatalf("the %s window's ticks are not strictly ascending at %q", window.Name, tick.Label)
			}
		}
		// A WEEK-LONG GAP LANDS ON THE WEEK BOUNDARY THE REST OF THE VIEW USES.
		// Aligning it to the epoch instead gives Thursdays — still midnights, still
		// distinct, still inside the window, so every other assertion here passes
		// and the axis silently disagrees with the Week chip beside it.
		const oneWeekMs = 7 * 24 * 60 * 60 * 1000
		if window.GapMs == oneWeekMs || window.GapMs == 2*oneWeekMs {
			weekGapsSeen++
			if !window.EveryTickOnAWeekBoundary {
				t.Fatalf("the %s window uses a %.0f-day gap but its ticks fall on weekdays %v; a "+
					"week-long gap has to land where timelinePeriodStart puts a week",
					window.Name, window.GapMs/float64(24*60*60*1000), window.TickWeekdays)
			}
		}
		// A LABEL WITH NO TIME IS A CLAIM OF MIDNIGHT. This is the assertion the old
		// version of this test did not make, and the whole reason the week axis could
		// print "9 Jul" for a tick at 9 Jul 12:00 with the suite green.
		if labelShape != axisLabelWithTime {
			for _, tick := range window.Ticks {
				if tick.Hour != 0 || tick.Minute != 0 {
					t.Fatalf("the %s window labels the tick at %02d:%02d on the %dth as %q, with no "+
						"time in it — a date-only label claims midnight, so a tick that is not at "+
						"midnight may not have one",
						window.Name, tick.Hour, tick.Minute, tick.DayOfMonth, tick.Label)
				}
			}
		}
		// Every number in the label has to be one the instant carries. Matching
		// the whole label also pins the shape, so a window cannot quietly gain
		// or lose a component the reader relies on.
		for _, tick := range window.Ticks {
			var wantLabelPattern string
			switch labelShape {
			case axisLabelWithTime:
				wantLabelPattern = fmt.Sprintf(`^%d [A-Z][a-z]{2} %02d:%02d$`, tick.DayOfMonth, tick.Hour, tick.Minute)
			case axisLabelDateOnly:
				wantLabelPattern = fmt.Sprintf(`^%d [A-Z][a-z]{2}$`, tick.DayOfMonth)
			case axisLabelWithYear:
				wantLabelPattern = fmt.Sprintf(`^%d [A-Z][a-z]{2} %d$`, tick.DayOfMonth, tick.Year)
			}
			if !regexp.MustCompile(wantLabelPattern).MatchString(tick.Label) {
				t.Fatalf("the %s window renders the tick at %d/%02d:%02d/%d as %q, want %s matching %s",
					window.Name, tick.DayOfMonth, tick.Hour, tick.Minute, tick.Year,
					tick.Label, labelShape, wantLabelPattern)
			}
		}
	}
	// The week-boundary assertion above is inside a conditional, so it is worth
	// nothing if no fixture window ever picks a week-long gap. The Month window is
	// the one that does; if the ladder is re-tuned so none of them do, this says so
	// instead of quietly passing.
	// Both rungs, counted separately: one fixture hitting the 7-day rung leaves the
	// 14-day rung's alignment unchecked, which is exactly how a mutation of it
	// passed the first time this ran.
	if weekGapsSeen < 2 {
		t.Errorf("only %d fixture window(s) chose a week-length gap; both the 7-day and the "+
			"14-day rung need one, or the alignment of the unvisited rung is never checked",
			weekGapsSeen)
	}
}

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
func TestJavaScriptBehaviorSpanFormattersCarryRoundedRemainders(t *testing.T) {
	roundingCases := []spanRoundingCase{
		{119.5, "2h 0m", "2h 0m", "the minute remainder rounds to a full hour"},
		{-119.5, "−2h 0m", "−2h 0m", "the same carry on a reversed span"},
		{59.96, "1h 0m", "1h 0m", "the sub-hour branch rounds up to the hour boundary"},
		{119.4, "1h 59m", "1h 59m", "just under the carry still splits normally"},
		{2879, "47h 59m", "2d 0h", "the hour remainder carries into the day field"},
		{1439.5, "24h 0m", "1d 0h", "a day boundary carries the same way"},
		{7.5, "7.5 min", "8 min", "each formatter keeps its own sub-hour precision"},
	}
	probeValues, encodeError := json.Marshal(roundingCases)
	if encodeError != nil {
		t.Fatalf("encode probe values: %v", encodeError)
	}

	indexHtml := generateLiveSite(t)
	javascriptProbe := sliceBalancedBlockAfter(t, indexHtml, "function formatDurationMinutes(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineFormatSpanMinutes(") + `
process.stdout.write(JSON.stringify(` + string(probeValues) + `.map(function (roundingCase) {
  return {
    durationsText: formatDurationMinutes(roundingCase.minutes),
    timelineText: timelineFormatSpanMinutes(roundingCase.minutes)
  };
})));`

	probeOutput := runJavaScriptBehaviorProbe(t, "span formatter rounding", javascriptProbe)
	var drawnTexts []struct {
		DurationsText string `json:"durationsText"`
		TimelineText  string `json:"timelineText"`
	}
	if decodeError := json.Unmarshal(probeOutput, &drawnTexts); decodeError != nil {
		t.Fatalf("decode span formatting: %v (output %q)", decodeError, probeOutput)
	}
	if len(drawnTexts) != len(roundingCases) {
		t.Fatalf("probe returned %d results, want %d", len(drawnTexts), len(roundingCases))
	}
	for caseIndex, roundingCase := range roundingCases {
		if drawnTexts[caseIndex].DurationsText != roundingCase.WantDurationsText {
			t.Errorf("%s: formatDurationMinutes(%.2f) drew %q, want %q",
				roundingCase.Requirement, roundingCase.Minutes,
				drawnTexts[caseIndex].DurationsText, roundingCase.WantDurationsText)
		}
		if drawnTexts[caseIndex].TimelineText != roundingCase.WantTimelineText {
			t.Errorf("%s: timelineFormatSpanMinutes(%.2f) drew %q, want %q",
				roundingCase.Requirement, roundingCase.Minutes,
				drawnTexts[caseIndex].TimelineText, roundingCase.WantTimelineText)
		}
	}
}

// TestJavaScriptBehaviorCalendarDayBreakdownGroupsStatuses executes the real
// calendarDayBreakdown from the assembled client. The count line it feeds ("12
// done · 2 cancelled") is the calendar's main at-a-glance signal, so a status
// landing in the wrong group misreports the day — counting abandoned or
// still-open work as "done" is the failure that matters. It also pins that only
// non-zero groups render, that the three blocked variants collapse into one
// group while an unrecognized status does NOT join them, and the fixed group
// order the colours in board.css are written against.
func TestJavaScriptBehaviorCalendarDayBreakdownGroupsStatuses(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := sliceBalancedBlockAfter(t, indexHtml, "function calendarDayBreakdown(") + `
var entries = [
  { id: "REQ-1", status: "cancelled" },
  { id: "REQ-2", status: "completed" },
  { id: "REQ-3", status: "blocked-archive-collision" },
  { id: "REQ-4", status: "completed" },
  { id: "REQ-5", status: "blocked" },
  { id: "REQ-6", status: "completed-with-issues" },
  { id: "REQ-7", status: "blockd-dependency-cycle" },
  { id: "REQ-8", status: "claimed" },
  { id: "REQ-9" }
];
process.stdout.write(JSON.stringify(calendarDayBreakdown(entries)));`
	probeOutput := runJavaScriptBehaviorProbe(t, "calendar day breakdown", javascriptProbe)

	var breakdown []struct {
		Group string `json:"group"`
		Label string `json:"label"`
		Count int    `json:"count"`
	}
	if decodeError := json.Unmarshal(probeOutput, &breakdown); decodeError != nil {
		t.Fatalf("decode calendar day breakdown: %v (output %q)", decodeError, probeOutput)
	}
	want := []struct {
		group string
		count int
	}{
		{"done", 2},
		{"with-issues", 1},
		{"claimed", 1},
		{"blocked", 2},      // `blocked` + `blocked-archive-collision`, one group
		{"cancelled", 1},    // never folded into done
		{"unrecognized", 2}, // the typo'd status and the one with no status at all
	}
	if len(breakdown) != len(want) {
		t.Fatalf("breakdown = %#v, want %d non-zero groups (empty groups must not render)", breakdown, len(want))
	}
	for index, wantPart := range want {
		if breakdown[index].Group != wantPart.group || breakdown[index].Count != wantPart.count {
			t.Fatalf("breakdown[%d] = %s×%d, want %s×%d (fixed group order, exact status matching)",
				index, breakdown[index].Group, breakdown[index].Count, wantPart.group, wantPart.count)
		}
	}
}

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
func TestJavaScriptBehaviorTimelineRefusesToRenderAgainstAnUnmeasurableHost(t *testing.T) {
	rendererFragment, readError := embeddedWebAssets.ReadFile("web/board-timeline.js")
	if readError != nil {
		t.Fatalf("read web/board-timeline.js: %v", readError)
	}

	// Twelve rows, so an eight-row viewport is visibly a truncation rather than a
	// coincidence, spread over four hours so a 120px plot is visibly a crush.
	timelinePayload := `{
	  "now": "2026-08-18T13:00:00Z",
	  "rangeStart": "2026-08-18T09:00:00Z",
	  "rangeEnd": "2026-08-18T13:00:00Z",
	  "rows": [
	    {"id":"REQ-901","createdTime":"2026-08-18T09:00:00Z","claimedTime":"2026-08-18T09:10:00Z","completedTime":"2026-08-18T09:40:00Z","waitMinutes":10,"workMinutes":30,"waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-902","createdTime":"2026-08-18T09:20:00Z","claimedTime":"2026-08-18T09:30:00Z","completedTime":"2026-08-18T10:00:00Z","waitMinutes":10,"workMinutes":30,"waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-903","createdTime":"2026-08-18T09:40:00Z","claimedTime":"2026-08-18T09:50:00Z","completedTime":"2026-08-18T10:20:00Z","waitMinutes":10,"workMinutes":30,"waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-904","createdTime":"2026-08-18T10:00:00Z","claimedTime":"2026-08-18T10:10:00Z","completedTime":"2026-08-18T10:40:00Z","waitMinutes":10,"workMinutes":30,"waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-905","createdTime":"2026-08-18T10:20:00Z","claimedTime":"2026-08-18T10:30:00Z","completedTime":"2026-08-18T11:00:00Z","waitMinutes":10,"workMinutes":30,"waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-906","createdTime":"2026-08-18T10:40:00Z","claimedTime":"2026-08-18T10:50:00Z","completedTime":"2026-08-18T11:20:00Z","waitMinutes":10,"workMinutes":30,"waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-907","createdTime":"2026-08-18T11:00:00Z","claimedTime":"2026-08-18T11:10:00Z","completedTime":"2026-08-18T11:40:00Z","waitMinutes":10,"workMinutes":30,"waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-908","createdTime":"2026-08-18T11:20:00Z","claimedTime":"2026-08-18T11:30:00Z","completedTime":"2026-08-18T12:00:00Z","waitMinutes":10,"workMinutes":30,"waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-909","createdTime":"2026-08-18T11:40:00Z","claimedTime":"2026-08-18T11:50:00Z","completedTime":"2026-08-18T12:20:00Z","waitMinutes":10,"workMinutes":30,"waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-910","createdTime":"2026-08-18T12:00:00Z","claimedTime":"2026-08-18T12:10:00Z","completedTime":"2026-08-18T12:40:00Z","waitMinutes":10,"workMinutes":30,"waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-911","createdTime":"2026-08-18T12:10:00Z","claimedTime":"2026-08-18T12:20:00Z","completedTime":"2026-08-18T12:50:00Z","waitMinutes":10,"workMinutes":30,"waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-912","createdTime":"2026-08-18T12:20:00Z","claimedTime":"2026-08-18T12:30:00Z","completedTime":"2026-08-18T12:55:00Z","waitMinutes":10,"workMinutes":25,"waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false}
	  ]
	}`

	// Render measurable, then unmeasurable, then measurable again — the three
	// states a reader passes through when they resize the browser on another view
	// and come back. Both renders after the first go through the same entry point
	// the resize listener uses.
	probeDriver := `
function drawnRowIds() {
  var ids = [];
  (function walk(node) {
    (node.children || []).forEach(function (childNode) {
      var attributes = childNode.attributes || {};
      if (childNode.stubName === "g" && attributes["data-detail-id"]) {
        ids.push(attributes["data-detail-id"]);
        return;
      }
      walk(childNode);
    });
  })(timelineStubHosts["timeline-scroll"]);
  return ids;
}
function countDescendants(node, stubName) {
  var found = 0;
  (node.children || []).forEach(function (childNode) {
    if (childNode.stubName === stubName) { found++; }
    found += countDescendants(childNode, stubName);
  });
  return found;
}
function hostSize(widthPx, heightPx) {
  var host = timelineStubHosts["timeline-scroll"];
  host.clientWidth = widthPx;
  host.clientHeight = heightPx;
  host.getBoundingClientRect = function () {
    return { width: widthPx, height: heightPx, left: 0, top: 0 };
  };
}
// The stub's textContent is a plain property, so the renderer's own
// "scrollHost.textContent = \"\"" does not drop its children — every other probe
// in this lane renders once and never notices. Three renders do, so the fixture
// clears what a real DOM would have cleared. This is stub bookkeeping, not a
// production behaviour: getting it wrong made the second render's row list read
// as the first render's.
function clearRenderedHosts() {
  ["timeline-scroll", "timeline-axis", "timeline-table-body"].forEach(function (hostId) {
    timelineStubHosts[hostId].children = [];
  });
  timelineStubHosts["timeline-summary"].textContent = "";
}
function snapshot() {
  return {
    rowIds: drawnRowIds(),
    summary: timelineStubHosts["timeline-summary"].textContent,
    // The axis SVG SHELL is created by renderTimelineView before renderAll runs,
    // so counting the host's children counts a container that is always there.
    // What matters is whether any tick was drawn INTO it.
    axisTicks: countDescendants(timelineStubHosts["timeline-axis"], "text")
  };
}
hostSize(900, 400);
renderTimelineView();
var measured = snapshot();

hostSize(0, 0);
clearRenderedHosts();
renderTimelineView();
var unmeasurable = snapshot();

hostSize(900, 400);
clearRenderedHosts();
renderTimelineView();
var remeasured = snapshot();

process.stdout.write(JSON.stringify({
  measured: measured, unmeasurable: unmeasurable, remeasured: remeasured
}));
`

	javascriptProbe := timelineRenderDomStubPreamble +
		"var boardData = { timeline: " + timelinePayload + " };\n" +
		string(rendererFragment) +
		probeDriver
	probeOutput := runJavaScriptBehaviorProbe(t, "timeline unmeasurable host", javascriptProbe)

	type renderSnapshot struct {
		RowIds    []string `json:"rowIds"`
		Summary   string   `json:"summary"`
		AxisTicks int      `json:"axisTicks"`
	}
	var hostResult struct {
		Measured     renderSnapshot `json:"measured"`
		Unmeasurable renderSnapshot `json:"unmeasurable"`
		Remeasured   renderSnapshot `json:"remeasured"`
	}
	if decodeError := json.Unmarshal(probeOutput, &hostResult); decodeError != nil {
		t.Fatalf("decode timeline unmeasurable-host behavior: %v (output %q)", decodeError, probeOutput)
	}

	// SETUP, ASSERTED: without a full first render there is nothing for the
	// unmeasurable render to be compared against.
	if hostResult.Measured.AxisTicks == 0 {
		t.Fatal("the measurable render drew no axis tick labels, so the zero-tick assertion below " +
			"would pass against any render at all")
	}
	if len(hostResult.Measured.RowIds) != 12 {
		t.Fatalf("the measurable render drew %d rows, want all 12 fixture rows; the probe is not "+
			"measuring a full chart", len(hostResult.Measured.RowIds))
	}
	if !strings.Contains(hostResult.Measured.Summary, "12 REQs in the window") {
		t.Fatalf("the measurable render's summary is %q, want it to name all 12 rows", hostResult.Measured.Summary)
	}

	// THE DEFECT. Eight rows and a 120px plot are not a smaller truth, they are a
	// measurement of a box that does not exist.
	if len(hostResult.Unmeasurable.RowIds) != 0 {
		t.Fatalf("rendering against a zero-width host drew %d rows (%v); an unmeasurable host is "+
			"\"not yet\", not a 120px plot with an eight-row viewport",
			len(hostResult.Unmeasurable.RowIds), hostResult.Unmeasurable.RowIds)
	}
	if hostResult.Unmeasurable.Summary != "" {
		t.Fatalf("rendering against a zero-width host wrote the summary %q; it must not describe a "+
			"window it could not lay out", hostResult.Unmeasurable.Summary)
	}
	if hostResult.Unmeasurable.AxisTicks != 0 {
		t.Fatalf("rendering against a zero-width host drew %d axis tick labels", hostResult.Unmeasurable.AxisTicks)
	}

	// And the numbers come back whole once the host has a box again, which is what
	// the ResizeObserver triggers in a real engine.
	if len(hostResult.Remeasured.RowIds) != len(hostResult.Measured.RowIds) {
		t.Fatalf("after the host regained its box the render drew %d rows, want the %d it drew before",
			len(hostResult.Remeasured.RowIds), len(hostResult.Measured.RowIds))
	}
	if hostResult.Remeasured.Summary != hostResult.Measured.Summary {
		t.Fatalf("after the host regained its box the summary reads %q, want the %q it read before",
			hostResult.Remeasured.Summary, hostResult.Measured.Summary)
	}
}

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
func TestJavaScriptBehaviorReversedWaitDrawsAsABreak(t *testing.T) {
	rendererFragment, readError := embeddedWebAssets.ReadFile("web/board-timeline.js")
	if readError != nil {
		t.Fatalf("read web/board-timeline.js: %v", readError)
	}

	// claimed_at precedes created_at on REQ-901 only. REQ-902 is an ordinary
	// closed wait; REQ-903 is an unclaimed REQ whose wait runs to the now-line.
	timelinePayload := `{
	  "now": "2026-08-18T12:00:00Z",
	  "rangeStart": "2026-08-18T09:00:00Z",
	  "rangeEnd": "2026-08-18T13:00:00Z",
	  "rows": [
	    {"id":"REQ-901","createdTime":"2026-08-18T11:00:00Z","claimedTime":"2026-08-18T10:00:00Z",
	     "completedTime":"2026-08-18T11:30:00Z","waitMinutes":-60,"workMinutes":90,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-902","createdTime":"2026-08-18T10:00:00Z","claimedTime":"2026-08-18T10:30:00Z",
	     "completedTime":"2026-08-18T11:00:00Z","waitMinutes":30,"workMinutes":30,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-903","createdTime":"2026-08-18T11:00:00Z","claimedTime":null,
	     "completedTime":null,"waitMinutes":60,"workMinutes":0,
	     "waitOpen":true,"workOpen":false,"hasWork":false,"anomaly":false}
	  ]
	}`

	javascriptProbe := timelineRenderDomStubPreamble +
		"var boardData = { timeline: " + timelinePayload + " };\n" +
		string(rendererFragment) +
		timelineRenderProbeDriver
	probeOutput := runJavaScriptBehaviorProbe(t, "timeline reversed wait", javascriptProbe)

	var drawn struct {
		Rows []struct {
			Id    string `json:"id"`
			Rects []struct {
				Class string  `json:"class"`
				Width float64 `json:"width"`
			} `json:"rects"`
		} `json:"rows"`
	}
	if decodeError := json.Unmarshal(probeOutput, &drawn); decodeError != nil {
		t.Fatalf("decode drawn timeline rows: %v (output starts %q)",
			decodeError, string(probeOutput[:min(len(probeOutput), 400)]))
	}
	if len(drawn.Rows) != 3 {
		t.Fatalf("want one drawn group per fixture row, got %d", len(drawn.Rows))
	}

	rowClasses := map[string][]string{}
	rowWidths := map[string]map[string]float64{}
	for _, drawnRow := range drawn.Rows {
		rowWidths[drawnRow.Id] = map[string]float64{}
		for _, rect := range drawnRow.Rects {
			rowClasses[drawnRow.Id] = append(rowClasses[drawnRow.Id], rect.Class)
			rowWidths[drawnRow.Id][rect.Class] = rect.Width
		}
	}

	rowDrewClassContaining := func(rowId string, classFragment string) bool {
		for _, drawnClass := range rowClasses[rowId] {
			if strings.Contains(drawnClass, classFragment) {
				return true
			}
		}
		return false
	}

	// (1) The reversed wait is a break, and there is no wait bar to misread.
	if !rowDrewClassContaining("REQ-901", "timeline-segment-broken") {
		t.Errorf("a wait whose claim precedes its capture must draw the break marker, got %v", rowClasses["REQ-901"])
	}
	if rowDrewClassContaining("REQ-901", "timeline-segment-wait") {
		t.Errorf("a reversed wait must draw NO wait bar — the table prints the negative value beside it, got %v",
			rowClasses["REQ-901"])
	}
	// Its work span is ordinary and must be untouched by the wait's branch.
	if !rowDrewClassContaining("REQ-901", "timeline-segment-work") {
		t.Errorf("a reversed wait must not suppress the row's ordinary work bar, got %v", rowClasses["REQ-901"])
	}

	// (2) An ordinary closed wait still draws its bar, with real width.
	if !rowDrewClassContaining("REQ-902", "timeline-segment-wait") {
		t.Errorf("an ordinary positive wait must still draw its bar, got %v", rowClasses["REQ-902"])
	}
	if rowDrewClassContaining("REQ-902", "timeline-segment-broken") {
		t.Errorf("an ordinary positive wait must not draw a break marker, got %v", rowClasses["REQ-902"])
	}

	// (3) The open wait keeps its is-open bar: it is measured to the now-line and
	// is never reversed, so the new branch must not reach it.
	if !rowDrewClassContaining("REQ-903", "timeline-segment-wait is-open") {
		t.Errorf("an unclaimed REQ must still draw its open wait bar, got %v", rowClasses["REQ-903"])
	}
	if rowDrewClassContaining("REQ-903", "timeline-segment-broken") {
		t.Errorf("an open wait must not draw a break marker, got %v", rowClasses["REQ-903"])
	}

	// The break marker is a fixed-width mark, not a measured span — a break whose
	// width tracked the reversed magnitude would be the same lie in a new shape.
	if brokenWidth := rowWidths["REQ-901"]["timeline-segment-broken"]; brokenWidth != 6 {
		t.Errorf("the break marker must be the same fixed 6-unit mark the work branch draws, got %v", brokenWidth)
	}
}

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
func TestJavaScriptBehaviorTimelineFallbackBoundsSpanTheWholeMatchedSet(t *testing.T) {
	rendererFragment, readError := embeddedWebAssets.ReadFile("web/board-timeline.js")
	if readError != nil {
		t.Fatalf("read web/board-timeline.js: %v", readError)
	}

	// rangeStart is deliberately unparseable, which is what takes the fallback. The
	// rows span four hours, and they are in the newest-first order the producer
	// emits, so a fallback anchored on [0] lands at the NEWEST end.
	// REQ-933's WORK RUNS EIGHT HOURS PAST the newest created_at on purpose. An
	// extent taken from created_at alone would end the window at 12:00 and clip that
	// bar off the right edge while still listing the row, so the assertions below
	// would pass; naming the window in the summary is what makes the difference
	// visible.
	brokenRangePayload := `{
	  "now": "2026-08-18T21:00:00Z",
	  "rangeStart": "not-a-timestamp",
	  "rangeEnd": "2026-08-18T21:00:00Z",
	  "rows": [
	    {"id":"REQ-933","createdTime":"2026-08-18T12:00:00Z","claimedTime":"2026-08-18T12:10:00Z",
	     "completedTime":"2026-08-18T20:00:00Z","waitMinutes":10,"workMinutes":470,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-932","createdTime":"2026-08-18T10:00:00Z","claimedTime":"2026-08-18T10:10:00Z",
	     "completedTime":"2026-08-18T10:30:00Z","waitMinutes":10,"workMinutes":20,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-931","createdTime":"2026-08-18T08:00:00Z","claimedTime":"2026-08-18T08:10:00Z",
	     "completedTime":"2026-08-18T08:30:00Z","waitMinutes":10,"workMinutes":20,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false}
	  ]
	}`

	// A payload whose range is unreadable AND whose only row is still OPEN, with no
	// projection to push the bound past the now-line. timelineRowSegments draws that
	// row's wait to `now`, so an extent built from stored stamps alone stops eight
	// hours short of it and the live part of the bar is unreachable in every window.
	// (Found by Codex on the pull request.)
	openRowPayload := `{
	  "now": "2026-08-18T20:00:00Z",
	  "rangeStart": "not-a-timestamp",
	  "rangeEnd": "not-a-timestamp",
	  "rows": [
	    {"id":"REQ-951","createdTime":"2026-08-18T12:00:00Z","claimedTime":null,
	     "completedTime":null,"waitMinutes":480,"workMinutes":0,
	     "waitOpen":true,"workOpen":false,"hasWork":false,"anomaly":false}
	  ]
	}`

	// And a payload where no row carries a readable instant at all: the fallback
	// must decline rather than invent a window. Doubly unreachable from the producer
	// — timelineRange always returns real instants, and buildTimelineAggregate drops
	// a ticket whose created_at does not parse — but a payload-integrity failure has
	// to end in a legible message rather than a window built from NaNs, which is the
	// same posture timelineRowSegments takes for its own unreachable case.
	unreadablePayload := `{
	  "now": "2026-08-18T13:00:00Z",
	  "rangeStart": "not-a-timestamp",
	  "rangeEnd": "also-not-a-timestamp",
	  "rows": [
	    {"id":"REQ-941","createdTime":"not-a-timestamp","claimedTime":null,
	     "completedTime":null,"waitMinutes":0,"workMinutes":0,
	     "waitOpen":true,"workOpen":false,"hasWork":false,"anomaly":false}
	  ]
	}`

	probeDriver := `
function drawnRowIds() {
  var ids = [];
  (function walk(node) {
    (node.children || []).forEach(function (childNode) {
      var attributes = childNode.attributes || {};
      if (childNode.stubName === "g" && attributes["data-detail-id"]) { ids.push(attributes["data-detail-id"]); return; }
      walk(childNode);
    });
  })(timelineStubHosts["timeline-scroll"]);
  return ids;
}
renderTimelineView();
process.stdout.write(JSON.stringify({
  rowIds: drawnRowIds(),
  summary: timelineStubHosts["timeline-summary"].textContent
}));
`

	renderWith := func(payload string) (drawnIds []string, summary string) {
		t.Helper()
		javascriptProbe := timelineRenderDomStubPreamble +
			"var boardData = { timeline: " + payload + " };\n" +
			string(rendererFragment) +
			probeDriver
		probeOutput := runJavaScriptBehaviorProbe(t, "timeline fallback bounds", javascriptProbe)
		var result struct {
			RowIds  []string `json:"rowIds"`
			Summary string   `json:"summary"`
		}
		if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
			t.Fatalf("decode timeline fallback bounds behavior: %v (output %q)", decodeError, probeOutput)
		}
		return result.RowIds, result.Summary
	}

	drawnIds, summary := renderWith(brokenRangePayload)

	// SETUP, ASSERTED: if the payload did not take the fallback branch, everything
	// below passes for the wrong reason.
	if !strings.Contains(summary, "3 REQs in the window") {
		t.Fatalf("the fallback render summarised %q; want all three rows, which is what says the "+
			"bounds cover the whole matched set rather than one hour around the newest row", summary)
	}
	if len(drawnIds) != 3 {
		t.Fatalf("the fallback render drew %v, want all three rows — the old fallback bounded the "+
			"view to one hour around REQ-933, the newest, leaving the other two unreachable", drawnIds)
	}
	// The oldest row specifically. A fallback anchored on the newest capture leaves
	// exactly this one off the chart.
	oldestIsDrawn := false
	for _, drawnId := range drawnIds {
		if drawnId == "REQ-931" {
			oldestIsDrawn = true
		}
	}
	if !oldestIsDrawn {
		t.Errorf("the oldest row REQ-931 is not on the chart (drawn: %v); the fallback bounds have "+
			"to reach the earliest instant the matched set carries", drawnIds)
	}
	// EVERY instant the rows carry, not just created_at. REQ-933's work ends at
	// 20:00, eight hours past the newest capture; an extent taken from created_at
	// alone ends the window around 12:04 and clips that bar off the frame while
	// still listing its row.
	if !strings.Contains(summary, "→ 2026-08-18 20:") {
		t.Errorf("the fallback window is %q; it has to reach REQ-933's completion at 20:00, so the "+
			"extent must read claimed and completed instants and not created_at alone", summary)
	}

	// AN OPEN ROW'S NOW-LINE IS PART OF THE EXTENT. The window has to reach 20:00,
	// where the bar is drawn to, not stop at the 12:00 the row has stored.
	openRowIds, openRowSummary := renderWith(openRowPayload)
	if len(openRowIds) != 1 {
		t.Fatalf("the open-row fallback render drew %v, want the single fixture row", openRowIds)
	}
	if !strings.Contains(openRowSummary, "→ 2026-08-18 20:") {
		t.Errorf("the fallback window for a still-open row is %q; it has to reach the now-line at "+
			"20:00, which is where timelineRowSegments draws that row's wait to — bounds are what "+
			"every control clamps against, so anything short of it is unreachable", openRowSummary)
	}

	// And with nothing readable, the view says so instead of fabricating a window.
	_, unreadableSummary := renderWith(unreadablePayload)
	if !strings.Contains(unreadableSummary, "nothing to place on a timeline") {
		t.Errorf("with no readable range and no readable row instants the summary reads %q; want the "+
			"existing decline rather than an invented window", unreadableSummary)
	}
}

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
func TestJavaScriptBehaviorTimelineNoMatchStateRetiresTheToolbar(t *testing.T) {
	rendererFragment, readError := embeddedWebAssets.ReadFile("web/board-timeline.js")
	if readError != nil {
		t.Fatalf("read web/board-timeline.js: %v", readError)
	}

	timelinePayload := `{
	  "now": "2026-08-18T12:00:00Z",
	  "rangeStart": "2026-08-18T09:00:00Z",
	  "rangeEnd": "2026-08-18T13:00:00Z",
	  "rows": [
	    {"id":"REQ-921","createdTime":"2026-08-18T09:00:00Z","claimedTime":"2026-08-18T09:30:00Z",
	     "completedTime":"2026-08-18T10:00:00Z","waitMinutes":30,"workMinutes":30,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-922","createdTime":"2026-08-18T10:00:00Z","claimedTime":"2026-08-18T10:30:00Z",
	     "completedTime":"2026-08-18T11:00:00Z","waitMinutes":30,"workMinutes":30,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false}
	  ]
	}`

	// The stub's querySelectorAll has to answer the control selectors, so the
	// controls are real stub nodes the driver can inspect and press.
	probeDriver := `
var stubControls = ["period-prev", "period-day", "zoom-fit", "range-start"].map(function (name) {
  var control = makeStubNode(name === "range-start" ? "input" : "button");
  control.controlName = name;
  control.disabled = false;
  control.onclick = null;
  return control;
});
document.querySelectorAll = function (selector) {
  if (String(selector).indexOf(".timeline-periods button") !== -1) { return stubControls; }
  if (String(selector).indexOf("[data-timeline-period]") !== -1) { return []; }
  return [];
};
document.getElementById = function (nodeId) {
  if (nodeId === "timeline-zoom-fit") { return stubControls[2]; }
  return timelineStubHosts[nodeId] || null;
};

function controlState() {
  return stubControls.map(function (control) {
    return { name: control.controlName, wired: typeof control.onclick === "function", disabled: !!control.disabled };
  });
}

timelineStubVisibleIds = null;
renderTimelineView();
var matched = { controls: controlState(), summary: timelineStubHosts["timeline-summary"].textContent };

// Nothing matches. The early return fires.
["timeline-summary", "timeline-axis", "timeline-scroll", "timeline-readout",
 "timeline-table-body", "timeline-forecast", "timeline-excluded", "timeline-period-state"
].forEach(function (hostId) { timelineStubHosts[hostId] = makeStubNode("div"); });
timelineStubVisibleIds = [];
renderTimelineView();
var noMatch = { controls: controlState(), summary: timelineStubHosts["timeline-summary"].textContent };

// And back: a filter that matches again must restore every control.
["timeline-summary", "timeline-axis", "timeline-scroll", "timeline-readout",
 "timeline-table-body", "timeline-forecast", "timeline-excluded", "timeline-period-state"
].forEach(function (hostId) { timelineStubHosts[hostId] = makeStubNode("div"); });
timelineStubVisibleIds = null;
renderTimelineView();
var matchedAgain = { controls: controlState(), summary: timelineStubHosts["timeline-summary"].textContent };

process.stdout.write(JSON.stringify({ matched: matched, noMatch: noMatch, matchedAgain: matchedAgain }));
`

	javascriptProbe := timelineRenderDomStubPreamble +
		"var boardData = { timeline: " + timelinePayload + " };\n" +
		string(rendererFragment) +
		probeDriver
	probeOutput := runJavaScriptBehaviorProbe(t, "timeline no-match toolbar", javascriptProbe)

	type controlState struct {
		Name     string `json:"name"`
		Wired    bool   `json:"wired"`
		Disabled bool   `json:"disabled"`
	}
	type renderState struct {
		Controls []controlState `json:"controls"`
		Summary  string         `json:"summary"`
	}
	var toolbarResult struct {
		Matched      renderState `json:"matched"`
		NoMatch      renderState `json:"noMatch"`
		MatchedAgain renderState `json:"matchedAgain"`
	}
	if decodeError := json.Unmarshal(probeOutput, &toolbarResult); decodeError != nil {
		t.Fatalf("decode timeline no-match toolbar behavior: %v (output %q)", decodeError, probeOutput)
	}

	// SETUP, ASSERTED: without these the checks below are measuring nothing.
	if len(toolbarResult.Matched.Controls) == 0 {
		t.Fatal("the probe found no controls, so nothing below is measured")
	}
	if !strings.Contains(toolbarResult.NoMatch.Summary, "No REQ matches the current filters") {
		t.Fatalf("the second render did not take the no-match path (summary %q)",
			toolbarResult.NoMatch.Summary)
	}
	if !strings.Contains(toolbarResult.Matched.Summary, "REQs in the window") {
		t.Fatalf("the first render did not draw a chart (summary %q)", toolbarResult.Matched.Summary)
	}
	// The "still carries a handler" check below is vacuous for any control this stub
	// never wires, so at least one has to arrive wired or that half proves nothing.
	wiredAfterAMatchingRender := 0
	for _, control := range toolbarResult.Matched.Controls {
		if control.Wired {
			wiredAfterAMatchingRender++
		}
	}
	if wiredAfterAMatchingRender == 0 {
		t.Fatal("no control was wired by the matching render, so the handler half of this test " +
			"cannot fail; give the stub's getElementById a control the renderer wires")
	}

	for _, control := range toolbarResult.NoMatch.Controls {
		if control.Wired {
			t.Errorf("after the no-match render the %s control still carries a handler; it belongs "+
				"to the previous render, whose rows the filter excluded", control.Name)
		}
		if !control.Disabled {
			t.Errorf("after the no-match render the %s control is still pressable; a control that "+
				"cannot act must say so rather than doing nothing", control.Name)
		}
	}
	for _, control := range toolbarResult.MatchedAgain.Controls {
		if control.Disabled {
			t.Errorf("the %s control is still disabled after the filter matched again", control.Name)
		}
	}
}

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
func TestJavaScriptBehaviorTimelineSummaryCountsRowsDrawnAsBreaks(t *testing.T) {
	rendererFragment, readError := embeddedWebAssets.ReadFile("web/board-timeline.js")
	if readError != nil {
		t.Fatalf("read web/board-timeline.js: %v", readError)
	}

	timelinePayload := `{
	  "now": "2026-08-18T12:00:00Z",
	  "rangeStart": "2026-08-18T08:00:00Z",
	  "rangeEnd": "2026-08-18T13:00:00Z",
	  "rows": [
	    {"id":"REQ-911","createdTime":"2026-08-18T09:00:00Z","claimedTime":"2026-08-18T09:30:00Z",
	     "completedTime":"2026-08-18T10:00:00Z","waitMinutes":30,"workMinutes":30,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":true},
	    {"id":"REQ-912","createdTime":"2026-08-18T11:00:00Z","claimedTime":"2026-08-18T10:00:00Z",
	     "completedTime":"2026-08-18T11:30:00Z","waitMinutes":-60,"workMinutes":90,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-913","createdTime":"2026-08-18T09:00:00Z","claimedTime":"2026-08-18T11:00:00Z",
	     "completedTime":"2026-08-18T10:00:00Z","waitMinutes":120,"workMinutes":-60,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-914","createdTime":"2026-08-18T11:00:00Z","claimedTime":"2026-08-18T10:00:00Z",
	     "completedTime":"2026-08-18T09:00:00Z","waitMinutes":-60,"workMinutes":-60,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-915","createdTime":"2026-08-18T09:00:00Z","claimedTime":"2026-08-18T10:00:00Z",
	     "completedTime":"2026-08-18T11:00:00Z","waitMinutes":60,"workMinutes":60,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-916","createdTime":"2026-08-18T09:00:00Z","claimedTime":"2026-08-18T09:45:00Z",
	     "waitMinutes":45,"workMinutes":0,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":true}
	  ]
	}`

	probeDriver := `
function timelineSummaryWithFilter(visibleIds) {
  [
    "timeline-summary", "timeline-axis", "timeline-scroll", "timeline-readout",
    "timeline-table-body", "timeline-forecast", "timeline-excluded", "timeline-period-state"
  ].forEach(function (hostId) { timelineStubHosts[hostId] = makeStubNode("div"); });
  timelineStubVisibleIds = visibleIds;
  renderTimelineView();
  return timelineStubHosts["timeline-summary"].textContent || "";
}
process.stdout.write(JSON.stringify({
  unfiltered: timelineSummaryWithFilter(null),
  filtered: timelineSummaryWithFilter(["REQ-912", "REQ-915"]),
  anomalyOnly: timelineSummaryWithFilter(["REQ-911"]),
  unresolvedOnly: timelineSummaryWithFilter(["REQ-916"]),
  unresolvedAndAnomalyOnly: timelineSummaryWithFilter(["REQ-911", "REQ-916"]),
  reversedPair: timelineSummaryWithFilter(["REQ-912", "REQ-913"]),
  reversedWaitOnly: timelineSummaryWithFilter(["REQ-912"]),
  reversedWorkOnly: timelineSummaryWithFilter(["REQ-913"]),
  combinedCausesOnly: timelineSummaryWithFilter(["REQ-914"]),
  healthyOnly: timelineSummaryWithFilter(["REQ-915"])
}));
`

	javascriptProbe := timelineRenderDomStubPreamble +
		"var boardData = { timeline: " + timelinePayload + " };\n" +
		string(rendererFragment) +
		probeDriver
	probeOutput := runJavaScriptBehaviorProbe(t, "timeline summary break count", javascriptProbe)

	var summaries struct {
		Unfiltered               string `json:"unfiltered"`
		Filtered                 string `json:"filtered"`
		AnomalyOnly              string `json:"anomalyOnly"`
		UnresolvedOnly           string `json:"unresolvedOnly"`
		UnresolvedAndAnomalyOnly string `json:"unresolvedAndAnomalyOnly"`
		ReversedPair             string `json:"reversedPair"`
		ReversedWaitOnly         string `json:"reversedWaitOnly"`
		ReversedWorkOnly         string `json:"reversedWorkOnly"`
		CombinedCausesOnly       string `json:"combinedCausesOnly"`
		HealthyOnly              string `json:"healthyOnly"`
	}
	if decodeError := json.Unmarshal(probeOutput, &summaries); decodeError != nil {
		t.Fatalf("decode rendered timeline summaries: %v (output starts %q)",
			decodeError, string(probeOutput[:min(len(probeOutput), 400)]))
	}

	wantBreakClause := func(caseName string, summary string, wantCount int) {
		t.Helper()
		wantClause := fmt.Sprintf(". %d with broken stamps, drawn as breaks.", wantCount)
		if !strings.Contains(summary, wantClause) {
			t.Errorf("%s summary must contain %q, got %q", caseName, wantClause, summary)
		}
	}

	wantNoBreakClause := func(caseName string, summary string) {
		t.Helper()
		if strings.Contains(summary, "with broken stamps") {
			t.Errorf("%s summary must omit the break clause, got %q", caseName, summary)
		}
	}

	wantBreakClause("unfiltered", summaries.Unfiltered, 4)
	wantBreakClause("filtered", summaries.Filtered, 1)
	wantBreakClause("reversed wait and work", summaries.ReversedPair, 2)
	wantBreakClause("reversed wait only", summaries.ReversedWaitOnly, 1)
	wantBreakClause("reversed work only", summaries.ReversedWorkOnly, 1)
	wantBreakClause("combined causes", summaries.CombinedCausesOnly, 1)
	// A REQ that stopped with no resolvable end instant IS drawn as a break, so it
	// is counted.
	wantBreakClause("unresolved only", summaries.UnresolvedOnly, 1)
	// And a row that is merely flagged anomalous, with every span drawn, is NOT —
	// which is the whole of REQ-328's change to this clause. Both of these fail if
	// `row.anomaly ||` comes back: the first would report 2, the second 1.
	wantBreakClause("unresolved beside an anomalous-but-drawn row", summaries.UnresolvedAndAnomalyOnly, 1)
	wantNoBreakClause("anomaly only", summaries.AnomalyOnly)
	wantNoBreakClause("healthy only", summaries.HealthyOnly)
}

// TestJavaScriptBehaviorTimelineForecastLabelsAFilteredView drives the WHOLE
// renderTimelineView, not renderTimelineForecast alone, because the defect lived
// in the wiring: rows were filtered, projection never was, and the call site
// handed the forecast a filtered row list it then ignored. A probe that calls
// the forecast function directly cannot tell a correct call site from one that
// always says "unfiltered" — this one can.
func TestJavaScriptBehaviorTimelineForecastLabelsAFilteredView(t *testing.T) {
	rendererFragment, readError := embeddedWebAssets.ReadFile("web/board-timeline.js")
	if readError != nil {
		t.Fatalf("read web/board-timeline.js: %v", readError)
	}

	timelinePayload := `{
	  "now": "2026-08-18T12:00:00Z",
	  "rangeStart": "2026-08-18T09:00:00Z",
	  "rangeEnd": "2026-08-18T13:00:00Z",
	  "rows": [
	    {"id":"REQ-901","createdTime":"2026-08-18T10:00:00Z","claimedTime":"2026-08-18T10:30:00Z",
	     "completedTime":"2026-08-18T11:00:00Z","waitMinutes":30,"workMinutes":30,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-902","createdTime":"2026-08-18T10:00:00Z","claimedTime":"2026-08-18T10:30:00Z",
	     "completedTime":"2026-08-18T11:00:00Z","waitMinutes":30,"workMinutes":30,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-903","createdTime":"2026-08-18T11:00:00Z","claimedTime":null,
	     "completedTime":null,"waitMinutes":60,"workMinutes":0,
	     "waitOpen":true,"workOpen":false,"hasWork":false,"anomaly":false}
	  ],
	  "projection": {
	    "confident": true,
	    "chainStart": "2026-08-18T12:00:00Z",
	    "queueEnd": "2026-08-18T14:30:00Z",
	    "windowSamples": 60, "windowSize": 60, "minimumSamples": 5,
	    "normalSamples": 55, "normalMinutes": 40,
	    "trivialSamples": 5, "trivialMinutes": 10,
	    "rows": [{"id":"REQ-903","startTime":"2026-08-18T12:00:00Z","endTime":"2026-08-18T12:40:00Z"}],
	    "excluded": [{"id":"REQ-904","reason":"waiting on an external condition"}],
	    "queueEndSource": "median"
	  }
	}`

	// One render per filter state, each from a fresh stub, so nothing carries over.
	probeDriver := `
function renderWithFilter(visibleIds) {
  [
    "timeline-summary", "timeline-axis", "timeline-scroll", "timeline-readout",
    "timeline-table-body", "timeline-forecast", "timeline-excluded", "timeline-period-state"
  ].forEach(function (hostId) { timelineStubHosts[hostId] = makeStubNode("div"); });
  timelineStubVisibleIds = visibleIds;
  renderTimelineView();
  return {
    summary: timelineStubHosts["timeline-summary"].textContent || "",
    forecast: collectStubText(timelineStubHosts["timeline-forecast"]),
    excluded: collectStubText(timelineStubHosts["timeline-excluded"])
  };
}
function collectStubText(node) {
  var text = node.textContent || "";
  (node.children || []).forEach(function (child) { text += " " + collectStubText(child); });
  return text;
}
// A DRAINED queue, which is the only state that reaches the chainCount === 0
// branch. The whole board's projection is replaced rather than a second fixture
// added, because the contradiction is between two clauses of ONE paragraph and
// both have to be produced by one render.
function renderDrainedWithFilter(visibleIds) {
  boardData.timeline.projection.rows = [];
  return renderWithFilter(visibleIds);
}
process.stdout.write(JSON.stringify({
  unfiltered: renderWithFilter(null),
  filtered: renderWithFilter(["REQ-901"]),
  drainedUnfiltered: renderDrainedWithFilter(null),
  drainedFiltered: renderDrainedWithFilter(["REQ-901"])
}));
`

	javascriptProbe := timelineRenderDomStubPreamble +
		"var boardData = { timeline: " + timelinePayload + " };\n" +
		string(rendererFragment) +
		probeDriver
	probeOutput := runJavaScriptBehaviorProbe(t, "timeline forecast filter label", javascriptProbe)

	type renderedView struct {
		Summary  string `json:"summary"`
		Forecast string `json:"forecast"`
		Excluded string `json:"excluded"`
	}
	var rendered struct {
		Unfiltered        renderedView `json:"unfiltered"`
		Filtered          renderedView `json:"filtered"`
		DrainedUnfiltered renderedView `json:"drainedUnfiltered"`
		DrainedFiltered   renderedView `json:"drainedFiltered"`
	}
	if decodeError := json.Unmarshal(probeOutput, &rendered); decodeError != nil {
		t.Fatalf("decode rendered timeline views: %v (output starts %q)",
			decodeError, string(probeOutput[:min(len(probeOutput), 400)]))
	}

	// The fixture has to actually produce the disagreement, or the assertions
	// below pass against a view that was never filtered.
	if !strings.Contains(rendered.Unfiltered.Summary, "3 REQs") {
		t.Fatalf("the unfiltered render must show all three fixture rows, got summary %q", rendered.Unfiltered.Summary)
	}
	if !strings.Contains(rendered.Filtered.Summary, "1 REQ ") {
		t.Fatalf("the filtered render must show one row, got summary %q", rendered.Filtered.Summary)
	}

	if strings.Contains(rendered.Unfiltered.Forecast, "whole queue") {
		t.Errorf("an unfiltered view has nothing to disambiguate, so the forecast must carry no label; got %q",
			rendered.Unfiltered.Forecast)
	}
	if !strings.Contains(rendered.Filtered.Forecast, "whole queue") {
		t.Errorf("a filtered view forecasts the whole queue and must say so — this is the wiring, not the copy; got %q",
			rendered.Filtered.Forecast)
	}
	if !strings.Contains(rendered.Filtered.Excluded, "whole queue") {
		t.Errorf("the excluded list names REQ-904, which no visible row carries, and must say whose queue it lists; got %q",
			rendered.Filtered.Excluded)
	}
	// Both renders still carry the estimate: the label is added, not substituted.
	for viewName, view := range map[string]renderedView{
		"unfiltered": rendered.Unfiltered,
		"filtered":   rendered.Filtered,
	} {
		if !strings.Contains(view.Forecast, "Queue empties around") {
			t.Errorf("the %s forecast lost its estimate: %q", viewName, view.Forecast)
		}
	}

	// A DRAINED QUEUE, where the paragraph used to contradict itself inside one
	// sentence pair: "This covers the whole queue, not the rows shown." followed by
	// "Nothing left to schedule — every remaining REQ is listed below.", above a
	// single row, with the excluded paragraph under it naming a REQ that was not
	// listed anywhere.
	//
	// SETUP, ASSERTED: without this the two assertions below pass against a forecast
	// that never took the drained branch at all.
	if !strings.Contains(rendered.DrainedFiltered.Forecast, "Nothing left to schedule") {
		t.Fatalf("the drained fixture did not reach the nothing-left branch: %q",
			rendered.DrainedFiltered.Forecast)
	}
	if strings.Contains(rendered.DrainedFiltered.Forecast, "listed below") {
		t.Errorf("the forecast claims every remaining REQ is \"listed below\" while also saying it "+
			"covers the whole queue and not the rows shown, above one row: %q",
			rendered.DrainedFiltered.Forecast)
	}
	if !strings.Contains(rendered.DrainedFiltered.Forecast, "whole queue") {
		t.Errorf("the drained filtered forecast dropped the whole-queue note, which is the half "+
			"that is TRUE and is what the figures under it depend on: %q",
			rendered.DrainedFiltered.Forecast)
	}
	// Unfiltered, the sentence is accurate and must be left alone: the rows on
	// screen really are the whole queue there.
	if !strings.Contains(rendered.DrainedUnfiltered.Forecast, "listed below") {
		t.Errorf("with nothing filtered the rows ARE the whole queue, so the forecast should still "+
			"say so: %q", rendered.DrainedUnfiltered.Forecast)
	}
}

// TestJavaScriptBehaviorDurationsReserveMatchesTheSentenceDrawn pins the reserve to
// a fixed point rather than to two passes.
//
// Holding the remainder sentence's room narrows the last label row, which can hide
// labels the unreserved pass had placed. The count therefore rises between the pass
// the reserve was measured from and the pass whose result is drawn — and across a
// digit boundary the sentence gets wider too, so `+10 more …` is painted into room
// reserved for `+8 more …`, over the last placed label. That is the exact collision
// the planner exists to prevent, and it was reachable before this: at 26 candidates
// of one width the two passes hid 8 then 10.
//
// The invariant is stated as a relationship, not a number: whatever the final count
// is, the reserve the final pass held must have been measured from that same count.
// It sweeps a range of densities rather than pinning one fixture, because the
// boundary moves with the face and the plot width.
func TestJavaScriptBehaviorDurationsReserveMatchesTheSentenceDrawn(t *testing.T) {
	rendererFragment, readError := embeddedWebAssets.ReadFile("web/board-durations.js")
	if readError != nil {
		t.Fatalf("read web/board-durations.js: %v", readError)
	}

	probeDriver := `
var lastMeasuredSentence = null;
function measureRemainderWidth(sentenceText) {
  lastMeasuredSentence = sentenceText;
  return sentenceText.length * 7;   // monotone in the digit count, like a real face
}
var mismatches = [], remainderCases = 0, multiPassCases = 0;
for (var candidateCount = 12; candidateCount <= 60; candidateCount += 1) {
  for (var markStep = 20; markStep <= 45; markStep += 5) {
    lastMeasuredSentence = null;
    var candidates = [];
    for (var index = 0; index < candidateCount; index += 1) {
      candidates.push({
        markX: DURATIONS_MARGIN_LEFT + 20 + index * markStep,
        textWidth: 80,
        id: "REQ-" + index
      });
    }
    var unreserved = placeDurationsLabelBand(candidates, 0);
    var band = packDurationsLabelBand(candidates, measureRemainderWidth, "not labelled");
    if (band.hiddenCount === 0) { continue; }
    remainderCases += 1;
    if (band.hiddenCount > unreserved.hiddenCount) { multiPassCases += 1; }
    var sentenceDrawn = composeDurationsRemainderText(band.hiddenCount, "not labelled");
    if (lastMeasuredSentence !== sentenceDrawn) {
      mismatches.push({
        candidateCount: candidateCount, markStep: markStep,
        reservedFor: lastMeasuredSentence, drawn: sentenceDrawn
      });
    }
  }
}
process.stdout.write(JSON.stringify({
  mismatches: mismatches, remainderCases: remainderCases, multiPassCases: multiPassCases
}));
`

	javascriptProbe := string(rendererFragment) + probeDriver
	probeOutput := runJavaScriptBehaviorProbe(t, "durations remainder reserve", javascriptProbe)

	var swept struct {
		Mismatches []struct {
			CandidateCount int    `json:"candidateCount"`
			MarkStep       int    `json:"markStep"`
			ReservedFor    string `json:"reservedFor"`
			Drawn          string `json:"drawn"`
		} `json:"mismatches"`
		RemainderCases int `json:"remainderCases"`
		MultiPassCases int `json:"multiPassCases"`
	}
	if decodeError := json.Unmarshal(probeOutput, &swept); decodeError != nil {
		t.Fatalf("decode durations reserve sweep: %v (output starts %q)",
			decodeError, string(probeOutput[:min(len(probeOutput), 400)]))
	}

	// Guard the sweep itself: a fixture range that never produces a remainder, or
	// never needs more than one reserved pass, would satisfy the invariant vacuously.
	if swept.RemainderCases == 0 {
		t.Fatal("no swept density produced a remainder — the sweep proves nothing about the reserve")
	}
	if swept.MultiPassCases == 0 {
		t.Fatal("no swept density hid more under the reserve than without it — the sweep never reaches the case this pins")
	}

	for _, mismatch := range swept.Mismatches {
		t.Errorf("%d candidates at step %d reserved room for %q but will draw %q",
			mismatch.CandidateCount, mismatch.MarkStep, mismatch.ReservedFor, mismatch.Drawn)
	}
}
