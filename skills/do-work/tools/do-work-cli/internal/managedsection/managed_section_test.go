package managedsection

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

const sectionFileBytes = "# >>> do-work:recipes >>>\nrun-kanban:\n    echo kanban\nalias rk := run-kanban\n# <<< do-work:recipes <<<\n"

func writeManagedFixture(t *testing.T, path string, contents []byte, mode os.FileMode) string {
	t.Helper()
	if err := os.WriteFile(path, contents, mode); err != nil {
		t.Fatalf("write fixture %s: %v", path, err)
	}
	if err := os.Chmod(path, goModeFromUnix(mode)); err != nil {
		t.Fatalf("chmod fixture %s: %v", path, err)
	}
	return path
}

func goModeFromUnix(mode os.FileMode) os.FileMode {
	goMode := mode.Perm()
	if mode&0o4000 != 0 {
		goMode |= os.ModeSetuid
	}
	if mode&0o2000 != 0 {
		goMode |= os.ModeSetgid
	}
	if mode&0o1000 != 0 {
		goMode |= os.ModeSticky
	}
	return goMode
}

func newSectionFile(t *testing.T, directory string) string {
	t.Helper()
	return writeManagedFixture(t, filepath.Join(directory, "section.just"), []byte(sectionFileBytes), 0o644)
}

func fileMode(t *testing.T, path string) os.FileMode {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	mode := info.Mode().Perm()
	if info.Mode()&os.ModeSetuid != 0 {
		mode |= 0o4000
	}
	if info.Mode()&os.ModeSetgid != 0 {
		mode |= 0o2000
	}
	if info.Mode()&os.ModeSticky != 0 {
		mode |= 0o1000
	}
	return mode
}

// Python's bytes.splitlines ends a line on a bare CR, and only on CR, LF and CRLF. A Go port
// that splits on LF alone mislocates every marker span in an old-Mac-terminated file, which
// is the single most likely silent byte-level regression in this port.
func TestLineSplittingMatchesPythonBytesSplitlines(t *testing.T) {
	tests := []struct {
		name          string
		data          string
		expectedLines []string
	}{
		{name: "lone carriage return ends a line", data: "a\rb\nc\r\nd", expectedLines: []string{"a\r", "b\n", "c\r\n", "d"}},
		{name: "vertical tab and form feed do not end a line", data: "a\vb\fc\x1cd\x1ee\x85f", expectedLines: []string{"a\vb\fc\x1cd\x1ee\x85f"}},
		{name: "trailing terminator produces no empty final line", data: "a\n", expectedLines: []string{"a\n"}},
		{name: "empty input produces no lines", data: "", expectedLines: nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lines, offsets := splitLinesKeepingEnds([]byte(test.data))
			if len(lines) != len(test.expectedLines) {
				t.Fatalf("line count = %d, want %d (lines %q)", len(lines), len(test.expectedLines), lines)
			}
			cursor := 0
			for index, expected := range test.expectedLines {
				if string(lines[index]) != expected {
					t.Errorf("line %d = %q, want %q", index, lines[index], expected)
				}
				if offsets[index] != cursor {
					t.Errorf("offset %d = %d, want %d", index, offsets[index], cursor)
				}
				cursor += len(lines[index])
			}
		})
	}
}

// rstrip(b"\r\n") strips a trailing SET of CR and LF bytes, not the "\r\n" suffix, so a
// marker terminated by CRLF, by a bare CR, or by LF all compare equal to the same constant.
func TestMarkerBodyStripsEveryTrailingCarriageReturnAndNewline(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected string
	}{
		{name: "mixed trailing terminators", line: "line\n\r\n\r", expected: "line"},
		{name: "carriage return only", line: "line\r", expected: "line"},
		{name: "no terminator", line: "line", expected: "line"},
		{name: "interior carriage return survives", line: "li\rne\n", expected: "li\rne"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := string(lineBody([]byte(test.line))); got != test.expected {
				t.Errorf("lineBody(%q) = %q, want %q", test.line, got, test.expected)
			}
		})
	}
}

