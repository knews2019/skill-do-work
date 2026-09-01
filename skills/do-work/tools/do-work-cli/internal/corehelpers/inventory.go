package corehelpers

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requestmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

type inventoryRow struct{ Classification, Path, Origin string }

func handleInventory(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	if len(arguments) != 0 {
		return usageResult(CommandInventory, "uncommitted-inventory accepts no options")
	}
	rows, err := readInventory(executionContext.RepositoryRoot)
	if err != nil {
		return usageResult(CommandInventory, err.Error())
	}
	findings := make([]resultmodel.CommandFinding, 0, len(rows))
	for _, row := range rows {
		evidence := row.Classification
		if row.Origin != "" {
			evidence += " from " + row.Origin
		}
		findings = append(findings, helperFinding("INVENTORY-"+row.Classification, resultmodel.SeverityInfo, []string{row.Path}, evidence, resultmodel.FixabilityManual, "", nil, []string{"git", "status", "--short", "--", row.Path}))
	}
	result := successResult(nil, findings)
	if len(rows) == 0 {
		result.Outcome = resultmodel.OutcomeFindings
		result.Findings = []resultmodel.CommandFinding{helperFinding("INVENTORY-CLEAN", resultmodel.SeverityInfo, nil, "no uncommitted files", resultmodel.FixabilityAutomatic, "", nil, []string{"git", "status", "--short"})}
	}
	return result
}

func readInventory(repositoryRoot string) ([]inventoryRow, error) {
	output, err := gitOutput(repositoryRoot, "-c", "status.renames=copies", "status", "--porcelain=v1", "--untracked-files=all", "-z")
	if err != nil {
		return nil, err
	}
	records := bytes.Split(output, []byte{0})
	rows := []inventoryRow{}
	hasExcluded := false
	for index := 0; index < len(records); index++ {
		record := records[index]
		if len(record) == 0 {
			continue
		}
		if len(record) < 4 {
			return nil, fmt.Errorf("short porcelain record")
		}
		status, path := string(record[:2]), string(record[3:])
		origin := ""
		if strings.ContainsAny(status, "RC") {
			index++
			if index >= len(records) || len(records[index]) == 0 {
				return nil, fmt.Errorf("rename origin missing")
			}
			origin = string(records[index])
		}
		if strings.Contains(status, "R") && secretPath(origin) {
			rows = append(rows, inventoryRow{"XD", origin, ""})
			hasExcluded = true
		}
		classification := classifyInventory(status, path, origin)
		if classification == "X" || classification == "XD" {
			hasExcluded = true
		}
		rows = append(rows, inventoryRow{classification, path, origin})
	}
	if hasExcluded {
		for index := range rows {
			if rows[index].Classification == "A" {
				rows[index].Classification = "X"
			}
		}
	}
	return rows, nil
}

func classifyInventory(status, path, origin string) string {
	if secretPath(path) || secretPath(origin) {
		if strings.Contains(status, "D") {
			return "XD"
		}
		return "X"
	}
	if status == "??" || strings.Contains(status, "A") {
		return "A"
	}
	if strings.Contains(status, "D") {
		return "D"
	}
	return "M"
}
func secretPath(path string) bool {
	if path == "" {
		return false
	}
	base := strings.ToLower(filepath.Base(path))
	if strings.HasPrefix(base, ".env") || strings.HasSuffix(base, ".env") || strings.HasSuffix(base, ".pem") || strings.HasSuffix(base, ".key") || strings.HasSuffix(base, ".p12") || strings.HasSuffix(base, ".pfx") || strings.Contains(base, "credential") || strings.Contains(base, "secret") {
		return true
	}
	return false
}

func handleAssociate(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	pathsFile, parseResult := requiredPathOption(arguments, "--paths-file", CommandAssociate)
	if parseResult != nil {
		return *parseResult
	}
	file, err := os.Open(absoluteFromRoot(executionContext.RepositoryRoot, pathsFile))
	if err != nil {
		return usageResult(CommandAssociate, err.Error())
	}
	defer file.Close()
	candidates := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if value := strings.TrimSpace(scanner.Text()); value != "" {
			candidates = append(candidates, value)
		}
	}
	if err := scanner.Err(); err != nil {
		return usageResult(CommandAssociate, err.Error())
	}
	associations, err := associatePaths(executionContext.RepositoryRoot, candidates)
	if err != nil {
		return usageResult(CommandAssociate, err.Error())
	}
	findings := []resultmodel.CommandFinding{}
	for _, candidate := range candidates {
		id := associations[candidate]
		code := "ASSOCIATION-UNOWNED"
		severity := resultmodel.SeverityWarning
		evidence := "no active or terminal request claims this path"
		if id != "" {
			code = "ASSOCIATION-FOUND"
			severity = resultmodel.SeverityInfo
			evidence = "owned by " + id
		}
		findings = append(findings, helperFinding(code, severity, []string{candidate}, evidence, resultmodel.FixabilityManual, map[bool]string{true: "assign or quarantine the path", false: ""}[id == ""], nil, nil))
	}
	return successResult(nil, findings)
}

