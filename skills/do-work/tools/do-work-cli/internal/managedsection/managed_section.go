// Package managedsection replaces do-work's managed marker span inside a consumer-owned
// text file without disturbing a single byte outside that span.
//
// It is a byte-exact port of the embedded Python that tools/replace-text-section.sh used to
// carry. Three of its rules are easy to "fix" into a regression:
//
//   - Lines end on a bare CR as well as on LF and CRLF, because Python's bytes.splitlines
//     does. A file written with old-Mac terminators has its span located correctly today.
//   - A marker's comparison body drops every trailing CR and LF byte, so the same constant
//     matches a marker terminated by LF, CRLF or a bare CR.
//   - The BOM asymmetry is deliberate (REQ-173): the Just scanner strips a BOM from line
//     zero's classification view, marker matching never does. A BOM before the begin marker
//     therefore reads as an unpaired end marker, and a shipped fixture depends on that.
package managedsection

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// The default pair owns the managed Just recipe section. BeginMarker/EndMarker retarget the
// replacer at any other single-line pair, such as the HTML-comment markers a consumer
// project's Markdown agent instructions use.
const (
	DefaultBeginMarker = "# >>> do-work:recipes >>>"
	DefaultEndMarker   = "# <<< do-work:recipes <<<"
)

const (
	ChangeKindCreated  = "created"
	ChangeKindModified = "modified"
)

// FailureReservedRecipeCollision is a refusal a consumer resolves by renaming a recipe;
// every other failure is a malformed input the caller must correct. The command layer maps
// the first to a refusal (exit 1) and the second to a failure (exit 2).
const (
	FailureInvalidInput            = "SECTION-INVALID-INPUT"
	FailureReservedRecipeCollision = "SECTION-RESERVED-RECIPE-COLLISION"
)

// SectionFailure carries the diagnostic text the shell tool printed verbatim, because three
// behavioural suites assert on those exact phrases.
type SectionFailure struct {
	Code    string
	Message string
}

func (failure *SectionFailure) Error() string { return failure.Message }

func invalidInput(format string, arguments ...any) *SectionFailure {
	return &SectionFailure{Code: FailureInvalidInput, Message: fmt.Sprintf(format, arguments...)}
}

type ReplaceRequest struct {
	TargetPath             string
	SectionFilePath        string
	TemplateFilePath       string
	RejectRecipeCollisions bool
	BeginMarker            string
	EndMarker              string
}

type ReplaceOutcome struct {
	Changed bool
	Kind    string
}

type markerPair struct {
	begin []byte
	end   []byte
}

// ReplaceSection creates, appends or replaces the managed span in one atomic rename, and
// writes nothing at all when the result already equals the target byte for byte.
func ReplaceSection(request ReplaceRequest) (ReplaceOutcome, error) {
	markers, err := resolveMarkers(request.BeginMarker, request.EndMarker)
	if err != nil {
		return ReplaceOutcome{}, err
	}
	if request.TargetPath == "" || request.SectionFilePath == "" {
		return ReplaceOutcome{}, invalidInput("--target and --section-file are both required")
	}

	sectionData, err := readRegularFile(request.SectionFilePath, "section file")
	if err != nil {
		return ReplaceOutcome{}, err
	}
	sectionSpan, err := findMarkerSpan(sectionData, "section file", markers, true)
	if err != nil {
		return ReplaceOutcome{}, err
	}
	if sectionSpan == nil || !bytes.HasSuffix(sectionData, []byte("\n")) {
		return ReplaceOutcome{}, invalidInput("section file must be one newline-terminated managed section")
	}

	// A dangling symlink counts as an existing target, so Lstat rather than Stat decides
	// creation; the symlink refusal below then catches it instead of writing through it.
	if _, lstatErr := os.Lstat(request.TargetPath); lstatErr != nil {
		return createFromTemplate(request, sectionData, markers)
	}

	if isSymlink(request.TargetPath) {
		return ReplaceOutcome{}, invalidInput("target must not be a symlink: %s", request.TargetPath)
	}
	targetData, err := readRegularFile(request.TargetPath, "target")
	if err != nil {
		return ReplaceOutcome{}, err
	}
	targetMode, err := permissionsOf(request.TargetPath, "target")
	if err != nil {
		return ReplaceOutcome{}, err
	}
	targetSpan, err := findMarkerSpan(targetData, "target", markers, false)
	if err != nil {
		return ReplaceOutcome{}, err
	}

	if request.RejectRecipeCollisions {
		if err := rejectReservedRecipeCollisions(sectionData, targetData, targetSpan); err != nil {
			return ReplaceOutcome{}, err
		}
	}

	replacementData := buildReplacement(targetData, sectionData, targetSpan)
	if bytes.Equal(replacementData, targetData) {
		return ReplaceOutcome{Changed: false}, nil
	}
	if err := atomicReplace(request.TargetPath, replacementData, targetMode); err != nil {
		return ReplaceOutcome{}, err
	}
	return ReplaceOutcome{Changed: true, Kind: ChangeKindModified}, nil
}