func TestReplaceSectionPreservesSurroundingBytesAndTargetMode(t *testing.T) {
	tests := []struct {
		name          string
		targetBytes   string
		targetMode    os.FileMode
		expectedBytes string
	}{
		{
			name:          "lone carriage return terminators survive outside the span",
			targetBytes:   "custom:\r# >>> do-work:recipes >>>\rOLD\r# <<< do-work:recipes <<<\rtail\r",
			targetMode:    0o644,
			expectedBytes: "custom:\r" + sectionFileBytes + "tail\r",
		},
		{
			name:          "CRLF terminators survive outside the span",
			targetBytes:   "custom:\r\n# >>> do-work:recipes >>>\r\nOLD\r\n# <<< do-work:recipes <<<\r\ntail\r\n",
			targetMode:    0o644,
			expectedBytes: "custom:\r\n" + sectionFileBytes + "tail\r\n",
		},
		{
			name:          "an embedded NUL outside the section is transparent",
			targetBytes:   "prefix\x00byte\n# >>> do-work:recipes >>>\nOLD\n# <<< do-work:recipes <<<\n",
			targetMode:    0o640,
			expectedBytes: "prefix\x00byte\n" + sectionFileBytes,
		},
		{
			name:          "target 02644 mode is preserved on replace",
			targetBytes:   "before\n# >>> do-work:recipes >>>\nOLD\n# <<< do-work:recipes <<<\nafter\n",
			targetMode:    0o2644,
			expectedBytes: "before\n" + sectionFileBytes + "after\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			directory := t.TempDir()
			sectionPath := newSectionFile(t, directory)
			targetPath := writeManagedFixture(t, filepath.Join(directory, "target"), []byte(test.targetBytes), test.targetMode)
			outcome, err := ReplaceSection(ReplaceRequest{TargetPath: targetPath, SectionFilePath: sectionPath})
			if err != nil {
				t.Fatalf("ReplaceSection: %v", err)
			}
			if outcome.Kind != ChangeKindModified {
				t.Errorf("change kind = %q, want %q", outcome.Kind, ChangeKindModified)
			}
			actual, readErr := os.ReadFile(targetPath)
			if readErr != nil {
				t.Fatalf("read target: %v", readErr)
			}
			if !bytes.Equal(actual, []byte(test.expectedBytes)) {
				t.Errorf("target bytes = %q, want %q", actual, test.expectedBytes)
			}
			if mode := fileMode(t, targetPath); mode != test.targetMode {
				t.Errorf("target mode = %o, want %o", mode, test.targetMode)
			}
		})
	}
}

func TestAppendChoosesTheSeparatorFromTheTargetTail(t *testing.T) {
	tests := []struct {
		name          string
		targetBytes   string
		expectedBytes string
	}{
		{name: "empty target takes no separator", targetBytes: "", expectedBytes: sectionFileBytes},
		{name: "target ending in a blank line takes no separator", targetBytes: "x\n\n", expectedBytes: "x\n\n" + sectionFileBytes},
		{name: "target ending in one newline takes one newline", targetBytes: "x\n", expectedBytes: "x\n\n" + sectionFileBytes},
		{name: "unterminated target takes two newlines", targetBytes: "x", expectedBytes: "x\n\n" + sectionFileBytes},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			directory := t.TempDir()
			sectionPath := newSectionFile(t, directory)
			targetPath := writeManagedFixture(t, filepath.Join(directory, "target"), []byte(test.targetBytes), 0o644)
			if _, err := ReplaceSection(ReplaceRequest{TargetPath: targetPath, SectionFilePath: sectionPath}); err != nil {
				t.Fatalf("ReplaceSection: %v", err)
			}
			actual, readErr := os.ReadFile(targetPath)
			if readErr != nil {
				t.Fatalf("read target: %v", readErr)
			}
			if !bytes.Equal(actual, []byte(test.expectedBytes)) {
				t.Errorf("target bytes = %q, want %q", actual, test.expectedBytes)
			}
		})
	}
}

func TestCreateFromTemplatePreservesTheTemplateMode(t *testing.T) {
	directory := t.TempDir()
	sectionPath := newSectionFile(t, directory)
	templatePath := writeManagedFixture(t, filepath.Join(directory, "template.just"), []byte(sectionFileBytes), 0o750)
	targetPath := filepath.Join(directory, "created")

	outcome, err := ReplaceSection(ReplaceRequest{
		TargetPath:       targetPath,
		SectionFilePath:  sectionPath,
		TemplateFilePath: templatePath,
	})
	if err != nil {
		t.Fatalf("ReplaceSection: %v", err)
	}
	if outcome.Kind != ChangeKindCreated {
		t.Errorf("change kind = %q, want %q", outcome.Kind, ChangeKindCreated)
	}
	actual, readErr := os.ReadFile(targetPath)
	if readErr != nil {
		t.Fatalf("read created target: %v", readErr)
	}
	if !bytes.Equal(actual, []byte(sectionFileBytes)) {
		t.Errorf("created bytes = %q, want %q", actual, sectionFileBytes)
	}
	if mode := fileMode(t, targetPath); mode != 0o750 {
		t.Errorf("created mode = %o, want 750", mode)
	}
}

