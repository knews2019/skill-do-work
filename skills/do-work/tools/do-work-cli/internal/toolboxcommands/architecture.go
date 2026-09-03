package toolboxcommands

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/gittransaction"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

var architectureNow = func() time.Time { return time.Now().UTC() }
var architectureAfterClaim = func(string) {}
var architectureClaimBundle = rootedMkdirExclusive
var architectureName = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}_\d{4})_architecture-report(?:-(\d+))?$`)
var architectureWatermark = regexp.MustCompile(`(?m)^\s*<meta name="architecture-report-verified-at" content="([0-9a-f]{7,40})">\s*$`)

func handleArchitecture(ctx commandruntime.ExecutionContext, args []string) resultmodel.CommandResult {
	rest, dryRun, commit, err := parseMutationFlags(args)
	if err != nil {
		return usageResult(CommandArchitecture, err.Error())
	}
	if len(rest) == 2 && rest[0] == "--scan" {
		if dryRun || commit {
			return usageResult(CommandArchitecture, "--scan is read-only")
		}
		return architectureScan(ctx, rest[1])
	}
	if len(rest) == 3 && rest[0] == "--publish" {
		return architecturePublish(ctx, rest[1], rest[2], dryRun, commit)
	}
	return usageResult(CommandArchitecture, "Usage: architecture-report-preflight [--dry-run|--commit] (--scan <reports-directory> | --publish <draft> <candidate>)")
}

func architectureScan(ctx commandruntime.ExecutionContext, reports string) resultmodel.CommandResult {
	headBytes, err := exec.Command("git", "-C", ctx.RepositoryRoot, "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		return architectureFailure("no resolvable HEAD commit to watermark against")
	}
	head := strings.TrimSpace(string(headBytes))
	if head == "" {
		return architectureFailure("no resolvable HEAD commit to watermark against")
	}
	slug := architectureNow().Format("2006-01-02_1504")
	candidate := filepath.ToSlash(filepath.Join(reports, slug+"_architecture-report"))
	rootReports := reports
	if !filepath.IsAbs(rootReports) {
		rootReports = filepath.Join(ctx.RepositoryRoot, rootReports)
	}
	type prior struct{ key, filesystemPath, displayPath string }
	priors := []prior{}
	entries, readDirectoryError := os.ReadDir(rootReports)
	if os.IsNotExist(readDirectoryError) {
		entries = nil
	} else if readDirectoryError != nil {
		return architectureFailure("reports directory is unreadable: " + readDirectoryError.Error())
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		match := architectureName.FindStringSubmatch(entry.Name())
		if match == nil {
			continue
		}
		sequence := 1
		if match[2] != "" {
			sequence, _ = strconv.Atoi(match[2])
		}
		index := filepath.Join(rootReports, entry.Name(), "index.html")
		info, statErr := os.Stat(index)
		if statErr != nil || !info.Mode().IsRegular() || info.Size() == 0 {
			continue
		}
		priors = append(priors, prior{fmt.Sprintf("%s %012d", match[1], sequence), index, filepath.ToSlash(filepath.Join(reports, entry.Name(), "index.html"))})
	}
	sort.Slice(priors, func(i, j int) bool { return priors[i].key > priors[j].key })
	priorPath, priorHash, resolves := "", "", "n/a"
	if len(priors) > 0 {
		priorFilesystemPath := priors[0].filesystemPath
		priorPath = priors[0].displayPath
		data, _ := os.ReadFile(priorFilesystemPath)
		match := architectureWatermark.FindSubmatch(data)
		priorHash = "unreadable"
		if match != nil {
			priorHash = string(match[1])
		}
		resolves = "no"
		if priorHash != "unreadable" && exec.Command("git", "-C", ctx.RepositoryRoot, "rev-parse", "--verify", "-q", priorHash+"^{commit}").Run() == nil {
			resolves = "yes"
		}
	}
	output := fmt.Sprintf("head_hash=%s\nreport_slug=%s\nreport_candidate=%s\nprior_report=%s\nprior_hash=%s\nprior_hash_resolves=%s\n", head, slug, candidate, priorPath, priorHash, resolves)
	f := toolboxFinding(CommandArchitecture, "ARCHITECTURE-SCAN", resultmodel.SeverityInfo, []string{candidate}, strings.TrimSpace(output), resultmodel.FixabilityAutomatic, "")
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Findings: []resultmodel.CommandFinding{f}, ExactTextOutput: &output}
}

func architectureFailure(evidence string) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFindings, Findings: []resultmodel.CommandFinding{
		toolboxFinding(CommandArchitecture, "ARCHITECTURE-PREFLIGHT-FAILED", resultmodel.SeverityError, nil, evidence, resultmodel.FixabilityManual, "architecture report preflight could not establish complete evidence"),
	}}
}

func architecturePublish(ctx commandruntime.ExecutionContext, draft, candidate string, dryRun, commit bool) resultmodel.CommandResult {
	draftPath := draft
	if !filepath.IsAbs(draftPath) {
		draftPath = filepath.Join(ctx.RepositoryRoot, draftPath)
	}
	data, err := os.ReadFile(draftPath)
	if err != nil {
		return usageResult(CommandArchitecture, "draft is not a regular readable file: "+draft)
	}
	if info, statErr := os.Stat(draftPath); statErr != nil || !info.Mode().IsRegular() {
		return usageResult(CommandArchitecture, "draft is not a regular file: "+draft)
	}
	if os.Getenv("DO_WORK_COMPATIBILITY_SHIM") == "1" {
		staged, stageError := os.CreateTemp("", "architecture-report-copy.*")
		if stageError != nil {
			return architectureFailure(stageError.Error())
		}
		stagedPath := staged.Name()
		_ = staged.Close()
		defer os.Remove(stagedPath)
		if output, copyError := exec.Command("cp", draftPath, stagedPath).CombinedOutput(); copyError != nil {
			return architectureFailure("draft copy failed: " + strings.TrimSpace(string(output)))
		}
		data, err = os.ReadFile(stagedPath)
		if err != nil {
			return architectureFailure(err.Error())
		}
	}
	candidateAbs := candidate
	if !filepath.IsAbs(candidateAbs) {
		candidateAbs = filepath.Join(ctx.RepositoryRoot, candidate)
	}
	baseRelative, _, pathErr := repositoryPath(ctx.RepositoryRoot, candidateAbs)
	if pathErr != nil {
		return usageResult(CommandArchitecture, pathErr.Error())
	}
	relative := baseRelative
	sequence := 1
	for {
		root, openErr := os.OpenRoot(ctx.RepositoryRoot)
		if openErr != nil {
			return architectureFailure(openErr.Error())
		}
		_, statErr := root.Lstat(filepath.FromSlash(relative))
		root.Close()
		if os.IsNotExist(statErr) {
			break
		}
		sequence++
		relative = fmt.Sprintf("%s-%d", baseRelative, sequence)
	}
	indexRel := filepath.ToSlash(filepath.Join(relative, "index.html"))
	preflight := runTransaction(CommandArchitecture, ctx.RepositoryRoot, []string{indexRel}, nil, true, false, "[do-work] Publish architecture report", nil)
	if preflight.Outcome != resultmodel.OutcomeSuccess || dryRun {
		if dryRun && preflight.Outcome == resultmodel.OutcomeSuccess {
			output := "would publish " + indexRel + "\n"
			preflight.ExactTextOutput = &output
		}
		return preflight
	}
	if parentErr := rootedMkdirAll(ctx.RepositoryRoot, filepath.ToSlash(filepath.Dir(relative)), 0o755); parentErr != nil {
		return architectureFailure(parentErr.Error())
	}
	var ownedBundle os.FileInfo
	for {
		var claimErr error
		if ownedBundle, claimErr = architectureClaimBundle(ctx.RepositoryRoot, relative, 0o755); claimErr == nil {
			break
		}
		if !errors.Is(claimErr, fs.ErrExist) {
			return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFindings, Findings: []resultmodel.CommandFinding{
				toolboxFinding(CommandArchitecture, "ARCHITECTURE-BUNDLE-CLAIM-FAILED", resultmodel.SeverityError, []string{relative}, "architecture bundle claim failed: "+claimErr.Error(), resultmodel.FixabilityManual, "architecture report bundle could not be claimed"),
			}}
		}
		sequence++
		relative = fmt.Sprintf("%s-%d", baseRelative, sequence)
		indexRel = filepath.ToSlash(filepath.Join(relative, "index.html"))
	}
	architectureAfterClaim(filepath.Join(ctx.RepositoryRoot, filepath.FromSlash(relative)))
	result := runTransaction(CommandArchitecture, ctx.RepositoryRoot, []string{indexRel}, nil, false, commit, "[do-work] Publish architecture report", func(recorder *gittransaction.MutationRecorder) error {
		if err := rootedPublishInOwnedDirectory(ctx.RepositoryRoot, relative, ownedBundle, "index.html", data, 0o644); err != nil {
			return err
		}
		return recorder.RecordCreated(indexRel)
	})
	if result.Outcome == resultmodel.OutcomeSuccess {
		output := indexRel + "\n"
		if dryRun {
			output = "would publish " + indexRel + "\n"
		}
		result.ExactTextOutput = &output
	}
	return result
}
