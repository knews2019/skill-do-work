// Package cleanup plans and applies deterministic do-work cleanup repairs.
package cleanup

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/schemanormalization"
)

type OperationKind string

const (
	OperationMove    OperationKind = "move"
	OperationReplace OperationKind = "replace"
	OperationDelete  OperationKind = "delete"
)

type CleanupOperation struct {
	Kind            OperationKind
	SourcePath      string
	DestinationPath string
	Contents        []byte
}

type OperationGroup struct {
	Code               string
	PassNumber         int
	AffectedID         string
	RequiredGroupCodes []string
	Operations         []CleanupOperation
}

type CleanupPlan struct {
	RepositoryRoot string
	Groups         []OperationGroup
	Findings       []resultmodel.CommandFinding
}

// BuildPlan derives Passes 0-4 from one repository snapshot. It never mutates.
func BuildPlan(snapshot *repositorymodel.RepositorySnapshot) CleanupPlan {
	plan := CleanupPlan{Findings: []resultmodel.CommandFinding{}}
	if snapshot == nil {
		plan.Findings = append(plan.Findings, manualFinding("CLEANUP-SNAPSHOT-REQUIRED", nil, nil, "repository discovery did not produce a snapshot"))
		return plan
	}
	plan.RepositoryRoot = snapshot.RepositoryRoot
	collidingIDs := map[string]bool{}
	for _, collision := range snapshot.CollisionEntries {
		collidingIDs[collision.RequestID] = true
		plan.Findings = append(plan.Findings, manualFinding("CLEANUP-ID-COLLISION", []string{collision.RequestID}, relativePaths(snapshot.RepositoryRoot, collision.ClaimPaths), "multiple paths claim the same request id"))
	}

	membersByUserRequest := map[string][]*repositorymodel.RequestFile{}
	for _, requestFile := range snapshot.RequestFiles {
		if strings.HasPrefix(requestFile.RelativePath, "archive/hold/") {
			continue
		}
		if requestFile.TypedRecord.UserRequestID != "" {
			membersByUserRequest[requestFile.TypedRecord.UserRequestID] = append(membersByUserRequest[requestFile.TypedRecord.UserRequestID], requestFile)
		}
	}
	resolvedUserRequests := map[string]bool{}
	userRequestSections := map[string]string{}
	memberGroupCodesByUserRequest := map[string]map[string]bool{}
	for _, userRequest := range snapshot.UserRequestFiles {
		if userRequest.TypedRecord.RequestID != "" {
			userRequestSections[userRequest.TypedRecord.RequestID] = userRequest.TreeSection
		}
		if userRequest.TreeSection != "user-requests" || userRequest.TypedRecord.RequestID == "" {
			continue
		}
		userRequestID := userRequest.TypedRecord.RequestID
		members := membersByUserRequest[userRequestID]
		resolved := true
		for _, member := range members {
			resolved = resolved && !collidingIDs[member.TypedRecord.RequestID] && schemanormalization.IsTerminalResolved(member.TypedRecord.RequestStatus)
		}
		resolvedUserRequests[userRequestID] = resolved
		if captureRequests, found := userRequest.ParsedDocument.FieldValue("requests"); found {
			for _, capturedID := range captureRequests.ListValues {
				if !requestIDExists(snapshot, capturedID) {
					plan.Findings = append(plan.Findings, manualFinding("CLEANUP-CAPTURE-ID-MISSING", []string{capturedID, userRequestID}, []string{repoRelative(snapshot.RepositoryRoot, userRequest.AbsolutePath)}, "request id listed in capture-time requests array was found nowhere"))
				}
			}
		}
	}

	for _, requestFile := range snapshot.RequestFiles {
		requestID := requestFile.TypedRecord.RequestID
		if requestID == "" {
			requestID = requestFile.FilenameID
		}
		if requestID == "" || collidingIDs[requestID] || requestFile.ParsedDocument == nil {
			continue
		}
		relativeSource := filepath.ToSlash(filepath.Join("do-work", requestFile.RelativePath))
		baseName := filepath.Base(requestFile.RelativePath)
		var relativeDestination string
		passNumber := 0
		switch {
		case (requestFile.TreeSection == "queue" || requestFile.TreeSection == "working") && schemanormalization.IsStopped(requestFile.TypedRecord.RequestStatus):
			if userRequestID := requestFile.TypedRecord.UserRequestID; userRequestID != "" && resolvedUserRequests[userRequestID] {
				relativeDestination = filepath.ToSlash(filepath.Join("do-work", "archive", userRequestID, baseName))
			} else {
				relativeDestination = filepath.ToSlash(filepath.Join("do-work", "archive", baseName))
			}
		case requestFile.TreeSection == "archive" && strings.Count(requestFile.RelativePath, "/") == 1:
			passNumber = 2
			if requestFile.TypedRecord.UserRequestID != "" {
				userRequestSection, exists := userRequestSections[requestFile.TypedRecord.UserRequestID]
				if !exists {
					plan.Findings = append(plan.Findings, manualFinding("CLEANUP-UR-NOT-FOUND", []string{requestFile.TypedRecord.UserRequestID}, []string{relativeSource}, "request references a user request folder that was not found"))
					continue
				}
				if userRequestSection == "user-requests" && !resolvedUserRequests[requestFile.TypedRecord.UserRequestID] {
					continue
				}
				relativeDestination = filepath.ToSlash(filepath.Join("do-work", "archive", requestFile.TypedRecord.UserRequestID, baseName))
			} else {
				relativeDestination = filepath.ToSlash(filepath.Join("do-work", "archive", "legacy", baseName))
			}
		default:
			continue
		}
		operation := CleanupOperation{Kind: OperationMove, SourcePath: relativeSource, DestinationPath: relativeDestination}
		if requestFile.TypedRecord.OriginalStatus != requestFile.TypedRecord.RequestStatus {
			document := requestFile.ParsedDocument
			if err := document.SetScalar("status", requestFile.TypedRecord.RequestStatus); err == nil {
				operation.Contents = document.DocumentBytes()
			}
		}
		operations := []CleanupOperation{operation}
		if requestFile.TreeSection == "working" {
			if checkpointOperation, found := ownedCheckpointRemoval(snapshot.RepositoryRoot, requestID); found {
				operations = append(operations, checkpointOperation)
			}
		}
		groupCode := "ARCHIVE-" + requestID
		plan.Groups = append(plan.Groups, OperationGroup{Code: groupCode, PassNumber: passNumber, AffectedID: requestID, Operations: operations})
		userRequestID := requestFile.TypedRecord.UserRequestID
		if userRequestID != "" && resolvedUserRequests[userRequestID] && userRequestSections[userRequestID] == "user-requests" {
			if memberGroupCodesByUserRequest[userRequestID] == nil {
				memberGroupCodesByUserRequest[userRequestID] = map[string]bool{}
			}
			memberGroupCodesByUserRequest[userRequestID][groupCode] = true
		}
	}

	for _, userRequest := range snapshot.UserRequestFiles {
		userRequestID := userRequest.TypedRecord.RequestID
		if userRequest.TreeSection != "user-requests" || !resolvedUserRequests[userRequestID] {
			continue
		}
		sourceDirectory := filepath.Dir(userRequest.AbsolutePath)
		operations, finding := planDirectoryMove(snapshot.RepositoryRoot, sourceDirectory, filepath.Join(snapshot.DoWorkRoot, "archive", userRequestID))
		if finding != nil {
			plan.Findings = append(plan.Findings, *finding)
			continue
		}
		if len(operations) > 0 {
			requiredGroupCodes := make([]string, 0, len(memberGroupCodesByUserRequest[userRequestID]))
			for groupCode := range memberGroupCodesByUserRequest[userRequestID] {
				requiredGroupCodes = append(requiredGroupCodes, groupCode)
			}
			sort.Strings(requiredGroupCodes)
			plan.Groups = append(plan.Groups, OperationGroup{Code: "CLOSE-" + userRequestID, PassNumber: 1, AffectedID: userRequestID, RequiredGroupCodes: requiredGroupCodes, Operations: operations})
		}
	}

	for _, manifest := range snapshot.RunManifestFiles {
		if manifest.Status != "consumed" {
			continue
		}
		operations, finding := planDirectoryDelete(snapshot.RepositoryRoot, filepath.Join(snapshot.DoWorkRoot, filepath.FromSlash(manifest.RunDirectory)))
		if finding != nil {
			plan.Findings = append(plan.Findings, *finding)
			continue
		}
		plan.Groups = append(plan.Groups, OperationGroup{Code: "SWEEP-RUN-" + filepath.Base(manifest.RunDirectory), PassNumber: 4, Operations: operations})
	}
	planLegacyContexts(snapshot, &plan)
	planMisplacedTrees(snapshot, &plan)
	sortPlan(&plan)
	return plan
}

