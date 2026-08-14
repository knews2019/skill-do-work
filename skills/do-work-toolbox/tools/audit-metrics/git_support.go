package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Git plumbing shared by the subcommands. Every call goes through `git -C
// <repoRoot>` — never cwd-relative, so the answers stay pinned to the repo the
// caller named even when the tool runs from somewhere else (the worktree
// trap). Git exit codes are three-valued in general — 0, 1 = "no", anything
// else (usually 128) = git declining to answer. None of the queries here use
// exit 1 as an answer (the shallow probe answers on stdout), so every non-zero
// exit is git declining and is reported as an error with git's own stderr —
// never folded into a default.

// runGitCommand runs one git command against repoRoot and returns its stdout.
func runGitCommand(repoRoot string, arguments ...string) (string, error) {
	fullArguments := append([]string{"-C", repoRoot}, arguments...)
	command := exec.Command("git", fullArguments...)
	var standardOutput, standardError bytes.Buffer
	command.Stdout = &standardOutput
	command.Stderr = &standardError
	if runError := command.Run(); runError != nil {
		return "", fmt.Errorf("git %s: %w\n%s", strings.Join(arguments, " "), runError, standardError.String())
	}
	return standardOutput.String(), nil
}

// resolveRepositoryRoot canonicalizes --repo-root to the repository's toplevel
// via `git rev-parse --show-toplevel`. This is load-bearing, not cosmetic:
// `git log` prints toplevel-relative paths while `git ls-files` prints
// cwd-relative ones, so a --repo-root pointing INSIDE the repo would make the
// churn join silently empty. It also turns "not a git repo" into a clean
// reported error before any measuring starts.
func resolveRepositoryRoot(repoRootFlag string) (string, error) {
	topLevelOutput, topLevelError := runGitCommand(repoRootFlag, "rev-parse", "--show-toplevel")
	if topLevelError != nil {
		return "", topLevelError
	}
	return strings.TrimSpace(topLevelOutput), nil
}

// listTrackedFiles returns the repo-relative paths of every tracked file, via
// NUL-separated `ls-files -z` so no path shape can split a record.
func listTrackedFiles(repoRoot string) ([]string, error) {
	listOutput, listError := runGitCommand(repoRoot, "ls-files", "-z")
	if listError != nil {
		return nil, listError
	}
	var trackedPaths []string
	for _, trackedPath := range strings.Split(listOutput, "\x00") {
		if trackedPath != "" {
			trackedPaths = append(trackedPaths, trackedPath)
		}
	}
	return trackedPaths, nil
}

// isShallowRepository reports whether repoRoot is a shallow clone. The answer
// is on stdout ("true"/"false"); a non-zero exit is an error, never "false" —
// silently reading a broken probe as "not shallow" is exactly the truncation
// the churn report must surface.
func isShallowRepository(repoRoot string) (bool, error) {
	probeOutput, probeError := runGitCommand(repoRoot, "rev-parse", "--is-shallow-repository")
	if probeError != nil {
		return false, probeError
	}
	return strings.TrimSpace(probeOutput) == "true", nil
}

// applyPathExcludes drops every path matching one of the exclude prefixes.
// Prefix match on the repo-relative slash path: "skills/" excludes the tree,
// "CHANGELOG.md" excludes the file. The default exclude list is EMPTY — the
// caller owns it (a built-in list would go stale; see the audit action).
func applyPathExcludes(paths []string, excludePrefixes []string) []string {
	if len(excludePrefixes) == 0 {
		return paths
	}
	var keptPaths []string
	for _, candidatePath := range paths {
		if !pathMatchesAnyPrefix(candidatePath, excludePrefixes) {
			keptPaths = append(keptPaths, candidatePath)
		}
	}
	return keptPaths
}

// pathMatchesAnyPrefix reports whether candidatePath starts with any prefix.
func pathMatchesAnyPrefix(candidatePath string, excludePrefixes []string) bool {
	for _, excludePrefix := range excludePrefixes {
		if strings.HasPrefix(candidatePath, excludePrefix) {
			return true
		}
	}
	return false
}

// fileMeasurement is one tracked file's size numbers. Binary files (NUL byte in
// the first 8 KiB) keep Lines/Words at zero and are excluded from line/word
// totals and distributions — word counts on binaries are noise, not data.
type fileMeasurement struct {
	Path   string
	Lines  int
	Words  int
	Binary bool
}

// binarySniffLimit bounds the NUL-byte sniff, matching the "read a header, not
// the file" convention diff tools use.
const binarySniffLimit = 8192

// measureTrackedFile reads one tracked file and computes its line/word counts.
func measureTrackedFile(repoRoot string, relativePath string) (fileMeasurement, error) {
	fileBytes, readError := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(relativePath)))
	if readError != nil {
		return fileMeasurement{}, readError
	}
	measurement := fileMeasurement{Path: relativePath}
	sniffWindow := fileBytes
	if len(sniffWindow) > binarySniffLimit {
		sniffWindow = sniffWindow[:binarySniffLimit]
	}
	if bytes.IndexByte(sniffWindow, 0) >= 0 {
		measurement.Binary = true
		return measurement, nil
	}
	measurement.Lines = countTextLines(fileBytes)
	measurement.Words = len(bytes.Fields(fileBytes))
	return measurement, nil
}

// countTextLines counts newline-terminated lines plus a final unterminated one.
func countTextLines(fileBytes []byte) int {
	lineCount := bytes.Count(fileBytes, []byte{'\n'})
	if len(fileBytes) > 0 && fileBytes[len(fileBytes)-1] != '\n' {
		lineCount++
	}
	return lineCount
}
