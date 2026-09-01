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

func handleCleanupReservations(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	filtered, dryRun, dryRunError := extractDryRun(arguments)
	if dryRunError != nil || len(filtered) != 0 {
		return usageResult(CommandCleanupReservations, "cleanup-req-reservations accepts no options")
	}
	doWorkRoot := filepath.Join(executionContext.RepositoryRoot, "do-work")
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
	claimed := claimedRequestNumbers(executionContext.RepositoryRoot)
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

func claimedRequestNumbers(repositoryRoot string) map[int]bool {
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
		return numbers
	}
	for _, rootName := range []string{"queue", "working", "archive"} {
		root := filepath.Join(repositoryRoot, "do-work", rootName)
		_ = filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return nil
			}
			if entry.Type()&os.ModeSymlink != 0 {
				if entry.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
				return nil
			}
			match := requestFilePattern.FindStringSubmatch(entry.Name())
			if match != nil {
				number, _ := strconv.Atoi(match[1])
				numbers[number] = true
			}
			return nil
		})
	}
	return numbers
}
