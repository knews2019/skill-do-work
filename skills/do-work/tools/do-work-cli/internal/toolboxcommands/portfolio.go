package toolboxcommands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/gittransaction"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func handlePortfolio(ctx commandruntime.ExecutionContext, args []string) resultmodel.CommandResult {
	rest, dryRun, commit, err := parseMutationFlags(args)
	if err != nil {
		return usageResult(CommandPortfolio, err.Error())
	}
	if len(rest) != 3 && len(rest) != 4 {
		return usageResult(CommandPortfolio, "Usage: publish-portfolio-summary [--dry-run|--commit] (--canonical-only <source> <canonical> | --with-snapshot <source> <canonical> <candidate>)")
	}
	mode := rest[0]
	if mode != "--canonical-only" && mode != "--with-snapshot" {
		return usageResult(CommandPortfolio, "unknown portfolio mode")
	}
	if mode == "--canonical-only" && len(rest) != 3 || mode == "--with-snapshot" && len(rest) != 4 {
		return usageResult(CommandPortfolio, "portfolio mode has wrong argument count")
	}
	source := absoluteFromRoot(ctx.RepositoryRoot, rest[1])
	data, readErr := os.ReadFile(source)
	if readErr != nil {
		return usageResult(CommandPortfolio, "Portfolio source is not a regular file: "+rest[1])
	}
	if info, statErr := os.Stat(source); statErr != nil || !info.Mode().IsRegular() {
		return usageResult(CommandPortfolio, "Portfolio source is not a regular file: "+rest[1])
	}
	canonicalRel, canonicalAbs, pathErr := repositoryPath(ctx.RepositoryRoot, rest[2])
	if pathErr != nil {
		return usageResult(CommandPortfolio, pathErr.Error())
	}
	_, canonicalStatError := os.Lstat(canonicalAbs)
	canonicalExisted := canonicalStatError == nil
	if canonicalStatError != nil && !os.IsNotExist(canonicalStatError) {
		return usageResult(CommandPortfolio, canonicalStatError.Error())
	}
	targets := []string{canonicalRel}
	directories := []string{filepath.ToSlash(filepath.Dir(canonicalRel))}
	snapshotRel := ""
	if mode == "--with-snapshot" {
		snapshotRel, _, pathErr = firstFreePortfolioPath(ctx.RepositoryRoot, rest[3])
		if pathErr != nil {
			return usageResult(CommandPortfolio, pathErr.Error())
		}
		targets = append(targets, snapshotRel)
		directories = append(directories, filepath.ToSlash(filepath.Dir(snapshotRel)))
	}
	createdDirectories := absentTransactionDirectories(ctx.RepositoryRoot, directories...)
	if !commit && !dryRun {
		return publishPortfolioDirect(ctx.RepositoryRoot, data, canonicalRel, canonicalAbs, canonicalExisted, mode, rest)
	}
	result := runTransaction(CommandPortfolio, ctx.RepositoryRoot, targets, createdDirectories, dryRun, commit, "[do-work] Publish portfolio summary", func(recorder *gittransaction.MutationRecorder) error {
		for _, dir := range createdDirectories {
			if err := rootedMkdirAll(ctx.RepositoryRoot, dir, 0o755); err != nil {
				return err
			}
			if err := recorder.RecordCreatedDirectory(dir); err != nil {
				return err
			}
		}
		if snapshotRel != "" {
			if err := rootedPublishFile(ctx.RepositoryRoot, snapshotRel, data, 0o644, false); err != nil {
				return err
			}
			if err := recorder.RecordCreated(snapshotRel); err != nil {
				return err
			}
		}
		if canonicalExisted {
			info, err := os.Lstat(canonicalAbs)
			if err != nil {
				return err
			}
			if info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
				return fmt.Errorf("Portfolio canonical path is unsafe: %s", canonicalRel)
			}
			if err := rootedPublishFile(ctx.RepositoryRoot, canonicalRel, data, 0o644, true); err != nil {
				return err
			}
			return recorder.RecordTouched(canonicalRel)
		}
		if err := rootedPublishFile(ctx.RepositoryRoot, canonicalRel, data, 0o644, false); err != nil {
			return err
		}
		return recorder.RecordCreated(canonicalRel)
	})
	if result.Outcome == resultmodel.OutcomeSuccess {
		paths := []string{canonicalRel}
		if snapshotRel != "" {
			paths = append(paths, snapshotRel)
		}
		output := strings.Join(paths, "\n") + "\n"
		if dryRun {
			output = "would publish\n" + output
		}
		result.ExactTextOutput = &output
	}
	return result
}

