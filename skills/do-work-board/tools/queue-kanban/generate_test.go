package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
		"web/board-filters.js",
		"web/board-testing.js",
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

	probeCommand := exec.Command(nodePath, "-e", javascriptProbe)
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
	javascriptProbe := `
Date.now = function () { return Date.parse("2026-08-15T12:00:00Z"); };
var boardData = { calendar: [
  { id: "REQ-recent", completionTime: "2026-08-15T06:00:00Z" },
  { id: "REQ-old", completionTime: "2026-08-13T06:00:00Z" }
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
		t.Fatalf("recentlyDoneIds(24) = %#v, want only REQ-recent", result.RecentIds)
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
    { id: "REQ-502", completionTime: "2026-08-15T06:00:00Z" },
    { id: "REQ-501", completionTime: "2026-08-13T06:00:00Z" }
  ]
};
var requestsById = boardData.requests;
var userRequestsById = boardData.userRequests;
var viewState = { windowHours: 24 };
var filterState = { searchText: "", domain: "", status: "", userRequestActivity: "active" };
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