func createFromTemplate(request ReplaceRequest, sectionData []byte, markers markerPair) (ReplaceOutcome, error) {
	if request.TemplateFilePath == "" {
		return ReplaceOutcome{}, invalidInput("target is absent; --template-file is required")
	}
	templateData, err := readRegularFile(request.TemplateFilePath, "template file")
	if err != nil {
		return ReplaceOutcome{}, err
	}
	templateSpan, err := findMarkerSpan(templateData, "template file", markers, false)
	if err != nil {
		return ReplaceOutcome{}, err
	}
	if templateSpan == nil || !bytes.Equal(templateData[templateSpan.start:templateSpan.end], sectionData) {
		return ReplaceOutcome{}, invalidInput("template file must contain the supplied managed section exactly once")
	}
	templateMode, err := permissionsOf(request.TemplateFilePath, "template file")
	if err != nil {
		return ReplaceOutcome{}, err
	}
	if err := atomicReplace(request.TargetPath, templateData, templateMode); err != nil {
		return ReplaceOutcome{}, err
	}
	return ReplaceOutcome{Changed: true, Kind: ChangeKindCreated}, nil
}

// buildReplacement replaces the located span, or appends the section with the separator the
// target's own tail chooses. All four separator cases are behaviour a consumer's justfile
// depends on, so they stay enumerated rather than normalised.
func buildReplacement(targetData, sectionData []byte, targetSpan *byteSpan) []byte {
	if targetSpan != nil {
		replacement := make([]byte, 0, targetSpan.start+len(sectionData)+len(targetData)-targetSpan.end)
		replacement = append(replacement, targetData[:targetSpan.start]...)
		replacement = append(replacement, sectionData...)
		return append(replacement, targetData[targetSpan.end:]...)
	}
	var separator []byte
	switch {
	case len(targetData) == 0:
		separator = nil
	case bytes.HasSuffix(targetData, []byte("\n\n")):
		separator = nil
	case bytes.HasSuffix(targetData, []byte("\n")):
		separator = []byte("\n")
	default:
		separator = []byte("\n\n")
	}
	replacement := make([]byte, 0, len(targetData)+len(separator)+len(sectionData))
	replacement = append(replacement, targetData...)
	replacement = append(replacement, separator...)
	return append(replacement, sectionData...)
}

func rejectReservedRecipeCollisions(sectionData, targetData []byte, targetSpan *byteSpan) error {
	reservedNames := JustDefinitionNames(sectionData)
	if len(reservedNames) == 0 {
		return invalidInput("section file defines no Just recipes or aliases for collision validation")
	}
	unmanagedData := targetData
	if targetSpan != nil {
		unmanagedData = make([]byte, 0, len(targetData)-(targetSpan.end-targetSpan.start))
		unmanagedData = append(unmanagedData, targetData[:targetSpan.start]...)
		unmanagedData = append(unmanagedData, targetData[targetSpan.end:]...)
	}
	collisionNames := make([]string, 0)
	for name := range JustDefinitionNames(unmanagedData) {
		if _, reserved := reservedNames[name]; reserved {
			collisionNames = append(collisionNames, name)
		}
	}
	if len(collisionNames) == 0 {
		return nil
	}
	sort.Strings(collisionNames)
	return &SectionFailure{
		Code: FailureReservedRecipeCollision,
		Message: "target defines reserved Just recipe or alias outside managed section: " +
			strings.Join(collisionNames, ", "),
	}
}

func resolveMarkers(beginMarker, endMarker string) (markerPair, error) {
	if (beginMarker == "") != (endMarker == "") {
		return markerPair{}, invalidInput("--begin-marker and --end-marker must be supplied together")
	}
	if beginMarker == "" {
		return markerPair{begin: []byte(DefaultBeginMarker), end: []byte(DefaultEndMarker)}, nil
	}
	for _, named := range []struct {
		flag  string
		value string
	}{{"--begin-marker", beginMarker}, {"--end-marker", endMarker}} {
		if strings.ContainsAny(named.value, "\n\r") {
			return markerPair{}, invalidInput("%s must be one non-empty line", named.flag)
		}
	}
	if beginMarker == endMarker {
		return markerPair{}, invalidInput("--begin-marker and --end-marker must differ")
	}
	return markerPair{begin: []byte(beginMarker), end: []byte(endMarker)}, nil
}

type byteSpan struct {
	start int
	end   int
}

