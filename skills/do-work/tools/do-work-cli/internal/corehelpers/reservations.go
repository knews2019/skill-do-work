package corehelpers

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

var reservationNamePattern = regexp.MustCompile(`^REQ-(\d{6})$`)
var requestFilePattern = regexp.MustCompile(`^REQ-(\d+)(?:-|\.md)`)
var beforeReservationRemoval = func(string) {}

// CleanupReservations runs the fail-soft reservation cleanup operation used by the
// public helper and SessionStart hook. Hook callers consume the typed changes directly.
func CleanupReservations(repositoryRoot string) resultmodel.CommandResult {
	return cleanupReservations(repositoryRoot, false)
}

func handleCleanupReservations(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	filtered, dryRun, dryRunError := extractDryRun(arguments)
	if dryRunError != nil || len(filtered) != 0 {
		return usageResult(CommandCleanupReservations, "cleanup-req-reservations accepts no options")
	}
	return cleanupReservations(executionContext.RepositoryRoot, dryRun)
}

func cleanupReservations(repositoryRoot string, dryRun bool) resultmodel.CommandResult {
	doWorkRoot := filepath.Join(repositoryRoot, "do-work")
	if info, err := os.Lstat(doWorkRoot); err == nil && (info.Mode()&os.ModeSymlink != 0 || !info.IsDir()) {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeRefused, Findings: []resultmodel.CommandFinding{helperFinding("RESERVATION-ROOT-UNSAFE", resultmodel.SeverityError, []string{"do-work"}, "do-work root is not a real directory", resultmodel.FixabilityRefused, "cleanup refuses link traversal", nil, nil)}}
	}
	root := filepath.Join(doWorkRoot, ".req-reservations")
	rootInfo, err := os.Lstat(root)
	if os.IsNotExist(err) {
		return successResult(nil, nil)
	}
	if err != nil {
		return usageResult(CommandCleanupReservations, err.Error())
	}
	if rootInfo.Mode()&os.ModeSymlink != 0 || !rootInfo.IsDir() {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeRefused, Findings: []resultmodel.CommandFinding{helperFinding("RESERVATION-ROOT-UNSAFE", resultmodel.SeverityError, []string{"do-work/.req-reservations"}, "reservation root is not a real directory", resultmodel.FixabilityRefused, "cleanup refuses link traversal", nil, nil)}}
	}
	rootHandle, openError := os.OpenRoot(root)
	if openError != nil {
		return usageResult(CommandCleanupReservations, openError.Error())
	}
	defer rootHandle.Close()
	openedInfo, statError := rootHandle.Stat(".")
	if statError != nil || !os.SameFile(rootInfo, openedInfo) {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeRefused, Findings: []resultmodel.CommandFinding{helperFinding("RESERVATION-ROOT-RACED", resultmodel.SeverityError, []string{"do-work/.req-reservations"}, "reservation root identity changed before inspection", resultmodel.FixabilityRefused, "cleanup has no authority over the replacement directory", nil, nil)}}
	}
	claimed, committedAuthority := claimedRequestNumbers(repositoryRoot)
	if !committedAuthority {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Findings: []resultmodel.CommandFinding{helperFinding(
			"RESERVATION-GIT-AUTHORITY-UNAVAILABLE", resultmodel.SeverityWarning,
			[]string{"do-work/.req-reservations"},
			"committed request authority could not be established with Git",
			resultmodel.FixabilityRefused,
			"all reservation markers are preserved when Git evidence is unavailable",
			[]string{"git", "rev-parse", "--verify", "HEAD"}, nil,
		)}}
	}
	entries, err := fs.ReadDir(rootHandle.FS(), ".")
	if err != nil {
		return usageResult(CommandCleanupReservations, err.Error())
	}
	changes := []resultmodel.RecordedChange{}
	findings := []resultmodel.CommandFinding{}
	now := time.Now()
	for _, entry := range entries {
		match := reservationNamePattern.FindStringSubmatch(entry.Name())
		if match == nil {
			findings = append(findings, helperFinding("RESERVATION-MALFORMED", resultmodel.SeverityWarning, []string{filepath.ToSlash(filepath.Join("do-work/.req-reservations", entry.Name()))}, "marker name is not canonical", resultmodel.FixabilityRefused, "malformed entries are preserved", nil, nil))
			continue
		}
		markerPath := filepath.Join(root, entry.Name())
		firstInfo, err := rootHandle.Lstat(entry.Name())
		if err != nil || !firstInfo.Mode().IsRegular() {
			continue
		}
		number, _ := strconv.Atoi(match[1])
		eligible := claimed[number] || now.Sub(firstInfo.ModTime()) >= 48*time.Hour
		if !eligible {
			continue
		}
		secondInfo, err := rootHandle.Lstat(entry.Name())
		if err != nil || !os.SameFile(firstInfo, secondInfo) || !secondInfo.Mode().IsRegular() {
			findings = append(findings, helperFinding("RESERVATION-RACED", resultmodel.SeverityWarning, []string{filepath.ToSlash(filepath.Join("do-work/.req-reservations", entry.Name()))}, "marker identity changed before removal", resultmodel.FixabilityRefused, "the current object is not the inspected marker", nil, nil))
			continue
		}
		if !dryRun {
			beforeReservationRemoval(markerPath)
			finalInfo, finalError := rootHandle.Lstat(entry.Name())
			finalClaims, finalAuthority := claimedRequestNumbers(repositoryRoot)
			finalClaimed := finalClaims[number]
			finalEligible := finalAuthority && finalError == nil && finalInfo.Mode().IsRegular() && os.SameFile(secondInfo, finalInfo) &&
				(finalClaimed || time.Since(finalInfo.ModTime()) >= 48*time.Hour)
			if !finalEligible {
				continue
			}
			if err := rootHandle.Remove(entry.Name()); err != nil {
				findings = append(findings, helperFinding("RESERVATION-REMOVE-FAILED", resultmodel.SeverityWarning, []string{markerPath}, err.Error(), resultmodel.FixabilityManual, "cleanup is fail-soft", nil, nil))
				continue
			}
		}
		reason := "stale for at least 48 hours"
		if claimed[number] {
			reason = "matching request exists"
		}
		if dryRun {
			reason = "would remove: " + reason
		}
		changes = append(changes, resultmodel.RecordedChange{Path: filepath.ToSlash(filepath.Join("do-work/.req-reservations", entry.Name())), Kind: "deleted", Detail: reason})
	}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Changes: changes, Findings: findings}
}

func claimedRequestNumbers(repositoryRoot string) (map[int]bool, bool) {
	numbers := map[int]bool{}
	tracked, err := gitOutput(repositoryRoot, "ls-tree", "-r", "--name-only", "-z", "HEAD", "--", "do-work/queue", "do-work/working", "do-work/archive")
	if err == nil {
		for _, path := range strings.Split(string(tracked), "\x00") {
			match := requestFilePattern.FindStringSubmatch(filepath.Base(path))
			if match != nil {
				number, _ := strconv.Atoi(match[1])
				numbers[number] = true
			}
		}
		return numbers, true
	}
	if _, gitRepositoryError := gitOutput(repositoryRoot, "rev-parse", "--is-inside-work-tree"); gitRepositoryError == nil {
		// An unborn repository has no landed authority. Do not reinterpret the
		// working tree as committed merely because HEAD does not exist yet.
		return numbers, true
	}
	// A failed repository probe is ambiguous: the directory may be outside Git, or
	// Git itself may be unavailable. Neither case supplies committed authority for
	// unattended deletion, so preserve every marker.
	return numbers, false
}
