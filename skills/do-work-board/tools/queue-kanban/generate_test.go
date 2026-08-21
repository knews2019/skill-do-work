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

// durationsRenderProbeDriver renders the fixture board and reports every drawn
// bar's x-interval, the slowest-day annotation's anchor x (the only mid-anchored
// text at the annotation baseline), and every Panel A mark centre in draw order.
const durationsRenderProbeDriver = `
renderDurationsView();
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
process.stdout.write(JSON.stringify({ bars: drawnBars, annotationXs: drawnAnnotationXs, markCxs: drawnMarkCxs }));
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
// annotation anchored inside it, and every Panel A mark at the x the Go-side
// label planner assumed. That last assertion is the drift pin: the renderer's
// axis domain and durations.go's durationLabelTimeRange must stay one
// definition, and it fails in whichever direction either side moves alone.
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

			// (4) Every Panel A mark at the x the Go-side label planner assumed.
			rangeStart, rangeEnd, hasRange := durationLabelTimeRange(aggregate.Samples)
			if !hasRange {
				t.Fatal("fixture produced no label time range")
			}
			if len(drawn.MarkCxs) != len(aggregate.Samples) {
				t.Fatalf("drew %d marks for %d samples", len(drawn.MarkCxs), len(aggregate.Samples))
			}
			for sampleIndex, sample := range aggregate.Samples {
				plannedMarkX := marginLeft + durationLabelPlotX(sample.CompletionTime, rangeStart, rangeEnd)
				if math.Abs(drawn.MarkCxs[sampleIndex]-plannedMarkX) > 0.06 {
					t.Errorf("%s mark drawn at x=%.2f but the label planner assumed %.2f — renderer and durations.go no longer share one axis domain",
						sample.RequestId, drawn.MarkCxs[sampleIndex], plannedMarkX)
				}
			}
		})
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
	javascriptProbe := timelineProbePreamble(t, "TIMELINE_MIN_SPAN_MS", "TIMELINE_DAY_MS", "TIMELINE_ROW_HEIGHT", "TIMELINE_NOW_JUMP_MARGIN_FRACTION") +
		rendererDeclarationLine(t, "web/board-timeline.js", "TIMELINE_PERIOD_LEVEL_NAMES") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineZoomedWindow(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelinePeriodStart(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineSteppedPeriodStart(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelinePeriodWindow(") + "\n" +
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
function anchorOf(movedWindow) {
  return (movedWindow.windowStartMs + movedWindow.windowEndMs) / 2;
}

var weekWindow = timelinePeriodWindow(anchorOf(fitted), "week", 0, boundStart, boundEnd);
var nextWeek = timelinePeriodWindow(anchorOf(weekWindow), "week", 1, boundStart, boundEnd);
var prevWeek = timelinePeriodWindow(anchorOf(nextWeek), "week", -1, boundStart, boundEnd);

// Held down against the end of the range: stepping has to stop, keeping the
// window on the data instead of running off it.
var atRangeEnd = weekWindow;
for (var step = 0; step < 60; step++) {
  atRangeEnd = timelinePeriodWindow(anchorOf(atRangeEnd), "week", 1, boundStart, boundEnd);
}
var pastRangeEnd = timelinePeriodWindow(anchorOf(atRangeEnd), "week", 1, boundStart, boundEnd);

var dayWindow = timelinePeriodWindow(anchorOf(fitted), "day", 0, boundStart, boundEnd);
var monthWindow = timelinePeriodWindow(anchorOf(fitted), "month", 0, boundStart, boundEnd);
var nextMonth = timelinePeriodWindow(anchorOf(monthWindow), "month", 1, boundStart, boundEnd);

// A free zoom through the pointer path's own transform: the level must stop
// reading as an exact week rather than keep claiming one.
var freelyZoomed = timelineZoomedWindow(
  weekWindow.windowStartMs, weekWindow.windowEndMs, 1.6, 0.5, boundStart, boundEnd);

// Three archived rows, then the still-open work — the 677-row board in miniature.
var rows = [
  { waitOpen: false, workOpen: false },
  { waitOpen: false, workOpen: false },
  { waitOpen: false, workOpen: false },
  { waitOpen: true, workOpen: false },
  { waitOpen: false, workOpen: true }
];
var scrollHostStub = { scrollTop: 0 };
var nowJump = timelineNowJump(nowMs, queueEndMs, rows, boundStart, boundEnd);
// The two assignments the Now button makes, run against a stub scroll host.
if (nowJump.scrollTop !== null) {
  scrollHostStub.scrollTop = nowJump.scrollTop;
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
  nowInsideWindow: nowMs >= nowJump.window.windowStartMs && nowMs <= nowJump.window.windowEndMs,
  queueEndInsideWindow: queueEndMs >= nowJump.window.windowStartMs && queueEndMs <= nowJump.window.windowEndMs,
  scrollTop: scrollHostStub.scrollTop,
  wantScrollTop: 3 * TIMELINE_ROW_HEIGHT
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

		NowWindowIso         string  `json:"nowWindowIso"`
		NowInsideWindow      bool    `json:"nowInsideWindow"`
		QueueEndInsideWindow bool    `json:"queueEndInsideWindow"`
		ScrollTop            float64 `json:"scrollTop"`
		WantScrollTop        float64 `json:"wantScrollTop"`
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
	if !periodResult.NowInsideWindow || !periodResult.QueueEndInsideWindow {
		t.Fatalf("the Now window %s does not cover both the now-line and the projected queue end", periodResult.NowWindowIso)
	}
	if periodResult.ScrollTop != periodResult.WantScrollTop {
		t.Fatalf("Now left the row list at scrollTop %.0f, want %.0f — the first still-open row",
			periodResult.ScrollTop, periodResult.WantScrollTop)
	}
}

// The axis is the one part of this view whose defect is pure text. The ticks are
// positioned correctly, so nothing about the drawing looks wrong while several
// labels read the same instant. REQ-227 wrote the minute as the literal ":00",
// and REQ-235's Now button made a sub-hour window the landing state of the
// view's most-used control: seven ticks, two distinct labels, measured live.
//
// This drives the formatter over every window the view can be left in and
// requires two things of each: no two ticks may render the same label, and every
// number in a label must be one the tick's own instant carries.
func TestJavaScriptBehaviorTimelineAxisLabelsNameTheirOwnInstant(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := timelineProbePreamble(t, "TIMELINE_MIN_SPAN_MS", "TIMELINE_DAY_MS", "TIMELINE_AXIS_TICK_COUNT") +
		rendererDeclarationLine(t, "web/board-timeline.js", "TIMELINE_YEAR_MS") + "\n" +
		rendererDeclarationLine(t, "web/board-timeline.js", "TIMELINE_MONTHS") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineFormatAxisTick(") + `