// findMarkerSpan reports exact byte offsets, and the end marker's own terminator is inside
// the span. A file with neither marker returns nil, which means "no managed section here",
// not an error.
func findMarkerSpan(data []byte, label string, markers markerPair, requireSectionOnly bool) (*byteSpan, error) {
	lines, offsets := splitLinesKeepingEnds(data)
	beginIndexes := make([]int, 0, 1)
	endIndexes := make([]int, 0, 1)
	for index, line := range lines {
		body := lineBody(line)
		if bytes.Equal(body, markers.begin) {
			beginIndexes = append(beginIndexes, index)
		}
		if bytes.Equal(body, markers.end) {
			endIndexes = append(endIndexes, index)
		}
	}
	if len(beginIndexes) == 0 && len(endIndexes) == 0 {
		return nil, nil
	}
	if len(beginIndexes) != 1 || len(endIndexes) != 1 {
		return nil, invalidInput("%s must contain exactly one begin marker and one end marker", label)
	}
	beginIndex, endIndex := beginIndexes[0], endIndexes[0]
	if beginIndex >= endIndex {
		return nil, invalidInput("%s has reversed or nested managed markers", label)
	}
	if requireSectionOnly && (beginIndex != 0 || endIndex != len(lines)-1) {
		return nil, invalidInput("%s must contain only the complete managed section", label)
	}
	return &byteSpan{start: offsets[beginIndex], end: offsets[endIndex] + len(lines[endIndex])}, nil
}

// splitLinesKeepingEnds reproduces Python's bytes.splitlines(keepends=True): a line ends on
// LF, on CRLF, or on a bare CR, and on nothing else. Vertical tab, form feed and the Unicode
// line separators are ordinary content bytes here.
func splitLinesKeepingEnds(data []byte) ([][]byte, []int) {
	lines := make([][]byte, 0, bytes.Count(data, []byte("\n"))+1)
	offsets := make([]int, 0, cap(lines))
	lineStart := 0
	for index := 0; index < len(data); index++ {
		var lineEnd int
		switch data[index] {
		case '\n':
			lineEnd = index + 1
		case '\r':
			lineEnd = index + 1
			if lineEnd < len(data) && data[lineEnd] == '\n' {
				lineEnd++
			}
		default:
			continue
		}
		lines = append(lines, data[lineStart:lineEnd])
		offsets = append(offsets, lineStart)
		index = lineEnd - 1
		lineStart = lineEnd
	}
	if lineStart < len(data) {
		lines = append(lines, data[lineStart:])
		offsets = append(offsets, lineStart)
	}
	return lines, offsets
}

// lineBody drops a trailing SET of CR and LF bytes, matching Python's rstrip(b"\r\n"). A
// TrimSuffix of "\r\n" is a different function and would mis-compare a CR-terminated marker.
func lineBody(line []byte) []byte {
	end := len(line)
	for end > 0 && (line[end-1] == '\r' || line[end-1] == '\n') {
		end--
	}
	return line[:end]
}

func readRegularFile(path, label string) ([]byte, error) {
	if isSymlink(path) {
		return nil, invalidInput("%s must not be a symlink: %s", label, path)
	}
	info, err := os.Stat(path)
	if err != nil {
		return nil, invalidInput("cannot stat %s %s: %v", label, path, err)
	}
	if !info.Mode().IsRegular() {
		return nil, invalidInput("%s must be a regular file: %s", label, path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, invalidInput("cannot read %s %s: %v", label, path, err)
	}
	return data, nil
}

func permissionsOf(path, label string) (os.FileMode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, invalidInput("cannot stat %s %s: %v", label, path, err)
	}
	return info.Mode().Perm() | (info.Mode() & (os.ModeSetuid | os.ModeSetgid | os.ModeSticky)), nil
}

func isSymlink(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.Mode()&os.ModeSymlink != 0
}

// atomicReplace stages the new bytes in the TARGET's own directory so the rename cannot
// cross a filesystem boundary, then renames over the target. The directory fsync afterwards
// is best effort: the visible target has already changed, so a filesystem that cannot fsync
// a directory must not turn a completed replacement into a reported failure.
func atomicReplace(path string, content []byte, mode os.FileMode) error {
	parentDirectory := filepath.Dir(mustAbsolute(path))
	if info, err := os.Stat(parentDirectory); err != nil || !info.IsDir() {
		return invalidInput("target parent directory does not exist: %s", parentDirectory)
	}
	temporaryFile, err := os.CreateTemp(parentDirectory, "."+filepath.Base(path)+".*.tmp")
	if err != nil {
		return invalidInput("atomic replacement failed for %s: %v", path, err)
	}
	temporaryPath := temporaryFile.Name()
	writeErr := func() error {
		if _, err := temporaryFile.Write(content); err != nil {
			return err
		}
		if err := temporaryFile.Sync(); err != nil {
			return err
		}
		return temporaryFile.Chmod(mode)
	}()
	closeErr := temporaryFile.Close()
	if writeErr == nil {
		writeErr = closeErr
	}
	if writeErr != nil {
		_ = os.Remove(temporaryPath)
		return invalidInput("atomic replacement failed for %s: %v", path, writeErr)
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		_ = os.Remove(temporaryPath)
		return invalidInput("atomic replacement failed for %s: %v", path, err)
	}
	if directory, err := os.Open(parentDirectory); err == nil {
		_ = directory.Sync()
		_ = directory.Close()
	}
	return nil
}

func mustAbsolute(path string) string {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return absolute
}