// A second run must not write at all: the installer's byte-idempotence and the unchanged
// mtime the shipped suite depends on both rest on the conditional write.
func TestSecondRunLeavesTheTargetUntouched(t *testing.T) {
	directory := t.TempDir()
	sectionPath := newSectionFile(t, directory)
	targetPath := writeManagedFixture(t, filepath.Join(directory, "target"), []byte(sectionFileBytes), 0o600)

	beforeInfo, err := os.Stat(targetPath)
	if err != nil {
		t.Fatalf("stat before: %v", err)
	}
	outcome, replaceErr := ReplaceSection(ReplaceRequest{TargetPath: targetPath, SectionFilePath: sectionPath})
	if replaceErr != nil {
		t.Fatalf("ReplaceSection: %v", replaceErr)
	}
	if outcome.Changed {
		t.Errorf("outcome reported a change for an already-matching target")
	}
	afterInfo, statErr := os.Stat(targetPath)
	if statErr != nil {
		t.Fatalf("stat after: %v", statErr)
	}
	if !afterInfo.ModTime().Equal(beforeInfo.ModTime()) {
		t.Errorf("modification time changed for an already-matching target")
	}
	if afterInfo.Mode().Perm() != 0o600 {
		t.Errorf("mode = %o, want 600", afterInfo.Mode().Perm())
	}
}

// The BOM asymmetry is deliberate (REQ-173): the Just scanner strips a BOM from line 0's
// classification view, while marker matching never does. A BOM before the begin marker
// therefore reads as one end marker and no begin marker.
func TestRefusalsCarryTheirCurrentDiagnosticVerbatim(t *testing.T) {
	tests := []struct {
		name             string
		targetBytes      string
		expectedCode     string
		expectedEvidence string
	}{
		{
			name:             "BOM before the begin marker breaks the pair",
			targetBytes:      "\xef\xbb\xbf# >>> do-work:recipes >>>\nOLD\n# <<< do-work:recipes <<<\n",
			expectedCode:     FailureInvalidInput,
			expectedEvidence: "target must contain exactly one begin marker and one end marker",
		},
		{
			name:             "duplicated begin marker",
			targetBytes:      "# >>> do-work:recipes >>>\n# >>> do-work:recipes >>>\nOLD\n# <<< do-work:recipes <<<\n",
			expectedCode:     FailureInvalidInput,
			expectedEvidence: "target must contain exactly one begin marker and one end marker",
		},
		{
			name:             "begin marker without an end marker",
			targetBytes:      "# >>> do-work:recipes >>>\nOLD\n",
			expectedCode:     FailureInvalidInput,
			expectedEvidence: "target must contain exactly one begin marker and one end marker",
		},
		{
			name:             "reversed markers",
			targetBytes:      "# <<< do-work:recipes <<<\nOLD\n# >>> do-work:recipes >>>\n",
			expectedCode:     FailureInvalidInput,
			expectedEvidence: "target has reversed or nested managed markers",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			directory := t.TempDir()
			sectionPath := newSectionFile(t, directory)
			targetPath := writeManagedFixture(t, filepath.Join(directory, "target"), []byte(test.targetBytes), 0o644)
			_, err := ReplaceSection(ReplaceRequest{TargetPath: targetPath, SectionFilePath: sectionPath})
			failure := asSectionFailure(t, err)
			if failure.Code != test.expectedCode {
				t.Errorf("failure code = %q, want %q", failure.Code, test.expectedCode)
			}
			if failure.Message != test.expectedEvidence {
				t.Errorf("failure message = %q, want %q", failure.Message, test.expectedEvidence)
			}
			actual, readErr := os.ReadFile(targetPath)
			if readErr != nil {
				t.Fatalf("read target: %v", readErr)
			}
			if !bytes.Equal(actual, []byte(test.targetBytes)) {
				t.Errorf("a refused replacement wrote the target")
			}
		})
	}
}

func TestSymlinkedTargetIsRefusedWithoutFollowingTheLink(t *testing.T) {
	directory := t.TempDir()
	sectionPath := newSectionFile(t, directory)
	realPath := writeManagedFixture(t, filepath.Join(directory, "real"), []byte("untouched\n"), 0o644)
	linkPath := filepath.Join(directory, "link")
	if err := os.Symlink(realPath, linkPath); err != nil {
		t.Fatalf("symlink: %v", err)
	}

	_, err := ReplaceSection(ReplaceRequest{TargetPath: linkPath, SectionFilePath: sectionPath})
	failure := asSectionFailure(t, err)
	if failure.Message != "target must not be a symlink: "+linkPath {
		t.Errorf("failure message = %q", failure.Message)
	}
	linked, readErr := os.ReadFile(realPath)
	if readErr != nil {
		t.Fatalf("read link destination: %v", readErr)
	}
	if string(linked) != "untouched\n" {
		t.Errorf("the symlink destination was written: %q", linked)
	}
}

