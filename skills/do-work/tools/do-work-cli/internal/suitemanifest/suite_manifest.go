// Package suitemanifest validates an extracted do-work suite archive before any byte of it
// reaches a consumer's project. It is a port of tools/validate-suite-manifest.sh and, like
// that script, it stops at the first violation: the manifest describes what the installer is
// about to copy, so the first thing wrong with it is the thing that must be fixed.
package suitemanifest

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// The four modules the suite ships, and the only destinations each may claim. The set is
// closed by design: an archive that maps a module anywhere else is not a do-work suite.
var requiredModuleDestinations = map[string]string{
	"skills/do-work":           ".claude/skills/do-work",
	"skills/do-work-board":     ".claude/skills/do-work-board",
	"skills/do-work-knowledge": ".claude/skills/do-work-knowledge",
	"skills/do-work-toolbox":   ".claude/skills/do-work-toolbox",
}

var semanticVersionPattern = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)

// currentVersionPattern matches the one line in actions/version.md that names the version.
// The whole remainder of the line is the version, matching the shell validator; the updater's
// looser leading-digits parse is deliberately not reused here.
var currentVersionPattern = regexp.MustCompile(`^\*\*Current version\*\*:[ \t]*(.*)$`)

// ModuleMapping is one validated manifest row: what to copy and where it goes.
type ModuleMapping struct {
	Source      string
	Destination string
}

// ValidationResult is what a valid suite looks like to the installer.
type ValidationResult struct {
	SuiteVersion string
	Modules      []ModuleMapping
}

// SummaryLine reproduces the shell validator's success line verbatim, because the updater
// and the behavioural suites both read it.
func (result ValidationResult) SummaryLine() string {
	return fmt.Sprintf("suite manifest valid: v%s (%d modules)", result.SuiteVersion, len(result.Modules))
}

// ValidateSuite checks an extracted archive root and returns its version and module plan.
// The error text carries the shell validator's phrasing verbatim.
func ValidateSuite(archiveRoot string) (ValidationResult, error) {
	if info, err := os.Stat(archiveRoot); err != nil || !info.IsDir() {
		return ValidationResult{}, fmt.Errorf("root does not exist: %s", archiveRoot)
	}
	resolvedRoot, err := filepath.EvalSymlinks(archiveRoot)
	if err != nil {
		return ValidationResult{}, fmt.Errorf("root does not exist: %s", archiveRoot)
	}
	resolvedRoot, err = filepath.Abs(resolvedRoot)
	if err != nil {
		return ValidationResult{}, fmt.Errorf("root does not exist: %s", archiveRoot)
	}

	suiteVersion, err := readSingleLineVersionFile(filepath.Join(resolvedRoot, "VERSION"),
		"VERSION", "VERSION must be a regular file in the suite root")
	if err != nil {
		return ValidationResult{}, err
	}

	manifestPath := filepath.Join(resolvedRoot, "suite", "modules.tsv")
	if !isRegularNonSymlink(manifestPath) {
		return ValidationResult{}, errors.New("suite/modules.tsv must be a regular file in the suite root")
	}
	modules, err := readManifestRows(resolvedRoot, manifestPath)
	if err != nil {
		return ValidationResult{}, err
	}

	coreVersion, err := readSingleLineVersionFile(filepath.Join(resolvedRoot, "skills", "do-work", "VERSION"),
		"skills/do-work/VERSION", "skills/do-work/VERSION must be a regular file")
	if err != nil {
		return ValidationResult{}, err
	}
	if coreVersion != suiteVersion {
		return ValidationResult{}, fmt.Errorf("skills/do-work/VERSION mismatch (expected %s, found %s)", suiteVersion, coreVersion)
	}

	if err := validateActionVersion(resolvedRoot, suiteVersion); err != nil {
		return ValidationResult{}, err
	}
	return ValidationResult{SuiteVersion: suiteVersion, Modules: modules}, nil
}

func readManifestRows(resolvedRoot, manifestPath string) ([]ModuleMapping, error) {
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, errors.New("suite/modules.tsv must be a regular file in the suite root")
	}
	manifestLines := splitOnNewlinesOnly(manifestData)
	if len(manifestLines) == 0 {
		return nil, errors.New("suite/modules.tsv is empty")
	}
	if manifestLines[0] != "source\tdestination" {
		return nil, errors.New("manifest header must be exactly: source<TAB>destination")
	}

	modules := make([]ModuleMapping, 0, len(requiredModuleDestinations))
	seenSources := map[string]struct{}{}
	seenDestinations := map[string]struct{}{}
	for offset, manifestLine := range manifestLines[1:] {
		lineNumber := offset + 2
		if manifestLine == "" {
			return nil, fmt.Errorf("line %d is blank", lineNumber)
		}
		if strings.Contains(manifestLine, "\r") {
			return nil, fmt.Errorf("line %d contains a carriage return", lineNumber)
		}
		source, destination, hasTab := strings.Cut(manifestLine, "\t")
		if !hasTab {
			return nil, fmt.Errorf("line %d must contain one tab", lineNumber)
		}
		if strings.Contains(destination, "\t") {
			return nil, fmt.Errorf("line %d has an unknown extra column", lineNumber)
		}
		if source == "" || destination == "" {
			return nil, fmt.Errorf("line %d has an empty column", lineNumber)
		}
		if traversesDirectories(source) {
			return nil, fmt.Errorf("line %d source traverses directories", lineNumber)
		}
		if traversesDirectories(destination) {
			return nil, fmt.Errorf("line %d destination traverses directories", lineNumber)
		}
		if strings.HasPrefix(source, "/") {
			return nil, fmt.Errorf("line %d source must be relative", lineNumber)
		}
		if strings.HasPrefix(destination, "/") {
			return nil, fmt.Errorf("line %d destination must be relative", lineNumber)
		}
		if _, duplicate := seenSources[source]; duplicate {
			return nil, fmt.Errorf("line %d duplicates source %s", lineNumber, source)
		}
		if _, duplicate := seenDestinations[destination]; duplicate {
			return nil, fmt.Errorf("line %d duplicates destination %s", lineNumber, destination)
		}
		seenSources[source] = struct{}{}
		seenDestinations[destination] = struct{}{}

		expectedDestination, required := requiredModuleDestinations[source]
		if !required {
			return nil, fmt.Errorf("line %d declares unexpected source %s", lineNumber, source)
		}
		if destination != expectedDestination {
			return nil, fmt.Errorf("line %d maps %s outside its required destination", lineNumber, source)
		}
		if err := validateModuleTree(resolvedRoot, source); err != nil {
			return nil, err
		}
		modules = append(modules, ModuleMapping{Source: source, Destination: destination})
	}
	if len(modules) != len(requiredModuleDestinations) {
		return nil, fmt.Errorf("manifest must contain exactly four modules (found %d)", len(modules))
	}
	return modules, nil
}