func publishPortfolioDirect(repositoryRoot string, data []byte, canonicalRel, canonicalAbs string, canonicalExisted bool, mode string, arguments []string) resultmodel.CommandResult {
	result := resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, RepositoryRoot: repositoryRoot}
	publishedSnapshotRel := ""
	publishedSnapshotDisplay := ""
	if mode == "--with-snapshot" {
		for {
			snapshotRel, snapshotAbs, err := firstFreePortfolioPath(repositoryRoot, arguments[3])
			if err != nil {
				return portfolioDirectFailure(repositoryRoot, "PORTFOLIO-SNAPSHOT-PATH-FAILED", arguments[3], err, "canonical was not changed", "")
			}
			if err := rootedMkdirAll(repositoryRoot, filepath.ToSlash(filepath.Dir(snapshotRel)), 0o755); err != nil {
				return portfolioDirectFailure(repositoryRoot, "PORTFOLIO-SNAPSHOT-DIRECTORY-FAILED", snapshotRel, err, "canonical was not changed", "")
			}
			if err := rootedPublishFile(repositoryRoot, snapshotRel, data, 0o644, false); err != nil {
				if _, statErr := os.Lstat(snapshotAbs); statErr == nil {
					continue
				}
				return portfolioDirectFailure(repositoryRoot, "PORTFOLIO-SNAPSHOT-PUBLISH-FAILED", snapshotRel, err, "canonical was not changed", "")
			}
			publishedSnapshotRel = snapshotRel
			publishedSnapshotDisplay = portfolioDisplayPath(arguments[3], snapshotAbs)
			result.Changes = append(result.Changes, resultmodel.RecordedChange{Path: snapshotRel, Kind: "created", Detail: "published immutable portfolio snapshot"})
			break
		}
	}
	if err := rootedMkdirAll(repositoryRoot, filepath.ToSlash(filepath.Dir(canonicalRel)), 0o755); err != nil {
		return portfolioDirectFailure(repositoryRoot, "PORTFOLIO-CANONICAL-DIRECTORY-FAILED", canonicalRel, err, "snapshot retained when present", publishedSnapshotDisplay)
	}
	if canonicalExisted {
		info, err := os.Lstat(canonicalAbs)
		if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			if err == nil {
				err = fmt.Errorf("Portfolio canonical path is unsafe: %s", canonicalRel)
			}
			return portfolioDirectFailure(repositoryRoot, "PORTFOLIO-CANONICAL-UNSAFE", canonicalRel, err, "snapshot retained when present", publishedSnapshotDisplay)
		}
	}
	if err := rootedPublishFile(repositoryRoot, canonicalRel, data, 0o644, canonicalExisted); err != nil {
		return portfolioDirectFailure(repositoryRoot, "PORTFOLIO-CANONICAL-PUBLISH-FAILED", canonicalRel, err, "snapshot retained when present", publishedSnapshotDisplay)
	}
	kind := "created"
	if canonicalExisted {
		kind = "modified"
	}
	result.Changes = append(result.Changes, resultmodel.RecordedChange{Path: canonicalRel, Kind: kind, Detail: "published canonical portfolio summary"})
	output := arguments[2] + "\n"
	if publishedSnapshotRel != "" {
		output += publishedSnapshotDisplay + "\n"
	}
	result.ExactTextOutput = &output
	return result
}

func portfolioDisplayPath(original, chosenAbsolute string) string {
	chosenBase := filepath.Base(chosenAbsolute)
	if separator := strings.LastIndexAny(original, `/\\`); separator >= 0 {
		return original[:separator+1] + chosenBase
	}
	return chosenBase
}

func portfolioDirectFailure(repositoryRoot, code, path string, err error, stopReason, retainedSnapshot string) resultmodel.CommandResult {
	result := resultmodel.CommandResult{Outcome: resultmodel.OutcomeFindings, RepositoryRoot: repositoryRoot,
		Findings: []resultmodel.CommandFinding{toolboxFinding(CommandPortfolio, code, resultmodel.SeverityWarning, []string{path}, err.Error(), resultmodel.FixabilityManual, stopReason)}}
	if retainedSnapshot != "" {
		result.Findings = append(result.Findings, toolboxFinding(CommandPortfolio, "PORTFOLIO-SNAPSHOT-RETAINED", resultmodel.SeverityWarning, []string{retainedSnapshot}, "immutable snapshot published before canonical failure and intentionally retained", resultmodel.FixabilityManual, "repair the canonical destination without deleting the snapshot"))
	}
	return result
}

func absoluteFromRoot(root, path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(root, path)
}
func repositoryPath(root, path string) (string, string, error) {
	absolute := absoluteFromRoot(root, path)
	relative, err := filepath.Rel(root, absolute)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", "", fmt.Errorf("path must stay inside repository root: %s", path)
	}
	relative = filepath.ToSlash(relative)
	if err := validateNoLinkedAncestors(root, relative, false); err != nil {
		return "", "", err
	}
	return relative, absolute, nil
}
func firstFreePortfolioPath(root, candidate string) (string, string, error) {
	_, absolute, err := repositoryPath(root, candidate)
	if err != nil {
		return "", "", err
	}
	extension := filepath.Ext(absolute)
	stem := strings.TrimSuffix(absolute, extension)
	for sequence := 1; ; sequence++ {
		chosen := absolute
		if sequence > 1 {
			chosen = fmt.Sprintf("%s-%d%s", stem, sequence, extension)
		}
		if _, err := os.Lstat(chosen); os.IsNotExist(err) {
			chosenRel, _ := filepath.Rel(root, chosen)
			return filepath.ToSlash(chosenRel), chosen, nil
		}
	}
}
