package toolboxcommands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/atomicfile"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/gittransaction"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

var architectureNow = func() time.Time { return time.Now().UTC() }
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
	return usageResult(CommandArchitecture, "Usage: architecture-report-preflight --scan <reports-directory> | --publish <draft> <candidate> [--dry-run|--commit]")
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
	type prior struct{ key, path string }
	priors := []prior{}
	entries, _ := os.ReadDir(rootReports)
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
		priors = append(priors, prior{fmt.Sprintf("%s %012d", match[1], sequence), filepath.ToSlash(filepath.Join(reports, entry.Name(), "index.html"))})
	}
	sort.Slice(priors, func(i, j int) bool { return priors[i].key > priors[j].key })
	priorPath, priorHash, resolves := "", "", "n/a"
	if len(priors) > 0 {
		priorPath = priors[0].path
		data, _ := os.ReadFile(filepath.Join(ctx.RepositoryRoot, filepath.FromSlash(priorPath)))
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
	candidateAbs := candidate
	if !filepath.IsAbs(candidateAbs) {
		candidateAbs = filepath.Join(ctx.RepositoryRoot, candidate)
	}
	chosen := candidateAbs
	for sequence := 2; ; sequence++ {
		if _, err := os.Lstat(chosen); os.IsNotExist(err) {
			break
		}
		chosen = candidateAbs + fmt.Sprintf("-%d", sequence)
	}
	relative, err := filepath.Rel(ctx.RepositoryRoot, chosen)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return usageResult(CommandArchitecture, "candidate must stay inside repository root")
	}
	relative = filepath.ToSlash(relative)
	indexRel := filepath.ToSlash(filepath.Join(relative, "index.html"))
	parentRel := filepath.ToSlash(filepath.Dir(relative))
	createdDirectories := absentTransactionDirectories(ctx.RepositoryRoot, parentRel, relative)
	result := runTransaction(CommandArchitecture, ctx.RepositoryRoot, []string{indexRel}, createdDirectories, dryRun, commit, "[do-work] Publish architecture report", func(recorder *gittransaction.MutationRecorder) error {
		if err := os.MkdirAll(filepath.Dir(chosen), 0o755); err != nil {
			return err
		}
		for _, directory := range createdDirectories {
			if directory == relative {
				continue
			}
			if err := recorder.RecordCreatedDirectory(directory); err != nil {
				return err
			}
		}
		if err := os.Mkdir(chosen, 0o755); err != nil {
			return err
		}
		if err := recorder.RecordCreatedDirectory(relative); err != nil {
			return err
		}
		if err := atomicfile.CreateExclusive(filepath.Join(chosen, "index.html"), data, 0o644); err != nil {
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