func requestIDExists(snapshot *repositorymodel.RepositorySnapshot, requestID string) bool {
	for _, requestFile := range snapshot.RequestFiles {
		if requestFile.TypedRecord.RequestID == requestID || requestFile.FilenameID == requestID {
			return true
		}
	}
	return false
}

func planLegacyContexts(snapshot *repositorymodel.RepositorySnapshot, plan *CleanupPlan) {
	entries, err := os.ReadDir(filepath.Join(snapshot.DoWorkRoot, "archive"))
	if err != nil {
		return
	}
	for _, entry := range entries {
		if entry.IsDir() || entry.Type()&os.ModeSymlink != 0 || !strings.HasPrefix(strings.ToUpper(entry.Name()), "CONTEXT-") || !strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
			continue
		}
		plan.Groups = append(plan.Groups, OperationGroup{Code: "LEGACY-" + entry.Name(), PassNumber: 2, Operations: []CleanupOperation{{Kind: OperationMove,
			SourcePath: filepath.ToSlash(filepath.Join("do-work", "archive", entry.Name())), DestinationPath: filepath.ToSlash(filepath.Join("do-work", "archive", "legacy", entry.Name()))}}})
	}
}

func ownedCheckpointRemoval(repositoryRoot, requestID string) (CleanupOperation, bool) {
	checkpointPath := filepath.Join(repositoryRoot, "do-work", "CHECKPOINT.md")
	contents, err := os.ReadFile(checkpointPath)
	if err != nil {
		return CleanupOperation{}, false
	}
	hostname, err := os.Hostname()
	if err != nil {
		return CleanupOperation{}, false
	}
	if dotIndex := strings.IndexByte(hostname, '.'); dotIndex >= 0 {
		hostname = hostname[:dotIndex]
	}
	writerToken := "writer: " + hostname + ":" + repositoryRoot
	lines := strings.SplitAfter(string(contents), "\n")
	updated := make([]string, 0, len(lines))
	removed := false
	for _, line := range lines {
		if strings.Contains(line, "- "+requestID+":") && strings.Contains(line, writerToken) {
			removed = true
			continue
		}
		updated = append(updated, line)
	}
	if !removed {
		return CleanupOperation{}, false
	}
	return CleanupOperation{Kind: OperationReplace, SourcePath: "do-work/CHECKPOINT.md", Contents: []byte(strings.Join(updated, ""))}, true
}