// The same evenly-spaced ticks renderAxis draws, formatted the same way.
function axisTicks(name, startMs, spanMs) {
  var ticks = [];
  for (var tickIndex = 0; tickIndex <= TIMELINE_AXIS_TICK_COUNT; tickIndex++) {
    var tickMs = startMs + (spanMs * tickIndex) / TIMELINE_AXIS_TICK_COUNT;
    var instant = new Date(tickMs);
    ticks.push({
      label: timelineFormatAxisTick(tickMs, spanMs),
      dayOfMonth: instant.getUTCDate(),
      hour: instant.getUTCHours(),
      minute: instant.getUTCMinutes(),
      year: instant.getUTCFullYear()
    });
  }
  return { name: name, ticks: ticks };
}

var mondayMs = Date.UTC(2026, 7, 17);      // 17 Aug 2026 is a Monday
process.stdout.write(JSON.stringify({
  windows: [
    // Where the Now button lands: the window covers the now-line and the
    // forecast's queue-empty instant, which on a healthy queue is well under an
    // hour, so the span settles on the view's floor and the start is wherever
    // "now" fell — 11:26, not the top of an hour.
    axisTicks("Now", Date.UTC(2026, 7, 18, 11, 26), TIMELINE_MIN_SPAN_MS),
    axisTicks("Day", Date.UTC(2026, 7, 18), TIMELINE_DAY_MS),
    // A free zoom between the period levels. Six ticks across four days sit 16h
    // apart, so a date-only label repeats itself twice over.
    axisTicks("free zoom, four days", Date.UTC(2026, 7, 15), 4 * TIMELINE_DAY_MS),
    axisTicks("Week", mondayMs, 7 * TIMELINE_DAY_MS),
    axisTicks("Month", Date.UTC(2026, 7, 1), Date.UTC(2026, 8, 1) - Date.UTC(2026, 7, 1)),
    axisTicks("Fit all", Date.UTC(2026, 3, 7), Date.UTC(2026, 7, 18) - Date.UTC(2026, 3, 7)),
    // Fit all is the whole capture history, and it only grows. Once it reaches
    // back past a year, one day-and-month comes round twice.
    axisTicks("Fit all across two years", Date.UTC(2025, 7, 18), 2 * TIMELINE_YEAR_MS)
  ]
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline axis labels", javascriptProbe)
	var axisResult struct {
		Windows []struct {
			Name  string `json:"name"`
			Ticks []struct {
				Label      string `json:"label"`
				DayOfMonth int    `json:"dayOfMonth"`
				Hour       int    `json:"hour"`
				Minute     int    `json:"minute"`
				Year       int    `json:"year"`
			} `json:"ticks"`
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
		"Now":                      axisLabelWithTime,
		"Day":                      axisLabelWithTime,
		"free zoom, four days":     axisLabelWithTime,
		"Week":                     axisLabelDateOnly,
		"Month":                    axisLabelDateOnly,
		"Fit all":                  axisLabelDateOnly,
		"Fit all across two years": axisLabelWithYear,
	}
	if len(axisResult.Windows) != len(wantAxisLabelShape) {
		t.Fatalf("the probe drove %d windows, want the %d named", len(axisResult.Windows), len(wantAxisLabelShape))
	}

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
	}
	if !sawStaleClaim {
		t.Errorf("verifyFindings carries no stale-claim finding: %+v", boardData.VerifyFindings)
	}
	// The anomaly must still exist as a finding — it is suppressed from the page,
	// not from verify. Otherwise this test would pass on a probe that stopped working.
	report := collectVerifyFindings(repoRoot, board, moment)
	if anomalies := findingsMentioning(report, verifyCategoryCompletionAnomaly); len(anomalies) == 0 {
		t.Error("the fixture produced no completion anomaly, so the suppression assertion proves nothing")
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
process.stdout.write(JSON.stringify({
  unfiltered: renderWithFilter(null),
  filtered: renderWithFilter(["REQ-901"])
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
		Unfiltered renderedView `json:"unfiltered"`
		Filtered   renderedView `json:"filtered"`
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
}