// validateModuleTree resolves the module root physically as well as textually, so a symlink
// inside the archive cannot redirect a module out of the suite.
func validateModuleTree(resolvedRoot, source string) error {
	moduleRoot := filepath.Join(resolvedRoot, source)
	info, err := os.Lstat(moduleRoot)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("%s must be a real directory, not a missing path or symlink", source)
	}
	physicalRoot, err := filepath.EvalSymlinks(moduleRoot)
	if err != nil {
		return fmt.Errorf("%s must be a real directory, not a missing path or symlink", source)
	}
	if !strings.HasPrefix(physicalRoot, filepath.Join(resolvedRoot, "skills")+string(filepath.Separator)) {
		return fmt.Errorf("%s resolves outside the suite root", source)
	}
	skillFile := filepath.Join(moduleRoot, "SKILL.md")
	skillInfo, err := os.Lstat(skillFile)
	if err != nil || !skillInfo.Mode().IsRegular() || skillInfo.Size() == 0 {
		return fmt.Errorf("%s/SKILL.md must be a non-empty regular file", source)
	}
	return nil
}

func validateActionVersion(resolvedRoot, suiteVersion string) error {
	actionDirectory := filepath.Join(resolvedRoot, "skills", "do-work", "actions")
	directoryInfo, err := os.Lstat(actionDirectory)
	if err != nil || !directoryInfo.IsDir() {
		return errors.New("skills/do-work/actions must be a real directory")
	}
	actionVersionFile := filepath.Join(actionDirectory, "version.md")
	if !isRegularNonSymlink(actionVersionFile) {
		return errors.New("skills/do-work/actions/version.md must be a regular file")
	}
	actionVersion, markerCount, err := ReadActionVersion(actionVersionFile)
	if err != nil {
		return errors.New("skills/do-work/actions/version.md must be a regular file")
	}
	if markerCount != 1 {
		return errors.New("skills/do-work/actions/version.md must contain exactly one Current version line")
	}
	if !semanticVersionPattern.MatchString(actionVersion) {
		return errors.New("skills/do-work/actions/version.md Current version must be a plain semantic version (X.Y.Z)")
	}
	if actionVersion != suiteVersion {
		return fmt.Errorf("skills/do-work/actions/version.md mismatch (expected %s, found %s)", suiteVersion, actionVersion)
	}
	return nil
}

// ReadActionVersion returns the version an actions/version.md declares plus how many
// Current-version lines it carries, so callers can enforce "exactly one" themselves.
func ReadActionVersion(path string) (string, int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", 0, err
	}
	version := ""
	markerCount := 0
	for _, line := range splitOnNewlinesOnly(data) {
		match := currentVersionPattern.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		markerCount++
		if markerCount == 1 {
			version = match[1]
		}
	}
	return version, markerCount, nil
}

// splitOnNewlinesOnly splits at LF and keeps any CR as content, matching how the shell
// validator's `IFS= read -r` loop saw a CRLF-terminated manifest row. A bufio.Scanner would
// strip the CR and silently accept a manifest the shell rejects.
func splitOnNewlinesOnly(data []byte) []string {
	if len(data) == 0 {
		return nil
	}
	text := string(data)
	lines := strings.Split(text, "\n")
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// ReadSuiteVersion reads a VERSION file's first line without imposing the validator's
// stricter single-line rule, matching `sed -n '1p'`.
func ReadSuiteVersion(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	firstLine, _, _ := strings.Cut(string(data), "\n")
	return firstLine, nil
}

// IsSemanticVersion reports whether a version string is a plain X.Y.Z.
func IsSemanticVersion(version string) bool {
	return semanticVersionPattern.MatchString(version)
}

func readSingleLineVersionFile(path, label, missingMessage string) (string, error) {
	if !isRegularNonSymlink(path) {
		return "", errors.New(missingMessage)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", errors.New(missingMessage)
	}
	firstLine, _, _ := strings.Cut(string(data), "\n")
	if !semanticVersionPattern.MatchString(firstLine) {
		return "", fmt.Errorf("%s must be a plain semantic version (X.Y.Z)", label)
	}
	if string(data) != firstLine+"\n" {
		return "", fmt.Errorf("%s must contain exactly one newline-terminated line", label)
	}
	return firstLine, nil
}

// traversesDirectories rejects `.` and `..` as whole path segments, matching the shell
// validator's `/$path/` case patterns.
func traversesDirectories(path string) bool {
	for _, segment := range strings.Split(path, "/") {
		if segment == "." || segment == ".." {
			return true
		}
	}
	return false
}

func isRegularNonSymlink(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.Mode().IsRegular()
}