func planDirectoryMove(repositoryRoot, sourceDirectory, destinationDirectory string) ([]CleanupOperation, *resultmodel.CommandFinding) {
	var operations []CleanupOperation
	err := filepath.WalkDir(sourceDirectory, func(path string, entry os.DirEntry, entryError error) error {
		if entryError != nil {
			return entryError
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink %s is not safe to move", path)
		}
		if entry.IsDir() {
			return nil
		}
		if !entry.Type().IsRegular() {
			return fmt.Errorf("special file %s is not safe to move", path)
		}
		relativeChild, _ := filepath.Rel(sourceDirectory, path)
		operations = append(operations, CleanupOperation{Kind: OperationMove, SourcePath: repoRelative(repositoryRoot, path), DestinationPath: repoRelative(repositoryRoot, filepath.Join(destinationDirectory, relativeChild))})
		return nil
	})
	if err != nil {
		finding := manualFinding("CLEANUP-UNSAFE-TREE", nil, []string{repoRelative(repositoryRoot, sourceDirectory)}, err.Error())
		return nil, &finding
	}
	return operations, nil
}

func planDirectoryDelete(repositoryRoot, directory string) ([]CleanupOperation, *resultmodel.CommandFinding) {
	var operations []CleanupOperation
	err := filepath.WalkDir(directory, func(path string, entry os.DirEntry, entryError error) error {
		if entryError != nil {
			return entryError
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink %s prevents safe run cleanup", path)
		}
		if entry.IsDir() {
			return nil
		}
		if !entry.Type().IsRegular() {
			return fmt.Errorf("special file %s prevents safe run cleanup", path)
		}
		operations = append(operations, CleanupOperation{Kind: OperationDelete, SourcePath: repoRelative(repositoryRoot, path)})
		return nil
	})
	if err != nil {
		finding := manualFinding("CLEANUP-RUN-REFUSED", nil, []string{repoRelative(repositoryRoot, directory)}, err.Error())
		return nil, &finding
	}
	return operations, nil
}