// A dangling symlink counts as an existing target, so it is refused as a symlink rather than
// silently created through the link.
func TestDanglingSymlinkTargetIsRefusedRatherThanCreated(t *testing.T) {
	directory := t.TempDir()
	sectionPath := newSectionFile(t, directory)
	linkPath := filepath.Join(directory, "dangling")
	if err := os.Symlink(filepath.Join(directory, "nowhere-at-all"), linkPath); err != nil {
		t.Fatalf("symlink: %v", err)
	}

	_, err := ReplaceSection(ReplaceRequest{
		TargetPath:       linkPath,
		SectionFilePath:  sectionPath,
		TemplateFilePath: sectionPath,
	})
	failure := asSectionFailure(t, err)
	if failure.Message != "target must not be a symlink: "+linkPath {
		t.Errorf("failure message = %q", failure.Message)
	}
}

func TestReservedRecipeCollisionNamesEverySortedCollidingDefinition(t *testing.T) {
	tests := []struct {
		name             string
		targetBytes      string
		expectedEvidence string
	}{
		{
			name:             "a recipe and an alias collide",
			targetBytes:      "alias rk := something\nrun-kanban:\n    echo dup\n",
			expectedEvidence: "target defines reserved Just recipe or alias outside managed section: rk, run-kanban",
		},
		{
			name:             "a BOM and CRLF collision is still detected",
			targetBytes:      "\xef\xbb\xbfrun-kanban:\r\n    echo custom\r\n",
			expectedEvidence: "target defines reserved Just recipe or alias outside managed section: run-kanban",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			directory := t.TempDir()
			sectionPath := newSectionFile(t, directory)
			targetPath := writeManagedFixture(t, filepath.Join(directory, "target"), []byte(test.targetBytes), 0o644)
			_, err := ReplaceSection(ReplaceRequest{
				TargetPath:             targetPath,
				SectionFilePath:        sectionPath,
				RejectRecipeCollisions: true,
			})
			failure := asSectionFailure(t, err)
			if failure.Code != FailureReservedRecipeCollision {
				t.Errorf("failure code = %q, want %q", failure.Code, FailureReservedRecipeCollision)
			}
			if failure.Message != test.expectedEvidence {
				t.Errorf("failure message = %q, want %q", failure.Message, test.expectedEvidence)
			}
			actual, readErr := os.ReadFile(targetPath)
			if readErr != nil {
				t.Fatalf("read target: %v", readErr)
			}
			if !bytes.Equal(actual, []byte(test.targetBytes)) {
				t.Errorf("a rejected collision wrote the target")
			}
		})
	}
}

// The reserved names live inside the managed span, so they must not collide with themselves.
func TestReservedNamesInsideTheManagedSpanAreNotCollisions(t *testing.T) {
	directory := t.TempDir()
	sectionPath := newSectionFile(t, directory)
	targetPath := writeManagedFixture(t, filepath.Join(directory, "target"),
		[]byte("other:\n    echo x\n"+sectionFileBytes), 0o644)

	if _, err := ReplaceSection(ReplaceRequest{
		TargetPath:             targetPath,
		SectionFilePath:        sectionPath,
		RejectRecipeCollisions: true,
	}); err != nil {
		t.Fatalf("ReplaceSection refused a self-collision: %v", err)
	}
}

func TestCustomMarkersRetargetTheManagedSpan(t *testing.T) {
	directory := t.TempDir()
	const beginMarker = "<!-- >>> do-work:communication-style >>> -->"
	const endMarker = "<!-- <<< do-work:communication-style <<< -->"
	sectionPath := writeManagedFixture(t, filepath.Join(directory, "section.md"),
		[]byte(beginMarker+"\nNEW\n"+endMarker+"\n"), 0o644)
	targetPath := writeManagedFixture(t, filepath.Join(directory, "CLAUDE.md"),
		[]byte("a\n"+beginMarker+"\nOLD\n"+endMarker+"\nb\n"), 0o644)

	if _, err := ReplaceSection(ReplaceRequest{
		TargetPath:      targetPath,
		SectionFilePath: sectionPath,
		BeginMarker:     beginMarker,
		EndMarker:       endMarker,
	}); err != nil {
		t.Fatalf("ReplaceSection: %v", err)
	}
	actual, readErr := os.ReadFile(targetPath)
	if readErr != nil {
		t.Fatalf("read target: %v", readErr)
	}
	expected := "a\n" + beginMarker + "\nNEW\n" + endMarker + "\nb\n"
	if string(actual) != expected {
		t.Errorf("target bytes = %q, want %q", actual, expected)
	}
}