func associatePaths(repositoryRoot string, candidates []string) (map[string]string, error) {
	claims := map[string]struct {
		id        string
		completed time.Time
		active    bool
	}{}
	roots := []string{filepath.Join(repositoryRoot, "do-work", "working"), filepath.Join(repositoryRoot, "do-work", "archive")}
	for _, root := range roots {
		walkError := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				if os.IsNotExist(walkErr) {
					return nil
				}
				return walkErr
			}
			if entry.Type()&os.ModeSymlink != 0 {
				if entry.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if entry.IsDir() || !strings.HasPrefix(entry.Name(), "REQ-") || !strings.HasSuffix(entry.Name(), ".md") {
				return nil
			}
			contents, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			document, err := requestmodel.ParseDocument(contents)
			if err != nil {
				return nil
			}
			record := document.TypedRecord()
			active := strings.Contains(filepath.ToSlash(path), "/working/")
			status := record.RequestStatus
			if !active && !terminalSuccessStatus(status) {
				return nil
			}
			paths, found, parseErr := allBacktickedPaths(string(contents), "Implementation Summary")
			if parseErr != nil {
				return parseErr
			}
			if !found {
				return nil
			}
			completed, _ := requestmodel.ParseTimestamp(record.CompletedAt)
			for _, claimed := range paths {
				current, exists := claims[claimed]
				if !exists || completed.After(current.completed) {
					claims[claimed] = struct {
						id        string
						completed time.Time
						active    bool
					}{record.RequestID, completed, active}
				}
			}
			return nil
		})
		if walkError != nil {
			return nil, walkError
		}
	}
	output := map[string]string{}
	for _, candidate := range candidates {
		output[candidate] = claims[candidate].id
	}
	return output, nil
}

func terminalSuccessStatus(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "completed", "completed-with-issues", "complete", "done", "finished", "closed":
		return true
	default:
		return false
	}
}

func handleProtectedInventory(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	if len(arguments) == 0 {
		return usageResult(CommandProtectedInventory, "start or associate is required")
	}
	mode := arguments[0]
	quarantineName := "do-work-protected-inventory"
	for index := 1; index < len(arguments); index++ {
		if arguments[index] == "--quarantine-name" || strings.HasPrefix(arguments[index], "--quarantine-name=") {
			value, err := optionValue(arguments, &index, "--quarantine-name")
			if err != nil {
				return usageResult(CommandProtectedInventory, err.Error())
			}
			if strings.ContainsAny(value, "/\\\n\r") {
				return usageResult(CommandProtectedInventory, "invalid quarantine name")
			}
			quarantineName = value
		} else {
			return usageResult(CommandProtectedInventory, "unknown option "+arguments[index])
		}
	}
	if mode != "start" && mode != "associate" {
		return usageResult(CommandProtectedInventory, "mode must be start or associate")
	}
	gitPath, err := gitOutput(executionContext.RepositoryRoot, "rev-parse", "--git-path", quarantineName)
	if err != nil {
		return usageResult(CommandProtectedInventory, err.Error())
	}
	quarantinePath := strings.TrimSpace(string(gitPath))
	if !filepath.IsAbs(quarantinePath) {
		quarantinePath = filepath.Join(executionContext.RepositoryRoot, quarantinePath)
	}
	rows, err := readInventory(executionContext.RepositoryRoot)
	if err != nil {
		return usageResult(CommandProtectedInventory, err.Error())
	}
	protected := []string{}
	ordinary := []string{}
	for _, row := range rows {
		if row.Classification == "X" || row.Classification == "XD" {
			protected = append(protected, row.Path)
		} else {
			ordinary = append(ordinary, row.Path)
		}
	}
	if mode == "start" {
		if len(protected) == 0 {
			_ = os.Remove(quarantinePath)
		} else {
			if err := writePrivateAtomic(quarantinePath, []byte(strings.Join(uniqueSorted(protected), "\n")+"\n"), 0o600); err != nil {
				return usageResult(CommandProtectedInventory, err.Error())
			}
		}
		return successResult([]resultmodel.RecordedChange{{Path: quarantinePath, Kind: "git-private", Detail: fmt.Sprintf("recorded %d protected paths", len(protected))}}, nil)
	}
	existing, _ := os.ReadFile(quarantinePath)
	quarantined := stringSet(append(nonblankLines(existing), protected...))
	candidateFile, err := os.CreateTemp("", "do-work-associate-*.txt")
	if err != nil {
		return usageResult(CommandProtectedInventory, err.Error())
	}
	candidatePath := candidateFile.Name()
	defer os.Remove(candidatePath)
	for _, path := range ordinary {
		if !quarantined[path] {
			fmt.Fprintln(candidateFile, path)
		}
	}
	_ = candidateFile.Close()
	return handleAssociate(executionContext, []string{"--paths-file", candidatePath})
}