func planMisplacedTrees(snapshot *repositorymodel.RepositorySnapshot, plan *CleanupPlan) {
	misplacedArchiveRoot := filepath.Join(snapshot.DoWorkRoot, "archive", "user-requests")
	if entries, err := os.ReadDir(misplacedArchiveRoot); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() || entry.Type()&os.ModeSymlink != 0 || !strings.HasPrefix(strings.ToUpper(entry.Name()), "UR-") {
				continue
			}
			sourceDirectory := filepath.Join(misplacedArchiveRoot, entry.Name())
			operations, finding := planDirectoryMove(snapshot.RepositoryRoot, sourceDirectory, filepath.Join(snapshot.DoWorkRoot, "archive", entry.Name()))
			if finding != nil {
				plan.Findings = append(plan.Findings, *finding)
			} else if len(operations) > 0 {
				plan.Groups = append(plan.Groups, OperationGroup{Code: "FIX-ARCHIVE-" + entry.Name(), PassNumber: 3, AffectedID: entry.Name(), Operations: operations})
			}
		}
	}
	_ = filepath.WalkDir(snapshot.RepositoryRoot, func(path string, entry os.DirEntry, entryError error) error {
		if entryError != nil {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !entry.IsDir() {
			return nil
		}
		if path == filepath.Join(snapshot.RepositoryRoot, ".git") || path == snapshot.DoWorkRoot {
			return filepath.SkipDir
		}
		if path != snapshot.RepositoryRoot {
			if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
				return filepath.SkipDir
			}
		}
		if entry.Name() != "do-work" {
			return nil
		}
		if _, err := os.Stat(filepath.Join(path, "SKILL.md")); err == nil {
			return filepath.SkipDir
		}
		operations, finding := planDirectoryMove(snapshot.RepositoryRoot, path, snapshot.DoWorkRoot)
		if finding != nil {
			plan.Findings = append(plan.Findings, *finding)
		} else {
			groupedOperations := map[string][]CleanupOperation{}
			for _, operation := range operations {
				parts := strings.Split(strings.TrimPrefix(operation.SourcePath, repoRelative(snapshot.RepositoryRoot, path)+"/"), "/")
				itemKey := operation.SourcePath
				if len(parts) >= 2 && (parts[0] == "user-requests" || parts[0] == "archive") {
					itemKey = parts[0] + "/" + parts[1]
				}
				groupedOperations[itemKey] = append(groupedOperations[itemKey], operation)
			}
			for itemKey, itemOperations := range groupedOperations {
				plan.Groups = append(plan.Groups, OperationGroup{Code: "RELOCATE-" + strings.ReplaceAll(itemKey, "/", "-"), PassNumber: 3, Operations: itemOperations})
			}
		}
		return filepath.SkipDir
	})
}

func sortPlan(plan *CleanupPlan) {
	sort.Slice(plan.Groups, func(left, right int) bool {
		if plan.Groups[left].PassNumber != plan.Groups[right].PassNumber {
			return plan.Groups[left].PassNumber < plan.Groups[right].PassNumber
		}
		return plan.Groups[left].Code < plan.Groups[right].Code
	})
	for index := range plan.Groups {
		sort.Slice(plan.Groups[index].Operations, func(left, right int) bool {
			return operationPath(plan.Groups[index].Operations[left]) < operationPath(plan.Groups[index].Operations[right])
		})
	}
}

func operationPath(operation CleanupOperation) string {
	if operation.DestinationPath != "" {
		return operation.DestinationPath
	}
	return operation.SourcePath
}

func manualFinding(code string, ids, paths []string, evidence string) resultmodel.CommandFinding {
	return resultmodel.CommandFinding{Code: code, Severity: resultmodel.SeverityWarning, AffectedIDs: ids, AffectedPaths: paths,
		Evidence: []string{evidence}, Fixability: resultmodel.FixabilityManual, AutomationStopReason: "cleanup cannot prove this repair is safe",
		NextArgv: []string{"do-work-cli", "cleanup", "--dry-run"}, VerificationArgv: []string{"do-work-cli", "--format", "json", "cleanup", "--dry-run"}}
}

func repoRelative(repositoryRoot, path string) string {
	relativePath, err := filepath.Rel(repositoryRoot, path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(relativePath)
}

func relativePaths(repositoryRoot string, paths []string) []string {
	result := make([]string, 0, len(paths))
	for _, path := range paths {
		result = append(result, repoRelative(repositoryRoot, path))
	}
	return result
}