func TestInvalidInputsAreRefusedBeforeAnyWrite(t *testing.T) {
	directory := t.TempDir()
	sectionPath := newSectionFile(t, directory)
	tests := []struct {
		name             string
		request          ReplaceRequest
		expectedEvidence string
	}{
		{
			name:             "an absent target without a template",
			request:          ReplaceRequest{TargetPath: filepath.Join(directory, "absent"), SectionFilePath: sectionPath},
			expectedEvidence: "target is absent; --template-file is required",
		},
		{
			name: "one-sided marker override",
			request: ReplaceRequest{
				TargetPath:      filepath.Join(directory, "absent"),
				SectionFilePath: sectionPath,
				BeginMarker:     "B",
			},
			expectedEvidence: "--begin-marker and --end-marker must be supplied together",
		},
		{
			name: "identical marker override",
			request: ReplaceRequest{
				TargetPath:      filepath.Join(directory, "absent"),
				SectionFilePath: sectionPath,
				BeginMarker:     "SAME",
				EndMarker:       "SAME",
			},
			expectedEvidence: "--begin-marker and --end-marker must differ",
		},
		{
			name: "a marker override carrying a newline",
			request: ReplaceRequest{
				TargetPath:      filepath.Join(directory, "absent"),
				SectionFilePath: sectionPath,
				BeginMarker:     "B\nB",
				EndMarker:       "E",
			},
			expectedEvidence: "--begin-marker must be one non-empty line",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := ReplaceSection(test.request)
			failure := asSectionFailure(t, err)
			if failure.Message != test.expectedEvidence {
				t.Errorf("failure message = %q, want %q", failure.Message, test.expectedEvidence)
			}
			if _, statErr := os.Lstat(test.request.TargetPath); statErr == nil {
				t.Errorf("a refused request created %s", test.request.TargetPath)
			}
		})
	}
}

func TestSectionFileMustBeExactlyOneNewlineTerminatedSection(t *testing.T) {
	tests := []struct {
		name         string
		sectionBytes string
	}{
		{name: "content before the begin marker", sectionBytes: "junk\n" + sectionFileBytes},
		{name: "content after the end marker", sectionBytes: sectionFileBytes + "junk\n"},
		{name: "no managed markers at all", sectionBytes: "just text\n"},
		{name: "not newline terminated", sectionBytes: "# >>> do-work:recipes >>>\n# <<< do-work:recipes <<<"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			directory := t.TempDir()
			sectionPath := writeManagedFixture(t, filepath.Join(directory, "section"), []byte(test.sectionBytes), 0o644)
			targetPath := writeManagedFixture(t, filepath.Join(directory, "target"), []byte("x\n"), 0o644)
			_, err := ReplaceSection(ReplaceRequest{TargetPath: targetPath, SectionFilePath: sectionPath})
			if err == nil {
				t.Fatalf("a malformed section file was accepted")
			}
		})
	}
}

// The temporary file must be created in the target's own directory so the rename is atomic
// on the target's filesystem, and it must not survive a completed replacement.
func TestAtomicReplaceLeavesNoTemporaryBesideTheTarget(t *testing.T) {
	directory := t.TempDir()
	sectionPath := newSectionFile(t, directory)
	targetDirectory := filepath.Join(directory, "nested")
	if err := os.Mkdir(targetDirectory, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	targetPath := writeManagedFixture(t, filepath.Join(targetDirectory, "justfile"), []byte("x\n"), 0o644)

	if _, err := ReplaceSection(ReplaceRequest{TargetPath: targetPath, SectionFilePath: sectionPath}); err != nil {
		t.Fatalf("ReplaceSection: %v", err)
	}
	entries, readErr := os.ReadDir(targetDirectory)
	if readErr != nil {
		t.Fatalf("read target directory: %v", readErr)
	}
	if len(entries) != 1 || entries[0].Name() != "justfile" {
		names := make([]string, 0, len(entries))
		for _, entry := range entries {
			names = append(names, entry.Name())
		}
		t.Errorf("target directory holds %v, want only the target", names)
	}
}

func asSectionFailure(t *testing.T, err error) *SectionFailure {
	t.Helper()
	if err == nil {
		t.Fatalf("expected a section failure, got nil")
	}
	failure, ok := err.(*SectionFailure)
	if !ok {
		t.Fatalf("error %v is not a *SectionFailure", err)
	}
	return failure
}
